package metadata_test

import (
	"github.com/jt05610/loppu/metadata"
	"testing"
)

func TestVersion_Update(t *testing.T) {
	for _, tc := range []struct {
		Name   string
		Start  metadata.Version
		Type   metadata.VersionType
		Expect metadata.Version
	}{
		{"Minor", "1.2.3", metadata.Minor, "1.3.0"},
		{"Major", "1.2.3", metadata.Major, "2.0.0"},
		{"Patch", "1.2.3", metadata.Patch, "1.2.4"},
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
