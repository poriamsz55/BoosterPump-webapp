import { HTTP_URL } from '../config.js';

class ProjectsManager {
    constructor() {
        this.projects = [];
        this.projectsGrid = document.getElementById('projectsGrid');
        this.searchInput = document.getElementById('searchProjects');

        this.init();
        
        // Add focus event listener
        window.addEventListener('focus', () => this.checkForUpdates());
    }

    async checkForUpdates() {
        if (localStorage.getItem('projectsListNeedsUpdate') === 'true') {
            await this.getProjectsFromDB();
            this.renderProjects(this.projects);
            localStorage.removeItem('projectsListNeedsUpdate');
        }
    }

    async init() {
        console.log('Initializing ProjectsManager...');
        await this.getProjectsFromDB();
        this.renderProjects(this.projects);
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.searchInput.addEventListener('input', (e) => this.handleSearch(e));
        document.getElementById('addProjectToDBButton').addEventListener('click', () =>{
            window.location.href = '/add/project/db';
        });
        // Remove the form submit event listener from here
        this.form.addEventListener('submit', async (e) => {
            e.preventDefault(); // Prevent form submission
            await this.saveProject();
        });
    }


    renderProjects(projectsList) {
        this.projectsGrid.innerHTML = projectsList.map(project => `
            <div class="card" data-id="${project.id}">
                <div class="card-header">
                    <span class="card-title">${this.escapeHtml(project.name)}</span>
                </div>
                <div class="card-price">${this.formatPrice(project.price)}</div>
                <div class="card-actions">
                    <button class="action-button delete-btn" data-id="${project.id}">
                        <i class="fas fa-trash"></i>
                    </button>
                    <button class="action-button copy-btn" data-id="${project.id}">
                        <i class="fas fa-copy"></i>
                    </button>
                </div>
            </div>
        `).join('');

        this.attachCardEventListeners();
    }

    attachCardEventListeners() {
        this.projectsGrid.querySelectorAll('.card').forEach(card => {
            card.addEventListener('click', (e) => {
                if (!e.target.closest('.action-button')) {
                    const projectId = card.getAttribute('data-id');
                    this.handleProjectClick(projectId);
                }
            });
        });

        this.projectsGrid.querySelectorAll('.delete-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.deleteProject(btn.dataset.id);
            });
        });

        this.projectsGrid.querySelectorAll('.copy-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.copyProject(btn.dataset.id);
            });
        });
    }

    handleProjectClick(id) {
        window.location.href = `/projects/details?id=${id}`;
    }

    handleSearch(e) {
        const searchTerm = e.target.value.toLowerCase();
        const filteredProjects = this.projects.filter(project =>
            project.name.toLowerCase().includes(searchTerm)
        );
        this.renderProjects(filteredProjects);
    }

    async getProjectsFromDB() {
        try {
            const response = await fetch(`${HTTP_URL}/project/getAll`);
            const data = await response.json();
            console.log(data);

            if (Array.isArray(data)) {
                this.projects = data;
            } else {
                console.error('Invalid response format:', data);
                this.projects = [];
            }
        } catch (error) {
            console.error('Error fetching projects:', error);
            this.projects = [];
        }
    }

    async deleteProject(id) {
        if (confirm('Are you sure you want to delete this project?')) {
            try {
                const response = await fetch(`${HTTP_URL}/project/delete/${id}`, {
                    method: 'POST'
                });

                if (response.ok) {
                    await this.getProjectsFromDB();
                    this.renderProjects(this.projects);
                } else {
                    alert('Failed to delete project.');
                }
            } catch (error) {
                console.error('Error deleting project:', error);
                alert('An error occurred while deleting the project.');
            }
        }
    }

    async copyProject(id) {
        const projectToCopy = this.projects.find(project => project.id === id);
        if (projectToCopy) {
            // Implementation for copying project
            console.log('Copying project:', projectToCopy);
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

    formatPrice(price) {
        return new Intl.NumberFormat('fa-IR').format(price);
    }
}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    new ProjectsManager();
});