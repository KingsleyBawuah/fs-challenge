package main

import "testing"

func Test_containsIssueCmd(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Notes without the command #issue",
			args: args{
				text: "Imagine a note with the word issue but not with a #",
			},
			want: false,
		},
		{
			name: "Note with the command #issue",
			args: args{
				text: "Imagine a note with some words and #issue",
			},
			want: true,
		},
		{
			name: "Note only the command #issue and no actual information",
			args: args{
				text: "#issue",
			},
			want: false,
		},
		{
			name: "Note with command #issue inside another word",
			args: args{
				text: "Imagine text he#issuere",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsIssueCmd(tt.args.text); got != tt.want {
				t.Errorf("containsIssueCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
