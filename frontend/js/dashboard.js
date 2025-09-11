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

    // Load initial view
    await this.loadView("overview");

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
    const tenantNav = document.querySelector(".tenant-nav");
    const gymSelector = document.getElementById("gym-selector");

    if (this.currentUser.user_type === "platform_admin") {
      platformAdminNav.style.display = "block";
      tenantNav.style.display = "none";
      gymSelector.style.display = "block";
      this.loadGyms();
    } else {
      platformAdminNav.style.display = "none";
      tenantNav.style.display = "block";
      gymSelector.style.display = "none";
      this.currentGym = this.currentUser.gym_id;
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
          option.textContent = gym.name;
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
      });
    });
  }

  async loadView(viewName) {
    this.currentView = viewName;

    // Update page title and action button
    this.updatePageHeader(viewName);

    // Show loading state
    this.showLoading();

    try {
      switch (viewName) {
        case "overview":
          await this.loadOverview();
          break;
        case "members":
          await this.loadMembers();
          break;
        case "exercises":
          await this.loadCustomExercises();
          break;
        case "equipment":
          await this.loadCustomEquipment();
          break;
        case "workouts":
          await this.loadWorkoutTemplates();
          break;
        case "workout-instances":
          await this.loadWorkoutInstances();
          break;
        case "analytics":
          await this.loadAnalytics();
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

  async loadOverview() {
    if (!this.currentGym && this.currentUser.role !== "platform_admin") {
      this.showError("No gym selected");
      return;
    }

    const html = `
            <div class="dashboard-grid">
                <div class="dashboard-card stat-card">
                    <div class="card-icon">
                        <i class="fas fa-users"></i>
                    </div>
                    <span class="stat-number" id="total-members">-</span>
                    <span class="stat-label">Total Members</span>
                </div>
                <div class="dashboard-card stat-card">
                    <div class="card-icon">
                        <i class="fas fa-dumbbell"></i>
                    </div>
                    <span class="stat-number" id="total-workouts">-</span>
                    <span class="stat-label">Workout Templates</span>
                </div>
                <div class="dashboard-card stat-card">
                    <div class="card-icon">
                        <i class="fas fa-running"></i>
                    </div>
                    <span class="stat-number" id="total-exercises">-</span>
                    <span class="stat-label">Custom Exercises</span>
                </div>
                <div class="dashboard-card stat-card">
                    <div class="card-icon">
                        <i class="fas fa-tools"></i>
                    </div>
                    <span class="stat-number" id="total-equipment">-</span>
                    <span class="stat-label">Equipment Items</span>
                </div>
            </div>
            
            <div class="dashboard-card">
                <div class="card-header">
                    <h2 class="card-title">Recent Workout Instances</h2>
                </div>
                <div id="recent-instances">Loading...</div>
            </div>
        `;

    this.setContent(html);
    await this.loadDashboardStats();
  }

  async loadDashboardStats() {
    try {
      // Load stats in parallel
      const [members, workouts, exercises, equipment, instances] =
        await Promise.all([
          this.apiCall(`/user`),
          this.apiCall(`/custom-workout-template`),
          this.apiCall(`/custom-exercise`),
          this.apiCall(`/custom-equipment`),
          this.apiCall(`/custom-workout-instance`),
        ]);

      const [
        membersData,
        workoutsData,
        exercisesData,
        equipmentData,
        instancesData,
      ] = await Promise.all([
        members.json(),
        workouts.json(),
        exercises.json(),
        equipment.json(),
        instances.json(),
      ]);

      // Update stats
      document.getElementById("total-members").textContent =
        membersData.length || 0;
      document.getElementById("total-workouts").textContent =
        workoutsData.length || 0;
      document.getElementById("total-exercises").textContent =
        exercisesData.length || 0;
      document.getElementById("total-equipment").textContent =
        equipmentData.length || 0;

      // Update recent instances
      this.renderRecentInstances(instancesData);
    } catch (error) {
      console.error("Failed to load dashboard stats:", error);
    }
  }

  async loadCustomExercises() {
    if (!this.currentGym) {
      this.showError("No gym selected");
      return;
    }

    try {
      const response = await this.apiCall(
        `/gyms/${this.currentGym}/custom-exercises`
      );
      const exercises = await response.json();

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

  // Placeholder methods for other views
  async loadMembers() {
    this.showComingSoon("Members management");
  }

  async loadCustomEquipment() {
    this.showComingSoon("Equipment management");
  }

  async loadWorkoutTemplates() {
    this.showComingSoon("Workout templates");
  }

  async loadWorkoutInstances() {
    this.showComingSoon("Workout instances");
  }

  async loadAnalytics() {
    this.showComingSoon("Analytics dashboard");
  }

  async loadGymsManagement() {
    this.showComingSoon("Gyms management");
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
    document.getElementById("content-body").innerHTML = html;
  }

  showLoading() {
    this.setContent(`
            <div class="loading-state">
                <i class="fas fa-spinner fa-spin"></i>
                <span>Loading...</span>
            </div>
        `);
  }

  showError(message) {
    this.setContent(`
            <div class="empty-state">
                <i class="fas fa-exclamation-triangle"></i>
                <h3>Error</h3>
                <p>${message}</p>
            </div>
        `);
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
