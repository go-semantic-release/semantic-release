package analyzer

import (
	"context"

	"github.com/go-semantic-release/semantic-release/pkg/semrel"
)

const PluginNameCommitAnalyzer = "commit_analyzer"

type CommitAnalyzerServer struct {
	Impl CommitAnalyzer
	UnimplementedCommitAnalyzerPluginServer
}

func (c *CommitAnalyzerServer) Analyze(ctx context.Context, in *CommitAnalyzerRequest) (*CommitAnalyzerResponse, error) {
	return &CommitAnalyzerResponse{
		Commits: c.Impl.Analyze(in.RawCommits),
	}, nil
}

type CommitAnalyzerClient struct {
	Impl CommitAnalyzerPluginClient
}

func (c *CommitAnalyzerClient) Analyze(commits []*semrel.RawCommit) []*semrel.Commit {
	res, err := c.Impl.Analyze(context.Background(), &CommitAnalyzerRequest{RawCommits: commits})
	if err != nil {
		panic(err)
	}
	return res.Commits
}
