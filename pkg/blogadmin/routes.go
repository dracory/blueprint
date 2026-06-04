package blogadmin

import (
	"errors"
	"net/http"
	"project/internal/links"
	"project/internal/app"

	"github.com/dracory/req"
	"github.com/dracory/rtr"

	aiPostContentUpdate "project/pkg/blogadmin/ai_post_content_update"
	aiPostEditor "project/pkg/blogadmin/ai_post_editor"
	aiPostGenerator "project/pkg/blogadmin/ai_post_generator"
	aiTest "project/pkg/blogadmin/ai_test"
	aiTitleGenerator "project/pkg/blogadmin/ai_title_generator"
	aiTools "project/pkg/blogadmin/ai_tools"
	blogSettings "project/pkg/blogadmin/blog_settings"
	"project/pkg/blogadmin/category_manager"
	"project/pkg/blogadmin/dashboard"
	"project/pkg/blogadmin/post_create"
	"project/pkg/blogadmin/post_delete"
	"project/pkg/blogadmin/post_manager"
	"project/pkg/blogadmin/post_update"
	"project/pkg/blogadmin/shared"
	"project/pkg/blogadmin/tag_manager"
)

func Routes(app app.AppInterface, opts ...AdminOptions) ([]rtr.RouteInterface, error) {
	_ = opts // Options available for future use
	if app == nil {
		return nil, errors.New("app cannot be nil")
	}
	handler := func(w http.ResponseWriter, r *http.Request) string {
		controller := req.GetStringTrimmed(r, "controller")

		switch controller {
		case shared.CONTROLLER_HOME:
			return post_manager.NewPostManagerController(app).Handler(w, r)
		case shared.CONTROLLER_POST_CREATE:
			return post_create.NewPostCreateController(app).Handler(w, r)
		case shared.CONTROLLER_POST_DELETE:
			return post_delete.NewPostDeleteController(app).Handler(w, r)
		case shared.CONTROLLER_POST_MANAGER:
			return post_manager.NewPostManagerController(app).Handler(w, r)
		case shared.CONTROLLER_POST_UPDATE:
			return post_update.NewPostUpdateController(app).Handler(w, r)
		case shared.CONTROLLER_AI_TOOLS:
			return aiTools.NewAiToolsController(app).Handler(w, r)
		case shared.CONTROLLER_BLOG_SETTINGS:
			return blogSettings.NewBlogSettingsController(app).Handler(w, r)
		case shared.CONTROLLER_AI_POST_CONTENT_UPDATE:
			return aiPostContentUpdate.NewController(app).Handler(w, r)
		case shared.CONTROLLER_AI_POST_GENERATOR:
			return aiPostGenerator.NewAiPostGeneratorController(app).Handler(w, r)
		case shared.CONTROLLER_AI_TITLE_GENERATOR:
			return aiTitleGenerator.NewAiTitleGeneratorController(app).Handler(w, r)
		case shared.CONTROLLER_AI_POST_EDITOR:
			return aiPostEditor.NewAiPostEditorController(app).Handler(w, r)
		case shared.CONTROLLER_AI_TEST:
			return aiTest.NewAiTestController(app).Handler(w, r)
		case shared.CONTROLLER_DASHBOARD:
			return dashboard.NewDashboardController(app).Handler(w, r)
		case shared.CONTROLLER_CATEGORY_MANAGER:
			return category_manager.NewCategoryManagerController(app).Handler(w, r)
		case shared.CONTROLLER_TAG_MANAGER:
			return tag_manager.NewTagManagerController(app).Handler(w, r)
		}

		// Default to dashboard
		return dashboard.NewDashboardController(app).Handler(w, r)
	}

	blog := rtr.NewRoute().
		SetName("Admin > Blog").
		SetPath(links.ADMIN_BLOG).
		SetHTMLHandler(handler)

	blogCatchAll := rtr.NewRoute().
		SetName("Admin > Blog > Catchall").
		SetPath(links.ADMIN_BLOG + links.CATCHALL).
		SetHTMLHandler(handler)

	return []rtr.RouteInterface{
		blog,
		blogCatchAll,
	}, nil
}
