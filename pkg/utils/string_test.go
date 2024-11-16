package utils

import (
	"testing"
)

func TestRemoveAccents(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "Case 1",
			input:  "HuyNd đã Test  trường hợp này",
			output: "huynd da test truong hop nay",
		},
		{
			name:   "Case 2",
			input:  "cộng hòa xã hội chủ nghĩa việt nam % * \\",
			output: "cong hoa xa hoi chu nghia viet nam %",
		},
		{
			name:   "Case 2",
			input:  "tag 2 222",
			output: "tag 2 222",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveAccents(tt.input); got != tt.output {
				t.Errorf("RemoveAccents() = %v, want %v", got, tt.output)
			}
		})
	}

}

func TestRemoveDuplicatedWhiteSpace(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "Case 1",
			input:  "đồng hồ  thời   trang",
			output: "đồng hồ thời trang",
		},
		{
			name:   "Case 2",
			input:  "đồng hồ  thời    trang",
			output: "đồng hồ thời trang",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicatedWhiteSpace(tt.input); got != tt.output {
				t.Errorf("RemoveDuplicatedWhiteSpace() = %v, want %v", got, tt.output)
			}
		})
	}
}
