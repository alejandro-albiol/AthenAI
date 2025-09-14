// API Configuration
const API_BASE_URL = "/api/v1";

// API Helper Class
class API {
  constructor() {
    this.baseURL = API_BASE_URL;
  }

  // Helper method to get auth headers
  getAuthHeaders() {
    // Check for both token storage formats for compatibility
    const token =
      localStorage.getItem("accessToken") || localStorage.getItem("auth_token");
    const gymId = localStorage.getItem("gymId");

    const headers = {
      "Content-Type": "application/json",
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    if (gymId) {
      headers["X-Gym-ID"] = gymId;
    }

    return headers;
  }

  // Generic request method
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      headers: this.getAuthHeaders(),
      ...options,
    };

    try {
      const response = await fetch(url, config);

      // Handle token refresh if needed
      if (response.status === 401) {
        const refreshed = await this.refreshToken();
        if (refreshed) {
          // Retry the request with new token
          config.headers = this.getAuthHeaders();
          return await fetch(url, config);
        } else {
          // Redirect to login if refresh fails
          window.location.href = "/";
          return;
        }
      }

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const contentType = response.headers.get("content-type");
      if (contentType && contentType.includes("application/json")) {
        return await response.json();
      }
      return response;
    } catch (error) {
      console.error("API request failed:", error);
      throw error;
    }
  }

  // Authentication endpoints
  async login(email, password, gymId = null) {
    const headers = { "Content-Type": "application/json" };
    if (gymId) {
      headers["X-Gym-ID"] = gymId;
    }

    return await fetch(`${this.baseURL}/auth/login`, {
      method: "POST",
      headers,
      body: JSON.stringify({ email, password }),
    });
  }

  async refreshToken() {
    const refreshToken = localStorage.getItem("refreshToken");
    if (!refreshToken) return false;

    try {
      const response = await fetch(`${this.baseURL}/auth/refresh`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (response.ok) {
        const data = await response.json();
        localStorage.setItem("accessToken", data.access_token);
        if (data.refresh_token) {
          localStorage.setItem("refreshToken", data.refresh_token);
        }
        return true;
      }
    } catch (error) {
      console.error("Token refresh failed:", error);
    }

    return false;
  }

  async logout() {
    try {
      await this.request("/auth/logout", { method: "POST" });
    } catch (error) {
      console.error("Logout request failed:", error);
    } finally {
      // Clear local storage regardless of API call success
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("gymId");
      localStorage.removeItem("userRole");
      window.location.href = "/";
    }
  }

  async validateToken() {
    return await this.request("/auth/validate");
  }

  // Get current user information
  async getMe() {
    return await fetch(`${this.baseURL}/auth/validate`, {
      method: "GET",
      headers: this.getAuthHeaders(),
    });
  }

  // Set authentication token
  setToken(token) {
    localStorage.setItem("auth_token", token);
  }

  // Clear authentication token
  clearToken() {
    localStorage.removeItem("auth_token");
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    localStorage.removeItem("gymId");
    localStorage.removeItem("userRole");
  }

  // Gym Management (Platform Admin only)
  async getGyms() {
    return await this.request("/gym");
  }

  async getGymById(id) {
    return await this.request(`/gym/${id}`);
  }

  async getGymByName(name) {
    return await this.request(`/gym/name/${name}`);
  }

  async createGym(gymData) {
    return await this.request("/gym", {
      method: "POST",
      body: JSON.stringify(gymData),
    });
  }

  async updateGym(id, gymData) {
    return await this.request(`/gym/${id}`, {
      method: "PUT",
      body: JSON.stringify(gymData),
    });
  }

  async activateGym(id) {
    return await this.request(`/gym/${id}/activate`, {
      method: "PUT",
    });
  }

  async deactivateGym(id) {
    return await this.request(`/gym/${id}/deactivate`, {
      method: "PUT",
    });
  }

  async deleteGym(id) {
    return await this.request(`/gym/${id}`, {
      method: "DELETE",
    });
  }

  // User Management
  async getUsers() {
    return await this.request("/user");
  }

  async getUserById(id) {
    return await this.request(`/user/${id}`);
  }

  async createUser(userData) {
    return await this.request("/user", {
      method: "POST",
      body: JSON.stringify(userData),
    });
  }

  async updateUser(id, userData) {
    return await this.request(`/user/${id}`, {
      method: "PUT",
      body: JSON.stringify(userData),
    });
  }

  async deleteUser(id) {
    return await this.request(`/user/${id}`, {
      method: "DELETE",
    });
  }

  // Public Exercise Management
  async getExercises() {
    return await this.request("/exercise");
  }

  async getExerciseById(id) {
    return await this.request(`/exercise/${id}`);
  }

  // Equipment Management
  async getEquipment() {
    return await this.request("/equipment");
  }

  async getEquipmentById(id) {
    return await this.request(`/equipment/${id}`);
  }

  // Muscular Groups
  async getMuscularGroups() {
    return await this.request("/muscular-group");
  }

  // Template Blocks
  async getTemplateBlocks() {
    return await this.request("/template-block");
  }

  // Workout Templates
  async getWorkoutTemplates() {
    return await this.request("/workout-template");
  }

  // Custom Exercise Management (Tenant-specific)
  async getCustomExercises() {
    return await this.request("/custom-exercise");
  }

  async getCustomExerciseById(id) {
    return await this.request(`/custom-exercise/${id}`);
  }

  async createCustomExercise(exerciseData) {
    return await this.request("/custom-exercise", {
      method: "POST",
      body: JSON.stringify(exerciseData),
    });
  }

  async updateCustomExercise(id, exerciseData) {
    return await this.request(`/custom-exercise/${id}`, {
      method: "PUT",
      body: JSON.stringify(exerciseData),
    });
  }

  async deleteCustomExercise(id) {
    return await this.request(`/custom-exercise/${id}`, {
      method: "DELETE",
    });
  }

  // Custom Equipment Management (Tenant-specific)
  async getCustomEquipment() {
    return await this.request("/custom-equipment");
  }

  async getCustomEquipmentById(id) {
    return await this.request(`/custom-equipment/${id}`);
  }

  async createCustomEquipment(equipmentData) {
    return await this.request("/custom-equipment", {
      method: "POST",
      body: JSON.stringify(equipmentData),
    });
  }

  async updateCustomEquipment(id, equipmentData) {
    return await this.request(`/custom-equipment/${id}`, {
      method: "PUT",
      body: JSON.stringify(equipmentData),
    });
  }

  async deleteCustomEquipment(id) {
    return await this.request(`/custom-equipment/${id}`, {
      method: "DELETE",
    });
  }

  // Custom Workout Instance Management (Tenant-specific)
  async getCustomWorkoutInstances() {
    return await this.request("/custom-workout-instance");
  }

  async getCustomWorkoutInstanceById(id) {
    return await this.request(`/custom-workout-instance/${id}`);
  }

  async createCustomWorkoutInstance(instanceData) {
    return await this.request("/custom-workout-instance", {
      method: "POST",
      body: JSON.stringify(instanceData),
    });
  }

  async updateCustomWorkoutInstance(id, instanceData) {
    return await this.request(`/custom-workout-instance/${id}`, {
      method: "PUT",
      body: JSON.stringify(instanceData),
    });
  }

  async deleteCustomWorkoutInstance(id) {
    return await this.request(`/custom-workout-instance/${id}`, {
      method: "DELETE",
    });
  }

  // Additional methods for relations and specific endpoints can be added here
}

// Export API instance
const api = new API();
