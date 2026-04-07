package admin

import (
	"project/internal/links"
)

// modalFileDelete creates the file deletion modal HTML
func (c *FileManagerController) modalFileDelete(currentDirectory string) string {
	url := links.Admin().FileManager()
	return `
	<!-- START: Modal File Delete -->
	<div class="modal fade" id="ModalFileDelete" role="dialog">
		<div class="modal-dialog" role="document">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title" id="myModalLabel">Confirm File Delete</h5>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<div class="modal-body">
					<p>
						Are you sure you want to delete this file?
					</p>
					<p class="text-danger">
						This operation is final and CANNOT BE undone!
					</p>
					<form id="FormFileDelete" name="FormFileDelete" method="POST">
					    <input type="hidden" name="action" value="` + JSON_ACTION_FILE_DELETE + `" />
						<input type="hidden" name="current_dir" value="` + currentDirectory + `" />
						<input type="hidden" name="delete_file" value="" />
						<input type="hidden" name="_token" value="<?php echo csrf_token(); ?>" />
					</form>
				</div>
				<div class="modal-footer" style="display:block;">
					<button type="button" class="btn btn-secondary float-start" data-bs-dismiss="modal">
						<i class="bi bi-chevron-left"></i>
						Close
					</button>
					<button type="button" class="btn btn-danger float-end" onclick="fileDelete()">
						<i class="bi bi-trash"></i>
						Delete File
					</button>
				</div>
			</div>
		</div>
	</div>
	<script>
		function modalFileDeleteShow(fileName) {
			document.getElementById('FormFileDelete').querySelector('input[name="delete_file"]').value = fileName;
			const modal = new bootstrap.Modal(document.getElementById('ModalFileDelete'), {})
			modal.show();
		}
		function fileDelete() {
			const form = document.getElementById('FormFileDelete');
			const formData = new FormData(form);
			
			fetch("` + url + `", {
				method: 'POST',
				body: formData
			})
			.then(response => response.json())
			.then(data => {
				setTimeout(() => {
					window.location.href = window.location.href;
				}, 1);
				
				const modalElement = document.getElementById('ModalFileDelete');
				const modal = bootstrap.Modal.getInstance(modalElement) || new bootstrap.Modal(modalElement);
				modal.hide();
				
				if (data.status === "success") {
					showNotification(data.message, "success");
				} else {
					showNotification(data.message, "error");
				}
			})
			.catch(error => {
				showNotification("IO Error", "error");
				console.error('File delete error:', error);
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
	<!-- END: Modal File Delete -->
	`
}
