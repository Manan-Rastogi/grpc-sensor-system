// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.3
// source: proto/sensor.proto

package sensor

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SensorService_SendSensorData_FullMethodName      = "/sensor.SensorService/SendSensorData"
	SensorService_GetSensorDataStream_FullMethodName = "/sensor.SensorService/GetSensorDataStream"
)

// SensorServiceClient is the client API for SensorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// gRPC Service
type SensorServiceClient interface {
	SendSensorData(ctx context.Context, in *SensorData, opts ...grpc.CallOption) (*ServerResponse, error)
	// server streaming-> Case here---> Build a live feed system, where I (client) ask you (server) to stream sensor updates for a given sensor ID.
	GetSensorDataStream(ctx context.Context, in *SensorRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SensorData], error)
}

type sensorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSensorServiceClient(cc grpc.ClientConnInterface) SensorServiceClient {
	return &sensorServiceClient{cc}
}

func (c *sensorServiceClient) SendSensorData(ctx context.Context, in *SensorData, opts ...grpc.CallOption) (*ServerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ServerResponse)
	err := c.cc.Invoke(ctx, SensorService_SendSensorData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sensorServiceClient) GetSensorDataStream(ctx context.Context, in *SensorRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SensorData], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &SensorService_ServiceDesc.Streams[0], SensorService_GetSensorDataStream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SensorRequest, SensorData]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SensorService_GetSensorDataStreamClient = grpc.ServerStreamingClient[SensorData]

// SensorServiceServer is the server API for SensorService service.
// All implementations must embed UnimplementedSensorServiceServer
// for forward compatibility.
//
// gRPC Service
type SensorServiceServer interface {
	SendSensorData(context.Context, *SensorData) (*ServerResponse, error)
	// server streaming-> Case here---> Build a live feed system, where I (client) ask you (server) to stream sensor updates for a given sensor ID.
	GetSensorDataStream(*SensorRequest, grpc.ServerStreamingServer[SensorData]) error
	mustEmbedUnimplementedSensorServiceServer()
}

// UnimplementedSensorServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSensorServiceServer struct{}

func (UnimplementedSensorServiceServer) SendSensorData(context.Context, *SensorData) (*ServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendSensorData not implemented")
}
func (UnimplementedSensorServiceServer) GetSensorDataStream(*SensorRequest, grpc.ServerStreamingServer[SensorData]) error {
	return status.Errorf(codes.Unimplemented, "method GetSensorDataStream not implemented")
}
func (UnimplementedSensorServiceServer) mustEmbedUnimplementedSensorServiceServer() {}
func (UnimplementedSensorServiceServer) testEmbeddedByValue()                       {}

// UnsafeSensorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SensorServiceServer will
// result in compilation errors.
type UnsafeSensorServiceServer interface {
	mustEmbedUnimplementedSensorServiceServer()
}

func RegisterSensorServiceServer(s grpc.ServiceRegistrar, srv SensorServiceServer) {
	// If the following call pancis, it indicates UnimplementedSensorServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SensorService_ServiceDesc, srv)
}

func _SensorService_SendSensorData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SensorData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SensorServiceServer).SendSensorData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SensorService_SendSensorData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SensorServiceServer).SendSensorData(ctx, req.(*SensorData))
	}
	return interceptor(ctx, in, info, handler)
}

func _SensorService_GetSensorDataStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SensorRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SensorServiceServer).GetSensorDataStream(m, &grpc.GenericServerStream[SensorRequest, SensorData]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SensorService_GetSensorDataStreamServer = grpc.ServerStreamingServer[SensorData]

// SensorService_ServiceDesc is the grpc.ServiceDesc for SensorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SensorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sensor.SensorService",
	HandlerType: (*SensorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendSensorData",
			Handler:    _SensorService_SendSensorData_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetSensorDataStream",
			Handler:       _SensorService_GetSensorDataStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/sensor.proto",
}
