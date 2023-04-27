package loppu_test

import (
	"github.com/jt05610/loppu"
	"testing"
)

func TestVersion_Update(t *testing.T) {
	for _, tc := range []struct {
		Name   string
		Start  loppu.Version
		Type   loppu.VersionType
		Expect loppu.Version
	}{
		{"Minor", "1.2.3", loppu.Minor, "1.3.0"},
		{"Major", "1.2.3", loppu.Major, "2.0.0"},
		{"Patch", "1.2.3", loppu.Patch, "1.2.4"},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			res := tc.Start.Update(tc.Type)
			if res != tc.Expect {
				t.Logf(
					`%s Fail
Expected: %s
Actual:   %s`, tc.Name, tc.Expect, res)
				t.Fail()
			}
		})
	}

}
