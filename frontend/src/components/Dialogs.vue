<template>
  <!-- 设置对话框 -->
  <div v-if="isSettingsModalOpen" class="modal-overlay">
    <div class="modal-container">
      <div class="modal-header">
        <h2>设置</h2>
        <button class="close-button" @click="$emit('close-settings-modal')">×</button>
      </div>
      
      <div class="modal-content">
        <div class="form-group">
          <label for="theme">主题</label>
          <select id="theme" v-model="localSettings.theme">
            <option value="light">浅色</option>
            <option value="dark">深色</option>
          </select>
        </div>
        
        <div class="form-group">
          <label for="language">语言</label>
          <select id="language" v-model="localSettings.language">
            <option value="zh-CN">简体中文</option>
            <option value="en-US">English</option>
          </select>
        </div>
      </div>
      
      <div class="modal-footer">
        <button class="cancel-button" @click="$emit('close-settings-modal')">取消</button>
        <button class="save-button" @click="saveSettings">保存</button>
      </div>
    </div>
  </div>
  
  <!-- 关于对话框 -->
  <div v-if="isAboutModalOpen" class="modal-overlay">
    <div class="modal-container">
      <div class="modal-header">
        <h2>关于</h2>
        <button class="close-button" @click="$emit('close-about-modal')">×</button>
      </div>
      
      <div class="modal-content">
        <div class="about-content">
          <div class="app-info">
            <h3>QuickCmd</h3>
            <p>版本: 1.0.0</p>
            <p>快速命令管理工具</p>
            <p>作者: longan55</p>
          </div>
          
          <div class="developer-info">
            <h4>开发者</h4>
            <p>© 2024 QuickCmd Team</p>
            <p>保留所有权利</p>
          </div>
          
          <div class="links">
            <a href="#" class="link-button">GitHub</a>
            <a href="#" class="link-button">文档</a>
            <a href="#" class="link-button">反馈</a>
          </div>
        </div>
      </div>
      
      <div class="modal-footer">
        <button class="close-button-full" @click="$emit('close-about-modal')">关闭</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';

const props = defineProps({
  isSettingsModalOpen: {
    type: Boolean,
    required: true
  },
  isAboutModalOpen: {
    type: Boolean,
    required: true
  },
  settings: {
    type: Object,
    default: () => ({
      apiEndpoint: '',
      apiKey: '',
      theme: 'light',
      language: 'zh-CN'
    })
  }
});

const emit = defineEmits([
  'close-settings-modal',
  'close-about-modal',
  'update:settings'
]);

// 设置数据
const localSettings = ref({ ...props.settings });

// 当props.settings变化时，更新本地设置
watch(() => props.settings, (newSettings) => {
  localSettings.value = { ...newSettings };
}, { deep: true });

// 保存设置
function saveSettings() {
  emit('update:settings', { ...localSettings.value });
  emit('close-settings-modal');
}
</script>

<style scoped>
/* 模态框遮罩 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

/* 模态框容器 */
.modal-container {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
}

/* 模态框头部 */
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #ddd;
}

.modal-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.close-button {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #7f8c8d;
  padding: 0;
  width: 24px;
  height: 24px;
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

/* 模态框内容 */
.modal-content {
  padding: 20px;
}

/* 关于内容 */
.about-content {
  text-align: center;
}

.app-info {
  margin-bottom: 20px;
}

.app-info h3 {
  font-size: 24px;
  font-weight: 600;
  color: #3498db;
  margin-bottom: 10px;
}

.app-info p {
  margin: 5px 0;
  color: #2c3e50;
}

.developer-info {
  margin-bottom: 20px;
}

.developer-info h4 {
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 10px;
}

.developer-info p {
  margin: 5px 0;
  color: #7f8c8d;
}

/* 链接按钮 */
.links {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-top: 20px;
}

.link-button {
  display: inline-block;
  padding: 8px 16px;
  background-color: #f8f9fa;
  color: #3498db;
  border: 1px solid #3498db;
  border-radius: 4px;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.link-button:hover {
  background-color: #3498db;
  color: white;
}

/* 模态框底部 */
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 20px;
  border-top: 1px solid #ddd;
}

/* 全宽关闭按钮 */
.close-button-full {
  width: 100%;
  padding: 10px 16px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.close-button-full:hover {
  background-color: #2980b9;
}

/* 表单组 */
.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-size: 14px;
  font-weight: 600;
  color: #2c3e50;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  outline: none;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus {
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.1);
}

/* 按钮样式 */
.cancel-button {
  padding: 10px 16px;
  background-color: #7f8c8d;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.cancel-button:hover {
  background-color: #6c757d;
}

.save-button {
  padding: 10px 16px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.save-button:hover {
  background-color: #2980b9;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .modal-container {
    width: 95%;
    margin: 10px;
  }
  
  .modal-header,
  .modal-content,
  .modal-footer {
    padding: 15px;
  }
  
  .modal-header h2 {
    font-size: 18px;
  }
  
  .links {
    flex-direction: column;
    align-items: center;
  }
  
  .link-button {
    width: 100%;
    text-align: center;
  }
  
  .modal-footer {
    flex-direction: column;
  }
  
  .cancel-button,
  .save-button {
    width: 100%;
  }
}
</style>