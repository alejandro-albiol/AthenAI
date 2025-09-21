/**
 * Authentication Utility Functions
 * Handles token storage, retrieval, and automatic refresh
 * Vanilla JS - no imports, works directly in browsers
 */

class AuthManager {
  constructor() {
    this.isRefreshing = false;
    this.refreshSubscribers = [];
  }

  /**
   * Store authentication tokens using appropriate storage strategy
   */
  storeTokens(accessToken, refreshToken, userInfo, rememberMe = false) {
    const storage = rememberMe ? localStorage : sessionStorage;

    storage.setItem("auth_token", accessToken);
    storage.setItem("refresh_token", refreshToken);
    storage.setItem("user_info", JSON.stringify(userInfo));

    // Always store remember preference in localStorage for consistency
    localStorage.setItem("remember_me", rememberMe ? "true" : "false");
  }

  /**
   * Get stored authentication token
   */
  getToken() {
    // Try localStorage first (remember me), then sessionStorage
    return (
      localStorage.getItem("auth_token") || sessionStorage.getItem("auth_token")
    );
  }

  /**
   * Get stored refresh token
   */
  getRefreshToken() {
    return (
      localStorage.getItem("refresh_token") ||
      sessionStorage.getItem("refresh_token")
    );
  }

  /**
   * Get stored user info
   */
  getUserInfo() {
    const userInfoStr =
      localStorage.getItem("user_info") || sessionStorage.getItem("user_info");
    try {
      return userInfoStr ? JSON.parse(userInfoStr) : null;
    } catch {
      return null;
    }
  }

  /**
   * Check if user chose to be remembered
   */
  isRemembered() {
    return localStorage.getItem("remember_me") === "true";
  }

  /**
   * Check if user is authenticated
   */
  isAuthenticated() {
    return !!this.getToken();
  }

  /**
   * Clear all authentication data
   */
  clearAuth() {
    // Clear from both storages
    localStorage.removeItem("auth_token");
    localStorage.removeItem("refresh_token");
    localStorage.removeItem("user_info");
    localStorage.removeItem("remember_me");

    sessionStorage.removeItem("auth_token");
    sessionStorage.removeItem("refresh_token");
    sessionStorage.removeItem("user_info");
  }

  /**
   * Refresh access token using refresh token
   */
  async refreshAccessToken() {
    const refreshToken = this.getRefreshToken();
    if (!refreshToken) {
      throw new Error("No refresh token available");
    }

    try {
      const api = new ApiClient();
      const response = await api.post("/auth/refresh", {
        refresh_token: refreshToken,
      });

      if (response.status === "success" && response.data) {
        // Update tokens using same storage strategy
        const rememberMe = this.isRemembered();
        this.storeTokens(
          response.data.access_token,
          response.data.refresh_token || refreshToken, // Use new refresh token if provided
          response.data.user_info || this.getUserInfo(),
          rememberMe
        );

        return response.data.access_token;
      } else {
        throw new Error(response.message || "Token refresh failed");
      }
    } catch (error) {
      // If refresh fails, clear auth and redirect to login
      this.clearAuth();
      throw error;
    }
  }

  /**
   * Add subscriber for token refresh events
   */
  subscribeToTokenRefresh(callback) {
    this.refreshSubscribers.push(callback);
  }

  /**
   * Notify all subscribers of token refresh
   */
  notifyTokenRefresh(token) {
    this.refreshSubscribers.forEach((callback) => callback(token));
  }

  /**
   * Handle token refresh with deduplication
   */
  async handleTokenRefresh() {
    if (this.isRefreshing) {
      // Wait for ongoing refresh
      return new Promise((resolve) => {
        this.subscribeToTokenRefresh(resolve);
      });
    }

    this.isRefreshing = true;

    try {
      const newToken = await this.refreshAccessToken();
      this.notifyTokenRefresh(newToken);
      return newToken;
    } catch (error) {
      this.notifyTokenRefresh(null);
      throw error;
    } finally {
      this.isRefreshing = false;
      this.refreshSubscribers = [];
    }
  }

  /**
   * Setup automatic token refresh
   */
  setupTokenRefresh() {
    // Check token expiry every 5 minutes
    setInterval(() => {
      if (this.isAuthenticated()) {
        // In a real implementation, you'd decode JWT to check expiry
        // For now, we'll rely on API calls to trigger refresh when needed
      }
    }, 5 * 60 * 1000); // 5 minutes
  }
}

// Create global instance
const authManager = new AuthManager();

// Make available globally for vanilla JS
window.AuthManager = AuthManager;
window.authManager = authManager;
