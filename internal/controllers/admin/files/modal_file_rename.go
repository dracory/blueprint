package admin

import (
	"project/internal/links"
)

// modalFileRename creates the file rename modal HTML
func (c *FileManagerController) modalFileRename(currentDirectory string) string {
	url := links.Admin().FileManager()
	return `
<!-- START: Modal File Rename -->
<div class="modal fade" id="ModalFileRename" role="dialog">
	<div class="modal-dialog" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="myModalLabel">File Rename</h5>
				<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
			</div>
			<div class="modal-body">
				<form id="FormFileRename" name="FormFileRename" onsubmit="event.preventDefault(); fileRename(); return false;">
					<div class="form-group">
						<label>New Name</label>
						<input name="new_file" value="" class="form-control" />
					</div>
					<input type="hidden" name="action" value="` + JSON_ACTION_FILE_RENAME + `" />
					<input type="hidden" name="current_dir" value="` + currentDirectory + `" />
					<input type="hidden" name="rename_file" value="" />
				</form>
			</div>
			<div class="modal-footer" style="display:block;">
				<button type="button" class="btn btn-secondary float-start" data-bs-dismiss="modal">
					<i class="bi bi-chevron-left"></i>
					Close
				</button>
				<button type="button" class="btn btn-success float-end" onclick="fileRename()">
					<i class="bi bi-check"></i>
					Rename File
				</button>
			</div>
		</div>
	</div>
</div>
<script>
	function modalFileRenameShow(fileName) {
		const form = document.getElementById('FormFileRename');
		form.querySelector('input[name="new_file"]').value = fileName;
		form.querySelector('input[name="rename_file"]').value = fileName;
		const modal = new bootstrap.Modal(document.getElementById('ModalFileRename'), {})
		modal.show();
	}
	function fileRename() {
		const form = document.getElementById('FormFileRename');
		const formData = new FormData(form);
		
		fetch("` + url + `", {
			method: 'POST',
			body: formData
		})
		.then(response => {
			if (response.redirected) {
				window.location.href = response.url;
				return;
			}
			return response.json();
		})
		.then(data => {
			if (!data) return;
			if (data.status === "success") {
				showNotification(data.message, "success");
				setTimeout(() => {
					window.location.href = window.location.href;
				}, 1000)
			} else {
				showNotification(data.message, "error");
			}
			const modal = bootstrap.Modal.getInstance(document.getElementById('ModalFileRename'));
			if (modal) {
				modal.hide();
			}
		})
		.catch(error => {
			showNotification("IO Error", "error");
			console.error('File rename error:', error);
		});
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
</script>
<!-- END: Modal File Rename -->
	`
}
