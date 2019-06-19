package git

// Hash of commit
type Hash string

// String return the string representation
func (h Hash) String() string {
	return string(h)
}

// Short convert the 7 first bytes of the Hash to String
func (h Hash) Short() Hash {
	return h[:7]
}
