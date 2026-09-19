# ---------- build stage ----------
FROM golang:1.22-alpine AS build
WORKDIR /app

# Устанавливаем всё, что нужно для сборки:
#   git        — для go mod download (replace-директивы на GitHub)
#   gcc        — для CGO-зависимостей (DataDog/zstd и т.п.)
#   musl-dev   — заголовки для gcc
RUN apk add --no-cache git gcc musl-dev

# Кешируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Собираем бинарник
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

# ---------- runtime stage ----------
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/server .
EXPOSE 8080
CMD ["./server"]