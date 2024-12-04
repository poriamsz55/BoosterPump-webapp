import { HTTP_URL } from '../config.js';

class AddDeviceManager {
    constructor() {
        this.parts = [];
        this.addedParts = [];
        this.partsGrid = document.getElementById('partsGrid');
        this.searchInput = document.getElementById('searchParts');
        this.modal = document.getElementById('addPartToDeviceModal');
        this.form = document.getElementById('addDeviceDBForm')
        this.partForm = document.getElementById('addPartToDeviceForm');

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
        this.renderParts(this.parts);
        this.renderAddedParts(this.addedParts);
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.searchInput.addEventListener('input', (e) => this.handleSearch(e));
        document.getElementById('addPartToDeviceBtn').addEventListener('click', () => this.openModal());
        document.getElementById('cancelBtn').addEventListener('click', () => this.closeModal());
        document.getElementById('addDeviceDBBtn').addEventListener('click', () => this.saveDevice());
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
        const partsGrid = document.getElementById('addedPartsGrid');
        partsGrid.innerHTML = '';
        addedParts.forEach(part => {
            let founded = this.parts.find(p => p.id.toString() === part.id);
            if (!founded) return;

            const partCard = document.createElement('div');
            partCard.classList.add('card');
            partCard.setAttribute('data-part-id', part.id);
            partCard.innerHTML = `
                <div class="card-title">${founded.name}</div>
                <div class="card-price">${founded.price}</div>
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

        const devicePart = {
            id: partId,
            count: count
        };

        const deviceParts = JSON.parse(localStorage.getItem('deviceParts')) || [];
        deviceParts.push(devicePart);
        localStorage.setItem('deviceParts', JSON.stringify(deviceParts));

        this.addedParts.push(devicePart);
        this.renderAddedParts(this.addedParts);
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

    async saveDevice() {
        const formData = new FormData();
        formData.append('deviceName', document.getElementById('deviceName').value);
        formData.append('converterType', document.getElementById('converterType').value); // Fixed field name
        formData.append('filter', document.getElementById('filterCheckbox').checked); // Fixed field name

        try {
            const response = await fetch(`${HTTP_URL}/device/add`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) throw new Error('Failed to save device');

            const deviceId = await response.json();
            const deviceParts = JSON.parse(localStorage.getItem('deviceParts')) || [];
            const partFormData = new FormData();
            partFormData.append('deviceId', deviceId.id);
            partFormData.append('parts', JSON.stringify(deviceParts));

            const partResponses = await fetch(`${HTTP_URL}/devicePart/add/list`, {
                method: 'POST',
                body: partFormData,
            });

            if (!partResponses.ok) throw new Error('Failed to save device parts');

            localStorage.removeItem('deviceParts');

            // Set update flag in localStorage
            localStorage.setItem('devicesListNeedsUpdate', 'true');

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

    formatPrice(price) {
        return new Intl.NumberFormat('fa-IR').format(price);
    }
}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    new AddDeviceManager();
});