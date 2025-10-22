<template>
  <div class="search-page">
    <!-- 搜索头部 -->
    <div class="search-header">
      <div class="container">
        <div class="search-main">
          <h1 class="search-title">搜索</h1>
          <div class="search-input-container">
            <el-input
              v-model="searchQuery"
              placeholder="搜索课程、笔记、用户..."
              size="large"
              class="search-input"
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
              <template #append>
                <el-button type="primary" @click="handleSearch" :loading="searching">
                  搜索
                </el-button>
              </template>
            </el-input>
          </div>
          
        </div>
      </div>
    </div>

    <div class="container">
      <el-row :gutter="20">
        <!-- 搜索过滤器 -->
        <el-col :span="6">
          <div class="search-filters">
            <el-card>
              <template #header>
                <div class="filter-header">
                  <span>筛选条件</span>
                  <el-button link size="small" @click="resetFilters">重置</el-button>
                </div>
              </template>
              
              <!-- 内容类型 -->
              <div class="filter-section">
                <h4 class="filter-title">内容类型</h4>
                <el-checkbox-group v-model="filters.types">
                  <el-checkbox label="notes">笔记</el-checkbox>
                  <el-checkbox label="courses">课程</el-checkbox>
                  <el-checkbox label="users">用户</el-checkbox>
                </el-checkbox-group>
              </div>
              
              <!-- 课程分类 -->
              <div class="filter-section">
                <h4 class="filter-title">课程分类</h4>
                <el-checkbox-group v-model="filters.categories">
                  <el-checkbox label="math">数学</el-checkbox>
                  <el-checkbox label="computer">计算机</el-checkbox>
                  <el-checkbox label="physics">物理</el-checkbox>
                  <el-checkbox label="chemistry">化学</el-checkbox>
                  <el-checkbox label="english">英语</el-checkbox>
                </el-checkbox-group>
              </div>
              
              <!-- 难度等级 -->
              <div class="filter-section">
                <h4 class="filter-title">难度等级</h4>
                <el-checkbox-group v-model="filters.difficulty">
                  <el-checkbox label="beginner">入门</el-checkbox>
                  <el-checkbox label="intermediate">中级</el-checkbox>
                  <el-checkbox label="advanced">高级</el-checkbox>
                </el-checkbox-group>
              </div>
              
              <!-- 排序方式 -->
              <div class="filter-section">
                <h4 class="filter-title">排序方式</h4>
                <el-radio-group v-model="filters.sortBy">
                  <el-radio value="time">时间</el-radio>
                  <el-radio value="popularity">热度</el-radio>
                </el-radio-group>
              </div>
            </el-card>
          </div>
        </el-col>

        <!-- 搜索结果 -->
        <el-col :span="18">
          <div class="search-results">
            <!-- 搜索结果头部 -->
            <div class="results-header" v-if="hasSearched">
              <div class="results-info">
                <span class="results-count">
                  找到 {{ totalResults }} 个结果
                  <span v-if="searchQuery">关于 "{{ searchQuery }}"</span>
                </span>
                <span class="search-time">用时 {{ searchTime }}ms</span>
              </div>
              
              <!-- 结果类型标签 -->
              <el-tabs v-model="activeResultType" @tab-change="handleTabChange">
                <el-tab-pane label="全部" name="all">
                  <template #label>
                    全部 <el-badge :value="totalResults" class="tab-badge" />
                  </template>
                </el-tab-pane>
                <el-tab-pane label="笔记" name="notes">
                  <template #label>
                    笔记 <el-badge :value="results.notes.length" class="tab-badge" />
                  </template>
                </el-tab-pane>
                <el-tab-pane label="课程" name="courses">
                  <template #label>
                    课程 <el-badge :value="results.courses.length" class="tab-badge" />
                  </template>
                </el-tab-pane>
                <el-tab-pane label="用户" name="users">
                  <template #label>
                    用户 <el-badge :value="results.users.length" class="tab-badge" />
                  </template>
                </el-tab-pane>
              </el-tabs>
            </div>

            <!-- 加载状态 -->
            <div v-if="searching" class="loading-container">
              <el-skeleton :rows="5" animated />
            </div>

            <!-- 无搜索状态 -->
            <div v-else-if="!hasSearched" class="no-search-state">
              <el-empty description="输入关键词开始搜索">
                <div class="search-suggestions">
                  <h4>搜索建议</h4>
                  <ul>
                    <li>使用具体的关键词，如"微积分"、"数据结构"</li>
                    <li>尝试使用同义词或相关词汇</li>
                    <li>检查拼写是否正确</li>
                    <li>使用筛选条件缩小搜索范围</li>
                  </ul>
                </div>
              </el-empty>
            </div>

            <!-- 无结果状态 -->
            <div v-else-if="totalResults === 0" class="no-results-state">
              <el-empty description="没有找到相关结果">
                <div class="no-results-suggestions">
                  <h4>搜索建议</h4>
                  <ul>
                    <li>尝试使用更通用的关键词</li>
                    <li>减少筛选条件</li>
                    <li>检查拼写是否正确</li>
                    <li>尝试搜索相关的主题</li>
                  </ul>
                  <el-button type="primary" @click="resetFilters">清除筛选条件</el-button>
                </div>
              </el-empty>
            </div>

            <!-- 搜索结果列表 -->
            <div v-else class="results-list">
              <!-- 笔记结果 -->
              <div v-if="activeResultType === 'all' || activeResultType === 'notes'">
                <div v-if="displayResults.notes.length > 0" class="result-section">
                  <h3 v-if="activeResultType === 'all'" class="section-title">
                    笔记 ({{ results.notes.length }})
                  </h3>
                  <div class="note-results">
                    <div
                      v-for="note in displayResults.notes"
                      :key="`note-${note.id}-${searchQuery}`"
                      class="result-item note-item"
                      @click="$router.push(`/notes/${note.id}`)"
                    >
                      <div class="item-content">
                        <h4 class="item-title" v-html="note && note.title ? highlightText(note.title) : ''"></h4>
                        <p class="item-meta">
                          <el-tag size="small" type="success">{{ note.course_name || '' }}</el-tag>
                          <span class="date">{{ note.created_at ? formatDate(note.created_at) : '' }}</span>
                        </p>
                        <p class="item-description" v-html="note && note.excerpt ? highlightText(note.excerpt) : ''"></p>
                        <div class="item-stats">
                          <span><el-icon><View /></el-icon> {{ note.view_count ?? 0 }}</span>
                          <span><el-icon><ChatDotRound /></el-icon> {{ note.comment_count ?? 0 }}</span>
                          <span><el-icon><Star /></el-icon> {{ note.like_count ?? 0 }}</span>
                        </div>
                      </div>
                      <div class="item-actions">
                        <el-button size="small" type="primary">查看</el-button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 课程结果 -->
              <div v-if="activeResultType === 'all' || activeResultType === 'courses'">
                <div v-if="displayResults.courses.length > 0" class="result-section">
                  <h3 v-if="activeResultType === 'all'" class="section-title">
                    课程 ({{ results.courses.length }})
                  </h3>
                  <div class="course-results">
                    <div
                      v-for="course in displayResults.courses"
                      :key="`course-${course.id}-${searchQuery}`"
                      class="result-item course-item"
                      @click="$router.push(`/courses/${course.id}`)"
                    >
                      <div class="item-content">
                        <h4 class="item-title" v-html="course && course.name ? highlightText(course.name) : ''"></h4>
                        <p class="item-meta">
                          <el-tag size="small">{{ course.school || '' }}</el-tag>
                          <el-tag size="small">{{ course.department || '' }}</el-tag>
                          <span class="instructor">{{ course.teacher || '' }}</span>
                        </p>
                        <p class="item-description" v-html="course && course.description ? highlightText(course.description) : ''"></p>
                        <div class="item-stats">
                          <!-- 语义检索结果暂不包含以下统计，展示为0或隐藏 -->
                          <span><el-icon><User /></el-icon> {{ course.student_count ?? 0 }} 学生</span>
                          <span><el-icon><Document /></el-icon> {{ course.note_count ?? 0 }} 笔记</span>
                        </div>
                      </div>
                      <div class="item-actions">
                        <el-button size="small" type="primary">查看课程</el-button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 用户结果 -->
              <div v-if="activeResultType === 'all' || activeResultType === 'users'">
                <div v-if="displayResults.users.length > 0" class="result-section">
                  <h3 v-if="activeResultType === 'all'" class="section-title">
                    用户 ({{ results.users.length }})
                  </h3>
                  <div class="user-results">
                    <div
                      v-for="user in displayResults.users"
                      :key="`user-${user.id}-${searchQuery}`"
                      class="result-item user-item"
                      @click="$router.push(`/users/${user.id}`)"
                    >
                      <div class="user-avatar">
                        <el-avatar :size="50" :src="user.avatar">
                          <el-icon><User /></el-icon>
                        </el-avatar>
                      </div>
                      <div class="item-content">
                        <h4 class="item-title" v-html="user && user.username ? highlightText(user.username) : ''"></h4>
                        <p class="item-meta">
                          <span class="school">{{ user.school || '' }}</span>
                          <span class="major">{{ user.major || '' }}</span>
                        </p>
                        <p class="item-description">{{ user.bio || '' }}</p>
                        <div class="item-stats">
                          <span><el-icon><Document /></el-icon> {{ user.note_count }} 笔记</span>
                          <span><el-icon><Star /></el-icon> {{ user.like_count }} 获赞</span>
                        </div>
                      </div>
                      <div class="item-actions">
                        <el-button size="small" type="primary">查看主页</el-button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 分页 -->
              <div class="pagination-container">
                <el-pagination
                  v-model:current-page="currentPage"
                  :page-size="pageSize"
                  :total="totalResults"
                  layout="prev, pager, next, jumper, total"
                  @current-change="handlePageChange"
                />
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { searchAPI } from '../api'
import {
  Search,
  View,
  ChatDotRound,
  Star,
  User,
  Document
} from '@element-plus/icons-vue'

export default {
  name: 'Search',
  components: {
    Search,
    View,
    ChatDotRound,
    Star,
    User,
    Document
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    
    const searchQuery = ref('')
    const searching = ref(false)
    const hasSearched = ref(false)
    const searchTime = ref(0)
    const activeResultType = ref('all')
    const currentPage = ref(1)
    const pageSize = ref(10)
    
    
    // 筛选条件
    const filters = reactive({
      types: ['notes', 'courses', 'users'],
      categories: [],
      difficulty: [],
      sortBy: 'time'
    })

    
    // 搜索结果
    const results = reactive({
      notes: [],
      courses: [],
      users: []
    })
    
    // 防止无限循环的标志
    let isUpdatingFilters = false
    
    
    // 计算属性
    const totalResults = computed(() => {
      const notesLength = Array.isArray(results.notes) ? results.notes.length : 0
      const coursesLength = Array.isArray(results.courses) ? results.courses.length : 0
      const usersLength = Array.isArray(results.users) ? results.users.length : 0
      return notesLength + coursesLength + usersLength
    })
    
    const displayResults = computed(() => {
      const start = (currentPage.value - 1) * pageSize.value
      const end = start + pageSize.value
      
      // 确保results对象的属性存在且为数组
      const safeResults = {
        notes: Array.isArray(results.notes) ? results.notes : [],
        courses: Array.isArray(results.courses) ? results.courses : [],
        users: Array.isArray(results.users) ? results.users : []
      }
      
      if (activeResultType.value === 'all') {
        return {
          notes: safeResults.notes.slice(0, 3), // 在全部结果中只显示前3个
          courses: safeResults.courses.slice(0, 3),
          users: safeResults.users.slice(0, 3)
        }
      } else {
        return {
          notes: activeResultType.value === 'notes' ? safeResults.notes.slice(start, end) : [],
          courses: activeResultType.value === 'courses' ? safeResults.courses.slice(start, end) : [],
          users: activeResultType.value === 'users' ? safeResults.users.slice(start, end) : []
        }
      }
    })
    
    // 搜索方法
    const handleSearch = async () => {
      if (!searchQuery.value.trim()) {
        ElMessage.warning('请输入搜索关键词')
        return
      }
      searching.value = true
      const startTime = Date.now()
      try {
        // 简单笔记搜索：统一按笔记查询
        const typeParam = 'notes'
        // 排序映射到后端可识别字段
        const sortMap = { relevance: 'created_at', time: 'created_at', popularity: 'like_count' }
        const sortByParam = sortMap[filters.sortBy] || 'created_at'
        
        const params = {
          keyword: searchQuery.value,
          page: currentPage.value,
          page_size: pageSize.value,
          sort_by: sortByParam,
          order: 'desc'
        }

        const response = await searchAPI.searchNotes(params)
        const notes = Array.isArray(response?.notes) ? response.notes : []
        // 兼容当前展示结构：补充 excerpt 字段
        results.notes = notes.map(n => ({
          ...n,
          excerpt: n.description || ''
        }))
        results.courses = []
        results.users = []
        hasSearched.value = true
        searchTime.value = Date.now() - startTime
        currentPage.value = 1
        router.push({ query: { q: searchQuery.value, type: typeParam, sort_by: sortByParam } })
      } catch (error) {
        ElMessage.error('搜索失败，请稍后重试')
      } finally {
        searching.value = false
      }
    }
    
    // 重置筛选条件
    const resetFilters = () => {
      // 设置标志防止watch触发
      isUpdatingFilters = true
      
      Object.assign(filters, {
        types: ['notes', 'courses', 'users'],
        categories: [],
        difficulty: [],
        sortBy: 'time'
      })
      
      // 重置标志并手动触发搜索
      isUpdatingFilters = false
      if (hasSearched.value) {
        handleSearch()
      }
    }
    
    // 标签页切换
    const handleTabChange = (tabName) => {
      activeResultType.value = tabName
      currentPage.value = 1
    }
    
    // 分页处理
    const handlePageChange = (page) => {
      currentPage.value = page
    }
    
    // 高亮搜索关键词
    const highlightText = (text) => {
      if (!searchQuery.value || !text || typeof text !== 'string') return text || ''
      
      try {
        // 转义HTML特殊字符以防止XSS
        const escapedText = text.replace(/[&<>"']/g, (match) => {
          const escapeMap = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#39;'
          }
          return escapeMap[match]
        })
        
        // 转义搜索关键词中的特殊正则字符
        const escapedQuery = searchQuery.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
        const regex = new RegExp(`(${escapedQuery})`, 'gi')
        return escapedText.replace(regex, '<mark>$1</mark>')
      } catch (error) {
        console.warn('Error in highlightText:', error)
        return text || ''
      }
    }
    
    // 获取难度类型
    const getDifficultyType = (difficulty) => {
      const types = {
        beginner: 'success',
        intermediate: 'warning',
        advanced: 'danger'
      }
      return types[difficulty] || 'info'
    }
    
    // 获取难度文本
    const getDifficultyText = (difficulty) => {
      const texts = {
        beginner: '入门',
        intermediate: '中级',
        advanced: '高级'
      }
      return texts[difficulty] || difficulty
    }
    
    // 格式化日期
    const formatDate = (dateStr) => {
      const date = new Date(dateStr)
      return date.toLocaleDateString('zh-CN')
    }
    
    // 监听筛选条件变化
    watch(filters, () => {
      if (hasSearched.value && !isUpdatingFilters) {
        handleSearch()
      }
    }, { deep: true })
    
    // 初始化
    onMounted(() => {
      // 从URL参数获取搜索关键词
      if (route.query.q) {
        searchQuery.value = route.query.q
        handleSearch()
      }
    })
    
    return {
      searchQuery,
      searching,
      hasSearched,
      searchTime,
      activeResultType,
      currentPage,
      pageSize,
      filters,
      results,
      totalResults,
      displayResults,
      handleSearch,
      resetFilters,
      handleTabChange,
      handlePageChange,
      highlightText,
      getDifficultyType,
      getDifficultyText,
      formatDate
    }
  }
}
</script>

<style scoped>
.search-header {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: white;
  padding: 40px 0;
}

.search-main {
  text-align: center;
}

.search-title {
  font-size: 2.5rem;
  font-weight: 600;
  margin-bottom: 30px;
}

.search-input-container {
  max-width: 600px;
  margin: 0 auto 20px;
}

.search-input {
  box-shadow: var(--shadow-medium);
}

.quick-tags {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.tags-label {
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.9rem;
}

.hot-tag {
  cursor: pointer;
  transition: all 0.3s;
}

.hot-tag:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-light);
}

.search-filters {
  position: sticky;
  top: 20px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-section {
  margin-bottom: 25px;
}

.filter-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 10px;
}

.filter-section .el-checkbox-group,
.filter-section .el-radio-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.results-header {
  margin-bottom: 20px;
}

.results-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 15px;
  border-bottom: 1px solid var(--border-color);
}

.results-count {
  font-size: 1.1rem;
  color: var(--text-primary);
}

.search-time {
  font-size: 0.9rem;
  color: var(--text-light);
}

.tab-badge {
  margin-left: 5px;
}

.loading-container {
  padding: 20px;
}

.no-search-state,
.no-results-state {
  text-align: center;
  padding: 40px 20px;
}

.search-suggestions,
.no-results-suggestions {
  max-width: 400px;
  margin: 20px auto;
  text-align: left;
}

.search-suggestions h4,
.no-results-suggestions h4 {
  margin-bottom: 15px;
  color: var(--text-primary);
}

.search-suggestions ul,
.no-results-suggestions ul {
  list-style: disc;
  padding-left: 20px;
  color: var(--text-secondary);
  line-height: 1.6;
}

.results-list {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.result-section {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: var(--shadow-light);
}

.section-title {
  font-size: 1.3rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 2px solid var(--primary-color);
}

.result-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin-bottom: 15px;
  cursor: pointer;
  transition: all 0.3s;
}

.result-item:hover {
  border-color: var(--primary-color);
  box-shadow: var(--shadow-light);
  transform: translateY(-2px);
}

.result-item:last-child {
  margin-bottom: 0;
}

.item-content {
  flex: 1;
}

.item-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
  line-height: 1.4;
}

.item-title :deep(mark) {
  background: var(--warning-color);
  color: var(--text-primary);
  padding: 2px 4px;
  border-radius: 3px;
}

.item-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
  font-size: 0.9rem;
  color: var(--text-light);
}

.item-description {
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 15px;
}

.item-description :deep(mark) {
  background: var(--warning-color);
  color: var(--text-primary);
  padding: 2px 4px;
  border-radius: 3px;
}

.item-stats {
  display: flex;
  gap: 20px;
  font-size: 0.9rem;
  color: var(--text-light);
}

.item-stats span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.item-actions {
  margin-left: 20px;
}

.user-item {
  align-items: center;
}

.user-avatar {
  margin-right: 15px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 30px;
  padding: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .search-title {
    font-size: 2rem;
  }
  
  .container .el-row {
    flex-direction: column;
  }
  
  .search-filters {
    position: static;
    margin-bottom: 20px;
  }
  
  .results-info {
    flex-direction: column;
    gap: 10px;
    align-items: flex-start;
  }
  
  .result-item {
    flex-direction: column;
    gap: 15px;
  }
  
  .item-actions {
    margin-left: 0;
    align-self: flex-start;
  }
  
  .user-item {
    flex-direction: row;
    align-items: flex-start;
  }
  
  .quick-tags {
    justify-content: flex-start;
  }
}
</style>