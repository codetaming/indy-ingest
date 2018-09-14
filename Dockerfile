FROM golang:1.11-stretch AS BUILD
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o ingestd .

FROM alpine:3.8
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN mkdir /app
COPY --from=BUILD /app/ingestd /app
RUN ls -l

EXPOSE 9000
ENTRYPOINT /app/ingestd