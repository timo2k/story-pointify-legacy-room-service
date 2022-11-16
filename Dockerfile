FROM golang:1 as builder

WORKDIR /go/src/app
COPY *.go .

RUN go mod download
RUN CGO_ENABLED=0 go build .

FROM alpine
COPY --from=builder /go/src/app/app .
CMD ["./app"]