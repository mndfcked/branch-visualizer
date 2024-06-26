package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func GetChildBranches(branchName string) ([]string, error) {
	cmd := exec.Command("git", "branch", "--list", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error listing branches: %w", err)
	}

	branches := strings.Split(strings.TrimSpace(string(output)), "\n")
	var children []string

	for _, b := range branches {
		if b != branchName {
			baseCmd := exec.Command("git", "merge-base", branchName, b)
			baseOutput, err := baseCmd.Output()
			if err != nil {
				continue
			}

			revParseCmd := exec.Command("git", "rev-parse", branchName)
			revParseOutput, err := revParseCmd.Output()
			if err != nil {
				continue
			}

			if strings.TrimSpace(string(baseOutput)) == strings.TrimSpace(string(revParseOutput)) {
				children = append(children, b)
			}
		}
	}

	return children, nil
}

func GetPRInfo(branchName string) (int, string, error) {
	cmd := exec.Command("git", "config", "--get", fmt.Sprintf("branch.%s.description", branchName))
	output, err := cmd.Output()
	if err != nil {
		return 0, "", nil
	}

	info := strings.TrimSpace(string(output))
	parts := strings.SplitN(info, ":", 2)
	if len(parts) != 2 {
		return 0, "", nil
	}

	var prNumber int
	_, err = fmt.Sscanf(parts[0], "%d", &prNumber)
	if err != nil {
		return 0, "", nil
	}

	return prNumber, strings.TrimSpace(parts[1]), nil
}
