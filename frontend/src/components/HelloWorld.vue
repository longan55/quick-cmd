<script setup>
import {reactive, ref, onMounted, onUnmounted} from 'vue'
import {Greet} from '../../wailsjs/go/main/App'

const data = reactive({
  name: "",
  resultText: "Please enter your name below 👇",
})

const activeMenu = ref(['home']) // 改为数组，支持多选
const menuType = ref('tags') // tags, collections, all
const systemType = ref(['windows']) // 改为数组，支持多选
const searchKeyword = ref('') // 搜索关键词

// 弹窗控制响应式数据
const isSettingsOpen = ref(false) // 控制设置页面的显示/隐藏
const isAboutOpen = ref(false) // 控制关于对话框的显示/隐藏

// 排序相关响应式数据
const sortOptions = reactive({
  time: false,
  name: false,
  copyCount: false,
  id: false,
  sortValue: false
})

const sortDirections = reactive({
  time: 'asc',
  name: 'asc',
  copyCount: 'asc',
  id: 'asc',
  sortValue: 'asc'
})

const isSortDropdownOpen = ref(false) // 控制排序下拉框的显示/隐藏

function greet() {
  Greet(data.name).then(result => {
    data.resultText = result
  })
}

// 系统类型选择方法
function toggleSystemType(type) {
  const index = systemType.value.indexOf(type)
  if (index === -1) {
    systemType.value.push(type) // 添加到数组
  } else {
    systemType.value.splice(index, 1) // 从数组中移除
  }
}

// 菜单选择方法
function toggleActiveMenu(menu) {
  const index = activeMenu.value.indexOf(menu)
  if (index === -1) {
    activeMenu.value.push(menu) // 添加到数组
  } else {
    activeMenu.value.splice(index, 1) // 从数组中移除
  }
}

// 排序相关方法
function toggleSortDropdown() {
  isSortDropdownOpen.value = !isSortDropdownOpen.value
}

function toggleSortOption(option) {
  sortOptions[option] = !sortOptions[option]
}

function toggleSortDirection(option) {
  sortDirections[option] = sortDirections[option] === 'asc' ? 'desc' : 'asc'
}

// 点击外部区域关闭下拉框
function handleClickOutside(event) {
  const dropdown = document.querySelector('.sort-dropdown')
  const trigger = document.querySelector('.sort-trigger')
  if (dropdown && trigger && !dropdown.contains(event.target) && !trigger.contains(event.target)) {
    isSortDropdownOpen.value = false
  }
}

// 添加点击外部区域关闭下拉框的事件监听
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

// 打开设置页面
function openSettings() {
  isSettingsOpen.value = true
}

// 关闭设置页面
function closeSettings() {
  isSettingsOpen.value = false
}

// 打开关于对话框
function openAbout() {
  isAboutOpen.value = true
}

// 关闭关于对话框
function closeAbout() {
  isAboutOpen.value = false
}

// 移除事件监听
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

</script>

<template>
  <div class="app-container">
    <!-- 顶部选项菜单栏 -->
    <div class="top-menu-bar">
      <div class="menu-options">
        <button class="menu-button" @click="openSettings">设置</button>
        <button class="menu-button" @click="openAbout">关于</button>
      </div>
    </div>
    
    <!-- 内容容器：将左侧菜单和右侧内容水平排列 -->
    <div class="content-container">
    <!-- 左侧菜单栏 -->
    <aside class="sidebar">
      <div class="sidebar-header">
        <!-- <h2>菜单</h2> -->
        <!-- 搜索框 -->
        <div class="search-container">
          <input 
            type="text" 
            v-model="searchKeyword" 
            class="search-input" 
            placeholder="搜索..." 
            autocomplete="off"
          />
          <!-- 清除按钮，只有输入内容时显示 -->
          <button 
            class="clear-button" 
            v-if="searchKeyword" 
            @click="searchKeyword = ''"
          >
            ✕
          </button>
          <button class="search-button">
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
          <button class="switcher-button" @click="menuType = 'tags'" :class="{ active: menuType === 'tags' }">
            标签
          </button>
          <button class="switcher-button" @click="menuType = 'collections'" :class="{ active: menuType === 'collections' }">
            集合
          </button>
          <button class="switcher-button" @click="menuType = 'all'" :class="{ active: menuType === 'all' }">
            指令
          </button>
        </div>

        <!-- 排序选项控件 -->
        <div class="sort-container">
          <button class="sort-trigger switcher-button" @click="toggleSortDropdown">
            排序选项
            <span class="sort-icon">▼</span>
          </button>
          
          <!-- 排序选项下拉框 -->
          <div class="sort-dropdown" v-if="isSortDropdownOpen">
            <div class="sort-option" @click="toggleSortOption('time')">
              <input type="checkbox" :checked="sortOptions.time" readonly />
              <span class="option-text">时间</span>
              <button class="sort-direction-button" @click.stop="toggleSortDirection('time')">
                {{ sortDirections.time === 'asc' ? '↑' : '↓' }}
              </button>
            </div>
            <div class="sort-option" @click="toggleSortOption('name')">
              <input type="checkbox" :checked="sortOptions.name" readonly />
              <span class="option-text">名称</span>
              <button class="sort-direction-button" @click.stop="toggleSortDirection('name')">
                {{ sortDirections.name === 'asc' ? '↑' : '↓' }}
              </button>
            </div>
            <div class="sort-option" @click="toggleSortOption('copyCount')">
              <input type="checkbox" :checked="sortOptions.copyCount" readonly />
              <span class="option-text">复制次数</span>
              <button class="sort-direction-button" @click.stop="toggleSortDirection('copyCount')">
                {{ sortDirections.copyCount === 'asc' ? '↑' : '↓' }}
              </button>
            </div>
            <div class="sort-option" @click="toggleSortOption('id')">
              <input type="checkbox" :checked="sortOptions.id" readonly />
              <span class="option-text">ID</span>
              <button class="sort-direction-button" @click.stop="toggleSortDirection('id')">
                {{ sortDirections.id === 'asc' ? '↑' : '↓' }}
              </button>
            </div>
            <div class="sort-option" @click="toggleSortOption('sortValue')">
              <input type="checkbox" :checked="sortOptions.sortValue" readonly />
              <span class="option-text">排序值</span>
              <button class="sort-direction-button" @click.stop="toggleSortDirection('sortValue')">
                {{ sortDirections.sortValue === 'asc' ? '↑' : '↓' }}
              </button>
            </div>
          </div>
        </div>
      </div>
      <nav class="sidebar-nav">
        <!-- 标签菜单 -->
        <div v-if="menuType === 'tags'" class="menu-buttons">
          <button class="menu-button" @click="toggleActiveMenu('home')" :class="{ active: activeMenu.includes('home') }">
            <span class="menu-icon">🏠</span>
            <span class="menu-text">首页</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('tag1')" :class="{ active: activeMenu.includes('tag1') }">
            <span class="menu-icon">🏷️</span>
            <span class="menu-text">工作</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('tag2')" :class="{ active: activeMenu.includes('tag2') }">
            <span class="menu-icon">🏷️</span>
            <span class="menu-text">学习</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('tag3')" :class="{ active: activeMenu.includes('tag3') }">
            <span class="menu-icon">🏷️</span>
            <span class="menu-text">生活</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('tag4')" :class="{ active: activeMenu.includes('tag4') }">
            <span class="menu-icon">🏷️</span>
            <span class="menu-text">娱乐</span>
          </button>
        </div>

        <!-- 集合菜单 -->
        <div v-else-if="menuType === 'collections'" class="menu-buttons">
          <button class="menu-button" @click="toggleActiveMenu('home')" :class="{ active: activeMenu.includes('home') }">
            <span class="menu-icon">🏠</span>
            <span class="menu-text">首页</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('collection1')" :class="{ active: activeMenu.includes('collection1') }">
            <span class="menu-icon">📁</span>
            <span class="menu-text">常用工具</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('collection2')" :class="{ active: activeMenu.includes('collection2') }">
            <span class="menu-icon">📁</span>
            <span class="menu-text">开发资源</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('collection3')" :class="{ active: activeMenu.includes('collection3') }">
            <span class="menu-icon">📁</span>
            <span class="menu-text">文档资料</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('collection4')" :class="{ active: activeMenu.includes('collection4') }">
            <span class="menu-icon">📁</span>
            <span class="menu-text">项目管理</span>
          </button>
        </div>

        <!-- 全部菜单 -->
        <div v-else-if="menuType === 'all'" class="menu-buttons">
          <button class="menu-button" @click="toggleActiveMenu('home')" :class="{ active: activeMenu.includes('home') }">
            <span class="menu-icon">🏠</span>
            <span class="menu-text">首页</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('tag1')" :class="{ active: activeMenu.includes('tag1') }">
            <span class="menu-icon">🏷️</span>
            <span class="menu-text">工作</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('tag2')" :class="{ active: activeMenu.includes('tag2') }">
            <span class="menu-icon">🏷️</span>
            <span class="menu-text">学习</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('tag3')" :class="{ active: activeMenu.includes('tag3') }">
            <span class="menu-icon">🏷️</span>
            <span class="menu-text">生活</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('tag4')" :class="{ active: activeMenu.includes('tag4') }">
            <span class="menu-icon">🏷️</span>
            <span class="menu-text">娱乐</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('collection1')" :class="{ active: activeMenu.includes('collection1') }">
            <span class="menu-icon">📁</span>
            <span class="menu-text">常用工具</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('collection2')" :class="{ active: activeMenu.includes('collection2') }">
            <span class="menu-icon">📁</span>
            <span class="menu-text">开发资源</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('collection3')" :class="{ active: activeMenu.includes('collection3') }">
            <span class="menu-icon">📁</span>
            <span class="menu-text">文档资料</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('collection4')" :class="{ active: activeMenu.includes('collection4') }">
            <span class="menu-icon">📁</span>
            <span class="menu-text">项目管理</span>
          </button>
          <button class="menu-button" @click="toggleActiveMenu('settings')" :class="{ active: activeMenu.includes('settings') }">
            <span class="menu-icon">⚙️</span>
            <span class="menu-text">设置</span>
          </button>
        </div>
      </nav>
    </aside>

    <!-- 右侧主内容区 -->
    <main class="main-content">
      <!-- 顶部导航栏 (20%高度) -->
      <header class="top-nav">
        <!-- <div class="top-nav-content">
          <div class="logo">
            <h1>Quick Cmd</h1>
          </div>
          <div class="user-info">
            <span>欢迎使用</span>
          </div>
        </div> -->
        <div class="top-nav-content">
          <div id="result" class="result">{{ data.resultText }}</div>
          <div id="input" class="input-box">
            <input id="name" v-model="data.name" autocomplete="off" class="input" type="text"/>
            <button class="btn" @click="greet">Greet</button>
          </div>
        </div>
      </header>

      <!-- 内容区域 -->
      <section class="content-area">
        
        
        <!-- 嵌套组件区域 -->
        <div class="nested-component-area">
          <slot></slot>
        </div>
      </section>
    </main>
    </div> <!-- 关闭内容容器 -->
    
    <!-- 设置页面 -->
    <div class="settings-overlay" v-if="isSettingsOpen">
      <div class="settings-container">
        <div class="settings-header">
          <h2>设置</h2>
          <button class="close-button" @click="closeSettings">×</button>
        </div>
        <div class="settings-content">
          <!-- 设置内容区域，目前留空 -->
        </div>
      </div>
    </div>
    
    <!-- 关于对话框 -->
    <div class="about-overlay" v-if="isAboutOpen">
      <div class="about-dialog">
        <div class="about-header">
          <h2>关于</h2>
          <button class="close-button" @click="closeAbout">×</button>
        </div>
        <div class="about-content">
          <p>Quick Cmd</p>
          <p>版本 1.0.0</p>
          <p>一个快速命令管理工具</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 全局布局样式 */
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100%;
  overflow: hidden;
}

/* 内容容器样式 */
.content-container {
  display: flex;
  flex-grow: 1;
  overflow: hidden;
}

/* 顶部选项菜单栏 */
.top-menu-bar {
  height: 30px;
  background-color: #2c3e50;
  color: white;
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
.top-menu-bar .menu-button {
  background: none;
  border: none;
  color: white;
  padding: 5px 10px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.3s ease;
}

/* 菜单按钮悬停效果 */
.top-menu-bar .menu-button:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

/* 左侧菜单栏 */
.sidebar {
  width: 300px;
  background-color: #ffffff;
  color: #333333;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
  border-right: 1px solid #e0e0e0;
  flex-grow: 1;
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
  max-width: calc(100% - 40px); /* 最大宽度为父容器宽度减去padding */
  box-sizing: border-box;
  overflow: hidden;
  margin-left: auto;
  margin-right: auto;
}

/* 搜索输入框 */
.search-input {
  width: 150px;
  padding: 8px 15px 8px 16px; /* 增加右侧padding，为清除按钮留出空间 */
  border: 1px solid #069b6e;
  border-radius: 20px 0 0 20px;
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
  color: white;
  box-shadow: 0 4px 8px rgba(52, 152, 219, 0.3);
}

/* 菜单按钮点击效果 */
.menu-button:active {
  transform: translateX(5px) translateY(2px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

/* 右侧主内容区 */
.main-content {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 顶部导航栏 (20%高度) */
.top-nav {
  height: 20%;
  background-color: #e8f4f8;
  color: #333;
  display: flex;
  align-items: center;
  padding: 0 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.top-nav-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo h1 {
  margin: 0;
  font-size: 1.8rem;
}

.user-info {
  font-size: 1.1rem;
}

/* 内容区域 */
.content-area {
  flex: 1;
  background-color: #ecf0f1;
  padding: 20px;
  overflow-y: auto;
}

/* 示例内容样式 */
.example-content {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
  text-align: center;
}

.input-box {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
}

.input-box .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  padding: 0 8px;
  cursor: pointer;
  background-color: #3498db;
  color: white;
  transition: background-color 0.3s;
}

.input-box .btn:hover {
  background-color: #2980b9;
}

.input-box .input {
  border: 1px solid #ddd;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(255, 255, 255, 1);
  -webkit-font-smoothing: antialiased;
  width: 200px;
}

.input-box .input:focus {
  border-color: #3498db;
}

/* 嵌套组件区域 */
.nested-component-area {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  min-height: 200px;
}

/* 设置页面遮罩层 */
.settings-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 2000;
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 设置页面容器 */
.settings-container {
  width: 100%;
  height: 100%;
  background-color: white;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 设置页面头部 */
.settings-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background-color: #3498db;
  color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.settings-header h2 {
  margin: 0;
  font-size: 1.5rem;
}

/* 关闭按钮样式 */
.close-button {
  background: none;
  border: none;
  color: white;
  font-size: 1.8rem;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background-color 0.3s ease;
}

.close-button:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

/* 设置页面内容区域 */
.settings-content {
  flex-grow: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #f9f9f9;
}

/* 关于对话框遮罩层 */
.about-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 2000;
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 关于对话框容器 */
.about-dialog {
  background-color: white;
  border-radius: 10px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
  width: 400px;
  max-width: 90%;
  overflow: hidden;
}

/* 关于对话框头部 */
.about-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background-color: #2c3e50;
  color: white;
}

.about-header h2 {
  margin: 0;
  font-size: 1.3rem;
}

/* 关于对话框内容区域 */
.about-content {
  padding: 20px;
  text-align: center;
  line-height: 1.6;
}

.about-content p {
  margin: 8px 0;
}
</style>
