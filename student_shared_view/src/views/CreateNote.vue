<template>
  <div class="create-note-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <h1 class="page-title">发布笔记</h1>
        <p class="page-subtitle">分享你的学习心得和知识总结</p>
      </div>
    </div>

    <div class="container">
      <el-row :gutter="20">
        <!-- 主要编辑区域 -->
        <el-col :span="18">
          <el-card class="editor-card">
            <el-form
              ref="noteFormRef"
              :model="noteForm"
              :rules="noteRules"
              label-width="80px"
            >
              <!-- 笔记标题 -->
              <el-form-item label="标题" prop="title">
                <el-input
                  v-model="noteForm.title"
                  placeholder="请输入笔记标题"
                  size="large"
                  maxlength="100"
                  show-word-limit
                />
              </el-form-item>

              <!-- 课程选择 -->
              <el-form-item label="课程" prop="course_id">
                <el-select
                  v-model="noteForm.course_id"
                  placeholder="请选择课程"
                  style="width: 100%"
                  :loading="coursesLoading"
                  clearable
                >
                  <el-option
                    label="无所属课程"
                    :value="0"
                  />
                  <el-option
                    v-for="course in courses"
                    :key="course.id"
                    :label="course.name"
                    :value="course.id"
                  />
                </el-select>
              </el-form-item>

              <!-- 笔记描述 -->
              <el-form-item label="描述">
                <el-input
                  v-model="noteForm.description"
                  type="textarea"
                  placeholder="请输入笔记描述（可选）"
                  :rows="3"
                  maxlength="500"
                  show-word-limit
                />
              </el-form-item>

              <!-- 内容编辑器 -->
              <el-form-item label="内容" prop="content">
                <div class="editor-container">
                  <!-- 工具栏 -->
                  <div class="editor-toolbar">
                    <el-button-group>
                      <el-button size="small" @click="insertText('**', '**')">
                        <el-icon><EditPen /></el-icon>
                      </el-button>
                      <el-button size="small" @click="insertText('*', '*')">
                        <el-icon><EditPen /></el-icon>
                      </el-button>
                      <el-button size="small" @click="insertText('`', '`')">
                        代码
                      </el-button>
                      <el-button size="small" @click="insertText('\n```\n', '\n```\n')">
                        代码块
                      </el-button>
                    </el-button-group>
                    
                    <el-button-group>
                      <el-button size="small" @click="insertText('# ', '')">
                        H1
                      </el-button>
                      <el-button size="small" @click="insertText('## ', '')">
                        H2
                      </el-button>
                      <el-button size="small" @click="insertText('### ', '')">
                        H3
                      </el-button>
                    </el-button-group>
                    
                    <el-button-group>
                      <el-button size="small" @click="insertText('- ', '')">
                        列表
                      </el-button>
                      <el-button size="small" @click="insertText('1. ', '')">
                        编号
                      </el-button>
                      <el-button size="small" @click="insertText('> ', '')">
                        引用
                      </el-button>
                    </el-button-group>
                    
                    <el-button size="small" @click="showImageUpload = true">
                      <el-icon><Picture /></el-icon>
                      插入图片
                    </el-button>
                    
                    <el-button size="small" @click="togglePreview">
                      <el-icon><View /></el-icon>
                      {{ showPreview ? '编辑' : '预览' }}
                    </el-button>
                  </div>
                  
                  <!-- 编辑区域 -->
                  <div class="editor-content">
                    <el-input
                      v-if="!showPreview"
                      ref="contentTextareaRef"
                      v-model="noteForm.content"
                      type="textarea"
                      :rows="20"
                      placeholder="请输入笔记内容，支持 Markdown 格式"
                      class="content-textarea"
                    />
                    
                    <!-- 预览区域 -->
                    <div v-else class="preview-content" v-html="previewContent"></div>
                  </div>
                </div>
              </el-form-item>

              <!-- 操作按钮 -->
              <el-form-item>
                <div class="action-buttons">
                  <el-button type="primary" @click="publishNote" :loading="publishLoading">
                    <el-icon><Upload /></el-icon>
                    发布笔记
                  </el-button>
                  <el-button @click="$router.go(-1)">
                    取消
                  </el-button>
                </div>
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>

        <!-- 侧边栏 -->
        <el-col :span="6">
          <div class="sidebar">
            <!-- 发布设置 -->
            <el-card class="settings-card">
              <template #header>
                <span>发布设置</span>
              </template>
              
              <el-form label-width="60px">
                <el-form-item label="可见性">
                  <el-radio-group v-model="noteForm.visibility">
                    <el-radio value="public">公开</el-radio>
                    <el-radio value="private">私有</el-radio>
                  </el-radio-group>
                </el-form-item>
                
                <el-form-item label="评论">
                  <el-switch v-model="noteForm.allow_comments" />
                </el-form-item>
                
                <el-form-item label="难度">
                  <el-select v-model="noteForm.difficulty" placeholder="选择难度">
                    <el-option label="入门" value="beginner" />
                    <el-option label="中级" value="intermediate" />
                    <el-option label="高级" value="advanced" />
                  </el-select>
                </el-form-item>
              </el-form>
            </el-card>

            <!-- 写作提示 -->
            <el-card class="tips-card">
              <template #header>
                <span>写作提示</span>
              </template>
              
              <div class="tips-content">
                <div class="tip-item">
                  <h4>📝 标题建议</h4>
                  <p>使用清晰、具体的标题，包含关键词</p>
                </div>
                
                <div class="tip-item">
                  <h4>🏷️ 标签使用</h4>
                  <p>添加相关标签，帮助其他同学找到你的笔记</p>
                </div>
                
                <div class="tip-item">
                  <h4>📖 内容结构</h4>
                  <p>使用标题、列表等格式，让内容更易读</p>
                </div>
                
                <div class="tip-item">
                  <h4>💡 分享经验</h4>
                  <p>分享学习心得和实际应用案例</p>
                </div>
              </div>
            </el-card>

            <!-- Markdown 语法帮助 -->
            <el-card class="markdown-help">
              <template #header>
                <span>Markdown 语法</span>
              </template>
              
              <div class="markdown-syntax">
                <div class="syntax-item">
                  <code># 标题</code>
                  <span>一级标题</span>
                </div>
                <div class="syntax-item">
                  <code>**粗体**</code>
                  <span>粗体文字</span>
                </div>
                <div class="syntax-item">
                  <code>*斜体*</code>
                  <span>斜体文字</span>
                </div>
                <div class="syntax-item">
                  <code>`代码`</code>
                  <span>行内代码</span>
                </div>
                <div class="syntax-item">
                  <code>- 列表</code>
                  <span>无序列表</span>
                </div>
                <div class="syntax-item">
                  <code>> 引用</code>
                  <span>引用文字</span>
                </div>
              </div>
            </el-card>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 图片上传对话框 -->
    <el-dialog v-model="showImageUpload" title="插入图片" width="400px">
      <el-upload
        class="image-uploader"
        action="#"
        :show-file-list="false"
        :before-upload="beforeImageUpload"
        :http-request="handleImageUpload"
      >
        <el-button type="primary">
          <el-icon><Upload /></el-icon>
          选择图片
        </el-button>
        <template #tip>
          <div class="el-upload__tip">
            支持 jpg/png 文件，且不超过 2MB
          </div>
        </template>
      </el-upload>
      
      <div class="image-url-input">
        <el-divider>或</el-divider>
        <el-input
          v-model="imageUrl"
          placeholder="输入图片链接"
          @keyup.enter="insertImage"
        >
          <template #append>
            <el-button @click="insertImage">插入</el-button>
          </template>
        </el-input>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, computed, nextTick, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { noteAPI, courseAPI } from '../api'
import {
  EditPen,
  Picture,
  View,
  Document,
  Upload
} from '@element-plus/icons-vue'

export default {
  name: 'CreateNote',
  components: {
    EditPen,
    Picture,
    View,
    Document,
    Upload
  },
  setup() {
    const router = useRouter()
    const noteFormRef = ref()
    const tagInputRef = ref()
    const contentTextareaRef = ref()
    
    const publishLoading = ref(false)
    const showPreview = ref(false)
    const showImageUpload = ref(false)

    const imageUrl = ref('')
    

    
    // 笔记表单
    const noteForm = reactive({
      title: '',
      content: '',
      description: '',
      course_id: 0,

      visibility: 'public',
      allow_comments: true,
      difficulty: 'intermediate'
    })
    
    // 课程列表
    const courses = ref([])
    const coursesLoading = ref(false)
    
    // 表单验证规则
    const noteRules = {
      title: [
        { required: true, message: '请输入笔记标题', trigger: 'blur' },
        { min: 5, max: 100, message: '标题长度在 5 到 100 个字符', trigger: 'blur' }
      ],
      content: [
        { required: true, message: '请输入笔记内容', trigger: 'blur' },
        { min: 50, message: '内容至少需要 50 个字符', trigger: 'blur' }
      ]
    }
    
    // 预览内容（简单的 Markdown 转换）
    const previewContent = computed(() => {
      let content = noteForm.content
      
      // 简单的 Markdown 转换
      content = content
        .replace(/^# (.*$)/gim, '<h1>$1</h1>')
        .replace(/^## (.*$)/gim, '<h2>$1</h2>')
        .replace(/^### (.*$)/gim, '<h3>$1</h3>')
        .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
        .replace(/\*(.*?)\*/g, '<em>$1</em>')
        .replace(/`(.*?)`/g, '<code>$1</code>')
        .replace(/^- (.*$)/gim, '<li>$1</li>')
        .replace(/^> (.*$)/gim, '<blockquote>$1</blockquote>')
        .replace(/\n/g, '<br>')
      
      return content
    })
    
    // 插入文本
    const insertText = (before, after) => {
      const textarea = contentTextareaRef.value?.textarea
      if (!textarea) return
      
      const start = textarea.selectionStart
      const end = textarea.selectionEnd
      const selectedText = noteForm.content.substring(start, end)
      
      const newText = before + selectedText + after
      noteForm.content = noteForm.content.substring(0, start) + newText + noteForm.content.substring(end)
      
      nextTick(() => {
        textarea.focus()
        textarea.setSelectionRange(start + before.length, start + before.length + selectedText.length)
      })
    }
    
    // 切换预览
    const togglePreview = () => {
      showPreview.value = !showPreview.value
    }
    

    
    // 图片上传相关
    const beforeImageUpload = (file) => {
      const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
      const isLt2M = file.size / 1024 / 1024 < 2
      
      if (!isJPG) {
        ElMessage.error('只能上传 JPG/PNG 格式的图片!')
      }
      if (!isLt2M) {
        ElMessage.error('图片大小不能超过 2MB!')
      }
      return isJPG && isLt2M
    }
    
    const handleImageUpload = (options) => {
      // 这里应该实现真实的图片上传逻辑
      // 模拟上传成功
      setTimeout(() => {
        const imageUrl = 'https://via.placeholder.com/400x300'
        insertImage(imageUrl)
        ElMessage.success('图片上传成功')
      }, 1000)
    }
    
    const insertImage = (url = imageUrl.value) => {
      if (!url) {
        ElMessage.warning('请输入图片链接')
        return
      }
      
      const imageMarkdown = `\n![图片描述](${url})\n`
      noteForm.content += imageMarkdown
      
      showImageUpload.value = false
      imageUrl.value = ''
    }
    
    // 发布笔记
    const publishNote = async () => {
      if (!noteFormRef.value) return
      
      try {
        await noteFormRef.value.validate()
        
        await ElMessageBox.confirm('确定要发布这篇笔记吗？', '确认发布', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'info'
        })
        
        publishLoading.value = true
        
        // 创建新笔记
        const noteData = {
          title: noteForm.title,
          content: noteForm.content,
          description: noteForm.description,
          course_id: noteForm.course_id,

          visibility: noteForm.visibility,
          allow_comments: noteForm.allow_comments,
          difficulty: noteForm.difficulty
        }
        
        const response = await noteAPI.createNote(noteData)
        ElMessage.success('笔记发布成功！')
        
        // 跳转到笔记详情页
        router.push(`/note/${response.id}`)
      } catch (error) {
        if (error !== 'cancel') {
          console.error('发布失败:', error)
          ElMessage.error('发布失败，请稍后重试')
        }
      } finally {
        publishLoading.value = false
      }
    }
    
    // 获取课程列表
    const fetchCourses = async () => {
      try {
        coursesLoading.value = true
        const response = await courseAPI.getCourses()
        courses.value = response.courses || response || []
      } catch (error) {
        console.error('获取课程列表失败:', error)
        ElMessage.error('获取课程列表失败')
      } finally {
        coursesLoading.value = false
      }
    }
    
    onMounted(() => {
      fetchCourses()
    })
    
    return {
      noteFormRef,
      tagInputRef,
      contentTextareaRef,
      publishLoading,
      showPreview,
      showImageUpload,

      imageUrl,
      courses,
      coursesLoading,

      noteForm,
      noteRules,
      previewContent,
      insertText,
      togglePreview,

      beforeImageUpload,
      handleImageUpload,
      insertImage,
      publishNote,
      fetchCourses
    }
  }
}
</script>

<style scoped>
.editor-card {
  margin-bottom: 20px;
}

.tags-input {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.tag-item {
  margin: 0;
}

.tag-input {
  width: 120px;
}

.button-new-tag {
  border-style: dashed;
}

.editor-container {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
}

.editor-toolbar {
  display: flex;
  gap: 10px;
  padding: 10px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  flex-wrap: wrap;
}

.editor-content {
  min-height: 500px;
}

.content-textarea {
  border: none;
  border-radius: 0;
}

.content-textarea :deep(.el-textarea__inner) {
  border: none;
  border-radius: 0;
  resize: vertical;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  line-height: 1.6;
}

.preview-content {
  padding: 15px;
  min-height: 465px;
  line-height: 1.6;
  background: white;
}

.preview-content :deep(h1),
.preview-content :deep(h2),
.preview-content :deep(h3) {
  margin: 20px 0 10px 0;
  color: var(--text-primary);
}

.preview-content :deep(h1) {
  font-size: 1.8rem;
  border-bottom: 2px solid var(--primary-color);
  padding-bottom: 10px;
}

.preview-content :deep(h2) {
  font-size: 1.5rem;
}

.preview-content :deep(h3) {
  font-size: 1.3rem;
}

.preview-content :deep(code) {
  background: var(--bg-secondary);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.preview-content :deep(blockquote) {
  border-left: 4px solid var(--primary-color);
  padding-left: 15px;
  margin: 15px 0;
  color: var(--text-secondary);
  background: var(--bg-secondary);
  padding: 10px 15px;
  border-radius: 0 4px 4px 0;
}

.preview-content :deep(li) {
  margin: 5px 0;
  list-style: disc;
  margin-left: 20px;
}

.action-buttons {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.sidebar {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.settings-card,
.tips-card,
.markdown-help {
  position: sticky;
  top: 20px;
}

.tips-content {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.tip-item h4 {
  margin: 0 0 5px 0;
  font-size: 0.9rem;
  color: var(--text-primary);
}

.tip-item p {
  margin: 0;
  font-size: 0.8rem;
  color: var(--text-secondary);
  line-height: 1.4;
}

.markdown-syntax {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.syntax-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
}

.syntax-item code {
  background: var(--bg-secondary);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.syntax-item span {
  color: var(--text-secondary);
}

.image-uploader {
  text-align: center;
}

.image-url-input {
  margin-top: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .container .el-row {
    flex-direction: column;
  }
  
  .sidebar {
    margin-top: 20px;
  }
  
  .editor-toolbar {
    gap: 5px;
  }
  
  .editor-toolbar .el-button-group {
    margin-bottom: 5px;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .settings-card,
  .tips-card,
  .markdown-help {
    position: static;
  }
}
</style>