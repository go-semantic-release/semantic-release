// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.13.0
// source: pkg/updater/updater.proto

package updater

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type FilesUpdaterInit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FilesUpdaterInit) Reset() {
	*x = FilesUpdaterInit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterInit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterInit) ProtoMessage() {}

func (x *FilesUpdaterInit) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterInit.ProtoReflect.Descriptor instead.
func (*FilesUpdaterInit) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{0}
}

type FilesUpdaterName struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FilesUpdaterName) Reset() {
	*x = FilesUpdaterName{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterName) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterName) ProtoMessage() {}

func (x *FilesUpdaterName) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterName.ProtoReflect.Descriptor instead.
func (*FilesUpdaterName) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{1}
}

type FilesUpdaterVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FilesUpdaterVersion) Reset() {
	*x = FilesUpdaterVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterVersion) ProtoMessage() {}

func (x *FilesUpdaterVersion) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterVersion.ProtoReflect.Descriptor instead.
func (*FilesUpdaterVersion) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{2}
}

type FilesUpdaterForFiles struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FilesUpdaterForFiles) Reset() {
	*x = FilesUpdaterForFiles{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterForFiles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterForFiles) ProtoMessage() {}

func (x *FilesUpdaterForFiles) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterForFiles.ProtoReflect.Descriptor instead.
func (*FilesUpdaterForFiles) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{3}
}

type FilesUpdaterApply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FilesUpdaterApply) Reset() {
	*x = FilesUpdaterApply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterApply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterApply) ProtoMessage() {}

func (x *FilesUpdaterApply) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterApply.ProtoReflect.Descriptor instead.
func (*FilesUpdaterApply) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{4}
}

type FilesUpdaterInit_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Config map[string]string `protobuf:"bytes,1,rep,name=config,proto3" json:"config,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *FilesUpdaterInit_Request) Reset() {
	*x = FilesUpdaterInit_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterInit_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterInit_Request) ProtoMessage() {}

func (x *FilesUpdaterInit_Request) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterInit_Request.ProtoReflect.Descriptor instead.
func (*FilesUpdaterInit_Request) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{0, 0}
}

func (x *FilesUpdaterInit_Request) GetConfig() map[string]string {
	if x != nil {
		return x.Config
	}
	return nil
}

type FilesUpdaterInit_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *FilesUpdaterInit_Response) Reset() {
	*x = FilesUpdaterInit_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterInit_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterInit_Response) ProtoMessage() {}

func (x *FilesUpdaterInit_Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterInit_Response.ProtoReflect.Descriptor instead.
func (*FilesUpdaterInit_Response) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{0, 1}
}

func (x *FilesUpdaterInit_Response) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type FilesUpdaterName_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FilesUpdaterName_Request) Reset() {
	*x = FilesUpdaterName_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterName_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterName_Request) ProtoMessage() {}

func (x *FilesUpdaterName_Request) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterName_Request.ProtoReflect.Descriptor instead.
func (*FilesUpdaterName_Request) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{1, 0}
}

type FilesUpdaterName_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *FilesUpdaterName_Response) Reset() {
	*x = FilesUpdaterName_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterName_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterName_Response) ProtoMessage() {}

func (x *FilesUpdaterName_Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterName_Response.ProtoReflect.Descriptor instead.
func (*FilesUpdaterName_Response) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{1, 1}
}

func (x *FilesUpdaterName_Response) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type FilesUpdaterVersion_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FilesUpdaterVersion_Request) Reset() {
	*x = FilesUpdaterVersion_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterVersion_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterVersion_Request) ProtoMessage() {}

func (x *FilesUpdaterVersion_Request) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterVersion_Request.ProtoReflect.Descriptor instead.
func (*FilesUpdaterVersion_Request) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{2, 0}
}

type FilesUpdaterVersion_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *FilesUpdaterVersion_Response) Reset() {
	*x = FilesUpdaterVersion_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterVersion_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterVersion_Response) ProtoMessage() {}

func (x *FilesUpdaterVersion_Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterVersion_Response.ProtoReflect.Descriptor instead.
func (*FilesUpdaterVersion_Response) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{2, 1}
}

func (x *FilesUpdaterVersion_Response) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type FilesUpdaterForFiles_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FilesUpdaterForFiles_Request) Reset() {
	*x = FilesUpdaterForFiles_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterForFiles_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterForFiles_Request) ProtoMessage() {}

func (x *FilesUpdaterForFiles_Request) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterForFiles_Request.ProtoReflect.Descriptor instead.
func (*FilesUpdaterForFiles_Request) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{3, 0}
}

type FilesUpdaterForFiles_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Files string `protobuf:"bytes,1,opt,name=files,proto3" json:"files,omitempty"`
}

func (x *FilesUpdaterForFiles_Response) Reset() {
	*x = FilesUpdaterForFiles_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterForFiles_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterForFiles_Response) ProtoMessage() {}

func (x *FilesUpdaterForFiles_Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterForFiles_Response.ProtoReflect.Descriptor instead.
func (*FilesUpdaterForFiles_Response) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{3, 1}
}

func (x *FilesUpdaterForFiles_Response) GetFiles() string {
	if x != nil {
		return x.Files
	}
	return ""
}

type FilesUpdaterApply_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File       string `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	NewVersion string `protobuf:"bytes,2,opt,name=new_version,json=newVersion,proto3" json:"new_version,omitempty"`
}

func (x *FilesUpdaterApply_Request) Reset() {
	*x = FilesUpdaterApply_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterApply_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterApply_Request) ProtoMessage() {}

func (x *FilesUpdaterApply_Request) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterApply_Request.ProtoReflect.Descriptor instead.
func (*FilesUpdaterApply_Request) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{4, 0}
}

func (x *FilesUpdaterApply_Request) GetFile() string {
	if x != nil {
		return x.File
	}
	return ""
}

func (x *FilesUpdaterApply_Request) GetNewVersion() string {
	if x != nil {
		return x.NewVersion
	}
	return ""
}

type FilesUpdaterApply_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *FilesUpdaterApply_Response) Reset() {
	*x = FilesUpdaterApply_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_updater_updater_proto_msgTypes[15]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilesUpdaterApply_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilesUpdaterApply_Response) ProtoMessage() {}

func (x *FilesUpdaterApply_Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_updater_updater_proto_msgTypes[15]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilesUpdaterApply_Response.ProtoReflect.Descriptor instead.
func (*FilesUpdaterApply_Response) Descriptor() ([]byte, []int) {
	return file_pkg_updater_updater_proto_rawDescGZIP(), []int{4, 1}
}

func (x *FilesUpdaterApply_Response) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_pkg_updater_updater_proto protoreflect.FileDescriptor

var file_pkg_updater_updater_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x6b, 0x67, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x2f, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xba, 0x01, 0x0a, 0x10,
	0x46, 0x69, 0x6c, 0x65, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x69, 0x74,
	0x1a, 0x83, 0x01, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3d, 0x0a, 0x06,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x69, 0x74, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x39, 0x0a, 0x0b, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x20, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x3d, 0x0a, 0x10, 0x46, 0x69, 0x6c, 0x65,
	0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x1a, 0x09, 0x0a, 0x07,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x46, 0x0a, 0x13, 0x46, 0x69, 0x6c, 0x65, 0x73,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x09,
	0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22,
	0x43, 0x0a, 0x14, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x46,
	0x6f, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x1a, 0x09, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x20, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x22, 0x75, 0x0a, 0x11, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x72, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x1a, 0x3e, 0x0a, 0x07, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x65, 0x77, 0x5f,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6e,
	0x65, 0x77, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x20, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0xe7, 0x02, 0x0a, 0x12,
	0x46, 0x69, 0x6c, 0x65, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x12, 0x3d, 0x0a, 0x04, 0x49, 0x6e, 0x69, 0x74, 0x12, 0x19, 0x2e, 0x46, 0x69, 0x6c,
	0x65, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x69, 0x74, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x69, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x3d, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x2e, 0x46, 0x69, 0x6c, 0x65,
	0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x46, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x2e, 0x46, 0x69,
	0x6c, 0x65, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x46, 0x69, 0x6c, 0x65,
	0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x49, 0x0a, 0x08, 0x46, 0x6f, 0x72, 0x46,
	0x69, 0x6c, 0x65, 0x73, 0x12, 0x1d, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x72, 0x46, 0x6f, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x72, 0x46, 0x6f, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x05, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x12, 0x1a, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x41, 0x70, 0x70, 0x6c, 0x79,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x73,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x40, 0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x2d, 0x73, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63, 0x2d,
	0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x2f, 0x73, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63,
	0x2d, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x2f, 0x76, 0x32, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_updater_updater_proto_rawDescOnce sync.Once
	file_pkg_updater_updater_proto_rawDescData = file_pkg_updater_updater_proto_rawDesc
)

func file_pkg_updater_updater_proto_rawDescGZIP() []byte {
	file_pkg_updater_updater_proto_rawDescOnce.Do(func() {
		file_pkg_updater_updater_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_updater_updater_proto_rawDescData)
	})
	return file_pkg_updater_updater_proto_rawDescData
}

var file_pkg_updater_updater_proto_msgTypes = make([]protoimpl.MessageInfo, 16)
var file_pkg_updater_updater_proto_goTypes = []interface{}{
	(*FilesUpdaterInit)(nil),              // 0: FilesUpdaterInit
	(*FilesUpdaterName)(nil),              // 1: FilesUpdaterName
	(*FilesUpdaterVersion)(nil),           // 2: FilesUpdaterVersion
	(*FilesUpdaterForFiles)(nil),          // 3: FilesUpdaterForFiles
	(*FilesUpdaterApply)(nil),             // 4: FilesUpdaterApply
	(*FilesUpdaterInit_Request)(nil),      // 5: FilesUpdaterInit.Request
	(*FilesUpdaterInit_Response)(nil),     // 6: FilesUpdaterInit.Response
	nil,                                   // 7: FilesUpdaterInit.Request.ConfigEntry
	(*FilesUpdaterName_Request)(nil),      // 8: FilesUpdaterName.Request
	(*FilesUpdaterName_Response)(nil),     // 9: FilesUpdaterName.Response
	(*FilesUpdaterVersion_Request)(nil),   // 10: FilesUpdaterVersion.Request
	(*FilesUpdaterVersion_Response)(nil),  // 11: FilesUpdaterVersion.Response
	(*FilesUpdaterForFiles_Request)(nil),  // 12: FilesUpdaterForFiles.Request
	(*FilesUpdaterForFiles_Response)(nil), // 13: FilesUpdaterForFiles.Response
	(*FilesUpdaterApply_Request)(nil),     // 14: FilesUpdaterApply.Request
	(*FilesUpdaterApply_Response)(nil),    // 15: FilesUpdaterApply.Response
}
var file_pkg_updater_updater_proto_depIdxs = []int32{
	7,  // 0: FilesUpdaterInit.Request.config:type_name -> FilesUpdaterInit.Request.ConfigEntry
	5,  // 1: FilesUpdaterPlugin.Init:input_type -> FilesUpdaterInit.Request
	8,  // 2: FilesUpdaterPlugin.Name:input_type -> FilesUpdaterName.Request
	10, // 3: FilesUpdaterPlugin.Version:input_type -> FilesUpdaterVersion.Request
	12, // 4: FilesUpdaterPlugin.ForFiles:input_type -> FilesUpdaterForFiles.Request
	14, // 5: FilesUpdaterPlugin.Apply:input_type -> FilesUpdaterApply.Request
	6,  // 6: FilesUpdaterPlugin.Init:output_type -> FilesUpdaterInit.Response
	9,  // 7: FilesUpdaterPlugin.Name:output_type -> FilesUpdaterName.Response
	11, // 8: FilesUpdaterPlugin.Version:output_type -> FilesUpdaterVersion.Response
	13, // 9: FilesUpdaterPlugin.ForFiles:output_type -> FilesUpdaterForFiles.Response
	15, // 10: FilesUpdaterPlugin.Apply:output_type -> FilesUpdaterApply.Response
	6,  // [6:11] is the sub-list for method output_type
	1,  // [1:6] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_updater_updater_proto_init() }
func file_pkg_updater_updater_proto_init() {
	if File_pkg_updater_updater_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_updater_updater_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterInit); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterName); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterVersion); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterForFiles); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterApply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterInit_Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterInit_Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterName_Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterName_Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterVersion_Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterVersion_Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterForFiles_Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterForFiles_Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterApply_Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_updater_updater_proto_msgTypes[15].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilesUpdaterApply_Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_updater_updater_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   16,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_updater_updater_proto_goTypes,
		DependencyIndexes: file_pkg_updater_updater_proto_depIdxs,
		MessageInfos:      file_pkg_updater_updater_proto_msgTypes,
	}.Build()
	File_pkg_updater_updater_proto = out.File
	file_pkg_updater_updater_proto_rawDesc = nil
	file_pkg_updater_updater_proto_goTypes = nil
	file_pkg_updater_updater_proto_depIdxs = nil
}
