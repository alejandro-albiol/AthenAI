/**
 * API Utility Module
 * Centralized API communication with error handling and authentication
 */
class ApiClient {
  constructor(baseUrl = "/api/v1") {
    this.baseUrl = baseUrl;
    this.defaultHeaders = {
      "Content-Type": "application/json",
    };
  }

  getAuthToken() {
    // Use authManager if available, fallback to direct localStorage access
    if (window.authManager) {
      return window.authManager.getToken();
    }
    return (
      localStorage.getItem("auth_token") || sessionStorage.getItem("auth_token")
    );
  }

  getHeaders(additionalHeaders = {}) {
    const headers = { ...this.defaultHeaders, ...additionalHeaders };

    const token = this.getAuthToken();
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }

    return headers;
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseUrl}${endpoint}`;
    const config = {
      headers: this.getHeaders(options.headers),
      ...options,
    };

    try {
      const response = await fetch(url, config);

      if (!response.ok) {
        // Handle 401 errors with automatic token refresh
        if (
          response.status === 401 &&
          window.authManager &&
          !endpoint.includes("/auth/")
        ) {
          try {
            console.log("Access token expired, attempting refresh...");
            await window.authManager.handleTokenRefresh();

            // Retry the original request with new token
            const retryHeaders = {
              "Content-Type": "application/json",
              ...headers,
            };

            const newToken = window.authManager.getToken();
            if (newToken) {
              retryHeaders.Authorization = `Bearer ${newToken}`;
            }

            const retryResponse = await fetch(url, {
              ...options,
              headers: retryHeaders,
            });

            if (retryResponse.ok) {
              const contentType = retryResponse.headers.get("content-type");
              if (contentType && contentType.includes("application/json")) {
                return await retryResponse.json();
              }
              return retryResponse;
            }
          } catch (refreshError) {
            console.error("Token refresh failed:", refreshError);
            // Clear auth and redirect to login
            if (window.authManager) {
              window.authManager.clearAuth();
            }
            window.location.href = "/";
            return;
          }
        }

        throw new ApiError(
          `HTTP ${response.status}: ${response.statusText}`,
          response.status,
          await this.getErrorBody(response)
        );
      }

      const contentType = response.headers.get("content-type");
      if (contentType && contentType.includes("application/json")) {
        return await response.json();
      }

      return response;
    } catch (error) {
      if (error instanceof ApiError) {
        throw error;
      }

      throw new ApiError(error.message || "Network error occurred", 0, null);
    }
  }

  async getErrorBody(response) {
    try {
      const contentType = response.headers.get("content-type");
      if (contentType && contentType.includes("application/json")) {
        return await response.json();
      }
      return await response.text();
    } catch {
      return null;
    }
  }

  // HTTP Methods
  async get(endpoint, params = {}) {
    let url = endpoint;
    if (Object.keys(params).length > 0) {
      const searchParams = new URLSearchParams();
      Object.keys(params).forEach((key) => {
        if (params[key] !== undefined && params[key] !== null) {
          searchParams.append(key, params[key]);
        }
      });
      url += `?${searchParams.toString()}`;
    }

    return this.request(url, {
      method: "GET",
    });
  }

  async post(endpoint, data = null) {
    return this.request(endpoint, {
      method: "POST",
      body: data ? JSON.stringify(data) : null,
    });
  }

  async put(endpoint, data = null) {
    return this.request(endpoint, {
      method: "PUT",
      body: data ? JSON.stringify(data) : null,
    });
  }

  async patch(endpoint, data = null) {
    return this.request(endpoint, {
      method: "PATCH",
      body: data ? JSON.stringify(data) : null,
    });
  }

  async delete(endpoint) {
    return this.request(endpoint, {
      method: "DELETE",
    });
  }

  // Domain-specific API methods
  async getEquipment() {
    return this.get("/equipment");
  }

  async createEquipment(data) {
    return this.post("/equipment", data);
  }

  async updateEquipment(id, data) {
    return this.put(`/equipment/${id}`, data);
  }

  async deleteEquipment(id) {
    return this.delete(`/equipment/${id}`);
  }

  async restoreEquipment(id) {
    return this.post(`/equipment/${id}/restore`);
  }

  async getExercises() {
    return this.get("/exercise");
  }

  async createExercise(data) {
    return this.post("/exercise", data);
  }

  async updateExercise(id, data) {
    return this.put(`/exercise/${id}`, data);
  }

  async deleteExercise(id) {
    return this.delete(`/exercise/${id}`);
  }

  async restoreExercise(id) {
    return this.post(`/exercise/${id}/restore`);
  }

  async getMuscularGroups() {
    return this.get("/muscular-group");
  }

  async createMuscularGroup(data) {
    return this.post("/muscular-group", data);
  }

  async updateMuscularGroup(id, data) {
    return this.put(`/muscular-group/${id}`, data);
  }

  async deleteMuscularGroup(id) {
    return this.delete(`/muscular-group/${id}`);
  }

  // Exercise-Muscular Group linking
  async createExerciseMuscularGroupLink(data) {
    return this.post("/exercise-muscular-groups/link", data);
  }

  async deleteExerciseMuscularGroupLink(id) {
    return this.delete(`/exercise-muscular-groups/link/${id}`);
  }

  async getExerciseMuscularGroupLinks(exerciseId) {
    return this.get(`/exercise-muscular-groups/exercise/${exerciseId}/links`);
  }

  // Exercise-Equipment linking
  async createExerciseEquipmentLink(data) {
    return this.post("/exercise-equipment/link", data);
  }

  async deleteExerciseEquipmentLink(id) {
    return this.delete(`/exercise-equipment/link/${id}`);
  }

  async getExerciseEquipmentLinks(exerciseId) {
    return this.get(`/exercise-equipment/exercise/${exerciseId}/links`);
  }

  async getGyms() {
    return this.get("/gym");
  }

  async createGym(data) {
    return this.post("/gym", data);
  }

  async updateGym(id, data) {
    return this.put(`/gym/${id}`, data);
  }

  async deleteGym(id) {
    return this.delete(`/gym/${id}`);
  }

  async restoreGym(id) {
    return this.patch(`/gym/${id}/restore`);
  }

  async getPlatformStats() {
    return this.get("/platform/stats");
  }

  // Authentication methods
  async login(email, password, inviteToken = null, rememberMe = false) {
    const credentials = {
      email,
      password,
    };

    if (inviteToken) {
      credentials.invite_token = inviteToken;
    }

    // Note: remember_me handled client-side only for now
    // TODO: Add remember_me to backend when DTO is updated

    return this.post("/auth/login", credentials);
  }

  async logout() {
    return this.post("/auth/logout");
  }

  async getCurrentUser() {
    return this.get("/auth/validate");
  }

  // User Management Endpoints
  async getUsers(gymId = null) {
    // If gymId is provided, use platform admin endpoint to get users for specific gym
    if (gymId) {
      return this.get(`/user/gym/${gymId}`);
    }
    // Otherwise, get users for current authenticated gym context
    return this.get("/user");
  }

  async getAllUsers() {
    // For platform admins - get users across all gyms
    return this.get("/user");
  }

  async getUserById(userId, gymId = null) {
    if (gymId) {
      return this.get(`/user/gym/${gymId}/${userId}`);
    }
    return this.get(`/user/${userId}`);
  }

  async createUser(userData, gymId = null) {
    if (gymId) {
      return this.post(`/user/gym/${gymId}`, userData);
    }
    return this.post("/user", userData);
  }

  async updateUser(userId, userData, gymId = null) {
    if (gymId) {
      return this.put(`/user/gym/${gymId}/${userId}`, userData);
    }
    return this.put(`/user/${userId}`, userData);
  }

  async deleteUser(userId, gymId = null) {
    if (gymId) {
      return this.delete(`/user/gym/${gymId}/${userId}`);
    }
    return this.delete(`/user/${userId}`);
  }

  async updateUserPassword(userId, passwordData, gymId = null) {
    if (gymId) {
      return this.put(`/user/gym/${gymId}/${userId}/password`, passwordData);
    }
    return this.put(`/user/${userId}/password`, passwordData);
  }

  async verifyUser(userId, gymId = null) {
    if (gymId) {
      return this.post(`/user/gym/${gymId}/${userId}/verify`);
    }
    return this.post(`/user/${userId}/verify`);
  }

  async setUserActive(userId, active, gymId = null) {
    if (gymId) {
      return this.post(`/user/gym/${gymId}/${userId}/active`, { active });
    }
    return this.post(`/user/${userId}/active`, { active });
  }

  // Invitation management methods
  async createInvitation(data) {
    return this.post("/invitation", data);
  }

  async getGymInvitations(gymId) {
    return this.get(`/gym/${gymId}/invitations`);
  }

  async decodeInvitation(token) {
    return this.get(`/invitation/decode/${token}`);
  }

  async acceptInvitation(token, data) {
    return this.post(`/invitation/accept/${token}`, data);
  }

  async resendInvitation(invitationId) {
    return this.post(`/invitation/${invitationId}/resend`);
  }

  async deleteInvitation(invitationId) {
    return this.delete(`/invitation/${invitationId}`);
  }
}

/**
 * API Error class for better error handling
 */
class ApiError extends Error {
  constructor(message, status, body) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.body = body;
  }

  isAuthError() {
    return this.status === 401 || this.status === 403;
  }

  isValidationError() {
    return this.status === 400 || this.status === 422;
  }

  isServerError() {
    return this.status >= 500;
  }

  getValidationErrors() {
    if (this.isValidationError() && this.body && this.body.errors) {
      return this.body.errors;
    }
    return {};
  }
}

// Make available globally for vanilla JS
window.ApiClient = ApiClient;
window.ApiError = ApiError;
