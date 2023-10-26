FROM docker.atl-paas.net/build/golang:1.20.6-fips AS BUILDER
RUN mkdir -p /go/src/github.com/asecurityteam/serverfull-gateway
WORKDIR /go/src/github.com/asecurityteam/serverfull-gateway
COPY . .
# RUN sdcli go dep
RUN go mod vendor
RUN go mod tidy
RUN GO111MODULE=on CGO_ENABLED=1 GOOS=linux go build -a -o /opt/app main.go
RUN go tool nm /opt/app | grep -i _Cfunc__goboringcrypto_
RUN go tool nm /opt/app | grep -i crypto/internal/boring/sig.BoringCrypto

##################################

FROM alpine:latest as CERTS
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .

###################################

FROM docker.atl-paas.net/sox/micros-golang:2.0.0
COPY --from=BUILDER /opt/app .
COPY --from=CERTS /zoneinfo.zip /
COPY --from=CERTS /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV ZONEINFO /zoneinfo.zip
USER root
RUN chown -R appuser:appuser /opt/service/*
USER appuser
ENTRYPOINT ["/app"]
