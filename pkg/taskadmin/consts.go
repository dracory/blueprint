package taskadmin

type contextKey string

const keyEndpoint = contextKey("endpoint")

const PathHome = "home"
const ActionModalQueuedTaskDeleteShow = "modal-queued-task-delete-show"
const ActionModalQueuedTaskDeleteSubmitted = "modal-queued-task-delete-submitted"
const ActionModalQueuedTaskDetailsShow = "modal-queued-task-details-show"
const ActionModalQueuedTaskFilterShow = "modal-queued-task-filter-show"
const ActionModalQueuedTaskParametersShow = "modal-queued-task-parameters-show"
const ActionModalQueuedTaskRequeueShow = "modal-queued-task-requeue-show"
const ActionModalQueuedTaskRequeueSubmitted = "modal-queued-task-requeue-submitted"
const ActionQueuedTaskRestart = "queue-task-restart"
const ActionModalQueuedTaskEnqueueShow = "modal-queued-task-enqueue-show"
const ActionModalQueuedTaskEnqueueSubmitted = "modal-task-enqueue-submitted"
