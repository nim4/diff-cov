package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// TODO: Fetch target not master
func Fetch() error {
	out, err := exec.Command(
		"git", "fetch", "origin", "master:refs/remotes/origin/master",
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error fetching: %v\n%v", string(out), err)
	}

	return nil
}

func CurrentBranch() (string, error) {
	out, err := exec.Command(
		"git", "branch", "--show-current",
	).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		return "", fmt.Errorf("error getting current branch: %v\n%v", string(out), err)
	}

	return string(out), nil
}

func Diff(targetBranch string) (string, error) {
	f, err := ioutil.TempFile(os.TempDir(), "diff-cov-")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %v", err)
	}
	_ = f.Close()

	output := fmt.Sprintf("--output=%s", f.Name())
	target := fmt.Sprintf("%s..HEAD", targetBranch)
	out, err := exec.Command(
		"git", "diff",
		"--ignore-all-space", "--ignore-blank-lines",
		"--no-color", "--no-ext-diff", "-U0", output, target,
	).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error getting diff: %v\n%v", string(out), err)
	}

	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return "", fmt.Errorf("error reading temp file: %v", err)
	}
	_ = os.Remove(f.Name())

	return string(b), nil
}
