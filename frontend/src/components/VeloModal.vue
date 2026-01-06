<template>
  <transition name="modal">
    <div v-if="show" class="modal-mask" @click="closeOnMask">
      <div class="modal-wrapper">
        <div class="modal-container" @click.stop>
          <div class="modal-header">
            <slot name="header">
              <h3>{{ title }}</h3>
            </slot>
            <button class="modal-close" @click="$emit('close')">&times;</button>
          </div>

          <div class="modal-body">
            <slot name="body">
              {{ message }}
            </slot>
          </div>

          <div class="modal-footer">
            <slot name="footer">
              <button class="btn-secondary" @click="$emit('close')">{{ cancelText }}</button>
              <button 
                class="btn-confirm" 
                :class="confirmClass" 
                @click="$emit('confirm')"
              >
                {{ confirmText }}
              </button>
            </slot>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script>
export default {
  name: 'VeloModal',
  props: {
    show: Boolean,
    title: {
      type: String,
      default: 'Confirm Action'
    },
    message: String,
    confirmText: {
      type: String,
      default: 'Confirm'
    },
    cancelText: {
      type: String,
      default: 'Cancel'
    },
    type: {
      type: String,
      default: 'primary' // primary, danger
    }
  },
  computed: {
    confirmClass() {
      return this.type === 'danger' ? 'btn-danger' : 'btn-primary';
    }
  },
  methods: {
    closeOnMask() {
      this.$emit('close');
    }
  }
}
</script>

<style scoped>
.modal-mask {
  position: fixed;
  z-index: 9998;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.3s ease;
}

.modal-wrapper {
  width: 100%;
  max-width: 450px;
  padding: 20px;
}

.modal-container {
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  transition: all 0.3s ease;
  overflow: hidden;
}

.modal-header {
  padding: 20px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #f1f3f5;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.15rem;
  font-weight: 700;
  color: #1a1b1e;
}

.modal-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #adb5bd;
  cursor: pointer;
  line-height: 1;
}

.modal-body {
  padding: 24px;
  font-size: 0.95rem;
  color: #495057;
  line-height: 1.5;
}

.modal-footer {
  padding: 16px 24px;
  background-color: #f8f9fa;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

button {
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-secondary {
  background-color: #fff;
  border: 1px solid #dee2e6;
  color: #495057;
}

.btn-secondary:hover {
  background-color: #f1f3f5;
  border-color: #ced4da;
}

.btn-primary {
  background-color: #5c7cfa;
  color: white;
}

.btn-primary:hover {
  background-color: #4c6ef5;
}

.btn-danger {
  background-color: #fa5252;
  color: white;
}

.btn-danger:hover {
  background-color: #e03131;
}

/* Animation */
.modal-enter-from {
  opacity: 0;
}

.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.95);
  opacity: 0;
}
</style>
