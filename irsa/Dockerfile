FROM golang:1.23-alpine AS builder

WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY cmd cmd
COPY internal internal
COPY pkg pkg

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制编译好的应用
COPY --from=builder /app/main .

# 暴露应用端口
EXPOSE 8080

# 运行应用
CMD ["./main"]