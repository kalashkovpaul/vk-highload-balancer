FROM golang:1.19 AS build

ADD . /app
WORKDIR /app
RUN go build  ./main.go

FROM ubuntu:20.04

WORKDIR /usr/src/app
COPY . .
COPY --from=build /app/main/ .

EXPOSE 80
CMD ./main