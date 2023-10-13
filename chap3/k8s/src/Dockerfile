FROM golang:1.20-alpine3.18 AS builder

WORKDIR /go/src/app

COPY . .

RUN go build -o bin/main main.go

FROM golang:1.20-alpine3.18

WORKDIR /go/src/app

COPY --from=builder /go/src/app/bin/ .

ENTRYPOINT [ "./main" ]
