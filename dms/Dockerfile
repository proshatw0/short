FROM golang:1.17-alpine

COPY . /
WORKDIR /

RUN go build -o dms

EXPOSE 6379

CMD ["./dms"]
