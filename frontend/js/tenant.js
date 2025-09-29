/**
 * Tenant Dashboard Manager
 * Handles gym-specific dashboard for gym_admin, trainer, and member users
 */

class TenantDashboardManager {
  constructor() {
    this.currentUser = null;
    this.currentGym = null;
    this.currentView = "overview";
    this.managers = {};

    this.init();
  }

  async init() {
    try {
      // Wait for all component scripts to load
      await this.waitForDependencies();

      // Initialize content area
      this.contentArea = document.getElementById("content-body");

      if (!this.contentArea) {
        throw new Error("Content area element not found");
      }

      // Show loading state
      this.showLoading("Initializing gym dashboard...");

      // Initialize managers
      this.initializeManagers();

      // Check authentication and verify tenant access
      await this.checkAuth();

      // Initialize navigation
      this.initNavigation();

      // Load initial view
      await this.loadView("overview");

      // Hide loading overlay
      this.hideLoading();
    } catch (error) {
      console.error("Failed to initialize tenant dashboard:", error);
      this.showError("Failed to initialize dashboard: " + error.message);
    }
  }

  async waitForDependencies() {
    const requiredClasses = [
      "Modal",
      "notifications",
      "ApiClient",
      "DataTable",
      "UserManager",
      "appState",
    ];

    let attempts = 0;
    const maxAttempts = 50;

    while (attempts < maxAttempts) {
      const allLoaded = requiredClasses.every((className) => {
        return window[className] !== undefined;
      });

      if (allLoaded) {
        return;
      }

      await new Promise((resolve) => setTimeout(resolve, 100));
      attempts++;
    }

    console.warn("Some tenant dashboard dependencies may not have loaded");
  }

  initializeManagers() {
    try {
      this.managers.user = new UserManager();
      this.managers.gym = window.GymManager ? new GymManager() : null;
      this.managers.workout = window.WorkoutManager
        ? new WorkoutManager()
        : null;
    } catch (error) {
      console.error("Error initializing managers:", error);
    }
  }

  async checkAuth() {
    try {
      // Check for stored auth token (both localStorage and sessionStorage)
      const authToken =
        localStorage.getItem("auth_token") ||
        sessionStorage.getItem("auth_token");
      const userInfo =
        localStorage.getItem("user_info") ||
        sessionStorage.getItem("user_info");

      if (!authToken || !userInfo) {
        console.error("No authentication data found");
        this.redirectToLogin();
        return;
      }

      this.currentUser = JSON.parse(userInfo);

      // Verify this is a tenant user
      if (this.currentUser.user_type !== "tenant_user") {
        console.error("Invalid user type for tenant dashboard");
        this.redirectToLogin();
        return;
      }

      // Verify user has gym_id
      if (!this.currentUser.gym_id) {
        console.error("No gym association found for user");
        this.showError("No gym association found. Please contact support.");
        return;
      }

      // Set user state
      appState.setState({ user: this.currentUser });

      // Load gym information
      await this.loadGymInfo();

      // Update UI with user info
      this.updateUserInfo();
      this.setupNavigation();
    } catch (error) {
      console.error("Authentication check failed:", error);
      this.redirectToLogin();
    }
  }

  async loadGymInfo() {
    try {
      const api = new ApiClient();
      const response = await api.getGym(this.currentUser.gym_id);

      if (response && response.data) {
        this.currentGym = response.data;
        appState.setState({ currentGym: this.currentGym });

        // Update gym branding in header
        this.updateGymBranding(response.data);

        // Update page title with gym name
        document.title = `${response.data.name} - AthenAI Dashboard`;
      }
    } catch (error) {
      console.error("Failed to load gym info:", error);
    }
  }

  updateGymBranding(gymData) {
    // Update gym name in header
    const gymNameEl = document.getElementById("gym-name");
    if (gymNameEl) {
      gymNameEl.innerHTML = `
        <div class="gym-brand-info">
          <span class="gym-name">${gymData.name}</span>
          ${
            gymData.description
              ? `<span class="gym-tagline">${gymData.description}</span>`
              : ""
          }
        </div>
      `;
    }

    // Add gym status indicator
    const navBrand = document.querySelector(".nav-brand");
    if (navBrand && !navBrand.querySelector(".gym-status")) {
      const statusIndicator = document.createElement("div");
      statusIndicator.className = "gym-status";
      statusIndicator.innerHTML = `
        <div class="status-dot active" title="Connected to ${gymData.name}"></div>
      `;
      navBrand.appendChild(statusIndicator);
    }
  }

  updateUserInfo() {
    const userNameEl = document.getElementById("user-name");
    const userRoleEl = document.getElementById("user-role");

    if (userNameEl) {
      userNameEl.textContent =
        this.currentUser.username || this.currentUser.email;
    }

    if (userRoleEl) {
      userRoleEl.textContent = this.currentUser.role || "Member";
    }
  }

  setupNavigation() {
    // Setup navigation based on user role
    const navItems = document.querySelectorAll(".nav-item");

    // Define role-based permissions
    const permissions = {
      member: {
        allowed: ["overview", "workouts"],
        description: "Member Access",
      },
      trainer: {
        allowed: ["overview", "members", "workouts", "equipment"],
        description: "Trainer Access",
      },
      gym_admin: {
        allowed: [
          "overview",
          "members",
          "trainers",
          "workouts",
          "equipment",
          "analytics",
        ],
        description: "Admin Access",
      },
    };

    const userPermissions =
      permissions[this.currentUser.role] || permissions.member;

    navItems.forEach((item) => {
      const view = item.dataset.view;

      if (!userPermissions.allowed.includes(view)) {
        item.style.display = "none";
      } else {
        item.style.display = "flex";

        // Add role indicator for restricted features
        if (
          this.currentUser.role !== "gym_admin" &&
          ["analytics", "trainers"].includes(view)
        ) {
          const roleIndicator = item.querySelector(".role-indicator");
          if (!roleIndicator) {
            const indicator = document.createElement("span");
            indicator.className = "role-indicator";
            indicator.innerHTML = '<i class="fas fa-crown"></i>';
            indicator.title = "Admin Only";
            item.appendChild(indicator);
          }
        }
      }
    });

    // Add role badge to user info
    this.addRoleBadge(userPermissions.description);

    // Setup logout functionality
    const logoutBtn = document.getElementById("logout-btn");
    if (logoutBtn) {
      logoutBtn.addEventListener("click", () => this.logout());
    }
  }

  addRoleBadge(roleDescription) {
    const userRole = document.getElementById("user-role");
    if (userRole) {
      userRole.innerHTML = `
        <span class="role-badge role-${this.currentUser.role}">
          ${this.formatRole(this.currentUser.role)}
        </span>
        <span class="role-description">${roleDescription}</span>
      `;
    }
  }

  formatRole(role) {
    const roleMap = {
      gym_admin: "Admin",
      trainer: "Trainer",
      member: "Member",
    };
    return roleMap[role] || role;
  }

  initNavigation() {
    const navItems = document.querySelectorAll(".nav-item");

    navItems.forEach((item) => {
      item.addEventListener("click", async (e) => {
        e.preventDefault();

        const view = item.dataset.view;
        await this.loadView(view);

        // Update active state
        navItems.forEach((nav) => nav.classList.remove("active"));
        item.classList.add("active");
      });
    });
  }

  async loadView(viewName) {
    try {
      this.currentView = viewName;

      // Update page title
      const pageTitleEl = document.getElementById("page-title");
      const titleMap = {
        overview: "Overview",
        members: "Members",
        trainers: "Trainers",
        workouts: "Workouts",
        equipment: "Equipment",
        analytics: "Analytics",
      };

      if (pageTitleEl) {
        pageTitleEl.textContent = titleMap[viewName] || "Dashboard";
      }

      // Load view content
      switch (viewName) {
        case "overview":
          await this.loadOverview();
          break;
        case "members":
          await this.loadMembers();
          break;
        case "trainers":
          await this.loadTrainers();
          break;
        case "workouts":
          await this.loadWorkouts();
          break;
        case "equipment":
          await this.loadEquipment();
          break;
        case "analytics":
          await this.loadAnalytics();
          break;
        default:
          this.showError("Unknown view: " + viewName);
      }
    } catch (error) {
      console.error(`Error loading view ${viewName}:`, error);
      this.showError(`Failed to load ${viewName}: ${error.message}`);
    }
  }

  async loadOverview() {
    // Check if this is a first-time user and show onboarding
    await this.checkAndShowOnboarding();

    const content = `
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">
            <i class="fas fa-users"></i>
          </div>
          <div class="stat-value" id="members-count">-</div>
          <div class="stat-label">Active Members</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <i class="fas fa-user-tie"></i>
          </div>
          <div class="stat-value" id="trainers-count">-</div>
          <div class="stat-label">Trainers</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <i class="fas fa-dumbbell"></i>
          </div>
          <div class="stat-value" id="workouts-count">-</div>
          <div class="stat-label">Total Workouts</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <i class="fas fa-chart-line"></i>
          </div>
          <div class="stat-value" id="activity-score">-</div>
          <div class="stat-label">Activity Score</div>
        </div>
      </div>
      
      <div class="tenant-card">
        <div class="card-header">
          <h3 class="card-title">
            <i class="fas fa-clock"></i>
            Recent Activity
          </h3>
        </div>
        <div class="card-body">
          <div class="empty-state">
            <i class="fas fa-history"></i>
            <h3>No Recent Activity</h3>
            <p>Recent gym activity will appear here once members start using the system.</p>
          </div>
        </div>
      </div>
    `;

    this.setContent(content);

    // Load statistics
    await this.loadStatistics();
  }

  async loadStatistics() {
    try {
      if (!this.managers.user || !this.currentUser.gym_id) return;

      // Load member count
      const users = await this.managers.user.loadUsers(this.currentUser.gym_id);
      const members = users.filter((user) => user.role === "member");
      const trainers = users.filter((user) => user.role === "trainer");

      document.getElementById("members-count").textContent = members.length;
      document.getElementById("trainers-count").textContent = trainers.length;

      // TODO: Load workout and activity statistics when available
      document.getElementById("workouts-count").textContent = "0";
      document.getElementById("activity-score").textContent = "100%";
    } catch (error) {
      console.error("Error loading statistics:", error);
    }
  }

  async loadMembers() {
    if (!this.hasPermission("manage_members")) {
      this.showError("You don't have permission to view members");
      return;
    }

    await this.loadUsersByRole("member", "Members");
  }

  async loadTrainers() {
    if (!this.hasPermission("manage_trainers")) {
      this.showError("You don't have permission to view trainers");
      return;
    }

    await this.loadUsersByRole("trainer", "Trainers");
  }

  async loadUsersByRole(role, title) {
    try {
      // Set gym context for user manager
      this.managers.user.setGymContext(this.currentUser.gym_id);

      const content = `
        <div class="tenant-card">
          <div class="card-header">
            <h3 class="card-title">
              <i class="fas fa-${role === "member" ? "users" : "user-tie"}"></i>
              ${title}
            </h3>
            <div class="card-actions">
              ${
                this.hasPermission("create_users")
                  ? `
                <button class="btn btn-primary" onclick="tenantDashboard.openUserModal('create', '${role}')">
                  <i class="fas fa-plus"></i> Add ${
                    role === "member" ? "Member" : "Trainer"
                  }
                </button>
              `
                  : ""
              }
            </div>
          </div>
          <div class="card-body" id="users-table-container">
            <div class="loading-state">
              <i class="fas fa-spinner fa-spin"></i>
              <p>Loading ${title.toLowerCase()}...</p>
            </div>
          </div>
        </div>
      `;

      this.setContent(content);

      // Load users
      const allUsers = await this.managers.user.loadUsers(
        this.currentUser.gym_id
      );
      const filteredUsers = allUsers.filter((user) => user.role === role);

      // Create table
      const tableContainer = document.getElementById("users-table-container");
      if (tableContainer && filteredUsers) {
        if (filteredUsers.length === 0) {
          tableContainer.innerHTML = `
            <div class="empty-state">
              <i class="fas fa-${role === "member" ? "users" : "user-tie"}"></i>
              <h3>No ${title}</h3>
              <p>No ${title.toLowerCase()} found in your gym.</p>
              ${
                this.hasPermission("create_users")
                  ? `
                <button class="btn btn-primary" onclick="tenantDashboard.openUserModal('create', '${role}')">
                  <i class="fas fa-plus"></i> Add First ${
                    role === "member" ? "Member" : "Trainer"
                  }
                </button>
              `
                  : ""
              }
            </div>
          `;
        } else {
          const dataTable = new DataTable(
            tableContainer,
            filteredUsers,
            this.managers.user.getTableColumns(),
            {
              title: title,
              actions: this.getTableActions(),
              searchable: true,
              sortable: true,
              pagination: true,
              pageSize: 20,
              emptyMessage: `No ${title.toLowerCase()} found`,
            }
          );
        }
      }
    } catch (error) {
      console.error(`Error loading ${title.toLowerCase()}:`, error);
      this.showError(`Failed to load ${title.toLowerCase()}: ${error.message}`);
    }
  }

  async loadWorkouts() {
    const content = `
      <div class="tenant-card">
        <div class="card-header">
          <h3 class="card-title">
            <i class="fas fa-dumbbell"></i>
            Workout Management
          </h3>
          <div class="card-actions">
            <button class="btn btn-primary">
              <i class="fas fa-plus"></i> Create Workout
            </button>
          </div>
        </div>
        <div class="card-body">
          <div class="empty-state">
            <i class="fas fa-dumbbell"></i>
            <h3>Workout Management</h3>
            <p>Workout management features will be available soon.</p>
          </div>
        </div>
      </div>
    `;

    this.setContent(content);
  }

  async loadEquipment() {
    const content = `
      <div class="tenant-card">
        <div class="card-header">
          <h3 class="card-title">
            <i class="fas fa-cog"></i>
            Equipment Management
          </h3>
          <div class="card-actions">
            <button class="btn btn-primary">
              <i class="fas fa-plus"></i> Add Equipment
            </button>
          </div>
        </div>
        <div class="card-body">
          <div class="empty-state">
            <i class="fas fa-cog"></i>
            <h3>Equipment Management</h3>
            <p>Equipment management features will be available soon.</p>
          </div>
        </div>
      </div>
    `;

    this.setContent(content);
  }

  async loadAnalytics() {
    const content = `
      <div class="tenant-card">
        <div class="card-header">
          <h3 class="card-title">
            <i class="fas fa-chart-bar"></i>
            Gym Analytics
          </h3>
        </div>
        <div class="card-body">
          <div class="empty-state">
            <i class="fas fa-chart-bar"></i>
            <h3>Analytics Dashboard</h3>
            <p>Advanced analytics and reporting features will be available soon.</p>
          </div>
        </div>
      </div>
    `;

    this.setContent(content);
  }

  hasPermission(permission) {
    const role = this.currentUser.role;

    const permissions = {
      gym_admin: [
        "manage_members",
        "manage_trainers",
        "create_users",
        "edit_users",
        "delete_users",
      ],
      trainer: ["manage_members", "create_users", "edit_users"],
      member: [],
    };

    return permissions[role]?.includes(permission) || false;
  }

  getTableActions() {
    const actions = [];

    if (this.hasPermission("edit_users")) {
      actions.push({
        label: "Edit",
        icon: "fas fa-edit",
        action: (user) => this.openUserModal("edit", null, user),
      });
    }

    actions.push({
      label: "View",
      icon: "fas fa-eye",
      action: (user) => this.openUserModal("view", null, user),
    });

    if (this.hasPermission("delete_users")) {
      actions.push({
        label: "Delete",
        icon: "fas fa-trash",
        action: (user) => this.deleteUser(user),
        variant: "danger",
      });
    }

    return actions;
  }

  async openUserModal(mode = "create", targetRole = null, userData = null) {
    try {
      const isEdit = mode === "edit" && userData;
      const isView = mode === "view" && userData;

      const title = isView
        ? `View User: ${userData.username}`
        : isEdit
        ? `Edit User: ${userData.username}`
        : `Add New ${
            targetRole === "member"
              ? "Member"
              : targetRole === "trainer"
              ? "Trainer"
              : "User"
          }`;

      const fields = this.managers.user.getFormFields();

      // Filter fields based on role and permissions
      const filteredFields = this.filterFormFields(fields, mode, targetRole);

      const formHtml = this.generateFormHtml(
        filteredFields,
        isEdit || isView ? userData : { role: targetRole },
        isView
      );

      const modal = new Modal({
        title: title,
        size: "lg",
        content: `
          <form id="user-form">
            ${formHtml}
          </form>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Cancel</button>
            ${
              !isView
                ? `
              <button type="submit" class="btn btn-primary" id="save-user-btn">
                <i class="fas fa-save"></i> ${isEdit ? "Update" : "Create"} User
              </button>
            `
                : ""
            }
          </div>
        `,
      });

      modal.show();

      if (!isView) {
        const form = modal.element.querySelector("#user-form");
        const saveBtn = modal.element.querySelector("#save-user-btn");

        const handleSubmit = async (e) => {
          e.preventDefault();

          try {
            saveBtn.disabled = true;
            saveBtn.innerHTML =
              '<i class="fas fa-spinner fa-spin"></i> Saving...';

            const formData = getFormData(form);
            const validation = this.managers.user.validateUserData(formData);

            if (validation.length > 0) {
              notifications.error("Please check the form for errors");
              return;
            }

            const preparedData = this.managers.user.prepareDataForSubmission(
              formData,
              isEdit
            );

            if (isEdit) {
              await this.managers.user.updateUser(userData.id, preparedData);
              notifications.success("User updated successfully");
            } else {
              await this.managers.user.createUser(preparedData);
              notifications.success("User created successfully");
            }

            modal.hide();

            // Reload current view
            await this.loadView(this.currentView);
          } catch (error) {
            console.error("Error saving user:", error);
            notifications.error(
              `Failed to ${isEdit ? "update" : "create"} user: ${error.message}`
            );
          } finally {
            saveBtn.disabled = false;
            saveBtn.innerHTML = `<i class="fas fa-save"></i> ${
              isEdit ? "Update" : "Create"
            } User`;
          }
        };

        form.addEventListener("submit", handleSubmit);
        saveBtn.addEventListener("click", handleSubmit);
      }
    } catch (error) {
      console.error("Error opening user modal:", error);
      notifications.error("Failed to open user form");
    }
  }

  filterFormFields(fields, mode, targetRole) {
    // Filter fields based on user permissions and target role
    return fields.filter((field) => {
      // Always show basic fields
      if (["username", "email", "phone"].includes(field.name)) {
        return true;
      }

      // Role field handling
      if (field.name === "role") {
        if (this.currentUser.role === "gym_admin") {
          // Gym admins can assign trainer and member roles
          field.options = field.options.filter((opt) =>
            ["trainer", "member"].includes(opt.value)
          );
          return true;
        } else if (this.currentUser.role === "trainer") {
          // Trainers can only create members
          field.options = field.options.filter((opt) => opt.value === "member");
          return true;
        }
        return false;
      }

      // Password field only for creation
      if (field.name === "password") {
        return mode === "create";
      }

      return true;
    });
  }

  generateFormHtml(fields, data = {}, readonly = false) {
    return fields
      .map((field) => {
        const value = data[field.name] || field.defaultValue || "";
        const readonlyAttr = readonly ? "readonly disabled" : "";

        switch (field.type) {
          case "select":
            const options = field.options
              .map(
                (opt) =>
                  `<option value="${opt.value}" ${
                    value === opt.value ? "selected" : ""
                  }>${opt.label}</option>`
              )
              .join("");

            return `
            <div class="form-group">
              <label for="${field.name}">${field.label}${
              field.required ? " *" : ""
            }</label>
              <select name="${field.name}" id="${
              field.name
            }" class="form-control" ${readonlyAttr} ${
              field.required ? "required" : ""
            }>
                <option value="">Select ${field.label}</option>
                ${options}
              </select>
            </div>
          `;

          case "textarea":
            return `
            <div class="form-group">
              <label for="${field.name}">${field.label}${
              field.required ? " *" : ""
            }</label>
              <textarea name="${field.name}" id="${
              field.name
            }" class="form-control" ${readonlyAttr} ${
              field.required ? "required" : ""
            }>${value}</textarea>
            </div>
          `;

          default:
            return `
            <div class="form-group">
              <label for="${field.name}">${field.label}${
              field.required ? " *" : ""
            }</label>
              <input type="${field.type}" name="${field.name}" id="${
              field.name
            }" class="form-control" value="${value}" ${readonlyAttr} ${
              field.required ? "required" : ""
            }>
            </div>
          `;
        }
      })
      .join("");
  }

  async deleteUser(user) {
    if (!this.hasPermission("delete_users")) {
      notifications.error("You don't have permission to delete users");
      return;
    }

    const confirmed = await this.confirmAction(
      "Delete User",
      `Are you sure you want to delete ${user.username}? This action cannot be undone.`,
      "danger"
    );

    if (confirmed) {
      try {
        await this.managers.user.deleteUser(user.id);
        notifications.success("User deleted successfully");
        await this.loadView(this.currentView);
      } catch (error) {
        console.error("Error deleting user:", error);
        notifications.error("Failed to delete user: " + error.message);
      }
    }
  }

  async confirmAction(title, message, variant = "primary") {
    return new Promise((resolve) => {
      const modal = new Modal({
        title: title,
        content: `
          <p>${message}</p>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal" onclick="this.closest('.modal').dispatchEvent(new CustomEvent('resolve', {detail: false}))">Cancel</button>
            <button type="button" class="btn btn-${variant}" onclick="this.closest('.modal').dispatchEvent(new CustomEvent('resolve', {detail: true}))">Confirm</button>
          </div>
        `,
      });

      modal.element.addEventListener("resolve", (e) => {
        modal.hide();
        resolve(e.detail);
      });

      modal.show();
    });
  }

  setContent(content) {
    if (this.contentArea) {
      this.contentArea.innerHTML = content;
    }
  }

  showLoading(message = "Loading...") {
    const overlay = document.getElementById("loading-overlay");
    if (overlay) {
      overlay.querySelector("p").textContent = message;
      overlay.style.display = "flex";
    }
  }

  hideLoading() {
    const overlay = document.getElementById("loading-overlay");
    if (overlay) {
      overlay.style.display = "none";
    }
  }

  showError(message) {
    if (this.contentArea) {
      this.contentArea.innerHTML = `
        <div class="tenant-card">
          <div class="card-body">
            <div class="empty-state">
              <i class="fas fa-exclamation-triangle" style="color: #ef4444;"></i>
              <h3>Error</h3>
              <p>${message}</p>
              <button class="btn btn-primary" onclick="location.reload()">
                <i class="fas fa-refresh"></i> Retry
              </button>
            </div>
          </div>
        </div>
      `;
    }
    notifications.error(message);
  }

  redirectToLogin() {
    // Clear stored auth data
    localStorage.removeItem("auth_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user_info");

    // Redirect to login
    const currentOrigin = window.location.origin;
    let currentPath = window.location.pathname.replace(
      /\/pages\/tenant\/.*/,
      ""
    );
    if (currentPath.endsWith("/")) {
      currentPath = currentPath.slice(0, -1);
    }
    window.location.href = `${currentOrigin}${currentPath}/`;
  }

  async logout() {
    try {
      // Call logout API if available
      const api = new ApiClient();
      await api.logout();
    } catch (error) {
      console.error("Logout API call failed:", error);
    } finally {
      this.redirectToLogin();
    }
  }

  async checkAndShowOnboarding() {
    // Check if user has seen onboarding
    const onboardingKey = `onboarding_seen_${this.currentUser.user_id}`;
    const hasSeenOnboarding = localStorage.getItem(onboardingKey);

    if (!hasSeenOnboarding) {
      this.showOnboardingWelcome();
      localStorage.setItem(onboardingKey, "true");
    }
  }

  showOnboardingWelcome() {
    const roleWelcomes = {
      gym_admin: {
        title: `Welcome to ${this.currentGym?.name || "your gym"}, Admin!`,
        message:
          "You have full access to manage members, trainers, equipment, and analytics.",
        features: [
          "Manage gym members and trainers",
          "View comprehensive analytics",
          "Configure gym equipment",
          "Monitor all gym activities",
        ],
      },
      trainer: {
        title: `Welcome to ${this.currentGym?.name || "your gym"}, Trainer!`,
        message:
          "You can manage members, create workouts, and track equipment usage.",
        features: [
          "Manage assigned members",
          "Create custom workouts",
          "Track equipment usage",
          "Monitor member progress",
        ],
      },
      member: {
        title: `Welcome to ${this.currentGym?.name || "your gym"}!`,
        message:
          "Track your fitness journey with personalized workouts and progress monitoring.",
        features: [
          "Access personalized workouts",
          "Track your progress",
          "View available equipment",
          "Monitor your fitness goals",
        ],
      },
    };

    const welcome = roleWelcomes[this.currentUser.role] || roleWelcomes.member;

    // Show welcome notification
    notifications.success(`${welcome.title} ${welcome.message}`);

    // Add onboarding overlay to overview (optional)
    setTimeout(() => {
      const contentBody = document.getElementById("content-body");
      if (contentBody) {
        const onboardingBanner = document.createElement("div");
        onboardingBanner.className = "onboarding-banner";
        onboardingBanner.innerHTML = `
          <div class="onboarding-content">
            <div class="onboarding-header">
              <h3>${welcome.title}</h3>
              <button class="close-onboarding" onclick="this.closest('.onboarding-banner').remove()">
                <i class="fas fa-times"></i>
              </button>
            </div>
            <p>${welcome.message}</p>
            <div class="onboarding-features">
              <h4>What you can do:</h4>
              <ul>
                ${welcome.features
                  .map(
                    (feature) =>
                      `<li><i class="fas fa-check"></i> ${feature}</li>`
                  )
                  .join("")}
              </ul>
            </div>
          </div>
        `;
        contentBody.insertBefore(onboardingBanner, contentBody.firstChild);
      }
    }, 2000);
  }
}

// Initialize tenant dashboard when page loads
let tenantDashboard;

document.addEventListener("DOMContentLoaded", () => {
  tenantDashboard = new TenantDashboardManager();
});

// Make tenantDashboard globally available for onclick handlers
window.tenantDashboard = tenantDashboard;
