<template>
  <div class="main-content">
    <!-- 新增指令界面 -->
    <div v-if="activeAddInterface === 'command'" class="add-interface">
      <div class="interface-header">
        <h2>新增指令</h2>
        <button class="close-button" @click="$emit('toggle-add-interface', 'command')">×</button>
      </div>
      
      <div class="form-container">
        <div class="form-group">
          <label for="command-name">指令名称</label>
          <input type="text" id="command-name" v-model="newCommand.name" placeholder="请输入指令名称" />
        </div>
        
        <div class="form-group">
          <label for="command-content">指令内容</label>
          <textarea id="command-content" v-model="newCommand.content" placeholder="请输入指令内容" rows="4"></textarea>
        </div>
        
        <div class="form-group">
          <label for="command-description">指令描述</label>
          <textarea id="command-description" v-model="newCommand.description" placeholder="请输入指令描述" rows="2"></textarea>
        </div>
        
        <div class="form-group">
          <label>所属标签</label>
          <el-select
            v-model="newCommand.tags"
            multiple
            filterable
            placeholder="请选择标签"
            style="width: 100%"
          >
            <el-option
              v-for="tag in allTags"
              :key="tag.id"
              :label="`[${tag.id}] ${tag.name}`"
              :value="tag.id"
            />
          </el-select>
        </div>
        
        <div class="form-group">
          <label>所属集合</label>
          <el-select
            v-model="newCommand.collections"
            multiple
            filterable
            placeholder="请选择集合"
            style="width: 100%"
          >
            <el-option
              v-for="collection in allCollections"
              :key="collection.id"
              :label="`[${collection.id}] ${collection.name}`"
              :value="collection.id"
            />
          </el-select>
        </div>
        
        <div class="form-group">
          <label>适用系统</label>
          <div class="system-selector">
            <label class="system-option">
              <input type="checkbox" value="windows" v-model="newCommand.systemType" />
              <span>Windows</span>
            </label>
            <label class="system-option">
              <input type="checkbox" value="linux" v-model="newCommand.systemType" />
              <span>Linux</span>
            </label>
            <label class="system-option">
              <input type="checkbox" value="mac" v-model="newCommand.systemType" />
              <span>Mac</span>
            </label>
          </div>
        </div>
        
        <div class="form-actions">
          <button class="cancel-button" @click="$emit('toggle-add-interface', 'command')">取消</button>
          <button class="save-button" @click="saveCommand">保存</button>
        </div>
      </div>
    </div>
    
    <!-- 新增集合界面 -->
    <div v-else-if="activeAddInterface === 'collection'" class="add-interface">
      <div class="interface-header">
        <h2>新增集合</h2>
        <button class="close-button" @click="$emit('toggle-add-interface', 'collection')">×</button>
      </div>
      
      <div class="form-container">
        <div class="form-group">
          <label for="collection-name">集合名称</label>
          <input type="text" id="collection-name" v-model="newCollection.name" placeholder="请输入集合名称" />
        </div>
        
        <div class="form-group">
          <label for="collection-description">集合描述</label>
          <textarea id="collection-description" v-model="newCollection.description" placeholder="请输入集合描述" rows="2"></textarea>
        </div>
        
        <div class="form-group">
          <label>关联指令</label>
          <el-select
            v-model="newCollection.commandIds"
            multiple
            filterable
            placeholder="请选择关联指令"
            style="width: 100%"
          >
            <el-option
              v-for="cmd in allCommands"
              :key="cmd.id"
              :label="`[${cmd.id}] ${cmd.name}`"
              :value="cmd.id"
            />
          </el-select>
        </div>
        
        <div class="form-group">
          <label>适用系统</label>
          <div class="system-selector">
            <label class="system-option">
              <input type="checkbox" value="windows" v-model="newCollection.systemType" />
              <span>Windows</span>
            </label>
            <label class="system-option">
              <input type="checkbox" value="linux" v-model="newCollection.systemType" />
              <span>Linux</span>
            </label>
            <label class="system-option">
              <input type="checkbox" value="mac" v-model="newCollection.systemType" />
              <span>Mac</span>
            </label>
          </div>
        </div>
        
        <div class="form-actions">
          <button class="cancel-button" @click="$emit('toggle-add-interface', 'collection')">取消</button>
          <button class="save-button" @click="saveCollection">保存</button>
        </div>
      </div>
    </div>
    
    <!-- 新增标签界面 -->
    <div v-else-if="activeAddInterface === 'tag'" class="add-interface">
      <div class="interface-header">
        <h2>新增标签</h2>
        <button class="close-button" @click="$emit('toggle-add-interface', 'tag')">×</button>
      </div>
      
      <div class="form-container">
        <div class="form-group">
          <label for="tag-name">标签名称</label>
          <input type="text" id="tag-name" v-model="newTag.name" placeholder="请输入标签名称" />
        </div>
      
        
        <div class="form-group">
          <label for="tag-description">标签描述</label>
          <textarea id="tag-description" v-model="newTag.description" placeholder="请输入标签描述" rows="2"></textarea>
        </div>
        
        <div class="form-group">
          <label>关联指令</label>
          <el-select
            v-model="newTag.commandIds"
            multiple
            filterable
            placeholder="请选择关联指令"
            style="width: 100%"
          >
            <el-option
              v-for="cmd in allCommands"
              :key="cmd.id"
              :label="`[${cmd.id}] ${cmd.name}`"
              :value="cmd.id"
            />
          </el-select>
        </div>
        
        <div class="form-group">
          <label>适用系统</label>
          <div class="system-selector">
            <label class="system-option">
              <input type="checkbox" value="windows" v-model="newTag.systemType" />
              <span>Windows</span>
            </label>
            <label class="system-option">
              <input type="checkbox" value="linux" v-model="newTag.systemType" />
              <span>Linux</span>
            </label>
            <label class="system-option">
              <input type="checkbox" value="mac" v-model="newTag.systemType" />
              <span>Mac</span>
            </label>
          </div>
        </div>
        
        <div class="form-actions">
          <button class="cancel-button" @click="$emit('toggle-add-interface', 'tag')">取消</button>
          <button class="save-button" @click="saveTag">保存</button>
        </div>
      </div>
    </div>
    
    <!-- 查询界面 -->
    <div v-else class="query-interface">
      <div class="interface-header">
        <h2>
          {{ activeMenu === '0' ? (menuType === 'tags' ? '标签列表' : menuType === 'collections' ? '集合列表' : '命令列表') : '命令列表' }}
        </h2>
        
        <div class="interface-actions">
          <button class="add-button" @click="$emit('toggle-add-interface', menuType === 'tags' ? 'tag' : menuType === 'collections' ? 'collection' : 'command')">
            + 新增{{ menuType === 'tags' ? '标签' : menuType === 'collections' ? '集合' : '命令' }}
          </button>
        </div>
      </div>
      
      <!-- 命令列表 -->
      <div v-if="menuType === 'all' || (menuType === 'tags' && activeMenu !== '0') || (menuType === 'collections' && activeMenu !== '0')" class="table-container">
        <table class="command-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>名称</th>
              <th>内容</th>
              <th>描述</th>
              <!-- <th>标签</th>
              <th>集合</th> -->
              <th>使用次数</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="command in filteredCommands" :key="command.id">
              <td>{{ command.id }}</td>
              <td>{{ command.name }}</td>
              <td class="command-content-cell">{{ command.content }}</td>
              <td>{{ command.description }}</td>
              <!-- <td>
                <span v-for="tagId in command.tags" :key="tagId" class="tag-badge">
                  {{ getTagName(tagId) }}
                </span>
              </td>
              <td>
                <span v-for="collId in command.collections" :key="collId" class="collection-badge">
                  {{ getCollectionName(collId) }}
                </span>
              </td> -->
              <td>{{ command.copyCount }}</td>
              <td class="action-buttons">
                <button class="copy-button" @click="$emit('copy-to-clipboard', command.content)">
                  复制
                </button>
                <button class="edit-button" @click="$emit('edit-item', command)">
                  编辑
                </button>
                <button class="delete-button" @click="$emit('delete-item', command)">
                  删除
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- 标签列表 -->
      <div v-else-if="menuType === 'tags' && activeMenu === '0'" class="table-container">
        <table class="item-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>名称</th>
              <th>描述</th>
              <th>命令数量</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in tags.filter(tag => tag.id !== 0)" :key="item.id">
              <td>{{ item.id }}</td>
              <td>{{ item.name }}</td>
              <td>{{ item.description }}</td>
              <td>{{ getItemCommandCount(item.id) }}</td>
              <td class="action-buttons">
                <button class="edit-button" @click="$emit('edit-item', item)">
                  编辑
                </button>
                <button class="delete-button" @click="$emit('delete-item', item)">
                  删除
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- 集合列表 -->
      <div v-else-if="menuType === 'collections' && activeMenu === '0'" class="table-container">
        <table class="item-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>名称</th>
              <th>描述</th>
              <th>命令数量</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in collections.filter(col => col.id !== 0)" :key="item.id">
              <td>{{ item.id }}</td>
              <td>{{ item.name }}</td>
              <td>{{ item.description }}</td>
              <td>{{ getItemCommandCount(item.id) }}</td>
              <td class="action-buttons">
                <button class="edit-button" @click="$emit('edit-item', item)">
                  编辑
                </button>
                <button class="delete-button" @click="$emit('delete-item', item)">
                  删除
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, defineProps, defineEmits, onMounted, watch } from 'vue';
import { CreateCommand, GetAllCommandsIDAndName, GetAllTagsIDAndName, GetAllCollectionsIDAndName } from '../../../wailsjs/go/main/App';

const props = defineProps({
  activeAddInterface: {
    type: String,
    required: true
  },
  menuType: {
    type: String,
    required: true
  },
  activeMenu: {
    type: String,
    required: true
  },
  tags: {
    type: Array,
    required: true
  },
  collections: {
    type: Array,
    required: true
  },
  commands: {
    type: Array,
    required: true
  },
  searchKeyword: {
    type: String,
    required: true
  }
});

const emit = defineEmits([
  'toggle-add-interface',
  'update:searchKeyword',
  'add-command',
  'add-collection',
  'add-tag',
  'edit-item',
  'delete-item',
  'copy-to-clipboard'
]);

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
  sortValue: 0,
  systemType: [],
  commandIds: []
});

const newTag = ref({
  id: '',
  name: '',
  description: '',
  sortValue: 0,
  systemType: [],
  commandIds: []
});

const allCommands = ref([]);
const allTags = ref([]);
const allCollections = ref([]);

onMounted(async () => {
  try {
    const [commands, tags, collections] = await Promise.all([
      GetAllCommandsIDAndName(),
      GetAllTagsIDAndName(),
      GetAllCollectionsIDAndName()
    ]);
    allCommands.value = commands;
    allTags.value = tags;
    allCollections.value = collections;
  } catch (error) {
    console.error('获取数据失败:', error);
  }
});

// 监听activeAddInterface的变化，当切换到不同界面时刷新对应数据
watch(
  () => props.activeAddInterface,
  (newVal) => {
    if (newVal === 'tag' || newVal === 'collection') {
      refreshAllCommands();
    } else if (newVal === 'command') {
      refreshAllTagsAndCollections();
    }
  }
);

// 计算属性：根据搜索关键词过滤命令
const filteredCommands = computed(() => {
  if (!props.searchKeyword) {
    return props.commands;
  }
  
  return props.commands.filter(command => 
    command.name.toLowerCase().includes(props.searchKeyword.toLowerCase()) ||
    command.content.toLowerCase().includes(props.searchKeyword.toLowerCase()) ||
    command.description.toLowerCase().includes(props.searchKeyword.toLowerCase()) ||
    command.tags.some(tagId => getTagName(tagId).toLowerCase().includes(props.searchKeyword.toLowerCase())) ||
    command.collections.some(collId => getCollectionName(collId).toLowerCase().includes(props.searchKeyword.toLowerCase()))
  );
});

// 获取标签名称
function getTagName(tagId) {
  const tag = props.tags.find(tag => tag.id === tagId);
  return tag ? tag.name : '';
}

// 获取集合名称
function getCollectionName(collId) {
  const collection = props.collections.find(collection => collection.id === collId);
  return collection ? collection.name : '';
}

// 获取标签/集合的命令数量
function getItemCommandCount(itemId) {
  return props.commands.filter(command => 
    props.menuType === 'tags' 
      ? command.tags.includes(itemId) 
      : command.collections.includes(itemId)
  ).length;
}

// 生成UUID
function generateUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

// 保存命令
async function saveCommand() {
  if (!newCommand.value.name || !newCommand.value.content) {
    alert('请填写指令名称和内容');
    return;
  }
  
  try {
    // 构造符合后端Command结构体的数据
    const commandData = {
      name: newCommand.value.name,
      content: newCommand.value.content,
      description: newCommand.value.description,
      os: newCommand.value.systemType,
      tagIDs: newCommand.value.tags.map(tag => typeof tag === 'object' ? tag.id : tag), // 提取标签ID
      collectionIDs: newCommand.value.collections.map(collection => typeof collection === 'object' ? collection.id : collection) // 提取集合ID
    };
    
    console.log('保存命令数据:', commandData);
    
    // 调用后端API创建命令
    await CreateCommand(commandData);
    
    // 显示成功消息
    alert('命令创建成功！');
    
    // 重置表单并关闭界面
    resetCommandForm();
    emit('toggle-add-interface', 'command');
    
    // 通知父组件刷新数据
    emit('refresh-data');
    
  } catch (error) {
    console.error('创建命令失败:', error);
    alert(`创建命令失败: ${error.message || error}`);
  }
}

// 保存集合
function saveCollection() {
  if (!newCollection.value.name) {
    alert('请填写集合名称');
    return;
  }
  
  // 构造符合后端Collection结构体的数据
  const collectionData = {
    name: newCollection.value.name,
    description: newCollection.value.description,
    os: newCollection.value.systemType,
    commandIds: newCollection.value.commandIds
  };
  
  emit('add-collection', collectionData);
  
  // 重置表单并关闭界面
  resetCollectionForm();
  emit('toggle-add-interface', 'collection');
  
  // 通知父组件刷新数据
  emit('refresh-data');
}

// 保存标签
function saveTag() {
  if (!newTag.value.name) {
    alert('请填写标签名称');
    return;
  }
  
  // 构造符合后端Tag结构体的数据
  const tagData = {
    name: newTag.value.name,
    description: newTag.value.description,
    searchCount: 0,
    os: newTag.value.systemType,
    commandIds: newTag.value.commandIds
  };
  
  emit('add-tag', tagData);
  
  // 重置表单并关闭界面
  resetTagForm();
  emit('toggle-add-interface', 'tag');
  
  // 通知父组件刷新数据
  emit('refresh-data');
}

// 重置表单
function resetCommandForm() {
  newCommand.value = {
    id: '',
    name: '',
    content: '',
    description: '',
    tags: [],
    collections: [],
    sortValue: 0,
    copyCount: 0,
    systemType: ['windows']
  };
}

function resetCollectionForm() {
  newCollection.value = {
    id: '',
    name: '',
    description: '',
    sortValue: 0,
    systemType: [],
    commandIds: []
  };
}

function resetTagForm() {
  newTag.value = {
    id: '',
    name: '',
    description: '',
    sortValue: 0,
    systemType: [],
    commandIds: []
  };
}

// 刷新所有指令数据
async function refreshAllCommands() {
  try {
    const commands = await GetAllCommandsIDAndName();
    allCommands.value = commands;
  } catch (error) {
    console.error('刷新指令数据失败:', error);
  }
}

// 刷新所有标签和集合数据
async function refreshAllTagsAndCollections() {
  try {
    const [tags, collections] = await Promise.all([
      GetAllTagsIDAndName(),
      GetAllCollectionsIDAndName()
    ]);
    allTags.value = tags;
    allCollections.value = collections;
  } catch (error) {
    console.error('刷新标签和集合数据失败:', error);
  }
}
</script>

<style scoped>
/* 主内容区域样式 */
.main-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #ecf0f1;
}

/* 界面头部 */
.interface-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #ddd;
}

.interface-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.close-button {
  background: none;
  border: none;
  font-size: 32px;
  cursor: pointer;
  color: #7f8c8d;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.close-button:hover {
  background-color: rgba(0, 0, 0, 0.1);
  color: #2c3e50;
}

/* 界面操作区 */
.interface-actions {
  display: flex;
  gap: 10px;
}

/* 新增界面 */
.add-interface {
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* 表单容器 */
.form-container {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

/* 表单组 */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.form-group label {
  font-size: 14px;
  font-weight: 600;
  color: #2c3e50;
}

.form-group input,
.form-group textarea {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  outline: none;
  transition: all 0.3s ease;
}

.form-group input:focus,
.form-group textarea:focus {
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.1);
}

/* 标签选择器 */
.tag-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 5px;
}

.tag-option {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background-color: #f8f9fa;
  border: 1px solid #ddd;
  border-radius: 16px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.tag-option:hover {
  background-color: #e9ecef;
}

/* 集合选择器 */
.collection-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 5px;
}

.collection-option {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background-color: #f8f9fa;
  border: 1px solid #ddd;
  border-radius: 16px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.collection-option:hover {
  background-color: #e9ecef;
}

/* 系统选择器 */
.system-selector {
  display: flex;
  gap: 15px;
  margin-top: 5px;
}

.system-option {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  cursor: pointer;
}

/* 表单操作区 */
.form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 10px;
}

/* 按钮样式 */
.add-button {
  background-color: #27ae60;
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.add-button:hover {
  background-color: #219a52;
}

.save-button {
  background-color: #3498db;
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.save-button:hover {
  background-color: #2980b9;
}

.cancel-button {
  background-color: #7f8c8d;
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.cancel-button:hover {
  background-color: #6c757d;
}

/* 查询界面 */
.query-interface {
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* 表格容器 */
.table-container {
  margin-top: 20px;
  overflow-x: auto;
}

/* 命令表格 */
.command-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.command-table th,
.command-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

.command-table th {
  background-color: #34495e;
  color: white;
  font-weight: 600;
}

.command-table tr:hover {
  background-color: #f5f5f5;
}

/* 命令内容单元格 */
.command-content-cell {
  max-width: 200px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 标签徽章 */
.tag-badge {
  display: inline-block;
  padding: 2px 8px;
  background-color: #3498db;
  color: white;
  border-radius: 12px;
  font-size: 12px;
  margin-right: 4px;
}

/* 集合徽章 */
.collection-badge {
  display: inline-block;
  padding: 2px 8px;
  background-color: #27ae60;
  color: white;
  border-radius: 12px;
  font-size: 12px;
  margin-right: 4px;
}

/* 操作按钮组 */
.action-buttons {
  display: flex;
  gap: 8px;
}

.copy-button,
.edit-button,
.delete-button {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.copy-button {
  background-color: #3498db;
  color: white;
}

.copy-button:hover {
  background-color: #2980b9;
}

.edit-button {
  background-color: #f39c12;
  color: white;
}

.edit-button:hover {
  background-color: #e67e22;
}

.delete-button {
  background-color: #e74c3c;
  color: white;
}

.delete-button:hover {
  background-color: #c0392b;
}

/* 项目表格 */
.item-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.item-table th,
.item-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

.item-table th {
  background-color: #34495e;
  color: white;
  font-weight: 600;
}

.item-table tr:hover {
  background-color: #f5f5f5;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-content {
    padding: 10px;
  }
  
  .interface-header h2 {
    font-size: 20px;
  }
  
  .form-container {
    gap: 10px;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .tag-selector,
  .collection-selector {
    justify-content: center;
  }
  
  .system-selector {
    justify-content: center;
  }
}
</style>