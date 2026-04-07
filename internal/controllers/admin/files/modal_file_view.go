package admin

// modalFileView creates the file view modal HTML
func (c *FileManagerController) modalFileView() string {
	return `
<div class="modal fade" id="ModalFileView" role="dialog">
	<div class="modal-dialog" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="myModalLabel">File View</h5>
				<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
			</div>
			<div class="modal-body" style="text-align:center;">
			    <img id="FileViewImage" src="" class="img-fluid" style="max-height:200px;;" />
			</div>
			<div class="modal-footer" style="display:block;">
				<button type="button" class="btn btn-secondary float-start" data-bs-dismiss="modal">
					<i class="bi bi-chevron-left"></i>
					Close
				</button>
			</div>
		</div>
	</div>
</div>
<script>
	function modalFileViewShow(fileURL) {
		document.getElementById('FileViewImage').src = fileURL;
		const modal = new bootstrap.Modal(document.getElementById('ModalFileView'), {})
		modal.show();
	}
</script>
	`
}
