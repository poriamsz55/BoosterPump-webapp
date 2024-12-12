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

document.getElementById("btnUploadDB").addEventListener("click", async function() {
    // Create file input element
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = '.db';
    
    // Create and style the button
    const btn = this;

    // When a file is selected
    fileInput.addEventListener('change', async function() {
        if (!fileInput.files || !fileInput.files[0]) {
            return;
        }

        try {
            // Ask user whether to replace or merge
            const isReplace = confirm(
                "Choose database upload mode:\n" +
                "OK = Replace existing database\n" +
                "Cancel = Merge with existing database"
            );

            // Show loading state
            btn.disabled = true;

            // Create form data
            const formData = new FormData();
            formData.append('database', fileInput.files[0]);
            formData.append('replace', isReplace);

            // Make the request
            const response = await fetch('/api/database/upload', {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                const text = await response.text();
                throw new Error(text || 'Failed to upload database');
            }

            const result = await response.text();
            alert('Success: ' + result);

            // Optionally reload the page if needed
            if (isReplace) {
                window.location.reload();
            }

        } catch (error) {
            alert('Error: ' + error.message);
        } finally {
            // Reset button state
            btn.disabled = false;
            
            // Clean up
            fileInput.value = '';
        }
    });

    // Trigger file picker
    fileInput.click();
});