package git

// Commit data
type Commit struct {
	// Hash of the commit object.
	Hash Hash
	// Author is the original author of the commit.
	Author Signature
	// Committer is the one performing the commit.
	// It might be different from Author.
	Committer Signature
	// Message is the commit message, contains arbitrary text.
	Message string
}