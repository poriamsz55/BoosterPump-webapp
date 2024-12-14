document.getElementById("btnProjectPage").addEventListener("click", function() {
    window.location.href = "/projects";
});

document.getElementById("btnDevicePage").addEventListener("click", function() {
    window.location.href = "/devices";
});

document.getElementById("btnPartPage").addEventListener("click", function() {
    window.location.href = "/parts";
});

document.getElementById("btnDownloadDB").addEventListener("click", async function() {
    // Get the button element
    const btn = this;

    try {
        // Ask for filename using prompt
        let fileName = prompt("Enter file name for the database:", "booster_pump.db");
        
        if (fileName === null) return;

        // Remove any slashes and ensure .db extension
        fileName = fileName.replace(/[/\\]/g, '');
        if (!fileName.endsWith('.db')) {
            fileName += '.db';
        }

        // Show loading state
        btn.disabled = true;

        // Create form data
        const formData = new FormData();
        formData.append('name', fileName);

        // Make the request
        const response = await fetch('http://localhost:8080/api/database/download', {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            const text = await response.text();
            throw new Error(text || 'Failed to download database');
        }

        const blob = await response.blob();
        
        // Create download link
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.style.display = 'none';
        a.href = url;
        a.download = fileName;

        // Add to document and trigger download
        document.body.appendChild(a);
        a.click();

        // Cleanup
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);

    } catch (error) {
        alert('Error: ' + error.message);
    } finally {
        // Reset button state
        btn.disabled = false;
    }
});


// Your updated event listener
document.getElementById("btnUploadDB").addEventListener("click", async function() {
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = '.db';
    
    const btn = this;

    fileInput.addEventListener('change', async function() {
        if (!fileInput.files || !fileInput.files[0]) {
            return;
        }
    
        try {
            const result = await customConfirm(
                "نحوه بارگزاری دیتابیس را انتخاب کنید:\n" +
                "جایگزینی = حذف داده های قبلی و جایگزینی داده های جدید\n" +
                "افزودن = حفظ داده های قبلی و اضافه کردن داده های جدید"
            );
    
            // If user cancelled or clicked outside
            if (result === 'cancel') {
                return;
            }
    
            btn.disabled = true;
    
            const formData = new FormData();
            formData.append('database', fileInput.files[0]);
            formData.append('replace', result === 'replace');
    
            const response = await fetch('/api/database/upload', {
                method: 'POST',
                body: formData
            });
    
            if (!response.ok) {
                const text = await response.text();
                throw new Error(text || 'Failed to upload database');
            }
    
            const responseText = await response.text();
            alert('Success: ' + responseText);
    
            if (result === 'replace') {
                window.location.reload();
            }
    
        } catch (error) {
            alert('Error: ' + error.message);
        } finally {
            btn.disabled = false;
            fileInput.value = '';
        }
    });

    fileInput.click();
});

function customConfirm(message) {
    return new Promise((resolve) => {
        const overlay = document.createElement('div');
        overlay.className = 'confirm-overlay';

        const dialog = document.createElement('div');
        dialog.className = 'confirm-dialog';

        const messageEl = document.createElement('div');
        messageEl.className = 'confirm-message';
        messageEl.textContent = message;

        const buttonsContainer = document.createElement('div');
        buttonsContainer.className = 'confirm-buttons';

        const cancelButton = document.createElement('button');
        cancelButton.className = 'confirm-button cancel';
        cancelButton.textContent = 'انصراف';

        const mergeButton = document.createElement('button');
        mergeButton.className = 'confirm-button';
        mergeButton.textContent = 'افزودن';

        const replaceButton = document.createElement('button');
        replaceButton.className = 'confirm-button replace';
        replaceButton.textContent = 'جایگزینی';

        const closeDialog = (result) => {
            document.body.removeChild(overlay);
            resolve(result);
        };

        replaceButton.onclick = () => closeDialog('replace');
        mergeButton.onclick = () => closeDialog('merge');
        cancelButton.onclick = () => closeDialog('cancel');
        
        // Close on overlay click
        overlay.onclick = (e) => {
            if (e.target === overlay) {
                closeDialog('cancel');
            }
        };

        // Close on Escape key
        document.addEventListener('keydown', function(e) {
            if (e.key === 'Escape' && document.contains(overlay)) {
                closeDialog('cancel');
            }
        });

        buttonsContainer.appendChild(cancelButton);
        buttonsContainer.appendChild(mergeButton);
        buttonsContainer.appendChild(replaceButton);
        dialog.appendChild(messageEl);
        dialog.appendChild(buttonsContainer);
        overlay.appendChild(dialog);
        document.body.appendChild(overlay);

        replaceButton.focus();
    });
}


// esc key
function handleEscKey(callback) {
    document.addEventListener('keydown', function(event) {
        if (event.key === 'Escape' || event.key === 'Esc' || event.keyCode === 27) {
            callback();
        }
    });
}

// use esc key
handleEscKey(function() {
    // if custom confirm dialog is open, close it
    if (document.querySelector('.confirm-overlay')) {
        document.querySelector('.confirm-overlay').remove();
    }
});