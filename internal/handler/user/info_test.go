package user_handler

import (
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

func TestHandler_GetUserInfo(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(rec *httptest.ResponseRecorder, ctx *gin.Context) *service_mocks.UserService
		expect    func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "claims missing",
			setup: func(rec *httptest.ResponseRecorder, ctx *gin.Context) *service_mocks.UserService {
				// do not set claims
				svc := service_mocks.NewUserService(t)
				return svc
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), ClaimsNotFound.Error())
			},
		},
		{
			name: "invalid claims type",
			setup: func(rec *httptest.ResponseRecorder, ctx *gin.Context) *service_mocks.UserService {
				ctx.Set("claims", "not-a-map")
				svc := service_mocks.NewUserService(t)
				return svc
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), ErrorExtractClaims.Error())
			},
		},
		{
			name: "service not found",
			setup: func(rec *httptest.ResponseRecorder, ctx *gin.Context) *service_mocks.UserService {
				claims := jwt.MapClaims{"sub": "missing-id"}
				ctx.Set("claims", claims)
				svc := service_mocks.NewUserService(t)
				svc.On("GetUserInfo", ctx, "missing-id").Return((*userModel.UserInfo)(nil), user_service.ServiceErr.NotFoundUserInfo)
				return svc
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), user_service.ServiceErr.NotFoundUserInfo.Error())
			},
		},
		{
			name: "service internal error",
			setup: func(rec *httptest.ResponseRecorder, ctx *gin.Context) *service_mocks.UserService {
				claims := jwt.MapClaims{"sub": "some-id"}
				ctx.Set("claims", claims)
				svc := service_mocks.NewUserService(t)
				svc.On("GetUserInfo", ctx, "some-id").Return((*userModel.UserInfo)(nil), assert.AnError)
				return svc
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
				assert.Contains(t, rec.Body.String(), "Internal server error")
			},
		},
		{
			name: "success",
			setup: func(rec *httptest.ResponseRecorder, ctx *gin.Context) *service_mocks.UserService {
				claims := jwt.MapClaims{"sub": "u-1"}
				ctx.Set("claims", claims)
				svc := service_mocks.NewUserService(t)
				info := &userModel.UserInfo{Id: "u-1", DisplayName: "Bob", Email: "bob@example.com", UserName: "bob"}
				svc.On("GetUserInfo", ctx, "u-1").Return(info, nil)
				return svc
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rec.Code)
				body := rec.Body.String()
				assert.Contains(t, body, "u-1")
				assert.Contains(t, body, "Bob")
				assert.Contains(t, body, "bob@example.com")
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			svc := tc.setup(rec, ctx)
			h := NewUserHandler(svc)
			h.GetUserInfo(ctx)
			tc.expect(t, rec)
		})
	}
}
