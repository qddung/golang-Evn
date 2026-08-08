package sqldb

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type globalServerManager struct {
	once       sync.Once
	server     *embeddedpostgres.EmbeddedPostgres
	config     embeddedpostgres.Config
	startError error
	dbLock     sync.Mutex

	// CẢI TIẾN: Giữ kết nối Root DB luôn mở để dùng chung, tránh mở/đóng liên tục
	rootDB *gorm.DB
}

var serverManager globalServerManager

func (sm *globalServerManager) startOnce() (embeddedpostgres.Config, *gorm.DB, error) {
	sm.once.Do(func() {
		port, err := sm.getFreePort()
		if err != nil {
			sm.startError = fmt.Errorf("failed to get free port: %w", err)
			return
		}

		sm.config = embeddedpostgres.DefaultConfig().
			CachePath(".cache/global/.embedded-postgres-go/").
			RuntimePath(".cache/global/.embedded-postgres-go/extracted").
			DataPath(".cache/global/.embedded-postgres-go/extracted/data").
			BinariesPath(".cache/global/.embedded-postgres-go/extracted").
			Port(port)

		sm.server = embeddedpostgres.NewDatabase(sm.config)
		if err := sm.server.Start(); err != nil {
			sm.startError = fmt.Errorf("cannot start embedded postgres: %w", err)
			return
		}

		// Khởi tạo kết nối Root DB duy nhất một lần tại đây
		rootDSN := sm.config.GetConnectionURL()
		sm.rootDB, err = gorm.Open(postgres.Open(rootDSN), &gorm.Config{
			SkipDefaultTransaction: true, // Tăng tốc độ thực thi lệnh đơn
		})
		if err != nil {
			sm.startError = fmt.Errorf("failed to connect to root db: %w", err)
			return
		}
	})
	return sm.config, sm.rootDB, sm.startError
}

func (sm *globalServerManager) getFreePort() (uint32, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return uint32(port), nil
}

func randomDBName() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return "db_" + hex.EncodeToString(bytes)
}

func NewMiniPostgresEmbedded(t *testing.T) (*gorm.DB, error) {
	// Lấy config và rootDB đã được khởi tạo sẵn
	config, rootDB, err := serverManager.startOnce()
	if err != nil {
		return nil, err
	}

	dbName := randomDBName()

	// Khóa để tạo database tuần tự tránh lỗi xung đột hệ thống của Postgres
	serverManager.dbLock.Lock()
	if err := rootDB.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName)).Error; err != nil {
		serverManager.dbLock.Unlock()
		return nil, fmt.Errorf("failed to create isolated database: %w", err)
	}
	serverManager.dbLock.Unlock()

	config.Database(dbName)
	// Kết nối vào DB độc lập vừa tạo
	isolatedDSN := config.GetConnectionURL()

	db, err := gorm.Open(postgres.Open(isolatedDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect isolated db via gorm: %w", err)
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return nil, fmt.Errorf("failed to install extension: %w", err)
	}

	// CẢI TIẾN: Chỉ đóng kết nối của chính test case đó khi kết thúc.
	// Bỏ hoàn toàn lệnh DROP DATABASE để tiết kiệm thời gian I/O.
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return db, nil
}
