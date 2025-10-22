<template>
  <div class="courses-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <h1 class="page-title">课程中心</h1>
        <p class="page-subtitle">发现优质课程，开启学习之旅</p>
      </div>
    </div>

    <div class="container">
      <!-- 搜索和筛选 -->
      <div class="search-section">
        <el-row :gutter="20">
          <el-col :span="20">
            <el-input
              v-model="searchQuery"
              placeholder="搜索课程名称、教师或描述"
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
          <el-col :span="4">
            <el-button size="large" @click="resetFilters">
              重置
            </el-button>
          </el-col>
        </el-row>
      </div>

      <!-- 课程列表 -->
      <div class="courses-content">
        <div class="courses-header">
          <div class="courses-count">
            共找到 <span class="count-number">{{ total }}</span> 门课程
          </div>
          <div class="header-actions">

          </div>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-container">
          <el-skeleton :rows="3" animated />
          <el-skeleton :rows="3" animated />
          <el-skeleton :rows="3" animated />
        </div>

        <!-- 课程卡片网格 -->
        <div v-else class="courses-grid">
          <transition-group name="course-list" tag="div" class="course-grid">
            <div v-for="course in courses" :key="course.id" class="course-card">
              <el-card shadow="hover" @click="goToCourseDetail(course.id)">
                <div class="course-header">
                  <div class="course-category">
                    <el-tag :type="getCategoryType(course.category)" size="small">
                      {{ course.category }}
                    </el-tag>
                  </div>
                  <div class="course-time">
                    {{ formatTime(course.created_at) }}
                  </div>
                </div>
                
                <div class="course-content">
                  <h3 class="course-title">{{ course.name }}</h3>
                  <p class="course-teacher">
                    <el-icon><User /></el-icon>
                    {{ course.teacher }}
                  </p>
                  <p class="course-description">{{ course.description }}</p>
                </div>
                
                <div class="course-footer">
                  <div class="course-stats">
                    <span class="stat-item">
                      <el-icon><Document /></el-icon>
                      {{ course.note_count || 0 }} 笔记
                    </span>
                    <span class="stat-item">
                      <el-icon><User /></el-icon>
                      {{ course.student_count || 0 }} 学习人数
                    </span>
                  </div>
                  <el-button type="primary" size="small">
                    查看详情
                  </el-button>
                </div>
              </el-card>
            </div>
          </transition-group>
        </div>

        <!-- 空状态 -->
        <div v-if="!loading && courses.length === 0" class="empty-state">
          <el-empty description="暂无课程数据">
            <el-button type="primary" @click="resetFilters">重置筛选</el-button>
          </el-empty>
        </div>

        <!-- 分页 -->
        <div v-if="!loading && courses.length > 0" class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[12, 24, 48]"
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
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { courseAPI, roleUtils } from '../api'
import { Search, User, Document, View } from '@element-plus/icons-vue'

export default {
  name: 'Courses',
  components: {
    Search,
    User,
    Document,
    View
  },
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const courses = ref([])
    const total = ref(0)
    const currentPage = ref(1)
    const pageSize = ref(12)
    const searchQuery = ref('')


    // 获取课程列表
    const fetchCourses = async () => {
      loading.value = true
      try {
        const params = {
          page: currentPage.value,
          page_size: pageSize.value,
          search: searchQuery.value
        }
        
        const response = await courseAPI.getCourses(params)
        courses.value = response.courses || []
        total.value = response.total || 0
      } catch (error) {
        // 使用模拟数据
        const mockCourses = [
          {
            id: 1,
            name: '高等数学A',
            teacher: '张教授',
            description: '微积分基础理论与应用，包括极限、导数、积分等核心概念',
            category: '数学',
            note_count: 25,
            view_count: 1200,
            created_at: '2024-01-15T10:00:00Z'
          },
          {
            id: 2,
            name: '数据结构与算法',
            teacher: '李教授',
            description: '计算机科学基础课程，涵盖线性表、树、图等数据结构',
            category: '计算机',
            note_count: 42,
            view_count: 2100,
            created_at: '2024-01-20T14:30:00Z'
          },
          {
            id: 3,
            name: '大学英语',
            teacher: '王教授',
            description: '英语听说读写综合训练，提高英语应用能力',
            category: '语言',
            note_count: 18,
            view_count: 800,
            created_at: '2024-01-25T09:15:00Z'
          },
          {
            id: 4,
            name: '线性代数',
            teacher: '赵教授',
            description: '矩阵理论、向量空间、线性变换等数学基础',
            category: '数学',
            note_count: 31,
            view_count: 1500,
            created_at: '2024-02-01T11:20:00Z'
          },
          {
            id: 5,
            name: '操作系统',
            teacher: '陈教授',
            description: '计算机操作系统原理，进程管理、内存管理、文件系统',
            category: '计算机',
            note_count: 28,
            view_count: 1800,
            created_at: '2024-02-05T16:45:00Z'
          },
          {
            id: 6,
            name: '微观经济学',
            teacher: '刘教授',
            description: '市场经济基本原理，供需关系、价格机制、市场结构',
            category: '经济',
            note_count: 22,
            view_count: 950,
            created_at: '2024-02-10T13:30:00Z'
          }
        ]
        
        // 简单的筛选和搜索逻辑
        let filteredCourses = mockCourses
        

        
        if (searchQuery.value) {
          const query = searchQuery.value.toLowerCase()
          filteredCourses = filteredCourses.filter(course => 
            course.name.toLowerCase().includes(query) ||
            course.teacher.toLowerCase().includes(query) ||
            course.description.toLowerCase().includes(query)
          )
        }
        
        // 分页处理
        const startIndex = (currentPage.value - 1) * pageSize.value
        const endIndex = startIndex + pageSize.value
        const paginatedCourses = filteredCourses.slice(startIndex, endIndex)
        
        courses.value = paginatedCourses
        total.value = filteredCourses.length
      } finally {
        loading.value = false
      }
    }

    // 搜索处理
    const handleSearch = () => {
      currentPage.value = 1
      fetchCourses()
    }





    // 分页大小变化
    const handleSizeChange = (size) => {
      pageSize.value = size
      currentPage.value = 1
      fetchCourses()
    }

    // 当前页变化
    const handleCurrentChange = (page) => {
      currentPage.value = page
      fetchCourses()
    }

    // 重置筛选
    const resetFilters = () => {
      searchQuery.value = ''
      currentPage.value = 1
      fetchCourses()
    }

    // 跳转到课程详情
    const goToCourseDetail = (courseId) => {
      router.push(`/courses/${courseId}`)
    }

    // 获取分类标签类型
    const getCategoryType = (category) => {
      const typeMap = {
        '数学': 'primary',
        '计算机': 'success',
        '语言': 'warning',
        '物理': 'info',
        '化学': 'danger',
        '经济': 'warning',
        '管理': 'success'
      }
      return typeMap[category] || 'info'
    }

    // 格式化时间
    const formatTime = (timeStr) => {
      const date = new Date(timeStr)
      return date.toLocaleDateString('zh-CN')
    }

    onMounted(() => {
      fetchCourses()
    })

    return {
      loading,
      courses,
      total,
      currentPage,
      pageSize,
      searchQuery,
      handleSearch,
      handleSizeChange,
      handleCurrentChange,
      resetFilters,
      goToCourseDetail,
      getCategoryType,
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

.courses-content {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: var(--shadow-light);
}

.courses-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid var(--border-color);
}

.courses-count {
  font-size: 1rem;
  color: var(--text-secondary);
}

.count-number {
  color: var(--primary-color);
  font-weight: 600;
}

.course-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.course-card {
  cursor: pointer;
  transition: transform 0.3s ease;
}

.course-card:hover {
  transform: translateY(-5px);
}

.course-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.course-time {
  font-size: 0.9rem;
  color: var(--text-light);
}

.course-content {
  margin-bottom: 20px;
}

.course-title {
  font-size: 1.3rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 10px;
  line-height: 1.4;
}

.course-teacher {
  display: flex;
  align-items: center;
  color: var(--primary-color);
  font-weight: 500;
  margin-bottom: 10px;
}

.course-teacher .el-icon {
  margin-right: 5px;
}

.course-description {
  color: var(--text-secondary);
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.course-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 15px;
  border-top: 1px solid var(--border-color);
}

.course-stats {
  display: flex;
  gap: 15px;
}

.stat-item {
  display: flex;
  align-items: center;
  font-size: 0.9rem;
  color: var(--text-light);
}

.stat-item .el-icon {
  margin-right: 4px;
}

.loading-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
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
.course-list-enter-active,
.course-list-leave-active {
  transition: all 0.3s ease;
}

.course-list-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.course-list-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .courses-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }
  
  .course-grid {
    grid-template-columns: 1fr;
  }
  
  .search-section .el-row {
    flex-direction: column;
  }
  
  .search-section .el-col {
    width: 100%;
    margin-bottom: 15px;
  }
}
</style>