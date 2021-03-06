// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: jarvis.v1.proto

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

// JarvisV1Client is the client API for JarvisV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JarvisV1Client interface {
	ListDailyClose(ctx context.Context, in *ListDailyCloseRequest, opts ...grpc.CallOption) (*ListDailyCloseResponse, error)
	ListStocks(ctx context.Context, in *ListStockRequest, opts ...grpc.CallOption) (*ListStockResponse, error)
	ListCategories(ctx context.Context, in *ListCategoriesRequest, opts ...grpc.CallOption) (*ListCategoriesResponse, error)
	GetStakeConcentration(ctx context.Context, in *GetStakeConcentrationRequest, opts ...grpc.CallOption) (*GetStakeConcentrationResponse, error)
	ListThreePrimary(ctx context.Context, in *ListThreePrimaryRequest, opts ...grpc.CallOption) (*ListThreePrimaryResponse, error)
}

type jarvisV1Client struct {
	cc grpc.ClientConnInterface
}

func NewJarvisV1Client(cc grpc.ClientConnInterface) JarvisV1Client {
	return &jarvisV1Client{cc}
}

func (c *jarvisV1Client) ListDailyClose(ctx context.Context, in *ListDailyCloseRequest, opts ...grpc.CallOption) (*ListDailyCloseResponse, error) {
	out := new(ListDailyCloseResponse)
	err := c.cc.Invoke(ctx, "/jarvis.v1.JarvisV1/ListDailyClose", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jarvisV1Client) ListStocks(ctx context.Context, in *ListStockRequest, opts ...grpc.CallOption) (*ListStockResponse, error) {
	out := new(ListStockResponse)
	err := c.cc.Invoke(ctx, "/jarvis.v1.JarvisV1/ListStocks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jarvisV1Client) ListCategories(ctx context.Context, in *ListCategoriesRequest, opts ...grpc.CallOption) (*ListCategoriesResponse, error) {
	out := new(ListCategoriesResponse)
	err := c.cc.Invoke(ctx, "/jarvis.v1.JarvisV1/ListCategories", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jarvisV1Client) GetStakeConcentration(ctx context.Context, in *GetStakeConcentrationRequest, opts ...grpc.CallOption) (*GetStakeConcentrationResponse, error) {
	out := new(GetStakeConcentrationResponse)
	err := c.cc.Invoke(ctx, "/jarvis.v1.JarvisV1/GetStakeConcentration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jarvisV1Client) ListThreePrimary(ctx context.Context, in *ListThreePrimaryRequest, opts ...grpc.CallOption) (*ListThreePrimaryResponse, error) {
	out := new(ListThreePrimaryResponse)
	err := c.cc.Invoke(ctx, "/jarvis.v1.JarvisV1/ListThreePrimary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JarvisV1Server is the server API for JarvisV1 service.
// All implementations should embed UnimplementedJarvisV1Server
// for forward compatibility
type JarvisV1Server interface {
	ListDailyClose(context.Context, *ListDailyCloseRequest) (*ListDailyCloseResponse, error)
	ListStocks(context.Context, *ListStockRequest) (*ListStockResponse, error)
	ListCategories(context.Context, *ListCategoriesRequest) (*ListCategoriesResponse, error)
	GetStakeConcentration(context.Context, *GetStakeConcentrationRequest) (*GetStakeConcentrationResponse, error)
	ListThreePrimary(context.Context, *ListThreePrimaryRequest) (*ListThreePrimaryResponse, error)
}

// UnimplementedJarvisV1Server should be embedded to have forward compatible implementations.
type UnimplementedJarvisV1Server struct {
}

func (UnimplementedJarvisV1Server) ListDailyClose(context.Context, *ListDailyCloseRequest) (*ListDailyCloseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDailyClose not implemented")
}
func (UnimplementedJarvisV1Server) ListStocks(context.Context, *ListStockRequest) (*ListStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListStocks not implemented")
}
func (UnimplementedJarvisV1Server) ListCategories(context.Context, *ListCategoriesRequest) (*ListCategoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCategories not implemented")
}
func (UnimplementedJarvisV1Server) GetStakeConcentration(context.Context, *GetStakeConcentrationRequest) (*GetStakeConcentrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStakeConcentration not implemented")
}
func (UnimplementedJarvisV1Server) ListThreePrimary(context.Context, *ListThreePrimaryRequest) (*ListThreePrimaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListThreePrimary not implemented")
}

// UnsafeJarvisV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JarvisV1Server will
// result in compilation errors.
type UnsafeJarvisV1Server interface {
	mustEmbedUnimplementedJarvisV1Server()
}

func RegisterJarvisV1Server(s grpc.ServiceRegistrar, srv JarvisV1Server) {
	s.RegisterService(&JarvisV1_ServiceDesc, srv)
}

func _JarvisV1_ListDailyClose_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDailyCloseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JarvisV1Server).ListDailyClose(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jarvis.v1.JarvisV1/ListDailyClose",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JarvisV1Server).ListDailyClose(ctx, req.(*ListDailyCloseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JarvisV1_ListStocks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListStockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JarvisV1Server).ListStocks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jarvis.v1.JarvisV1/ListStocks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JarvisV1Server).ListStocks(ctx, req.(*ListStockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JarvisV1_ListCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCategoriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JarvisV1Server).ListCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jarvis.v1.JarvisV1/ListCategories",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JarvisV1Server).ListCategories(ctx, req.(*ListCategoriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JarvisV1_GetStakeConcentration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStakeConcentrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JarvisV1Server).GetStakeConcentration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jarvis.v1.JarvisV1/GetStakeConcentration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JarvisV1Server).GetStakeConcentration(ctx, req.(*GetStakeConcentrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JarvisV1_ListThreePrimary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListThreePrimaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JarvisV1Server).ListThreePrimary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jarvis.v1.JarvisV1/ListThreePrimary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JarvisV1Server).ListThreePrimary(ctx, req.(*ListThreePrimaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// JarvisV1_ServiceDesc is the grpc.ServiceDesc for JarvisV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var JarvisV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "jarvis.v1.JarvisV1",
	HandlerType: (*JarvisV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListDailyClose",
			Handler:    _JarvisV1_ListDailyClose_Handler,
		},
		{
			MethodName: "ListStocks",
			Handler:    _JarvisV1_ListStocks_Handler,
		},
		{
			MethodName: "ListCategories",
			Handler:    _JarvisV1_ListCategories_Handler,
		},
		{
			MethodName: "GetStakeConcentration",
			Handler:    _JarvisV1_GetStakeConcentration_Handler,
		},
		{
			MethodName: "ListThreePrimary",
			Handler:    _JarvisV1_ListThreePrimary_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "jarvis.v1.proto",
}
