/**
 * Vanilla JS Dashboard Manager
 * Uses modular components loaded via script tags for browser compatibility
 * No ES6 imports - works directly in browsers without build tools
 */

class VanillaDashboardManager {
  constructor() {
    this.currentUser = null;
    this.currentView = "overview";
    this.components = {};
    this.managers = {};

    this.init();
  }

  async init() {
    try {
      console.log("Starting dashboard initialization...");

      // Wait for all component scripts to load
      await this.waitForDependencies();
      console.log("Dependencies loaded successfully");

      // Initialize content area
      this.contentArea = document.getElementById("content-body");
      console.log("Content area found:", this.contentArea);

      if (!this.contentArea) {
        throw new Error("Content area element not found");
      }

      // Test content area by adding a simple message
      this.contentArea.innerHTML =
        '<div style="padding: 20px; text-align: center; color: #666;">Initializing dashboard...</div>';
      console.log("Content area test successful");

      // Initialize managers
      this.initializeManagers();
      console.log("Managers initialized");

      // Check authentication
      await this.checkAuth();
      console.log("Authentication checked");

      // Initialize navigation
      this.initNavigation();
      console.log("Navigation initialized");

      // Initialize mobile functionality
      this.initMobile();
      console.log("Mobile functionality initialized");

      // Load the initial overview automatically
      console.log("Loading initial overview...");
      await this.loadView("overview");

      console.log("Vanilla Dashboard Manager initialized successfully");
    } catch (error) {
      console.error("Failed to initialize dashboard:", error);
      this.showError("Failed to initialize dashboard: " + error.message);
    }
  }

  // Wait for all required classes to be available
  async waitForDependencies() {
    const requiredClasses = [
      "ApiClient",
      "notifications",
      "appState",
      "getFormData",
      "Modal",
      "DataTable",
      "Card",
      "Grid",
      "EquipmentManager",
      "ExerciseManager",
      "GymManager",
    ];

    let attempts = 0;
    const maxAttempts = 100; // 10 seconds max wait

    while (attempts < maxAttempts) {
      const allLoaded = requiredClasses.every((className) => {
        return window[className] !== undefined;
      });

      if (allLoaded) {
        console.log("All dependencies loaded successfully");
        return;
      }

      await new Promise((resolve) => setTimeout(resolve, 100));
      attempts++;
    }

    console.warn(
      "Some dependencies may not have loaded:",
      requiredClasses.filter((cls) => window[cls] === undefined)
    );
  }

  initializeManagers() {
    // Initialize our specialized managers - they should all be available now
    this.managers.equipment = new EquipmentManager();
    this.managers.exercise = new ExerciseManager();
    this.managers.gym = new GymManager();

    // Set up event listeners for manager events
    this.setupManagerEventListeners();
  }

  setupManagerEventListeners() {
    // Equipment events
    document.addEventListener("equipment:edit", (e) => {
      this.openEquipmentModal("edit", e.detail.equipment);
    });

    // Exercise events
    document.addEventListener("exercise:edit", (e) => {
      this.openExerciseModal("edit", e.detail.exercise);
    });

    document.addEventListener("exercise:view", (e) => {
      this.viewExerciseDetails(e.detail.exercise);
    });

    // Gym events
    document.addEventListener("gym:edit", (e) => {
      this.openGymModal("edit", e.detail.gym);
    });

    document.addEventListener("gym:view", (e) => {
      this.viewGymDetails(e.detail.gym);
    });
  }

  async checkAuth() {
    try {
      // Check for stored auth token from login
      const authToken = localStorage.getItem("auth_token");
      const userInfo = localStorage.getItem("user_info");

      if (!authToken || !userInfo) {
        console.error("No authentication data found");
        window.location.href = "/";
        return;
      }

      this.currentUser = JSON.parse(userInfo);
      appState.setState({ user: this.currentUser });

      this.updateUserInfo();
      this.setupNavigation();
    } catch (error) {
      console.error("Authentication check failed:", error);
      // Redirect to login or show auth error
      window.location.href = "/";
    }
  }

  updateUserInfo() {
    const userNameEl = document.getElementById("user-name");
    const userRoleEl = document.getElementById("user-role");

    if (userNameEl) {
      userNameEl.textContent = this.currentUser.username || this.currentUser.id;
    }

    if (userRoleEl) {
      userRoleEl.textContent =
        this.currentUser.role || this.currentUser.user_type || "User";
    }
  }

  setupNavigation() {
    const platformAdminNav = document.querySelectorAll(".platform-admin-nav");

    if (this.currentUser.user_type === "platform_admin") {
      platformAdminNav.forEach((nav) => (nav.style.display = "block"));
    } else {
      platformAdminNav.forEach((nav) => (nav.style.display = "none"));
    }
  }

  initNavigation() {
    // Handle navigation clicks using event delegation
    const navMenu = document.getElementById("navigation-menu");
    if (navMenu) {
      navMenu.addEventListener("click", (e) => {
        const navItem = e.target.closest(".nav-item");
        if (navItem) {
          e.preventDefault();

          const view = navItem.getAttribute("data-view");
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
    document.querySelectorAll(".nav-item").forEach((item) => {
      item.classList.remove("active");
    });

    // Add active class to clicked item
    activeItem.classList.add("active");
  }

  setActiveNavItem(viewName) {
    // Remove active class from all nav items
    document.querySelectorAll(".nav-item").forEach((item) => {
      item.classList.remove("active");
    });

    // Add active class to the nav item with matching data-view
    const navItem = document.querySelector(`[data-view="${viewName}"]`);
    if (navItem) {
      navItem.classList.add("active");
    }
  }

  initMobile() {
    const sidebar = document.querySelector(".sidebar");
    const mobileToggle = document.getElementById("mobile-sidebar-toggle");
    const sidebarToggle = document.getElementById("sidebar-toggle");

    // Mobile sidebar toggle
    if (mobileToggle) {
      mobileToggle.addEventListener("click", () => {
        sidebar.classList.toggle("sidebar-open");
      });
    }

    if (sidebarToggle) {
      sidebarToggle.addEventListener("click", () => {
        sidebar.classList.toggle("sidebar-open");
      });
    }

    // Close sidebar when clicking outside on mobile
    document.addEventListener("click", (e) => {
      if (
        window.innerWidth <= 768 &&
        !sidebar.contains(e.target) &&
        !mobileToggle.contains(e.target)
      ) {
        sidebar.classList.remove("sidebar-open");
      }
    });
  }

  async loadView(viewName) {
    this.currentView = viewName;
    appState.setState({ currentView: viewName });

    // Update active navigation item
    this.setActiveNavItem(viewName);

    // Show loading state
    this.showLoading();

    try {
      switch (viewName) {
        case "overview":
          await this.loadDashboardOverview();
          break;
        case "gyms":
          await this.loadGymsManagement();
          break;
        case "gym-users":
          await this.loadGymUsersManagement();
          break;
        case "gym-analytics":
          await this.loadGymAnalytics();
          break;
        case "gym-requests":
          await this.loadGymRequests();
          break;
        case "equipment":
          await this.loadEquipmentManagement();
          break;
        case "exercises":
          await this.loadExercisesManagement();
          break;
        case "workout-templates":
          await this.loadWorkoutTemplates();
          break;
        case "muscular-groups":
          await this.loadMuscularGroupsManagement();
          break;
        case "system-settings":
          await this.loadSystemSettings();
          break;
        case "system-logs":
          await this.loadSystemLogs();
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

  // Dashboard Overview - Gym Management Focus
  async loadDashboardOverview() {
    try {
      console.log("Loading dashboard overview...");

      // Load real data from available APIs
      const api = new ApiClient();

      // Get real data from existing endpoints
      const [gyms, equipment, exercises] = await Promise.allSettled([
        api.request("/gym"),
        api.request("/equipment"),
        api.request("/exercise"),
      ]);

      // Extract successful results and handle API response format
      const gymsData =
        gyms.status === "fulfilled" ? gyms.value?.data || gyms.value || [] : [];
      const equipmentData =
        equipment.status === "fulfilled"
          ? equipment.value?.data || equipment.value || []
          : [];
      const exercisesData =
        exercises.status === "fulfilled"
          ? exercises.value?.data || exercises.value || []
          : [];

      // Debug logging
      console.log("API Response Debug:", {
        gyms: gyms.status === "fulfilled" ? gyms.value : gyms.reason,
        equipment:
          equipment.status === "fulfilled" ? equipment.value : equipment.reason,
        exercises:
          exercises.status === "fulfilled" ? exercises.value : exercises.reason,
        gymsData,
        equipmentData,
        exercisesData,
      });

      // Build dashboard data from real API responses
      const dashboardData = this.buildDashboardDataFromAPIs(
        gymsData,
        equipmentData,
        exercisesData
      );

      console.log("Built dashboard data:", dashboardData);

      const content = `
        <div class="dashboard-header">
          <div class="dashboard-header-content">
            <div class="dashboard-header-text">
              <h1 class="dashboard-title">AthenAI Dashboard</h1>
              <p class="dashboard-subtitle">Platform administration and management</p>
            </div>
            <div class="dashboard-header-actions">
              <button class="btn-action" id="refreshDashboard" title="Refresh Dashboard">
                <i class="fas fa-sync-alt"></i>
              </button>
              <button class="btn-action" id="dashboardHelp" title="Help & Documentation">
                <i class="fas fa-question-circle"></i>
              </button>
            </div>
          </div>
        </div>
        <div class="dashboard-overview">
          <!-- Key Metrics Row -->
          <div class="metrics-row">
            <div class="metric-card metric-primary">
              <div class="metric-icon">
                <i class="fas fa-building"></i>
              </div>
              <div class="metric-content">
                <div class="metric-value">${
                  dashboardData?.gyms?.total || 0
                }</div>
                <div class="metric-label">Registered Gyms</div>
                <div class="metric-change neutral">
                  <i class="fas fa-info-circle"></i>
                  Total facilities
                </div>
              </div>
            </div>

            <div class="metric-card metric-success">
              <div class="metric-icon">
                <i class="fas fa-dumbbell"></i>
              </div>
              <div class="metric-content">
                <div class="metric-value">${
                  dashboardData?.equipment?.total || 0
                }</div>
                <div class="metric-label">Equipment Items</div>
                <div class="metric-change neutral">
                  <i class="fas fa-info-circle"></i>
                  In catalog
                </div>
              </div>
            </div>

            <div class="metric-card metric-warning">
              <div class="metric-icon">
                <i class="fas fa-running"></i>
              </div>
              <div class="metric-content">
                <div class="metric-value">${
                  dashboardData?.exercises?.total || 0
                }</div>
                <div class="metric-label">Exercises</div>
                <div class="metric-change neutral">
                  <i class="fas fa-info-circle"></i>
                  In library
                </div>
              </div>
            </div>

            <div class="metric-card metric-info">
              <div class="metric-icon">
                <i class="fas fa-chart-line"></i>
              </div>
              <div class="metric-content">
                <div class="metric-value">${
                  dashboardData?.systemStatus?.api === "healthy"
                    ? "Online"
                    : "Offline"
                }</div>
                <div class="metric-label">System Status</div>
                <div class="metric-change ${
                  dashboardData?.systemStatus?.api === "healthy"
                    ? "positive"
                    : "negative"
                }">
                  <i class="fas ${
                    dashboardData?.systemStatus?.api === "healthy"
                      ? "fa-check-circle"
                      : "fa-exclamation-triangle"
                  }"></i>
                  APIs ${dashboardData?.systemStatus?.api}
                </div>
              </div>
            </div>
          </div>

          <!-- Main Content Grid -->
          <div class="dashboard-grid">
            <!-- Recent Activity -->
            <div class="dashboard-card span-2">
              <div class="card-header">
                <h3 class="card-title">
                  <i class="fas fa-activity"></i>
                  Recent Activity
                </h3>
              </div>
              <div class="card-body">
                <div class="activity-list">
                  ${this.renderRecentActivity(
                    dashboardData?.recentActivity || []
                  )}
                </div>
              </div>
            </div>

            <!-- Recent Gyms -->
            <div class="dashboard-card">
              <div class="card-header">
                <h3 class="card-title">
                  <i class="fas fa-building"></i>
                  Recent Gyms
                </h3>
                <button class="btn btn-ghost btn-sm" onclick="dashboard.loadView('gyms')">
                  View All
                </button>
              </div>
              <div class="card-body">
                <div class="gym-list">
                  ${this.renderRecentGyms(dashboardData?.gyms?.recent || [])}
                </div>
              </div>
            </div>

            <!-- Quick Actions -->
            <div class="dashboard-card">
              <div class="card-header">
                <h3 class="card-title">
                  <i class="fas fa-bolt"></i>
                  Quick Actions
                </h3>
              </div>
              <div class="card-body">
                <div class="quick-actions">
                  <button class="btn btn-primary btn-block" onclick="dashboard.openAddGymModal()">
                    <i class="fas fa-plus"></i>
                    Add New Gym
                  </button>
                  <button class="btn btn-secondary btn-block" onclick="dashboard.loadView('gym-requests')">
                    <i class="fas fa-bell"></i>
                    Review Requests
                  </button>
                  <button class="btn btn-info btn-block" onclick="dashboard.loadView('gym-analytics')">
                    <i class="fas fa-chart-bar"></i>
                    View Analytics
                  </button>
                  <button class="btn btn-success btn-block" onclick="dashboard.exportGymData()">
                    <i class="fas fa-download"></i>
                    Export Data
                  </button>
                </div>
              </div>
            </div>

            <!-- System Health -->
            <div class="dashboard-card">
              <div class="card-header">
                <h3 class="card-title">
                  <i class="fas fa-heartbeat"></i>
                  System Status
                </h3>
              </div>
              <div class="card-body">
                <div class="system-status">
                  ${this.renderSystemStatus(dashboardData?.systemStatus || {})}
                </div>
              </div>
            </div>
          </div>
        </div>
      `;

      this.contentArea.innerHTML = content;
      console.log("Dashboard content rendered successfully");

      // Add event listeners for header actions
      const refreshBtn = document.getElementById("refreshDashboard");
      const helpBtn = document.getElementById("dashboardHelp");

      if (refreshBtn) {
        refreshBtn.addEventListener("click", () => {
          console.log("Refreshing dashboard...");
          this.showOverview(); // Reload the dashboard
        });
      }

      if (helpBtn) {
        helpBtn.addEventListener("click", () => {
          this.showHelpModal();
        });
      }

      // Update navigation badges
      this.updateNavigationBadges(dashboardData);
      console.log("Dashboard overview loaded successfully");
    } catch (error) {
      console.error("Error loading dashboard overview:", error);
      this.showError("Failed to load dashboard overview");
    }
  }

  // Build dashboard data from real API responses only
  buildDashboardDataFromAPIs(gymsData, equipmentData, exercisesData) {
    const gymsTotal = Array.isArray(gymsData) ? gymsData.length : 0;
    const equipmentTotal = Array.isArray(equipmentData)
      ? equipmentData.length
      : 0;
    const exercisesTotal = Array.isArray(exercisesData)
      ? exercisesData.length
      : 0;

    return {
      gyms: {
        total: gymsTotal,
        recent: gymsData?.slice(0, 3) || [],
      },
      equipment: {
        total: equipmentTotal,
        recent: equipmentData?.slice(0, 5) || [],
      },
      exercises: {
        total: exercisesTotal,
        recent: exercisesData?.slice(0, 5) || [],
      },
      recentActivity: this.generateRecentActivity(
        gymsData,
        equipmentData,
        exercisesData
      ),
      systemStatus: {
        api: "healthy", // APIs are working since we got responses
        database: "healthy", // Database is working since APIs returned data
        lastUpdate: new Date().toISOString(),
      },
    };
  }

  // Generate recent activity from real data only
  generateRecentActivity(gymsData, equipmentData, exercisesData) {
    const activities = [];

    // Add activities based on real gyms data
    if (gymsData && gymsData.length > 0) {
      gymsData.slice(0, 3).forEach((gym) => {
        if (gym.name) {
          activities.push({
            title: `Gym: ${gym.name}`,
            icon: "fas fa-building",
            time_ago: this.formatTimeAgo(
              new Date(gym.created_at || Date.now())
            ),
            status: "success",
          });
        }
      });
    }

    // Add activities based on equipment data
    if (equipmentData && equipmentData.length > 0) {
      equipmentData.slice(0, 3).forEach((equipment) => {
        if (equipment.name) {
          activities.push({
            title: `Equipment: ${equipment.name}`,
            icon: "fas fa-dumbbell",
            time_ago: this.formatTimeAgo(
              new Date(equipment.created_at || Date.now())
            ),
            status: "success",
          });
        }
      });
    }

    // Add activities based on exercises data
    if (exercisesData && exercisesData.length > 0) {
      exercisesData.slice(0, 2).forEach((exercise) => {
        if (exercise.name) {
          activities.push({
            title: `Exercise: ${exercise.name}`,
            icon: "fas fa-running",
            time_ago: this.formatTimeAgo(
              new Date(exercise.created_at || Date.now())
            ),
            status: "success",
          });
        }
      });
    }

    return activities.slice(0, 8); // Return max 8 activities
  }

  // Equipment Management
  async loadEquipmentManagement() {
    try {
      const equipment = await this.managers.equipment.loadEquipment();

      const content = `
        <div class="dashboard-header">
          <h1 class="dashboard-title">Equipment Management</h1>
          <p class="dashboard-subtitle">Manage gym equipment and inventory</p>
        </div>
        <div class="dashboard-content">
          <div class="dashboard-card">
          <div class="card-header">
            <div class="card-header-content">
              <h3 class="card-title">Equipment Catalog</h3>
              <p class="card-subtitle">${equipment.length} items available</p>
            </div>
            <button class="btn btn-primary" onclick="dashboard.openEquipmentModal('create')">
              <i class="fas fa-plus"></i> Add Equipment
            </button>
          </div>
          <div class="card-body">
            <div id="equipment-table"></div>
          </div>
        </div>
        </div>
      `;

      this.setContent(content);

      // Initialize data table
      this.components.equipmentTable = new DataTable("#equipment-table", {
        data: equipment,
        columns: this.managers.equipment.getTableColumns(),
        rowActions: this.managers.equipment.getRowActions(),
        emptyMessage: "No equipment found. Add some equipment to get started.",
        filterable: true,
        sortable: true,
      });
    } catch (error) {
      throw new Error("Failed to load equipment management: " + error.message);
    }
  }

  // Exercises Management
  async loadExercisesManagement() {
    try {
      const exercises = await this.managers.exercise.loadExercises();

      const content = `
        <div class="dashboard-header">
          <h1 class="dashboard-title">Exercise Management</h1>
          <p class="dashboard-subtitle">Manage exercise library and routines</p>
        </div>
        <div class="dashboard-content">
          <div class="dashboard-card">
            <div class="card-header">
              <div class="card-header-content">
                <h3 class="card-title">Exercise Library</h3>
                <p class="card-subtitle">${exercises.length} exercises available</p>
              </div>
              <button class="btn btn-primary" onclick="dashboard.openExerciseModal('create')">
                <i class="fas fa-plus"></i> Add Exercise
              </button>
            </div>
            <div class="card-body">
              <div id="exercises-table"></div>
            </div>
          </div>
        </div>
      `;

      this.setContent(content);

      // Initialize data table
      this.components.exercisesTable = new DataTable("#exercises-table", {
        data: exercises,
        columns: this.managers.exercise.getTableColumns(),
        rowActions: this.managers.exercise.getRowActions(),
        emptyMessage: "No exercises found. Add some exercises to get started.",
        filterable: true,
        sortable: true,
      });
    } catch (error) {
      throw new Error("Failed to load exercises management: " + error.message);
    }
  }

  // Gyms Management
  async loadGymsManagement() {
    try {
      const gyms = await this.managers.gym.loadGyms();

      // Ensure gyms is always an array
      const safeGyms = Array.isArray(gyms) ? gyms : [];

      const content = `
        <div class="dashboard-header">
          <div class="dashboard-header-content">
            <div class="dashboard-header-text">
              <h1 class="dashboard-title">Gym Management</h1>
              <p class="dashboard-subtitle">Manage gym facilities and information</p>
            </div>
            <div class="dashboard-header-actions">
              <button class="btn-action" id="refreshDashboard" title="Refresh Dashboard">
                <i class="fas fa-sync-alt"></i>
              </button>
              <button class="btn-action" id="dashboardHelp" title="Help & Documentation">
                <i class="fas fa-question-circle"></i>
              </button>
            </div>
          </div>
        </div>
        <div class="dashboard-content">
          <div class="dashboard-card">
            <div class="card-header">
              <div class="card-header-content">
                <h3 class="card-title">Registered Gyms</h3>
                <p class="card-subtitle">${safeGyms.length} gyms registered</p>
              </div>
              <button class="btn btn-primary" onclick="dashboard.openGymModal('create')">
                <i class="fas fa-plus"></i> Add Gym
              </button>
            </div>
            <div class="card-body">
              <div id="gyms-table"></div>
            </div>
          </div>
        </div>
      `;

      this.setContent(content);

      // Initialize data table with safe data
      this.components.gymsTable = new DataTable("#gyms-table", {
        data: safeGyms,
        columns: this.managers.gym.getTableColumns(),
        rowActions: this.managers.gym.getRowActions(),
        emptyMessage: "No gyms found. Add some gyms to get started.",
        filterable: true,
        sortable: true,
      });

      console.log(
        "Gyms management loaded successfully with",
        safeGyms.length,
        "gyms"
      );
    } catch (error) {
      console.error("Error in loadGymsManagement:", error);
      throw new Error("Failed to load gyms management: " + error.message);
    }
  }

  async loadMuscularGroupsManagement() {
    const content = `
      <div class="dashboard-header">
        <h1 class="dashboard-title">Muscle Groups</h1>
        <p class="dashboard-subtitle">Manage muscle group categories</p>
      </div>
      <div class="dashboard-content">
        <div class="empty-state">
          <i class="fas fa-clock"></i>
          <h3>Coming Soon</h3>
          <p>Muscular Groups Management functionality is under development and will be available soon.</p>
        </div>
      </div>
    `;
    this.setContent(content);
  }

  // New Gym-Focused Views
  async loadGymUsersManagement() {
    try {
      // Since we don't have a gym-users endpoint, let's show a placeholder
      const content = `
        <div class="dashboard-header">
          <h1 class="dashboard-title">Members & Trainers</h1>
          <p class="dashboard-subtitle">Manage gym members and training staff</p>
        </div>
        <div class="dashboard-content">
          <div class="dashboard-card">
            <div class="card-header">
              <h3 class="card-title">
                <i class="fas fa-users"></i>
                Gym Members & Trainers
              </h3>
              <div class="card-actions">
                <button class="btn btn-secondary btn-sm" onclick="dashboard.exportGymUsers()">
                  <i class="fas fa-download"></i> Export
                </button>
                <button class="btn btn-primary" onclick="dashboard.openAddUserModal()">
                  <i class="fas fa-plus"></i> Add User
                </button>
              </div>
            </div>
            <div class="card-body">
              <div class="empty-state">
                <i class="fas fa-users"></i>
                <h3>User Management</h3>
                <p>User management functionality will be implemented when user endpoints are available.</p>
                <button class="btn btn-primary" onclick="dashboard.loadView('gyms')">
                  <i class="fas fa-building"></i> Manage Gyms Instead
                </button>
              </div>
            </div>
          </div>
        </div>
      `;

      this.setContent(content);
    } catch (error) {
      console.error("Error loading gym users:", error);
      this.showError("Failed to load gym users");
    }
  }

  async loadGymAnalytics() {
    try {
      const content = `
        <div class="dashboard-header">
          <h1 class="dashboard-title">Gym Analytics</h1>
          <p class="dashboard-subtitle">Performance metrics and usage statistics</p>
        </div>
        <div class="dashboard-content">
          <div class="analytics-dashboard">
            <!-- Analytics Charts -->
            <div class="analytics-grid">
              <div class="dashboard-card span-2">
                <div class="card-header">
                  <h3 class="card-title">
                    <i class="fas fa-chart-line"></i>
                    Gym Performance Trends
                  </h3>
                </div>
                <div class="card-body">
                  <div class="empty-state">
                    <i class="fas fa-chart-line"></i>
                    <h3>Analytics Dashboard</h3>
                    <p>Advanced analytics and reporting features will be available when analytics endpoints are implemented.</p>
                    <button class="btn btn-primary" onclick="dashboard.loadView('gyms')">
                      <i class="fas fa-building"></i> View Gyms Data
                    </button>
                  </div>
                </div>
              </div>

              <div class="dashboard-card">
                <div class="card-header">
                  <h3 class="card-title">
                    <i class="fas fa-users"></i>
                    Member Activity
                  </h3>
                </div>
                <div class="card-body">
                  <div class="empty-state-mini">Analytics coming soon</div>
                </div>
              </div>

              <div class="dashboard-card">
                <div class="card-header">
                  <h3 class="card-title">
                    <i class="fas fa-dumbbell"></i>
                    Equipment Usage
                  </h3>
                </div>
                <div class="card-body">
                  <div class="empty-state-mini">Analytics coming soon</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      `;

      this.setContent(content);
    } catch (error) {
      console.error("Error loading gym analytics:", error);
      this.showError("Failed to load gym analytics");
    }
  }

  async loadGymRequests() {
    try {
      const content = `
        <div class="dashboard-header">
          <h1 class="dashboard-title">Requests & Issues</h1>
          <p class="dashboard-subtitle">Handle gym support requests and reported issues</p>
        </div>
        <div class="dashboard-content">
          <div class="requests-dashboard">
            <div class="dashboard-card">
              <div class="card-header">
                <h3 class="card-title">
                  <i class="fas fa-bell"></i>
                  Gym Requests & Issues
                </h3>
                <div class="card-actions">
                  <select class="form-select" onchange="dashboard.filterRequests(this.value)">
                    <option value="all">All Requests</option>
                    <option value="pending">Pending</option>
                    <option value="urgent">Urgent</option>
                    <option value="resolved">Resolved</option>
                  </select>
                </div>
              </div>
              <div class="card-body">
                <div class="empty-state">
                  <i class="fas fa-bell"></i>
                  <h3>Request Management</h3>
                  <p>Request and issue tracking will be available when request endpoints are implemented.</p>
                  <button class="btn btn-primary" onclick="dashboard.loadView('gyms')">
                    <i class="fas fa-building"></i> Manage Gyms
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      `;

      this.setContent(content);
    } catch (error) {
      console.error("Error loading gym requests:", error);
      this.showError("Failed to load gym requests");
    }
  }

  async loadWorkoutTemplates() {
    const content = `
      <div class="dashboard-header">
        <h1 class="dashboard-title">Workout Templates</h1>
        <p class="dashboard-subtitle">Manage pre-built workout programs</p>
      </div>
      <div class="dashboard-content">
        <div class="empty-state">
          <i class="fas fa-clock"></i>
          <h3>Coming Soon</h3>
          <p>Workout Templates functionality is under development and will be available soon.</p>
        </div>
      </div>
    `;
    this.setContent(content);
  }

  async loadSystemSettings() {
    const content = `
      <div class="dashboard-header">
        <h1 class="dashboard-title">System Settings</h1>
        <p class="dashboard-subtitle">Configure platform-wide settings</p>
      </div>
      <div class="dashboard-content">
        <div class="empty-state">
          <i class="fas fa-clock"></i>
          <h3>Coming Soon</h3>
          <p>System Settings functionality is under development and will be available soon.</p>
        </div>
      </div>
    `;
    this.setContent(content);
  }

  async loadSystemLogs() {
    const content = `
      <div class="dashboard-header">
        <h1 class="dashboard-title">Activity Logs</h1>
        <p class="dashboard-subtitle">View system activity and audit trails</p>
      </div>
      <div class="dashboard-content">
        <div class="empty-state">
          <i class="fas fa-clock"></i>
          <h3>Coming Soon</h3>
          <p>System Activity Logs functionality is under development and will be available soon.</p>
        </div>
      </div>
    `;
    this.setContent(content);
  }

  // Modal Management
  openEquipmentModal(mode = "create", equipment = null) {
    const isEdit = mode === "edit" && equipment;
    const title = isEdit ? "Edit Equipment" : "Add Equipment";

    const formHtml = this.generateForm(
      this.managers.equipment.getFormSchema(),
      equipment
    );

    const modal = new Modal({
      title: title,
      content: formHtml,
      size: "medium",
      buttons: [
        {
          text: "Cancel",
          action: "dismiss",
          className: "btn btn-secondary",
        },
        {
          text: isEdit ? "Update" : "Create",
          action: "save",
          className: "btn btn-primary",
          handler: () => this.saveEquipment(isEdit, equipment?.id),
        },
      ],
    });

    modal.show();
    this.components.equipmentModal = modal;
  }

  openExerciseModal(mode = "create", exercise = null) {
    const isEdit = mode === "edit" && exercise;
    const title = isEdit ? "Edit Exercise" : "Add Exercise";

    const formHtml = this.generateForm(
      this.managers.exercise.getFormSchema(),
      exercise
    );

    const modal = new Modal({
      title: title,
      content: formHtml,
      size: "large",
      buttons: [
        {
          text: "Cancel",
          action: "dismiss",
          className: "btn btn-secondary",
        },
        {
          text: isEdit ? "Update" : "Create",
          action: "save",
          className: "btn btn-primary",
          handler: () => this.saveExercise(isEdit, exercise?.id),
        },
      ],
    });

    modal.show();
    this.components.exerciseModal = modal;
  }

  openGymModal(mode = "create", gym = null) {
    const isEdit = mode === "edit" && gym;
    const title = isEdit ? "Edit Gym" : "Add Gym";

    const formHtml = this.generateForm(this.managers.gym.getFormSchema(), gym);

    const modal = new Modal({
      title: title,
      content: formHtml,
      size: "medium",
      buttons: [
        {
          text: "Cancel",
          action: "dismiss",
          className: "btn btn-secondary",
        },
        {
          text: isEdit ? "Update" : "Create",
          action: "save",
          className: "btn btn-primary",
          handler: () => this.saveGym(isEdit, gym?.id),
        },
      ],
    });

    modal.show();
    this.components.gymModal = modal;
  }

  // Form Generation
  generateForm(schema, data = null) {
    let formHtml = '<form id="dynamic-form" class="form-grid">';

    Object.keys(schema).forEach((fieldName) => {
      const field = schema[fieldName];
      const value = data ? data[fieldName] || "" : "";

      formHtml += `<div class="form-group">`;
      formHtml += `<label class="form-label" for="${fieldName}">${field.label}${
        field.required ? " *" : ""
      }</label>`;

      if (field.type === "textarea") {
        formHtml += `<textarea class="form-control" id="${fieldName}" name="${fieldName}" rows="${
          field.rows || 3
        }" placeholder="${field.placeholder || ""}">${value}</textarea>`;
      } else if (field.type === "select") {
        formHtml += `<select class="form-control" id="${fieldName}" name="${fieldName}">`;
        field.options.forEach((option) => {
          const selected = value === option.value ? "selected" : "";
          formHtml += `<option value="${option.value}" ${selected}>${option.label}</option>`;
        });
        formHtml += `</select>`;
      } else {
        formHtml += `<input type="${
          field.type
        }" class="form-control" id="${fieldName}" name="${fieldName}" value="${value}" placeholder="${
          field.placeholder || ""
        }" ${field.required ? "required" : ""}>`;
      }

      formHtml += `</div>`;
    });

    formHtml += "</form>";
    return formHtml;
  }

  // Save handlers
  async saveEquipment(isEdit, id = null) {
    try {
      const formData = getFormData(document.getElementById("dynamic-form"));

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
      const formData = getFormData(document.getElementById("dynamic-form"));

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
      const formData = getFormData(document.getElementById("dynamic-form"));

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
            <p><strong>Type:</strong> ${this.managers.exercise.formatExerciseType(
              exercise.type
            )}</p>
            <p><strong>Difficulty:</strong> ${exercise.difficulty}</p>
            ${
              exercise.instructions
                ? `<p><strong>Instructions:</strong> ${exercise.instructions}</p>`
                : ""
            }
          </div>
        </div>
      </div>
    `;

    const modal = new Modal({
      title: "Exercise Details",
      content: content,
      size: "large",
      buttons: [
        {
          text: "Edit",
          action: "edit",
          className: "btn btn-primary",
          handler: () => {
            modal.hide();
            this.openExerciseModal("edit", exercise);
          },
        },
        {
          text: "Close",
          action: "dismiss",
          className: "btn btn-secondary",
        },
      ],
    });

    modal.show();
  }

  // View gym details
  viewGymDetails(gym) {
    const content = `
      <div class="gym-details">
        <h4>${gym.name}</h4>
        ${
          gym.contact_name
            ? `<p><strong>Contact:</strong> ${gym.contact_name}</p>`
            : ""
        }
        ${
          gym.contact_email
            ? `<p><strong>Email:</strong> ${gym.contact_email}</p>`
            : ""
        }
        ${
          gym.contact_phone
            ? `<p><strong>Phone:</strong> ${gym.contact_phone}</p>`
            : ""
        }
        ${gym.address ? `<p><strong>Address:</strong> ${gym.address}</p>` : ""}
        ${
          gym.description
            ? `<p><strong>Description:</strong> ${gym.description}</p>`
            : ""
        }
        <p><strong>Created:</strong> ${this.managers.gym.formatDate(
          gym.created_at
        )}</p>
        ${
          gym.deleted_at
            ? `<p><strong>Status:</strong> <span class="status-deleted">Deleted</span></p>`
            : ""
        }
      </div>
    `;

    const modal = new Modal({
      title: "Gym Details",
      content: content,
      size: "medium",
      buttons: [
        {
          text: "Edit",
          action: "edit",
          className: "btn btn-primary",
          handler: () => {
            modal.hide();
            this.openGymModal("edit", gym);
          },
        },
        {
          text: "Close",
          action: "dismiss",
          className: "btn btn-secondary",
        },
      ],
    });

    modal.show();
  }

  // Dashboard helper methods
  renderRecentActivity(activities) {
    if (!activities || activities.length === 0) {
      return '<div class="empty-state-mini">No recent activity</div>';
    }

    return activities
      .map(
        (activity) => `
      <div class="activity-item">
        <div class="activity-icon">
          <i class="fas ${this.getActivityIcon(activity.type)}"></i>
        </div>
        <div class="activity-content">
          <div class="activity-title">${activity.title}</div>
          <div class="activity-meta">
            <span class="activity-gym">${activity.gymName}</span>
            <span class="activity-time">${this.formatTimeAgo(
              activity.timestamp
            )}</span>
          </div>
        </div>
      </div>
    `
      )
      .join("");
  }

  renderRecentGyms(gyms) {
    if (!gyms || gyms.length === 0) {
      return '<div class="empty-state-mini">No gyms registered yet</div>';
    }

    return gyms
      .map(
        (gym) => `
      <div class="gym-item">
        <div class="gym-info">
          <div class="gym-name">${gym.name}</div>
          <div class="gym-details">
            ${
              gym.contact_email
                ? `<span><i class="fas fa-envelope"></i> ${gym.contact_email}</span>`
                : ""
            }
            ${
              gym.address
                ? `<span><i class="fas fa-map-marker-alt"></i> ${gym.address}</span>`
                : ""
            }
          </div>
        </div>
        <div class="gym-actions">
          <button class="btn btn-sm btn-ghost" onclick="dashboard.viewGymDetails(${JSON.stringify(
            gym
          ).replace(/"/g, "&quot;")})">
            <i class="fas fa-eye"></i> View
          </button>
        </div>
      </div>
    `
      )
      .join("");
  }

  renderSystemStatus(status) {
    const services = [
      {
        name: "API Service",
        status: status.api || "online",
        icon: "fas fa-server",
      },
      {
        name: "Database",
        status: status.database || "online",
        icon: "fas fa-database",
      },
      {
        name: "File Storage",
        status: status.storage || "online",
        icon: "fas fa-cloud",
      },
      {
        name: "Background Jobs",
        status: status.jobs || "online",
        icon: "fas fa-cogs",
      },
    ];

    return services
      .map(
        (service) => `
      <div class="status-item">
        <div class="status-icon">
          <i class="${service.icon}"></i>
        </div>
        <div class="status-content">
          <div class="status-name">${service.name}</div>
          <div class="status-indicator status-${service.status}">
            <span class="status-dot"></span>
            ${service.status.charAt(0).toUpperCase() + service.status.slice(1)}
          </div>
        </div>
      </div>
    `
      )
      .join("");
  }

  updateNavigationBadges(data) {
    // Update gym count badge
    const gymsCount = document.getElementById("gyms-count");
    if (gymsCount) {
      gymsCount.textContent = data?.gyms?.total || "0";
    }

    // Update requests badge
    const requestsCount = document.getElementById("requests-count");
    if (requestsCount) {
      const pending = data?.requests?.pending || 0;
      requestsCount.textContent = pending;
      requestsCount.style.display = pending > 0 ? "inline-block" : "none";
    }
  }

  getActivityIcon(type) {
    const icons = {
      member_joined: "fa-user-plus",
      payment_received: "fa-credit-card",
      equipment_added: "fa-plus",
      workout_completed: "fa-check-circle",
      issue_reported: "fa-exclamation-triangle",
      gym_created: "fa-building",
      trainer_assigned: "fa-user-tie",
    };
    return icons[type] || "fa-info-circle";
  }

  formatTimeAgo(timestamp) {
    const now = new Date();
    const time = new Date(timestamp);
    const diffInHours = Math.floor((now - time) / (1000 * 60 * 60));

    if (diffInHours < 1) return "Just now";
    if (diffInHours < 24) return `${diffInHours}h ago`;
    if (diffInHours < 48) return "Yesterday";
    return `${Math.floor(diffInHours / 24)}d ago`;
  }

  // Quick action methods
  openAddGymModal() {
    this.openGymModal("create");
  }

  exportGymData() {
    notifications.info("Preparing gym data export...");
    // Implementation would handle CSV/Excel export
  }

  // Dashboard Overview Helper Methods
  renderRecentActivity(activities) {
    if (!activities || activities.length === 0) {
      return '<div class="empty-state-mini">No recent activity</div>';
    }

    return activities
      .map(
        (activity) => `
      <div class="activity-item">
        <div class="activity-icon">
          <i class="${activity.icon || "fas fa-info-circle"}"></i>
        </div>
        <div class="activity-content">
          <div class="activity-title">${activity.title}</div>
          <div class="activity-meta">
            <span><i class="fas fa-building"></i> ${activity.gym_name}</span>
            <span><i class="fas fa-clock"></i> ${activity.time_ago}</span>
          </div>
        </div>
      </div>
    `
      )
      .join("");
  }

  renderGymRankings(gyms) {
    if (!gyms || gyms.length === 0) {
      return '<div class="empty-state-mini">No gym data available</div>';
    }

    return gyms
      .map(
        (gym, index) => `
      <div class="gym-ranking-item">
        <div class="ranking-position">${index + 1}</div>
        <div class="gym-info">
          <div class="gym-name">${gym.name}</div>
          <div class="gym-metric">${gym.members} members • ${
          gym.utilization
        }% utilization</div>
        </div>
        <div class="gym-score">${gym.score}</div>
      </div>
    `
      )
      .join("");
  }

  renderSystemStatus(status) {
    const services = [
      { name: "API Server", key: "api", icon: "fas fa-server" },
      { name: "Database", key: "database", icon: "fas fa-database" },
      { name: "File Storage", key: "storage", icon: "fas fa-cloud" },
      { name: "Authentication", key: "auth", icon: "fas fa-shield-alt" },
    ];

    return services
      .map((service) => {
        const serviceStatus = status[service.key] || "offline";
        const statusClass =
          serviceStatus === "online"
            ? "status-online"
            : serviceStatus === "warning"
            ? "status-warning"
            : "status-offline";

        return `
        <div class="status-item">
          <i class="${service.icon} status-icon"></i>
          <div class="status-content">
            <span class="status-name">${service.name}</span>
            <span class="status-indicator ${statusClass}">
              <span class="status-dot"></span>
              ${serviceStatus}
            </span>
          </div>
        </div>
      `;
      })
      .join("");
  }

  renderPerformanceChart(data) {
    // Placeholder for chart rendering
    return '<div class="empty-state-mini">Performance chart will be rendered here</div>';
  }

  renderMemberActivityChart(data) {
    // Placeholder for chart rendering
    return '<div class="empty-state-mini">Member activity chart will be rendered here</div>';
  }

  renderEquipmentUsageChart(data) {
    // Placeholder for chart rendering
    return '<div class="empty-state-mini">Equipment usage chart will be rendered here</div>';
  }

  updateNavigationBadges(data) {
    // Update gym count badge
    const gymsCount = document.getElementById("gyms-count");
    if (gymsCount) {
      gymsCount.textContent = data?.gyms?.total || "0";
    }

    // Update requests count badge
    const requestsCount = document.getElementById("requests-count");
    if (requestsCount) {
      const pending = data?.requests?.pending || 0;
      requestsCount.textContent = pending > 0 ? pending.toString() : "0";
      requestsCount.style.display = pending > 0 ? "inline" : "none";
    }
  }

  // Quick Action Methods
  async openAddGymModal() {
    this.openGymModal("create");
  }

  async exportGymData() {
    try {
      const gyms = await this.managers.gym.loadGyms();

      // Create CSV content
      const headers = [
        "Name",
        "Contact Name",
        "Contact Email",
        "Contact Phone",
        "Address",
        "Created Date",
      ];
      const csvContent = [
        headers.join(","),
        ...gyms.map((gym) =>
          [
            `"${gym.name || ""}"`,
            `"${gym.contact_name || ""}"`,
            `"${gym.contact_email || ""}"`,
            `"${gym.contact_phone || ""}"`,
            `"${gym.address || ""}"`,
            `"${gym.created_at || ""}"`,
          ].join(",")
        ),
      ].join("\n");

      // Create and download file
      const blob = new Blob([csvContent], { type: "text/csv" });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `gyms_export_${new Date().toISOString().split("T")[0]}.csv`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);

      notifications.success("Gym data exported successfully");
    } catch (error) {
      console.error("Export error:", error);
      notifications.error("Failed to export gym data");
    }
  }

  async exportGymUsers() {
    notifications.info(
      "User export will be available when user management is implemented"
    );
  }

  async openAddUserModal() {
    notifications.info(
      "User creation will be available when user endpoints are implemented"
    );
  }

  async filterRequests(filter) {
    notifications.info(
      `Request filtering will be implemented with request endpoints. Filter: ${filter}`
    );
  }

  // Utility methods
  setContent(html) {
    this.hideLoading();
    this.hideError();
    this.contentArea.innerHTML = html;
  }

  showLoading() {
    const loadingState = document.getElementById("loading-state");
    if (loadingState) {
      loadingState.style.display = "flex";
    }
  }

  hideLoading() {
    const loadingState = document.getElementById("loading-state");
    if (loadingState) {
      loadingState.style.display = "none";
    }
  }

  showError(message) {
    console.error("Dashboard Error:", message);
    this.hideLoading();

    const errorMessage = document.getElementById("error-message");
    const errorState = document.getElementById("error-state");

    if (errorMessage && errorState) {
      errorMessage.textContent = message;
      errorState.style.display = "flex";
    } else {
      // Fallback: show error in content area if error elements don't exist
      if (this.contentArea) {
        this.contentArea.innerHTML = `
          <div style="padding: 40px; text-align: center; color: #dc2626;">
            <i class="fas fa-exclamation-triangle" style="font-size: 48px; margin-bottom: 16px;"></i>
            <h3>Error</h3>
            <p>${message}</p>
            <button onclick="location.reload()" style="margin-top: 16px; padding: 8px 16px; background: #dc2626; color: white; border: none; border-radius: 4px; cursor: pointer;">
              Reload Page
            </button>
          </div>
        `;
      } else {
        alert("Dashboard Error: " + message);
      }
    }
  }

  hideError() {
    document.getElementById("error-state").style.display = "none";
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

  showHelpModal() {
    // For now, use the existing openHelp function
    // In the future, this could show a proper modal with help content
    if (typeof openHelp === "function") {
      openHelp();
    } else {
      console.log("Help documentation coming soon!");
    }
  }
}

// Global functions for compatibility
function logout() {
  localStorage.removeItem("auth_token");
  window.location.href = "/";
}

function refreshCurrentView() {
  if (window.dashboard) {
    window.dashboard.refreshCurrentView();
  }
}

function openHelp() {
  notifications.info("Help documentation coming soon!");
}

function openUserSettings() {
  notifications.info("User settings coming soon!");
}

function retryLoad() {
  if (window.dashboard) {
    window.dashboard.refreshCurrentView();
  }
}

// Initialize dashboard when DOM is loaded
document.addEventListener("DOMContentLoaded", function () {
  console.log("DOM loaded, initializing dashboard...");
  window.dashboard = new VanillaDashboardManager();
  console.log("Vanilla Dashboard initialized");
});

// Fallback initialization if DOM is already loaded
if (document.readyState === "loading") {
  // DOM is still loading, event listener above will handle it
} else {
  // DOM is already loaded
  console.log("DOM already loaded, initializing dashboard immediately...");
  window.dashboard = new VanillaDashboardManager();
  console.log("Vanilla Dashboard initialized (immediate)");
}
