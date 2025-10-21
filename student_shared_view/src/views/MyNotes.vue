<template>
  <div class="my-notes-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <h1 class="page-title">我的笔记</h1>
        <p class="page-subtitle">管理您的学习笔记</p>
      </div>
    </div>

    <div class="container">
      <!-- 操作栏 -->
      <div class="action-bar">
        <div class="left-actions">
          <el-button type="primary" @click="$router.push('/create-note')">
            <el-icon><Plus /></el-icon>
            新建笔记
          </el-button>
          <el-button v-if="selectedNotes.length > 0" type="danger" @click="handleBatchDelete">
            <el-icon><Delete /></el-icon>
            批量删除 ({{ selectedNotes.length }})
          </el-button>
        </div>
        <div class="right-actions">
          <div class="search-container">
            <el-input
              v-model="searchQuery"
              placeholder="搜索笔记标题或内容"
              style="width: 300px"
              clearable
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
              <template #append>
                <el-button type="primary" @click="handleSearch">
                  搜索
                </el-button>
              </template>
            </el-input>
            
            <!-- 搜索历史下拉菜单 -->
            <el-dropdown v-if="searchHistory.length > 0" trigger="click" class="search-history-dropdown">
              <el-button size="small" text>
                <el-icon><Clock /></el-icon>
                历史
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item 
                    v-for="(item, index) in searchHistory" 
                    :key="index"
                    @click="selectSearchHistory(item)"
                  >
                    {{ item }}
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="clearSearchHistory">
                    <span style="color: #f56c6c;">清除历史</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          
          <el-button style="margin-left: 10px" @click="handleReset">
            重置
          </el-button>
        </div>
      </div>

      <!-- 统计信息 -->
      <div class="stats-section">
        <el-row :gutter="20">
          <el-col :span="6">
            <div class="stat-card">
              <div class="stat-icon">
                <el-icon><Document /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-value">{{ stats.total }}</div>
                <div class="stat-title">总笔记数</div>
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-card">
              <div class="stat-icon">
                <el-icon><View /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-value">{{ stats.totalViews }}</div>
                <div class="stat-title">总浏览量</div>
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-card">
              <div class="stat-icon">
                <el-icon><Star /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-value">{{ stats.totalLikes }}</div>
                <div class="stat-title">总点赞数</div>
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-card">
              <div class="stat-icon">
                <el-icon><ChatDotRound /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-value">{{ stats.totalComments }}</div>
                <div class="stat-title">总评论数</div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- 搜索结果提示 -->
      <div v-if="searchQuery" class="search-result-info">
        <el-alert
          :title="`搜索 '${searchQuery}' 找到 ${filteredNotes.length} 条结果`"
          type="info"
          :closable="false"
          show-icon
        />
      </div>

      <!-- 笔记列表 -->
      <div class="notes-section">
        <div v-if="loading" class="loading-container">
          <el-skeleton :rows="4" animated />
          <el-skeleton :rows="4" animated />
          <el-skeleton :rows="4" animated />
        </div>

        <div v-else-if="filteredNotes.length === 0" class="empty-container">
          <el-empty description="暂无笔记">
            <el-button type="primary" @click="$router.push('/create-note')">
              发布第一篇笔记
            </el-button>
          </el-empty>
        </div>

        <div v-else class="notes-table">
          <el-table
            :data="filteredNotes"
            style="width: 100%"
            @selection-change="handleSelectionChange"
          >
            <el-table-column type="selection" width="55" />
            <el-table-column prop="title" label="标题" min-width="200">
              <template #default="scope">
                <div class="note-title-cell">
                  <el-link 
                    type="primary" 
                    @click="$router.push(`/notes/${scope.row.id}`)"
                    :underline="false"
                  >
                    {{ scope.row.title }}
                  </el-link>
                  <div class="note-description">{{ scope.row.description || '暂无描述' }}</div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="course_name" label="课程" width="120">
              <template #default="scope">
                <el-tag size="small" type="success">{{ scope.row.course_name }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="scope">
                <el-tag 
                  :type="scope.row.status === 'public' ? 'success' : 'warning'"
                  size="small"
                >
                  {{ scope.row.status === 'public' ? '公开' : '私有' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="统计" width="200">
              <template #default="scope">
                <div class="note-stats">
                  <span><el-icon><View /></el-icon> {{ scope.row.view_count }}</span>
                  <span><el-icon><Star /></el-icon> {{ scope.row.like_count }}</span>
                  <span><el-icon><ChatDotRound /></el-icon> {{ scope.row.comment_count || 0 }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="120">
              <template #default="scope">
                {{ formatDate(scope.row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <div class="action-buttons">
                  <el-button size="small" @click="$router.push(`/notes/${scope.row.id}`)">
                    查看
                  </el-button>
                  <el-button size="small" @click="handleEditNote(scope.row.id)">
                    编辑
                  </el-button>
                  <el-button size="small" type="danger" @click="handleDeleteNote(scope.row.id)">
                    删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 分页 -->
        <div v-if="total > pageSize" class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Delete, Search, Document, View, Star, ChatDotRound, Clock
} from '@element-plus/icons-vue'
import { noteAPI } from '@/api'

export default {
  name: 'MyNotes',
  components: {
    Plus, Delete, Search, Document, View, Star, ChatDotRound, Clock
  },
  setup() {
    const router = useRouter()
    
    // 响应式数据
    const loading = ref(false)
    const notes = ref([])
    const selectedNotes = ref([])
    const searchQuery = ref('')
    const searchTimer = ref(null)
    const searchHistory = ref(JSON.parse(localStorage.getItem('noteSearchHistory') || '[]'))

    const currentPage = ref(1)
    const pageSize = ref(20)
    const total = ref(0)

    // 统计数据
    const stats = reactive({
      total: 0,
      totalViews: 0,
      totalLikes: 0,
      totalComments: 0
    })

    // 计算属性
    const filteredNotes = computed(() => {
      let filtered = notes.value

      // 搜索筛选
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(note => {
          // 搜索标题
          const titleMatch = note.title && note.title.toLowerCase().includes(query)
          // 搜索描述
          const descMatch = note.description && note.description.toLowerCase().includes(query)
          // 搜索课程名称
          const courseMatch = note.course_name && note.course_name.toLowerCase().includes(query)
          // 搜索内容（如果有的话）
          const contentMatch = note.content && note.content.toLowerCase().includes(query)
          
          return titleMatch || descMatch || courseMatch || contentMatch
        })
      }

      return filtered
    })

    // 获取用户笔记
    const fetchMyNotes = async () => {
      try {
        loading.value = true
        const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
        const response = await noteAPI.getNotes({
          user_id: userInfo.id,
          status: 'all',
          page: currentPage.value,
          page_size: pageSize.value
        })
        
        notes.value = response.notes || []
        total.value = response.total || 0
        
        // 计算统计数据
        calculateStats()
      } catch (error) {
        console.error('获取笔记失败:', error)
        ElMessage.error('获取笔记失败')
      } finally {
        loading.value = false
      }
    }

    // 计算统计数据
    const calculateStats = () => {
      stats.total = notes.value.length
      stats.totalViews = notes.value.reduce((sum, note) => sum + (note.view_count || 0), 0)
      stats.totalLikes = notes.value.reduce((sum, note) => sum + (note.like_count || 0), 0)
      stats.totalComments = notes.value.reduce((sum, note) => sum + (note.comment_count || 0), 0)
    }

    // 处理选择变化
    const handleSelectionChange = (selection) => {
      selectedNotes.value = selection
    }

    // 处理搜索（带防抖）
    const handleSearch = () => {
      // 清除之前的定时器
      if (searchTimer.value) {
        clearTimeout(searchTimer.value)
      }
      
      // 设置新的定时器，300ms后执行搜索
      searchTimer.value = setTimeout(() => {
        const query = searchQuery.value.trim()
        if (query) {
          // 添加到搜索历史
          addToSearchHistory(query)
          console.log('搜索关键词:', query)
        }
      }, 300)
    }

    // 添加搜索历史
    const addToSearchHistory = (query) => {
      // 移除重复项
      const filtered = searchHistory.value.filter(item => item !== query)
      // 添加到开头
      searchHistory.value = [query, ...filtered].slice(0, 10) // 最多保存10条
      // 保存到localStorage
      localStorage.setItem('noteSearchHistory', JSON.stringify(searchHistory.value))
    }

    // 选择搜索历史
    const selectSearchHistory = (query) => {
      searchQuery.value = query
      handleSearch()
    }

    // 清除搜索历史
    const clearSearchHistory = () => {
      searchHistory.value = []
      localStorage.removeItem('noteSearchHistory')
    }

    // 重置搜索
    const handleReset = () => {
      searchQuery.value = ''
      console.log('已重置搜索条件')
    }

    // 编辑笔记
    const handleEditNote = (noteId) => {
      // 这里应该跳转到编辑页面
      router.push(`/notes/${noteId}/edit`)
    }


    // 删除单个笔记
    const handleDeleteNote = async (noteId) => {
      try {
        await ElMessageBox.confirm('确定要删除这篇笔记吗？', '确认删除', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        await noteAPI.deleteNote(noteId)
        
        // 从列表中移除
        const index = notes.value.findIndex(note => note.id === noteId)
        if (index > -1) {
          notes.value.splice(index, 1)
          stats.total--
        }
        
        ElMessage.success('删除成功')
      } catch (error) {
        if (error !== 'cancel') {
          console.error('删除失败:', error)
          ElMessage.error('删除失败')
        }
      }
    }

    // 批量删除
    const handleBatchDelete = async () => {
      try {
        await ElMessageBox.confirm(
          `确定要删除选中的 ${selectedNotes.value.length} 篇笔记吗？`,
          '确认批量删除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        // 批量删除
        const deletePromises = selectedNotes.value.map(note => 
          noteAPI.deleteNote(note.id)
        )
        
        await Promise.all(deletePromises)
        
        // 从列表中移除
        const deletedIds = selectedNotes.value.map(note => note.id)
        notes.value = notes.value.filter(note => !deletedIds.includes(note.id))
        selectedNotes.value = []
        
        // 重新计算统计
        calculateStats()
        
        ElMessage.success('批量删除成功')
      } catch (error) {
        if (error !== 'cancel') {
          console.error('批量删除失败:', error)
          ElMessage.error('批量删除失败')
        }
      }
    }

    // 分页处理
    const handleSizeChange = (newSize) => {
      pageSize.value = newSize
      currentPage.value = 1
      fetchMyNotes()
    }

    const handleCurrentChange = (newPage) => {
      currentPage.value = newPage
      fetchMyNotes()
    }

    // 格式化日期
    const formatDate = (dateStr) => {
      const date = new Date(dateStr)
      return date.toLocaleDateString('zh-CN')
    }

    // 组件挂载时获取数据
    onMounted(() => {
      fetchMyNotes()
    })

    return {
      loading,
      notes,
      selectedNotes,
      searchQuery,
      searchHistory,
      currentPage,
      pageSize,
      total,
      stats,
      filteredNotes,
      handleSelectionChange,
      handleSearch,
      handleReset,
      selectSearchHistory,
      clearSearchHistory,
      handleEditNote,
      handleDeleteNote,
      handleBatchDelete,
      handleSizeChange,
      handleCurrentChange,
      formatDate
    }
  }
}
</script>

<style scoped>
.my-notes-page {
  min-height: 100vh;
  background: var(--bg-color);
}

.page-header {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-dark) 100%);
  color: white;
  padding: 60px 0;
  text-align: center;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 600;
  margin-bottom: 10px;
}

.page-subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 30px 0;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.left-actions,
.right-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.right-actions .el-input {
  transition: all 0.3s ease;
}

.right-actions .el-input:focus-within {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

.right-actions .el-button {
  transition: all 0.3s ease;
}

.right-actions .el-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.search-container {
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-history-dropdown {
  margin-left: 5px;
}

.search-history-dropdown .el-button {
  padding: 8px 12px;
  font-size: 12px;
  color: #909399;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.search-history-dropdown .el-button:hover {
  color: #409eff;
  border-color: #c6e2ff;
  background-color: #ecf5ff;
}

.search-result-info {
  margin-bottom: 20px;
}

.search-result-info .el-alert {
  border-radius: 8px;
  border: none;
  background: linear-gradient(135deg, #e3f2fd 0%, #f3e5f5 100%);
}

.stats-section {
  margin-bottom: 30px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  gap: 15px;
}

.stat-icon {
  font-size: 2rem;
  color: var(--primary-color);
}

.stat-value {
  font-size: 1.8rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 5px;
}

.stat-title {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.notes-section {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.loading-container,
.empty-container {
  padding: 40px;
  text-align: center;
}

.note-title-cell {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.note-description {
  font-size: 0.85rem;
  color: var(--text-secondary);
  line-height: 1.4;
}

.note-stats {
  display: flex;
  gap: 15px;
  font-size: 0.9rem;
  color: var(--text-light);
}

.note-stats span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .action-bar {
    flex-direction: column;
    gap: 15px;
  }
  
  .left-actions,
  .right-actions {
    width: 100%;
    justify-content: center;
  }
  
  .stats-section .el-row {
    flex-direction: column;
  }
  
  .stats-section .el-col {
    margin-bottom: 15px;
  }
}
</style>