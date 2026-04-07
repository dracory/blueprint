package admin

import (
	"project/internal/links"
	"strings"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

// tableFileList creates the file list table HTML
func (c *FileManagerController) tableFileList(currentDirectory, parentDirectory string, directoryList, fileList []FileEntry) string {
	table := hb.Table().Class("table table-bordered table-striped").Children([]hb.TagInterface{
		hb.Thead().Children([]hb.TagInterface{
			hb.TR().Children([]hb.TagInterface{
				hb.TH().Style("width:30px;").Child(
					hb.Input().Type("checkbox").Class("form-check-input").Attr("onclick", "toggleSelectAll(this)"),
				),
				hb.TH().Style("width:1px;").Text(""),
				hb.TH().Text("Directory/File Name"),
				hb.TH().Style("width:100px;").Text("Size"),
				hb.TH().Style("width:100px;").Text("Modified"),
				hb.TH().Style("width:220px;").Text("Actions"),
			}),
		}),
		hb.Tbody().
			// Parent DIrectory
			ChildIfF(currentDirectory != "", func() hb.TagInterface {
				parentDirectoryURL := links.Admin().FileManager(map[string]string{"current_dir": parentDirectory})

				return hb.TR().Children([]hb.TagInterface{
					hb.TD(), // Empty checkbox cell for parent directory
					hb.TD().Children([]hb.TagInterface{
						hb.I().Class("bi bi-folder").Text(""),
					}),
					hb.TD().Children([]hb.TagInterface{
						hb.Hyperlink().Href(parentDirectoryURL).Children([]hb.TagInterface{
							hb.I().Class("bi bi-arrow-90deg-up").Text("").Style("margin-right: 5px;"),
							hb.Span().Text("parent"),
						}),
					}),
					hb.TD().Children([]hb.TagInterface{}),
					hb.TD().Children([]hb.TagInterface{}),
					hb.TD().Children([]hb.TagInterface{}),
				})
			}).
			// Directory List
			ChildIfF(len(directoryList) > 0, func() hb.TagInterface {
				return hb.Wrap().Children(lo.Map(directoryList, func(dir FileEntry, _ int) hb.TagInterface {
					name := dir.Name
					if dir.Name == "." || dir.Name == ".." {
						return nil
					}
					path := strings.TrimRight(dir.Path, "/")
					pathURL := links.Admin().FileManager(map[string]string{"current_dir": path})
					size := dir.SizeHuman

					buttonDelete := hb.Button().Class("btn btn-danger btn-sm").OnClick(`modalDirectoryDeleteShow('` + name + `')`).Children([]hb.TagInterface{
						hb.I().Class("bi bi-trash").Text("").Style("margin-right: 5px;"),
						hb.Span().Text("Delete"),
					})

					buttonRename := hb.Button().Class("btn btn-primary btn-sm").OnClick(`modalFileRenameShow('` + name + `')`).Children([]hb.TagInterface{
						hb.I().Class("bi bi-pencil").Text("").Style("margin-right: 5px;"),
						hb.Span().Text("Rename"),
					})

					return hb.TR().Children([]hb.TagInterface{
						hb.TD().Children([]hb.TagInterface{
							hb.Input().Type("checkbox").Class("form-check-input file-select").Attr("data-path", path).Attr("data-type", "directory"),
						}),
						hb.TD().Children([]hb.TagInterface{
							hb.I().Class("bi bi-folder").Text(""),
						}),
						hb.TD().Children([]hb.TagInterface{
							hb.Hyperlink().Href(pathURL).Children([]hb.TagInterface{
								hb.Span().Text(name).Style("font-weight: bold;"),
							}),
						}),
						hb.TD().Children([]hb.TagInterface{
							hb.Span().Text(size).Style("font-size: 12px;"),
						}),
						hb.TD().Children([]hb.TagInterface{
							hb.Span().Text("").Style("font-size: 11px;"),
						}),
						hb.TD().Children([]hb.TagInterface{
							buttonRename,
							buttonDelete,
						}),
					})
				}))
			}).
			// File List
			ChildIfF(len(fileList) > 0, func() hb.TagInterface {
				return hb.Wrap().Children(lo.Map(fileList, func(file FileEntry, _ int) hb.TagInterface {
					buttonDelete := hb.Button().Class("btn btn-danger btn-sm").OnClick(`modalFileDeleteShow('` + file.Name + `')`).Children([]hb.TagInterface{
						hb.I().Class("bi bi-trash").Text("").Style("margin-right: 5px;"),
						hb.Span().Text("Delete"),
					})

					buttonRename := hb.Button().Class("btn btn-primary btn-sm").OnClick(`modalFileRenameShow('` + file.Name + `')`).Children([]hb.TagInterface{
						hb.I().Class("bi bi-pencil").Text("").Style("margin-right: 5px;"),
						hb.Span().Text("Rename"),
					})

					buttonView := hb.Button().Class("btn btn-info btn-sm").OnClick(`modalFileViewShow('` + file.Name + `')`).Children([]hb.TagInterface{
						hb.I().Class("bi bi-eye").Text("").Style("margin-right: 5px;"),
						hb.Span().Text("View"),
					})

					buttonSelect := hb.Button().Class("btn btn-success btn-sm .btn-select").OnClick(`fileSelectedUrl('` + file.URL + `')`).Children([]hb.TagInterface{
						hb.I().Class("bi bi-chevron-right").Text("").Style("margin-right: 5px;"),
						hb.Span().Text("Select"),
					})

					return hb.TR().Children([]hb.TagInterface{
						hb.TD().Children([]hb.TagInterface{
							hb.Input().Type("checkbox").Class("form-check-input file-select").Attr("data-path", file.Path).Attr("data-type", "file"),
						}),
						hb.TD().Children([]hb.TagInterface{
							hb.I().Class("bi bi-file").Text(""),
						}),
						hb.TD().Children([]hb.TagInterface{
							hb.Span().Text(file.Name).Style("font-weight: bold;"),
							hb.Div().
								Children([]hb.TagInterface{
									hb.Span().Text("URL: "),
									hb.Hyperlink().
										Href(file.URL).
										Target("_blank").
										Children([]hb.TagInterface{
											hb.Span().Text(file.URL),
										}),
								}).
								Style("font-size: 12px;"),
						}),
						hb.TD().Children([]hb.TagInterface{
							hb.Span().Text(file.SizeHuman).Style("font-size: 12px;"),
						}),
						hb.TD().Children([]hb.TagInterface{
							hb.Span().
								HTML(lo.Substring(file.LastModifiedHuman, 0, 10)).
								Attr("title", file.LastModifiedHuman).
								Style("font-size: 11px;"),
						}),
						hb.TD().Children([]hb.TagInterface{
							buttonView,
							buttonRename,
							buttonDelete,
							buttonSelect,
						}),
					})
				}))
			}),
	})
	return table.ToHTML()
}
