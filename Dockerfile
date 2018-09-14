FROM golang:1.11-stretch AS BUILD
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ingestd .

FROM alpine:3.8
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN mkdir /app
ENV SERVER_PORT=9000
COPY --from=BUILD /app/ingestd .
CMD ["./ingestd"]
EXPOSE 9000