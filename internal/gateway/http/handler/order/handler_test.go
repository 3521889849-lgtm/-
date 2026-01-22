package order

import "testing"

func TestIsConnRefused(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"connectex", errString("dial tcp 127.0.0.1:8890: connectex: No connection could be made because the target machine actively refused it."), true},
		{"refused", errString("dial tcp 127.0.0.1:8890: connection refused"), true},
		{"other", errString("timeout"), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isConnRefused(c.err); got != c.want {
				t.Fatalf("got=%v want=%v", got, c.want)
			}
		})
	}
}

type errString string

func (e errString) Error() string { return string(e) }

