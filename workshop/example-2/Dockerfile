FROM 289208114389.dkr.ecr.us-east-1.amazonaws.com/golang:1.14.4-stretch as builder

ENV LIB_RDKAFKA_VERSION='v1.5.2'

RUN git clone --depth 1 --branch ${LIB_RDKAFKA_VERSION} https://github.com/edenhill/librdkafka.git \
    && cd librdkafka \
    && ./configure \
    && make \
    && make install \
    && ldconfig \
    && rm -rf librdkafka

WORKDIR /app

COPY . .

RUN go build -a -ldflags="-s -w" -o kafka-producer

FROM 289208114389.dkr.ecr.us-east-1.amazonaws.com/debian:sid-slim

ENV TZ=America/Sao_Paulo
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update \
     && apt-get install -y curl \
     && apt-get install -y --no-install-recommends ca-certificates \
     && update-ca-certificates

WORKDIR /app

COPY --from=builder ./app/kafka-producer .
