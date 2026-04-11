package post_update

import (
	"context"
	"errors"

	"project/internal/registry"

	"github.com/dracory/blogstore"
	"github.com/dracory/sb"
	"github.com/dracory/versionstore"
	"github.com/samber/lo"
)

func createPostVersioning(ctx context.Context, registry registry.RegistryInterface, post blogstore.PostInterface) error {
	if registry == nil || registry.GetBlogStore() == nil {
		return errors.New("blog store not available")
	}

	if post == nil {
		return errors.New("post is nil")
	}

	if !registry.GetBlogStore().VersioningEnabled() {
		return nil
	}

	lastVersioningList, err := registry.GetBlogStore().VersioningList(ctx, blogstore.NewVersioningQuery().
		SetEntityType(blogstore.VERSIONING_TYPE_POST).
		SetEntityID(post.GetID()).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder(sb.DESC).
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

	return registry.GetBlogStore().VersioningCreate(ctx, blogstore.NewVersioning().
		SetEntityID(post.GetID()).
		SetEntityType(blogstore.VERSIONING_TYPE_POST).
		SetContent(content))
}

