document.getElementById('partDetailsForm').addEventListener('submit', function(e) {
    e.preventDefault();
    // Add your save functionality here
    console.log('Form submitted');
});

// Add smooth hover effect to inputs
const inputs = document.querySelectorAll('input');
inputs.forEach(input => {
    input.addEventListener('focus', function() {
        this.style.transform = 'scale(1.02)';
    });
    input.addEventListener('blur', function() {
        this.style.transform = 'scale(1)';
    });
});