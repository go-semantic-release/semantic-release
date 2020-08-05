// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: pkg/semrel/structs.proto

package semrel

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

type RawCommit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SHA        string `protobuf:"bytes,1,opt,name=SHA,proto3" json:"SHA,omitempty"`
	RawMessage string `protobuf:"bytes,2,opt,name=RawMessage,proto3" json:"RawMessage,omitempty"`
}

func (x *RawCommit) Reset() {
	*x = RawCommit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_semrel_structs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RawCommit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawCommit) ProtoMessage() {}

func (x *RawCommit) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_semrel_structs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawCommit.ProtoReflect.Descriptor instead.
func (*RawCommit) Descriptor() ([]byte, []int) {
	return file_pkg_semrel_structs_proto_rawDescGZIP(), []int{0}
}

func (x *RawCommit) GetSHA() string {
	if x != nil {
		return x.SHA
	}
	return ""
}

func (x *RawCommit) GetRawMessage() string {
	if x != nil {
		return x.RawMessage
	}
	return ""
}

type Change struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Major bool `protobuf:"varint,1,opt,name=Major,proto3" json:"Major,omitempty"`
	Minor bool `protobuf:"varint,2,opt,name=Minor,proto3" json:"Minor,omitempty"`
	Patch bool `protobuf:"varint,3,opt,name=Patch,proto3" json:"Patch,omitempty"`
}

func (x *Change) Reset() {
	*x = Change{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_semrel_structs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Change) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Change) ProtoMessage() {}

func (x *Change) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_semrel_structs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Change.ProtoReflect.Descriptor instead.
func (*Change) Descriptor() ([]byte, []int) {
	return file_pkg_semrel_structs_proto_rawDescGZIP(), []int{1}
}

func (x *Change) GetMajor() bool {
	if x != nil {
		return x.Major
	}
	return false
}

func (x *Change) GetMinor() bool {
	if x != nil {
		return x.Minor
	}
	return false
}

func (x *Change) GetPatch() bool {
	if x != nil {
		return x.Patch
	}
	return false
}

type Commit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SHA     string   `protobuf:"bytes,1,opt,name=SHA,proto3" json:"SHA,omitempty"`
	Raw     []string `protobuf:"bytes,2,rep,name=Raw,proto3" json:"Raw,omitempty"`
	Type    string   `protobuf:"bytes,3,opt,name=Type,proto3" json:"Type,omitempty"`
	Scope   string   `protobuf:"bytes,4,opt,name=Scope,proto3" json:"Scope,omitempty"`
	Message string   `protobuf:"bytes,5,opt,name=Message,proto3" json:"Message,omitempty"`
	Change  *Change  `protobuf:"bytes,6,opt,name=Change,proto3" json:"Change,omitempty"`
}

func (x *Commit) Reset() {
	*x = Commit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_semrel_structs_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Commit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Commit) ProtoMessage() {}

func (x *Commit) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_semrel_structs_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Commit.ProtoReflect.Descriptor instead.
func (*Commit) Descriptor() ([]byte, []int) {
	return file_pkg_semrel_structs_proto_rawDescGZIP(), []int{2}
}

func (x *Commit) GetSHA() string {
	if x != nil {
		return x.SHA
	}
	return ""
}

func (x *Commit) GetRaw() []string {
	if x != nil {
		return x.Raw
	}
	return nil
}

func (x *Commit) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Commit) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *Commit) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Commit) GetChange() *Change {
	if x != nil {
		return x.Change
	}
	return nil
}

type Release struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SHA     string `protobuf:"bytes,1,opt,name=SHA,proto3" json:"SHA,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=Version,proto3" json:"Version,omitempty"`
}

func (x *Release) Reset() {
	*x = Release{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_semrel_structs_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Release) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Release) ProtoMessage() {}

func (x *Release) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_semrel_structs_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Release.ProtoReflect.Descriptor instead.
func (*Release) Descriptor() ([]byte, []int) {
	return file_pkg_semrel_structs_proto_rawDescGZIP(), []int{3}
}

func (x *Release) GetSHA() string {
	if x != nil {
		return x.SHA
	}
	return ""
}

func (x *Release) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

var File_pkg_semrel_structs_proto protoreflect.FileDescriptor

var file_pkg_semrel_structs_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x65, 0x6d, 0x72, 0x65, 0x6c, 0x2f, 0x73, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3d, 0x0a, 0x09, 0x52, 0x61,
	0x77, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x53, 0x48, 0x41, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x53, 0x48, 0x41, 0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x61, 0x77,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x52,
	0x61, 0x77, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x4a, 0x0a, 0x06, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x4d, 0x61, 0x6a, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x05, 0x4d, 0x61, 0x6a, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x4d, 0x69, 0x6e,
	0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x4d, 0x69, 0x6e, 0x6f, 0x72, 0x12,
	0x14, 0x0a, 0x05, 0x50, 0x61, 0x74, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05,
	0x50, 0x61, 0x74, 0x63, 0x68, 0x22, 0x91, 0x01, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x53, 0x48, 0x41, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x53,
	0x48, 0x41, 0x12, 0x10, 0x0a, 0x03, 0x52, 0x61, 0x77, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x03, 0x52, 0x61, 0x77, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x63, 0x6f, 0x70,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x0a, 0x06, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x52, 0x06, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x22, 0x35, 0x0a, 0x07, 0x52, 0x65, 0x6c,
	0x65, 0x61, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x53, 0x48, 0x41, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x53, 0x48, 0x41, 0x12, 0x18, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67,
	0x6f, 0x2d, 0x73, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63, 0x2d, 0x72, 0x65, 0x6c, 0x65, 0x61,
	0x73, 0x65, 0x2f, 0x73, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63, 0x2d, 0x72, 0x65, 0x6c, 0x65,
	0x61, 0x73, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x65, 0x6d, 0x72, 0x65, 0x6c, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_semrel_structs_proto_rawDescOnce sync.Once
	file_pkg_semrel_structs_proto_rawDescData = file_pkg_semrel_structs_proto_rawDesc
)

func file_pkg_semrel_structs_proto_rawDescGZIP() []byte {
	file_pkg_semrel_structs_proto_rawDescOnce.Do(func() {
		file_pkg_semrel_structs_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_semrel_structs_proto_rawDescData)
	})
	return file_pkg_semrel_structs_proto_rawDescData
}

var file_pkg_semrel_structs_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_semrel_structs_proto_goTypes = []interface{}{
	(*RawCommit)(nil), // 0: RawCommit
	(*Change)(nil),    // 1: Change
	(*Commit)(nil),    // 2: Commit
	(*Release)(nil),   // 3: Release
}
var file_pkg_semrel_structs_proto_depIdxs = []int32{
	1, // 0: Commit.Change:type_name -> Change
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_semrel_structs_proto_init() }
func file_pkg_semrel_structs_proto_init() {
	if File_pkg_semrel_structs_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_semrel_structs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RawCommit); i {
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
		file_pkg_semrel_structs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Change); i {
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
		file_pkg_semrel_structs_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Commit); i {
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
		file_pkg_semrel_structs_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Release); i {
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
			RawDescriptor: file_pkg_semrel_structs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_semrel_structs_proto_goTypes,
		DependencyIndexes: file_pkg_semrel_structs_proto_depIdxs,
		MessageInfos:      file_pkg_semrel_structs_proto_msgTypes,
	}.Build()
	File_pkg_semrel_structs_proto = out.File
	file_pkg_semrel_structs_proto_rawDesc = nil
	file_pkg_semrel_structs_proto_goTypes = nil
	file_pkg_semrel_structs_proto_depIdxs = nil
}
