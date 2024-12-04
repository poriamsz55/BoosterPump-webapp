import { HTTP_URL } from '../config.js';
import { formatPrice } from '../format-price.js';

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
            formatPrice(this);
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
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.searchInput.addEventListener('input', (e) => this.handleSearch(e));
        document.getElementById('addPartToDeviceBtn').addEventListener('click', () => this.openModal());
        document.getElementById('cancelBtn').addEventListener('click', () => this.closeModal());
        document.getElementById('saveDeviceDBBtn').addEventListener('click', () => this.saveDevice());
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
                this.addedParts = deviceDetails.device_part;
                console.log(deviceDetails)

                if (deviceDetails) {
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
                        formatPrice(this.priceInput);
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

        this.partsGrid.innerHTML = partsList.map(part =>

            `
                    <div class="card" data-id="${part.id}">
                        <div class="card-header">
                            <span class="card-title">${this.escapeHtml(part.name)}</span>
                        </div>
                        <div class="card-price">${this.formatPrice(part.price)}</div>
                    </div>
                `

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

    removeInputsFromCard(card) {
        // Remove existing inputs if any
        card.querySelectorAll('input, button').forEach(el => el.remove());
    }

    // render added parts in modal
    renderAddedParts(addedParts) {
        console.log(addedParts)
        const partsGrid = document.getElementById('addedPartsGrid');
        partsGrid.innerHTML = '';
        addedParts.forEach(part => {

            const partCard = document.createElement('div');
            partCard.classList.add('card');
            partCard.setAttribute('data-part-id', part.id);
            partCard.innerHTML = `
                <div class="card-title">${this.escapeHtml(part.part.name)}</div>
                <div class="card-price">${this.formatPrice(part.price)}</div>
                <div class="card-count">${part.count}</div>
                <button type="button" class="action-button delete-btn" data-id="delete-${part.id}">
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
                console.log(partId);
                this.deletePart(partId);
            });
        });
    }

    deletePart(partId) {
        this.addedParts = this.addedParts.filter(part => part.id !== partId);
        localStorage.setItem('deviceParts', JSON.stringify(this.addedParts));

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

    formatPrice(price) {
        return new Intl.NumberFormat('fa-IR').format(price);
    }
}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    new AddDeviceDetailsManager();
});