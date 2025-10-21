<template>
  <div class="notes-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <h1 class="page-title">学习笔记</h1>
        <p class="page-subtitle">分享知识，共同成长</p>
      </div>
    </div>

    <div class="container">
      <!-- 搜索和筛选 -->
      <div class="search-section">
        <el-row :gutter="20">
          <el-col :span="15">
            <el-input
              v-model="searchQuery"
              placeholder="搜索笔记标题或内容"
              size="large"
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
          </el-col>

          <el-col :span="3">
            <el-button size="large" @click="resetFilters">
              重置
            </el-button>
          </el-col>

          <el-col :span="6">
            <el-button 
              v-if="isLoggedIn" 
              type="primary" 
              size="large" 
              @click="$router.push('/create-note')"
            >
              <el-icon><EditPen /></el-icon>
              发布笔记
            </el-button>
            <el-button 
              v-else 
              type="primary" 
              size="large" 
              @click="$router.push('/login')"
            >
              <el-icon><EditPen /></el-icon>
              登录后发布
            </el-button>
          </el-col>
        </el-row>
      </div>

      <!-- 笔记列表 -->
      <div class="notes-content">
        <div class="notes-header">
          <div class="notes-count">
            共找到 <span class="count-number">{{ total }}</span> 篇笔记
          </div>

        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-container">
          <el-skeleton :rows="4" animated />
          <el-skeleton :rows="4" animated />
          <el-skeleton :rows="4" animated />
        </div>

        <!-- 笔记列表 -->
        <div v-else class="notes-list">
          <transition-group name="note-list" tag="div">
            <div v-for="note in notes" :key="note.id" class="note-item">
              <el-card shadow="hover" @click="goToNoteDetail(note.id)">
                <div class="note-header">
                  <div class="note-meta">
                    <el-tag size="small" type="success">{{ note.course_name }}</el-tag>
                    <span class="note-time">{{ formatTime(note.created_at) }}</span>
                  </div>
                  <div class="note-actions">
                    <el-dropdown trigger="click" @command="handleNoteAction">
                      <el-icon class="action-icon"><MoreFilled /></el-icon>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item :command="{action: 'view', id: note.id}">查看详情</el-dropdown-item>

                          <el-dropdown-item v-if="canEditNote(note)" divided :command="{action: 'edit', id: note.id}">编辑</el-dropdown-item>
                          <el-dropdown-item v-if="canDeleteNote(note)" :command="{action: 'delete', id: note.id}">删除</el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                </div>
                
                <div class="note-content">
                  <h3 class="note-title">{{ note.title }}</h3>
                  <div class="note-author">
                    <el-avatar :size="24" :src="note.author_avatar">
                      <el-icon><User /></el-icon>
                    </el-avatar>
                    <span class="author-name">{{ note.username }}</span>
                  </div>
                  <p class="note-preview">{{ getPreviewText(note) }}</p>
                </div>
                
                <div class="note-footer">
                  <div class="note-stats">
                    <span class="stat-item">
                      <el-icon><View /></el-icon>
                      {{ note.view_count || 0 }}
                    </span>
                    <span class="stat-item">
                      <el-icon><ChatDotRound /></el-icon>
                      {{ note.comment_count || 0 }}
                    </span>
                    <span class="stat-item">
                      <el-icon><Star /></el-icon>
                      {{ note.like_count || 0 }}
                    </span>
                  </div>
                  <div class="note-tags">
                    <el-tag
                      v-for="tag in note.tags"
                      :key="tag"
                      size="small"
                      type="info"
                      effect="plain"
                    >
                      {{ tag }}
                    </el-tag>
                  </div>
                </div>
              </el-card>
            </div>
          </transition-group>
        </div>

        <!-- 空状态 -->
        <div v-if="!loading && notes.length === 0" class="empty-state">
          <el-empty description="暂无笔记数据">
            <el-button type="primary" @click="resetFilters">重置筛选</el-button>
          </el-empty>
        </div>

        <!-- 分页 -->
        <div v-if="!loading && notes.length > 0" class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50]"
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
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { noteAPI, roleUtils } from '../api'
import { Search, EditPen, User, View, ChatDotRound, Star, MoreFilled } from '@element-plus/icons-vue'

export default {
  name: 'Notes',
  components: {
    Search,
    EditPen,
    User,
    View,
    ChatDotRound,
    Star,
    MoreFilled
  },
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const notes = ref([])
    const total = ref(0)
    const currentPage = ref(1)
    const pageSize = ref(10)
    const searchQuery = ref('')


    const currentUserId = ref(1) // 模拟当前用户ID
    
    // 角色权限检查
    const canEditNote = (note) => {
      return noteAPI.canEditNote(note)
    }
    
    const canDeleteNote = (note) => {
      return noteAPI.canDeleteNote(note)
    }
    
    const isLoggedIn = computed(() => {
      return roleUtils.isLoggedIn()
    })

    // 获取笔记列表
    const fetchNotes = async () => {
      loading.value = true
      try {
        const params = {
          page: currentPage.value,
          limit: pageSize.value,
          search: searchQuery.value
        }
        
        const response = await noteAPI.getNotes(params)
        notes.value = response.notes || []
        total.value = response.total || 0
      } catch (error) {
        console.error('获取笔记列表失败:', error)
        ElMessage.error('获取笔记列表失败')
        // 使用模拟数据作为后备
        const mockNotes = [
          {
            id: 1,
            title: '微积分重点知识总结',
            content: '本笔记总结了微积分的重点知识点，包括导数、积分的基本概念和计算方法。导数是函数在某一点的瞬时变化率，积分是导数的逆运算...',
            username: '小明',
            author_id: 1,
            author_avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
            course_name: '高等数学',
            view_count: 245,
            comment_count: 12,
            like_count: 38,

            created_at: '2024-02-15T10:30:00Z'
          },
          {
            id: 2,
            title: '二叉树遍历算法详解',
            content: '详细介绍了二叉树的前序、中序、后序遍历算法，包含代码实现和时间复杂度分析。前序遍历：根-左-右，中序遍历：左-根-右...',
            username: '小红',
            author_id: 2,
            author_avatar: 'https://cube.elemecdn.com/9/c2/f0ee8a3c7c9638a54940382568c9dpng.png',
            course_name: '数据结构',
            view_count: 189,
            comment_count: 8,
            like_count: 25,

            created_at: '2024-02-14T14:20:00Z'
          },
          {
            id: 3,
            title: '英语语法要点整理',
            content: '整理了大学英语常用语法要点，包括时态、语态、从句等重要语法知识。现在时表示经常性动作，过去时表示过去发生的动作...',
            username: '小李',
            author_id: 3,
            author_avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
            course_name: '大学英语',
            view_count: 156,
            comment_count: 15,
            like_count: 42,

            created_at: '2024-02-13T09:15:00Z'
          },
          {
            id: 4,
            title: '矩阵运算技巧总结',
            content: '线性代数中矩阵运算的各种技巧和方法，包括矩阵乘法、求逆、特征值计算等。矩阵乘法满足结合律但不满足交换律...',
            username: '小王',
            author_id: 4,
            author_avatar: 'https://cube.elemecdn.com/6/94/4d3ea53c084bad6931a56d5158a48jpeg.jpeg',
            course_name: '线性代数',
            view_count: 203,
            comment_count: 6,
            like_count: 31,

            created_at: '2024-02-12T16:45:00Z'
          },
          {
            id: 5,
            title: '进程调度算法比较',
            content: '操作系统中各种进程调度算法的比较分析，包括FCFS、SJF、RR等算法的优缺点。FCFS算法简单但可能导致长作业等待...',
            username: '小张',
            author_id: 5,
            author_avatar: 'https://cube.elemecdn.com/d/e6/c4d93a3805b3ce3f323f7974e6f78jpeg.jpeg',
            course_name: '操作系统',
            view_count: 178,
            comment_count: 10,
            like_count: 28,

            created_at: '2024-02-11T11:30:00Z'
          }
        ]
        
        // 简单的筛选和搜索逻辑
        let filteredNotes = mockNotes
        

        
        if (searchQuery.value) {
          const query = searchQuery.value.toLowerCase()
          filteredNotes = filteredNotes.filter(note => 
            note.title.toLowerCase().includes(query) ||
            note.content.toLowerCase().includes(query)
          )
        }
        
        notes.value = filteredNotes
        total.value = filteredNotes.length
      } finally {
        loading.value = false
      }
    }

    // 搜索处理
    const handleSearch = () => {
      currentPage.value = 1
      fetchNotes()
    }





    // 分页大小变化
    const handleSizeChange = (size) => {
      pageSize.value = size
      currentPage.value = 1
      fetchNotes()
    }

    // 当前页变化
    const handleCurrentChange = (page) => {
      currentPage.value = page
      fetchNotes()
    }

    // 重置筛选
    const resetFilters = () => {
      searchQuery.value = ''
      currentPage.value = 1
      fetchNotes()
    }

    // 跳转到笔记详情
    const goToNoteDetail = (noteId) => {
      router.push(`/notes/${noteId}`)
    }

    // 处理笔记操作
    const handleNoteAction = async (command) => {
      const { action, id } = command
      
      switch (action) {
        case 'view':
          goToNoteDetail(id)
          break

        case 'edit':
          router.push(`/notes/${id}/edit`)
          break
        case 'delete':
          try {
            await ElMessageBox.confirm('确定要删除这篇笔记吗？', '确认删除', {
              confirmButtonText: '确定',
              cancelButtonText: '取消',
              type: 'warning'
            })
            
            await noteAPI.deleteNote(id)
            ElMessage.success('删除成功')
            fetchNotes()
          } catch (error) {
            if (error !== 'cancel') {
              console.error('删除失败:', error)
              ElMessage.error('删除失败')
            }
          }
          break
      }
    }

    // 获取预览文本
    const getPreviewText = (note) => {
      // 优先使用description字段
      if (note.description && typeof note.description === 'string' && note.description.trim()) {
        return note.description.length > 150 ? note.description.substring(0, 150) + '...' : note.description
      }
      // 如果没有description，使用content字段
      if (note.content && typeof note.content === 'string' && note.content.trim()) {
        return note.content.length > 150 ? note.content.substring(0, 150) + '...' : note.content
      }
      return '暂无内容预览'
    }

    // 格式化时间
    const formatTime = (timeStr) => {
      const date = new Date(timeStr)
      const now = new Date()
      const diff = now - date
      const days = Math.floor(diff / (1000 * 60 * 60 * 24))
      
      if (days === 0) {
        const hours = Math.floor(diff / (1000 * 60 * 60))
        if (hours === 0) {
          const minutes = Math.floor(diff / (1000 * 60))
          return minutes <= 0 ? '刚刚' : `${minutes}分钟前`
        }
        return `${hours}小时前`
      } else if (days < 7) {
        return `${days}天前`
      } else {
        return date.toLocaleDateString('zh-CN')
      }
    }

    onMounted(() => {
      fetchNotes()
    })

    return {
      loading,
      notes,
      total,
      currentPage,
      pageSize,
      searchQuery,
      currentUserId,
      isLoggedIn,
      canEditNote,
      canDeleteNote,
      handleSearch,
      handleSizeChange,
      handleCurrentChange,
      resetFilters,
      goToNoteDetail,
      handleNoteAction,
      getPreviewText,
      formatTime
    }
  }
}
</script>

<style scoped>
.search-section {
  margin-bottom: 30px;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: var(--shadow-light);
}

.notes-content {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: var(--shadow-light);
}

.notes-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid var(--border-color);
}

.notes-count {
  font-size: 1rem;
  color: var(--text-secondary);
}

.count-number {
  color: var(--primary-color);
  font-weight: 600;
}

.notes-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.note-item {
  cursor: pointer;
  transition: transform 0.3s ease;
}

.note-item:hover {
  transform: translateY(-2px);
}

.note-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.note-meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

.note-time {
  font-size: 0.9rem;
  color: var(--text-light);
}

.action-icon {
  color: var(--text-light);
  cursor: pointer;
  padding: 5px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.action-icon:hover {
  background-color: var(--bg-secondary);
  color: var(--primary-color);
}

.note-content {
  margin-bottom: 20px;
}

.note-title {
  font-size: 1.4rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 12px;
  line-height: 1.4;
}

.note-author {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.author-name {
  color: var(--primary-color);
  font-weight: 500;
  font-size: 0.95rem;
}

.note-preview {
  color: var(--text-secondary);
  line-height: 1.6;
  font-size: 0.95rem;
}

.note-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 15px;
  border-top: 1px solid var(--border-color);
}

.note-stats {
  display: flex;
  gap: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.9rem;
  color: var(--text-light);
}

.note-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.loading-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
}

/* 动画效果 */
.note-list-enter-active,
.note-list-leave-active {
  transition: all 0.3s ease;
}

.note-list-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.note-list-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .notes-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }
  
  .search-section .el-row {
    flex-direction: column;
  }
  
  .search-section .el-col {
    width: 100%;
    margin-bottom: 15px;
  }
  
  .note-footer {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }
  
  .note-stats {
    width: 100%;
  }
  
  .note-tags {
    width: 100%;
  }
}
</style>