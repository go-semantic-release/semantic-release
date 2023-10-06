package hooks

import (
	"fmt"
)

type Hooks interface {
	Name() string
	Version() string
	Success(*SuccessHookConfig) error
	NoRelease(*NoReleaseConfig) error
}

type ChainedHooksExecutor struct {
	HooksChain []Hooks
}

func (c *ChainedHooksExecutor) Success(config *SuccessHookConfig) error {
	for _, h := range c.HooksChain {
		name := h.Name()
		err := h.Success(config)
		if err != nil {
			return fmt.Errorf("%s hook has failed: %w", name, err)
		}
	}
	return nil
}

func (c *ChainedHooksExecutor) NoRelease(config *NoReleaseConfig) error {
	for _, h := range c.HooksChain {
		name := h.Name()
		err := h.NoRelease(config)
		if err != nil {
			return fmt.Errorf("%s hook has failed: %w", name, err)
		}
	}
	return nil
}

func (c *ChainedHooksExecutor) Init(hooksConf map[string]map[string]interface{}) error {
	for _, h := range c.HooksChain {
		var hookClient *Client = h.(*Client)
		conf := hooksConf[h.Name()]
		err := hookClient.Init(conf)
		if err != nil {
			return fmt.Errorf("[hook %s]: %v", h.Name(), err)
		}
	}
	return nil
}

func (c *ChainedHooksExecutor) GetNameVersionPairs() []string {
	ret := make([]string, len(c.HooksChain))
	for i, h := range c.HooksChain {
		ret[i] = fmt.Sprintf("%s@%s", h.Name(), h.Version())
	}
	return ret
}
