FROM golang:1.17.0-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY config.yaml .

RUN go mod download

COPY . .

RUN go build -o ./out/app ./client
RUN go build -o ./out/app ./server

EXPOSE 8080

CMD [ "./out/app" ]

