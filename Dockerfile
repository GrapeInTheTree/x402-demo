FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM builder AS facilitator-build
RUN CGO_ENABLED=0 go build -o /facilitator ./cmd/facilitator

FROM builder AS resource-build
RUN CGO_ENABLED=0 go build -o /resource ./cmd/resource

FROM builder AS client-build
RUN CGO_ENABLED=0 go build -o /client ./cmd/client

FROM alpine:3.21 AS facilitator
RUN apk add --no-cache ca-certificates
COPY --from=facilitator-build /facilitator /facilitator
ENTRYPOINT ["/facilitator"]

FROM alpine:3.21 AS resource
RUN apk add --no-cache ca-certificates
COPY --from=resource-build /resource /resource
ENTRYPOINT ["/resource"]

FROM alpine:3.21 AS client
RUN apk add --no-cache ca-certificates
COPY --from=client-build /client /client
ENTRYPOINT ["/client"]
