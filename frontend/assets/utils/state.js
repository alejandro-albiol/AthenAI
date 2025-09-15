/**
 * State Management Utility
 * Simple state management for application data
 */
class StateManager {
  constructor(initialState = {}) {
    this.state = { ...initialState };
    this.subscribers = new Map();
    this.middleware = [];
  }

  getState() {
    return { ...this.state };
  }

  setState(updates, notify = true) {
    const prevState = { ...this.state };
    this.state = { ...this.state, ...updates };

    if (notify) {
      this.notifySubscribers(prevState, this.state, updates);
    }
  }

  subscribe(key, callback) {
    if (!this.subscribers.has(key)) {
      this.subscribers.set(key, []);
    }
    this.subscribers.get(key).push(callback);

    // Return unsubscribe function
    return () => {
      const callbacks = this.subscribers.get(key);
      if (callbacks) {
        const index = callbacks.indexOf(callback);
        if (index > -1) {
          callbacks.splice(index, 1);
        }
      }
    };
  }

  notifySubscribers(prevState, currentState, updates) {
    Object.keys(updates).forEach((key) => {
      const callbacks = this.subscribers.get(key);
      if (callbacks) {
        callbacks.forEach((callback) => {
          callback(currentState[key], prevState[key], currentState);
        });
      }
    });

    // Notify global subscribers
    const globalCallbacks = this.subscribers.get("*");
    if (globalCallbacks) {
      globalCallbacks.forEach((callback) => {
        callback(currentState, prevState, updates);
      });
    }
  }

  addMiddleware(middlewareFn) {
    this.middleware.push(middlewareFn);
  }

  dispatch(action) {
    let result = action;

    // Apply middleware
    this.middleware.forEach((middleware) => {
      result = middleware(result, this.getState(), this);
    });

    return result;
  }
}

// Create global state instance
const appState = new StateManager({
  user: null,
  currentView: "overview",
  loading: false,
  error: null,
  notifications: [],
  gyms: [],
  equipment: [],
  exercises: [],
  muscularGroups: [],
});

// Make available globally for vanilla JS
window.StateManager = StateManager;
window.appState = appState;
