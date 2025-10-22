<template>
  <div class="profile-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="container">
        <h1 class="page-title">个人中心</h1>
        <p class="page-subtitle">管理你的个人信息和学习记录</p>
      </div>
    </div>

    <div class="container">
      <el-row :gutter="20">
        <!-- 左侧个人信息 -->
        <el-col :span="8">
          <div class="profile-sidebar">
            <!-- 用户基本信息卡片 -->
            <el-card class="user-info-card">
              <div class="user-avatar-section">
                <el-avatar :size="80" :src="userInfo.avatar">
                  <el-icon><User /></el-icon>
                </el-avatar>
                <el-button link class="change-avatar-btn" @click="handleChangeAvatar">
                  更换头像
                </el-button>
                
                <!-- 头像上传组件 -->
                <AvatarUpload
                  v-model="showAvatarDialog"
                  :current-avatar="userInfo.avatar"
                  @success="handleAvatarSuccess"
                />
              </div>
              
              <div class="user-basic-info">
                <h3 class="username">{{ userInfo.username }}</h3>
                <p class="user-email">{{ userInfo.email }}</p>
                <p class="join-date">加入时间：{{ formatDate(userInfo.created_at) }}</p>
              </div>
              
              <div class="user-stats">
                <div class="stat-item">
                  <div class="stat-number">{{ userStats.notes }}</div>
                  <div class="stat-label">发布笔记</div>
                </div>
                <div class="stat-item">
                  <div class="stat-number">{{ userStats.comments }}</div>
                  <div class="stat-label">评论数</div>
                </div>
                <div class="stat-item">
                  <div class="stat-number">{{ userStats.likes }}</div>
                  <div class="stat-label">获得点赞</div>
                </div>
              </div>
            </el-card>

            <!-- 快捷操作 -->
            <el-card class="quick-actions-card">
              <template #header>
                <span>快捷操作</span>
              </template>
              <div class="quick-actions">
                <el-button type="primary" @click="$router.push('/create-note')">
                  <el-icon><EditPen /></el-icon>
                  发布笔记
                </el-button>
                <el-button @click="$router.push('/notes')">
                  <el-icon><Document /></el-icon>
                  浏览笔记
                </el-button>
                <el-button @click="$router.push('/courses')">
                  <el-icon><Reading /></el-icon>
                  查看课程
                </el-button>
              </div>
            </el-card>
          </div>
        </el-col>

        <!-- 右侧主要内容 -->
        <el-col :span="16">
          <div class="profile-content">
            <!-- 标签页 -->
            <el-tabs v-model="activeTab" class="profile-tabs">
              <!-- 个人信息编辑 -->
              <el-tab-pane label="个人信息" name="info">
                <el-card>
                  <template #header>
                    <span>编辑个人信息</span>
                  </template>
                  
                  <el-form
                    ref="profileFormRef"
                    :model="profileForm"
                    :rules="profileRules"
                    label-width="100px"
                  >
                    <el-form-item label="昵称" prop="nickname">
                      <el-input v-model="profileForm.nickname" placeholder="请输入昵称" />
                    </el-form-item>
                    
                    <el-form-item label="头像链接">
                      <el-input v-model="profileForm.avatar" placeholder="请输入头像链接" />
                    </el-form-item>
                    
                    <el-form-item label="个人简介">
                      <el-input
                        v-model="profileForm.introduction"
                        type="textarea"
                        :rows="4"
                        placeholder="介绍一下自己吧"
                      />
                    </el-form-item>
                    
                    <el-form-item label="学校">
                      <el-input v-model="profileForm.school" placeholder="请输入学校名称" />
                    </el-form-item>
                    
                    <el-form-item label="院系">
                      <el-input v-model="profileForm.department" placeholder="请输入院系" />
                    </el-form-item>
                    
                    <el-form-item label="专业">
                      <el-input v-model="profileForm.major" placeholder="请输入专业" />
                    </el-form-item>
                    
                    <el-form-item>
                      <el-button type="primary" :loading="updateLoading" @click="handleUpdateProfile">
                        保存修改
                      </el-button>
                      <el-button @click="resetProfileForm">重置</el-button>
                    </el-form-item>
                  </el-form>
                </el-card>
              </el-tab-pane>

              <!-- 我的课程 -->
              <el-tab-pane label="我的课程" name="courses">
                <el-card>
                  <template #header>
                    <div class="card-header">
                      <span>我的课程</span>
                      <div class="header-actions">
                        <el-button type="primary" size="small" @click="loadMyCourses">
                          <el-icon><Refresh /></el-icon>
                          刷新
                        </el-button>
                        <el-button type="primary" size="small" @click="$router.push('/courses')">
                          <el-icon><Plus /></el-icon>
                          加入课程
                        </el-button>
                      </div>
                    </div>
                  </template>
                  
                  <!-- 课程统计 -->
                  <div class="courses-overview">
                    <el-row :gutter="20">
                      <el-col :span="8">
                        <div class="overview-card">
                          <div class="overview-icon">
                            <el-icon><Reading /></el-icon>
                          </div>
                          <div class="overview-content">
                            <div class="overview-value">{{ myCourses.length }}</div>
                            <div class="overview-title">已加入课程</div>
                          </div>
                        </div>
                      </el-col>
                      <el-col :span="8">
                        <div class="overview-card">
                          <div class="overview-icon">
                            <el-icon><User /></el-icon>
                          </div>
                          <div class="overview-content">
                            <div class="overview-value">{{ myCourses.filter(c => c.role === 'teacher').length }}</div>
                            <div class="overview-title">教授课程</div>
                          </div>
                        </div>
                      </el-col>
                      <el-col :span="8">
                        <div class="overview-card">
                          <div class="overview-icon">
                            <el-icon><School /></el-icon>
                          </div>
                          <div class="overview-content">
                            <div class="overview-value">{{ myCourses.filter(c => c.role === 'student').length }}</div>
                            <div class="overview-title">学习课程</div>
                          </div>
                        </div>
                      </el-col>
                    </el-row>
                  </div>
                  
                  <!-- 课程列表 -->
                  <div class="courses-list">
                    <div v-if="coursesLoading" class="loading-courses">
                      <el-skeleton :rows="3" animated />
                    </div>
                    <div v-else-if="myCourses.length === 0" class="empty-courses">
                      <el-empty description="还没有加入任何课程">
                        <el-button type="primary" @click="$router.push('/courses')">去发现课程</el-button>
                      </el-empty>
                    </div>
                    <div v-else class="my-courses-list">
                      <div v-for="course in myCourses" :key="course.id" class="course-item">
                        <div class="course-info">
                          <div class="course-header">
                            <h5 class="course-title" @click="$router.push(`/courses/${course.id}`)">{{ course.name }}</h5>
                            <el-tag :type="course.role === 'teacher' ? 'warning' : 'success'" size="small">
                              {{ course.role === 'teacher' ? '教师' : '学生' }}
                            </el-tag>
                          </div>
                          <p class="course-code">课程代码：{{ course.code }}</p>
                          <p class="course-description">{{ course.description || '暂无描述' }}</p>
                          <p class="course-meta">
                            <span class="course-school">{{ course.school }}</span>
                            <span class="course-department">{{ course.department }}</span>
                            <span class="course-teacher">教师：{{ course.teacher }}</span>
                            <span class="course-semester">{{ course.semester }}</span>
                          </p>
                          <p class="course-join-date">加入时间：{{ formatDate(course.joined_at) }}</p>
                        </div>
                        <div class="course-stats">
                          <span><el-icon><User /></el-icon> {{ course.student_count }} 人</span>
                          <span><el-icon><Document /></el-icon> {{ course.note_count }} 笔记</span>
                        </div>
                        <div class="course-actions">
                          <el-button size="small" type="primary" @click="$router.push(`/courses/${course.id}`)">查看详情</el-button>
                          <el-button size="small" @click="$router.push(`/notes?course=${course.id}`)">查看笔记</el-button>
                        </div>
                      </div>
                    </div>
                  </div>
                </el-card>
              </el-tab-pane>

              <!-- 我的笔记 -->
              <el-tab-pane label="我的笔记" name="notes">
                <el-card>
                  <template #header>
                    <div class="card-header">
                      <span>我的笔记概览</span>
                      <div class="header-actions">
                        <el-button type="primary" size="small" @click="$router.push('/my-notes')">
                          <el-icon><Setting /></el-icon>
                          管理笔记
                        </el-button>
                        <el-button type="primary" size="small" @click="$router.push('/create-note')">
                          <el-icon><Plus /></el-icon>
                          新建笔记
                        </el-button>
                      </div>
                    </div>
                  </template>
                  
                  <!-- 笔记统计 -->
                  <div class="notes-overview">
                    <el-row :gutter="20">
                      <el-col :span="6">
                        <div class="overview-card">
                          <div class="overview-icon">
                            <el-icon><Document /></el-icon>
                          </div>
                          <div class="overview-content">
                            <div class="overview-value">{{ myNotes.length }}</div>
                            <div class="overview-title">总笔记数</div>
                          </div>
                        </div>
                      </el-col>
                      <el-col :span="6">
                        <div class="overview-card">
                          <div class="overview-icon">
                            <el-icon><View /></el-icon>
                          </div>
                          <div class="overview-content">
                            <div class="overview-value">{{ getTotalViews() }}</div>
                            <div class="overview-title">总浏览量</div>
                          </div>
                        </div>
                      </el-col>
                      <el-col :span="6">
                        <div class="overview-card">
                          <div class="overview-icon">
                            <el-icon><CirclePlus /></el-icon>
                          </div>
                          <div class="overview-content">
                            <div class="overview-value">{{ getTotalLikes() }}</div>
                            <div class="overview-title">总点赞数</div>
                          </div>
                        </div>
                      </el-col>
                      <el-col :span="6">
                        <div class="overview-card">
                          <div class="overview-icon">
                            <el-icon><ChatDotRound /></el-icon>
                          </div>
                          <div class="overview-content">
                            <div class="overview-value">{{ getTotalComments() }}</div>
                            <div class="overview-title">总评论数</div>
                          </div>
                        </div>
                      </el-col>
                    </el-row>
                  </div>
                  
                  <!-- 最近笔记 -->
                  <div class="recent-notes">
                    <h4>最近笔记</h4>
                    <div v-if="notesLoading" class="loading-notes">
                      <el-skeleton :rows="3" animated />
                    </div>
                    <div v-else-if="myNotes.length === 0" class="empty-notes">
                      <el-empty description="还没有发布任何笔记">
                        <el-button type="primary" @click="$router.push('/create-note')">发布第一篇笔记</el-button>
                      </el-empty>
                    </div>
                    <div v-else class="recent-notes-list">
                      <div v-for="note in myNotes.slice(0, 3)" :key="note.id" class="recent-note-item">
                        <div class="note-info">
                          <h5 class="note-title" @click="$router.push(`/notes/${note.id}`)">{{ note.title }}</h5>
                          <p class="note-meta">
                            <el-tag size="small" type="success">{{ note.course_name }}</el-tag>
                            <span class="note-date">{{ formatDate(note.created_at) }}</span>
                          </p>
                        </div>
                        <div class="note-stats">
                          <span><el-icon><View /></el-icon> {{ note.view_count }}</span>
                          <span><el-icon><Star /></el-icon> {{ note.like_count }}</span>
                        </div>
                      </div>
                      <div v-if="myNotes.length > 3" class="view-all">
                        <el-button type="text" @click="$router.push('/my-notes')">
                          查看全部 {{ myNotes.length }} 篇笔记 <el-icon><ArrowRight /></el-icon>
                        </el-button>
                      </div>
                    </div>
                  </div>
                </el-card>
              </el-tab-pane>

              <!-- 我的收藏 -->
              <!-- 我的点赞 -->
              <el-tab-pane label="我的点赞" name="likes">
                <el-card>
                  <template #header>
                    <div class="card-header">
                      <span>我的点赞</span>
                      <div class="header-actions">
                        <el-button type="primary" size="small" @click="loadMyLikes">
                          <el-icon><CirclePlus /></el-icon>
                          刷新
                        </el-button>
                      </div>
                    </div>
                  </template>
                  
                  <!-- 点赞统计 -->
                  <div class="likes-overview">
                    <el-row :gutter="20">
                      <el-col :span="8">
                        <div class="overview-card">
                          <div class="overview-icon">
                            <el-icon><CirclePlus /></el-icon>
                          </div>
                          <div class="overview-content">
                            <div class="overview-value">{{ myLikes.length }}</div>
                            <div class="overview-title">点赞笔记数</div>
                          </div>
                        </div>
                      </el-col>
                    </el-row>
                  </div>
                  
                  <!-- 点赞笔记列表 -->
                  <div class="likes-list">
                    <div v-if="likesLoading" class="loading-likes">
                      <el-skeleton :rows="3" animated />
                    </div>
                    <div v-else-if="myLikes.length === 0" class="empty-likes">
                      <el-empty description="还没有点赞任何笔记">
                        <el-button type="primary" @click="$router.push('/notes')">去发现好笔记</el-button>
                      </el-empty>
                    </div>
                    <div v-else class="likes-notes-list">
                      <div v-for="note in myLikes" :key="note.id" class="like-note-item">
                        <div class="note-info">
                          <h5 class="note-title" @click="$router.push(`/notes/${note.id}`)">{{ note.title }}</h5>
                          <p class="note-description">{{ note.description || '暂无描述' }}</p>
                          <p class="note-meta">
                            <el-tag size="small" type="success">{{ note.course_name }}</el-tag>
                            <span class="note-author">作者：{{ note.username }}</span>
                            <span class="note-date">点赞于：{{ formatDate(note.liked_at) }}</span>
                          </p>
                        </div>
                        <div class="note-stats">
                          <span><el-icon><View /></el-icon> {{ note.view_count }}</span>
                          <span><el-icon><Star /></el-icon> {{ note.like_count }}</span>
                          <span><el-icon><ChatDotRound /></el-icon> {{ note.comment_count }}</span>
                        </div>
                        <div class="note-actions">
                          <el-button size="small" type="primary" @click="$router.push(`/notes/${note.id}`)">查看</el-button>
                          <el-button size="small" type="danger" @click="unlikeNote(note.id)">取消点赞</el-button>
                        </div>
                      </div>
                    </div>
                  </div>
                </el-card>
              </el-tab-pane>

              <el-tab-pane label="我的收藏" name="favorites">
                <el-card>
                  <template #header>
                    <div class="card-header">
                      <span>我的收藏</span>
                      <div class="header-actions">
                        <el-button type="primary" size="small" @click="loadMyFavorites">
                          <el-icon><Star /></el-icon>
                          刷新
                        </el-button>
                      </div>
                    </div>
                  </template>
                  
                  <!-- 收藏统计 -->
                  <div class="favorites-overview">
                    <el-row :gutter="20">
                      <el-col :span="8">
                        <div class="overview-card">
                            <div class="overview-icon">
                              <el-icon><Star /></el-icon>
                            </div>
                            <div class="overview-content">
                              <div class="overview-value">{{ myFavorites.length }}</div>
                              <div class="overview-title">收藏笔记数</div>
                            </div>
                          </div>
                      </el-col>
                    </el-row>
                  </div>
                  
                  <!-- 收藏笔记列表 -->
                  <div class="favorites-list">
                    <div v-if="favoritesLoading" class="loading-favorites">
                      <el-skeleton :rows="3" animated />
                    </div>
                    <div v-else-if="myFavorites.length === 0" class="empty-favorites">
                      <el-empty description="还没有收藏任何笔记">
                        <el-button type="primary" @click="$router.push('/notes')">去发现好笔记</el-button>
                      </el-empty>
                    </div>
                    <div v-else class="favorites-notes-list">
                      <div v-for="note in myFavorites" :key="note.id" class="favorite-note-item">
                        <div class="note-info">
                          <h5 class="note-title" @click="$router.push(`/notes/${note.id}`)">{{ note.title }}</h5>
                          <p class="note-description">{{ note.description || '暂无描述' }}</p>
                          <p class="note-meta">
                            <el-tag size="small" type="success">{{ note.course_name }}</el-tag>
                            <span class="note-author">作者：{{ note.username }}</span>
                            <span class="note-date">收藏于：{{ formatDate(note.favorited_at) }}</span>
                          </p>
                        </div>
                        <div class="note-stats">
                          <span><el-icon><View /></el-icon> {{ note.view_count }}</span>
                          <span><el-icon><Star /></el-icon> {{ note.like_count }}</span>
                          <span><el-icon><ChatDotRound /></el-icon> {{ note.comment_count }}</span>
                        </div>
                        <div class="note-actions">
                          <el-button size="small" type="primary" @click="$router.push(`/notes/${note.id}`)">查看</el-button>
                          <el-button size="small" type="danger" @click="unfavoriteNote(note.id)">取消收藏</el-button>
                        </div>
                      </div>
                    </div>
                  </div>
                </el-card>
              </el-tab-pane>

            </el-tabs>
          </div>
        </el-col>
      </el-row>
    </div>


  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { userAPI, noteAPI, courseAPI } from '../api'
import { User, EditPen, Document, Reading, Plus, View, ChatDotRound, Star, Clock, Setting, ArrowRight, StarFilled, CirclePlus } from '@element-plus/icons-vue'
import AvatarUpload from '@/components/AvatarUpload.vue'

export default {
  name: 'Profile',
  components: {
      User,
      EditPen,
      Document,
      Reading,
      Plus,
      View,
      ChatDotRound,
      Star,
      Clock,
      Setting,
      ArrowRight,
      StarFilled,
      CirclePlus,
      AvatarUpload
    },
  setup() {
    const router = useRouter()
    const activeTab = ref('info')
    const updateLoading = ref(false)
    const profileFormRef = ref()
    
    // 用户信息
    const userInfo = ref({
      id: 1,
      username: '小明',
      email: 'xiaoming@example.com',
      avatar: '',
      bio: '热爱学习的大学生',
      school: '某某大学',
      major: '计算机科学与技术',
      email_verified: true,
      created_at: '2024-01-01T00:00:00Z'
    })
    
    // 用户统计
    const userStats = ref({
      notes: 12,
      comments: 45,
      likes: 128
    })
    

    
    // 我的笔记
    const myNotes = ref([])
    const notesLoading = ref(false)
    
    // 我的收藏
    const myFavorites = ref([])
    const favoritesLoading = ref(false)
    
    // 我的点赞
    const myLikes = ref([])
    const likesLoading = ref(false)
    
    // 我的课程
    const myCourses = ref([])
    const coursesLoading = ref(false)
    const coursesPage = ref(1)
    const coursesPageSize = ref(10)
    const coursesTotal = ref(0)
    
    // 个人信息表单
    const profileForm = reactive({
      nickname: '',
      avatar: '',
      school: '',
      department: '',
      major: '',
      introduction: ''
    })
    

    
    // 表单验证规则
    const profileRules = {
      nickname: [
        { required: true, message: '请输入昵称', trigger: 'blur' },
        { min: 2, max: 20, message: '昵称长度在 2 到 20 个字符', trigger: 'blur' }
      ]
    }
    

    
    // 初始化表单数据
    const initProfileForm = () => {
      Object.assign(profileForm, {
        nickname: userInfo.value.nickname || userInfo.value.username,
        avatar: userInfo.value.avatar || '',
        school: userInfo.value.school || '',
        department: userInfo.value.department || '',
        major: userInfo.value.major || '',
        introduction: userInfo.value.introduction || userInfo.value.bio || ''
      })
    }
    
    // 更新个人信息
    const handleUpdateProfile = async () => {
      if (!profileFormRef.value) return
      
      try {
        await profileFormRef.value.validate()
        updateLoading.value = true
        
        const response = await userAPI.updateProfile(profileForm)
        
        if (response.user) {
          // 更新本地用户信息
          Object.assign(userInfo.value, response.user)
          // 同时更新localStorage中的用户信息
          localStorage.setItem('userInfo', JSON.stringify(response.user))
          ElMessage.success('个人信息更新成功')
        } else {
          ElMessage.error(response.message || '更新失败')
        }
      } catch (error) {
        if (error.response && error.response.data) {
          ElMessage.error(error.response.data.error || '更新失败，请稍后重试')
        } else {
          ElMessage.error('更新失败，请稍后重试')
        }
      } finally {
        updateLoading.value = false
      }
    }
    
    // 重置表单
    const resetProfileForm = () => {
      initProfileForm()
    }
    

    
    // 更换头像
    const showAvatarDialog = ref(false)
    const handleChangeAvatar = () => {
      showAvatarDialog.value = true
    }
    
    // 头像上传成功处理
    const handleAvatarSuccess = async (avatarUrl) => {
      try {
        // 更新用户资料中的头像
        await userAPI.updateProfile({
          avatar: avatarUrl
        })
        
        // 更新本地用户信息
        userInfo.value.avatar = avatarUrl
        profileForm.avatar = avatarUrl
        
        ElMessage.success('头像更新成功')
        
        // 重新加载用户信息
        await loadUserInfo()
      } catch (error) {
        ElMessage.error('头像更新失败')
        console.error('Avatar update error:', error)
      }
    }
    
    // 删除笔记
    const handleDeleteNote = async (noteId) => {
      try {
        await ElMessageBox.confirm('确定要删除这篇笔记吗？', '确认删除', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        // 这里应该调用删除API
        // await noteAPI.deleteNote(noteId)
        
        // 从列表中移除
        const index = myNotes.value.findIndex(note => note.id === noteId)
        if (index > -1) {
          myNotes.value.splice(index, 1)
          userStats.value.notes--
        }
        
        ElMessage.success('删除成功')
      } catch (error) {
        // 用户取消删除
      }
    }
    

    
    // 加载用户信息
    const loadUserInfo = () => {
      const storedUserInfo = localStorage.getItem('userInfo')
      if (storedUserInfo) {
        try {
          const parsedUserInfo = JSON.parse(storedUserInfo)
          Object.assign(userInfo.value, parsedUserInfo)
        } catch (error) {
          console.error('解析用户信息失败:', error)
        }
      }
    }
    
    // 格式化日期
    const formatDate = (dateStr) => {
      const date = new Date(dateStr)
      return date.toLocaleDateString('zh-CN')
    }
    
    // 计算总浏览量
    const getTotalViews = () => {
      return myNotes.value.reduce((total, note) => total + (note.view_count || 0), 0)
    }
    
    // 计算总点赞数
    const getTotalLikes = () => {
      return myNotes.value.reduce((total, note) => total + (note.like_count || 0), 0)
    }
    
    // 计算总评论数
    const getTotalComments = () => {
      return myNotes.value.reduce((total, note) => total + (note.comment_count || 0), 0)
    }
    
    // 获取用户笔记数据
    const fetchMyNotes = async () => {
      try {
        notesLoading.value = true
        const token = localStorage.getItem('token')
        const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
        
        if (!token) {
          console.warn('用户未登录，无法获取笔记')
          return
        }
        
        if (!userInfo.id) {
          console.warn('用户信息不完整，无法获取笔记')
          return
        }
        
        const response = await noteAPI.getNotes({
          user_id: userInfo.id,
          status: 'all',
          page: 1,
          page_size: 10 // 只获取最近10篇笔记用于展示
        })
        
        myNotes.value = response.notes || []
        
        // 更新用户统计数据
        userStats.value.notes = myNotes.value.length
        userStats.value.likes = getTotalLikes()
        userStats.value.comments = getTotalComments()
        
      } catch (error) {
        console.error('获取笔记失败:', error)
        if (error.response) {
          console.error('API错误:', error.response.data)
        }
        // 如果是认证错误，清除本地存储并跳转到登录页
        if (error.response && error.response.status === 401) {
          localStorage.removeItem('token')
          localStorage.removeItem('userInfo')
          ElMessage.error('登录已过期，请重新登录')
          // 可以选择跳转到登录页
          // this.$router.push('/login')
        }
      } finally {
        notesLoading.value = false
      }
    }
    
    // 加载我的收藏
    const loadMyFavorites = async () => {
      favoritesLoading.value = true
      try {
        const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
        if (!userInfo.id) {
          ElMessage.error('请先登录')
          return
        }
        
        // 调用获取收藏笔记的API
        const response = await noteAPI.getMyFavorites()
        myFavorites.value = response.data || []
      } catch (error) {
        console.error('加载收藏失败:', error)
        ElMessage.error('加载收藏失败')
        myFavorites.value = []
      } finally {
        favoritesLoading.value = false
      }
    }
    
    // 取消收藏笔记
    const unfavoriteNote = async (noteId) => {
      try {
        await ElMessageBox.confirm(
          '确定要取消收藏这篇笔记吗？',
          '确认操作',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await noteAPI.unfavoriteNote(noteId)
        ElMessage.success('取消收藏成功')
        
        // 从列表中移除
        myFavorites.value = myFavorites.value.filter(note => note.id !== noteId)
      } catch (error) {
        if (error !== 'cancel') {
          console.error('取消收藏失败:', error)
          ElMessage.error('取消收藏失败')
        }
      }
    }
    
    // 加载我的点赞
    const loadMyLikes = async () => {
      likesLoading.value = true
      try {
        const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
        if (!userInfo.id) {
          ElMessage.error('请先登录')
          return
        }
        
        // 调用获取点赞笔记的API
        const response = await noteAPI.getMyLikes()
        myLikes.value = response.data || []
      } catch (error) {
        console.error('加载点赞失败:', error)
        ElMessage.error('加载点赞失败')
        myLikes.value = []
      } finally {
        likesLoading.value = false
      }
    }
    
    // 取消点赞笔记
    const unlikeNote = async (noteId) => {
      try {
        await ElMessageBox.confirm(
          '确定要取消点赞这篇笔记吗？',
          '确认操作',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )
        
        await noteAPI.unlikeNote(noteId)
        ElMessage.success('取消点赞成功')
        
        // 从列表中移除
        myLikes.value = myLikes.value.filter(note => note.id !== noteId)
      } catch (error) {
        if (error !== 'cancel') {
          console.error('取消点赞失败:', error)
          ElMessage.error('取消点赞失败')
        }
      }
    }
    
    // 加载我的课程
    const loadMyCourses = async () => {
      coursesLoading.value = true
      try {
        const response = await courseAPI.getMyCourses({
          page: coursesPage.value,
          page_size: coursesPageSize.value
        })
        
        myCourses.value = response.data || []
        coursesTotal.value = response.total || 0
      } catch (error) {
        console.error('加载我的课程失败:', error)
        ElMessage.error('加载我的课程失败')
        myCourses.value = []
      } finally {
        coursesLoading.value = false
      }
    }
    
    // 刷新我的课程
    const refreshMyCourses = () => {
      coursesPage.value = 1
      loadMyCourses()
    }
    
    // 课程分页变化
    const handleCoursesPageChange = (page) => {
      coursesPage.value = page
      loadMyCourses()
    }
    
    onMounted(() => {
      loadUserInfo()
      initProfileForm()
      fetchMyNotes()
      loadMyFavorites()
      loadMyLikes()
      loadMyCourses()
    })
    
    return {
      activeTab,
      updateLoading,
      profileFormRef,
      userInfo,
      userStats,
      myNotes,
      notesLoading,
      fetchMyNotes,
      myFavorites,
      favoritesLoading,
      loadMyFavorites,
      unfavoriteNote,
      myLikes,
      likesLoading,
      loadMyLikes,
      unlikeNote,
      myCourses,
      coursesLoading,
      coursesPage,
      coursesPageSize,
      coursesTotal,
      loadMyCourses,
      refreshMyCourses,
      handleCoursesPageChange,
      profileForm,
      profileRules,
      handleUpdateProfile,
      resetProfileForm,
      handleChangeAvatar,
      handleAvatarSuccess,
      showAvatarDialog,
      handleDeleteNote,
      formatDate,
      getTotalViews,
      getTotalLikes,
      getTotalComments
    }
  }
}
</script>

<style scoped>
.profile-sidebar {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.user-info-card {
  text-align: center;
}

.user-avatar-section {
  margin-bottom: 20px;
}

.change-avatar-btn {
  display: block;
  margin: 10px auto 0;
  color: var(--primary-color);
}

.user-basic-info {
  margin-bottom: 25px;
}

.username {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.user-email {
  color: var(--text-secondary);
  margin-bottom: 5px;
}

.join-date {
  font-size: 0.9rem;
  color: var(--text-light);
}

.user-stats {
  display: flex;
  justify-content: space-around;
  padding-top: 20px;
  border-top: 1px solid var(--border-color);
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--primary-color);
  margin-bottom: 5px;
}

.stat-label {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.quick-actions .el-button {
  justify-content: flex-start;
}

.profile-content {
  background: white;
  border-radius: 8px;
  box-shadow: var(--shadow-light);
}

.profile-tabs {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.notes-overview {
  margin-bottom: 30px;
}

.overview-card {
  display: flex;
  align-items: center;
  padding: 20px;
  background: var(--bg-secondary);
  border-radius: 8px;
  gap: 15px;
  transition: all 0.3s ease;
}

.overview-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.overview-icon {
  font-size: 2rem;
  color: var(--primary-color);
}

.overview-value {
  font-size: 1.8rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 5px;
}

.overview-title {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.recent-notes {
  margin-top: 20px;
}

.recent-notes h4 {
  margin-bottom: 20px;
  color: var(--text-primary);
  font-size: 1.2rem;
}

.recent-notes-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.recent-note-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background: var(--bg-secondary);
  border-radius: 8px;
  transition: all 0.3s ease;
}

.recent-note-item:hover {
  background: var(--bg-hover);
}

.recent-note-item .note-title {
  font-size: 1rem;
  margin-bottom: 5px;
}

.view-all {
  text-align: center;
  margin-top: 15px;
  padding: 10px;
  border-top: 1px solid var(--border-color);
}

.empty-notes {
  text-align: center;
  padding: 40px 20px;
}

.notes-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.note-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 15px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  transition: border-color 0.3s;
}

.note-item:hover {
  border-color: var(--primary-color);
}

.note-info {
  flex: 1;
}

.note-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
  cursor: pointer;
}

.note-title:hover {
  color: var(--primary-color);
}

.note-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.note-date {
  font-size: 0.9rem;
  color: var(--text-light);
}

.note-preview {
  color: var(--text-secondary);
  line-height: 1.5;
}

.note-actions {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 10px;
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

/* 收藏相关样式 */
.favorites-overview {
  margin-bottom: 20px;
}

.favorites-list {
  margin-top: 20px;
}

.favorites-notes-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.favorite-note-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 15px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  transition: border-color 0.3s;
}

.favorite-note-item:hover {
  border-color: var(--primary-color);
}

.favorite-note-item .note-info {
  flex: 1;
}

.favorite-note-item .note-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
  cursor: pointer;
}

.favorite-note-item .note-title:hover {
  color: var(--primary-color);
}

.favorite-note-item .note-description {
  color: var(--text-secondary);
  margin-bottom: 8px;
  line-height: 1.5;
}

.favorite-note-item .note-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.favorite-note-item .note-author,
.favorite-note-item .note-date {
  font-size: 0.9rem;
  color: var(--text-light);
}

.favorite-note-item .note-stats {
  display: flex;
  gap: 15px;
  font-size: 0.9rem;
  color: var(--text-light);
  margin-bottom: 10px;
}

.favorite-note-item .note-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

/* 点赞相关样式 */
.likes-overview {
  margin-bottom: 20px;
}

.likes-list {
  margin-top: 20px;
}

.likes-notes-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.like-note-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 15px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  transition: border-color 0.3s;
}

.like-note-item:hover {
  border-color: var(--primary-color);
}

.like-note-item .note-info {
  flex: 1;
}

.like-note-item .note-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
  cursor: pointer;
}

.like-note-item .note-title:hover {
  color: var(--primary-color);
}

.like-note-item .note-description {
  color: var(--text-secondary);
  margin-bottom: 8px;
  line-height: 1.5;
}

.like-note-item .note-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.like-note-item .note-author,
.like-note-item .note-date {
  font-size: 0.9rem;
  color: var(--text-light);
}

.like-note-item .note-stats {
  display: flex;
  gap: 15px;
  font-size: 0.9rem;
  color: var(--text-light);
  margin-bottom: 10px;
}

.like-note-item .note-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.loading-favorites,
.empty-favorites {
  text-align: center;
  padding: 40px 20px;
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



/* 响应式设计 */
@media (max-width: 768px) {
  .container .el-row {
    flex-direction: column;
  }
  
  .profile-sidebar {
    margin-bottom: 20px;
  }
  
  .note-item {
    flex-direction: column;
    gap: 15px;
  }
  
  .note-actions {
    align-items: flex-start;
    width: 100%;
  }
  

}
</style>