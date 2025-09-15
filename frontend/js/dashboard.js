/**
 * Vanilla JS Dashboard Manager
 * Uses modular components loaded via script tags for browser compatibility
 * No ES6 imports - works directly in browsers without build tools
 */

class VanillaDashboardManager {
  constructor() {
    this.currentUser = null;
    this.currentView = 'overview';
    this.components = {};
    this.managers = {};
    
    this.init();
  }

  async init() {
    try {
      // Wait for all component scripts to load
      await this.waitForDependencies();
      
      // Initialize content area
      this.contentArea = document.getElementById('content-body');
      
      // Initialize managers
      this.initializeManagers();
      
      // Check authentication
      await this.checkAuth();
      
      // Initialize navigation
      this.initNavigation();
      
      // Initialize mobile functionality
      this.initMobile();
      
      console.log('Vanilla Dashboard Manager initialized successfully');
    } catch (error) {
      console.error('Failed to initialize dashboard:', error);
      this.showError('Failed to initialize dashboard');
    }
  }

  // Wait for all required classes to be available
  async waitForDependencies() {
    const requiredClasses = [
      'ApiClient', 'notifications', 'appState', 'getFormData',
      'Modal', 'DataTable', 'Card', 'Grid',
      'EquipmentManager', 'ExerciseManager', 'GymManager'
    ];
    
    let attempts = 0;
    const maxAttempts = 100; // 10 seconds max wait
    
    while (attempts < maxAttempts) {
      const allLoaded = requiredClasses.every(className => {
        return window[className] !== undefined;
      });
      
      if (allLoaded) {
        console.log('All dependencies loaded successfully');
        return;
      }
      
      await new Promise(resolve => setTimeout(resolve, 100));
      attempts++;
    }
    
    console.warn('Some dependencies may not have loaded:', 
      requiredClasses.filter(cls => window[cls] === undefined));
  }

  initializeManagers() {
    // Initialize our specialized managers (they're now global)
    this.managers.equipment = new EquipmentManager();
    this.managers.exercise = new ExerciseManager();
    this.managers.gym = new GymManager();
    
    // Set up event listeners for manager events
    this.setupManagerEventListeners();
  }

  setupManagerEventListeners() {
    // Equipment events
    document.addEventListener('equipment:edit', (e) => {
      this.openEquipmentModal('edit', e.detail.equipment);
    });
    
    // Exercise events
    document.addEventListener('exercise:edit', (e) => {
      this.openExerciseModal('edit', e.detail.exercise);
    });
    
    document.addEventListener('exercise:view', (e) => {
      this.viewExerciseDetails(e.detail.exercise);
    });
    
    // Gym events
    document.addEventListener('gym:edit', (e) => {
      this.openGymModal('edit', e.detail.gym);
    });
    
    document.addEventListener('gym:view', (e) => {
      this.viewGymDetails(e.detail.gym);
    });
  }

  async checkAuth() {
    try {
      const api = new ApiClient();
      const userData = await api.getCurrentUser();
      
      this.currentUser = userData;
      appState.setState({ user: userData });
      
      this.updateUserInfo();
      this.setupNavigation();
    } catch (error) {
      console.error('Authentication check failed:', error);
      // Redirect to login or show auth error
      window.location.href = '/';
    }
  }

  updateUserInfo() {
    const userNameEl = document.getElementById('user-name');
    const userRoleEl = document.getElementById('user-role');
    
    if (userNameEl) {
      userNameEl.textContent = this.currentUser.username || this.currentUser.id;
    }
    
    if (userRoleEl) {
      userRoleEl.textContent = this.currentUser.role || this.currentUser.user_type || 'User';
    }
  }

  setupNavigation() {
    const platformAdminNav = document.querySelectorAll('.platform-admin-nav');
    
    if (this.currentUser.user_type === 'platform_admin') {
      platformAdminNav.forEach(nav => nav.style.display = 'block');
    } else {
      platformAdminNav.forEach(nav => nav.style.display = 'none');
    }
  }

  initNavigation() {
    // Handle navigation clicks using event delegation
    const navMenu = document.getElementById('navigation-menu');
    if (navMenu) {
      navMenu.addEventListener('click', (e) => {
        const navItem = e.target.closest('.nav-item');
        if (navItem) {
          e.preventDefault();
          
          const view = navItem.getAttribute('data-view');
          if (view) {
            this.loadView(view);
            this.updateActiveNavItem(navItem);
          }
        }
      });
    }
  }

  updateActiveNavItem(activeItem) {
    // Remove active class from all nav items
    document.querySelectorAll('.nav-item').forEach(item => {
      item.classList.remove('active');
    });
    
    // Add active class to clicked item
    activeItem.classList.add('active');
  }

  initMobile() {
    const sidebar = document.querySelector('.sidebar');
    const mobileToggle = document.getElementById('mobile-sidebar-toggle');
    const sidebarToggle = document.getElementById('sidebar-toggle');
    
    // Mobile sidebar toggle
    if (mobileToggle) {
      mobileToggle.addEventListener('click', () => {
        sidebar.classList.toggle('sidebar-open');
      });
    }
    
    if (sidebarToggle) {
      sidebarToggle.addEventListener('click', () => {
        sidebar.classList.toggle('sidebar-open');
      });
    }
    
    // Close sidebar when clicking outside on mobile
    document.addEventListener('click', (e) => {
      if (window.innerWidth <= 768 && 
          !sidebar.contains(e.target) && 
          !mobileToggle.contains(e.target)) {
        sidebar.classList.remove('sidebar-open');
      }
    });
  }

  async loadView(viewName) {
    this.currentView = viewName;
    appState.setState({ currentView: viewName });
    
    // Update page header
    this.updatePageHeader(viewName);
    
    // Show loading state
    this.showLoading();
    
    try {
      switch (viewName) {
        case 'platform-overview':
          await this.loadPlatformOverview();
          break;
        case 'gyms':
          await this.loadGymsManagement();
          break;
        case 'equipment':
          await this.loadEquipmentManagement();
          break;
        case 'exercises':
          await this.loadExercisesManagement();
          break;
        case 'muscular-groups':
          await this.loadMuscularGroupsManagement();
          break;
        default:
          this.showComingSoon(viewName);
      }
    } catch (error) {
      console.error(`Error loading view ${viewName}:`, error);
      this.showError(`Failed to load ${viewName}`);
    } finally {
      this.hideLoading();
    }
  }

  updatePageHeader(viewName) {
    const titles = {
      'platform-overview': {
        title: 'Platform Overview',
        subtitle: 'Monitor platform-wide statistics and activity.'
      },
      'gyms': {
        title: 'Gyms Management',
        subtitle: 'Manage all gyms in the platform.'
      },
      'equipment': {
        title: 'Equipment Management',
        subtitle: 'Manage global equipment catalog.'
      },
      'exercises': {
        title: 'Exercise Management',
        subtitle: 'Manage global exercise library.'
      },
      'muscular-groups': {
        title: 'Muscular Groups',
        subtitle: 'Manage muscle group categories.'
      }
    };

    const config = titles[viewName] || titles['platform-overview'];
    
    document.getElementById('page-title').textContent = config.title;
    document.getElementById('page-subtitle').textContent = config.subtitle;
  }

  // Platform Overview
  async loadPlatformOverview() {
    try {
      // Load platform statistics
      const api = new ApiClient();
      const stats = await api.getPlatformStats();
      
      const content = `
        <div class="dashboard-grid grid-auto">
          <div class="dashboard-card">
            <div class="card-header">
              <div class="card-header-content">
                <div class="card-icon">
                  <i class="fas fa-building"></i>
                </div>
                <div>
                  <h3 class="card-title">Total Gyms</h3>
                  <p class="card-subtitle">Active gym locations</p>
                </div>
              </div>
            </div>
            <div class="card-body">
              <div style="font-size: 2rem; font-weight: bold; color: var(--color-primary);">
                ${stats?.totalGyms || 0}
              </div>
            </div>
          </div>
          
          <div class="dashboard-card">
            <div class="card-header">
              <div class="card-header-content">
                <div class="card-icon">
                  <i class="fas fa-dumbbell"></i>
                </div>
                <div>
                  <h3 class="card-title">Equipment Items</h3>
                  <p class="card-subtitle">Global equipment catalog</p>
                </div>
              </div>
            </div>
            <div class="card-body">
              <div style="font-size: 2rem; font-weight: bold; color: var(--color-success);">
                ${stats?.totalEquipment || 0}
              </div>
            </div>
          </div>
          
          <div class="dashboard-card">
            <div class="card-header">
              <div class="card-header-content">
                <div class="card-icon">
                  <i class="fas fa-list-check"></i>
                </div>
                <div>
                  <h3 class="card-title">Exercises</h3>
                  <p class="card-subtitle">Exercise library</p>
                </div>
              </div>
            </div>
            <div class="card-body">
              <div style="font-size: 2rem; font-weight: bold; color: var(--color-info);">
                ${stats?.totalExercises || 0}
              </div>
            </div>
          </div>
        </div>
        
        <div class="dashboard-card" style="margin-top: var(--spacing-lg);">
          <div class="card-header">
            <h3 class="card-title">Quick Actions</h3>
          </div>
          <div class="card-body">
            <div class="grid-container grid-3">
              <button class="btn btn-primary" onclick="dashboard.loadView('gyms')">
                <i class="fas fa-building"></i>
                Manage Gyms
              </button>
              <button class="btn btn-success" onclick="dashboard.loadView('equipment')">
                <i class="fas fa-tools"></i>
                Manage Equipment
              </button>
              <button class="btn btn-info" onclick="dashboard.loadView('exercises')">
                <i class="fas fa-list-check"></i>
                Manage Exercises
              </button>
            </div>
          </div>
        </div>
      `;
      
      this.setContent(content);
    } catch (error) {
      throw new Error('Failed to load platform overview: ' + error.message);
    }
  }

  // Equipment Management
  async loadEquipmentManagement() {
    try {
      const equipment = await this.managers.equipment.loadEquipment();
      
      const content = `
        <div class="dashboard-card">
          <div class="card-header">
            <div class="card-header-content">
              <h3 class="card-title">Equipment Management</h3>
              <p class="card-subtitle">${equipment.length} items in catalog</p>
            </div>
            <button class="btn btn-primary" onclick="dashboard.openEquipmentModal('create')">
              <i class="fas fa-plus"></i> Add Equipment
            </button>
          </div>
          <div class="card-body">
            <div id="equipment-table"></div>
          </div>
        </div>
      `;
      
      this.setContent(content);
      
      // Initialize data table
      this.components.equipmentTable = new DataTable('#equipment-table', {
        data: equipment,
        columns: this.managers.equipment.getTableColumns(),
        rowActions: this.managers.equipment.getRowActions(),
        emptyMessage: 'No equipment found. Add some equipment to get started.',
        filterable: true,
        sortable: true
      });
      
    } catch (error) {
      throw new Error('Failed to load equipment management: ' + error.message);
    }
  }

  // Exercises Management
  async loadExercisesManagement() {
    try {
      const exercises = await this.managers.exercise.loadExercises();
      
      const content = `
        <div class="dashboard-card">
          <div class="card-header">
            <div class="card-header-content">
              <h3 class="card-title">Exercise Management</h3>
              <p class="card-subtitle">${exercises.length} exercises in library</p>
            </div>
            <button class="btn btn-primary" onclick="dashboard.openExerciseModal('create')">
              <i class="fas fa-plus"></i> Add Exercise
            </button>
          </div>
          <div class="card-body">
            <div id="exercises-table"></div>
          </div>
        </div>
      `;
      
      this.setContent(content);
      
      // Initialize data table
      this.components.exercisesTable = new DataTable('#exercises-table', {
        data: exercises,
        columns: this.managers.exercise.getTableColumns(),
        rowActions: this.managers.exercise.getRowActions(),
        emptyMessage: 'No exercises found. Add some exercises to get started.',
        filterable: true,
        sortable: true
      });
      
    } catch (error) {
      throw new Error('Failed to load exercises management: ' + error.message);
    }
  }

  // Gyms Management
  async loadGymsManagement() {
    try {
      const gyms = await this.managers.gym.loadGyms();
      
      const content = `
        <div class="dashboard-card">
          <div class="card-header">
            <div class="card-header-content">
              <h3 class="card-title">Gyms Management</h3>
              <p class="card-subtitle">${gyms.length} total gyms</p>
            </div>
            <button class="btn btn-primary" onclick="dashboard.openGymModal('create')">
              <i class="fas fa-plus"></i> Add Gym
            </button>
          </div>
          <div class="card-body">
            <div id="gyms-table"></div>
          </div>
        </div>
      `;
      
      this.setContent(content);
      
      // Initialize data table
      this.components.gymsTable = new DataTable('#gyms-table', {
        data: gyms,
        columns: this.managers.gym.getTableColumns(),
        rowActions: this.managers.gym.getRowActions(),
        emptyMessage: 'No gyms found. Add some gyms to get started.',
        filterable: true,
        sortable: true
      });
      
    } catch (error) {
      throw new Error('Failed to load gyms management: ' + error.message);
    }
  }

  async loadMuscularGroupsManagement() {
    this.showComingSoon('Muscular Groups Management');
  }

  // Modal Management
  openEquipmentModal(mode = 'create', equipment = null) {
    const isEdit = mode === 'edit' && equipment;
    const title = isEdit ? 'Edit Equipment' : 'Add Equipment';
    
    const formHtml = this.generateForm(this.managers.equipment.getFormSchema(), equipment);
    
    const modal = new Modal({
      title: title,
      content: formHtml,
      size: 'medium',
      buttons: [
        {
          text: 'Cancel',
          action: 'dismiss',
          className: 'btn btn-secondary'
        },
        {
          text: isEdit ? 'Update' : 'Create',
          action: 'save',
          className: 'btn btn-primary',
          handler: () => this.saveEquipment(isEdit, equipment?.id)
        }
      ]
    });
    
    modal.show();
    this.components.equipmentModal = modal;
  }

  openExerciseModal(mode = 'create', exercise = null) {
    const isEdit = mode === 'edit' && exercise;
    const title = isEdit ? 'Edit Exercise' : 'Add Exercise';
    
    const formHtml = this.generateForm(this.managers.exercise.getFormSchema(), exercise);
    
    const modal = new Modal({
      title: title,
      content: formHtml,
      size: 'large',
      buttons: [
        {
          text: 'Cancel',
          action: 'dismiss',
          className: 'btn btn-secondary'
        },
        {
          text: isEdit ? 'Update' : 'Create',
          action: 'save',
          className: 'btn btn-primary',
          handler: () => this.saveExercise(isEdit, exercise?.id)
        }
      ]
    });
    
    modal.show();
    this.components.exerciseModal = modal;
  }

  openGymModal(mode = 'create', gym = null) {
    const isEdit = mode === 'edit' && gym;
    const title = isEdit ? 'Edit Gym' : 'Add Gym';
    
    const formHtml = this.generateForm(this.managers.gym.getFormSchema(), gym);
    
    const modal = new Modal({
      title: title,
      content: formHtml,
      size: 'medium',
      buttons: [
        {
          text: 'Cancel',
          action: 'dismiss',
          className: 'btn btn-secondary'
        },
        {
          text: isEdit ? 'Update' : 'Create',
          action: 'save',
          className: 'btn btn-primary',
          handler: () => this.saveGym(isEdit, gym?.id)
        }
      ]
    });
    
    modal.show();
    this.components.gymModal = modal;
  }

  // Form Generation
  generateForm(schema, data = null) {
    let formHtml = '<form id="dynamic-form" class="form-grid">';
    
    Object.keys(schema).forEach(fieldName => {
      const field = schema[fieldName];
      const value = data ? (data[fieldName] || '') : '';
      
      formHtml += `<div class="form-group">`;
      formHtml += `<label class="form-label" for="${fieldName}">${field.label}${field.required ? ' *' : ''}</label>`;
      
      if (field.type === 'textarea') {
        formHtml += `<textarea class="form-control" id="${fieldName}" name="${fieldName}" rows="${field.rows || 3}" placeholder="${field.placeholder || ''}">${value}</textarea>`;
      } else if (field.type === 'select') {
        formHtml += `<select class="form-control" id="${fieldName}" name="${fieldName}">`;
        field.options.forEach(option => {
          const selected = value === option.value ? 'selected' : '';
          formHtml += `<option value="${option.value}" ${selected}>${option.label}</option>`;
        });
        formHtml += `</select>`;
      } else {
        formHtml += `<input type="${field.type}" class="form-control" id="${fieldName}" name="${fieldName}" value="${value}" placeholder="${field.placeholder || ''}" ${field.required ? 'required' : ''}>`;
      }
      
      formHtml += `</div>`;
    });
    
    formHtml += '</form>';
    return formHtml;
  }

  // Save handlers
  async saveEquipment(isEdit, id = null) {
    try {
      const formData = getFormData(document.getElementById('dynamic-form'));
      
      if (isEdit) {
        await this.managers.equipment.updateEquipment(id, formData);
      } else {
        await this.managers.equipment.createEquipment(formData);
      }
      
      this.components.equipmentModal.hide();
      
      // Refresh the equipment table
      if (this.components.equipmentTable) {
        const equipment = await this.managers.equipment.loadEquipment();
        this.components.equipmentTable.updateData(equipment);
      }
      
      return false; // Prevent modal from closing automatically
    } catch (error) {
      // Error is already handled by the manager
      return false;
    }
  }

  async saveExercise(isEdit, id = null) {
    try {
      const formData = getFormData(document.getElementById('dynamic-form'));
      
      if (isEdit) {
        await this.managers.exercise.updateExercise(id, formData);
      } else {
        await this.managers.exercise.createExercise(formData);
      }
      
      this.components.exerciseModal.hide();
      
      // Refresh the exercises table
      if (this.components.exercisesTable) {
        const exercises = await this.managers.exercise.loadExercises();
        this.components.exercisesTable.updateData(exercises);
      }
      
      return false;
    } catch (error) {
      return false;
    }
  }

  async saveGym(isEdit, id = null) {
    try {
      const formData = getFormData(document.getElementById('dynamic-form'));
      
      if (isEdit) {
        await this.managers.gym.updateGym(id, formData);
      } else {
        await this.managers.gym.createGym(formData);
      }
      
      this.components.gymModal.hide();
      
      // Refresh the gyms table
      if (this.components.gymsTable) {
        const gyms = await this.managers.gym.loadGyms();
        this.components.gymsTable.updateData(gyms);
      }
      
      return false;
    } catch (error) {
      return false;
    }
  }

  // View exercise details
  viewExerciseDetails(exercise) {
    const content = `
      <div class="exercise-details">
        <div style="display: flex; gap: var(--spacing-lg); margin-bottom: var(--spacing-lg);">
          <div style="flex: 1;">
            <h4>${exercise.name}</h4>
            <p><strong>Type:</strong> ${this.managers.exercise.formatExerciseType(exercise.type)}</p>
            <p><strong>Difficulty:</strong> ${exercise.difficulty}</p>
            ${exercise.instructions ? `<p><strong>Instructions:</strong> ${exercise.instructions}</p>` : ''}
          </div>
        </div>
      </div>
    `;
    
    const modal = new Modal({
      title: 'Exercise Details',
      content: content,
      size: 'large',
      buttons: [
        {
          text: 'Edit',
          action: 'edit',
          className: 'btn btn-primary',
          handler: () => {
            modal.hide();
            this.openExerciseModal('edit', exercise);
          }
        },
        {
          text: 'Close',
          action: 'dismiss',
          className: 'btn btn-secondary'
        }
      ]
    });
    
    modal.show();
  }

  // View gym details
  viewGymDetails(gym) {
    const content = `
      <div class="gym-details">
        <h4>${gym.name}</h4>
        ${gym.contact_name ? `<p><strong>Contact:</strong> ${gym.contact_name}</p>` : ''}
        ${gym.contact_email ? `<p><strong>Email:</strong> ${gym.contact_email}</p>` : ''}
        ${gym.contact_phone ? `<p><strong>Phone:</strong> ${gym.contact_phone}</p>` : ''}
        ${gym.address ? `<p><strong>Address:</strong> ${gym.address}</p>` : ''}
        ${gym.description ? `<p><strong>Description:</strong> ${gym.description}</p>` : ''}
        <p><strong>Created:</strong> ${this.managers.gym.formatDate(gym.created_at)}</p>
        ${gym.deleted_at ? `<p><strong>Status:</strong> <span class="status-deleted">Deleted</span></p>` : ''}
      </div>
    `;
    
    const modal = new Modal({
      title: 'Gym Details',
      content: content,
      size: 'medium',
      buttons: [
        {
          text: 'Edit',
          action: 'edit',
          className: 'btn btn-primary',
          handler: () => {
            modal.hide();
            this.openGymModal('edit', gym);
          }
        },
        {
          text: 'Close',
          action: 'dismiss',
          className: 'btn btn-secondary'
        }
      ]
    });
    
    modal.show();
  }

  // Utility methods
  setContent(html) {
    this.hideLoading();
    this.hideError();
    this.contentArea.innerHTML = html;
  }

  showLoading() {
    document.getElementById('loading-state').style.display = 'flex';
  }

  hideLoading() {
    document.getElementById('loading-state').style.display = 'none';
  }

  showError(message) {
    this.hideLoading();
    document.getElementById('error-message').textContent = message;
    document.getElementById('error-state').style.display = 'flex';
  }

  hideError() {
    document.getElementById('error-state').style.display = 'none';
  }

  showComingSoon(feature) {
    this.setContent(`
      <div class="empty-state">
        <i class="fas fa-clock"></i>
        <h3>Coming Soon</h3>
        <p>${feature} functionality is under development and will be available soon.</p>
      </div>
    `);
  }

  // Refresh current view
  async refreshCurrentView() {
    if (this.currentView) {
      await this.loadView(this.currentView);
    }
  }
}

// Global functions for compatibility
function logout() {
  localStorage.removeItem('auth_token');
  window.location.href = '/';
}

function refreshCurrentView() {
  if (window.dashboard) {
    window.dashboard.refreshCurrentView();
  }
}

function openHelp() {
  notifications.info('Help documentation coming soon!');
}

function openUserSettings() {
  notifications.info('User settings coming soon!');
}

function retryLoad() {
  if (window.dashboard) {
    window.dashboard.refreshCurrentView();
  }
}

// Initialize dashboard when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
  window.dashboard = new VanillaDashboardManager();
  console.log('Vanilla Dashboard initialized');
});