import { HTTP_URL } from '../config.js';
import { formatPriceInput } from '../format-price.js';
import { handleEscKey } from '../keyboard-utils.js';

class PartDetailsManager {
    constructor() {
        this.hasChanged = false;
        this.partId = new URLSearchParams(window.location.search).get('id');
        this.priceInput = document.getElementById('partPrice');
        this.backBtnDiv = document.getElementById('backBtn');
        this.form = document.getElementById('partDetailsForm');

        this.init();
        this.setupEventListeners();

        handleEscKey(() => {
            this.handleBackButton();
        });
    }

    init() {
        if (this.partId) {
            this.fetchPartDetails();
        }
    }

    setupEventListeners() {
        // Price input formatting
        this.priceInput.addEventListener('input', function () {
            formatPriceInput(this);
        });

        // Form change detection
        const formElements = document.querySelectorAll('input, textarea');
        formElements.forEach(element => {
            element.addEventListener('input', () => {
                this.hasChanged = true;
                console.log('Form has changed');
            });
        });

        // Back button
        this.backBtnDiv.addEventListener('click', () => this.handleBackButton());

        // Form submission
        this.form.addEventListener('submit', (e) => this.handleSubmit(e));
    }

    handleBackButton() {
        if (this.hasChanged) {
            if (confirm('Are you sure you want to leave without saving?')) {
                window.history.back();
            }
        } else {
            window.history.back();
        }
    }

    async fetchPartDetails() {
        try {
            const formData = new FormData();
            formData.append('partId', this.partId);

            const response = await fetch(`${HTTP_URL}/part/getById`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                throw new Error('Failed to fetch part details');
            }

            const partDetails = await response.json();

            if (partDetails) {
                this.populateForm(partDetails);
            } else {
                throw new Error('Part details not found');
            }

        } catch (error) {
            console.error('Error fetching part details:', error);
            alert(error.message || 'An error occurred. Please try again later.');
        }
    }

    populateForm(partDetails) {
        document.getElementById('partName').value = partDetails.name || '';
        document.getElementById('partSize').value = this.formatSize(partDetails.size) || '';
        document.getElementById('partMaterial').value = partDetails.material || '';
        document.getElementById('partBrand').value = partDetails.brand || '';

        if (partDetails.price) {
            this.priceInput.value = partDetails.price;
            formatPriceInput(this.priceInput);
        }
    }

    async handleSubmit(e) {
        e.preventDefault();

        const formData = new FormData();
        formData.append('partId', this.partId);
        formData.append('partName', document.getElementById('partName').value);
        formData.append('partSize', document.getElementById('partSize').value);
        formData.append('partMaterial', document.getElementById('partMaterial').value);
        formData.append('partBrand', document.getElementById('partBrand').value);

        const priceValue = this.priceInput.value.replace(/,/g, '');
        formData.append('partPrice', priceValue);

        if (document.getElementById('partName').value.trim() === '') {
            alert('Please enter a part name.');
            return;
        }

        if (this.priceInput.value.trim() === '') {
            formData.set('partPrice', '0');
        }

        try {
            const response = await fetch(`${HTTP_URL}/part/update`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                throw new Error('Failed to update part details');
            }

            localStorage.setItem('partsListNeedsUpdate', 'true');
            window.history.back();

        } catch (error) {
            console.error('Error updating part details:', error);
            alert(error.message || 'An error occurred while saving. Please try again later.');
        }
    }

    formatSize(size) {
        return `\u202A${size}\u202C`;
    }
}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    new PartDetailsManager();
});