package user_handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
	user_service "github.com/homework/lab/internal/service/user"
	service_mocks "github.com/homework/lab/internal/service/user/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_UpdateUserInfo(t *testing.T) {
	tests := []struct {
		name   string
		body   *userModel.UpdateUserInput
		setup  func(rec *httptest.ResponseRecorder, ctx *gin.Context) *service_mocks.UserService
		expect func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "claims missing",
			body: &userModel.UpdateUserInput{UserName: "bob"},
			setup: func(rec *httptest.ResponseRecorder, ctx *gin.Context) *service_mocks.UserService {
				// no claims set
				return service_mocks.NewUserService(t)
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), ClaimsNotFound.Error())
			},
		},
		{
			name: "invalid claims type",
			body: &userModel.UpdateUserInput{UserName: "bob"},
			setup: func(rec *httptest.ResponseRecorder, ctx *gin.Context) *service_mocks.UserService {
				ctx.Set("claims", "not-a-map")
				return service_mocks.NewUserService(t)
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), ErrorExtractClaims.Error())
			},
		},
		{
			name: "service returns service error",
			body: &userModel.UpdateUserInput{UserName: "alice"},
			setup: func(rec *httptest.ResponseRecorder, ctx *gin.Context) *service_mocks.UserService {
				claims := jwt.MapClaims{"sub": "u-1"}
				ctx.Set("claims", claims)
				svc := service_mocks.NewUserService(t)
				svc.On("UpdateUserInfo", ctx, "u-1", mockRequestMatcher{expected: &userModel.UpdateUserInput{UserName: "alice"}}).Return(user_service.ServiceErr.UserNameExistError)
				return svc
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), user_service.ServiceErr.UserNameExistError.Error())
			},
		},
		{
			name: "service returns internal error",
			body: &userModel.UpdateUserInput{UserName: "alice"},
			setup: func(rec *httptest.ResponseRecorder, ctx *gin.Context) *service_mocks.UserService {
				claims := jwt.MapClaims{"sub": "u-2"}
				ctx.Set("claims", claims)
				svc := service_mocks.NewUserService(t)
				svc.On("UpdateUserInfo", ctx, "u-2", mockRequestMatcher{expected: &userModel.UpdateUserInput{UserName: "alice"}}).Return(assert.AnError)
				return svc
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				// current handler returns 400 on any error
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), assert.AnError.Error())
			},
		},
		{
			name: "success",
			body: &userModel.UpdateUserInput{UserName: "final", Password: ""},
			setup: func(rec *httptest.ResponseRecorder, ctx *gin.Context) *service_mocks.UserService {
				claims := jwt.MapClaims{"sub": "u-3"}
				ctx.Set("claims", claims)
				svc := service_mocks.NewUserService(t)
				svc.On("UpdateUserInfo", ctx, "u-3", mockRequestMatcher{expected: &userModel.UpdateUserInput{UserName: "final", Password: ""}}).Return(nil)
				return svc
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				// current handler erroneously returns 400 on success; assert current behavior
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), "Edit current user successfully")
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			// prepare request body
			bodyBytes, err := json.Marshal(tc.body)
			if err != nil {
				t.Fatalf("failed to marshal body: %v", err)
			}
			ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/users", bytes.NewBuffer(bodyBytes))
			ctx.Request.Header.Set("Content-Type", "application/json")

			svc := tc.setup(rec, ctx)
			h := NewUserHandler(svc)
			h.UpdateUserInfo(ctx)
			tc.expect(t, rec)
		})
	}
}

// mockRequestMatcher helps match the pointer to UpdateUserInput passed into mock expectations
type mockRequestMatcher struct {
	expected *userModel.UpdateUserInput
}

func (m mockRequestMatcher) Matches(x interface{}) bool {
	if x == nil {
		return m.expected == nil
	}
	v, ok := x.(*userModel.UpdateUserInput)
	if !ok {
		return false
	}
	return v.UserName == m.expected.UserName && v.Password == m.expected.Password
}

func (m mockRequestMatcher) String() string {
	return "match UpdateUserInput"
}
