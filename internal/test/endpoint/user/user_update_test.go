package user_endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	user_entity "github.com/homework/lab/internal/models/entity"
	user_repository "github.com/homework/lab/internal/repository/user"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
	"github.com/homework/lab/pkg/helpers"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateUserInfo_Integration(t *testing.T) {
	eng, mockJwt, conn, _ := setupEngineWithDB(t)
	db := conn.GetSqlDB()

	// create a user in DB
	hasher := helpers.NewHasher()
	pw, err := hasher.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user := &user_entity.User{DisplayName: "Integration", Email: "upd@example.com", Password: pw, UserName: "upduser"}
	if err := user_repository.NewUserRepository(db).CreateUser(context.Background(), user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	// generate token
	claims := map[string]interface{}{"sub": user.Id}
	token, err := mockJwt.JwtGenarate.GenerateToken(claims)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	// prepare update body
	updateBody := map[string]string{"username": "updated-name"}
	bodyBytes, _ := json.Marshal(updateBody)

	// send PUT request
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	eng.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Logf("update response body: %s", rec.Body.String())
	}

	// do a GET to verify change
	rec2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/v1/self/info", nil)
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
	assert.Equal(t, "updated-name", resp.Data.UserName)
}
