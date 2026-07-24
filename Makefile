.PHONY: dev-run swagger run-app test

dev-run:
	go run github.com/swaggo/swag/cmd/swag init -g ./cmd/api/main.go --output docs
	go run cmd/api/main.go

swagger:
	go run github.com/swaggo/swag/cmd/swag init -g ./cmd/api/main.go --output docs

run-app:
	go run cmd/api/main.go

COVERAGE_EXCLUDE=mocks|main.go|test
COVERAGE_THRESHOLD = 80

IMG_NAME=dungi3/golang-learn-bookmark_service
IMG_TAG := latest

test:
	go test ./... -coverprofile=./test/coverage_tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" ./test/coverage_tmp > ./test/coverage_out
	go tool cover -html=./test/coverage_out -o coverage.html
	@total=$$(go tool cover -func=./test/coverage_out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi



docker-build:
	docker build -t $(IMG_NAME):$(IMG_TAG) .

docker-release:
	docker push $(IMG_NAME):$(IMG_TAG)