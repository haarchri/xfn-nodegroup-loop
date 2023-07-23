FROM golang:1.20-alpine3.17 as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /function .

###
FROM alpine:3.17.3

LABEL org.opencontainers.image.source="https://github.com/haarchri/xfn-nodegroup-loop"

COPY crossplane.yaml /crossplane.yaml
COPY --from=builder /function /function

ENTRYPOINT ["/function"]
