<template>
  <div class="card">
    <div class="card-header" v-if="$slots.header || title">
      <slot name="header">
        <h4 class="card-title">{{title}}</h4>
        <p class="card-category" v-if="subTitle">{{subTitle}}</p>
      </slot>
    </div>
    <div class="card-body">
      <slot></slot>
    </div>
    <div class="card-footer" v-if="$slots.footer">
      <hr>
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: "VeloCard",
  props: {
    title: String,
    subTitle: String
  }
};
</script>

<style scoped>
.card {
  border-radius: 8px;
  background-color: #FFFFFF;
  color: #2c3e50;
  margin-bottom: 24px;
  position: relative;
  border: 1px solid #e9ecef;
  display: flex;
  flex-direction: column;
  word-wrap: break-word;
  transition: border-color 0.2s ease;
}

.card:hover {
  border-color: #dee2e6;
}

.card-header {
  background-color: transparent;
  padding: 20px 24px 0;
}

.card-title {
  margin: 0;
  color: #1a1b1e;
  font-weight: 600;
  font-size: 1.1rem;
}

.card-category {
  color: #909296;
  font-size: 13px;
  margin-top: 4px;
  margin-bottom: 12px;
}

.card-body {
  padding: 12px 24px 24px 24px;
  flex: 1 1 auto;
}

.card-footer {
  padding: 0 24px 16px 24px;
  background-color: transparent;
  line-height: 30px;
  border-top: none !important;
  font-size: 13px;
  color: #adb5bd;
}

hr {
  margin-top: 0;
  margin-bottom: 12px;
  border: 0;
  border-top: 1px solid #f1f3f5;
}
</style>