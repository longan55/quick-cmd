<template>
  <aside class="sidebar">
    <div class="sidebar-header">
      <!-- 搜索框 -->
      <div class="search-container">
        <input 
          type="text" 
          v-model="searchKeywordLocal" 
          class="search-input" 
          placeholder="搜索..." 
          autocomplete="off"
          @keyup.enter="performSearch"
        />
        <!-- 清除按钮，只有输入内容时显示 -->
        <button 
          class="clear-button" 
          v-if="searchKeywordLocal" 
          @click="clearSearch"
        >
          ✕
        </button>
        <button class="search-button" @click="performSearch">
          🔍
        </button>
      </div>
      
      <!-- 系统类型选择控件 -->
      <div class="system-type-switcher">
        <button class="switcher-button" @click="toggleSystemType('windows')" :class="{ active: systemType.includes('windows') }">
          Windows
        </button>
        <button class="switcher-button" @click="toggleSystemType('linux')" :class="{ active: systemType.includes('linux') }">
          Linux
        </button>
        <button class="switcher-button" @click="toggleSystemType('mac')" :class="{ active: systemType.includes('mac') }">
          Mac
        </button>
      </div>
      
      <!-- 菜单类型切换控件 -->
      <div class="menu-type-switcher">
        <button class="switcher-button" @click="changeMenuType('tags')" :class="{ active: menuType === 'tags' }">
          标签
        </button>
        <button class="switcher-button" @click="changeMenuType('collections')" :class="{ active: menuType === 'collections' }">
          集合
        </button>
        <button class="switcher-button" @click="changeMenuType('all')" :class="{ active: menuType === 'all' }">
          指令
        </button>
      </div>

      <!-- 排序选项控件 -->
      <div class="sort-container">
        <button class="sort-trigger switcher-button" @click="emit('toggle-sort-dropdown')">
          排序选项
          <span class="sort-icon">▼</span>
        </button>
        
        <!-- 排序选项下拉框 -->
        <div class="sort-dropdown" v-if="isSortDropdownOpen">
          <div class="sort-option" @click="emit('toggle-sort-option', 'time')">
            <input type="checkbox" :checked="sortOptions.time" readonly />
            <span class="option-text">时间</span>
            <button class="sort-direction-button" @click.stop="emit('toggle-sort-direction', 'time')">
              {{ sortDirections.time === 'asc' ? '🔼' : '🔽'}}
            </button>
          </div>
          <div class="sort-option" @click="emit('toggle-sort-option', 'name')">
            <input type="checkbox" :checked="sortOptions.name" readonly />
            <span class="option-text">名称</span>
            <button class="sort-direction-button" @click.stop="emit('toggle-sort-direction', 'name')">
              {{ sortDirections.name === 'asc' ? '🔼' : '🔽' }}
            </button>
          </div>
          <div class="sort-option" @click="emit('toggle-sort-option', 'copyCount')">
            <input type="checkbox" :checked="sortOptions.copyCount" readonly />
            <span class="option-text">复制次数</span>
            <button class="sort-direction-button" @click.stop="emit('toggle-sort-direction', 'copyCount')">
              {{ sortDirections.copyCount === 'asc' ? '🔼' : '🔽' }}
            </button>
          </div>
          <div class="sort-option" @click="emit('toggle-sort-option', 'id')">
            <input type="checkbox" :checked="sortOptions.id" readonly />
            <span class="option-text">ID</span>
            <button class="sort-direction-button" @click.stop="emit('toggle-sort-direction', 'id')">
              {{ sortDirections.id === 'asc' ? '🔼' : '🔽' }}
            </button>
          </div>
          <div class="sort-option" @click="emit('toggle-sort-option', 'sortValue')">
            <input type="checkbox" :checked="sortOptions.sortValue" readonly />
            <span class="option-text">排序值</span>
            <button class="sort-direction-button" @click.stop="emit('toggle-sort-direction', 'sortValue')">
              {{ sortDirections.sortValue === 'asc' ? '🔼' : '🔽' }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <nav class="sidebar-nav">
      <!-- 标签菜单 -->
      <div v-if="menuType === 'tags'" class="menu-buttons">
        <button 
          v-for="tag in tags || []" 
          :key="tag.id" 
          class="menu-button" 
          @click="emit('toggle-active-menu', tag.id)" 
          :class="{ active: activeMenu === String(tag.id) }"
        >
          <!-- <span class="menu-icon">{{ tag.icon }}</span> -->
          <span class="menu-text">{{ tag.name }}</span>
        </button>
      </div>

      <!-- 集合菜单 -->
      <div v-else-if="menuType === 'collections'" class="menu-buttons">
        <button 
          v-for="collection in collections || []" 
          :key="collection.id" 
          class="menu-button" 
          @click="emit('toggle-active-menu', collection.id)" 
          :class="{ active: activeMenu === String(collection.id) }"
        >
          <span class="menu-text">{{ collection.name }}</span>
        </button>
      </div>

      <!-- 全部菜单 -->
      <div v-else-if="menuType === 'all'" class="menu-buttons">
        <button 
          v-for="menu in menuItems.all || []" 
          :key="menu.id" 
          class="menu-button" 
          @click="emit('toggle-active-menu', menu.id)" 
          :class="{ active: activeMenu === String(menu.id) }"
        >
          <span class="menu-icon">{{ menu.icon }}</span>
          <span class="menu-text">{{ menu.name }}</span>
        </button>
      </div>
    </nav>
  </aside>
</template>

<script setup>
// 导入Vue 3的响应式API和生命周期钩子
import { ref, defineProps, defineEmits, watch } from 'vue';
// 导入后端API函数

// 定义组件属性
const props = defineProps({
  tags:{
    type: Array,
    default: () => []
  },
  collections:{
    type: Array,
    default: () => []
  },
  cmdType:{
    type: Array,
    default: () => []
  },
  menuItems: {
    type: Object,
    default: () => ({})
  },
  activeMenu: {
    type: String,
    default: 'home'
  },
  menuType: {
    type: String,
    default: 'tags'
  },
  systemType: {
    type: Array,
    default: () => ['']
  },
  searchKeyword: {
    type: String,
    default: ''
  },
  isSortDropdownOpen: {
    type: Boolean,
    default: false
  },
  sortOptions: {
    type: Object,
    default: () => ({
      time: false,
      name: false,
      copyCount: false,
      id: false,
      sortValue: false
    })
  },
  sortDirections: {
    type: Object,
    default: () => ({
      time: 'asc',
      name: 'asc',
      copyCount: 'asc',
      id: 'asc',
      sortValue: 'asc'
    })
  }
});

// 定义组件事件
const emit = defineEmits([
  'toggle-active-menu', 
  'toggle-system-type',
  'update:menuType',
  'update:searchKeyword',
  'toggle-sort-dropdown',
  'toggle-sort-option',
  'toggle-sort-direction',
  'update:systemType'
]);

// 本地响应式数据
// 本地搜索关键词，用于双向绑定
const searchKeywordLocal = ref(props.searchKeyword);

// 监听外部传入的searchKeyword变化
watch(() => props.searchKeyword, (newValue) => {
  searchKeywordLocal.value = newValue;
});

// 清除搜索
function clearSearch() {
  // 清空本地搜索关键词
  searchKeywordLocal.value = '';
  // 触发更新事件
  emit('update:searchKeyword', '');
  // 执行搜索
  performSearch();
}

// 执行搜索
function performSearch() {
  // 触发搜索关键词更新事件
  emit('update:searchKeyword', searchKeywordLocal.value);
}

// 切换系统类型
function toggleSystemType(type) {
  // 触发系统类型切换事件
  emit('toggle-system-type', type);
}

// 改变菜单类型
function changeMenuType(type) {
  // 触发菜单类型更新事件
  emit('update:menuType', type);
  // 无论菜单类型是否变化，都重置为第一个选项（ID=0）
  emit('toggle-active-menu', 0);
}
</script>

<style scoped>
/* 左侧菜单栏 */
.sidebar {
  width: 250px;
  background-color: #ffffff;
  color: #333333;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
  border-right: 1px solid #e0e0e0;
  flex-shrink: 0; /* 改为flex-shrink: 0，防止侧边栏被压缩 */
}

/* 左侧菜单栏标题分割线 */
.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
  background-color: #fafafa;
  box-sizing: border-box;
}

.sidebar-header h2 {
  margin: 0 0 15px 0;
  font-size: 1.5rem;
  text-align: center;
}

/* 系统类型切换器 */
.system-type-switcher {
  display: flex;
  gap: 8px;
  justify-content: center;
  margin-bottom: 15px;
}

/* 搜索框容器 */
.search-container {
  margin-bottom: 15px;
  display: flex;
  position: relative;
  width: auto; /* 改为auto，让容器宽度由内容决定 */
  max-width: calc(100% - 15px); /* 最大宽度为父容器宽度减去padding */
  box-sizing: border-box;
  overflow: hidden;
  margin-left: auto;
  margin-right: auto;
}

/* 搜索输入框 */
.search-input {
  width: 200px;
  padding: 8px 15px 8px 8px; /* 增加右侧padding，为清除按钮留出空间 */
  border: 1px solid #c3ff00;
  border-radius: 18px 0 0 18px;
  font-size: 0.9rem;
  outline: none;
  transition: all 0.3s ease;
  background-color: #fff;
  color: #333;
  box-sizing: border-box;
  min-width: 0;
}

/* 搜索输入框聚焦效果 */
.search-input:focus {
  border-color: #3498db;
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.1);
}

/* 搜索输入框悬停效果 */
.search-input:hover {
  border-color: #bbb;
}

/* 清除按钮 */
.clear-button {
  padding: 8px 10px;
  border: 1px solid #ddd;
  border-left: none;
  border-right: none;
  background-color: #fff;
  color: #999;
  cursor: pointer;
  transition: all 0.3s ease;
  outline: none;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
}

/* 清除按钮悬停效果 */
.clear-button:hover {
  background-color: #f5f5f5;
  color: #333;
}

/* 清除按钮点击效果 */
.clear-button:active {
  background-color: #e0e0e0;
}

/* 搜索按钮 */
.search-button {
  padding: 8px 16px;
  border: 1px solid #ddd;
  border-left: none;
  border-radius: 0 20px 20px 0;
  background-color: #f0f0f0;
  color: #333;
  cursor: pointer;
  transition: all 0.3s ease;
  outline: none;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 搜索按钮悬停效果 */
.search-button:hover {
  background-color: #e0e0e0;
  color: #3498db;
}

/* 搜索按钮点击效果 */
.search-button:active {
  background-color: #d0d0d0;
  transform: translateY(1px);
}

/* 菜单类型切换器 */
.menu-type-switcher {
  display: flex;
  gap: 8px;
  justify-content: center;
}

/* 切换按钮样式 */
.switcher-button {
  padding: 6px 12px;
  background-color: #f0f0f0;
  color: #333;
  border: none;
  border-radius: 16px;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.3s ease;
  outline: none;
}

/* 切换按钮悬停效果 */
.switcher-button:hover {
  background-color: #e0e0e0;
  transform: translateY(-2px);
}

/* 切换按钮激活状态 */
.switcher-button.active {
  background-color: #3498db;
  color: white;
  box-shadow: 0 2px 8px rgba(52, 152, 219, 0.4);
  transform: translateY(-2px);
}

/* 切换按钮点击效果 */
.switcher-button:active {
  transform: translateY(0);
}

/* 排序选项容器 */
.sort-container {
  margin-top: 10px;
  position: relative;
  width: 100%;
  display: flex;
  justify-content: center;
}

/* 排序触发器按钮 */
.sort-trigger {
  display: flex;
  align-items: center;
  gap: 5px;
}

.sort-icon {
  font-size: 0.7rem;
  transition: transform 0.3s ease;
}

/* 排序下拉框 */
.sort-dropdown {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-top: 5px;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  width: 200px;
  z-index: 100;
  overflow: hidden;
}

/* 排序选项 */
.sort-option {
  display: flex;
  align-items: center;
  padding: 10px 15px;
  cursor: pointer;
  transition: background-color 0.3s ease;
  user-select: none;
}

.sort-option:hover {
  background-color: #f5f5f5;
}

.sort-option input[type="checkbox"] {
  margin-right: 10px;
  cursor: pointer;
  accent-color: #3498db;
}

.option-text {
  flex: 1;
}

/* 排序方向按钮 */
.sort-direction-button {
  background-color: transparent;
  border: none;
  cursor: pointer;
  font-size: 0.8rem;
  color: #666;
  padding: 2px 5px;
  border-radius: 3px;
  transition: all 0.3s ease;
}

.sort-direction-button:hover {
  background-color: #e0e0e0;
  color: #333;
}

.sidebar-nav {
  flex: 1;
  padding: 20px;
}

/* 菜单按钮容器 */
.menu-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 菜单按钮样式 */
.menu-button {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 14px 20px;
  background-color: transparent;
  color: #333;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: left;
  outline: none;
}

/* 菜单图标 */
.menu-icon {
  font-size: 1.2rem;
  width: 24px;
  text-align: center;
}

/* 菜单按钮悬停效果 */
.menu-button:hover {
  background-color: rgba(52, 152, 219, 0.1);
  transform: translateX(5px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

/* 菜单按钮激活状态 */
.menu-button.active {
  background-color: #3498db;
  color: black;
  box-shadow: 0 4px 8px rgba(52, 152, 219, 0.3);
}

/* 菜单按钮点击效果 */
.menu-button:active {
  transform: translateX(5px) translateY(2px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}
</style>