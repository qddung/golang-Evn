# GolangEVN

Đây là một ứng dụng Go đơn giản sử dụng Gin để cung cấp API health check. Dự án được tổ chức theo cấu trúc phân lớp nhằm dễ mở rộng và kiểm thử.

## Mục tiêu

- Cung cấp endpoint `/ping` để kiểm tra trạng thái dịch vụ.
- Tách riêng logic điều hướng, handler, service và cấu hình.
- Có test cho unit test và integration test.

## Cấu trúc thư mục

```text
.
├── cmd/
│   └── api/
│       └── main.go          # Điểm khởi động ứng dụng
├── internal/
│   ├── api/
│   │   └── health_check.go  # Cấu hình và khởi tạo router API
│   ├── config/
│   │   └── config.go        # Đọc cấu hình từ môi trường
│   ├── handler/
│   │   └── health_check.go  # Xử lý HTTP cho endpoint health check
│   ├── service/
│   │   └── health_check.go  # Logic nghiệp vụ cho health check
│   └── intergration_test/
│       └── endpoint/
│           └── health_check_test.go
```

## Công nghệ sử dụng

- Go
- Gin Web Framework
- Testify
- envconfig
- UUID

## Cài đặt

Yêu cầu:
- Go 1.25 hoặc mới hơn

Cài đặt dependency:

```bash
go mod tidy
```

## Chạy ứng dụng

```bash
go run ./cmd/api
```

Sau khi chạy, bạn có thể gọi endpoint:

```bash
curl http://localhost:8080/ping
```

Kết quả trả về sẽ có dạng JSON như:

```json
{
  "message": "OK",
  "service_name": "bookmark_service",
  "instance_id": "..."
}
```

## Test

Chạy unit test:

```bash
go test ./...
```

Chạy riêng integration test:

```bash
go test -v -run TestHealthCheck_Integration ./internal/intergration_test/endpoint
```

## Cấu hình môi trường

Dự án hỗ trợ một số biến môi trường như:

- `APP_PORT`: cổng chạy ứng dụng, mặc định `8080`
- `SERVICE_NAME`: tên dịch vụ, mặc định `bookmark_service`
- `INSTANCE_ID`: ID instance, nếu không cung cấp sẽ tự sinh ngẫu nhiên

## Ghi chú

Dự án này phù hợp để làm ví dụ về cách xây dựng một API Go nhỏ theo kiến trúc rõ ràng, có tách layer và có test. 
