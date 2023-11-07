// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.0
// source: .proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ShortURLService_CreateShortURL_FullMethodName = "/shorturl.ShortURLService/CreateShortURL"
	ShortURLService_GetOriginalURL_FullMethodName = "/shorturl.ShortURLService/GetOriginalURL"
)

// ShortURLServiceClient is the client API for ShortURLService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortURLServiceClient interface {
	CreateShortURL(ctx context.Context, in *OriginalURL, opts ...grpc.CallOption) (*ShortURL, error)
	GetOriginalURL(ctx context.Context, in *ShortURL, opts ...grpc.CallOption) (*OriginalURL, error)
}

type shortURLServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShortURLServiceClient(cc grpc.ClientConnInterface) ShortURLServiceClient {
	return &shortURLServiceClient{cc}
}

func (c *shortURLServiceClient) CreateShortURL(ctx context.Context, in *OriginalURL, opts ...grpc.CallOption) (*ShortURL, error) {
	out := new(ShortURL)
	err := c.cc.Invoke(ctx, ShortURLService_CreateShortURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortURLServiceClient) GetOriginalURL(ctx context.Context, in *ShortURL, opts ...grpc.CallOption) (*OriginalURL, error) {
	out := new(OriginalURL)
	err := c.cc.Invoke(ctx, ShortURLService_GetOriginalURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortURLServiceServer is the server API for ShortURLService service.
// All implementations must embed UnimplementedShortURLServiceServer
// for forward compatibility
type ShortURLServiceServer interface {
	CreateShortURL(context.Context, *OriginalURL) (*ShortURL, error)
	GetOriginalURL(context.Context, *ShortURL) (*OriginalURL, error)
	mustEmbedUnimplementedShortURLServiceServer()
}

// UnimplementedShortURLServiceServer must be embedded to have forward compatible implementations.
type UnimplementedShortURLServiceServer struct {
}

func (UnimplementedShortURLServiceServer) CreateShortURL(context.Context, *OriginalURL) (*ShortURL, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShortURL not implemented")
}
func (UnimplementedShortURLServiceServer) GetOriginalURL(context.Context, *ShortURL) (*OriginalURL, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOriginalURL not implemented")
}
func (UnimplementedShortURLServiceServer) mustEmbedUnimplementedShortURLServiceServer() {}

// UnsafeShortURLServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortURLServiceServer will
// result in compilation errors.
type UnsafeShortURLServiceServer interface {
	mustEmbedUnimplementedShortURLServiceServer()
}

func RegisterShortURLServiceServer(s grpc.ServiceRegistrar, srv ShortURLServiceServer) {
	s.RegisterService(&ShortURLService_ServiceDesc, srv)
}

func _ShortURLService_CreateShortURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OriginalURL)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServiceServer).CreateShortURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortURLService_CreateShortURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServiceServer).CreateShortURL(ctx, req.(*OriginalURL))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortURLService_GetOriginalURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortURL)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServiceServer).GetOriginalURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortURLService_GetOriginalURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServiceServer).GetOriginalURL(ctx, req.(*ShortURL))
	}
	return interceptor(ctx, in, info, handler)
}

// ShortURLService_ServiceDesc is the grpc.ServiceDesc for ShortURLService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShortURLService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shorturl.ShortURLService",
	HandlerType: (*ShortURLServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateShortURL",
			Handler:    _ShortURLService_CreateShortURL_Handler,
		},
		{
			MethodName: "GetOriginalURL",
			Handler:    _ShortURLService_GetOriginalURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: ".proto",
}
