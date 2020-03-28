// Code generated by protoc-gen-go. DO NOT EDIT.
// source: import/service.proto

package _import

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() {
	proto.RegisterFile("import/service.proto", fileDescriptor_6746b50af5597e59)
}

var fileDescriptor_6746b50af5597e59 = []byte{
	// 108 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc9, 0xcc, 0x2d, 0xc8,
	0x2f, 0x2a, 0xd1, 0x2f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x62, 0x83, 0x88, 0x4a, 0x09, 0x43, 0x65, 0x33, 0xf3, 0x4a, 0x52, 0xf3, 0x4a, 0x20, 0x92,
	0x46, 0x01, 0x5c, 0xbc, 0x9e, 0x60, 0xe1, 0x60, 0x88, 0x1e, 0x21, 0x7b, 0x2e, 0xee, 0xe0, 0x92,
	0xc4, 0xa2, 0x12, 0x88, 0xa8, 0x90, 0x88, 0x1e, 0x44, 0x97, 0x1e, 0x84, 0xef, 0x09, 0xd6, 0x2b,
	0x25, 0x85, 0x4d, 0x34, 0x28, 0xb5, 0xb8, 0x34, 0xa7, 0x24, 0x89, 0x0d, 0x6c, 0xb0, 0x31, 0x20,
	0x00, 0x00, 0xff, 0xff, 0x68, 0x0a, 0x8d, 0xfa, 0x8d, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ImportServiceClient is the client API for ImportService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ImportServiceClient interface {
	// Start import process of raw data and returns result of the operation when finished
	StartImport(ctx context.Context, in *ImportIntent, opts ...grpc.CallOption) (*ImportIntentResult, error)
}

type importServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewImportServiceClient(cc grpc.ClientConnInterface) ImportServiceClient {
	return &importServiceClient{cc}
}

func (c *importServiceClient) StartImport(ctx context.Context, in *ImportIntent, opts ...grpc.CallOption) (*ImportIntentResult, error) {
	out := new(ImportIntentResult)
	err := c.cc.Invoke(ctx, "/import.ImportService/StartImport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ImportServiceServer is the server API for ImportService service.
type ImportServiceServer interface {
	// Start import process of raw data and returns result of the operation when finished
	StartImport(context.Context, *ImportIntent) (*ImportIntentResult, error)
}

// UnimplementedImportServiceServer can be embedded to have forward compatible implementations.
type UnimplementedImportServiceServer struct {
}

func (*UnimplementedImportServiceServer) StartImport(ctx context.Context, req *ImportIntent) (*ImportIntentResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartImport not implemented")
}

func RegisterImportServiceServer(s *grpc.Server, srv ImportServiceServer) {
	s.RegisterService(&_ImportService_serviceDesc, srv)
}

func _ImportService_StartImport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportIntent)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImportServiceServer).StartImport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/import.ImportService/StartImport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImportServiceServer).StartImport(ctx, req.(*ImportIntent))
	}
	return interceptor(ctx, in, info, handler)
}

var _ImportService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "import.ImportService",
	HandlerType: (*ImportServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartImport",
			Handler:    _ImportService_StartImport_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "import/service.proto",
}
