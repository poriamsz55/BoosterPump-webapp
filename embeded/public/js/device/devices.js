import { HTTP_URL } from '../config.js';
import { formatPriceValue } from '../format-price.js';

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
        await this.getDevicesFromDB();
        this.renderDevices(this.devices);
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.searchInput.addEventListener('input', (e) => this.handleSearch(e));
        document.getElementById('addDeviceToDBButton').addEventListener('click', () => {
            window.location.href = '/add/device/db';
        });
    }


    renderDevices(devicesList) {
        this.devicesGrid.innerHTML = devicesList.map(device => `
            <div class="card" data-id="${device.id}">
                <div class="card-header">
                    <span class="card-title">${this.escapeHtml(device.name)}</span>
                </div>
                <div class="card-price">${formatPriceValue(device.price)}</div>
                <div class="card-actions">
                    <button class="action-button delete-btn" data-id="delete-${device.id}">
                        <i class="fas fa-trash"></i>
                    </button>
                    <button class="action-button copy-btn" data-id="copy-${device.id}">
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
                this.deleteDevice(btn.dataset.id.toString().replace('delete-', ''));
            });
        });

        this.devicesGrid.querySelectorAll('.copy-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.copyDevice(btn.dataset.id.toString().replace('copy-', ''));
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

            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            const data = await response.json();

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
                const formData = new FormData();
                formData.append('deviceId', id); 
                const response = await fetch(`${HTTP_URL}/device/delete`, {
                    method: 'POST',
                    body: formData
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
        const deviceToCopy = this.devices.find(device => device.id.toString() === id.toString());
        if (deviceToCopy) {
            try {
                const formData = new FormData();
                formData.append('deviceId', id);

                const response = await fetch(`${HTTP_URL}/device/copy`, {
                    method: 'POST',
                    body: formData
                });

                if (response.ok) {
                    await this.getDevicesFromDB();
                    this.renderDevices(this.devices);
                } else {
                    alert('Failed to copy device.');
                }
            } catch (error) {
                console.error('Error copying device:', error);
                alert('An error occurred while copying the device.');
            }
        } else {
            alert('Device not found.');
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

}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    new DevicesManager();
});