package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/keilerkonzept/terraform-module-versions/pkg/httputil"
	"github.com/keilerkonzept/terraform-module-versions/pkg/modulecall"
	"github.com/keilerkonzept/terraform-module-versions/pkg/output"
	"github.com/keilerkonzept/terraform-module-versions/pkg/registry"
	"github.com/keilerkonzept/terraform-module-versions/pkg/scan"
	"github.com/keilerkonzept/terraform-module-versions/pkg/update"

	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/sgreben/flagvar"

	ghttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

var (
	appName       = "terraform-module-versions"
	version       = "3-SNAPSHOT"
	updatesClient = update.Client{
		Registry: registry.Client{
			HTTP: http.DefaultClient,
		},
	}
	auth = ghttp.BasicAuth{}

	config struct {
		Paths                           []string
		ModuleNames                     flagvar.StringSet
		Output                          flagvar.Enum
		OutputFormat                    output.Format
		RegistryHeaders                 flagvar.Assignments
		Quiet                           bool
		MatchingUpdatesFoundNonzeroExit bool
		AnyUpdatesFoundNonzeroExit      bool
		All                             bool
		GenerateSed                     bool
		IncludePrereleaseVersions       bool
	}
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("[" + appName + "] ")
	log.SetOutput(os.Stderr)

	config.Output.Choices = output.FormatNames
	config.Output.Value = string(output.FormatMarkdown)
	config.OutputFormat = output.FormatMarkdown
	config.RegistryHeaders.Separator = ":"

	rootFlagSet := flag.NewFlagSet(appName, flag.ExitOnError)
	listFlagSet := flag.NewFlagSet(appName+" list", flag.ExitOnError)
	checkFlagSet := flag.NewFlagSet(appName+" check", flag.ExitOnError)

	rootFlagSet.BoolVar(&config.Quiet, "quiet", false, "suppress log output (stderr)")
	rootFlagSet.BoolVar(&config.Quiet, "q", false, "(alias for -quiet)")
	rootFlagSet.Var(&config.Output, "output", "output format, "+config.Output.Help())
	rootFlagSet.Var(&config.Output, "o", "(alias for -output)")
	listFlagSet.Var(&config.Output, "output", "output format, "+config.Output.Help())
	listFlagSet.Var(&config.Output, "o", "(alias for -output)")
	checkFlagSet.Var(&config.Output, "output", "output format, "+config.Output.Help())
	checkFlagSet.Var(&config.Output, "o", "(alias for -output)")
	checkFlagSet.BoolVar(&config.MatchingUpdatesFoundNonzeroExit, "e", config.MatchingUpdatesFoundNonzeroExit, "(alias for -updates-found-nonzero-exit)")
	checkFlagSet.BoolVar(&config.MatchingUpdatesFoundNonzeroExit, "updates-found-nonzero-exit", config.MatchingUpdatesFoundNonzeroExit, "exit with a nonzero code when modules with updates matching are found (respecting version constraints)")
	checkFlagSet.BoolVar(&config.AnyUpdatesFoundNonzeroExit, "n", config.AnyUpdatesFoundNonzeroExit, "(alias for -any-updates-found-nonzero-exit)")
	checkFlagSet.BoolVar(&config.AnyUpdatesFoundNonzeroExit, "any-updates-found-nonzero-exit", config.AnyUpdatesFoundNonzeroExit, "exit with a nonzero code when modules with updates are found (ignoring version constraints)")
	checkFlagSet.BoolVar(&config.IncludePrereleaseVersions, "pre-release", config.IncludePrereleaseVersions, "include pre-release versions")
	checkFlagSet.BoolVar(&config.All, "a", config.All, "(alias for -all)")
	checkFlagSet.BoolVar(&config.All, "all", config.All, "include modules without updates")
	listFlagSet.Var(&config.ModuleNames, "module", "include this module (may be specified repeatedly. by default, all modules are included)")
	checkFlagSet.Var(&config.ModuleNames, "module", "include this module (may be specified repeatedly. by default, all modules are included)")
	checkFlagSet.Var(&config.RegistryHeaders, "H", "(alias for -registry-header)")
	checkFlagSet.Var(&config.RegistryHeaders, "registry-header", fmt.Sprintf("extra HTTP headers for requests to Terraform module registries (%s, may be specified repeatedly)", config.RegistryHeaders.Help()))
	checkFlagSet.BoolVar(&config.GenerateSed, "sed", config.GenerateSed, "generate sed statements for upgrade")

	cmdList := &ffcli.Command{
		Name:       "list",
		ShortUsage: appName + " list [options] [<path> ...]",
		ShortHelp:  "List referenced terraform modules with their detected versions",
		FlagSet:    listFlagSet,
		Exec: func(_ context.Context, args []string) error {
			config.Paths = args
			list(scanForModuleCalls())
			return nil
		},
	}
	cmdList.LongHelp = cmdList.ShortHelp

	cmdCheck := &ffcli.Command{
		Name:       "check",
		ShortUsage: appName + " check [options] [<path> ...]",
		ShortHelp:  "Check referenced terraform modules' sources for newer versions",
		FlagSet:    checkFlagSet,
		Exec: func(_ context.Context, args []string) error {
			config.Paths = args
			updates(scanForModuleCalls())
			return nil
		},
	}
	cmdCheck.LongHelp = cmdCheck.ShortHelp

	cmdVersion := &ffcli.Command{
		Name:       "version",
		ShortUsage: appName + " version",
		ShortHelp:  "Print version and exit",
		Exec: func(_ context.Context, args []string) error {
			fmt.Println(version)
			os.Exit(0)
			return nil
		},
	}
	cmdVersion.LongHelp = cmdVersion.ShortHelp

	cmdRoot := &ffcli.Command{
		ShortUsage:  appName + " [options] <subcommand>",
		FlagSet:     rootFlagSet,
		Subcommands: []*ffcli.Command{cmdList, cmdCheck, cmdVersion},
		Exec: func(_ context.Context, args []string) error {
			return flag.ErrHelp
		},
	}


	if err := cmdRoot.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if config.Quiet {
		log.SetOutput(ioutil.Discard)
	}
	if f, ok := output.ParseFormatName(config.Output.Value); ok {
		config.OutputFormat = f
	}
	if len(config.RegistryHeaders.Values) > 0 {
		headers := make(http.Header, len(config.RegistryHeaders.Values))
		for _, kv := range config.RegistryHeaders.Values {
			headers.Add(kv.Key, strings.TrimLeftFunc(kv.Value, unicode.IsSpace))
		}
		updatesClient.Registry.HTTP.Transport = httputil.AddHeadersRoundtripper{
			Headers: headers,
			Nested:  http.DefaultTransport,
		}
	}
	if githubToken := os.Getenv("GITHUB_TOKEN"); githubToken != "" {
		updatesClient.GitAuth = &githttp.BasicAuth{
			Username: githubToken,
		}
	}
	basicAuthUsername := os.Getenv("BASICAUTH_USERNAME")
	basicAuthPassword := os.Getenv("BASICAUTH_PASSWORD")
	if basicAuthUsername != "" && basicAuthPassword != "" {
		updatesClient.GitAuth = &githttp.BasicAuth{
			Username: basicAuthUsername,
			Password: basicAuthPassword,
		}
	}
	if err := cmdRoot.Run(context.Background()); err != nil && !errors.Is(err, flag.ErrHelp) {
		log.Fatal(err)
	}
}

func scanForModuleCalls() []scan.Result {
	scanResults, err := scan.Scan(config.Paths)
	if err != nil {
		log.Fatal(err)
	}
	moduleNamesFilter := config.ModuleNames.Value
	moduleNamesFilterEmpty := len(moduleNamesFilter) == 0
	scanResultsFiltered := make([]scan.Result, 0, len(scanResults))
	for _, r := range scanResults {
		include := moduleNamesFilterEmpty || moduleNamesFilter[r.ModuleCall.Name]
		if !include {
			continue
		}
		scanResultsFiltered = append(scanResultsFiltered, r)
	}
	return scanResultsFiltered
}

func list(scanResults []scan.Result) {
	var out output.Modules
	for _, m := range scanResults {
		parsed, err := modulecall.Parse(m.ModuleCall, config.UseUrl)
		if err != nil {
			log.Printf("error: %v", err)
			out = append(out, output.Module{
				Path:              m.Path,
				Name:              m.ModuleCall.Name,
				Source:            m.ModuleCall.Source,
				VersionConstraint: m.ModuleCall.Version,
			})
			continue
		}
		out = append(out, output.Module{
			Path:              m.Path,
			Name:              m.ModuleCall.Name,
			Source:            m.ModuleCall.Source,
			VersionConstraint: parsed.ConstraintsString,
			Version:           parsed.VersionString,
			Type:              parsed.Source.Type(),
		})
	}
	sort.Sort(out)
	if err := out.Write(os.Stdout, config.OutputFormat); err != nil {
		log.Fatal(err)
	}
}

func updates(scanResults []scan.Result) {
	var (
		out                  output.Updates
		foundMatchingUpdates bool
		foundAnyUpdates      bool
	)
	if len(config.BasicAuth) > 0 {
		var authParts = strings.Split(config.BasicAuth, ":")
		if len(authParts) != 2 {
			log.Fatal("IllegalValueException: basic-auth should be [username]:[password]!")
			os.Exit(2)
		}
		auth.Username = authParts[0]
		auth.Password = authParts[1]
		updatesClient.GitAuth = &auth
	}
	for _, m := range scanResults {
		parsed, err := modulecall.Parse(m.ModuleCall, config.UseUrl)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		update, err := updatesClient.Update(*parsed.Source, parsed.Version, parsed.Constraints, config.IncludePrereleaseVersions)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		updateOutput := output.Update{
			Path:              m.Path,
			Name:              m.ModuleCall.Name,
			Source:            m.ModuleCall.Source,
			VersionConstraint: parsed.ConstraintsString,
			Version:           parsed.VersionString,
			LatestMatching:    update.LatestMatchingVersion,
			MatchingUpdate:    update.LatestMatchingUpdate != "",
			LatestOverall:     update.LatestOverallVersion,
			NonMatchingUpdate: update.LatestOverallUpdate != "" && update.LatestOverallUpdate != update.LatestMatchingVersion,
		}
		hasUpdate := false
		if updateOutput.MatchingUpdate {
			foundMatchingUpdates = true
			foundAnyUpdates = true
			hasUpdate = true
		}
		if updateOutput.NonMatchingUpdate {
			foundAnyUpdates = true
			hasUpdate = true
			foundNonMatchingUpdates = true
		}
		if !config.All && !hasUpdate {
			continue
		}
		out = append(out, updateOutput)
	}
	sort.Sort(out)
	if err := out.Format(os.Stdout, config.OutputFormat); err != nil {
		log.Fatal(err)
	}

	if config.GenerateSed {
		out.GenerateSed()
	}

	if config.MatchingUpdatesFoundNonzeroExit {
		if foundMatchingUpdates {
			os.Exit(1)
		}
	}
	if config.AnyUpdatesFoundNonzeroExit {
		if foundAnyUpdates {
			os.Exit(1)
		}
	}
}
