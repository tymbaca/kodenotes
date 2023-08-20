FROM golang:1.21 AS build
WORKDIR /go/src
COPY ./go.mod .
RUN go mod download
COPY . .
# RUN mkdir /go/bin
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app

FROM alpine:latest AS run
WORKDIR /go
COPY --from=build /go/bin/app .
RUN mkdir logs
CMD ./app
