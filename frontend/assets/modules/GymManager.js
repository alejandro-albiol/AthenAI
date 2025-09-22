/**
 * Gym Management Module
 * Handles all gym-related functionality
 * Uses global utilities (loaded via script tags)
 */
class GymManager {
  constructor() {
    this.api = new ApiClient();
    this.currentGyms = [];
  }

  async loadGyms() {
    try {
      appState.setState({ loading: true });
      const response = await this.api.getGyms();

      // Handle different response structures
      let gyms = [];
      if (response && Array.isArray(response)) {
        gyms = response;
      } else if (response && response.data && Array.isArray(response.data)) {
        gyms = response.data;
      } else if (response && Array.isArray(response.gyms)) {
        gyms = response.gyms;
      } else {
        console.warn("Unexpected response structure:", response);
      }

      this.currentGyms = gyms;
      appState.setState({ gyms, loading: false });

      return gyms;
    } catch (error) {
      console.error("Error loading gyms:", error);
      appState.setState({ loading: false, error: error.message });
      notifications.error("Failed to load gyms");
      // Return empty array on error instead of throwing
      return [];
    }
  }

  async createGym(gymData) {
    try {
      const response = await this.api.createGym(gymData);
      notifications.success("Gym created successfully");

      await this.loadGyms();
      return response;
    } catch (error) {
      console.error("Error creating gym:", error);
      notifications.error("Failed to create gym: " + error.message);
      throw error;
    }
  }

  async updateGym(id, gymData) {
    try {
      const response = await this.api.updateGym(id, gymData);
      notifications.success("Gym updated successfully");

      await this.loadGyms();
      return response;
    } catch (error) {
      console.error("Error updating gym:", error);
      notifications.error("Failed to update gym: " + error.message);
      throw error;
    }
  }

  async deleteGym(id) {
    try {
      await this.api.deleteGym(id);
      notifications.success("Gym deleted successfully");

      await this.loadGyms();
      return true;
    } catch (error) {
      console.error("Error deleting gym:", error);
      notifications.error("Failed to delete gym: " + error.message);
      throw error;
    }
  }

  async restoreGym(id) {
    try {
      await this.api.restoreGym(id);
      notifications.success("Gym restored successfully");

      await this.loadGyms();
      return true;
    } catch (error) {
      console.error("Error restoring gym:", error);
      notifications.error("Failed to restore gym: " + error.message);
      throw error;
    }
  }

  getTableColumns() {
    return [
      {
        key: "name",
        title: "Gym Name",
        sortable: true,
        render: (value, item) => `
          <div class="gym-info">
            <strong class="${
              item.deleted_at ? "deleted-text" : ""
            }">${value}</strong>
            ${
              item.deleted_at
                ? '<span class="deleted-badge">Deleted</span>'
                : ""
            }
          </div>
        `,
      },
      {
        key: "contact_info",
        title: "Contact",
        render: (value, item) => `
          <div class="contact-info">
            ${
              item.email
                ? `<div><i class="fas fa-envelope"></i> ${item.email}</div>`
                : ""
            }
            ${
              item.phone
                ? `<div><i class="fas fa-phone"></i> ${item.phone}</div>`
                : ""
            }
          </div>
        `,
      },
      {
        key: "address",
        title: "Address",
        render: (value) =>
          value
            ? `<span title="${value}">${this.truncateText(value, 30)}</span>`
            : "No address",
      },
      {
        key: "created_at",
        title: "Created",
        sortable: true,
        render: (value) => this.formatDate(value),
      },
    ];
  }

  getRowActions() {
    return [
      {
        action: "view",
        icon: "fas fa-eye",
        title: "View Details",
        className: "btn btn-info",
        handler: (gym) => this.viewGymDetails(gym),
      },
      {
        action: "edit",
        icon: "fas fa-edit",
        title: "Edit Gym",
        className: "btn btn-outline",
        disabled: (gym) => Boolean(gym.deleted_at),
        handler: (gym) => this.editGym(gym),
      },
      {
        action: "delete",
        icon: "fas fa-trash",
        title: "Delete Gym",
        className: "btn btn-danger",
        disabled: (gym) => Boolean(gym.deleted_at),
        handler: (gym) => this.requestDeleteGym(gym),
      },
      {
        action: "restore",
        icon: "fas fa-undo",
        title: "Restore Gym",
        className: "btn btn-warning",
        disabled: (gym) => !Boolean(gym.deleted_at),
        handler: (gym) => this.requestRestoreGym(gym),
      },
    ];
  }

  async viewGymDetails(gym) {
    document.dispatchEvent(
      new CustomEvent("gym:view", {
        detail: { gym },
      })
    );
  }

  async editGym(gym) {
    document.dispatchEvent(
      new CustomEvent("gym:edit", {
        detail: { gym },
      })
    );
  }

  async requestDeleteGym(gym) {
    document.dispatchEvent(
      new CustomEvent("gym:delete", {
        detail: { gym },
      })
    );
  }

  async requestRestoreGym(gym) {
    document.dispatchEvent(
      new CustomEvent("gym:restore", {
        detail: { gym },
      })
    );
  }

  // Legacy methods for backwards compatibility (can be removed later)
  async confirmDeleteGym(gym) {
    return this.requestDeleteGym(gym);
  }

  async confirmRestoreGym(gym) {
    return this.requestRestoreGym(gym);
  }

  truncateText(text, maxLength) {
    if (!text || text.length <= maxLength) return text;
    return text.substring(0, maxLength) + "...";
  }

  formatDate(dateString) {
    if (!dateString) return "Unknown";

    try {
      const date = new Date(dateString);
      return date.toLocaleDateString("en-US", {
        year: "numeric",
        month: "short",
        day: "numeric",
      });
    } catch (error) {
      return "Invalid date";
    }
  }

  getFormSchema() {
    return {
      name: {
        type: "text",
        label: "Gym Name",
        required: true,
        placeholder: "Enter gym name",
      },
      email: {
        type: "email",
        label: "Contact Email",
        required: true,
        placeholder: "Enter contact email",
      },
      phone: {
        type: "tel",
        label: "Contact Phone",
        required: true,
        placeholder: "Enter contact phone number",
      },
      address: {
        type: "textarea",
        label: "Address",
        required: true,
        rows: 3,
        placeholder: "Enter gym address",
      },
    };
  }

  validateGymData(data) {
    const errors = {};

    if (!data.name || data.name.trim() === "") {
      errors.name = "Gym name is required";
    }

    if (!data.email || data.email.trim() === "") {
      errors.email = "Contact email is required";
    } else if (!this.isValidEmail(data.email)) {
      errors.email = "Please enter a valid email address";
    }

    if (!data.phone || data.phone.trim() === "") {
      errors.phone = "Contact phone is required";
    } else if (!this.isValidPhone(data.phone)) {
      errors.phone = "Please enter a valid phone number";
    }

    if (!data.address || data.address.trim() === "") {
      errors.address = "Address is required";
    }

    return {
      isValid: Object.keys(errors).length === 0,
      errors,
    };
  }

  isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  isValidPhone(phone) {
    // Simple phone validation - accepts various formats
    const phoneRegex = /^[\+]?[1-9][\d]{0,15}$/;
    const cleanPhone = phone.replace(/[\s\-\(\)\.]/g, "");
    return phoneRegex.test(cleanPhone);
  }

  getStats() {
    const gyms = this.currentGyms;
    const activeGyms = gyms.filter((gym) => !gym.deleted_at);
    const deletedGyms = gyms.filter((gym) => gym.deleted_at);

    return {
      total: gyms.length,
      active: activeGyms.length,
      deleted: deletedGyms.length,
      withContact: gyms.filter((gym) => gym.contact_email).length,
      withAddress: gyms.filter((gym) => gym.address).length,
      recentlyCreated: this.getRecentlyCreatedCount(gyms, 30), // Last 30 days
    };
  }

  getRecentlyCreatedCount(gyms, days) {
    const cutoffDate = new Date();
    cutoffDate.setDate(cutoffDate.getDate() - days);

    return gyms.filter((gym) => {
      if (!gym.created_at) return false;
      const createdDate = new Date(gym.created_at);
      return createdDate >= cutoffDate;
    }).length;
  }

  // Filter gyms by criteria
  filterGyms(criteria) {
    let filtered = [...this.currentGyms];

    if (criteria.status) {
      if (criteria.status === "active") {
        filtered = filtered.filter((gym) => !gym.deleted_at);
      } else if (criteria.status === "deleted") {
        filtered = filtered.filter((gym) => gym.deleted_at);
      }
    }

    if (criteria.search) {
      const searchTerm = criteria.search.toLowerCase();
      filtered = filtered.filter(
        (gym) =>
          gym.name.toLowerCase().includes(searchTerm) ||
          (gym.contact_name &&
            gym.contact_name.toLowerCase().includes(searchTerm)) ||
          (gym.contact_email &&
            gym.contact_email.toLowerCase().includes(searchTerm)) ||
          (gym.address && gym.address.toLowerCase().includes(searchTerm))
      );
    }

    return filtered;
  }

  // Get gym activity summary for platform overview
  getGymActivitySummary() {
    const gyms = this.currentGyms;
    const now = new Date();
    const lastWeek = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
    const lastMonth = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);

    return {
      totalGyms: gyms.length,
      activeGyms: gyms.filter((gym) => !gym.deleted_at).length,
      newThisWeek: gyms.filter(
        (gym) => gym.created_at && new Date(gym.created_at) >= lastWeek
      ).length,
      newThisMonth: gyms.filter(
        (gym) => gym.created_at && new Date(gym.created_at) >= lastMonth
      ).length,
    };
  }
}

// Make GymManager globally available
window.GymManager = GymManager;
