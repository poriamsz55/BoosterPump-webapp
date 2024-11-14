let darkMode = false;

function toggleTheme() {
    darkMode = !darkMode;
    document.documentElement.setAttribute('data-theme', darkMode ? 'dark' : 'light');
    document.querySelector('.theme-toggle i').className = darkMode ? 'fas fa-sun' : 'fas fa-moon';
}

// Sample parts data
const parts = [
    {
        name: "قطعه 1",
        size: "10x20",
        material: "فولاد",
        brand: "برند A",
        price: "1000000"
    },
    // Add more sample parts as needed
];

function createPartCard(part) {
    return `
        <div class="part-card">
            <div class="part-info">
                <span>نام: ${part.name}</span>
            </div>
            <div class="part-info">
                <span>سایز: ${part.size}</span>
            </div>
            <div class="part-info">
                <span>جنس: ${part.material}</span>
            </div>
            <div class="part-info">
                <span>برند: ${part.brand}</span>
            </div>
            <div class="part-info">
                <span>قیمت: ${part.price} تومان</span>
            </div>
            <div class="card-actions">
                <button class="action-button" onclick="deletePart(this)">
                    <i class="fas fa-trash"></i>
                </button>
            </div>
        </div>
    `;
}

function initializeParts() {
    const container = document.getElementById('partsContainer');
    parts.forEach(part => {
        container.innerHTML += createPartCard(part);
    });
}

function addNewPart() {
    const container = document.getElementById('partsContainer');
    const newPart = {
        name: "قطعه جدید",
        size: "0x0",
        material: "نامشخص",
        brand: "نامشخص",
        price: "0"
    };
    container.innerHTML += createPartCard(newPart);
}

function deletePart(button) {
    button.closest('.part-card').remove();
}

function saveDetails() {
    // Implement save functionality
    alert('ذخیره شد!');
}

// Search functionality
document.getElementById('searchInput').addEventListener('input', function(e) {
    const searchTerm = e.target.value.toLowerCase();
    const cards = document.querySelectorAll('.part-card');
    
    cards.forEach(card => {
        const text = card.textContent.toLowerCase();
        card.style.display = text.includes(searchTerm) ? 'block' : 'none';
    });
});

function showOverlay() {
    document.getElementById('overlay').style.display = 'flex';
    document.body.style.overflow = 'hidden';
}

function hideOverlay() {
    document.getElementById('overlay').style.display = 'none';
    document.body.style.overflow = 'auto';
}

function savePart() {
    // Add your save logic here
    console.log('Saving part...');
    hideOverlay();
}

// Close overlay when clicking outside the content
document.getElementById('overlay').addEventListener('click', function(e) {
    if (e.target === this) {
        hideOverlay();
    }
});

// Prevent closing when clicking inside the content
document.querySelector('.overlay-content').addEventListener('click', function(e) {
    e.stopPropagation();
});

// Initialize the page
document.addEventListener('DOMContentLoaded', initializeParts);
