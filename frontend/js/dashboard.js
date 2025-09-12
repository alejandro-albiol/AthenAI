// Dashboard Management
class DashboardManager {
  constructor() {
    this.currentUser = null;
    this.currentGym = null;
    this.currentView = "overview";
    this.apiBase = "/api/v1";

    this.init();
  }

  async init() {
    // Check authentication
    await this.checkAuth();

    // Initialize navigation
    this.initNavigation();

    console.log("Dashboard initialized");
  }

  async checkAuth() {
    try {
      const token = localStorage.getItem("auth_token");
      if (!token) {
        window.location.href = "/";
        return;
      }

      // Verify token and get user info
      const response = await fetch(`${this.apiBase}/auth/validate`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        localStorage.removeItem("auth_token");
        window.location.href = "/";
        return;
      }

      const validationData = await response.json();

      // Check if token is valid and extract user info from claims
      if (!validationData.data || !validationData.data.valid) {
        localStorage.removeItem("auth_token");
        window.location.href = "/";
        return;
      }

      // Extract user info from the claims
      this.currentUser = {
        id: validationData.data.claims.user_id,
        username: validationData.data.claims.username,
        user_type: validationData.data.claims.user_type,
        gym_id: validationData.data.claims.gym_id,
        role: validationData.data.claims.role,
        is_active: validationData.data.claims.is_active,
      };

      this.updateUserInfo();
      this.setupNavigation();
    } catch (error) {
      console.error("Auth check failed:", error);
      window.location.href = "/";
    }
  }

  updateUserInfo() {
    document.getElementById("user-name").textContent =
      this.currentUser.username || this.currentUser.id;
    document.getElementById("user-role").textContent =
      this.currentUser.role || this.currentUser.user_type || "User";
  }

  setupNavigation() {
    const platformAdminNav = document.querySelector(".platform-admin-nav");
    const gymSelector = document.getElementById("gym-selector");

    if (this.currentUser.user_type === "platform_admin") {
      platformAdminNav.style.display = "block";
      gymSelector.style.display = "none";
      // Load platform overview by default for platform admins
      this.loadView("platform-overview");
    } else {
      // Regular gym members and admins should not access this dashboard
      // Redirect them to their gym-specific page
      console.log("Redirecting non-platform user to gym page");
      window.location.href = `/gym/${this.currentUser.gym_id}`;
      return;
    }
  }

  async loadGyms() {
    try {
      const response = await this.apiCall("/gym");
      const gymsData = await response.json();

      // Extract gyms from the API response format
      const gyms = gymsData.data || [];

      const select = document.getElementById("gym-select");

      if (gyms.length === 0) {
        select.innerHTML =
          '<option value="">No gyms available - Create one first</option>';
        select.disabled = true;
      } else {
        select.innerHTML = '<option value="">Select a gym...</option>';
        select.disabled = false;

        gyms.forEach((gym) => {
          const option = document.createElement("option");
          option.value = gym.id;

          // Show deleted gyms with special styling and label
          if (gym.deleted_at) {
            option.textContent = `${gym.name} (Deleted)`;
            option.style.color = "#999";
            option.style.fontStyle = "italic";
            option.disabled = true;
          } else {
            option.textContent = gym.name;
          }

          select.appendChild(option);
        });
      }

      select.addEventListener("change", (e) => {
        this.currentGym = e.target.value;
        this.loadView(this.currentView);
      });
    } catch (error) {
      console.error("Failed to load gyms:", error);
      const select = document.getElementById("gym-select");
      select.innerHTML = '<option value="">Error loading gyms</option>';
      select.disabled = true;
    }
  }

  initNavigation() {
    // Handle navigation clicks
    document.querySelectorAll(".nav-item").forEach((item) => {
      item.addEventListener("click", (e) => {
        e.preventDefault();

        // Update active state
        document
          .querySelectorAll(".nav-item")
          .forEach((nav) => nav.classList.remove("active"));
        item.classList.add("active");

        // Load view
        const view = item.dataset.view;
        this.loadView(view);

        // Close mobile sidebar on navigation (mobile)
        this.closeMobileSidebar();
      });
    });

    // Initialize mobile sidebar toggle
    this.initMobileSidebar();
  }

  initMobileSidebar() {
    const sidebar = document.querySelector(".sidebar");
    const sidebarToggle = document.getElementById("sidebar-toggle");
    const mobileSidebarToggle = document.getElementById(
      "mobile-sidebar-toggle"
    );

    // Mobile sidebar toggle functionality
    const toggleSidebar = () => {
      sidebar.classList.toggle("sidebar-open");
      document.body.classList.toggle("sidebar-open");
    };

    const closeSidebar = () => {
      sidebar.classList.remove("sidebar-open");
      document.body.classList.remove("sidebar-open");
    };

    // Add event listeners
    if (sidebarToggle) {
      sidebarToggle.addEventListener("click", toggleSidebar);
    }

    if (mobileSidebarToggle) {
      mobileSidebarToggle.addEventListener("click", toggleSidebar);
    }

    // Close sidebar when clicking outside on mobile
    document.addEventListener("click", (e) => {
      if (
        window.innerWidth <= 768 &&
        sidebar.classList.contains("sidebar-open") &&
        !sidebar.contains(e.target) &&
        !sidebarToggle?.contains(e.target) &&
        !mobileSidebarToggle?.contains(e.target)
      ) {
        closeSidebar();
      }
    });

    // Close sidebar on window resize if switching to desktop
    window.addEventListener("resize", () => {
      if (window.innerWidth > 768) {
        closeSidebar();
      }
    });
  }

  closeMobileSidebar() {
    if (window.innerWidth <= 768) {
      const sidebar = document.querySelector(".sidebar");
      sidebar.classList.remove("sidebar-open");
      document.body.classList.remove("sidebar-open");
    }
  }

  async loadView(viewName) {
    this.currentView = viewName;

    // Update page title and action button
    this.updatePageHeader(viewName);

    // Show loading state
    this.showLoading();

    try {
      switch (viewName) {
        case "platform-overview":
          await this.loadPlatformOverview();
          break;
        case "gyms":
          await this.loadGymsManagement();
          break;
        default:
          this.showError("View not found");
      }
    } catch (error) {
      console.error("Failed to load view:", error);
      this.showError("Failed to load content");
    }
  }

  updatePageHeader(viewName) {
    const titles = {
      overview: {
        title: "Dashboard Overview",
        subtitle: "Welcome back! Here's what's happening at your gym.",
        action: "Add Member",
      },
      members: {
        title: "Members",
        subtitle: "Manage your gym members and their memberships.",
        action: "Add Member",
      },
      exercises: {
        title: "Custom Exercises",
        subtitle: "Create and manage your custom exercises.",
        action: "Add Exercise",
      },
      equipment: {
        title: "Equipment",
        subtitle: "Manage your gym equipment inventory.",
        action: "Add Equipment",
      },
      workouts: {
        title: "Workout Templates",
        subtitle: "Create and manage workout templates.",
        action: "Add Template",
      },
      "workout-instances": {
        title: "Workout Instances",
        subtitle: "Track member workout sessions.",
        action: "Add Instance",
      },
      analytics: {
        title: "Analytics",
        subtitle: "View insights and performance metrics.",
        action: "Export Data",
      },
      "platform-overview": {
        title: "Platform Overview",
        subtitle: "Monitor platform-wide statistics and activity.",
        action: "Add Gym",
      },
      gyms: {
        title: "Gyms Management",
        subtitle: "Manage all gyms in the platform.",
        action: "Add Gym",
      },
    };

    const config = titles[viewName] || titles["overview"];

    document.getElementById("page-title").textContent = config.title;
    document.getElementById("page-subtitle").textContent = config.subtitle;
    document.getElementById("action-text").textContent = config.action;
  }

  async loadCustomExercises() {
    if (!this.currentGym && this.currentUser.user_type !== "platform_admin") {
      this.showError("No gym selected");
      return;
    }

    try {
      const response = await this.apiCall("/custom-exercise");
      const exercisesData = await response.json();
      const exercises = exercisesData.data || [];

      const html = `
                <div class="data-table">
                    <table>
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Type</th>
                                <th>Difficulty</th>
                                <th>Muscular Groups</th>
                                <th>Created By</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            ${exercises
                              .map(
                                (exercise) => `
                                <tr>
                                    <td><strong>${exercise.name}</strong></td>
                                    <td>${exercise.exercise_type || "-"}</td>
                                    <td><span class="status-badge status-${
                                      exercise.difficulty_level
                                    }">${
                                  exercise.difficulty_level || "-"
                                }</span></td>
                                    <td>${
                                      exercise.muscular_groups
                                        ? exercise.muscular_groups.join(", ")
                                        : "-"
                                    }</td>
                                    <td>${exercise.created_by || "-"}</td>
                                    <td>
                                        <div class="action-buttons">
                                            <button class="btn-icon btn-edit" onclick="dashboard.editExercise('${
                                              exercise.id
                                            }')" title="Edit">
                                                <i class="fas fa-edit"></i>
                                            </button>
                                            <button class="btn-icon btn-delete" onclick="dashboard.deleteExercise('${
                                              exercise.id
                                            }')" title="Delete">
                                                <i class="fas fa-trash"></i>
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            `
                              )
                              .join("")}
                        </tbody>
                    </table>
                </div>
            `;

      if (exercises.length === 0) {
        this.showEmptyState(
          "exercises",
          "No custom exercises found",
          "Create your first custom exercise to get started."
        );
      } else {
        this.setContent(html);
      }
    } catch (error) {
      console.error("Failed to load exercises:", error);
      this.showError("Failed to load exercises");
    }
  }

  async loadPlatformOverview() {
    try {
      // Load platform-wide statistics
      const [gymsResponse] = await Promise.all([this.apiCall("/gym")]);

      const gymsData = await gymsResponse.json();
      const gyms = gymsData.data || [];

      // Calculate platform metrics
      const activeGyms = gyms.filter((gym) => !gym.deleted_at);
      const deletedGyms = gyms.filter((gym) => gym.deleted_at);
      const recentGyms = gyms.filter((gym) => {
        const created = new Date(gym.created_at);
        const thirtyDaysAgo = new Date();
        thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);
        return created > thirtyDaysAgo && !gym.deleted_at;
      });

      const content = `
        <div class="platform-overview">
          <!-- Platform Metrics Cards -->
          <div class="dashboard-grid">
            <div class="dashboard-card stat-card primary">
              <div class="card-icon">
                <i class="fas fa-building"></i>
              </div>
              <span class="stat-number">${activeGyms.length}</span>
              <span class="stat-label">Active Gyms</span>
              <span class="stat-change positive">+${
                recentGyms.length
              } this month</span>
            </div>
            
            <div class="dashboard-card stat-card">
              <div class="card-icon">
                <i class="fas fa-chart-line"></i>
              </div>
              <span class="stat-number">${(
                (activeGyms.length / (activeGyms.length + deletedGyms.length)) *
                100
              ).toFixed(1)}%</span>
              <span class="stat-label">Platform Health</span>
              <span class="stat-change neutral">Active gym ratio</span>
            </div>
            
            <div class="dashboard-card stat-card">
              <div class="card-icon">
                <i class="fas fa-calendar-plus"></i>
              </div>
              <span class="stat-number">${recentGyms.length}</span>
              <span class="stat-label">New This Month</span>
              <span class="stat-change neutral">Last 30 days</span>
            </div>
            
            <div class="dashboard-card stat-card warning">
              <div class="card-icon">
                <i class="fas fa-exclamation-triangle"></i>
              </div>
              <span class="stat-number">${deletedGyms.length}</span>
              <span class="stat-label">Deleted Gyms</span>
              <span class="stat-change negative">Require attention</span>
            </div>
          </div>

          <!-- Recent Activity -->
          <div class="dashboard-section">
            <div class="section-header">
              <h3><i class="fas fa-clock"></i> Recent Activity</h3>
              <a href="#" onclick="dashboard.loadView('gyms')" class="view-all-link">View All Gyms</a>
            </div>
            
            <div class="activity-cards">
              ${gyms
                .slice(0, 5)
                .map(
                  (gym) => `
                <div class="activity-card ${gym.deleted_at ? "deleted" : ""}">
                  <div class="activity-icon">
                    <i class="fas fa-${
                      gym.deleted_at ? "trash" : "building"
                    }"></i>
                  </div>
                  <div class="activity-content">
                    <div class="activity-title">${gym.name}</div>
                    <div class="activity-meta">
                      ${gym.deleted_at ? "Deleted" : "Created"} ${new Date(
                    gym.deleted_at || gym.created_at
                  ).toLocaleDateString()}
                    </div>
                    <div class="activity-location">${
                      gym.address || "No address provided"
                    }</div>
                  </div>
                  <div class="activity-actions">
                    ${
                      !gym.deleted_at
                        ? `
                      <button class="btn btn-sm btn-primary" onclick="dashboard.manageGym('${gym.id}')">
                        <i class="fas fa-cog"></i> Manage
                      </button>
                    `
                        : `
                      <button class="btn btn-sm btn-secondary" onclick="dashboard.restoreGym('${gym.id}')">
                        <i class="fas fa-undo"></i> Restore
                      </button>
                    `
                    }
                  </div>
                </div>
              `
                )
                .join("")}
            </div>
          </div>

          <!-- Platform Statistics -->
          <div class="dashboard-section">
            <div class="section-header">
              <h3><i class="fas fa-chart-bar"></i> Platform Statistics</h3>
            </div>
            
            <div class="stats-grid">
              <div class="stat-item">
                <div class="stat-chart">
                  <div class="chart-bar">
                    <div class="bar-fill" style="width: ${
                      (activeGyms.length / Math.max(gyms.length, 1)) * 100
                    }%"></div>
                  </div>
                </div>
                <div class="stat-info">
                  <span class="stat-title">Gym Retention Rate</span>
                  <span class="stat-value">${(
                    (activeGyms.length / Math.max(gyms.length, 1)) *
                    100
                  ).toFixed(1)}%</span>
                </div>
              </div>
              
              <div class="stat-item">
                <div class="stat-chart">
                  <div class="chart-bar">
                    <div class="bar-fill" style="width: ${
                      (recentGyms.length / Math.max(activeGyms.length, 1)) * 100
                    }%"></div>
                  </div>
                </div>
                <div class="stat-info">
                  <span class="stat-title">Growth Rate (30d)</span>
                  <span class="stat-value">${(
                    (recentGyms.length / Math.max(activeGyms.length, 1)) *
                    100
                  ).toFixed(1)}%</span>
                </div>
              </div>
              
              <div class="stat-item">
                <div class="stat-chart">
                  <div class="chart-bar">
                    <div class="bar-fill" style="width: ${Math.min(
                      (activeGyms.length / 10) * 100,
                      100
                    )}%"></div>
                  </div>
                </div>
                <div class="stat-info">
                  <span class="stat-title">Platform Capacity</span>
                  <span class="stat-value">${activeGyms.length}/∞</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Quick Actions -->
          <div class="dashboard-section">
            <div class="section-header">
              <h3><i class="fas fa-bolt"></i> Quick Actions</h3>
            </div>
            
            <div class="quick-actions">
              <button class="action-button primary" onclick="dashboard.showCreateGymModal()">
                <div class="action-icon">
                  <i class="fas fa-plus"></i>
                </div>
                <div class="action-content">
                  <span class="action-title">Add New Gym</span>
                  <span class="action-subtitle">Onboard a new fitness center</span>
                </div>
              </button>
              
              <button class="action-button" onclick="dashboard.loadView('gyms')">
                <div class="action-icon">
                  <i class="fas fa-building"></i>
                </div>
                <div class="action-content">
                  <span class="action-title">Manage Gyms</span>
                  <span class="action-subtitle">View and manage all gyms</span>
                </div>
              </button>
              
              <button class="action-button" onclick="dashboard.generatePlatformReport()">
                <div class="action-icon">
                  <i class="fas fa-chart-line"></i>
                </div>
                <div class="action-content">
                  <span class="action-title">Platform Report</span>
                  <span class="action-subtitle">Generate analytics report</span>
                </div>
              </button>
              
              <button class="action-button warning" onclick="dashboard.showSystemSettings()">
                <div class="action-icon">
                  <i class="fas fa-cogs"></i>
                </div>
                <div class="action-content">
                  <span class="action-title">System Settings</span>
                  <span class="action-subtitle">Configure platform settings</span>
                </div>
              </button>
            </div>
          </div>
        </div>
      `;

      this.setContent(content);
    } catch (error) {
      console.error("Failed to load platform overview:", error);
      this.showError("Failed to load platform overview");
    }
  }

  // Platform overview helper methods
  generatePlatformReport() {
    // TODO: Implement platform report generation
    console.log("Generate platform report");
  }

  showSystemSettings() {
    // TODO: Implement system settings modal
    console.log("Show system settings");
  }

  async loadGymsManagement() {
    try {
      const response = await this.apiCall("/gym");
      const gymsData = await response.json();
      const gyms = gymsData.data || [];

      const content = `
        <div class="management-view">
          <div class="view-header">
            <div class="search-filters">
              <input type="text" id="gym-search" placeholder="Search gyms..." class="search-input">
              <select id="gym-status-filter" class="filter-select">
                <option value="">All Status</option>
                <option value="active">Active Only</option>
                <option value="deleted">Deleted Only</option>
              </select>
            </div>
            <button class="btn btn-primary" onclick="dashboard.showCreateGymModal()">
              <i class="fas fa-plus"></i> Add Gym
            </button>
          </div>
          
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Gym Name</th>
                  <th>Contact</th>
                  <th>Address</th>
                  <th>Status</th>
                  <th>Created</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody id="gyms-table-body">
                ${gyms
                  .map(
                    (gym) => `
                  <tr class="${
                    gym.deleted_at ? "deleted-row" : ""
                  }" data-gym-id="${gym.id}">
                    <td>
                      <strong>${gym.name}</strong>
                      ${
                        gym.deleted_at
                          ? '<span class="deleted-badge">Deleted</span>'
                          : ""
                      }
                    </td>
                    <td>
                      <div class="contact-info">
                        ${
                          gym.email
                            ? `<div><i class="fas fa-envelope"></i> ${gym.email}</div>`
                            : ""
                        }
                        ${
                          gym.phone
                            ? `<div><i class="fas fa-phone"></i> ${gym.phone}</div>`
                            : ""
                        }
                      </div>
                    </td>
                    <td>${gym.address || "-"}</td>
                    <td>
                      <span class="status-badge ${
                        gym.deleted_at
                          ? "status-deleted"
                          : gym.is_active
                          ? "status-active"
                          : "status-inactive"
                      }">
                        ${
                          gym.deleted_at
                            ? "Deleted"
                            : gym.is_active
                            ? "Active"
                            : "Inactive"
                        }
                      </span>
                    </td>
                    <td>${new Date(gym.created_at).toLocaleDateString()}</td>
                    <td>
                      <div class="action-buttons">
                        ${
                          !gym.deleted_at
                            ? `
                          <button class="btn btn-sm btn-outline" onclick="dashboard.viewGymDetails('${gym.id}')">
                            <i class="fas fa-eye"></i> View
                          </button>
                          <button class="btn btn-sm btn-primary" onclick="dashboard.manageGym('${gym.id}')">
                            <i class="fas fa-cog"></i> Manage
                          </button>
                        `
                            : `
                          <button class="btn btn-sm btn-secondary" onclick="dashboard.restoreGym('${gym.id}')">
                            <i class="fas fa-undo"></i> Restore
                          </button>
                        `
                        }
                      </div>
                    </td>
                  </tr>
                `
                  )
                  .join("")}
              </tbody>
            </table>
            
            ${
              gyms.length === 0
                ? `
              <div class="empty-state">
                <i class="fas fa-building"></i>
                <h3>No Gyms Found</h3>
                <p>Get started by creating your first gym.</p>
                <button class="btn btn-primary" onclick="dashboard.showCreateGymModal()">
                  <i class="fas fa-plus"></i> Create First Gym
                </button>
              </div>
            `
                : ""
            }
          </div>
        </div>
      `;

      this.setContent(content);
      this.setupGymManagementListeners();
    } catch (error) {
      console.error("Failed to load gyms:", error);
      this.showError("Failed to load gyms management");
    }
  }

  setupGymManagementListeners() {
    // Search functionality
    const searchInput = document.getElementById("gym-search");
    const statusFilter = document.getElementById("gym-status-filter");

    if (searchInput && statusFilter) {
      const filterGyms = () => {
        const searchTerm = searchInput.value.toLowerCase();
        const statusValue = statusFilter.value;
        const rows = document.querySelectorAll("#gyms-table-body tr");

        rows.forEach((row) => {
          const gymName = row
            .querySelector("td strong")
            .textContent.toLowerCase();
          const isDeleted = row.classList.contains("deleted-row");

          const matchesSearch = gymName.includes(searchTerm);
          const matchesStatus =
            !statusValue ||
            (statusValue === "active" && !isDeleted) ||
            (statusValue === "deleted" && isDeleted);

          row.style.display = matchesSearch && matchesStatus ? "" : "none";
        });
      };

      searchInput.addEventListener("input", filterGyms);
      statusFilter.addEventListener("change", filterGyms);
    }
  }

  // Gym management actions
  viewGymDetails(gymId) {
    // TODO: Implement gym details modal/view
    console.log("View gym details:", gymId);
  }

  manageGym(gymId) {
    // Switch to gym-specific management context
    this.currentGym = gymId;
    this.loadView("members"); // Load the tenant view for this specific gym
  }

  restoreGym(gymId) {
    // TODO: Implement gym restoration
    console.log("Restore gym:", gymId);
  }

  showCreateGymModal() {
    // TODO: Implement gym creation modal
    console.log("Show create gym modal");
  }

  // Utility methods
  async apiCall(endpoint, options = {}) {
    const token = localStorage.getItem("auth_token");
    const config = {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
        ...options.headers,
      },
      ...options,
    };

    const response = await fetch(`${this.apiBase}${endpoint}`, config);

    if (!response.ok) {
      throw new Error(`API call failed: ${response.status}`);
    }

    return response;
  }

  setContent(html) {
    // Hide loading state overlay
    const loadingState = document.getElementById("loading-state");
    if (loadingState) {
      loadingState.style.display = "none";
    }

    // Hide error state
    const errorState = document.getElementById("error-state");
    if (errorState) {
      errorState.style.display = "none";
    }

    // Update content
    document.getElementById("content-body").innerHTML = html;
  }

  showLoading() {
    // Show loading state overlay
    const loadingState = document.getElementById("loading-state");
    if (loadingState) {
      loadingState.style.display = "flex";
    }

    // Hide error state
    const errorState = document.getElementById("error-state");
    if (errorState) {
      errorState.style.display = "none";
    }

    // Clear content body
    document.getElementById("content-body").innerHTML = "";
  }

  showError(message) {
    // Hide loading state
    const loadingState = document.getElementById("loading-state");
    if (loadingState) {
      loadingState.style.display = "none";
    }

    // Show error state overlay
    const errorState = document.getElementById("error-state");
    const errorMessage = document.getElementById("error-message");
    if (errorState && errorMessage) {
      errorMessage.textContent = message;
      errorState.style.display = "flex";
    }

    // Clear content body
    document.getElementById("content-body").innerHTML = "";
  }

  showEmptyState(type, title, description) {
    this.setContent(`
            <div class="empty-state">
                <i class="fas fa-inbox"></i>
                <h3>${title}</h3>
                <p>${description}</p>
                <button class="btn btn-primary" onclick="openCreateModal()">
                    <i class="fas fa-plus"></i>
                    Create ${type.charAt(0).toUpperCase() + type.slice(1)}
                </button>
            </div>
        `);
  }

  showComingSoon(feature) {
    this.setContent(`
            <div class="empty-state">
                <i class="fas fa-clock"></i>
                <h3>Coming Soon</h3>
                <p>${feature} feature is under development and will be available soon.</p>
            </div>
        `);
  }

  renderRecentInstances(instances) {
    const container = document.getElementById("recent-instances");

    if (!instances || instances.length === 0) {
      container.innerHTML =
        '<p class="text-muted">No recent workout instances found.</p>';
      return;
    }

    const html = `
            <div class="data-table">
                <table>
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>Difficulty</th>
                            <th>Duration</th>
                            <th>Exercises</th>
                            <th>Created</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${instances
                          .map(
                            (instance) => `
                            <tr>
                                <td><strong>${instance.name}</strong></td>
                                <td><span class="status-badge status-${
                                  instance.difficulty_level
                                }">${instance.difficulty_level}</span></td>
                                <td>${
                                  instance.estimated_duration_minutes || "-"
                                } min</td>
                                <td>${instance.total_exercises || 0}</td>
                                <td>${new Date(
                                  instance.created_at
                                ).toLocaleDateString()}</td>
                            </tr>
                        `
                          )
                          .join("")}
                    </tbody>
                </table>
            </div>
        `;

    container.innerHTML = html;
  }
}

// Global functions for button actions
function openCreateModal() {
  // This will be implemented based on the current view
  console.log("Open create modal for:", dashboard.currentView);
}

function logout() {
  localStorage.removeItem("auth_token");
  window.location.href = "/";
}

// Initialize dashboard when DOM is loaded
let dashboard;
document.addEventListener("DOMContentLoaded", function () {
  dashboard = new DashboardManager();
});
