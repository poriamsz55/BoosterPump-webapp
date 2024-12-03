import { HTTP_URL } from '../config.js';
import { formatPrice } from '../format-price.js';

document.addEventListener('DOMContentLoaded', async function () {
    // Remove e.preventDefault() as it's not needed here
    const urlParams = new URLSearchParams(window.location.search);
    const deviceId = urlParams.get('id');
    // Add event listener for price input
    const priceInput = document.getElementById('devicePrice');
    priceInput.addEventListener('input', function () {
        formatPrice(this);
    });

    if (deviceId) {
        try {
            const formData = new FormData();
            formData.append('deviceId', deviceId);

            const response = await fetch(`${HTTP_URL}/device/getById`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                throw new Error('Failed to fetch device details');
            }

            const deviceDetails = await response.json();
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
                    priceInput.value = deviceDetails.price;
                    formatPrice(priceInput);
                }
            } else {
                throw new Error('Device details not found');
            }

        } catch (error) {
            console.error('Error fetching device details:', error);
            alert(error.message || 'An error occurred. Please try again later.');
        }
    }
});

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