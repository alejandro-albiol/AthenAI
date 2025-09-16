/**
 * Grid Component
 * Reusable grid layout for displaying cards or items
 * Extends BaseComponent (loaded via script tag)
 */
class Grid extends BaseComponent {
  constructor(selector, options = {}) {
    super(selector, options);
  }

  getDefaultOptions() {
    return {
      ...super.getDefaultOptions(),
      items: [],
      columns: "auto-fit",
      minColumnWidth: "300px",
      gap: "20px",
      itemTemplate: null,
      itemRenderer: null,
      emptyMessage: "No items to display",
      emptyIcon: "fas fa-inbox",
      className: "grid-container",
    };
  }

  render() {
    if (!this.element) return;

    this.element.className = this.options.className;
    this.element.style.cssText = `
      display: grid;
      grid-template-columns: repeat(${this.options.columns}, minmax(${this.options.minColumnWidth}, 1fr));
      gap: ${this.options.gap};
      width: 100%;
    `;

    if (this.options.items.length === 0) {
      this.renderEmptyState();
    } else {
      this.renderItems();
    }
  }

  renderItems() {
    const itemsHtml = this.options.items
      .map((item) => this.renderItem(item))
      .join("");
    this.element.innerHTML = itemsHtml;
  }

  renderItem(item) {
    if (this.options.itemRenderer) {
      return this.options.itemRenderer(item);
    }

    if (this.options.itemTemplate) {
      return this.templateReplace(this.options.itemTemplate, item);
    }

    return `<div class="grid-item">${JSON.stringify(item)}</div>`;
  }

  renderEmptyState() {
    this.element.innerHTML = `
      <div class="empty-state" style="grid-column: 1 / -1;">
        <i class="${this.options.emptyIcon}"></i>
        <h3>No Items</h3>
        <p>${this.options.emptyMessage}</p>
      </div>
    `;
  }

  templateReplace(template, data) {
    return template.replace(/\{\{(\w+(?:\.\w+)*)\}\}/g, (match, path) => {
      const value = this.getNestedValue(data, path);
      return value !== undefined ? value : match;
    });
  }

  getNestedValue(obj, path) {
    return path.split(".").reduce((current, key) => {
      return current && current[key] !== undefined ? current[key] : undefined;
    }, obj);
  }

  setupEventListeners() {
    this.addEventListener("click", ".grid-item", (e, target) => {
      const index = Array.from(this.element.children).indexOf(target);
      const item = this.options.items[index];

      if (item) {
        this.trigger("itemClick", { item, index, target });
      }
    });

    this.addEventListener("click", "[data-action]", (e, target) => {
      const action = target.getAttribute("data-action");
      const itemElement = target.closest(".grid-item");
      const index = Array.from(this.element.children).indexOf(itemElement);
      const item = this.options.items[index];

      if (item) {
        this.trigger("itemAction", { action, item, index, target });
      }
    });
  }

  updateItems(items) {
    this.options.items = items;
    this.render();
  }

  addItem(item) {
    this.options.items.push(item);
    this.render();
  }

  removeItem(index) {
    this.options.items.splice(index, 1);
    this.render();
  }

  updateItem(index, item) {
    if (index >= 0 && index < this.options.items.length) {
      this.options.items[index] = item;
      this.render();
    }
  }

  findItemIndex(predicate) {
    return this.options.items.findIndex(predicate);
  }

  getItem(index) {
    return this.options.items[index];
  }

  getAllItems() {
    return [...this.options.items];
  }

  trigger(eventName, data = {}) {
    const event = new CustomEvent(`grid:${eventName}`, {
      detail: { grid: this, ...data },
    });
    this.element.dispatchEvent(event);
  }
}

// Make Grid globally available
window.Grid = Grid;
