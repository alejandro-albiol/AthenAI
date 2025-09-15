/**
 * Form Utility Module
 * Reusable form handling with validation and error display
 */
class FormHandler {
  constructor(formSelector, options = {}) {
    this.form =
      typeof formSelector === "string"
        ? document.querySelector(formSelector)
        : formSelector;

    this.options = {
      validateOnSubmit: true,
      validateOnBlur: false,
      showErrors: true,
      clearErrorsOnFocus: true,
      ...options,
    };

    this.validators = new Map();
    this.errors = new Map();

    if (this.form) {
      this.init();
    }
  }

  init() {
    this.setupEventListeners();
  }

  setupEventListeners() {
    // Form submission
    this.form.addEventListener("submit", (e) => {
      if (this.options.validateOnSubmit) {
        e.preventDefault();
        this.handleSubmit(e);
      }
    });

    // Field validation on blur
    if (this.options.validateOnBlur) {
      this.form.addEventListener(
        "blur",
        (e) => {
          if (e.target.matches("input, select, textarea")) {
            this.validateField(e.target);
          }
        },
        true
      );
    }

    // Clear errors on focus
    if (this.options.clearErrorsOnFocus) {
      this.form.addEventListener(
        "focus",
        (e) => {
          if (e.target.matches("input, select, textarea")) {
            this.clearFieldError(e.target);
          }
        },
        true
      );
    }
  }

  async handleSubmit(e) {
    const isValid = await this.validate();

    if (isValid) {
      if (this.options.onSubmit) {
        try {
          await this.options.onSubmit(this.getFormData(), this.form);
        } catch (error) {
          this.handleSubmitError(error);
        }
      }
    }
  }

  async validate() {
    this.clearAllErrors();
    let isValid = true;

    const fields = this.form.querySelectorAll("input, select, textarea");

    for (const field of fields) {
      const fieldValid = await this.validateField(field);
      if (!fieldValid) {
        isValid = false;
      }
    }

    return isValid;
  }

  async validateField(field) {
    const name = field.name;
    const value = field.value;
    const validators = this.validators.get(name) || [];

    for (const validator of validators) {
      try {
        const result = await validator(value, field, this.getFormData());
        if (result !== true) {
          this.setFieldError(field, result);
          return false;
        }
      } catch (error) {
        this.setFieldError(field, error.message);
        return false;
      }
    }

    this.clearFieldError(field);
    return true;
  }

  addValidator(fieldName, validatorFn) {
    if (!this.validators.has(fieldName)) {
      this.validators.set(fieldName, []);
    }
    this.validators.get(fieldName).push(validatorFn);
  }

  setFieldError(field, message) {
    const name = field.name;
    this.errors.set(name, message);

    if (this.options.showErrors) {
      this.displayFieldError(field, message);
    }
  }

  clearFieldError(field) {
    const name = field.name;
    this.errors.delete(name);
    this.removeFieldErrorDisplay(field);
  }

  clearAllErrors() {
    this.errors.clear();
    const errorElements = this.form.querySelectorAll(".field-error");
    errorElements.forEach((el) => el.remove());

    const errorFields = this.form.querySelectorAll(".field-invalid");
    errorFields.forEach((el) => el.classList.remove("field-invalid"));
  }

  displayFieldError(field, message) {
    // Remove existing error
    this.removeFieldErrorDisplay(field);

    // Add error class to field
    field.classList.add("field-invalid");

    // Create error element
    const errorElement = document.createElement("div");
    errorElement.className = "field-error";
    errorElement.textContent = message;

    // Insert error element after field or field container
    const container = field.closest(".form-group") || field.parentElement;
    container.appendChild(errorElement);
  }

  removeFieldErrorDisplay(field) {
    field.classList.remove("field-invalid");

    const container = field.closest(".form-group") || field.parentElement;
    const errorElement = container.querySelector(".field-error");
    if (errorElement) {
      errorElement.remove();
    }
  }

  handleSubmitError(error) {
    if (error.isValidationError && error.isValidationError()) {
      const validationErrors = error.getValidationErrors();
      Object.keys(validationErrors).forEach((fieldName) => {
        const field = this.form.querySelector(`[name="${fieldName}"]`);
        if (field) {
          this.setFieldError(field, validationErrors[fieldName]);
        }
      });
    } else {
      // Show general error
      if (this.options.onError) {
        this.options.onError(error);
      }
    }
  }

  getFormData() {
    const formData = new FormData(this.form);
    const data = {};

    for (const [key, value] of formData.entries()) {
      if (data[key]) {
        // Handle multiple values (checkboxes, etc.)
        if (Array.isArray(data[key])) {
          data[key].push(value);
        } else {
          data[key] = [data[key], value];
        }
      } else {
        data[key] = value;
      }
    }

    return data;
  }

  setFormData(data) {
    Object.keys(data).forEach((key) => {
      const field = this.form.querySelector(`[name="${key}"]`);
      if (field) {
        if (field.type === "checkbox") {
          field.checked = Boolean(data[key]);
        } else {
          field.value = data[key] || "";
        }
      }
    });
  }

  reset() {
    this.form.reset();
    this.clearAllErrors();
  }

  isValid() {
    return this.errors.size === 0;
  }

  getErrors() {
    return Object.fromEntries(this.errors);
  }
}

// Common validators
const validators = {
  required: (message = "This field is required") => {
    return (value) => {
      if (!value || value.trim() === "") {
        return message;
      }
      return true;
    };
  },

  email: (message = "Please enter a valid email address") => {
    return (value) => {
      if (value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
        return message;
      }
      return true;
    };
  },

  minLength: (min, message) => {
    return (value) => {
      if (value && value.length < min) {
        return message || `Must be at least ${min} characters`;
      }
      return true;
    };
  },

  maxLength: (max, message) => {
    return (value) => {
      if (value && value.length > max) {
        return message || `Must be no more than ${max} characters`;
      }
      return true;
    };
  },

  pattern: (regex, message = "Invalid format") => {
    return (value) => {
      if (value && !regex.test(value)) {
        return message;
      }
      return true;
    };
  },

  custom: (validatorFn) => validatorFn,
};

// Global utility function for easy form data extraction
function getFormData(form) {
  const formData = new FormData(form);
  const data = {};

  for (let [key, value] of formData.entries()) {
    if (data[key]) {
      // Handle multiple values (checkboxes, etc.)
      if (Array.isArray(data[key])) {
        data[key].push(value);
      } else {
        data[key] = [data[key], value];
      }
    } else {
      data[key] = value;
    }
  }

  return data;
}

// Make available globally for vanilla JS
window.FormHandler = FormHandler;
window.validators = validators;
window.getFormData = getFormData;
