package main

import (
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "github.com/octago/polygon/api"
	"github.com/octago/polygon/server"
)

func main() {
	srv, err := server.New(os.Args[1])
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPolygonServerServer(grpcServer, srv)
	fmt.Println("adam and steve ok")
	grpcServer.Serve(lis)
}
