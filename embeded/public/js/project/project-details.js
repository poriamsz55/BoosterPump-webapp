import { HTTP_URL } from '../config.js';
import { formatPrice } from '../format-price.js';

document.getElementById('projectDetailsForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const urlParams = new URLSearchParams(window.location.search);
    const projectId = urlParams.get('id');

    const formData = new FormData();
    formData.append('projectId', projectId);
    formData.append('projectName', document.getElementById('projectName').value);

    formData.append('projectConverter', document.getElementById('projectConverter').value);
    formData.append('projectFilter', document.getElementById('projectFilter').value);

    try {
        const response = await fetch(`${HTTP_URL}/project/update`, {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error('Failed to update project details');
        }

        // Set update flag in localStorage
        localStorage.setItem('projectsListNeedsUpdate', 'true');

        alert('اطلاعات با موفقیت ذخیره شد');
        window.history.back();

    } catch (error) {
        console.error('Error updating project details:', error);
        alert(error.message || 'An error occurred while saving. Please try again later.');
    }
});


// Project Details Manager
class AddProjectDetailsManager {
    constructor() {
        this.devices = [];
        this.addedDevices = [];
        this.devicesGrid = document.getElementById('devicesGrid');
        this.searchInput = document.getElementById('searchDevices');
        this.modal = document.getElementById('addDeviceToProjectModal');
        this.deviceForm = document.getElementById('addDeviceToProjectForm');
        // Remove e.preventDefault() as it's not needed here
        const urlParams = new URLSearchParams(window.location.search);
        this.projectId = urlParams.get('id');
        // Add event listener for price input
        this.priceInput = document.getElementById('projectPrice');
        this.priceInput.addEventListener('input', function () {
            formatPrice(this);
        });

        this.init();

        // Add this to handle back button
        window.addEventListener('popstate', () => {
            localStorage.removeItem('projectDevices');
            this.devices = [];
            this.addedDevices = [];
            this.renderAddedDevices([]);
        });
    }

    async init() {
        this.closeModal();
        await this.getDevicesFromDB();
        await this.getProjectDetails();
        this.renderAddedDevices(this.addedDevices);
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.searchInput.addEventListener('input', (e) => this.handleSearch(e));
        document.getElementById('addDeviceToProjectBtn').addEventListener('click', () => this.openModal());
        document.getElementById('cancelBtn').addEventListener('click', () => this.closeModal());
        document.getElementById('saveProjectDBBtn').addEventListener('click', () => this.saveProject());
    }

    async saveProject() {
        const formData = new FormData();
        formData.append('projectName', document.getElementById('projectName').value);
        formData.append('converterType', document.getElementById('converterType').value); // Fixed field name
        formData.append('filter', document.getElementById('filterCheckbox').checked); // Fixed field name
        formData.append('devices', JSON.stringify(projectDevices));

        try {
            const response = await fetch(`${HTTP_URL}/project/update`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) throw new Error('Failed to save project');

            localStorage.removeItem('projectDevices');

            // Set update flag in localStorage
            localStorage.setItem('projectsListNeedsUpdate', 'true');

            window.history.back();
        } catch (error) {
            console.error('Error:', error);
            alert('An error occurred. Please try again later.');
        }
    }

    async getProjectDetails() {
        if (this.projectId) {
            try {
                const formData = new FormData();
                formData.append('projectId', this.projectId);

                const response = await fetch(`${HTTP_URL}/project/getById`, {
                    method: 'POST',
                    body: formData
                });

                if (!response.ok) {
                    throw new Error('Failed to fetch project details');
                }

                const projectDetails = await response.json();
                console.log(projectDetails)

                if (projectDetails) {
                    this.addedDevices = projectDetails.project_device;
                    // add addedDevices to localStorage
                    localStorage.setItem('projectDevices', JSON.stringify(this.addedDevices));

                    // Use value instead of textContent for input elements
                    document.getElementById('projectName').value = projectDetails.name || '';

                    // converter
                    const converterSelect = document.getElementById('converterType');
                    const converterOption = converterSelect.querySelector(`option[value="${projectDetails.converter}"]`);
                    if (converterOption) {
                        converterOption.selected = true;
                    }

                    // filter check box
                    const filterCheckbox = document.getElementById('filterCheckbox');
                    filterCheckbox.checked = projectDetails.filter === 'true';

                    // Format the initial price value
                    if (projectDetails.price) {
                        this.priceInput.value = projectDetails.price;
                        formatPrice(this.priceInput);
                    }
                } else {
                    throw new Error('Project details not found');
                }

            } catch (error) {
                console.error('Error fetching project details:', error);
                alert(error.message || 'An error occurred. Please try again later.');
            }
        }
    }

    renderDevices(devicesList) {
        console.log("this.addedDevices => ", this.addedDevices);
        console.log("devicesList => ", devicesList);

        this.devicesGrid.innerHTML = devicesList.map(device => {

            // check if device is already added 
            let added = false;
            for (let i = 0; i < this.addedDevices.length; i++) {
                console.log("this.addedDevices[i].device.id => ", this.addedDevices[i].device.id);
                console.log("device.id => ", device.id);
                if (this.addedDevices[i].device.id === device.id) {
                    added = true;
                    break;
                }
            }

            console.log("added => ", added);
            if (added) {
                return `
                        <div class="card disabled" data-id="${device.id}">
                            <div class="card-header">
                                <span class="card-title">${this.escapeHtml(device.name)}</span>
                            </div>
                            <div class="card-price">${this.formatPrice(device.price)}</div>
                        </div>
                    `;
                }else{
                    return `
                        <div class="card" data-id="${device.id}">
                            <div class="card-header">
                                <span class="card-title">${this.escapeHtml(device.name)}</span>
                            </div>
                            <div class="card-price">${this.formatPrice(device.price)}</div>
                        </div>
                    `;
                }
            }

        ).join(''); 
        
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

        this.addedDevices.push(projectDevice);
        this.renderAddedDevices(this.addedDevices);
    }


    removeInputsFromCard(card) {
        // Remove existing inputs if any
        card.querySelectorAll('input, button').forEach(el => el.remove());
    }

    // render added devices in modal
    renderAddedDevices(addedDevices) {
        console.log(addedDevices)
        const devicesGrid = document.getElementById('addedDevicesGrid');
        devicesGrid.innerHTML = '';
        addedDevices.forEach(device => {

            const deviceCard = document.createElement('div');
            deviceCard.classList.add('card');
            deviceCard.setAttribute('data-device-id', device.id);
            deviceCard.innerHTML = `
                <div class="card-title">${this.escapeHtml(device.device.name)}</div>
                <div class="card-price">${this.formatPrice(device.price)}</div>
                <div class="card-count">${device.count}</div>
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
        devicesGrid.querySelectorAll('.delete-btn').forEach(button => {
            button.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                // get id data-id="delete-${device.id}
                const deviceId = button.dataset.id.replace('delete-', '');
                console.log(deviceId);
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


    async getDevicesFromDB() {
        try {
            const response = await fetch(`${HTTP_URL}/device/getAll`);
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

    handleSearch(e) {
        const searchTerm = e.target.value.toLowerCase();
        const filteredDevices = this.devices.filter(device =>
            device.name.toLowerCase().includes(searchTerm)
        );
        this.renderDevices(filteredDevices);
    }

    openModal() {
        this.renderDevices(this.devices);
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
    new AddProjectDetailsManager();
});