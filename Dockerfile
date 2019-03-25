FROM golang:1.12 AS build
ADD . /sevcon
WORKDIR /sevcon
RUN go build -v

FROM gcr.io/distroless/base
WORKDIR /app
ADD ./site /app/site
COPY --from=build /sevcon/sevcon /app/sevcon
ENTRYPOINT ["/app/sevcon"]
