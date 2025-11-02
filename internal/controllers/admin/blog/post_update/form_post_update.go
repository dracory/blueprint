package post_update

import (
	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogblocks"

	"github.com/dracory/blockeditor"
	"github.com/dracory/blogstore"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func formPostUpdate(data postUpdateControllerData) hb.TagInterface {

	formPostUpdate := form.NewForm(form.FormOptions{
		ID: "FormPostUpdate",
	})

	if data.view == VIEW_DETAILS {
		formPostUpdate.SetFields(fieldsDetails(data))
	}

	if data.view == VIEW_CONTENT {
		fieldsContent, err := fieldsContent(data)
		if err != nil {
			return hb.Div().Class("alert alert-danger").Text("Error creating blockeditor: ").Text(err.Error())
		}
		formPostUpdate.SetFields(fieldsContent)
	}

	if data.view == VIEW_SEO {
		formPostUpdate.SetFields(fieldsSEO(data))
	}

	if data.formErrorMessage != "" {
		formPostUpdate.AddField(form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "error", Text: data.formErrorMessage}).ToHTML(),
		}))
	}

	if data.formSuccessMessage != "" {
		formPostUpdate.AddField(form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "success", Text: data.formSuccessMessage}).ToHTML(),
		}))
	}

	return formPostUpdate.Build()
}

func fieldsContent(data postUpdateControllerData) ([]form.FieldInterface, error) {
	editor := lo.IfF(data.post != nil, func() string { return data.post.Editor() }).Else("")

	fieldContent := form.NewField(form.FieldOptions{
		Label:   "Content",
		Name:    "post_content",
		Type:    form.FORM_FIELD_TYPE_TEXTAREA,
		Value:   data.formContent,
		Help:    "The content of this blog post to display on the post details page.",
		Options: []form.FieldOption{},
	})

	// For HTML Area editor, configure the Trumbowyg editor
	if editor == blogstore.POST_EDITOR_HTMLAREA {
		htmlAreaFieldOptions := []form.FieldOption{
			{
				Key: "config",
				Value: `{
	btns: [
		['viewHTML'],
		['undo', 'redo'],
		['formatting'],
		['strong', 'em', 'del'],
		['superscript', 'subscript'],
		['link','justifyLeft','justifyRight','justifyCenter','justifyFull'],
		['unorderedList', 'orderedList'],
		['insertImage'],
		['removeformat'],
		['horizontalRule'],
		['fullscreen'],
	],
	autogrow: true,
	removeformatPasted: true,
	tagsToRemove: ['script', 'link', 'embed', 'iframe', 'input'],
	tagsToKeep: ['hr', 'img', 'i'],
	autogrowOnEnter: true,
	linkTargets: ['_blank'],
	}`,
			}}
		fieldContent.Type = form.FORM_FIELD_TYPE_HTMLAREA
		fieldContent.Options = htmlAreaFieldOptions
	}

	if editor == blogstore.POST_EDITOR_BLOCKEDITOR {
		value := fieldContent.Value

		if value == "" {
			value = `[]`
		}

		editor, err := blockeditor.NewEditor(blockeditor.NewEditorOptions{
			// ID:    "blockeditor" + uid.HumanUid(),
			Name:  fieldContent.Name,
			Value: value,
			HandleEndpoint: shared.NewLinks().PostUpdate(map[string]string{
				"post_id": data.postID,
				"action":  ACTION_BLOCKEDITOR_HANDLE,
			}),
			BlockDefinitions: blogblocks.BlockEditorDefinitions(),
		})

		if err != nil {
			return nil, err
		}

		fieldContent.Type = form.FORM_FIELD_TYPE_BLOCKEDITOR
		fieldContent.CustomInput = editor
	}

	fieldsContent := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Title",
			Name:  "post_title",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formTitle,
			Help:  "The title of this blog as will be seen everywhere",
		}),
		form.NewField(form.FieldOptions{
			Label: "Summary",
			Name:  "post_summary",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formSummary,
			Help:  "A short summary of this blog post to display on the post listing page.",
		}),
		fieldContent,
		form.NewField(form.FieldOptions{
			Label:    "Post ID",
			Name:     "post_id",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.postID,
			Readonly: true,
		}),
		form.NewField(form.FieldOptions{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    VIEW_CONTENT,
			Readonly: true,
		}),
	}

	if editor == blogstore.POST_EDITOR_BLOCKAREA {
		contentScript := hb.Script(`
setTimeout(() => {
	const textArea = document.querySelector('textarea[name="post_content"]');
	// get ID
	console.log(textArea.id);
	const blockArea = new BlockArea(textArea.id);
	blockArea.setParentId('` + data.postID + `');
	blockArea.registerBlock(BlockAreaHeading);
	blockArea.registerBlock(BlockAreaText);
	blockArea.registerBlock(BlockAreaImage);
	blockArea.registerBlock(BlockAreaCode);
	blockArea.registerBlock(BlockAreaRawHtml);
	blockArea.init();
}, 2000)
			`).
			ToHTML()

		fieldsContent = append(fieldsContent, form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: contentScript,
		}))
	}

	if editor == blogstore.POST_EDITOR_MARKDOWN {
		contentScript := hb.Script(`
setTimeout(() => {
	const textArea = document.querySelector('textarea[name="post_content"]');
	textArea.style.height = '300px';
}, 2000)
			`).
			ToHTML()

		fieldsContent = append(fieldsContent, form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: contentScript,
		}))
	}

	return fieldsContent, nil
}

func fieldsDetails(data postUpdateControllerData) []form.FieldInterface {
	fieldsDetails := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Status",
			Name:  "post_status",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formStatus,
			Help:  "The status of this blog post. Published posts will be displayed on the website.",
			Options: []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
				{
					Value: "Draft",
					Key:   blogstore.POST_STATUS_DRAFT,
				},
				{
					Value: "Published",
					Key:   blogstore.POST_STATUS_PUBLISHED,
				},
				{
					Value: "Unpublished",
					Key:   blogstore.POST_STATUS_UNPUBLISHED,
				},
				{
					Value: "In Trash Bin",
					Key:   blogstore.POST_STATUS_TRASH,
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Image URL",
			Name:  "post_image_url",
			Type:  form.FORM_FIELD_TYPE_IMAGE,
			Value: data.formImageUrl,
			Help:  "The image that will be displayed on the blog post. If left empty, the default image will be used.",
		}),
		form.NewField(form.FieldOptions{
			Label: "Featured",
			Name:  "post_featured",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formFeatured,
			Help:  "Is this blog post featured? Featured posts will be displayed on the home page.",
			Options: []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
				{
					Value: "No",
					Key:   "no",
				},
				{
					Value: "Yes",
					Key:   "yes",
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Published At",
			Name:  "post_published_at",
			Type:  form.FORM_FIELD_TYPE_DATETIME,
			Value: data.formPublishedAt,
			Help:  "The date this blog post was published.",
		}),
		form.NewField(form.FieldOptions{
			Label: "Editor",
			Name:  "post_editor",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formEditor,
			Help:  "The editor that will be used while editing this blogpost. Once set, this should not be changed, or the content may be lost.",
			Options: []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
				{
					Value: "BlockArea",
					Key:   blogstore.POST_EDITOR_BLOCKAREA,
				},
				{
					Value: "BlockEditor",
					Key:   blogstore.POST_EDITOR_BLOCKEDITOR,
				},
				{
					Value: "Markdown",
					Key:   blogstore.POST_EDITOR_MARKDOWN,
				},
				{
					Value: "HTML Area (WYSIWYG)",
					Key:   blogstore.POST_EDITOR_HTMLAREA,
				},
				{
					Value: "Text Area",
					Key:   blogstore.POST_EDITOR_TEXTAREA,
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Admin Notes",
			Name:  "post_memo",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: data.formMemo,
			Help:  "Admin notes for this blogpost. These notes will not be visible to the public.",
		}),
		form.NewField(form.FieldOptions{
			Label:    "Post ID",
			Name:     "post_id",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    data.postID,
			Readonly: true,
			Help:     "The reference number (ID) of the post.",
		}),
		form.NewField(form.FieldOptions{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.view,
			Readonly: true,
		}),
	}

	return fieldsDetails
}

func fieldsSEO(data postUpdateControllerData) []form.FieldInterface {
	fieldsSEO := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Meta Description",
			Name:  "post_meta_description",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formMetaDescription,
			Help:  "The description of this blog as will be seen in search engines.",
		}),
		form.NewField(form.FieldOptions{
			Label: "Meta Keywords",
			Name:  "post_meta_keywords",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formMetaKeywords,
			Help:  "Specifies the keywords that will be used by the search engines to find this blog. Separate keywords with commas.",
		}),
		form.NewField(form.FieldOptions{
			Label: "Meta Robots",
			Name:  "post_meta_robots",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formMetaRobots,
			Help:  "Specifies if this blog should be indexed by the search engines. Index, Follow, means all. NoIndex, NoFollow means none.",
			Options: []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
				{
					Value: "INDEX, FOLLOW",
					Key:   "INDEX, FOLLOW",
				},
				{
					Value: "NOINDEX, FOLLOW",
					Key:   "NOINDEX, FOLLOW",
				},
				{
					Value: "INDEX, NOFOLLOW",
					Key:   "INDEX, NOFOLLOW",
				},
				{
					Value: "NOINDEX, NOFOLLOW",
					Key:   "NOINDEX, NOFOLLOW",
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Canonical URL",
			Name:  "post_canonical_url",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formCanonicalURL,
			Help:  "The canonical URL for this blog. This is used by the search engines to display the preferred version of the blog post in search results.",
		}),
		form.NewField(form.FieldOptions{
			Label:    "Post ID",
			Name:     "post_id",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    data.postID,
			Readonly: true,
		}),
		form.NewField(form.FieldOptions{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    VIEW_SEO,
			Readonly: true,
		}),
	}

	return fieldsSEO
}
