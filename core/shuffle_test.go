package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuffle(t *testing.T) {
	as := assert.New(t)
	cases := []struct {
		name   string
		input  []string
		expect []string
	}{
		{
			"range",
			[]string{"1-3"},
			[]string{"1", "2", "3"},
		},
		{
			"negative range",
			[]string{"-1-1"},
			[]string{"-1", "0", "1"},
		},
		{
			"reversed range",
			[]string{"-2--3"},
			[]string{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := extractArgs(c.input)
			as.NoError(err)
			as.Equal(c.expect, result)
		})
	}
}
