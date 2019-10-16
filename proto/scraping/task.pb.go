// Code generated by protoc-gen-go. DO NOT EDIT.
// source: scraping/task.proto

package scraping

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
	common "shitty.moe/satelit-project/satelit-scraper/proto/common"
	data "shitty.moe/satelit-project/satelit-scraper/proto/data"
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

// Represents a task for anime pages scraping
type Task struct {
	// Task ID
	Id *common.UUID `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// External DB from where to scrape info
	Source data.Source `protobuf:"varint,2,opt,name=source,proto3,enum=data.Source" json:"source,omitempty"`
	// Scraping jobs
	Jobs                 []*Job   `protobuf:"bytes,3,rep,name=jobs,proto3" json:"jobs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Task) Reset()         { *m = Task{} }
func (m *Task) String() string { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()    {}
func (*Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_69f65971c76f182a, []int{0}
}

func (m *Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Task.Unmarshal(m, b)
}
func (m *Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Task.Marshal(b, m, deterministic)
}
func (m *Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Task.Merge(m, src)
}
func (m *Task) XXX_Size() int {
	return xxx_messageInfo_Task.Size(m)
}
func (m *Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Task proto.InternalMessageInfo

func (m *Task) GetId() *common.UUID {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Task) GetSource() data.Source {
	if m != nil {
		return m.Source
	}
	return data.Source_UNKNOWN
}

func (m *Task) GetJobs() []*Job {
	if m != nil {
		return m.Jobs
	}
	return nil
}

// Represents a single scraping job for an anime page
type Job struct {
	// Job ID
	Id *common.UUID `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Anime ID
	AnimeId              int32    `protobuf:"zigzag32,2,opt,name=anime_id,json=animeId,proto3" json:"anime_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Job) Reset()         { *m = Job{} }
func (m *Job) String() string { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()    {}
func (*Job) Descriptor() ([]byte, []int) {
	return fileDescriptor_69f65971c76f182a, []int{1}
}

func (m *Job) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Job.Unmarshal(m, b)
}
func (m *Job) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Job.Marshal(b, m, deterministic)
}
func (m *Job) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Job.Merge(m, src)
}
func (m *Job) XXX_Size() int {
	return xxx_messageInfo_Job.Size(m)
}
func (m *Job) XXX_DiscardUnknown() {
	xxx_messageInfo_Job.DiscardUnknown(m)
}

var xxx_messageInfo_Job proto.InternalMessageInfo

func (m *Job) GetId() *common.UUID {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Job) GetAnimeId() int32 {
	if m != nil {
		return m.AnimeId
	}
	return 0
}

// Scrape task creation request
type TaskCreate struct {
	// Maximum number of entities to scrape
	Limit int32 `protobuf:"zigzag32,1,opt,name=limit,proto3" json:"limit,omitempty"`
	// External data source to scrape data from
	Source               data.Source `protobuf:"varint,2,opt,name=source,proto3,enum=data.Source" json:"source,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *TaskCreate) Reset()         { *m = TaskCreate{} }
func (m *TaskCreate) String() string { return proto.CompactTextString(m) }
func (*TaskCreate) ProtoMessage()    {}
func (*TaskCreate) Descriptor() ([]byte, []int) {
	return fileDescriptor_69f65971c76f182a, []int{2}
}

func (m *TaskCreate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskCreate.Unmarshal(m, b)
}
func (m *TaskCreate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskCreate.Marshal(b, m, deterministic)
}
func (m *TaskCreate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskCreate.Merge(m, src)
}
func (m *TaskCreate) XXX_Size() int {
	return xxx_messageInfo_TaskCreate.Size(m)
}
func (m *TaskCreate) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskCreate.DiscardUnknown(m)
}

var xxx_messageInfo_TaskCreate proto.InternalMessageInfo

func (m *TaskCreate) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *TaskCreate) GetSource() data.Source {
	if m != nil {
		return m.Source
	}
	return data.Source_UNKNOWN
}

// Intermediate result of a parse task
type TaskYield struct {
	// ID of the related task
	TaskId *common.UUID `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	// ID of the related job
	JobId *common.UUID `protobuf:"bytes,2,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	// Parsed anime entity
	Anime                *data.Anime `protobuf:"bytes,3,opt,name=anime,proto3" json:"anime,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *TaskYield) Reset()         { *m = TaskYield{} }
func (m *TaskYield) String() string { return proto.CompactTextString(m) }
func (*TaskYield) ProtoMessage()    {}
func (*TaskYield) Descriptor() ([]byte, []int) {
	return fileDescriptor_69f65971c76f182a, []int{3}
}

func (m *TaskYield) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskYield.Unmarshal(m, b)
}
func (m *TaskYield) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskYield.Marshal(b, m, deterministic)
}
func (m *TaskYield) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskYield.Merge(m, src)
}
func (m *TaskYield) XXX_Size() int {
	return xxx_messageInfo_TaskYield.Size(m)
}
func (m *TaskYield) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskYield.DiscardUnknown(m)
}

var xxx_messageInfo_TaskYield proto.InternalMessageInfo

func (m *TaskYield) GetTaskId() *common.UUID {
	if m != nil {
		return m.TaskId
	}
	return nil
}

func (m *TaskYield) GetJobId() *common.UUID {
	if m != nil {
		return m.JobId
	}
	return nil
}

func (m *TaskYield) GetAnime() *data.Anime {
	if m != nil {
		return m.Anime
	}
	return nil
}

// Signals that a task has been finished
type TaskFinish struct {
	// ID of the related task
	TaskId               *common.UUID `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *TaskFinish) Reset()         { *m = TaskFinish{} }
func (m *TaskFinish) String() string { return proto.CompactTextString(m) }
func (*TaskFinish) ProtoMessage()    {}
func (*TaskFinish) Descriptor() ([]byte, []int) {
	return fileDescriptor_69f65971c76f182a, []int{4}
}

func (m *TaskFinish) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskFinish.Unmarshal(m, b)
}
func (m *TaskFinish) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskFinish.Marshal(b, m, deterministic)
}
func (m *TaskFinish) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskFinish.Merge(m, src)
}
func (m *TaskFinish) XXX_Size() int {
	return xxx_messageInfo_TaskFinish.Size(m)
}
func (m *TaskFinish) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskFinish.DiscardUnknown(m)
}

var xxx_messageInfo_TaskFinish proto.InternalMessageInfo

func (m *TaskFinish) GetTaskId() *common.UUID {
	if m != nil {
		return m.TaskId
	}
	return nil
}

func init() {
	proto.RegisterType((*Task)(nil), "scraping.Task")
	proto.RegisterType((*Job)(nil), "scraping.Job")
	proto.RegisterType((*TaskCreate)(nil), "scraping.TaskCreate")
	proto.RegisterType((*TaskYield)(nil), "scraping.TaskYield")
	proto.RegisterType((*TaskFinish)(nil), "scraping.TaskFinish")
}

func init() { proto.RegisterFile("scraping/task.proto", fileDescriptor_69f65971c76f182a) }

var fileDescriptor_69f65971c76f182a = []byte{
	// 290 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0xd5, 0xa6, 0x4d, 0xcb, 0x15, 0x10, 0x31, 0x0c, 0xa1, 0x53, 0x1b, 0x18, 0x3a, 0xb9,
	0x22, 0xac, 0x2c, 0x08, 0x84, 0x48, 0x47, 0x43, 0x07, 0xa6, 0xca, 0x8e, 0x23, 0x70, 0xd2, 0xc4,
	0x55, 0x9c, 0xf0, 0xfb, 0x91, 0xcf, 0xc9, 0x80, 0x84, 0x50, 0xc7, 0x7b, 0xef, 0xfc, 0xde, 0x67,
	0x1b, 0x2e, 0x4d, 0x5a, 0xf3, 0x83, 0xaa, 0x3e, 0xd7, 0x0d, 0x37, 0x05, 0x3d, 0xd4, 0xba, 0xd1,
	0x64, 0xda, 0x8b, 0xf3, 0x20, 0xd5, 0x65, 0xa9, 0xab, 0x75, 0xdb, 0x2a, 0xe9, 0xcc, 0xf9, 0x85,
	0xe4, 0x0d, 0x5f, 0xf3, 0x4a, 0x95, 0x59, 0xa7, 0x04, 0xa8, 0x18, 0xdd, 0xd6, 0x69, 0x27, 0x45,
	0x05, 0x8c, 0xde, 0xb9, 0x29, 0xc8, 0x1c, 0x86, 0x4a, 0x86, 0x83, 0xc5, 0x60, 0x35, 0x8b, 0x81,
	0x62, 0xca, 0x76, 0x9b, 0x3c, 0xb3, 0xa1, 0x92, 0xe4, 0x16, 0x7c, 0x77, 0x26, 0x1c, 0x2e, 0x06,
	0xab, 0xf3, 0xf8, 0x94, 0xda, 0x1c, 0xfa, 0x86, 0x1a, 0xeb, 0x3c, 0xb2, 0x84, 0x51, 0xae, 0x85,
	0x09, 0xbd, 0x85, 0xb7, 0x9a, 0xc5, 0x67, 0xb4, 0x47, 0xa3, 0x1b, 0x2d, 0x18, 0x5a, 0xd1, 0x03,
	0x78, 0x1b, 0x2d, 0xfe, 0xed, 0xba, 0x86, 0x29, 0x12, 0xef, 0x94, 0xc4, 0xb6, 0x80, 0x4d, 0x70,
	0x4e, 0x64, 0xf4, 0x0a, 0x60, 0x51, 0x9f, 0xea, 0x8c, 0x37, 0x19, 0xb9, 0x82, 0xf1, 0x5e, 0x95,
	0xaa, 0xc1, 0x9c, 0x80, 0xb9, 0xe1, 0x38, 0xd4, 0xe8, 0x1b, 0x4e, 0x6c, 0xd2, 0x87, 0xca, 0xf6,
	0x92, 0xdc, 0xc0, 0xc4, 0xbe, 0xe8, 0xee, 0x4f, 0x24, 0xdf, 0x5a, 0x89, 0x24, 0x4b, 0xf0, 0x73,
	0x2d, 0x7a, 0xa8, 0xdf, 0x3b, 0xe3, 0x5c, 0x0b, 0x5c, 0x19, 0x23, 0x69, 0xe8, 0xe1, 0xc6, 0xcc,
	0x35, 0x3f, 0x5a, 0x89, 0x39, 0x27, 0xba, 0x73, 0x37, 0x78, 0x51, 0x95, 0x32, 0x5f, 0x47, 0x15,
	0x0b, 0x1f, 0xbf, 0xe9, 0xfe, 0x27, 0x00, 0x00, 0xff, 0xff, 0x8f, 0x8a, 0x0d, 0xe3, 0xff, 0x01,
	0x00, 0x00,
}
