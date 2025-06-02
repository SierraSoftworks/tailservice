FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.24.3 as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app/
ADD . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o tailservice main.go

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.22.0

RUN mkdir -p /app/data/ && adduser -D -u 1000 tailservice -h /app/data
VOLUME /app/data
USER nonroot

WORKDIR /app/
COPY --from=builder /app/tailservice /app/tailservice
ENTRYPOINT ["/app/tailservice"]
