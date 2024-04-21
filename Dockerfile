FROM golang:1.19 as builder

RUN apt install make

COPY . /src
WORKDIR /src
RUN make


FROM docker:dind as release

COPY --from=builder /src/bin /app

WORKDIR /app

RUN touch .env

ENTRYPOINT ["./ctfinstancer"]
