package taskadmin

import "net/http"

func (a *admin) taskEnqueueModalShow(r *http.Request) string {
	return a.modalTaskEnqueue(r).ToHTML()
}
