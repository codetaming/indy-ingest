# build stage
FROM golang:1.10 AS builder
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR $GOPATH/src/github.com/codetaming/indyingest
COPY Gopkg.toml Gopkg.lock ./
COPY . ./
RUN dep ensure
RUN cd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o ingestd

# final stage
FROM alpine
EXPOSE 8000
WORKDIR /app
COPY --from=builder /cmd/ingestd /app/
ENTRYPOINT ./ingestd