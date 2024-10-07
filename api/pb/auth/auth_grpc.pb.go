// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.18.0
// source: auth.proto

package auth

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

// AuthenticateClient is the client API for Authenticate service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthenticateClient interface {
	GenerateJWT(ctx context.Context, in *GenerateJWTReq, opts ...grpc.CallOption) (*GenerateJWTResp, error)
	ValidateJWT(ctx context.Context, in *ValidateJWTReq, opts ...grpc.CallOption) (*ValidateJWTResp, error)
}

type authenticateClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthenticateClient(cc grpc.ClientConnInterface) AuthenticateClient {
	return &authenticateClient{cc}
}

func (c *authenticateClient) GenerateJWT(ctx context.Context, in *GenerateJWTReq, opts ...grpc.CallOption) (*GenerateJWTResp, error) {
	out := new(GenerateJWTResp)
	err := c.cc.Invoke(ctx, "/Authenticate/GenerateJWT", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authenticateClient) ValidateJWT(ctx context.Context, in *ValidateJWTReq, opts ...grpc.CallOption) (*ValidateJWTResp, error) {
	out := new(ValidateJWTResp)
	err := c.cc.Invoke(ctx, "/Authenticate/ValidateJWT", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthenticateServer is the server API for Authenticate service.
// All implementations must embed UnimplementedAuthenticateServer
// for forward compatibility
type AuthenticateServer interface {
	GenerateJWT(context.Context, *GenerateJWTReq) (*GenerateJWTResp, error)
	ValidateJWT(context.Context, *ValidateJWTReq) (*ValidateJWTResp, error)
	mustEmbedUnimplementedAuthenticateServer()
}

// UnimplementedAuthenticateServer must be embedded to have forward compatible implementations.
type UnimplementedAuthenticateServer struct {
}

func (UnimplementedAuthenticateServer) GenerateJWT(context.Context, *GenerateJWTReq) (*GenerateJWTResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateJWT not implemented")
}
func (UnimplementedAuthenticateServer) ValidateJWT(context.Context, *ValidateJWTReq) (*ValidateJWTResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateJWT not implemented")
}
func (UnimplementedAuthenticateServer) mustEmbedUnimplementedAuthenticateServer() {}

// UnsafeAuthenticateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthenticateServer will
// result in compilation errors.
type UnsafeAuthenticateServer interface {
	mustEmbedUnimplementedAuthenticateServer()
}

func RegisterAuthenticateServer(s grpc.ServiceRegistrar, srv AuthenticateServer) {
	s.RegisterService(&Authenticate_ServiceDesc, srv)
}

func _Authenticate_GenerateJWT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateJWTReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticateServer).GenerateJWT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Authenticate/GenerateJWT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticateServer).GenerateJWT(ctx, req.(*GenerateJWTReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Authenticate_ValidateJWT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateJWTReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticateServer).ValidateJWT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Authenticate/ValidateJWT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticateServer).ValidateJWT(ctx, req.(*ValidateJWTReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Authenticate_ServiceDesc is the grpc.ServiceDesc for Authenticate service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Authenticate_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Authenticate",
	HandlerType: (*AuthenticateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenerateJWT",
			Handler:    _Authenticate_GenerateJWT_Handler,
		},
		{
			MethodName: "ValidateJWT",
			Handler:    _Authenticate_ValidateJWT_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}