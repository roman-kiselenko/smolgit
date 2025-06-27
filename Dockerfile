FROM --platform=$BUILDPLATFORM golang:1.23-alpine3.20 AS builder
ARG TARGETOS TARGETARCH

RUN mkdir /app && mkdir -p /usr/local/src/smolgit
WORKDIR /usr/local/src/smolgit

ADD ./go.mod ./go.sum ./
RUN go mod download
ADD . ./

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -v -o /build/smolgit

FROM --platform=$BUILDPLATFORM alpine/git:v2.49.0 AS runner

COPY --from=builder /build/smolgit /usr/bin/smolgit

EXPOSE 3080
EXPOSE 3081
ENTRYPOINT ["/usr/bin/smolgit", "--config", "/etc/smolgit/config.yaml"]
