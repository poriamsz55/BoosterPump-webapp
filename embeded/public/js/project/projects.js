// Sample projects data
const projects = [
    { id: 1, name: "Project 1", price: "$1000" },
    { id: 2, name: "Project 2", price: "$2000" },
    { id: 3, name: "Project 3", price: "$3000" },
];

// Initialize the projects page
function initProjectsPage() {
    const projectsGrid = document.getElementById('projectsGrid');
    const searchInput = document.getElementById('searchProjects');
    const projectDetails = document.getElementById('projectDetails');

    function renderProjects(projectsList) {
        projectsGrid.innerHTML = projectsList.map(project => `
            <div class="project-card" data-id="${project.id}">
                <div class="project-header">
                    <span class="project-title">${project.name}</span>
                </div>
                <div class="project-price">${project.price}</div>
                <div class="project-actions">
                    <button class="action-button" onclick="deleteProject(${project.id})">
                        <i class="fas fa-trash"></i>
                    </button>
                    <button class="action-button" onclick="copyProject(${project.id})">
                        <i class="fas fa-copy"></i>
                    </button>
                </div>
            </div>
        `).join('');

        // Add event listeners to each project card for click handling
        projectsGrid.querySelectorAll('.project-card').forEach(card => {
            card.addEventListener('click', (e) => {
                const projectId = card.getAttribute('data-id');
                handleProjectClick(projectId);
            });
        });
    }

    
    function handleProjectClick(id) {
        // Redirect to project-details.html with the project ID as a query parameter
        window.location.href = `/projects/details?id=${id}`;
    }

    // Initial render
    renderProjects(projects);

    // Search functionality
    searchInput.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        const filteredProjects = projects.filter(project => 
            project.name.toLowerCase().includes(searchTerm)
        );
        renderProjects(filteredProjects);
    });

    // Add project button
    document.getElementById('addProjectButton').addEventListener('click', () => {
        // Implement add project functionality
        console.log('Add project clicked');
    });
}

// Navigation functions
function showPage(pageId) {
    document.querySelectorAll('.page').forEach(page => {
        page.classList.remove('active');
    });
    document.getElementById(pageId).classList.add('active');
}

// Project actions
function deleteProject(id) {
    console.log('Delete project:', id);
    // Implement delete functionality
}

function copyProject(id) {
    console.log('Copy project:', id);
    // Implement copy functionality
}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    initProjectsPage();
    showPage('projectsPage'); // Show projects page by default
});
