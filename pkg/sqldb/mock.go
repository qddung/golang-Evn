package sqldb

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"gorm.io/gorm"
)

type MiniPostgres interface {
	DB() (*gorm.DB, error)
}

type miniPostgres struct {
	mu           sync.Mutex
	db           *gorm.DB // config postgres driver
	portListener uint32
}

func NewMiniPostgres(t *testing.T) (*gorm.DB, error) {
	db := &miniPostgres{db: nil, portListener: 0}
	err := db.Start(t)
	if err != nil {
		return nil, err
	}
	err = db.InstallExtensions()
	if err != nil {
		return nil, err
	}
	return db.DB()
}

func (_m *miniPostgres) InstallExtensions() error {
	_m.mu.Lock()
	defer _m.mu.Unlock()
	if _m.db == nil {
		return errors.New("not start")
	}
	tx := _m.db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	return tx.Error
}

func (_m *miniPostgres) DB() (*gorm.DB, error) {
	_m.mu.Lock()
	defer _m.mu.Unlock()
	if _m.db == nil {
		return nil, errors.New("not start")
	}
	return _m.db, nil
}

// Hàm helper dùng để tìm một cổng ngẫu nhiên còn trống trên Windows/Linux/Mac
func (_m *miniPostgres) getFreePort() (uint32, error) {
	// Dùng đúng câu lệnh net.Listen với cổng trống ":"
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return 0, err
	}
	// Lấy số cổng thực tế mà hệ điều hành vừa cấp
	port := listener.Addr().(*net.TCPAddr).Port

	// Phải đóng listener ngay lập tức để giải phóng cổng này cho Postgres sử dụng
	listener.Close()

	return uint32(port), nil
}

// Khởi động server Postgres ngầm, cho mock test schema
func (_m *miniPostgres) Start(t *testing.T) error {
	_m.mu.Lock()
	if _m.db != nil {
		_m.mu.Unlock()
		return nil
	}
	port, err := _m.getFreePort()
	if err != nil {
		_m.mu.Unlock()
		return err
	}
	dbLocation := rand.Text()
	config := embeddedpostgres.DefaultConfig().
		CachePath(fmt.Sprintf(".cache/%s/.embedded-postgres-go/", dbLocation)).
		RuntimePath(fmt.Sprintf(".cache/%s/.embedded-postgres-go/extracted", dbLocation)).
		DataPath(fmt.Sprintf(".cache/%s/.embedded-postgres-go/extracted/data", dbLocation)).
		BinariesPath(fmt.Sprintf(".cache/%s/.embedded-postgres-go/extracted", dbLocation)).
		Port(port)
	// config.RuntimePath();
	// 2. Khởi động Server Postgres ngầm
	pgServer := embeddedpostgres.NewDatabase(config)

	if err := pgServer.Start(); err != nil {
		_m.mu.Unlock()
		log.Fatalf("Không thể chạy Embedded Postgres: %v", err)
		message := err.Error()
		return errors.New(message)
	}
	_m.mu.Unlock()

	// 3. Sử dụng hàm có sẵn của thư viện để lấy chuỗi Connection URL tự động
	// Hàm này tự trả về định dạng: "postgres://my_user:my_password@localhost:PORT/my_db?sslmode=disable"
	dsn := config.GetConnectionURL()
	// 4. Dùng GORM kết nối thẳng vào chuỗi dsn vừa lấy được
	db, err := buildClient(dsn)
	if err != nil {
		return err
	}
	_m.db = db

	// end testcase
	t.Cleanup(func() {
		if err := os.RemoveAll(".cache"); err != nil {
			log.Printf("Error occurred while removing cache directory: %v", err)
		}
		pgServer.Stop()
	})
	return nil

}
