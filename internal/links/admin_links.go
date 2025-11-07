package links

import "github.com/samber/lo"

type adminLinks struct{}

func Admin() *adminLinks {
	return &adminLinks{}
}

func (l *adminLinks) Home(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_HOME, p)
}

func (l *adminLinks) Blog(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_BLOG, p)
}

// Cms is the cms old manager
func (l *adminLinks) Cms(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_CMS, p)
}

// CmsNew is the new cms manager
func (l *adminLinks) CmsNew(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_CMS_NEW, p)
}

// FileManager is the file manager
func (l *adminLinks) FileManager(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_FILE_MANAGER, p)
}

func (l *adminLinks) MediaManager(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_MEDIA, p)
}

func (l *adminLinks) Shop(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_SHOP, p)
}

func (l *adminLinks) Stats(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_STATS, p)
}

func (l *adminLinks) Tasks(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_TASKS, p)
}

func (l *adminLinks) Users(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_USERS, p)
}

func (l *adminLinks) UsersUserCreate(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_USERS_USER_CREATE, p)
}

func (l *adminLinks) UsersUserDelete(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_USERS_USER_DELETE, p)
}

func (l *adminLinks) UsersUserImpersonate(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_USERS_USER_IMPERSONATE, p)
}

func (l *adminLinks) UsersUserManager(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_USERS_USER_MANAGER, p)
}

func (l *adminLinks) UsersUserUpdate(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	return URL(ADMIN_USERS_USER_UPDATE, p)
}
