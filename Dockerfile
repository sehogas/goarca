FROM alpine:3.21.3 AS root-certs
RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1001 app
RUN adduser app -u 1001 -D -G app 

FROM golang:1.24.1-alpine@sha256:43c094ad24b6ac0546c62193baeb3e6e49ce14d3250845d166c77c25f64b0386 AS builder
WORKDIR /build
COPY --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=0.1.0
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
ENV FLAGS="-s -w -X 'main.Version=${VERSION}'"
RUN go build -ldflags="${FLAGS}" -o ./api-arca ./cmd/api/.

FROM scratch
COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:1001 --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=1001:1001 --from=builder /build/api-arca /
COPY --chown=1001:1001 --from=builder /build/data /data
USER app
EXPOSE 4433
ENTRYPOINT ["/api-arca"]