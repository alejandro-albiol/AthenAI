// DOM Content Loaded
document.addEventListener("DOMContentLoaded", function () {
  // Initialize the application
  initializeApp();
});

// Initialize application
function initializeApp() {
  // Setup form handlers
  setupFormHandlers();

  // Setup modal event listeners
  setupModalEventListeners();

  // Setup smooth scrolling
  setupSmoothScrolling();

  console.log("AthenAI Frontend initialized");
}

// Modal Management
function openRegisterModal() {
  const modal = document.getElementById("registerModal");
  modal.style.display = "block";
  document.body.style.overflow = "hidden";
}

function closeRegisterModal() {
  const modal = document.getElementById("registerModal");
  modal.style.display = "none";
  document.body.style.overflow = "auto";

  // Reset form
  document.getElementById("registerForm").reset();
}

function openLoginModal() {
  const modal = document.getElementById("loginModal");
  modal.style.display = "block";
  document.body.style.overflow = "hidden";
}

function closeLoginModal() {
  const modal = document.getElementById("loginModal");
  modal.style.display = "none";
  document.body.style.overflow = "auto";

  // Reset form
  document.getElementById("loginForm").reset();
}

// Setup modal event listeners
function setupModalEventListeners() {
  // Close modals when clicking outside
  window.onclick = function (event) {
    const registerModal = document.getElementById("registerModal");
    const loginModal = document.getElementById("loginModal");

    if (event.target === registerModal) {
      closeRegisterModal();
    }

    if (event.target === loginModal) {
      closeLoginModal();
    }
  };

  // Close modals on escape key
  document.addEventListener("keydown", function (event) {
    if (event.key === "Escape") {
      closeRegisterModal();
      closeLoginModal();
    }
  });
}

// Form Handlers
function setupFormHandlers() {
  // Register form handler
  const registerForm = document.getElementById("registerForm");
  registerForm.addEventListener("submit", handleGymRegistration);

  // Login form handler
  const loginForm = document.getElementById("loginForm");
  loginForm.addEventListener("submit", handleLogin);
}

// Handle gym registration
async function handleGymRegistration(event) {
  event.preventDefault();

  const formData = new FormData(event.target);
  const gymData = {
    name: formData.get("gymName"),
    owner_name: formData.get("ownerName"),
    email: formData.get("email"),
    phone: formData.get("phone"),
    address: formData.get("address"),
    domain: formData.get("domain"),
  };

  try {
    showLoading("Creating your gym account...");

    const response = await fetch("/api/v1/gym", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(gymData),
    });

    hideLoading();

    if (response.ok) {
      const result = await response.json();
      showSuccess("Gym registered successfully! Welcome to AthenAI!");
      closeRegisterModal();

      // Redirect to gym dashboard (you can implement this later)
      setTimeout(() => {
        window.location.href = `/gym/${result.id}/dashboard`;
      }, 2000);
    } else {
      const error = await response.json();
      showError(error.message || "Failed to register gym. Please try again.");
    }
  } catch (error) {
    hideLoading();
    showError("Network error. Please check your connection and try again.");
    console.error("Registration error:", error);
  }
}

// Handle login
async function handleLogin(event) {
  event.preventDefault();

  const formData = new FormData(event.target);
  const loginData = {
    email: formData.get("email"),
    password: formData.get("password"),
  };

  try {
    showLoading("Signing you in...");

    const response = await fetch("/api/v1/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(loginData),
    });

    hideLoading();

    if (response.ok) {
      const result = await response.json();
      showSuccess("Welcome back!");
      closeLoginModal();

      // Store auth token (implement proper token management)
      localStorage.setItem("auth_token", result.data.access_token);

      // Redirect to dashboard
      setTimeout(() => {
        window.location.href = "/pages/dashboard/";
      }, 1000);
    } else {
      const error = await response.json();
      showError(error.message || "Invalid credentials. Please try again.");
    }
  } catch (error) {
    hideLoading();
    showError("Network error. Please check your connection and try again.");
    console.error("Login error:", error);
  }
}

// Schedule demo functionality
function scheduleDemo() {
  showInfo(
    "Demo scheduling feature coming soon! Please contact us at demo@athenai.com"
  );
}

// Notification System
function showSuccess(message) {
  showNotification(message, "success");
}

function showError(message) {
  showNotification(message, "error");
}

function showInfo(message) {
  showNotification(message, "info");
}

function showNotification(message, type) {
  // Remove existing notifications
  const existingNotification = document.querySelector(".notification");
  if (existingNotification) {
    existingNotification.remove();
  }

  // Create notification element
  const notification = document.createElement("div");
  notification.className = `notification notification-${type}`;
  notification.innerHTML = `
        <div class="notification-content">
            <i class="fas fa-${getIconForType(type)}"></i>
            <span>${message}</span>
            <button class="notification-close" onclick="this.parentElement.parentElement.remove()">
                <i class="fas fa-times"></i>
            </button>
        </div>
    `;

  // Add styles
  notification.style.cssText = `
        position: fixed;
        top: 90px;
        right: 20px;
        z-index: 3000;
        min-width: 300px;
        max-width: 500px;
        padding: 16px;
        border-radius: 8px;
        color: white;
        font-weight: 500;
        box-shadow: 0 10px 25px rgba(0,0,0,0.2);
        animation: slideInRight 0.3s ease-out;
        background: ${getColorForType(type)};
    `;

  // Add to page
  document.body.appendChild(notification);

  // Auto remove after 5 seconds
  setTimeout(() => {
    if (notification.parentElement) {
      notification.style.animation = "slideOutRight 0.3s ease-in";
      setTimeout(() => {
        if (notification.parentElement) {
          notification.remove();
        }
      }, 300);
    }
  }, 5000);
}

function getIconForType(type) {
  switch (type) {
    case "success":
      return "check-circle";
    case "error":
      return "exclamation-circle";
    case "info":
      return "info-circle";
    default:
      return "bell";
  }
}

function getColorForType(type) {
  switch (type) {
    case "success":
      return "#10b981";
    case "error":
      return "#ef4444";
    case "info":
      return "#3b82f6";
    default:
      return "#6b7280";
  }
}

// Loading overlay
function showLoading(message = "Loading...") {
  // Remove existing loading overlay
  hideLoading();

  const overlay = document.createElement("div");
  overlay.id = "loadingOverlay";
  overlay.innerHTML = `
        <div class="loading-content">
            <div class="spinner"></div>
            <p>${message}</p>
        </div>
    `;

  overlay.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0,0,0,0.7);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 4000;
        backdrop-filter: blur(5px);
    `;

  // Add spinner styles
  const style = document.createElement("style");
  style.textContent = `
        .loading-content {
            text-align: center;
            color: white;
        }
        .spinner {
            width: 40px;
            height: 40px;
            border: 4px solid rgba(255,255,255,0.3);
            border-top: 4px solid white;
            border-radius: 50%;
            animation: spin 1s linear infinite;
            margin: 0 auto 16px;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        @keyframes slideInRight {
            0% { transform: translateX(100%); opacity: 0; }
            100% { transform: translateX(0); opacity: 1; }
        }
        @keyframes slideOutRight {
            0% { transform: translateX(0); opacity: 1; }
            100% { transform: translateX(100%); opacity: 0; }
        }
    `;

  document.head.appendChild(style);
  document.body.appendChild(overlay);
}

function hideLoading() {
  const overlay = document.getElementById("loadingOverlay");
  if (overlay) {
    overlay.remove();
  }
}

// Smooth scrolling for navigation links
function setupSmoothScrolling() {
  document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
    anchor.addEventListener("click", function (e) {
      e.preventDefault();
      const target = document.querySelector(this.getAttribute("href"));
      if (target) {
        target.scrollIntoView({
          behavior: "smooth",
          block: "start",
        });
      }
    });
  });
}

// Form validation helpers
function validateEmail(email) {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return re.test(email);
}

function validateForm(formData) {
  const errors = [];

  if (!formData.get("gymName") || formData.get("gymName").trim().length < 2) {
    errors.push("Gym name must be at least 2 characters long");
  }

  if (
    !formData.get("ownerName") ||
    formData.get("ownerName").trim().length < 2
  ) {
    errors.push("Owner name must be at least 2 characters long");
  }

  if (!formData.get("email") || !validateEmail(formData.get("email"))) {
    errors.push("Please enter a valid email address");
  }

  return errors;
}

// API Helper functions
class ApiClient {
  constructor() {
    this.baseURL = "/api/v1";
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    };

    // Add auth token if available
    const token = localStorage.getItem("auth_token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    const response = await fetch(url, config);
    return response;
  }

  async get(endpoint) {
    return this.request(endpoint, { method: "GET" });
  }

  async post(endpoint, data) {
    return this.request(endpoint, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async put(endpoint, data) {
    return this.request(endpoint, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async delete(endpoint) {
    return this.request(endpoint, { method: "DELETE" });
  }
}

// Create global API client instance
window.apiClient = new ApiClient();
