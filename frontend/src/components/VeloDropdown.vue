<template>
  <div class="dropdown" :class="{ show: isOpen }" ref="dropdownRef">
    <button 
      class="dropdown-toggle" 
      @click="toggleDropDown" 
      type="button"
      :class="{ 'has-value': modelValue && modelValue !== 'all' }"
    >
      <span class="title">{{ currentLabel }}</span>
      <i class="arrow-down"></i>
    </button>
    <div class="dropdown-menu" :class="{ show: isOpen }">
      <a 
        v-for="option in options" 
        :key="option.value" 
        href="#" 
        class="dropdown-item"
        :class="{ active: modelValue === option.value }"
        @click.prevent="selectOption(option)"
      >
        {{ option.label }}
      </a>
    </div>
  </div>
</template>

<script>
export default {
  name: 'VeloDropdown',
  props: {
    placeholder: {
      type: String,
      default: 'Select...'
    },
    options: {
      type: Array,
      default: () => []
    },
    modelValue: {
      type: [String, Number],
      default: ''
    }
  },
  data() {
    return {
      isOpen: false
    }
  },
  computed: {
    currentLabel() {
      const selected = this.options.find(o => o.value === this.modelValue);
      return selected ? selected.label : this.placeholder;
    }
  },
  mounted() {
    document.addEventListener('click', this.handleClickOutside)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.handleClickOutside)
  },
  methods: {
    toggleDropDown() {
      this.isOpen = !this.isOpen
    },
    closeDropDown() {
      this.isOpen = false
    },
    selectOption(option) {
      this.$emit('update:modelValue', option.value)
      this.$emit('change', option.value)
      this.closeDropDown()
    },
    handleClickOutside(event) {
      if (this.$refs.dropdownRef && !this.$refs.dropdownRef.contains(event.target)) {
        this.closeDropDown()
      }
    }
  }
}
</script>

<style scoped>
.dropdown {
  position: relative;
  display: inline-block;
  min-width: 160px;
}

.dropdown-toggle {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 8px 16px;
  font-size: 0.95rem;
  font-weight: 500;
  color: #555;
  background-color: #fff;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  outline: none;
}

.dropdown-toggle:hover {
  background-color: #f8f9fa;
  border-color: #ccc;
}

.dropdown-toggle.has-value {
  color: #333;
  border-color: #bbb;
}

.arrow-down {
  border: solid #888;
  border-width: 0 2px 2px 0;
  display: inline-block;
  padding: 3px;
  transform: rotate(45deg);
  margin-left: 10px;
  margin-bottom: 2px;
  transition: transform 0.2s ease;
}

.dropdown.show .arrow-down {
  transform: rotate(-135deg);
  margin-bottom: -2px;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  z-index: 1000;
  display: none;
  min-width: 100%;
  padding: 5px 0;
  margin-top: 4px;
  font-size: 0.95rem;
  color: #212529;
  text-align: left;
  list-style: none;
  background-color: #fff;
  background-clip: padding-box;
  border: 1px solid rgba(0, 0, 0, 0.15);
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  max-height: 300px;
  overflow-y: auto;
  opacity: 0;
  transform: translateY(-10px);
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.dropdown-menu.show {
  display: block;
  opacity: 1;
  transform: translateY(0);
}

.dropdown-item {
  display: block;
  width: 100%;
  padding: 8px 16px;
  clear: both;
  font-weight: 400;
  color: #212529;
  text-align: inherit;
  white-space: nowrap;
  background-color: transparent;
  border: 0;
  text-decoration: none;
  transition: background-color 0.15s ease;
}

.dropdown-item:hover, .dropdown-item:focus {
  color: #16181b;
  text-decoration: none;
  background-color: #f8f9fa;
}

.dropdown-item.active {
  color: #fff;
  text-decoration: none;
  background-color: #41B883; /* Vue Green */
}
</style>
