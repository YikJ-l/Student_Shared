<template>
  <div class="teacher-container">
    <el-card class="teacher-header">
      <h1>教师工作台</h1>
      <p>欢迎，{{ userInfo.nickname || userInfo.username }}！您当前的角色是：{{ userInfo.role }}</p>
    </el-card>

    <el-row :gutter="20" class="teacher-stats">
      <el-col :span="8">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon"><Reading /></el-icon>
            <div class="stat-info">
              <h3>{{ stats.myCourses }}</h3>
              <p>我的课程</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon"><User /></el-icon>
            <div class="stat-info">
              <h3>{{ stats.totalStudents }}</h3>
              <p>学生总数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon"><Document /></el-icon>
            <div class="stat-info">
              <h3>{{ stats.courseNotes }}</h3>
              <p>课程笔记</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="teacher-actions">
      <el-col :span="12">
        <el-card>
          <template #header>
            <h3>课程管理</h3>
          </template>
          <div class="action-buttons">
            <el-button type="primary" @click="createCourse">
              <el-icon><Plus /></el-icon>
              创建新课程
            </el-button>
            <el-button type="success" @click="viewMyCourses">
              <el-icon><Reading /></el-icon>
              我的课程
            </el-button>
            <el-button type="info" @click="manageStudents">
              <el-icon><User /></el-icon>
              学生管理
            </el-button>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <h3>教学资源</h3>
          </template>
          <div class="action-buttons">
            <el-button type="warning" @click="uploadMaterials">
              <el-icon><Upload /></el-icon>
              上传资料
            </el-button>
            <el-button type="danger" @click="viewAssignments">
              <el-icon><EditPen /></el-icon>
              作业管理
            </el-button>
            <el-button type="primary" @click="viewAnalytics">
              <el-icon><DataAnalysis /></el-icon>
              学习分析
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="my-courses">
      <template #header>
        <div class="courses-header">
          <h3>我的课程</h3>
          <el-button type="primary" size="small" @click="createCourse">
            <el-icon><Plus /></el-icon>
            新建课程
          </el-button>
        </div>
      </template>
      <el-table :data="myCourses" style="width: 100%">
        <el-table-column prop="name" label="课程名称" />
        <el-table-column prop="code" label="课程代码" width="120" />
        <el-table-column prop="students" label="学生数" width="100">
          <template #default="scope">
            <el-tag>{{ scope.row.students }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === '进行中' ? 'success' : 'info'">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="editCourse(scope.row)">
              编辑
            </el-button>
            <el-button size="small" type="info" @click="viewCourseDetail(scope.row)">
              查看
            </el-button>
            <el-button size="small" type="danger" @click="deleteCourse(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card class="recent-activities">
      <template #header>
        <h3>最近活动</h3>
      </template>
      <el-timeline>
        <el-timeline-item
          v-for="activity in recentActivities"
          :key="activity.id"
          :timestamp="activity.time"
          placement="top"
        >
          <el-card>
            <h4>{{ activity.title }}</h4>
            <p>{{ activity.description }}</p>
          </el-card>
        </el-timeline-item>
      </el-timeline>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Reading, 
  User, 
  Document, 
  Plus, 
  Upload, 
  EditPen, 
  DataAnalysis 
} from '@element-plus/icons-vue'
import { roleUtils, courseAPI } from '@/api/index.js'
import { useRouter } from 'vue-router'

const router = useRouter()

// 用户信息
const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

// 统计数据
const stats = ref({
  myCourses: 0,
  totalStudents: 0,
  courseNotes: 0
})

// 我的课程
const myCourses = ref([
  {
    id: 1,
    name: 'Vue.js 进阶开发',
    code: 'CS301',
    students: 45,
    status: '进行中',
    created_at: '2024-01-10 09:00:00'
  },
  {
    id: 2,
    name: 'JavaScript 高级编程',
    code: 'CS302',
    students: 38,
    status: '进行中',
    created_at: '2024-01-08 14:30:00'
  },
  {
    id: 3,
    name: 'React 实战项目',
    code: 'CS303',
    students: 52,
    status: '已结束',
    created_at: '2023-12-15 10:15:00'
  }
])

// 最近活动
const recentActivities = ref([
  {
    id: 1,
    title: '新学生加入课程',
    description: '张三同学加入了《Vue.js 进阶开发》课程',
    time: '2024-01-15 10:30:00'
  },
  {
    id: 2,
    title: '课程资料更新',
    description: '更新了《JavaScript 高级编程》的第三章资料',
    time: '2024-01-15 09:15:00'
  },
  {
    id: 3,
    title: '作业提交',
    description: '收到15份《React 实战项目》的作业提交',
    time: '2024-01-14 16:45:00'
  }
])

// 检查教师权限
const checkTeacherPermission = () => {
  if (!roleUtils.hasManagePermission()) {
    ElMessage.error('权限不足，只有教师和管理员可以访问此页面')
    router.push('/')
    return false
  }
  return true
}

// 加载统计数据
const loadStats = async () => {
  try {
    // 这里应该调用实际的API获取统计数据
    stats.value = {
      myCourses: myCourses.value.length,
      totalStudents: myCourses.value.reduce((sum, course) => sum + course.students, 0),
      courseNotes: 156
    }
  } catch (error) {
    ElMessage.error('加载统计数据失败')
  }
}

// 课程管理功能
const createCourse = () => {
  ElMessage.info('创建课程功能开发中...')
  // router.push('/create-course')
}

const viewMyCourses = () => {
  router.push('/courses')
}

const manageStudents = () => {
  ElMessage.info('学生管理功能开发中...')
}

const editCourse = (course) => {
  ElMessage.info(`编辑课程：${course.name}`)
}

const viewCourseDetail = (course) => {
  router.push(`/courses/${course.id}`)
}

const deleteCourse = async (course) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除课程「${course.name}」吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 调用删除API
    await courseAPI.deleteCourse(course.id)
    
    // 从列表中移除
    const index = myCourses.value.findIndex(c => c.id === course.id)
    if (index > -1) {
      myCourses.value.splice(index, 1)
    }
    
    ElMessage.success('课程删除成功')
    loadStats() // 重新加载统计数据
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除课程失败')
    }
  }
}

// 教学资源功能
const uploadMaterials = () => {
  ElMessage.info('上传资料功能开发中...')
}

const viewAssignments = () => {
  ElMessage.info('作业管理功能开发中...')
}

const viewAnalytics = () => {
  ElMessage.info('学习分析功能开发中...')
}

onMounted(() => {
  if (checkTeacherPermission()) {
    loadStats()
  }
})
</script>

<style scoped>
.teacher-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.teacher-header {
  margin-bottom: 20px;
  text-align: center;
}

.teacher-header h1 {
  color: #67c23a;
  margin-bottom: 10px;
}

.teacher-stats {
  margin-bottom: 20px;
}

.stat-card {
  height: 120px;
}

.stat-content {
  display: flex;
  align-items: center;
  height: 100%;
}

.stat-icon {
  font-size: 40px;
  color: #67c23a;
  margin-right: 15px;
}

.stat-info h3 {
  font-size: 28px;
  font-weight: bold;
  margin: 0;
  color: #303133;
}

.stat-info p {
  margin: 5px 0 0 0;
  color: #909399;
  font-size: 14px;
}

.teacher-actions {
  margin-bottom: 20px;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.action-buttons .el-button {
  justify-content: flex-start;
}

.my-courses {
  margin-bottom: 20px;
}

.courses-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 0;
}

.courses-header h3 {
  margin: 0;
  color: #303133;
}

.recent-activities h3 {
  margin: 0;
  color: #303133;
}

.el-timeline-item .el-card {
  margin-bottom: 0;
}

.el-timeline-item .el-card h4 {
  margin: 0 0 10px 0;
  color: #303133;
}

.el-timeline-item .el-card p {
  margin: 0;
  color: #606266;
}
</style>