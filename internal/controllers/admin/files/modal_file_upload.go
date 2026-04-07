package admin

import (
	"project/internal/links"
)

// modalFileUpload creates the file upload modal HTML
func (c *FileManagerController) modalFileUpload(currentDirectory string) string {
	url := links.Admin().FileManager()
	return `
<!-- START: Modal Upload File -->
<div class="modal fade" id="ModalUploadFile" role="dialog">
	<div class="modal-dialog" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="myModalLabel">File Upload</h5>
				<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
			</div>
			<div class="modal-body">
				<form id="FormFileUpload" name="FormFileUpload" method="POST" enctype="multipart/form-data">
				    <input type="hidden" name="action" value="file_upload" />
					<input type="hidden" name="current_dir" value="` + currentDirectory + `" />
					<input type="file" id="file-input" name="upload_file" value="" />
					<input type="hidden" name="_token" value="<?php echo csrf_token(); ?>" />
				</form>
			</div>
			<div class="modal-footer" style="display:block;">
				<button type="button" class="btn btn-secondary float-start" data-bs-dismiss="modal">
					<i class="bi bi-chevron-left"></i>
					Close
				</button>
				<button type="button" class="btn btn-primary float-end" onclick="fileUpload();/*FormFileUpload.submit();*/">
					<i class="bi bi-check"></i>
					Start Upload
				</button>
			</div>
		</div>
	</div>
</div>

<script>
function fileUpload() {
	const file = document.getElementById('file-input').files[0];
	const formData = new FormData();
	formData.append('action', 'file_upload');
	formData.append('current_dir', '` + currentDirectory + `');
	formData.append('upload_file', file);

	try {
		fetch("` + url + `", { method: 'POST', body: formData })
		.then((response) => response.json())
		.then((data) => {
			if (data.status === "success") {
				showNotification(data.message, "success");
			} else {
				showNotification(data.message, "error");
			}
			const modalElement = document.getElementById('ModalUploadFile');
			const modal = bootstrap.Modal.getInstance(modalElement) || new bootstrap.Modal(modalElement);
			modal.hide();
			setTimeout(() => {
				window.location.href = window.location.href;
			}, 1000)
		})
	} catch (error) {
		showNotification("IO Error", "error");
		console.error('File upload error:', error);
	}
	
	// Simple notification function to replace jQuery notify
	function showNotification(message, type) {
		const notification = document.createElement('div');
		const alertClass = type === 'success' ? 'success' : 'danger';
		notification.className = 'alert alert-' + alertClass + ' alert-dismissible fade show position-fixed';
		notification.style.cssText = 'top: 20px; right: 20px; z-index: 9999; min-width: 300px;';
		notification.innerHTML = 
			'<div>' + message + '</div>' +
			'<button type="button" class="btn-close" data-bs-dismiss="alert"></button>';
		
		document.body.appendChild(notification);
		
		setTimeout(() => {
			if (notification.parentNode) {
				notification.parentNode.removeChild(notification);
			}
		}, 3000);
	}
}
</script>
<!-- END: Modal Upload File -->
	`
}
