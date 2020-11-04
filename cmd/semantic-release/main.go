package main

import (
	"fmt"
	"github.com/go-semantic-release/semantic-release/v2/pkg/runner"
	"log"
	"os"
	"strings"

	"github.com/go-semantic-release/semantic-release/v2/pkg/config"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/manager"
	"github.com/spf13/cobra"
)

// SRVERSION is the semantic-release version (added at compile time)
var SRVERSION string

var exitHandler func()

func errorHandler(logger *log.Logger) func(error, ...int) {
	return func(err error, exitCode ...int) {
		if err != nil {
			logger.Println(err)
			if exitHandler != nil {
				exitHandler()
			}
			if len(exitCode) == 1 {
				os.Exit(exitCode[0])
				return
			}
			os.Exit(1)
		}
	}
}

func main() {
	cmd := &cobra.Command{
		Use:     "semantic-release",
		Short:   "semantic-release - fully automated package/module/image publishing",
		Run:     cliHandler,
		Version: SRVERSION,
	}

	err := config.InitConfig(cmd)
	if err != nil {
		fmt.Printf("\nConfig error: %s\n", err.Error())
		os.Exit(1)
		return
	}
	err = cmd.Execute()
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(1)
	}
}

func newSemanticRelease(pluginManager *manager.PluginManager, logger *log.Logger) (*runner.SemanticRelease, error) {
	ci, err := pluginManager.GetCICondition()

	if err != nil {
		return nil, err
	}

	logger.Printf("CI-condition plugin: %s@%s\n", ci.Name(), ci.Version())

	prov, err := pluginManager.GetProvider()
	logger.Printf("provider plugin: %s@%s\n", prov.Name(), prov.Version())

	if err != nil {
		return nil, err
	}

	hooksExecutor, err := pluginManager.GetChainedHooksExecutor()

	if err != nil {
		return nil, err
	}

	hooksNames := hooksExecutor.GetNameVersionPairs()
	if len(hooksNames) > 0 {
		logger.Printf("hooks SemanticRelease: %s\n", strings.Join(hooksNames, ", "))
	}

	commitAnalyzer, err := pluginManager.GetCommitAnalyzer()

	if err != nil {
		return nil, err
	}

	logger.Printf("commit-analyzer plugin: %s@%s\n", commitAnalyzer.Name(), commitAnalyzer.Version())

	changelogGenerator, err := pluginManager.GetChangelogGenerator()

	if err != nil {
		return nil, err
	}

	logger.Printf("changelog-generator plugin: %s@%s\n", changelogGenerator.Name(), changelogGenerator.Version())

	updater, err := pluginManager.GetChainedUpdater()

	if err != nil {
		return nil, err
	}

	logger.Printf("files-updater SemanticRelease: %s\n", strings.Join(updater.GetNameVersionPairs(), ", "))

	return &runner.SemanticRelease{
		Prov:               prov,
		CI:                 ci,
		HooksExecutor:      hooksExecutor,
		CommitAnalyzer:     commitAnalyzer,
		ChangelogGenerator: changelogGenerator,
		Updater:            updater,
	}, nil
}

func cliHandler(cmd *cobra.Command, args []string) {
	logger := log.New(os.Stderr, "[go-semantic-release]: ", 0)
	exitIfError := errorHandler(logger)

	logger.Printf("version: %s\n", SRVERSION)

	conf, err := config.NewConfig(cmd)
	exitIfError(err)

	pluginManager, err := manager.New(conf)
	exitIfError(err)

	semantic, err := newSemanticRelease(pluginManager, logger)
	exitIfError(err)

	semantic.Run(conf)

	exitHandler = func() {
		pluginManager.Stop()
	}

}
