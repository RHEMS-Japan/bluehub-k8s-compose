package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "../proto"
)

const (
    port = ":10001"
)

// server is used to implement HelloWorldServer.
type server struct{}

// SayHello implements HelloWorldServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    log.Printf("Received: %v", in.Name)
    return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

// GetUser
func (s *server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.User, error) {
    log.Printf("Received: %v", in.Id)
    return &pb.User{
        Id: in.Id,
        Name: "SampleUser"}, nil
}

// CreateUser
func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.User, error) {
    log.Printf("Received: %v", in.Name)
    return &pb.User{
        Id: "123",
        Name: in.Name}, nil
}

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterHelloWorldServiceServer(s, &server{})
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
