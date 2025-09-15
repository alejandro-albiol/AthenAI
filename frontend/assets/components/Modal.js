/**
 * Modal Component
 * Reusable modal with standard structure and behavior
 */
class Modal extends window.BaseComponent {
  constructor(options = {}) {
    super(null, options);
    this.createModal();
  }

  getDefaultOptions() {
    return {
      ...super.getDefaultOptions(),
      title: "Modal",
      content: "",
      size: "medium", // small, medium, large
      closeable: true,
      backdrop: true,
      keyboard: true,
      appendTo: document.body,
      className: "modal",
      buttons: [],
    };
  }

  createModal() {
    this.element = document.createElement("div");
    this.element.className = `${this.options.className} modal-${this.options.size}`;
    this.element.innerHTML = this.getTemplate();

    if (this.options.appendTo) {
      this.options.appendTo.appendChild(this.element);
    }

    this.init();
  }

  getTemplate() {
    return `
      <div class="modal-backdrop" ${
        this.options.backdrop ? 'data-dismiss="modal"' : ""
      }></div>
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">${this.options.title}</h3>
          ${
            this.options.closeable
              ? '<button type="button" class="modal-close" data-dismiss="modal">&times;</button>'
              : ""
          }
        </div>
        <div class="modal-body">
          ${this.options.content}
        </div>
        ${this.options.buttons.length > 0 ? this.renderButtons() : ""}
      </div>
    `;
  }

  renderButtons() {
    const buttonsHtml = this.options.buttons
      .map((button) => {
        const className = button.className || "btn btn-secondary";
        const action = button.action || "dismiss";
        return `<button type="button" class="${className}" data-action="${action}">${button.text}</button>`;
      })
      .join("");

    return `
      <div class="modal-footer">
        ${buttonsHtml}
      </div>
    `;
  }

  setupEventListeners() {
    // Close modal handlers
    this.addEventListener("click", '[data-dismiss="modal"]', () => this.hide());

    // Button actions
    this.addEventListener("click", "[data-action]", (e, target) => {
      const action = target.getAttribute("data-action");
      const button = this.options.buttons.find((btn) => btn.action === action);

      if (button && button.handler) {
        const result = button.handler.call(this, e);
        if (result !== false && action !== "custom") {
          this.hide();
        }
      } else if (action === "dismiss") {
        this.hide();
      }
    });

    // Keyboard support
    if (this.options.keyboard) {
      document.addEventListener("keydown", this.handleKeydown.bind(this));
    }
  }

  handleKeydown(e) {
    if (e.key === "Escape" && this.isVisible()) {
      this.hide();
    }
  }

  show() {
    super.show();
    this.element.style.display = "flex";
    document.body.style.overflow = "hidden";

    // Animation
    setTimeout(() => {
      this.addClass("modal-show");
    }, 10);

    this.trigger("show");
  }

  hide() {
    this.removeClass("modal-show");

    setTimeout(() => {
      super.hide();
      document.body.style.overflow = "";
      this.trigger("hide");
    }, 300);
  }

  isVisible() {
    return this.element && this.element.style.display !== "none";
  }

  setTitle(title) {
    const titleElement = this.element.querySelector(".modal-title");
    if (titleElement) {
      titleElement.textContent = title;
    }
  }

  setContent(content) {
    const bodyElement = this.element.querySelector(".modal-body");
    if (bodyElement) {
      bodyElement.innerHTML = content;
    }
  }

  updateButtons(buttons) {
    this.options.buttons = buttons;
    const footerElement = this.element.querySelector(".modal-footer");
    if (footerElement) {
      footerElement.outerHTML = this.renderButtons();
    } else if (buttons.length > 0) {
      const contentElement = this.element.querySelector(".modal-content");
      contentElement.insertAdjacentHTML("beforeend", this.renderButtons());
    }
  }

  trigger(eventName, data = {}) {
    const event = new CustomEvent(`modal:${eventName}`, {
      detail: { modal: this, ...data },
    });
    this.element.dispatchEvent(event);
  }

  destroy() {
    if (this.options.keyboard) {
      document.removeEventListener("keydown", this.handleKeydown.bind(this));
    }

    if (this.element && this.element.parentNode) {
      this.element.parentNode.removeChild(this.element);
    }

    super.destroy();
    document.body.style.overflow = "";
  }
}

// Make available globally for vanilla JS
window.Modal = Modal;
