FROM golang:alpine as build

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .
RUN apk add --no-cache tzdata
RUN apk add --no-cache make
RUN apk add --no-cache openssl
RUN apk add build-base
RUN go mod download && \
GOOS=linux go build -tags musl -ldflags "-w -s" -o bookmark_service cmd/api/main.go
RUN make generated-key

FROM alpine:latest

ARG app_name=app
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=build /opt/app/privatekey.pem /app/privatekey.pem
COPY --from=build /opt/app/publickey.pem /app/publickey.pem
COPY --from=build /opt/app/bookmark_service /app/bookmark_service
COPY --from=build /opt/app/docs /app/docs
COPY --from=build /opt/app/migrations /app/migrations
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

RUN ln -snf /usr/share/zoneinfo/"$TZ" /etc/localtime && echo "$TZ" > /etc/timezone

CMD ["/app/bookmark_service"]