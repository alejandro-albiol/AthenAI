/**
 * Notification Utility Module
 * Centralized notification system with different types and positions
 */
class NotificationManager {
  constructor() {
    this.container = null;
    this.notifications = new Map();
    this.createContainer();
  }

  createContainer() {
    this.container = document.createElement("div");
    this.container.id = "notification-container";
    this.container.className = "notification-container";
    document.body.appendChild(this.container);
  }

  show(message, type = "info", options = {}) {
    const config = {
      duration: 5000,
      closeable: true,
      position: "top-right",
      ...options,
    };

    const notification = this.createNotification(message, type, config);
    this.container.appendChild(notification.element);

    // Store notification
    this.notifications.set(notification.id, notification);

    // Auto-remove after duration
    if (config.duration > 0) {
      setTimeout(() => {
        this.remove(notification.id);
      }, config.duration);
    }

    // Animation
    setTimeout(() => {
      notification.element.classList.add("notification-show");
    }, 10);

    return notification.id;
  }

  createNotification(message, type, config) {
    const id = Date.now() + Math.random();
    const element = document.createElement("div");

    element.className = `notification notification-${type} notification-${config.position}`;
    element.setAttribute("data-id", id);

    element.innerHTML = `
      <div class="notification-content">
        <i class="notification-icon fas fa-${this.getIcon(type)}"></i>
        <span class="notification-message">${message}</span>
        ${
          config.closeable
            ? '<button class="notification-close"><i class="fas fa-times"></i></button>'
            : ""
        }
      </div>
    `;

    // Close button handler
    if (config.closeable) {
      const closeBtn = element.querySelector(".notification-close");
      closeBtn.addEventListener("click", () => {
        this.remove(id);
      });
    }

    return { id, element, type, config };
  }

  remove(id) {
    const notification = this.notifications.get(id);
    if (!notification) return;

    notification.element.classList.add("notification-hide");

    setTimeout(() => {
      if (notification.element.parentNode) {
        notification.element.parentNode.removeChild(notification.element);
      }
      this.notifications.delete(id);
    }, 300);
  }

  removeAll() {
    this.notifications.forEach((notification, id) => {
      this.remove(id);
    });
  }

  getIcon(type) {
    const icons = {
      success: "check-circle",
      error: "exclamation-circle",
      warning: "exclamation-triangle",
      info: "info-circle",
    };
    return icons[type] || "info-circle";
  }

  // Convenience methods
  success(message, options = {}) {
    return this.show(message, "success", options);
  }

  error(message, options = {}) {
    return this.show(message, "error", { duration: 8000, ...options });
  }

  warning(message, options = {}) {
    return this.show(message, "warning", options);
  }

  info(message, options = {}) {
    return this.show(message, "info", options);
  }
}

// Create global instance
const notifications = new NotificationManager();

// Make available globally for vanilla JS
window.notifications = notifications;
window.NotificationManager = NotificationManager;
