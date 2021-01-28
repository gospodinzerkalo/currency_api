FROM golang:latest

WORKDIR /app

RUN go mod init github.com/gospodinzerkalo/currencyapi

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["make","run"]