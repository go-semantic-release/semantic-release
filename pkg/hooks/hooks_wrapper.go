package hooks

import (
	"context"
	"errors"
)

const PluginName = "hooks"

type Server struct {
	Impl Hooks
	UnimplementedHooksPluginServer
}

func (h *Server) Init(_ context.Context, request *HooksInit_Request) (*HooksInit_Response, error) {
	err := h.Impl.Init(request.Config)
	res := &HooksInit_Response{}
	if err != nil {
		res.Error = err.Error()
	}
	return res, nil
}

func (h *Server) Name(_ context.Context, _ *HooksName_Request) (*HooksName_Response, error) {
	return &HooksName_Response{Name: h.Impl.Name()}, nil
}

func (h *Server) Version(_ context.Context, _ *HooksVersion_Request) (*HooksVersion_Response, error) {
	return &HooksVersion_Response{Version: h.Impl.Version()}, nil
}

func (h *Server) Success(_ context.Context, request *SuccessHook_Request) (*SuccessHook_Response, error) {
	err := h.Impl.Success(request.Config)
	res := &SuccessHook_Response{}
	if err != nil {
		res.Error = err.Error()
	}
	return res, nil
}

func (h *Server) NoRelease(_ context.Context, request *NoReleaseHook_Request) (*NoReleaseHook_Response, error) {
	err := h.Impl.NoRelease(request.Config)
	res := &NoReleaseHook_Response{}
	if err != nil {
		res.Error = err.Error()
	}
	return res, nil
}

type Client struct {
	Impl HooksPluginClient
}

func (h *Client) Init(m map[string]string) error {
	res, err := h.Impl.Init(context.Background(), &HooksInit_Request{
		Config: m,
	})
	if err != nil {
		return err
	}
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

func (h *Client) Name() string {
	res, err := h.Impl.Name(context.Background(), &HooksName_Request{})
	if err != nil {
		panic(err)
	}
	return res.Name
}

func (h *Client) Version() string {
	res, err := h.Impl.Version(context.Background(), &HooksVersion_Request{})
	if err != nil {
		panic(err)
	}
	return res.Version
}

func (h *Client) Success(config *SuccessHookConfig) error {
	res, err := h.Impl.Success(context.Background(), &SuccessHook_Request{
		Config: config,
	})
	if err != nil {
		return err
	}
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

func (h *Client) NoRelease(config *NoReleaseConfig) error {
	res, err := h.Impl.NoRelease(context.Background(), &NoReleaseHook_Request{
		Config: config,
	})
	if err != nil {
		return err
	}
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}
