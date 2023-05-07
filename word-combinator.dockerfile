# # base go image
FROM golang:1.20-alpine as builder
RUN apk --no-cache add gcc g++ make git
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN mkdir /app
# COPY . /app
WORKDIR /app

# Copying all the files
COPY . .

#build new app 
RUN CGO_ENABLED=0 go build -o app
#run script
# RUN "./script.sh"

# COPY --from=builder /app/prod/app /app/prod/app

CMD [ "./app"]