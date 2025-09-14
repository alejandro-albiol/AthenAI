// Dashboard Management
class DashboardManager {
  constructor() {
    this.currentUser = null;
    this.currentView = "overview";
    this.apiBase = "/api/v1";
    this.api = new API(); // Initialize API instance

    this.init();
  }

  async init() {
    // Initialize content area
    this.contentArea = document.getElementById("content-body");

    // Check authentication
    await this.checkAuth();

    // Initialize navigation
    this.initNavigation();
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
      this.setupModalHandlers();
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
    const platformAdminNav = document.querySelectorAll(".platform-admin-nav");

    if (this.currentUser.user_type === "platform_admin") {
      // Show platform admin navigation
      platformAdminNav.forEach((nav) => (nav.style.display = "block"));

      // Update page header for platform admin
      document.getElementById("page-title").textContent = "Platform Dashboard";
      document.getElementById("page-subtitle").textContent =
        "Manage your AthenAI platform and gym network.";

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

  setupModalHandlers() {
    // Close modals when clicking outside
    window.addEventListener("click", (event) => {
      const muscularGroupModal = document.getElementById("muscularGroupModal");
      const equipmentModal = document.getElementById("equipmentModal");
      const helpModal = document.getElementById("helpModal");
      const exerciseModal = document.getElementById("exerciseModal");
      const exerciseDetailsModal = document.getElementById(
        "exerciseDetailsModal"
      );

      if (event.target === muscularGroupModal) {
        this.closeMuscularGroupModal();
      }
      if (event.target === equipmentModal) {
        this.closeEquipmentModal();
      }
      if (event.target === helpModal) {
        this.closeHelp();
      }
      if (event.target === exerciseModal) {
        this.closeExerciseModal();
      }
      if (event.target === exerciseDetailsModal) {
        this.closeExerciseDetailsModal();
      }
    });
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
        case "equipment":
          await this.loadEquipmentManagement();
          break;
        case "exercises":
          await this.loadExercisesManagement();
          break;
        case "muscular-groups":
          await this.loadMuscularGroupsManagement();
          break;
        case "workout-templates":
          await this.loadWorkoutTemplatesManagement();
          break;
        case "platform-analytics":
          await this.loadPlatformAnalytics();
          break;
        case "platform-settings":
          await this.loadPlatformSettings();
          break;
        default:
          this.showComingSoon(viewName);
      }

      // Hide loading state after successful content load
      this.hideLoading();
    } catch (error) {
      console.error("Failed to load view:", error);
      this.hideLoading();
      this.showError("Failed to load content");
    }
  }

  updatePageHeader(viewName) {
    const titles = {
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
      equipment: {
        title: "Equipment Management",
        subtitle: "Manage global equipment catalog.",
        action: "Add Equipment",
      },
      exercises: {
        title: "Exercise Management",
        subtitle: "Manage global exercise library.",
        action: "Add Exercise",
      },
      "muscular-groups": {
        title: "Muscular Groups",
        subtitle: "Manage muscle group categories.",
        action: "Add Group",
      },
      "workout-templates": {
        title: "Workout Templates",
        subtitle: "Manage global workout templates.",
        action: "Add Template",
      },
      "platform-analytics": {
        title: "Platform Analytics",
        subtitle: "View platform-wide insights and metrics.",
        action: "Export Data",
      },
      "platform-settings": {
        title: "Platform Settings",
        subtitle: "Configure system-wide settings.",
        action: "Save Settings",
      },
    };

    const config = titles[viewName] || titles["platform-overview"];

    document.getElementById("page-title").textContent = config.title;
    document.getElementById("page-subtitle").textContent = config.subtitle;
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
          <!-- Gym Status Overview -->
          <div class="dashboard-section">
            <div class="section-header">
              <h3><i class="fas fa-building"></i> Gym Network Status</h3>
              <a href="#" onclick="dashboard.loadView('gyms')" class="view-all-link">Manage All Gyms</a>
            </div>
            
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
              
              <div class="dashboard-card stat-card success">
                <div class="card-icon">
                  <i class="fas fa-users"></i>
                </div>
                <span class="stat-number">${this.calculateTotalMembers(
                  gyms
                )}</span>
                <span class="stat-label">Total Members</span>
                <span class="stat-change positive">Across all gyms</span>
              </div>
              
              <div class="dashboard-card stat-card ${
                deletedGyms.length > 0 ? "warning" : "success"
              }">
                <div class="card-icon">
                  <i class="fas fa-${
                    deletedGyms.length > 0
                      ? "exclamation-triangle"
                      : "shield-alt"
                  }"></i>
                </div>
                <span class="stat-number">${deletedGyms.length}</span>
                <span class="stat-label">Inactive Gyms</span>
                <span class="stat-change ${
                  deletedGyms.length > 0 ? "negative" : "neutral"
                }">${
        deletedGyms.length > 0 ? "Require attention" : "All systems operational"
      }</span>
              </div>
            </div>
          </div>

          <!-- Recent Gym Activity -->
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
                      <span class="status-badge ${
                        gym.deleted_at ? "inactive" : "active"
                      }">
                        ${gym.deleted_at ? "Inactive" : "Active"}
                      </span>
                      • ${gym.deleted_at ? "Deleted" : "Created"} ${new Date(
                    gym.deleted_at || gym.created_at
                  ).toLocaleDateString()}
                    </div>
                    <div class="activity-location">
                      <i class="fas fa-map-marker-alt"></i> ${
                        gym.address || "No address provided"
                      }
                    </div>
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

  // Helper method to calculate total members
  calculateTotalMembers(gyms) {
    // For now, return 0 since we don't have member data
    // In the future, this could make an API call to get actual member counts
    // or sum up member counts from each gym if that data is available
    const activeGyms = gyms.filter((gym) => !gym.deleted_at);

    // TODO: Implement actual member count API call
    // This should call something like: GET /api/v1/platform/members/count
    // or aggregate member counts from each gym
    return 0;
  }

  // Platform overview helper methods
  generatePlatformReport() {
    // TODO: Implement platform report generation
    this.showToast("Platform report generation coming soon!", "info");
  }

  showSystemSettings() {
    // TODO: Implement system settings modal
    this.showToast("System settings panel coming soon!", "info");
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
  async viewGymDetails(gymId) {
    try {
      const response = await this.apiCall(`/gym/${gymId}`);
      const gymData = await response.json();
      const gym = gymData.data;

      const modal = document.createElement("div");
      modal.className = "modal";
      modal.innerHTML = `
        <div class="modal-content large">
          <div class="modal-header">
            <h2><i class="fas fa-building"></i> ${gym.name}</h2>
            <span class="close" onclick="this.closest('.modal').remove()">&times;</span>
          </div>
          <div class="modal-body">
            <div class="gym-details-grid">
              <div class="detail-section">
                <h4><i class="fas fa-info-circle"></i> General Information</h4>
                <div class="detail-item">
                  <strong>Name:</strong> ${gym.name}
                </div>
                <div class="detail-item">
                  <strong>Address:</strong> ${gym.address || "Not provided"}
                </div>
                <div class="detail-item">
                  <strong>Description:</strong> ${
                    gym.description || "No description"
                  }
                </div>
                <div class="detail-item">
                  <strong>Status:</strong> 
                  <span class="status-badge ${
                    gym.deleted_at ? "inactive" : "active"
                  }">
                    ${gym.deleted_at ? "Inactive" : "Active"}
                  </span>
                </div>
              </div>
              
              <div class="detail-section">
                <h4><i class="fas fa-user"></i> Contact Information</h4>
                <div class="detail-item">
                  <strong>Contact Name:</strong> ${
                    gym.contact_name || "Not provided"
                  }
                </div>
                <div class="detail-item">
                  <strong>Email:</strong> ${gym.contact_email || "Not provided"}
                </div>
                <div class="detail-item">
                  <strong>Phone:</strong> ${gym.contact_phone || "Not provided"}
                </div>
              </div>
              
              <div class="detail-section">
                <h4><i class="fas fa-calendar"></i> Timeline</h4>
                <div class="detail-item">
                  <strong>Created:</strong> ${new Date(
                    gym.created_at
                  ).toLocaleDateString()}
                </div>
                ${
                  gym.updated_at
                    ? `
                <div class="detail-item">
                  <strong>Last Updated:</strong> ${new Date(
                    gym.updated_at
                  ).toLocaleDateString()}
                </div>
                `
                    : ""
                }
                ${
                  gym.deleted_at
                    ? `
                <div class="detail-item">
                  <strong>Deleted:</strong> ${new Date(
                    gym.deleted_at
                  ).toLocaleDateString()}
                </div>
                `
                    : ""
                }
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" onclick="this.closest('.modal').remove()">Close</button>
            ${
              gym.deleted_at
                ? `
              <button type="button" class="btn btn-success" onclick="dashboard.restoreGym('${gym.id}')">
                <i class="fas fa-undo"></i> Restore Gym
              </button>
            `
                : `
              <button type="button" class="btn btn-warning" onclick="dashboard.deleteGym('${gym.id}')">
                <i class="fas fa-trash"></i> Delete Gym
              </button>
            `
            }
          </div>
        </div>
      `;

      document.body.appendChild(modal);
      modal.style.display = "block";
    } catch (error) {
      console.error("Failed to load gym details:", error);
      this.showErrorMessage("Failed to load gym details");
    }
  }

  manageGym(gymId) {
    // Show gym management modal/details for platform admin
    this.viewGymDetails(gymId);
  }

  async restoreGym(gymId) {
    if (!confirm("Are you sure you want to restore this gym?")) {
      return;
    }

    try {
      const response = await this.apiCall(`/gym/${gymId}/restore`, {
        method: "PUT",
      });

      if (response.ok) {
        // Close any open modals
        const modals = document.querySelectorAll(".modal");
        modals.forEach((modal) => modal.remove());

        this.showSuccessMessage("Gym restored successfully!");

        // Reload current view
        if (this.currentView === "gyms") {
          this.loadGymsManagement();
        } else {
          this.loadPlatformOverview();
        }
      } else {
        const error = await response.json();
        throw new Error(error.message || "Failed to restore gym");
      }
    } catch (error) {
      console.error("Restore gym error:", error);
      this.showErrorMessage(error.message || "Failed to restore gym");
    }
  }

  async deleteGym(gymId) {
    if (
      !confirm(
        "Are you sure you want to delete this gym? This action can be undone later."
      )
    ) {
      return;
    }

    try {
      const response = await this.apiCall(`/gym/${gymId}`, {
        method: "DELETE",
      });

      if (response.ok) {
        // Close any open modals
        const modals = document.querySelectorAll(".modal");
        modals.forEach((modal) => modal.remove());

        this.showSuccessMessage("Gym deleted successfully!");

        // Reload current view
        if (this.currentView === "gyms") {
          this.loadGymsManagement();
        } else {
          this.loadPlatformOverview();
        }
      } else {
        const error = await response.json();
        throw new Error(error.message || "Failed to delete gym");
      }
    } catch (error) {
      console.error("Delete gym error:", error);
      this.showErrorMessage(error.message || "Failed to delete gym");
    }
  }

  showCreateGymModal() {
    const modal = document.createElement("div");
    modal.className = "modal";
    modal.innerHTML = `
      <div class="modal-content">
        <div class="modal-header">
          <h2><i class="fas fa-plus"></i> Add New Gym</h2>
          <span class="close" onclick="this.closest('.modal').remove()">&times;</span>
        </div>
        <div class="modal-body">
          <form id="create-gym-form">
            <div class="form-group">
              <label for="gym-name">Gym Name *</label>
              <input type="text" id="gym-name" name="name" required class="form-control">
            </div>
            
            <div class="form-group">
              <label for="gym-contact-name">Contact Name</label>
              <input type="text" id="gym-contact-name" name="contact_name" class="form-control">
            </div>
            
            <div class="form-group">
              <label for="gym-contact-email">Contact Email</label>
              <input type="email" id="gym-contact-email" name="contact_email" class="form-control">
            </div>
            
            <div class="form-group">
              <label for="gym-contact-phone">Contact Phone</label>
              <input type="tel" id="gym-contact-phone" name="contact_phone" class="form-control">
            </div>
            
            <div class="form-group">
              <label for="gym-address">Address</label>
              <textarea id="gym-address" name="address" rows="3" class="form-control"></textarea>
            </div>
            
            <div class="form-group">
              <label for="gym-description">Description</label>
              <textarea id="gym-description" name="description" rows="4" class="form-control" placeholder="Brief description of the gym..."></textarea>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" onclick="this.closest('.modal').remove()">Cancel</button>
          <button type="button" class="btn btn-primary" onclick="dashboard.createGym()">Create Gym</button>
        </div>
      </div>
    `;

    document.body.appendChild(modal);
    modal.style.display = "block";

    // Focus on the first input
    setTimeout(() => {
      document.getElementById("gym-name").focus();
    }, 100);
  }

  async createGym() {
    const form = document.getElementById("create-gym-form");
    const formData = new FormData(form);
    const gymData = Object.fromEntries(formData);

    try {
      const response = await this.apiCall("/gym", {
        method: "POST",
        body: JSON.stringify(gymData),
      });

      if (response.ok) {
        // Close modal
        document.querySelector(".modal").remove();

        // Show success message
        this.showSuccessMessage("Gym created successfully!");

        // Reload the view if we're on the gyms page
        if (this.currentView === "gyms") {
          this.loadGymsManagement();
        } else {
          this.loadPlatformOverview();
        }
      } else {
        const error = await response.json();
        throw new Error(error.message || "Failed to create gym");
      }
    } catch (error) {
      console.error("Create gym error:", error);
      this.showErrorMessage(error.message || "Failed to create gym");
    }
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

  hideLoading() {
    // Hide loading state overlay
    const loadingState = document.getElementById("loading-state");
    if (loadingState) {
      loadingState.style.display = "none";
    }
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

  showSuccessMessage(message) {
    this.showToast(message, "success");
  }

  showErrorMessage(message) {
    this.showToast(message, "error");
  }

  showToast(message, type = "info") {
    // Remove any existing toast
    const existingToast = document.querySelector(".toast");
    if (existingToast) {
      existingToast.remove();
    }

    // Create toast element
    const toast = document.createElement("div");
    toast.className = `toast toast-${type}`;
    toast.innerHTML = `
      <div class="toast-content">
        <i class="fas fa-${
          type === "success"
            ? "check-circle"
            : type === "error"
            ? "exclamation-circle"
            : "info-circle"
        }"></i>
        <span>${message}</span>
      </div>
      <button class="toast-close" onclick="this.parentElement.remove()">
        <i class="fas fa-times"></i>
      </button>
    `;

    // Add to page
    document.body.appendChild(toast);

    // Show with animation
    setTimeout(() => toast.classList.add("show"), 100);

    // Auto-remove after 5 seconds
    setTimeout(() => {
      if (toast.parentNode) {
        toast.classList.remove("show");
        setTimeout(() => toast.remove(), 300);
      }
    }, 5000);
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

  // Platform Admin CRUD Functions (Placeholders)
  async loadEquipmentManagement() {
    try {
      const response = await this.api.getEquipment();
      const equipmentData = response?.data;
      const isDataNull = equipmentData === null;
      const equipmentArray = equipmentData || [];

      const content = `
        <div class="content-header">
          <h2>Equipment Management</h2>
          <p>Manage gym equipment and tools</p>
          <button class="btn btn-primary" onclick="dashboard.openCreateEquipmentModal()">
            <i class="fas fa-plus"></i> Add Equipment
          </button>
        </div>
        
        <div class="equipment-grid">
          ${
            equipmentArray.length > 0
              ? equipmentArray
                  .map(
                    (equipment) => `
            <div class="equipment-card" data-id="${equipment.id}">
              <div class="equipment-icon">
                ${this.getEquipmentIcon(equipment.category)}
              </div>
              <div class="equipment-info">
                <h3>${equipment.name}</h3>
                <p class="equipment-category">${equipment.category}</p>
                ${
                  equipment.description
                    ? `<p class="equipment-description">${equipment.description}</p>`
                    : ""
                }
              </div>
              <div class="equipment-actions">
                <button class="btn btn-outline btn-sm" onclick="dashboard.editEquipment(${
                  equipment.id
                })">
                  <i class="fas fa-edit"></i>
                </button>
                <button class="btn btn-danger btn-sm" onclick="dashboard.deleteEquipment(${
                  equipment.id
                })">
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </div>
          `
                  )
                  .join("")
              : isDataNull
              ? `
            <div class="empty-state">
              <i class="fas fa-exclamation-triangle warning-icon"></i>
              <h3>No Equipment Available</h3>
              <p>There is currently no equipment in the platform database.</p>
              <button class="btn-icon btn-refresh btn-sm" onclick="dashboard.loadEquipmentManagement()" title="Refresh">
                <i class="fas fa-sync-alt"></i>
              </button>
            </div>
          `
              : `
            <div class="empty-state">
              <i class="fas fa-dumbbell"></i>
              <h3>No Equipment Found</h3>
              <p>Your equipment catalog is empty. Use the "Add Equipment" button above to get started.</p>
            </div>
          `
          }
        </div>
        
        <!-- Create/Edit Modal -->
        <div id="equipmentModal" class="modal">
          <div class="modal-content">
            <div class="modal-header">
              <h3 class="modal-title" id="equipmentModalTitle">Add Equipment</h3>
              <button class="modal-close" onclick="dashboard.closeEquipmentModal()">&times;</button>
            </div>
            <div class="modal-body">
              <form id="equipmentForm">
                <div class="form-group">
                  <label for="equipmentName">Name *</label>
                  <input type="text" class="form-control" id="equipmentName" name="name" required>
                </div>
                <div class="form-group">
                  <label for="equipmentDescription">Description</label>
                  <textarea class="form-control" id="equipmentDescription" name="description" rows="3"></textarea>
                </div>
                <div class="form-group">
                  <label for="equipmentCategory">Category *</label>
                  <select class="form-control" id="equipmentCategory" name="category" required>
                    <option value="">Select category...</option>
                    <option value="cardio">Cardio</option>
                    <option value="free_weights">Free Weights</option>
                    <option value="machines">Machines</option>
                    <option value="accessories">Accessories</option>
                    <option value="bodyweight">Bodyweight</option>
                  </select>
                </div>
                <div class="form-actions">
                  <button type="button" class="btn btn-secondary" onclick="dashboard.closeEquipmentModal()">Cancel</button>
                  <button type="submit" class="btn btn-primary">Save</button>
                </div>
              </form>
            </div>
          </div>
        </div>
      `;

      this.setContent(content);
    } catch (error) {
      console.error("Error loading equipment:", error);
      this.hideLoading();
      this.showToast("Error loading equipment", "error");
    }
  }

  async loadExercisesManagement() {
    try {
      const [exercisesResponse, muscularGroupsResponse, equipmentResponse] =
        await Promise.all([
          this.api.getExercises(),
          this.api.getMuscularGroups(),
          this.api.getEquipment(),
        ]);

      // Validate API responses and provide fallbacks
      const exercises = exercisesResponse?.data;
      const muscularGroups = muscularGroupsResponse?.data || [];
      const equipment = equipmentResponse?.data || [];

      // Handle null data case (API returned null instead of empty array)
      const isExercisesDataNull = exercises === null;
      const exercisesArray = exercises || [];

      const content = `
        <div class="content-header">
          <h2>Exercise Management</h2>
          <p>Manage global exercise library</p>
          <div class="header-actions">
            ${
              muscularGroups.length === 0
                ? `
              <button class="btn btn-warning btn-sm" onclick="dashboard.openCreateMuscularGroupModal()">
                <i class="fas fa-plus"></i> Create Muscular Group
              </button>
            `
                : ""
            }
            ${
              equipment.length === 0
                ? `
              <button class="btn btn-info btn-sm" onclick="dashboard.openCreateEquipmentModal()">
                <i class="fas fa-plus"></i> Create Equipment
              </button>
            `
                : ""
            }
            <button class="btn btn-primary" onclick="dashboard.openCreateExerciseModal()">
              <i class="fas fa-plus"></i> Add Exercise
            </button>
          </div>
        </div>
        
        <div class="exercises-filters">
          <div class="filter-group">
            <label for="exerciseTypeFilter">Exercise Type:</label>
            <select id="exerciseTypeFilter" onchange="dashboard.filterExercises()">
              <option value="">All Types</option>
              <option value="strength">Strength</option>
              <option value="cardio">Cardio</option>
              <option value="flexibility">Flexibility</option>
              <option value="balance">Balance</option>
              <option value="functional">Functional</option>
            </select>
          </div>
          <div class="filter-group">
            <label for="difficultyFilter">Difficulty:</label>
            <select id="difficultyFilter" onchange="dashboard.filterExercises()">
              <option value="">All Levels</option>
              <option value="beginner">Beginner</option>
              <option value="intermediate">Intermediate</option>
              <option value="advanced">Advanced</option>
            </select>
          </div>
          <div class="filter-group">
            <label for="exerciseSearch">Search:</label>
            <input type="text" id="exerciseSearch" placeholder="Search exercises..." oninput="dashboard.filterExercises()">
          </div>
        </div>
        
        <div class="exercises-grid" id="exercisesGrid">
          ${
            exercisesArray.length > 0
              ? exercisesArray
                  .map(
                    (exercise) => `
            <div class="exercise-card" data-type="${
              exercise.exercise_type
            }" data-difficulty="${
                      exercise.difficulty_level
                    }" data-name="${exercise.name.toLowerCase()}">
              <div class="exercise-header">
                <div class="exercise-image">
                  ${
                    exercise.image_url
                      ? `<img src="${exercise.image_url}" alt="${exercise.name}" onerror="this.style.display='none'">`
                      : `<div class="exercise-placeholder">${this.getExerciseTypeIcon(
                          exercise.exercise_type
                        )}</div>`
                  }
                </div>
                <div class="exercise-info">
                  <h3>${exercise.name}</h3>
                  <div class="exercise-tags">
                    <span class="exercise-type-tag ${
                      exercise.exercise_type
                    }">${this.formatExerciseType(exercise.exercise_type)}</span>
                    <span class="difficulty-tag ${
                      exercise.difficulty_level
                    }">${this.formatDifficulty(
                      exercise.difficulty_level
                    )}</span>
                  </div>
                </div>
              </div>
              
              <div class="exercise-details">
                <div class="muscular-groups">
                  <strong>Muscles:</strong> ${exercise.muscular_groups.join(
                    ", "
                  )}
                </div>
                ${
                  exercise.equipment_needed &&
                  exercise.equipment_needed.length > 0
                    ? `<div class="equipment-needed">
                    <strong>Equipment:</strong> ${exercise.equipment_needed.join(
                      ", "
                    )}
                  </div>`
                    : ""
                }
                <div class="exercise-instructions">
                  <strong>Instructions:</strong> 
                  <p>${exercise.instructions.substring(0, 100)}${
                      exercise.instructions.length > 100 ? "..." : ""
                    }</p>
                </div>
                <div class="status-badge ${
                  exercise.is_active ? "active" : "inactive"
                }">
                  ${exercise.is_active ? "Active" : "Inactive"}
                </div>
              </div>
              
              <div class="exercise-actions">
                <button class="btn btn-sm btn-info" onclick="dashboard.viewExerciseDetails('${
                  exercise.id
                }')">
                  <i class="fas fa-eye"></i>
                </button>
                <button class="btn btn-sm btn-secondary" onclick="dashboard.editExercise('${
                  exercise.id
                }')">
                  <i class="fas fa-edit"></i>
                </button>
                <button class="btn btn-sm btn-danger" onclick="dashboard.deleteExercise('${
                  exercise.id
                }')">
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </div>
          `
                  )
                  .join("")
              : isExercisesDataNull
              ? `
            <div class="empty-state">
              <i class="fas fa-exclamation-triangle warning-icon"></i>
              <h3>No Exercises Available</h3>
              <p>There are currently no exercises in the platform database. This might be because the platform is newly set up or there was an issue loading exercise data.</p>
              <button class="btn-icon btn-refresh btn-sm" onclick="dashboard.loadExercisesManagement()" title="Refresh">
                <i class="fas fa-sync-alt"></i>
              </button>
            </div>
          `
              : `
            <div class="empty-state">
              <i class="fas fa-list-check"></i>
              <h3>No Exercises Found</h3>
              <p>Your exercise library is empty. Use the "Add Exercise" button above to get started.</p>
            </div>
          `
          }
        </div>
        
        <!-- Create/Edit Exercise Modal -->
        <div id="exerciseModal" class="modal">
          <div class="modal-content large-modal">
            <div class="modal-header">
              <h3 class="modal-title" id="exerciseModalTitle">Add Exercise</h3>
              <button class="modal-close" onclick="dashboard.closeExerciseModal()">&times;</button>
            </div>
            <div class="modal-body">
              <form id="exerciseForm">
              <div class="form-row">
                <div class="form-group">
                  <label for="exerciseName">Name *</label>
                  <input type="text" class="form-control" id="exerciseName" name="name" required>
                </div>
                <div class="form-group">
                  <label for="exerciseType">Exercise Type *</label>
                  <select class="form-control" id="exerciseType" name="exercise_type" required>
                    <option value="">Select type...</option>
                    <option value="strength">Strength</option>
                    <option value="cardio">Cardio</option>
                    <option value="flexibility">Flexibility</option>
                    <option value="balance">Balance</option>
                    <option value="functional">Functional</option>
                  </select>
                </div>
              </div>
              
              <div class="form-row">
                <div class="form-group">
                  <label for="difficultyLevel">Difficulty Level *</label>
                  <select class="form-control" id="difficultyLevel" name="difficulty_level" required>
                    <option value="">Select difficulty...</option>
                    <option value="beginner">Beginner</option>
                    <option value="intermediate">Intermediate</option>
                    <option value="advanced">Advanced</option>
                  </select>
                </div>
                <div class="form-group">
                  <label for="exerciseSynonyms">Synonyms (comma-separated)</label>
                  <input type="text" class="form-control" id="exerciseSynonyms" name="synonyms" placeholder="Alternative names...">
                </div>
              </div>
              
              <div class="form-group">
                <label for="muscularGroups">Muscular Groups *</label>
                <div class="checkbox-group" id="muscularGroupsCheckboxes">
                  ${
                    muscularGroups.length > 0
                      ? muscularGroups
                          .map(
                            (group) => `
                    <label class="checkbox-item">
                      <input type="checkbox" name="muscular_groups" value="${group.name}">
                      <span>${group.name}</span>
                    </label>
                  `
                          )
                          .join("")
                      : `
                    <div class="empty-dependency">
                      <i class="fas fa-exclamation-triangle"></i>
                      <p>No muscular groups available. <a href="#" onclick="dashboard.loadView('muscular-groups')">Create muscular groups</a> first.</p>
                    </div>
                  `
                  }
                </div>
              </div>
              
              <div class="form-group">
                <label for="equipmentNeeded">Equipment Needed</label>
                <div class="checkbox-group" id="equipmentCheckboxes">
                  ${
                    equipment.length > 0
                      ? equipment
                          .map(
                            (eq) => `
                    <label class="checkbox-item">
                      <input type="checkbox" name="equipment_needed" value="${eq.name}">
                      <span>${eq.name}</span>
                    </label>
                  `
                          )
                          .join("")
                      : `
                    <div class="empty-dependency">
                      <i class="fas fa-info-circle"></i>
                      <p>No equipment available. <a href="#" onclick="dashboard.loadView('equipment')">Create equipment</a> first, or leave empty for bodyweight exercises.</p>
                    </div>
                  `
                  }
                </div>
              </div>
              
              <div class="form-group">
                <label for="exerciseInstructions">Instructions *</label>
                <textarea class="form-control" id="exerciseInstructions" name="instructions" rows="4" required placeholder="Detailed exercise instructions..."></textarea>
              </div>
              
              <div class="form-row">
                <div class="form-group">
                  <label for="videoUrl">Video URL</label>
                  <input type="url" class="form-control" id="videoUrl" name="video_url" placeholder="https://...">
                </div>
                <div class="form-group">
                  <label for="imageUrl">Image URL</label>
                  <input type="url" class="form-control" id="imageUrl" name="image_url" placeholder="https://...">
                </div>
              </div>
              
              <div class="form-actions">
                <button type="button" class="btn btn-secondary" onclick="dashboard.closeExerciseModal()">Cancel</button>
                <button type="submit" class="btn btn-primary">Save Exercise</button>
              </div>
            </form>
            </div>
          </div>
        </div>
        
        <!-- Exercise Details Modal -->
        <div id="exerciseDetailsModal" class="modal">
          <div class="modal-content large-modal">
            <span class="close" onclick="dashboard.closeExerciseDetailsModal()">&times;</span>
            <div id="exerciseDetailsContent">
              <!-- Content will be loaded dynamically -->
            </div>
          </div>
        </div>
        
        <!-- Muscular Group Modal for dependency creation -->
        <div id="muscularGroupModal" class="modal">
          <div class="modal-content">
            <div class="modal-header">
              <h3 class="modal-title" id="modalTitle">Add Muscular Group</h3>
              <button class="modal-close" onclick="dashboard.closeMuscularGroupModal()">&times;</button>
            </div>
            <div class="modal-body">
              <form id="muscularGroupForm">
                <div class="form-group">
                  <label for="groupName">Name *</label>
                  <input type="text" class="form-control" id="groupName" name="name" required>
                </div>
                <div class="form-group">
                  <label for="groupDescription">Description</label>
                  <textarea class="form-control" id="groupDescription" name="description" rows="3"></textarea>
                </div>
                <div class="form-group">
                  <label for="bodyPart">Body Part *</label>
                  <select class="form-control" id="bodyPart" name="body_part" required>
                    <option value="">Select body part...</option>
                    <option value="upper_body">Upper Body</option>
                    <option value="lower_body">Lower Body</option>
                    <option value="core">Core</option>
                    <option value="full_body">Full Body</option>
                  </select>
                </div>
                <div class="form-actions">
                  <button type="button" class="btn btn-secondary" onclick="dashboard.closeMuscularGroupModal()">Cancel</button>
                  <button type="submit" class="btn btn-primary">Save</button>
                </div>
              </form>
            </div>
          </div>
        </div>
        
        <!-- Equipment Modal for dependency creation -->
        <div id="equipmentModal" class="modal">
          <div class="modal-content">
            <div class="modal-header">
              <h3 class="modal-title" id="equipmentModalTitle">Add Equipment</h3>
              <button class="modal-close" onclick="dashboard.closeEquipmentModal()">&times;</button>
            </div>
            <div class="modal-body">
              <form id="equipmentForm">
                <div class="form-group">
                  <label for="equipmentName">Name *</label>
                  <input type="text" class="form-control" id="equipmentName" name="name" required>
                </div>
                <div class="form-group">
                  <label for="equipmentDescription">Description</label>
                  <textarea class="form-control" id="equipmentDescription" name="description" rows="3"></textarea>
                </div>
                <div class="form-group">
                  <label for="equipmentCategory">Category *</label>
                  <select class="form-control" id="equipmentCategory" name="category" required>
                    <option value="">Select category...</option>
                    <option value="cardio">Cardio</option>
                    <option value="free_weights">Free Weights</option>
                    <option value="machines">Machines</option>
                    <option value="accessories">Accessories</option>
                    <option value="bodyweight">Bodyweight</option>
                  </select>
                </div>
                <div class="form-actions">
                  <button type="button" class="btn btn-secondary" onclick="dashboard.closeEquipmentModal()">Cancel</button>
                  <button type="submit" class="btn btn-primary">Save</button>
                </div>
              </form>
            </div>
          </div>
        </div>
      `;

      this.setContent(content);
    } catch (error) {
      console.error("Error loading exercises:", error);
      this.hideLoading();
      this.showToast("Failed to load exercises", "error");
    }
  }

  getExerciseTypeIcon(type) {
    const icons = {
      strength: '<i class="fas fa-dumbbell" style="color: #e74c3c;"></i>',
      cardio: '<i class="fas fa-heartbeat" style="color: #e91e63;"></i>',
      flexibility: '<i class="fas fa-hand-paper" style="color: #9c27b0;"></i>',
      balance: '<i class="fas fa-balance-scale" style="color: #673ab7;"></i>',
      functional: '<i class="fas fa-user-ninja" style="color: #3f51b5;"></i>',
    };
    return (
      icons[type] ||
      '<i class="fas fa-question-circle" style="color: #95a5a6;"></i>'
    );
  }

  formatExerciseType(type) {
    const formatted = {
      strength: "Strength",
      cardio: "Cardio",
      flexibility: "Flexibility",
      balance: "Balance",
      functional: "Functional",
    };
    return formatted[type] || type;
  }

  formatDifficulty(difficulty) {
    const formatted = {
      beginner: "Beginner",
      intermediate: "Intermediate",
      advanced: "Advanced",
    };
    return formatted[difficulty] || difficulty;
  }

  filterExercises() {
    const typeFilter = document.getElementById("exerciseTypeFilter").value;
    const difficultyFilter = document.getElementById("difficultyFilter").value;
    const searchFilter = document
      .getElementById("exerciseSearch")
      .value.toLowerCase();

    const exerciseCards = document.querySelectorAll(".exercise-card");

    exerciseCards.forEach((card) => {
      const type = card.dataset.type;
      const difficulty = card.dataset.difficulty;
      const name = card.dataset.name;

      const typeMatch = !typeFilter || type === typeFilter;
      const difficultyMatch =
        !difficultyFilter || difficulty === difficultyFilter;
      const searchMatch = !searchFilter || name.includes(searchFilter);

      if (typeMatch && difficultyMatch && searchMatch) {
        card.style.display = "block";
      } else {
        card.style.display = "none";
      }
    });
  }

  async openCreateExerciseModal() {
    try {
      // Fetch fresh data for muscular groups and equipment
      const [muscularGroupsResponse, equipmentResponse] = await Promise.all([
        this.api.getMuscularGroups(),
        this.api.getEquipment(),
      ]);

      const muscularGroups = muscularGroupsResponse?.data || [];
      const equipment = equipmentResponse?.data || [];

      const modal = document.getElementById("exerciseModal");
      const form = document.getElementById("exerciseForm");
      const title = document.getElementById("exerciseModalTitle");

      title.textContent = "Add Exercise";
      form.reset();
      form.removeAttribute("data-edit-id");

      // Update muscular groups checkboxes
      const muscularGroupsContainer = document.getElementById(
        "muscularGroupsCheckboxes"
      );
      muscularGroupsContainer.innerHTML =
        muscularGroups.length > 0
          ? muscularGroups
              .map(
                (group) => `
            <label class="checkbox-item">
              <input type="checkbox" name="muscular_groups" value="${group.name}">
              <span>${group.name}</span>
            </label>
          `
              )
              .join("")
          : `
          <div class="empty-dependency">
            <i class="fas fa-exclamation-triangle"></i>
            <p>No muscular groups available. <a href="#" onclick="dashboard.loadView('muscular-groups')">Create muscular groups</a> first.</p>
          </div>
        `;

      // Update equipment checkboxes
      const equipmentContainer = document.getElementById("equipmentCheckboxes");
      equipmentContainer.innerHTML =
        equipment.length > 0
          ? equipment
              .map(
                (eq) => `
            <label class="checkbox-item">
              <input type="checkbox" name="equipment_needed" value="${eq.name}">
              <span>${eq.name}</span>
            </label>
          `
              )
              .join("")
          : `
          <div class="empty-dependency">
            <i class="fas fa-info-circle"></i>
            <p>No equipment available. <a href="#" onclick="dashboard.loadView('equipment')">Create equipment</a> first, or leave empty for bodyweight exercises.</p>
          </div>
        `;

      modal.style.display = "block";

      form.onsubmit = async (e) => {
        e.preventDefault();
        await this.saveExercise();
      };
    } catch (error) {
      console.error("Error loading exercise modal:", error);
      this.showToast("Error loading exercise form", "error");
    }
  }

  async editExercise(id) {
    try {
      const response = await this.api.getExerciseById(id);
      const exercise = response.data;

      const modal = document.getElementById("exerciseModal");
      const form = document.getElementById("exerciseForm");
      const title = document.getElementById("exerciseModalTitle");

      title.textContent = "Edit Exercise";
      form.setAttribute("data-edit-id", id);

      // Fill form fields
      document.getElementById("exerciseName").value = exercise.name;
      document.getElementById("exerciseType").value = exercise.exercise_type;
      document.getElementById("difficultyLevel").value =
        exercise.difficulty_level;
      document.getElementById("exerciseSynonyms").value = exercise.synonyms
        ? exercise.synonyms.join(", ")
        : "";
      document.getElementById("exerciseInstructions").value =
        exercise.instructions;
      document.getElementById("videoUrl").value = exercise.video_url || "";
      document.getElementById("imageUrl").value = exercise.image_url || "";

      // Check muscular groups
      exercise.muscular_groups.forEach((group) => {
        const checkbox = document.querySelector(
          `input[name="muscular_groups"][value="${group}"]`
        );
        if (checkbox) checkbox.checked = true;
      });

      // Check equipment
      if (exercise.equipment_needed) {
        exercise.equipment_needed.forEach((eq) => {
          const checkbox = document.querySelector(
            `input[name="equipment_needed"][value="${eq}"]`
          );
          if (checkbox) checkbox.checked = true;
        });
      }

      modal.style.display = "block";

      form.onsubmit = async (e) => {
        e.preventDefault();
        await this.saveExercise();
      };
    } catch (error) {
      console.error("Error loading exercise:", error);
      this.showToast("Failed to load exercise details", "error");
    }
  }

  async saveExercise() {
    const form = document.getElementById("exerciseForm");
    const editId = form.getAttribute("data-edit-id");

    // Get selected muscular groups
    const muscularGroups = Array.from(
      document.querySelectorAll('input[name="muscular_groups"]:checked')
    ).map((cb) => cb.value);

    // Get selected equipment
    const equipment = Array.from(
      document.querySelectorAll('input[name="equipment_needed"]:checked')
    ).map((cb) => cb.value);

    // Parse synonyms
    const synonymsText = document.getElementById("exerciseSynonyms").value;
    const synonyms = synonymsText
      ? synonymsText
          .split(",")
          .map((s) => s.trim())
          .filter((s) => s)
      : [];

    const formData = {
      name: document.getElementById("exerciseName").value,
      exercise_type: document.getElementById("exerciseType").value,
      difficulty_level: document.getElementById("difficultyLevel").value,
      synonyms: synonyms,
      muscular_groups: muscularGroups,
      equipment: equipment,
      instructions: document.getElementById("exerciseInstructions").value,
      video_url: document.getElementById("videoUrl").value || null,
      image_url: document.getElementById("imageUrl").value || null,
      created_by: this.currentUser.id,
    };

    try {
      if (editId) {
        await this.api.updateExercise(editId, formData);
        this.showToast("Exercise updated successfully", "success");
      } else {
        await this.api.createExercise(formData);
        this.showToast("Exercise created successfully", "success");
      }

      this.closeExerciseModal();
      await this.loadExercisesManagement();
    } catch (error) {
      console.error("Error saving exercise:", error);
      this.showToast("Failed to save exercise", "error");
    }
  }

  async deleteExercise(id) {
    if (!confirm("Are you sure you want to delete this exercise?")) {
      return;
    }

    try {
      await this.api.deleteExercise(id);
      this.showToast("Exercise deleted successfully", "success");
      await this.loadExercisesManagement();
    } catch (error) {
      console.error("Error deleting exercise:", error);
      this.showToast("Failed to delete exercise", "error");
    }
  }

  async viewExerciseDetails(id) {
    try {
      const response = await this.api.getExerciseById(id);
      const exercise = response.data;

      const modal = document.getElementById("exerciseDetailsModal");
      const content = document.getElementById("exerciseDetailsContent");

      content.innerHTML = `
        <div class="exercise-details-view">
          <div class="exercise-header-details">
            <div class="exercise-image-large">
              ${
                exercise.image_url
                  ? `<img src="${exercise.image_url}" alt="${exercise.name}">`
                  : `<div class="exercise-placeholder-large">${this.getExerciseTypeIcon(
                      exercise.exercise_type
                    )}</div>`
              }
            </div>
            <div class="exercise-meta">
              <h2>${exercise.name}</h2>
              <div class="exercise-tags-large">
                <span class="exercise-type-tag ${
                  exercise.exercise_type
                }">${this.formatExerciseType(exercise.exercise_type)}</span>
                <span class="difficulty-tag ${
                  exercise.difficulty_level
                }">${this.formatDifficulty(exercise.difficulty_level)}</span>
                <span class="status-badge ${
                  exercise.is_active ? "active" : "inactive"
                }">${exercise.is_active ? "Active" : "Inactive"}</span>
              </div>
              ${
                exercise.synonyms && exercise.synonyms.length > 0
                  ? `<div class="synonyms"><strong>Also known as:</strong> ${exercise.synonyms.join(
                      ", "
                    )}</div>`
                  : ""
              }
            </div>
          </div>
          
          <div class="exercise-details-grid">
            <div class="detail-section">
              <h4><i class="fas fa-user-doctor"></i> Muscular Groups</h4>
              <div class="tags-list">
                ${exercise.muscular_groups
                  .map((group) => `<span class="tag">${group}</span>`)
                  .join("")}
              </div>
            </div>
            
            ${
              exercise.equipment_needed && exercise.equipment_needed.length > 0
                ? `
              <div class="detail-section">
                <h4><i class="fas fa-tools"></i> Equipment Needed</h4>
                <div class="tags-list">
                  ${exercise.equipment_needed
                    .map((eq) => `<span class="tag">${eq}</span>`)
                    .join("")}
                </div>
              </div>
            `
                : ""
            }
            
            <div class="detail-section full-width">
              <h4><i class="fas fa-list-ol"></i> Instructions</h4>
              <p class="instructions-text">${exercise.instructions}</p>
            </div>
            
            ${
              exercise.video_url
                ? `
              <div class="detail-section full-width">
                <h4><i class="fas fa-video"></i> Video</h4>
                <a href="${exercise.video_url}" target="_blank" class="video-link">
                  <i class="fas fa-external-link-alt"></i> Watch Video
                </a>
              </div>
            `
                : ""
            }
          </div>
          
          <div class="exercise-actions-details">
            <button class="btn btn-secondary" onclick="dashboard.editExercise('${
              exercise.id
            }'); dashboard.closeExerciseDetailsModal();">
              <i class="fas fa-edit"></i> Edit Exercise
            </button>
            <button class="btn btn-danger" onclick="dashboard.deleteExercise('${
              exercise.id
            }'); dashboard.closeExerciseDetailsModal();">
              <i class="fas fa-trash"></i> Delete Exercise
            </button>
          </div>
        </div>
      `;

      modal.style.display = "block";
    } catch (error) {
      console.error("Error loading exercise details:", error);
      this.showToast("Failed to load exercise details", "error");
    }
  }

  closeExerciseModal() {
    const modal = document.getElementById("exerciseModal");
    modal.style.display = "none";
  }

  closeExerciseDetailsModal() {
    const modal = document.getElementById("exerciseDetailsModal");
    modal.style.display = "none";
  }

  async loadMuscularGroupsManagement() {
    try {
      const muscularGroups = await this.api.getMuscularGroups();

      // Validate API response and provide fallback
      const groupsData = muscularGroups?.data;
      const isDataNull = groupsData === null;
      const groupsArray = groupsData || [];

      const content = `
        <div class="content-header">
          <h2>Muscular Groups Management</h2>
          <p>Manage muscle group categories for exercises</p>
          <button class="btn btn-primary" onclick="dashboard.openCreateMuscularGroupModal()">
            <i class="fas fa-plus"></i> Add Muscular Group
          </button>
        </div>
        
        <div class="muscular-groups-grid">
          ${
            groupsArray.length > 0
              ? groupsArray
                  .map(
                    (group) => `
            <div class="muscular-group-card" data-id="${group.id}">
              <div class="muscular-group-icon">
                ${this.getMuscularGroupIcon(group.body_part)}
              </div>
              <div class="muscular-group-info">
                <h3>${group.name}</h3>
                <p class="body-part-tag ${
                  group.body_part
                }">${this.formatBodyPart(group.body_part)}</p>
                ${
                  group.description
                    ? `<p class="description">${group.description}</p>`
                    : ""
                }
                <div class="status-badge ${
                  group.is_active ? "active" : "inactive"
                }">
                  ${group.is_active ? "Active" : "Inactive"}
                </div>
              </div>
              <div class="muscular-group-actions">
                <button class="btn btn-sm btn-secondary" onclick="dashboard.editMuscularGroup('${
                  group.id
                }')">
                  <i class="fas fa-edit"></i>
                </button>
                <button class="btn btn-sm btn-danger" onclick="dashboard.deleteMuscularGroup('${
                  group.id
                }')">
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </div>
          `
                  )
                  .join("")
              : isDataNull
              ? `
            <div class="empty-state">
              <i class="fas fa-exclamation-triangle" style="color: #f59e0b;"></i>
              <h3>No Muscular Groups Available</h3>
              <p>There are currently no muscular groups in the platform database.</p>
              <div style="display: flex; gap: 10px; justify-content: center; flex-wrap: wrap;">
                <button class="btn-icon btn-refresh btn-sm" onclick="dashboard.loadMuscularGroupsManagement()" title="Refresh">
                  <i class="fas fa-sync-alt"></i>
                </button>
              </div>
            </div>
          `
              : `
            <div class="empty-state">
              <i class="fas fa-user-doctor"></i>
              <h3>No Muscular Groups Found</h3>
              <p>Your muscular groups library is empty. Use the "Add Muscular Group" button above to get started.</p>
            </div>
          `
          }
        </div>
        
        <!-- Create/Edit Modal -->
        <div id="muscularGroupModal" class="modal">
          <div class="modal-content">
            <div class="modal-header">
              <h3 class="modal-title" id="modalTitle">Add Muscular Group</h3>
              <button class="modal-close" onclick="dashboard.closeMuscularGroupModal()">&times;</button>
            </div>
            <div class="modal-body">
              <form id="muscularGroupForm">
                <div class="form-group">
                  <label for="groupName">Name *</label>
                  <input type="text" class="form-control" id="groupName" name="name" required>
                </div>
                <div class="form-group">
                  <label for="groupDescription">Description</label>
                  <textarea class="form-control" id="groupDescription" name="description" rows="3"></textarea>
                </div>
                <div class="form-group">
                  <label for="bodyPart">Body Part *</label>
                  <select class="form-control" id="bodyPart" name="body_part" required>
                    <option value="">Select body part...</option>
                    <option value="upper_body">Upper Body</option>
                    <option value="lower_body">Lower Body</option>
                    <option value="core">Core</option>
                    <option value="full_body">Full Body</option>
                  </select>
                </div>
                <div class="form-actions">
                  <button type="button" class="btn btn-secondary" onclick="dashboard.closeMuscularGroupModal()">Cancel</button>
                  <button type="submit" class="btn btn-primary">Save</button>
                </div>
              </form>
            </div>
          </div>
        </div>
      `;

      this.setContent(content);
    } catch (error) {
      console.error("Error loading muscular groups:", error);
      this.hideLoading();
      this.showToast("Failed to load muscular groups", "error");
    }
  }

  getMuscularGroupIcon(bodyPart) {
    const icons = {
      upper_body: '<i class="fas fa-arrow-up" style="color: #3498db;"></i>',
      lower_body: '<i class="fas fa-running" style="color: #e74c3c;"></i>',
      core: '<i class="fas fa-circle-dot" style="color: #f39c12;"></i>',
      full_body: '<i class="fas fa-user" style="color: #9b59b6;"></i>',
    };
    return (
      icons[bodyPart] ||
      '<i class="fas fa-question-circle" style="color: #95a5a6;"></i>'
    );
  }

  formatBodyPart(bodyPart) {
    const formatted = {
      upper_body: "Upper Body",
      lower_body: "Lower Body",
      core: "Core",
      full_body: "Full Body",
    };
    return formatted[bodyPart] || bodyPart;
  }

  getEquipmentIcon(category) {
    const icons = {
      cardio: '<i class="fas fa-heartbeat" style="color: #e74c3c;"></i>',
      free_weights:
        '<i class="fas fa-weight-hanging" style="color: #2c3e50;"></i>',
      machines: '<i class="fas fa-cogs" style="color: #7f8c8d;"></i>',
      accessories: '<i class="fas fa-tools" style="color: #f39c12;"></i>',
      bodyweight: '<i class="fas fa-male" style="color: #27ae60;"></i>',
    };
    return (
      icons[category] ||
      '<i class="fas fa-question-circle" style="color: #95a5a6;"></i>'
    );
  }

  openCreateMuscularGroupModal() {
    const modal = document.getElementById("muscularGroupModal");
    const form = document.getElementById("muscularGroupForm");
    const title = document.getElementById("modalTitle");

    title.textContent = "Add Muscular Group";
    form.reset();
    form.removeAttribute("data-edit-id");
    modal.style.display = "block";

    form.onsubmit = async (e) => {
      e.preventDefault();
      await this.saveMuscularGroup();
    };
  }

  async editMuscularGroup(id) {
    try {
      const response = await this.api.getMuscularGroupById(id);
      const group = response.data;

      const modal = document.getElementById("muscularGroupModal");
      const form = document.getElementById("muscularGroupForm");
      const title = document.getElementById("modalTitle");

      title.textContent = "Edit Muscular Group";
      form.setAttribute("data-edit-id", id);

      document.getElementById("groupName").value = group.name;
      document.getElementById("groupDescription").value =
        group.description || "";
      document.getElementById("bodyPart").value = group.body_part;

      modal.style.display = "block";

      form.onsubmit = async (e) => {
        e.preventDefault();
        await this.saveMuscularGroup();
      };
    } catch (error) {
      console.error("Error loading muscular group:", error);
      this.showToast("Failed to load muscular group details", "error");
    }
  }

  async saveMuscularGroup() {
    const form = document.getElementById("muscularGroupForm");
    const editId = form.getAttribute("data-edit-id");

    const formData = {
      name: document.getElementById("groupName").value,
      description: document.getElementById("groupDescription").value,
      body_part: document.getElementById("bodyPart").value,
    };

    try {
      if (editId) {
        await this.api.updateMuscularGroup(editId, formData);
        this.showToast("Muscular group updated successfully", "success");
      } else {
        await this.api.createMuscularGroup(formData);
        this.showToast("Muscular group created successfully", "success");
      }

      this.closeMuscularGroupModal();
      await this.loadMuscularGroupsManagement();
    } catch (error) {
      console.error("Error saving muscular group:", error);
      this.showToast("Failed to save muscular group", "error");
    }
  }

  async deleteMuscularGroup(id) {
    if (!confirm("Are you sure you want to delete this muscular group?")) {
      return;
    }

    try {
      await this.api.deleteMuscularGroup(id);
      this.showToast("Muscular group deleted successfully", "success");
      await this.loadMuscularGroupsManagement();
    } catch (error) {
      console.error("Error deleting muscular group:", error);
      this.showToast("Failed to delete muscular group", "error");
    }
  }

  closeMuscularGroupModal() {
    const modal = document.getElementById("muscularGroupModal");
    modal.style.display = "none";
  }

  async loadWorkoutTemplatesManagement() {
    this.showComingSoon("Workout Templates Management");
  }

  async loadPlatformAnalytics() {
    this.showComingSoon("Platform Analytics");
  }

  async loadPlatformSettings() {
    this.showComingSoon("Platform Settings");
  }

  // Refresh current view
  async refreshCurrentView() {
    // Show loading state
    this.showLoading();

    try {
      // Reload the current view
      switch (this.currentView) {
        case "platform-overview":
          await this.loadPlatformOverview();
          break;
        case "gyms":
          await this.loadGymsManagement();
          break;
        case "exercises":
          await this.loadExercisesManagement();
          break;
        case "equipment":
          await this.loadEquipmentManagement();
          break;
        case "muscular-groups":
          await this.loadMuscularGroupsManagement();
          break;
        case "workout-templates":
          await this.loadWorkoutTemplates();
          break;
        case "platform-analytics":
          await this.loadPlatformAnalytics();
          break;
        case "platform-settings":
          await this.loadPlatformSettings();
          break;
        default:
          await this.loadPlatformOverview();
      }

      // Hide loading state after successful refresh
      this.hideLoading();
    } catch (error) {
      console.error("Failed to refresh view:", error);
      this.hideLoading();
      this.showError("Failed to refresh content");
    }
  }

  // Open help modal
  openHelp() {
    const modal = document.createElement("div");
    modal.className = "modal";
    modal.innerHTML = `
      <div class="modal-content">
        <div class="modal-header">
          <h2><i class="fas fa-question-circle"></i> Help & Support</h2>
          <span class="close" onclick="this.closest('.modal').remove()">&times;</span>
        </div>
        <div class="modal-body">
          <div class="help-sections">
            <div class="help-section">
              <h4><i class="fas fa-tachometer-alt"></i> Platform Overview</h4>
              <p>View comprehensive statistics about your gym network, including active gyms, member counts, and recent activity.</p>
            </div>
            
            <div class="help-section">
              <h4><i class="fas fa-building"></i> Gym Management</h4>
              <p>Create, edit, and manage gyms in your network. Use search and filters to find specific gyms quickly.</p>
            </div>
            
            <div class="help-section">
              <h4><i class="fas fa-plus"></i> Quick Actions</h4>
              <p>Use quick action buttons to create new gyms, generate reports, and access system settings.</p>
            </div>
            
            <div class="help-section">
              <h4><i class="fas fa-sync-alt"></i> Refresh Data</h4>
              <p>Click the refresh button to reload the current view with the latest data from the server.</p>
            </div>
            
            <div class="help-section">
              <h4><i class="fas fa-phone"></i> Contact Support</h4>
              <p>For technical support or questions about platform features, contact our support team.</p>
              <div class="contact-info">
                <p><strong>Email:</strong> support@athenai.com</p>
                <p><strong>Documentation:</strong> <a href="#" target="_blank">Help Center</a></p>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" onclick="this.closest('.modal').remove()">Close</button>
          <button type="button" class="btn btn-primary" onclick="window.open('mailto:support@athenai.com', '_blank')">
            <i class="fas fa-envelope"></i> Contact Support
          </button>
        </div>
      </div>
    `;

    document.body.appendChild(modal);
    modal.style.display = "block";
  }

  // Equipment Modal Methods
  openCreateEquipmentModal() {
    const modal = document.getElementById("equipmentModal");
    const form = document.getElementById("equipmentForm");
    const title = document.getElementById("equipmentModalTitle");

    title.textContent = "Add Equipment";
    form.reset();
    form.removeAttribute("data-edit-id");
    modal.style.display = "block";

    form.onsubmit = async (e) => {
      e.preventDefault();
      await this.saveEquipment();
    };
  }

  async editEquipment(id) {
    try {
      const response = await this.api.getEquipmentById(id);
      const equipment = response.data;

      const modal = document.getElementById("equipmentModal");
      const form = document.getElementById("equipmentForm");
      const title = document.getElementById("equipmentModalTitle");

      title.textContent = "Edit Equipment";
      form.setAttribute("data-edit-id", id);

      document.getElementById("equipmentName").value = equipment.name;
      document.getElementById("equipmentDescription").value =
        equipment.description || "";
      document.getElementById("equipmentCategory").value = equipment.category;

      modal.style.display = "block";

      form.onsubmit = async (e) => {
        e.preventDefault();
        await this.saveEquipment();
      };
    } catch (error) {
      console.error("Error loading equipment:", error);
      this.showToast("Error loading equipment", "error");
    }
  }

  async saveEquipment() {
    const form = document.getElementById("equipmentForm");
    const editId = form.getAttribute("data-edit-id");

    // Debug: Check if form elements exist
    const nameElement = document.getElementById("equipmentName");
    const descElement = document.getElementById("equipmentDescription");
    const categoryElement = document.getElementById("equipmentCategory");

    console.log("Form elements:", {
      nameElement: nameElement?.value,
      descElement: descElement?.value,
      categoryElement: categoryElement?.value,
    });

    const equipmentData = {
      name: document.getElementById("equipmentName").value,
      description:
        document.getElementById("equipmentDescription").value || null,
      category: document.getElementById("equipmentCategory").value,
    };

    try {
      if (editId) {
        await this.api.updateEquipment(editId, equipmentData);
        this.showToast("Equipment updated successfully!", "success");
      } else {
        await this.api.createEquipment(equipmentData);
        this.showToast("Equipment created successfully!", "success");
      }

      this.closeEquipmentModal();
      await this.loadEquipmentManagement();
    } catch (error) {
      console.error("Error saving equipment:", error);
      console.error("Equipment data that failed:", equipmentData);
      this.showToast("Error saving equipment: " + error.message, "error");
    }
  }

  async deleteEquipment(id) {
    if (!confirm("Are you sure you want to delete this equipment?")) {
      return;
    }

    try {
      await this.api.deleteEquipment(id);
      this.showToast("Equipment deleted successfully!", "success");
      await this.loadEquipmentManagement();
    } catch (error) {
      console.error("Error deleting equipment:", error);
      this.showToast("Error deleting equipment", "error");
    }
  }

  closeEquipmentModal() {
    const modal = document.getElementById("equipmentModal");
    modal.style.display = "none";
  }
}

// Global function for logout
function logout() {
  localStorage.removeItem("auth_token");
  window.location.href = "/";
}

// Global function for refresh button
function refreshCurrentView() {
  if (dashboard) {
    dashboard.refreshCurrentView();
  }
}

// Global function for help button
function openHelp() {
  if (dashboard) {
    dashboard.openHelp();
  }
}

// Initialize dashboard when DOM is loaded
let dashboard;
document.addEventListener("DOMContentLoaded", function () {
  dashboard = new DashboardManager();
});
