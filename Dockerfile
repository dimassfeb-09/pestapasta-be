ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app

# Salin dan verifikasi modul Go
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Salin kode aplikasi dan lakukan build
COPY . .
RUN go build -v -o /run-app .

FROM debian:bookworm

# Instal paket sertifikat CA
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Salin aplikasi dari builder
COPY --from=builder /run-app /usr/local/bin/

# Jalankan aplikasi
CMD ["run-app"]
