// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package token_sync

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

// TokenSyncClient is the client API for TokenSync service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokenSyncClient interface {
	Replicate(ctx context.Context, in *TokenBroadcast, opts ...grpc.CallOption) (*Response, error)
	GetTimeStamp(ctx context.Context, in *TokenId, opts ...grpc.CallOption) (*TimeStamp, error)
	GetTokenDetails(ctx context.Context, in *TokenId, opts ...grpc.CallOption) (*TokenBroadcast, error)
}

type tokenSyncClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenSyncClient(cc grpc.ClientConnInterface) TokenSyncClient {
	return &tokenSyncClient{cc}
}

func (c *tokenSyncClient) Replicate(ctx context.Context, in *TokenBroadcast, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/token_sync.TokenSync/Replicate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenSyncClient) GetTimeStamp(ctx context.Context, in *TokenId, opts ...grpc.CallOption) (*TimeStamp, error) {
	out := new(TimeStamp)
	err := c.cc.Invoke(ctx, "/token_sync.TokenSync/GetTimeStamp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenSyncClient) GetTokenDetails(ctx context.Context, in *TokenId, opts ...grpc.CallOption) (*TokenBroadcast, error) {
	out := new(TokenBroadcast)
	err := c.cc.Invoke(ctx, "/token_sync.TokenSync/GetTokenDetails", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenSyncServer is the server API for TokenSync service.
// All implementations must embed UnimplementedTokenSyncServer
// for forward compatibility
type TokenSyncServer interface {
	Replicate(context.Context, *TokenBroadcast) (*Response, error)
	GetTimeStamp(context.Context, *TokenId) (*TimeStamp, error)
	GetTokenDetails(context.Context, *TokenId) (*TokenBroadcast, error)
	mustEmbedUnimplementedTokenSyncServer()
}

// UnimplementedTokenSyncServer must be embedded to have forward compatible implementations.
type UnimplementedTokenSyncServer struct {
}

func (UnimplementedTokenSyncServer) Replicate(context.Context, *TokenBroadcast) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Replicate not implemented")
}
func (UnimplementedTokenSyncServer) GetTimeStamp(context.Context, *TokenId) (*TimeStamp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTimeStamp not implemented")
}
func (UnimplementedTokenSyncServer) GetTokenDetails(context.Context, *TokenId) (*TokenBroadcast, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTokenDetails not implemented")
}
func (UnimplementedTokenSyncServer) mustEmbedUnimplementedTokenSyncServer() {}

// UnsafeTokenSyncServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenSyncServer will
// result in compilation errors.
type UnsafeTokenSyncServer interface {
	mustEmbedUnimplementedTokenSyncServer()
}

func RegisterTokenSyncServer(s grpc.ServiceRegistrar, srv TokenSyncServer) {
	s.RegisterService(&TokenSync_ServiceDesc, srv)
}

func _TokenSync_Replicate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenBroadcast)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenSyncServer).Replicate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token_sync.TokenSync/Replicate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenSyncServer).Replicate(ctx, req.(*TokenBroadcast))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenSync_GetTimeStamp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenSyncServer).GetTimeStamp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token_sync.TokenSync/GetTimeStamp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenSyncServer).GetTimeStamp(ctx, req.(*TokenId))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenSync_GetTokenDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenSyncServer).GetTokenDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token_sync.TokenSync/GetTokenDetails",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenSyncServer).GetTokenDetails(ctx, req.(*TokenId))
	}
	return interceptor(ctx, in, info, handler)
}

// TokenSync_ServiceDesc is the grpc.ServiceDesc for TokenSync service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TokenSync_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "token_sync.TokenSync",
	HandlerType: (*TokenSyncServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Replicate",
			Handler:    _TokenSync_Replicate_Handler,
		},
		{
			MethodName: "GetTimeStamp",
			Handler:    _TokenSync_GetTimeStamp_Handler,
		},
		{
			MethodName: "GetTokenDetails",
			Handler:    _TokenSync_GetTokenDetails_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "token_sync/token_sync.proto",
}
