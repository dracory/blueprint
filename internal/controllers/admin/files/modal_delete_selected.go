package admin

import (
	"project/internal/links"
)

// modalDeleteSelected creates the bulk delete modal HTML
func (c *FileManagerController) modalDeleteSelected(currentDirectory string) string {
	url := links.Admin().FileManager()

	return `
<!-- START: Modal Delete Selected -->
<div class="modal fade" id="ModalDeleteSelected" role="dialog">
	<div class="modal-dialog" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">Delete Selected Items</h5>
				<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
			</div>
			<div class="modal-body">
				<p class="text-danger">
					<strong>Warning:</strong> This operation is final and CANNOT BE undone!
				</p>
				<form id="FormDeleteSelected" name="FormDeleteSelected" onsubmit="event.preventDefault(); executeDeleteSelected(); return false;">
					<div class="form-group mb-3">
						<label class="form-label">Selected Items to Delete</label>
						<div id="DeleteSelectedItemsList" class="border p-2 rounded" style="max-height: 150px; overflow-y: auto;">
							<em class="text-muted">No items selected</em>
						</div>
					</div>
					<input type="hidden" name="action" value="bulk_delete" />
					<input type="hidden" name="current_dir" value="` + currentDirectory + `" />
					<input type="hidden" name="selected_items" id="DeleteSelectedItemsInput" value="" />
				</form>
			</div>
			<div class="modal-footer" style="display:block;">
				<button type="button" class="btn btn-secondary float-start" data-bs-dismiss="modal">
					<i class="bi bi-chevron-left"></i>
					Cancel
				</button>
				<button type="button" class="btn btn-danger float-end" onclick="executeDeleteSelected()">
					<i class="bi bi-trash"></i>
					Delete Selected Items
				</button>
			</div>
		</div>
	</div>
</div>
<script>
	// Update the modal when it's shown
	document.getElementById('ModalDeleteSelected').addEventListener('show.bs.modal', function () {
		const selected = getSelectedItems();
		const listContainer = document.getElementById('DeleteSelectedItemsList');
		const itemsInput = document.getElementById('DeleteSelectedItemsInput');
		
		if (selected.length === 0) {
			listContainer.innerHTML = '<em class="text-muted">No items selected</em>';
			itemsInput.value = '';
		} else {
			const ul = document.createElement('ul');
			ul.className = 'list-unstyled mb-0';
			selected.forEach(item => {
				const li = document.createElement('li');
				const icon = item.type === 'directory' ? '<i class="bi bi-folder me-1"></i>' : '<i class="bi bi-file me-1"></i>';
				li.innerHTML = icon + item.path.split('/').pop();
				ul.appendChild(li);
			});
			listContainer.innerHTML = '';
			listContainer.appendChild(ul);
			itemsInput.value = JSON.stringify(selected);
		}
	});

	function executeDeleteSelected() {
		const form = document.getElementById('FormDeleteSelected');
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
			const modal = bootstrap.Modal.getInstance(document.getElementById('ModalDeleteSelected'));
			if (modal) {
				modal.hide();
			}
		})
		.catch(error => {
			showNotification("IO Error: " + error.message, "error");
			console.error('Delete selected error:', error);
		});
	}
</script>
<!-- END: Modal Delete Selected -->
	`
}
