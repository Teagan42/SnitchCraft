FROM golang:1.23
WORKDIR /app

COPY . .

ENTRYPOINT [ "/app/entrypoint.sh" ]