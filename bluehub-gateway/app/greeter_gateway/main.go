package main

import (
  "flag"
  "fmt"
  "net/http"

  "github.com/golang/glog"
  "golang.org/x/net/context"
  "github.com/grpc-ecosystem/grpc-gateway/runtime"
  "google.golang.org/grpc"

  gw "../proto"
)

func run() error {
  ctx := context.Background()
  ctx, cancel := context.WithCancel(ctx)
  defer cancel()

  mux := runtime.NewServeMux()
  opts := []grpc.DialOption{grpc.WithInsecure()}
  endpoint := fmt.Sprintf("bluehub-router-service:10001")
  err := gw.RegisterHelloWorldServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
  if err != nil {
    return err
  }

  return http.ListenAndServe(":10000", mux)
}

func main() {
  flag.Parse()
  defer glog.Flush()

  if err := run(); err != nil {
    glog.Fatal(err)
  }
}
