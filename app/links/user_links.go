package links

type userLinks struct{}

// User is a shortcut for NewUserLinks
func User() *userLinks {
	return NewUserLinks()
}

// Deprecated: Use User() instead. NewUserLinks will be removed in the next major version.
func NewUserLinks() *userLinks {
	return &userLinks{}
}

func (l *userLinks) Home(params map[string]string) string {
	return URL(USER_HOME, params)
}

func (l *userLinks) Profile(params map[string]string) string {
	return URL(USER_PROFILE, params)
}

func (l *userLinks) ProfileSave() string {
	return URL(USER_PROFILE_UPDATE, nil)
}
