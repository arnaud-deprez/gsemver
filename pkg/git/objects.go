package git

import (
	"encoding/hex"
	"time"
)

// Hash of commit
type Hash [20]byte

// NewHash return a new Hash from a hexadecimal hash representation
func NewHash(s string) Hash {
	b, _ := hex.DecodeString(s)

	var h Hash
	copy(h[:], b)

	return h
}

// Signature is used to identify who and when created a commit or tag.
type Signature struct {
	// Name represents a person name. It is an arbitrary string.
	Name string
	// Email is an email, but it cannot be assumed to be well-formed.
	Email string
	// When is the timestamp of the signature.
	When time.Time
}

// Tag is data of git-tag
type Tag struct {
	// Hash of the tag.
	Hash Hash
	// Name of the tag.
	Name string
	// Tagger is the one who created the tag.
	Tagger Signature
	// Message is an arbitrary text message.
	Message string
}

// Commit data
type Commit struct {
	// Hash of the commit object.
	Hash Hash
	// Author is the original author of the commit.
	Author Signature
	// Committer is the one performing the commit, might be different from
	// Author.
	Committer Signature
	// Message is the commit message, contains arbitrary text.
	Message string
}
