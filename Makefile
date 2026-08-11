.PHONY: dev-run swagger run-app test



dev-run:
	go run github.com/swaggo/swag/cmd/swag init -g ./cmd/api/main.go --output docs
	go run cmd/api/main.go

swagger:
	go run github.com/swaggo/swag/cmd/swag init -g ./cmd/api/main.go --output docs

run-app:
	go run cmd/api/main.go

IMG_NAME=dungi3/golang-learn-bookmark_service
GIT_TAG := $(shell git describe --tags --exact-match --abbrev=0 2>/dev/null)
IMG_TAG := latest
ifneq ($(GIT_TAG),)
	IMG_TAG := $(GIT_TAG)
endif

testdir := ./test

testdir-check:
	@if [ ! -d "$(testdir)" ]; then \
		echo "Directory does not exist. Creating..."; \
		mkdir -p "$(testdir)"; \
	fi
# 1. Định nghĩa giá trị mặc định (người dùng có thể đè bằng: make test OPTION=cache)
OPTION ?= cache
COVERAGE_EXCLUDE=mocks|main.go|test|config.go
COVERAGE_THRESHOLD ?= 80

# 2. Xử lý logic điều kiện (Sử dụng cú pháp ifeq/ifneq chuẩn của Make)
ifeq ($(OPTION),cache)
    # Giữ nguyên cache mặc định của Go
    CACHECMD := go test ./... -coverprofile=./test/coverage_tmp -covermode=atomic -coverpkg=./... -p 1
else
    # Thêm biến môi trường tắt cache (Goflags hoặc -count=1)
    CACHECMD := GOFLAGS="-count=1" go test ./... -coverprofile=./test/coverage_tmp -covermode=atomic -coverpkg=./... -p 1
endif

testdir-check:
	@mkdir -p ./test

test: testdir-check
	mkdir -p .cache
	chmod 700 .cache

	# 3. Thực thi lệnh test đã cấu hình
	$(CACHECMD)
	
	grep -vE "$(COVERAGE_EXCLUDE)" ./test/coverage_tmp > ./test/coverage_out
	go tool cover -html=./test/coverage_out -o ./test/coverage.html
	@total=$$(go tool cover -func=./test/coverage_out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi


docker-up:
	docker compose -f deployment/docker-compose.yml up -d 

docker-build:
	docker build -f deployment/Dockerfile -t $(IMG_NAME):$(IMG_TAG) .

docker-release:
	docker push $(IMG_NAME):$(IMG_TAG)