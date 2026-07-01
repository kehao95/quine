package main

import "testing"

func TestMissionFromArgsOptional(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want string
	}{
		{name: "absent", args: nil, want: ""},
		{name: "blank", args: []string{"   "}, want: ""},
		{name: "words", args: []string{"answer", "the", "question"}, want: "answer the question"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := missionFromArgs(tt.args); got != tt.want {
				t.Fatalf("missionFromArgs(%#v) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}
