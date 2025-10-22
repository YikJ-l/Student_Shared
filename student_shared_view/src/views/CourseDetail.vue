<template>
  <div class="course-detail-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <el-breadcrumb separator=">">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/courses' }">课程中心</el-breadcrumb-item>
          <el-breadcrumb-item>{{ course.name }}</el-breadcrumb-item>
        </el-breadcrumb>
      </div>
    </div>

    <div class="container">
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="8" animated />
      </div>

      <div v-else-if="course.id" class="course-detail">
        <!-- 课程基本信息 -->
        <el-card class="course-info-card">
          <div class="course-header">
            <div class="course-main-info">
              <h1 class="course-title">{{ course.name }}</h1>
              <div class="course-meta">
                <el-tag :type="getCategoryType(course.category)" size="large">
                  {{ course.category }}
                </el-tag>
                <span class="course-teacher">
                  <el-icon><User /></el-icon>
                  {{ course.teacher }}
                </span>
                <span class="course-time">
                  <el-icon><Clock /></el-icon>
                  {{ formatTime(course.created_at) }}
                </span>
              </div>
              <p class="course-description">{{ course.description }}</p>
            </div>
            
            <div class="course-stats">
              <div class="stat-item">
                <div class="stat-number">{{ course.note_count || 0 }}</div>
                <div class="stat-label">相关笔记</div>
              </div>             
              <div class="stat-item">
                <div class="stat-number">{{ course.student_count || 0 }}</div>
                <div class="stat-label">学习人数</div>
              </div>
            </div>
          </div>
        </el-card>

        <!-- 课程详细内容 -->
        <el-row :gutter="20">
          <el-col :span="16">
            <!-- 课程介绍 -->
            <el-card class="course-content-card">
              <template #header>
                <span>课程介绍</span>
              </template>
              <div class="course-content">
                <div v-html="course.content || course.description"></div>
              </div>
            </el-card>

            <!-- 相关笔记 -->
            <el-card class="related-notes-card">
              <template #header>
                <div class="card-header">
                  <span>相关笔记 ({{ relatedNotes.length }})</span>
                  <el-button type="primary" size="small" @click="$router.push('/create-note')">
                    <el-icon><EditPen /></el-icon>
                    发布笔记
                  </el-button>
                </div>
              </template>
              
              <div v-if="relatedNotes.length === 0" class="empty-notes">
                <el-empty description="暂无相关笔记">
                  <el-button type="primary" @click="$router.push('/create-note')">发布第一篇笔记</el-button>
                </el-empty>
              </div>
              
              <div v-else class="notes-list">
                <div v-for="note in relatedNotes" :key="note.id" class="note-item">
                  <el-card shadow="hover" @click="$router.push(`/notes/${note.id}`)">
                    <div class="note-header">
                      <h4 class="note-title">{{ note.title }}</h4>
                      <div class="note-meta">
                        <span class="note-author">
                          <el-icon><User /></el-icon>
                          {{ note.username || '未知作者' }}
                        </span>
                        <span class="note-time">{{ formatTime(note.created_at) }}</span>
                      </div>
                    </div>
                    <p class="note-preview">{{ getPreviewText(note.content) }}</p>
                    <div class="note-footer">
                      <div class="note-stats">

                        <span><el-icon><ChatDotRound /></el-icon> {{ note.comment_count || 0 }}</span>
                        <span><el-icon><Star /></el-icon> {{ note.like_count || 0 }}</span>
                      </div>
                      <div class="note-tags">
                        <el-tag v-for="tag in (note.tags || [])" :key="tag" size="small" type="info" effect="plain">
                          {{ tag }}
                        </el-tag>
                      </div>
                    </div>
                  </el-card>
                </div>
              </div>
            </el-card>
          </el-col>

          <el-col :span="8">
            <!-- 课程信息侧边栏 -->
            <el-card class="course-sidebar">
              <template #header>
                <span>课程信息</span>
              </template>
              
              <div class="sidebar-content">
                <div class="info-item">
                  <label>课程分类：</label>
                  <span>{{ course.category }}</span>
                </div>
                <div class="info-item">
                  <label>授课教师：</label>
                  <span>{{ course.teacher }}</span>
                </div>
                <div class="info-item">
                  <label>创建时间：</label>
                  <span>{{ formatDate(course.created_at) }}</span>
                </div>
                <div class="info-item">
                  <label>更新时间：</label>
                  <span>{{ formatDate(course.updated_at || course.created_at) }}</span>
                </div>
              </div>
              
              <div class="sidebar-actions">
                <el-button 
                  v-if="!isJoined" 
                  type="primary" 
                  size="large" 
                  @click="handleJoinCourse"
                >
                  <el-icon><Plus /></el-icon>
                  加入学习
                </el-button>
                <el-button 
                  v-else 
                  type="success" 
                  size="large" 
                  disabled
                >
                  <el-icon><Check /></el-icon>
                  已加入
                </el-button>

              </div>
            </el-card>

            <!-- 推荐课程 -->
            <el-card class="recommended-courses">
              <template #header>
                <span>推荐课程</span>
              </template>
              
              <div class="recommended-list">
                <div v-for="recommendedCourse in recommendedCourses" :key="recommendedCourse.id" class="recommended-item" @click="$router.push(`/courses/${recommendedCourse.id}`)">
                  <div class="recommended-info">
                    <h5>{{ recommendedCourse.name }}</h5>
                    <p>{{ recommendedCourse.teacher }}</p>
                    <div class="recommended-stats">
                      <span><el-icon><Document /></el-icon> {{ recommendedCourse.note_count || 0 }}</span>
                      <span><el-icon><User /></el-icon> {{ recommendedCourse.student_count || 0 }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 错误状态 -->
      <div v-else class="error-state">
        <el-result icon="warning" title="课程不存在" sub-title="抱歉，您访问的课程不存在或已被删除">
          <template #extra>
            <el-button type="primary" @click="$router.push('/courses')">返回课程列表</el-button>
          </template>
        </el-result>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { courseAPI, noteAPI } from '../api'
import { User, Clock, EditPen, View, ChatDotRound, Star, Plus, Document, Check } from '@element-plus/icons-vue'

export default {
  name: 'CourseDetail',
  components: {
    User,
    Clock,
    EditPen,
    View,
    ChatDotRound,
    Star,
    Plus,

    Document,
    Check
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const loading = ref(true)
    const course = ref({})
    const relatedNotes = ref([])
    const recommendedCourses = ref([])
    const isJoined = ref(false)

    // 获取课程详情
    const fetchCourseDetail = async () => {
      try {
        const courseId = route.params.id
        const response = await courseAPI.getCourseDetail(courseId)
        course.value = response
        isJoined.value = response.is_joined || false
        console.log('课程详情获取成功:', response)
        
        // 获取相关笔记
        await fetchRelatedNotes(courseId)
        
        // 获取推荐课程
        await fetchRecommendedCourses()
      } catch (error) {
        console.error('获取课程详情失败:', error)
        // 使用模拟数据
        const courseId = parseInt(route.params.id)
        course.value = {
          id: courseId,
          name: '高等数学A',
          teacher: '张教授',
          description: '微积分基础理论与应用，包括极限、导数、积分等核心概念。本课程旨在培养学生的数学思维能力和解决实际问题的能力。',
          category: '数学',
          content: '<h3>课程目标</h3><p>通过本课程的学习，学生将掌握微积分的基本理论和方法，能够运用微积分知识解决实际问题。</p><h3>主要内容</h3><ul><li>函数与极限</li><li>导数与微分</li><li>积分学</li><li>微分方程</li></ul>',
          note_count: 25,
          view_count: 1200,
          student_count: 89,
          created_at: '2024-01-15T10:00:00Z',
          updated_at: '2024-02-20T15:30:00Z'
        }
        
        // 模拟相关笔记
        relatedNotes.value = [
          {
            id: 1,
            title: '微积分重点知识总结',
            content: '本笔记总结了微积分的重点知识点，包括导数、积分的基本概念和计算方法...',
            author: '小明',
            view_count: 245,
            comment_count: 12,
            like_count: 38,

            created_at: '2024-02-15T10:30:00Z'
          },
          {
            id: 2,
            title: '极限计算技巧分享',
            content: '分享一些极限计算的常用技巧和方法，帮助大家更好地理解和掌握极限概念...',
            author: '小红',
            view_count: 189,
            comment_count: 8,
            like_count: 25,

            created_at: '2024-02-10T14:20:00Z'
          }
        ]
        
        // 模拟推荐课程
        recommendedCourses.value = [
          {
            id: 2,
            name: '线性代数',
            teacher: '李教授',
            note_count: 18,
            view_count: 890
          },
          {
            id: 3,
            name: '概率论与数理统计',
            teacher: '王教授',
            note_count: 22,
            view_count: 756
          }
        ]
      } finally {
        loading.value = false
      }
    }

    // 获取相关笔记
    const fetchRelatedNotes = async (courseId) => {
      try {
        const response = await noteAPI.getNotes({ course_id: Number(courseId), page_size: 10, sort_by: 'created_at', order: 'desc' })
        relatedNotes.value = response.notes || []
      } catch (error) {
        console.error('获取相关笔记失败:', error)
      }
    }

    // 获取推荐课程
    const fetchRecommendedCourses = async () => {
      try {
        const response = await courseAPI.getCourses({ page_size: 5 })
        recommendedCourses.value = response.courses || []
      } catch (error) {
        console.error('获取推荐课程失败:', error)
      }
    }

    // 加入课程
    const handleJoinCourse = async () => {
      try {
        const courseId = route.params.id
        await courseAPI.joinCourse(courseId)
        ElMessage.success('成功加入课程学习')
        isJoined.value = true
        // 可以在这里更新课程状态或重新获取课程信息
        await fetchCourseDetail()
      } catch (error) {
        console.error('加入课程失败:', error)
        if (error.response?.status === 409) {
          ElMessage.warning('您已加入该课程')
          isJoined.value = true
        } else if (error.response?.status === 400) {
          ElMessage.error('课程不可用')
        } else {
          ElMessage.error('加入课程失败，请稍后重试')
        }
      }
    }



    // 获取分类类型
    const getCategoryType = (category) => {
      const typeMap = {
        '数学': 'primary',
        '计算机': 'success',
        '语言': 'warning',
        '物理': 'info',
        '化学': 'danger',
        '经济': 'primary',
        '管理': 'success'
      }
      return typeMap[category] || 'info'
    }

    // 格式化时间
    const formatTime = (timeString) => {
      const date = new Date(timeString)
      const now = new Date()
      const diff = now - date
      const days = Math.floor(diff / (1000 * 60 * 60 * 24))
      
      if (days === 0) {
        const hours = Math.floor(diff / (1000 * 60 * 60))
        if (hours === 0) {
          const minutes = Math.floor(diff / (1000 * 60))
          return `${minutes}分钟前`
        }
        return `${hours}小时前`
      } else if (days < 7) {
        return `${days}天前`
      } else {
        return date.toLocaleDateString()
      }
    }

    // 格式化日期
    const formatDate = (timeString) => {
      return new Date(timeString).toLocaleDateString()
    }

    // 获取预览文本
    const getPreviewText = (content) => {
      if (!content || typeof content !== 'string') {
        return '暂无内容预览'
      }
      return content.replace(/<[^>]*>/g, '').substring(0, 100) + '...'
    }

    // 监听路由参数变化
    watch(() => route.params.id, (newId, oldId) => {
      if (newId !== oldId) {
        fetchCourseDetail()
      }
    })

    onMounted(() => {
      fetchCourseDetail()
    })

    return {
      loading,
      course,
      relatedNotes,
      recommendedCourses,
      isJoined,
      handleJoinCourse,

      getCategoryType,
      formatTime,
      formatDate,
      getPreviewText
    }
  }
}
</script>

<style scoped>
.page-header {
  background: var(--bg-secondary);
  padding: 20px 0;
  border-bottom: 1px solid var(--border-color);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.loading-container {
  padding: 40px 0;
}

.course-detail {
  padding: 30px 0;
}

.course-info-card {
  margin-bottom: 30px;
}

.course-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 40px;
}

.course-main-info {
  flex: 1;
}

.course-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 20px;
  line-height: 1.2;
}

.course-meta {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.course-teacher,
.course-time {
  display: flex;
  align-items: center;
  gap: 5px;
  color: var(--text-secondary);
  font-size: 1rem;
}

.course-description {
  font-size: 1.1rem;
  line-height: 1.6;
  color: var(--text-secondary);
}

.course-stats {
  display: flex;
  gap: 30px;
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 2rem;
  font-weight: 700;
  color: var(--primary-color);
  margin-bottom: 5px;
}

.stat-label {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.course-content-card,
.related-notes-card {
  margin-bottom: 30px;
}

.course-content {
  font-size: 1rem;
  line-height: 1.8;
  color: var(--text-primary);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.empty-notes {
  text-align: center;
  padding: 40px 20px;
}

.notes-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.note-item {
  cursor: pointer;
}

.note-header {
  margin-bottom: 15px;
}

.note-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.note-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.note-author {
  display: flex;
  align-items: center;
  gap: 5px;
  color: var(--primary-color);
  font-weight: 500;
}

.note-time {
  font-size: 0.9rem;
  color: var(--text-light);
}

.note-preview {
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 15px;
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
  gap: 15px;
  font-size: 0.9rem;
  color: var(--text-light);
}

.note-stats span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.note-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.course-sidebar {
  margin-bottom: 30px;
}

.sidebar-content {
  margin-bottom: 30px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--border-color);
}

.info-item:last-child {
  border-bottom: none;
}

.info-item label {
  font-weight: 500;
  color: var(--text-secondary);
}

.sidebar-actions {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.recommended-courses {
  margin-bottom: 30px;
}

.recommended-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.recommended-item {
  cursor: pointer;
  padding: 15px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  transition: border-color 0.3s;
}

.recommended-item:hover {
  border-color: var(--primary-color);
}

.recommended-info h5 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 5px;
}

.recommended-info p {
  color: var(--text-secondary);
  margin-bottom: 10px;
}

.recommended-stats {
  display: flex;
  gap: 15px;
  font-size: 0.9rem;
  color: var(--text-light);
}

.recommended-stats span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.error-state {
  padding: 60px 0;
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .course-header {
    flex-direction: column;
    gap: 20px;
  }
  
  .course-stats {
    justify-content: space-around;
    width: 100%;
  }
  
  .course-title {
    font-size: 2rem;
  }
  
  .course-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .el-row {
    flex-direction: column;
  }
  
  .el-col {
    width: 100%;
    margin-bottom: 20px;
  }
}
</style>