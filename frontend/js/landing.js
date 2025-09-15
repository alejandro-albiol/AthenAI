/**
 * Landing Page Manager
 * Handles landing page functionality using modular components
 * Vanilla JS - no imports, works directly in browsers
 */

class LandingPageManager {
  constructor() {
    this.config = {
      CONTACT_EMAIL: "contact@athenai.com",
      SUPPORT_EMAIL: "support@athenai.com",
      DEMO_EMAIL: "demo@athenai.com",
      API_BASE_URL: "/api/v1",
      MODAL_ANIMATION_DURATION: 300,
      SUCCESS_MESSAGE_DURATION: 5000,
      REDIRECT_DELAY: 1500,
      COMPANY_NAME: "AthenAI",
      SUPPORT_RESPONSE_TIME: "24 hours",
    };

    this.currentGym = null;
    this.inviteToken = null;

    this.init();
  }

  async init() {
    try {
      // Wait for dependencies (Modal, notifications, etc.)
      await this.waitForDependencies();

      // Check for invitation token in URL
      this.checkInvitationToken();

      // Setup event listeners
      this.setupEventListeners();

      // Setup smooth scrolling
      this.setupSmoothScrolling();

      // Initialize components
      this.initializeComponents();

      console.log("Landing Page Manager initialized");
    } catch (error) {
      console.error("Failed to initialize landing page:", error);
    }
  }

  async waitForDependencies() {
    const requiredClasses = ["Modal", "notifications", "ApiClient"];

    let attempts = 0;
    const maxAttempts = 50;

    while (attempts < maxAttempts) {
      const allLoaded = requiredClasses.every((className) => {
        return window[className] !== undefined;
      });

      if (allLoaded) {
        console.log("Landing page dependencies loaded");
        return;
      }

      await new Promise((resolve) => setTimeout(resolve, 100));
      attempts++;
    }

    console.warn("Some landing page dependencies may not have loaded");
  }

  checkInvitationToken() {
    const urlParams = new URLSearchParams(window.location.search);
    this.inviteToken = urlParams.get("invite");

    if (this.inviteToken) {
      this.showGymContext(this.inviteToken);
    } else {
      this.updateLoginForPlatformAdmin();
    }
  }

  async showGymContext(inviteToken) {
    try {
      // Try to decode gym info from token
      const gymInfo = await this.fetchGymInfoFromToken(inviteToken);

      if (gymInfo) {
        this.currentGym = gymInfo;

        // Show gym-specific messaging
        this.updateUIForGymContext(gymInfo);

        // Automatically show login modal after a short delay
        setTimeout(() => {
          this.openLoginModal();
        }, 1000);

        notifications.info(
          `Welcome to ${gymInfo.name}! Please login to continue.`
        );
      }
    } catch (error) {
      console.error("Error fetching gym info:", error);
      notifications.error("Invalid invitation link");
    }
  }

  updateLoginForPlatformAdmin() {
    this.currentGym = null;
    this.inviteToken = null;

    // Update UI for platform admin context
    this.updateUIForPlatformContext();
  }

  async fetchGymInfoFromToken(inviteToken) {
    try {
      const api = new ApiClient();
      const response = await api.request(`/invitations/decode/${inviteToken}`, {
        method: "GET",
      });

      return response.data;
    } catch (error) {
      console.error("Failed to decode invitation token:", error);
      return null;
    }
  }

  updateUIForGymContext(gymInfo) {
    // Update hero section for gym context
    const heroTitle = document.querySelector(".hero-content h1");
    const heroSubtitle = document.querySelector(".hero-subtitle");

    if (heroTitle) {
      heroTitle.innerHTML = `Welcome to <span class="highlight">${gymInfo.name}</span>`;
    }

    if (heroSubtitle) {
      heroSubtitle.textContent = "Access your gym management platform.";
    }

    // Update CTA section
    const ctaTitle = document.querySelector(".cta h2");
    const ctaText = document.querySelector(".cta p");

    if (ctaTitle) {
      ctaTitle.textContent = "Ready to Start?";
    }

    if (ctaText) {
      ctaText.textContent = "Login to access your gym dashboard.";
    }
  }

  updateUIForPlatformContext() {
    // Reset to default platform admin messaging
    const heroTitle = document.querySelector(".hero-content h1");
    const heroSubtitle = document.querySelector(".hero-subtitle");

    if (heroTitle) {
      heroTitle.innerHTML =
        '<span class="highlight">AthenAI</span> Fitness Management';
    }

    if (heroSubtitle) {
      heroSubtitle.textContent =
        "A modern platform for gym management and personalized workout planning.";
    }
  }

  setupEventListeners() {
    // Login button clicks
    const loginButtons = document.querySelectorAll(
      "#loginBtn, #heroLoginBtn, .login-trigger"
    );
    loginButtons.forEach((btn) => {
      btn.addEventListener("click", (e) => {
        e.preventDefault();
        this.openLoginModal();
      });
    });

    // Smooth scroll for navigation links
    const navLinks = document.querySelectorAll('.nav-link[href^="#"]');
    navLinks.forEach((link) => {
      link.addEventListener("click", (e) => {
        e.preventDefault();
        const target = document.querySelector(link.getAttribute("href"));
        if (target) {
          target.scrollIntoView({
            behavior: "smooth",
            block: "start",
          });
        }
      });
    });

    // Contact form submission (if exists)
    const contactForm = document.getElementById("contactForm");
    if (contactForm) {
      contactForm.addEventListener("submit", (e) => {
        e.preventDefault();
        this.handleContactForm(e.target);
      });
    }
  }

  setupSmoothScrolling() {
    // Enable smooth scrolling for the entire page
    document.documentElement.style.scrollBehavior = "smooth";
  }

  initializeComponents() {
    // Initialize any landing page specific components
    this.setupFeatureCards();
  }

  setupFeatureCards() {
    // Add interactive behavior to feature cards
    const featureCards = document.querySelectorAll(".feature-card");

    featureCards.forEach((card) => {
      card.addEventListener("mouseenter", () => {
        card.style.transform = "translateY(-8px)";
      });

      card.addEventListener("mouseleave", () => {
        card.style.transform = "translateY(0)";
      });
    });
  }

  openLoginModal() {
    const loginFormHtml = this.generateLoginForm();

    const modal = new Modal({
      title: this.currentGym
        ? `Login to ${this.currentGym.name}`
        : "Platform Admin Login",
      content: loginFormHtml,
      size: "medium",
      buttons: [
        {
          text: "Cancel",
          action: "dismiss",
          className: "btn btn-secondary",
        },
        {
          text: "Login",
          action: "login",
          className: "btn btn-primary",
          handler: () => this.handleLogin(),
        },
      ],
    });

    modal.show();
    this.loginModal = modal;
  }

  generateLoginForm() {
    return `
      <form id="loginForm" class="form-grid">
        <div class="form-group">
          <label class="form-label" for="email">Email Address</label>
          <input 
            type="email" 
            class="form-control" 
            id="email" 
            name="email" 
            required 
            placeholder="Enter your email"
            autocomplete="email"
          >
        </div>
        
        <div class="form-group">
          <label class="form-label" for="password">Password</label>
          <input 
            type="password" 
            class="form-control" 
            id="password" 
            name="password" 
            required 
            placeholder="Enter your password"
            autocomplete="current-password"
          >
        </div>
        
        ${
          this.inviteToken
            ? `<input type="hidden" name="invite_token" value="${this.inviteToken}">`
            : ""
        }
        
        <div class="form-group">
          <label class="checkbox-container">
            <input type="checkbox" name="remember" id="remember">
            <span class="checkmark"></span>
            Remember me
          </label>
        </div>
      </form>
      
      <div class="login-help">
        <p>
          <a href="#" onclick="landingPage.showForgotPassword()">Forgot your password?</a>
        </p>
        ${
          !this.currentGym
            ? "<p><small>Platform administrators only. Need gym access? Contact your gym manager.</small></p>"
            : ""
        }
      </div>
    `;
  }

  async handleLogin() {
    try {
      const form = document.getElementById("loginForm");
      const formData = getFormData(form);

      // Validate form
      if (!formData.email || !formData.password) {
        notifications.error("Please fill in all required fields");
        return false;
      }

      // Show loading
      this.showLoginLoading(true);

      const api = new ApiClient();
      const response = await api.login(
        formData.email,
        formData.password,
        formData.invite_token
      );

      if (response.status === "success" && response.data) {
        // Store auth token
        localStorage.setItem("auth_token", response.data.access_token);
        localStorage.setItem("refresh_token", response.data.refresh_token);
        localStorage.setItem(
          "user_info",
          JSON.stringify(response.data.user_info)
        );

        notifications.success("Login successful! Redirecting...");

        this.loginModal.hide();

        // Redirect to dashboard with absolute path
        setTimeout(() => {
          const currentOrigin = window.location.origin;
          const currentPath = window.location.pathname.replace(
            "/index.html",
            ""
          );
          const redirectUrl = `${currentOrigin}${currentPath}/pages/dashboard/index.html`;
          window.location.href = redirectUrl;
        }, this.config.REDIRECT_DELAY);

        return true;
      } else {
        notifications.error(response.message || "Login failed");
        return false;
      }
    } catch (error) {
      console.error("Login error:", error);
      notifications.error("Login failed. Please try again.");
      return false;
    } finally {
      this.showLoginLoading(false);
    }
  }

  showLoginLoading(loading) {
    const submitBtn = document.querySelector(".btn-primary");
    if (submitBtn) {
      if (loading) {
        submitBtn.disabled = true;
        submitBtn.innerHTML =
          '<i class="fas fa-spinner fa-spin"></i> Logging in...';
      } else {
        submitBtn.disabled = false;
        submitBtn.innerHTML = "Login";
      }
    }
  }

  showForgotPassword() {
    notifications.info(
      "Password reset functionality coming soon. Please contact support."
    );
  }

  async handleContactForm(form) {
    try {
      const formData = getFormData(form);

      // Basic validation
      if (!formData.name || !formData.email || !formData.message) {
        notifications.error("Please fill in all required fields");
        return;
      }

      // Show loading
      const submitBtn = form.querySelector('button[type="submit"]');
      const originalText = submitBtn.textContent;
      submitBtn.disabled = true;
      submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Sending...';

      const api = new ApiClient();
      await api.request("/contact", {
        method: "POST",
        body: JSON.stringify(formData),
      });

      notifications.success(
        "Message sent successfully! We'll get back to you soon."
      );
      form.reset();
    } catch (error) {
      console.error("Contact form error:", error);
      notifications.error("Failed to send message. Please try again.");
    } finally {
      const submitBtn = form.querySelector('button[type="submit"]');
      submitBtn.disabled = false;
      submitBtn.textContent = originalText;
    }
  }
}

// Initialize landing page when DOM is loaded
document.addEventListener("DOMContentLoaded", function () {
  window.landingPage = new LandingPageManager();
  console.log("Landing page initialized");
});
