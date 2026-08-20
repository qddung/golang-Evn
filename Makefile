.PHONY: dev-run swagger run-app test


swagger:
	go run github.com/swaggo/swag/cmd/swag init --parseDependency --parseInternal -g ./cmd/api/main.go --output docs

dev-run: swagger
	go run github.com/joho/godotenv/cmd/godotenv go run cmd/api/main.go

run-app:
	go run cmd/api/main.go


# 1. Định nghĩa giá trị mặc định (người dùng có thể đè bằng: make test OPTION=cache)
OPTION ?= cache
COVERAGE_EXCLUDE=mocks|main.go|test|config.go|infrastructure/**
COVERAGE_THRESHOLD ?= 80

# 2. Xử lý logic điều kiện (Sử dụng cú pháp ifeq/ifneq chuẩn của Make)
CACHECMD := 
ifeq ($(OPTION),nocache)
    CACHECMD := go clean -testcache    
endif

test: 
	@mkdir -p ./test
	# run clean test cache
	$(CACHECMD)

	go test ./... -coverprofile=./test/coverage_tmp -covermode=atomic -coverpkg=./... -p 1	
	grep -vE "$(COVERAGE_EXCLUDE)" ./test/coverage_tmp > ./test/coverage_out
	go tool cover -html=./test/coverage_out -o ./test/coverage.html
	@total=$$(go tool cover -func=./test/coverage_out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

test-nocache: test OPTION=nocache

docker-up:
	docker compose -f deployment/docker-compose.yml up -d 

rebuild:
	docker compose -f deployment/docker-compose.yml build


IMG_NAME=dungi3/golang-learn-bookmark_service
GIT_TAG := $(shell git describe --tags --exact-match --abbrev=0 2>/dev/null)
IMG_TAG := latest
ifneq ($(GIT_TAG),)
	IMG_TAG := $(GIT_TAG)
endif

docker-build:
	docker build -f deployment/Dockerfile -t $(IMG_NAME):$(IMG_TAG) .

docker-release:
	docker push $(IMG_NAME):$(IMG_TAG)

generated-key:
	openssl genrsa -out privatekey.pem 2048
	openssl rsa -in privatekey.pem -pubout -out publickey.pem
