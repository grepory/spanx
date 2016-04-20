// Code generated by protoc-gen-gogo.
// source: opsee.proto
// DO NOT EDIT!

/*
Package opseeproto is a generated protocol buffer package.

It is generated from these files:
	opsee.proto

It has these top-level messages:
*/
package opseeproto

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.GoGoProtoPackageIsVersion1

var E_Graphql = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64040,
	Name:          "opseeproto.graphql",
	Tag:           "varint,64040,opt,name=graphql",
}

var E_Required = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         60000,
	Name:          "opseeproto.required",
	Tag:           "varint,60000,opt,name=required",
}

func init() {
	proto.RegisterExtension(E_Graphql)
	proto.RegisterExtension(E_Required)
}

var fileDescriptorOpsee = []byte{
	// 157 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x2f, 0x28, 0x4e,
	0x4d, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x02, 0x73, 0xc0, 0x6c, 0x29, 0x85, 0xf4,
	0xfc, 0xfc, 0xf4, 0x9c, 0x54, 0x7d, 0x30, 0x2f, 0xa9, 0x34, 0x4d, 0x3f, 0x25, 0xb5, 0x38, 0xb9,
	0x28, 0xb3, 0xa0, 0x24, 0xbf, 0x08, 0xa2, 0xda, 0xca, 0x82, 0x8b, 0x3d, 0xbd, 0x28, 0xb1, 0x20,
	0xa3, 0x30, 0x47, 0x48, 0x46, 0x0f, 0xa2, 0x5a, 0x0f, 0xa6, 0x5a, 0xcf, 0x2d, 0x33, 0x27, 0xd5,
	0xbf, 0xa0, 0x24, 0x33, 0x3f, 0xaf, 0x58, 0x62, 0xc5, 0x17, 0x66, 0x05, 0x46, 0x0d, 0x8e, 0x20,
	0x98, 0x72, 0x2b, 0x6b, 0x2e, 0x8e, 0xa2, 0xd4, 0xc2, 0xd2, 0xcc, 0xa2, 0xd4, 0x14, 0x21, 0x59,
	0x2c, 0x5a, 0x53, 0x73, 0x52, 0x60, 0x7a, 0x1f, 0x5c, 0x81, 0xe8, 0x85, 0x6b, 0x70, 0x52, 0xe2,
	0x12, 0x4e, 0xce, 0xcf, 0x45, 0xd7, 0xe5, 0xc4, 0xed, 0x0f, 0x72, 0x7b, 0x00, 0x88, 0x5b, 0x0c,
	0x08, 0x00, 0x00, 0xff, 0xff, 0x56, 0x50, 0xbe, 0x31, 0xd6, 0x00, 0x00, 0x00,
}
