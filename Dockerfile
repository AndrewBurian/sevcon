FROM golang:1.11 AS build
ADD . /go/src/sevcon
WORKDIR /go/src/sevcon
RUN go build -v

FROM gcr.io/distroless/base
WORKDIR /app
ADD ./site /app/site
COPY --from=build /go/src/sevcon/sevcon /app/sevcon
ENTRYPOINT ["/app/sevcon"]