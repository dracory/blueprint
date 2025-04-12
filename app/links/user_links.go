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
	p := lo.IfF(len(params) > 0, func() map[string]string { return params[0] }).Else(map[string]string{})
	return URL(USER_HOME, p)
}

func (l *userLinks) Profile(params ...map[string]string) string {
	p := lo.IfF(len(params) > 0, func() map[string]string { return params[0] }).Else(map[string]string{})
	return URL(USER_PROFILE, p)
}

func (l *userLinks) ProfileSave() string {
	return URL(USER_PROFILE_UPDATE, nil)
}
