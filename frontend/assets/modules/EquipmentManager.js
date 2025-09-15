import { ApiClient } from "../utils/api.js";
import notifications from "../utils/notifications.js";
import { appState } from "../utils/state.js";

/**
 * Equipment Management Module
 * Handles all equipment-related functionality
 */
class EquipmentManager {
  constructor() {
    this.api = new ApiClient();
    this.currentEquipment = [];
  }

  async loadEquipment() {
    try {
      appState.setState({ loading: true });
      const response = await this.api.getEquipment();
      const equipment = response?.data || [];

      this.currentEquipment = equipment;
      appState.setState({ equipment, loading: false });

      return equipment;
    } catch (error) {
      console.error("Error loading equipment:", error);
      appState.setState({ loading: false, error: error.message });
      notifications.error("Failed to load equipment");
      throw error;
    }
  }

  async createEquipment(equipmentData) {
    try {
      const response = await this.api.createEquipment(equipmentData);
      notifications.success("Equipment created successfully");

      // Refresh equipment list
      await this.loadEquipment();

      return response;
    } catch (error) {
      console.error("Error creating equipment:", error);
      notifications.error("Failed to create equipment: " + error.message);
      throw error;
    }
  }

  async updateEquipment(id, equipmentData) {
    try {
      const response = await this.api.updateEquipment(id, equipmentData);
      notifications.success("Equipment updated successfully");

      // Refresh equipment list
      await this.loadEquipment();

      return response;
    } catch (error) {
      console.error("Error updating equipment:", error);
      notifications.error("Failed to update equipment: " + error.message);
      throw error;
    }
  }

  async deleteEquipment(id) {
    try {
      await this.api.deleteEquipment(id);
      notifications.success("Equipment deleted successfully");

      // Refresh equipment list
      await this.loadEquipment();

      return true;
    } catch (error) {
      console.error("Error deleting equipment:", error);
      notifications.error("Failed to delete equipment: " + error.message);
      throw error;
    }
  }

  getEquipmentIcon(category) {
    const icons = {
      cardio: '<i class="fas fa-heartbeat" style="color: #e74c3c;"></i>',
      free_weights:
        '<i class="fas fa-weight-hanging" style="color: #2c3e50;"></i>',
      machines: '<i class="fas fa-cogs" style="color: #7f8c8d;"></i>',
      accessories: '<i class="fas fa-tools" style="color: #f39c12;"></i>',
      bodyweight: '<i class="fas fa-male" style="color: #27ae60;"></i>',
    };
    return (
      icons[category] ||
      '<i class="fas fa-question-circle" style="color: #95a5a6;"></i>'
    );
  }

  formatCategory(category) {
    const formatted = {
      cardio: "Cardio",
      free_weights: "Free Weights",
      machines: "Machines",
      accessories: "Accessories",
      bodyweight: "Bodyweight",
    };
    return formatted[category] || category;
  }

  getTableColumns() {
    return [
      {
        key: "name",
        title: "Name",
        sortable: true,
        render: (value, item) => `
          <div class="equipment-info">
            <div class="equipment-icon-inline">${this.getEquipmentIcon(
              item.category
            )}</div>
            <div>
              <strong>${value}</strong>
              <div class="equipment-category-small">${this.formatCategory(
                item.category
              )}</div>
            </div>
          </div>
        `,
      },
      {
        key: "category",
        title: "Category",
        sortable: true,
        render: (value) => this.formatCategory(value),
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
        title: "Edit Equipment",
        className: "btn btn-sm btn-outline",
        handler: (equipment) => this.editEquipment(equipment),
      },
      {
        action: "delete",
        icon: "fas fa-trash",
        title: "Delete Equipment",
        className: "btn btn-sm btn-danger",
        handler: (equipment) => this.confirmDeleteEquipment(equipment),
      },
    ];
  }

  async editEquipment(equipment) {
    // This would typically open a modal with pre-filled form
    // For now, we'll emit an event that the main dashboard can handle
    document.dispatchEvent(
      new CustomEvent("equipment:edit", {
        detail: { equipment },
      })
    );
  }

  async confirmDeleteEquipment(equipment) {
    if (confirm(`Are you sure you want to delete "${equipment.name}"?`)) {
      await this.deleteEquipment(equipment.id);
    }
  }

  truncateText(text, maxLength) {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + "...";
  }

  // Generate equipment form data
  getFormSchema() {
    return {
      name: {
        type: "text",
        label: "Equipment Name",
        required: true,
        placeholder: "Enter equipment name",
      },
      category: {
        type: "select",
        label: "Category",
        required: true,
        options: [
          { value: "", label: "Select category..." },
          { value: "cardio", label: "Cardio" },
          { value: "free_weights", label: "Free Weights" },
          { value: "machines", label: "Machines" },
          { value: "accessories", label: "Accessories" },
          { value: "bodyweight", label: "Bodyweight" },
        ],
      },
      description: {
        type: "textarea",
        label: "Description",
        rows: 3,
        placeholder: "Enter equipment description (optional)",
      },
    };
  }

  validateEquipmentData(data) {
    const errors = {};

    if (!data.name || data.name.trim() === "") {
      errors.name = "Equipment name is required";
    }

    if (!data.category) {
      errors.category = "Category is required";
    }

    if (data.description && data.description.length > 500) {
      errors.description = "Description must be less than 500 characters";
    }

    return {
      isValid: Object.keys(errors).length === 0,
      errors,
    };
  }

  // Get equipment statistics
  getStats() {
    const equipment = this.currentEquipment;
    const categoryStats = {};

    equipment.forEach((item) => {
      categoryStats[item.category] = (categoryStats[item.category] || 0) + 1;
    });

    return {
      total: equipment.length,
      categories: Object.keys(categoryStats).length,
      categoryBreakdown: categoryStats,
      mostCommonCategory: Object.keys(categoryStats).reduce(
        (a, b) => (categoryStats[a] > categoryStats[b] ? a : b),
        "none"
      ),
    };
  }
}

export default EquipmentManager;
