package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"github.com/jbcpollak/strcase"
	"github.com/Masterminds/semver/v3"
	"github.com/go-semantic-release/semantic-release/v2/pkg/config"
	"github.com/go-semantic-release/semantic-release/v2/pkg/generator"
	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin/manager"
	"github.com/go-semantic-release/semantic-release/v2/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
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

	config.SetFlags(cmd)
	cobra.OnInitialize(func() {
		err := config.InitConfig(cmd)
		if err != nil {
			fmt.Printf("\nConfig error: %s\n", err.Error())
			os.Exit(1)
			return
		}
	})
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "github token")
	slug := flag.String("slug", os.Getenv("TRAVIS_REPO_SLUG"), "slug of the repository")
	ghr := flag.Bool("ghr", false, "create a .ghr file with the parameters for ghr")
	noci := flag.Bool("noci", false, "run semantic-release locally")
	nochange := flag.Bool("nochange", false, "don't return an error code when the calculated version has already been tagged")
	dry := flag.Bool("dry", false, "do not create release")
	flow := flag.Bool("flow", false, "follow branch naming conventions")
	vFile := flag.Bool("vf", false, "create a .version file")
	showVersion := flag.Bool("version", false, "outputs the semantic-release version")
	updateFile := flag.String("update", "", "updates the version of a certain file")
	branchEnv := flag.Bool("branch_env", false, "use GIT_BRANCH environment variable with branch information")
	defaultBranchFlag := flag.String(
		"default_branch",
		os.Getenv("GIT_DEFAULT_BRANCH"),
		"override the branch to consider the default for creating non-pre-release tags",
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("semantic-release v%s", SRVERSION)
		return
	}

	logger := log.New(os.Stderr, "[semantic-release]: ", 0)
	exitIfError := errorHandler(logger)

	logger.Printf("version: %s\n", SRVERSION)

	conf, err := config.NewConfig(cmd)
	exitIfError(err)

	pluginManager, err := manager.New(conf)
	exitIfError(err)
	exitHandler = func() {
		logger.Println("stopping plugins...")
		pluginManager.Stop()
	}
	defer exitHandler()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		exitIfError(errors.New("terminating..."))
	}()

	if conf.DownloadPlugins {
		exitIfError(pluginManager.FetchAllPlugins())
		logger.Println("all plugins are downloaded")
		os.Exit(0)
	}

	ci, err := pluginManager.GetCICondition()
	exitIfError(err)
	ciName := ci.Name()
	logger.Printf("ci-condition plugin: %s@%s\n", ciName, ci.Version())

	prov, err := pluginManager.GetProvider()
	exitIfError(err)
	provName := prov.Name()
	logger.Printf("provider plugin: %s@%s\n", provName, prov.Version())

	if conf.ProviderOpts["token"] == "" {
		conf.ProviderOpts["token"] = conf.Token
	}
	err = prov.Init(conf.ProviderOpts)
	exitIfError(err)

	logger.Println("getting default branch...")
	repoInfo, err := prov.GetInfo()
	exitIfError(err)
	logger.Println("found default branch: " + repoInfo.DefaultBranch)
	if repoInfo.Private {
		logger.Println("repo is private")
	}

	currentBranch := ""
	if *branchEnv {
		envBranch, present := os.LookupEnv("GIT_BRANCH")
		currentBranch = envBranch
		if !present {
			exitIfError(errors.New("branch not present in env var: GIT_BRANCH"))
		}
	} else {
		curCommitInfo, err := condition.GetCurCommitInfo()
		if err == git.ErrRepositoryNotExists {
			logger.Println(`Repository (.git directory) does not exist in local directory. Be sure to
run go-semantic-release in a git repository`)
		}
		exitIfError(err)
		currentBranch = curCommitInfo.Branch
	}

	logger.Println("found current branch: " + currentBranch)

	if !conf.AllowMaintainedVersionOnDefaultBranch && conf.MaintainedVersion != "" && currentBranch == repoInfo.DefaultBranch {
		exitIfError(fmt.Errorf("maintained version not allowed on default branch"))
	}

	if conf.MaintainedVersion != "" {
		logger.Println("found maintained version: " + conf.MaintainedVersion)
		repoInfo.DefaultBranch = "*"
	}

	prerelease := ""
	if *flow && config.MaintainedVersion == "" {
		switch currentBranch {
		// If branch is defaultBranch -> no pre-latestRelease version
		case defaultBranch:
			prerelease = ""
		// If branch is develop -> beta latestRelease
		case "develop":
			prerelease = "beta"
		default:
			branchPath := strings.Split(currentBranch, "/")
			prerelease = branchPath[len(branchPath)-1]
			prerelease = strcase.ToLowerCamel(prerelease)
		}
	}

	if prerelease != "" {
		logger.Println("Determined prerelease version: " + prerelease)
	}

	if !*noci {
		logger.Println("running CI condition...")
		conditionConfig := map[string]string{
			"token":         conf.Token,
			"defaultBranch": repoInfo.DefaultBranch,
			"private":       fmt.Sprintf("%t", repoInfo.Private),
		}
		for k, v := range conf.CIConditionOpts {
			conditionConfig[k] = v
		}
		err = ci.RunCondition(conditionConfig)
		if err != nil {
			herr := hooksExecutor.NoRelease(&hooks.NoReleaseConfig{
				Reason:  hooks.NoReleaseReason_CONDITION,
				Message: err.Error(),
			})
			if herr != nil {
				logger.Printf("there was an error executing the hooks plugins: %s", herr.Error())
			}
			exitIfError(err, 66)
		}

	}

	logger.Println("getting latest release...")
	latestRelease, err := repo.GetLatestRelease(config.MaintainedVersion, prerelease)
	exitIfError(err)
	logger.Println("found version: " + latestRelease.Version.String())

	if strings.Contains(config.MaintainedVersion, "-") && latestRelease.Version.Prerelease() == "" {
		exitIfError(fmt.Errorf("no pre-release for this version possible"))
	}

	logger.Println("getting commits...")
	rawCommits, err := prov.GetCommits(release.SHA, currentSha)
	exitIfError(err)

	logger.Println("analyzing commits...")
	commitAnalyzer, err := pluginManager.GetCommitAnalyzer()
	exitIfError(err)
	logger.Printf("commit-analyzer plugin: %s@%s\n", commitAnalyzer.Name(), commitAnalyzer.Version())
	exitIfError(commitAnalyzer.Init(conf.ChangelogGeneratorOpts))

	commits := commitAnalyzer.Analyze(rawCommits)

	logger.Println("calculating new version...")
	newVer := semrel.GetNewVersion(commits, latestRelease, prerelease)

	if *nochange && newVer == latestRelease.Version {
		logger.Println("Latest version tag is equal to current commit using version: " + newVer.String())
	} else {
		if newVer == nil {
			exitIfError(errors.New("no change"))
		}
		logger.Println("new version: " + newVer.String())
	}

	logger.Println("generating changelog...")
	changelogGenerator, err := pluginManager.GetChangelogGenerator()
	exitIfError(err)
	logger.Printf("changelog-generator plugin: %s@%s\n", changelogGenerator.Name(), changelogGenerator.Version())
	exitIfError(changelogGenerator.Init(conf.ChangelogGeneratorOpts))

	changelogRes := changelogGenerator.Generate(&generator.ChangelogGeneratorConfig{
		Commits:       commits,
		LatestRelease: release,
		NewVersion:    newVer,
	})
	if conf.Changelog != "" {
		oldFile := make([]byte, 0)
		if conf.PrependChangelog {
			oldFileData, err := ioutil.ReadFile(conf.Changelog)
			if err == nil {
				oldFile = append([]byte("\n"), oldFileData...)
			}
		}
		changelogData := append([]byte(changelogRes), oldFile...)
		exitIfError(ioutil.WriteFile(conf.Changelog, changelogData, 0644))
	}

	if newVer != latestRelease.Version {
		logger.Println("creating release...")
		exitIfError(repo.CreateRelease(commits, latestRelease, newVer, currentBranch))
	}

	if conf.Ghr {
		exitIfError(ioutil.WriteFile(".ghr", []byte(fmt.Sprintf("-u %s -r %s v%s", repoInfo.Owner, repoInfo.Repo, newVer)), 0644))
	}

	if conf.VersionFile {
		exitIfError(ioutil.WriteFile(".version", []byte(newVer), 0644))
	}

	if len(conf.UpdateFiles) == 0 && len(conf.FilesUpdaterPlugins) > 0 {
		logger.Println("warning: file update plugins found but no files marked for update. You may be missing the update flag, e.g. --update package.json")
	}

	if len(conf.UpdateFiles) > 0 {
		logger.Println("updating files...")
		updater, err := pluginManager.GetChainedUpdater()
		exitIfError(err)
		logger.Printf("files-updater plugins: %s\n", strings.Join(updater.GetNameVersionPairs(), ", "))
		exitIfError(updater.Init(conf.FilesUpdaterOpts))

		for _, f := range conf.UpdateFiles {
			exitIfError(updater.Apply(f, newVer))
		}
	}

	herr := hooksExecutor.Success(&hooks.SuccessHookConfig{
		Commits:     commits,
		PrevRelease: release,
		NewRelease: &semrel.Release{
			SHA:     currentSha,
			Version: newVer,
		},
		Changelog: changelogRes,
		RepoInfo:  repoInfo,
	})

	if herr != nil {
		logger.Printf("there was an error executing the hooks plugins: %s", herr.Error())
	}

	logger.Println("done.")
}
