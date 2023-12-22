// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.2
// source: api/liveclass/liveclass.proto

package liveclass

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	LiveClass_LiveClassCommonCallBack_FullMethodName = "/api.liveclass.LiveClass/LiveClassCommonCallBack"
	LiveClass_LiveClassCallBack_FullMethodName       = "/api.liveclass.LiveClass/LiveClassCallBack"
	LiveClass_LiveClassDelay_FullMethodName          = "/api.liveclass.LiveClass/LiveClassDelay"
	LiveClass_LiveClassBegin_FullMethodName          = "/api.liveclass.LiveClass/LiveClassBegin"
	LiveClass_LiveClassEnd_FullMethodName            = "/api.liveclass.LiveClass/LiveClassEnd"
	LiveClass_InviteGroup_FullMethodName             = "/api.liveclass.LiveClass/InviteGroup"
	LiveClass_GetTeachingUser_FullMethodName         = "/api.liveclass.LiveClass/GetTeachingUser"
	LiveClass_VideoTranscodeStatus_FullMethodName    = "/api.liveclass.LiveClass/VideoTranscodeStatus"
	LiveClass_GenTrReport_FullMethodName             = "/api.liveclass.LiveClass/GenTrReport"
	LiveClass_LiveClassRelatedGroup_FullMethodName   = "/api.liveclass.LiveClass/LiveClassRelatedGroup"
)

// LiveClassClient is the client API for LiveClass service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LiveClassClient interface {
	// 实时课堂通用回调
	LiveClassCommonCallBack(ctx context.Context, in *LiveClassCommonCallBackRequest, opts ...grpc.CallOption) (*LiveClassCommonCallBackReply, error)
	// * 历史课堂回调
	LiveClassCallBack(ctx context.Context, in *LiveClassCallBackRequest, opts ...grpc.CallOption) (*LiveClassCallBackReply, error)
	// * 课堂延时
	LiveClassDelay(ctx context.Context, in *LiveClassDelayRequest, opts ...grpc.CallOption) (*LiveClassDelayReply, error)
	// * 课堂开始
	LiveClassBegin(ctx context.Context, in *LiveClassBeginRequest, opts ...grpc.CallOption) (*LiveClassBeginReply, error)
	// * 课堂结束
	LiveClassEnd(ctx context.Context, in *LiveClassEndRequest, opts ...grpc.CallOption) (*LiveClassEndReply, error)
	// 联播班级邀请
	InviteGroup(ctx context.Context, in *InviteGroupRequest, opts ...grpc.CallOption) (*InviteGroupReply, error)
	// 获取教研活动人员
	GetTeachingUser(ctx context.Context, in *GetTeachingUserRequest, opts ...grpc.CallOption) (*GetTeachingUserReply, error)
	// 视频转码成功回调
	VideoTranscodeStatus(ctx context.Context, in *VideoTranscodeStatusRequest, opts ...grpc.CallOption) (*VideoTranscodeStatusRequestReply, error)
	// 生成教研活动报告定时任务调用
	GenTrReport(ctx context.Context, in *GenTrReportRequest, opts ...grpc.CallOption) (*GenTrReportResponse, error)
	// 生成教研活动报告定时任务调用
	LiveClassRelatedGroup(ctx context.Context, in *LiveClassRelatedGroupRequest, opts ...grpc.CallOption) (*LiveClassRelatedGroupResponse, error)
}

type liveClassClient struct {
	cc grpc.ClientConnInterface
}

func NewLiveClassClient(cc grpc.ClientConnInterface) LiveClassClient {
	return &liveClassClient{cc}
}

func (c *liveClassClient) LiveClassCommonCallBack(ctx context.Context, in *LiveClassCommonCallBackRequest, opts ...grpc.CallOption) (*LiveClassCommonCallBackReply, error) {
	out := new(LiveClassCommonCallBackReply)
	err := c.cc.Invoke(ctx, LiveClass_LiveClassCommonCallBack_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClassClient) LiveClassCallBack(ctx context.Context, in *LiveClassCallBackRequest, opts ...grpc.CallOption) (*LiveClassCallBackReply, error) {
	out := new(LiveClassCallBackReply)
	err := c.cc.Invoke(ctx, LiveClass_LiveClassCallBack_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClassClient) LiveClassDelay(ctx context.Context, in *LiveClassDelayRequest, opts ...grpc.CallOption) (*LiveClassDelayReply, error) {
	out := new(LiveClassDelayReply)
	err := c.cc.Invoke(ctx, LiveClass_LiveClassDelay_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClassClient) LiveClassBegin(ctx context.Context, in *LiveClassBeginRequest, opts ...grpc.CallOption) (*LiveClassBeginReply, error) {
	out := new(LiveClassBeginReply)
	err := c.cc.Invoke(ctx, LiveClass_LiveClassBegin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClassClient) LiveClassEnd(ctx context.Context, in *LiveClassEndRequest, opts ...grpc.CallOption) (*LiveClassEndReply, error) {
	out := new(LiveClassEndReply)
	err := c.cc.Invoke(ctx, LiveClass_LiveClassEnd_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClassClient) InviteGroup(ctx context.Context, in *InviteGroupRequest, opts ...grpc.CallOption) (*InviteGroupReply, error) {
	out := new(InviteGroupReply)
	err := c.cc.Invoke(ctx, LiveClass_InviteGroup_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClassClient) GetTeachingUser(ctx context.Context, in *GetTeachingUserRequest, opts ...grpc.CallOption) (*GetTeachingUserReply, error) {
	out := new(GetTeachingUserReply)
	err := c.cc.Invoke(ctx, LiveClass_GetTeachingUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClassClient) VideoTranscodeStatus(ctx context.Context, in *VideoTranscodeStatusRequest, opts ...grpc.CallOption) (*VideoTranscodeStatusRequestReply, error) {
	out := new(VideoTranscodeStatusRequestReply)
	err := c.cc.Invoke(ctx, LiveClass_VideoTranscodeStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClassClient) GenTrReport(ctx context.Context, in *GenTrReportRequest, opts ...grpc.CallOption) (*GenTrReportResponse, error) {
	out := new(GenTrReportResponse)
	err := c.cc.Invoke(ctx, LiveClass_GenTrReport_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveClassClient) LiveClassRelatedGroup(ctx context.Context, in *LiveClassRelatedGroupRequest, opts ...grpc.CallOption) (*LiveClassRelatedGroupResponse, error) {
	out := new(LiveClassRelatedGroupResponse)
	err := c.cc.Invoke(ctx, LiveClass_LiveClassRelatedGroup_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LiveClassServer is the server API for LiveClass service.
// All implementations must embed UnimplementedLiveClassServer
// for forward compatibility
type LiveClassServer interface {
	// 实时课堂通用回调
	LiveClassCommonCallBack(context.Context, *LiveClassCommonCallBackRequest) (*LiveClassCommonCallBackReply, error)
	// * 历史课堂回调
	LiveClassCallBack(context.Context, *LiveClassCallBackRequest) (*LiveClassCallBackReply, error)
	// * 课堂延时
	LiveClassDelay(context.Context, *LiveClassDelayRequest) (*LiveClassDelayReply, error)
	// * 课堂开始
	LiveClassBegin(context.Context, *LiveClassBeginRequest) (*LiveClassBeginReply, error)
	// * 课堂结束
	LiveClassEnd(context.Context, *LiveClassEndRequest) (*LiveClassEndReply, error)
	// 联播班级邀请
	InviteGroup(context.Context, *InviteGroupRequest) (*InviteGroupReply, error)
	// 获取教研活动人员
	GetTeachingUser(context.Context, *GetTeachingUserRequest) (*GetTeachingUserReply, error)
	// 视频转码成功回调
	VideoTranscodeStatus(context.Context, *VideoTranscodeStatusRequest) (*VideoTranscodeStatusRequestReply, error)
	// 生成教研活动报告定时任务调用
	GenTrReport(context.Context, *GenTrReportRequest) (*GenTrReportResponse, error)
	// 生成教研活动报告定时任务调用
	LiveClassRelatedGroup(context.Context, *LiveClassRelatedGroupRequest) (*LiveClassRelatedGroupResponse, error)
	mustEmbedUnimplementedLiveClassServer()
}

// UnimplementedLiveClassServer must be embedded to have forward compatible implementations.
type UnimplementedLiveClassServer struct {
}

func (UnimplementedLiveClassServer) LiveClassCommonCallBack(context.Context, *LiveClassCommonCallBackRequest) (*LiveClassCommonCallBackReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LiveClassCommonCallBack not implemented")
}
func (UnimplementedLiveClassServer) LiveClassCallBack(context.Context, *LiveClassCallBackRequest) (*LiveClassCallBackReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LiveClassCallBack not implemented")
}
func (UnimplementedLiveClassServer) LiveClassDelay(context.Context, *LiveClassDelayRequest) (*LiveClassDelayReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LiveClassDelay not implemented")
}
func (UnimplementedLiveClassServer) LiveClassBegin(context.Context, *LiveClassBeginRequest) (*LiveClassBeginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LiveClassBegin not implemented")
}
func (UnimplementedLiveClassServer) LiveClassEnd(context.Context, *LiveClassEndRequest) (*LiveClassEndReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LiveClassEnd not implemented")
}
func (UnimplementedLiveClassServer) InviteGroup(context.Context, *InviteGroupRequest) (*InviteGroupReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InviteGroup not implemented")
}
func (UnimplementedLiveClassServer) GetTeachingUser(context.Context, *GetTeachingUserRequest) (*GetTeachingUserReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTeachingUser not implemented")
}
func (UnimplementedLiveClassServer) VideoTranscodeStatus(context.Context, *VideoTranscodeStatusRequest) (*VideoTranscodeStatusRequestReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VideoTranscodeStatus not implemented")
}
func (UnimplementedLiveClassServer) GenTrReport(context.Context, *GenTrReportRequest) (*GenTrReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenTrReport not implemented")
}
func (UnimplementedLiveClassServer) LiveClassRelatedGroup(context.Context, *LiveClassRelatedGroupRequest) (*LiveClassRelatedGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LiveClassRelatedGroup not implemented")
}
func (UnimplementedLiveClassServer) mustEmbedUnimplementedLiveClassServer() {}

// UnsafeLiveClassServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LiveClassServer will
// result in compilation errors.
type UnsafeLiveClassServer interface {
	mustEmbedUnimplementedLiveClassServer()
}

func RegisterLiveClassServer(s grpc.ServiceRegistrar, srv LiveClassServer) {
	s.RegisterService(&LiveClass_ServiceDesc, srv)
}

func _LiveClass_LiveClassCommonCallBack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LiveClassCommonCallBackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveClassServer).LiveClassCommonCallBack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LiveClass_LiveClassCommonCallBack_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveClassServer).LiveClassCommonCallBack(ctx, req.(*LiveClassCommonCallBackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LiveClass_LiveClassCallBack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LiveClassCallBackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveClassServer).LiveClassCallBack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LiveClass_LiveClassCallBack_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveClassServer).LiveClassCallBack(ctx, req.(*LiveClassCallBackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LiveClass_LiveClassDelay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LiveClassDelayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveClassServer).LiveClassDelay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LiveClass_LiveClassDelay_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveClassServer).LiveClassDelay(ctx, req.(*LiveClassDelayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LiveClass_LiveClassBegin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LiveClassBeginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveClassServer).LiveClassBegin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LiveClass_LiveClassBegin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveClassServer).LiveClassBegin(ctx, req.(*LiveClassBeginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LiveClass_LiveClassEnd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LiveClassEndRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveClassServer).LiveClassEnd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LiveClass_LiveClassEnd_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveClassServer).LiveClassEnd(ctx, req.(*LiveClassEndRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LiveClass_InviteGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InviteGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveClassServer).InviteGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LiveClass_InviteGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveClassServer).InviteGroup(ctx, req.(*InviteGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LiveClass_GetTeachingUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTeachingUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveClassServer).GetTeachingUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LiveClass_GetTeachingUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveClassServer).GetTeachingUser(ctx, req.(*GetTeachingUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LiveClass_VideoTranscodeStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VideoTranscodeStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveClassServer).VideoTranscodeStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LiveClass_VideoTranscodeStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveClassServer).VideoTranscodeStatus(ctx, req.(*VideoTranscodeStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LiveClass_GenTrReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenTrReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveClassServer).GenTrReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LiveClass_GenTrReport_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveClassServer).GenTrReport(ctx, req.(*GenTrReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LiveClass_LiveClassRelatedGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LiveClassRelatedGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveClassServer).LiveClassRelatedGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LiveClass_LiveClassRelatedGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveClassServer).LiveClassRelatedGroup(ctx, req.(*LiveClassRelatedGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LiveClass_ServiceDesc is the grpc.ServiceDesc for LiveClass service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LiveClass_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.liveclass.LiveClass",
	HandlerType: (*LiveClassServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LiveClassCommonCallBack",
			Handler:    _LiveClass_LiveClassCommonCallBack_Handler,
		},
		{
			MethodName: "LiveClassCallBack",
			Handler:    _LiveClass_LiveClassCallBack_Handler,
		},
		{
			MethodName: "LiveClassDelay",
			Handler:    _LiveClass_LiveClassDelay_Handler,
		},
		{
			MethodName: "LiveClassBegin",
			Handler:    _LiveClass_LiveClassBegin_Handler,
		},
		{
			MethodName: "LiveClassEnd",
			Handler:    _LiveClass_LiveClassEnd_Handler,
		},
		{
			MethodName: "InviteGroup",
			Handler:    _LiveClass_InviteGroup_Handler,
		},
		{
			MethodName: "GetTeachingUser",
			Handler:    _LiveClass_GetTeachingUser_Handler,
		},
		{
			MethodName: "VideoTranscodeStatus",
			Handler:    _LiveClass_VideoTranscodeStatus_Handler,
		},
		{
			MethodName: "GenTrReport",
			Handler:    _LiveClass_GenTrReport_Handler,
		},
		{
			MethodName: "LiveClassRelatedGroup",
			Handler:    _LiveClass_LiveClassRelatedGroup_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/liveclass/liveclass.proto",
}
