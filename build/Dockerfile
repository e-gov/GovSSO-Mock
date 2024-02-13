# syntax=docker/dockerfile:1
FROM golang:1.19-bookworm as builder

WORKDIR /govsso-mock
COPY go.mod go.sum ./

RUN go mod download

COPY ./app ./app
RUN go build -o /tmp/ ./app

# Use Debian 12 image for runtime to match builder image's Debian 12, in case dynamic linking to glibc is used.
# We are not using `nonroot` version of this image, because it would make some generated
# keys inaccessible from the container on some operating systems.
FROM gcr.io/distroless/base-nossl-debian12

WORKDIR /govsso-mock
COPY ./ui ./ui
COPY --from=builder /tmp/app /govsso-mock/govsso-mock

EXPOSE 10443
ENTRYPOINT ["/govsso-mock/govsso-mock"]
