package taskadmin

func (p *pageQueueManager) onModalQueuedTaskFilterShow(data pageQueueManagerData) string {
	return p.modalQueuedTaskFilters(data).ToHTML()
}
