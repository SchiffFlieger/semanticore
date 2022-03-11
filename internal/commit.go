package internal

import (
	"regexp"
	"strconv"
	"strings"
)

type CommitType string

const (
	TypeFix      CommitType = "fix"
	TypeFeat     CommitType = "feat"
	TypeTest     CommitType = "test"
	TypeChore    CommitType = "chore"
	TypeOps      CommitType = "ops"
	TypeDocs     CommitType = "docs"
	TypePerf     CommitType = "perf"
	TypeRefactor CommitType = "refactor"
	TypeSecurity CommitType = "security"
	TypeOther    CommitType = "other"
)

var commitRegexp = regexp.MustCompile(`([a-zA-Z]*)\s*([\(\[]([^\]\)]*)[\]\)])?\s*?(!?):?\s*(.*)`)

func ParseCommitMessage(msg string) (CommitType, string, string, bool) {
	match := commitRegexp.FindStringSubmatch(msg)
	var commitType, scope, description string
	var typ CommitType
	var major = false

	if len(match) == 6 {
		commitType, scope, description = strings.ToLower(match[1]), strings.ToLower(match[3]), strings.TrimSpace(match[5])
		if match[4] == "!" {
			major = true
		}
	}
	if len(description) == 0 {
		commitType = ""
	}

	switch {
	case strings.HasPrefix(commitType, "fix"), strings.HasPrefix(commitType, "bug"):
		typ = TypeFix
	case strings.HasPrefix(commitType, "feat"):
		typ = TypeFeat
	case strings.HasPrefix(commitType, "test"):
		typ = TypeTest
	case strings.HasPrefix(commitType, "chore"), strings.HasPrefix(commitType, "update"):
		typ = TypeChore
	case strings.HasPrefix(commitType, "ops"), strings.HasPrefix(commitType, "ci"), strings.HasPrefix(commitType, "cd"), strings.HasPrefix(commitType, "build"):
		typ = TypeOps
	case strings.HasPrefix(commitType, "doc"):
		typ = TypeDocs
	case strings.HasPrefix(commitType, "perf"):
		typ = TypePerf
	case strings.HasPrefix(commitType, "refactor"), strings.HasPrefix(commitType, "rework"):
		typ = TypeRefactor
	case strings.HasPrefix(commitType, "sec"):
		typ = TypeSecurity
	default:
		typ = TypeOther
		scope = ""
		description = msg
	}

	scope = strings.TrimSpace(scope)
	commitDescription := ""
	for _, line := range strings.Split(description, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 0 && commitDescription == "" {
			commitDescription = line
			break
		}
	}
	for _, line := range strings.Split(msg, "\n") {
		if strings.HasPrefix(line, "BREAKING CHANGE:") {
			major = true
		}
	}

	return typ, scope, commitDescription, major
}

var releaseCommitRegex = regexp.MustCompile(`^Release v(\d+).(\d+).(\d+)( \(.*\))?$`)

func DetectReleaseCommit(commit string, merge bool) (int, int, int) {
	candidates := []string{strings.SplitN(commit, "\n\n", 2)[0]}
	if merge {
		candidates = strings.Split(commit, "\n")
	}
	for _, candidate := range candidates {
		matches := releaseCommitRegex.FindStringSubmatch(candidate)
		if matches == nil {
			continue
		}
		major, _ := strconv.Atoi(matches[1])
		minor, _ := strconv.Atoi(matches[2])
		patch, _ := strconv.Atoi(matches[3])
		return major, minor, patch
	}
	return 0, 0, 0
}
