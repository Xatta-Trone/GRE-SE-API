# # base go image
FROM golang:1.20-alpine as builder

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build -o wcapp

RUN chmod +x /app/wcapp


COPY --from=builder /app/wcapp /app

CMD [ "/app/wcapp" ]