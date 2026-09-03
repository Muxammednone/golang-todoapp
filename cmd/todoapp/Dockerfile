FROM golang:1.26.5-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOMAXPROCS=1 go build -p=1 \
    -o /app/cmd/todoapp/exe /app/cmd/todoapp

FROM alpine:3.24
WORKDIR /app
COPY --from=builder /app/cmd/todoapp/exe /app
COPY --from=builder /app/public /app/public
ENV PROJECT_ROOT=/app
CMD [ "/app/exe" ]
