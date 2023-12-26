// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.2
// source: profile/profile.proto

package profile

import (
	context "context"
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
	ProfileService_GetUserProfile_FullMethodName       = "/bellis.backend.mobile.profile.ProfileService/GetUserProfile"
	ProfileService_ChangePassword_FullMethodName       = "/bellis.backend.mobile.profile.ProfileService/ChangePassword"
	ProfileService_ChangeEmail_FullMethodName          = "/bellis.backend.mobile.profile.ProfileService/ChangeEmail"
	ProfileService_ChangeAlert_FullMethodName          = "/bellis.backend.mobile.profile.ProfileService/ChangeAlert"
	ProfileService_ChangeSensitive_FullMethodName      = "/bellis.backend.mobile.profile.ProfileService/ChangeSensitive"
	ProfileService_UseGotify_FullMethodName            = "/bellis.backend.mobile.profile.ProfileService/UseGotify"
	ProfileService_UseEmail_FullMethodName             = "/bellis.backend.mobile.profile.ProfileService/UseEmail"
	ProfileService_UseWebhook_FullMethodName           = "/bellis.backend.mobile.profile.ProfileService/UseWebhook"
	ProfileService_GetEnvoyTelegramLink_FullMethodName = "/bellis.backend.mobile.profile.ProfileService/GetEnvoyTelegramLink"
	ProfileService_GetUserLoginLogs_FullMethodName     = "/bellis.backend.mobile.profile.ProfileService/GetUserLoginLogs"
)

// ProfileServiceClient is the client API for ProfileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProfileServiceClient interface {
	GetUserProfile(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserProfile, error)
	ChangePassword(ctx context.Context, in *NewPassword, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ChangeEmail(ctx context.Context, in *NewEmail, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ChangeAlert(ctx context.Context, in *Alert, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ChangeSensitive(ctx context.Context, in *Sensitive, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UseGotify(ctx context.Context, in *Gotify, opts ...grpc.CallOption) (*EnvoyPolicy, error)
	UseEmail(ctx context.Context, in *Email, opts ...grpc.CallOption) (*EnvoyPolicy, error)
	UseWebhook(ctx context.Context, in *Webhook, opts ...grpc.CallOption) (*EnvoyPolicy, error)
	GetEnvoyTelegramLink(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*EnvoyTelegramLink, error)
	GetUserLoginLogs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserLoginLogs, error)
}

type profileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProfileServiceClient(cc grpc.ClientConnInterface) ProfileServiceClient {
	return &profileServiceClient{cc}
}

func (c *profileServiceClient) GetUserProfile(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserProfile, error) {
	out := new(UserProfile)
	err := c.cc.Invoke(ctx, ProfileService_GetUserProfile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) ChangePassword(ctx context.Context, in *NewPassword, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ProfileService_ChangePassword_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) ChangeEmail(ctx context.Context, in *NewEmail, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ProfileService_ChangeEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) ChangeAlert(ctx context.Context, in *Alert, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ProfileService_ChangeAlert_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) ChangeSensitive(ctx context.Context, in *Sensitive, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ProfileService_ChangeSensitive_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) UseGotify(ctx context.Context, in *Gotify, opts ...grpc.CallOption) (*EnvoyPolicy, error) {
	out := new(EnvoyPolicy)
	err := c.cc.Invoke(ctx, ProfileService_UseGotify_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) UseEmail(ctx context.Context, in *Email, opts ...grpc.CallOption) (*EnvoyPolicy, error) {
	out := new(EnvoyPolicy)
	err := c.cc.Invoke(ctx, ProfileService_UseEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) UseWebhook(ctx context.Context, in *Webhook, opts ...grpc.CallOption) (*EnvoyPolicy, error) {
	out := new(EnvoyPolicy)
	err := c.cc.Invoke(ctx, ProfileService_UseWebhook_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) GetEnvoyTelegramLink(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*EnvoyTelegramLink, error) {
	out := new(EnvoyTelegramLink)
	err := c.cc.Invoke(ctx, ProfileService_GetEnvoyTelegramLink_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) GetUserLoginLogs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserLoginLogs, error) {
	out := new(UserLoginLogs)
	err := c.cc.Invoke(ctx, ProfileService_GetUserLoginLogs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileServiceServer is the server API for ProfileService service.
// All implementations should embed UnimplementedProfileServiceServer
// for forward compatibility
type ProfileServiceServer interface {
	GetUserProfile(context.Context, *emptypb.Empty) (*UserProfile, error)
	ChangePassword(context.Context, *NewPassword) (*emptypb.Empty, error)
	ChangeEmail(context.Context, *NewEmail) (*emptypb.Empty, error)
	ChangeAlert(context.Context, *Alert) (*emptypb.Empty, error)
	ChangeSensitive(context.Context, *Sensitive) (*emptypb.Empty, error)
	UseGotify(context.Context, *Gotify) (*EnvoyPolicy, error)
	UseEmail(context.Context, *Email) (*EnvoyPolicy, error)
	UseWebhook(context.Context, *Webhook) (*EnvoyPolicy, error)
	GetEnvoyTelegramLink(context.Context, *emptypb.Empty) (*EnvoyTelegramLink, error)
	GetUserLoginLogs(context.Context, *emptypb.Empty) (*UserLoginLogs, error)
}

// UnimplementedProfileServiceServer should be embedded to have forward compatible implementations.
type UnimplementedProfileServiceServer struct {
}

func (UnimplementedProfileServiceServer) GetUserProfile(context.Context, *emptypb.Empty) (*UserProfile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserProfile not implemented")
}
func (UnimplementedProfileServiceServer) ChangePassword(context.Context, *NewPassword) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (UnimplementedProfileServiceServer) ChangeEmail(context.Context, *NewEmail) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeEmail not implemented")
}
func (UnimplementedProfileServiceServer) ChangeAlert(context.Context, *Alert) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeAlert not implemented")
}
func (UnimplementedProfileServiceServer) ChangeSensitive(context.Context, *Sensitive) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeSensitive not implemented")
}
func (UnimplementedProfileServiceServer) UseGotify(context.Context, *Gotify) (*EnvoyPolicy, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseGotify not implemented")
}
func (UnimplementedProfileServiceServer) UseEmail(context.Context, *Email) (*EnvoyPolicy, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseEmail not implemented")
}
func (UnimplementedProfileServiceServer) UseWebhook(context.Context, *Webhook) (*EnvoyPolicy, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseWebhook not implemented")
}
func (UnimplementedProfileServiceServer) GetEnvoyTelegramLink(context.Context, *emptypb.Empty) (*EnvoyTelegramLink, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEnvoyTelegramLink not implemented")
}
func (UnimplementedProfileServiceServer) GetUserLoginLogs(context.Context, *emptypb.Empty) (*UserLoginLogs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserLoginLogs not implemented")
}

// UnsafeProfileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProfileServiceServer will
// result in compilation errors.
type UnsafeProfileServiceServer interface {
	mustEmbedUnimplementedProfileServiceServer()
}

func RegisterProfileServiceServer(s grpc.ServiceRegistrar, srv ProfileServiceServer) {
	s.RegisterService(&ProfileService_ServiceDesc, srv)
}

func _ProfileService_GetUserProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetUserProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_GetUserProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetUserProfile(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewPassword)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_ChangePassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).ChangePassword(ctx, req.(*NewPassword))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_ChangeEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewEmail)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).ChangeEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_ChangeEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).ChangeEmail(ctx, req.(*NewEmail))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_ChangeAlert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Alert)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).ChangeAlert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_ChangeAlert_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).ChangeAlert(ctx, req.(*Alert))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_ChangeSensitive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Sensitive)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).ChangeSensitive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_ChangeSensitive_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).ChangeSensitive(ctx, req.(*Sensitive))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_UseGotify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Gotify)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).UseGotify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_UseGotify_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).UseGotify(ctx, req.(*Gotify))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_UseEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Email)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).UseEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_UseEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).UseEmail(ctx, req.(*Email))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_UseWebhook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Webhook)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).UseWebhook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_UseWebhook_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).UseWebhook(ctx, req.(*Webhook))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_GetEnvoyTelegramLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetEnvoyTelegramLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_GetEnvoyTelegramLink_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetEnvoyTelegramLink(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_GetUserLoginLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetUserLoginLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_GetUserLoginLogs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetUserLoginLogs(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// ProfileService_ServiceDesc is the grpc.ServiceDesc for ProfileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProfileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bellis.backend.mobile.profile.ProfileService",
	HandlerType: (*ProfileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserProfile",
			Handler:    _ProfileService_GetUserProfile_Handler,
		},
		{
			MethodName: "ChangePassword",
			Handler:    _ProfileService_ChangePassword_Handler,
		},
		{
			MethodName: "ChangeEmail",
			Handler:    _ProfileService_ChangeEmail_Handler,
		},
		{
			MethodName: "ChangeAlert",
			Handler:    _ProfileService_ChangeAlert_Handler,
		},
		{
			MethodName: "ChangeSensitive",
			Handler:    _ProfileService_ChangeSensitive_Handler,
		},
		{
			MethodName: "UseGotify",
			Handler:    _ProfileService_UseGotify_Handler,
		},
		{
			MethodName: "UseEmail",
			Handler:    _ProfileService_UseEmail_Handler,
		},
		{
			MethodName: "UseWebhook",
			Handler:    _ProfileService_UseWebhook_Handler,
		},
		{
			MethodName: "GetEnvoyTelegramLink",
			Handler:    _ProfileService_GetEnvoyTelegramLink_Handler,
		},
		{
			MethodName: "GetUserLoginLogs",
			Handler:    _ProfileService_GetUserLoginLogs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "profile/profile.proto",
}
