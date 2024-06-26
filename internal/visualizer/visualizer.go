package visualizer

import (
	"fmt"

	"github.com/mndfcked/branch-visualizer/internal/git"
)

type Branch struct {
	Name     string
	PRNumber int
	Title    string
	Children []*Branch
}

func BuildBranchTree(branch *Branch) error {
	children, err := git.GetChildBranches(branch.Name)
	if err != nil {
		return fmt.Errorf("error getting child branches: %w", err)
	}

	for _, childName := range children {
		child := &Branch{Name: childName}
		branch.Children = append(branch.Children, child)
		if err := BuildBranchTree(child); err != nil {
			return err
		}
	}

	prNumber, title, err := git.GetPRInfo(branch.Name)
	if err != nil {
		return fmt.Errorf("error getting PR info: %w", err)
	}

	branch.PRNumber = prNumber
	branch.Title = title

	return nil
}

func PrintBranch(branch *Branch, prefix string, isLast bool, isRoot bool) {
	currentPrefix := prefix
	if !isRoot {
		if isLast {
			currentPrefix += "  "
		} else {
			currentPrefix += "│ "
		}
	}

	fmt.Print(currentPrefix)
	if !isRoot {
		if isLast {
			fmt.Print("- ")
		} else {
			fmt.Print("- ")
		}
	}

	if branch.PRNumber == 0 {
		fmt.Printf("**%s**", branch.Name)
	} else {
		fmt.Printf("**PR** 🦊 %s #%d 🔄", branch.Title, branch.PRNumber)
	}

	if isRoot {
		fmt.Print(" 👈")
	}
	fmt.Println()

	for i, child := range branch.Children {
		PrintBranch(child, currentPrefix, i == len(branch.Children)-1, false)
	}
}
