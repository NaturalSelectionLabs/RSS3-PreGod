package defers_test

import (
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/defers"
)

func TestDefers(t *testing.T) {
	t.Parallel()

	var str string

	type args struct {
		fns []func() error
	}

	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "register",
			args: args{
				fns: []func() error{
					func() error { str += "1,"; return nil }, //nolint:nlreturn // ignore
					func() error { str += "2,"; return nil }, //nolint:nlreturn // ignore
					func() error { str += "3,"; return nil }, //nolint:nlreturn // ignore
					func() error { str += "4,"; return nil }, //nolint:nlreturn // ignore
					nil,
				},
			},
			expected: "1,2,3,4,",
		},
	}

	//nolint:paralleltest // has to be sequential here
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defers.Register(tc.args.fns...)
			defers.Clean()
		})
	}
}
