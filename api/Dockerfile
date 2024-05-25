FROM golang:alpine as builder

WORKDIR /app


COPY go.mod ./

RUN go mod download


COPY . ./


RUN go build -o /godo/main .

FROM alpine

WORKDIR /app



COPY --from=builder /godo /app


COPY .env /app

EXPOSE 3000

CMD ["/app/main"]


