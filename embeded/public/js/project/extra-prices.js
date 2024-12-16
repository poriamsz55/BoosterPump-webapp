import { HTTP_URL } from '../config.js';
import { formatPriceValue } from '../format-price.js';
import { handleEscKey } from '../keyboard-utils.js';

class ExtraPricesManager {
    constructor() {
        this.extraPrices = [];
        this.extraPricesGrid = document.getElementById('extraPricesGrid');
        this.searchInput = document.getElementById('searchExtraPrices');
        this.modal = document.getElementById('addExtraPriceToDBModal');
        this.detailModal = document.getElementById('extraPriceDetailModal');
        this.form = document.getElementById('addExtraPriceToDBForm');
        this.detailForm = document.getElementById('extraPriceDetailForm');

        const urlParams = new URLSearchParams(window.location.search);
        this.projectId = urlParams.get('id');

        this.selectedExtraPrice = null;

        this.hasChanged = false;

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
        await this.getExtraPricesFromDB();
        this.renderExtraPrices(this.extraPrices);
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.searchInput.addEventListener('input', (e) => this.handleSearch(e));
        document.getElementById('addExtraPriceToDBButton').addEventListener('click', () => this.openModal());
        document.getElementById('cancelBtn').addEventListener('click', () => this.closeModal());
        document.getElementById('cancelDetailBtn').addEventListener('click', () => this.closeDetailModal());
        document.getElementById('saveDetailBtn').addEventListener('click', async () => await this.saveDetailExtraPrice());
        document.getElementById('backBtn').addEventListener('click', () => this.handleBackButton());
        // Remove the form submit event listener from here
        this.form.addEventListener('submit', async (e) => {
            e.preventDefault(); // Prevent form submission
            await this.saveExtraPrice();
        });

        // Add preventDefault to the detail form submission
        this.detailForm.addEventListener('submit', (e) => e.preventDefault());
    }


    handleBackButton() {
        if (this.hasChanged) {
            localStorage.setItem('projectDetailNeedsUpdate', 'true');
            localStorage.setItem('projectsListNeedsUpdate', 'true');
        }
        window.history.back();
    }

    renderExtraPrices(extraPricesList) {

        this.extraPricesGrid.innerHTML = '';

        this.extraPricesGrid.innerHTML = extraPricesList.map(extraPrice => `
            <div class="card" data-id="${extraPrice.id}">
                <div class="card-header">
                    <span class="card-title">${this.escapeHtml(extraPrice.name)}</span>
                </div>
                <div class="card-price">قیمت: ${formatPriceValue(extraPrice.price)}</div>
                <div class="card-actions">
                    <button class="action-button delete-btn" data-id="delete-${extraPrice.id}">
                        <i class="fas fa-trash"></i>
                    </button>
                    <button class="action-button copy-btn" data-id="copy-${extraPrice.id}">
                        <i class="fas fa-copy"></i>
                    </button>
                </div>
            </div>
        `).join('');

        this.attachCardEventListeners();
    }

    attachCardEventListeners() {
        this.extraPricesGrid.querySelectorAll('.card').forEach(card => {
            card.addEventListener('click', (e) => {
                if (!e.target.closest('.action-button')) {
                    const extraPriceId = card.getAttribute('data-id');
                    this.handleExtraPriceClick(extraPriceId);
                }
            });
        });

        this.extraPricesGrid.querySelectorAll('.delete-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.deleteExtraPrice(btn.dataset.id.toString().replace('delete-', ''));
            });
        });

        this.extraPricesGrid.querySelectorAll('.copy-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.copyExtraPrice(btn.dataset.id.toString().replace('copy-', ''));
            });
        });
    }

    handleExtraPriceClick(id) {
        this.selectedExtraPrice = this.extraPrices.find(extraPrice => extraPrice.id.toString() === id.toString());
        this.openDetailsModal();
    }

    handleSearch(e) {
        const searchTerm = e.target.value.toLowerCase();
        const filteredExtraPrices = this.extraPrices.filter(extraPrice =>
            extraPrice.name.toLowerCase().includes(searchTerm)
        );
        this.renderExtraPrices(filteredExtraPrices);
    }

    async handleFormSubmit(e) {
        e.preventDefault();
        await this.saveExtraPrice();
    }

    openModal() {
        this.modal.style.display = 'flex';
        this.form.reset();
    }

    closeModal() {
        this.modal.style.display = 'none';
        this.form.reset();
    }

    openDetailsModal() {
        this.detailModal.style.display = 'flex';
        this.detailForm.reset();

        if (this.selectedExtraPrice) {
            document.getElementById('detailExtraPriceName').value = this.selectedExtraPrice.name;
            document.getElementById('detailExtraPriceValue').value = formatPriceValue(this.selectedExtraPrice.price);
        }
    }

    // save detail btn for extraPrice
    async saveDetailExtraPrice() {
        // Get form data directly from the form elements
        const formData = new FormData();
        formData.append('extraPriceId', this.selectedExtraPrice.id);
        formData.append('extraPriceName', document.getElementById('detailExtraPriceName').value);
        formData.append('extraPriceValue', document.getElementById('detailExtraPriceValue').value.replace(/,/g, ''));
        
        if (document.getElementById('detailExtraPriceName').value.trim() === '') {
            alert('Please enter a extraPrice name.');
            return;
        }
        
        if (document.getElementById('detailExtraPriceValue').value.trim() === '') {
            // set value to 0 if empty
            formData.append('extraPriceValue', '0');    
        }

        this.hasChanged = true;

        try {
            const response = await fetch(`${HTTP_URL}/extraPrice/update`, {
                method: 'POST',
                body: formData
            });

            if (response.ok) {
                await this.getExtraPricesFromDB();
                this.renderExtraPrices(this.extraPrices);
                this.closeDetailModal();
            } else {
                alert('Failed to save extraPrice.');
            }
        } catch (error) {
            console.error('Error saving extraPrice:', error);
            alert('An error occurred. Please try again later.');
        }
    }


    closeDetailModal() {
        this.detailModal.style.display = 'none';
        this.detailForm.reset();
    }

    async saveExtraPrice() {
        // Get form data directly from the form elements
        const formData = new FormData();
        formData.append('projectId', this.projectId);
        formData.append('extraPriceName', document.getElementById('extraPriceName').value);
        formData.append('extraPriceValue', document.getElementById('extraPriceValue').value.replace(/,/g, ''));

        if (document.getElementById('extraPriceName').value.trim() === '') {
            alert('Please enter a extraPrice name.');
            return;
        }
        
        if (document.getElementById('extraPriceValue').value.trim() === '') {
            // set value to 0 if empty
            formData.append('extraPriceValue', '0');    
        }

        this.hasChanged = true;

        try {
            const response = await fetch(`${HTTP_URL}/extraPrice/add`, {
                method: 'POST',
                body: formData
            });

            if (response.ok) {
                await this.getExtraPricesFromDB();
                this.renderExtraPrices(this.extraPrices);
                this.closeModal();
            } else {
                alert('Failed to save extraPrice.');
            }
        } catch (error) {
            console.error('Error saving extraPrice:', error);
            alert('An error occurred. Please try again later.');
        }
    }

    async getExtraPricesFromDB() {
        try {
            const formData = new FormData();
            formData.append('projectId', this.projectId);
            const response = await fetch(`${HTTP_URL}/extraPrice/getAll`,
                {
                    method: 'POST',
                    body: formData,
                }
            );
            const data = await response.json();

            if (Array.isArray(data)) {
                this.extraPrices = data;
                // reverse the array
                this.extraPrices.reverse();
            } else {
                console.error('Invalid response format:', data);
                this.extraPrices = [];
            }
        } catch (error) {
            console.error('Error fetching extraPrices:', error);
            this.extraPrices = [];
        }
    }

    async deleteExtraPrice(id) {
        if (confirm('Are you sure you want to delete this extraPrice?')) {
            this.hasChanged = true;
            try {
                const formData = new FormData();
                formData.append('extraPriceId', id);
                const response = await fetch(`${HTTP_URL}/extraPrice/delete`, {
                    method: 'POST',
                    body: formData
                });

                if (response.ok) {
                    await this.getExtraPricesFromDB();
                    this.renderExtraPrices(this.extraPrices);
                } else {
                    alert('Failed to delete extraPrice.');
                }
            } catch (error) {
                console.error('Error deleting extraPrice:', error);
                alert('An error occurred while deleting the extraPrice.');
            }
        }
    }

    async copyExtraPrice(id) {
        const extraPriceToCopy = this.extraPrices.find(extraPrice => extraPrice.id.toString() === id.toString());
        if (extraPriceToCopy) {
            this.hasChanged = true;
            try {
                const formData = new FormData();
                formData.append('extraPriceId', id);

                const response = await fetch(`${HTTP_URL}/extraPrice/copy`, {
                    method: 'POST',
                    body: formData
                });

                if (response.ok) {
                    await this.getExtraPricesFromDB();
                    this.renderExtraPrices(this.extraPrices);
                } else {
                    alert('Failed to copy extraPrice.');
                }
            } catch (error) {
                console.error('Error copying extraPrice:', error);
                alert('An error occurred while copying the extraPrice.');
            }
        } else {
            alert('ExtraPrice not found.');
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
    new ExtraPricesManager();
});