/*
Copyright ApeCloud Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: node.proto

package protocol

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

// NodeClient is the client API for Node service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeClient interface {
	DescribeAllENIs(ctx context.Context, in *DescribeAllENIsRequest, opts ...grpc.CallOption) (*DescribeAllENIsResponse, error)
	SetupNetworkForENI(ctx context.Context, in *SetupNetworkForENIRequest, opts ...grpc.CallOption) (*SetupNetworkForENIResponse, error)
	CleanNetworkForENI(ctx context.Context, in *CleanNetworkForENIRequest, opts ...grpc.CallOption) (*CleanNetworkForENIResponse, error)
	SetupNetworkForService(ctx context.Context, in *SetupNetworkForServiceRequest, opts ...grpc.CallOption) (*SetupNetworkForServiceResponse, error)
	CleanNetworkForService(ctx context.Context, in *CleanNetworkForServiceRequest, opts ...grpc.CallOption) (*CleanNetworkForServiceResponse, error)
}

type nodeClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeClient(cc grpc.ClientConnInterface) NodeClient {
	return &nodeClient{cc}
}

func (c *nodeClient) DescribeAllENIs(ctx context.Context, in *DescribeAllENIsRequest, opts ...grpc.CallOption) (*DescribeAllENIsResponse, error) {
	out := new(DescribeAllENIsResponse)
	err := c.cc.Invoke(ctx, "/protocol.Node/DescribeAllENIs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) SetupNetworkForENI(ctx context.Context, in *SetupNetworkForENIRequest, opts ...grpc.CallOption) (*SetupNetworkForENIResponse, error) {
	out := new(SetupNetworkForENIResponse)
	err := c.cc.Invoke(ctx, "/protocol.Node/SetupNetworkForENI", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) CleanNetworkForENI(ctx context.Context, in *CleanNetworkForENIRequest, opts ...grpc.CallOption) (*CleanNetworkForENIResponse, error) {
	out := new(CleanNetworkForENIResponse)
	err := c.cc.Invoke(ctx, "/protocol.Node/CleanNetworkForENI", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) SetupNetworkForService(ctx context.Context, in *SetupNetworkForServiceRequest, opts ...grpc.CallOption) (*SetupNetworkForServiceResponse, error) {
	out := new(SetupNetworkForServiceResponse)
	err := c.cc.Invoke(ctx, "/protocol.Node/SetupNetworkForService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) CleanNetworkForService(ctx context.Context, in *CleanNetworkForServiceRequest, opts ...grpc.CallOption) (*CleanNetworkForServiceResponse, error) {
	out := new(CleanNetworkForServiceResponse)
	err := c.cc.Invoke(ctx, "/protocol.Node/CleanNetworkForService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServer is the server API for Node service.
// All implementations must embed UnimplementedNodeServer
// for forward compatibility
type NodeServer interface {
	DescribeAllENIs(context.Context, *DescribeAllENIsRequest) (*DescribeAllENIsResponse, error)
	SetupNetworkForENI(context.Context, *SetupNetworkForENIRequest) (*SetupNetworkForENIResponse, error)
	CleanNetworkForENI(context.Context, *CleanNetworkForENIRequest) (*CleanNetworkForENIResponse, error)
	SetupNetworkForService(context.Context, *SetupNetworkForServiceRequest) (*SetupNetworkForServiceResponse, error)
	CleanNetworkForService(context.Context, *CleanNetworkForServiceRequest) (*CleanNetworkForServiceResponse, error)
	mustEmbedUnimplementedNodeServer()
}

// UnimplementedNodeServer must be embedded to have forward compatible implementations.
type UnimplementedNodeServer struct {
}

func (UnimplementedNodeServer) DescribeAllENIs(context.Context, *DescribeAllENIsRequest) (*DescribeAllENIsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeAllENIs not implemented")
}
func (UnimplementedNodeServer) SetupNetworkForENI(context.Context, *SetupNetworkForENIRequest) (*SetupNetworkForENIResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetupNetworkForENI not implemented")
}
func (UnimplementedNodeServer) CleanNetworkForENI(context.Context, *CleanNetworkForENIRequest) (*CleanNetworkForENIResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CleanNetworkForENI not implemented")
}
func (UnimplementedNodeServer) SetupNetworkForService(context.Context, *SetupNetworkForServiceRequest) (*SetupNetworkForServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetupNetworkForService not implemented")
}
func (UnimplementedNodeServer) CleanNetworkForService(context.Context, *CleanNetworkForServiceRequest) (*CleanNetworkForServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CleanNetworkForService not implemented")
}
func (UnimplementedNodeServer) mustEmbedUnimplementedNodeServer() {}

// UnsafeNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeServer will
// result in compilation errors.
type UnsafeNodeServer interface {
	mustEmbedUnimplementedNodeServer()
}

func RegisterNodeServer(s grpc.ServiceRegistrar, srv NodeServer) {
	s.RegisterService(&Node_ServiceDesc, srv)
}

func _Node_DescribeAllENIs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeAllENIsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).DescribeAllENIs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Node/DescribeAllENIs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).DescribeAllENIs(ctx, req.(*DescribeAllENIsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_SetupNetworkForENI_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetupNetworkForENIRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).SetupNetworkForENI(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Node/SetupNetworkForENI",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).SetupNetworkForENI(ctx, req.(*SetupNetworkForENIRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_CleanNetworkForENI_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CleanNetworkForENIRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).CleanNetworkForENI(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Node/CleanNetworkForENI",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).CleanNetworkForENI(ctx, req.(*CleanNetworkForENIRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_SetupNetworkForService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetupNetworkForServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).SetupNetworkForService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Node/SetupNetworkForService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).SetupNetworkForService(ctx, req.(*SetupNetworkForServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_CleanNetworkForService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CleanNetworkForServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).CleanNetworkForService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.Node/CleanNetworkForService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).CleanNetworkForService(ctx, req.(*CleanNetworkForServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Node_ServiceDesc is the grpc.ServiceDesc for Node service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Node_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.Node",
	HandlerType: (*NodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DescribeAllENIs",
			Handler:    _Node_DescribeAllENIs_Handler,
		},
		{
			MethodName: "SetupNetworkForENI",
			Handler:    _Node_SetupNetworkForENI_Handler,
		},
		{
			MethodName: "CleanNetworkForENI",
			Handler:    _Node_CleanNetworkForENI_Handler,
		},
		{
			MethodName: "SetupNetworkForService",
			Handler:    _Node_SetupNetworkForService_Handler,
		},
		{
			MethodName: "CleanNetworkForService",
			Handler:    _Node_CleanNetworkForService_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node.proto",
}
