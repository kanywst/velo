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
  border-radius: 6px;
  box-shadow: 0 2px 2px rgba(204, 197, 185, 0.5);
  background-color: #FFFFFF;
  color: #252422;
  margin-bottom: 20px;
  position: relative;
  border: 1px solid rgba(0,0,0,.125);
  display: flex;
  flex-direction: column;
  word-wrap: break-word;
}

.card-header {
  background-color: transparent;
  border-bottom: 1px solid rgba(0,0,0,.125);
  padding: 15px 15px 0;
}

.card-title {
  margin: 0;
  color: #252422;
  font-weight: 300;
  font-size: 1.5rem;
}

.card-category {
  color: #9A9A9A;
  font-size: 14px;
  margin-top: 5px;
  margin-bottom: 10px;
}

.card-body {
  padding: 15px 15px 10px 15px;
  flex: 1 1 auto;
}

.card-footer {
  padding: 0 15px 10px 15px;
  background-color: transparent;
  line-height: 30px;
  border-top: none !important;
  font-size: 14px;
}

hr {
  margin-top: 5px;
  margin-bottom: 5px;
  border: 0;
  border-top: 1px solid rgba(0,0,0,.1);
}
</style>