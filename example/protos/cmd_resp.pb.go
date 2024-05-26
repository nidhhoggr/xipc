// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.0
// source: example/protobuf/protos/cmd_resp.proto

package protos

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

type CmdResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	ValueStr string `protobuf:"bytes,3,opt,name=value_str,json=valueStr,proto3" json:"value_str,omitempty"`
	Error    string `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *CmdResp) Reset() {
	*x = CmdResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_protobuf_protos_cmd_resp_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CmdResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CmdResp) ProtoMessage() {}

func (x *CmdResp) ProtoReflect() protoreflect.Message {
	mi := &file_example_protobuf_protos_cmd_resp_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CmdResp.ProtoReflect.Descriptor instead.
func (*CmdResp) Descriptor() ([]byte, []int) {
	return file_example_protobuf_protos_cmd_resp_proto_rawDescGZIP(), []int{0}
}

func (x *CmdResp) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CmdResp) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CmdResp) GetValueStr() string {
	if x != nil {
		return x.ValueStr
	}
	return ""
}

func (x *CmdResp) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_example_protobuf_protos_cmd_resp_proto protoreflect.FileDescriptor

var file_example_protobuf_protos_cmd_resp_proto_rawDesc = []byte{
	0x0a, 0x26, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x63, 0x6d, 0x64, 0x5f, 0x72, 0x65,
	0x73, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x63, 0x6d, 0x64, 0x5f, 0x72, 0x65,
	0x73, 0x70, 0x22, 0x60, 0x0a, 0x07, 0x43, 0x6d, 0x64, 0x52, 0x65, 0x73, 0x70, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x74, 0x72, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x42, 0x19, 0x5a, 0x17, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_example_protobuf_protos_cmd_resp_proto_rawDescOnce sync.Once
	file_example_protobuf_protos_cmd_resp_proto_rawDescData = file_example_protobuf_protos_cmd_resp_proto_rawDesc
)

func file_example_protobuf_protos_cmd_resp_proto_rawDescGZIP() []byte {
	file_example_protobuf_protos_cmd_resp_proto_rawDescOnce.Do(func() {
		file_example_protobuf_protos_cmd_resp_proto_rawDescData = protoimpl.X.CompressGZIP(file_example_protobuf_protos_cmd_resp_proto_rawDescData)
	})
	return file_example_protobuf_protos_cmd_resp_proto_rawDescData
}

var file_example_protobuf_protos_cmd_resp_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_example_protobuf_protos_cmd_resp_proto_goTypes = []interface{}{
	(*CmdResp)(nil), // 0: cmd_resp.CmdResp
}
var file_example_protobuf_protos_cmd_resp_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_example_protobuf_protos_cmd_resp_proto_init() }
func file_example_protobuf_protos_cmd_resp_proto_init() {
	if File_example_protobuf_protos_cmd_resp_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_example_protobuf_protos_cmd_resp_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CmdResp); i {
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
			RawDescriptor: file_example_protobuf_protos_cmd_resp_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_example_protobuf_protos_cmd_resp_proto_goTypes,
		DependencyIndexes: file_example_protobuf_protos_cmd_resp_proto_depIdxs,
		MessageInfos:      file_example_protobuf_protos_cmd_resp_proto_msgTypes,
	}.Build()
	File_example_protobuf_protos_cmd_resp_proto = out.File
	file_example_protobuf_protos_cmd_resp_proto_rawDesc = nil
	file_example_protobuf_protos_cmd_resp_proto_goTypes = nil
	file_example_protobuf_protos_cmd_resp_proto_depIdxs = nil
}
