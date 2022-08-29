FROM golang:1.19

ENV GO111MODULE=on

WORKDIR /govsso-mock
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build ./...
RUN go install ./...

EXPOSE 10443
CMD $GOPATH/bin/app