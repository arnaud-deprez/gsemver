package git

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashString(t *testing.T) {
	hashValue := "12345678901234567890"
	hash := Hash(hashValue)
	fmt.Println(hash[:])
	assert.Equal(t, hashValue, hash.String())
}

func TestHashShortString(t *testing.T) {
	hashValue := "12345678901234567890"
	hash := Hash(hashValue)
	fmt.Println(hash[:])
	assert.Equal(t, "1234567", hash.Short().String())
}
