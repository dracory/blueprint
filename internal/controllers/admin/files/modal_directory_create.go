package admin

import (
	"project/internal/links"
)

// modalDirectoryCreate creates the directory creation modal HTML
func (c *FileManagerController) modalDirectoryCreate(currentDirectory string) string {
	url := links.Admin().FileManager()
	if currentDirectory == "" {
		currentDirectory = "/"
	}
	return `
<div class="modal fade" id="ModalDirectoryCreate" role="dialog">
	<div class="modal-dialog" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="myModalLabel">New Directory</h5>
				<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
			</div>
			<div class="modal-body">
				<form id="FormDirectoryCreate" name="FormDirectoryCreate"" method="POST">
					<div class="form-group">
						<label>Directory name</label>
						<input type="text" class="form-control" name="create_dir" value="" />
					</div>
					<input type="hidden" name="action" value="` + JSON_ACTION_DIRECTORY_CREATE + `" />
					<input type="hidden" name="current_dir" value="` + currentDirectory + `" />
					<input type="hidden" name="_token" value="<?php echo csrf_token(); ?>" />
				</form>
			</div>
			<div class="modal-footer" style="display:block;">
				<button type="button" class="btn btn-secondary float-start" data-bs-dismiss="modal">
					<i class="bi bi-chevron-left"></i>
					Close
				</button>
				<button type="button" class="btn btn-primary float-end" onclick="directoryCreate();">
					<i class="bi bi-check"></i>
					Create Directory
				</button>
			</div>
		</div>
	</div>
</div>
<script>
	function directoryCreate() {
		const form = document.getElementById('FormDirectoryCreate');
		const formData = new FormData(form);
		
		fetch("` + url + `", {
			method: 'POST',
			body: formData
		})
		.then(response => response.json())
		.then(data => {
			if (data.status === "success") {
				// Show success notification using native browser notification or custom implementation
				showNotification(data.message, "success");
			} else {
				showNotification(data.message, "error");
			}
			
			// Close modal using Bootstrap's native JavaScript API
			const modalElement = document.getElementById('ModalDirectoryCreate');
			const modal = bootstrap.Modal.getInstance(modalElement) || new bootstrap.Modal(modalElement);
			modal.hide();
			
			// Refresh page after 1 second
			setTimeout(() => {
				window.location.href = window.location.href;
			}, 1000);
		})
		.catch(error => {
			showNotification("IO Error", "error");
			console.error('Directory creation error:', error);
		});
	}
	
	// Simple notification function to replace jQuery notify
	function showNotification(message, type) {
		// Create notification element
		const notification = document.createElement('div');
		const alertClass = type === 'success' ? 'success' : 'danger';
		notification.className = 'alert alert-' + alertClass + ' alert-dismissible fade show position-fixed';
		notification.style.cssText = 'top: 20px; right: 20px; z-index: 9999; min-width: 300px;';
		notification.innerHTML = 
			'<div>' + message + '</div>' +
			'<button type="button" class="btn-close" data-bs-dismiss="alert"></button>';
		
		// Add to body
		document.body.appendChild(notification);
		
		// Auto remove after 3 seconds
		setTimeout(() => {
			if (notification.parentNode) {
				notification.parentNode.removeChild(notification);
			}
		}, 3000);
	}
</script>
	`
}
