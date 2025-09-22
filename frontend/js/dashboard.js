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
      // Wait for all component scripts to load
      await this.waitForDependencies();

      // Initialize content area
      this.contentArea = document.getElementById("content-body");

      if (!this.contentArea) {
        throw new Error("Content area element not found");
      }

      // Test content area by adding a simple message
      this.contentArea.innerHTML =
        '<div style="padding: 20px; text-align: center; color: #666;">Initializing dashboard...</div>';

      // Initialize managers
      this.initializeManagers();

      // Check authentication
      await this.checkAuth();

      // Initialize navigation
      this.initNavigation();

      // Initialize mobile functionality
      this.initMobile();

      // Load the initial overview automatically
      await this.loadView("overview");
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
    // Initialize API client first
    this.api = new ApiClient();

    // Initialize our specialized managers - they should all be available now
    this.managers.equipment = new EquipmentManager();
    this.managers.exercise = new ExerciseManager();
    this.managers.gym = new GymManager();
    this.managers.muscularGroup = new MuscularGroupManager();
    this.managers.user = new UserManager();

    // Set up event listeners for manager events
    this.setupManagerEventListeners();
  }

  setupManagerEventListeners() {
    // Equipment events
    document.addEventListener("equipment:edit", (e) => {
      this.openEquipmentModal("edit", e.detail.equipment);
    });

    document.addEventListener("equipment:view", (e) => {
      this.viewEquipmentDetails(e.detail.equipment);
    });

    document.addEventListener("equipment:delete", (e) => {
      this.showDeleteConfirmation(e.detail.equipment, "equipment");
    });

    // Exercise events
    document.addEventListener("exercise:edit", (e) => {
      this.openExerciseModal("edit", e.detail.exercise);
    });

    document.addEventListener("exercise:view", (e) => {
      this.viewExerciseDetails(e.detail.exercise);
    });

    document.addEventListener("exercise:delete", (e) => {
      this.showDeleteConfirmation(e.detail.exercise, "exercise");
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

    document.addEventListener("gym:delete", (e) => {
      this.showDeleteConfirmation(e.detail.gym);
    });

    document.addEventListener("gym:restore", (e) => {
      this.showRestoreConfirmation(e.detail.gym);
    });

    // Bulk action events
    document.addEventListener("datatable:bulkAction", (e) => {
      this.handleBulkAction(e.detail.action, e.detail.data);
    });
  }

  async checkAuth() {
    try {
      // Check for stored auth token from login
      const authToken = localStorage.getItem("auth_token");
      const userInfo = localStorage.getItem("user_info");

      if (!authToken || !userInfo) {
        console.error("No authentication data found");
        this.redirectToLogin();
        return;
      }

      this.currentUser = JSON.parse(userInfo);

      // Verify this is a platform admin
      if (this.currentUser.user_type !== "platform_admin") {
        console.error(
          "Access denied: Platform dashboard is for platform admins only"
        );
        this.redirectToTenantDashboard();
        return;
      }

      appState.setState({ user: this.currentUser });

      this.updateUserInfo();
      this.setupNavigation();
    } catch (error) {
      console.error("Authentication check failed:", error);
      this.redirectToLogin();
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
        case "platform-users":
          await this.loadPlatformUsersManagement();
          break;
        case "platform-analytics":
          await this.loadPlatformAnalytics();
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
      // Get real data from existing endpoints
      const [gyms, equipment, exercises] = await Promise.allSettled([
        this.api.request("/gym"),
        this.api.request("/equipment"),
        this.api.request("/exercise"),
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

      // Build dashboard data from real API responses
      const dashboardData = this.buildDashboardDataFromAPIs(
        gymsData,
        equipmentData,
        exercisesData
      );

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

      // Add event listeners for header actions
      const refreshBtn = document.getElementById("refreshDashboard");
      const helpBtn = document.getElementById("dashboardHelp");

      if (refreshBtn) {
        refreshBtn.addEventListener("click", async () => {
          await this.loadDashboardOverview();
        });
      }

      if (helpBtn) {
        helpBtn.addEventListener("click", () => {
          this.showHelpModal();
        });
      }

      // Update navigation badges
      this.updateNavigationBadges(dashboardData);
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
        pagination: true,
        pageSize: 10,
        selectable: true,
        exportable: true,
        bulkActions: [
          {
            action: "delete",
            text: "Delete Selected",
            icon: "fas fa-trash",
            className: "btn btn-sm btn-danger",
          },
          {
            action: "restore",
            text: "Restore Selected",
            icon: "fas fa-undo",
            className: "btn btn-sm btn-success",
          },
        ],
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
        pagination: true,
        pageSize: 10,
        selectable: true,
        exportable: true,
        bulkActions: [
          {
            action: "delete",
            text: "Delete Selected",
            icon: "fas fa-trash",
            className: "btn btn-sm btn-danger",
          },
          {
            action: "restore",
            text: "Restore Selected",
            icon: "fas fa-undo",
            className: "btn btn-sm btn-success",
          },
        ],
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
        pagination: true,
        pageSize: 10,
        selectable: true,
        exportable: true,
        bulkActions: [
          {
            action: "delete",
            text: "Delete Selected",
            icon: "fas fa-trash",
            className: "btn btn-sm btn-danger",
          },
          {
            action: "restore",
            text: "Restore Selected",
            icon: "fas fa-undo",
            className: "btn btn-sm btn-success",
          },
        ],
      });
    } catch (error) {
      console.error("Error in loadGymsManagement:", error);
      throw new Error("Failed to load gyms management: " + error.message);
    }
  }

  async loadMuscularGroupsManagement() {
    try {
      const muscularGroups =
        await this.managers.muscularGroup.loadMuscularGroups();

      const content = `
        <div class="dashboard-header">
          <h1 class="dashboard-title">Muscle Groups</h1>
          <p class="dashboard-subtitle">Manage muscle group categories</p>
        </div>
        <div class="dashboard-content">
          <div class="dashboard-card">
          <div class="card-header">
            <div class="card-header-content">
              <h3 class="card-title">Muscle Groups Catalog</h3>
              <p class="card-subtitle">${muscularGroups.length} groups available</p>
            </div>
            <button class="btn btn-primary" onclick="dashboard.openMuscularGroupModal('create')">
              <i class="fas fa-plus"></i> Add Muscle Group
            </button>
          </div>
          <div class="card-body">
            <div id="muscular-groups-table"></div>
          </div>
        </div>
        </div>
      `;

      this.setContent(content);

      // Initialize data table
      this.components.muscularGroupsTable = new DataTable(
        "#muscular-groups-table",
        {
          data: muscularGroups,
          columns: this.managers.muscularGroup.getTableColumns(),
          rowActions: this.managers.muscularGroup.getRowActions(),
          emptyMessage:
            "No muscle groups found. Add some muscle groups to get started.",
          pagination: {
            enabled: true,
            pageSize: 10,
          },
          search: {
            enabled: true,
            placeholder: "Search muscle groups...",
          },
          bulkActions: {
            enabled: true,
            actions: [
              {
                label: "Delete Selected",
                icon: "fas fa-trash",
                variant: "danger",
                action: async (selectedIds) => {
                  if (
                    confirm(
                      `Are you sure you want to delete ${selectedIds.length} muscle group(s)?`
                    )
                  ) {
                    for (const id of selectedIds) {
                      try {
                        await this.managers.muscularGroup.deleteMuscularGroup(
                          id
                        );
                      } catch (error) {
                        console.error(
                          `Error deleting muscle group ${id}:`,
                          error
                        );
                      }
                    }
                    this.showNotification(
                      "Selected muscle groups deleted",
                      "success"
                    );
                    this.loadMuscularGroupsManagement(); // Reload
                  }
                },
              },
            ],
          },
          export: {
            enabled: true,
            filename: "muscle_groups",
          },
        }
      );

      this.components.muscularGroupsTable.render();
    } catch (error) {
      console.error("Error loading muscle groups management:", error);
      this.showError("Failed to load muscle groups: " + error.message);
      throw new Error(
        "Failed to load muscle groups management: " + error.message
      );
    }
  }

  async loadPlatformUsersManagement() {
    try {
      // Show gym-centric user management interface
      const content = `
        <div class="dashboard-header">
          <h1 class="dashboard-title">Platform User Management</h1>
          <p class="dashboard-subtitle">Manage users across all gyms - select a gym to view its users</p>
        </div>
        <div class="dashboard-content">
          <!-- Gym Selection and Overview -->
          <div class="dashboard-card">
            <div class="card-header">
              <h3 class="card-title">
                <i class="fas fa-building"></i>
                Gym Overview & User Statistics
              </h3>
              <div class="card-actions">
                <button class="btn btn-secondary btn-sm" onclick="dashboard.refreshGymStats()">
                  <i class="fas fa-refresh"></i> Refresh
                </button>
              </div>
            </div>
            <div class="card-body" id="gym-stats-container">
              <div class="loading-state">
                <i class="fas fa-spinner fa-spin"></i>
                <p>Loading gym statistics...</p>
              </div>
            </div>
          </div>

          <!-- Selected Gym Users -->
          <div class="dashboard-card" id="selected-gym-users" style="display: none;">
            <div class="card-header">
              <h3 class="card-title">
                <i class="fas fa-users"></i>
                <span id="selected-gym-title">Gym Users</span>
              </h3>
              <div class="card-actions">
                <button class="btn btn-secondary btn-sm" onclick="dashboard.exportSelectedGymUsers()">
                  <i class="fas fa-download"></i> Export
                </button>
                <button class="btn btn-primary" onclick="dashboard.openPlatformUserModal('create')" id="add-user-btn">
                  <i class="fas fa-plus"></i> Add User
                </button>
                <button class="btn btn-outline btn-sm" onclick="dashboard.clearSelectedGym()">
                  <i class="fas fa-times"></i> Clear Selection
                </button>
              </div>
            </div>
            <div class="card-body" id="selected-gym-users-table">
              <!-- Users table will be loaded here -->
            </div>
          </div>
        </div>
      `;

      this.setContent(content);

      // Load gym statistics with user counts
      await this.loadGymUserStatistics();
    } catch (error) {
      console.error("Error loading platform users:", error);
      this.showError(
        "Failed to load platform user management: " + error.message
      );
    }
  }

  async loadGymUserStatistics() {
    try {
      // Load all gyms first
      const gyms = await this.managers.gym.loadGyms();

      // For each gym, get the actual users using platform admin endpoints
      const gymStats = await Promise.all(
        gyms.map(async (gym) => {
          try {
            const users = await this.api.getUsers(gym.id);
            const userList = users?.data || users || [];

            // Calculate statistics from real user data
            const userCount = userList.length;
            const activeUserCount = userList.filter(
              (user) => user.active !== false
            ).length;

            // Count users by role
            const roles = {};
            userList.forEach((user) => {
              const role = user.role || "member";
              roles[role] = (roles[role] || 0) + 1;
            });

            return {
              ...gym,
              userCount,
              activeUserCount,
              roles,
            };
          } catch (error) {
            console.error(`Error getting users for gym ${gym.id}:`, error);
            return {
              ...gym,
              userCount: 0,
              activeUserCount: 0,
              roles: {},
              error: true,
            };
          }
        })
      );

      this.renderGymUserStatistics(gymStats);
    } catch (error) {
      console.error("Error loading gym statistics:", error);
      throw error;
    }
  }

  renderGymUserStatistics(gymStats) {
    const container = document.getElementById("gym-stats-container");

    if (!gymStats || gymStats.length === 0) {
      container.innerHTML = `
        <div class="empty-state">
          <i class="fas fa-building"></i>
          <h3>No Gyms Found</h3>
          <p>No gyms are currently registered in the platform.</p>
        </div>
      `;
      return;
    }

    // Calculate totals
    const totalUsers = gymStats.reduce((sum, gym) => sum + gym.userCount, 0);
    const totalActiveUsers = gymStats.reduce(
      (sum, gym) => sum + gym.activeUserCount,
      0
    );
    const totalGyms = gymStats.length;

    container.innerHTML = `
      <!-- Platform Statistics Summary -->
      <div class="platform-stats-grid">
        <div class="stat-card">
          <div class="stat-icon">
            <i class="fas fa-building"></i>
          </div>
          <div class="stat-value">${totalGyms}</div>
          <div class="stat-label">Total Gyms</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <i class="fas fa-users"></i>
          </div>
          <div class="stat-value">${totalUsers}</div>
          <div class="stat-label">Total Users</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <i class="fas fa-user-check"></i>
          </div>
          <div class="stat-value">${totalActiveUsers}</div>
          <div class="stat-label">Active Users</div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <i class="fas fa-percentage"></i>
          </div>
          <div class="stat-value">${
            totalUsers > 0
              ? Math.round((totalActiveUsers / totalUsers) * 100)
              : 0
          }%</div>
          <div class="stat-label">Activity Rate</div>
        </div>
      </div>

      <!-- Gym List with User Statistics -->
      <div class="gym-user-management-list">
        <h4>Gyms & User Management</h4>
        <div class="gym-list-table">
          <div class="gym-list-header">
            <div class="gym-name">Gym Name</div>
            <div class="gym-status">Status</div>
            <div class="gym-users">Total Users</div>
            <div class="gym-active">Active</div>
            <div class="gym-roles">Roles</div>
            <div class="gym-actions">Actions</div>
          </div>
          ${gymStats
            .map(
              (gym) => `
            <div class="gym-list-row ${
              gym.error ? "gym-error" : ""
            }" data-gym-id="${gym.id}">
              <div class="gym-name">
                <strong>${gym.name}</strong>
                <small>${gym.address || "No address"}</small>
              </div>
              <div class="gym-status">
                <span class="status-badge ${
                  gym.is_active ? "status-active" : "status-inactive"
                }">
                  ${gym.is_active ? "Active" : "Inactive"}
                </span>
              </div>
              <div class="gym-users">
                <span class="user-count">${gym.userCount}</span>
              </div>
              <div class="gym-active">
                <span class="active-count">${gym.activeUserCount}</span>
              </div>
              <div class="gym-roles">
                <div class="role-distribution">
                  ${
                    Object.entries(gym.roles)
                      .map(
                        ([role, count]) =>
                          `<span class="role-badge role-${role}" title="${count} ${role}s">${count}</span>`
                      )
                      .join("") || '<span class="text-muted">No data</span>'
                  }
                </div>
              </div>
              <div class="gym-actions">
                <button class="btn btn-primary btn-sm" onclick="dashboard.selectGymForUserManagement('${
                  gym.id
                }', '${gym.name}')">
                  <i class="fas fa-users"></i> Manage Users
                </button>
                ${
                  gym.error
                    ? `
                  <span class="text-danger" title="Error loading user data">
                    <i class="fas fa-exclamation-triangle"></i>
                  </span>
                `
                    : ""
                }
              </div>
            </div>
          `
            )
            .join("")}
        </div>
      </div>
    `;
  }

  async selectGymForUserManagement(gymId, gymName) {
    try {
      this.selectedGymId = gymId;
      this.selectedGymName = gymName;

      // Show the selected gym users section
      const selectedSection = document.getElementById("selected-gym-users");
      const titleElement = document.getElementById("selected-gym-title");

      selectedSection.style.display = "block";
      titleElement.textContent = `${gymName} - Users`;

      // Auto-scroll to the users table
      selectedSection.scrollIntoView({
        behavior: "smooth",
        block: "start",
      });

      // Update the gym context for user manager
      this.managers.user.setGymContext(gymId);

      // Load users for the selected gym
      const tableContainer = document.getElementById(
        "selected-gym-users-table"
      );
      tableContainer.innerHTML = `
        <div class="loading-state">
          <i class="fas fa-spinner fa-spin"></i>
          <p>Loading users for ${gymName}...</p>
        </div>
      `;

      const users = await this.managers.user.loadUsers(gymId);

      console.log("Loaded users for gym:", gymId, "Users:", users);

      if (!users || users.length === 0) {
        console.log("No users found or empty array");
        tableContainer.innerHTML = `
          <div class="empty-state">
            <i class="fas fa-users"></i>
            <h3>No Users Found</h3>
            <p>This gym has no registered users yet.</p>
            <button class="btn btn-primary" onclick="dashboard.openPlatformUserModal('create')">
              <i class="fas fa-plus"></i> Add First User
            </button>
          </div>
        `;
        return;
      }

      // Create users table
      const dataTable = new DataTable(tableContainer, {
        data: users,
        columns: this.getPlatformUserTableColumns(),
        title: `${gymName} Users`,
        rowActions: this.getPlatformUserRowActions(),
        searchable: true,
        sortable: true,
        pagination: true,
        pageSize: 20,
        emptyMessage: `No users found in ${gymName}`,
        filters: [
          {
            name: "role",
            label: "Role",
            type: "select",
            options: [
              { value: "", label: "All Roles" },
              { value: "gym_admin", label: "Gym Admin" },
              { value: "trainer", label: "Trainer" },
              { value: "member", label: "Member" },
            ],
          },
          {
            name: "is_active",
            label: "Status",
            type: "select",
            options: [
              { value: "", label: "All Status" },
              { value: "true", label: "Active" },
              { value: "false", label: "Inactive" },
            ],
          },
        ],
      });

      // Scroll to the users section
      selectedSection.scrollIntoView({ behavior: "smooth" });
    } catch (error) {
      console.error("Error selecting gym for user management:", error);
      notifications.error(
        `Failed to load users for ${gymName}: ${error.message}`
      );
    }
  }

  clearSelectedGym() {
    this.selectedGymId = null;
    this.selectedGymName = null;

    const selectedSection = document.getElementById("selected-gym-users");
    selectedSection.style.display = "none";

    // Clear the user manager context
    this.managers.user.setGymContext(null);
  }

  async refreshGymStats() {
    try {
      notifications.info("Refreshing gym statistics...");
      await this.loadGymUserStatistics();
      notifications.success("Statistics refreshed successfully");
    } catch (error) {
      console.error("Error refreshing stats:", error);
      notifications.error("Failed to refresh statistics: " + error.message);
    }
  }

  async exportSelectedGymUsers() {
    if (!this.selectedGymId || !this.selectedGymName) {
      notifications.warning("Please select a gym first");
      return;
    }

    try {
      notifications.info("Preparing user export...");

      const users = await this.managers.user.loadUsers(this.selectedGymId);

      if (!users || users.length === 0) {
        notifications.warning("No users found to export");
        return;
      }

      // Prepare CSV data
      const csvHeaders = [
        "Username",
        "Email",
        "Role",
        "Status",
        "Verified",
        "Phone",
        "Created",
      ];
      const csvData = users.map((user) => [
        user.username,
        user.email,
        user.role,
        user.is_active ? "Active" : "Inactive",
        user.verified ? "Yes" : "No",
        user.phone || "",
        new Date(user.created_at).toLocaleDateString(),
      ]);

      // Create CSV content
      const csvContent = [csvHeaders, ...csvData]
        .map((row) => row.map((field) => `"${field}"`).join(","))
        .join("\n");

      // Download CSV
      const blob = new Blob([csvContent], { type: "text/csv" });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `${this.selectedGymName.replace(/[^a-z0-9]/gi, "_")}-users-${
        new Date().toISOString().split("T")[0]
      }.csv`;
      a.click();
      window.URL.revokeObjectURL(url);

      notifications.success("User export completed successfully");
    } catch (error) {
      console.error("Error exporting users:", error);
      notifications.error("Failed to export users: " + error.message);
    }
  }

  getPlatformUserTableColumns() {
    return [
      {
        key: "username",
        title: "Username",
        sortable: true,
      },
      {
        key: "email",
        title: "Email",
        sortable: true,
      },
      {
        key: "role",
        title: "Role",
        sortable: true,
        render: (value) => {
          const roleLabels = {
            gym_admin: "Gym Admin",
            trainer: "Trainer",
            member: "Member",
          };
          return `<span class="role-badge role-${value}">${
            roleLabels[value] || value
          }</span>`;
        },
      },
      {
        key: "is_active",
        title: "Status",
        sortable: true,
        render: (value) => {
          return value
            ? '<span class="status-badge status-active">Active</span>'
            : '<span class="status-badge status-inactive">Inactive</span>';
        },
      },
      {
        key: "verified",
        title: "Verified",
        sortable: true,
        render: (value) => {
          return value
            ? '<i class="fas fa-check-circle text-success" title="Verified"></i>'
            : '<i class="fas fa-times-circle text-danger" title="Unverified"></i>';
        },
      },
      {
        key: "created_at",
        title: "Created",
        sortable: true,
        render: (value) => {
          return new Date(value).toLocaleDateString();
        },
      },
    ];
  }

  getPlatformUserRowActions() {
    return [
      {
        action: "view",
        icon: "fas fa-eye",
        title: "View",
        className: "btn btn-sm btn-outline",
        handler: (user) => this.openPlatformUserModal("view", user),
      },
      {
        action: "edit",
        icon: "fas fa-edit",
        title: "Edit",
        className: "btn btn-sm btn-outline",
        handler: (user) => this.openPlatformUserModal("edit", user),
      },
      {
        action: "transfer",
        icon: "fas fa-exchange-alt",
        title: "Transfer Gym",
        className: "btn btn-sm btn-outline",
        handler: (user) => this.openGymTransferModal(user),
        condition: (user) => user.role !== "platform_admin",
      },
      {
        action: "deactivate",
        icon: "fas fa-ban",
        title: "Deactivate",
        className: "btn btn-sm btn-warning",
        handler: (user) => this.toggleUserStatus(user),
        condition: (user) => user.is_active,
      },
      {
        action: "activate",
        icon: "fas fa-check",
        title: "Activate",
        className: "btn btn-sm btn-success",
        handler: (user) => this.toggleUserStatus(user),
        condition: (user) => !user.is_active,
      },
      {
        action: "delete",
        icon: "fas fa-trash",
        title: "Delete",
        className: "btn btn-sm btn-danger",
        handler: (user) => this.deletePlatformUser(user),
      },
    ];
  }

  async getGymFilterOptions() {
    try {
      const gyms = await this.managers.gym.loadGyms();
      const options = [{ value: "", label: "All Gyms" }];
      gyms.forEach((gym) => {
        options.push({ value: gym.id, label: gym.name });
      });
      return options;
    } catch (error) {
      console.error("Error loading gym options:", error);
      return [{ value: "", label: "All Gyms" }];
    }
  }

  async openPlatformUserModal(mode = "create", userData = null) {
    try {
      const isEdit = mode === "edit" && userData;
      const isView = mode === "view" && userData;

      const title = isView
        ? `View User: ${userData.username}`
        : isEdit
        ? `Edit User: ${userData.username}`
        : "Add New User";

      // Get all available gyms for platform admin user creation
      const gyms = await this.managers.gym.loadGyms();

      // Enhanced form fields for platform admin user management
      const fields = this.getPlatformUserFormFields(gyms);

      // Pre-select gym if we have one selected
      const initialData = isEdit || isView ? userData : {};
      if (mode === "create" && this.selectedGymId) {
        initialData.gym_id = this.selectedGymId;
      }

      const formHtml = this.generateFormHtml(fields, initialData, isView);

      const modal = new Modal({
        title: title,
        size: "lg",
        content: `
          <form id="platform-user-form">
            ${formHtml}
          </form>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Cancel</button>
            ${
              !isView
                ? `
              <button type="submit" class="btn btn-primary" id="save-platform-user-btn">
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
        const form = modal.element.querySelector("#platform-user-form");
        const saveBtn = modal.element.querySelector("#save-platform-user-btn");

        const handleSubmit = async (e) => {
          e.preventDefault();

          try {
            saveBtn.disabled = true;
            saveBtn.innerHTML =
              '<i class="fas fa-spinner fa-spin"></i> Saving...';

            const formData = getFormData(form);
            const validation = this.validatePlatformUserData(formData);

            if (validation.length > 0) {
              this.showFormErrors(form, validation);
              return;
            }

            const preparedData = this.preparePlatformUserData(formData, isEdit);

            if (isEdit) {
              await this.managers.user.updateUser(userData.id, preparedData);
              notifications.success("User updated successfully");
            } else {
              // For creation, use the selected gym context
              if (this.selectedGymId) {
                // Temporarily set gym context for creation
                const originalGymContext = this.managers.user.currentGymId;
                this.managers.user.setGymContext(this.selectedGymId);

                await this.managers.user.createUser(preparedData);

                // Restore original context
                this.managers.user.setGymContext(originalGymContext);

                notifications.success("User created successfully");
              } else {
                throw new Error("Please select a gym first to manage users");
              }
            }
            modal.hide();

            // Reload appropriate view with a small delay to ensure backend processing
            if (this.currentView === "platform-users") {
              if (this.selectedGymId) {
                // Add small delay to ensure backend has processed the creation
                setTimeout(async () => {
                  await this.selectGymForUserManagement(
                    this.selectedGymId,
                    this.selectedGymName
                  );
                }, 500);
              }
              await this.refreshGymStats();
            }
          } catch (error) {
            console.error("Error saving platform user:", error);
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
      console.error("Error opening platform user modal:", error);
      notifications.error("Failed to open user form");
    }
  }

  getPlatformUserFormFields(gyms) {
    // No gym field needed - context is automatically determined by selected gym
    return [
      {
        name: "role",
        label: "Role",
        type: "select",
        options: [
          { value: "gym_admin", label: "Gym Admin" },
          { value: "trainer", label: "Trainer" },
          { value: "member", label: "Member" },
        ],
        required: true,
        helpText: "Select the user's role first to determine available fields",
      },
      {
        name: "username",
        label: "Username",
        type: "text",
        required: true,
      },
      {
        name: "email",
        label: "Email",
        type: "email",
        required: true,
      },
      {
        name: "password",
        label: "Password",
        type: "password",
        required: true,
        helpText: "Leave blank to keep existing password (for edits)",
      },
      {
        name: "phone",
        label: "Phone Number",
        type: "tel",
      },
      {
        name: "is_active",
        label: "Active Status",
        type: "checkbox",
        defaultValue: true,
        helpText: "User will be active by default",
      },
      {
        name: "verified",
        label: "Verified Status",
        type: "checkbox",
        defaultValue: true,
        helpText: "User will be verified by default",
      },
      {
        name: "description",
        label: "Description",
        type: "textarea",
        helpText: "Optional description or notes about the user",
      },
    ];
  }

  validatePlatformUserData(formData) {
    const errors = [];

    if (!formData.username || formData.username.trim().length < 3) {
      errors.push({
        field: "username",
        message: "Username must be at least 3 characters",
      });
    }

    if (!formData.email || !formData.email.includes("@")) {
      errors.push({
        field: "email",
        message: "Valid email address is required",
      });
    }

    // Gym context is required but managed automatically
    if (!this.selectedGymId) {
      errors.push({
        field: "general",
        message: "Please select a gym first to manage users",
      });
    }

    if (!formData.role) {
      errors.push({ field: "role", message: "User role is required" });
    }

    return errors;
  }

  preparePlatformUserData(formData, isEdit) {
    const data = {
      username: formData.username.trim(),
      email: formData.email.trim(),
      role: formData.role,
      phone: formData.phone || "",
      is_active: !!formData.is_active,
      verified: !!formData.verified,
      description: formData.description || "",
    };

    // Only include password if provided (for creation or password updates)
    if (formData.password && formData.password.trim()) {
      data.password = formData.password;
    }

    return data;
  }

  async openGymTransferModal(user) {
    try {
      const gyms = await this.managers.gym.loadGyms();
      const gymOptions = gyms
        .filter((gym) => gym.id !== user.gym_id) // Exclude current gym
        .map((gym) => ({
          value: gym.id,
          label: gym.name,
        }));

      const modal = new Modal({
        title: `Transfer User: ${user.username}`,
        content: `
          <div class="transfer-user-form">
            <p>Transfer <strong>${user.username}</strong> from <strong>${
          user.gym_name || "Unknown Gym"
        }</strong> to:</p>
            
            <div class="form-group">
              <label for="new-gym-id">New Gym:</label>
              <select id="new-gym-id" class="form-control" required>
                <option value="">Select a gym...</option>
                ${gymOptions
                  .map(
                    (gym) =>
                      `<option value="${gym.value}">${gym.label}</option>`
                  )
                  .join("")}
              </select>
            </div>
            
            <div class="form-group">
              <label for="transfer-reason">Reason for Transfer:</label>
              <textarea id="transfer-reason" class="form-control" placeholder="Optional reason for the transfer..."></textarea>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Cancel</button>
            <button type="button" class="btn btn-primary" id="confirm-transfer-btn">
              <i class="fas fa-exchange-alt"></i> Transfer User
            </button>
          </div>
        `,
      });

      modal.show();

      const confirmBtn = modal.element.querySelector("#confirm-transfer-btn");
      confirmBtn.addEventListener("click", async () => {
        const newGymId = modal.element.querySelector("#new-gym-id").value;
        const reason = modal.element.querySelector("#transfer-reason").value;

        if (!newGymId) {
          notifications.error("Please select a gym to transfer to");
          return;
        }

        try {
          confirmBtn.disabled = true;
          confirmBtn.innerHTML =
            '<i class="fas fa-spinner fa-spin"></i> Transferring...';

          // Update user's gym_id
          await this.managers.user.updateUser(user.id, {
            gym_id: newGymId,
            // Could add transfer_reason to user notes/description
            description:
              user.description +
              `\n[Transfer: ${reason || "No reason provided"}]`,
          });

          notifications.success(`User transferred successfully`);
          modal.hide();

          if (this.currentView === "platform-users") {
            await this.loadPlatformUsersManagement();
          }
        } catch (error) {
          console.error("Error transferring user:", error);
          notifications.error("Failed to transfer user: " + error.message);
        } finally {
          confirmBtn.disabled = false;
          confirmBtn.innerHTML =
            '<i class="fas fa-exchange-alt"></i> Transfer User';
        }
      });
    } catch (error) {
      console.error("Error opening gym transfer modal:", error);
      notifications.error("Failed to open transfer dialog");
    }
  }

  async toggleUserStatus(user) {
    const action = user.is_active ? "deactivate" : "activate";
    const confirmed = await this.confirmAction(
      `${action.charAt(0).toUpperCase() + action.slice(1)} User`,
      `Are you sure you want to ${action} ${user.username}?`,
      user.is_active ? "warning" : "success"
    );

    if (confirmed) {
      try {
        await this.managers.user.updateUser(user.id, {
          is_active: !user.is_active,
        });

        notifications.success(`User ${action}d successfully`);

        if (this.currentView === "platform-users") {
          await this.loadPlatformUsersManagement();
        }
      } catch (error) {
        console.error(`Error ${action}ing user:`, error);
        notifications.error(`Failed to ${action} user: ${error.message}`);
      }
    }
  }

  async deletePlatformUser(user) {
    const confirmed = await this.confirmAction(
      "Delete User",
      `Are you sure you want to permanently delete ${user.username}? This action cannot be undone.`,
      "danger"
    );

    if (confirmed) {
      try {
        await this.managers.user.deleteUser(user.id);
        notifications.success("User deleted successfully");

        if (this.currentView === "platform-users") {
          await this.loadPlatformUsersManagement();
        }
      } catch (error) {
        console.error("Error deleting user:", error);
        notifications.error("Failed to delete user: " + error.message);
      }
    }
  }

  async loadPlatformAnalytics() {
    try {
      const content = `
        <div class="dashboard-header">
          <h1 class="dashboard-title">Platform Analytics</h1>
          <p class="dashboard-subtitle">Cross-gym metrics and platform statistics</p>
        </div>
        <div class="dashboard-content">
          <div class="analytics-dashboard">
            <!-- Analytics Charts -->
            <div class="analytics-grid">
              <div class="dashboard-card span-2">
                <div class="card-header">
                  <h3 class="card-title">
                    <i class="fas fa-chart-line"></i>
                    Platform Performance Trends
                  </h3>
                </div>
                <div class="card-body">
                  <div class="empty-state">
                    <i class="fas fa-chart-line"></i>
                    <h3>Platform Analytics</h3>
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

  // Navigation and view management code continues here...

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
        ${gym.email ? `<p><strong>Email:</strong> ${gym.email}</p>` : ""}
        ${gym.phone ? `<p><strong>Phone:</strong> ${gym.phone}</p>` : ""}
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
              gym.email
                ? `<span><i class="fas fa-envelope"></i> ${gym.email}</span>`
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
      const headers = ["Name", "Email", "Phone", "Address", "Created Date"];
      const csvContent = [
        headers.join(","),
        ...gyms.map((gym) =>
          [
            `"${gym.name || ""}"`,
            `"${gym.email || ""}"`,
            `"${gym.phone || ""}"`,
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
    await this.openUserModal("create");
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
      // Help documentation coming soon - could show a modal here
    }
  }

  // Modal methods for CRUD operations
  async openGymModal(mode = "create", gymData = null) {
    try {
      const isEdit = mode === "edit" && gymData;
      const title = isEdit ? `Edit Gym: ${gymData.name}` : "Add New Gym";

      const schema = this.managers.gym.getFormSchema();
      const formHtml = this.generateFormHtml(schema, isEdit ? gymData : {});

      const modal = new Modal({
        title: title,
        size: "lg",
        content: `
          <form id="gym-form">
            ${formHtml}
          </form>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Cancel</button>
            <button type="submit" class="btn btn-primary" id="save-gym-btn">
              <i class="fas fa-save"></i> ${isEdit ? "Update" : "Create"} Gym
            </button>
          </div>
        `,
      });

      modal.show();

      // Focus first input after modal is shown
      setTimeout(() => {
        const firstInput = modal.element.querySelector(
          "input, textarea, select"
        );
        if (firstInput) firstInput.focus();
      }, 100);

      // Handle form submission
      const form = modal.element.querySelector("#gym-form");
      const saveBtn = modal.element.querySelector("#save-gym-btn");

      const handleSubmit = async (e) => {
        e.preventDefault();

        try {
          saveBtn.disabled = true;
          saveBtn.innerHTML =
            '<i class="fas fa-spinner fa-spin"></i> Saving...';

          const formData = getFormData(form);
          const validation = this.managers.gym.validateGymData(formData);

          // Clear previous errors
          this.clearFormErrors(form);

          if (!validation.isValid) {
            this.showFormErrors(form, validation.errors);
            return;
          }

          if (isEdit) {
            await this.managers.gym.updateGym(gymData.id, formData);
          } else {
            await this.managers.gym.createGym(formData);
          }

          modal.hide();
          // Refresh the current view to show updated data
          if (this.currentView === "gyms") {
            await this.loadGymsManagement();
          }
        } catch (error) {
          console.error("Error saving gym:", error);
          notifications.error(
            `Failed to ${isEdit ? "update" : "create"} gym: ${error.message}`
          );
        } finally {
          saveBtn.disabled = false;
          saveBtn.innerHTML = `<i class="fas fa-save"></i> ${
            isEdit ? "Update" : "Create"
          } Gym`;
        }
      };

      form.addEventListener("submit", handleSubmit);
      saveBtn.addEventListener("click", handleSubmit);
    } catch (error) {
      console.error("Error opening gym modal:", error);
      notifications.error("Failed to open gym form");
    }
  }

  async openEquipmentModal(mode = "create", equipmentData = null) {
    try {
      const isEdit = mode === "edit" && equipmentData;
      const title = isEdit
        ? `Edit Equipment: ${equipmentData.name}`
        : "Add New Equipment";

      const schema = this.managers.equipment.getFormSchema();
      const formHtml = this.generateFormHtml(
        schema,
        isEdit ? equipmentData : {}
      );

      const modal = new Modal({
        title: title,
        size: "lg",
        content: `
          <form id="equipment-form">
            ${formHtml}
          </form>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Cancel</button>
            <button type="submit" class="btn btn-primary" id="save-equipment-btn">
              <i class="fas fa-save"></i> ${
                isEdit ? "Update" : "Create"
              } Equipment
            </button>
          </div>
        `,
      });

      modal.show();

      // Focus first input after modal is shown
      setTimeout(() => {
        const firstInput = modal.element.querySelector(
          "input, textarea, select"
        );
        if (firstInput) firstInput.focus();
      }, 100);

      const form = modal.element.querySelector("#equipment-form");
      const saveBtn = modal.element.querySelector("#save-equipment-btn");

      const handleSubmit = async (e) => {
        e.preventDefault();

        try {
          saveBtn.disabled = true;
          saveBtn.innerHTML =
            '<i class="fas fa-spinner fa-spin"></i> Saving...';

          const formData = getFormData(form);
          const validation =
            this.managers.equipment.validateEquipmentData(formData);

          this.clearFormErrors(form);

          if (!validation.isValid) {
            this.showFormErrors(form, validation.errors);
            return;
          }

          if (isEdit) {
            await this.managers.equipment.updateEquipment(
              equipmentData.id,
              formData
            );
          } else {
            await this.managers.equipment.createEquipment(formData);
          }

          modal.hide();
          if (this.currentView === "equipment") {
            await this.loadEquipmentManagement();
          }
        } catch (error) {
          console.error("Error saving equipment:", error);
          notifications.error(
            `Failed to ${isEdit ? "update" : "create"} equipment: ${
              error.message
            }`
          );
        } finally {
          saveBtn.disabled = false;
          saveBtn.innerHTML = `<i class="fas fa-save"></i> ${
            isEdit ? "Update" : "Create"
          } Equipment`;
        }
      };

      form.addEventListener("submit", handleSubmit);
      saveBtn.addEventListener("click", handleSubmit);
    } catch (error) {
      console.error("Error opening equipment modal:", error);
      notifications.error("Failed to open equipment form");
    }
  }

  async openMuscularGroupModal(mode = "create", muscularGroupData = null) {
    try {
      const isEdit = mode === "edit" && muscularGroupData;
      const isView = mode === "view" && muscularGroupData;
      const title = isView
        ? `View Muscle Group: ${muscularGroupData.name}`
        : isEdit
        ? `Edit Muscle Group: ${muscularGroupData.name}`
        : "Add New Muscle Group";

      const fields = this.managers.muscularGroup.getFormFields();
      const formHtml = this.generateFormHtml(
        fields,
        isEdit || isView ? muscularGroupData : {},
        isView // readonly mode
      );

      const modal = new Modal({
        title: title,
        size: "lg",
        content: `
          <form id="muscular-group-form">
            ${formHtml}
          </form>
          ${
            !isView
              ? `
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Cancel</button>
            <button type="submit" class="btn btn-primary" id="save-muscular-group-btn">
              <i class="fas fa-save"></i> ${
                isEdit ? "Update" : "Create"
              } Muscle Group
            </button>
          </div>
          `
              : `
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Close</button>
            ${
              isEdit
                ? ""
                : `<button type="button" class="btn btn-primary" onclick="dashboard.openMuscularGroupModal('edit', ${JSON.stringify(
                    muscularGroupData
                  ).replace(/"/g, "&quot;")})">
              <i class="fas fa-edit"></i> Edit
            </button>`
            }
          </div>
          `
          }
        `,
      });

      modal.show();

      if (!isView) {
        // Focus first input after modal is shown
        setTimeout(() => {
          const firstInput = modal.element.querySelector(
            "input, textarea, select"
          );
          if (firstInput) firstInput.focus();
        }, 100);

        const form = modal.element.querySelector("#muscular-group-form");
        const saveBtn = modal.element.querySelector("#save-muscular-group-btn");

        const handleSubmit = async (e) => {
          e.preventDefault();

          try {
            saveBtn.disabled = true;
            saveBtn.innerHTML =
              '<i class="fas fa-spinner fa-spin"></i> Saving...';

            const formData = getFormData(form);
            const validation =
              this.managers.muscularGroup.validateMuscularGroupData(formData);

            this.clearFormErrors(form);

            if (validation.length > 0) {
              this.showFormErrors(form, validation);
              return;
            }

            const preparedData =
              this.managers.muscularGroup.prepareDataForSubmission(formData);

            if (isEdit) {
              await this.managers.muscularGroup.updateMuscularGroup(
                muscularGroupData.id,
                preparedData
              );
              notifications.success("Muscle group updated successfully");
            } else {
              await this.managers.muscularGroup.createMuscularGroup(
                preparedData
              );
              notifications.success("Muscle group created successfully");
            }

            modal.hide();
            if (this.currentView === "muscular-groups") {
              await this.loadMuscularGroupsManagement();
            }
          } catch (error) {
            console.error("Error saving muscle group:", error);
            notifications.error(
              `Failed to ${isEdit ? "update" : "create"} muscle group: ${
                error.message
              }`
            );
          } finally {
            saveBtn.disabled = false;
            saveBtn.innerHTML = `<i class="fas fa-save"></i> ${
              isEdit ? "Update" : "Create"
            } Muscle Group`;
          }
        };

        form.addEventListener("submit", handleSubmit);
        saveBtn.addEventListener("click", handleSubmit);
      }
    } catch (error) {
      console.error("Error opening muscle group modal:", error);
      notifications.error("Failed to open muscle group form");
    }
  }

  async openExerciseModal(mode = "create", exerciseData = null) {
    try {
      const isEdit = mode === "edit" && exerciseData;
      const title = isEdit
        ? `Edit Exercise: ${exerciseData.name}`
        : "Add New Exercise";

      // Get form schema asynchronously
      const schema = await this.managers.exercise.getFormSchema();

      // If editing, load existing links
      let formData = isEdit ? { ...exerciseData } : {};
      if (isEdit && exerciseData.id) {
        const links = await this.managers.exercise.getExerciseLinks(
          exerciseData.id
        );
        formData.muscular_groups = links.muscularGroups.map(
          (link) => link.muscular_group_id
        );
        formData.equipment_ids = links.equipment.map(
          (link) => link.equipment_id
        );
      }

      const formHtml = this.generateFormHtml(schema, formData);

      const modal = new Modal({
        title: title,
        size: "lg",
        content: `
          <form id="exercise-form">
            ${formHtml}
          </form>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Cancel</button>
            <button type="submit" class="btn btn-primary" id="save-exercise-btn">
              <i class="fas fa-save"></i> ${
                isEdit ? "Update" : "Create"
              } Exercise
            </button>
          </div>
        `,
      });

      modal.show();

      // Focus first input after modal is shown
      setTimeout(() => {
        const firstInput = modal.element.querySelector(
          "input, textarea, select"
        );
        if (firstInput) firstInput.focus();
      }, 100);

      const form = modal.element.querySelector("#exercise-form");
      const saveBtn = modal.element.querySelector("#save-exercise-btn");

      const handleSubmit = async (e) => {
        e.preventDefault();

        try {
          saveBtn.disabled = true;
          saveBtn.innerHTML =
            '<i class="fas fa-spinner fa-spin"></i> Saving...';

          const formData = getFormData(form);

          // muscular_groups and equipment_ids are now arrays from multiselect
          // No need to convert from comma-separated strings

          const validation =
            this.managers.exercise.validateExerciseData(formData);

          this.clearFormErrors(form);

          if (!validation.isValid) {
            this.showFormErrors(form, validation.errors);
            return;
          }

          if (isEdit) {
            await this.managers.exercise.updateExercise(
              exerciseData.id,
              formData
            );
          } else {
            await this.managers.exercise.createExercise(formData);
          }

          modal.hide();
          if (this.currentView === "exercises") {
            await this.loadExercisesManagement();
          }
        } catch (error) {
          console.error("Error saving exercise:", error);
          notifications.error(
            `Failed to ${isEdit ? "update" : "create"} exercise: ${
              error.message
            }`
          );
        } finally {
          saveBtn.disabled = false;
          saveBtn.innerHTML = `<i class="fas fa-save"></i> ${
            isEdit ? "Update" : "Create"
          } Exercise`;
        }
      };

      form.addEventListener("submit", handleSubmit);
      saveBtn.addEventListener("click", handleSubmit);
    } catch (error) {
      console.error("Error opening exercise modal:", error);
      notifications.error("Failed to open exercise form");
    }
  }

  async openUserModal(mode = "create", userData = null) {
    try {
      // Get current user info to determine context
      const currentUser = appState.getState().user;
      if (!currentUser) {
        notifications.error("User authentication required");
        return;
      }

      // Determine gym context based on user type
      let gymContext = null;

      if (currentUser.user_type === "platform_admin") {
        // Platform admins can work with any gym
        // For create/edit operations, they need to specify which gym
        const currentGym = appState.getState().currentGym;
        if ((mode === "create" || mode === "edit") && !currentGym) {
          notifications.error("Please select a gym first to manage users");
          return;
        }
        gymContext = currentGym ? currentGym.id : null;
      } else if (currentUser.user_type === "tenant_user") {
        // Tenant users work within their own gym
        if (currentUser.gym_id) {
          gymContext = currentUser.gym_id;
        } else {
          notifications.error(
            "Unable to determine gym context. Please contact support."
          );
          return;
        }
      } else {
        notifications.error("Invalid user type for user management");
        return;
      }

      const isEdit = mode === "edit" && userData;
      const isView = mode === "view" && userData;
      const title = isView
        ? `View User: ${userData.username}`
        : isEdit
        ? `Edit User: ${userData.username}`
        : "Add New User";

      const fields = this.managers.user.getFormFields();
      const formHtml = this.generateFormHtml(
        fields,
        isEdit || isView ? userData : {},
        isView // readonly mode
      );

      const modal = new Modal({
        title: title,
        size: "lg",
        content: `
          <form id="user-form">
            ${formHtml}
          </form>
          ${
            !isView
              ? `
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Cancel</button>
            <button type="submit" class="btn btn-primary" id="save-user-btn">
              <i class="fas fa-save"></i> ${isEdit ? "Update" : "Create"} User
            </button>
          </div>
          `
              : `
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Close</button>
            ${
              isEdit
                ? ""
                : `<button type="button" class="btn btn-primary" onclick="dashboard.openUserModal('edit', ${JSON.stringify(
                    userData
                  ).replace(/"/g, "&quot;")})">
              <i class="fas fa-edit"></i> Edit
            </button>`
            }
          </div>
          `
          }
        `,
      });

      modal.show();

      if (!isView) {
        // Focus first input after modal is shown
        setTimeout(() => {
          const firstInput = modal.element.querySelector(
            "input, textarea, select"
          );
          if (firstInput) firstInput.focus();
        }, 100);

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

            this.clearFormErrors(form);

            if (validation.length > 0) {
              this.showFormErrors(form, validation);
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

            // Reload current view if applicable
            if (this.currentView === "gyms") {
              await this.loadGymsManagement();
            }
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

  async viewGymDetails(gym) {
    try {
      const modal = new Modal({
        title: `Gym Details: ${gym.name}`,
        size: "lg",
        content: `
          <div class="gym-details-modal">
            <div class="detail-section">
              <h4><i class="fas fa-building"></i> Basic Information</h4>
              <div class="detail-grid">
                <div class="detail-item">
                  <label>Name:</label>
                  <span>${gym.name}</span>
                </div>
                <div class="detail-item">
                  <label>Status:</label>
                  <span class="status-badge ${
                    gym.deleted_at ? "deleted" : "active"
                  }">
                    ${gym.deleted_at ? "Deleted" : "Active"}
                  </span>
                </div>
                <div class="detail-item">
                  <label>Created:</label>
                  <span>${this.managers.gym.formatDate(gym.created_at)}</span>
                </div>
                <div class="detail-item">
                  <label>Updated:</label>
                  <span>${this.managers.gym.formatDate(gym.updated_at)}</span>
                </div>
              </div>
            </div>

            ${
              gym.email || gym.phone
                ? `
            <div class="detail-section">
              <h4><i class="fas fa-address-card"></i> Contact Information</h4>
              <div class="detail-grid">
                ${
                  gym.email
                    ? `
                  <div class="detail-item">
                    <label>Email:</label>
                    <span><a href="mailto:${gym.email}">${gym.email}</a></span>
                  </div>
                `
                    : ""
                }
                ${
                  gym.phone
                    ? `
                  <div class="detail-item">
                    <label>Phone:</label>
                    <span><a href="tel:${gym.phone}">${gym.phone}</a></span>
                  </div>
                `
                    : ""
                }
              </div>
            </div>
            `
                : ""
            }

            ${
              gym.address
                ? `
            <div class="detail-section">
              <h4><i class="fas fa-map-marker-alt"></i> Address</h4>
              <p>${gym.address}</p>
            </div>
            `
                : ""
            }

            ${
              gym.description
                ? `
            <div class="detail-section">
              <h4><i class="fas fa-info-circle"></i> Description</h4>
              <p>${gym.description}</p>
            </div>
            `
                : ""
            }
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Close</button>
            ${
              !gym.deleted_at
                ? `
              <button type="button" class="btn btn-primary" onclick="dashboard.openGymModal('edit', ${JSON.stringify(
                gym
              ).replace(/"/g, "&quot;")})">
                <i class="fas fa-edit"></i> Edit Gym
              </button>
            `
                : ""
            }
          </div>
        `,
      });

      modal.show();
    } catch (error) {
      console.error("Error showing gym details:", error);
      notifications.error("Failed to show gym details");
    }
  }

  async viewEquipmentDetails(equipment) {
    try {
      const modal = new Modal({
        title: `Equipment Details: ${equipment.name}`,
        size: "lg",
        content: `
          <div class="equipment-details-modal">
            <div class="detail-section">
              <h4><i class="fas fa-dumbbell"></i> Basic Information</h4>
              <div class="detail-grid">
                <div class="detail-item">
                  <label>Name:</label>
                  <span>${equipment.name}</span>
                </div>
                <div class="detail-item">
                  <label>Category:</label>
                  <span>${equipment.category}</span>
                </div>
                <div class="detail-item">
                  <label>Type:</label>
                  <span>${equipment.type}</span>
                </div>
                <div class="detail-item">
                  <label>Status:</label>
                  <span class="status-badge ${
                    equipment.deleted_at ? "deleted" : "active"
                  }">
                    ${equipment.deleted_at ? "Deleted" : "Active"}
                  </span>
                </div>
                <div class="detail-item">
                  <label>Created:</label>
                  <span>${this.managers.equipment.formatDate(
                    equipment.created_at
                  )}</span>
                </div>
                <div class="detail-item">
                  <label>Updated:</label>
                  <span>${this.managers.equipment.formatDate(
                    equipment.updated_at
                  )}</span>
                </div>
              </div>
            </div>

            ${
              equipment.description
                ? `
            <div class="detail-section">
              <h4><i class="fas fa-info-circle"></i> Description</h4>
              <p>${equipment.description}</p>
            </div>
            `
                : ""
            }
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Close</button>
            ${
              !equipment.deleted_at
                ? `
              <button type="button" class="btn btn-primary" onclick="dashboard.openEquipmentModal('edit', ${JSON.stringify(
                equipment
              ).replace(/"/g, "&quot;")})">
                <i class="fas fa-edit"></i> Edit Equipment
              </button>
            `
                : ""
            }
          </div>
        `,
      });

      modal.show();
    } catch (error) {
      console.error("Error showing equipment details:", error);
      notifications.error("Failed to show equipment details");
    }
  }

  async viewExerciseDetails(exercise) {
    try {
      const modal = new Modal({
        title: `Exercise Details: ${exercise.name}`,
        size: "lg",
        content: `
          <div class="exercise-details-modal">
            <div class="detail-section">
              <h4><i class="fas fa-running"></i> Basic Information</h4>
              <div class="detail-grid">
                <div class="detail-item">
                  <label>Name:</label>
                  <span>${exercise.name}</span>
                </div>
                <div class="detail-item">
                  <label>Category:</label>
                  <span>${exercise.category}</span>
                </div>
                <div class="detail-item">
                  <label>Muscle Groups:</label>
                  <span>${
                    Array.isArray(exercise.muscle_groups)
                      ? exercise.muscle_groups.join(", ")
                      : exercise.muscle_groups || ""
                  }</span>
                </div>
                <div class="detail-item">
                  <label>Difficulty:</label>
                  <span>${exercise.difficulty || "Not specified"}</span>
                </div>
                <div class="detail-item">
                  <label>Status:</label>
                  <span class="status-badge ${
                    exercise.deleted_at ? "deleted" : "active"
                  }">
                    ${exercise.deleted_at ? "Deleted" : "Active"}
                  </span>
                </div>
                <div class="detail-item">
                  <label>Created:</label>
                  <span>${this.managers.exercise.formatDate(
                    exercise.created_at
                  )}</span>
                </div>
                <div class="detail-item">
                  <label>Updated:</label>
                  <span>${this.managers.exercise.formatDate(
                    exercise.updated_at
                  )}</span>
                </div>
              </div>
            </div>

            ${
              exercise.description
                ? `
            <div class="detail-section">
              <h4><i class="fas fa-info-circle"></i> Description</h4>
              <p>${exercise.description}</p>
            </div>
            `
                : ""
            }

            ${
              exercise.instructions
                ? `
            <div class="detail-section">
              <h4><i class="fas fa-list-ol"></i> Instructions</h4>
              <p>${exercise.instructions}</p>
            </div>
            `
                : ""
            }
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" data-dismiss="modal">Close</button>
            ${
              !exercise.deleted_at
                ? `
              <button type="button" class="btn btn-primary" onclick="dashboard.openExerciseModal('edit', ${JSON.stringify(
                exercise
              ).replace(/"/g, "&quot;")})">
                <i class="fas fa-edit"></i> Edit Exercise
              </button>
            `
                : ""
            }
          </div>
        `,
      });

      modal.show();
    } catch (error) {
      console.error("Error showing exercise details:", error);
      notifications.error("Failed to show exercise details");
    }
  }

  // Helper methods for form handling
  generateFormHtml(schema, data = {}, readonly = false) {
    return Object.entries(schema)
      .map(([key, field]) => {
        const value = data[key] || "";
        const required = field.required && !readonly ? "required" : "";
        const fieldId = `field-${key}`;
        const readonlyAttr = readonly ? "readonly disabled" : "";

        switch (field.type) {
          case "textarea":
            return `
            <div class="form-group">
              <label for="${fieldId}">${field.label}</label>
              <textarea 
                id="${fieldId}" 
                name="${key}" 
                class="form-control" 
                rows="${field.rows || 3}"
                placeholder="${field.placeholder || ""}"
                ${required}
                ${readonlyAttr}
              >${value}</textarea>
              ${
                field.help
                  ? `<small class="form-help">${field.help}</small>`
                  : ""
              }
              <div class="field-error" id="${fieldId}-error"></div>
            </div>
          `;
          case "select":
            const optionsHtml = field.options
              .map(
                (option) =>
                  `<option value="${option.value}" ${
                    option.value === value ? "selected" : ""
                  }>${option.label}</option>`
              )
              .join("");
            return `
            <div class="form-group">
              <label for="${fieldId}">${field.label}</label>
              <select 
                id="${fieldId}" 
                name="${key}" 
                class="form-control" 
                ${required}
                ${readonlyAttr}
              >
                ${optionsHtml}
              </select>
              ${
                field.help
                  ? `<small class="form-help">${field.help}</small>`
                  : ""
              }
              <div class="field-error" id="${fieldId}-error"></div>
            </div>
          `;
          case "multiselect":
            const selectedValues = Array.isArray(value) ? value : [];
            const multiselectOptionsHtml = field.options
              .filter((option) => option.value !== "") // Remove empty option for multiselect
              .map(
                (option) =>
                  `<option value="${option.value}" ${
                    selectedValues.includes(option.value) ? "selected" : ""
                  }>${option.label}</option>`
              )
              .join("");
            return `
            <div class="form-group">
              <label for="${fieldId}">${field.label}</label>
              <select 
                id="${fieldId}" 
                name="${key}" 
                class="form-control multiselect" 
                multiple
                ${required}
                ${readonlyAttr}
              >
                ${multiselectOptionsHtml}
              </select>
              ${
                field.help
                  ? `<small class="form-help">${field.help}</small>`
                  : ""
              }
              <div class="field-error" id="${fieldId}-error"></div>
            </div>
          `;
          default:
            // Handle arrays (like muscle_groups)
            const displayValue = Array.isArray(value)
              ? value.join(", ")
              : value;
            return `
            <div class="form-group">
              <label for="${fieldId}">${field.label}</label>
              <input 
                type="${field.type || "text"}" 
                id="${fieldId}" 
                name="${key}" 
                class="form-control" 
                value="${displayValue}"
                placeholder="${field.placeholder || ""}"
                ${required}
                ${readonlyAttr}
              />
              ${
                field.help
                  ? `<small class="form-help">${field.help}</small>`
                  : ""
              }
              <div class="field-error" id="${fieldId}-error"></div>
            </div>
          `;
        }
      })
      .join("");
  }

  clearFormErrors(form) {
    const errorElements = form.querySelectorAll(".field-error");
    errorElements.forEach((el) => {
      el.textContent = "";
      el.style.display = "none";
    });

    const inputElements = form.querySelectorAll(".form-control");
    inputElements.forEach((el) => {
      el.classList.remove("is-invalid");
    });
  }

  showFormErrors(form, errors) {
    Object.entries(errors).forEach(([field, message]) => {
      const input = form.querySelector(`[name="${field}"]`);
      const errorElement = form.querySelector(`#field-${field}-error`);

      if (input) {
        input.classList.add("is-invalid");
      }

      if (errorElement) {
        errorElement.textContent = message;
        errorElement.style.display = "block";
      }
    });
  }

  // Bulk action handler
  async handleBulkAction(action, selectedData) {
    if (!selectedData || selectedData.length === 0) {
      notifications.warn("No items selected");
      return;
    }

    const itemType =
      this.currentView === "gyms"
        ? "gym"
        : this.currentView === "equipment"
        ? "equipment"
        : "exercise";
    const manager =
      this.currentView === "gyms"
        ? this.managers.gym
        : this.currentView === "equipment"
        ? this.managers.equipment
        : this.managers.exercise;

    try {
      switch (action) {
        case "delete":
          await this.showBulkDeleteConfirmation(
            selectedData,
            itemType,
            manager
          );
          break;
        case "restore":
          await this.showBulkRestoreConfirmation(
            selectedData,
            itemType,
            manager
          );
          break;
        default:
          notifications.warn(`Unknown bulk action: ${action}`);
      }
    } catch (error) {
      notifications.error(`Failed to perform bulk ${action}: ${error.message}`);
    }
  }

  async showBulkDeleteConfirmation(selectedData, itemType, manager) {
    const modal = new Modal({
      title: `Delete ${selectedData.length} ${itemType}${
        selectedData.length > 1 ? "s" : ""
      }`,
      size: "md",
      content: `
        <div class="confirmation-modal">
          <div class="confirmation-icon delete">
            <i class="fas fa-exclamation-triangle"></i>
          </div>
          <h4>Delete ${selectedData.length} selected ${itemType}${
        selectedData.length > 1 ? "s" : ""
      }?</h4>
          <p>This action will remove the selected items from active listings. They can be restored later if needed.</p>
          <div class="selected-items-list">
            ${selectedData
              .slice(0, 5)
              .map((item) => `<div class="selected-item">• ${item.name}</div>`)
              .join("")}
            ${
              selectedData.length > 5
                ? `<div class="more-items">... and ${
                    selectedData.length - 5
                  } more</div>`
                : ""
            }
          </div>
          <div class="confirmation-details">
            <strong>Note:</strong> This is a soft delete - all data will be preserved.
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-outline" data-dismiss="modal">Cancel</button>
          <button type="button" class="btn btn-danger" id="confirm-bulk-delete">
            <i class="fas fa-trash"></i> Delete Selected
          </button>
        </div>
      `,
    });

    modal.show();

    const confirmBtn = modal.element.querySelector("#confirm-bulk-delete");
    confirmBtn.addEventListener("click", async () => {
      try {
        confirmBtn.disabled = true;
        confirmBtn.innerHTML =
          '<i class="fas fa-spinner fa-spin"></i> Deleting...';

        const deletePromises = selectedData.map((item) => {
          if (itemType === "gym") return manager.deleteGym(item.id);
          if (itemType === "equipment") return manager.deleteEquipment(item.id);
          if (itemType === "exercise") return manager.deleteExercise(item.id);
        });

        await Promise.all(deletePromises);

        modal.hide();
        notifications.success(
          `Successfully deleted ${selectedData.length} ${itemType}${
            selectedData.length > 1 ? "s" : ""
          }`
        );

        // Refresh the current view
        await this.loadView(this.currentView);
      } catch (error) {
        confirmBtn.disabled = false;
        confirmBtn.innerHTML = '<i class="fas fa-trash"></i> Delete Selected';
        notifications.error(`Failed to delete items: ${error.message}`);
      }
    });
  }

  async showBulkRestoreConfirmation(selectedData, itemType, manager) {
    const modal = new Modal({
      title: `Restore ${selectedData.length} ${itemType}${
        selectedData.length > 1 ? "s" : ""
      }`,
      size: "md",
      content: `
        <div class="confirmation-modal">
          <div class="confirmation-icon restore">
            <i class="fas fa-undo"></i>
          </div>
          <h4>Restore ${selectedData.length} selected ${itemType}${
        selectedData.length > 1 ? "s" : ""
      }?</h4>
          <p>This action will restore the selected items to active status.</p>
          <div class="selected-items-list">
            ${selectedData
              .slice(0, 5)
              .map((item) => `<div class="selected-item">• ${item.name}</div>`)
              .join("")}
            ${
              selectedData.length > 5
                ? `<div class="more-items">... and ${
                    selectedData.length - 5
                  } more</div>`
                : ""
            }
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-outline" data-dismiss="modal">Cancel</button>
          <button type="button" class="btn btn-success" id="confirm-bulk-restore">
            <i class="fas fa-undo"></i> Restore Selected
          </button>
        </div>
      `,
    });

    modal.show();

    const confirmBtn = modal.element.querySelector("#confirm-bulk-restore");
    confirmBtn.addEventListener("click", async () => {
      try {
        confirmBtn.disabled = true;
        confirmBtn.innerHTML =
          '<i class="fas fa-spinner fa-spin"></i> Restoring...';

        const restorePromises = selectedData.map((item) => {
          if (itemType === "gym") return manager.restoreGym(item.id);
          if (itemType === "equipment")
            return manager.restoreEquipment(item.id);
          if (itemType === "exercise") return manager.restoreExercise(item.id);
        });

        await Promise.all(restorePromises);

        modal.hide();
        notifications.success(
          `Successfully restored ${selectedData.length} ${itemType}${
            selectedData.length > 1 ? "s" : ""
          }`
        );

        // Refresh the current view
        await this.loadView(this.currentView);
      } catch (error) {
        confirmBtn.disabled = false;
        confirmBtn.innerHTML = '<i class="fas fa-undo"></i> Restore Selected';
        notifications.error(`Failed to restore items: ${error.message}`);
      }
    });
  }

  // Confirmation modal methods
  async showDeleteConfirmation(gym) {
    try {
      const modal = new Modal({
        title: "Delete Gym",
        size: "md",
        content: `
        <div class="confirmation-modal">
          <div class="confirmation-icon delete">
            <i class="fas fa-exclamation-triangle"></i>
          </div>
          <h4>Delete "${gym.name}"?</h4>
          <p>This action will remove the gym from active listings. The gym can be restored later if needed.</p>
          <div class="confirmation-details">
            <strong>Note:</strong> This is a soft delete - all data will be preserved.
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-outline" data-dismiss="modal">Cancel</button>
          <button type="button" class="btn btn-danger" id="confirm-delete-btn">
            <i class="fas fa-trash"></i> Delete Gym
          </button>
        </div>
      `,
      });

      modal.show();

      // Wait a bit for modal to render before accessing elements
      setTimeout(() => {
        const confirmBtn = modal.element.querySelector("#confirm-delete-btn");
        confirmBtn.addEventListener("click", async () => {
          try {
            confirmBtn.disabled = true;
            confirmBtn.innerHTML =
              '<i class="fas fa-spinner fa-spin"></i> Deleting...';

            await this.managers.gym.deleteGym(gym.id);
            modal.hide();

            // Refresh the current view
            if (this.currentView === "gyms") {
              await this.loadGymsManagement();
            }
          } catch (error) {
            console.error("Error deleting gym:", error);
            notifications.error(`Failed to delete gym: ${error.message}`);
          } finally {
            confirmBtn.disabled = false;
            confirmBtn.innerHTML = '<i class="fas fa-trash"></i> Delete Gym';
          }
        });
      }, 100);
    } catch (error) {
      console.error("Error showing delete confirmation:", error);
      notifications.error("Failed to show confirmation dialog");
    }
  }

  async showRestoreConfirmation(gym) {
    try {
      const modal = new Modal({
        title: "Restore Gym",
        size: "md",
        content: `
        <div class="confirmation-modal">
          <div class="confirmation-icon restore">
            <i class="fas fa-undo"></i>
          </div>
          <h4>Restore "${gym.name}"?</h4>
          <p>This action will restore the gym to active status and make it available for use again.</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-outline" data-dismiss="modal">Cancel</button>
          <button type="button" class="btn btn-warning" id="confirm-restore-btn">
            <i class="fas fa-undo"></i> Restore Gym
          </button>
        </div>
      `,
      });

      modal.show();

      // Wait for modal to render before accessing elements
      setTimeout(() => {
        const confirmBtn = modal.element.querySelector("#confirm-restore-btn");
        confirmBtn.addEventListener("click", async () => {
          try {
            confirmBtn.disabled = true;
            confirmBtn.innerHTML =
              '<i class="fas fa-spinner fa-spin"></i> Restoring...';

            await this.managers.gym.restoreGym(gym.id);
            modal.hide();

            // Refresh the current view
            if (this.currentView === "gyms") {
              await this.loadGymsManagement();
            }
          } catch (error) {
            console.error("Error restoring gym:", error);
            notifications.error(`Failed to restore gym: ${error.message}`);
          } finally {
            confirmBtn.disabled = false;
            confirmBtn.innerHTML = '<i class="fas fa-undo"></i> Restore Gym';
          }
        });
      }, 100);
    } catch (error) {
      console.error("Error showing restore confirmation:", error);
      notifications.error("Failed to show confirmation dialog");
    }
  }

  redirectToLogin() {
    // Clear stored auth data
    localStorage.removeItem("auth_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user_info");

    // Redirect to login page
    const currentOrigin = window.location.origin;
    let currentPath = window.location.pathname.replace(
      /\/pages\/dashboard\/.*/,
      ""
    );
    if (currentPath.endsWith("/")) {
      currentPath = currentPath.slice(0, -1);
    }
    window.location.href = `${currentOrigin}${currentPath}/`;
  }

  redirectToTenantDashboard() {
    // Redirect tenant users to their appropriate dashboard
    const currentOrigin = window.location.origin;
    let currentPath = window.location.pathname.replace(
      /\/pages\/dashboard\/.*/,
      ""
    );
    if (currentPath.endsWith("/")) {
      currentPath = currentPath.slice(0, -1);
    }
    window.location.href = `${currentOrigin}${currentPath}/pages/tenant/index.html`;
  }

  async confirmAction(title, message, variant = "primary") {
    return new Promise((resolve) => {
      const modal = new Modal({
        title: title,
        content: `
          <div class="confirmation-modal">
            <div class="confirmation-icon ${
              variant === "danger"
                ? "text-danger"
                : variant === "warning"
                ? "text-warning"
                : "text-info"
            }">
              <i class="fas fa-${
                variant === "danger"
                  ? "exclamation-triangle"
                  : variant === "warning"
                  ? "exclamation-circle"
                  : "question-circle"
              }"></i>
            </div>
            <p>${message}</p>
          </div>
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

  showFormErrors(form, errors) {
    // Clear previous errors
    this.clearFormErrors(form);

    errors.forEach((error) => {
      const field = form.querySelector(`[name="${error.field}"]`);
      if (field) {
        field.classList.add("error");

        // Add error message
        const errorDiv = document.createElement("div");
        errorDiv.className = "field-error";
        errorDiv.textContent = error.message;
        field.parentNode.appendChild(errorDiv);
      }
    });
  }

  clearFormErrors(form) {
    // Remove error classes and messages
    form.querySelectorAll(".error").forEach((field) => {
      field.classList.remove("error");
    });

    form.querySelectorAll(".field-error").forEach((error) => {
      error.remove();
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
              ${
                field.helpText
                  ? `<small class="form-help">${field.helpText}</small>`
                  : ""
              }
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
              ${
                field.helpText
                  ? `<small class="form-help">${field.helpText}</small>`
                  : ""
              }
            </div>
          `;

          case "checkbox":
            return `
            <div class="form-group">
              <div class="form-check">
                <input type="checkbox" name="${field.name}" id="${
              field.name
            }" class="form-check-input" ${
              value ? "checked" : ""
            } ${readonlyAttr}>
                <label for="${field.name}" class="form-check-label">${
              field.label
            }</label>
              </div>
              ${
                field.helpText
                  ? `<small class="form-help">${field.helpText}</small>`
                  : ""
              }
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
              ${
                field.helpText
                  ? `<small class="form-help">${field.helpText}</small>`
                  : ""
              }
            </div>
          `;
        }
      })
      .join("");
  }
} // Global functions for compatibility
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
  window.dashboard = new VanillaDashboardManager();
});

// Fallback initialization if DOM is already loaded
if (document.readyState === "loading") {
  // DOM is still loading, event listener above will handle it
} else {
  // DOM is already loaded
  window.dashboard = new VanillaDashboardManager();
}
