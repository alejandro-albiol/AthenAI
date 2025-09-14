// AthenAI Configuration
const ATHENAI_CONFIG = {
  // Email Configuration
  CONTACT_EMAIL: "contact@athenai.com",
  SUPPORT_EMAIL: "support@athenai.com",
  DEMO_EMAIL: "demo@athenai.com",

  // API Configuration
  API_BASE_URL: "/api/v1",

  // UI Configuration
  MODAL_ANIMATION_DURATION: 300,
  SUCCESS_MESSAGE_DURATION: 5000,
  REDIRECT_DELAY: 1500,

  // Company Information
  COMPANY_NAME: "AthenAI",
  SUPPORT_RESPONSE_TIME: "24 hours",
};

// DOM Content Loaded
document.addEventListener("DOMContentLoaded", function () {
  // Initialize the application
  initializeApp();
});

// Initialize application
function initializeApp() {
  // Check for invitation token in URL
  checkInvitationToken();

  // Setup form handlers
  setupFormHandlers();

  // Setup modal event listeners
  setupModalEventListeners();

  // Setup smooth scrolling
  setupSmoothScrolling();

  // Setup contact form (if it exists)
  setupContactForm();

  console.log("AthenAI Frontend initialized");
}

// Check for invitation token in URL
function checkInvitationToken() {
  const urlParams = new URLSearchParams(window.location.search);
  const inviteToken = urlParams.get("invite");

  if (inviteToken) {
    // Show gym context
    showGymContext(inviteToken);
  } else {
    // Platform admin login
    updateLoginForPlatformAdmin();
  }
}

function showGymContext(inviteToken) {
  // Store the invitation token
  document.getElementById("gymToken").value = inviteToken;

  // Show gym info section
  document.getElementById("gymInfo").style.display = "block";

  // Try to decode gym info from token (this would need backend endpoint)
  fetchGymInfoFromToken(inviteToken);

  // Automatically open the login modal for invited users
  setTimeout(() => {
    openLoginModal();
  }, 1000);
}

function updateLoginForPlatformAdmin() {
  // Hide gym info
  const gymInfo = document.getElementById("gymInfo");
  if (gymInfo) {
    gymInfo.style.display = "none";
  }

  // Clear gym token
  const gymToken = document.getElementById("gymToken");
  if (gymToken) {
    gymToken.value = "";
  }
}

async function fetchGymInfoFromToken(inviteToken) {
  try {
    // This would be a new endpoint that decodes the invitation token
    const response = await fetch(`/api/v1/invitations/decode/${inviteToken}`);
    if (response.ok) {
      const data = await response.json();
      document.getElementById("gymName").textContent =
        data.gym_name || "Your Gym";
    } else {
      document.getElementById("gymName").textContent = "Invited Gym";
    }
  } catch (error) {
    console.error("Failed to fetch gym info:", error);
    document.getElementById("gymName").textContent = "Invited Gym";
  }
}

// Modal Management
function openRegisterModal() {
  const modal = document.getElementById("registerModal");
  modal.style.display = "block";
  document.body.style.overflow = "hidden";
}

function closeRegisterModal() {
  const modal = document.getElementById("registerModal");
  if (modal) {
    modal.style.display = "none";
    document.body.style.overflow = "auto";

    // Reset form
    const registerForm = document.getElementById("registerForm");
    if (registerForm) {
      registerForm.reset();
    }
  }
}

function openLoginModal() {
  console.log("Opening login modal...");
  const modal = document.getElementById("loginModal");
  if (modal) {
    modal.style.display = "block";
    document.body.style.overflow = "hidden";
    console.log("Login modal opened successfully");
  } else {
    console.error("Login modal not found!");
  }
}

function closeLoginModal() {
  const modal = document.getElementById("loginModal");
  modal.style.display = "none";
  document.body.style.overflow = "auto";

  // Reset form
  const loginForm = document.getElementById("loginForm");
  if (loginForm) {
    loginForm.reset();
  }
}

// Get Started Modal Functions
function openGetStartedModal() {
  console.log("Opening Get Started modal...");
  const modal = document.getElementById("getStartedModal");
  if (modal) {
    modal.style.display = "block";
    document.body.style.overflow = "hidden";
    console.log("Get Started modal opened successfully");
  } else {
    console.error("Get Started modal not found!");
  }
}

function closeGetStartedModal() {
  console.log("Closing Get Started modal...");
  const modal = document.getElementById("getStartedModal");
  if (modal) {
    modal.style.display = "none";
    document.body.style.overflow = "auto";

    // Reset form
    const getStartedForm = document.getElementById("getStartedForm");
    if (getStartedForm) {
      getStartedForm.reset();
    }
    console.log("Get Started modal closed successfully");
  }
}

// Platform Admin Login Function
function openPlatformAdminLogin() {
  console.log("Opening Platform Admin login...");
  // Ensure we're in platform admin mode (no gym context)
  updateLoginForPlatformAdmin();

  // Open the login modal
  openLoginModal();
}

// Setup modal event listeners
function setupModalEventListeners() {
  // Close modals when clicking outside
  window.onclick = function (event) {
    const registerModal = document.getElementById("registerModal");
    const loginModal = document.getElementById("loginModal");
    const getStartedModal = document.getElementById("getStartedModal");

    if (event.target === registerModal) {
      closeRegisterModal();
    }

    if (event.target === loginModal) {
      closeLoginModal();
    }

    if (event.target === getStartedModal) {
      closeGetStartedModal();
    }
  };

  // Close modals on escape key
  document.addEventListener("keydown", function (event) {
    if (event.key === "Escape") {
      closeRegisterModal();
      closeLoginModal();
      closeGetStartedModal();
    }
  });
}

// Form Handlers
function setupFormHandlers() {
  console.log("Setting up form handlers...");

  // Register form handler (for backward compatibility)
  const registerForm = document.getElementById("registerForm");
  if (registerForm) {
    registerForm.addEventListener("submit", handleGymRegistration);
  }

  // Get Started form handler
  const getStartedForm = document.getElementById("getStartedForm");
  if (getStartedForm) {
    getStartedForm.addEventListener("submit", handleGetStartedSubmission);
    console.log("Get Started form handler attached");
  }

  // Login form handler
  const loginForm = document.getElementById("loginForm");
  if (loginForm) {
    loginForm.addEventListener("submit", handleLogin);
    console.log("Login form handler attached");
  }

  // Platform Admin Login button (top navigation)
  const loginBtn = document.getElementById("loginBtn");
  if (loginBtn) {
    loginBtn.addEventListener("click", openPlatformAdminLogin);
    console.log("Login button handler attached");
  } else {
    console.error("Login button not found!");
  }

  // Get Started button (hero section)
  const heroLoginBtn = document.getElementById("heroLoginBtn");
  if (heroLoginBtn) {
    heroLoginBtn.addEventListener("click", openGetStartedModal);
    console.log("Get Started button handler attached");
  } else {
    console.error("Get Started button not found!");
  }

  // Close button handlers
  const closeButtons = document.querySelectorAll(".close");
  closeButtons.forEach((button) => {
    button.addEventListener("click", function () {
      closeLoginModal();
      closeGetStartedModal();
      closeRegisterModal();
    });
  });
  console.log(
    `Close button handlers attached: ${closeButtons.length} buttons found`
  );
}

// Handle gym registration
async function handleGymRegistration(event) {
  event.preventDefault();

  const formData = new FormData(event.target);
  const accessRequestData = {
    gym_name: formData.get("gymName"),
    owner_name: formData.get("ownerName"),
    email: formData.get("email"),
    phone: formData.get("phone"),
    member_count: formData.get("memberCount"),
    current_system: formData.get("currentSystem"),
    timeline: formData.get("timeline"),
    address: formData.get("address"),
    request_type: "access_request",
    source: "landing_page",
  };

  try {
    showLoading("Submitting your access request...");

    // Try to send to backend first
    let response;
    try {
      response = await fetch("/api/v1/access-requests", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(accessRequestData),
      });
    } catch (backendError) {
      // If backend is not available, send via contact form method
      console.log("Backend not available, using email fallback");
      await sendAccessRequestEmail(accessRequestData);
      hideLoading();
      showModalSuccess(
        "Thank you! Your demo request has been submitted. We'll contact you within 24 hours to schedule your personalized demo."
      );

      setTimeout(() => {
        closeRegisterModal();
      }, 3000);
      return;
    }

    hideLoading();

    if (response && response.ok) {
      const result = await response.json();
      showModalSuccess(
        "Thank you! Your demo request has been submitted. We'll contact you within 24 hours to schedule your personalized demo."
      );

      setTimeout(() => {
        closeRegisterModal();
      }, 3000);
    } else {
      const error = response ? await response.json() : {};
      showModalError(
        error.message ||
          "Failed to submit demo request. Please try contacting us directly."
      );
    }
  } catch (error) {
    hideLoading();
    showModalError(
      `Network error. Please check your connection and try again, or contact us directly at ${ATHENAI_CONFIG.CONTACT_EMAIL}`
    );
    console.error("Access request error:", error);
  }
}

async function sendAccessRequestEmail(accessRequestData) {
  // Create email body with all the access request information
  const subject = encodeURIComponent(
    `AthenAI Access Request - ${accessRequestData.gym_name}`
  );
  const body = encodeURIComponent(`
NEW ACCESS REQUEST FOR ATHENAI

Gym Information:
- Gym Name: ${accessRequestData.gym_name}
- Owner: ${accessRequestData.owner_name}
- Email: ${accessRequestData.email}
- Phone: ${accessRequestData.phone}
- Address: ${accessRequestData.address || "Not provided"}

Business Details:
- Current Members: ${accessRequestData.member_count || "Not specified"}
- Current System: ${accessRequestData.current_system || "Not specified"}
- Timeline: ${accessRequestData.timeline || "Not specified"}

Request Details:
- Source: Landing Page
- Date: ${new Date().toLocaleDateString()}
- Time: ${new Date().toLocaleTimeString()}

Next Steps:
1. Schedule demo call
2. Discuss pricing and features
3. Setup onboarding process

---
This is an automated message from the AthenAI landing page.
  `);

  // Open email client with pre-filled information
  window.open(
    `mailto:${ATHENAI_CONFIG.DEMO_EMAIL}?subject=${subject}&body=${body}`
  );

  return { success: true };
}

// Handle Get Started form submission
async function handleGetStartedSubmission(event) {
  event.preventDefault();

  const formData = new FormData(event.target);
  const gymRequestData = {
    gym_name: formData.get("gymName"),
    owner_name: formData.get("ownerName"),
    email: formData.get("email"),
    phone: formData.get("phone"),
    member_count: formData.get("memberCount"),
    request_type: "gym_setup_request",
    source: "get_started_form",
    status: "pending_review",
  };

  try {
    showLoading("Submitting your request...");

    // Try to send to backend API (this could create a record in a "gym_requests" table)
    try {
      const response = await fetch("/api/v1/gym-requests", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(gymRequestData),
      });

      hideLoading();

      if (response.ok) {
        showModalSuccess(
          "Thank you! Your request has been submitted. We'll review your information and contact you within 24 hours to set up your gym on AthenAI."
        );

        setTimeout(() => {
          closeGetStartedModal();
        }, 3000);
      } else {
        const error = await response.json();
        showModalError(
          error.message ||
            "Failed to submit request. Please try contacting us directly."
        );
      }
    } catch (backendError) {
      // Fallback: send email if backend is not available
      console.log("Backend not available, using email fallback");
      await sendGymRequestEmail(gymRequestData);
      hideLoading();
      showModalSuccess(
        "Thank you! Your request has been submitted. We'll contact you within 24 hours to set up your gym on AthenAI."
      );

      setTimeout(() => {
        closeGetStartedModal();
      }, 3000);
    }
  } catch (error) {
    hideLoading();
    showModalError(
      `Network error. Please check your connection and try again, or contact us directly at ${ATHENAI_CONFIG.CONTACT_EMAIL}`
    );
    console.error("Gym request error:", error);
  }
}

// Send gym request via email (fallback)
async function sendGymRequestEmail(gymRequestData) {
  const subject = encodeURIComponent(
    `New Gym Setup Request - ${gymRequestData.gym_name}`
  );
  const body = encodeURIComponent(`
NEW GYM SETUP REQUEST FOR ATHENAI

Gym Name: ${gymRequestData.gym_name}
Owner/Manager: ${gymRequestData.owner_name}
Email: ${gymRequestData.email}
Phone: ${gymRequestData.phone || "Not provided"}
Member Count: ${gymRequestData.member_count || "Not specified"}

Request Type: Gym Setup Request
Source: Get Started Form
Status: Pending Review

Please contact this gym to set up their AthenAI account.
  `);

  window.open(
    `mailto:${ATHENAI_CONFIG.CONTACT_EMAIL}?subject=${subject}&body=${body}`
  );

  return { success: true };
}

// Handle login
async function handleLogin(event) {
  event.preventDefault();

  const formData = new FormData(event.target);
  const gymToken = formData.get("gymToken");

  const loginData = {
    email: formData.get("email"),
    password: formData.get("password"),
  };

  // Prepare headers
  const headers = {
    "Content-Type": "application/json",
  };

  // Add gym context if we have an invitation token
  if (gymToken) {
    try {
      // Decode the invitation token to get gym ID
      // This would need a backend endpoint or client-side decoding
      const response = await fetch(`/api/v1/invitations/decode/${gymToken}`);
      if (response.ok) {
        const data = await response.json();
        headers["X-Gym-ID"] = data.gym_id;
      }
    } catch (error) {
      console.error("Failed to decode invitation token:", error);
      showModalError(
        "Invalid invitation link. Please contact your gym administrator."
      );
      return;
    }
  }

  try {
    showLoading("Signing you in...");

    const response = await fetch("/api/v1/auth/login", {
      method: "POST",
      headers: headers,
      body: JSON.stringify(loginData),
    });

    hideLoading();

    if (response.ok) {
      const result = await response.json();

      // Show success message in modal instead of top-right notification
      const welcomeMessage = gymToken
        ? "Welcome to your gym! Redirecting to dashboard..."
        : "Welcome back! Redirecting to platform dashboard...";
      showModalSuccess(welcomeMessage);

      // Store auth token (implement proper token management)
      localStorage.setItem("auth_token", result.data.access_token);

      // Redirect to dashboard after a short delay
      setTimeout(() => {
        closeLoginModal();
        window.location.href = "/pages/dashboard/";
      }, 1500);
    } else {
      // Show red error banner for wrong credentials
      const error = await response.json();
      showModalError(
        error.message || "Invalid email or password. Please try again."
      );
    }
  } catch (error) {
    hideLoading();
    // Show error for network issues
    showModalError(
      "Network error. Please check your connection and try again."
    );
    console.error("Login error:", error);
  }
}

// Schedule demo functionality
function scheduleDemo() {
  // Open the access request modal for demo scheduling
  openRegisterModal();

  // Pre-fill the form to indicate this is a demo request
  setTimeout(() => {
    const timelineSelect = document.getElementById("timeline");
    if (timelineSelect && !timelineSelect.value) {
      timelineSelect.value = "immediately";
    }
  }, 100);
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

// Modal-specific notification functions
function showModalSuccess(message) {
  showModalMessage(message, "success");
}

function showModalError(message) {
  showModalMessage(message, "error");
}

function showModalMessage(message, type) {
  // Find the active modal
  const activeModal = document.querySelector('.modal[style*="block"]');
  if (!activeModal) {
    // Fallback to regular notification if no modal is open
    showNotification(message, type);
    return;
  }

  // Remove existing modal message
  const existingMessage = activeModal.querySelector(".modal-message");
  if (existingMessage) {
    existingMessage.remove();
  }

  // Create modal message element
  const modalMessage = document.createElement("div");
  modalMessage.className = `modal-message modal-message-${type}`;
  modalMessage.innerHTML = `
    <i class="fas fa-${getIconForType(type)}"></i>
    <span>${message}</span>
  `;

  // Add styles
  modalMessage.style.cssText = `
    padding: 12px 16px;
    margin: 20px 0;
    border-radius: 6px;
    display: flex;
    align-items: center;
    gap: 10px;
    font-weight: 500;
    animation: fadeIn 0.3s ease-out;
    background: ${getColorForType(type)}15;
    color: ${getColorForType(type)};
    border: 1px solid ${getColorForType(type)}30;
  `;

  // Insert the message before form actions or at the end of modal body
  const modalBody = activeModal.querySelector(".modal-body");
  const formActions = modalBody.querySelector(".form-actions");

  if (formActions) {
    modalBody.insertBefore(modalMessage, formActions);
  } else {
    modalBody.appendChild(modalMessage);
  }

  // Auto remove success messages after 3 seconds
  if (type === "success") {
    setTimeout(() => {
      if (modalMessage.parentElement) {
        modalMessage.style.animation = "fadeOut 0.3s ease-in";
        setTimeout(() => {
          if (modalMessage.parentElement) {
            modalMessage.remove();
          }
        }, 300);
      }
    }, 3000);
  }
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

// Pricing Section Functions
function togglePricing() {
  const toggle = document.getElementById("billingToggle");
  const amounts = document.querySelectorAll(".amount");

  amounts.forEach((amount) => {
    const monthly = amount.getAttribute("data-monthly");
    const annual = amount.getAttribute("data-annual");

    if (toggle.checked) {
      // Annual pricing
      amount.textContent = annual;
    } else {
      // Monthly pricing
      amount.textContent = monthly;
    }
  });
}

function selectPlan(planType) {
  const toggle = document.getElementById("billingToggle");
  const isAnnual = toggle.checked;

  // Store plan selection for pre-filling the access request form
  const planData = {
    type: planType,
    billing: isAnnual ? "annual" : "monthly",
  };

  localStorage.setItem("selectedPlan", JSON.stringify(planData));

  if (planType === "enterprise") {
    // For enterprise, redirect to contact form
    openContactSales();
  } else {
    // For other plans, open access request modal
    openRegisterModal();

    // Pre-fill some context based on the selected plan
    console.log(
      `Selected plan: ${planType} (${isAnnual ? "annual" : "monthly"})`
    );

    // Pre-select member count range based on the plan
    setTimeout(() => {
      const memberCountSelect = document.getElementById("memberCount");
      if (memberCountSelect && planType === "starter") {
        memberCountSelect.value = "1-100";
      } else if (memberCountSelect && planType === "professional") {
        memberCountSelect.value = "101-500";
      }
    }, 100);
  }
}

function openContactSales() {
  // For enterprise plans, scroll to contact form with pre-filled subject
  const contactSection = document.getElementById("contact");
  if (contactSection) {
    contactSection.scrollIntoView({ behavior: "smooth" });

    // Pre-fill contact form for enterprise inquiry
    setTimeout(() => {
      const subjectSelect = document.getElementById("contactSubject");
      if (subjectSelect) {
        subjectSelect.value = "partnership";
      }

      const messageTextarea = document.getElementById("contactMessage");
      if (messageTextarea && !messageTextarea.value) {
        messageTextarea.value =
          "Hi! I'm interested in the Enterprise plan for my gym. Please contact me to discuss pricing and features for large-scale implementation.";
      }
    }, 500);
  }
}

// Add smooth scrolling enhancement for pricing link
function enhanceSmoothScrolling() {
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

// Contact Form Functions
function setupContactForm() {
  const form = document.getElementById("contactForm");
  if (form) {
    form.addEventListener("submit", handleContactSubmit);
  }
}

function validateContactForm() {
  const name = document.getElementById("contactName").value.trim();
  const email = document.getElementById("contactEmail").value.trim();
  const message = document.getElementById("contactMessage").value.trim();

  let isValid = true;

  // Clear previous errors
  document.getElementById("nameError").textContent = "";
  document.getElementById("emailError").textContent = "";
  document.getElementById("messageError").textContent = "";
  document.getElementById("recaptchaError").textContent = "";

  // Validate name
  if (name.length < 2) {
    document.getElementById("nameError").textContent =
      "Name must be at least 2 characters long";
    isValid = false;
  }

  // Validate email
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(email)) {
    document.getElementById("emailError").textContent =
      "Please enter a valid email address";
    isValid = false;
  }

  // Validate message
  if (message.length < 10) {
    document.getElementById("messageError").textContent =
      "Message must be at least 10 characters long";
    isValid = false;
  }

  // Validate reCAPTCHA (when implemented)
  if (typeof grecaptcha !== "undefined") {
    const recaptchaResponse = grecaptcha.getResponse();
    if (!recaptchaResponse) {
      document.getElementById("recaptchaError").textContent =
        "Please complete the reCAPTCHA verification";
      isValid = false;
    }
  }

  return isValid;
}

async function handleContactSubmit(e) {
  e.preventDefault();

  if (!validateContactForm()) {
    return;
  }

  const submitBtn = document.getElementById("contactSubmit");
  const btnText = submitBtn.querySelector(".btn-text");
  const btnLoading = submitBtn.querySelector(".btn-loading");
  const statusDiv = document.getElementById("contactStatus");

  // Show loading state
  btnText.style.display = "none";
  btnLoading.style.display = "inline";
  submitBtn.disabled = true;
  statusDiv.style.display = "none";

  // Collect form data
  const formData = {
    name: document.getElementById("contactName").value.trim(),
    email: document.getElementById("contactEmail").value.trim(),
    subject:
      document.getElementById("contactSubject").value || "General Inquiry",
    message: document.getElementById("contactMessage").value.trim(),
    recaptcha:
      typeof grecaptcha !== "undefined" ? grecaptcha.getResponse() : null,
    timestamp: new Date().toISOString(),
    source: "AthenAI Landing Page",
  };

  try {
    // For now, we'll use a simple email service or backend endpoint
    // You can replace this with your preferred email service
    const response = await sendContactEmail(formData);

    if (response.success) {
      showStatus(
        "success",
        "Thank you! Your message has been sent successfully. We'll get back to you within 24 hours."
      );
      document.getElementById("contactForm").reset();
      if (typeof grecaptcha !== "undefined") {
        grecaptcha.reset();
      }
    } else {
      throw new Error(response.message || "Failed to send message");
    }
  } catch (error) {
    console.error("Contact form error:", error);
    showStatus(
      "error",
      `Sorry, there was an error sending your message. Please try again or email us directly at ${ATHENAI_CONFIG.CONTACT_EMAIL}`
    );
  } finally {
    // Reset button state
    btnText.style.display = "inline";
    btnLoading.style.display = "none";
    submitBtn.disabled = false;
  }
}

function showStatus(type, message) {
  const statusDiv = document.getElementById("contactStatus");
  statusDiv.className = `form-status ${type}`;
  statusDiv.textContent = message;
  statusDiv.style.display = "block";

  // Auto-hide success messages after 5 seconds
  if (type === "success") {
    setTimeout(() => {
      statusDiv.style.display = "none";
    }, 5000);
  }
}

async function sendContactEmail(formData) {
  // Option 1: Use a backend endpoint (recommended)
  try {
    const response = await fetch("/api/contact", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(formData),
    });

    if (response.ok) {
      return { success: true };
    } else {
      throw new Error("Server error");
    }
  } catch (error) {
    console.log("Backend not available, using alternative method");

    // Option 2: Use EmailJS or similar service
    // You'll need to set up EmailJS account and get your service ID
    /*
    try {
      const result = await emailjs.send(
        'your_service_id',
        'your_template_id',
        {
          to_email: ATHENAI_CONFIG.CONTACT_EMAIL,
          from_name: formData.name,
          from_email: formData.email,
          subject: `AthenAI Contact: ${formData.subject}`,
          message: formData.message,
        }
      );
      return { success: true };
    } catch (emailError) {
      return { success: false, message: emailError.text };
    }
    */

    // Option 3: Fallback - open email client (for development)
    const subject = encodeURIComponent(`AthenAI Contact: ${formData.subject}`);
    const body = encodeURIComponent(`
Name: ${formData.name}
Email: ${formData.email}
Subject: ${formData.subject}

Message:
${formData.message}

---
Sent from AthenAI Landing Page at ${formData.timestamp}
    `);

    window.open(
      `mailto:${ATHENAI_CONFIG.CONTACT_EMAIL}?subject=${subject}&body=${body}`
    );
    return { success: true };
  }
}
