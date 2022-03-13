ARG GOPROXY=https://goproxy.cn,direct

# STEP 1 build executable binary from alpine:golang

FROM golang:1.16-alpine AS builder
LABEL maintainer="EthanWang <git@ethan.wang>"

ENV GOPROXY=${GOPROXY}

RUN apk add --update --no-cache build-base musl-dev linux-headers git \
    pkgconf g++ gmp-dev openssl-dev vips vips-dev
ENV CGO_CFLAGS_ALLOW=-Xpreprocessor

WORKDIR /src

ADD ./ /src

RUN go mod download

RUN make trojan-box

FROM alpine:latest
LABEL maintainer="EthanWang <git@ethan.wang>"

RUN apk add --update --no-cache tzdata ca-certificates bash py-pip groff \
    musl gmp openssl libstdc++ fish shadow vips vips-dev \
    && usermod -s `which fish` root
ENV DEBIAN_FRONTEND=noninteractive
ENV CGO_CFLAGS_ALLOW=-Xpreprocessor

COPY --from=builder /src/build/trojan-box /opt/trojan-box
COPY --from=builder /src/config /opt/config

# copy timeinfo.zip
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip


EXPOSE 5577

WORKDIR /opt
#CMD ["/opt/warehouse"]
