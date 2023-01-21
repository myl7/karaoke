// Copyright (C) 2023 myl7
// SPDX-License-Identifier: Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: server.proto

package rpc

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

// Pre-encoding format of msgs transmitted between servers
type Onion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body     []byte `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
	NextHop  string `protobuf:"bytes,2,opt,name=nextHop,proto3" json:"nextHop,omitempty"`
	DeadDrop string `protobuf:"bytes,3,opt,name=deadDrop,proto3" json:"deadDrop,omitempty"`
}

func (x *Onion) Reset() {
	*x = Onion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Onion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Onion) ProtoMessage() {}

func (x *Onion) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Onion.ProtoReflect.Descriptor instead.
func (*Onion) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{0}
}

func (x *Onion) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *Onion) GetNextHop() string {
	if x != nil {
		return x.NextHop
	}
	return ""
}

func (x *Onion) GetDeadDrop() string {
	if x != nil {
		return x.DeadDrop
	}
	return ""
}

type OnionMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Msg:
	//
	//	*OnionMsg_Body
	//	*OnionMsg_Meta
	Msg isOnionMsg_Msg `protobuf_oneof:"msg"`
}

func (x *OnionMsg) Reset() {
	*x = OnionMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnionMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnionMsg) ProtoMessage() {}

func (x *OnionMsg) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnionMsg.ProtoReflect.Descriptor instead.
func (*OnionMsg) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{1}
}

func (m *OnionMsg) GetMsg() isOnionMsg_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (x *OnionMsg) GetBody() []byte {
	if x, ok := x.GetMsg().(*OnionMsg_Body); ok {
		return x.Body
	}
	return nil
}

func (x *OnionMsg) GetMeta() *OnionMsgMeta {
	if x, ok := x.GetMsg().(*OnionMsg_Meta); ok {
		return x.Meta
	}
	return nil
}

type isOnionMsg_Msg interface {
	isOnionMsg_Msg()
}

type OnionMsg_Body struct {
	Body []byte `protobuf:"bytes,1,opt,name=body,proto3,oneof"`
}

type OnionMsg_Meta struct {
	Meta *OnionMsgMeta `protobuf:"bytes,2,opt,name=meta,proto3,oneof"`
}

func (*OnionMsg_Body) isOnionMsg_Msg() {}

func (*OnionMsg_Meta) isOnionMsg_Msg() {}

type OnionMsgMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Leave empty to indicate the onion is from a client
	// TODO: Signature
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *OnionMsgMeta) Reset() {
	*x = OnionMsgMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnionMsgMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnionMsgMeta) ProtoMessage() {}

func (x *OnionMsgMeta) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnionMsgMeta.ProtoReflect.Descriptor instead.
func (*OnionMsgMeta) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{2}
}

func (x *OnionMsgMeta) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type FwdOnionsRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FwdOnionsRes) Reset() {
	*x = FwdOnionsRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FwdOnionsRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FwdOnionsRes) ProtoMessage() {}

func (x *FwdOnionsRes) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FwdOnionsRes.ProtoReflect.Descriptor instead.
func (*FwdOnionsRes) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{3}
}

type CheckBloomReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bloom []byte `protobuf:"bytes,1,opt,name=bloom,proto3" json:"bloom,omitempty"`
	// TODO: Signature
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CheckBloomReq) Reset() {
	*x = CheckBloomReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckBloomReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckBloomReq) ProtoMessage() {}

func (x *CheckBloomReq) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckBloomReq.ProtoReflect.Descriptor instead.
func (*CheckBloomReq) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{4}
}

func (x *CheckBloomReq) GetBloom() []byte {
	if x != nil {
		return x.Bloom
	}
	return nil
}

func (x *CheckBloomReq) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CheckBloomRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *CheckBloomRes) Reset() {
	*x = CheckBloomRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckBloomRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckBloomRes) ProtoMessage() {}

func (x *CheckBloomRes) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckBloomRes.ProtoReflect.Descriptor instead.
func (*CheckBloomRes) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{5}
}

func (x *CheckBloomRes) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

var File_server_proto protoreflect.FileDescriptor

var file_server_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03,
	0x72, 0x70, 0x63, 0x22, 0x51, 0x0a, 0x05, 0x4f, 0x6e, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x62, 0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x78, 0x74, 0x48, 0x6f, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6e, 0x65, 0x78, 0x74, 0x48, 0x6f, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x65,
	0x61, 0x64, 0x44, 0x72, 0x6f, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x65,
	0x61, 0x64, 0x44, 0x72, 0x6f, 0x70, 0x22, 0x50, 0x0a, 0x08, 0x4f, 0x6e, 0x69, 0x6f, 0x6e, 0x4d,
	0x73, 0x67, 0x12, 0x14, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x48, 0x00, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x27, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x4f, 0x6e, 0x69,
	0x6f, 0x6e, 0x4d, 0x73, 0x67, 0x4d, 0x65, 0x74, 0x61, 0x48, 0x00, 0x52, 0x04, 0x6d, 0x65, 0x74,
	0x61, 0x42, 0x05, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x1e, 0x0a, 0x0c, 0x4f, 0x6e, 0x69, 0x6f,
	0x6e, 0x4d, 0x73, 0x67, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x0e, 0x0a, 0x0c, 0x46, 0x77, 0x64, 0x4f,
	0x6e, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x22, 0x35, 0x0a, 0x0d, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x6c, 0x6f,
	0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x62, 0x6c, 0x6f, 0x6f, 0x6d, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x1f, 0x0a, 0x0d, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x73,
	0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b,
	0x32, 0x72, 0x0a, 0x09, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x50, 0x43, 0x12, 0x2f, 0x0a,
	0x09, 0x46, 0x77, 0x64, 0x4f, 0x6e, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x0d, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x4f, 0x6e, 0x69, 0x6f, 0x6e, 0x4d, 0x73, 0x67, 0x1a, 0x11, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x46, 0x77, 0x64, 0x4f, 0x6e, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x28, 0x01, 0x12, 0x34,
	0x0a, 0x0a, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x12, 0x12, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x71,
	0x1a, 0x12, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x42, 0x6c, 0x6f, 0x6f,
	0x6d, 0x52, 0x65, 0x73, 0x42, 0x21, 0x5a, 0x1f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6d, 0x79, 0x6c, 0x37, 0x2f, 0x6b, 0x61, 0x72, 0x61, 0x6f, 0x6b, 0x65, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_server_proto_rawDescOnce sync.Once
	file_server_proto_rawDescData = file_server_proto_rawDesc
)

func file_server_proto_rawDescGZIP() []byte {
	file_server_proto_rawDescOnce.Do(func() {
		file_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_proto_rawDescData)
	})
	return file_server_proto_rawDescData
}

var file_server_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_server_proto_goTypes = []interface{}{
	(*Onion)(nil),         // 0: rpc.Onion
	(*OnionMsg)(nil),      // 1: rpc.OnionMsg
	(*OnionMsgMeta)(nil),  // 2: rpc.OnionMsgMeta
	(*FwdOnionsRes)(nil),  // 3: rpc.FwdOnionsRes
	(*CheckBloomReq)(nil), // 4: rpc.CheckBloomReq
	(*CheckBloomRes)(nil), // 5: rpc.CheckBloomRes
}
var file_server_proto_depIdxs = []int32{
	2, // 0: rpc.OnionMsg.meta:type_name -> rpc.OnionMsgMeta
	1, // 1: rpc.ServerRPC.FwdOnions:input_type -> rpc.OnionMsg
	4, // 2: rpc.ServerRPC.CheckBloom:input_type -> rpc.CheckBloomReq
	3, // 3: rpc.ServerRPC.FwdOnions:output_type -> rpc.FwdOnionsRes
	5, // 4: rpc.ServerRPC.CheckBloom:output_type -> rpc.CheckBloomRes
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_server_proto_init() }
func file_server_proto_init() {
	if File_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Onion); i {
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
		file_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnionMsg); i {
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
		file_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnionMsgMeta); i {
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
		file_server_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FwdOnionsRes); i {
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
		file_server_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckBloomReq); i {
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
		file_server_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckBloomRes); i {
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
	file_server_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*OnionMsg_Body)(nil),
		(*OnionMsg_Meta)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_server_proto_goTypes,
		DependencyIndexes: file_server_proto_depIdxs,
		MessageInfos:      file_server_proto_msgTypes,
	}.Build()
	File_server_proto = out.File
	file_server_proto_rawDesc = nil
	file_server_proto_goTypes = nil
	file_server_proto_depIdxs = nil
}
