// Copyright 2020 by cgs.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.1
// source: piano.proto

package rpcpiano

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	rpcmsgType "hcc/piano/action/grpc/pb/rpcmsgType"
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

// Symbols defined in public import of msgType.proto.

type Empty = rpcmsgType.Empty
type HccError = rpcmsgType.HccError
type Node = rpcmsgType.Node
type NodeDetail = rpcmsgType.NodeDetail
type Server = rpcmsgType.Server
type ServerNode = rpcmsgType.ServerNode
type Quota = rpcmsgType.Quota
type VNC = rpcmsgType.VNC
type Volume = rpcmsgType.Volume
type VolumeAttachment = rpcmsgType.VolumeAttachment
type Pool = rpcmsgType.Pool
type AdaptiveIPSetting = rpcmsgType.AdaptiveIPSetting
type AdaptiveIPAvailableIPList = rpcmsgType.AdaptiveIPAvailableIPList
type AdaptiveIPServer = rpcmsgType.AdaptiveIPServer
type Subnet = rpcmsgType.Subnet
type MetricInfo = rpcmsgType.MetricInfo
type MonitoringData = rpcmsgType.MonitoringData
type NormalAction = rpcmsgType.NormalAction
type HccAction = rpcmsgType.HccAction
type Action = rpcmsgType.Action
type Control = rpcmsgType.Control
type Controls = rpcmsgType.Controls
type ScheduledNodes = rpcmsgType.ScheduledNodes
type ScheduleServer = rpcmsgType.ScheduleServer
type ServerAction = rpcmsgType.ServerAction

type ReqMetricInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MetricInfo *rpcmsgType.MetricInfo `protobuf:"bytes,1,opt,name=metricInfo,proto3" json:"metricInfo,omitempty"`
}

func (x *ReqMetricInfo) Reset() {
	*x = ReqMetricInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_piano_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqMetricInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqMetricInfo) ProtoMessage() {}

func (x *ReqMetricInfo) ProtoReflect() protoreflect.Message {
	mi := &file_piano_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqMetricInfo.ProtoReflect.Descriptor instead.
func (*ReqMetricInfo) Descriptor() ([]byte, []int) {
	return file_piano_proto_rawDescGZIP(), []int{0}
}

func (x *ReqMetricInfo) GetMetricInfo() *rpcmsgType.MetricInfo {
	if x != nil {
		return x.MetricInfo
	}
	return nil
}

type ResMonitoringData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MonitoringData *rpcmsgType.MonitoringData `protobuf:"bytes,1,opt,name=monitoringData,proto3" json:"monitoringData,omitempty"`
	HccErrorStack  []*rpcmsgType.HccError     `protobuf:"bytes,2,rep,name=hccErrorStack,proto3" json:"hccErrorStack,omitempty"`
}

func (x *ResMonitoringData) Reset() {
	*x = ResMonitoringData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_piano_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResMonitoringData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResMonitoringData) ProtoMessage() {}

func (x *ResMonitoringData) ProtoReflect() protoreflect.Message {
	mi := &file_piano_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResMonitoringData.ProtoReflect.Descriptor instead.
func (*ResMonitoringData) Descriptor() ([]byte, []int) {
	return file_piano_proto_rawDescGZIP(), []int{1}
}

func (x *ResMonitoringData) GetMonitoringData() *rpcmsgType.MonitoringData {
	if x != nil {
		return x.MonitoringData
	}
	return nil
}

func (x *ResMonitoringData) GetHccErrorStack() []*rpcmsgType.HccError {
	if x != nil {
		return x.HccErrorStack
	}
	return nil
}

var File_piano_proto protoreflect.FileDescriptor

var file_piano_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x70, 0x69, 0x61, 0x6e, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x52,
	0x70, 0x63, 0x50, 0x69, 0x61, 0x6e, 0x6f, 0x1a, 0x0d, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x44, 0x0a, 0x0d, 0x52, 0x65, 0x71, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x33, 0x0a, 0x0a, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x4d, 0x73,
	0x67, 0x54, 0x79, 0x70, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x0a, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x8d, 0x01, 0x0a,
	0x11, 0x52, 0x65, 0x73, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x3f, 0x0a, 0x0e, 0x6d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67,
	0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x4d, 0x73, 0x67,
	0x54, 0x79, 0x70, 0x65, 0x2e, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x0e, 0x6d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x44,
	0x61, 0x74, 0x61, 0x12, 0x37, 0x0a, 0x0d, 0x68, 0x63, 0x63, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x53,
	0x74, 0x61, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x4d, 0x73, 0x67,
	0x54, 0x79, 0x70, 0x65, 0x2e, 0x48, 0x63, 0x63, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x0d, 0x68,
	0x63, 0x63, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x63, 0x6b, 0x32, 0x4c, 0x0a, 0x05,
	0x50, 0x69, 0x61, 0x6e, 0x6f, 0x12, 0x43, 0x0a, 0x09, 0x54, 0x65, 0x6c, 0x65, 0x67, 0x72, 0x61,
	0x70, 0x68, 0x12, 0x17, 0x2e, 0x52, 0x70, 0x63, 0x50, 0x69, 0x61, 0x6e, 0x6f, 0x2e, 0x52, 0x65,
	0x71, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x1b, 0x2e, 0x52, 0x70,
	0x63, 0x50, 0x69, 0x61, 0x6e, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f,
	0x72, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x74, 0x61, 0x22, 0x00, 0x42, 0x23, 0x5a, 0x21, 0x68, 0x63,
	0x63, 0x2f, 0x70, 0x69, 0x61, 0x6e, 0x6f, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x72, 0x70, 0x63, 0x70, 0x69, 0x61, 0x6e, 0x6f, 0x50,
	0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_piano_proto_rawDescOnce sync.Once
	file_piano_proto_rawDescData = file_piano_proto_rawDesc
)

func file_piano_proto_rawDescGZIP() []byte {
	file_piano_proto_rawDescOnce.Do(func() {
		file_piano_proto_rawDescData = protoimpl.X.CompressGZIP(file_piano_proto_rawDescData)
	})
	return file_piano_proto_rawDescData
}

var file_piano_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_piano_proto_goTypes = []interface{}{
	(*ReqMetricInfo)(nil),             // 0: RpcPiano.ReqMetricInfo
	(*ResMonitoringData)(nil),         // 1: RpcPiano.ResMonitoringData
	(*rpcmsgType.MetricInfo)(nil),     // 2: MsgType.MetricInfo
	(*rpcmsgType.MonitoringData)(nil), // 3: MsgType.MonitoringData
	(*rpcmsgType.HccError)(nil),       // 4: MsgType.HccError
}
var file_piano_proto_depIdxs = []int32{
	2, // 0: RpcPiano.ReqMetricInfo.metricInfo:type_name -> MsgType.MetricInfo
	3, // 1: RpcPiano.ResMonitoringData.monitoringData:type_name -> MsgType.MonitoringData
	4, // 2: RpcPiano.ResMonitoringData.hccErrorStack:type_name -> MsgType.HccError
	0, // 3: RpcPiano.Piano.Telegraph:input_type -> RpcPiano.ReqMetricInfo
	1, // 4: RpcPiano.Piano.Telegraph:output_type -> RpcPiano.ResMonitoringData
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_piano_proto_init() }
func file_piano_proto_init() {
	if File_piano_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_piano_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqMetricInfo); i {
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
		file_piano_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResMonitoringData); i {
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
			RawDescriptor: file_piano_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_piano_proto_goTypes,
		DependencyIndexes: file_piano_proto_depIdxs,
		MessageInfos:      file_piano_proto_msgTypes,
	}.Build()
	File_piano_proto = out.File
	file_piano_proto_rawDesc = nil
	file_piano_proto_goTypes = nil
	file_piano_proto_depIdxs = nil
}
