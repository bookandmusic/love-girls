FROM node:24-alpine AS frontend-builder

WORKDIR /build/frontend

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile

COPY frontend/ ./
RUN pnpm build

FROM golang:1.25-alpine AS backend-builder

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

WORKDIR /build/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
COPY --from=frontend-builder /build/frontend/dist ./internal/handler/assets/dist
RUN apk add --no-cache gcc musl-dev sqlite-dev && \
    CGO_ENABLED=1 GOOS=linux \
    go build -a -trimpath \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" \
    -o love-girl .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata sqlite-libs

WORKDIR /app

COPY --from=backend-builder /build/backend/love-girl .

EXPOSE 8182

ENV APP_ENV=prod
ENV TZ=Asia/Shanghai
ENV SERVER_ADDR=:8182
ENV SERVER_HOST_NAME=localhost:8182

CMD ["/app/love-girl"]
