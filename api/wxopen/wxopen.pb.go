// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.23.2
// source: api/wxopen/wxopen.proto

package wxopen

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type GetAccessTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppId     string `protobuf:"bytes,1,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"`
	AppSecret string `protobuf:"bytes,2,opt,name=app_secret,json=appSecret,proto3" json:"app_secret,omitempty"`
}

func (x *GetAccessTokenRequest) Reset() {
	*x = GetAccessTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_wxopen_wxopen_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAccessTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccessTokenRequest) ProtoMessage() {}

func (x *GetAccessTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_wxopen_wxopen_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccessTokenRequest.ProtoReflect.Descriptor instead.
func (*GetAccessTokenRequest) Descriptor() ([]byte, []int) {
	return file_api_wxopen_wxopen_proto_rawDescGZIP(), []int{0}
}

func (x *GetAccessTokenRequest) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

func (x *GetAccessTokenRequest) GetAppSecret() string {
	if x != nil {
		return x.AppSecret
	}
	return ""
}

type GetAccessTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	ExpiresIn   int64  `protobuf:"varint,2,opt,name=expires_in,json=expiresIn,proto3" json:"expires_in,omitempty"`
}

func (x *GetAccessTokenResponse) Reset() {
	*x = GetAccessTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_wxopen_wxopen_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAccessTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccessTokenResponse) ProtoMessage() {}

func (x *GetAccessTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_wxopen_wxopen_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccessTokenResponse.ProtoReflect.Descriptor instead.
func (*GetAccessTokenResponse) Descriptor() ([]byte, []int) {
	return file_api_wxopen_wxopen_proto_rawDescGZIP(), []int{1}
}

func (x *GetAccessTokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *GetAccessTokenResponse) GetExpiresIn() int64 {
	if x != nil {
		return x.ExpiresIn
	}
	return 0
}

type LoginQrCodeCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExpireSeconds int64  `protobuf:"varint,1,opt,name=expire_seconds,json=expireSeconds,proto3" json:"expire_seconds,omitempty"`
	SceneStr      string `protobuf:"bytes,2,opt,name=scene_str,json=sceneStr,proto3" json:"scene_str,omitempty"`
}

func (x *LoginQrCodeCreateRequest) Reset() {
	*x = LoginQrCodeCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_wxopen_wxopen_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginQrCodeCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginQrCodeCreateRequest) ProtoMessage() {}

func (x *LoginQrCodeCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_wxopen_wxopen_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginQrCodeCreateRequest.ProtoReflect.Descriptor instead.
func (*LoginQrCodeCreateRequest) Descriptor() ([]byte, []int) {
	return file_api_wxopen_wxopen_proto_rawDescGZIP(), []int{2}
}

func (x *LoginQrCodeCreateRequest) GetExpireSeconds() int64 {
	if x != nil {
		return x.ExpireSeconds
	}
	return 0
}

func (x *LoginQrCodeCreateRequest) GetSceneStr() string {
	if x != nil {
		return x.SceneStr
	}
	return ""
}

type LoginQrCodeCreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ticket        string `protobuf:"bytes,1,opt,name=ticket,proto3" json:"ticket,omitempty"`
	ExpireSeconds int64  `protobuf:"varint,2,opt,name=expire_seconds,json=expireSeconds,proto3" json:"expire_seconds,omitempty"`
	Url           string `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *LoginQrCodeCreateResponse) Reset() {
	*x = LoginQrCodeCreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_wxopen_wxopen_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginQrCodeCreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginQrCodeCreateResponse) ProtoMessage() {}

func (x *LoginQrCodeCreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_wxopen_wxopen_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginQrCodeCreateResponse.ProtoReflect.Descriptor instead.
func (*LoginQrCodeCreateResponse) Descriptor() ([]byte, []int) {
	return file_api_wxopen_wxopen_proto_rawDescGZIP(), []int{3}
}

func (x *LoginQrCodeCreateResponse) GetTicket() string {
	if x != nil {
		return x.Ticket
	}
	return ""
}

func (x *LoginQrCodeCreateResponse) GetExpireSeconds() int64 {
	if x != nil {
		return x.ExpireSeconds
	}
	return 0
}

func (x *LoginQrCodeCreateResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type GetWxUserInfoByCodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *GetWxUserInfoByCodeRequest) Reset() {
	*x = GetWxUserInfoByCodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_wxopen_wxopen_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWxUserInfoByCodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWxUserInfoByCodeRequest) ProtoMessage() {}

func (x *GetWxUserInfoByCodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_wxopen_wxopen_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWxUserInfoByCodeRequest.ProtoReflect.Descriptor instead.
func (*GetWxUserInfoByCodeRequest) Descriptor() ([]byte, []int) {
	return file_api_wxopen_wxopen_proto_rawDescGZIP(), []int{4}
}

func (x *GetWxUserInfoByCodeRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type GetWxUserInfoByCodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Openid     string `protobuf:"bytes,1,opt,name=openid,proto3" json:"openid,omitempty"`
	Nickname   string `protobuf:"bytes,2,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Headimgurl string `protobuf:"bytes,3,opt,name=headimgurl,proto3" json:"headimgurl,omitempty"`
}

func (x *GetWxUserInfoByCodeResponse) Reset() {
	*x = GetWxUserInfoByCodeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_wxopen_wxopen_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWxUserInfoByCodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWxUserInfoByCodeResponse) ProtoMessage() {}

func (x *GetWxUserInfoByCodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_wxopen_wxopen_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWxUserInfoByCodeResponse.ProtoReflect.Descriptor instead.
func (*GetWxUserInfoByCodeResponse) Descriptor() ([]byte, []int) {
	return file_api_wxopen_wxopen_proto_rawDescGZIP(), []int{5}
}

func (x *GetWxUserInfoByCodeResponse) GetOpenid() string {
	if x != nil {
		return x.Openid
	}
	return ""
}

func (x *GetWxUserInfoByCodeResponse) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *GetWxUserInfoByCodeResponse) GetHeadimgurl() string {
	if x != nil {
		return x.Headimgurl
	}
	return ""
}

type LoginUrlCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueryString string `protobuf:"bytes,1,opt,name=queryString,proto3" json:"queryString,omitempty"`
}

func (x *LoginUrlCreateRequest) Reset() {
	*x = LoginUrlCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_wxopen_wxopen_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginUrlCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginUrlCreateRequest) ProtoMessage() {}

func (x *LoginUrlCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_wxopen_wxopen_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginUrlCreateRequest.ProtoReflect.Descriptor instead.
func (*LoginUrlCreateRequest) Descriptor() ([]byte, []int) {
	return file_api_wxopen_wxopen_proto_rawDescGZIP(), []int{6}
}

func (x *LoginUrlCreateRequest) GetQueryString() string {
	if x != nil {
		return x.QueryString
	}
	return ""
}

type LoginUrlCreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url   string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	State string `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
	AppId string `protobuf:"bytes,3,opt,name=appId,proto3" json:"appId,omitempty"`
}

func (x *LoginUrlCreateResponse) Reset() {
	*x = LoginUrlCreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_wxopen_wxopen_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginUrlCreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginUrlCreateResponse) ProtoMessage() {}

func (x *LoginUrlCreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_wxopen_wxopen_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginUrlCreateResponse.ProtoReflect.Descriptor instead.
func (*LoginUrlCreateResponse) Descriptor() ([]byte, []int) {
	return file_api_wxopen_wxopen_proto_rawDescGZIP(), []int{7}
}

func (x *LoginUrlCreateResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *LoginUrlCreateResponse) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *LoginUrlCreateResponse) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

var File_api_wxopen_wxopen_proto protoreflect.FileDescriptor

var file_api_wxopen_wxopen_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x77, 0x78, 0x6f, 0x70, 0x65, 0x6e, 0x2f, 0x77, 0x78, 0x6f,
	0x70, 0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x61, 0x70, 0x69, 0x2e, 0x77,
	0x78, 0x6f, 0x70, 0x65, 0x6e, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x4d, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06,
	0x61, 0x70, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x70,
	0x70, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x70, 0x70, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x70, 0x70, 0x53, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x22, 0x5a, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x1d, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x69, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x49, 0x6e, 0x22, 0x5e,
	0x0a, 0x18, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x51, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x78,
	0x70, 0x69, 0x72, 0x65, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0d, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x73, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x63, 0x65, 0x6e, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x63, 0x65, 0x6e, 0x65, 0x53, 0x74, 0x72, 0x22, 0x6c,
	0x0a, 0x19, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x51, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74,
	0x69, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x69, 0x63,
	0x6b, 0x65, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x5f, 0x73, 0x65,
	0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x65, 0x78, 0x70,
	0x69, 0x72, 0x65, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72,
	0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x30, 0x0a, 0x1a,
	0x47, 0x65, 0x74, 0x57, 0x78, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x79, 0x43,
	0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x71,
	0x0a, 0x1b, 0x47, 0x65, 0x74, 0x57, 0x78, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x42,
	0x79, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x6f, 0x70, 0x65, 0x6e, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f,
	0x70, 0x65, 0x6e, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x68, 0x65, 0x61, 0x64, 0x69, 0x6d, 0x67, 0x75, 0x72, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x68, 0x65, 0x61, 0x64, 0x69, 0x6d, 0x67, 0x75, 0x72,
	0x6c, 0x22, 0x39, 0x0a, 0x15, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x71, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0x56, 0x0a, 0x16,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61,
	0x70, 0x70, 0x49, 0x64, 0x32, 0xfd, 0x03, 0x0a, 0x06, 0x57, 0x78, 0x6f, 0x70, 0x65, 0x6e, 0x12,
	0x73, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x21, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x77, 0x78, 0x6f, 0x70, 0x65, 0x6e, 0x2e, 0x47,
	0x65, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x77, 0x78, 0x6f, 0x70, 0x65,
	0x6e, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14,
	0x22, 0x0f, 0x2f, 0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x3a, 0x01, 0x2a, 0x12, 0x7f, 0x0a, 0x11, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x51, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x24, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x77, 0x78, 0x6f, 0x70, 0x65, 0x6e, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x51, 0x72, 0x43, 0x6f,
	0x64, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x25, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x77, 0x78, 0x6f, 0x70, 0x65, 0x6e, 0x2e, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x51, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x22, 0x12,
	0x2f, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x51, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x87, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x57, 0x78, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x26, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x77, 0x78, 0x6f, 0x70, 0x65, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x78,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x77, 0x78, 0x6f, 0x70,
	0x65, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x78, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x42, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1f,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x22, 0x14, 0x2f, 0x47, 0x65, 0x74, 0x57, 0x78, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x3a, 0x01, 0x2a, 0x12,
	0x73, 0x0a, 0x0e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x12, 0x21, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x77, 0x78, 0x6f, 0x70, 0x65, 0x6e, 0x2e, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x77, 0x78, 0x6f, 0x70, 0x65,
	0x6e, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14,
	0x22, 0x0f, 0x2f, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x3a, 0x01, 0x2a, 0x42, 0x21, 0x0a, 0x0a, 0x61, 0x70, 0x69, 0x2e, 0x77, 0x78, 0x6f, 0x70,
	0x65, 0x6e, 0x50, 0x01, 0x5a, 0x11, 0x61, 0x70, 0x69, 0x2f, 0x77, 0x78, 0x6f, 0x70, 0x65, 0x6e,
	0x3b, 0x77, 0x78, 0x6f, 0x70, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_wxopen_wxopen_proto_rawDescOnce sync.Once
	file_api_wxopen_wxopen_proto_rawDescData = file_api_wxopen_wxopen_proto_rawDesc
)

func file_api_wxopen_wxopen_proto_rawDescGZIP() []byte {
	file_api_wxopen_wxopen_proto_rawDescOnce.Do(func() {
		file_api_wxopen_wxopen_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_wxopen_wxopen_proto_rawDescData)
	})
	return file_api_wxopen_wxopen_proto_rawDescData
}

var file_api_wxopen_wxopen_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_wxopen_wxopen_proto_goTypes = []interface{}{
	(*GetAccessTokenRequest)(nil),       // 0: api.wxopen.GetAccessTokenRequest
	(*GetAccessTokenResponse)(nil),      // 1: api.wxopen.GetAccessTokenResponse
	(*LoginQrCodeCreateRequest)(nil),    // 2: api.wxopen.LoginQrCodeCreateRequest
	(*LoginQrCodeCreateResponse)(nil),   // 3: api.wxopen.LoginQrCodeCreateResponse
	(*GetWxUserInfoByCodeRequest)(nil),  // 4: api.wxopen.GetWxUserInfoByCodeRequest
	(*GetWxUserInfoByCodeResponse)(nil), // 5: api.wxopen.GetWxUserInfoByCodeResponse
	(*LoginUrlCreateRequest)(nil),       // 6: api.wxopen.LoginUrlCreateRequest
	(*LoginUrlCreateResponse)(nil),      // 7: api.wxopen.LoginUrlCreateResponse
}
var file_api_wxopen_wxopen_proto_depIdxs = []int32{
	0, // 0: api.wxopen.Wxopen.GetAccessToken:input_type -> api.wxopen.GetAccessTokenRequest
	2, // 1: api.wxopen.Wxopen.LoginQrCodeCreate:input_type -> api.wxopen.LoginQrCodeCreateRequest
	4, // 2: api.wxopen.Wxopen.GetWxUserInfoByCode:input_type -> api.wxopen.GetWxUserInfoByCodeRequest
	6, // 3: api.wxopen.Wxopen.LoginUrlCreate:input_type -> api.wxopen.LoginUrlCreateRequest
	1, // 4: api.wxopen.Wxopen.GetAccessToken:output_type -> api.wxopen.GetAccessTokenResponse
	3, // 5: api.wxopen.Wxopen.LoginQrCodeCreate:output_type -> api.wxopen.LoginQrCodeCreateResponse
	5, // 6: api.wxopen.Wxopen.GetWxUserInfoByCode:output_type -> api.wxopen.GetWxUserInfoByCodeResponse
	7, // 7: api.wxopen.Wxopen.LoginUrlCreate:output_type -> api.wxopen.LoginUrlCreateResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_wxopen_wxopen_proto_init() }
func file_api_wxopen_wxopen_proto_init() {
	if File_api_wxopen_wxopen_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_wxopen_wxopen_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAccessTokenRequest); i {
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
		file_api_wxopen_wxopen_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAccessTokenResponse); i {
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
		file_api_wxopen_wxopen_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginQrCodeCreateRequest); i {
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
		file_api_wxopen_wxopen_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginQrCodeCreateResponse); i {
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
		file_api_wxopen_wxopen_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWxUserInfoByCodeRequest); i {
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
		file_api_wxopen_wxopen_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWxUserInfoByCodeResponse); i {
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
		file_api_wxopen_wxopen_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginUrlCreateRequest); i {
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
		file_api_wxopen_wxopen_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginUrlCreateResponse); i {
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
			RawDescriptor: file_api_wxopen_wxopen_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_wxopen_wxopen_proto_goTypes,
		DependencyIndexes: file_api_wxopen_wxopen_proto_depIdxs,
		MessageInfos:      file_api_wxopen_wxopen_proto_msgTypes,
	}.Build()
	File_api_wxopen_wxopen_proto = out.File
	file_api_wxopen_wxopen_proto_rawDesc = nil
	file_api_wxopen_wxopen_proto_goTypes = nil
	file_api_wxopen_wxopen_proto_depIdxs = nil
}
