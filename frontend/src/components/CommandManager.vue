<template>
  <div class="app-container">
    <!-- 顶部菜单栏 -->
    <TopMenuBar 
      :menuItems="menuItems"
      :activeMenu="activeMenu"
      @toggle-active-menu="toggleActiveMenu"
      @toggle-add-interface="toggleAddInterface"
      @toggle-settings-modal="toggleSettingsModal"
      @toggle-about-modal="toggleAboutModal"
    />
    
    <div class="content-container">
      <!-- 左侧边栏 -->
      <Sidebar 
        :menuItems="menuItems"
        :activeMenu="activeMenu"
        :menuType="menuType"
        :systemType="systemType"
        :searchKeyword="searchKeyword"
        :isSortDropdownOpen="isSortDropdownOpen"
        :sortOptions="sortOptions"
        :sortDirections="sortDirections"
        :tags="tags"
        :collections="collections"
        @toggle-active-menu="toggleActiveMenu"
        @toggle-system-type="toggleSystemType"
        @update:menuType="menuType = $event"
        @update:searchKeyword="searchKeyword = $event"
        @update:systemType="systemType = $event"
        @toggle-sort-dropdown="toggleSortDropdown"
        @toggle-sort-option="toggleSortOption"
        @toggle-sort-direction="toggleSortDirection"
      />
      
      <!-- 右侧主内容区 -->
      <MainContent 
        :activeAddInterface="activeAddInterface"
        :menuType="menuType"
        :activeMenu="activeMenu"
        :tags="tags"
        :collections="collections"
        :commands="commands"
        :searchKeyword="searchKeyword"
        @toggle-add-interface="toggleAddInterface"
        @update:searchKeyword="searchKeyword = $event"
        @add-command="addCommand"
        @add-collection="addCollection"
        @add-tag="addTag"
        @edit-item="editItem"
        @delete-item="deleteItem"
        @copy-to-clipboard="copyToClipboard"
      />
    </div>
    
    <!-- 对话框组件 -->
    <Dialogs 
      :isSettingsModalOpen="isSettingsModalOpen"
      :isAboutModalOpen="isAboutModalOpen"
      :settings="settings"
      @update:settings="settings = $event"
      @close-settings-modal="closeSettingsModal"
      @close-about-modal="closeAboutModal"
    />
  </div>
</template>

<script setup>
// 导入Vue 3的响应式API和生命周期钩子
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
// 导入后端API函数
import { GetMenuItems, GetOptions, CreateTag, CreateCollection } from '../../wailsjs/go/main/App';
// 导入子组件
import TopMenuBar from './layout/TopMenuBar.vue';
import Sidebar from './layout/Sidebar.vue';
import MainContent from './layout/MainContent.vue';
import Dialogs from './Dialogs.vue';

// 响应式数据
// 激活的添加界面类型（'command' | 'collection' | 'tag' | ''）
const activeAddInterface = ref('');
// 选中的项目
const selectedItem = ref(null);
// 当前激活的菜单ID
const activeMenu = ref('home');
// 当前菜单类型（'tags' | 'collections' | 'all'）
const menuType = ref('tags');
// 系统类型筛选数组
const systemType = ref(['linux']);
// 搜索关键词
const searchKeyword = ref('');
// 排序下拉框是否打开
const isSortDropdownOpen = ref(false);
// 排序选项配置
const sortOptions = ref({
  time: false,      // 按时间排序
  name: false,      // 按名称排序
  copyCount: false, // 按复制次数排序
  id: false,        // 按ID排序
  sortValue: false  // 按排序值排序
});
// 排序方向配置
const sortDirections = ref({
  time: 'asc',      // 时间排序方向
  name: 'asc',      // 名称排序方向
  copyCount: 'asc', // 复制次数排序方向
  id: 'asc',        // ID排序方向
  sortValue: 'asc'  // 排序值排序方向
});
// 设置模态框是否打开
const isSettingsModalOpen = ref(false);
// 关于模态框是否打开
const isAboutModalOpen = ref(false);
// 设置数据
const settings = ref({
  apiEndpoint: '',   // API端点
  apiKey: '',        // API密钥
  theme: 'light',    // 主题
  language: 'zh-CN'  // 语言
});

// 表单数据
// 新命令表单
const newCommand = ref({
  id: '',           // 命令ID
  name: '',         // 命令名称
  content: '',      // 命令内容
  description: '',  // 命令描述
  tags: [],         // 关联的标签ID数组
  collections: [],  // 关联的集合ID数组
  sortValue: 0,     // 排序值
  copyCount: 0,     // 复制次数
  systemType: ['windows'] // 适用系统类型
});

// 新集合表单
const newCollection = ref({
  id: '',           // 集合ID
  name: '',         // 集合名称
  description: '',  // 集合描述
  sortValue: 0      // 排序值
});

// 新标签表单
const newTag = ref({
  id: '',           // 标签ID
  name: '',         // 标签名称
  description: '',  // 标签描述
  sortValue: 0      // 排序值
});

// 模拟数据 - 标签列表
const tags = ref([
  { id: '1', name: '开发', description: '开发相关指令', sortValue: 1 },
  { id: '2', name: '运维', description: '运维相关指令', sortValue: 2 },
  { id: '3', name: '测试', description: '测试相关指令', sortValue: 3 },
  { id: '4', name: '数据库', description: '数据库相关指令', sortValue: 4 },
  { id: '5', name: '网络', description: '网络相关指令', sortValue: 5 }
]);

// 模拟数据 - 集合列表
const collections = ref([
  { id: '1', name: '常用命令', description: '常用的命令集合', sortValue: 1 },
  { id: '2', name: 'Git命令', description: 'Git版本控制相关命令', sortValue: 2 },
  { id: '3', name: 'Docker命令', description: 'Docker容器相关命令', sortValue: 3 },
  { id: '4', name: 'Linux命令', description: 'Linux系统相关命令', sortValue: 4 },
  { id: '5', name: 'Windows命令', description: 'Windows系统相关命令', sortValue: 5 }
]);

// 模拟数据 - 命令列表（当前注释掉，使用后端API获取）
const commands = ref([]);

// 模拟菜单项数据
const menuItems = ref({
  topMenu: [  // 顶部菜单
    { id: 'home', name: '首页', icon: '🏠' },
    { id: 'commands', name: '命令管理', icon: '⚡' },
    { id: 'collections', name: '集合管理', icon: '📁' },
    { id: 'tags', name: '标签管理', icon: '🏷️' }
  ],
  tags: [     // 标签菜单
    { id: 'all-tags', name: '全部标签', icon: '🏷️' },
    { id: 'dev', name: '开发', icon: '💻' },
    { id: 'ops', name: '运维', icon: '🔧' },
    { id: 'test', name: '测试', icon: '🧪' },
    { id: 'db', name: '数据库', icon: '🗃️' },
    { id: 'network', name: '网络', icon: '🌐' }
  ],
  collections: [  // 集合菜单
    { id: 'all-collections', name: '全部集合', icon: '📁' },
    { id: 'common', name: '常用命令', icon: '⭐' },
    { id: 'git', name: 'Git命令', icon: '🔖' },
    { id: 'docker', name: 'Docker命令', icon: '🐳' },
    { id: 'linux', name: 'Linux命令', icon: '🐧' },
    { id: 'windows', name: 'Windows命令', icon: '🪟' }
  ],
  all: [     // 全部命令菜单
    { id: 'all-commands', name: '全部命令', icon: '⚡' },
    { id: 'recent', name: '最近使用', icon: '🕒' },
    { id: 'frequent', name: '高频使用', icon: '🔥' },
    { id: 'newest', name: '最新添加', icon: '🆕' }
  ]
});

// 计算属性：根据系统类型过滤命令
const filteredCommands = computed(() => {
  return commands.value.filter(command => {
    return command.systemType.some(type => systemType.value.includes(type));
  });
});

// 构建排序参数的辅助函数
function buildSortParams() {
  const sort = {};
  
  // 只传递勾选的排序选项
  if (sortOptions.value.name) {
    sort.name = sortDirections.value.name; // asc 或 desc
  } else {
    sort.name = null;
  }
  
  if (sortOptions.value.time) {
    sort.create_time = sortDirections.value.time; // asc 或 desc
  } else {
    sort.create_time = null;
  }
  
  if (sortOptions.value.copyCount) {
    sort.copy_counts = sortDirections.value.copyCount; // asc 或 desc
  } else {
    sort.copy_counts = null;
  }
  
  return sort;
}

// 切换活动菜单
function toggleActiveMenu(menuId) {
  // 更新当前激活的菜单ID
  activeMenu.value = menuId;
  
  // 构建Option参数
  const option = {
    Name: searchKeyword.value,
    Os: systemType.value,
    Type: menuType.value,
    ID: parseInt(menuId),
    Sort: buildSortParams()
  };
  
  // 调用GetOptions获取数据
  GetOptions(option).then((result) => {
    console.log("获取数据成功:", result);
    
    // 更新数据
    if (result.tags) {
      tags.value = result.tags;
    }
    if (result.collections) {
      collections.value = result.collections;
    }
    if (result.options) {
      commands.value = result.options;
    }
  }).catch((error) => {
    console.error("获取数据失败:", error);
  });
}

// 切换系统类型
function toggleSystemType(type) {
  // 查找系统类型在数组中的索引
  const index = systemType.value.indexOf(type);
  if (index === -1) {
    // 如果不存在，添加到数组中
    systemType.value.push(type);
  } else {
    // 如果存在，从数组中移除
    systemType.value.splice(index, 1);
  }
  
  // 调用GetOptions获取数据
  const option = {
    Name: searchKeyword.value,
    Os: systemType.value,
    Type: menuType.value,
    ID: parseInt(activeMenu.value),
    Sort: buildSortParams()
  };
  
  GetOptions(option).then((result) => {
    console.log("获取数据成功:", result);
    
    // 更新数据
    if (result.tags) {
      tags.value = result.tags;
    }
    if (result.collections) {
      collections.value = result.collections;
    }
    if (result.options) {
      commands.value = result.options;
    }
  }).catch((error) => {
    console.error("获取数据失败:", error);
  });
}

// 切换排序下拉框
function toggleSortDropdown() {
  isSortDropdownOpen.value = !isSortDropdownOpen.value;
}

// 切换排序选项
function toggleSortOption(option) {
  sortOptions.value[option] = !sortOptions.value[option];
}

// 切换排序方向
function toggleSortDirection(option) {
  sortDirections.value[option] = sortDirections.value[option] === 'asc' ? 'desc' : 'asc';
}

// 切换新增界面
function toggleAddInterface(type) {
  // 如果当前界面是目标类型，则关闭；否则打开目标类型界面
  activeAddInterface.value = activeAddInterface.value === type ? '' : type;
}

// 切换设置模态框
function toggleSettingsModal() {
  isSettingsModalOpen.value = !isSettingsModalOpen.value;
}

// 切换关于模态框
function toggleAboutModal() {
  isAboutModalOpen.value = !isAboutModalOpen.value;
}

// 关闭设置模态框
function closeSettingsModal() {
  isSettingsModalOpen.value = false;
}

// 关闭关于模态框
function closeAboutModal() {
  isAboutModalOpen.value = false;
}

// 新增命令
function addCommand(command) {
  // 将新命令添加到命令列表
  commands.value.push(command);
  // 这里可以添加保存到后端的逻辑
}

// 新增集合
async function addCollection(collection) {
  try {
    // 调用后端接口创建集合
    await CreateCollection(collection);
    // 将新集合添加到集合列表
    collections.value.push(collection);
  } catch (error) {
    console.error('创建集合失败:', error);
    alert('创建集合失败: ' + error.message);
  }
}

// 新增标签
async function addTag(tag) {
  try {
    // 调用后端接口创建标签
    await CreateTag(tag);
    // 将新标签添加到标签列表
    tags.value.push(tag);
  } catch (error) {
    console.error('创建标签失败:', error);
    alert('创建标签失败: ' + error.message);
  }
}

// 编辑项目
function editItem(item) {
  // 设置选中的项目
  selectedItem.value = item;
  // 这里可以添加编辑逻辑
}

// 删除项目
function deleteItem(item) {
  // 弹出确认对话框
  if (confirm(`确定要删除 ${item.name} 吗？`)) {
    if (item.tags) { // 是命令
      // 查找命令在数组中的索引
      const index = commands.value.findIndex(cmd => cmd.id === item.id);
      if (index !== -1) {
        // 从数组中移除
        commands.value.splice(index, 1);
      }
    } else if (item.collections) { // 是集合
      // 查找集合在数组中的索引
      const index = collections.value.findIndex(col => col.id === item.id);
      if (index !== -1) {
        // 从数组中移除
        collections.value.splice(index, 1);
      }
    } else { // 是标签
      // 查找标签在数组中的索引
      const index = tags.value.findIndex(tag => tag.id === item.id);
      if (index !== -1) {
        // 从数组中移除
        tags.value.splice(index, 1);
      }
    }
  }
}

// 复制到剪贴板
function copyToClipboard(content) {
  // 使用浏览器的剪贴板API复制内容
  navigator.clipboard.writeText(content).then(() => {
    // 显示复制成功提示
    showCopySuccess();
  }).catch(err => {
    console.error('复制失败:', err);
  });
}

// 显示复制成功提示
function showCopySuccess() {
  // 这里可以添加复制成功的动画或提示
  console.log('复制成功！');
}

// 生成UUID
function generateUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

// 点击空白处关闭排序下拉框
function handleClickOutside(event) {
  // 如果点击的元素不是排序容器的子元素，则关闭下拉框
  if (!event.target.closest('.sort-container')) {
    isSortDropdownOpen.value = false;
  }
}

// 组件挂载时
onMounted(() => {
  // 绑定点击事件
  document.addEventListener('click', handleClickOutside);
  
  // 获取菜单项
  GetMenuItems().then((result) => {
    console.log("获取菜单项成功:", result);
    // 这里可以根据后端返回的数据更新menuItems
  }).catch((error) => {
    console.error("获取菜单项失败:", error);
  });
  
  // 构建初始Option参数
  const option = {
    Name: searchKeyword.value,
    Os: systemType.value,
    Type: menuType.value,
    ID: parseInt(activeMenu.value),
    Sort: buildSortParams()
  };
  
  // 调用GetOptions获取初始数据
  GetOptions(option).then((result) => {
    console.log("获取初始数据成功:", result);
    
    // 更新数据
    if (result.tags) {
      tags.value = result.tags;
    }
    if (result.collections) {
      collections.value = result.collections;
    }
    if (result.options) {
      commands.value = result.options;
    }
  }).catch((error) => {
    console.error("获取初始数据失败:", error);
  });
});

// 组件卸载时
onUnmounted(() => {
  // 解绑点击事件
  document.removeEventListener('click', handleClickOutside);
});

// 监听菜单类型变化
watch(() => menuType.value, () => {
  // 构建Option参数
  const option = {
    Name: searchKeyword.value,
    Os: systemType.value,
    Type: menuType.value,
    ID: parseInt(activeMenu.value),
    Sort: buildSortParams()
  };
  
  // 调用GetOptions获取数据
  GetOptions(option).then((result) => {
    console.log("获取数据成功:", result);
    
    // 更新数据
    if (result.tags) {
      tags.value = result.tags;
    }
    if (result.collections) {
      collections.value = result.collections;
    }
    if (result.options) {
      commands.value = result.options;
    }
  }).catch((error) => {
    console.error("获取数据失败:", error);
  });
});

// 监听搜索关键词变化
watch(() => searchKeyword.value, () => {
  // 构建Option参数
  const option = {
    Name: searchKeyword.value,
    Os: systemType.value,
    Type: menuType.value,
    ID: parseInt(activeMenu.value),
    Sort: buildSortParams()
  };
  
  // 调用GetOptions获取数据
  GetOptions(option).then((result) => {
    console.log("获取数据成功:", result);
    
    // 更新数据
    if (result.tags) {
      tags.value = result.tags;
    }
    if (result.collections) {
      collections.value = result.collections;
    }
    if (result.options) {
      commands.value = result.options;
    }
  }).catch((error) => {
    console.error("获取数据失败:", error);
  });
});

// 监听排序选项变化
watch([sortOptions, sortDirections], () => {
  // 构建Option参数
  const option = {
    Name: searchKeyword.value,
    Os: systemType.value,
    Type: menuType.value,
    ID: parseInt(activeMenu.value),
    Sort: buildSortParams()
  };
  
  // 调用GetOptions获取数据
  GetOptions(option).then((result) => {
    console.log("获取数据成功:", result);
    
    // 更新数据
    if (result.tags) {
      tags.value = result.tags;
    }
    if (result.collections) {
      collections.value = result.collections;
    }
    if (result.options) {
      commands.value = result.options;
    }
  }).catch((error) => {
    console.error("获取数据失败:", error);
  });
}, { deep: true });
</script>

<style scoped>
/* 应用容器 */
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  background-color: #f5f5f5;
  color: #333;
  font-family: 'Arial', sans-serif;
}

/* 内容容器 */
.content-container {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .content-container {
    flex-direction: column;
  }
}
</style>