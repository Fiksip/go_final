FROM golang:1.25 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM ubuntu:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

VOLUME /data

CMD ["./main"]