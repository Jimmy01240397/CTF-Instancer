FROM golang:1.23-alpine as builder

RUN apk add --no-cache make build-base

WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /src
RUN make clean && make

FROM docker:dind as release

RUN wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/bin/yq && chmod +x /usr/bin/yq

COPY --from=builder /src/bin /app

WORKDIR /app

RUN touch .env && mkdir images

COPY ./docker-entrypoint.sh ./docker-entrypoint.sh

RUN chmod +x ./docker-entrypoint.sh

ENTRYPOINT ["./docker-entrypoint.sh"]
