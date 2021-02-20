// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/socialite/socialite.proto

package socialite

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Socialites service

type SocialitesService interface {
	// 小程序获取授权
	Auth(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	// rpc Register(Request) returns (Response) {} // 授权后注册【可用于增加新账号】
	// 授权网址
	AuthURL(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	// 绑定用户
	BuildUser(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
}

type socialitesService struct {
	c    client.Client
	name string
}

func NewSocialitesService(name string, c client.Client) SocialitesService {
	return &socialitesService{
		c:    c,
		name: name,
	}
}

func (c *socialitesService) Auth(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Socialites.Auth", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialitesService) AuthURL(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Socialites.AuthURL", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialitesService) BuildUser(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Socialites.BuildUser", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Socialites service

type SocialitesHandler interface {
	// 小程序获取授权
	Auth(context.Context, *Request, *Response) error
	// rpc Register(Request) returns (Response) {} // 授权后注册【可用于增加新账号】
	// 授权网址
	AuthURL(context.Context, *Request, *Response) error
	// 绑定用户
	BuildUser(context.Context, *Request, *Response) error
}

func RegisterSocialitesHandler(s server.Server, hdlr SocialitesHandler, opts ...server.HandlerOption) error {
	type socialites interface {
		Auth(ctx context.Context, in *Request, out *Response) error
		AuthURL(ctx context.Context, in *Request, out *Response) error
		BuildUser(ctx context.Context, in *Request, out *Response) error
	}
	type Socialites struct {
		socialites
	}
	h := &socialitesHandler{hdlr}
	return s.Handle(s.NewHandler(&Socialites{h}, opts...))
}

type socialitesHandler struct {
	SocialitesHandler
}

func (h *socialitesHandler) Auth(ctx context.Context, in *Request, out *Response) error {
	return h.SocialitesHandler.Auth(ctx, in, out)
}

func (h *socialitesHandler) AuthURL(ctx context.Context, in *Request, out *Response) error {
	return h.SocialitesHandler.AuthURL(ctx, in, out)
}

func (h *socialitesHandler) BuildUser(ctx context.Context, in *Request, out *Response) error {
	return h.SocialitesHandler.BuildUser(ctx, in, out)
}
