import { ApiClient } from "../utils/api.js";
import notifications from "../utils/notifications.js";
import { appState } from "../utils/state.js";

/**
 * Exercise Management Module
 * Handles all exercise-related functionality
 */
class ExerciseManager {
  constructor() {
    this.api = new ApiClient();
    this.currentExercises = [];
  }

  async loadExercises() {
    try {
      appState.setState({ loading: true });
      const response = await this.api.getExercises();
      const exercises = response?.data || [];

      this.currentExercises = exercises;
      appState.setState({ exercises, loading: false });

      return exercises;
    } catch (error) {
      console.error("Error loading exercises:", error);
      appState.setState({ loading: false, error: error.message });
      notifications.error("Failed to load exercises");
      throw error;
    }
  }

  async createExercise(exerciseData) {
    try {
      const response = await this.api.createExercise(exerciseData);
      notifications.success("Exercise created successfully");

      await this.loadExercises();
      return response;
    } catch (error) {
      console.error("Error creating exercise:", error);
      notifications.error("Failed to create exercise: " + error.message);
      throw error;
    }
  }

  async updateExercise(id, exerciseData) {
    try {
      const response = await this.api.updateExercise(id, exerciseData);
      notifications.success("Exercise updated successfully");

      await this.loadExercises();
      return response;
    } catch (error) {
      console.error("Error updating exercise:", error);
      notifications.error("Failed to update exercise: " + error.message);
      throw error;
    }
  }

  async deleteExercise(id) {
    try {
      await this.api.deleteExercise(id);
      notifications.success("Exercise deleted successfully");

      await this.loadExercises();
      return true;
    } catch (error) {
      console.error("Error deleting exercise:", error);
      notifications.error("Failed to delete exercise: " + error.message);
      throw error;
    }
  }

  getExerciseTypeIcon(type) {
    const icons = {
      strength: '<i class="fas fa-dumbbell" style="color: #2c3e50;"></i>',
      cardio: '<i class="fas fa-heartbeat" style="color: #e74c3c;"></i>',
      flexibility: '<i class="fas fa-leaf" style="color: #27ae60;"></i>',
      balance: '<i class="fas fa-balance-scale" style="color: #f39c12;"></i>',
      plyometrics: '<i class="fas fa-bolt" style="color: #9b59b6;"></i>',
      endurance: '<i class="fas fa-clock" style="color: #3498db;"></i>',
    };
    return (
      icons[type] ||
      '<i class="fas fa-question-circle" style="color: #95a5a6;"></i>'
    );
  }

  formatExerciseType(type) {
    const formatted = {
      strength: "Strength",
      cardio: "Cardio",
      flexibility: "Flexibility",
      balance: "Balance",
      plyometrics: "Plyometrics",
      endurance: "Endurance",
    };
    return formatted[type] || type;
  }

  formatDifficulty(difficulty) {
    const levels = {
      beginner: { label: "Beginner", color: "#27ae60" },
      intermediate: { label: "Intermediate", color: "#f39c12" },
      advanced: { label: "Advanced", color: "#e74c3c" },
      expert: { label: "Expert", color: "#8e44ad" },
    };

    const level = levels[difficulty] || { label: difficulty, color: "#95a5a6" };
    return `<span class="difficulty-badge" style="color: ${level.color};">${level.label}</span>`;
  }

  getTableColumns() {
    return [
      {
        key: "name",
        title: "Exercise",
        sortable: true,
        render: (value, item) => `
          <div class="exercise-info">
            <div class="exercise-icon-inline">${this.getExerciseTypeIcon(
              item.type
            )}</div>
            <div>
              <strong>${value}</strong>
              <div class="exercise-type-small">${this.formatExerciseType(
                item.type
              )}</div>
            </div>
          </div>
        `,
      },
      {
        key: "difficulty",
        title: "Difficulty",
        sortable: true,
        render: (value) => this.formatDifficulty(value),
      },
      {
        key: "muscle_groups",
        title: "Muscle Groups",
        render: (value) => {
          if (!value || value.length === 0) return "None specified";
          const groups = Array.isArray(value) ? value : [value];
          return (
            groups.slice(0, 2).join(", ") +
            (groups.length > 2 ? ` +${groups.length - 2} more` : "")
          );
        },
      },
      {
        key: "instructions",
        title: "Instructions",
        render: (value) =>
          value
            ? `<span title="${value}">${this.truncateText(value, 40)}</span>`
            : "No instructions",
      },
    ];
  }

  getRowActions() {
    return [
      {
        action: "view",
        icon: "fas fa-eye",
        title: "View Details",
        className: "btn btn-sm btn-info",
        handler: (exercise) => this.viewExerciseDetails(exercise),
      },
      {
        action: "edit",
        icon: "fas fa-edit",
        title: "Edit Exercise",
        className: "btn btn-sm btn-outline",
        handler: (exercise) => this.editExercise(exercise),
      },
      {
        action: "delete",
        icon: "fas fa-trash",
        title: "Delete Exercise",
        className: "btn btn-sm btn-danger",
        handler: (exercise) => this.confirmDeleteExercise(exercise),
      },
    ];
  }

  async viewExerciseDetails(exercise) {
    document.dispatchEvent(
      new CustomEvent("exercise:view", {
        detail: { exercise },
      })
    );
  }

  async editExercise(exercise) {
    document.dispatchEvent(
      new CustomEvent("exercise:edit", {
        detail: { exercise },
      })
    );
  }

  async confirmDeleteExercise(exercise) {
    if (confirm(`Are you sure you want to delete "${exercise.name}"?`)) {
      await this.deleteExercise(exercise.id);
    }
  }

  truncateText(text, maxLength) {
    if (!text || text.length <= maxLength) return text;
    return text.substring(0, maxLength) + "...";
  }

  getFormSchema() {
    return {
      name: {
        type: "text",
        label: "Exercise Name",
        required: true,
        placeholder: "Enter exercise name",
      },
      type: {
        type: "select",
        label: "Exercise Type",
        required: true,
        options: [
          { value: "", label: "Select type..." },
          { value: "strength", label: "Strength" },
          { value: "cardio", label: "Cardio" },
          { value: "flexibility", label: "Flexibility" },
          { value: "balance", label: "Balance" },
          { value: "plyometrics", label: "Plyometrics" },
          { value: "endurance", label: "Endurance" },
        ],
      },
      difficulty: {
        type: "select",
        label: "Difficulty Level",
        required: true,
        options: [
          { value: "", label: "Select difficulty..." },
          { value: "beginner", label: "Beginner" },
          { value: "intermediate", label: "Intermediate" },
          { value: "advanced", label: "Advanced" },
          { value: "expert", label: "Expert" },
        ],
      },
      instructions: {
        type: "textarea",
        label: "Instructions",
        rows: 4,
        placeholder: "Enter detailed exercise instructions",
      },
      equipment_ids: {
        type: "multiselect",
        label: "Required Equipment",
        placeholder: "Select equipment (optional)",
      },
      muscle_group_ids: {
        type: "multiselect",
        label: "Target Muscle Groups",
        placeholder: "Select muscle groups",
      },
      video_url: {
        type: "url",
        label: "Video URL",
        placeholder: "https://example.com/video (optional)",
      },
    };
  }

  validateExerciseData(data) {
    const errors = {};

    if (!data.name || data.name.trim() === "") {
      errors.name = "Exercise name is required";
    }

    if (!data.type) {
      errors.type = "Exercise type is required";
    }

    if (!data.difficulty) {
      errors.difficulty = "Difficulty level is required";
    }

    if (data.video_url && !this.isValidUrl(data.video_url)) {
      errors.video_url = "Please enter a valid URL";
    }

    if (data.instructions && data.instructions.length > 1000) {
      errors.instructions = "Instructions must be less than 1000 characters";
    }

    return {
      isValid: Object.keys(errors).length === 0,
      errors,
    };
  }

  isValidUrl(string) {
    try {
      new URL(string);
      return true;
    } catch (_) {
      return false;
    }
  }

  getStats() {
    const exercises = this.currentExercises;
    const typeStats = {};
    const difficultyStats = {};

    exercises.forEach((exercise) => {
      typeStats[exercise.type] = (typeStats[exercise.type] || 0) + 1;
      difficultyStats[exercise.difficulty] =
        (difficultyStats[exercise.difficulty] || 0) + 1;
    });

    return {
      total: exercises.length,
      types: Object.keys(typeStats).length,
      typeBreakdown: typeStats,
      difficultyBreakdown: difficultyStats,
      mostCommonType: Object.keys(typeStats).reduce(
        (a, b) => (typeStats[a] > typeStats[b] ? a : b),
        "none"
      ),
      averageDifficulty: this.getAverageDifficulty(exercises),
    };
  }

  getAverageDifficulty(exercises) {
    if (exercises.length === 0) return "N/A";

    const difficultyValues = {
      beginner: 1,
      intermediate: 2,
      advanced: 3,
      expert: 4,
    };

    const total = exercises.reduce((sum, exercise) => {
      return sum + (difficultyValues[exercise.difficulty] || 0);
    }, 0);

    const average = total / exercises.length;

    if (average <= 1.5) return "Beginner";
    if (average <= 2.5) return "Intermediate";
    if (average <= 3.5) return "Advanced";
    return "Expert";
  }

  // Filter exercises by criteria
  filterExercises(criteria) {
    let filtered = [...this.currentExercises];

    if (criteria.type && criteria.type !== "all") {
      filtered = filtered.filter((exercise) => exercise.type === criteria.type);
    }

    if (criteria.difficulty && criteria.difficulty !== "all") {
      filtered = filtered.filter(
        (exercise) => exercise.difficulty === criteria.difficulty
      );
    }

    if (criteria.muscleGroup) {
      filtered = filtered.filter(
        (exercise) =>
          exercise.muscle_groups &&
          exercise.muscle_groups.includes(criteria.muscleGroup)
      );
    }

    if (criteria.search) {
      const searchTerm = criteria.search.toLowerCase();
      filtered = filtered.filter(
        (exercise) =>
          exercise.name.toLowerCase().includes(searchTerm) ||
          (exercise.instructions &&
            exercise.instructions.toLowerCase().includes(searchTerm))
      );
    }

    return filtered;
  }
}

export default ExerciseManager;
