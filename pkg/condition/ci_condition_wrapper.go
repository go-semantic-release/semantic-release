package condition

import (
	"context"
	"errors"
)

const CIConditionPluginName = "ci_condition"

type CIConditionServer struct {
	Impl CICondition
	UnimplementedCIConditionPluginServer
}

func (c *CIConditionServer) Name(ctx context.Context, request *CIName_Request) (*CIName_Response, error) {
	return &CIName_Response{Name: c.Impl.Name()}, nil
}

func (c *CIConditionServer) Version(ctx context.Context, request *CIVersion_Request) (*CIVersion_Response, error) {
	return &CIVersion_Response{Version: c.Impl.Version()}, nil
}

func (c *CIConditionServer) RunCondition(ctx context.Context, request *RunCondition_Request) (*RunCondition_Response, error) {
	err := c.Impl.RunCondition(request.Value)
	ret := &RunCondition_Response{}
	if err != nil {
		ret.Error = err.Error()
	}
	return ret, nil
}

func (c *CIConditionServer) GetCurrentBranch(ctx context.Context, request *GetCurrentBranch_Request) (*GetCurrentBranch_Response, error) {
	return &GetCurrentBranch_Response{Value: c.Impl.GetCurrentBranch()}, nil
}

func (c *CIConditionServer) GetCurrentSHA(ctx context.Context, request *GetCurrentSHA_Request) (*GetCurrentSHA_Response, error) {
	return &GetCurrentSHA_Response{Value: c.Impl.GetCurrentSHA()}, nil
}

type CIConditionClient struct {
	Impl CIConditionPluginClient
}

func (c *CIConditionClient) Name() string {
	res, err := c.Impl.Name(context.Background(), &CIName_Request{})
	if err != nil {
		panic(err)
	}
	return res.Name
}

func (c *CIConditionClient) Version() string {
	res, err := c.Impl.Version(context.Background(), &CIVersion_Request{})
	if err != nil {
		panic(err)
	}
	return res.Version
}

func (c *CIConditionClient) RunCondition(m map[string]string) error {
	res, err := c.Impl.RunCondition(context.Background(), &RunCondition_Request{Value: m})
	if err != nil {
		return err
	}
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

func (c *CIConditionClient) GetCurrentBranch() string {
	res, err := c.Impl.GetCurrentBranch(context.Background(), &GetCurrentBranch_Request{})
	if err != nil {
		panic(err)
	}
	return res.Value
}

func (c *CIConditionClient) GetCurrentSHA() string {
	res, err := c.Impl.GetCurrentSHA(context.Background(), &GetCurrentSHA_Request{})
	if err != nil {
		panic(err)
	}
	return res.Value
}
