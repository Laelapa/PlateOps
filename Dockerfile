ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -v -o /plateops ./cmd/api

FROM alpine:latest
ENV SERVER_PORT=8080
RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -s /bin/sh appuser
RUN mkdir -p /docs
COPY --from=builder /plateops /usr/local/bin/
COPY --from=builder /usr/src/app/docs/openapi.json /docs/openapi.json
USER appuser
EXPOSE ${SERVER_PORT}

CMD ["plateops"]

# Licensing information for the included software is available in the LICENSE and NOTICE files in the project repository.