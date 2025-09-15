import BaseComponent from "./BaseComponent.js";

/**
 * Card Component
 * Reusable card container with header, body, and footer
 */
class Card extends BaseComponent {
  constructor(selector, options = {}) {
    super(selector, options);
  }

  getDefaultOptions() {
    return {
      ...super.getDefaultOptions(),
      title: "",
      subtitle: "",
      icon: "",
      headerActions: [],
      className: "dashboard-card",
      content: "",
      footer: "",
      variant: "default", // default, primary, warning, success, error
    };
  }

  render() {
    if (!this.element) return;

    this.element.className = `${this.options.className} card-${this.options.variant}`;
    this.element.innerHTML = this.getTemplate();
  }

  getTemplate() {
    return `
      ${
        this.options.title ||
        this.options.icon ||
        this.options.headerActions.length > 0
          ? this.getHeaderTemplate()
          : ""
      }
      <div class="card-body">
        ${this.options.content}
      </div>
      ${
        this.options.footer
          ? `<div class="card-footer">${this.options.footer}</div>`
          : ""
      }
    `;
  }

  getHeaderTemplate() {
    return `
      <div class="card-header">
        <div class="card-header-content">
          ${
            this.options.icon
              ? `<div class="card-icon">${this.options.icon}</div>`
              : ""
          }
          <div class="card-title-section">
            ${
              this.options.title
                ? `<h3 class="card-title">${this.options.title}</h3>`
                : ""
            }
            ${
              this.options.subtitle
                ? `<p class="card-subtitle">${this.options.subtitle}</p>`
                : ""
            }
          </div>
        </div>
        ${
          this.options.headerActions.length > 0 ? this.getActionsTemplate() : ""
        }
      </div>
    `;
  }

  getActionsTemplate() {
    return `
      <div class="card-actions">
        ${this.options.headerActions
          .map(
            (action) => `
          <button class="${action.className || "btn btn-sm btn-outline"}" 
                  data-action="${action.action}" 
                  title="${action.title || action.text}">
            ${action.icon ? `<i class="${action.icon}"></i>` : action.text}
          </button>
        `
          )
          .join("")}
      </div>
    `;
  }

  setupEventListeners() {
    this.addEventListener("click", "[data-action]", (e, target) => {
      const action = target.getAttribute("data-action");
      const actionConfig = this.options.headerActions.find(
        (a) => a.action === action
      );

      if (actionConfig && actionConfig.handler) {
        actionConfig.handler.call(this, e, target);
      }

      this.trigger("action", { action, target });
    });
  }

  setTitle(title) {
    this.options.title = title;
    const titleElement = this.element.querySelector(".card-title");
    if (titleElement) {
      titleElement.textContent = title;
    }
  }

  setContent(content) {
    this.options.content = content;
    const bodyElement = this.element.querySelector(".card-body");
    if (bodyElement) {
      bodyElement.innerHTML = content;
    }
  }

  setIcon(icon) {
    this.options.icon = icon;
    const iconElement = this.element.querySelector(".card-icon");
    if (iconElement) {
      iconElement.innerHTML = icon;
    }
  }

  addAction(action) {
    this.options.headerActions.push(action);
    this.render();
  }

  removeAction(actionName) {
    this.options.headerActions = this.options.headerActions.filter(
      (a) => a.action !== actionName
    );
    this.render();
  }

  trigger(eventName, data = {}) {
    const event = new CustomEvent(`card:${eventName}`, {
      detail: { card: this, ...data },
    });
    this.element.dispatchEvent(event);
  }
}

export default Card;
