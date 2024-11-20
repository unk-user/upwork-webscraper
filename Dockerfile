FROM golang:1.23 AS build
WORKDIR /helloworld

COPY go.mod go.sum ./

COPY . ./
RUN go build -tags lambda.norpc -o main main.go

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /helloworld/main ./main

ENTRYPOINT [ "./main" ]