// Sample devices data
const devices = [
    { id: 1, name: "Device 1", price: "$1000" },
    { id: 2, name: "Device 2", price: "$2000" },
    { id: 3, name: "Device 3", price: "$3000" },
];

// Initialize the devices page
function initDevicesPage() {
    const devicesGrid = document.getElementById('devicesGrid');
    const searchInput = document.getElementById('searchDevices');
    const deviceDetails = document.getElementById('deviceDetails');

    function renderDevices(devicesList) {
        devicesGrid.innerHTML = devicesList.map(device => `
            <div class="device-card" data-id="${device.id}">
                <div class="device-header">
                    <span class="device-title">${device.name}</span>
                </div>
                <div class="device-price">${device.price}</div>
                <div class="device-actions">
                    <button class="action-button" onclick="deleteDevice(${device.id})">
                        <i class="fas fa-trash"></i>
                    </button>
                    <button class="action-button" onclick="copyDevice(${device.id})">
                        <i class="fas fa-copy"></i>
                    </button>
                </div>
            </div>
        `).join('');

        // Add event listeners to each device card for click handling
        devicesGrid.querySelectorAll('.device-card').forEach(card => {
            card.addEventListener('click', (e) => {
                const deviceId = card.getAttribute('data-id');
                handleDeviceClick(deviceId);
            });
        });
    }

    
    function handleDeviceClick(id) {
        // Redirect to device-details.html with the device ID as a query parameter
        window.location.href = `/templates/device-details.html?id=${id}`;
    }

    // Initial render
    renderDevices(devices);

    // Search functionality
    searchInput.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        const filteredDevices = devices.filter(device => 
            device.name.toLowerCase().includes(searchTerm)
        );
        renderDevices(filteredDevices);
    });

    // Add device button
    document.getElementById('addDeviceButton').addEventListener('click', () => {
        // Implement add device functionality
        console.log('Add device clicked');
    });
}

// Navigation functions
function showPage(pageId) {
    document.querySelectorAll('.page').forEach(page => {
        page.classList.remove('active');
    });
    document.getElementById(pageId).classList.add('active');
}

// Device actions
function deleteDevice(id) {
    console.log('Delete device:', id);
    // Implement delete functionality
}

function copyDevice(id) {
    console.log('Copy device:', id);
    // Implement copy functionality
}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    initDevicesPage();
    showPage('devicesPage'); // Show devices page by default
});
