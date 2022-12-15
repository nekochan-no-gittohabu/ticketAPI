FROM golang:1.19-alpine

RUN mkdir /app
WORKDIR /app

COPY . .
COPY .env .

RUN apk add git

# Download all the dependencies
RUN go get -d -v ./...
RUN go install -v ./...

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

RUN go build -o /build

EXPOSE 8080

CMD [ "/build" ]