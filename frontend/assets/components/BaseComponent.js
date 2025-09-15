/**
 * Base Component Class
 * Provides common functionality for all UI components
 */
class BaseComponent {
  constructor(selector, options = {}) {
    this.selector = selector;
    this.element =
      typeof selector === "string"
        ? document.querySelector(selector)
        : selector;
    this.options = { ...this.getDefaultOptions(), ...options };
    this.eventHandlers = new Map();

    if (this.element) {
      this.init();
    }
  }

  getDefaultOptions() {
    return {
      className: "",
      template: "",
      data: {},
    };
  }

  init() {
    this.setupEventListeners();
    this.render();
  }

  render() {
    if (this.options.template && this.element) {
      this.element.innerHTML = this.options.template;
    }
  }

  setupEventListeners() {
    // Override in child classes
  }

  addEventListener(event, selector, handler) {
    const wrappedHandler = (e) => {
      const target = e.target.closest(selector);
      if (target) {
        handler.call(this, e, target);
      }
    };

    this.element.addEventListener(event, wrappedHandler);
    this.eventHandlers.set(`${event}-${selector}`, wrappedHandler);
  }

  removeEventListener(event, selector) {
    const handler = this.eventHandlers.get(`${event}-${selector}`);
    if (handler) {
      this.element.removeEventListener(event, handler);
      this.eventHandlers.delete(`${event}-${selector}`);
    }
  }

  destroy() {
    // Remove all event listeners
    this.eventHandlers.forEach((handler, key) => {
      const [event] = key.split("-");
      this.element.removeEventListener(event, handler);
    });
    this.eventHandlers.clear();

    // Clear element
    if (this.element) {
      this.element.innerHTML = "";
    }
  }

  // Utility methods
  show() {
    if (this.element) {
      this.element.style.display = "block";
    }
  }

  hide() {
    if (this.element) {
      this.element.style.display = "none";
    }
  }

  toggle() {
    if (this.element) {
      const isVisible = this.element.style.display !== "none";
      this.element.style.display = isVisible ? "none" : "block";
    }
  }

  addClass(className) {
    if (this.element) {
      this.element.classList.add(className);
    }
  }

  removeClass(className) {
    if (this.element) {
      this.element.classList.remove(className);
    }
  }

  hasClass(className) {
    return this.element ? this.element.classList.contains(className) : false;
  }

  setData(key, value) {
    this.options.data[key] = value;
  }

  getData(key) {
    return this.options.data[key];
  }

  updateTemplate(template) {
    this.options.template = template;
    this.render();
  }
}

// Make available globally for vanilla JS
window.BaseComponent = BaseComponent;
