package git

import (
	"fmt"
	"time"
)

// Signature is used to identify who and when created a commit or tag.
type Signature struct {
	// Name represents a person name. It is an arbitrary string.
	Name string
	// Email is an email, but it cannot be assumed to be well-formed.
	Email string
	// When is the timestamp of the signature.
	When time.Time
}

// GoString makes Signature satisfy the GoStringer interface.
func (s Signature) GoString() string {
	return fmt.Sprintf("git.Signature{Name: %q, Email: %q, When: %q}", s.Name, s.Email, s.When)
}
