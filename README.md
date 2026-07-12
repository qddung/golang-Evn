# GolangEVN - URL Shortener & Health Check Service

Đây là một dịch vụ backend viết bằng **Go** sử dụng **Gin Framework** và **Redis**, theo kiến trúc phân lớp (Clean Architecture / Layered Architecture). Dịch vụ cung cấp các API để rà soát sức khỏe hệ thống (Health Check), rút gọn liên kết (URL Shortener), chuyển hướng liên kết, cùng tài liệu API trực quan với **Swagger**.

---

## 🎯 Tính năng chính

1. **Health Check (`/ping`)**:
   - Kiểm tra tình trạng hoạt động của dịch vụ và kết nối Redis.
   - Trả về thông tin `service_name` và `instance_id`.

2. **URL Shortener (`/v1/links/shorten`)**:
   - Tạo mã rút gọn 6 ký tự ngẫu nhiên cho một URL dài.
   - Hỗ trợ thiết lập thời gian hết hạn (`exp`) tối đa 7 ngày (`604800` giây).
   - Tự động kiểm tra và xử lý xung đột mã ngẫu nhiên trong Redis.

3. **URL Redirect (`/v1/links/redirect/:code`)**:
   - Chuyển hướng người dùng từ mã rút gọn sang URL gốc.

4. **Tài liệu Swagger (`/swagger/*any`)**:
   - Xem và kiểm thử API trực tiếp trên giao diện Swagger UI.

---

## 📁 Cấu trúc thư mục

```text
.
├── cmd/
│   └── api/
│       └── main.go                  # Điểm khởi động ứng dụng (Entry point)
├── docs/                            # Tài liệu Swagger tự tạo từ chú thích API
├── internal/
│   ├── api/
│   │   └── api.go                   # Khởi tạo Engine, định tuyến Routes và Dependency Injection
│   ├── config/
│   │   └── config.go                # Quản lý cấu hình ứng dụng từ biến môi trường
│   ├── handler/
│   │   ├── health_check.go          # HTTP Handler cho Health Check
│   │   └── shorten_url.go           # HTTP Handler cho URL Shortener
│   ├── repository/
│   │   ├── ping.go                  # Giao tiếp Redis (Ping)
│   │   ├── url_storage.go           # Giao tiếp Redis (Lưu trữ và truy xuất URL rút gọn)
│   │   └── mocks/                   # Mocks tự tạo cho Repository (sử dụng mockery)
│   ├── service/
│   │   ├── health_check.go          # Logic nghiệp vụ Health Check
│   │   ├── shortern_url.go          # Logic nghiệp vụ rút gọn URL
│   │   └── mocks/                   # Mocks cho Service layer
│   └── intergration_test/
│       └── endpoint/
│           ├── health_check_test.go # Integration test cho endpoint /ping
│           └── shorten_test.go      # Integration test cho endpoint /v1/links/shorten & redirect
├── pkg/
│   ├── helpers/
│   │   └── code_gen.go              # Bộ sinh mã ký tự ngẫu nhiên
│   ├── redis/
│   │   ├── config.go                # Cấu hình kết nối Redis
│   │   ├── conn.go                  # Khởi tạo Redis Client
│   │   └── mock.go                  # Khởi tạo In-memory Redis (miniredis) phục vụ kiểm thử
│   └── response/
│       └── response.go              # Mẫu cấu trúc phản hồi lỗi chuẩn
└── Makefile                         # Các lệnh hỗ trợ tiện ích (dev, swagger, test, coverage)
```

---

## 🛠 Công nghệ sử dụng

- **Ngôn ngữ**: Go 1.25+
- **Web Framework**: Gin Web Framework (`github.com/gin-gonic/gin`)
- **Database / Cache**: Redis (`github.com/redis/go-redis/v9`)
- **Khởi tạo và Quản lý cấu hình**: `github.com/kelseyhightower/envconfig` & `github.com/google/uuid`
- **Tài liệu API**: Swaggo (`swag`, `gin-swagger`, `swagger-files`)
- **Kiểm thử & Mocking**:
  - `github.com/stretchr/testify` (Assertions & Mocks)
  - `github.com/alicebob/miniredis/v2` (In-memory Redis cho Integration & Unit Tests)
  - `mockery` (Tự động sinh mock code)

---

## 🚀 Hướng dẫn cài đặt & chạy ứng dụng

### 1. Yêu cầu hệ thống
- **Go**: 1.25 trở lên
- **Redis**: Chạy sẵn trên cổng mặc định `6379` (hoặc cấu hình qua biến môi trường)

### 2. Cài đặt các gói phụ thuộc
```bash
go mod tidy
```

### 3. Chạy ứng dụng

#### **Chạy nhanh bằng Makefile** (Tự động sinh lại Swagger và khởi động Server):
```bash
make dev-run
```

#### **Chạy bằng lệnh Go cơ bản**:
```bash
go run ./cmd/api/main.go
```
*Mặc định ứng dụng sẽ khởi động trên `http://localhost:8080`.*

---

## 📖 Tài liệu API & Ví dụ sử dụng

Sau khi khởi động server, bạn có thể truy cập giao diện **Swagger UI** tại:
👉 `http://localhost:8080/swagger/index.html`

### 1. Health Check
```bash
curl -X GET http://localhost:8080/ping
```
**Phản hồi `200 OK`**:
```json
{
  "message": "OK",
  "service_name": "book",
  "instance_id": "c1f2b3a4-..."
}
```

### 2. Rút gọn liên kết (Shorten URL)
```bash
curl -X POST http://localhost:8080/v1/links/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.google.com",
    "exp": 3600
  }'
```
**Phản hồi `200 OK`**:
```json
{
  "message": "Shorten URL generated successfully!",
  "code": "aBcDeF"
}
```

### 3. Chuyển hướng liên kết (Redirect)
Truy cập trực tiếp trên trình duyệt hoặc qua `curl`:
```bash
curl -i http://localhost:8080/v1/links/redirect/aBcDeF
```
**Phản hồi `302 Found`**: Chuyển hướng đến `https://www.google.com`.

---

## 🧪 Kiểm thử (Testing & Coverage)

Dự án được bao phủ tốt bởi cả **Unit Tests** và **Integration Tests** độc lập, sử dụng `miniredis` trong bộ nhớ nên **không đòi hỏi phải khởi động Redis thật khi chạy test**.

### Chạy toàn bộ các bộ test:
```bash
go test ./...
```

### Chạy riêng các Integration Test cho Endpoints:
```bash
go test -v ./internal/intergration_test/endpoint/...
```

### Chạy kiểm thử kèm đánh giá độ bao phủ mã nguồn (Coverage):
Có thể dùng Makefile:
```bash
make test
```

---

## ⚙️ Cấu hình biến môi trường

Các thông số có thể thay đổi bằng cách đặt biến môi trường (hoặc dùng file `.env` nếu tích hợp):

| Biến môi trường | Mô tả | Mặc định |
| :--- | :--- | :--- |
| `APP_PORT` | Cổng dịch vụ HTTP Server | `8080` |
| `SERVICE_NAME` | Tên dịch vụ hiển thị khi gọi health check | `book` |
| `INSTANCE_ID` | Mã định danh của instance (tự sinh ngẫu nhiên nếu để trống) | `UUID v4 ngẫu nhiên` |
| `REDIS_ADDR` | Địa chỉ máy chủ Redis | `localhost:6379` |
| `REDIS_PWD` | Mật khẩu xác thực Redis | `""` |
| `REDIS_DB` | Chỉ số Database Redis sử dụng | `0` |

---

## 💡 Ghi chú kiến trúc
- **Layered Clean Architecture**:
  - `API Layer (`internal/api`)`: Khởi tạo và liên kết các layer (Dependency Injection).
  - `Handler Layer (`internal/handler`)`: Tiếp nhận yêu cầu HTTP, kiểm tra tính hợp lệ dữ liệu vào (`binding`), gọi xuống Service.
  - `Service Layer (`internal/service`)`: Xử lý logic nghiệp vụ, xử lý xung đột, quy định hành vi.
  - `Repository Layer (`internal/repository`)`: Giao tiếp trực tiếp với cơ sở dữ liệu / cache (Redis).
