package metadata

import "os/user"

// Username gets the current user's full name or returns a blank string in
// case of error.
func Username() string {
	uName := ""
	u, err := user.Current()
	if err == nil {
		uName = u.Name
	}
	return uName
}
