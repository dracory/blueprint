package admin

import (
	"github.com/dracory/hb"
)

// uiManager creates the main file manager UI
func (c *FileManagerController) uiManager(currentDirectory, parentDirectory string, directoryList, fileList []FileEntry) string {
	buttonUpload := hb.Button().
		Class("btn btn-secondary float-end me-2").
		Data("bs-toggle", "modal").
		Data("bs-target", "#ModalUploadFile").
		// OnClick(`modalFileUploadShow()`).
		Child(hb.I().Class("bi bi-upload").Style("margin-right: 5px;")).
		HTML("Upload File")

	buttonUploadMultiple := hb.Button().
		Class("btn btn-primary float-end me-2").
		Data("bs-toggle", "modal").
		Data("bs-target", "#ModalUploadMultipleFiles").
		Child(hb.I().Class("bi bi-cloud-upload").Style("margin-right: 5px;")).
		HTML("Upload Multiple")

	buttonDirectoryCreate := hb.Button().
		Class("btn btn-info float-end me-2").
		Data("bs-toggle", "modal").
		Data("bs-target", "#ModalDirectoryCreate").
		// OnClick(`modalDirectoryCreateShow()`).
		Child(hb.I().Class("bi bi-plus-circle").Style("margin-right: 5px;")).
		HTML("New directory")

	// Bulk actions dropdown
	dropdownMenu := hb.Div().
		Class("dropdown-menu").
		Child(hb.Button().
			Class("dropdown-item text-danger").
			Attr("onclick", "selectAllItems()").
			Child(hb.I().Class("bi bi-check-square me-2")).
			HTML("Select All")).
		Child(hb.Button().
			Class("dropdown-item text-danger").
			Attr("onclick", "unselectAllItems()").
			Child(hb.I().Class("bi bi-square me-2")).
			HTML("Un-select All")).
		Child(hb.Div().Class("dropdown-divider")).
		Child(hb.Button().
			Class("dropdown-item text-danger").
			Attr("onclick", "showMoveSelectedModal()").
			Child(hb.I().Class("bi bi-folder-symlink me-2")).
			HTML("Move Selected")).
		Child(hb.Button().
			Class("dropdown-item text-danger").
			Attr("onclick", "showDeleteSelectedModal()").
			Child(hb.I().Class("bi bi-trash me-2")).
			HTML("Delete Selected"))

	title := hb.Heading3().
		HTML("File Manager").
		Child(buttonUpload).
		Child(buttonUploadMultiple).
		Child(buttonDirectoryCreate)

	// Bulk actions widget above the table
	selectAllCheckbox := hb.Input().
		Type("checkbox").
		Class("form-check-input").
		ID("select-all-checkbox").
		Attr("onclick", "toggleSelectAll(this)").
		Style("margin-right: 5px;")

	dropdownToggle := hb.Button().
		Class("btn btn-sm btn-outline-secondary dropdown-toggle").
		Attr("type", "button").
		Data("bs-toggle", "dropdown").
		Aria("expanded", "false")

	bulkActionsWidget := hb.Div().
		Class("mb-2").
		Child(hb.Div().Class("dropdown d-inline-block").
			Child(selectAllCheckbox).
			Child(dropdownToggle).
			Child(dropdownMenu))

	script := hb.Script(`
// Hide select buttons initially
const selectButtons = document.querySelectorAll('.btn-select');
selectButtons.forEach(btn => btn.style.display = 'none');

var openerArgs = {};

function fileSelectedUrl(selectedFileUrl) {
	if (window.opener === null) {
		return true;
	}
	window.opener.postMessage({msg: 'media-manager-file-selected', url: selectedFileUrl, args: openerArgs}, '*');
	window.close();
}

function messageReceived(event) {
	var data = event.data;
	openerArgs = data;
	console.log(data);
	// Show select buttons when message is received
	const selectButtons = document.querySelectorAll('.btn-select');
	selectButtons.forEach(btn => btn.style.display = 'inline-block');
}

window.addEventListener("message", messageReceived, false);

if (window.opener !== null) {
	window.opener.postMessage({msg: 'media-manager-loaded'}, '*');
}

// Toggle select all checkboxes
function toggleSelectAll(masterCheckbox) {
	const checkboxes = document.querySelectorAll('.file-select');
	checkboxes.forEach(cb => cb.checked = masterCheckbox.checked);
}

// Select all items
function selectAllItems() {
	const checkboxes = document.querySelectorAll('.file-select');
	checkboxes.forEach(cb => cb.checked = true);
	const masterCheckbox = document.getElementById('select-all-checkbox');
	if (masterCheckbox) masterCheckbox.checked = true;
}

// Un-select all items
function unselectAllItems() {
	const checkboxes = document.querySelectorAll('.file-select');
	checkboxes.forEach(cb => cb.checked = false);
	const masterCheckbox = document.getElementById('select-all-checkbox');
	if (masterCheckbox) masterCheckbox.checked = false;
}

// Get selected items
function getSelectedItems() {
	const checkboxes = document.querySelectorAll('.file-select:checked');
	return Array.from(checkboxes).map(cb => ({
		path: cb.getAttribute('data-path'),
		type: cb.getAttribute('data-type')
	}));
}

// Show move selected modal
function showMoveSelectedModal() {
	const selected = getSelectedItems();
	if (selected.length === 0) {
		showNotification('Please select at least one file or folder to move', 'error');
		return;
	}
	const modal = new bootstrap.Modal(document.getElementById('ModalMoveSelected'), {});
	modal.show();
}

// Show delete selected modal
function showDeleteSelectedModal() {
	const selected = getSelectedItems();
	if (selected.length === 0) {
		showNotification('Please select at least one file or folder to delete', 'error');
		return;
	}
	const modal = new bootstrap.Modal(document.getElementById('ModalDeleteSelected'), {});
	modal.show();
}

// Simple notification function
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
	`)

	html := title.ToHTML() + `
	` + bulkActionsWidget.ToHTML() + `
	` + c.tableFileList(currentDirectory, parentDirectory, directoryList, fileList) + `
	` + c.modalDirectoryCreate(currentDirectory) + `
	` + c.modalDirectoryDelete(currentDirectory) + `
	` + c.modalFileDelete(currentDirectory) + `
	` + c.modalFileRename(currentDirectory) + `
	` + c.modalFileUpload(currentDirectory) + `
	` + c.modalFileUploadMultiple(currentDirectory) + `
	` + c.modalFileView() + `
	` + c.modalMoveSelected(currentDirectory, directoryList) + `
	` + c.modalDeleteSelected(currentDirectory) +
		script.ToHTML()

	return html
}
