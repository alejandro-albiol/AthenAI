/**
 * Muscular Group Management Module
 * Handles all muscular group-related functionality
 */
class MuscularGroupManager {
  constructor() {
    this.api = new ApiClient();
    this.currentMuscularGroups = [];
  }

  async loadMuscularGroups() {
    try {
      const response = await this.api.getMuscularGroups();
      return response?.data || [];
    } catch (error) {
      console.error("Error loading muscular groups:", error);
      return [];
    }
  }

  async createMuscularGroup(muscularGroupData) {
    try {
      const response = await this.api.createMuscularGroup(muscularGroupData);
      return response;
    } catch (error) {
      console.error("Error creating muscular group:", error);
      throw error;
    }
  }

  async updateMuscularGroup(id, muscularGroupData) {
    try {
      const response = await this.api.updateMuscularGroup(
        id,
        muscularGroupData
      );
      return response;
    } catch (error) {
      console.error("Error updating muscular group:", error);
      throw error;
    }
  }

  async deleteMuscularGroup(id) {
    try {
      const response = await this.api.deleteMuscularGroup(id);
      return response;
    } catch (error) {
      console.error("Error deleting muscular group:", error);
      throw error;
    }
  }

  getMuscularGroupIcon(bodyPart) {
    const icons = {
      upper_body: '<i class="fas fa-user-tie" style="color: #3498db;"></i>',
      lower_body: '<i class="fas fa-walking" style="color: #27ae60;"></i>',
      core: '<i class="fas fa-circle" style="color: #e67e22;"></i>',
      full_body: '<i class="fas fa-male" style="color: #9b59b6;"></i>',
    };
    return (
      icons[bodyPart] ||
      '<i class="fas fa-question-circle" style="color: #95a5a6;"></i>'
    );
  }

  formatBodyPart(bodyPart) {
    const formatted = {
      upper_body: "Upper Body",
      lower_body: "Lower Body",
      core: "Core",
      full_body: "Full Body",
    };
    return formatted[bodyPart] || bodyPart;
  }

  truncateText(text, maxLength) {
    if (!text) return "";
    return text.length > maxLength
      ? text.substring(0, maxLength) + "..."
      : text;
  }

  getTableColumns() {
    return [
      {
        key: "name",
        title: "Name",
        sortable: true,
        render: (value, item) => `
          <div class="equipment-info">
            <div class="equipment-icon-inline">${this.getMuscularGroupIcon(
              item.body_part
            )}</div>
            <div>
              <strong>${value}</strong>
              <div class="equipment-category-small">${this.formatBodyPart(
                item.body_part
              )}</div>
            </div>
          </div>
        `,
      },
      {
        key: "body_part",
        title: "Body Part",
        sortable: true,
        render: (value) => this.formatBodyPart(value),
      },
      {
        key: "description",
        title: "Description",
        render: (value) =>
          value
            ? `<span title="${value}">${this.truncateText(value, 50)}</span>`
            : "No description",
      },
    ];
  }

  getRowActions() {
    return [
      {
        action: "edit",
        icon: "fas fa-edit",
        title: "Edit Muscular Group",
        className: "btn btn-sm btn-outline",
        handler: (muscularGroup) => this.editMuscularGroup(muscularGroup),
      },
      {
        action: "delete",
        icon: "fas fa-trash",
        title: "Delete Muscular Group",
        className: "btn btn-sm btn-danger",
        handler: (muscularGroup) =>
          this.confirmDeleteMuscularGroup(muscularGroup),
      },
    ];
  }

  async editMuscularGroup(muscularGroup) {
    // Open the edit modal using the dashboard's modal system
    dashboard.openMuscularGroupModal("edit", muscularGroup);
  }

  async confirmDeleteMuscularGroup(muscularGroup) {
    if (
      confirm(
        `Are you sure you want to delete the muscle group "${muscularGroup.name}"?`
      )
    ) {
      try {
        await this.deleteMuscularGroup(muscularGroup.id);
        notifications.success("Muscle group deleted successfully");
        dashboard.loadMuscularGroupsManagement();
      } catch (error) {
        notifications.error("Error deleting muscle group");
      }
    }
  }

  getFormFields() {
    return {
      name: {
        type: "text",
        label: "Name",
        placeholder: "Enter muscle group name",
        required: true,
        help: "Enter a descriptive name for the muscle group",
      },
      body_part: {
        type: "select",
        label: "Body Part",
        placeholder: "Select body part",
        required: true,
        options: [
          { value: "", label: "Select body part" },
          { value: "upper_body", label: "Upper Body" },
          { value: "lower_body", label: "Lower Body" },
          { value: "core", label: "Core" },
          { value: "full_body", label: "Full Body" },
        ],
        help: "Select the main body part for this muscle group",
      },
      description: {
        type: "textarea",
        label: "Description",
        placeholder: "Enter description (optional)",
        required: false,
        rows: 3,
        help: "Provide additional details about this muscle group",
      },
    };
  }

  validateMuscularGroupData(data) {
    const errors = [];

    if (!data.name || data.name.trim().length === 0) {
      errors.push("Name is required");
    }

    if (!data.body_part || data.body_part.trim().length === 0) {
      errors.push("Body part is required");
    }

    if (data.name && data.name.length > 100) {
      errors.push("Name must be less than 100 characters");
    }

    if (data.description && data.description.length > 500) {
      errors.push("Description must be less than 500 characters");
    }

    const validBodyParts = ["upper_body", "lower_body", "core", "full_body"];
    if (data.body_part && !validBodyParts.includes(data.body_part)) {
      errors.push("Invalid body part selection");
    }

    return errors;
  }

  prepareDataForSubmission(formData) {
    const data = { ...formData };

    if (data.name) data.name = data.name.trim();
    if (data.description) data.description = data.description.trim();

    if (!data.description) delete data.description;

    return data;
  }
}
