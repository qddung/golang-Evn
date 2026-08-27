package user_endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	userModel "github.com/homework/lab/internal/models/dto/api/user"
	user_entity "github.com/homework/lab/internal/models/entity"
	user_repository "github.com/homework/lab/internal/repository/user"
	"github.com/homework/lab/pkg/helpers"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateUserInfo_Integration(t *testing.T) {
	testCases := []struct {
		name           string
		updateBody     map[string]string
		expectedStatus int
		expectedUser   string
	}{
		{
			name:           "update username only",
			updateBody:     map[string]string{"username": "updated-name"},
			expectedStatus: http.StatusOK,
			expectedUser:   "updated-name",
		},
		{
			name:           "update username and password",
			updateBody:     map[string]string{"username": "updated-name-2", "password": "newpassword123"},
			expectedStatus: http.StatusOK,
			expectedUser:   "updated-name-2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			eng, mockJwt, conn, _ := setupEngineWithDB(t)
			db := conn.GetSqlDB()

			hasher := helpers.NewHasher()
			pw, err := hasher.HashPassword("password123")
			if err != nil {
				t.Fatalf("hash password: %v", err)
			}
			user := &user_entity.User{DisplayName: "Integration", Email: "upd@example.com", Password: pw, UserName: "upduser"}
			if err := user_repository.NewUserRepository(db).CreateUser(context.Background(), user); err != nil {
				t.Fatalf("create user: %v", err)
			}

			claims := map[string]interface{}{"sub": user.Id}
			token, err := mockJwt.JwtGenarate.GenerateToken(claims)
			if err != nil {
				t.Fatalf("generate token: %v", err)
			}

			bodyBytes, err := json.Marshal(tc.updateBody)
			if err != nil {
				t.Fatalf("marshal update body: %v", err)
			}

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatalf("new put request: %v", err)
			}
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			eng.ServeHTTP(rec, req)
			if rec.Code != tc.expectedStatus {
				t.Logf("update response body: %s", rec.Body.String())
			}
			assert.Equal(t, tc.expectedStatus, rec.Code)
			if rec.Code != http.StatusOK {
				return
			}

			rec2 := httptest.NewRecorder()
			req2, err := http.NewRequest(http.MethodGet, "/v1/self/info", nil)
			if err != nil {
				t.Fatalf("new get request: %v", err)
			}
			req2.Header.Set("Authorization", "Bearer "+token)

			eng.ServeHTTP(rec2, req2)
			if rec2.Code != http.StatusOK {
				t.Logf("get response body: %s", rec2.Body.String())
			}
			assert.Equal(t, http.StatusOK, rec2.Code)

			var resp struct {
				Data userModel.UserInfo `json:"data"`
			}
			if err := json.Unmarshal(rec2.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal resp: %v", err)
			}
			assert.Equal(t, tc.expectedUser, resp.Data.UserName)
		})
	}
}
