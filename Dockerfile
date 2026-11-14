ARG GO_VERSION=1.27.4
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app

# Salin dan verifikasi modul Go
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Salin kode aplikasi dan lakukan build
COPY . .
RUN go build -v -o /run-app .

FROM debian:bookworm-slim

# Salin sertifikat CA dari builder (menghindari apt-get saat build)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Salin aplikasi dari builder
COPY --from=builder /run-app /usr/local/bin/

EXPOSE 8080

# Jalankan aplikasi
CMD ["run-app"]
