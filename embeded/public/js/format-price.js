export function formatPriceInput(input) {

    // Remove non-numeric characters except for "."
    let value = input.value.replace(/[^0-9.]/g, '');

    // Format with commas
    value = value.replace(/\B(?=(\d{3})+(?!\d))/g, ',');

    // Update the input value
    input.value = value;
}

export function formatPriceValue(price) {

    // Remove non-numeric characters except for "."
    let value = price.toString().replace(/[^0-9.]/g, '');

    // Format with commas
    value = value.replace(/\B(?=(\d{3})+(?!\d))/g, ',');

    // Update the input value
    return value;
}

export function convertPriceToNumber(value) {
    let val = value.toString().replace(/,/g, '');
    // // convert val to number
    // // how to ceil value
    // if (val.includes('.')) {
    //     val = val.split('.');
    //     val[0] = val[0] + '.' + val[1];
    //     val = val[0];
    // }

    val = parseInt(val);
    return val;
}