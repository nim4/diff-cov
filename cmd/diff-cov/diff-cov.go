package main

import (
	"flag"
	"fmt"
	"github.com/waigani/diffparser"
	"golang.org/x/tools/cover"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	flagPackage      string
	flagIgnore       string
	flagTargetBranch string
	flagCoverProfile string
	flagMinimumCov   float64
	ignore           []string
)

func diff() ([]byte, error) {
	f, err := ioutil.TempFile(os.TempDir(), "diff-")
	if err != nil {
		return nil, err
	}

	output := fmt.Sprintf("--output=%s", f.Name())

	err = exec.Command(
		"git", "diff",
		"--ignore-all-space", "--ignore-blank-lines",
		"--no-color", "--no-ext-diff", "-U0", output, flagTargetBranch,
	).Run()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(f.Name())
}

func shouldCount(file string) bool {
	if !strings.HasSuffix(file, ".go") {
		return false
	}

	for _, suffix := range ignore {
		if strings.HasSuffix(file, suffix) {
			return false
		}
	}

	return true
}

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
		"Path of coverage profile file")
	flag.StringVar(&flagIgnore,
		"ignore", "_test.go,_gen.go,_mock.go",
		"Ignore files with given suffix")
	flag.StringVar(&flagPackage,
		"package", "",
		"Package import path(if not set, will try to extract it from the current working directory)")
	flag.StringVar(&flagTargetBranch,
		"target", "origin/master",
		"Target branch")
	flag.Float64Var(&flagMinimumCov,
		"min", 50,
		"Minimum required test coverage")
	flag.Parse()

	if !setPackage() {
		fmt.Println("provide package import path: ex. github.com/nim4/example")
		os.Exit(1)
	}

	ignore = strings.Split(flagIgnore, ",")
	ps, err := cover.ParseProfiles(flagCoverProfile)
	if err != nil {
		log.Fatal(err)
	}

	coverage := make(map[string]map[int]bool, len(ps))
	for _, p := range ps {
		coverage[p.FileName] = make(map[int]bool)
		for _, block := range p.Blocks {
			for line := block.StartLine; line <= block.EndLine; line++ {
				coverage[p.FileName][line] = true
			}
		}
	}

	b, err := diff()
	if err != nil {
		log.Fatal(err)
	}

	p, err := diffparser.Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}

	goDiff := 0
	goTestedDiff := 0
	for file, changes := range p.Changed() {
		if !shouldCount(file) {
			continue
		}
		fmt.Printf("%s %d\n", flagPackage+file, len(changes))
		goDiff += len(changes)

		// file has any test coverage?
		if _, ok := coverage[flagPackage+file]; !ok {
			continue
		}

		for _, line := range changes {
			if _, ok := coverage[flagPackage+file][line]; ok {
				goTestedDiff++
			}
		}
	}
	difCov := float64(goTestedDiff)/float64(goDiff)*100
	fmt.Printf("%d/%d = %.2f%%\n", goTestedDiff, goDiff, difCov)
	if difCov < flagMinimumCov {
		os.Exit(1)
	}
}
