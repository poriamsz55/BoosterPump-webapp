import { HTTP_URL } from '../config.js';
import { formatPriceInput } from '../format-price.js';

document.addEventListener('DOMContentLoaded', async function () {
    // Remove e.preventDefault() as it's not needed here
    const urlParams = new URLSearchParams(window.location.search);
    const partId = urlParams.get('id');
    // Add event listener for price input
    const priceInput = document.getElementById('partPrice');
    priceInput.addEventListener('input', function () {
        formatPriceInput(this);
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
            console.log(partDetails)

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