FROM golang:1.21 AS builder

WORKDIR $GOPATH/src/rinha

COPY . ./

RUN go get -u

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rinha .


FROM alpine:latest

COPY --from=builder /go/src/rinha/rinha ./
COPY migrations ./migrations


ENTRYPOINT ["./rinha"]