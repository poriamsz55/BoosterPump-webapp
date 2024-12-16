import { HTTP_URL } from '../config.js';
import { convertPriceToNumber, formatPriceValue, formatPriceInput } from '../format-price.js';
import { ProjectDevice } from './project-device.js';
import { converterToString, filterToString } from '../convert2str.js';
import { handleEscKey } from '../keyboard-utils.js';

// Project Details Manager
class AddProjectDetailsManager {
    constructor() {

        this.devices = [];
        this.addedDevices = [];
        this.devicesGrid = document.getElementById('devicesGrid');
        this.searchInput = document.getElementById('searchDevices');
        this.modal = document.getElementById('addDeviceToProjectModal');
        // Remove e.preventDefault() as it's not needed here
        const urlParams = new URLSearchParams(window.location.search);
        this.projectId = urlParams.get('id');
        // Add event listener for price input
        this.priceInput = document.getElementById('projectPrice');
        this.priceInput.addEventListener('input', function () {
            formatPriceInput(this);
        });

        this.hasChanged = false;

        // check projectName div is changed
        const projectNameDiv = document.getElementById('projectName');
        projectNameDiv.addEventListener('input', () => {
            this.hasChanged = true;
        });


        window.addEventListener('focus', () => this.checkForUpdates());

        this.init();

        handleEscKey(() => {
            if (window.getComputedStyle(this.modal).display !== 'none') {
                this.closeModal();
            } else {
                this.handleBackButton();
            }
        });
    }

    async checkForUpdates() {
        // extra price
        if (localStorage.getItem('projectDetailNeedsUpdate') != null && localStorage.getItem('projectDetailNeedsUpdate') === 'true') {
            await this.getExtraPricesByProjectId();
            localStorage.removeItem('projectDetailNeedsUpdate');
        }
    }

    async init() {
        this.closeModal();
        await this.getDevicesFromDB();
        await this.getProjectDetails();
        this.renderAddedDevices(this.addedDevices);
        this.renderDevices(this.devices);
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.searchInput.addEventListener('input', (e) => this.handleSearch(e));
        document.getElementById('addDeviceToProjectBtn').addEventListener('click', () => this.openModal());
        document.getElementById('cancelBtn').addEventListener('click', () => this.closeModal());
        document.getElementById('saveProjectDBBtn').addEventListener('click', () => this.saveProject());
        document.getElementById('extraPricesBtn').addEventListener('click', () => this.navigateToExtraPrices());
        document.getElementById('exportBtn').addEventListener('click', () => this.export());
        document.getElementById('backBtn').addEventListener('click', () => this.handleBackButton());

    }



    // id="backBtn" handle back button
    handleBackButton() {
        // ask if user wants to save changes
        if (this.hasChanged) {
            if (confirm('Are you sure you want to discard changes?')) {
                localStorage.removeItem('projectDevices');
                this.devices = [];
                this.addedDevices = [];
                this.renderAddedDevices([]);
                window.history.back();
            }
        } else {
            localStorage.removeItem('projectDevices');
            this.devices = [];
            this.addedDevices = [];
            this.renderAddedDevices([]);
            window.history.back();
        }
    }

    navigateToExtraPrices() {
        window.location.href = `/extra-prices?id=${this.projectId}`;
    }

    export() {
        // ask a name for the file
        const fileName = prompt('Enter a name for the file:');
        if (fileName) {
            try {
                const formData = new FormData();
                formData.append('projectId', this.projectId);
                formData.append('fileName', fileName);
                // send request to server
                fetch(`${HTTP_URL}/project/export`, {
                    method: 'POST',
                    body: formData
                })
                    .then(response => {
                        if (response.ok) {
                            alert('File exported successfully');
                        } else {
                            throw new Error('Failed to export file');
                        }
                    })
                    .catch(error => {
                        console.error('Error:', error);
                        alert('An error occurred. Please try again later.');
                    })
            } catch (error) {
                console.error('Error:', error);
                alert('An error occurred. Please try again later.');
            }
        }
    }

    async saveProject() {
        const formData = new FormData();
        formData.append('projectId', this.projectId);
        formData.append('projectName', document.getElementById('projectName').value);

        // check if project name is empty or contains only spaces
        if (document.getElementById('projectName').value.trim() === '') {
            alert('Please enter a project name.');
            return;
        }


        let devicesJson = [];
        for (const device of this.addedDevices) {
            devicesJson.push({
                id: device.deviceId.toString(),
                count: device.count.toString(),
            });
        }

        formData.append('devices', JSON.stringify(devicesJson));

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

                if (projectDetails) {
                    if (projectDetails.project_device === null || projectDetails.project_device === undefined || projectDetails.project_device.length === 0) {
                        this.addedDevices = [];
                    } else {
                        for (const device of projectDetails.project_device) {
                            // Create new array instead of pushing to existing
                            this.addedDevices = projectDetails.project_device.map(device =>
                                new ProjectDevice(
                                    device.id,
                                    device.device.id,
                                    this.projectId,
                                    device.device.name,
                                    device.device.price,
                                    device.count
                                )
                            );
                        }
                    }

                    this.addedDevices.reverse();

                    // add addedDevices to localStorage
                    localStorage.setItem('projectDevices', JSON.stringify(this.addedDevices));

                    // Use value instead of textContent for input elements
                    document.getElementById('projectName').value = projectDetails.name || '';

                    // Format the initial price value
                    if (projectDetails.price) {
                        this.priceInput.value = projectDetails.price;
                        formatPriceInput(this.priceInput);
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


    async getExtraPricesByProjectId() {
        if (this.projectId) {
            try {
                const formData = new FormData();
                formData.append('projectId', this.projectId);

                const response = await fetch(`${HTTP_URL}/extraPrice/getAll`, {
                    method: 'POST',
                    body: formData
                });

                console.log(response);

                if (!response.ok) {
                    throw new Error('Failed to fetch project details');
                }

                const extraPrices = await response.json();
                console.log(extraPrices);
                let price = 0;
                for (const device of this.addedDevices) {
                    price += device.price * device.count;
                }
                if (extraPrices) {
                    if (extraPrices === null || extraPrices === undefined || extraPrices.length === 0) {
                        // pass
                    } else {
                        // loop through extraPrices
                        for (const extraPrice of extraPrices) {
                            // update price
                            price += extraPrice.price;
                        }
                    }
                    
                    // Format the initial price value
                    this.priceInput.value = price.toString();
                    formatPriceInput(this.priceInput);
                } else {
                    // Format the initial price value
                    this.priceInput.value = price.toString();
                    formatPriceInput(this.priceInput);
                }

            } catch (error) {
                console.error('Error fetching project details:', error);
                alert(error.message || 'An error occurred. Please try again later.');
            }
        }
    }


    renderDevices(devicesList) {

        this.devicesGrid.innerHTML = devicesList.map(device => {

            const converterStr = converterToString(device.converter);
            const filterStr = filterToString(device.filter);

            // check if device is already added 
            let added = false;
            for (let i = 0; i < this.addedDevices.length; i++) {
                if (this.addedDevices[i].deviceId === device.id) {
                    added = true;
                    break;
                }
            }

            if (added) {
                return `
                        <div class="card disabled" data-id="${device.id}">
                            <div class="card-header">
                                <span class="card-title">${this.escapeHtml(device.name)}</span>
                            </div>
                             <div class="card-header">
                                <div class="card-sub-title">نوع تبدیل: ${converterStr}</div>
                            </div>
                            <div class="card-header">
                                <div class="card-sub-title">صافی ${filterStr}</div>
                            </div>
                            <div class="card-header">
                                <div class="card-price">قیمت: ${formatPriceValue(device.price)}</div>
                            </div>
                            </div>
                    `;
            } else {
                return `
                        <div class="card" data-id="${device.id}">
                           <div class="card-header">
                                <span class="card-title">${this.escapeHtml(device.name)}</span>
                            </div>
                             <div class="card-header">
                                <div class="card-sub-title">نوع تبدیل: ${converterStr}</div>
                            </div>
                            <div class="card-header">
                                <div class="card-sub-title">صافی ${filterStr}</div>
                            </div>
                            <div class="card-header">
                                <div class="card-price">قیمت: ${formatPriceValue(device.price)}</div>
                            </div>
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
        countInput.className = 'count-input';
        countInput.value = 1;
        countInput.min = '1';
        countInput.id = `count-${deviceId}`;

        const addButton = document.createElement('button');
        addButton.className = 'card-button';
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

        // update price
        let device;
        for (let i = 0; i < this.devices.length; i++) {
            if (this.devices[i].id.toString() === deviceId.toString()) {
                device = this.devices[i];
                break;
            }
        }

        const projectDevice = new ProjectDevice(-1, deviceId, this.projectId, device.name, device.price, count);

        const projectDevices = JSON.parse(localStorage.getItem('projectDevices')) || [];
        projectDevices.push(projectDevice);
        localStorage.setItem('projectDevices', JSON.stringify(projectDevices));


        // this.priceInput.value is string and (device.price * count) is number
        // convert this.priceInput.value to number
        this.priceInput.value = formatPriceValue(convertPriceToNumber(this.priceInput.value) + (device.price * count));

        this.hasChanged = true;
        this.addedDevices.push(projectDevice);
        this.renderAddedDevices(this.addedDevices);
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

            const deviceCard = document.createElement('div');
            deviceCard.classList.add('card');
            deviceCard.setAttribute('data-device-id', device.deviceId);
            const converterStr = converterToString(device.converter);
            const filterStr = filterToString(device.filter);

            deviceCard.innerHTML = `

                            <div class="card-header">
                                <span class="card-title">${this.escapeHtml(device.name)}</span>
                            </div>
                             <div class="card-header">
                                <div class="card-sub-title">نوع تبدیل: ${converterStr}</div>
                            </div>
                            <div class="card-header">
                                <div class="card-sub-title">صافی ${filterStr}</div>
                            </div>
                            <div class="card-header">
                                <div class="card-price">قیمت: ${formatPriceValue(device.price)}</div>
                            </div>
                             <div class="card-header">
                                <div class="card-count-container">
                                    <div class="card-count-title">تعداد:</div>
                                    <button class="count-btn minus-btn" id="minus-${device.deviceId}">-</button>
                                    <div class="card-count">${device.count}</div>
                                    <button class="count-btn plus-btn" id="plus-${device.deviceId}">+</button>
                                </div>
                            </div>
                <div class="card-actions">
                    <button type="button" class="action-button delete-btn" data-id="delete-${device.deviceId}">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>

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
                const device = this.addedDevices.find(device => device.deviceId.toString() === deviceId.toString());
                if (!device) return;

                if (button.classList.contains('minus-btn')) {
                    device.count = Math.max(device.count - 1, 1);
                } else if (button.classList.contains('plus-btn')) {
                    // convert device.count to number
                    device.count = (parseInt(device.count) + 1).toString();
                }

                // get localstorage and update it
                const projectDevices = JSON.parse(localStorage.getItem('projectDevices')) || [];
                const existingDevice = projectDevices.find(device => device.deviceId.toString() === deviceId.toString());
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

        // update price
        let device;
        for (let i = 0; i < this.addedDevices.length; i++) {
            if (this.addedDevices[i].deviceId.toString() === deviceId.toString()) {
                device = this.addedDevices[i];
                // remove from this.addedDevices
                this.addedDevices.splice(i, 1);
                localStorage.setItem('projectDevices', JSON.stringify(this.addedDevices));
                break;
            }
        }
        this.priceInput.value = formatPriceValue(convertPriceToNumber(this.priceInput.value) - (device.price * device.count));

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

    handleSearch(e) {
        const searchTerm = e.target.value.toLowerCase();
        const filteredDevices = this.devices.filter(device =>
            device.name.toLowerCase().includes(searchTerm)
        );
        this.renderDevices(filteredDevices);
    }

    openModal() {
        this.modal.style.display = 'flex';
    }

    closeModal() {
        // remove selected class from all cards
        this.devicesGrid.querySelectorAll('.card').forEach(card => {
            card.classList.remove('selected');
        });
        this.modal.style.display = 'none';
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
    new AddProjectDetailsManager();
});