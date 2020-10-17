<template>
  <div :class="className">
    <div class="edit-area" v-if="isEditMode">
      <form @submit="toggleEdit">
        <textarea class="float-left form-control" v-if="multi" v-model="text" required></textarea>
        <input type="text" class="float-left" v-if="!multi" v-model="text" required>

        <button class="float-right icon-btn" type="submit">
          <v-icon>save</v-icon>
        </button>
      </form>
    </div>
    <div class="row">
      <p class="editable text-area" v-if="!isEditMode" @click="toggleEdit">{{text}}</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'TextEditable',
  props: ['value', 'is-edit', 'class-name', 'on-change', 'multi'],
  data () {
    return {
      text: this.value,
      isEditMode: this.isEdit
    }
  },
  methods: {
    toggleEdit (e) {
      e.preventDefault()
      this.isEditMode = !this.isEditMode
      if (!this.isEditMode && typeof this.onChange === 'function') {
        this.onChange(this.text)
      }
    }
  }
}
</script>

<style>
.float-left {
  float: left;
  width: 90%;
  border-style: solid;
}

.float-right {
  float: right;
  margin-top: 4px;
}

span.float-right {
  cursor: pointer;
}

.icon-btn {
  border: none;
  background: none;
}
.row {
  display: flex;
  flex-direction: row;
}

.column {
  display: flex;
  flex-direction: column;
  flex: 10;
}
.editable {
}
.editable::after {
  content: "";
  display: block;
  font-family: "Material Icons";
  content: "edit";
  float: right;
  padding-left: 5px;
  margin: 0 6px 0 0;
}
</style>
