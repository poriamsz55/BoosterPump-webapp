import { HTTP_URL } from '../config.js';
import {formatPriceValue} from '../format-price.js';

class PartsManager {
    constructor() {
        this.parts = [];
        this.partsGrid = document.getElementById('partsGrid');
        this.searchInput = document.getElementById('searchParts');
        this.modal = document.getElementById('addPartToDBModal');
        this.form = document.getElementById('addPartToDBForm');

        this.init();
        
        // Add focus event listener
        window.addEventListener('focus', () => this.checkForUpdates());
    }

    async checkForUpdates() {
        if (localStorage.getItem('partsListNeedsUpdate') === 'true') {
            await this.getPartsFromDB();
            this.renderParts(this.parts);
            localStorage.removeItem('partsListNeedsUpdate');
        }
    }

    async init() {
        this.closeModal();
        await this.getPartsFromDB();
        this.renderParts(this.parts);
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.searchInput.addEventListener('input', (e) => this.handleSearch(e));
        document.getElementById('addPartToDBButton').addEventListener('click', () => this.openModal());
        document.getElementById('cancelBtn').addEventListener('click', () => this.closeModal());
        // Remove the form submit event listener from here
        this.form.addEventListener('submit', async (e) => {
            e.preventDefault(); // Prevent form submission
            await this.savePart();
        });
    }


    renderParts(partsList) {
        this.partsGrid.innerHTML = partsList.map(part => `
            <div class="card" data-id="${part.id}">
                <div class="card-header">
                    <span class="card-title">${this.escapeHtml(part.name)}</span>
                </div>
                <div class="card-price">${formatPriceValue(part.price)}</div>
                <div class="card-actions">
                    <button class="action-button delete-btn" data-id="${part.id}">
                        <i class="fas fa-trash"></i>
                    </button>
                    <button class="action-button copy-btn" data-id="${part.id}">
                        <i class="fas fa-copy"></i>
                    </button>
                </div>
            </div>
        `).join('');

        this.attachCardEventListeners();
    }

    attachCardEventListeners() {
        this.partsGrid.querySelectorAll('.card').forEach(card => {
            card.addEventListener('click', (e) => {
                if (!e.target.closest('.action-button')) {
                    const partId = card.getAttribute('data-id');
                    this.handlePartClick(partId);
                }
            });
        });

        this.partsGrid.querySelectorAll('.delete-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.deletePart(btn.dataset.id);
            });
        });

        this.partsGrid.querySelectorAll('.copy-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.copyPart(btn.dataset.id);
            });
        });
    }

    handlePartClick(id) {
        window.location.href = `/parts/details?id=${id}`;
    }

    handleSearch(e) {
        const searchTerm = e.target.value.toLowerCase();
        const filteredParts = this.parts.filter(part =>
            part.name.toLowerCase().includes(searchTerm)
        );
        this.renderParts(filteredParts);
    }

    async handleFormSubmit(e) {
        e.preventDefault();
        await this.savePart();
    }

    openModal() {
        this.modal.style.display = 'flex';
        this.form.reset();
    }

    closeModal() {
        this.modal.style.display = 'none';
        this.form.reset();
    }

    async savePart() {
        // Get form data directly from the form elements
        const formData = new FormData();
        formData.append('partName', document.getElementById('partName').value);
        formData.append('partSize', document.getElementById('partSize').value);
        formData.append('partMaterial', document.getElementById('partMaterial').value);
        formData.append('partBrand', document.getElementById('partBrand').value);
        formData.append('partPrice', document.getElementById('partPrice').value.replace(/,/g, ''));

        try {
            const response = await fetch(`${HTTP_URL}/part/add`, {
                method: 'POST',
                body: formData
            });

            if (response.ok) {
                await this.getPartsFromDB();
                this.renderParts(this.parts);
                this.closeModal();
            } else {
                alert('Failed to save part.');
            }
        } catch (error) {
            console.error('Error saving part:', error);
            alert('An error occurred. Please try again later.');
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

    async deletePart(id) {
        if (confirm('Are you sure you want to delete this part?')) {
            try {
                const response = await fetch(`${HTTP_URL}/part/delete/${id}`, {
                    method: 'DELETE'
                });

                if (response.ok) {
                    await this.getPartsFromDB();
                    this.renderParts(this.parts);
                } else {
                    alert('Failed to delete part.');
                }
            } catch (error) {
                console.error('Error deleting part:', error);
                alert('An error occurred while deleting the part.');
            }
        }
    }

    async copyPart(id) {
        const partToCopy = this.parts.find(part => part.id === id);
        if (partToCopy) {
            // Implementation for copying part
            console.log('Copying part:', partToCopy);
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
    new PartsManager();
});