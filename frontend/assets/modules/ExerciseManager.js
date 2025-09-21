/**
 * Exercise Management Module
 * Handles all exercise-related functionality
 * Uses global utilities (loaded via script tags)
 */
class ExerciseManager {
  constructor() {
    this.api = new ApiClient();
    this.currentExercises = [];
    this.muscularGroups = [];
    this.equipment = [];
    this.loadedReferenceData = false;
  }

  async loadReferenceData() {
    if (this.loadedReferenceData) return;

    try {
      const [muscularGroupsResponse, equipmentResponse] = await Promise.all([
        this.api.getMuscularGroups(),
        this.api.getEquipment(),
      ]);

      this.muscularGroups = muscularGroupsResponse?.data || [];
      this.equipment = equipmentResponse?.data || [];
      this.loadedReferenceData = true;
    } catch (error) {
      console.error("Error loading reference data:", error);
      // Don't throw here as it's not critical for basic functionality
    }
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
      // Extract linking data
      const { muscular_groups, equipment_ids, ...baseExerciseData } =
        exerciseData;

      // Create the exercise first
      const response = await this.api.createExercise(baseExerciseData);
      const exerciseId = response.data?.id;

      if (exerciseId) {
        // Create muscular group links
        if (muscular_groups && muscular_groups.length > 0) {
          await this.createMuscularGroupLinks(exerciseId, muscular_groups);
        }

        // Create equipment links
        if (equipment_ids && equipment_ids.length > 0) {
          await this.createEquipmentLinks(exerciseId, equipment_ids);
        }
      }

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
      // Extract linking data
      const { muscular_groups, equipment_ids, ...baseExerciseData } =
        exerciseData;

      // Update the exercise first
      const response = await this.api.updateExercise(id, baseExerciseData);

      // Update muscular group links
      await this.updateMuscularGroupLinks(id, muscular_groups || []);

      // Update equipment links
      await this.updateEquipmentLinks(id, equipment_ids || []);

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

  async restoreExercise(id) {
    try {
      await this.api.restoreExercise(id);
      notifications.success("Exercise restored successfully");

      await this.loadExercises();
      return true;
    } catch (error) {
      console.error("Error restoring exercise:", error);
      notifications.error("Failed to restore exercise: " + error.message);
      throw error;
    }
  }

  // Helper methods for managing exercise links
  async createMuscularGroupLinks(exerciseId, muscularGroupIds) {
    const linkPromises = muscularGroupIds.map((muscularGroupId) =>
      this.api.createExerciseMuscularGroupLink({
        exercise_id: exerciseId,
        muscular_group_id: muscularGroupId,
      })
    );
    await Promise.all(linkPromises);
  }

  async createEquipmentLinks(exerciseId, equipmentIds) {
    const linkPromises = equipmentIds.map((equipmentId) =>
      this.api.createExerciseEquipmentLink({
        exercise_id: exerciseId,
        equipment_id: equipmentId,
      })
    );
    await Promise.all(linkPromises);
  }

  async updateMuscularGroupLinks(exerciseId, newMuscularGroupIds) {
    try {
      // Get existing links
      const existingLinksResponse =
        await this.api.getExerciseMuscularGroupLinks(exerciseId);
      const existingLinks = existingLinksResponse?.data || [];
      const existingMuscularGroupIds = existingLinks.map(
        (link) => link.muscular_group_id
      );

      // Determine which links to add and remove
      const toAdd = newMuscularGroupIds.filter(
        (id) => !existingMuscularGroupIds.includes(id)
      );
      const toRemove = existingLinks.filter(
        (link) => !newMuscularGroupIds.includes(link.muscular_group_id)
      );

      // Remove old links
      const removePromises = toRemove.map((link) =>
        this.api.deleteExerciseMuscularGroupLink(link.id)
      );
      await Promise.all(removePromises);

      // Add new links
      await this.createMuscularGroupLinks(exerciseId, toAdd);
    } catch (error) {
      console.error("Error updating muscular group links:", error);
      // Don't throw here to avoid breaking the main exercise update
    }
  }

  async updateEquipmentLinks(exerciseId, newEquipmentIds) {
    try {
      // Get existing links
      const existingLinksResponse = await this.api.getExerciseEquipmentLinks(
        exerciseId
      );
      const existingLinks = existingLinksResponse?.data || [];
      const existingEquipmentIds = existingLinks.map(
        (link) => link.equipment_id
      );

      // Determine which links to add and remove
      const toAdd = newEquipmentIds.filter(
        (id) => !existingEquipmentIds.includes(id)
      );
      const toRemove = existingLinks.filter(
        (link) => !newEquipmentIds.includes(link.equipment_id)
      );

      // Remove old links
      const removePromises = toRemove.map((link) =>
        this.api.deleteExerciseEquipmentLink(link.id)
      );
      await Promise.all(removePromises);

      // Add new links
      await this.createEquipmentLinks(exerciseId, toAdd);
    } catch (error) {
      console.error("Error updating equipment links:", error);
      // Don't throw here to avoid breaking the main exercise update
    }
  }

  async getExerciseLinks(exerciseId) {
    try {
      const [muscularGroupLinksResponse, equipmentLinksResponse] =
        await Promise.all([
          this.api.getExerciseMuscularGroupLinks(exerciseId),
          this.api.getExerciseEquipmentLinks(exerciseId),
        ]);

      return {
        muscularGroups: muscularGroupLinksResponse?.data || [],
        equipment: equipmentLinksResponse?.data || [],
      };
    } catch (error) {
      console.error("Error loading exercise links:", error);
      return { muscularGroups: [], equipment: [] };
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
        handler: (exercise) => this.requestDeleteExercise(exercise),
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

  async requestDeleteExercise(exercise) {
    document.dispatchEvent(
      new CustomEvent("exercise:delete", {
        detail: { exercise },
      })
    );
  }

  truncateText(text, maxLength) {
    if (!text || text.length <= maxLength) return text;
    return text.substring(0, maxLength) + "...";
  }

  async getFormSchema() {
    // Ensure reference data is loaded
    await this.loadReferenceData();

    return {
      name: {
        type: "text",
        label: "Exercise Name",
        required: true,
        placeholder: "Enter exercise name",
      },
      category: {
        type: "select",
        label: "Category",
        required: true,
        options: [
          { value: "", label: "Select category" },
          { value: "cardio", label: "Cardio" },
          { value: "strength", label: "Strength" },
          { value: "flexibility", label: "Flexibility" },
          { value: "balance", label: "Balance" },
          { value: "sports", label: "Sports" },
        ],
      },
      muscular_groups: {
        type: "multiselect",
        label: "Muscle Groups",
        placeholder: "Select muscle groups",
        options: [
          { value: "", label: "Select muscle groups" },
          ...this.muscularGroups.map((mg) => ({
            value: mg.id,
            label: mg.name,
          })),
        ],
        help: "Select one or more muscle groups targeted by this exercise",
      },
      equipment_ids: {
        type: "multiselect",
        label: "Equipment",
        placeholder: "Select equipment (optional)",
        options: [
          { value: "", label: "Select equipment" },
          ...this.equipment.map((eq) => ({
            value: eq.id,
            label: eq.name,
          })),
        ],
        help: "Select equipment needed for this exercise (leave empty for bodyweight exercises)",
      },
      difficulty: {
        type: "select",
        label: "Difficulty",
        options: [
          { value: "", label: "Select difficulty" },
          { value: "beginner", label: "Beginner" },
          { value: "intermediate", label: "Intermediate" },
          { value: "advanced", label: "Advanced" },
        ],
      },
      description: {
        type: "textarea",
        label: "Description",
        rows: 3,
        placeholder: "Enter exercise description (optional)",
      },
      instructions: {
        type: "textarea",
        label: "Instructions",
        rows: 4,
        placeholder: "Enter detailed exercise instructions",
      },
    };
  }

  validateExerciseData(data) {
    const errors = {};

    if (!data.name || data.name.trim() === "") {
      errors.name = "Exercise name is required";
    }

    if (!data.category) {
      errors.category = "Category is required";
    }

    // Parse muscle groups if provided as string
    if (data.muscle_groups && typeof data.muscle_groups === "string") {
      data.muscle_groups = data.muscle_groups
        .split(",")
        .map((m) => m.trim())
        .filter((m) => m.length > 0);
    }

    if (data.description && data.description.length > 500) {
      errors.description = "Description must be less than 500 characters";
    }

    if (data.instructions && data.instructions.length > 1000) {
      errors.instructions = "Instructions must be less than 1000 characters";
    }

    return {
      isValid: Object.keys(errors).length === 0,
      errors,
    };
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

// Make ExerciseManager globally available
window.ExerciseManager = ExerciseManager;
