FROM golang:1.8

COPY . /go/src/github.com/drmjo/heimdall

RUN go get -v ./...
RUN go install -v ./...

CMD heimdall --serve --listen 0.0.0.0:80
