// Sample parts data
const parts = [
    { id: 1, name: "Part 1", price: "$1000" },
    { id: 2, name: "Part 2", price: "$2000" },
    { id: 3, name: "Part 3", price: "$3000" },
];

// Initialize the parts page
function initPartsPage() {
    const partsGrid = document.getElementById('partsGrid');
    const searchInput = document.getElementById('searchParts');
    const partDetails = document.getElementById('partDetails');

    function renderParts(partsList) {
        partsGrid.innerHTML = partsList.map(part => `
            <div class="part-card" data-id="${part.id}">
                <div class="part-header">
                    <span class="part-title">${part.name}</span>
                </div>
                <div class="part-price">${part.price}</div>
                <div class="part-actions">
                    <button class="action-button" onclick="deletePart(${part.id})">
                        <i class="fas fa-trash"></i>
                    </button>
                    <button class="action-button" onclick="copyPart(${part.id})">
                        <i class="fas fa-copy"></i>
                    </button>
                </div>
            </div>
        `).join('');

        // Add event listeners to each part card for click handling
        partsGrid.querySelectorAll('.part-card').forEach(card => {
            card.addEventListener('click', (e) => {
                const partId = card.getAttribute('data-id');
                handlePartClick(partId);
            });
        });
    }

    
    function handlePartClick(id) {
        // Redirect to part-details.html with the part ID as a query parameter
        window.location.href = `/parts/details?id=${id}`;
    }

    // Initial render
    renderParts(parts);

    // Search functionality
    searchInput.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        const filteredParts = parts.filter(part => 
            part.name.toLowerCase().includes(searchTerm)
        );
        renderParts(filteredParts);
    });

    // Add part button
    document.getElementById('addPartButton').addEventListener('click', () => {
        // Implement add part functionality
        console.log('Add part clicked');
    });
}

// Navigation functions
function showPage(pageId) {
    document.querySelectorAll('.page').forEach(page => {
        page.classList.remove('active');
    });
    document.getElementById(pageId).classList.add('active');
}

// Part actions
function deletePart(id) {
    console.log('Delete part:', id);
    // Implement delete functionality
}

function copyPart(id) {
    console.log('Copy part:', id);
    // Implement copy functionality
}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    initPartsPage();
    showPage('partsPage'); // Show parts page by default
});
