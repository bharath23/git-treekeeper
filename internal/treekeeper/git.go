package treekeeper

import (
	"github.com/bharath23/git-treekeeper/internal/git"
)

func PassThrough(args []string) (Response, error) {
	Verbose("Pass-through: git %v", args)
	err := git.RawRun(args...)
	return Response{
		Kind: ResponsePassThrough,
		PassThrough: &PassThroughResult{
			Args: args,
		},
	}, err
}
