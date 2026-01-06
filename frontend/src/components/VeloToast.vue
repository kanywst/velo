<template>
  <transition-group name="toast" tag="div" class="toast-container">
    <div 
      v-for="toast in toasts" 
      :key="toast.id" 
      class="toast" 
      :class="toast.type"
    >
      <div class="toast-icon">{{ getIcon(toast.type) }}</div>
      <div class="toast-content">{{ toast.message }}</div>
      <button class="toast-close" @click="remove(toast.id)">&times;</button>
    </div>
  </transition-group>
</template>

<script>
export default {
  name: 'VeloToast',
  data() {
    return {
      toasts: [],
      counter: 0
    }
  },
  methods: {
    add(message, type = 'success', duration = 4000) {
      const id = this.counter++;
      this.toasts.push({ id, message, type });
      if (duration > 0) {
        setTimeout(() => this.remove(id), duration);
      }
    },
    remove(id) {
      const index = this.toasts.findIndex(t => t.id === id);
      if (index !== -1) {
        this.toasts.splice(index, 1);
      }
    },
    getIcon(type) {
      switch (type) {
        case 'success': return '✅';
        case 'error': return '❌';
        case 'info': return 'ℹ️';
        default: return '🔔';
      }
    }
  }
}
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 24px;
  right: 24px;
  z-index: 10000;
  display: flex;
  flex-direction: column;
  gap: 12px;
  pointer-events: none;
}

.toast {
  min-width: 300px;
  max-width: 450px;
  background-color: #fff;
  border-radius: 10px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  pointer-events: auto;
  border-left: 4px solid #dee2e6;
  transition: all 0.3s cubic-bezier(0.68, -0.55, 0.265, 1.55);
}

.toast.success { border-left-color: #20c997; }
.toast.error { border-left-color: #fa5252; }
.toast.info { border-left-color: #5c7cfa; }

.toast-icon {
  font-size: 1.2rem;
}

.toast-content {
  flex: 1;
  font-size: 0.9rem;
  font-weight: 500;
  color: #1a1b1e;
}

.toast-close {
  background: none;
  border: none;
  color: #adb5bd;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 0;
  line-height: 1;
}

.toast-close:hover {
  color: #495057;
}

/* Animation */
.toast-enter-from {
  transform: translateX(100%);
  opacity: 0;
}
.toast-leave-to {
  transform: translateX(100%);
  opacity: 0;
}
</style>
