FROM golang:1.23-bookworm as base

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  small-user

WORKDIR $GOPATH/src/gitpebble/

COPY . .

RUN go mod download
RUN go mod verify

RUN GOOS=linux GOARCH=amd64 go build -o .

FROM gcr.io/distroless/static-debian11

COPY --from=base $GOPATH/src/gitpebble .

USER small-user:small-user

CMD ["./gitpebble"]
