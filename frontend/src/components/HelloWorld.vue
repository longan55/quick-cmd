<script setup>
import {reactive, ref, onMounted, onUnmounted} from 'vue'
import {Greet, GetMenuItems} from '../../wailsjs/go/main/App'

// 新增界面状态管理
const activeAddInterface = ref('') // 空字符串表示显示查询界面，否则显示对应的新增界面（command, collection, tag）
const selectedItem = ref(null) // 当前选中的集合或标签

// 新增表单数据
const newCommand = reactive({
  name: '',
  content: '',
  description: '',
  tagIDs: [],
  collectionIDs: []
})

const newCollection = reactive({
  name: '',
  description: '',
  commands: [] // 存储命令ID列表
})

const newTag = reactive({
  name: '',
  description: '',
  commands: [] // 存储命令ID列表
})

// 用于选择已存在命令的临时数据
const selectExistingCommand = ref(false)
const searchExistingCommands = ref('')

// 模拟数据 - 实际应该从后端获取
const mockCommands = ref([])
const mockCollections = ref([])
const mockTags = ref([])

// 复制提示状态
const showCopySuccess = ref(false)

const activeMenu = ref('home') // 改为字符串，支持单选
const menuType = ref('tags') // tags, collections, all
const systemType = ref(['windows']) // 改为数组，支持多选
const searchKeyword = ref('') // 搜索关键词
const menuItems = ref({}) // 存储从后端获取的菜单项

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

// 获取菜单项数据
function fetchMenuItems() {
  GetMenuItems().then(result => {
    menuItems.value = result
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
  activeMenu.value = menu // 直接设置为当前点击的菜单，实现单选
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
  fetchMenuItems() // 组件挂载时获取菜单项
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

// 重置表单数据
function resetCommandForm() {
  newCommand.name = ''
  newCommand.content = ''
  newCommand.description = ''
  newCommand.tagIDs = []
  newCommand.collectionIDs = []
}

function resetCollectionForm() {
  newCollection.name = ''
  newCollection.description = ''
  newCollection.commands = []
}

function resetTagForm() {
  newTag.name = ''
  newTag.description = ''
  newTag.commands = []
}

// 新增指令相关方法
function saveCommand() {
  // 验证表单
  if (!newCommand.name || !newCommand.content) {
    alert('请填写指令名称和内容')
    return
  }
  
  // 创建新指令（实际应该调用后端API）
  const cmdId = 'cmd-' + Date.now()
  const newCmd = {
    id: cmdId,
    name: newCommand.name,
    content: newCommand.content,
    description: newCommand.description,
    tagIDs: newCommand.tagIDs,
    collectionIDs: newCommand.collectionIDs,
    copyCount: 0,
    searchCount: 0,
    createdAt: new Date(),
    updatedAt: new Date(),
    DeletedAt: null
  }
  
  mockCommands.value.push(newCmd)
  
  // 更新标签和集合中的命令列表
  updateTagsAndCollectionsWithCommand(cmdId)
  
  // 重置表单并关闭界面
  resetCommandForm()
  activeAddInterface.value = ''
}

// 新增集合相关方法
function addNewCommandToCollection() {
  if (!newCommand.content) {
    alert('请输入指令内容')
    return
  }
  
  // 创建新指令
  const cmdId = 'cmd-' + Date.now()
  const newCmd = {
    id: cmdId,
    name: '新指令',
    content: newCommand.content,
    description: '',
    tagIDs: [],
    collectionIDs: [],
    copyCount: 0,
    searchCount: 0,
    createdAt: new Date(),
    updatedAt: new Date(),
    DeletedAt: null
  }
  
  mockCommands.value.push(newCmd)
  newCollection.commands.push(cmdId)
  newCommand.content = ''
}

function addExistingCommandToCollection(cmdId) {
  if (!newCollection.commands.includes(cmdId)) {
    newCollection.commands.push(cmdId)
  }
}

function removeCommandFromCollection(cmdId) {
  const index = newCollection.commands.indexOf(cmdId)
  if (index > -1) {
    newCollection.commands.splice(index, 1)
  }
}

function saveCollection() {
  if (!newCollection.name) {
    alert('请填写集合名称')
    return
  }
  
  // 创建新集合
  const collId = 'coll-' + Date.now()
  const newColl = {
    id: collId,
    name: newCollection.name,
    description: newCollection.description,
    commands: [...newCollection.commands],
    searchCount: 0,
    createdAt: new Date(),
    updatedAt: new Date(),
    DeletedAt: null
  }
  
  mockCollections.value.push(newColl)
  
  // 更新命令中的集合ID
  updateCommandsWithCollection(collId)
  
  // 重置表单并关闭界面
  resetCollectionForm()
  activeAddInterface.value = ''
}

// 新增标签相关方法
function addNewCommandToTag() {
  if (!newCommand.content) {
    alert('请输入指令内容')
    return
  }
  
  // 创建新指令
  const cmdId = 'cmd-' + Date.now()
  const newCmd = {
    id: cmdId,
    name: '新指令',
    content: newCommand.content,
    description: '',
    tagIDs: [],
    collectionIDs: [],
    copyCount: 0,
    searchCount: 0,
    createdAt: new Date(),
    updatedAt: new Date(),
    DeletedAt: null
  }
  
  mockCommands.value.push(newCmd)
  newTag.commands.push(cmdId)
  newCommand.content = ''
}

function addExistingCommandToTag(cmdId) {
  if (!newTag.commands.includes(cmdId)) {
    newTag.commands.push(cmdId)
  }
}

function removeCommandFromTag(cmdId) {
  const index = newTag.commands.indexOf(cmdId)
  if (index > -1) {
    newTag.commands.splice(index, 1)
  }
}

function saveTag() {
  if (!newTag.name) {
    alert('请填写标签名称')
    return
  }
  
  // 创建新标签
  const tagId = 'tag-' + Date.now()
  const newTagItem = {
    id: tagId,
    name: newTag.name,
    description: newTag.description,
    commands: [...newTag.commands],
    searchCount: 0,
    createdAt: new Date(),
    updatedAt: new Date(),
    DeletedAt: null
  }
  
  mockTags.value.push(newTagItem)
  
  // 更新命令中的标签ID
  updateCommandsWithTag(tagId)
  
  // 重置表单并关闭界面
  resetTagForm()
  activeAddInterface.value = ''
}

// 更新命令中的标签和集合信息
function updateCommandsWithTag(tagId) {
  const tag = mockTags.value.find(t => t.id === tagId)
  if (tag) {
    tag.commands.forEach(cmdId => {
      const cmd = mockCommands.value.find(c => c.id === cmdId)
      if (cmd && !cmd.tagIDs.includes(tagId)) {
        cmd.tagIDs.push(tagId)
      }
    })
  }
}

function updateCommandsWithCollection(collId) {
  const coll = mockCollections.value.find(c => c.id === collId)
  if (coll) {
    coll.commands.forEach(cmdId => {
      const cmd = mockCommands.value.find(c => c.id === cmdId)
      if (cmd && !cmd.collectionIDs.includes(collId)) {
        cmd.collectionIDs.push(collId)
      }
    })
  }
}

// 更新标签和集合中的命令列表
function updateTagsAndCollectionsWithCommand(cmdId) {
  newCommand.tagIDs.forEach(tagId => {
    const tag = mockTags.value.find(t => t.id === tagId)
    if (tag && !tag.commands.includes(cmdId)) {
      tag.commands.push(cmdId)
    }
  })
  
  newCommand.collectionIDs.forEach(collId => {
    const coll = mockCollections.value.find(c => c.id === collId)
    if (coll && !coll.commands.includes(cmdId)) {
      coll.commands.push(cmdId)
    }
  })
}

// 删除相关方法
function deleteCommand(cmdId) {
  if (confirm('确定要删除这条指令吗？')) {
    // 实际应该调用后端API进行软删除
    const cmdIndex = mockCommands.value.findIndex(cmd => cmd.id === cmdId)
    if (cmdIndex > -1) {
      mockCommands.value[cmdIndex].deletedAt = new Date()
    }
    
    // 从所有标签和集合中移除该命令
    mockTags.value.forEach(tag => {
      const cmdPos = tag.commands.indexOf(cmdId)
      if (cmdPos > -1) {
        tag.commands.splice(cmdPos, 1)
      }
    })
    
    mockCollections.value.forEach(coll => {
      const cmdPos = coll.commands.indexOf(cmdId)
      if (cmdPos > -1) {
        coll.commands.splice(cmdPos, 1)
      }
    })
  }
}

function deleteCollectionOrTag(itemId) {
  const item = menuType.value === 'collections' 
    ? mockCollections.value.find(coll => coll.id === itemId)
    : mockTags.value.find(tag => tag.id === itemId)
    
  if (item && confirm(`确定要删除${item.name}吗？`)) {
    // 实际应该调用后端API进行软删除
    item.deletedAt = new Date()
    
    // 从命令中移除关联
    item.commands.forEach(cmdId => {
      const cmd = mockCommands.value.find(c => c.id === cmdId)
      if (cmd) {
        if (menuType.value === 'collections') {
          const collIndex = cmd.collectionIDs.indexOf(itemId)
          if (collIndex > -1) {
            cmd.collectionIDs.splice(collIndex, 1)
          }
        } else {
          const tagIndex = cmd.tagIDs.indexOf(itemId)
          if (tagIndex > -1) {
            cmd.tagIDs.splice(tagIndex, 1)
          }
        }
      }
    })
  }
}

// 复制相关方法
function copyToClipboard(text) {
  if (!text) return
  
  navigator.clipboard.writeText(text)
    .then(() => {
      // 显示复制成功提示
      showCopySuccess.value = true
      setTimeout(() => {
        showCopySuccess.value = false
      }, 2000)
      
      // 更新命令的复制次数
      const cmd = mockCommands.value.find(c => c.content === text)
      if (cmd) {
        cmd.copyCount++
      }
    })
    .catch(err => {
      console.error('复制失败:', err)
    })
}

function copyAllCommands(item) {
  if (!item.commands || item.commands.length === 0) {
    alert('该' + (menuType.value === 'collections' ? '集合' : '标签') + '中没有指令')
    return
  }
  
  const commands = item.commands.map(cmdId => {
    const cmd = mockCommands.value.find(c => c.id === cmdId)
    return cmd ? cmd.content : ''
  }).filter(Boolean).join('\n')
  
  copyToClipboard(commands)
}

// 模拟初始化数据
function initMockData() {
  // 初始化模拟标签
  mockTags.value = [
    {
      id: 'tag-1',
      name: '工作',
      description: '工作相关指令',
      commands: [],
      searchCount: 0,
      createdAt: new Date(Date.now() - 86400000),
      updatedAt: new Date(Date.now() - 86400000),
      deletedAt: null
    },
    {
      id: 'tag-2',
      name: '学习',
      description: '学习相关指令',
      commands: [],
      searchCount: 0,
      createdAt: new Date(Date.now() - 86400000 * 2),
      updatedAt: new Date(Date.now() - 86400000 * 2),
      deletedAt: null
    }
  ]
  
  // 初始化模拟集合
  mockCollections.value = [
    {
      id: 'coll-1',
      name: 'Git常用命令',
      description: 'Git版本控制常用指令',
      commands: [],
      searchCount: 0,
      createdAt: new Date(Date.now() - 86400000 * 3),
      updatedAt: new Date(Date.now() - 86400000 * 3),
      deletedAt: null
    },
    {
      id: 'coll-2',
      name: '项目构建',
      description: '项目构建相关指令',
      commands: [],
      searchCount: 0,
      createdAt: new Date(Date.now() - 86400000 * 4),
      updatedAt: new Date(Date.now() - 86400000 * 4),
      deletedAt: null
    }
  ]
  
  // 初始化模拟指令
  const mockCmdData = [
    {
      id: 'cmd-1',
      name: 'Git提交',
      content: 'git commit -m "feat: add new feature"',
      description: '提交代码到本地仓库',
      tagIDs: ['tag-1'],
      collectionIDs: ['coll-1'],
      copyCount: 5,
      searchCount: 3,
      createdAt: new Date(Date.now() - 86400000 * 5),
      updatedAt: new Date(Date.now() - 86400000 * 5),
      deletedAt: null
    },
    {
      id: 'cmd-2',
      name: 'Git推送到远程',
      content: 'git push origin main',
      description: '将本地提交推送到远程仓库',
      tagIDs: ['tag-1'],
      collectionIDs: ['coll-1'],
      copyCount: 3,
      searchCount: 2,
      createdAt: new Date(Date.now() - 86400000 * 6),
      updatedAt: new Date(Date.now() - 86400000 * 6),
      deletedAt: null
    },
    {
      id: 'cmd-3',
      name: 'npm安装依赖',
      content: 'npm install',
      description: '安装项目依赖',
      tagIDs: ['tag-1'],
      collectionIDs: ['coll-2'],
      copyCount: 8,
      searchCount: 5,
      createdAt: new Date(Date.now() - 86400000 * 7),
      updatedAt: new Date(Date.now() - 86400000 * 7),
      deletedAt: null
    },
    {
      id: 'cmd-4',
      name: 'npm构建',
      content: 'npm run build',
      description: '构建项目',
      tagIDs: ['tag-1'],
      collectionIDs: ['coll-2'],
      copyCount: 6,
      searchCount: 4,
      createdAt: new Date(Date.now() - 86400000 * 8),
      updatedAt: new Date(Date.now() - 86400000 * 8),
      deletedAt: null
    },
    {
      id: 'cmd-5',
      name: 'Go运行',
      content: 'go run main.go',
      description: '运行Go程序',
      tagIDs: ['tag-2'],
      collectionIDs: [],
      copyCount: 4,
      searchCount: 3,
      createdAt: new Date(Date.now() - 86400000 * 9),
      updatedAt: new Date(Date.now() - 86400000 * 9),
      deletedAt: null
    }
  ]
  
  mockCmdData.forEach(cmd => {
    mockCommands.value.push(cmd)
    
    // 更新标签和集合中的命令列表
    cmd.tagIDs.forEach(tagId => {
      const tag = mockTags.value.find(t => t.id === tagId)
      if (tag && !tag.commands.includes(cmd.id)) {
        tag.commands.push(cmd.id)
      }
    })
    
    cmd.collectionIDs.forEach(collId => {
      const coll = mockCollections.value.find(c => c.id === collId)
      if (coll && !coll.commands.includes(cmd.id)) {
        coll.commands.push(cmd.id)
      }
    })
  })
}

// 组件挂载时初始化模拟数据
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  fetchMenuItems() // 组件挂载时获取菜单项
  initMockData() // 初始化模拟数据
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
          <button 
            v-for="menu in menuItems.tags" 
            :key="menu.id" 
            class="menu-button" 
            @click="toggleActiveMenu(menu.id)" 
            :class="{ active: activeMenu === menu.id }"
          >
            <span class="menu-icon">{{ menu.icon }}</span>
            <span class="menu-text">{{ menu.name }}</span>
          </button>
        </div>

        <!-- 集合菜单 -->
        <div v-else-if="menuType === 'collections'" class="menu-buttons">
          <button 
            v-for="menu in menuItems.collections" 
            :key="menu.id" 
            class="menu-button" 
            @click="toggleActiveMenu(menu.id)" 
            :class="{ active: activeMenu === menu.id }"
          >
            <span class="menu-icon">{{ menu.icon }}</span>
            <span class="menu-text">{{ menu.name }}</span>
          </button>
        </div>

        <!-- 全部菜单 -->
        <div v-else-if="menuType === 'all'" class="menu-buttons">
          <button 
            v-for="menu in menuItems.all" 
            :key="menu.id" 
            class="menu-button" 
            @click="toggleActiveMenu(menu.id)" 
            :class="{ active: activeMenu === menu.id }"
          >
            <span class="menu-icon">{{ menu.icon }}</span>
            <span class="menu-text">{{ menu.name }}</span>
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
          <button class="add-button" @click="activeAddInterface = 'command'">
            🎇 新增指令
          </button>
          <button class="add-button" @click="activeAddInterface = 'collection'">
            { } 新增集合
          </button>
          <button class="add-button" @click="activeAddInterface = 'tag'">
            🏷️ 新增标签
          </button>
        </div>
      </header>

      <!-- 内容区域 -->
      <section class="content-area">
        <!-- 复制成功提示 -->
        <div v-if="showCopySuccess" class="copy-success">
          已复制到剪贴板！
        </div>
        
        <!-- 新增指令界面 -->
        <div v-if="activeAddInterface === 'command'" class="add-interface">
          <h2>🎇 新增指令</h2>
          <div class="form-group">
            <label for="cmd-name">指令名称：</label>
            <input type="text" id="cmd-name" v-model="newCommand.name" class="form-input" placeholder="请输入指令名称" />
          </div>
          <div class="form-group">
            <label for="cmd-content">指令内容：</label>
            <textarea id="cmd-content" v-model="newCommand.content" class="form-textarea" placeholder="请输入指令内容" rows="3"></textarea>
          </div>
          <div class="form-group">
            <label for="cmd-desc">指令说明：</label>
            <textarea id="cmd-desc" v-model="newCommand.description" class="form-textarea" placeholder="请输入指令说明" rows="2"></textarea>
          </div>
          <div class="form-group">
            <label>所属标签：</label>
            <div class="tag-selector">
              <div v-for="tag in mockTags" :key="tag.id" class="tag-option">
                <input 
                  type="checkbox" 
                  :id="`tag-${tag.id}`" 
                  :value="tag.id" 
                  v-model="newCommand.tagIDs"
                />
                <label :for="`tag-${tag.id}`">{{ tag.name }}</label>
              </div>
            </div>
          </div>
          <div class="form-group">
            <label>所属集合：</label>
            <div class="collection-selector">
              <div v-for="collection in mockCollections" :key="collection.id" class="collection-option">
                <input 
                  type="checkbox" 
                  :id="`collection-${collection.id}`" 
                  :value="collection.id" 
                  v-model="newCommand.collectionIDs"
                />
                <label :for="`collection-${collection.id}`">{{ collection.name }}</label>
              </div>
            </div>
          </div>
          <div class="form-actions">
            <button class="btn btn-primary" @click="saveCommand">保存指令</button>
            <button class="btn btn-cancel" @click="activeAddInterface = ''">取消</button>
          </div>
        </div>

        <!-- 新增集合界面 -->
        <div v-else-if="activeAddInterface === 'collection'" class="add-interface">
          <h2>{ } 新增集合</h2>
          <div class="form-group">
            <label for="coll-name">集合名称：</label>
            <input type="text" id="coll-name" v-model="newCollection.name" class="form-input" placeholder="请输入集合名称" />
          </div>
          <div class="form-group">
            <label for="coll-desc">集合说明：</label>
            <textarea id="coll-desc" v-model="newCollection.description" class="form-textarea" placeholder="请输入集合说明" rows="2"></textarea>
          </div>
          <div class="form-group">
            <label>添加指令：</label>
            <div class="command-selector">
              <div class="add-command-options">
                <button 
                  class="btn btn-secondary" 
                  :class="{ active: !selectExistingCommand }" 
                  @click="selectExistingCommand = false"
                >
                  新增指令
                </button>
                <button 
                  class="btn btn-secondary" 
                  :class="{ active: selectExistingCommand }" 
                  @click="selectExistingCommand = true"
                >
                  选择已存在指令
                </button>
              </div>
              <div v-if="!selectExistingCommand" class="new-command-form">
                <input type="text" v-model="newCommand.content" class="form-input" placeholder="请输入指令内容" />
                <button class="btn btn-small" @click="addNewCommandToCollection">添加</button>
              </div>
              <div v-else class="existing-command-selector">
                <input 
                  type="text" 
                  v-model="searchExistingCommands" 
                  class="form-input" 
                  placeholder="搜索已存在的指令..." 
                />
                <div class="command-list">
                  <div 
                    v-for="cmd in mockCommands.filter(cmd => cmd.content.includes(searchExistingCommands))" 
                    :key="cmd.id" 
                    class="command-item"
                    @click="addExistingCommandToCollection(cmd.id)"
                  >
                    <button class="btn btn-small remove-button" @click.stop="removeCommandFromCollection(cmd.id)">×</button>
                    {{ cmd.content }}
                  </div>
                </div>
              </div>
              <div class="selected-commands">
                <div 
                  v-for="cmdId in newCollection.commands" 
                  :key="cmdId" 
                  class="selected-command-item"
                >
                  {{ mockCommands.find(c => c.id === cmdId)?.content || '未找到' }}
                  <button class="remove-button" @click="removeCommandFromCollection(cmdId)">×</button>
                </div>
              </div>
            </div>
          </div>
          <div class="form-actions">
            <button class="btn btn-primary" @click="saveCollection">保存集合</button>
            <button class="btn btn-cancel" @click="activeAddInterface = ''">取消</button>
          </div>
        </div>

        <!-- 新增标签界面 -->
        <div v-else-if="activeAddInterface === 'tag'" class="add-interface">
          <h2>🏷️ 新增标签</h2>
          <div class="form-group">
            <label for="tag-name">标签名称：</label>
            <input type="text" id="tag-name" v-model="newTag.name" class="form-input" placeholder="请输入标签名称" />
          </div>
          <div class="form-group">
            <label for="tag-desc">标签说明：</label>
            <textarea id="tag-desc" v-model="newTag.description" class="form-textarea" placeholder="请输入标签说明" rows="2"></textarea>
          </div>
          <div class="form-group">
            <label>添加指令：</label>
            <div class="command-selector">
              <div class="add-command-options">
                <button 
                  class="btn btn-secondary" 
                  :class="{ active: !selectExistingCommand }" 
                  @click="selectExistingCommand = false"
                >
                  新增指令
                </button>
                <button 
                  class="btn btn-secondary" 
                  :class="{ active: selectExistingCommand }" 
                  @click="selectExistingCommand = true"
                >
                  选择已存在指令
                </button>
              </div>
              <div v-if="!selectExistingCommand" class="new-command-form">
                <input type="text" v-model="newCommand.content" class="form-input" placeholder="请输入指令内容" />
                <button class="btn btn-small" @click="addNewCommandToTag">添加</button>
              </div>
              <div v-else class="existing-command-selector">
                <input 
                  type="text" 
                  v-model="searchExistingCommands" 
                  class="form-input" 
                  placeholder="搜索已存在的指令..." 
                />
                <div class="command-list">
                  <div 
                    v-for="cmd in mockCommands.filter(cmd => cmd.content.includes(searchExistingCommands))" 
                    :key="cmd.id" 
                    class="command-item"
                    @click="addExistingCommandToTag(cmd.id)"
                  >
                    <button class="btn btn-small remove-button" @click.stop="removeCommandFromTag(cmd.id)">×</button>
                    {{ cmd.content }}
                  </div>
                </div>
              </div>
              <div class="selected-commands">
                <div 
                  v-for="cmdId in newTag.commands" 
                  :key="cmdId" 
                  class="selected-command-item"
                >
                  {{ mockCommands.find(c => c.id === cmdId)?.content || '未找到' }}
                  <button class="remove-button" @click="removeCommandFromTag(cmdId)">×</button>
                </div>
              </div>
            </div>
          </div>
          <div class="form-actions">
            <button class="btn btn-primary" @click="saveTag">保存标签</button>
            <button class="btn btn-cancel" @click="activeAddInterface = ''">取消</button>
          </div>
        </div>

        <!-- 查询界面 -->
        <div v-else class="query-interface">
          <h2>{{ menuType === 'tags' ? '🏷️ 标签' : menuType === 'collections' ? '{ } 集合' : '🎇 指令' }}</h2>
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>{{ menuType === 'tags' ? '标签名称' : menuType === 'collections' ? '集合名称' : '指令' }}</th>
                  <th v-if="menuType !== 'all'">包含指令</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <!-- 集合/标签表格 -->
                <template v-if="menuType === 'collections' || menuType === 'tags'">
                  <tr v-for="item in (menuType === 'collections' ? mockCollections : mockTags)" :key="item.id">
                    <td>{{ item.name }}</td>
                    <td>
                      <div class="command-list-in-cell">
                        <div 
                          v-for="cmdId in (item.commands || [])" 
                          :key="cmdId" 
                          class="command-in-cell"
                          @click="copyToClipboard(mockCommands.find(c => c.id === cmdId)?.content)"
                        >
                          {{ mockCommands.find(c => c.id === cmdId)?.content || '已删除' }}
                        </div>
                      </div>
                    </td>
                    <td>
                      <button class="btn btn-action" @click="deleteCollectionOrTag(item.id)">删除</button>
                      <button class="btn btn-action" @click="copyAllCommands(item)">复制全部</button>
                    </td>
                  </tr>
                </template>
                <!-- 指令表格 -->
                <template v-else>
                  <tr v-for="cmd in mockCommands" :key="cmd.id">
                    <td>{{ cmd.content }}</td>
                    <td>
                      <button class="btn btn-action" @click="deleteCommand(cmd.id)">删除</button>
                      <button class="btn btn-action" @click="copyToClipboard(cmd.content)">复制</button>
                    </td>
                  </tr>
                </template>
              </tbody>
            </table>
          </div>
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
          <p>作者: longan55</p>
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
.top-menu-bar .menu-button {
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
.top-menu-bar .menu-button:hover {
  background-color: rgba(0, 0, 0, 0.1);
}

/* 左侧菜单栏 */
.sidebar {
  width: 200px;
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
  color: black;
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

/* 顶部导航栏内容 */
.top-nav-content {
  text-align: center;
}

.logo h1 {
  margin: 0;
  font-size: 2rem;
  color: #3498db;
}

.user-info {
  margin-top: 10px;
  font-size: 1rem;
  color: #666;
}

/* 内容区域 */
.content-area {
  flex: 1;
  background-color: #ecf0f1;
  padding: 0;
  margin: 0;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

/* 嵌套组件区域 */
.nested-component-area {
  background-color: white;
  padding: 0;
  margin: 0;
  border-radius: 0;
  box-shadow: none;
  flex: 1;
  height: 100%;
  width: 100%;
  overflow: hidden;
}

/* 设置页面样式 */
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
  transition: all 0.3s ease;
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
  background-color: rgba(255, 245, 245, 0.5);
  z-index: 2000;
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 关于对话框容器 */
.about-dialog {
  background-color: white;
  border-radius: 10px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  width: 400px;
  max-width: 90%;
  overflow: hidden;
}

/* 关于对话框头部 */
.about-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 15px;
  background-color: #f8f9fa;
  color: #333;
}

/* 关于对话框关闭按钮样式 */
.about-dialog .close-button {
  color: #333;
}

.about-dialog .close-button:hover {
  background-color: rgba(0, 0, 0, 0.1);
}

.about-header h2 {
  margin: 0;
  font-size: 1.2rem;
  font-weight: normal;
}

/* 关于对话框内容区域 */
.about-content {
  padding: 20px;
  text-align: center;
  line-height: 1.6;
}

.about-content p {
  margin: 8px 0;
  color: #333;
}

/* 新增按钮样式 */
.add-button {
  padding: 10px 20px;
  margin: 0 10px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 5px;
}

.add-button:hover {
  background-color: #2980b9;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(52, 152, 219, 0.3);
}

.add-button:active {
  transform: translateY(0);
  box-shadow: 0 2px 4px rgba(52, 152, 219, 0.2);
}

/* 顶部导航栏内容调整 */
.top-nav-content {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  gap: 20px;
}

/* 新增界面样式 */
.add-interface {
  background-color: white;
  margin: 20px;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  color: black;
}

.add-interface h2 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
  font-size: 1.5rem;
}

/* 表单样式 */
.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  color: #555;
  font-weight: 500;
}

.form-input, .form-textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
  box-sizing: border-box;
}

.form-input:focus, .form-textarea:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

/* 按钮样式统一 */
.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 5px;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 5px;
}

.btn-primary {
  background-color: #3498db;
  color: white;
}

.btn-primary:hover {
  background-color: #2980b9;
}

.btn-cancel {
  background-color: #e74c3c;
  color: white;
}

.btn-cancel:hover {
  background-color: #c0392b;
}

.btn-secondary {
  background-color: #95a5a6;
  color: white;
}

.btn-secondary:hover {
  background-color: #7f8c8d;
}

.btn-secondary.active {
  background-color: #3498db;
}

.btn-small {
  padding: 5px 10px;
  font-size: 0.8rem;
}

.btn-action {
  padding: 5px 10px;
  font-size: 0.8rem;
  margin-right: 5px;
  background-color: #3498db;
  color: white;
}

.btn-action:hover {
  background-color: #2980b9;
}

/* 表单操作按钮容器 */
.form-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

/* 查询界面样式 */
.query-interface {
  background-color: white;
  margin: 20px;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  color: black;
}

.query-interface h2 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #333;
  font-size: 1.5rem;
}

/* 表格样式 */
.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.data-table th, .data-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ddd;
  color: black;
}

.data-table th {
  background-color: #f8f9fa;
  font-weight: 600;
  color: #555;
}

.data-table tr:hover {
  background-color: #f5f5f5;
}

/* 复制成功提示 */
.copy-success {
  position: fixed;
  bottom: 20px;
  right: 20px;
  background-color: #2ecc71;
  color: white;
  padding: 10px 20px;
  border-radius: 5px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  z-index: 1000;
  animation: fadeInOut 3s ease;
}

@keyframes fadeInOut {
  0% { opacity: 0; transform: translateY(20px); }
  10% { opacity: 1; transform: translateY(0); }
  90% { opacity: 1; transform: translateY(0); }
  100% { opacity: 0; transform: translateY(20px); }
}

/* 标签和集合选择器样式 */
.tag-selector, .collection-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.tag-option, .collection-option {
  display: flex;
  align-items: center;
  gap: 5px;
}

.tag-option input[type="checkbox"], .collection-option input[type="checkbox"] {
  accent-color: #3498db;
}

/* 命令选择器样式 */
.command-selector {
  margin-top: 10px;
}

.add-command-options {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.new-command-form {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
}

.existing-command-selector {
  margin-bottom: 10px;
}

.command-list {
  max-height: 150px;
  overflow-y: auto;
  border: 1px solid #ddd;
  border-radius: 5px;
  padding: 10px;
  margin-top: 5px;
}

.command-item {
  padding: 8px;
  cursor: pointer;
  border-radius: 3px;
  transition: background-color 0.2s ease;
  color: black;
}

.command-item:hover {
  background-color: rgba(52, 152, 219, 0.1);
}

/* 已选择命令样式 */
.selected-commands {
  margin-top: 10px;
  border-top: 1px solid #eee;
  padding-top: 10px;
}

.selected-command-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px 10px;
  background-color: #f0f8ff;
  border: 1px solid #3498db;
  border-radius: 5px;
  margin-bottom: 5px;
  color: black;
}

.remove-button {
  background-color: #e74c3c;
  color: white;
  border: none;
  border-radius: 3px;
  padding: 2px 6px;
  font-size: 0.8rem;
  cursor: pointer;
}

.remove-button:hover {
  background-color: #c0392b;
}

/* 表格中命令列表样式 */
.command-list-in-cell {
  max-height: 100px;
  overflow-y: auto;
}

.command-in-cell {
  padding: 5px 10px;
  background-color: #f0f8ff;
  border-radius: 3px;
  margin-bottom: 5px;
  cursor: pointer;
  font-family: monospace;
  font-size: 0.8rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 300px;
  color: black;
}

.command-in-cell:hover {
  background-color: #d4e6f1;
}
</style>







