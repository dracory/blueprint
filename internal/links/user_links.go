package links

import "github.com/samber/lo"

type userLinks struct{}

// User is a shortcut for NewUserLinks
func User() *userLinks {
	return NewUserLinks()
}

// Deprecated: Use User() instead. NewUserLinks will be removed in the next major version.
func NewUserLinks() *userLinks {
	return &userLinks{}
}

func (l *userLinks) Home(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(USER_HOME, p)
}

func (l *userLinks) Profile(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(USER_PROFILE, p)
}
