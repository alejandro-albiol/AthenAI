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
    return localStorage.getItem("auth_token");
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
    const url = new URL(`${this.baseUrl}${endpoint}`, window.location.origin);
    Object.keys(params).forEach((key) => {
      if (params[key] !== undefined && params[key] !== null) {
        url.searchParams.append(key, params[key]);
      }
    });

    return this.request(url.pathname + url.search, {
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

  async getExercises() {
    return this.get("/exercises");
  }

  async createExercise(data) {
    return this.post("/exercises", data);
  }

  async updateExercise(id, data) {
    return this.put(`/exercises/${id}`, data);
  }

  async deleteExercise(id) {
    return this.delete(`/exercises/${id}`);
  }

  async getMuscularGroups() {
    return this.get("/muscular-groups");
  }

  async createMuscularGroup(data) {
    return this.post("/muscular-groups", data);
  }

  async updateMuscularGroup(id, data) {
    return this.put(`/muscular-groups/${id}`, data);
  }

  async deleteMuscularGroup(id) {
    return this.delete(`/muscular-groups/${id}`);
  }

  async getGyms() {
    return this.get("/gyms");
  }

  async createGym(data) {
    return this.post("/gyms", data);
  }

  async updateGym(id, data) {
    return this.put(`/gyms/${id}`, data);
  }

  async deleteGym(id) {
    return this.delete(`/gyms/${id}`);
  }

  async restoreGym(id) {
    return this.patch(`/gyms/${id}/restore`);
  }

  async getPlatformStats() {
    return this.get("/platform/stats");
  }

  // Authentication methods
  async login(email, password, inviteToken = null) {
    const credentials = {
      email,
      password,
    };

    if (inviteToken) {
      credentials.invite_token = inviteToken;
    }

    return this.post("/auth/login", credentials);
  }

  async logout() {
    return this.post("/auth/logout");
  }

  async getCurrentUser() {
    return this.get("/auth/me");
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
