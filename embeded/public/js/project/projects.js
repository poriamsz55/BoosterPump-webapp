import { HTTP_URL } from '../config.js';
import { formatPriceValue } from '../format-price.js';
import { handleEscKey } from '../keyboard-utils.js';
import { gregorianToJalali, formatTime } from '../jalali.js';

class ProjectsManager {
    constructor() {
        this.projects = [];
        this.projectsGrid = document.getElementById('projectsGrid');
        this.searchInput = document.getElementById('searchProjects');

        this.init();

        // Add focus event listener
        window.addEventListener('focus', () => this.checkForUpdates());

        this.sortIcon = document.querySelector('.sort-icon');
        this.sortDropdown = document.querySelector('.sort-dropdown');
        this.sortOptions = document.querySelectorAll('.sort-option');
        this.setupSortEventListeners();

        handleEscKey(() => {
            // if sort dropdown is open, close it
            if (this.sortDropdown.classList.contains('show')) {
                this.sortDropdown.classList.remove('show');
            }else{
                window.history.back();
            }

        });
    }

    setupSortEventListeners() {
        // Toggle dropdown
        this.sortIcon.addEventListener('click', (e) => {
            e.stopPropagation();
            this.sortDropdown.classList.toggle('show');
        });

        // Close dropdown when clicking outside
        document.addEventListener('click', () => {
            this.sortDropdown.classList.remove('show');
        });

        // Prevent dropdown from closing when clicking inside it
        this.sortDropdown.addEventListener('click', (e) => {
            e.stopPropagation();
        });

        // Handle sort options
        this.sortOptions.forEach(option => {
            option.addEventListener('click', () => {
                const sortType = option.dataset.sort;
                const directionIcon = option.querySelector('.direction');
                const isDescending = directionIcon.classList.contains('fa-arrow-down');
                
                // Update active state and direction for all options
                this.sortOptions.forEach(opt => {
                    opt.classList.remove('active');
                    opt.querySelector('.direction').className = 'fas fa-arrow-up direction';
                });
                
                // Update clicked option
                option.classList.add('active');
                directionIcon.className = isDescending 
                    ? 'fas fa-arrow-up direction'
                    : 'fas fa-arrow-down direction';

                // Sort the projects
                this.sortProjects(sortType, !isDescending);
            });
        });
    }

    sortProjects(sortType, isDescending) {
        const sortedProjects = [...this.projects].sort((a, b) => {
            let valueA, valueB;
            
            if (sortType === 'name') {
                valueA = a.name.toLowerCase();
                valueB = b.name.toLowerCase();
            } else if (sortType === 'date') {
                valueA = new Date(a.modified_at);
                valueB = new Date(b.modified_at);
            }
            
            if (isDescending) {
                [valueA, valueB] = [valueB, valueA];
            }
            
            if (valueA < valueB) return -1;
            if (valueA > valueB) return 1;
            return 0;
        });

        this.renderProjects(sortedProjects);
    }


    async checkForUpdates() {
        if (localStorage.getItem('projectsListNeedsUpdate') === 'true') {
            await this.getProjectsFromDB();
            this.renderProjects(this.projects);
            localStorage.removeItem('projectsListNeedsUpdate');
        }
    }

    async init() {
        await this.getProjectsFromDB();
        this.renderProjects(this.projects);
        this.setupEventListeners();
    }

    setupEventListeners() {
        this.searchInput.addEventListener('input', (e) => this.handleSearch(e));
        document.getElementById('addProjectToDBButton').addEventListener('click', () => {
            window.location.href = '/add/project/db';
        });
    }

    renderProjects(projectsList) {
        this.projectsGrid.innerHTML = projectsList.map(project => `
            <div class="card" data-id="${project.id}" data-modified-at="${project.modified_at}">
                <div class="card-header">
                    <span class="card-title">${this.escapeHtml(project.name)}</span>
                </div>
                <div class="card-header">
                <div class="card-price">قیمت: ${formatPriceValue(project.price)}</div>
                </div>
                <div class="card-actions">
                    <button class="action-button delete-btn" data-id="delete-${project.id}">
                        <i class="fas fa-trash"></i>
                    </button>
                    <button class="action-button copy-btn" data-id="copy-${project.id}">
                        <i class="fas fa-copy"></i>
                    </button>
                </div>

                 <div class="card-datetime">
                    <div class="datetime-container">
                        <div>${gregorianToJalali(new Date(project.modified_at))}</div>
                        <div>${formatTime(project.modified_at)}</div>
                    </div>
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
                this.deleteProject(btn.dataset.id.toString().replace('delete-', ''));
            });
        });

        this.projectsGrid.querySelectorAll('.copy-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.stopPropagation();
                this.copyProject(btn.dataset.id.toString().replace('copy-', ''));
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

            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            const data = await response.json();

            if (Array.isArray(data)) {
                this.projects = data;
                this.projects.reverse();
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
                const formData = new FormData();
                formData.append('projectId', id); 
                const response = await fetch(`${HTTP_URL}/project/delete`, {
                    method: 'POST',
                    body: formData
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
        const projectToCopy = this.projects.find(project => project.id.toString() === id.toString());
        if (projectToCopy) {
            try {
                const formData = new FormData();
                formData.append('projectId', id);

                const response = await fetch(`${HTTP_URL}/project/copy`, {
                    method: 'POST',
                    body: formData
                });

                if (response.ok) {
                    await this.getProjectsFromDB();
                    this.renderProjects(this.projects);
                } else {
                    alert('Failed to copy project.');
                }
            } catch (error) {
                console.error('Error copying project:', error);
                alert('An error occurred while copying the project.');
            }
        } else {
            alert('Project not found.');
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

}

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    new ProjectsManager();
});