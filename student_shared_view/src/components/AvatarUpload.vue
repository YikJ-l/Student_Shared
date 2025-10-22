<template>
  <div class="avatar-upload">
    <!-- 头像上传对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="更换头像"
      width="600px"
      :before-close="handleClose"
      center
      :modal="true"
      :close-on-click-modal="false"
      :close-on-press-escape="true"
      append-to-body
    >
      <el-tabs v-model="activeTab" class="avatar-tabs">
        <!-- 本地上传 -->
        <el-tab-pane label="本地上传" name="upload">
          <div class="upload-section">
            <el-upload
              ref="uploadRef"
              class="avatar-uploader"
              :action="uploadUrl"
              :headers="uploadHeaders"
              :show-file-list="false"
              :on-success="handleUploadSuccess"
              :on-error="handleUploadError"
              :before-upload="beforeUpload"
              :on-progress="handleProgress"
              accept="image/*"
            >
              <div class="upload-area">
                <div v-if="uploading" class="uploading">
                  <el-progress
                    type="circle"
                    :percentage="uploadProgress"
                    :width="80"
                  />
                  <p>上传中...</p>
                </div>
                <div v-else class="upload-placeholder">
                  <el-icon class="upload-icon"><Plus /></el-icon>
                  <p>点击上传头像</p>
                  <p class="upload-tip">支持 JPG、PNG 格式，文件大小不超过 2MB</p>
                </div>
              </div>
            </el-upload>
            
            <!-- 预览区域 -->
            <div v-if="previewUrl" class="preview-section">
              <h4>预览 <span class="preview-tip">(点击头像可放大查看)</span></h4>
              <div class="preview-avatars">
                <el-avatar 
                  :size="80" 
                  :src="previewUrl" 
                  class="preview-avatar-clickable"
                  @click="showImagePreview(previewUrl)"
                />
                <el-avatar 
                  :size="60" 
                  :src="previewUrl" 
                  class="preview-avatar-clickable"
                  @click="showImagePreview(previewUrl)"
                />
                <el-avatar 
                  :size="40" 
                  :src="previewUrl" 
                  class="preview-avatar-clickable"
                  @click="showImagePreview(previewUrl)"
                />
              </div>
            </div>
          </div>
        </el-tab-pane>
        
        <!-- URL链接 -->
        <el-tab-pane label="网络图片" name="url">
          <div class="url-section">
            <el-form :model="urlForm" :rules="urlRules" ref="urlFormRef">
              <el-form-item label="图片链接" prop="url">
                <el-input
                  v-model="urlForm.url"
                  placeholder="请输入图片链接地址"
                  @input="handleUrlChange"
                />
              </el-form-item>
            </el-form>
            
            <!-- URL预览 -->
            <div v-if="urlForm.url && isValidUrl" class="preview-section">
              <h4>预览 <span class="preview-tip">(点击头像可放大查看)</span></h4>
              <div class="preview-avatars">
                <el-avatar 
                  :size="80" 
                  :src="urlForm.url" 
                  @error="handleImageError" 
                  class="preview-avatar-clickable"
                  @click="showImagePreview(urlForm.url)"
                />
                <el-avatar 
                  :size="60" 
                  :src="urlForm.url" 
                  @error="handleImageError" 
                  class="preview-avatar-clickable"
                  @click="showImagePreview(urlForm.url)"
                />
                <el-avatar 
                  :size="40" 
                  :src="urlForm.url" 
                  @error="handleImageError" 
                  class="preview-avatar-clickable"
                  @click="showImagePreview(urlForm.url)"
                />
              </div>
            </div>
            
            <div v-if="urlForm.url && !isValidUrl" class="error-tip">
              <el-alert
                title="图片链接无效"
                description="请检查图片链接是否正确，确保图片可以正常访问"
                type="error"
                :closable="false"
              />
            </div>
          </div>
        </el-tab-pane>
        
        <!-- 预设头像 -->
        <el-tab-pane label="预设头像" name="preset">
          <div class="preset-section">
            <div class="preset-avatars">
              <div
                v-for="(avatar, index) in presetAvatars"
                :key="index"
                class="preset-avatar-item"
                :class="{ active: selectedPreset === avatar }"
                @click="selectPreset(avatar)"
              >
                <el-avatar 
                  :size="60" 
                  :src="avatar" 
                  class="preset-avatar-clickable"
                  @click.stop="showImagePreview(avatar)"
                />
                <div class="preset-avatar-overlay" @click="selectPreset(avatar)">
                  <span v-if="selectedPreset === avatar" class="selected-icon">✓</span>
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="handleClose">取消</el-button>
          <el-button
            type="primary"
            :loading="saving"
            :disabled="!canSave"
            @click="handleSave"
          >
            保存
          </el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 图片预览对话框 -->
    <el-dialog
      v-model="imagePreviewVisible"
      title="头像预览"
      width="80%"
      :modal="true"
      :close-on-click-modal="true"
      :close-on-press-escape="true"
      append-to-body
      center
    >
      <div class="image-preview-container">
        <el-image
          :src="previewImageUrl"
          fit="contain"
          class="preview-image"
          :preview-src-list="[previewImageUrl]"
          :initial-index="0"
          hide-on-click-modal
        />
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import api from '@/api'

export default {
  name: 'AvatarUpload',
  components: {
    Plus
  },
  props: {
    modelValue: {
      type: Boolean,
      default: false
    },
    currentAvatar: {
      type: String,
      default: ''
    }
  },
  emits: ['update:modelValue', 'success'],
  setup(props, { emit }) {
    const dialogVisible = ref(false)
    const activeTab = ref('upload')
    const uploading = ref(false)
    const uploadProgress = ref(0)
    const saving = ref(false)
    const previewUrl = ref('')
    const selectedPreset = ref('')
    const isValidUrl = ref(true)
    const uploadRef = ref()
    const urlFormRef = ref()
    
    // 图片预览相关
    const imagePreviewVisible = ref(false)
    const previewImageUrl = ref('')
    
    // URL表单
    const urlForm = reactive({
      url: ''
    })
    
    // URL验证规则
    const urlRules = {
      url: [
        { required: true, message: '请输入图片链接', trigger: 'blur' },
        { type: 'url', message: '请输入有效的URL地址', trigger: 'blur' }
      ]
    }
    
    // 预设头像列表
    const presetAvatars = [
      'http://localhost:8081/img/avatars/avatar1.svg',
      'http://localhost:8081/img/avatars/avatar2.svg',
      'http://localhost:8081/img/avatars/avatar3.svg',
      'http://localhost:8081/img/avatars/avatar4.svg',
      'http://localhost:8081/img/avatars/avatar5.svg',
      'http://localhost:8081/img/avatars/avatar6.svg'
    ]
    
    // 从全局API实例获取服务端origin
    const apiBase = new URL(api.defaults.baseURL)
    const serverOrigin = `${apiBase.protocol}//${apiBase.host}`

    // 上传配置（统一使用全局API baseURL）
    const uploadUrl = computed(() => {
      return `${api.defaults.baseURL}/upload/avatar`
    })
    
    const uploadHeaders = computed(() => {
      const token = localStorage.getItem('token')
      return {
        token
      }
    })
    
    // 是否可以保存
    const canSave = computed(() => {
      if (activeTab.value === 'upload') {
        return previewUrl.value !== ''
      } else if (activeTab.value === 'url') {
        return urlForm.url !== '' && isValidUrl.value
      } else if (activeTab.value === 'preset') {
        return selectedPreset.value !== ''
      }
      return false
    })
    
    // 监听对话框显示状态
    watch(() => props.modelValue, (val) => {
      dialogVisible.value = val
      if (val) {
        resetForm()
      }
    })
    
    watch(dialogVisible, (val) => {
      emit('update:modelValue', val)
    })
    
    // 重置表单
    const resetForm = () => {
      activeTab.value = 'upload'
      uploading.value = false
      uploadProgress.value = 0
      previewUrl.value = ''
      selectedPreset.value = ''
      urlForm.url = ''
      isValidUrl.value = true
    }
    
    // 上传前检查
    const beforeUpload = (file) => {
      const isImage = file.type.startsWith('image/')
      const isLt2M = file.size / 1024 / 1024 < 2
      
      if (!isImage) {
        ElMessage.error('只能上传图片文件!')
        return false
      }
      if (!isLt2M) {
        ElMessage.error('图片大小不能超过 2MB!')
        return false
      }
      
      // 创建预览
      const reader = new FileReader()
      reader.onload = (e) => {
        previewUrl.value = e.target.result
      }
      reader.readAsDataURL(file)
      
      return true
    }
    
    // 上传进度
    const handleProgress = (event) => {
      uploading.value = true
      uploadProgress.value = Math.round(event.percent)
    }
    
    // 上传成功（后端返回 { url: '/uploads/avatars/xxx', size: n }）
    const handleUploadSuccess = (response) => {
      uploading.value = false
      uploadProgress.value = 0
      
      if (response && response.url) {
        previewUrl.value = `${serverOrigin}${response.url}`
        ElMessage.success('上传成功')
      } else {
        ElMessage.error('上传失败')
      }
    }
    
    // 上传失败
    const handleUploadError = (error) => {
      uploading.value = false
      uploadProgress.value = 0
      ElMessage.error('上传失败，请重试')
      console.error('Upload error:', error)
    }
    
    // URL变化处理
    const handleUrlChange = () => {
      isValidUrl.value = true
    }
    
    // 图片加载错误
    const handleImageError = () => {
      isValidUrl.value = false
    }
    
    // 选择预设头像
    const selectPreset = (avatar) => {
      selectedPreset.value = avatar
    }
    
    // 保存头像
    const handleSave = async () => {
      let avatarUrl = ''
      
      if (activeTab.value === 'upload') {
        avatarUrl = previewUrl.value
      } else if (activeTab.value === 'url') {
        // 验证URL表单
        try {
          await urlFormRef.value.validate()
          avatarUrl = urlForm.url
        } catch (error) {
          return
        }
      } else if (activeTab.value === 'preset') {
        avatarUrl = selectedPreset.value
      }
      
      if (!avatarUrl) {
        ElMessage.error('请选择头像')
        return
      }
      
      saving.value = true
      
      try {
        // 触发保存事件
        emit('success', avatarUrl)
        ElMessage.success('头像更新成功')
        handleClose()
      } catch (error) {
        ElMessage.error('保存失败，请重试')
      } finally {
        saving.value = false
      }
    }
    
    // 关闭对话框
    const handleClose = () => {
      dialogVisible.value = false
    }
    
    // 显示图片预览
    const showImagePreview = (imageUrl) => {
      previewImageUrl.value = imageUrl
      imagePreviewVisible.value = true
    }
    
    return {
      dialogVisible,
      activeTab,
      uploading,
      uploadProgress,
      saving,
      previewUrl,
      selectedPreset,
      isValidUrl,
      uploadRef,
      urlFormRef,
      imagePreviewVisible,
      previewImageUrl,
      urlForm,
      urlRules,
      presetAvatars,
      uploadUrl,
      uploadHeaders,
      canSave,
      beforeUpload,
      handleProgress,
      handleUploadSuccess,
      handleUploadError,
      handleUrlChange,
      handleImageError,
      selectPreset,
      handleSave,
      handleClose,
      showImagePreview
    }
  }
}
</script>

<style scoped>
.avatar-upload {
  display: inline-block;
}

.avatar-tabs {
  margin-top: 20px;
}

.el-tabs__content {
  padding: 20px 0;
}

.upload-section {
  text-align: center;
}

.avatar-uploader {
  margin-bottom: 20px;
}

.upload-area {
  border: 2px dashed #d9d9d9;
  border-radius: 6px;
  width: 100%;
  max-width: 400px;
  height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.3s;
  margin: 0 auto;
}

.upload-area:hover {
  border-color: #409eff;
}

.upload-icon {
  font-size: 28px;
  color: #8c939d;
  margin-bottom: 16px;
}

.upload-tip {
  color: #8c939d;
  font-size: 12px;
  margin-top: 8px;
}

.uploading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.preview-section {
  margin-top: 20px;
}

.preview-section h4 {
  margin-bottom: 16px;
  color: #303133;
}

.preview-tip {
  font-size: 12px;
  color: #909399;
  font-weight: normal;
}

.preview-avatars {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

.preview-avatar-clickable {
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.preview-avatar-clickable:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

.url-section {
  padding: 0 20px;
}

.error-tip {
  margin-top: 16px;
}

.preset-section {
  padding: 0 20px;
}

.preset-avatars {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-top: 16px;
  padding: 0 20px;
}

.preset-avatar-item {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 8px;
  border: 2px solid transparent;
  border-radius: 8px;
  cursor: pointer;
  transition: border-color 0.2s ease, background-color 0.2s ease;
  position: relative;
}

.preset-avatar-clickable {
  cursor: pointer;
  transition: transform 0.2s ease;
}

.preset-avatar-clickable:hover {
  transform: scale(1.1);
}

.preset-avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  pointer-events: none;
}

.selected-icon {
  position: absolute;
  top: 4px;
  right: 4px;
  background: #409eff;
  color: white;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
}

.preset-avatar-item:hover {
  border-color: #409eff;
  background-color: #f0f9ff;
}

.preset-avatar-item.active {
  border-color: #409eff;
  background-color: #e1f3ff;
}

.preset-avatar-item .el-avatar {
  pointer-events: none;
  transition: none;
}

.dialog-footer {
  text-align: right;
}

/* 确保对话框正确定位 */
.avatar-upload .el-dialog {
  margin: 0 auto;
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.avatar-upload .el-dialog__wrapper {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 2000;
}

/* 图片预览样式 */
.image-preview-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.preview-image {
  max-width: 100%;
  max-height: 70vh;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}
</style>