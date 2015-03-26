package test_tools

import (
    "fmt"
    "sort"
    "testing"

    "github.com/stretchr/testify/assert"
)

type Json map[string]interface{}
type JsonList []Json

func (j Json) keys() []string {
    keys := make([]string, 0)

    for key := range j {
        keys = append(keys, key)
    }
    return keys
}

func (actual Json) Equal(t *testing.T, expected Json) {
    expectedKeys := sort.StringSlice(expected.keys())
    expectedKeys.Sort()
    actualKeys := sort.StringSlice(actual.keys())
    actualKeys.Sort()

    assert.Equal(t, expectedKeys, actualKeys)
    for _, key := range expectedKeys {
        assert.Equal(t, expected[key], actual[key], fmt.Sprintf("Comparing key %s", key))
    }
}

func (actualList JsonList) Equal(t *testing.T, expectedList JsonList) {
    assert.Equal(t, len(expectedList), len(actualList))
    for i, expected := range expectedList {
        actualList[i].Equal(t, expected)
    }
}
