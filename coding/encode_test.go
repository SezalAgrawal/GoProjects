package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "case1",
			input: "aaa",
			want:  "3a",
		},
		{
			name:  "case2",
			input: "aaabaaacd",
			want:  "3ab3acd",
		},
		{
			name:  "case3",
			input: "",
			want:  "",
		},
		{
			name:  "case4",
			input: "###$$",
			want:  "3#2$",
		},
		// {
		// 	name:  "case4",
		// 	input: "333",
		// 	want:  "33",
		// },
		// {
		// 	name:  "case4",
		// 	input: "333",
		// 	want:  "33",
		// },
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := encode(c.input)
			assert.Equal(t, c.want, got)
		})
	}
}

func TestDecode(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "case1",
			input: "3a",
			want:  "aaa",
		},
		{
			name:  "case2",
			input: "3ab3acd",
			want:  "aaabaaacd",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := decode(c.input)
			assert.Equal(t, c.want, got)
		})
	}
}
