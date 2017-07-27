package main

import (
	"context"
	"fmt"
	"io"
	"os"

	pb "github.com/octago/polygon/api"
	"google.golang.org/grpc"
)

func startStdout(stream pb.PolygonServer_AttachClient) {
	fmt.Println("start stdout")
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(chunk.Chunk)
	}
}

func startStdin(stream pb.PolygonServer_AttachClient) {
	fmt.Println("start stdin")
	for {
		b := make([]byte, 4096)
		read, err := os.Stdin.Read(b)
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		err = stream.Send(&pb.StreamChunk{Chunk: b[:read]})
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := pb.NewPolygonServerClient(conn)
	ctx := context.Background()

	createReq := &pb.CreateRequest{
		TemplateId: "ulimit",
	}
	resp, err := client.Create(ctx, createReq)
	if err != nil {
		panic(err)
	}
	id := resp.Stand.Id

	fmt.Println("container ID", id)

	stream, err := client.Attach(ctx)
	if err != nil {
		panic(err)
	}
	firstChunk := &pb.StreamChunk{
		StandId: id,
	}
	if err := stream.Send(firstChunk); err != nil {
		panic(err)
	}
	go startStdin(stream)
	go startStdout(stream)
	select {}
}
