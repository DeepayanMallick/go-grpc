package main

import (
	"context"
	"log"
	"net"

	pb "github.com/DeepayanMallick/go-grpc/gunk/v1/profile"
	"google.golang.org/grpc"
)

const (
	// Port for gRPC server to listen to
	PORT = ":50051"
)

type ProfileServer struct {
	pb.UnimplementedProfileServiceServer
}

func (s *ProfileServer) Profile(ctx context.Context, req *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	log.Printf("Req: %+v", req)
	response := &pb.ProfileResponse{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@graphlogic.com",
	}

	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)

	if err != nil {
		log.Fatalf("failed connection: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterProfileServiceServer(s, &ProfileServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
