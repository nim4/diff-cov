package coverprofile

import (
	"fmt"
	"github.com/nim4/diff-cov/profile"
	"github.com/nim4/diff-cov/utils"
	"golang.org/x/tools/cover"
	"strings"
)

type GoProfile struct {
	packageName string
	profilePath string
}

func NewGoProfile(packageName string, profilePath string) *GoProfile {
	return &GoProfile{
		packageName: packageName,
		profilePath: profilePath,
	}
}

func (g GoProfile) GetCoverage(ignoreFiles []string, verbose bool) (profile.Coverage, error) {
	ps, err := cover.ParseProfiles(g.profilePath)
	if err != nil {
		return nil, fmt.Errorf("Error parsing %q: %v\n", g.profilePath, err)
	}

	coverage := make(profile.Coverage, len(ps))
	for _, p := range ps {
		if !utils.ShouldCountFile(p.FileName, ignoreFiles) {
			if verbose {
				fmt.Printf("Ignoring %q\n", p.FileName)
			}
			continue
		}

		file := strings.TrimPrefix(p.FileName, g.packageName)
		coverage[file] = make(map[int]bool)
		for _, block := range p.Blocks {
			for line := block.StartLine; line <= block.EndLine; line++ {
				coverage[file][line] = block.Count > 0
			}
		}
	}

	return coverage, nil
}