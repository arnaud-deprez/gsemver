package git

import (
	"strconv"
	"strings"
	"time"

	"github.com/arnaud-deprez/gsemver/pkg/git"
)

var (
	// constants
	separator = "-->8--"
	delimiter = "$_$"

	// fields
	hashField      = "HASH"
	authorField    = "AUTHOR"
	committerField = "COMMITTER"
	messageField   = "MESSAGE"
	// subjectField   = "SUBJECT"
	// bodyField      = "BODY"

	// formats
	hashFormat      = hashField + ":%H"
	authorFormat    = authorField + ":%an\t%ae\t%at"
	committerFormat = committerField + ":%cn\t%ce\t%ct"
	messageFormat   = messageField + ":%B"

	// log
	logFormat = separator + strings.Join([]string{
		hashFormat,
		authorFormat,
		committerFormat,
		messageFormat,
	}, delimiter)
)

type commitParser struct {
	logFormat string
}

func (p *commitParser) Parse(out string) []git.Commit {
	if p == nil {
		p = &commitParser{
			logFormat: logFormat,
		}
	}

	lines := strings.Split(out, separator)
	lines = lines[1:]
	commits := make([]git.Commit, len(lines))

	for i, line := range lines {
		commit := p.parseCommit(line)
		commits[i] = *commit
	}

	return commits
}

func (p *commitParser) parseCommit(input string) *git.Commit {
	commit := &git.Commit{}
	tokens := strings.Split(input, delimiter)

	for _, token := range tokens {
		firstSep := strings.Index(token, ":")
		field := token[0:firstSep]
		value := strings.TrimSpace(token[firstSep+1:])

		switch field {
		case hashField:
			commit.Hash = git.NewHash(value)
		case authorField:
			commit.Author = p.parseSignature(value)
		case committerField:
			commit.Committer = p.parseSignature(value)
		case messageField:
			commit.Message = value
		}
	}

	return commit
}

func (p *commitParser) parseSignature(input string) git.Signature {
	arr := strings.Split(input, "\t")
	ts, err := strconv.Atoi(arr[2])
	if err != nil {
		ts = 0
	}

	return git.Signature{
		Name:  arr[0],
		Email: arr[1],
		When:  time.Unix(int64(ts), 0),
	}
}
