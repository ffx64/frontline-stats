export GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    TZ=America/Sao_Paulo

go mod download
go build -o gamestats_backend ./cmd/main.go
