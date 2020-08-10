package generator

import "context"

const ChangelogGeneratorPluginName = "changelog_generator"

type ChangelogGeneratorServer struct {
	Impl ChangelogGenerator
	UnimplementedChangelogGeneratorPluginServer
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
