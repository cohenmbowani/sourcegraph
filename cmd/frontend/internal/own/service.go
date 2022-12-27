package own

import (
	"bytes"
	"context"

	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/authz"
	"github.com/sourcegraph/sourcegraph/internal/gitserver"
	"github.com/sourcegraph/sourcegraph/internal/own/codeowners"

	codeownerspb "github.com/sourcegraph/sourcegraph/internal/own/codeowners/proto"
)

// Service gives access to code ownership data.
// At this point only data from CODEOWNERS file is presented, if available.
type Service interface {
	// Owners returns a CODEOWNERS file from a given repository at given commit ID.
	// In the case the file can not be found, `nil` `*codeownerspb.File` and `nil` `error` is returned.
	Owners(context.Context, api.RepoName, api.CommitID) (*codeownerspb.File, error)
}

var _ Service = service{}

func NewService(git gitserver.Client) Service {
	return service{git: git}
}

type service struct {
	git gitserver.Client
}

// codeownersLocations contain all the locations where CODEOWNERS file
// is expected to be found relative to the repository root directory.
// These are in line with GitHub and GitLab documentation.
var codeownersLocations = []string{
	"CODEOWNERS",
	".github/CODEOWNERS",
	".gitlab/CODEOWNERS",
	"docs/CODEOWNERS",
}

func (s service) Owners(ctx context.Context, repoName api.RepoName, commitID api.CommitID) (*codeownerspb.File, error) {
	var content []byte
	var err error
	for _, path := range codeownersLocations {
		content, err = s.git.ReadFile(
			ctx,
			authz.DefaultSubRepoPermsChecker,
			repoName,
			commitID,
			path,
		)
		if err == nil && content != nil {
			break
		}
	}
	if content == nil {
		return nil, nil
	}
	return codeowners.Parse(bytes.NewReader(content))
}
