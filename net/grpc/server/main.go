package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc/networkstuff"

	"google.golang.org/grpc"
)

type routerServiceServer struct {
	pb.UnimplementedRouterServiceServer
	localRouters []*pb.Router
}

func (s *routerServiceServer) GetRouter(ctx context.Context, req *pb.RouterRequest) (*pb.Router, error) {
	for _, r := range s.localRouters {
		if r.Id == req.Id {
			return r, nil
		}
	}
	return &pb.Router{}, nil
}

var server = &routerServiceServer{
	localRouters: []*pb.Router{
		&pb.Router{
			Id:       1,
			Hostname: "Router A",
			Interfaces: []*pb.Interface{
				&pb.Interface{
					Id:          1000,
					Description: "Gi 0/0/0",
				},
				&pb.Interface{
					Id:          1001,
					Description: "Gi 0/0/1",
				},
			},
		},
		&pb.Router{
			Id:       2,
			Hostname: "Router B",
			Interfaces: []*pb.Interface{
				&pb.Interface{
					Id:          2000,
					Description: "Gi 0/0/0",
				},
				&pb.Interface{
					Id:          2001,
					Description: "Gi 0/0/1",
				},
			},
		},
	},
}

func main() {
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterRouterServiceServer(grpcServer, server)
	grpcServer.Serve(ln)
}
