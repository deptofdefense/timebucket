# build stage
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git make gcc g++ ca-certificates && update-ca-certificates

WORKDIR /src

COPY . .

RUN make tidy

RUN make bin/timebucket_linux_amd64

# final stage
FROM gcr.io/distroless/base:latest
COPY --from=builder /src/bin/timebucket_linux_amd64 /bin/timebucket
ENTRYPOINT ["/bin/timebucket"]
CMD ["--help"]
