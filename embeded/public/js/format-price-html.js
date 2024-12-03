function formatPrice(input) {
    // Remove non-numeric characters except for "."
    let value = input.value.replace(/[^0-9.]/g, '');

    // Format with commas
    value = value.replace(/\B(?=(\d{3})+(?!\d))/g, ',');

    // Update the input value
    input.value = value;
}