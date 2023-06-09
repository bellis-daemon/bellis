// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.22.2
// source: profile.proto

package profile

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UserProfile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email     string       `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	CreatedAt string       `protobuf:"bytes,2,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	IsVip     bool         `protobuf:"varint,3,opt,name=isVip,proto3" json:"isVip,omitempty"`
	Envoy     *EnvoyPolicy `protobuf:"bytes,4,opt,name=envoy,proto3" json:"envoy,omitempty"`
}

func (x *UserProfile) Reset() {
	*x = UserProfile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserProfile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserProfile) ProtoMessage() {}

func (x *UserProfile) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserProfile.ProtoReflect.Descriptor instead.
func (*UserProfile) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{0}
}

func (x *UserProfile) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UserProfile) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *UserProfile) GetIsVip() bool {
	if x != nil {
		return x.IsVip
	}
	return false
}

func (x *UserProfile) GetEnvoy() *EnvoyPolicy {
	if x != nil {
		return x.Envoy
	}
	return nil
}

type EnvoyPolicy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PolicyID      string              `protobuf:"bytes,1,opt,name=policyID,proto3" json:"policyID,omitempty"`
	PolicyType    int32               `protobuf:"varint,2,opt,name=policyType,proto3" json:"policyType,omitempty"`
	OfflineAlert  bool                `protobuf:"varint,3,opt,name=offlineAlert,proto3" json:"offlineAlert,omitempty"`
	PredictAlert  bool                `protobuf:"varint,4,opt,name=predictAlert,proto3" json:"predictAlert,omitempty"`
	PolicyContent *EnvoyPolicyContent `protobuf:"bytes,5,opt,name=policyContent,proto3" json:"policyContent,omitempty"`
}

func (x *EnvoyPolicy) Reset() {
	*x = EnvoyPolicy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvoyPolicy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvoyPolicy) ProtoMessage() {}

func (x *EnvoyPolicy) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvoyPolicy.ProtoReflect.Descriptor instead.
func (*EnvoyPolicy) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{1}
}

func (x *EnvoyPolicy) GetPolicyID() string {
	if x != nil {
		return x.PolicyID
	}
	return ""
}

func (x *EnvoyPolicy) GetPolicyType() int32 {
	if x != nil {
		return x.PolicyType
	}
	return 0
}

func (x *EnvoyPolicy) GetOfflineAlert() bool {
	if x != nil {
		return x.OfflineAlert
	}
	return false
}

func (x *EnvoyPolicy) GetPredictAlert() bool {
	if x != nil {
		return x.PredictAlert
	}
	return false
}

func (x *EnvoyPolicy) GetPolicyContent() *EnvoyPolicyContent {
	if x != nil {
		return x.PolicyContent
	}
	return nil
}

type NewPassword struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Password string `protobuf:"bytes,1,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *NewPassword) Reset() {
	*x = NewPassword{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewPassword) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewPassword) ProtoMessage() {}

func (x *NewPassword) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewPassword.ProtoReflect.Descriptor instead.
func (*NewPassword) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{2}
}

func (x *NewPassword) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type NewEmail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *NewEmail) Reset() {
	*x = NewEmail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewEmail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewEmail) ProtoMessage() {}

func (x *NewEmail) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewEmail.ProtoReflect.Descriptor instead.
func (*NewEmail) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{3}
}

func (x *NewEmail) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type Alert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OfflineAlert bool `protobuf:"varint,1,opt,name=OfflineAlert,proto3" json:"OfflineAlert,omitempty"`
	PredictAlert bool `protobuf:"varint,2,opt,name=PredictAlert,proto3" json:"PredictAlert,omitempty"`
}

func (x *Alert) Reset() {
	*x = Alert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Alert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Alert) ProtoMessage() {}

func (x *Alert) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Alert.ProtoReflect.Descriptor instead.
func (*Alert) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{4}
}

func (x *Alert) GetOfflineAlert() bool {
	if x != nil {
		return x.OfflineAlert
	}
	return false
}

func (x *Alert) GetPredictAlert() bool {
	if x != nil {
		return x.PredictAlert
	}
	return false
}

type Gotify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url   string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *Gotify) Reset() {
	*x = Gotify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Gotify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Gotify) ProtoMessage() {}

func (x *Gotify) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Gotify.ProtoReflect.Descriptor instead.
func (*Gotify) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{5}
}

func (x *Gotify) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Gotify) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type Email struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *Email) Reset() {
	*x = Email{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Email) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Email) ProtoMessage() {}

func (x *Email) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Email.ProtoReflect.Descriptor instead.
func (*Email) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{6}
}

func (x *Email) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type EnvoyPolicyContent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Content:
	//
	//	*EnvoyPolicyContent_Gotify
	//	*EnvoyPolicyContent_Email
	Content isEnvoyPolicyContent_Content `protobuf_oneof:"content"`
}

func (x *EnvoyPolicyContent) Reset() {
	*x = EnvoyPolicyContent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvoyPolicyContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvoyPolicyContent) ProtoMessage() {}

func (x *EnvoyPolicyContent) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvoyPolicyContent.ProtoReflect.Descriptor instead.
func (*EnvoyPolicyContent) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{7}
}

func (m *EnvoyPolicyContent) GetContent() isEnvoyPolicyContent_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (x *EnvoyPolicyContent) GetGotify() *Gotify {
	if x, ok := x.GetContent().(*EnvoyPolicyContent_Gotify); ok {
		return x.Gotify
	}
	return nil
}

func (x *EnvoyPolicyContent) GetEmail() *Email {
	if x, ok := x.GetContent().(*EnvoyPolicyContent_Email); ok {
		return x.Email
	}
	return nil
}

type isEnvoyPolicyContent_Content interface {
	isEnvoyPolicyContent_Content()
}

type EnvoyPolicyContent_Gotify struct {
	Gotify *Gotify `protobuf:"bytes,1,opt,name=gotify,proto3,oneof"`
}

type EnvoyPolicyContent_Email struct {
	Email *Email `protobuf:"bytes,2,opt,name=email,proto3,oneof"`
}

func (*EnvoyPolicyContent_Gotify) isEnvoyPolicyContent_Content() {}

func (*EnvoyPolicyContent_Email) isEnvoyPolicyContent_Content() {}

var File_profile_proto protoreflect.FileDescriptor

var file_profile_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x1d, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e,
	0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x99, 0x01, 0x0a, 0x0b,
	0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x69, 0x73, 0x56, 0x69, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05,
	0x69, 0x73, 0x56, 0x69, 0x70, 0x12, 0x40, 0x0a, 0x05, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x2e, 0x45, 0x6e, 0x76, 0x6f, 0x79, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x52, 0x05, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x22, 0xea, 0x01, 0x0a, 0x0b, 0x45, 0x6e, 0x76, 0x6f,
	0x79, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x6f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x6c,
	0x65, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x6f, 0x66, 0x66, 0x6c, 0x69,
	0x6e, 0x65, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x70, 0x72, 0x65, 0x64, 0x69,
	0x63, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x70,
	0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x12, 0x57, 0x0a, 0x0d, 0x70,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x31, 0x2e, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x62, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x2e, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x2e, 0x45, 0x6e, 0x76, 0x6f, 0x79, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x52, 0x0d, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x43, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x22, 0x29, 0x0a, 0x0b, 0x4e, 0x65, 0x77, 0x50, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22,
	0x20, 0x0a, 0x08, 0x4e, 0x65, 0x77, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x22, 0x4f, 0x0a, 0x05, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x4f, 0x66,
	0x66, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0c, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x12, 0x22,
	0x0a, 0x0c, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x41, 0x6c, 0x65,
	0x72, 0x74, 0x22, 0x30, 0x0a, 0x06, 0x47, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x21, 0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x18, 0x0a,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x9e, 0x01, 0x0a, 0x12, 0x45, 0x6e, 0x76, 0x6f,
	0x79, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x3f,
	0x0a, 0x06, 0x67, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25,
	0x2e, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e,
	0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x47,
	0x6f, 0x74, 0x69, 0x66, 0x79, 0x48, 0x00, 0x52, 0x06, 0x67, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x12,
	0x3c, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24,
	0x2e, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e,
	0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x45,
	0x6d, 0x61, 0x69, 0x6c, 0x48, 0x00, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x42, 0x09, 0x0a,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x32, 0x97, 0x04, 0x0a, 0x0e, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x54, 0x0a, 0x0e, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2a, 0x2e, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x62,
	0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x12, 0x54, 0x0a, 0x0e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x12, 0x2a, 0x2e, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x62, 0x61, 0x63,
	0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x2e, 0x4e, 0x65, 0x77, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x4e, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x27, 0x2e, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e,
	0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x4e, 0x65, 0x77, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x4b, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x12, 0x24, 0x2e, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e,
	0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x12, 0x5e, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x47, 0x6f, 0x74, 0x69, 0x66,
	0x79, 0x12, 0x25, 0x2e, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x62, 0x61, 0x63, 0x6b, 0x65,
	0x6e, 0x64, 0x2e, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x2e, 0x47, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x1a, 0x2a, 0x2e, 0x62, 0x65, 0x6c, 0x6c, 0x69,
	0x73, 0x2e, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x45, 0x6e, 0x76, 0x6f, 0x79, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x12, 0x5c, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x45, 0x6d, 0x61, 0x69, 0x6c,
	0x12, 0x24, 0x2e, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e,
	0x64, 0x2e, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x2e, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x1a, 0x2a, 0x2e, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2e,
	0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x45, 0x6e, 0x76, 0x6f, 0x79, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x62, 0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2d, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2f, 0x62,
	0x65, 0x6c, 0x6c, 0x69, 0x73, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x2f, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x6d, 0x6f, 0x62, 0x69, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_profile_proto_rawDescOnce sync.Once
	file_profile_proto_rawDescData = file_profile_proto_rawDesc
)

func file_profile_proto_rawDescGZIP() []byte {
	file_profile_proto_rawDescOnce.Do(func() {
		file_profile_proto_rawDescData = protoimpl.X.CompressGZIP(file_profile_proto_rawDescData)
	})
	return file_profile_proto_rawDescData
}

var file_profile_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_profile_proto_goTypes = []interface{}{
	(*UserProfile)(nil),        // 0: bellis.backend.mobile.profile.UserProfile
	(*EnvoyPolicy)(nil),        // 1: bellis.backend.mobile.profile.EnvoyPolicy
	(*NewPassword)(nil),        // 2: bellis.backend.mobile.profile.NewPassword
	(*NewEmail)(nil),           // 3: bellis.backend.mobile.profile.NewEmail
	(*Alert)(nil),              // 4: bellis.backend.mobile.profile.Alert
	(*Gotify)(nil),             // 5: bellis.backend.mobile.profile.Gotify
	(*Email)(nil),              // 6: bellis.backend.mobile.profile.Email
	(*EnvoyPolicyContent)(nil), // 7: bellis.backend.mobile.profile.EnvoyPolicyContent
	(*emptypb.Empty)(nil),      // 8: google.protobuf.Empty
}
var file_profile_proto_depIdxs = []int32{
	1,  // 0: bellis.backend.mobile.profile.UserProfile.envoy:type_name -> bellis.backend.mobile.profile.EnvoyPolicy
	7,  // 1: bellis.backend.mobile.profile.EnvoyPolicy.policyContent:type_name -> bellis.backend.mobile.profile.EnvoyPolicyContent
	5,  // 2: bellis.backend.mobile.profile.EnvoyPolicyContent.gotify:type_name -> bellis.backend.mobile.profile.Gotify
	6,  // 3: bellis.backend.mobile.profile.EnvoyPolicyContent.email:type_name -> bellis.backend.mobile.profile.Email
	8,  // 4: bellis.backend.mobile.profile.ProfileService.GetUserProfile:input_type -> google.protobuf.Empty
	2,  // 5: bellis.backend.mobile.profile.ProfileService.ChangePassword:input_type -> bellis.backend.mobile.profile.NewPassword
	3,  // 6: bellis.backend.mobile.profile.ProfileService.ChangeEmail:input_type -> bellis.backend.mobile.profile.NewEmail
	4,  // 7: bellis.backend.mobile.profile.ProfileService.ChangeAlert:input_type -> bellis.backend.mobile.profile.Alert
	5,  // 8: bellis.backend.mobile.profile.ProfileService.UseGotify:input_type -> bellis.backend.mobile.profile.Gotify
	6,  // 9: bellis.backend.mobile.profile.ProfileService.UseEmail:input_type -> bellis.backend.mobile.profile.Email
	0,  // 10: bellis.backend.mobile.profile.ProfileService.GetUserProfile:output_type -> bellis.backend.mobile.profile.UserProfile
	8,  // 11: bellis.backend.mobile.profile.ProfileService.ChangePassword:output_type -> google.protobuf.Empty
	8,  // 12: bellis.backend.mobile.profile.ProfileService.ChangeEmail:output_type -> google.protobuf.Empty
	8,  // 13: bellis.backend.mobile.profile.ProfileService.ChangeAlert:output_type -> google.protobuf.Empty
	1,  // 14: bellis.backend.mobile.profile.ProfileService.UseGotify:output_type -> bellis.backend.mobile.profile.EnvoyPolicy
	1,  // 15: bellis.backend.mobile.profile.ProfileService.UseEmail:output_type -> bellis.backend.mobile.profile.EnvoyPolicy
	10, // [10:16] is the sub-list for method output_type
	4,  // [4:10] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_profile_proto_init() }
func file_profile_proto_init() {
	if File_profile_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_profile_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserProfile); i {
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
		file_profile_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvoyPolicy); i {
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
		file_profile_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewPassword); i {
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
		file_profile_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewEmail); i {
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
		file_profile_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Alert); i {
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
		file_profile_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Gotify); i {
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
		file_profile_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Email); i {
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
		file_profile_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvoyPolicyContent); i {
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
	file_profile_proto_msgTypes[7].OneofWrappers = []interface{}{
		(*EnvoyPolicyContent_Gotify)(nil),
		(*EnvoyPolicyContent_Email)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_profile_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_profile_proto_goTypes,
		DependencyIndexes: file_profile_proto_depIdxs,
		MessageInfos:      file_profile_proto_msgTypes,
	}.Build()
	File_profile_proto = out.File
	file_profile_proto_rawDesc = nil
	file_profile_proto_goTypes = nil
	file_profile_proto_depIdxs = nil
}
