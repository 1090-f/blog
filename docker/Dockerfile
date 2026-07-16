# syntax=docker/dockerfile:1

FROM node:20-alpine AS frontend-builder
WORKDIR /src/web

COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.24-alpine AS backend-builder
WORKDIR /src

# Use the Go module proxy when building from networks with limited access to proxy.golang.org.
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
COPY --from=frontend-builder /src/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/blog .

FROM alpine:3.20
RUN addgroup -S app && adduser -S -G app app \
    && mkdir -p /app/uploads \
    && chown -R app:app /app
WORKDIR /app
COPY --from=backend-builder /out/blog ./blog
COPY --from=backend-builder /src/web/dist ./web/dist

USER app
EXPOSE 8080 8081
ENTRYPOINT ["./blog"]
