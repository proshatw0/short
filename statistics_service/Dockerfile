FROM golang:latest

COPY . /
WORKDIR /

RUN go build -o stats-service

EXPOSE 1333

CMD ["./stats-service"]