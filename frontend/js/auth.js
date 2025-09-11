// Authentication Management
class AuthManager {
  constructor() {
    this.currentUser = null;
    this.token = localStorage.getItem("auth_token");
  }

  // Check if user is authenticated
  isAuthenticated() {
    return !!this.token;
  }

  // Get current user
  getCurrentUser() {
    return this.currentUser;
  }

  // Login function
  async login(email, password) {
    try {
      const response = await api.login(email, password);
      const data = await response.json();

      if (response.ok) {
        this.token = data.token;
        this.currentUser = data.user;
        api.setToken(this.token);

        return { success: true, user: this.currentUser };
      } else {
        return { success: false, error: data.error || "Login failed" };
      }
    } catch (error) {
      console.error("Login error:", error);
      return { success: false, error: "Network error. Please try again." };
    }
  }

  // Register function
  async register(gymData) {
    try {
      const response = await api.register(gymData);
      const data = await response.json();

      if (response.ok) {
        this.token = data.token;
        this.currentUser = data.user;
        api.setToken(this.token);

        return { success: true, user: this.currentUser };
      } else {
        return { success: false, error: data.error || "Registration failed" };
      }
    } catch (error) {
      console.error("Registration error:", error);
      return { success: false, error: "Network error. Please try again." };
    }
  }

  // Logout function
  logout() {
    this.token = null;
    this.currentUser = null;
    api.clearToken();
    window.location.href = "/";
  }

  // Verify current session
  async verifySession() {
    if (!this.token) {
      return false;
    }

    try {
      const response = await api.getMe();

      if (response.ok) {
        this.currentUser = await response.json();
        return true;
      } else {
        this.logout();
        return false;
      }
    } catch (error) {
      console.error("Session verification failed:", error);
      this.logout();
      return false;
    }
  }

  // Redirect based on authentication status
  async redirectIfNeeded() {
    const isAuth = await this.verifySession();
    const currentPath = window.location.pathname;

    if (isAuth && currentPath === "/") {
      // User is authenticated and on landing page, redirect to dashboard
      window.location.href = "/pages/dashboard/";
    } else if (!isAuth && currentPath.includes("/pages/dashboard/")) {
      // User is not authenticated and trying to access dashboard
      window.location.href = "/";
    }
  }
}

// Form validation helpers
function validateEmail(email) {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return re.test(email);
}

function validatePassword(password) {
  return password.length >= 8;
}

function validateGymName(name) {
  return name.trim().length >= 2;
}

// Form submission handlers
async function handleLogin(event) {
  event.preventDefault();

  const form = event.target;
  const email = form.email.value.trim();
  const password = form.password.value;

  // Clear previous errors
  clearFormErrors(form);

  // Validate inputs
  const errors = [];
  if (!validateEmail(email)) {
    errors.push({
      field: "email",
      message: "Please enter a valid email address",
    });
  }
  if (!validatePassword(password)) {
    errors.push({
      field: "password",
      message: "Password must be at least 8 characters",
    });
  }

  if (errors.length > 0) {
    showFormErrors(form, errors);
    return;
  }

  // Show loading state
  const submitBtn = form.querySelector('button[type="submit"]');
  const originalText = submitBtn.innerHTML;
  submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Signing in...';
  submitBtn.disabled = true;

  try {
    const result = await auth.login(email, password);

    if (result.success) {
      // Redirect to dashboard
      window.location.href = "/pages/dashboard/";
    } else {
      showFormErrors(form, [{ field: "general", message: result.error }]);
    }
  } catch (error) {
    showFormErrors(form, [
      { field: "general", message: "An unexpected error occurred" },
    ]);
  } finally {
    // Reset button
    submitBtn.innerHTML = originalText;
    submitBtn.disabled = false;
  }
}

async function handleRegister(event) {
  event.preventDefault();

  const form = event.target;
  const gymData = {
    gym_name: form.gym_name.value.trim(),
    admin_name: form.admin_name.value.trim(),
    admin_email: form.admin_email.value.trim(),
    admin_password: form.admin_password.value,
    subscription_type: form.subscription_type.value,
  };

  // Clear previous errors
  clearFormErrors(form);

  // Validate inputs
  const errors = [];
  if (!validateGymName(gymData.gym_name)) {
    errors.push({
      field: "gym_name",
      message: "Gym name must be at least 2 characters",
    });
  }
  if (!gymData.admin_name.trim()) {
    errors.push({ field: "admin_name", message: "Admin name is required" });
  }
  if (!validateEmail(gymData.admin_email)) {
    errors.push({
      field: "admin_email",
      message: "Please enter a valid email address",
    });
  }
  if (!validatePassword(gymData.admin_password)) {
    errors.push({
      field: "admin_password",
      message: "Password must be at least 8 characters",
    });
  }

  if (errors.length > 0) {
    showFormErrors(form, errors);
    return;
  }

  // Show loading state
  const submitBtn = form.querySelector('button[type="submit"]');
  const originalText = submitBtn.innerHTML;
  submitBtn.innerHTML =
    '<i class="fas fa-spinner fa-spin"></i> Creating account...';
  submitBtn.disabled = true;

  try {
    const result = await auth.register(gymData);

    if (result.success) {
      // Show success message and redirect
      showSuccessMessage(
        "Account created successfully! Redirecting to dashboard..."
      );
      setTimeout(() => {
        window.location.href = "/pages/dashboard/";
      }, 2000);
    } else {
      showFormErrors(form, [{ field: "general", message: result.error }]);
    }
  } catch (error) {
    showFormErrors(form, [
      { field: "general", message: "An unexpected error occurred" },
    ]);
  } finally {
    // Reset button
    submitBtn.innerHTML = originalText;
    submitBtn.disabled = false;
  }
}

// Form error helpers
function clearFormErrors(form) {
  // Remove error classes and messages
  form.querySelectorAll(".form-group").forEach((group) => {
    group.classList.remove("has-error");
    const errorMsg = group.querySelector(".error-message");
    if (errorMsg) {
      errorMsg.remove();
    }
  });

  // Remove general error message
  const generalError = form.querySelector(".general-error");
  if (generalError) {
    generalError.remove();
  }
}

function showFormErrors(form, errors) {
  errors.forEach((error) => {
    if (error.field === "general") {
      // Show general error at top of form
      const errorDiv = document.createElement("div");
      errorDiv.className = "general-error";
      errorDiv.style.cssText =
        "background: #fee2e2; color: #991b1b; padding: 10px; border-radius: 4px; margin-bottom: 15px; font-size: 0.875rem;";
      errorDiv.innerHTML = `<i class="fas fa-exclamation-triangle"></i> ${error.message}`;
      form.insertBefore(errorDiv, form.firstChild);
    } else {
      // Show field-specific error
      const field = form.querySelector(`[name="${error.field}"]`);
      if (field) {
        const formGroup = field.closest(".form-group");
        formGroup.classList.add("has-error");

        const errorMsg = document.createElement("div");
        errorMsg.className = "error-message";
        errorMsg.style.cssText =
          "color: #ef4444; font-size: 0.75rem; margin-top: 5px;";
        errorMsg.textContent = error.message;
        formGroup.appendChild(errorMsg);
      }
    }
  });
}

function showSuccessMessage(message) {
  const successDiv = document.createElement("div");
  successDiv.style.cssText =
    "position: fixed; top: 20px; right: 20px; background: #dcfce7; color: #166534; padding: 15px 20px; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.1); z-index: 9999; font-weight: 500;";
  successDiv.innerHTML = `<i class="fas fa-check-circle"></i> ${message}`;
  document.body.appendChild(successDiv);

  setTimeout(() => {
    successDiv.remove();
  }, 5000);
}

// Initialize auth manager
const auth = new AuthManager();

// Setup form handlers when DOM is loaded
document.addEventListener("DOMContentLoaded", function () {
  // Check if we need to redirect based on auth status
  auth.redirectIfNeeded();

  // Setup login form
  const loginForm = document.getElementById("loginForm");
  if (loginForm) {
    loginForm.addEventListener("submit", handleLogin);
  }

  // Setup register form
  const registerForm = document.getElementById("registerForm");
  if (registerForm) {
    registerForm.addEventListener("submit", handleRegister);
  }
});
