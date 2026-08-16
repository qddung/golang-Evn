package user_handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/models/dto/api/user"
	user_service "github.com/homework/lab/internal/service/user"
	"github.com/homework/lab/internal/service/user/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_Login(t *testing.T) {
	defaultUser := &user.UserLogin{
		UserName: "testuser@gmail.com",
		Password: "123131242",
	}
	testCases := []struct {
		name string

		setupRepo    func(ctx *gin.Context) *mocks.UserService
		expectedFunc func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "UserName Not Exits",
			setupRepo: func(ctx *gin.Context) *mocks.UserService {
				userServiceMocks := mocks.NewUserService(t)
				userServiceMocks.On("Login", ctx, *defaultUser).Return("", user_service.ServiceErr.UserNameNotExistError)
				return userServiceMocks
			},
			expectedFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				code := rec.Code
				body := rec.Body.String()
				assert.Equal(t, http.StatusBadRequest, code)
				assert.Contains(t, body, user_service.ServiceErr.UserNameNotExistError.Error())
			},
		},

		{
			name: "Password Not Match",
			setupRepo: func(ctx *gin.Context) *mocks.UserService {
				userServiceMocks := mocks.NewUserService(t)
				userServiceMocks.On("Login", ctx, *defaultUser).Return("", user_service.ServiceErr.PasswordError)
				return userServiceMocks
			},
			expectedFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				code := rec.Code
				body := rec.Body.String()
				assert.Equal(t, http.StatusBadRequest, code)
				assert.Contains(t, body, user_service.ServiceErr.PasswordError.Error())
			},
		},

		{
			name: "Success",
			setupRepo: func(ctx *gin.Context) *mocks.UserService {
				userServiceMocks := mocks.NewUserService(t)
				userServiceMocks.On("Login", ctx, *defaultUser).Return("token_123", nil)
				return userServiceMocks
			},
			expectedFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				code := rec.Code
				body := rec.Body.String()
				assert.Equal(t, http.StatusOK, code)
				assert.Contains(t, body, "token_123")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			var bodyBytes, errDecode = json.Marshal(defaultUser)
			if errDecode != nil {
				testItem.Fatal(errDecode)
			}
			ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(bodyBytes))
			ctx.Request.Header.Set("Content-Type", "application/json")
			service := tc.setupRepo(ctx)
			handler := NewUserHandler(service)
			handler.Login(ctx)
			tc.expectedFunc(testItem, rec)
		})
	}

}
