<template>
  <div class="top-menu-bar">
    <div class="menu-options">
      <button class="menu-button" @click="$emit('toggle-settings-modal')">设置</button>
      <button class="menu-button" @click="$emit('toggle-about-modal')">关于</button>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';

const props = defineProps({
  menuItems: {
    type: Array,
    default: () => [
      { id: 'commands', name: '指令管理' },
      { id: 'collections', name: '集合管理' },
      { id: 'tags', name: '标签管理' }
    ]
  },
  activeMenu: {
    type: String,
    default: 'commands'
  }
});

const emit = defineEmits([
  'change-menu',
  'open-add-interface',
  'toggle-settings-modal',
  'toggle-about-modal'
]);
</script>

<style scoped>
/* 顶部选项菜单栏 */
.top-menu-bar {
  height: 30px;
  background-color: #f8f9fa;
  color: #333;
  display: flex;
  align-items: center;
  padding: 0 15px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  user-select: none;
}

/* 菜单选项容器 */
.menu-options {
  display: flex;
  gap: 15px;
}

/* 菜单按钮 */
.menu-button {
  background: none;
  border: none;
  color: #333;
  padding: 5px 10px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.3s ease;
}

/* 菜单按钮悬停效果 */
.menu-button:hover {
  background-color: rgba(0, 0, 0, 0.1);
}
</style>