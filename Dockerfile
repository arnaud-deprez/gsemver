FROM golang:alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o gsemver

FROM alpine:3.19
RUN apk --no-cache add ca-certificates git
COPY --from=builder /src/gsemver /usr/local/bin/gsemver
ENTRYPOINT ["gsemver"]
