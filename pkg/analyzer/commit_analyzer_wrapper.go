package analyzer

import (
	"context"

	"github.com/go-semantic-release/semantic-release/pkg/semrel"
)

const CommitAnalyzerPluginName = "commit_analyzer"

type CommitAnalyzerServer struct {
	Impl CommitAnalyzer
	UnimplementedCommitAnalyzerPluginServer
}

func (c *CommitAnalyzerServer) Analyze(ctx context.Context, request *AnalyzeCommits_Request) (*AnalyzeCommits_Response, error) {
	return &AnalyzeCommits_Response{
		Commits: c.Impl.Analyze(request.RawCommits),
	}, nil
}

type CommitAnalyzerClient struct {
	Impl CommitAnalyzerPluginClient
}

func (c *CommitAnalyzerClient) Analyze(commits []*semrel.RawCommit) []*semrel.Commit {
	res, err := c.Impl.Analyze(context.Background(), &AnalyzeCommits_Request{RawCommits: commits})
	if err != nil {
		panic(err)
	}
	return res.Commits
}
