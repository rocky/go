// Package for getting status of a repo under vcs.
package vcs

import "os/exec"

type Type uint8

const (
	Git Type = iota
	Hg
)

// VcsType returns a vcsType string compatible with github.com/sourcegraph/go-vcs notation.
func (t Type) VcsType() (vcsType string) {
	switch t {
	case Git:
		return "git"
	case Hg:
		return "hg"
	default:
		panic("bad vcs.Type")
	}
}

// TODO: Add comments.
type Vcs interface {
	RootPath() string
	Type() Type

	GetStatus() string
	GetStash() string

	GetRemote() string // Get primary remote repository Url.

	GetDefaultBranch() string
	GetLocalBranch() string

	GetLocalRev() string
	GetRemoteRev() string

	// Returns true if given commit is contained in the default local branch.
	IsContained(rev string) bool
}

type commonVcs struct {
	rootPath string
}

func (this *commonVcs) RootPath() string {
	return this.rootPath
}

// New returns a new Vcs if path is under version control, otherwise nil.
// It should be a valid path.
// TODO: Use a better type for path, e.g., github.com/shurcooL/go/path.
func New(path string) Vcs {
	// TODO: Try to figure out vcs provider with a more constant-time operation.
	// TODO: Potentially check in parallel.
	for _, vcsProvider := range vcsProviders {
		if vcs := vcsProvider(path); vcs != nil {
			return vcs
		}
	}

	return nil
}

// Experimental, NewFromType returns a Vcs repository of the specified type without a local representation.
// Operations that require a local repository will fail.
func NewFromType(t Type) Vcs {
	switch t {
	case Git:
		return &gitVcs{}
	case Hg:
		return &hgVcs{}
	default:
		panic("bad vcs.Type")
	}
}

type vcsProvider func(path string) Vcs

var vcsProviders []vcsProvider

func addVcsProvider(s vcsProvider) {
	vcsProviders = append(vcsProviders, s)
}

func init() {
	// As an optimization, add Vcs providers sorted by the most likely first.

	// git
	if _, err := exec.LookPath("git"); err == nil {
		addVcsProvider(func(path string) Vcs {
			if isRepo, rootPath := getGitRepoRoot(path); isRepo {
				return &gitVcs{commonVcs{rootPath: rootPath}}
			}
			return nil
		})
	}

	// hg
	if _, err := exec.LookPath("hg"); err == nil {
		addVcsProvider(func(path string) Vcs {
			if isRepo, rootPath := getHgRepoRoot(path); isRepo {
				return &hgVcs{commonVcs{rootPath: rootPath}}
			}
			return nil
		})
	}
}
