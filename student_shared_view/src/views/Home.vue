<template>
  <div class="home">
    <!-- 英雄区域 -->
    <section class="hero-section">
      <div class="container">
        <div class="hero-content">
          <h1 class="hero-title">高校课程互助与笔记分享平台</h1>
          <p class="hero-subtitle">连接同学，分享知识，共同进步</p>
          <div class="hero-actions">
            <el-button type="primary" size="large" @click="$router.push('/courses')">
              <el-icon><Reading /></el-icon>
              浏览课程
            </el-button>
            <el-button size="large" @click="$router.push('/notes')">
              <el-icon><Document /></el-icon>
              查看笔记
            </el-button>
          </div>
        </div>
        <div class="hero-image">
          <el-icon class="hero-icon"><School /></el-icon>
        </div>
      </div>
    </section>

    <!-- 功能特色 -->
    <section class="features-section">
      <div class="container">
        <h2 class="section-title">平台特色</h2>
        <div class="features-grid">
          <div class="feature-card">
            <el-card shadow="hover">
              <div class="feature-icon">
                <el-icon><Reading /></el-icon>
              </div>
              <h3>课程资源</h3>
              <p>丰富的课程资源库，涵盖各个专业领域，帮助你找到所需的学习材料。</p>
            </el-card>
          </div>
          <div class="feature-card">
            <el-card shadow="hover">
              <div class="feature-icon">
                <el-icon><Document /></el-icon>
              </div>
              <h3>笔记分享</h3>
              <p>分享你的学习笔记，与同学交流学习心得，共同提高学习效率。</p>
            </el-card>
          </div>
          <div class="feature-card">
            <el-card shadow="hover">
              <div class="feature-icon">
                <el-icon><ChatDotRound /></el-icon>
              </div>
              <h3>互动交流</h3>
              <p>在线评论和讨论，与其他同学实时交流，解决学习中的疑问。</p>
            </el-card>
          </div>
          <div class="feature-card">
            <el-card shadow="hover">
              <div class="feature-icon">
                <el-icon><Search /></el-icon>
              </div>
              <h3>智能搜索</h3>
              <p>强大的搜索功能，快速找到你需要的课程和笔记内容。</p>
            </el-card>
          </div>
        </div>
      </div>
    </section>

    <!-- 最新内容 -->
    <section class="latest-section">
      <div class="container">
        <div class="latest-content">
          <div class="latest-courses">
            <h3>最新课程</h3>
            <el-card v-for="course in latestCourses" :key="course.id" class="course-item" shadow="hover">
              <div class="course-info">
                <h4>{{ course.name }}</h4>
                <p class="course-teacher">{{ course.teacher }}</p>
                <p class="course-desc">{{ course.description }}</p>
                <div class="course-meta">
                  <el-tag size="small">{{ course.department || course.category }}</el-tag>
                  <el-tag size="small" type="info">{{ course.note_count || 0 }} 笔记</el-tag>
                  <el-tag size="small" type="warning">{{ course.student_count || 0 }} 学生</el-tag>
                  <span class="course-time">{{ formatTime(course.created_at) }}</span>
                </div>
              </div>
            </el-card>
            <div class="view-more">
              <el-button link @click="$router.push('/courses')">查看更多课程 →</el-button>
            </div>
          </div>

          <div class="latest-notes">
            <h3>热门笔记</h3>
            <el-card v-for="note in latestNotes" :key="note.id" class="note-item" shadow="hover">
              <div class="note-info">
                <h4>{{ note.title }}</h4>
                <p class="note-author">{{ note.username || note.author }}</p>
                <p class="note-content">{{ note.description || (note.content && typeof note.content === 'string' ? note.content.substring(0, 100) + '...' : '暂无内容预览') }}</p>
                <div class="note-meta">
                  <el-tag size="small" type="success">{{ note.course_name }}</el-tag>
                  <el-tag size="small" type="info">{{ note.like_count || 0 }} 点赞</el-tag>
                  <span class="note-time">{{ formatTime(note.created_at) }}</span>
                </div>
              </div>
            </el-card>
            <div class="view-more">
              <el-button link @click="$router.push('/notes')">查看更多笔记 →</el-button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 统计数据 -->
    <section class="stats-section">
      <div class="container">
        <div class="stats-grid">
          <div class="stat-item">
            <div class="stat-number">{{ stats.courses }}</div>
            <div class="stat-label">课程数量</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">{{ stats.notes }}</div>
            <div class="stat-label">笔记数量</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">{{ stats.users }}</div>
            <div class="stat-label">注册用户</div>
          </div>
          <div class="stat-item">
            <div class="stat-number">{{ stats.comments }}</div>
            <div class="stat-label">评论数量</div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { courseAPI, noteAPI, homeAPI } from '../api'
import { Reading, Document, ChatDotRound, Search, School } from '@element-plus/icons-vue'

export default {
  name: 'Home',
  components: {
    Reading,
    Document,
    ChatDotRound,
    Search,
    School
  },
  setup() {
    const latestCourses = ref([])
    const latestNotes = ref([])
    const stats = ref({
      courses: 0,
      notes: 0,
      users: 0,
      comments: 0
    })

    // 获取最新课程
    const fetchLatestCourses = async () => {
      try {
        const response = await courseAPI.getLatestCourses()
        // 后端直接返回课程数组，不是包含courses字段的对象
        latestCourses.value = Array.isArray(response) ? response : []
        console.log('获取到的最新课程:', latestCourses.value)
        
        // 如果需要总数统计，可以单独调用获取总数的API
        try {
          const allCoursesResponse = await courseAPI.getCourses({ page: 1, page_size: 1 })
          stats.value.courses = allCoursesResponse.total || 0
        } catch (statsError) {
          console.warn('获取课程总数失败:', statsError)
        }
      } catch (error) {
        console.error('获取最新课程失败:', error)
      }
    }

    // 获取热门笔记
    const fetchLatestNotes = async () => {
      try {
        const response = await noteAPI.getPopularNotes()
        latestNotes.value = response.notes || []
        // 如果需要总数统计，可以单独调用获取总数的API
        const allNotesResponse = await noteAPI.getNotes({ page: 1, page_size: 1 })
        stats.value.notes = allNotesResponse.total || 0
      } catch (error) {
        console.error('获取热门笔记失败:', error)
      }
    }

    // 获取统计数据
    const fetchStats = async () => {
      try {
        const response = await homeAPI.getPublicStats()
        stats.value = {
          courses: response.courses || 0,
          notes: response.notes || 0,
          users: response.users || 0,
          comments: response.comments || 0
        }
      } catch (error) {
        console.error('获取统计数据失败:', error)
      }
    }

    // 格式化时间
    const formatTime = (timeStr) => {
      const date = new Date(timeStr)
      return date.toLocaleDateString('zh-CN')
    }

    onMounted(() => {
      fetchLatestCourses()
      fetchLatestNotes()
      fetchStats() // 获取真实统计数据
    })

    return {
      latestCourses,
      latestNotes,
      stats,
      formatTime
    }
  }
}
</script>

<style scoped>
.hero-section {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: white;
  padding: 100px 20px;
  min-height: 600px;
  display: flex;
  align-items: center;
  margin: 20px;
  border-radius: 24px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  position: relative;
  overflow: hidden;
}

.hero-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="50" cy="50" r="1" fill="%23ffffff" opacity="0.05"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>') repeat;
  pointer-events: none;
}

.hero-section .container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  position: relative;
  z-index: 1;
}

.hero-content {
  flex: 1;
  max-width: 600px;
  padding-right: 40px;
}

.hero-title {
  font-size: 3.5rem;
  font-weight: 800;
  margin-bottom: 24px;
  line-height: 1.1;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  background: linear-gradient(45deg, #ffffff, #f0f8ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: 1.4rem;
  margin-bottom: 48px;
  opacity: 0.95;
  line-height: 1.6;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.hero-actions {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
}

.hero-actions .el-button {
  padding: 16px 32px;
  font-size: 1.1rem;
  font-weight: 600;
  border-radius: 12px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.hero-actions .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
}

.hero-actions .el-button--primary {
  background: rgba(255, 255, 255, 0.2);
  border: 2px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
}

.hero-actions .el-button--primary:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
}

.hero-image {
  flex: 1;
  text-align: center;
  position: relative;
}

.hero-icon {
  font-size: 200px;
  opacity: 0.2;
}

.features-section {
  padding: 60px 20px;
}

.section-title {
  text-align: center;
  font-size: 2rem;
  margin-bottom: 30px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.feature-card .el-card {
  height: 100%;
}

.latest-section {
  padding: 60px 20px;
}

.latest-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
}

.course-item,
.note-item {
  margin-bottom: 16px;
}

.course-meta,
.note-meta {
  display: flex;
  gap: 8px;
  align-items: center;
}

.view-more {
  text-align: right;
}

.stats-section {
  padding: 60px 20px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.stat-item {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.stat-number {
  font-size: 2rem;
  font-weight: bold;
}

.stat-label {
  color: #666;
}
</style>