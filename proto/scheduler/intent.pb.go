// Code generated by protoc-gen-go. DO NOT EDIT.
// source: scheduler/intent.proto

package scheduler

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// External DB
type ImportIntent_Source int32

const (
	ImportIntent_UNKNOWN ImportIntent_Source = 0
	ImportIntent_ANIDB   ImportIntent_Source = 1
)

var ImportIntent_Source_name = map[int32]string{
	0: "UNKNOWN",
	1: "ANIDB",
}

var ImportIntent_Source_value = map[string]int32{
	"UNKNOWN": 0,
	"ANIDB":   1,
}

func (x ImportIntent_Source) String() string {
	return proto.EnumName(ImportIntent_Source_name, int32(x))
}

func (ImportIntent_Source) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_85717b989cca20ad, []int{0, 0}
}

// Asks to import anime titles index and schedule new titles for scraping
type ImportIntent struct {
	// Intent ID
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Represents an external DB from where anime titles index should be imported
	Source ImportIntent_Source `protobuf:"varint,2,opt,name=source,proto3,enum=scheduler.ImportIntent_Source" json:"source,omitempty"`
	// Identifiers of anime titles that should be re-imported
	ReimportIds []int32 `protobuf:"zigzag32,3,rep,packed,name=reimport_ids,json=reimportIds,proto3" json:"reimport_ids,omitempty"`
	// URL to send request with import result
	CallbackUrl          string   `protobuf:"bytes,4,opt,name=callback_url,json=callbackUrl,proto3" json:"callback_url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImportIntent) Reset()         { *m = ImportIntent{} }
func (m *ImportIntent) String() string { return proto.CompactTextString(m) }
func (*ImportIntent) ProtoMessage()    {}
func (*ImportIntent) Descriptor() ([]byte, []int) {
	return fileDescriptor_85717b989cca20ad, []int{0}
}

func (m *ImportIntent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImportIntent.Unmarshal(m, b)
}
func (m *ImportIntent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImportIntent.Marshal(b, m, deterministic)
}
func (m *ImportIntent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImportIntent.Merge(m, src)
}
func (m *ImportIntent) XXX_Size() int {
	return xxx_messageInfo_ImportIntent.Size(m)
}
func (m *ImportIntent) XXX_DiscardUnknown() {
	xxx_messageInfo_ImportIntent.DiscardUnknown(m)
}

var xxx_messageInfo_ImportIntent proto.InternalMessageInfo

func (m *ImportIntent) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ImportIntent) GetSource() ImportIntent_Source {
	if m != nil {
		return m.Source
	}
	return ImportIntent_UNKNOWN
}

func (m *ImportIntent) GetReimportIds() []int32 {
	if m != nil {
		return m.ReimportIds
	}
	return nil
}

func (m *ImportIntent) GetCallbackUrl() string {
	if m != nil {
		return m.CallbackUrl
	}
	return ""
}

type ImportIntentResult struct {
	// Intent ID
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// If import succeeded then `true`, `false` otherwise
	Succeeded bool `protobuf:"varint,2,opt,name=succeeded,proto3" json:"succeeded,omitempty"`
	// IDs of anime titles that was not imported
	SkippedIds []int32 `protobuf:"zigzag32,3,rep,packed,name=skipped_ids,json=skippedIds,proto3" json:"skipped_ids,omitempty"`
	// Description of the error if import failed
	ErrorDescription     string   `protobuf:"bytes,4,opt,name=error_description,json=errorDescription,proto3" json:"error_description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImportIntentResult) Reset()         { *m = ImportIntentResult{} }
func (m *ImportIntentResult) String() string { return proto.CompactTextString(m) }
func (*ImportIntentResult) ProtoMessage()    {}
func (*ImportIntentResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_85717b989cca20ad, []int{1}
}

func (m *ImportIntentResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImportIntentResult.Unmarshal(m, b)
}
func (m *ImportIntentResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImportIntentResult.Marshal(b, m, deterministic)
}
func (m *ImportIntentResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImportIntentResult.Merge(m, src)
}
func (m *ImportIntentResult) XXX_Size() int {
	return xxx_messageInfo_ImportIntentResult.Size(m)
}
func (m *ImportIntentResult) XXX_DiscardUnknown() {
	xxx_messageInfo_ImportIntentResult.DiscardUnknown(m)
}

var xxx_messageInfo_ImportIntentResult proto.InternalMessageInfo

func (m *ImportIntentResult) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ImportIntentResult) GetSucceeded() bool {
	if m != nil {
		return m.Succeeded
	}
	return false
}

func (m *ImportIntentResult) GetSkippedIds() []int32 {
	if m != nil {
		return m.SkippedIds
	}
	return nil
}

func (m *ImportIntentResult) GetErrorDescription() string {
	if m != nil {
		return m.ErrorDescription
	}
	return ""
}

func init() {
	proto.RegisterEnum("scheduler.ImportIntent_Source", ImportIntent_Source_name, ImportIntent_Source_value)
	proto.RegisterType((*ImportIntent)(nil), "scheduler.ImportIntent")
	proto.RegisterType((*ImportIntentResult)(nil), "scheduler.ImportIntentResult")
}

func init() { proto.RegisterFile("scheduler/intent.proto", fileDescriptor_85717b989cca20ad) }

var fileDescriptor_85717b989cca20ad = []byte{
	// 265 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0x41, 0x4b, 0xc3, 0x30,
	0x14, 0x80, 0xcd, 0xa6, 0xd5, 0xbe, 0x8e, 0xd1, 0xe5, 0x20, 0x3d, 0x88, 0xd6, 0x9e, 0x0a, 0x42,
	0x05, 0x05, 0xef, 0xca, 0x2e, 0x45, 0xa8, 0x50, 0x19, 0x1e, 0xcb, 0x96, 0x3c, 0x30, 0x2c, 0x36,
	0xe5, 0xa5, 0xf9, 0x1f, 0xfe, 0x22, 0x7f, 0x9b, 0x2c, 0x76, 0x5b, 0xc1, 0xeb, 0x97, 0xef, 0x25,
	0xdf, 0x0b, 0x5c, 0x5a, 0xf1, 0x89, 0xd2, 0x69, 0xa4, 0x7b, 0xd5, 0xf6, 0xd8, 0xf6, 0x45, 0x47,
	0xa6, 0x37, 0x3c, 0x3c, 0xf0, 0xec, 0x87, 0xc1, 0xac, 0xfc, 0xea, 0x0c, 0xf5, 0xa5, 0x37, 0xf8,
	0x1c, 0x26, 0x4a, 0x26, 0x2c, 0x65, 0x79, 0x58, 0x4f, 0x94, 0xe4, 0x4f, 0x10, 0x58, 0xe3, 0x48,
	0x60, 0x32, 0x49, 0x59, 0x3e, 0x7f, 0xb8, 0x2e, 0x0e, 0xc3, 0xc5, 0x78, 0xb0, 0x78, 0xf7, 0x56,
	0x3d, 0xd8, 0xfc, 0x16, 0x66, 0x84, 0xca, 0x0b, 0x8d, 0x92, 0x36, 0x99, 0xa6, 0xd3, 0x7c, 0x51,
	0x47, 0x7b, 0x56, 0x4a, 0xbb, 0x53, 0xc4, 0x5a, 0xeb, 0xcd, 0x5a, 0x6c, 0x1b, 0x47, 0x3a, 0x39,
	0xf5, 0x8f, 0x46, 0x7b, 0xb6, 0x22, 0x9d, 0xa5, 0x10, 0xfc, 0xdd, 0xcb, 0x23, 0x38, 0x5f, 0x55,
	0xaf, 0xd5, 0xdb, 0x47, 0x15, 0x9f, 0xf0, 0x10, 0xce, 0x9e, 0xab, 0x72, 0xf9, 0x12, 0xb3, 0xec,
	0x9b, 0x01, 0x1f, 0x77, 0xd4, 0x68, 0x9d, 0xfe, 0xbf, 0xc6, 0x15, 0x84, 0xd6, 0x09, 0x81, 0x28,
	0x51, 0xfa, 0x4d, 0x2e, 0xea, 0x23, 0xe0, 0x37, 0x10, 0xd9, 0xad, 0xea, 0x3a, 0x94, 0xa3, 0x56,
	0x18, 0xd0, 0x2e, 0xf5, 0x0e, 0x16, 0x48, 0x64, 0xa8, 0x91, 0x68, 0x05, 0xa9, 0xae, 0x57, 0xa6,
	0x1d, 0x7a, 0x63, 0x7f, 0xb0, 0x3c, 0xf2, 0x4d, 0xe0, 0x7f, 0xf9, 0xf1, 0x37, 0x00, 0x00, 0xff,
	0xff, 0xfd, 0xb0, 0x59, 0x79, 0x7f, 0x01, 0x00, 0x00,
}
