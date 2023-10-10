// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.2
// source: profile/tls.proto

package profile

import (
	context "context"
	public "github.com/bellis-daemon/bellis/modules/backend/app/mobile/public"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	TLSService_GetUserTLS_FullMethodName = "/bellis.backend.mobile.profile.TLSService/GetUserTLS"
	TLSService_CreateTLS_FullMethodName  = "/bellis.backend.mobile.profile.TLSService/CreateTLS"
	TLSService_UpdateTLS_FullMethodName  = "/bellis.backend.mobile.profile.TLSService/UpdateTLS"
	TLSService_DeleteTLS_FullMethodName  = "/bellis.backend.mobile.profile.TLSService/DeleteTLS"
)

// TLSServiceClient is the client API for TLSService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TLSServiceClient interface {
	GetUserTLS(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TLSList, error)
	CreateTLS(ctx context.Context, in *TLS, opts ...grpc.CallOption) (*TLS, error)
	UpdateTLS(ctx context.Context, in *TLS, opts ...grpc.CallOption) (*TLS, error)
	DeleteTLS(ctx context.Context, in *public.PrimitiveID, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type tLSServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTLSServiceClient(cc grpc.ClientConnInterface) TLSServiceClient {
	return &tLSServiceClient{cc}
}

func (c *tLSServiceClient) GetUserTLS(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TLSList, error) {
	out := new(TLSList)
	err := c.cc.Invoke(ctx, TLSService_GetUserTLS_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tLSServiceClient) CreateTLS(ctx context.Context, in *TLS, opts ...grpc.CallOption) (*TLS, error) {
	out := new(TLS)
	err := c.cc.Invoke(ctx, TLSService_CreateTLS_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tLSServiceClient) UpdateTLS(ctx context.Context, in *TLS, opts ...grpc.CallOption) (*TLS, error) {
	out := new(TLS)
	err := c.cc.Invoke(ctx, TLSService_UpdateTLS_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tLSServiceClient) DeleteTLS(ctx context.Context, in *public.PrimitiveID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, TLSService_DeleteTLS_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TLSServiceServer is the server API for TLSService service.
// All implementations should embed UnimplementedTLSServiceServer
// for forward compatibility
type TLSServiceServer interface {
	GetUserTLS(context.Context, *emptypb.Empty) (*TLSList, error)
	CreateTLS(context.Context, *TLS) (*TLS, error)
	UpdateTLS(context.Context, *TLS) (*TLS, error)
	DeleteTLS(context.Context, *public.PrimitiveID) (*emptypb.Empty, error)
}

// UnimplementedTLSServiceServer should be embedded to have forward compatible implementations.
type UnimplementedTLSServiceServer struct {
}

func (UnimplementedTLSServiceServer) GetUserTLS(context.Context, *emptypb.Empty) (*TLSList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserTLS not implemented")
}
func (UnimplementedTLSServiceServer) CreateTLS(context.Context, *TLS) (*TLS, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTLS not implemented")
}
func (UnimplementedTLSServiceServer) UpdateTLS(context.Context, *TLS) (*TLS, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTLS not implemented")
}
func (UnimplementedTLSServiceServer) DeleteTLS(context.Context, *public.PrimitiveID) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTLS not implemented")
}

// UnsafeTLSServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TLSServiceServer will
// result in compilation errors.
type UnsafeTLSServiceServer interface {
	mustEmbedUnimplementedTLSServiceServer()
}

func RegisterTLSServiceServer(s grpc.ServiceRegistrar, srv TLSServiceServer) {
	s.RegisterService(&TLSService_ServiceDesc, srv)
}

func _TLSService_GetUserTLS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TLSServiceServer).GetUserTLS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TLSService_GetUserTLS_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TLSServiceServer).GetUserTLS(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TLSService_CreateTLS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TLS)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TLSServiceServer).CreateTLS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TLSService_CreateTLS_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TLSServiceServer).CreateTLS(ctx, req.(*TLS))
	}
	return interceptor(ctx, in, info, handler)
}

func _TLSService_UpdateTLS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TLS)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TLSServiceServer).UpdateTLS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TLSService_UpdateTLS_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TLSServiceServer).UpdateTLS(ctx, req.(*TLS))
	}
	return interceptor(ctx, in, info, handler)
}

func _TLSService_DeleteTLS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(public.PrimitiveID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TLSServiceServer).DeleteTLS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TLSService_DeleteTLS_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TLSServiceServer).DeleteTLS(ctx, req.(*public.PrimitiveID))
	}
	return interceptor(ctx, in, info, handler)
}

// TLSService_ServiceDesc is the grpc.ServiceDesc for TLSService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TLSService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bellis.backend.mobile.profile.TLSService",
	HandlerType: (*TLSServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserTLS",
			Handler:    _TLSService_GetUserTLS_Handler,
		},
		{
			MethodName: "CreateTLS",
			Handler:    _TLSService_CreateTLS_Handler,
		},
		{
			MethodName: "UpdateTLS",
			Handler:    _TLSService_UpdateTLS_Handler,
		},
		{
			MethodName: "DeleteTLS",
			Handler:    _TLSService_DeleteTLS_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "profile/tls.proto",
}
