package main

import (
	"fmt"
	"os"

	"github.com/mndfcked/branch-visualizer/internal/git"
	"github.com/mndfcked/branch-visualizer/internal/visualizer"
)

func main() {
	currentBranch, err := git.GetCurrentBranch()
	if err != nil {
		fmt.Println("Error getting current branch:", err)
		os.Exit(1)
	}

	root := &visualizer.Branch{Name: currentBranch}
	visualizer.BuildBranchTree(root)

	fmt.Println("**Current dependencies on/for this PR:**")
	visualizer.PrintBranch(root, "", true, true)
}
