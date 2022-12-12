# AWS S3 Proxy v2.0
# docker run -d -p 8080:80 -e AWS_REGION -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e AWS_S3_BUCKET pottava/s3-proxy

FROM golang:1.13.7-alpine3.10 AS builder
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache gcc musl-dev git
RUN git config --global http.version HTTP/1.1
RUN go env -w GOPROXY="https://goproxy.cn,direct"
ADD . aws-s3-proxy
WORKDIR aws-s3-proxy
ENV APP_VERSION=v2.3.0
#RUN git checkout "${APP_VERSION}" > /dev/null 2>&1
RUN go mod download
RUN go mod verify
RUN githash=$(git rev-parse --short HEAD 2>/dev/null) \
    && today=$(date +%Y-%m-%d --utc) \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags '-s -w -X main.ver=${APP_VERSION} -X main.commit=${githash} -X main.date=${today}' \
    -o /app

FROM alpine:3.11 AS libs
RUN apk --no-cache add ca-certificates

FROM alpine:3.10

ADD mc_client/mc /bin
RUN chmod +x /bin/mc
COPY --from=libs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /aws-s3-proxy
CMD ["/aws-s3-proxy"]
