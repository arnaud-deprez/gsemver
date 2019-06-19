package git

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