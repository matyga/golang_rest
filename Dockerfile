FROM golang

ADD . /go/src/github.com/matyga/golang_rest
RUN go install github.com/matyga/golang_rest
RUN go get github.com/gorilla/mux
WORKDIR /go/src/github.com/matyga/golang_rest

COPY . .

EXPOSE 8080

ENTRYPOINT ["go","run","rest1.go"] 

