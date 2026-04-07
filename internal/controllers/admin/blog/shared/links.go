package shared

import (
	"project/internal/links"

	"github.com/samber/lo"
)

type Links struct{}

func NewLinks() *Links {
	return &Links{}
}

func (*Links) Home(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_POST_MANAGER
	return links.Admin().Blog(p)
}

func (*Links) PostCreate(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_POST_CREATE
	return links.Admin().Blog(p)
}

func (*Links) PostDelete(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_POST_DELETE
	return links.Admin().Blog(p)
}

func (*Links) PostManager(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_POST_MANAGER
	return links.Admin().Blog(p)
}

func (*Links) PostUpdate(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_POST_UPDATE
	return links.Admin().Blog(p)
}

func (*Links) PostUpdateV1(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_POST_UPDATE_V1
	return links.Admin().Blog(p)
}

func (*Links) BlogSettings(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_BLOG_SETTINGS
	return links.Admin().Blog(p)
}

func (*Links) AiTools(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_AI_TOOLS
	return links.Admin().Blog(p)
}

func (*Links) AiPostContentUpdate(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_AI_POST_CONTENT_UPDATE
	return links.Admin().Blog(p)
}

func (*Links) AiPostGenerator(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_AI_POST_GENERATOR
	return links.Admin().Blog(p)
}

func (*Links) AiTitleGenerator(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_AI_TITLE_GENERATOR
	return links.Admin().Blog(p)
}

func (*Links) AiPostEditor(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_AI_POST_EDITOR
	return links.Admin().Blog(p)
}

func (*Links) AiTest(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_AI_TEST
	return links.Admin().Blog(p)
}

func (*Links) Dashboard(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_DASHBOARD
	return links.Admin().Blog(p)
}

func (*Links) CategoryManager(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_CATEGORY_MANAGER
	return links.Admin().Blog(p)
}

func (*Links) CategoryCreate(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_CATEGORY_CREATE
	return links.Admin().Blog(p)
}

func (*Links) CategoryUpdate(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_CATEGORY_UPDATE
	return links.Admin().Blog(p)
}

func (*Links) CategoryDelete(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_CATEGORY_DELETE
	return links.Admin().Blog(p)
}

func (*Links) TagManager(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_TAG_MANAGER
	return links.Admin().Blog(p)
}

func (*Links) TagCreate(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_TAG_CREATE
	return links.Admin().Blog(p)
}

func (*Links) TagUpdate(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_TAG_UPDATE
	return links.Admin().Blog(p)
}

func (*Links) TagDelete(params ...map[string]string) string {
	p := lo.FirstOr(params, map[string]string{})
	p["controller"] = CONTROLLER_TAG_DELETE
	return links.Admin().Blog(p)
}
