package main

import (
	"flag"
	"fmt"
	"github.com/nim4/diff-cov/diffCoverage"
	"github.com/nim4/diff-cov/git"
	"github.com/nim4/diff-cov/profile/coverprofile"
	"os"
	"strings"
)

var (
	flagPackage      string
	flagIgnoreFiles  string
	flagIgnoreBranch string
	flagTargetBranch string
	flagCoverProfile string
	flagMinimumLine  int
	flagMinimumCov   float64
	flagFetchTarget  bool
	flagVerbose      bool
)

func setPackage() bool {
	if flagPackage == "" {
		dir, err := os.Getwd()
		if err != nil {
			return false
		}

		p := strings.SplitN(dir, "/go/src/", 2)
		if len(p) != 2 {
			return false
		}
		flagPackage = p[1]
	}

	if !strings.HasSuffix(flagPackage, "/") {
		flagPackage += "/"
	}
	return true
}

func main() {
	flag.StringVar(&flagCoverProfile,
		"coverprofile", "cover.out",
		"Path of 'coverprofile' file")
	flag.StringVar(&flagIgnoreFiles,
		"ignore-files", "_test.go,_gen.go,_mock.go",
		"Ignore files with given suffix")
	flag.StringVar(&flagIgnoreBranch,
		"ignore-branches", "hotfix,bugfix",
		"Ignore branches which contains given words(case-insensitive)")
	flag.StringVar(&flagPackage,
		"package", "",
		"Package import path(if not set, will try to extract it from the current working directory)")
	flag.StringVar(&flagTargetBranch,
		"target", "origin/master",
		"Target branch")
	flag.BoolVar(&flagFetchTarget,
		"fetch", true,
		"Fetch the target branch")
	flag.IntVar(&flagMinimumLine,
		"min-diff", 0,
		"Minimum diff size to trigger coverage check")
	flag.Float64Var(&flagMinimumCov,
		"min-cov", 0,
		"Minimum required test coverage")
	flag.BoolVar(&flagVerbose,
		"v", false,
		"Verbose")
	flag.Parse()

	cb, err := git.CurrentBranch()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, ignoreBranch := range strings.Split(flagIgnoreBranch, ",") {
		if strings.Contains(strings.ToLower(cb), ignoreBranch) {
			fmt.Printf("Ignoring branch %q(contains %q)! Good bye!\n", cb, ignoreBranch)
			os.Exit(0)
		}
	}

	if !setPackage() {
		fmt.Println("provide package import path: ex. github.com/nim4/example")
		os.Exit(1)
	}

	if flagFetchTarget {
		err := git.Fetch()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	ignoreFiles := strings.Split(flagIgnoreFiles, ",")
	profile := coverprofile.NewGoProfile(flagPackage, flagCoverProfile)
	coverage, err := profile.GetCoverage(ignoreFiles, flagVerbose)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	diff, err := git.Diff(flagTargetBranch)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	goDiff, goTestedDiff, err := diffCoverage.Calculate(coverage, diff, ignoreFiles, flagVerbose)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	difCov := 0.0
	if goDiff > 0 {
		difCov = float64(goTestedDiff) / float64(goDiff) * 100
	}
	fmt.Printf("%d/%d = %.2f%%\n", goTestedDiff, goDiff, difCov)
	if goDiff > flagMinimumLine && difCov < flagMinimumCov {
		os.Exit(1)
	}
}
