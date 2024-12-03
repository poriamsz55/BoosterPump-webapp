import { HTTP_URL } from '../config.js';

class DevicesManager {
    constructor() {
        this.devices = [];
        this.devicesGrid = document.getElementById('devicesGrid');
        this.searchInput = document.getElementById('searchDevices');

        this.init();
        
        // Add focus event listener
        window.addEventListener('focus', () => this.checkForUpdates());
    }

    async checkForUpdates() {
        if (localStorage.getItem('devicesListNeedsUpdate') === 'true') {
            await this.getDevicesFromDB();
            this.renderDevices(this.devices);
            localStorage.removeItem('devicesListNeedsUpdate');
        }
    }

    async init() {
        console.log('Initializing DevicesManager...');
        await this.getDevicesFromDB();
        this.renderDevices(this.devices);
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.searchInput.addEventListener('input', (e) => this.handleSearch(e));
        document.getElementById('addDeviceToDBButton').addEventListener('click', () =>{
            window.location.href = '/add/device/db';
        });
        // Remove the form submit event listener from here
        this.form.addEventListener('submit', async (e) => {
            e.preventDefault(); // Prevent form submission
            await this.saveDevice();
        });
    }


    renderDevices(devicesList) {
        this.devicesGrid.innerHTML = devicesList.map(device => `
            <div class="card" data-id="${device.id}">
                <div class="card-header">
                    <span class="card-title">${this.escapeHtml(device.name)}</span>
                </div>
                <div class="card-price">${this.formatPrice(device.price)}</div>
                <div class="card-actions">
                    <button class="action-button delete-btn" data-id="${device.id}">
                        <i class="fas fa-trash"></i>
                    </button>
                    <button class="action-button copy-btn" data-id="${device.id}">
                        <i class="fas fa-copy"></i>
                    </button>
                </div>
            </div>
        `).join('');

        this.attachCardEventListeners();
    }

    attachCardEventListeners() {
        this.devicesGrid.querySelectorAll('.card').forEach(card => {
            card.addEventListener('click', (e) => {
                if (!e.target.closest('.action-button')) {
                    const deviceId = card.getAttribute('data-id');
                    this.handleDeviceClick(deviceId);
                }
            });
        });

        this.devicesGrid.querySelectorAll('.delete-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.deleteDevice(btn.dataset.id);
            });
        });

        this.devicesGrid.querySelectorAll('.copy-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.copyDevice(btn.dataset.id);
            });
        });
    }

    handleDeviceClick(id) {
        window.location.href = `/devices/details?id=${id}`;
    }

    handleSearch(e) {
        const searchTerm = e.target.value.toLowerCase();
        const filteredDevices = this.devices.filter(device =>
            device.name.toLowerCase().includes(searchTerm)
        );
        this.renderDevices(filteredDevices);
    }

    async getDevicesFromDB() {
        try {
            const response = await fetch(`${HTTP_URL}/device/getAll`);
            const data = await response.json();
            console.log(data);

            if (Array.isArray(data)) {
                this.devices = data;
            } else {
                console.error('Invalid response format:', data);
                this.devices = [];
            }
        } catch (error) {
            console.error('Error fetching devices:', error);
            this.devices = [];
        }
    }

    async deleteDevice(id) {
        if (confirm('Are you sure you want to delete this device?')) {
            try {
                const response = await fetch(`${HTTP_URL}/device/delete/${id}`, {
                    method: 'POST'
                });

                if (response.ok) {
                    await this.getDevicesFromDB();
                    this.renderDevices(this.devices);
                } else {
                    alert('Failed to delete device.');
                }
            } catch (error) {
                console.error('Error deleting device:', error);
                alert('An error occurred while deleting the device.');
            }
        }
    }

    async copyDevice(id) {
        const deviceToCopy = this.devices.find(device => device.id === id);
        if (deviceToCopy) {
            // Implementation for copying device
            console.log('Copying device:', deviceToCopy);
        }
    }

    escapeHtml(unsafe) {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }

    formatPrice(price) {
        return new Intl.NumberFormat('fa-IR').format(price);
    }
}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    new DevicesManager();
});