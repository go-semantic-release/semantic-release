package analyzer

import (
	"context"
	"errors"

	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
)

const CommitAnalyzerPluginName = "commit_analyzer"

type CommitAnalyzerServer struct {
	Impl CommitAnalyzer
	UnimplementedCommitAnalyzerPluginServer
}

func (c *CommitAnalyzerServer) Init(_ context.Context, request *CommitAnalyzerInit_Request) (*CommitAnalyzerInit_Response, error) {
	err := c.Impl.Init(request.Config)
	res := &CommitAnalyzerInit_Response{}
	if err != nil {
		res.Error = err.Error()
	}
	return res, nil
}

func (c *CommitAnalyzerServer) Name(_ context.Context, _ *CommitAnalyzerName_Request) (*CommitAnalyzerName_Response, error) {
	return &CommitAnalyzerName_Response{Name: c.Impl.Name()}, nil
}

func (c *CommitAnalyzerServer) Version(_ context.Context, _ *CommitAnalyzerVersion_Request) (*CommitAnalyzerVersion_Response, error) {
	return &CommitAnalyzerVersion_Response{Version: c.Impl.Version()}, nil
}

func (c *CommitAnalyzerServer) Analyze(_ context.Context, request *AnalyzeCommits_Request) (*AnalyzeCommits_Response, error) {
	return &AnalyzeCommits_Response{
		Commits: c.Impl.Analyze(request.RawCommits),
	}, nil
}

type CommitAnalyzerClient struct {
	Impl CommitAnalyzerPluginClient
}

func (c *CommitAnalyzerClient) Init(m map[string]string) error {
	res, err := c.Impl.Init(context.Background(), &CommitAnalyzerInit_Request{
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

func (c *CommitAnalyzerClient) Name() string {
	res, err := c.Impl.Name(context.Background(), &CommitAnalyzerName_Request{})
	if err != nil {
		panic(err)
	}
	return res.Name
}

func (c *CommitAnalyzerClient) Version() string {
	res, err := c.Impl.Version(context.Background(), &CommitAnalyzerVersion_Request{})
	if err != nil {
		panic(err)
	}
	return res.Version
}

func (c *CommitAnalyzerClient) Analyze(commits []*semrel.RawCommit) []*semrel.Commit {
	res, err := c.Impl.Analyze(context.Background(), &AnalyzeCommits_Request{RawCommits: commits})
	if err != nil {
		panic(err)
	}
	return res.Commits
}
