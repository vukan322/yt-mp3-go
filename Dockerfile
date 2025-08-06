FROM node:20-alpine AS build-frontend
WORKDIR /app
COPY web/static ./web/static
RUN npm install -g esbuild && \
    esbuild web/static/css/main.css --bundle --minify --outfile=web/static/css/style.css && \
    esbuild web/static/js/main.js --bundle --minify --outfile=web/static/js/bundle.js

FROM golang:1.23-alpine AS build-go
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=build-frontend /app/web/static/css/style.css ./web/static/css/
COPY --from=build-frontend /app/web/static/js/bundle.js ./web/static/js/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates python3 py3-pip ffmpeg && \
    pip3 install --break-system-packages yt-dlp
WORKDIR /app
COPY --from=build-go /app/main .
COPY --from=build-go /app/web ./web
COPY --from=build-go /app/locales ./locales
RUN mkdir -p downloads
EXPOSE 8080
ENV APP_ENV=production DOMAIN=localhost PORT=8080 BASE_PATH=/yt-downloader
CMD ["./main"]
