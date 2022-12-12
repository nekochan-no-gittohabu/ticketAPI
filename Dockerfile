FROM golang:1.19-alpine

RUN mkdir /app
WORKDIR /app

COPY . .
COPY .env .

# Download all the dependencies
RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o /build

EXPOSE 8080

CMD [ "/build" ]