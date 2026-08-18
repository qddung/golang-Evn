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
	eng, mockJwt, conn, _ := setupEngineWithDB(t)
	db := conn.GetSqlDB()

	// create a user in DB
	hasher := helpers.NewHasher()
	pw, err := hasher.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user := &user_entity.User{DisplayName: "Integration", Email: "int@example.com", Password: pw, UserName: "intuser"}
	if err := user_repository.NewUserRepository(db).CreateUser(context.Background(), user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	// generate token with subject = user.Id
	claims := map[string]interface{}{"sub": user.Id}
	token, err := mockJwt.JwtGenarate.GenerateToken(claims)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/self/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	eng.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Logf("response body: %s", rec.Body.String())
	}
	assert.Equal(t, http.StatusOK, rec.Code)

	// parse response into DTO
	var resp struct {
		Data userModel.UserInfo `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal resp: %v", err)
	}
	assert.Equal(t, user.Id, resp.Data.Id)
	assert.Equal(t, user.Email, resp.Data.Email)
}
