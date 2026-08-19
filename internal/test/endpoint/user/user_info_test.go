package user_endpoint

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/connection"
	user_entity "github.com/homework/lab/internal/models/entity"
	user_repository "github.com/homework/lab/internal/repository/user"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
	"github.com/homework/lab/pkg/sqldb"
	"github.com/homework/lab/pkg/helpers"
	"github.com/stretchr/testify/assert"
)

func setupEngineWithDB(t *testing.T) (api.Engine, *jwt_pkg.MockJwt, connection.DBConnector, *gin.Engine) {
	// prepare in-memory DB
	db, err := sqldb.NewMiniPostgres(t)
	if err != nil {
		t.Fatalf("failed to create mini db: %v", err)
	}
	// migrate user table
	if err := db.AutoMigrate(&user_entity.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// connector
	conn := connection.NewDBConnector(nil, db)

	// jwt mock
	mockJwt := jwt_pkg.NewMockJwt()

	// engine
	r := gin.New()
	cfg := &config.Config{AppPort: "8080", BasePath: "/", ServiceName: "testsvc"}
	eng := api.NewEngine(&api.EnginOpt{App: r, Cfg: cfg, Connector: conn, JwtGenerator: mockJwt.JwtGenarate, JwtValidator: mockJwt.JwtValidate})
	return eng, mockJwt, conn, r
}

func Test_GetUserInfo_Integration(t *testing.T) {
	testCases := []struct {
		name         string
		displayName  string
		email        string
		userName     string
		expectedID   string
		expectedMail string
	}{
		{
			name:         "user info success",
			displayName:  "Integration",
			email:        "int@example.com",
			userName:     "intuser",
			expectedMail: "int@example.com",
		},
		{
			name:         "user info success alt user",
			displayName:  "Another User",
			email:        "alt@example.com",
			userName:     "altuser",
			expectedMail: "alt@example.com",
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
			user := &user_entity.User{
				DisplayName: tc.displayName,
				Email:       tc.email,
				Password:    pw,
				UserName:    tc.userName,
			}
			if err := user_repository.NewUserRepository(db).CreateUser(context.Background(), user); err != nil {
				t.Fatalf("create user: %v", err)
			}

			claims := map[string]interface{}{"sub": user.Id}
			token, err := mockJwt.JwtGenarate.GenerateToken(claims)
			if err != nil {
				t.Fatalf("generate token: %v", err)
			}

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/v1/self/info", nil)
			if err != nil {
				t.Fatalf("new get request: %v", err)
			}
			req.Header.Set("Authorization", "Bearer "+token)

			eng.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Logf("response body: %s", rec.Body.String())
			}
			assert.Equal(t, http.StatusOK, rec.Code)

			var resp struct {
				Data userModel.UserInfo `json:"data"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("unmarshal resp: %v", err)
			}
			assert.Equal(t, user.Id, resp.Data.Id)
			assert.Equal(t, tc.expectedMail, resp.Data.Email)
		})
	}
}
