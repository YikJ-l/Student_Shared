<template>
  <div class="edit-note-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <el-breadcrumb separator=">">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/my-notes' }">我的笔记</el-breadcrumb-item>
          <el-breadcrumb-item>编辑笔记</el-breadcrumb-item>
        </el-breadcrumb>
        <h1 class="page-title">编辑笔记</h1>
      </div>
    </div>

    <div class="container">
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="8" animated />
      </div>

      <div v-else class="edit-note-content">
        <el-form
          ref="noteFormRef"
          :model="noteForm"
          :rules="noteRules"
          label-width="100px"
          class="note-form"
        >
          <!-- 基本信息 -->
          <el-card class="form-section">
            <template #header>
              <span>基本信息</span>
            </template>
            
            <el-form-item label="笔记标题" prop="title">
              <el-input
                v-model="noteForm.title"
                placeholder="请输入笔记标题"
                maxlength="200"
                show-word-limit
              />
            </el-form-item>
            
            <el-form-item label="笔记描述" prop="description">
              <el-input
                v-model="noteForm.description"
                type="textarea"
                :rows="3"
                placeholder="请输入笔记描述"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>
            
            <el-form-item label="所属课程" prop="course_id">
              <el-select
                v-model="noteForm.course_id"
                placeholder="请选择课程"
                style="width: 100%"
                filterable
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
            
            <el-form-item label="可见性" prop="status">
              <el-radio-group v-model="noteForm.status">
                <el-radio label="public">公开</el-radio>
                <el-radio label="private">私有</el-radio>
              </el-radio-group>
            </el-form-item>
            

          </el-card>

          <!-- 笔记内容 -->
          <el-card class="form-section">
            <template #header>
              <div class="content-header">
                <span>笔记内容</span>
                <div class="editor-actions">
                  <el-button
                    :type="showPreview ? 'default' : 'primary'"
                    size="small"
                    @click="showPreview = false"
                  >
                    编辑
                  </el-button>
                  <el-button
                    :type="showPreview ? 'primary' : 'default'"
                    size="small"
                    @click="showPreview = true"
                  >
                    预览
                  </el-button>
                </div>
              </div>
            </template>
            
            <div class="editor-container">
              <!-- 编辑模式 -->
              <div v-if="!showPreview" class="editor-content">
                <el-input
                  v-model="noteForm.content"
                  type="textarea"
                  :rows="20"
                  placeholder="请输入笔记内容，支持Markdown格式"
                  class="content-editor"
                />
              </div>
              
              <!-- 预览模式 -->
              <div v-else class="preview-content">
                <div v-if="noteForm.content" v-html="renderedContent" class="markdown-content"></div>
                <div v-else class="empty-preview">
                  <el-empty description="暂无内容" />
                </div>
              </div>
            </div>
          </el-card>

          <!-- 操作按钮 -->
          <div class="form-actions">
            <el-button size="large" @click="handleCancel">取消</el-button>
            <el-button type="primary" size="large" :loading="saving" @click="handleSave">
              保存修改
            </el-button>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { noteAPI, courseAPI } from '../api'

export default {
  name: 'EditNote',
  setup() {
    const route = useRoute()
    const router = useRouter()
    
    const loading = ref(true)
    const saving = ref(false)
    const showPreview = ref(false)
    const noteFormRef = ref()
    const courses = ref([])
    const noteAuthorId = ref(null)
    
    // 表单数据
    const noteForm = reactive({
      title: '',
      description: '',
      content: '',
      course_id: null,
      status: 'public',

    })
    
    // 表单验证规则
    const noteRules = {
      title: [
        { required: true, message: '请输入笔记标题', trigger: 'blur' },
        { min: 1, max: 200, message: '标题长度在 1 到 200 个字符', trigger: 'blur' }
      ],
      description: [
        { max: 500, message: '描述长度不能超过 500 个字符', trigger: 'blur' }
      ],
      course_id: [
        { required: true, message: '请选择所属课程', trigger: 'change' }
      ]
    }
    
    // 渲染的内容（简单的换行处理）
    const renderedContent = computed(() => {
      if (!noteForm.content) return ''
      return noteForm.content.replace(/\n/g, '<br>')
    })
    
    // 获取笔记详情
    const fetchNoteDetail = async () => {
      try {
        const noteId = route.params.id
        const response = await noteAPI.getNoteDetail(noteId)
        
        // 保存作者ID用于权限验证
        noteAuthorId.value = response.author_id || response.user_id
        
        // 填充表单数据
        Object.assign(noteForm, {
          title: response.title || '',
          description: response.description || '',
          content: response.content || '',
          course_id: response.course_id,
          status: response.status || 'public',

        })
        
      } catch (error) {
        console.error('获取笔记详情失败:', error)
        ElMessage.error('获取笔记详情失败')
        router.push('/my-notes')
      }
    }
    
    // 获取课程列表
    const fetchCourses = async () => {
      try {
        const response = await courseAPI.getCourses({ page_size: 100 })
        courses.value = response.courses || []
      } catch (error) {
        console.error('获取课程列表失败:', error)
        ElMessage.error('获取课程列表失败')
      }
    }
    
    // 保存笔记
    const handleSave = async () => {
      try {
        // 表单验证
        const valid = await noteFormRef.value.validate()
        if (!valid) return
        
        saving.value = true
        const noteId = route.params.id
        
        await noteAPI.updateNote(noteId, {
          title: noteForm.title,
          description: noteForm.description,
          content: noteForm.content,
          course_id: noteForm.course_id,
          status: noteForm.status,

        }, noteAuthorId.value)
        
        ElMessage.success('笔记更新成功')
        router.push('/my-notes')
        
      } catch (error) {
        console.error('保存笔记失败:', error)
        ElMessage.error('保存笔记失败')
      } finally {
        saving.value = false
      }
    }
    
    // 取消编辑
    const handleCancel = async () => {
      try {
        await ElMessageBox.confirm('确定要取消编辑吗？未保存的修改将丢失。', '确认取消', {
          confirmButtonText: '确定',
          cancelButtonText: '继续编辑',
          type: 'warning'
        })
        
        const noteId = route.params.id
        router.push(`/notes/${noteId}`)
      } catch {
        // 用户取消
      }
    }
    
    // 初始化
    onMounted(async () => {
      try {
        await Promise.all([
          fetchNoteDetail(),
          fetchCourses()
        ])
      } finally {
        loading.value = false
      }
    })
    
    return {
      loading,
      saving,
      showPreview,
      noteFormRef,
      noteForm,
      noteRules,
      courses,
      renderedContent,
      handleSave,
      handleCancel
    }
  }
}
</script>

<style scoped>
.edit-note-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.page-header {
  background: white;
  border-bottom: 1px solid var(--border-color);
  padding: 20px 0;
}

.page-title {
  margin: 10px 0 0 0;
  font-size: 24px;
  font-weight: 600;
  color: var(--text-color-primary);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.loading-container {
  padding: 40px 0;
}

.edit-note-content {
  padding: 30px 0;
}

.note-form {
  max-width: 800px;
  margin: 0 auto;
}

.form-section {
  margin-bottom: 20px;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.editor-actions {
  display: flex;
  gap: 10px;
}

.editor-container {
  min-height: 400px;
}

.content-editor {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.6;
}

.preview-content {
  min-height: 400px;
  padding: 20px;
  background-color: #fafafa;
  border: 1px solid var(--border-color);
  border-radius: 4px;
}

.markdown-content {
  line-height: 1.8;
  color: var(--text-color-primary);
}

.empty-preview {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 300px;
}

.form-actions {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .note-form {
    max-width: 100%;
  }
  
  .content-header {
    flex-direction: column;
    gap: 10px;
    align-items: flex-start;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .form-actions .el-button {
    width: 100%;
  }
}
</style>