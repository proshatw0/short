FROM golang:1.17-alpine

COPY . /
WORKDIR /

RUN go build -o stats-service

EXPOSE 1333

CMD ["./stats-service"]
