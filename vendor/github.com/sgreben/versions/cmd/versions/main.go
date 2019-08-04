package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/sgreben/versions/pkg/versions"

	"github.com/sgreben/versions/pkg/versionscmd"

	"github.com/posener/complete"
	"github.com/posener/complete/cmd/install"

	"github.com/sgreben/flagvar"

	cli "github.com/jawher/mow.cli"
)

const name = "versions"

var configuration struct {
	JSONIndent int
	Quiet      bool
	Silent     bool
}

func jsonEncode(value interface{}, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", strings.Repeat(" ", configuration.JSONIndent))
	return enc.Encode(value)
}

var exit struct {
	NonzeroBecause []string
}

func main() {
	log.SetOutput(os.Stderr)
	app := cli.App(name, "do things with versions")

	var ( // Global flags
		jsonIndent = app.IntOpt("indent", 0, "Set the indentation of JSON output")
		quiet      = app.BoolOpt("q quiet", false, "Disable all log output (stderr)")
		silent     = app.BoolOpt("s silent", false, "Disable all log output (stderr) and all normal output (stdout)")
	)
	app.Before = func() { // Copy global flag values to `configuration` struct
		configuration.JSONIndent = *jsonIndent
		configuration.Quiet = *quiet
		configuration.Silent = *silent

		if configuration.Silent {
			os.Stdout, _ = os.Open(os.DevNull)
			configuration.Quiet = true
		}
		if configuration.Quiet {
			log.SetOutput(ioutil.Discard)
		}
	}

	completeCmd := complete.Command{
		Sub: complete.Commands{},
		GlobalFlags: complete.Flags{
			"-h":     complete.PredictNothing,
			"--help": complete.PredictNothing,
		},
	}
	completer := complete.New(name, completeCmd)

	completeCmd.Sub["sort"] = complete.Command{
		Flags: complete.Flags{
			"-l":       complete.PredictAnything,
			"--latest": complete.PredictAnything,
		},
		Args: complete.PredictAnything,
	}
	app.Command("sort", "Sort versions", func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [VERSIONS...]"
		var (
			latest   = cmd.IntOpt("l latest", 0, "Print only the latest `N` versions")
			versions = cmd.StringsArg("VERSIONS", nil, "Versions to sort")
		)
		cmd.Action = func() {
			sortCmd(*versions, *latest)
		}
	})

	completeCmd.Sub["compare"] = complete.Command{
		Args: versionscmd.PredictSet1("later", "earlier"),
	}
	app.Command("compare", "Compare versions", func(cmd *cli.Cmd) {
		var (
			nonzeroExitOnFalse = cmd.BoolOpt("fail", false, "Exit with non-zero code if the result is 'false'")
		)
		cmd.Command("later", "Check if a version is strictly later than another version", func(cmd *cli.Cmd) {
			cmd.Spec = "LATER_VERSION EARLIER_VERSION"
			var (
				laterVersion   = cmd.StringArg("LATER_VERSION", "", "The version asserted to be the strictly later version")
				earlierVersion = cmd.StringArg("EARLIER_VERSION", "", "The version asserted to be the strictly earlier version")
			)
			cmd.Action = func() {
				laterCmd(*laterVersion, *earlierVersion, *nonzeroExitOnFalse)
			}
		})
		cmd.Command("earlier", "Check if a version is strictly earlier than another version", func(cmd *cli.Cmd) {
			cmd.Spec = "EARLIER_VERSION LATER_VERSION"
			var (
				earlierVersion = cmd.StringArg("EARLIER_VERSION", "", "The version asserted to be the strictly earlier version")
				laterVersion   = cmd.StringArg("LATER_VERSION", "", "The version asserted to be the strictly later version")
			)
			cmd.Action = func() {
				laterCmd(*laterVersion, *earlierVersion, *nonzeroExitOnFalse)
			}
		})
	})

	completeCmd.Sub["fetch"] = complete.Command{
		Args: versionscmd.PredictSet1("git", "docker"),
		Flags: complete.Flags{
			"-l":       complete.PredictAnything,
			"--latest": complete.PredictAnything,
		},
	}
	app.Command("fetch", "Fetch versions", func(cmd *cli.Cmd) {
		var (
			latest = cmd.IntOpt("l latest", 0, "Print only the latest `N` versions")
		)
		cmd.Command("git", "Fetch versions from Git tags", func(cmd *cli.Cmd) {
			cmd.Spec = "REPOSITORY"
			var (
				url = cmd.StringArg("REPOSITORY", "", "Git repository")
			)
			cmd.Action = func() {
				fetchFromGitCmd(*url, *latest)
			}
		})
		cmd.Command("docker", "Fetch versions from Docker image tags", func(cmd *cli.Cmd) {
			cmd.Spec = "REPOSITORY"
			var (
				repository = cmd.StringArg("REPOSITORY", "", "Docker repository")
			)
			cmd.Action = func() {
				fetchFromDockerCmd(*repository, *latest)
			}
		})
	})

	completeCmd.Sub["select"] = complete.Command{
		Args: versionscmd.PredictSet1("single", "all", "mvs"),
		Flags: complete.Flags{
			"--from-git":    complete.PredictAnything,
			"--from-docker": complete.PredictAnything,
		},
	}
	app.Command("select", "Select versions given constraints", func(cmd *cli.Cmd) {
		var (
			fromGit    = cmd.StringOpt("from-git", "", "Fetch candidate versions from Git tags")
			fromDocker = cmd.StringOpt("from-docker", "", "Fetch candidate versions from Docker tags")
		)
		cmd.Command("single", "Select a single version", func(cmd *cli.Cmd) {
			cmd.Spec = "CONSTRAINT [VERSIONS...]"
			var (
				constraint = cmd.StringArg("CONSTRAINT", "", "Version constraint")
				versions   = cmd.StringsArg("VERSIONS", nil, "Candidate versions")
			)
			cmd.Action = func() {
				if *fromGit != "" {
					vs, err := fetchFromGit(*fromGit, 0)
					if err != nil {
						exit.NonzeroBecause = append(exit.NonzeroBecause, err.Error())
					}
					for _, v := range vs {
						*versions = append(*versions, v.Version.String())
					}
				}
				if *fromDocker != "" {
					vs, err := fetchFromDocker(*fromDocker, 0)
					if err != nil {
						exit.NonzeroBecause = append(exit.NonzeroBecause, err.Error())
					}
					for _, v := range vs {
						*versions = append(*versions, v.Version.String())
					}
				}
				selectSingleCmd(*constraint, *versions)
			}
		})
		cmd.Command("all", "Select all matching versions", func(cmd *cli.Cmd) {
			cmd.Spec = "CONSTRAINT [VERSIONS...]"
			var (
				constraint = cmd.StringArg("CONSTRAINT", "", "Version constraint")
				versions   = cmd.StringsArg("VERSIONS", nil, "Candidate versions")
			)
			cmd.Action = func() {
				if *fromGit != "" {
					vs, err := fetchFromGit(*fromGit, 0)
					if err != nil {
						exit.NonzeroBecause = append(exit.NonzeroBecause, err.Error())
					}
					for _, v := range vs {
						*versions = append(*versions, v.Version.String())
					}
				}
				if *fromDocker != "" {
					vs, err := fetchFromDocker(*fromDocker, 0)
					if err != nil {
						exit.NonzeroBecause = append(exit.NonzeroBecause, err.Error())
					}
					for _, v := range vs {
						*versions = append(*versions, v.Version.String())
					}
				}
				selectAllCmd(*constraint, *versions)
			}
		})
		cmd.Command("mvs", "Select versions to satisfy a constraint graph using MVS (https://research.swtch.com/vgo-mvs)", func(cmd *cli.Cmd) {
			cmd.Spec = "CONSTRAINTS..."
			constraints := flagvar.JSONs{
				Value: func() interface{} {
					return &versions.ConstraintStringGraph{}
				},
			}
			cmd.VarArg("CONSTRAINTS", &constraints, `constraint graph (JSON structure: {"my-package":{"1.0": {"other-package":"~0.0.1"}}})`)
			cmd.Action = func() {
				graph := versions.ConstraintGraph{}
				for _, val := range constraints.Values {
					err := graph.Add(*val.(*versions.ConstraintStringGraph))
					if err != nil {
						log.Fatal(err)
					}
				}
				selectMvsCmd(graph)
			}
		})
	})

	app.Command("complete", "Shell completion (zsh, fish, bash)", func(cmd *cli.Cmd) {
		cmd.Command("install", "Install all completions", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				err := install.Install(name)
				if err != nil {
					log.Fatal(err)
				}
			}
		})
		cmd.Command("uninstall", "Uninstall all completions", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				err := install.Uninstall(name)
				if err != nil {
					log.Fatal(err)
				}
			}
		})
	})

	completeCmd.Sub["help"] = complete.Command{
		Args: versionscmd.PredictSet1("fetch", "sort", "compare", "complete"),
	}
	app.Command("help", "Display help for a command", func(cmd *cli.Cmd) {
		cmd.Spec = "[COMMAND...]"
		var (
			command = cmd.StringsArg("COMMAND", nil, "Command to show help for")
		)
		cmd.Action = func() {
			args := append([]string{name}, *command...)
			args = append(args, "--help")
			app.Run(args)
		}
	})

	if completer.Complete() {
		return
	}

	app.Run(os.Args)

	if len(exit.NonzeroBecause) > 0 {
		log.Printf("non-zero exit: %s", strings.Join(exit.NonzeroBecause, ", "))
		os.Exit(1)
	}
}
