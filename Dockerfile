FROM golang:1.22-alpine

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o proxy

EXPOSE 8080
CMD ["./proxy"]