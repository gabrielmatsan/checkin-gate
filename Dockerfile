FROM golang:1.25.6-alpine AS base

FROM base AS dependencies
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

FROM base AS builder
WORKDIR /usr/src/app
COPY --from=dependencies /go/pkg /go/pkg
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -o ./api ./cmd/api/main.go

FROM alpine:3.21 AS runner
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -u 1001 appuser
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/api .
USER appuser
EXPOSE 8080
ENTRYPOINT ["./api"]
