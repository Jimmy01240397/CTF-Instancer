FROM golang:1.19-alpine as builder

RUN apk add --no-cache make build-base

COPY . /src
WORKDIR /src
RUN make clean && make


FROM docker:dind as release

COPY --from=builder /src/bin /app

WORKDIR /app

RUN touch .env

COPY ./docker-entrypoint.sh ./docker-entrypoint.sh

RUN chmod +x ./docker-entrypoint.sh

ENTRYPOINT ["./docker-entrypoint.sh"]
