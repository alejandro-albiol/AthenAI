/**
 * User Management Module
 * Handles multitenant user management with proper role-based access
 */
class UserManager {
  constructor() {
    this.api = new ApiClient();
    this.currentUsers = [];
    this.currentGymId = null;
  }

  setGymContext(gymId) {
    this.currentGymId = gymId;
  }

  async loadUsers(gymId = null) {
    try {
      const targetGymId = gymId || this.currentGymId;

      if (!targetGymId) {
        // Platform admin view - get all users
        const response = await this.api.getAllUsers();
        const users = response?.data || [];
        this.currentUsers = users;
        return users;
      } else {
        // Gym-specific view
        const response = await this.api.getUsers(targetGymId);

        // Handle different possible response structures
        let users = [];
        if (response) {
          if (Array.isArray(response)) {
            users = response;
          } else if (response.data && Array.isArray(response.data)) {
            users = response.data;
          } else if (response.users && Array.isArray(response.users)) {
            users = response.users;
          } else {
            console.log("Unexpected response structure:", response);
            users = [];
          }
        }

        this.currentUsers = users;
        return users;
      }
    } catch (error) {
      console.error("Error loading users:", error);
      this.currentUsers = [];
      return [];
    }
  }

  async createUser(userData) {
    try {
      let targetGymId = this.currentGymId;

      // For platform admin operations, allow userData to specify gym_id
      if (userData.gym_id && !targetGymId) {
        targetGymId = userData.gym_id;
      }

      if (!targetGymId) {
        throw new Error("Please select a gym first to create users");
      }

      const response = await this.api.createUser(userData, targetGymId);

      // Clear cache to force refresh
      this.currentUsers = [];

      return response;
    } catch (error) {
      console.error("Error creating user:", error);
      throw error;
    }
  }

  async updateUser(userId, userData) {
    try {
      if (!this.currentGymId) {
        throw new Error("Please select a gym first to update users");
      }
      const response = await this.api.updateUser(
        userId,
        userData,
        this.currentGymId
      );
      return response;
    } catch (error) {
      console.error("Error updating user:", error);
      throw error;
    }
  }

  async deleteUser(userId) {
    try {
      if (!this.currentGymId) {
        throw new Error("Please select a gym first to delete users");
      }
      const response = await this.api.deleteUser(userId, this.currentGymId);
      return response;
    } catch (error) {
      console.error("Error deleting user:", error);
      throw error;
    }
  }

  async verifyUser(userId) {
    try {
      if (!this.currentGymId) {
        throw new Error("Please select a gym first to verify users");
      }
      const response = await this.api.verifyUser(userId, this.currentGymId);
      return response;
    } catch (error) {
      console.error("Error verifying user:", error);
      throw error;
    }
  }

  async setUserActive(userId, active) {
    try {
      if (!this.currentGymId) {
        throw new Error("Please select a gym first to update user status");
      }
      const response = await this.api.setUserActive(
        userId,
        active,
        this.currentGymId
      );
      return response;
    } catch (error) {
      console.error("Error updating user status:", error);
      throw error;
    }
  }

  getUserRoleIcon(role) {
    const icons = {
      platform_admin: '<i class="fas fa-crown" style="color: #f39c12;"></i>',
      gym_admin: '<i class="fas fa-user-shield" style="color: #e74c3c;"></i>',
      trainer: '<i class="fas fa-dumbbell" style="color: #3498db;"></i>',
      member: '<i class="fas fa-user-check" style="color: #27ae60;"></i>',
      guest: '<i class="fas fa-user-clock" style="color: #f39c12;"></i>',
    };
    return (
      icons[role] ||
      '<i class="fas fa-user-question" style="color: #95a5a6;"></i>'
    );
  }

  formatRole(role) {
    const formatted = {
      platform_admin: "Platform Admin",
      gym_admin: "Gym Administrator",
      trainer: "Personal Trainer",
      member: "Gym Member",
      guest: "Guest Access",
    };
    return formatted[role] || role;
  }

  getUserStatusBadge(user) {
    if (!user.is_active) {
      return '<span class="badge badge-danger">Inactive</span>';
    }

    // Special handling for different roles
    if (user.role === "guest") {
      if (user.verified) {
        return '<span class="badge badge-info">Guest (Verified)</span>';
      } else {
        return '<span class="badge badge-warning">Guest (Trial)</span>';
      }
    }

    if (user.role === "member") {
      if (user.verified) {
        return '<span class="badge badge-success">Active Member</span>';
      } else {
        return '<span class="badge badge-warning">Pending Member</span>';
      }
    }

    // For staff roles (gym_admin, trainer, platform_admin)
    if (user.verified) {
      return '<span class="badge badge-primary">Staff</span>';
    } else {
      return '<span class="badge badge-warning">Unverified Staff</span>';
    }
  }

  getTableColumns() {
    return [
      {
        key: "username",
        title: "User",
        sortable: true,
        render: (value, user) => `
          <div class="user-info">
            <div class="user-icon-inline">${this.getUserRoleIcon(
              user.role
            )}</div>
            <div>
              <strong>${value}</strong>
              <div class="user-email-small">${user.email}</div>
            </div>
          </div>
        `,
      },
      {
        key: "contact_info",
        title: "Contact",
        render: (value, user) => `
          <div class="contact-info">
            ${
              user.phone
                ? `<div><i class="fas fa-phone"></i> ${user.phone}</div>`
                : "<div class='text-muted'>No phone</div>"
            }
            <div class="user-role-small">${this.formatRole(user.role)}</div>
          </div>
        `,
      },
      {
        key: "created_at",
        title: "Joined",
        sortable: true,
        render: (value) => {
          return value ? new Date(value).toLocaleDateString() : "Unknown";
        },
      },
    ];
  }

  getRowActions() {
    return [
      {
        action: "edit",
        icon: "fas fa-edit",
        title: "Edit User",
        className: "btn btn-sm btn-outline",
        handler: (user) => this.editUser(user),
      },
      {
        action: "verify",
        icon: "fas fa-check-circle",
        title: "Verify User",
        className: "btn btn-sm btn-success",
        handler: (user) => this.confirmVerifyUser(user),
        show: (user) => !user.verified,
      },
      {
        action: "activate",
        icon: "fas fa-toggle-on",
        title: "Activate User",
        className: "btn btn-sm btn-info",
        handler: (user) => this.toggleUserActive(user),
        show: (user) => !user.is_active,
      },
      {
        action: "deactivate",
        icon: "fas fa-toggle-off",
        title: "Deactivate User",
        className: "btn btn-sm btn-warning",
        handler: (user) => this.toggleUserActive(user),
        show: (user) => user.is_active,
      },
      {
        action: "delete",
        icon: "fas fa-trash",
        title: "Delete User",
        className: "btn btn-sm btn-danger",
        handler: (user) => this.confirmDeleteUser(user),
      },
    ];
  }

  async editUser(user) {
    dashboard.openUserModal("edit", user);
  }

  async confirmVerifyUser(user) {
    if (confirm(`Are you sure you want to verify "${user.username}"?`)) {
      try {
        await this.verifyUser(user.id);
        notifications.success("User verified successfully");
        dashboard.loadUsersManagement();
      } catch (error) {
        notifications.error("Error verifying user");
      }
    }
  }

  async toggleUserActive(user) {
    const newStatus = !user.is_active;
    const action = newStatus ? "activate" : "deactivate";

    if (confirm(`Are you sure you want to ${action} "${user.username}"?`)) {
      try {
        await this.setUserActive(user.id, newStatus);
        notifications.success(`User ${action}d successfully`);
        dashboard.loadUsersManagement();
      } catch (error) {
        notifications.error(`Error ${action}ing user`);
      }
    }
  }

  async confirmDeleteUser(user) {
    if (
      confirm(
        `Are you sure you want to delete "${user.username}"? This action cannot be undone.`
      )
    ) {
      try {
        await this.deleteUser(user.id);
        notifications.success("User deleted successfully");
        dashboard.loadUsersManagement();
      } catch (error) {
        notifications.error("Error deleting user");
      }
    }
  }

  getFormFields() {
    const baseFields = {
      username: {
        type: "text",
        label: "Username",
        placeholder: "Enter username",
        required: true,
        help: "Unique username for login",
      },
      email: {
        type: "email",
        label: "Email",
        placeholder: "Enter email address",
        required: true,
        help: "User's email address",
      },
      password: {
        type: "password",
        label: "Password",
        placeholder: "Enter password",
        required: true,
        help: "Minimum 8 characters",
        showOnEdit: false, // Don't show password field when editing
      },
      role: {
        type: "select",
        label: "Role",
        placeholder: "Select role",
        required: true,
        options: [
          { value: "", label: "Select role" },
          { value: "guest", label: "Guest Access" },
          { value: "member", label: "Gym Member" },
          { value: "trainer", label: "Personal Trainer" },
          { value: "gym_admin", label: "Gym Administrator" },
        ],
        help: "User's role and permissions in the gym",
      },
      phone: {
        type: "tel",
        label: "Phone",
        placeholder: "Enter phone number",
        required: false,
        help: "Contact phone number",
      },
      description: {
        type: "textarea",
        label: "Description",
        placeholder: "Enter description (optional)",
        required: false,
        rows: 3,
        help: "Additional notes about the user",
      },
    };

    return baseFields;
  }

  getMemberSpecificFields() {
    return {
      training_phase: {
        type: "select",
        label: "Training Phase",
        placeholder: "Select training phase",
        required: false,
        options: [
          { value: "", label: "Select phase" },
          { value: "weight_loss", label: "Weight Loss" },
          { value: "muscle_gain", label: "Muscle Gain" },
          { value: "cardio_improve", label: "Cardio Improvement" },
          { value: "maintenance", label: "Maintenance" },
        ],
        help: "Current training goal/phase",
      },
      motivation: {
        type: "select",
        label: "Motivation",
        placeholder: "Select motivation",
        required: false,
        options: [
          { value: "", label: "Select motivation" },
          { value: "medical_recommendation", label: "Medical Recommendation" },
          { value: "self_improvement", label: "Self Improvement" },
          { value: "competition", label: "Competition" },
          { value: "rehabilitation", label: "Rehabilitation" },
          { value: "wellbeing", label: "Wellbeing" },
        ],
        help: "Primary motivation for training",
      },
      special_situation: {
        type: "select",
        label: "Special Situation",
        placeholder: "Select if applicable",
        required: false,
        options: [
          { value: "", label: "None" },
          { value: "pregnancy", label: "Pregnancy" },
          { value: "post_partum", label: "Post Partum" },
          { value: "injury_recovery", label: "Injury Recovery" },
          { value: "chronic_condition", label: "Chronic Condition" },
          { value: "elderly_population", label: "Elderly Population" },
          { value: "physical_limitation", label: "Physical Limitation" },
          { value: "none", label: "No Special Situation" },
        ],
        help: "Any special medical or physical considerations",
      },
    };
  }

  getAllFormFields() {
    return {
      ...this.getFormFields(),
      ...this.getMemberSpecificFields(),
    };
  }

  validateUserData(data) {
    const errors = [];

    if (!data.username || data.username.trim().length === 0) {
      errors.push("Username is required");
    }

    if (!data.email || data.email.trim().length === 0) {
      errors.push("Email is required");
    }

    if (!data.role || data.role.trim().length === 0) {
      errors.push("Role is required");
    }

    if (data.password && data.password.length < 8) {
      errors.push("Password must be at least 8 characters");
    }

    if (data.username && data.username.length > 50) {
      errors.push("Username must be less than 50 characters");
    }

    if (data.email && !this.isValidEmail(data.email)) {
      errors.push("Please enter a valid email address");
    }

    const validRoles = [
      "guest",
      "member",
      "trainer",
      "gym_admin",
      "platform_admin",
    ];
    if (data.role && !validRoles.includes(data.role)) {
      errors.push("Invalid role selection");
    }

    return errors;
  }

  isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  prepareDataForSubmission(formData, isEdit = false) {
    const data = { ...formData };

    // Trim string fields
    if (data.username) data.username = data.username.trim();
    if (data.email) data.email = data.email.trim();
    if (data.phone) data.phone = data.phone.trim();
    if (data.description) data.description = data.description.trim();
    if (data.special_situation)
      data.special_situation = data.special_situation.trim();

    // Remove empty optional fields
    if (!data.phone) delete data.phone;
    if (!data.description) delete data.description;

    // Handle member-specific fields based on role
    if (data.role === "member") {
      // For members, clean up empty member-specific fields
      if (!data.training_phase) delete data.training_phase;
      if (!data.motivation) delete data.motivation;
      if (!data.special_situation) delete data.special_situation;
    } else {
      // For admin, trainer, and other roles, remove member-specific fields entirely
      // The backend will set appropriate defaults
      delete data.training_phase;
      delete data.motivation;
      delete data.special_situation;
    }

    // Don't include password in edit requests unless it's being changed
    if (isEdit && (!data.password || data.password.trim() === "")) {
      delete data.password;
    }

    return data;
  }

  // Get user statistics for dashboard
  getUserStats(users = this.currentUsers) {
    const stats = {
      total: users.length,
      active: 0,
      verified: 0,
      unverified: 0,
      inactive: 0,
      roleBreakdown: {},
    };

    users.forEach((user) => {
      if (user.is_active) stats.active++;
      else stats.inactive++;

      if (user.verified) stats.verified++;
      else stats.unverified++;

      stats.roleBreakdown[user.role] =
        (stats.roleBreakdown[user.role] || 0) + 1;
    });

    return stats;
  }
}
