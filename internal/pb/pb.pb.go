// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.0
// source: .proto

package pb

import (
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

type OriginalURL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *OriginalURL) Reset() {
	*x = OriginalURL{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OriginalURL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OriginalURL) ProtoMessage() {}

func (x *OriginalURL) ProtoReflect() protoreflect.Message {
	mi := &file___proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OriginalURL.ProtoReflect.Descriptor instead.
func (*OriginalURL) Descriptor() ([]byte, []int) {
	return file___proto_rawDescGZIP(), []int{0}
}

func (x *OriginalURL) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type ShortURL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *ShortURL) Reset() {
	*x = ShortURL{}
	if protoimpl.UnsafeEnabled {
		mi := &file___proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortURL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortURL) ProtoMessage() {}

func (x *ShortURL) ProtoReflect() protoreflect.Message {
	mi := &file___proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortURL.ProtoReflect.Descriptor instead.
func (*ShortURL) Descriptor() ([]byte, []int) {
	return file___proto_rawDescGZIP(), []int{1}
}

func (x *ShortURL) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File___proto protoreflect.FileDescriptor

var file___proto_rawDesc = []byte{
	0x0a, 0x06, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x75,
	0x72, 0x6c, 0x22, 0x1f, 0x0a, 0x0b, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52,
	0x4c, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x72, 0x6c, 0x22, 0x1c, 0x0a, 0x08, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x32, 0x8f, 0x01, 0x0a, 0x0f, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3d, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x12, 0x15, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x75,
	0x72, 0x6c, 0x2e, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x1a, 0x12,
	0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x75, 0x72, 0x6c, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55,
	0x52, 0x4c, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x12, 0x12, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x75, 0x72,
	0x6c, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x1a, 0x15, 0x2e, 0x73, 0x68, 0x6f,
	0x72, 0x74, 0x75, 0x72, 0x6c, 0x2e, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52,
	0x4c, 0x22, 0x00, 0x42, 0x10, 0x5a, 0x0e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x70, 0x62, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file___proto_rawDescOnce sync.Once
	file___proto_rawDescData = file___proto_rawDesc
)

func file___proto_rawDescGZIP() []byte {
	file___proto_rawDescOnce.Do(func() {
		file___proto_rawDescData = protoimpl.X.CompressGZIP(file___proto_rawDescData)
	})
	return file___proto_rawDescData
}

var file___proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file___proto_goTypes = []interface{}{
	(*OriginalURL)(nil), // 0: shorturl.OriginalURL
	(*ShortURL)(nil),    // 1: shorturl.ShortURL
}
var file___proto_depIdxs = []int32{
	0, // 0: shorturl.ShortURLService.CreateShortURL:input_type -> shorturl.OriginalURL
	1, // 1: shorturl.ShortURLService.GetOriginalURL:input_type -> shorturl.ShortURL
	1, // 2: shorturl.ShortURLService.CreateShortURL:output_type -> shorturl.ShortURL
	0, // 3: shorturl.ShortURLService.GetOriginalURL:output_type -> shorturl.OriginalURL
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file___proto_init() }
func file___proto_init() {
	if File___proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file___proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OriginalURL); i {
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
		file___proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShortURL); i {
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
			RawDescriptor: file___proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file___proto_goTypes,
		DependencyIndexes: file___proto_depIdxs,
		MessageInfos:      file___proto_msgTypes,
	}.Build()
	File___proto = out.File
	file___proto_rawDesc = nil
	file___proto_goTypes = nil
	file___proto_depIdxs = nil
}
