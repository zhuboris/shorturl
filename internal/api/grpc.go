package api

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"shorturl/internal/pb"
)

// GRPCServer is gRPC server implementation that processing requests to short URL service.
// It implements pb.ShortURLServiceServer interface from generated protobuf.
// Object must be initialized with NewGRPCServer
type GRPCServer struct {
	pb.UnimplementedShortURLServiceServer

	server     *grpc.Server
	listener   net.Listener
	urlService ShortURLService
}

// NewGRPCServer initializes GRPCServer with its address to listen, and short URL service.
// It returns a pointer to object.
func NewGRPCServer(listenAddress string, urlService ShortURLService) (*GRPCServer, error) {
	server := initGRPCServer(urlService)
	err := server.initListener(listenAddress)
	if err != nil {
		return nil, err
	}

	return server, nil
}

// Run is calling method Serve of object's grpc.Server with object's listener. It will return its error.
func (s *GRPCServer) Run() error {
	return s.server.Serve(s.listener)
}

// Stop is calling method GracefulStop of object's grpc.Server.
func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

// CreateShortURL is an implementation of rpc CreateShortURL method. It is
// processing a request and handles its error. If error is not nil, it responds
// with corresponded error codes.Code and writes an error message.
func (s *GRPCServer) CreateShortURL(ctx context.Context, req *pb.OriginalURL) (*pb.ShortURL, error) {
	shortURL, err := handleCreationShortURL(ctx, req.Url, s.urlService)
	if err != nil {
		_, code := errorStatusCodes(err)
		return nil, status.Error(code, err.Error())
	}

	resp := &pb.ShortURL{Url: shortURL}
	return resp, nil
}

// GetOriginalURL is an implementation of rpc GetOriginalURL method. It is
// processing a request and handles its error. If error is not nil, it responds
// with corresponded error codes.Code and writes an error message.
func (s *GRPCServer) GetOriginalURL(ctx context.Context, req *pb.ShortURL) (*pb.OriginalURL, error) {
	original, err := handleGetOriginalURL(ctx, req.Url, s.urlService)
	if err != nil {
		_, code := errorStatusCodes(err)
		return nil, status.Error(code, err.Error())
	}

	resp := &pb.OriginalURL{Url: original}
	return resp, nil
}

// initGRPCServer initializes grpc.Server and registers it to serve requests with
// GRPCServer object as pb.ShortURLServiceServer. Also, it registers reflection.
//
// This function initializes and returns a pointer to GRPCServer
// that is ready to start serving requests.
func initGRPCServer(urlService ShortURLService) *GRPCServer {
	server := grpc.NewServer(grpc.UnaryInterceptor(loggingUnaryInterceptor))

	serviceServer := &GRPCServer{
		server:     server,
		urlService: urlService,
	}

	pb.RegisterShortURLServiceServer(server, serviceServer)
	reflection.Register(server)
	return serviceServer
}

func (s *GRPCServer) initListener(listenAddress string) error {
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.listener = listener
	return nil
}
