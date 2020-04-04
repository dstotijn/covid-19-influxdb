FROM golang:1.14-alpine

WORKDIR /go/src/covid-19-influxdb
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["covid-19-influxdb"]