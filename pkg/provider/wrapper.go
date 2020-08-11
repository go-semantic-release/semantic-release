package provider

import (
	"context"
	"errors"

	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
)

const PluginName = "provider"

type Server struct {
	Impl Provider
	UnimplementedProviderPluginServer
}

func (s *Server) Init(ctx context.Context, request *ProviderInit_Request) (*ProviderInit_Response, error) {
	err := s.Impl.Init(request.Config)
	res := &ProviderInit_Response{}
	if err != nil {
		res.Error = err.Error()
	}
	return res, nil
}

func (s *Server) Name(ctx context.Context, request *ProviderName_Request) (*ProviderName_Response, error) {
	return &ProviderName_Response{Name: s.Impl.Name()}, nil
}

func (s *Server) GetInfo(ctx context.Context, request *GetInfo_Request) (*GetInfo_Response, error) {
	info, err := s.Impl.GetInfo()
	if err != nil {
		return &GetInfo_Response{Error: err.Error()}, nil
	}
	return &GetInfo_Response{Info: info}, nil
}

func (s *Server) GetCommits(ctx context.Context, request *GetCommits_Request) (*GetCommits_Response, error) {
	commits, err := s.Impl.GetCommits(request.SHA)
	if err != nil {
		return &GetCommits_Response{Error: err.Error()}, nil
	}
	return &GetCommits_Response{RawCommits: commits}, nil
}

func (s *Server) GetReleases(ctx context.Context, request *GetReleases_Request) (*GetReleases_Response, error) {
	releases, err := s.Impl.GetReleases(request.Regexp)
	if err != nil {
		return &GetReleases_Response{Error: err.Error()}, nil
	}
	return &GetReleases_Response{Releases: releases}, nil
}

func (s *Server) CreateRelease(ctx context.Context, request *CreateRelease_Request) (*CreateRelease_Response, error) {
	err := s.Impl.CreateRelease(request.Config)
	res := &CreateRelease_Response{}
	if err != nil {
		res.Error = err.Error()
	}
	return res, nil
}

type Client struct {
	Impl ProviderPluginClient
}

func (c *Client) Init(m map[string]string) error {
	res, err := c.Impl.Init(context.Background(), &ProviderInit_Request{
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

func (c *Client) GetInfo() (*RepositoryInfo, error) {
	res, err := c.Impl.GetInfo(context.Background(), &GetInfo_Request{})
	if err != nil {
		return nil, err
	}
	if res.Error != "" {
		return nil, errors.New(res.Error)
	}
	return res.Info, nil
}

func (c *Client) GetCommits(sha string) ([]*semrel.RawCommit, error) {
	res, err := c.Impl.GetCommits(context.Background(), &GetCommits_Request{
		SHA: sha,
	})
	if err != nil {
		return nil, err
	}
	if res.Error != "" {
		return nil, errors.New(res.Error)
	}
	return res.RawCommits, nil
}

func (c *Client) GetReleases(re string) ([]*semrel.Release, error) {
	res, err := c.Impl.GetReleases(context.Background(), &GetReleases_Request{
		Regexp: re,
	})
	if err != nil {
		return nil, err
	}
	if res.Error != "" {
		return nil, errors.New(res.Error)
	}
	return res.Releases, nil
}

func (c *Client) CreateRelease(config *CreateReleaseConfig) error {
	res, err := c.Impl.CreateRelease(context.Background(), &CreateRelease_Request{
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

func (c *Client) Name() string {
	res, err := c.Impl.Name(context.Background(), &ProviderName_Request{})
	if err != nil {
		panic(err)
	}
	return res.Name
}
