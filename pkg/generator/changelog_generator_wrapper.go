package generator

import (
	"context"
	"errors"
)

const ChangelogGeneratorPluginName = "changelog_generator"

type ChangelogGeneratorServer struct {
	Impl ChangelogGenerator
	UnimplementedChangelogGeneratorPluginServer
}

func (c *ChangelogGeneratorServer) Init(ctx context.Context, request *ChangelogGeneratorInit_Request) (*ChangelogGeneratorInit_Response, error) {
	err := c.Impl.Init(request.Config)
	res := &ChangelogGeneratorInit_Response{}
	if err != nil {
		res.Error = err.Error()
	}
	return res, nil
}

func (c *ChangelogGeneratorServer) Name(ctx context.Context, request *ChangelogGeneratorName_Request) (*ChangelogGeneratorName_Response, error) {
	return &ChangelogGeneratorName_Response{Name: c.Impl.Name()}, nil
}

func (c *ChangelogGeneratorServer) Version(ctx context.Context, request *ChangelogGeneratorVersion_Request) (*ChangelogGeneratorVersion_Response, error) {
	return &ChangelogGeneratorVersion_Response{Version: c.Impl.Version()}, nil
}

func (c *ChangelogGeneratorServer) Generate(ctx context.Context, request *GenerateChangelog_Request) (*GenerateChangelog_Response, error) {
	return &GenerateChangelog_Response{
		Changelog: c.Impl.Generate(request.Config),
	}, nil
}

type ChangelogGeneratorClient struct {
	Impl ChangelogGeneratorPluginClient
}

func (c *ChangelogGeneratorClient) Generate(config *ChangelogGeneratorConfig) string {
	res, err := c.Impl.Generate(context.Background(), &GenerateChangelog_Request{Config: config})
	if err != nil {
		panic(err)
	}
	return res.Changelog
}

func (c *ChangelogGeneratorClient) Init(m map[string]string) error {
	res, err := c.Impl.Init(context.Background(), &ChangelogGeneratorInit_Request{
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

func (c *ChangelogGeneratorClient) Name() string {
	res, err := c.Impl.Name(context.Background(), &ChangelogGeneratorName_Request{})
	if err != nil {
		panic(err)
	}
	return res.Name
}

func (c *ChangelogGeneratorClient) Version() string {
	res, err := c.Impl.Version(context.Background(), &ChangelogGeneratorVersion_Request{})
	if err != nil {
		panic(err)
	}
	return res.Version
}
