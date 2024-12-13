import { HTTP_URL } from '../config.js';
import { formatPriceValue } from '../format-price.js';
import { converterToString, filterToString } from '../convert2str.js';
import { handleEscKey } from '../keyboard-utils.js';

class AddProjectManager {
    constructor() {
        this.devices = [];
        this.addedDevices = [];
        this.devicesGrid = document.getElementById('devicesGrid');
        this.searchInput = document.getElementById('searchDevices');
        this.modal = document.getElementById('addDeviceToProjectModal');
        this.form = document.getElementById('addProjectDBForm')
        this.deviceForm = document.getElementById('addDeviceToProjectForm');
        
        this.hasChanged = false;

        // check projectName div is changed
        const projectNameDiv = document.getElementById('projectName');
        projectNameDiv.addEventListener('input', () => {
            this.hasChanged = true;
        });

        this.init();

        handleEscKey(() => {
            if (window.getComputedStyle(this.modal).display !== 'none') {
                this.closeModal();
            } else {
                this.handleBackButton();
            }
        });
    }

    async init() {
        this.closeModal();
        await this.getDevicesFromDB();
        this.renderDevices(this.devices);
        this.renderAddedDevices(this.addedDevices);
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.searchInput.addEventListener('input', (e) => this.handleSearch(e));
        document.getElementById('addDeviceToProjectBtn').addEventListener('click', () => this.openModal());
        document.getElementById('cancelBtn').addEventListener('click', () => this.closeModal());
        document.getElementById('addProjectDBBtn').addEventListener('click', () => this.saveProject());
        document.getElementById('backBtn').addEventListener('click', () => this.handleBackButton());

       
    }

    // handle back button
    handleBackButton() {
        if (this.hasChanged) {
            if (confirm('You have unsaved changes. Are you sure you want to leave?')) {
                localStorage.removeItem('projectDevices');
                this.devices = [];
                this.addedDevices = [];
                this.renderAddedDevices([]);
                window.history.back();
            }
        }else{
            localStorage.removeItem('projectDevices');
            this.devices = [];
            this.addedDevices = [];
            this.renderAddedDevices([]);
            window.history.back();
        }
    }

    renderDevices(devicesList) {

        this.devicesGrid.innerHTML = devicesList.map(device => {
            const converterStr = converterToString(device.converter);
            const filterStr = filterToString(device.filter);
            return `
                    <div class="card" data-id="${device.id}">
                          <div class="card-header">
                                <span class="card-title">${this.escapeHtml(device.name)}</span>
                            </div>
                             <div class="card-header">
                                <div class="card-title">نوع تبدیل: ${converterStr}</div>
                            </div>
                            <div class="card-header">
                                <div class="card-title">صافی ${filterStr}</div>
                            </div>
                            <div class="card-header">
                                <div class="card-price">قیمت: ${formatPriceValue(device.price)}</div>
                            </div>
                    </div>
                `;
        }).join('');

        this.attachCardEventListeners();
    }


    attachCardEventListeners() {

        this.devicesGrid.querySelectorAll('.card').forEach(card => {
            if (card.hasEventListener) return; // Prevent duplicate listeners
            card.hasEventListener = true;

            card.addEventListener('click', (e) => {
                if (card.classList.contains('disabled') ||
                    card.classList.contains('selected') ||
                    e.target.closest('.action-button')) {
                    return;
                }

                this.devicesGrid.querySelectorAll('.card').forEach(c => {
                    if (c.classList.contains('selected')) {
                        c.classList.remove('selected');
                        const deviceId = c.getAttribute('data-id');
                        this.removeInputsFromCard(c, deviceId);
                    }
                });
                card.classList.add('selected');

                const deviceId = card.getAttribute('data-id');
                this.addInputsToCard(card, deviceId);
            });
        });
    }

    addInputsToCard(card, deviceId) {
        // Remove existing inputs if any
        card.querySelectorAll('input, button').forEach(el => el.remove());

        const countInput = document.createElement('input');
        countInput.type = 'number';
        countInput.value = 1;
        countInput.min = '1';
        countInput.id = `count-${deviceId}`;

        const addButton = document.createElement('button');
        addButton.type = 'button';
        addButton.textContent = 'افزودن به دستگاه';
        addButton.id = `add-to-project-${deviceId}`;

        addButton.addEventListener('click', () =>
            this.addToProject(deviceId, countInput.value));

        card.appendChild(countInput);
        card.appendChild(addButton);
    }

    removeInputsFromCard(card) {
        // Remove existing inputs if any
        card.querySelectorAll('input, button').forEach(el => el.remove());
    }

    // render added devices in modal
    renderAddedDevices(addedDevices) {
        const devicesGrid = document.getElementById('addedDevicesGrid');
        devicesGrid.innerHTML = '';
        addedDevices.forEach(device => {
            let founded = this.devices.find(p => p.id.toString() === device.id);
            if (!founded) return;

            const deviceCard = document.createElement('div');
            deviceCard.classList.add('card');
            deviceCard.setAttribute('data-device-id', device.id);
            const converterStr = converterToString(device.converter);
            const filterStr = filterToString(device.filter);
            deviceCard.innerHTML = `
                            <div class="card-header">
                                <span class="card-title">${this.escapeHtml(founded.name)}</span>
                            </div>
                             <div class="card-header">
                                <div class="card-title">نوع تبدیل: ${converterStr}</div>
                            </div>
                            <div class="card-header">
                                <div class="card-title">صافی ${filterStr}</div>
                            </div>
                            <div class="card-header">
                                <div class="card-price">قیمت: ${formatPriceValue(founded.price)}</div>
                            </div>
                            <div class="card-header">
                                <div class="card-count-container">
                                    <div class="card-count-title">تعداد:</div>
                                    <button class="count-btn minus-btn" id="minus-${device.id}">-</button>
                                    <div class="card-count">${device.count}</div>
                                    <button class="count-btn plus-btn" id="plus-${device.id}">+</button>
                                </div>
                            </div>
                <button type="button" class="action-button delete-btn" data-id="delete-${device.id}">
                    <i class="fas fa-trash"></i>
                </button>
            `;
            devicesGrid.appendChild(deviceCard);
        });

        this.attachAddedCardEventListeners();
    }

    // attachAddedCardEventListeners
    attachAddedCardEventListeners() {
        const devicesGrid = document.getElementById('addedDevicesGrid');

        devicesGrid.querySelectorAll('.count-btn').forEach(button => {
            button.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                this.hasChanged = true;
                const deviceId = button.id.split('-')[1];
                const device = this.addedDevices.find(device => device.id.toString() === deviceId.toString());
                if (!device) return;

                if (button.classList.contains('minus-btn')) {
                    device.count = Math.max(device.count - 1, 1);
                } else if (button.classList.contains('plus-btn')) {
                    // convert device.count to number
                    device.count = (parseInt(device.count) + 1).toString();
                }

                // get localstorage and update it
                const projectDevices = JSON.parse(localStorage.getItem('projectDevices')) || [];
                const existingDevice = projectDevices.find(device => device.id.toString() === deviceId.toString());
                if (existingDevice) {
                    existingDevice.count = device.count;
                } else {
                    projectDevices.push(device);
                }
                localStorage.setItem('projectDevices', JSON.stringify(projectDevices));

                // update card count
                const card = devicesGrid.querySelector(`[data-device-id="${deviceId}"]`);
                if (card) {
                    card.querySelector('.card-count').textContent = `${device.count}`;
                }
            });
        });

        devicesGrid.querySelectorAll('.delete-btn').forEach(button => {
            button.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                // get id data-id="delete-${device.id}
                const deviceId = button.dataset.id.replace('delete-', '');
                this.hasChanged = true;
                this.deleteDevice(deviceId);
            });
        });
    }

    deleteDevice(deviceId) {
        this.addedDevices = this.addedDevices.filter(device => device.id !== deviceId);
        localStorage.setItem('projectDevices', JSON.stringify(this.addedDevices));

        // Remove just the specific card instead of re-rendering everything
        const cardToRemove = document.querySelector(`[data-device-id="${deviceId}"]`);
        if (cardToRemove) {
            cardToRemove.remove();
        }

        // Update the disabled state of the corresponding device in the devices grid
        const originalCard = this.devicesGrid.querySelector(`[data-id="${deviceId}"]`);
        if (originalCard) {
            originalCard.classList.remove('disabled');
            originalCard.classList.remove('selected');
        }
    }


    handleSearch(e) {
        const searchTerm = e.target.value.toLowerCase();
        const filteredDevices = this.devices.filter(device =>
            device.name.toLowerCase().includes(searchTerm)
        );
        this.renderDevices(filteredDevices);
    }

    openModal() {
        this.modal.style.display = 'flex';
        this.deviceForm.reset();
    }

    closeModal() {
        // remove selected class from all cards
        this.devicesGrid.querySelectorAll('.card').forEach(card => {
            card.classList.remove('selected');
        });
        this.modal.style.display = 'none';
        this.deviceForm.reset();
    }

    async addToProject(deviceId, count) {

        // remove selected class from card
        const card = this.devicesGrid.querySelector(`.card[data-id="${deviceId}"]`);
        card.classList.remove('selected');

        // add disabled class to card
        card.classList.add('disabled');

        // remove count input field
        const countInput = document.getElementById(`count-${deviceId}`);
        countInput.remove();

        // remove add to project button
        const addToProjectBtn = document.getElementById(`add-to-project-${deviceId}`);
        addToProjectBtn.remove();

        const projectDevice = {
            id: deviceId,
            count: count
        };

        const projectDevices = JSON.parse(localStorage.getItem('projectDevices')) || [];
        projectDevices.push(projectDevice);
        localStorage.setItem('projectDevices', JSON.stringify(projectDevices));

        this.hasChanged = true;

        this.addedDevices.push(projectDevice);
        this.renderAddedDevices(this.addedDevices);
    }


    async getDevicesFromDB() {
        try {
            const response = await fetch(`${HTTP_URL}/device/getAll`);
            const data = await response.json();

            if (Array.isArray(data)) {
                this.devices = data;
                this.devices.reverse();
            } else {
                console.error('Invalid response format:', data);
                this.devices = [];
            }
        } catch (error) {
            console.error('Error fetching devices:', error);
            this.devices = [];
        }
    }

    async saveProject() {
        const formData = new FormData();
        formData.append('projectName', document.getElementById('projectName').value);

        // check if project name is empty or contains only spaces
        if (document.getElementById('projectName').value.trim() === '') {
            alert('Please enter a project name.');
            return;
        }

        try {
            const response = await fetch(`${HTTP_URL}/project/add`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) throw new Error('Failed to save project');

            const projectId = await response.json();
            const projectDevices = JSON.parse(localStorage.getItem('projectDevices')) || [];
            const deviceFormData = new FormData();
            deviceFormData.append('projectId', projectId.id);
            deviceFormData.append('devices', JSON.stringify(projectDevices));

            const deviceResponses = await fetch(`${HTTP_URL}/projectDevice/add/list`, {
                method: 'POST',
                body: deviceFormData,
            });

            if (!deviceResponses.ok) throw new Error('Failed to save project devices');

            localStorage.removeItem('projectDevices');

            // Set update flag in localStorage
            localStorage.setItem('projectsListNeedsUpdate', 'true');

            window.history.back();
        } catch (error) {
            console.error('Error:', error);
            alert('An error occurred. Please try again later.');
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
    new AddProjectManager();
});