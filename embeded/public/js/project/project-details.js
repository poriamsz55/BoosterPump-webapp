// Sample data
const devices = [
    {
        name: "دستگاه ۱",
        converter: "مبدل A",
        filter: "فیلتر X",
        count: 2,
        price: "۵۰۰,۰۰۰ تومان"
    },
    {
        name: "دستگاه ۲",
        converter: "مبدل B",
        filter: "فیلتر Y",
        count: 1,
        price: "۱,۰۰۰,۰۰۰ تومان"
    }
];

const devicesContainer = document.querySelector('.devices-container');

// Render devices
function renderDevices() {
    devicesContainer.innerHTML = devices.map((device, index) => `
        <div class="device-card">
            <div class="device-info">
                <div class="device-title">${device.name}</div>
                <div class="device-detail">
                    <span>مبدل:</span>
                    <span>${device.converter}</span>
                </div>
                <div class="device-detail">
                    <span>فیلتر:</span>
                    <span>${device.filter}</span>
                </div>
                <div class="device-detail">
                    <span>تعداد:</span>
                    <div class="count-controls">
                        <button class="count-btn" onclick="decrementCount(${index})">
                            <i class="fas fa-minus"></i>
                        </button>
                        <span class="count-display">${device.count}</span>
                        <button class="count-btn" onclick="incrementCount(${index})">
                            <i class="fas fa-plus"></i>
                        </button>
                    </div>
                </div>
                <div class="device-detail">
                    <span>قیمت:</span>
                    <span>${device.price}</span>
                </div>
                <button class="delete-btn" onclick="deleteDevice(${index})">
                    <i class="fas fa-trash" style="color: white;"></i>
                </button>
            </div>
        </div>
    `).join('');
}

function deleteDevice(index) {
    devices.splice(index, 1);
    renderDevices();
}

function incrementCount(index) {
    devices[index].count += 1;
    renderDevices();
    animateCount(index);
}

function decrementCount(index) {
    if (devices[index].count > 1) {
        devices[index].count -= 1;
        renderDevices();
        animateCount(index);
    }
}

function animateCount(index) {
    const countDisplays = document.querySelectorAll('.count-display');
    countDisplays[index].classList.add('pulse');
    setTimeout(() => {
        countDisplays[index].classList.remove('pulse');
    }, 300);
}

function getProjectIdFromUrl() {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get('id');
}

function openModal() {
    document.getElementById('addDeviceModal').style.display = 'flex';
}

function closeModal() {
    document.getElementById('addDeviceModal').style.display = 'none';
}

function handleAddDevice(event) {
    event.preventDefault();
    
    const newDevice = {
        name: document.getElementById('deviceName').value,
        converter: document.getElementById('converterType').value,
        filter: document.getElementById('filterType').value,
        count: parseInt(document.getElementById('deviceCount').value),
        price: "۵۰۰,۰۰۰ تومان" // You might want to calculate this based on the selected device
    };

    devices.push(newDevice);
    renderDevices();
    closeModal();
    event.target.reset();
}

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('addDeviceModal');
    if (event.target === modal) {
        closeModal();
    }
}

// Initialize the page
document.addEventListener('DOMContentLoaded', () => {
    const projectId = getProjectIdFromUrl();
    renderDevices();
});