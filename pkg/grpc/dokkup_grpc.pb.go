// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: grpc/dokkup.proto

package dokkup

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DokkupClient is the client API for Dokkup service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DokkupClient interface {
	DeployContainer(ctx context.Context, in *DeployContainerRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type dokkupClient struct {
	cc grpc.ClientConnInterface
}

func NewDokkupClient(cc grpc.ClientConnInterface) DokkupClient {
	return &dokkupClient{cc}
}

func (c *dokkupClient) DeployContainer(ctx context.Context, in *DeployContainerRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/Dokkup/DeployContainer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DokkupServer is the server API for Dokkup service.
// All implementations must embed UnimplementedDokkupServer
// for forward compatibility
type DokkupServer interface {
	DeployContainer(context.Context, *DeployContainerRequest) (*empty.Empty, error)
	mustEmbedUnimplementedDokkupServer()
}

// UnimplementedDokkupServer must be embedded to have forward compatible implementations.
type UnimplementedDokkupServer struct {
}

func (UnimplementedDokkupServer) DeployContainer(context.Context, *DeployContainerRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeployContainer not implemented")
}
func (UnimplementedDokkupServer) mustEmbedUnimplementedDokkupServer() {}

// UnsafeDokkupServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DokkupServer will
// result in compilation errors.
type UnsafeDokkupServer interface {
	mustEmbedUnimplementedDokkupServer()
}

func RegisterDokkupServer(s grpc.ServiceRegistrar, srv DokkupServer) {
	s.RegisterService(&Dokkup_ServiceDesc, srv)
}

func _Dokkup_DeployContainer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployContainerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DokkupServer).DeployContainer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Dokkup/DeployContainer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DokkupServer).DeployContainer(ctx, req.(*DeployContainerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Dokkup_ServiceDesc is the grpc.ServiceDesc for Dokkup service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Dokkup_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Dokkup",
	HandlerType: (*DokkupServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeployContainer",
			Handler:    _Dokkup_DeployContainer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/dokkup.proto",
}