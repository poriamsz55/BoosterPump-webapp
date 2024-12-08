import { HTTP_URL } from '../config.js';
import { formatPriceInput } from '../format-price.js';

let hasChanged = false;
let backBtnDiv;

document.addEventListener('DOMContentLoaded', async function () {
    // Remove e.preventDefault() as it's not needed here
    const urlParams = new URLSearchParams(window.location.search);
    const partId = urlParams.get('id');
    // Add event listener for price input
    const priceInput = document.getElementById('partPrice');
    priceInput.addEventListener('input', function () {
        formatPriceInput(this);
    });

    // check if the form has changed
    const formElements = document.querySelectorAll('input, textarea');
    formElements.forEach(element => {
        element.addEventListener('input', () => {
            hasChanged = true;
            console.log('Form has changed');
        });
    });

    backBtnDiv = document.getElementById('backBtn');
    // handle back button
    backBtnDiv.addEventListener('click', () => {
        if (hasChanged) {
            if (confirm('Are you sure you want to leave without saving?')) {
                window.history.back();
            }
        } else {
            window.history.back();
        }
    });

    if (partId) {
        try {
            const formData = new FormData();
            formData.append('partId', partId);

            const response = await fetch(`${HTTP_URL}/part/getById`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                throw new Error('Failed to fetch part details');
            }

            const partDetails = await response.json();

            if (partDetails) {
                // Use value instead of textContent for input elements
                document.getElementById('partName').value = partDetails.name || '';
                document.getElementById('partSize').value = partDetails.size || '';
                document.getElementById('partMaterial').value = partDetails.material || '';
                document.getElementById('partBrand').value = partDetails.brand || '';

                // Format the initial price value
                if (partDetails.price) {
                    priceInput.value = partDetails.price;
                    formatPriceInput(priceInput);
                }
            } else {
                throw new Error('Part details not found');
            }

        } catch (error) {
            console.error('Error fetching part details:', error);
            alert(error.message || 'An error occurred. Please try again later.');
        }
    }
});

document.getElementById('partDetailsForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const urlParams = new URLSearchParams(window.location.search);
    const partId = urlParams.get('id');

    const formData = new FormData();
    formData.append('partId', partId);
    formData.append('partName', document.getElementById('partName').value);
    formData.append('partSize', document.getElementById('partSize').value);
    formData.append('partMaterial', document.getElementById('partMaterial').value);
    formData.append('partBrand', document.getElementById('partBrand').value);

    // Remove commas from price before sending to server
    const priceValue = document.getElementById('partPrice').value.replace(/,/g, '');
    formData.append('partPrice', priceValue);

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
        const response = await fetch(`${HTTP_URL}/part/update`, {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error('Failed to update part details');
        }

        // Set update flag in localStorage
        localStorage.setItem('partsListNeedsUpdate', 'true');

        window.history.back();

    } catch (error) {
        console.error('Error updating part details:', error);
        alert(error.message || 'An error occurred while saving. Please try again later.');
    }
});