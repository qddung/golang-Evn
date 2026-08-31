package infrastructure

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/constant"
	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/connection"

	"github.com/homework/lab/internal/models/entity"
	"github.com/homework/lab/pkg/sqldb"
)

func CreateApi() api.Engine {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// create redis client
	rdClient := CreateRedisClient()

	// create SQL client
	sqlClient := CreateSqlClient()

	if !sqlClient.Migrator().HasTable(&entity.User{}) {
		log.Println("Table does not exist. Running AutoMigrate...")
		if errMigrate := sqlClient.AutoMigrate(&entity.User{}); errMigrate != nil {
			log.Fatalf("Migration failed: %v", errMigrate)
		}
	} else {
		log.Println("Table already exists. Skipping migration pass.")
	}

	migrator := sqldb.BuildMigrate(sqlClient, constant.MigrationPath)
	err = migrator.MigrateUp()
	if err != nil {
		panic("start migration failed " + err.Error())
	}

	// connector
	connector := connection.NewDBConnector(rdClient, sqlClient)
	jwtGenerator, jwtValidator := CreateJwtProvider()
	apiEngine := api.NewEngine(&api.EnginOpt{
		App:          gin.New(),
		Cfg:          cfg,
		Connector:    connector,
		JwtGenerator: jwtGenerator,
		JwtValidator: jwtValidator,
	})
	return apiEngine
}
