package json_util_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/verify/json_util"
)

func TestSortJsonByKeys(t *testing.T) {
	tests := []struct {
		name string
		json string
		opt  *json_util.SortOptions
		want string
	}{
		{
			name: "empty",
			json: `{}`,
			opt:  nil,
			want: `{}`,
		},
		{
			name: "abc",
			json: `{"a":1,"b":2,"c":3}`,
			opt:  nil,
			want: `{"a":1,"b":2,"c":3}`,
		},
		{
			name: "cba",
			json: `{"c":1,"b":2,"a":3}`,
			opt:  nil,
			want: `{"a":3,"b":2,"c":1}`,
		},
		{
			name: "nested",
			json: `{"d":{"c":1,"a":2},"b":3}`,
			opt:  nil,
			want: `{"b":3,"d":{"a":2,"c":1}}`,
		},
		{
			name: "json_with_sig_field",
			json: `{"signature":"abc","d":{"c":1,"a":2,"agents":"abc"},"b":3}`,
			opt:  &json_util.SortOptions{NoSignProperties: true},
			want: `{"b":3,"d":{"a":2,"c":1}}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := json_util.SortJsonByKeys([]byte(tt.json), tt.opt)
			if err != nil {
				t.Errorf("SortJsonByKeys() error = %v", err)

				return
			}

			if string(got) != tt.want {
				t.Errorf("SortJsonByKeys() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
