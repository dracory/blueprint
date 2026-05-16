package taskadmin

import (
	"context"
	"net/http"
	"project/internal/links"

	"strconv"

	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dracory/taskstore"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func newPageQueueManager(admin *admin) *pageQueueManager {
	return &pageQueueManager{
		admin: admin,
	}
}

type pageQueueManager struct {
	admin *admin
}

func (p *pageQueueManager) Handler(w http.ResponseWriter, r *http.Request) {
	// _ = r.Context().Value(keyEndpoint).(string)
	data, errorMessage := p.pageQueueManagerPrepareData(r)

	if errorMessage != "" {
		p.admin.layout(w, LayoutOptions{
			Title:   "Home",
			Content: errorMessage,
		})
	}

	if data.action == ActionModalQueuedTaskDeleteShow {
		content := p.admin.onModalQueuedTaskDeleteShow(r)
		w.Write([]byte(content))
		return
	}

	if data.action == ActionModalQueuedTaskDeleteSubmitted {
		content := p.admin.onModalQueuedTaskDeleteSubmitted(r)
		w.Write([]byte(content))
		return
	}

	if data.action == ActionModalQueuedTaskEnqueueShow {
		content := p.admin.taskEnqueueModalShow(r)
		w.Write([]byte(content))
		return
	}

	if data.action == ActionModalQueuedTaskEnqueueSubmitted {
		content := p.admin.onModaltaskEnqueueSubmitted(r)
		w.Write([]byte(content))
		return
	}

	if data.action == ActionModalQueuedTaskDetailsShow {
		content := p.admin.onModalQueuedTaskDetailsShow(data.queueID)
		w.Write([]byte(content))
		return
	}

	if data.action == ActionModalQueuedTaskFilterShow {
		content := p.onModalQueuedTaskFilterShow(data)
		w.Write([]byte(content))
		return
	}

	if data.action == ActionModalQueuedTaskParametersShow {
		content := p.admin.onModalQueuedTaskParametersShow(data.queueID)
		w.Write([]byte(content))
		return
	}

	if data.action == ActionModalQueuedTaskRequeueShow {
		content := p.admin.onModalQueuedTaskRequeueShow(data.queueID)
		w.Write([]byte(content))
		return
	}

	if data.action == ActionModalQueuedTaskRequeueSubmitted {
		content := p.admin.onModalQueuedTaskRequeueSubmitted(r)
		w.Write([]byte(content))
		return
	}

	if data.action == "queue-task-restart" {
		content := p.queueTaskRestart(w, r, data)
		w.Write([]byte(content))
		return
	}

	buttonTaskEnqueue := hb.NewButton().
		Text("Enqueue Task").
		Class("btn btn-primary float-end").
		HxPost(link(data.request, PathHome, map[string]string{
			"action": ActionModalQueuedTaskRequeueShow,
			"page":   data.page,
			"by":     data.sortBy,
			"sort":   data.sortOrder,
		})).
		HxTarget("#QueuedTasksListTable").
		HxSwap("outerHTML")

	title := hb.NewHeading1().
		Text("Queue Manager").
		Child(buttonTaskEnqueue)

	breadcrumbs := bs.Breadcrumbs([]bs.Breadcrumb{
		{
			Name: "Home",
			URL:  links.Admin().Home(map[string]string{}),
		},
		{
			Name: "Tasks",
			URL:  links.Admin().Tasks(map[string]string{}),
		},
		{
			Name: "Task Manager",
			URL:  links.Admin().Tasks(map[string]string{}),
		},
	})

	webpage := hb.NewWrap().Children([]hb.TagInterface{
		hb.NewSection().Children([]hb.TagInterface{
			hb.NewDiv().Class("container").
				Style("padding:40px;").
				Children([]hb.TagInterface{
					title,
					hb.NewDiv().
						Style("margin-top:20px;margin-bottom:20px;").
						Children([]hb.TagInterface{
							breadcrumbs,
						}),
					p.table(data),
				}),
		}),
	})

	// content := p.table(data).ToHTML()
	p.admin.layout(w, LayoutOptions{
		Title:   "Home",
		Content: webpage.ToHTML(),
	})
}

type pageQueueManagerData struct {
	request   http.Request
	action    string
	page      string
	sortOrder string
	sortBy    string
	status    string
	search    string
	dateFrom  string
	dateTo    string
	queueID   string
	taskID    string
}

func (p *pageQueueManager) queueTaskRestart(_ http.ResponseWriter, _ *http.Request, data pageQueueManagerData) string {
	if data.queueID == "" {
		return hb.NewDiv().Class("alert alert-danger").Text("queue id is required").ToHTML()
	}

	queue, err := p.admin.taskStore.TaskQueueFindByID(context.Background(), data.queueID)

	if err != nil {
		p.admin.logger.Error("At taskadmin > queueTaskRestart", "error", err.Error())
		return hb.NewDiv().Class("alert alert-danger").Text("Task failed to be deleted").ToHTML()
	}

	if queue == nil {
		return hb.NewDiv().Class("alert alert-danger").Text("queue not found").ToHTML()
	}

	queue.SetStatus(taskstore.TaskQueueStatusQueued)

	err = p.admin.taskStore.TaskQueueUpdate(context.Background(), queue)

	if err != nil {
		p.admin.logger.Error("At taskadmin > queueTaskRestart", "error", err.Error())
		return hb.NewDiv().Class("alert alert-danger").Text("Task failed to be queued").ToHTML()
	}

	return p.table(data).ToHTML()
}

func (p *pageQueueManager) table(data pageQueueManagerData) *hb.Tag {
	allTasks, err := p.admin.taskStore.TaskDefinitionList(context.Background(), nil)

	if err != nil {
		p.admin.logger.Error("At taskadmin > table", "error", err.Error())
		return hb.NewDiv().Text("error retrieving tasks")
	}

	pageInt, _ := strconv.Atoi(data.page)
	perPage := 20

	queuedTaskList, err := p.admin.taskStore.TaskQueueList(context.Background(), nil)

	if err != nil {
		p.admin.logger.Error("At taskadmin > table", "error", err.Error())
		return hb.NewDiv().Text("error retrieving queued tasks")
	}

	queuedTaskCount, errCount := p.admin.taskStore.TaskQueueCount(context.Background(), nil)

	if errCount != nil {
		p.admin.logger.Error("At taskadmin > table", "error", errCount.Error())
		return hb.NewDiv().Text("error retrieving queued tasks count")
	}

	table := hb.NewTable().
		Class("table table-striped table-hover table-bordered").
		Children([]hb.TagInterface{
			hb.NewThead().Children([]hb.TagInterface{
				hb.NewTR().Children([]hb.TagInterface{
					hb.NewTH().
						Child(sortableColumnLabel("Name, Alias, ID", "id", data.request, data.sortBy, data.sortOrder)).
						Style(""),
					hb.NewTH().
						Child(sortableColumnLabel("Start Time", "started_at", data.request, data.sortBy, data.sortOrder)).
						Style("width: 1px;"),
					hb.NewTH().
						Child(sortableColumnLabel("End Time", "completed_at", data.request, data.sortBy, data.sortOrder)).
						Style("width: 1px;"),
					hb.NewTH().
						HTML("Duration").
						Style("width: 1px;"),
					hb.NewTH().
						Child(sortableColumnLabel("Status", "status", data.request, data.sortBy, data.sortOrder)).
						Style("width: 100px;"),
					hb.NewTH().
						HTML("Actions").
						Style("width: 160px;"),
				}),
			}),
			hb.NewTbody().Children(lo.Map(queuedTaskList, func(queuedTask taskstore.TaskQueueInterface, _ int) hb.TagInterface {
				task, taskExists := lo.Find(allTasks, func(t taskstore.TaskDefinitionInterface) bool {
					return t.ID() == queuedTask.TaskID()
				})

				taskName := lo.IfF(taskExists, func() string { return task.Title() }).Else("Unknown")

				buttonDelete := hb.NewButton().
					Class("btn btn-sm btn-danger").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					HTML("Delete").
					HxPost(link(data.request, PathHome, map[string]string{
						"action":   ActionModalQueuedTaskDeleteShow,
						"queue_id": queuedTask.ID(),
						"page":     data.page,
						"by":       data.sortBy,
						"sort":     data.sortOrder,
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonParameters := hb.NewButton().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					HTML("Parameters").
					HxPost(link(data.request, PathHome, map[string]string{
						"action":   ActionModalQueuedTaskParametersShow,
						"queue_id": queuedTask.ID(),
						"page":     data.page,
						"by":       data.sortBy,
						"sort":     data.sortOrder,
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonDetails := hb.NewButton().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					HTML("Details").
					HxPost(link(data.request, PathHome, map[string]string{
						"action":   ActionModalQueuedTaskDetailsShow,
						"queue_id": queuedTask.ID(),
						"page":     data.page,
						"by":       data.sortBy,
						"sort":     data.sortOrder,
					})).
					// HxTarget("#QueuedTasksListTable").
					// HxSelectOob("#ModalMessage").
					HxTarget("body").
					HxSwap("beforeend")

				buttonRequeue := hb.NewButton().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					HTML("Requeue").
					HxPost(link(data.request, PathHome, map[string]string{
						"action":   ActionModalQueuedTaskRequeueShow,
						"queue_id": queuedTask.ID(),
						"page":     data.page,
						"by":       data.sortBy,
						"sort":     data.sortOrder,
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonRestart := hb.NewButton().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					HTML("Restart").
					HxPost(link(data.request, PathHome, map[string]string{
						"action":   "queue-task-restart",
						"queue_id": queuedTask.ID(),
						"page":     data.page,
						"by":       data.sortBy,
						"sort":     data.sortOrder,
					})).
					HxTarget("#QueuedTasksListTable").
					HxSwap("outerHTML")

				// linkTask := hb.NewHyperlink().
				// 	HTML(queuedTask.ID()).
				// 	Href(links.Admin().Tasks(map[string]string{
				// 		"task_id": queuedTask.ID(),
				// 	}))

				startedAtDate := lo.IfF(queuedTask.StartedAt() != "", func() string {
					return carbon.Parse(queuedTask.StartedAt()).Format("d M Y")
				}).Else("-")
				startedAtTime := lo.IfF(queuedTask.StartedAt() != "", func() string {
					return carbon.Parse(queuedTask.StartedAt()).ToTimeString()
				}).Else("-")
				completeddAtDate := lo.IfF(queuedTask.CompletedAt() != "", func() string {
					return carbon.Parse(queuedTask.CompletedAt()).Format("d M Y")
				}).Else("-")
				completeddAtTime := lo.IfF(queuedTask.CompletedAt() != "", func() string {
					return carbon.Parse(queuedTask.CompletedAt()).ToTimeString()
				}).Else("-")

				elapsedTime := lo.IfF(queuedTask.StartedAt() != "" && queuedTask.CompletedAt() != "", func() string {
					startedAt := carbon.Parse(queuedTask.StartedAt())
					completedAt := carbon.Parse(queuedTask.CompletedAt())
					return cast.ToString(completedAt.DiffInMinutes(startedAt)) + " minutes"
				}).Else("-")

				return hb.NewTR().Children([]hb.TagInterface{
					hb.NewTD().
						Child(hb.NewDiv().Text(taskName)).
						Child(hb.NewDiv().Text("Alias: ").Text(task.Alias())).
						Child(hb.NewDiv().Text("ID: ").Text(queuedTask.ID())),
					hb.NewTD().
						Child(hb.NewDiv().Text(startedAtDate)).
						Child(hb.NewDiv().Text(startedAtTime)).
						Style("white-space: nowrap;"),
					hb.NewTD().
						Child(hb.NewDiv().Text(completeddAtDate)).
						Child(hb.NewDiv().Text(completeddAtTime)).
						Style("white-space: nowrap;"),
					hb.NewTD().
						Child(hb.NewDiv().Text(elapsedTime)).
						Style("white-space: nowrap;"),
					hb.NewTD().
						Text(queuedTask.GetStatus()),
					hb.NewTD().
						Style("text-align: center;").
						Child(buttonParameters).
						Child(buttonDetails).
						Child(buttonRequeue).
						Child(buttonRestart).
						Child(buttonDelete),
				})
			})),
		})

	return hb.NewDiv().
		ID("QueuedTasksListTable").
		Children([]hb.TagInterface{
			p.tableFilter(data),
			table,
			p.tablePagination(int(queuedTaskCount), pageInt, perPage, data),
		})
}

func (p *pageQueueManager) tableFilter(data pageQueueManagerData) *hb.Tag {
	buttonFilter := hb.NewButton().
		Class("btn btn-sm btn-info me-2").
		Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
		Child(hb.NewI().Class("bi bi-filter me-2")).
		Text("Filters").
		HxPost(link(data.request, PathHome, map[string]string{
			"action": ActionModalQueuedTaskFilterShow,
		})).
		HxTarget("body").
		HxSwap("beforeend")

	return hb.NewDiv().
		Class(`d-flex justify-content-left align-items-center pagination-primary-soft rounded mb-3 p-2`).
		Style(`background-color: #f8f9fa;`).
		Child(
			buttonFilter,
		).
		ChildIf(
			data.status != "",
			hb.NewSpan().
				Text("Status: ").
				Text(data.status),
		)
}

func (p *pageQueueManager) tablePagination(count int, page int, perPage int, data pageQueueManagerData) *hb.Tag {
	url := link(data.request, PathHome, map[string]string{
		"page": "",
	})

	pagination := bs.Pagination(bs.PaginationOptions{
		NumberItems:       count,
		CurrentPageNumber: page,
		PagesToShow:       20,
		PerPage:           perPage,
		URL:               url,
	})

	return hb.NewDiv().
		Class(`d-flex justify-content-left mt-5 pagination-primary-soft rounded mb-0`).
		HTML(pagination)
}

func (p *pageQueueManager) pageQueueManagerPrepareData(r *http.Request) (data pageQueueManagerData, errorMessage string) {
	data.request = *r
	data.action = req.GetStringTrimmed(r, "action")
	data.page = req.GetStringTrimmedOr(r, "page", "0")
	data.sortOrder = req.GetStringTrimmed(r, "sort")
	data.sortBy = req.GetStringTrimmed(r, "by")
	data.status = req.GetStringTrimmed(r, "status")
	data.search = req.GetStringTrimmed(r, "search")
	data.dateFrom = req.GetStringTrimmed(r, "date_from")
	data.dateTo = req.GetStringTrimmed(r, "date_to")
	data.queueID = req.GetStringTrimmed(r, "queue_id")
	data.taskID = req.GetStringTrimmed(r, "task_id")

	if !lo.Contains([]string{sb.ASC, sb.DESC}, data.sortOrder) {
		data.sortOrder = sb.DESC
	}

	if !lo.Contains([]string{
		taskstore.COLUMN_STARTED_AT,
		taskstore.COLUMN_COMPLETED_AT,
		taskstore.COLUMN_ID,
		taskstore.COLUMN_TASK_ID,
		taskstore.COLUMN_STATUS,
	}, data.sortBy) {
		data.sortBy = taskstore.COLUMN_CREATED_AT
	}

	p.admin.logger.Info("queueManagerController > prepareData > ", data.sortBy, data.sortOrder)

	return data, ""
}
