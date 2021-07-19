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

const (
	ServiceBAddr string = "service-b:50051"
	// ServiceBAddr string = "127.0.0.1:50052"
)

type server struct {
	pb.UnimplementedGreeterServer

	clientB pb.GreeterClient
}

func newServer(addrB string) *server {
	conn, err := grpc.Dial(addrB, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &server{
		clientB: pb.NewGreeterClient(conn),
	}
}

func (s *server) SayHello(c context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	reply := s.SayHelloB(ServiceBAddr)
	// reply := "fack"
	log.Printf("Greeting: %s from %s", reply, ServiceBAddr)

	return &pb.HelloReply{
		Message: reply,
	}, nil
}

func (s *server) SayHelloB(addr string) string {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r, err := s.clientB.SayHello(ctx, &pb.HelloRequest{Name: "service-a"})
	if err != nil {
		log.Printf("could not greet: %v", err)
	}

	return r.GetMessage()
}

func GrpcServe(addr, grpcAddr string, stopCh <-chan struct{}) (<-chan struct{}, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("serve: %s", err)
		return nil, err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, newServer(grpcAddr))

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
		fmt.Printf("Stopped http listening on %s", addr)
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
	httpAddr := ":8000"

	timeout := 3 * time.Second

	stopCh := shutdownsignal.SetupSignalHandler()

	grpcStoppedCh, _ := GrpcServe(grpcAddr, ServiceBAddr, stopCh)
	httpStoppedCh, _ := HTTPServe(httpAddr, grpcAddr, timeout, stopCh)

	<-stopCh

	<-grpcStoppedCh
	<-httpStoppedCh
}
