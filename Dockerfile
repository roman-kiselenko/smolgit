FROM golang:1.22-alpine3.20 AS builder

RUN mkdir /app && mkdir -p /usr/local/src/smolgit
WORKDIR /usr/local/src/smolgit

ADD ./go.mod ./go.sum ./
RUN go mod download
ADD . ./

RUN go build -v -o /build/smolgit

FROM alpine/git:2.45.2 AS runner

COPY --from=builder /build/smolgit /usr/bin/smolgit

EXPOSE 3080
EXPOSE 3081
ENTRYPOINT ["/usr/bin/smolgit", "--config", "/etc/smolgit/config.yaml"]