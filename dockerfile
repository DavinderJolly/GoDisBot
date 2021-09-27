FROM golang:alpine

WORKDIR /GoDisBot

ENV TOKEN=value

COPY . .

RUN : \
    && go mod download \
    && go build -o build/Bot github.com/DavinderJolly/GoDisBot \
    && :

CMD ["./build/Bot"]
