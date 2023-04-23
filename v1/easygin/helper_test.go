package easygin

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestSliceContains(t *testing.T) {
	tests := []struct {
		Needle         string
		Haystack       []string
		ExpectedResult bool
		Description    string
	}{
		{
			Needle:         "something",
			Haystack:       []string{"something", "something_else", "where"},
			ExpectedResult: true,
			Description:    "if the haystack has the needle , must be true",
		},
		{
			Needle:         "not_there",
			Haystack:       []string{"something", "something_else", "where"},
			ExpectedResult: false,
			Description:    "if the haystack does not have the needle , must be false",
		},
		{
			Needle:         "",
			Haystack:       []string{"something", "something_else", "where"},
			ExpectedResult: false,
			Description:    "if the needle is empty and haystack does not have an empty el , must be false",
		},
		{
			Needle:         "something",
			Haystack:       []string{},
			ExpectedResult: false,
			Description:    "if the needle is populated and haystack is empty, must be false",
		},
	}

	for _, v := range tests {
		assert.Equal(t, SliceContains(v.Haystack, v.Needle), v.ExpectedResult)
	}
}
