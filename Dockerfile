# syntax=docker/dockerfile:1

FROM golang:1.21.3-alpine3.18 AS BUILDER

RUN go version

COPY . /github.com/movie_lib/
WORKDIR /github.com/movie_lib/

RUN go mod download
RUN GOOS=linux go build -o ./bin/server ./cmd/api/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=BUILDER /github.com/movie_lib/bin/server .
COPY --from=BUILDER /github.com/movie_lib/configs/ configs/

EXPOSE 8080

CMD ["./server", "-ConfigPath", "configs/app/api/deploy.yaml"]
