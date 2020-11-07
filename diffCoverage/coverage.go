package diffCoverage

import (
	"fmt"
	"github.com/nim4/diff-cov/profile"
	"github.com/nim4/diff-cov/utils"
	"github.com/waigani/diffparser"
)

func Calculate(coverage profile.Coverage, diff string, ignoreFiles []string, verbose bool) (int, int, error) {
	goDiff := 0
	goTestedDiff := 0

	p, err := diffparser.Parse(diff)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing diff: %v", err)
	}

	for file, changes := range p.Changed() {
		if !utils.ShouldCountFile(file, ignoreFiles) {
			continue
		}
		if verbose {
			fmt.Printf("%s %d\n", file, len(changes))
		}

		// file has any test diffCoverage?
		if _, ok := coverage[file]; !ok {
			fmt.Printf("No coverage found for %q! is coverprofile updated?\n", file)
			goDiff += len(changes)
			continue
		}

		for _, line := range changes {
			tested, ok := coverage[file][line]
			if ok {
				goDiff++
				if tested {
					goTestedDiff++
				}
			}
		}
	}

	return goDiff, goTestedDiff, nil
}