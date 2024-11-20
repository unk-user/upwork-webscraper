FROM golang:1.23 AS build

WORKDIR /usr/app/

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ./main

CMD [ "./main" ]