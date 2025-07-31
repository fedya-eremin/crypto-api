FROM golang:1.24-alpine3.22 AS builder

WORKDIR /app
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GIT_TERMINAL_PROMPT=0

COPY . .

RUN go build -mod=vendor -ldflags "-s -w" -o app cmd/main/main.go

FROM alpine:3.22
COPY --from=builder /app/app /app
COPY --from=builder /app/database/migrations /migrations
COPY --from=builder /app/spec/openapi.yaml /openapi.yaml

ENTRYPOINT ["/app"]
