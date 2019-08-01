package version

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshallBumpStrategyType(t *testing.T) {
	assert := assert.New(t)

	testData := []struct {
		value    BumpStrategyType
		expected string
	}{
		{PATCH, `"PATCH"`},
		{MINOR, `"MINOR"`},
		{MAJOR, `"MAJOR"`},
		{AUTO, `"AUTO"`},
	}

	for _, tc := range testData {
		bytes, err := json.Marshal(tc.value)
		assert.Nil(err)
		assert.Equal(tc.expected, string(bytes))
	}
}

func TestUnMarshallBumpStrategyType(t *testing.T) {
	assert := assert.New(t)

	testData := []struct {
		value    string
		expected BumpStrategyType
	}{
		{`"PATCH"`, PATCH},
		{`"MINOR"`, MINOR},
		{`"MAJOR"`, MAJOR},
		{`"AUTO"`, AUTO},
		{`"patch"`, PATCH},
		{`"minor"`, MINOR},
		{`"major"`, MAJOR},
		{`"auto"`, AUTO},
		{`"foo"`, AUTO}, // fallback to AUTO if unknown value
	}

	for _, tc := range testData {
		var value BumpStrategyType
		err := json.Unmarshal([]byte(tc.value), &value)
		assert.Nil(err)
		assert.Equal(tc.expected, value)
	}
}
