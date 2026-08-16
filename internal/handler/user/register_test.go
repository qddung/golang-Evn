package user_handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/homework/lab/internal/models/dto/api"
	"github.com/homework/lab/internal/models/dto/api/user"
	user_service "github.com/homework/lab/internal/service/user"
	"github.com/homework/lab/internal/service/user/mocks"
	"github.com/stretchr/testify/assert"
)

var testErr = errors.New("test error")

func TestService_Register(t *testing.T) {
	testCases := []struct {
		name string

		setupRepo    func(ctx *gin.Context, info user.UserRegister) *mocks.UserService
		input        *user.UserRegister
		expectedFunc func(t *testing.T, rec *httptest.ResponseRecorder, registerInput *user.UserRegister)
	}{
		{
			// Duplicate email
			name: "Duplicate email",
			input: &user.UserRegister{
				DisplayName: "test",
				Email:       "test@example.com",
				Password:    "123131242",
				UserName:    "testuser",
			},
			setupRepo: func(ctx *gin.Context, info user.UserRegister) *mocks.UserService {
				userServiceMocks := mocks.NewUserService(t)
				userServiceMocks.On("Register", ctx, info).Return(nil, user_service.ServiceErr.EmailExistError)

				return userServiceMocks
			},
			expectedFunc: func(t *testing.T, rec *httptest.ResponseRecorder, registerInput *user.UserRegister) {
				code := rec.Code
				body := rec.Body.String()
				assert.Equal(t, http.StatusConflict, code)
				assert.Contains(t, body, user_service.ServiceErr.EmailExistError.Error())
			},
		},
		{
			// Duplicate username
			name: "Duplicate username",
			input: &user.UserRegister{
				DisplayName: "test",
				Email:       "test@example.com",
				Password:    "123131242",
				UserName:    "testuser",
			},
			setupRepo: func(ctx *gin.Context, info user.UserRegister) *mocks.UserService {
				userServiceMocks := mocks.NewUserService(t)
				userServiceMocks.On("Register", ctx, info).Return(nil, user_service.ServiceErr.UserNameExistError)
				return userServiceMocks
			},
			expectedFunc: func(t *testing.T, rec *httptest.ResponseRecorder, registerInput *user.UserRegister) {
				code := rec.Code
				body := rec.Body.String()

				assert.Equal(t, http.StatusConflict, code)
				assert.Contains(t, body, user_service.ServiceErr.UserNameExistError.Error())
			},
		},
		{
			name: "Create user error",
			input: &user.UserRegister{
				DisplayName: "test",
				Email:       "test@example.com",
				Password:    "123131242",
				UserName:    "testuser",
			},
			setupRepo: func(ctx *gin.Context, info user.UserRegister) *mocks.UserService {
				userServiceMocks := mocks.NewUserService(t)
				userServiceMocks.On("Register", ctx, info).Return(nil, testErr)
				return userServiceMocks
			},
			expectedFunc: func(t *testing.T, rec *httptest.ResponseRecorder, registerInput *user.UserRegister) {
				code := rec.Code
				body := rec.Body.String()

				assert.Equal(t, http.StatusInternalServerError, code)
				assert.Contains(t, body, testErr.Error())
			},
		},
		{
			name: "Create user successfully",
			input: &user.UserRegister{
				DisplayName: "test",
				Email:       "test@example.com",
				Password:    "123131242",
				UserName:    "testuser",
			},
			setupRepo: func(ctx *gin.Context, info user.UserRegister) *mocks.UserService {
				userServiceMocks := mocks.NewUserService(t)
				userServiceMocks.On("Register", ctx, info).Return(&user.UserInfo{
					Id:          uuid.NewString(),
					DisplayName: "test",
					Email:       "test@example.com",
					UserName:    "testuser",
					UpdateAt:    time.Now().String(),
					CreateAt:    time.Now().String(),
				}, nil)
				return userServiceMocks
			},
			expectedFunc: func(t *testing.T, rec *httptest.ResponseRecorder, registerInput *user.UserRegister) {
				code := rec.Code
				body := rec.Body.String()
				res := &api.Response[user.UserInfo]{}
				err := json.Unmarshal([]byte(body), res)
				if err != nil {
					t.Fatal(err)
					return
				}
				userInfo := res.Data
				assert.Equal(t, http.StatusOK, code)
				assert.Contains(t, res.Message, createSuccessMessage)
				assert.Equal(t, registerInput.DisplayName, userInfo.DisplayName)
				assert.Equal(t, registerInput.UserName, userInfo.UserName)
				assert.Equal(t, registerInput.Email, userInfo.Email)
				assert.NotEmpty(t, userInfo.Id)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			var bodyBytes, errDecode = json.Marshal(tc.input)
			if errDecode != nil {
				testItem.Fatal(errDecode)
				return
			}
			// userDecode := &user.UserRegister{}
			// json.Unmarshal(bodyBytes, userDecode)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(bodyBytes))
			ctx.Request.Header.Set("Content-Type", "application/json")
			u := tc.input
			service := tc.setupRepo(ctx, *u)
			handler := NewUserHandler(service)
			handler.Register(ctx)
			tc.expectedFunc(testItem, rec, u)
		})
	}

}
