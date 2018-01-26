FROM golang:1.9.2-alpine

ADD . /go/src/github.com/matyga/golang_rest
RUN go install github.com/matyga/golang_rest

WORKDIR /go/src/github.com/matyga/golang_rest

EXPOSE 8080

ENTRYPOINT ["go","run","rest1.go"] 

