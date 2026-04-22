package admin

import (
	"project/internal/links"
)

// modalFileUploadMultiple creates the multiple file upload modal with drag and drop
func (c *FileManagerController) modalFileUploadMultiple(currentDirectory string) string {
	url := links.Admin().FileManager()
	return `
<!-- START: Modal Upload Multiple Files -->
<div class="modal fade" id="ModalUploadMultipleFiles" role="dialog">
	<div class="modal-dialog modal-lg" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<h5 class="modal-title" id="myModalLabel">Upload Multiple Files</h5>
				<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
			</div>
			<div class="modal-body">
				<form id="FormFileUploadMultiple" name="FormFileUploadMultiple" method="POST" enctype="multipart/form-data">
				    <input type="hidden" name="action" value="file_upload_multiple" />
					<input type="hidden" name="current_dir" value="` + currentDirectory + `" />
					<input type="hidden" name="_token" value="<?php echo csrf_token(); ?>" />
					
					<!-- Drag and Drop Area -->
					<div id="drop-area" class="border border-2 border-dashed rounded-3 p-5 text-center mb-3" style="border-color: #dee2e6 !important;">
						<div id="drop-area-content">
							<i class="bi bi-cloud-upload display-4 text-muted mb-3"></i>
							<h5 class="mb-2">Drag & Drop Files Here</h5>
							<p class="text-muted mb-3">or click to select files</p>
							<button type="button" class="btn btn-outline-primary" onclick="document.getElementById('file-input-multiple').click()">
								<i class="bi bi-folder2-open me-2"></i>Select Files
							</button>
						</div>
						
						<!-- File Input -->
						<input type="file" id="file-input-multiple" name="upload_files[]" multiple style="display: none;" />
						
						<!-- Preview Area -->
						<div id="file-preview" class="mt-3" style="display: none;">
							<hr>
							<h6>Selected Files:</h6>
							<div id="file-list" class="row"></div>
						</div>
					</div>
					
					<!-- Upload Progress -->
					<div id="upload-progress" class="progress mb-3" style="display: none;">
						<div class="progress-bar" role="progressbar" style="width: 0%" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100">0%</div>
					</div>
					
					<!-- Upload Status -->
					<div id="upload-status" class="mb-3"></div>
				</form>
			</div>
			<div class="modal-footer" style="display:block;">
				<button type="button" class="btn btn-secondary float-start" data-bs-dismiss="modal">
					<i class="bi bi-chevron-left"></i>
					Close
				</button>
				<button type="button" class="btn btn-primary float-end" id="btn-start-upload" onclick="startMultipleUpload()" disabled>
					<i class="bi bi-upload"></i>
					Start Upload
				</button>
				<button type="button" class="btn btn-warning float-end me-2" id="btn-clear-files" onclick="clearFiles()" style="display: none;">
					<i class="bi bi-x-circle"></i>
					Clear Files
				</button>
			</div>
		</div>
	</div>
</div>

<script>
// Drag and Drop functionality
const dropArea = document.getElementById('drop-area');
const fileInput = document.getElementById('file-input-multiple');
const fileList = document.getElementById('file-list');
const filePreview = document.getElementById('file-preview');
const uploadBtn = document.getElementById('btn-start-upload');
const clearBtn = document.getElementById('btn-clear-files');
const uploadProgress = document.getElementById('upload-progress');
const uploadStatus = document.getElementById('upload-status');

// Prevent default drag behaviors
['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
    dropArea.addEventListener(eventName, preventDefaults, false);
    document.body.addEventListener(eventName, preventDefaults, false);
});

// Highlight drop area when dragging over it
['dragenter', 'dragover'].forEach(eventName => {
    dropArea.addEventListener(eventName, highlight, false);
});

['dragleave', 'drop'].forEach(eventName => {
    dropArea.addEventListener(eventName, unhighlight, false);
});

// Handle dropped files
dropArea.addEventListener('drop', handleDrop, false);

function preventDefaults(e) {
    e.preventDefault();
    e.stopPropagation();
}

function highlight(e) {
    dropArea.classList.add('border-primary', 'bg-light');
}

function unhighlight(e) {
    dropArea.classList.remove('border-primary', 'bg-light');
}

function handleDrop(e) {
    const dt = e.dataTransfer;
    const files = dt.files;
    handleFiles(files);
}

function handleFiles(files) {
	if (files.length > 0) {
		// Use DataTransfer to set files on the input element
		const dt = new DataTransfer();
		Array.from(files).forEach(file => dt.items.add(file));
		fileInput.files = dt.files;

		displayFiles(files);
		uploadBtn.disabled = false;
		clearBtn.style.display = 'inline-block';
	}
}

// File input change handler
fileInput.addEventListener('change', function(e) {
    const files = e.target.files;
    if (files.length > 0) {
        displayFiles(files);
        uploadBtn.disabled = false;
        clearBtn.style.display = 'inline-block';
    }
});

function displayFiles(files) {
    fileList.innerHTML = '';
    filePreview.style.display = 'block';
    
    Array.from(files).forEach((file, index) => {
        const fileItem = document.createElement('div');
        fileItem.className = 'col-md-6 mb-2';
        
        const fileSize = formatFileSize(file.size);
        const fileIcon = getFileIcon(file.type);
        
        fileItem.innerHTML = 
            '<div class="card h-100">' +
                '<div class="card-body p-2">' +
                    '<div class="d-flex align-items-center">' +
                        '<div class="me-2">' + fileIcon + '</div>' +
                        '<div class="flex-grow-1">' +
                            '<div class="fw-bold text-truncate" style="max-width: 200px;" title="' + file.name + '">' + file.name + '</div>' +
                            '<small class="text-muted">' + fileSize + '</small>' +
                        '</div>' +
                        '<button type="button" class="btn btn-sm btn-outline-danger" onclick="removeFile(' + index + ')">' +
                            '<i class="bi bi-x"></i>' +
                        '</button>' +
                    '</div>' +
                '</div>' +
            '</div>';
        
        fileList.appendChild(fileItem);
    });
}

function getFileIcon(fileType) {
    if (fileType.startsWith('image/')) {
        return '<i class="bi bi-file-image text-primary"></i>';
    } else if (fileType.startsWith('video/')) {
        return '<i class="bi bi-file-play text-warning"></i>';
    } else if (fileType.startsWith('audio/')) {
        return '<i class="bi bi-file-music text-info"></i>';
    } else if (fileType === 'application/pdf') {
        return '<i class="bi bi-file-pdf text-danger"></i>';
    } else if (fileType.includes('document') || fileType.includes('sheet')) {
        return '<i class="bi bi-file-earmark-text text-success"></i>';
    } else {
        return '<i class="bi bi-file-earmark"></i>';
    }
}

function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

function removeFile(index) {
    // This would need more complex implementation to actually remove files from the FileList
    // For now, we'll just refresh the display
    fileInput.value = '';
    filePreview.style.display = 'none';
    uploadBtn.disabled = true;
    clearBtn.style.display = 'none';
}

function clearFiles() {
    fileInput.value = '';
    filePreview.style.display = 'none';
    uploadBtn.disabled = true;
    clearBtn.style.display = 'none';
    uploadStatus.innerHTML = '';
}

// Multiple file upload functionality
async function startMultipleUpload() {
    const files = fileInput.files;
    if (files.length === 0) {
        $.notify("No files selected", "error");
        return;
    }
    
    uploadBtn.disabled = true;
    clearBtn.disabled = true;
    uploadProgress.style.display = 'block';
    uploadStatus.innerHTML = '';
    
    let successCount = 0;
    let errorCount = 0;
    
    for (let i = 0; i < files.length; i++) {
        const file = files[i];
        const progress = Math.round(((i + 1) / files.length) * 100);
        
        uploadProgress.querySelector('.progress-bar').style.width = progress + '%';
        uploadProgress.querySelector('.progress-bar').textContent = progress + '%';
        
        const result = await uploadSingleFile(file);
        if (result.success) {
            successCount++;
        } else {
            errorCount++;
        }
    }
    
    // Final status
    if (errorCount === 0) {
        uploadStatus.innerHTML = '<div class="alert alert-success"><i class="bi bi-check-circle"></i> Successfully uploaded ' + successCount + ' files!</div>';
        $.notify('Successfully uploaded ' + successCount + ' files!', "success");
        
        // Refresh the page to show newly uploaded files
        setTimeout(() => {
            location.reload();
        }, 1500);
    } else {
        uploadStatus.innerHTML = '<div class="alert alert-warning"><i class="bi bi-exclamation-triangle"></i> Uploaded ' + successCount + ' files successfully, ' + errorCount + ' files failed.</div>';
        $.notify('Uploaded ' + successCount + ' files successfully, ' + errorCount + ' files failed.', "warning");
    }
    
    // Reset after 2 seconds
    setTimeout(() => {
        uploadBtn.disabled = false;
        clearBtn.disabled = false;
        uploadProgress.style.display = 'none';
        uploadProgress.querySelector('.progress-bar').style.width = '0%';
        uploadProgress.querySelector('.progress-bar').textContent = '0%';
    }, 2000);
}

async function uploadSingleFile(file) {
    const formData = new FormData();
    formData.append('action', 'file_upload');
    formData.append('current_dir', '` + currentDirectory + `');
    formData.append('upload_file', file);
    
    try {
        const response = await fetch("` + url + `", {
            method: 'POST',
            body: formData
        });
        
        const result = await response.json();
        
        if (result.status === "success") {
            return { success: true, message: result.message };
        } else {
            return { success: false, message: result.message };
        }
    } catch (error) {
        return { success: false, message: "Upload failed: " + error.message };
    }
}

// Initialize modal event handlers
document.addEventListener('DOMContentLoaded', function() {
    const modal = document.getElementById('ModalUploadMultipleFiles');
    if (modal) {
        modal.addEventListener('hidden.bs.modal', function() {
            clearFiles();
        });
    }
});
</script>
<!-- END: Modal Upload Multiple Files -->
	`
}
