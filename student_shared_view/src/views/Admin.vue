<template>
  <div class="admin-container">
    <!-- 页面标题 -->
    <div class="admin-header">
      <h1>管理员控制台</h1>
      <p>欢迎，{{ userInfo.nickname || userInfo.username }}！您当前的角色是：{{ userInfo.role }}</p>
    </div>

    <!-- 统计数据卡片 -->
    <el-row :gutter="20" class="admin-stats">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon"><User /></el-icon>
            <div class="stat-info">
              <h3>{{ stats.totalUsers }}</h3>
              <p>总用户数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon"><Reading /></el-icon>
            <div class="stat-info">
              <h3>{{ stats.totalCourses }}</h3>
              <p>总课程数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon"><Document /></el-icon>
            <div class="stat-info">
              <h3>{{ stats.totalNotes }}</h3>
              <p>总笔记数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon"><ChatDotRound /></el-icon>
            <div class="stat-info">
              <h3>{{ stats.totalComments }}</h3>
              <p>总评论数</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 管理功能选项卡 -->
    <el-card class="admin-tabs">
      <el-tabs v-model="activeTab" type="border-card">
        <!-- 用户管理 -->
        <el-tab-pane label="用户管理" name="users">
          <div class="tab-content">
            <div class="tab-header">
              <h3>用户管理</h3>
              <el-button type="primary" @click="refreshUsers">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
            
            <div class="search-bar">
              <el-input
                v-model="userSearch"
                placeholder="搜索用户名或邮箱"
                style="width: 300px; margin-right: 10px;"
                @input="searchUsers"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-select v-model="userRoleFilter" placeholder="筛选角色" style="width: 120px;" @change="filterUsers">
                <el-option label="全部" value="" />
                <el-option label="学生" value="student" />
                <el-option label="教师" value="teacher" />
                <el-option label="管理员" value="admin" />
              </el-select>
            </div>

            <el-table :data="filteredUsers" style="width: 100%" v-loading="usersLoading">
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="username" label="用户名" width="120" />
              <el-table-column prop="email" label="邮箱" width="200" />
              <el-table-column prop="nickname" label="昵称" width="120" />
              <el-table-column prop="role" label="角色" width="100">
                <template #default="scope">
                  <el-tag v-if="scope.row" :type="getRoleTagType(scope.row.role)">{{ getRoleDisplayName(scope.row.role) }}</el-tag>
                  <span v-else>-</span>
                </template>
              </el-table-column>
              <el-table-column prop="school" label="学校" width="150" />
              <el-table-column prop="created_at" label="注册时间" width="180">
                <template #default="scope">
                  {{ scope.row ? formatDate(scope.row.created_at) : '-' }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="200">
                <template #default="scope">
                  <template v-if="scope.row">
                    <el-button size="small" @click="editUser(scope.row)">编辑</el-button>
                    <el-button size="small" type="warning" @click="changeUserRole(scope.row)">角色</el-button>
                    <el-button size="small" type="danger" @click="deleteUser(scope.row)" :disabled="scope.row.role === 'admin'">删除</el-button>
                  </template>
                  <span v-else>-</span>
                </template>
              </el-table-column>
            </el-table>

            <el-pagination
              v-model:current-page="userCurrentPage"
              v-model:page-size="userPageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="userTotal"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleUserSizeChange"
              @current-change="handleUserCurrentChange"
              style="margin-top: 20px; justify-content: center;"
            />
          </div>
        </el-tab-pane>

        <!-- 课程管理 -->
        <el-tab-pane label="课程管理" name="courses">
          <div class="tab-content">
            <div class="tab-header">
              <h3>课程管理</h3>
              <div>
                <el-button type="success" @click="createCourse">
                  <el-icon><Plus /></el-icon>
                  新建课程
                </el-button>
                <el-button type="primary" @click="refreshCourses">
                  <el-icon><Refresh /></el-icon>
                  刷新
                </el-button>
              </div>
            </div>

            <div class="search-bar">
              <el-input
                v-model="courseSearch"
                placeholder="搜索课程名称或代码"
                style="width: 300px; margin-right: 10px;"
                @input="searchCourses"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-select v-model="courseStatusFilter" placeholder="筛选状态" style="width: 120px;" @change="filterCourses">
                <el-option label="全部" value="" />
                <el-option label="活跃" value="active" />
                <el-option label="已归档" value="archived" />
              </el-select>
            </div>

            <el-table :data="filteredCourses" style="width: 100%" v-loading="coursesLoading">
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="code" label="课程代码" width="120" />
              <el-table-column prop="name" label="课程名称" width="200" />
              <el-table-column prop="teacher" label="教师" width="120" />
              <el-table-column prop="school" label="学校" width="150" />
              <el-table-column prop="department" label="院系" width="120" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="scope">
                  <el-tag :type="scope.row.status === 'active' ? 'success' : 'info'">{{ scope.row.status === 'active' ? '活跃' : '已归档' }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="note_count" label="笔记数" width="100" />
              <el-table-column prop="student_count" label="学生数" width="100" />
              <el-table-column label="操作" width="200">
                <template #default="scope">
                  <template v-if="scope.row">
                    <el-button size="small" @click="editCourse(scope.row)">编辑</el-button>
                    <el-button size="small" type="warning" @click="toggleCourseStatus(scope.row)">
                      {{ scope.row.status === 'active' ? '归档' : '激活' }}
                    </el-button>
                    <el-button size="small" type="danger" @click="deleteCourse(scope.row)">删除</el-button>
                  </template>
                  <span v-else>-</span>
                </template>
              </el-table-column>
            </el-table>

            <el-pagination
              v-model:current-page="courseCurrentPage"
              v-model:page-size="coursePageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="courseTotal"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleCourseSizeChange"
              @current-change="handleCourseCurrentChange"
              style="margin-top: 20px; justify-content: center;"
            />
          </div>
        </el-tab-pane>

        <!-- 笔记管理 -->
        <el-tab-pane label="笔记管理" name="notes">
          <div class="tab-content">
            <div class="tab-header">
              <h3>笔记管理</h3>
              <el-button type="primary" @click="refreshNotes">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>

            <div class="search-bar">
              <el-input
                v-model="noteSearch"
                placeholder="搜索笔记标题或内容"
                style="width: 300px; margin-right: 10px;"
                @input="searchNotes"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-select v-model="noteStatusFilter" placeholder="筛选状态" style="width: 120px;" @change="filterNotes">
                <el-option label="全部" value="" />
                <el-option label="公开" value="public" />
                <el-option label="私有" value="private" />
                <el-option label="草稿" value="draft" />
              </el-select>
            </div>

            <el-table :data="filteredNotes" style="width: 100%" v-loading="notesLoading">
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="title" label="标题" width="200" />
              <el-table-column prop="author" label="作者" width="120" />
              <el-table-column prop="course_name" label="课程" width="150" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="scope">
                  <el-tag v-if="scope.row" :type="getNoteStatusTagType(scope.row.status)">{{ getNoteStatusDisplayName(scope.row.status) }}</el-tag>
                  <span v-else>-</span>
                </template>
              </el-table-column>
              <el-table-column prop="view_count" label="浏览量" width="100" />
              <el-table-column prop="like_count" label="点赞数" width="100" />
              <el-table-column prop="created_at" label="创建时间" width="180">
                <template #default="scope">
                  {{ scope.row ? formatDate(scope.row.created_at) : '-' }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="200">
                <template #default="scope">
                  <template v-if="scope.row">
                    <el-button size="small" @click="viewNote(scope.row)">查看</el-button>
                    <!-- <el-button size="small" type="warning" @click="toggleNoteStatus(scope.row)">状态</el-button> -->
                    <el-button size="small" type="danger" @click="deleteNote(scope.row)">删除</el-button>
                  </template>
                  <span v-else>-</span>
                </template>
              </el-table-column>
            </el-table>

            <el-pagination
              v-model:current-page="noteCurrentPage"
              v-model:page-size="notePageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="noteTotal"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleNoteSizeChange"
              @current-change="handleNoteCurrentChange"
              style="margin-top: 20px; justify-content: center;"
            />
          </div>
        </el-tab-pane>


      </el-tabs>
    </el-card>

    <!-- 用户编辑对话框 -->
    <el-dialog v-model="userEditDialogVisible" title="编辑用户" width="500px">
      <el-form :model="editingUser" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="editingUser.username" disabled />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="editingUser.email" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="editingUser.nickname" />
        </el-form-item>
        <el-form-item label="学校">
          <el-input v-model="editingUser.school" />
        </el-form-item>
        <el-form-item label="院系">
          <el-input v-model="editingUser.department" />
        </el-form-item>
        <el-form-item label="专业">
          <el-input v-model="editingUser.major" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="userEditDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveUser">保存</el-button>
      </template>
    </el-dialog>

    <!-- 角色修改对话框 -->
    <el-dialog v-model="roleChangeDialogVisible" title="修改用户角色" width="400px">
      <el-form :model="roleChangeForm" label-width="80px">
        <el-form-item label="用户">
          <el-input v-model="roleChangeForm.username" disabled />
        </el-form-item>
        <el-form-item label="当前角色">
          <el-tag :type="getRoleTagType(roleChangeForm.currentRole)">{{ getRoleDisplayName(roleChangeForm.currentRole) }}</el-tag>
        </el-form-item>
        <el-form-item label="新角色">
          <el-select v-model="roleChangeForm.newRole" placeholder="选择新角色">
            <el-option label="学生" value="student" />
            <el-option label="教师" value="teacher" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="roleChangeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveUserRole">确认修改</el-button>
      </template>
    </el-dialog>

    <!-- 新建课程对话框 -->
    <el-dialog v-model="courseCreateDialogVisible" :title="isEditMode ? '编辑课程' : '创建新课程'" width="600px">
      <el-alert
        title="提示"
        description="课程代码必须唯一，建议使用学校标准格式（如：CS101、MATH201等）"
        type="info"
        :closable="false"
        style="margin-bottom: 20px;"
      />
      <el-form :model="courseForm" :rules="courseRules" ref="courseFormRef" label-width="100px">
        <el-form-item label="课程代码" prop="code">
          <el-input v-model="courseForm.code" placeholder="请输入唯一的课程代码（如：CS101）" />
        </el-form-item>
        <el-form-item label="课程名称" prop="name">
          <el-input v-model="courseForm.name" placeholder="请输入课程名称" />
        </el-form-item>
        <el-form-item label="授课教师" prop="teacher">
          <el-select 
            v-model="courseForm.teacher" 
            filterable 
            remote 
            :remote-method="fetchTeacherOptions" 
            :loading="teacherLoading"
            placeholder="请输入关键词搜索教师">
            <el-option
              v-for="t in teacherOptions"
              :key="t.value"
              :label="t.label"
              :value="t.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="所属学校" prop="school">
          <el-input v-model="courseForm.school" placeholder="请输入所属学校" />
        </el-form-item>
        <el-form-item label="所属院系" prop="department">
          <el-input v-model="courseForm.department" placeholder="请输入所属院系" />
        </el-form-item>
        <el-form-item label="课程描述" prop="description">
          <el-input 
            v-model="courseForm.description" 
            type="textarea" 
            :rows="4" 
            placeholder="请输入课程描述" 
          />
        </el-form-item>
        <el-form-item label="课程内容">
          <el-input 
            v-model="courseForm.content" 
            type="textarea" 
            :rows="6" 
            placeholder="请输入详细的课程内容介绍（可选）" 
          />
        </el-form-item>
        <el-form-item label="课程状态">
          <el-radio-group v-model="courseForm.status">
            <el-radio value="active">活跃</el-radio>
            <el-radio value="archived">已归档</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="courseCreateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCourse" :loading="courseSubmitting">
          {{ isEditMode ? '更新课程' : '创建课程' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  User, 
  Reading, 
  Document, 
  ChatDotRound, 
  Setting, 
  Search, 
  Refresh, 
  Plus, 
  Delete, 
  Download, 
  FolderOpened 
} from '@element-plus/icons-vue'
import { roleUtils, userAPI, courseAPI, noteAPI } from '@/api/index.js'
import { useRouter } from 'vue-router'

const router = useRouter()

// 用户信息
const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

// 当前选中的标签页
const activeTab = ref('users')

// 统计数据
const stats = ref({
  totalUsers: 0,
  totalCourses: 0,
  totalNotes: 0,
  totalComments: 0,
  onlineUsers: 0
})



// 用户管理相关数据
const users = ref([])
const filteredUsers = ref([])
const usersLoading = ref(false)
const userSearch = ref('')
const userRoleFilter = ref('')
const userCurrentPage = ref(1)
const userPageSize = ref(10)
const userTotal = ref(0)

// 课程管理相关数据
const courses = ref([])
const filteredCourses = ref([])
const coursesLoading = ref(false)
const courseSearch = ref('')
const courseStatusFilter = ref('')
const courseCurrentPage = ref(1)
const coursePageSize = ref(10)
const courseTotal = ref(0)

// 笔记管理相关数据
const notes = ref([])
const filteredNotes = ref([])
const notesLoading = ref(false)
const noteSearch = ref('')
const noteStatusFilter = ref('')
const noteCurrentPage = ref(1)
const notePageSize = ref(10)
const noteTotal = ref(0)

// 对话框相关
const userEditDialogVisible = ref(false)
const roleChangeDialogVisible = ref(false)
const courseCreateDialogVisible = ref(false)
// 新增：教师选项列表与加载状态
const teacherOptions = ref([])
const teacherLoading = ref(false)
const teacherSearchTimer = ref(null)
const editingUser = ref({})
const roleChangeForm = ref({
  username: '',
  currentRole: '',
  newRole: ''
})

// 新建课程相关
const courseFormRef = ref(null)
const courseSubmitting = ref(false)
const isEditMode = ref(false)
const courseForm = ref({
  code: '',
  name: '',
  teacher: '',
  school: '',
  department: '',
  description: '',
  content: '',
  status: 'active'
})
const loadTeacherOptions = async () => {
  if (!roleUtils.isAdmin()) return
  teacherLoading.value = true
  try {
    const resp = await userAPI.getAllUsers({ role: 'teacher', page: 1, limit: 50 })
    const list = Array.isArray(resp.data) ? resp.data : Array.isArray(resp.users) ? resp.users : Array.isArray(resp) ? resp : []
    teacherOptions.value = list.map(u => ({
      value: u.username,
      label: (u.nickname ? `${u.nickname} (${u.username})` : u.username)
    }))
  } catch (e) {
    console.error('加载教师列表失败:', e)
    teacherOptions.value = []
  } finally {
    teacherLoading.value = false
  }
}

const fetchTeacherOptions = (query) => {
  if (!roleUtils.isAdmin()) return
  if (teacherSearchTimer.value) {
    clearTimeout(teacherSearchTimer.value)
  }
  teacherSearchTimer.value = setTimeout(async () => {
    teacherLoading.value = true
    try {
      const resp = await userAPI.getAllUsers({ role: 'teacher', search: (query || '').trim(), page: 1, limit: 20 })
      const list = Array.isArray(resp.data) ? resp.data : Array.isArray(resp.users) ? resp.users : Array.isArray(resp) ? resp : []
      teacherOptions.value = list.map(u => ({
        value: u.username,
        label: (u.nickname ? `${u.nickname} (${u.username})` : u.username)
      }))
    } catch (e) {
      console.error('搜索教师失败:', e)
      teacherOptions.value = []
    } finally {
      teacherLoading.value = false
    }
  }, 300)
}
// 课程表单验证规则
const courseRules = {
  code: [
    { required: true, message: '请输入课程代码', trigger: 'blur' },
    { min: 2, max: 20, message: '课程代码长度在 2 到 20 个字符', trigger: 'blur' },
    { pattern: /^[A-Za-z0-9]+$/, message: '课程代码只能包含字母和数字', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入课程名称', trigger: 'blur' },
    { min: 2, max: 100, message: '课程名称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  teacher: [
    { required: true, message: '请选择授课教师', trigger: 'change' }
  ],

  school: [
    { required: true, message: '请输入所属学校', trigger: 'blur' },
    { min: 2, max: 100, message: '学校名称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  department: [
    { required: true, message: '请输入所属院系', trigger: 'blur' },
    { min: 2, max: 100, message: '院系名称长度在 2 到 100 个字符', trigger: 'blur' }
  ]
}

// 检查管理员权限
const checkAdminPermission = () => {
  if (!roleUtils.isAdmin()) {
    ElMessage.error('权限不足，只有管理员可以访问此页面')
    router.push('/')
    return false
  }
  return true
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString || dateString === 'null' || dateString === 'undefined') return '-'
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return '-'
    return date.toLocaleString('zh-CN')
  } catch (error) {
    console.warn('日期格式化错误:', dateString, error)
    return '-'
  }
}

// 获取角色标签类型
const getRoleTagType = (role) => {
  if (!role) return 'info'
  const typeMap = {
    'admin': 'danger',
    'teacher': 'warning',
    'student': 'success'
  }
  return typeMap[role] || 'info'
}

// 获取角色显示名称
const getRoleDisplayName = (role) => {
  if (!role) return '未知'
  const nameMap = {
    'admin': '管理员',
    'teacher': '教师',
    'student': '学生'
  }
  return nameMap[role] || '未知'
}

// 获取笔记状态标签类型
const getNoteStatusTagType = (status) => {
  if (!status) return 'info'
  const typeMap = {
    'public': 'success',
    'private': 'warning',
    'draft': 'info'
  }
  return typeMap[status] || 'info'
}

// 获取笔记状态显示名称
const getNoteStatusDisplayName = (status) => {
  if (!status) return '未知'
  const nameMap = {
    'public': '公开',
    'private': '私有',
    'draft': '草稿'
  }
  return nameMap[status] || '未知'
}

// 加载统计数据
const loadStats = async () => {
  try {
    // 调用管理员统计API获取数据
    const statsRes = await userAPI.getAdminStats()
    
    stats.value = {
      totalUsers: statsRes.total_users || 0,
      totalCourses: statsRes.total_courses || 0,
      totalNotes: statsRes.total_notes || 0,
      totalComments: statsRes.total_comments || 0,
      onlineUsers: statsRes.online_users || 0
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
    ElMessage.error('加载统计数据失败')
    // 使用默认值
    stats.value = {
      totalUsers: 0,
      totalCourses: 0,
      totalNotes: 0,
      totalComments: 0,
      onlineUsers: 0
    }
  }
}

// 用户管理功能
const loadUsers = async () => {
  usersLoading.value = true
  try {
    const response = await userAPI.getAllUsers({
      page: userCurrentPage.value,
      limit: userPageSize.value,
      search: userSearch.value,
      role: userRoleFilter.value || undefined
    })
    
    // 确保数据是数组格式
    const userData = Array.isArray(response.data) ? response.data : 
                    Array.isArray(response) ? response : []
    
    users.value = userData
    filteredUsers.value = userData
    userTotal.value = response.total || userData.length
  } catch (error) {
    console.error('加载用户数据失败:', error)
    ElMessage.error('加载用户数据失败')
    users.value = []
    filteredUsers.value = []
    userTotal.value = 0
  } finally {
    usersLoading.value = false
  }
}

const refreshUsers = () => {
  userCurrentPage.value = 1
  loadUsers()
}

const searchUsers = () => {
  userCurrentPage.value = 1
  loadUsers()
}

const filterUsers = () => {
  userCurrentPage.value = 1
  loadUsers()
}

const editUser = (user) => {
  editingUser.value = { ...user }
  userEditDialogVisible.value = true
}

const saveUser = async () => {
  try {
    await userAPI.updateProfile(editingUser.value.id, {
      nickname: editingUser.value.nickname,
      email: editingUser.value.email,
      school: editingUser.value.school,
      department: editingUser.value.department,
      major: editingUser.value.major,
      introduction: editingUser.value.introduction
    })
    
    userEditDialogVisible.value = false
    ElMessage.success('用户信息更新成功')
    loadUsers() // 重新加载用户列表
  } catch (error) {
    console.error('更新用户信息失败:', error)
    ElMessage.error('更新用户信息失败')
  }
}

const changeUserRole = (user) => {
  roleChangeForm.value = {
    userId: user.id,
    username: user.username,
    currentRole: user.role,
    newRole: user.role
  }
  roleChangeDialogVisible.value = true
}

const saveUserRole = async () => {
  try {
    await userAPI.updateUserRole(roleChangeForm.value.userId, {
      role: roleChangeForm.value.newRole
    })
    
    roleChangeDialogVisible.value = false
    ElMessage.success('用户角色更新成功')
    loadUsers() // 重新加载用户列表
  } catch (error) {
    console.error('更新用户角色失败:', error)
    ElMessage.error('更新用户角色失败')
  }
}

const deleteUser = async (user) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${user.username}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await userAPI.deleteUser(user.id)
    
    ElMessage.success('用户删除成功')
    loadUsers() // 重新加载用户列表
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除用户失败:', error)
      ElMessage.error('删除用户失败')
    }
  }
}

const handleUserSizeChange = (size) => {
  userPageSize.value = size
  loadUsers()
}

const handleUserCurrentChange = (page) => {
  userCurrentPage.value = page
  loadUsers()
}

// 课程管理功能
const loadCourses = async () => {
  coursesLoading.value = true
  try {
    const response = await courseAPI.getCourses({
      page: courseCurrentPage.value,
      page_size: coursePageSize.value,
      search: courseSearch.value,
      status: courseStatusFilter.value || undefined
    })
    
    // 确保数据是数组格式
    const courseData = Array.isArray(response.courses) ? response.courses : 
                      Array.isArray(response.data) ? response.data : 
                      Array.isArray(response) ? response : []
    
    courses.value = courseData
    filteredCourses.value = courseData
    courseTotal.value = response.total || courseData.length
  } catch (error) {
    console.error('加载课程数据失败:', error)
    ElMessage.error('加载课程数据失败')
    courses.value = []
    filteredCourses.value = []
    courseTotal.value = 0
  } finally {
    coursesLoading.value = false
  }
}

const refreshCourses = () => {
  courseCurrentPage.value = 1
  loadCourses()
}

const searchCourses = () => {
  courseCurrentPage.value = 1
  loadCourses()
}

const filterCourses = () => {
  courseCurrentPage.value = 1
  loadCourses()
}

const createCourse = () => {
  // 重置表单
  courseForm.value = {
    code: '',
    name: '',
    teacher: '',
    category: '',
    school: '',
    department: '',
    description: '',
    content: '',
    status: 'active'
  }
  // 设置为创建模式
  isEditMode.value = false
  // 加载教师列表
  loadTeacherOptions()
  // 打开新建课程对话框
  courseCreateDialogVisible.value = true
}

// 提交课程（创建或编辑）
const submitCourse = async () => {
  if (!courseFormRef.value) return
  
  try {
    // 表单验证
    await courseFormRef.value.validate()
    
    courseSubmitting.value = true
    
    if (isEditMode.value) {
      // 编辑模式：更新课程
      const { id, ...updateData } = courseForm.value
      await courseAPI.updateCourse(id, updateData)
      ElMessage.success('课程更新成功')
    } else {
      // 创建模式：创建新课程
      await courseAPI.createCourse(courseForm.value)
      ElMessage.success('课程创建成功')
    }
    
    courseCreateDialogVisible.value = false
    
    // 刷新课程列表
    refreshCourses()
    
  } catch (error) {
    if (error.message && error.message.includes('验证失败')) {
      // 表单验证失败，不显示错误消息
      return
    }
    console.error('创建课程失败:', error)
    
    // 处理不同类型的错误
    let errorMessage = '创建课程失败'
    if (error.response) {
      const status = error.response.status
      const responseData = error.response.data
      
      if (status === 409) {
        errorMessage = responseData.error || '课程代码已存在，请使用其他代码'
      } else if (status === 403) {
        errorMessage = '权限不足：只有教师和管理员可以创建课程'
      } else if (status === 400) {
        errorMessage = responseData.error || '请求参数无效，请检查表单内容'
      } else {
        errorMessage = responseData.error || error.message || '创建课程失败'
      }
    } else if (error.message) {
      errorMessage = error.message
    }
    
    ElMessage.error(errorMessage)
  } finally {
    courseSubmitting.value = false
  }
}

const editCourse = (course) => {
  // 填充表单数据
  courseForm.value = {
    id: course.id,
    code: course.code,
    name: course.name,
    teacher: course.teacher,
    category: course.category,
    school: course.school,
    department: course.department,
    description: course.description,
    content: course.content || '',
    status: course.status
  }
  // 加载教师列表并确保当前教师可选
  loadTeacherOptions().then(() => {
    const exists = teacherOptions.value.some(opt => opt.value === courseForm.value.teacher)
    if (!exists && courseForm.value.teacher) {
      teacherOptions.value.unshift({ value: courseForm.value.teacher, label: courseForm.value.teacher })
    }
  })
  // 设置为编辑模式
  isEditMode.value = true
  courseCreateDialogVisible.value = true
}

const toggleCourseStatus = async (course) => {
  try {
    const newStatus = course.status === 'active' ? 'archived' : 'active'
    await courseAPI.updateCourse(course.id, { status: newStatus })
    
    ElMessage.success(`课程已${newStatus === 'active' ? '激活' : '归档'}`)
    refreshCourses()
  } catch (error) {
    console.error('修改课程状态失败:', error)
    ElMessage.error('修改课程状态失败')
  }
}

const deleteCourse = async (course) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除课程 "${course.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await courseAPI.deleteCourse(course.id)
    
    ElMessage.success('课程删除成功')
    refreshCourses()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除课程失败:', error)
      ElMessage.error('删除课程失败')
    }
  }
}

const handleCourseSizeChange = (size) => {
  coursePageSize.value = size
  loadCourses()
}

const handleCourseCurrentChange = (page) => {
  courseCurrentPage.value = page
  loadCourses()
}

// 笔记管理功能
const loadNotes = async () => {
  notesLoading.value = true
  try {
    const response = await noteAPI.getNotes({
      page: noteCurrentPage.value,
      page_size: notePageSize.value,
      search: noteSearch.value,
      status: noteStatusFilter.value || undefined
    })
    
    // 确保数据是数组格式
    const noteData = Array.isArray(response.notes) ? response.notes : 
                    Array.isArray(response.data) ? response.data : 
                    Array.isArray(response) ? response : []
    
    notes.value = noteData
    filteredNotes.value = noteData
    noteTotal.value = response.total || noteData.length
  } catch (error) {
    console.error('加载笔记数据失败:', error)
    ElMessage.error('加载笔记数据失败')
    notes.value = []
    filteredNotes.value = []
    noteTotal.value = 0
  } finally {
    notesLoading.value = false
  }
}

const refreshNotes = () => {
  noteCurrentPage.value = 1
  loadNotes()
}

const searchNotes = () => {
  noteCurrentPage.value = 1
  loadNotes()
}

const filterNotes = () => {
  noteCurrentPage.value = 1
  loadNotes()
}

const viewNote = (note) => {
  router.push(`/notes/${note.id}`)
}

// 管理员不能更改笔记状态 - 功能已禁用
// const toggleNoteStatus = async (note) => {
//   try {
//     const newStatus = note.status === 'public' ? 'private' : 'public'
//     await noteAPI.updateNote(note.id, { status: newStatus })
//     
//     note.status = newStatus
//     ElMessage.success(`笔记已${note.status === 'public' ? '公开' : '私有'}`)
//     refreshNotes()
//   } catch (error) {
//     console.error('切换笔记状态失败:', error)
//     ElMessage.error('修改笔记状态失败')
//   }
// }

const deleteNote = async (note) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除笔记 "${note.title}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await noteAPI.deleteNote(note.id)
    
    ElMessage.success('笔记删除成功')
    refreshNotes()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除笔记失败:', error)
      ElMessage.error('删除笔记失败')
    }
  }
}

const handleNoteSizeChange = (size) => {
  notePageSize.value = size
  loadNotes()
}

const handleNoteCurrentChange = (page) => {
  noteCurrentPage.value = page
  loadNotes()
}



onMounted(() => {
  if (checkAdminPermission()) {
    loadStats()
    loadUsers()
    loadCourses()
    loadNotes()
  }
})
</script>

<style scoped>
.admin-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.admin-header {
  margin-bottom: 20px;
  padding: 15px 0;
  border-bottom: 1px solid #ebeef5;
}

.admin-header h1 {
  color: #409eff;
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
}

.admin-header p {
  margin: 0;
  color: #606266;
  font-size: 14px;
}

.admin-stats {
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
  color: #409eff;
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

.admin-actions {
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

.recent-activities {
  margin-top: 20px;
}

.recent-activities h3 {
  margin: 0;
  color: #303133;
}
</style>