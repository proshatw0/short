FROM golang:latest

COPY . /
WORKDIR /

RUN go build -o dms

EXPOSE 6379

CMD ["./dms"]