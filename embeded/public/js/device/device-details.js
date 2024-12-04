import { HTTP_URL } from '../config.js';
import { convertPriceToNumber, formatPriceValue, formatPriceInput } from '../format-price.js';
import { DevicePart } from './device-part.js';

document.getElementById('deviceDetailsForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const urlParams = new URLSearchParams(window.location.search);
    const deviceId = urlParams.get('id');

    const formData = new FormData();
    formData.append('deviceId', deviceId);
    formData.append('deviceName', document.getElementById('deviceName').value);

    formData.append('deviceConverter', document.getElementById('deviceConverter').value);
    formData.append('deviceFilter', document.getElementById('deviceFilter').value);

    try {
        const response = await fetch(`${HTTP_URL}/device/update`, {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error('Failed to update device details');
        }

        // Set update flag in localStorage
        localStorage.setItem('devicesListNeedsUpdate', 'true');

        alert('اطلاعات با موفقیت ذخیره شد');
        window.history.back();

    } catch (error) {
        console.error('Error updating device details:', error);
        alert(error.message || 'An error occurred while saving. Please try again later.');
    }
});


// Device Details Manager
class AddDeviceDetailsManager {
    constructor() {

        this.parts = [];
        this.addedParts = [];
        this.partsGrid = document.getElementById('partsGrid');
        this.searchInput = document.getElementById('searchParts');
        this.modal = document.getElementById('addPartToDeviceModal');
        this.partForm = document.getElementById('addPartToDeviceForm');
        // Remove e.preventDefault() as it's not needed here
        const urlParams = new URLSearchParams(window.location.search);
        this.deviceId = urlParams.get('id');
        // Add event listener for price input
        this.priceInput = document.getElementById('devicePrice');
        this.priceInput.addEventListener('input', function () {
            formatPriceInput(this);
        });

        this.init();

        // Add this to handle back button
        window.addEventListener('popstate', () => {
            localStorage.removeItem('deviceParts');
            this.parts = [];
            this.addedParts = [];
            this.renderAddedParts([]);
        });
    }

    async init() {
        this.closeModal();
        await this.getPartsFromDB();
        await this.getDeviceDetails();
        this.renderAddedParts(this.addedParts);
        this.renderParts(this.parts);
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.searchInput.addEventListener('input', (e) => this.handleSearch(e));
        document.getElementById('addPartToDeviceBtn').addEventListener('click', () => this.openModal());
        document.getElementById('cancelBtn').addEventListener('click', () => this.closeModal());
        document.getElementById('saveDeviceDBBtn').addEventListener('click', () => this.saveDevice());
    }

    async saveDevice() {
        const formData = new FormData();
        formData.append('deviceName', document.getElementById('deviceName').value);
        formData.append('converterType', document.getElementById('converterType').value); // Fixed field name
        formData.append('filter', document.getElementById('filterCheckbox').checked); // Fixed field name
        formData.append('parts', JSON.stringify(deviceParts));

        try {
            const response = await fetch(`${HTTP_URL}/device/update`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) throw new Error('Failed to save device');

            localStorage.removeItem('deviceParts');

            // Set update flag in localStorage
            localStorage.setItem('devicesListNeedsUpdate', 'true');

            window.history.back();
        } catch (error) {
            console.error('Error:', error);
            alert('An error occurred. Please try again later.');
        }
    }

    async getDeviceDetails() {
        if (this.deviceId) {
            try {
                const formData = new FormData();
                formData.append('deviceId', this.deviceId);

                const response = await fetch(`${HTTP_URL}/device/getById`, {
                    method: 'POST',
                    body: formData
                });

                if (!response.ok) {
                    throw new Error('Failed to fetch device details');
                }

                const deviceDetails = await response.json();

                if (deviceDetails) {
                    for (const part of deviceDetails.device_part) {
                        console.log(part)
                        this.addedParts.push(
                            new DevicePart(
                                part.id,
                                part.part.id,
                                this.deviceId,
                                part.part.name,
                                part.part.price,
                                part.count
                            )
                        );
                    }

                    // add addedParts to localStorage
                    localStorage.setItem('deviceParts', JSON.stringify(this.addedParts));

                    // Use value instead of textContent for input elements
                    document.getElementById('deviceName').value = deviceDetails.name || '';

                    // converter
                    const converterSelect = document.getElementById('converterType');
                    const converterOption = converterSelect.querySelector(`option[value="${deviceDetails.converter}"]`);
                    if (converterOption) {
                        converterOption.selected = true;
                    }

                    // filter check box
                    const filterCheckbox = document.getElementById('filterCheckbox');
                    filterCheckbox.checked = deviceDetails.filter === 'true';

                    // Format the initial price value
                    if (deviceDetails.price) {
                        this.priceInput.value = deviceDetails.price;
                        formatPriceInput(this.priceInput);
                    }
                } else {
                    throw new Error('Device details not found');
                }

            } catch (error) {
                console.error('Error fetching device details:', error);
                alert(error.message || 'An error occurred. Please try again later.');
            }
        }
    }

    renderParts(partsList) {

        this.partsGrid.innerHTML = partsList.map(part => {

            // check if part is already added 
            let added = false;
            for (let i = 0; i < this.addedParts.length; i++) {
                if (this.addedParts[i].partId === part.id) {
                    added = true;
                    break;
                }
            }

            if (added) {
                return `
                        <div class="card disabled" data-id="${part.id}">
                            <div class="card-header">
                                <span class="card-title">${this.escapeHtml(part.name)}</span>
                            </div>
                            <div class="card-price">${formatPriceValue(part.price)}</div>
                        </div>
                    `;
            } else {
                return `
                        <div class="card" data-id="${part.id}">
                            <div class="card-header">
                                <span class="card-title">${this.escapeHtml(part.name)}</span>
                            </div>
                            <div class="card-price">${formatPriceValue(part.price)}</div>
                        </div>
                    `;
            }
        }

        ).join('');

        this.attachCardEventListeners();
    }

    attachCardEventListeners() {

        this.partsGrid.querySelectorAll('.card').forEach(card => {
            if (card.hasEventListener) return; // Prevent duplicate listeners
            card.hasEventListener = true;

            card.addEventListener('click', (e) => {
                if (card.classList.contains('disabled') ||
                    card.classList.contains('selected') ||
                    e.target.closest('.action-button')) {
                    return;
                }

                this.partsGrid.querySelectorAll('.card').forEach(c => {
                    if (c.classList.contains('selected')) {
                        c.classList.remove('selected');
                        const partId = c.getAttribute('data-id');
                        this.removeInputsFromCard(c, partId);
                    }
                });
                card.classList.add('selected');

                const partId = card.getAttribute('data-id');
                this.addInputsToCard(card, partId);
            });
        });
    }

    addInputsToCard(card, partId) {
        // Remove existing inputs if any
        card.querySelectorAll('input, button').forEach(el => el.remove());

        const countInput = document.createElement('input');
        countInput.type = 'number';
        countInput.value = 1;
        countInput.min = '1';
        countInput.id = `count-${partId}`;

        const addButton = document.createElement('button');
        addButton.type = 'button';
        addButton.textContent = 'افزودن به دستگاه';
        addButton.id = `add-to-device-${partId}`;

        addButton.addEventListener('click', () =>
            this.addToDevice(partId, countInput.value));

        card.appendChild(countInput);
        card.appendChild(addButton);
    }


    async addToDevice(partId, count) {

        // remove selected class from card
        const card = this.partsGrid.querySelector(`.card[data-id="${partId}"]`);
        card.classList.remove('selected');

        // add disabled class to card
        card.classList.add('disabled');

        // remove count input field
        const countInput = document.getElementById(`count-${partId}`);
        countInput.remove();

        // remove add to device button
        const addToDeviceBtn = document.getElementById(`add-to-device-${partId}`);
        addToDeviceBtn.remove();

        // update price
        let part;
        for (let i = 0; i < this.parts.length; i++) {
            if (this.parts[i].id.toString() === partId.toString()) {
                part = this.parts[i];
                break;
            }
        }

        const devicePart = new DevicePart(-1, partId, this.deviceId, part.name, part.price, count);

        const deviceParts = JSON.parse(localStorage.getItem('deviceParts')) || [];
        deviceParts.push(devicePart);
        localStorage.setItem('deviceParts', JSON.stringify(deviceParts));


        // this.priceInput.value is string and (part.price * count) is number
        // convert this.priceInput.value to number
        console.log("new value : ", formatPriceValue(part.price * count));
        this.priceInput.value = formatPriceValue(convertPriceToNumber(this.priceInput.value) + (part.price * count));


        this.addedParts.push(devicePart);
        this.renderAddedParts(this.addedParts);
    }


    removeInputsFromCard(card) {
        // Remove existing inputs if any
        card.querySelectorAll('input, button').forEach(el => el.remove());
    }

    // render added parts in modal
    renderAddedParts(addedParts) {
        const partsGrid = document.getElementById('addedPartsGrid');
        partsGrid.innerHTML = '';
        addedParts.forEach(part => {

            const partCard = document.createElement('div');
            partCard.classList.add('card');
            partCard.setAttribute('data-part-id', part.partId);
            partCard.innerHTML = `
                <div class="card-title">${this.escapeHtml(part.name)}</div>
                <div class="card-price">${formatPriceValue(part.price)}</div>
                <div class="card-count">${part.count}</div>
                <button type="button" class="action-button delete-btn" data-id="delete-${part.partId}">
                    <i class="fas fa-trash"></i>
                </button>
            `;
            partsGrid.appendChild(partCard);
        });

        this.attachAddedCardEventListeners();
    }

    // attachAddedCardEventListeners
    attachAddedCardEventListeners() {
        const partsGrid = document.getElementById('addedPartsGrid');
        partsGrid.querySelectorAll('.delete-btn').forEach(button => {
            button.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                // get id data-id="delete-${part.id}
                const partId = button.dataset.id.replace('delete-', '');
                this.deletePart(partId);
            });
        });
    }

    deletePart(partId) {

        // Remove just the specific card instead of re-rendering everything
        const cardToRemove = document.querySelector(`[data-part-id="${partId}"]`);
        if (cardToRemove) {
            cardToRemove.remove();
        }

        // Update the disabled state of the corresponding part in the parts grid
        const originalCard = this.partsGrid.querySelector(`[data-id="${partId}"]`);
        if (originalCard) {
            originalCard.classList.remove('disabled');
            originalCard.classList.remove('selected');
        }

        // update price
        let part;
        for (let i = 0; i < this.addedParts.length; i++) {
            if (this.addedParts[i].partId.toString() === partId.toString()) {
                part = this.addedParts[i];
                // remove from this.addedParts
                this.addedParts.splice(i, 1);
                localStorage.setItem('deviceParts', JSON.stringify(this.addedParts));
                break;
            }
        }
        console.log("delete part : ", part);
        console.log("delete value : ", formatPriceValue(part.price * part.count));
        this.priceInput.value = formatPriceValue(convertPriceToNumber(this.priceInput.value) - (part.price * part.count));

    }


    async getPartsFromDB() {
        try {
            const response = await fetch(`${HTTP_URL}/part/getAll`);
            const data = await response.json();

            if (Array.isArray(data)) {
                this.parts = data;
            } else {
                console.error('Invalid response format:', data);
                this.parts = [];
            }
        } catch (error) {
            console.error('Error fetching parts:', error);
            this.parts = [];
        }
    }

    handleSearch(e) {
        const searchTerm = e.target.value.toLowerCase();
        const filteredParts = this.parts.filter(part =>
            part.name.toLowerCase().includes(searchTerm)
        );
        this.renderParts(filteredParts);
    }

    openModal() {
        this.modal.style.display = 'flex';
        this.partForm.reset();
    }

    closeModal() {
        // remove selected class from all cards
        this.partsGrid.querySelectorAll('.card').forEach(card => {
            card.classList.remove('selected');
        });
        this.modal.style.display = 'none';
        this.partForm.reset();
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
    new AddDeviceDetailsManager();
});