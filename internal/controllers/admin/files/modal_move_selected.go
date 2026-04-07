package admin

import (
	"project/internal/links"
)

// modalMoveSelected creates the bulk move modal HTML
func (c *FileManagerController) modalMoveSelected(currentDirectory string, directoryList []FileEntry) string {
	url := links.Admin().FileManager()

	// Build directory options for the dropdown
	dirOptions := `<option value="">-- Root Directory --</option>`
	for _, dir := range directoryList {
		if dir.Name == "." || dir.Name == ".." {
			continue
		}
		dirOptions += `<option value="` + dir.Path + `">` + dir.Name + `</option>`
	}

	return `
<!-- START: Modal Move Selected -->
<div class="modal fade" id="ModalMoveSelected" role="dialog">
	<div class="modal-dialog" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title">Move Selected Items</h5>
				<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
			</div>
			<div class="modal-body">
				<form id="FormMoveSelected" name="FormMoveSelected" onsubmit="event.preventDefault(); executeMoveSelected(); return false;">
					<div class="form-group mb-3">
						<label class="form-label">Selected Items</label>
						<div id="SelectedItemsList" class="border p-2 rounded" style="max-height: 150px; overflow-y: auto;">
							<em class="text-muted">No items selected</em>
						</div>
					</div>
					<div class="form-group mb-3">
						<label class="form-label">Destination Directory</label>
						<select name="destination_dir" class="form-select" id="DestinationDirSelect">
							<option value="` + currentDirectory + `">-- Current Directory --</option>` + dirOptions + `
						</select>
						<small class="form-text text-muted">Select the folder where you want to move the selected items</small>
					</div>
					<input type="hidden" name="action" value="` + JSON_ACTION_BULK_MOVE + `" />
					<input type="hidden" name="current_dir" value="` + currentDirectory + `" />
					<input type="hidden" name="selected_items" id="SelectedItemsInput" value="" />
				</form>
			</div>
			<div class="modal-footer" style="display:block;">
				<button type="button" class="btn btn-secondary float-start" data-bs-dismiss="modal">
					<i class="bi bi-chevron-left"></i>
					Cancel
				</button>
				<button type="button" class="btn btn-warning float-end" onclick="executeMoveSelected()">
					<i class="bi bi-folder-symlink"></i>
					Move Selected Items
				</button>
			</div>
		</div>
	</div>
</div>
<script>
	// Update the modal when it's shown
	document.getElementById('ModalMoveSelected').addEventListener('show.bs.modal', function () {
		const selected = getSelectedItems();
		const listContainer = document.getElementById('SelectedItemsList');
		const itemsInput = document.getElementById('SelectedItemsInput');
		
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

	function executeMoveSelected() {
		const form = document.getElementById('FormMoveSelected');
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
			const modal = bootstrap.Modal.getInstance(document.getElementById('ModalMoveSelected'));
			if (modal) {
				modal.hide();
			}
		})
		.catch(error => {
			showNotification("IO Error: " + error.message, "error");
			console.error('Move selected error:', error);
		});
	}
</script>
<!-- END: Modal Move Selected -->
	`
}
