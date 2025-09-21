/**
 * DataTable Component
 * Reusable table with sorting, filtering, and actions
 * Extends BaseComponent (loaded via script tag)
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
      pagination: true,
      pageSize: 10,
      emptyMessage: "No data available",
      className: "data-table",
      rowActions: [],
      bulkActions: [],
      selectable: false,
      exportable: false,
    };
  }

  init() {
    // Initialize filteredData BEFORE calling super.init() which triggers render()
    this.currentPage = 1;
    this.sortColumn = null;
    this.sortDirection = "asc";
    this.selectedRows = new Set();
    this.filteredData = Array.isArray(this.options.data)
      ? [...this.options.data]
      : [];

    // Now call super.init() which will call render()
    super.init();
  }

  render() {
    if (!this.element) return;

    this.element.className = this.options.className;
    this.element.innerHTML = this.getTemplate();
    this.renderTable();
  }

  getTemplate() {
    return `
      ${
        this.options.filterable ||
        this.options.selectable ||
        this.options.exportable
          ? this.getToolbarTemplate()
          : ""
      }
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              ${
                this.options.selectable
                  ? '<th class="select-column"><input type="checkbox" class="select-all"></th>'
                  : ""
              }
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

  getToolbarTemplate() {
    return `
      <div class="table-toolbar">
        <div class="toolbar-left">
          ${
            this.options.filterable
              ? `
            <div class="search-input-container">
              <i class="fas fa-search"></i>
              <input type="text" class="search-input" placeholder="Search...">
            </div>
          `
              : ""
          }
          ${
            this.options.selectable && this.options.bulkActions.length > 0
              ? `
            <div class="bulk-actions" style="display: none;">
              <span class="bulk-selected-count">0 selected</span>
              ${this.options.bulkActions
                .map(
                  (action) => `
                <button class="btn btn-sm ${
                  action.className || "btn-outline"
                }" data-bulk-action="${action.action}">
                  ${action.icon ? `<i class="${action.icon}"></i>` : ""} ${
                    action.text
                  }
                </button>
              `
                )
                .join("")}
            </div>
          `
              : ""
          }
        </div>
        <div class="toolbar-right">
          ${
            this.options.exportable
              ? `
            <button class="btn btn-outline btn-sm export-btn" data-action="export">
              <i class="fas fa-download"></i> Export
            </button>
          `
              : ""
          }
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

      // Handle export action
      if (action === "export") {
        this.exportData();
        return;
      }

      const rowId = target.closest("tr")?.getAttribute("data-id");
      const rowData = this.options.data.find((item) => item.id == rowId);

      if (action && rowData) {
        this.handleRowAction(action, rowData, target);
      }
    });

    // Bulk selection
    if (this.options.selectable) {
      // Select all checkbox
      this.addEventListener("change", ".select-all", (e, target) => {
        const isChecked = target.checked;
        const pageData = this.getCurrentPageData();

        pageData.forEach((item) => {
          if (isChecked) {
            this.selectedRows.add(item.id);
          } else {
            this.selectedRows.delete(item.id);
          }
        });

        this.updateBulkActions();
        this.renderTable(); // Re-render to update checkboxes
      });

      // Individual row checkboxes
      this.addEventListener("change", ".row-select", (e, target) => {
        const rowId = target.closest("tr").getAttribute("data-id");

        if (target.checked) {
          this.selectedRows.add(rowId);
        } else {
          this.selectedRows.delete(rowId);
        }

        this.updateBulkActions();
        this.updateSelectAllCheckbox();
      });

      // Bulk actions
      this.addEventListener("click", "[data-bulk-action]", (e, target) => {
        const action = target.getAttribute("data-bulk-action");
        const selectedData = Array.from(this.selectedRows)
          .map((id) => this.options.data.find((item) => item.id == id))
          .filter(Boolean);

        if (selectedData.length > 0) {
          this.handleBulkAction(action, selectedData);
        }
      });
    }
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

    // Safety check - ensure filteredData is an array before sorting
    if (!this.filteredData || !Array.isArray(this.filteredData)) {
      this.filteredData = Array.isArray(this.options.data)
        ? [...this.options.data]
        : [];
    }

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

    // Ensure we have valid data to filter
    const data = Array.isArray(this.options.data) ? this.options.data : [];

    if (!term) {
      this.filteredData = [...data];
    } else {
      this.filteredData = data.filter((item) => {
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

    // Safety check - ensure filteredData is always an array
    if (!this.filteredData || !Array.isArray(this.filteredData)) {
      this.filteredData = Array.isArray(this.options.data)
        ? [...this.options.data]
        : [];
    }

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
    const selectCell = this.options.selectable
      ? `<td class="select-cell">
          <input type="checkbox" class="row-select" ${
            this.selectedRows.has(item.id) ? "checked" : ""
          }>
         </td>`
      : "";

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

    return `<tr data-id="${item.id}">${selectCell}${cells}${actions}</tr>`;
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
      // Safety check for filteredData
      const safeFilteredData = this.filteredData || [];
      const start = (this.currentPage - 1) * this.options.pageSize + 1;
      const end = Math.min(
        this.currentPage * this.options.pageSize,
        safeFilteredData.length
      );
      info.textContent = `Showing ${start}-${end} of ${safeFilteredData.length}`;
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
    const safeFilteredData = this.filteredData || [];
    return Math.ceil(safeFilteredData.length / this.options.pageSize);
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
    const safeNewData = Array.isArray(newData) ? newData : [];
    this.options.data = safeNewData;
    this.filteredData = [...safeNewData];
    this.currentPage = 1;
    this.renderTable();
  }

  refresh() {
    this.renderTable();
  }

  // Bulk selection methods
  updateBulkActions() {
    const bulkActionsContainer = this.element.querySelector(".bulk-actions");
    const selectedCount = this.selectedRows.size;

    if (bulkActionsContainer) {
      if (selectedCount > 0) {
        bulkActionsContainer.style.display = "flex";
        const countElement = bulkActionsContainer.querySelector(
          ".bulk-selected-count"
        );
        if (countElement) {
          countElement.textContent = `${selectedCount} selected`;
        }
      } else {
        bulkActionsContainer.style.display = "none";
      }
    }
  }

  updateSelectAllCheckbox() {
    const selectAllCheckbox = this.element.querySelector(".select-all");
    if (selectAllCheckbox) {
      const pageData = this.getCurrentPageData();
      const allSelected =
        pageData.length > 0 &&
        pageData.every((item) => this.selectedRows.has(item.id));
      const someSelected = pageData.some((item) =>
        this.selectedRows.has(item.id)
      );

      selectAllCheckbox.checked = allSelected;
      selectAllCheckbox.indeterminate = someSelected && !allSelected;
    }
  }

  getCurrentPageData() {
    if (!this.options.pagination) {
      return this.filteredData;
    }

    const startIndex = (this.currentPage - 1) * this.options.pageSize;
    const endIndex = startIndex + this.options.pageSize;
    return this.filteredData.slice(startIndex, endIndex);
  }

  // Export functionality
  exportData() {
    const dataToExport =
      this.selectedRows.size > 0
        ? Array.from(this.selectedRows)
            .map((id) => this.options.data.find((item) => item.id == id))
            .filter(Boolean)
        : this.filteredData;

    const csvContent = this.convertToCSV(dataToExport);
    this.downloadCSV(csvContent, "table-export.csv");
  }

  convertToCSV(data) {
    if (data.length === 0) return "";

    const headers = this.options.columns.map((col) => col.title);
    const rows = data.map((item) =>
      this.options.columns.map((col) => {
        const getValue = col.getValue || ((item) => item[col.key]);
        const value = getValue(item);
        // Escape quotes and wrap in quotes if contains comma
        return typeof value === "string" &&
          (value.includes(",") || value.includes('"'))
          ? `"${value.replace(/"/g, '""')}"`
          : value;
      })
    );

    return [headers, ...rows].map((row) => row.join(",")).join("\n");
  }

  downloadCSV(csvContent, filename) {
    const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
    const link = document.createElement("a");

    if (link.download !== undefined) {
      const url = URL.createObjectURL(blob);
      link.setAttribute("href", url);
      link.setAttribute("download", filename);
      link.style.visibility = "hidden";
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  }

  // Event handlers
  handleBulkAction(action, selectedData) {
    this.trigger("bulkAction", { action, data: selectedData });
  }

  trigger(eventName, data = {}) {
    const event = new CustomEvent(`datatable:${eventName}`, {
      detail: { table: this, ...data },
    });
    this.element.dispatchEvent(event);
  }
}

// Make DataTable globally available
window.DataTable = DataTable;
