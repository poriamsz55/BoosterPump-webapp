// Export
export function exportProject() {
    // Get the project details from local storage
    const projectDetails = JSON.parse(localStorage.getItem('projectDetails'));
    const projectDevices = JSON.parse(localStorage.getItem('projectDevices'));

    // Create a JSON object with the project details and devices
    const projectData = {
        projectDetails,
        projectDevices
    };

    // Convert the project data to a JSON string
    const jsonData = JSON.stringify(projectData);

    // Create a Blob with the JSON data
    const blob = new Blob([jsonData], { type: 'application/json' });

    // Create a URL for the Blob
    const url = URL.createObjectURL(blob);

}