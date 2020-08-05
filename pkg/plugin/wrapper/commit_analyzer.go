package wrapper

import (
	"context"

	"github.com/go-semantic-release/semantic-release/pkg/analyzer/commit"
	"github.com/go-semantic-release/semantic-release/pkg/semrel"
)

const TypeCommitAnalyzer = "commit_analyzer"

type CommitAnalyzerServer struct {
	Impl commit.Analyzer
	commit.UnimplementedCommitAnalyzerPluginServer
}

func (c *CommitAnalyzerServer) Analyze(ctx context.Context, in *commit.CommitAnalyzerRequest) (*commit.CommitAnalyzerResponse, error) {
	return &commit.CommitAnalyzerResponse{
		Commits: c.Impl.Analyze(in.RawCommits),
	}, nil
}

type CommitAnalyzerClient struct {
	Impl commit.CommitAnalyzerPluginClient
}

func (c *CommitAnalyzerClient) Analyze(commits []*semrel.RawCommit) []*semrel.Commit {
	res, err := c.Impl.Analyze(context.Background(), &commit.CommitAnalyzerRequest{RawCommits: commits})
	if err != nil {
		panic(err)
	}
	return res.Commits
}
