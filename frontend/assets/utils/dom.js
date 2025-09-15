/**
 * DOM Utility Functions
 * Common DOM manipulation and query helpers
 */

// Element creation helpers
function createElement(tag, attributes = {}, children = []) {
  const element = document.createElement(tag);

  // Set attributes
  Object.keys(attributes).forEach((key) => {
    if (key === "className") {
      element.className = attributes[key];
    } else if (key === "innerHTML") {
      element.innerHTML = attributes[key];
    } else if (key === "textContent") {
      element.textContent = attributes[key];
    } else if (key.startsWith("data-")) {
      element.setAttribute(key, attributes[key]);
    } else if (key === "style" && typeof attributes[key] === "object") {
      Object.assign(element.style, attributes[key]);
    } else {
      element.setAttribute(key, attributes[key]);
    }
  });

  // Append children
  children.forEach((child) => {
    if (typeof child === "string") {
      element.appendChild(document.createTextNode(child));
    } else if (child instanceof Node) {
      element.appendChild(child);
    }
  });

  return element;
}

// Query helpers
function $(selector, context = document) {
  return context.querySelector(selector);
}

function $$(selector, context = document) {
  return Array.from(context.querySelectorAll(selector));
}

// Class manipulation
function addClass(element, className) {
  if (element) {
    element.classList.add(className);
  }
}

function removeClass(element, className) {
  if (element) {
    element.classList.remove(className);
  }
}

function toggleClass(element, className, force) {
  if (element) {
    return element.classList.toggle(className, force);
  }
  return false;
}

function hasClass(element, className) {
  return element ? element.classList.contains(className) : false;
}

// Event helpers
function on(element, event, handler, options = {}) {
  if (element) {
    element.addEventListener(event, handler, options);
  }
}

function off(element, event, handler, options = {}) {
  if (element) {
    element.removeEventListener(event, handler, options);
  }
}

function delegate(container, selector, event, handler) {
  if (!container) return;

  const delegatedHandler = (e) => {
    const target = e.target.closest(selector);
    if (target && container.contains(target)) {
      handler.call(target, e);
    }
  };

  container.addEventListener(event, delegatedHandler);
  return () => container.removeEventListener(event, delegatedHandler);
}

// Style helpers
function show(element) {
  if (element) {
    element.style.display = "";
  }
}

function hide(element) {
  if (element) {
    element.style.display = "none";
  }
}

function toggle(element, force) {
  if (!element) return;

  const isHidden =
    element.style.display === "none" ||
    getComputedStyle(element).display === "none";

  if (force === undefined) {
    element.style.display = isHidden ? "" : "none";
  } else {
    element.style.display = force ? "" : "none";
  }
}

// Animation helpers
function fadeIn(element, duration = 300) {
  if (!element) return Promise.resolve();

  return new Promise((resolve) => {
    element.style.opacity = "0";
    element.style.display = "";

    const start = performance.now();

    function animate(currentTime) {
      const elapsed = currentTime - start;
      const progress = Math.min(elapsed / duration, 1);

      element.style.opacity = progress;

      if (progress < 1) {
        requestAnimationFrame(animate);
      } else {
        element.style.opacity = "";
        resolve();
      }
    }

    requestAnimationFrame(animate);
  });
}

function fadeOut(element, duration = 300) {
  if (!element) return Promise.resolve();

  return new Promise((resolve) => {
    const start = performance.now();
    const startOpacity = parseFloat(getComputedStyle(element).opacity) || 1;

    function animate(currentTime) {
      const elapsed = currentTime - start;
      const progress = Math.min(elapsed / duration, 1);

      element.style.opacity = startOpacity * (1 - progress);

      if (progress < 1) {
        requestAnimationFrame(animate);
      } else {
        element.style.display = "none";
        element.style.opacity = "";
        resolve();
      }
    }

    requestAnimationFrame(animate);
  });
}

// Form helpers
function getFormData(form) {
  if (!form) return {};

  const formData = new FormData(form);
  const data = {};

  for (const [key, value] of formData.entries()) {
    if (data[key]) {
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

function setFormData(form, data) {
  if (!form || !data) return;

  Object.keys(data).forEach((key) => {
    const field = form.querySelector(`[name="${key}"]`);
    if (field) {
      if (field.type === "checkbox" || field.type === "radio") {
        field.checked = Boolean(data[key]);
      } else {
        field.value = data[key] || "";
      }
    }
  });
}

// Scroll helpers
function scrollTo(element, options = {}) {
  if (!element) return;

  const defaultOptions = {
    behavior: "smooth",
    block: "start",
    inline: "nearest",
  };

  element.scrollIntoView({ ...defaultOptions, ...options });
}

function isInViewport(element) {
  if (!element) return false;

  const rect = element.getBoundingClientRect();
  return (
    rect.top >= 0 &&
    rect.left >= 0 &&
    rect.bottom <=
      (window.innerHeight || document.documentElement.clientHeight) &&
    rect.right <= (window.innerWidth || document.documentElement.clientWidth)
  );
}

// Debounce and throttle
function debounce(func, wait, immediate = false) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      timeout = null;
      if (!immediate) func.apply(this, args);
    };
    const callNow = immediate && !timeout;
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
    if (callNow) func.apply(this, args);
  };
}

function throttle(func, limit) {
  let inThrottle;
  return function executedFunction(...args) {
    if (!inThrottle) {
      func.apply(this, args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
}

// Text helpers
function escapeHtml(text) {
  const map = {
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    '"': "&quot;",
    "'": "&#039;",
  };
  return text.replace(/[&<>"']/g, (m) => map[m]);
}

function truncate(text, length, suffix = "...") {
  if (text.length <= length) return text;
  return text.substring(0, length) + suffix;
}

// Template helpers
function template(str, data) {
  return str.replace(/\{\{(\w+)\}\}/g, (match, key) => {
    return data.hasOwnProperty(key) ? data[key] : match;
  });
}

// Make available globally for vanilla JS
window.createElement = createElement;
window.$ = $;
window.$$ = $$;
window.addClass = addClass;
window.removeClass = removeClass;
window.toggleClass = toggleClass;
window.hasClass = hasClass;
window.on = on;
window.off = off;
window.delegate = delegate;
window.show = show;
window.hide = hide;
window.toggle = toggle;
window.fadeIn = fadeIn;
window.fadeOut = fadeOut;
window.getFormData = getFormData;
window.setFormData = setFormData;
window.scrollTo = scrollTo;
window.isInViewport = isInViewport;
window.debounce = debounce;
window.throttle = throttle;
window.escapeHtml = escapeHtml;
window.truncate = truncate;
window.template = template;
