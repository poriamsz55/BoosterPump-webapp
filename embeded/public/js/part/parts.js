import { HTTP_URL } from '../config.js';
import { formatPriceValue } from '../format-price.js';
import { handleEscKey } from '../keyboard-utils.js';
import { gregorianToJalali, formatTime } from '../jalali.js';

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

        this.sortIcon = document.querySelector('.sort-icon');
        this.sortDropdown = document.querySelector('.sort-dropdown');
        this.sortOptions = document.querySelectorAll('.sort-option');
        this.setupSortEventListeners();

        handleEscKey(() => {
            if (window.getComputedStyle(this.modal).display !== 'none') {
                this.closeModal();
            }else if (this.sortDropdown.classList.contains('show')) {
                this.sortDropdown.classList.remove('show');
            } else {
                window.history.back();
            }
        });
    }

    setupSortEventListeners() {
        // Toggle dropdown
        this.sortIcon.addEventListener('click', (e) => {
            e.stopPropagation();
            this.sortDropdown.classList.toggle('show');
        });

        // Close dropdown when clicking outside
        document.addEventListener('click', () => {
            this.sortDropdown.classList.remove('show');
        });

        // Prevent dropdown from closing when clicking inside it
        this.sortDropdown.addEventListener('click', (e) => {
            e.stopPropagation();
        });

        // Handle sort options
        this.sortOptions.forEach(option => {
            option.addEventListener('click', () => {
                const sortType = option.dataset.sort;
                const directionIcon = option.querySelector('.direction');
                const isDescending = directionIcon.classList.contains('fa-arrow-down');
                
                // Update active state and direction for all options
                this.sortOptions.forEach(opt => {
                    opt.classList.remove('active');
                    opt.querySelector('.direction').className = 'fas fa-arrow-up direction';
                });
                
                // Update clicked option
                option.classList.add('active');
                directionIcon.className = isDescending 
                    ? 'fas fa-arrow-up direction'
                    : 'fas fa-arrow-down direction';

                // Sort the parts
                this.sortParts(sortType, !isDescending);
            });
        });
    }

    sortParts(sortType, isDescending) {
        const sortedParts = [...this.parts].sort((a, b) => {
            let valueA, valueB;
            
            if (sortType === 'name') {
                valueA = a.name.toLowerCase();
                valueB = b.name.toLowerCase();
            } else if (sortType === 'date') {
                valueA = new Date(a.modified_at);
                valueB = new Date(b.modified_at);
            }
            
            if (isDescending) {
                [valueA, valueB] = [valueB, valueA];
            }
            
            if (valueA < valueB) return -1;
            if (valueA > valueB) return 1;
            return 0;
        });

        this.renderParts(sortedParts);
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
            <div class="card" data-id="${part.id}" data-modified-at="${part.modified_at}">
                <div class="card-header">
                    <span class="card-title">${this.escapeHtml(part.name)}</span>
                </div>
                  <div class="card-header">
                    <span class="card-sub-title">سایز: ${this.formatSize(part.size)}</span>
                </div>
                  <div class="card-header">
                    <span class="card-sub-title">جنس: ${this.escapeHtml(part.material)}</span>
                </div>
                   <div class="card-header">
                    <span class="card-sub-title">برند: ${this.escapeHtml(part.brand)}</span>
                </div>
                   <div class="card-header">
                <div class="card-price">قیمت: ${formatPriceValue(part.price)}</div>
                </div>
                <div class="card-actions">
                    <button class="action-button delete-btn" data-id="delete-${part.id}">
                        <i class="fas fa-trash"></i>
                    </button>
                    <button class="action-button copy-btn" data-id="copy-${part.id}">
                        <i class="fas fa-copy"></i>
                    </button>
                </div>

                <div class="card-datetime">
                    <div class="datetime-container">
                        <div>${gregorianToJalali(new Date(part.modified_at))}</div>
                        <div>${formatTime(part.modified_at)}</div>
                    </div>
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
                this.deletePart(btn.dataset.id.toString().replace('delete-', ''));
            });
        });

        this.partsGrid.querySelectorAll('.copy-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.copyPart(btn.dataset.id.toString().replace('copy-', ''));
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
        formData.append('partSize', document.getElementById('partSize').value) || '';
        formData.append('partMaterial', document.getElementById('partMaterial').value) || '';
        formData.append('partBrand', document.getElementById('partBrand').value) || '';
        formData.append('partPrice', document.getElementById('partPrice').value.replace(/,/g, ''));

        // Check if any of the required fields are empty
        if (document.getElementById('partName').value.trim() === '') {
            alert('Please enter a part name.');
            return;
        }

        if (document.getElementById('partPrice').value.trim() === '') {
            // update partPrice in formData with value 0 if empty
            formData.append('partPrice', '0');
        }

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
                this.parts.reverse();
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
                const formData = new FormData();
                formData.append('partId', id);
                const response = await fetch(`${HTTP_URL}/part/delete`, {
                    method: 'POST',
                    body: formData
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
        const partToCopy = this.parts.find(part => part.id.toString() === id.toString());
        if (partToCopy) {
            try {
                const formData = new FormData();
                formData.append('partId', id);

                const response = await fetch(`${HTTP_URL}/part/copy`, {
                    method: 'POST',
                    body: formData
                });

                if (response.ok) {
                    await this.getPartsFromDB();
                    this.renderParts(this.parts);
                } else {
                    alert('Failed to copy part.');
                }
            } catch (error) {
                console.error('Error copying part:', error);
                alert('An error occurred while copying the part.');
            }
        } else {
            alert('Part not found.');
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

    // formatSize(size) {
    //     return `<span style="direction: ltr">${size}</span>`;
    // }

    formatSize(size) {
        return `\u202A${size}\u202C`;
    }

}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    new PartsManager();
});