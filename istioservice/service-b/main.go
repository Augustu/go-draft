package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"

	pb "github.com/Augustu/go-draft/istioservice/proto"
	"github.com/Augustu/go-draft/istioservice/shutdownsignal"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(c context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Got request: %s", req.GetName())

	return &pb.HelloReply{
		Message: "world b",
	}, nil
}

func GrpcServe(addr string, stopCh <-chan struct{}) (<-chan struct{}, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("serve: %s", err)
		return nil, err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &server{})

	stoppedCh := make(chan struct{})

	go func() {
		defer close(stoppedCh)
		<-stopCh
		grpcServer.GracefulStop()
		fmt.Printf("Stopped grpc listening on %s", addr)
	}()

	go func() {
		fmt.Printf("Grpc serve at: %s\n", addr)
		err := grpcServer.Serve(listener)

		msg := fmt.Sprintf("Stopped grpc listening on %s", addr)
		select {
		case <-stopCh:
			fmt.Println(msg)
		default:
			panic(fmt.Sprintf("%s due to error: %s", msg, err))
		}
	}()

	return stoppedCh, nil
}

func HTTPServe(addr, grpcAddr string, shutdownTimeout time.Duration, stopCh <-chan struct{}) (<-chan struct{}, error) {
	gwmux := runtime.NewServeMux()

	dopts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterGreeterHandlerFromEndpoint(context.Background(), gwmux, grpcAddr, dopts)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return nil, err
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:    addr,
		Handler: gwmux,
	}

	stoppedCh := make(chan struct{})
	go func() {
		defer close(stoppedCh)
		<-stopCh
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		server.Shutdown(ctx)
		cancel()
	}()

	go func() {
		fmt.Printf("Http serve at: %s\n", addr)
		err := server.Serve(listener)

		msg := fmt.Sprintf("Stopped http listening on %s", addr)
		select {
		case <-stopCh:
			fmt.Println(msg)
		default:
			panic(fmt.Sprintf("%s due to error: %s", msg, err))
		}
	}()

	return stoppedCh, nil
}

func main() {

	grpcAddr := ":50051"

	stopCh := shutdownsignal.SetupSignalHandler()

	grpcStopCh, _ := GrpcServe(grpcAddr, stopCh)

	<-stopCh

	<-grpcStopCh
}
