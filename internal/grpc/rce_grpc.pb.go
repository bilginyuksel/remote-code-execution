// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: internal/grpc/rce.proto

package grpc

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

// CodeExecutorServiceClient is the client API for CodeExecutorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CodeExecutorServiceClient interface {
	Exec(ctx context.Context, in *CodeExecutionRequest, opts ...grpc.CallOption) (*CodeExecutionResponse, error)
}

type codeExecutorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCodeExecutorServiceClient(cc grpc.ClientConnInterface) CodeExecutorServiceClient {
	return &codeExecutorServiceClient{cc}
}

func (c *codeExecutorServiceClient) Exec(ctx context.Context, in *CodeExecutionRequest, opts ...grpc.CallOption) (*CodeExecutionResponse, error) {
	out := new(CodeExecutionResponse)
	err := c.cc.Invoke(ctx, "/grpc.codeExecutorService/Exec", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CodeExecutorServiceServer is the server API for CodeExecutorService service.
// All implementations must embed UnimplementedCodeExecutorServiceServer
// for forward compatibility
type CodeExecutorServiceServer interface {
	Exec(context.Context, *CodeExecutionRequest) (*CodeExecutionResponse, error)
	mustEmbedUnimplementedCodeExecutorServiceServer()
}

// UnimplementedCodeExecutorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCodeExecutorServiceServer struct {
}

func (UnimplementedCodeExecutorServiceServer) Exec(context.Context, *CodeExecutionRequest) (*CodeExecutionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exec not implemented")
}
func (UnimplementedCodeExecutorServiceServer) mustEmbedUnimplementedCodeExecutorServiceServer() {}

// UnsafeCodeExecutorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CodeExecutorServiceServer will
// result in compilation errors.
type UnsafeCodeExecutorServiceServer interface {
	mustEmbedUnimplementedCodeExecutorServiceServer()
}

func RegisterCodeExecutorServiceServer(s grpc.ServiceRegistrar, srv CodeExecutorServiceServer) {
	s.RegisterService(&CodeExecutorService_ServiceDesc, srv)
}

func _CodeExecutorService_Exec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CodeExecutionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CodeExecutorServiceServer).Exec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.codeExecutorService/Exec",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CodeExecutorServiceServer).Exec(ctx, req.(*CodeExecutionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CodeExecutorService_ServiceDesc is the grpc.ServiceDesc for CodeExecutorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CodeExecutorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.codeExecutorService",
	HandlerType: (*CodeExecutorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Exec",
			Handler:    _CodeExecutorService_Exec_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/grpc/rce.proto",
}
