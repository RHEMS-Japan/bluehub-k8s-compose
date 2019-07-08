FROM golang:1.12

RUN mkdir /go/src/app
COPY ./ /go/src/app

RUN go get -u -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go get -u -v github.com/golang/protobuf/protoc-gen-go
RUN go get -u -v google.golang.org/grpc

CMD [ "go", "run", "/go/src/app/greeter_gateway/main.go" ]
