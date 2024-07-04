// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.21.12
// source: api/proto/visualization.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	VisualizationService_CreateVisualization_FullMethodName = "/visualization.VisualizationService/CreateVisualization"
	VisualizationService_UpdateVisualization_FullMethodName = "/visualization.VisualizationService/UpdateVisualization"
	VisualizationService_GetVisualization_FullMethodName    = "/visualization.VisualizationService/GetVisualization"
	VisualizationService_ListVisualizations_FullMethodName  = "/visualization.VisualizationService/ListVisualizations"
	VisualizationService_DeleteVisualization_FullMethodName = "/visualization.VisualizationService/DeleteVisualization"
	VisualizationService_ExportVisualization_FullMethodName = "/visualization.VisualizationService/ExportVisualization"
)

// VisualizationServiceClient is the client API for VisualizationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VisualizationServiceClient interface {
	CreateVisualization(ctx context.Context, in *CreateVisualizationRequest, opts ...grpc.CallOption) (*VisualizationResponse, error)
	UpdateVisualization(ctx context.Context, in *UpdateVisualizationRequest, opts ...grpc.CallOption) (*VisualizationResponse, error)
	GetVisualization(ctx context.Context, in *GetVisualizationRequest, opts ...grpc.CallOption) (*VisualizationResponse, error)
	ListVisualizations(ctx context.Context, in *ListVisualizationsRequest, opts ...grpc.CallOption) (*ListVisualizationsResponse, error)
	DeleteVisualization(ctx context.Context, in *DeleteVisualizationRequest, opts ...grpc.CallOption) (*DeleteVisualizationResponse, error)
	ExportVisualization(ctx context.Context, in *ExportVisualizationRequest, opts ...grpc.CallOption) (*ExportVisualizationResponse, error)
}

type visualizationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVisualizationServiceClient(cc grpc.ClientConnInterface) VisualizationServiceClient {
	return &visualizationServiceClient{cc}
}

func (c *visualizationServiceClient) CreateVisualization(ctx context.Context, in *CreateVisualizationRequest, opts ...grpc.CallOption) (*VisualizationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VisualizationResponse)
	err := c.cc.Invoke(ctx, VisualizationService_CreateVisualization_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *visualizationServiceClient) UpdateVisualization(ctx context.Context, in *UpdateVisualizationRequest, opts ...grpc.CallOption) (*VisualizationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VisualizationResponse)
	err := c.cc.Invoke(ctx, VisualizationService_UpdateVisualization_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *visualizationServiceClient) GetVisualization(ctx context.Context, in *GetVisualizationRequest, opts ...grpc.CallOption) (*VisualizationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VisualizationResponse)
	err := c.cc.Invoke(ctx, VisualizationService_GetVisualization_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *visualizationServiceClient) ListVisualizations(ctx context.Context, in *ListVisualizationsRequest, opts ...grpc.CallOption) (*ListVisualizationsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListVisualizationsResponse)
	err := c.cc.Invoke(ctx, VisualizationService_ListVisualizations_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *visualizationServiceClient) DeleteVisualization(ctx context.Context, in *DeleteVisualizationRequest, opts ...grpc.CallOption) (*DeleteVisualizationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteVisualizationResponse)
	err := c.cc.Invoke(ctx, VisualizationService_DeleteVisualization_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *visualizationServiceClient) ExportVisualization(ctx context.Context, in *ExportVisualizationRequest, opts ...grpc.CallOption) (*ExportVisualizationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExportVisualizationResponse)
	err := c.cc.Invoke(ctx, VisualizationService_ExportVisualization_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VisualizationServiceServer is the server API for VisualizationService service.
// All implementations must embed UnimplementedVisualizationServiceServer
// for forward compatibility
type VisualizationServiceServer interface {
	CreateVisualization(context.Context, *CreateVisualizationRequest) (*VisualizationResponse, error)
	UpdateVisualization(context.Context, *UpdateVisualizationRequest) (*VisualizationResponse, error)
	GetVisualization(context.Context, *GetVisualizationRequest) (*VisualizationResponse, error)
	ListVisualizations(context.Context, *ListVisualizationsRequest) (*ListVisualizationsResponse, error)
	DeleteVisualization(context.Context, *DeleteVisualizationRequest) (*DeleteVisualizationResponse, error)
	ExportVisualization(context.Context, *ExportVisualizationRequest) (*ExportVisualizationResponse, error)
	mustEmbedUnimplementedVisualizationServiceServer()
}

// UnimplementedVisualizationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedVisualizationServiceServer struct {
}

func (UnimplementedVisualizationServiceServer) CreateVisualization(context.Context, *CreateVisualizationRequest) (*VisualizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateVisualization not implemented")
}
func (UnimplementedVisualizationServiceServer) UpdateVisualization(context.Context, *UpdateVisualizationRequest) (*VisualizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVisualization not implemented")
}
func (UnimplementedVisualizationServiceServer) GetVisualization(context.Context, *GetVisualizationRequest) (*VisualizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVisualization not implemented")
}
func (UnimplementedVisualizationServiceServer) ListVisualizations(context.Context, *ListVisualizationsRequest) (*ListVisualizationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVisualizations not implemented")
}
func (UnimplementedVisualizationServiceServer) DeleteVisualization(context.Context, *DeleteVisualizationRequest) (*DeleteVisualizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVisualization not implemented")
}
func (UnimplementedVisualizationServiceServer) ExportVisualization(context.Context, *ExportVisualizationRequest) (*ExportVisualizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportVisualization not implemented")
}
func (UnimplementedVisualizationServiceServer) mustEmbedUnimplementedVisualizationServiceServer() {}

// UnsafeVisualizationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VisualizationServiceServer will
// result in compilation errors.
type UnsafeVisualizationServiceServer interface {
	mustEmbedUnimplementedVisualizationServiceServer()
}

func RegisterVisualizationServiceServer(s grpc.ServiceRegistrar, srv VisualizationServiceServer) {
	s.RegisterService(&VisualizationService_ServiceDesc, srv)
}

func _VisualizationService_CreateVisualization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateVisualizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VisualizationServiceServer).CreateVisualization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VisualizationService_CreateVisualization_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VisualizationServiceServer).CreateVisualization(ctx, req.(*CreateVisualizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VisualizationService_UpdateVisualization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateVisualizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VisualizationServiceServer).UpdateVisualization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VisualizationService_UpdateVisualization_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VisualizationServiceServer).UpdateVisualization(ctx, req.(*UpdateVisualizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VisualizationService_GetVisualization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVisualizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VisualizationServiceServer).GetVisualization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VisualizationService_GetVisualization_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VisualizationServiceServer).GetVisualization(ctx, req.(*GetVisualizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VisualizationService_ListVisualizations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListVisualizationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VisualizationServiceServer).ListVisualizations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VisualizationService_ListVisualizations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VisualizationServiceServer).ListVisualizations(ctx, req.(*ListVisualizationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VisualizationService_DeleteVisualization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteVisualizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VisualizationServiceServer).DeleteVisualization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VisualizationService_DeleteVisualization_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VisualizationServiceServer).DeleteVisualization(ctx, req.(*DeleteVisualizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VisualizationService_ExportVisualization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportVisualizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VisualizationServiceServer).ExportVisualization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VisualizationService_ExportVisualization_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VisualizationServiceServer).ExportVisualization(ctx, req.(*ExportVisualizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VisualizationService_ServiceDesc is the grpc.ServiceDesc for VisualizationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VisualizationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "visualization.VisualizationService",
	HandlerType: (*VisualizationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateVisualization",
			Handler:    _VisualizationService_CreateVisualization_Handler,
		},
		{
			MethodName: "UpdateVisualization",
			Handler:    _VisualizationService_UpdateVisualization_Handler,
		},
		{
			MethodName: "GetVisualization",
			Handler:    _VisualizationService_GetVisualization_Handler,
		},
		{
			MethodName: "ListVisualizations",
			Handler:    _VisualizationService_ListVisualizations_Handler,
		},
		{
			MethodName: "DeleteVisualization",
			Handler:    _VisualizationService_DeleteVisualization_Handler,
		},
		{
			MethodName: "ExportVisualization",
			Handler:    _VisualizationService_ExportVisualization_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/visualization.proto",
}