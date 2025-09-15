import BaseComponent from "./BaseComponent.js";

/**
 * DataTable Component
 * Reusable table with sorting, filtering, and actions
 */
class DataTable extends BaseComponent {
  constructor(selector, options = {}) {
    super(selector, options);
  }

  getDefaultOptions() {
    return {
      ...super.getDefaultOptions(),
      data: [],
      columns: [],
      sortable: true,
      filterable: true,
      pagination: false,
      pageSize: 10,
      emptyMessage: "No data available",
      className: "data-table",
      rowActions: [],
    };
  }

  init() {
    super.init();
    this.currentPage = 1;
    this.sortColumn = null;
    this.sortDirection = "asc";
    this.filteredData = [...this.options.data];
  }

  render() {
    if (!this.element) return;

    this.element.className = this.options.className;
    this.element.innerHTML = this.getTemplate();
    this.renderTable();
  }

  getTemplate() {
    return `
      ${this.options.filterable ? this.getFilterTemplate() : ""}
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              ${this.options.columns
                .map((col) => this.getHeaderCell(col))
                .join("")}
              ${
                this.options.rowActions.length > 0
                  ? '<th class="actions-column">Actions</th>'
                  : ""
              }
            </tr>
          </thead>
          <tbody class="table-body">
            <!-- Table rows will be rendered here -->
          </tbody>
        </table>
      </div>
      ${this.options.pagination ? this.getPaginationTemplate() : ""}
      <div class="table-empty" style="display: none;">
        <div class="empty-state">
          <i class="fas fa-inbox"></i>
          <h3>No Data</h3>
          <p>${this.options.emptyMessage}</p>
        </div>
      </div>
    `;
  }

  getFilterTemplate() {
    return `
      <div class="table-filters">
        <div class="search-input-container">
          <i class="fas fa-search"></i>
          <input type="text" class="search-input" placeholder="Search...">
        </div>
      </div>
    `;
  }

  getPaginationTemplate() {
    return `
      <div class="table-pagination">
        <div class="pagination-info">
          <span class="showing-text"></span>
        </div>
        <div class="pagination-controls">
          <button class="btn btn-sm pagination-btn" data-action="prev">
            <i class="fas fa-chevron-left"></i>
          </button>
          <span class="page-numbers"></span>
          <button class="btn btn-sm pagination-btn" data-action="next">
            <i class="fas fa-chevron-right"></i>
          </button>
        </div>
      </div>
    `;
  }

  getHeaderCell(column) {
    const sortable = this.options.sortable && column.sortable !== false;
    const sortClass =
      this.sortColumn === column.key
        ? `sorted sorted-${this.sortDirection}`
        : "";

    return `
      <th class="table-header ${sortable ? "sortable" : ""} ${sortClass}" 
          ${sortable ? `data-sort="${column.key}"` : ""}>
        <span>${column.title}</span>
        ${sortable ? '<i class="sort-icon fas fa-sort"></i>' : ""}
      </th>
    `;
  }

  setupEventListeners() {
    // Sorting
    if (this.options.sortable) {
      this.addEventListener("click", ".sortable", (e, target) => {
        const column = target.getAttribute("data-sort");
        this.sort(column);
      });
    }

    // Filtering
    if (this.options.filterable) {
      this.addEventListener("input", ".search-input", (e, target) => {
        this.filter(target.value);
      });
    }

    // Pagination
    if (this.options.pagination) {
      this.addEventListener("click", ".pagination-btn", (e, target) => {
        const action = target.getAttribute("data-action");
        if (action === "prev" && this.currentPage > 1) {
          this.currentPage--;
          this.renderTable();
        } else if (
          action === "next" &&
          this.currentPage < this.getTotalPages()
        ) {
          this.currentPage++;
          this.renderTable();
        }
      });
    }

    // Row actions
    this.addEventListener("click", "[data-action]", (e, target) => {
      const action = target.getAttribute("data-action");
      const rowId = target.closest("tr").getAttribute("data-id");
      const rowData = this.options.data.find((item) => item.id == rowId);

      if (action && rowData) {
        this.handleRowAction(action, rowData, target);
      }
    });
  }

  sort(column) {
    if (this.sortColumn === column) {
      this.sortDirection = this.sortDirection === "asc" ? "desc" : "asc";
    } else {
      this.sortColumn = column;
      this.sortDirection = "asc";
    }

    const columnConfig = this.options.columns.find((col) => col.key === column);
    const getValue = columnConfig?.getValue || ((item) => item[column]);

    this.filteredData.sort((a, b) => {
      const aVal = getValue(a);
      const bVal = getValue(b);

      if (aVal < bVal) return this.sortDirection === "asc" ? -1 : 1;
      if (aVal > bVal) return this.sortDirection === "asc" ? 1 : -1;
      return 0;
    });

    this.currentPage = 1;
    this.render();
  }

  filter(searchTerm) {
    const term = searchTerm.toLowerCase().trim();

    if (!term) {
      this.filteredData = [...this.options.data];
    } else {
      this.filteredData = this.options.data.filter((item) => {
        return this.options.columns.some((column) => {
          const getValue = column.getValue || ((item) => item[column.key]);
          const value = getValue(item);
          return String(value).toLowerCase().includes(term);
        });
      });
    }

    this.currentPage = 1;
    this.renderTable();
  }

  renderTable() {
    const tbody = this.element.querySelector(".table-body");
    const emptyState = this.element.querySelector(".table-empty");

    if (this.filteredData.length === 0) {
      tbody.innerHTML = "";
      emptyState.style.display = "block";
      return;
    }

    emptyState.style.display = "none";

    const startIndex = this.options.pagination
      ? (this.currentPage - 1) * this.options.pageSize
      : 0;
    const endIndex = this.options.pagination
      ? startIndex + this.options.pageSize
      : this.filteredData.length;
    const pageData = this.filteredData.slice(startIndex, endIndex);

    tbody.innerHTML = pageData.map((item) => this.renderRow(item)).join("");

    if (this.options.pagination) {
      this.updatePagination();
    }

    this.updateSortHeaders();
  }

  renderRow(item) {
    const cells = this.options.columns
      .map((column) => {
        const getValue = column.getValue || ((item) => item[column.key]);
        const render = column.render || ((value) => value);
        const value = getValue(item);

        return `<td class="table-cell">${render(value, item)}</td>`;
      })
      .join("");

    const actions =
      this.options.rowActions.length > 0
        ? `<td class="actions-cell">${this.renderRowActions(item)}</td>`
        : "";

    return `<tr data-id="${item.id}">${cells}${actions}</tr>`;
  }

  renderRowActions(item) {
    return this.options.rowActions
      .map((action) => {
        const className = action.className || "btn btn-sm btn-outline";
        const disabled =
          action.disabled && action.disabled(item) ? "disabled" : "";

        return `
        <button class="${className}" data-action="${
          action.action
        }" ${disabled} title="${action.title || action.text}">
          ${action.icon ? `<i class="${action.icon}"></i>` : action.text}
        </button>
      `;
      })
      .join("");
  }

  updateSortHeaders() {
    // Remove existing sort classes
    this.element.querySelectorAll(".table-header").forEach((header) => {
      header.classList.remove("sorted", "sorted-asc", "sorted-desc");
    });

    // Add current sort class
    if (this.sortColumn) {
      const header = this.element.querySelector(
        `[data-sort="${this.sortColumn}"]`
      );
      if (header) {
        header.classList.add("sorted", `sorted-${this.sortDirection}`);
      }
    }
  }

  updatePagination() {
    const totalPages = this.getTotalPages();
    const info = this.element.querySelector(".showing-text");
    const pageNumbers = this.element.querySelector(".page-numbers");
    const prevBtn = this.element.querySelector('[data-action="prev"]');
    const nextBtn = this.element.querySelector('[data-action="next"]');

    if (info) {
      const start = (this.currentPage - 1) * this.options.pageSize + 1;
      const end = Math.min(
        this.currentPage * this.options.pageSize,
        this.filteredData.length
      );
      info.textContent = `Showing ${start}-${end} of ${this.filteredData.length}`;
    }

    if (pageNumbers) {
      pageNumbers.textContent = `Page ${this.currentPage} of ${totalPages}`;
    }

    if (prevBtn) {
      prevBtn.disabled = this.currentPage === 1;
    }

    if (nextBtn) {
      nextBtn.disabled = this.currentPage === totalPages;
    }
  }

  getTotalPages() {
    return Math.ceil(this.filteredData.length / this.options.pageSize);
  }

  handleRowAction(action, rowData, target) {
    const actionConfig = this.options.rowActions.find(
      (a) => a.action === action
    );
    if (actionConfig && actionConfig.handler) {
      actionConfig.handler.call(this, rowData, target);
    }

    // Emit event
    this.trigger("rowAction", { action, rowData, target });
  }

  updateData(newData) {
    this.options.data = newData;
    this.filteredData = [...newData];
    this.currentPage = 1;
    this.renderTable();
  }

  refresh() {
    this.renderTable();
  }

  trigger(eventName, data = {}) {
    const event = new CustomEvent(`datatable:${eventName}`, {
      detail: { table: this, ...data },
    });
    this.element.dispatchEvent(event);
  }
}

export default DataTable;
