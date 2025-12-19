<template>
  <div class="app-container">
    <!-- 左侧上菜单栏 -->
    <TopMenuBar 
      :menuItems="menuItems"
      :activeMenu="activeMenu"
      @toggle-active-menu="toggleActiveMenu"
      @toggle-add-interface="toggleAddInterface"
      @toggle-settings-modal="toggleSettingsModal"
      @toggle-about-modal="toggleAboutModal"
    />
    
    <div class="content-container">
      <!-- 左侧下边栏 -->
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
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { GetMenuItems,GetCommands,GetTags,GetCollections,GetCommandsByTagID,GetCommandsByCollectionID } from '../../wailsjs/go/main/App';
import TopMenuBar from './layout/TopMenuBar.vue';
import Sidebar from './layout/Sidebar.vue';
import MainContent from './layout/MainContent.vue';
import Dialogs from './Dialogs.vue';

// 响应式数据
const activeAddInterface = ref('');
const selectedItem = ref(null);
const activeMenu = ref('home');
const menuType = ref('tags');
const systemType = ref([]);
const searchKeyword = ref('');
const isSortDropdownOpen = ref(false);
const sortOptions = ref({
  time: false,
  name: false,
  copyCount: false,
  id: false,
  sortValue: false
});
const sortDirections = ref({
  time: 'asc',
  name: 'asc',
  copyCount: 'asc',
  id: 'asc',
  sortValue: 'asc'
});
const isSettingsModalOpen = ref(false);
const isAboutModalOpen = ref(false);
const settings = ref({
  apiEndpoint: '',
  apiKey: '',
  theme: 'light',
  language: 'zh-CN'
});

// 表单数据
const newCommand = ref({
  id: '',
  name: '',
  content: '',
  description: '',
  tags: [],
  collections: [],
  sortValue: 0,
  copyCount: 0,
  systemType: ['windows']
});

const newCollection = ref({
  id: '',
  name: '',
  description: '',
  sortValue: 0
});

const newTag = ref({
  id: '',
  name: '',
  description: '',
  sortValue: 0
});

// 设置数据
// const settings = ref({
//   theme: 'light',
//   autoUpdate: true,
//   language: 'zh-CN'
// });

// 模拟数据
const tags = ref([
  { id: '1', name: '开发', description: '开发相关指令', sortValue: 1 },
  { id: '2', name: '运维', description: '运维相关指令', sortValue: 2 },
  { id: '3', name: '测试', description: '测试相关指令', sortValue: 3 },
  { id: '4', name: '数据库', description: '数据库相关指令', sortValue: 4 },
  { id: '5', name: '网络', description: '网络相关指令', sortValue: 5 }
]);

const collections = ref([
  { id: '1', name: '常用命令', description: '常用的命令集合', sortValue: 1 },
  { id: '2', name: 'Git命令', description: 'Git版本控制相关命令', sortValue: 2 },
  { id: '3', name: 'Docker命令', description: 'Docker容器相关命令', sortValue: 3 },
  { id: '4', name: 'Linux命令', description: 'Linux系统相关命令', sortValue: 4 },
  { id: '5', name: 'Windows命令', description: 'Windows系统相关命令', sortValue: 5 }
]);

const commands = ref([
  // { 
  //   id: '1', 
  //   name: '查看当前目录', 
  //   content: 'ls -la', 
  //   description: '查看当前目录下的所有文件和文件夹，包括隐藏文件', 
  //   tags: ['1', '4'], 
  //   collections: ['1', '4'], 
  //   sortValue: 1, 
  //   copyCount: 10, 
  //   systemType: ['linux', 'mac'] 
  // },
  // { 
  //   id: '2', 
  //   name: '查看当前目录', 
  //   content: 'dir', 
  //   description: '查看当前目录下的所有文件和文件夹', 
  //   tags: ['1', '5'], 
  //   collections: ['1', '5'], 
  //   sortValue: 1, 
  //   copyCount: 8, 
  //   systemType: ['windows'] 
  // },
  // { 
  //   id: '3', 
  //   name: 'Git提交', 
  //   content: 'git commit -m "提交信息"', 
  //   description: '提交代码到Git仓库', 
  //   tags: ['1', '2'], 
  //   collections: ['1', '2'], 
  //   sortValue: 2, 
  //   copyCount: 15, 
  //   systemType: ['windows', 'linux', 'mac'] 
  // },
  // { 
  //   id: '4', 
  //   name: 'Docker启动容器', 
  //   content: 'docker start container_name', 
  //   description: '启动指定的Docker容器', 
  //   tags: ['1', '2'], 
  //   collections: ['1', '3'], 
  //   sortValue: 3, 
  //   copyCount: 7, 
  //   systemType: ['windows', 'linux', 'mac'] 
  // },
  // { 
  //   id: '5', 
  //   name: '查看IP地址', 
  //   content: 'ipconfig', 
  //   description: '查看Windows系统的IP地址信息', 
  //   tags: ['3', '5'], 
  //   collections: ['1', '5'], 
  //   sortValue: 4, 
  //   copyCount: 12, 
  //   systemType: ['windows'] 
  // },
  // { 
  //   id: '6', 
  //   name: '查看IP地址', 
  //   content: 'ifconfig', 
  //   description: '查看Linux/Mac系统的IP地址信息', 
  //   tags: ['3', '4'], 
  //   collections: ['1', '4'], 
  //   sortValue: 4, 
  //   copyCount: 9, 
  //   systemType: ['linux', 'mac'] 
  // },
  // { 
  //   id: '7', 
  //   name: '创建目录', 
  //   content: 'mkdir directory_name', 
  //   description: '创建新的目录', 
  //   tags: ['1'], 
  //   collections: ['1', '4'], 
  //   sortValue: 5, 
  //   copyCount: 6, 
  //   systemType: ['linux', 'mac'] 
  // },
  // { 
  //   id: '8', 
  //   name: '创建目录', 
  //   content: 'mkdir directory_name', 
  //   description: '创建新的目录', 
  //   tags: ['1'], 
  //   collections: ['1', '5'], 
  //   sortValue: 5, 
  //   copyCount: 5, 
  //   systemType: ['windows'] 
  // },
  // { 
  //   id: '9', 
  //   name: '删除文件', 
  //   content: 'rm file_name', 
  //   description: '删除指定的文件', 
  //   tags: ['1', '2'], 
  //   collections: ['1', '4'], 
  //   sortValue: 6, 
  //   copyCount: 4, 
  //   systemType: ['linux', 'mac'] 
  // },
  // { 
  //   id: '10', 
  //   name: '删除文件', 
  //   content: 'del file_name', 
  //   description: '删除指定的文件', 
  //   tags: ['1', '2'], 
  //   collections: ['1', '5'], 
  //   sortValue: 6, 
  //   copyCount: 3, 
  //   systemType: ['windows'] 
  // }
]);

// 模拟菜单项数据
const menuItems = ref({
  topMenu: [
    { id: 'home', name: '首页', icon: '🏠' },
    { id: 'commands', name: '命令管理', icon: '⚡' },
    { id: 'collections', name: '集合管理', icon: '📁' },
    { id: 'tags', name: '标签管理', icon: '🏷️' }
  ],
  tags: [
    { id: 'all-tags', name: '全部标签', icon: '🏷️' },
    { id: 'dev', name: '开发', icon: '💻' },
    { id: 'ops', name: '运维', icon: '🔧' },
    { id: 'test', name: '测试', icon: '🧪' },
    { id: 'db', name: '数据库', icon: '🗃️' },
    { id: 'network', name: '网络', icon: '🌐' }
  ],
  collections: [
    { id: 'all-collections', name: '全部集合', icon: '📁' },
    { id: 'common', name: '常用命令', icon: '⭐' },
    { id: 'git', name: 'Git命令', icon: '🔖' },
    { id: 'docker', name: 'Docker命令', icon: '🐳' },
    { id: 'linux', name: 'Linux命令', icon: '🐧' },
    { id: 'windows', name: 'Windows命令', icon: '🪟' }
  ],
  all: [
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

// 切换活动菜单
function toggleActiveMenu(menuId) {
  activeMenu.value = menuId;
  
  // 根据menuType调用不同的API获取命令列表
  if (menuType.value === 'tags') {
    // 调用GetCommandsByTagID获取该标签下的命令
    GetCommandsByTagID(menuId).then((result) => {
      commands.value = result;
      console.log("根据标签获取命令成功:", result);
    }).catch((error) => {
      console.error("根据标签获取命令失败:", error);
    });
  } else if (menuType.value === 'collections') {
    // 调用GetCommandsByCollectionID获取该集合下的命令
    GetCommandsByCollectionID(menuId).then((result) => {
      commands.value = result;
      console.log("根据集合获取命令成功:", result);
    }).catch((error) => {
      console.error("根据集合获取命令失败:", error);
    });
  } else {
    // 默认获取所有命令
    GetCommands().then((result) => {
      commands.value = result;
      console.log("获取所有命令成功:", result);
    }).catch((error) => {
      console.error("获取所有命令失败:", error);
    });
  }
}

// 切换系统类型
function toggleSystemType(type) {
  const index = systemType.value.indexOf(type);
  if (index === -1) {
    systemType.value.push(type);
  } else {
    systemType.value.splice(index, 1);
  }
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
  commands.value.push(command);
  // 这里可以添加保存到后端的逻辑
}

// 新增集合
function addCollection(collection) {
  collections.value.push(collection);
  // 这里可以添加保存到后端的逻辑
}

// 新增标签
function addTag(tag) {
  tags.value.push(tag);
  // 这里可以添加保存到后端的逻辑
}

// 编辑项目
function editItem(item) {
  selectedItem.value = item;
  // 这里可以添加编辑逻辑
}

// 删除项目
function deleteItem(item) {
  if (confirm(`确定要删除 ${item.name} 吗？`)) {
    if (item.tags) { // 是命令
      const index = commands.value.findIndex(cmd => cmd.id === item.id);
      if (index !== -1) {
        commands.value.splice(index, 1);
      }
    } else if (item.collections) { // 是集合
      const index = collections.value.findIndex(col => col.id === item.id);
      if (index !== -1) {
        collections.value.splice(index, 1);
      }
    } else { // 是标签
      const index = tags.value.findIndex(tag => tag.id === item.id);
      if (index !== -1) {
        tags.value.splice(index, 1);
      }
    }
  }
}

// 复制到剪贴板
function copyToClipboard(content) {
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
  
  // 获取命令
  GetCommands().then((result) => {
    commands.value = result;
  }).catch((error) => {
    console.error("获取命令失败:", error);
  });

  // 获取标签
  GetTags().then((result) => {
    tags.value = result;
    console.log("获取标签成功:", result);
  }).catch((error) => {
    console.error("获取标签失败:", error);
  });

  // 获取集合
  GetCollections().then((result) => {
    collections.value = result;
    console.log("获取集合成功:", result);
  }).catch((error) => {
    console.error("获取集合失败:", error);
  });

  // 初始化模拟数据
  console.log('初始化模拟数据');
});

// 组件卸载时
onUnmounted(() => {
  // 解绑点击事件
  document.removeEventListener('click', handleClickOutside);
});
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