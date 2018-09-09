FROM alpine:3.7
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY /bin/ingest /app/ingestd
WORKDIR /app

EXPOSE 9000
ENTRYPOINT ./ingestd