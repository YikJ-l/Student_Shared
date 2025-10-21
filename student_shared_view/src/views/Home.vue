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
        // 使用模拟数据作为后备
        latestCourses.value = [
          {
            id: 1,
            name: '高等数学',
            teacher: '张教授',
            description: '微积分基础理论与应用',
            category: '数学',
            created_at: new Date().toISOString()
          },
          {
            id: 2,
            name: '数据结构',
            teacher: '李教授',
            description: '计算机科学基础课程',
            category: '计算机',
            created_at: new Date().toISOString()
          },
          {
            id: 3,
            name: '大学英语',
            teacher: '王教授',
            description: '英语听说读写综合训练',
            category: '语言',
            created_at: new Date().toISOString()
          }
        ]
        stats.value.courses = 156
      }
    }

    // 获取热门笔记
    const fetchLatestNotes = async () => {
      try {
        const response = await noteAPI.getPopularNotes()
        latestNotes.value = response.notes || []
        // 如果需要总数统计，可以单独调用获取总数的API
        const allNotesResponse = await noteAPI.getNotes({ page: 1, limit: 1 })
        stats.value.notes = allNotesResponse.total || 0
      } catch (error) {
        console.error('获取热门笔记失败:', error)
        // 使用模拟数据作为后备
        latestNotes.value = [
          {
            id: 1,
            title: '微积分重点知识总结',
            author: '小明',
            content: '本笔记总结了微积分的重点知识点，包括导数、积分的基本概念和计算方法...',
            course_name: '高等数学',
            created_at: new Date().toISOString()
          },
          {
            id: 2,
            title: '二叉树遍历算法详解',
            author: '小红',
            content: '详细介绍了二叉树的前序、中序、后序遍历算法，包含代码实现和时间复杂度分析...',
            course_name: '数据结构',
            created_at: new Date().toISOString()
          },
          {
            id: 3,
            title: '英语语法要点整理',
            author: '小李',
            content: '整理了大学英语常用语法要点，包括时态、语态、从句等重要语法知识...',
            course_name: '大学英语',
            created_at: new Date().toISOString()
          }
        ]
        stats.value.notes = 892
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
        // 使用默认值作为后备
        stats.value = {
          courses: 0,
          notes: 0,
          users: 0,
          comments: 0
        }
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
  font-size: 16rem;
  opacity: 0.15;
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.1));
  animation: float 6s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-20px);
  }
}

.features-section {
  padding: 80px 0;
  background-color: var(--bg-secondary);
}

.section-title {
  text-align: center;
  font-size: 2.5rem;
  margin-bottom: 50px;
  color: var(--text-primary);
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 30px;
}

.feature-card .el-card {
  text-align: center;
  padding: 20px;
  height: 100%;
}

.feature-icon {
  font-size: 3rem;
  color: var(--primary-color);
  margin-bottom: 20px;
}

.feature-card h3 {
  font-size: 1.5rem;
  margin-bottom: 15px;
  color: var(--text-primary);
}

.feature-card p {
  color: var(--text-secondary);
  line-height: 1.6;
}

.latest-section {
  padding: 80px 0;
}

.latest-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
}

.latest-courses h3,
.latest-notes h3 {
  font-size: 1.8rem;
  margin-bottom: 30px;
  color: var(--text-primary);
  text-align: center;
}

.course-item,
.note-item {
  margin-bottom: 20px;
}

.course-info h4,
.note-info h4 {
  font-size: 1.2rem;
  margin-bottom: 8px;
  color: var(--text-primary);
}

.course-teacher,
.note-author {
  color: var(--primary-color);
  font-weight: 500;
  margin-bottom: 8px;
}

.course-desc,
.note-content {
  color: var(--text-secondary);
  margin-bottom: 15px;
  line-height: 1.5;
}

.course-meta,
.note-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.course-time,
.note-time {
  font-size: 0.9rem;
  color: var(--text-light);
}

.view-more {
  text-align: center;
  margin-top: 20px;
}

.stats-section {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: white;
  padding: 60px 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 40px;
  text-align: center;
}

.stat-number {
  font-size: 3rem;
  font-weight: 700;
  margin-bottom: 10px;
}

.stat-label {
  font-size: 1.1rem;
  opacity: 0.9;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .hero-section {
    margin: 10px;
    padding: 60px 16px;
    min-height: 500px;
    border-radius: 16px;
  }

  .hero-section .container {
    flex-direction: column;
    text-align: center;
  }

  .hero-content {
    margin-bottom: 40px;
    padding-right: 0;
    max-width: 100%;
  }

  .hero-title {
    font-size: 2.5rem;
    margin-bottom: 20px;
  }

  .hero-subtitle {
    font-size: 1.2rem;
    margin-bottom: 32px;
  }

  .hero-actions {
    justify-content: center;
    gap: 16px;
    flex-direction: column;
    align-items: center;
  }

  .hero-actions .el-button {
    padding: 14px 28px;
    font-size: 1rem;
    width: 200px;
  }

  .hero-icon {
    font-size: 10rem;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .latest-content {
    grid-template-columns: 1fr;
    gap: 30px;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>