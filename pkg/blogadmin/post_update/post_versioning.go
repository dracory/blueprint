package post_update

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/blogstore"
	"github.com/dracory/neat"
	"github.com/dracory/versionstore"
	"github.com/samber/lo"
)

func createPostVersioning(ctx context.Context, app app.AppInterface, post blogstore.PostInterface) error {
	if app == nil {
		return errors.New("blog store not available")
	}
	if app.GetBlogStore() == nil {
		return errors.New("blog store not available")
	}

	if post == nil {
		return errors.New("post is nil")
	}

	if !app.GetBlogStore().VersioningEnabled() {
		return nil
	}

	lastVersioningList, err := app.GetBlogStore().VersioningList(ctx, blogstore.NewVersioningQuery().
		SetEntityType(blogstore.VERSIONING_TYPE_POST).
		SetEntityID(post.GetID()).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder(neat.SortDesc).
		SetLimit(1))
	if err != nil {
		return err
	}

	content, err := post.MarshalToVersioning()
	if err != nil {
		return err
	}

	lastVersioning := lo.IfF[blogstore.VersioningInterface](len(lastVersioningList) > 0, func() blogstore.VersioningInterface {
		return lastVersioningList[0]
	}).ElseF(func() blogstore.VersioningInterface {
		return nil
	})
	if lastVersioning != nil {
		if lastVersioning.Content() == content {
			return nil
		}
	}

	return app.GetBlogStore().VersioningCreate(ctx, blogstore.NewVersioning().
		SetEntityID(post.GetID()).
		SetEntityType(blogstore.VERSIONING_TYPE_POST).
		SetContent(content))
}
