import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 10000
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.token = token
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    if (error.response) {
      const { status, data } = error.response
      if (status === 401) {
        localStorage.removeItem('token')
        window.location.href = '/login'
        ElMessage.error('登录已过期，请重新登录')
      } else {
        ElMessage.error(data.error || data.message || '请求失败')
      }
    } else {
      ElMessage.error('网络错误')
    }
    return Promise.reject(error)
  }
)

// 角色权限工具函数
export const roleUtils = {
  // 检查是否已登录
  isLoggedIn() {
    const token = localStorage.getItem('token')
    const userInfo = localStorage.getItem('userInfo')
    return !!(token && userInfo)
  },
  
  // 获取当前用户角色
  getCurrentUserRole() {
    const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
    return userInfo.role || 'student'
  },
  
  // 检查是否为管理员
  isAdmin() {
    return this.getCurrentUserRole() === 'admin'
  },
  
  // 检查是否为教师
  isTeacher() {
    return this.getCurrentUserRole() === 'teacher'
  },
  
  // 检查是否为学生
  isStudent() {
    return this.getCurrentUserRole() === 'student'
  },
  
  // 检查是否有指定角色权限
  hasRole(roles) {
    const currentRole = this.getCurrentUserRole()
    return Array.isArray(roles) ? roles.includes(currentRole) : roles === currentRole
  },
  
  // 检查是否有管理权限（管理员或教师）
  hasManagePermission() {
    return this.hasRole(['admin', 'teacher'])
  }
}

// 用户相关API
export const userAPI = {
  // 用户注册
  register(userData) {
    return api.post('/users/register', userData)
  },
  // 用户登录
  login(credentials) {
    return api.post('/users/login', credentials).then(response => {
      // 登录成功后保存用户信息到localStorage
      if (response.token) {
        localStorage.setItem('token', response.token)
        if (response.user) {
          localStorage.setItem('userInfo', JSON.stringify(response.user))
        }
      }
      return response
    })
  },
  // 获取用户信息
  getProfile() {
    return api.post('/users/profile').then(response => {
      // 更新本地用户信息
      if (response) {
        localStorage.setItem('userInfo', JSON.stringify(response))
      }
      return response
    })
  },
  // 更新用户信息
  updateProfile(userIdOrProfileData, profileData) {
    // 如果传入了两个参数且第一个是数字，说明是管理员更新其他用户的资料
    if (arguments.length === 2 && typeof userIdOrProfileData === 'number') {
      return api.post('/admin/users/update', { id: Number(userIdOrProfileData), ...(profileData || {}) })
    }
    // 否则是用户更新自己的资料，第一个参数就是profileData
    const actualProfileData = arguments.length === 1 ? userIdOrProfileData : profileData
    return api.post('/users/profile/update', actualProfileData).then(response => {
      // 更新本地用户信息（保存为用户对象，而不是整个响应）
      if (response) {
        localStorage.setItem('userInfo', JSON.stringify(response.user || response))
      }
      return response
    })
  },
  // 管理员获取所有用户
  getAllUsers(params = {}) {
    return api.post('/admin/users/list', params)
  },
  // 管理员更新用户角色
  updateUserRole(userId, roleData) {
    return api.post('/admin/users/update-role', { id: Number(userId), ...(roleData || {}) })
  },
  // 管理员删除用户
  deleteUser(userId) {
    return api.post('/admin/users/delete', { id: Number(userId) })
  },
  // 获取管理员统计数据
  getAdminStats() {
    return api.post('/admin/stats')
  },
  // 上传头像
  uploadAvatar(file) {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/upload/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },
  // 删除头像
  deleteAvatar(path) {
    return api.post('/upload/avatar/delete', { path })
  },
  // 登出
  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    window.location.href = '/login'
  }
}

// 课程相关API
export const courseAPI = {
  // 获取课程列表
  getCourses(params) {
    return api.post('/courses', params || {})
  },
  // 获取课程详情
  getCourseDetail(id) {
    return api.post('/courses/detail', { id: Number(id) })
  },
  // 创建课程（仅教师和管理员）
  createCourse(data) {
    if (!roleUtils.hasManagePermission()) {
      return Promise.reject(new Error('权限不足：只有教师和管理员可以创建课程'))
    }
    return api.post('/courses/create', data)
  },
  // 更新课程（仅教师和管理员）
  updateCourse(id, data) {
    if (!roleUtils.hasManagePermission()) {
      return Promise.reject(new Error('权限不足：只有教师和管理员可以更新课程'))
    }
    return api.post('/courses/update', { id: Number(id), ...(data || {}) })
  },
  // 删除课程（教师或管理员）
  deleteCourse(id) {
    if (!roleUtils.hasManagePermission()) {
      return Promise.reject(new Error('权限不足：只有教师或管理员可以删除课程'))
    }
    return api.post('/courses/delete', { id: Number(id) })
  },
  // 加入课程
  joinCourse(id) {
    return api.post('/courses/join', { id: Number(id) })
  },

  getMyCourses(params) {
    return api.post('/courses/my', params || {})
  },
  // 检查课程操作权限
  canCreateCourse() {
    return roleUtils.hasManagePermission()
  },
  canUpdateCourse() {
    return roleUtils.hasManagePermission()
  },
  canDeleteCourse() {
    return roleUtils.hasManagePermission()
  },
  // 获取最新课程（首页用）
  getLatestCourses() {
    return api.post('/home/latest-courses')
  }
}

// 笔记相关API
export const noteAPI = {
  // 获取笔记列表
  getNotes(params) {
    // 改为 POST，JSON 传筛选参数
    return api.post('/notes', params || {})
  },
  // 获取笔记详情
  getNoteDetail(id, params = {}) {
    // 改为 POST，JSON 传 id；忽略缓存控制
    const payload = { id: Number(id) }
    return api.post('/notes/detail', payload)
  },
  // 创建笔记（需要登录）
  createNote(data) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/notes/create', data)
  },
  // 更新笔记（仅作者或管理员）
  updateNote(id, data, noteAuthorId = null) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    // 如果提供了作者ID，检查权限
    if (noteAuthorId) {
      const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
      const isOwner = userInfo.id === noteAuthorId
      const isAdmin = roleUtils.isAdmin()
      if (!isOwner && !isAdmin) {
        return Promise.reject(new Error('权限不足：只能编辑自己的笔记或管理员权限'))
      }
    }
    // 改为 POST /notes/update，JSON 传 id + data
    return api.post('/notes/update', { id: Number(id), ...data })
  },
  // 删除笔记（仅作者或管理员）
  deleteNote(id, noteAuthorId = null) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    // 如果提供了作者ID，检查权限
    if (noteAuthorId) {
      const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
      const isOwner = userInfo.id === noteAuthorId
      const isAdmin = roleUtils.isAdmin()
      if (!isOwner && !isAdmin) {
        return Promise.reject(new Error('权限不足：只能删除自己的笔记或管理员权限'))
      }
    }
    // 改为 POST /notes/delete，JSON 传 id
    return api.post('/notes/delete', { id: Number(id) })
  },
  // 点赞笔记
  likeNote(id) {
    // 改为 POST /notes/like，JSON 传 id
    return api.post('/notes/like', { id: Number(id) })
  },
  // 取消点赞
  unlikeNote(id) {
    // 改为 POST /notes/unlike，JSON 传 id
    return api.post('/notes/unlike', { id: Number(id) })
  },
  // 收藏笔记
  favoriteNote(id) {
    // 改为 POST /notes/favorite，JSON 传 id
    return api.post('/notes/favorite', { id: Number(id) })
  },
  // 取消收藏
  unfavoriteNote(id) {
    // 改为 POST /notes/unfavorite，JSON 传 id
    return api.post('/notes/unfavorite', { id: Number(id) })
  },
  // 获取我的收藏笔记
  getMyFavorites(params = {}) {
    // 改为 POST /notes/favorites，可选分页
    return api.post('/notes/favorites', params)
  },
  // 获取我的点赞笔记
  getMyLikes(params = {}) {
    // 改为 POST /notes/likes，可选分页
    return api.post('/notes/likes', params)
  },
  // 获取作者其他笔记
  getNotesByAuthor(authorId) {
    // 复用 /notes POST 接口
    return api.post('/notes', { user_id: Number(authorId), page_size: 5 })
  },
  // 获取相关笔记（后端未提供专用接口，回退为列表）
  getRelatedNotes(id) {
    // 取最新公开笔记作为相关笔记
    return api.post('/notes', { page_size: 5, sort_by: 'created_at', order: 'desc' })
  },
  // 检查笔记操作权限
  canEditNote(noteAuthorId) {
    const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
    const isOwner = userInfo.id === noteAuthorId
    const isAdmin = roleUtils.isAdmin()
    return isOwner || isAdmin
  },
  canDeleteNote(noteAuthorId) {
    const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
    const isOwner = userInfo.id === noteAuthorId
    const isAdmin = roleUtils.isAdmin()
    return isOwner || isAdmin
  },
  // 获取热门笔记（首页用）
  getPopularNotes() {
    return api.post('/home/popular-notes')
  }
}

// 评论相关API
export const commentAPI = {
  // 获取评论列表
  getCommentsByNote(noteId, params) {
    // 改为 POST /comment/list，JSON 传 note_id 及分页
    const payload = { note_id: Number(noteId), ...(params || {}) }
    return api.post('/comment/list', payload)
  },
  // 创建评论（需要登录）
  createComment(data) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    // 改为 POST /comment/create
    return api.post('/comment/create', data)
  },
  // 更新评论（仅作者或管理员）
  updateComment(id, data, commentAuthorId = null) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    // 如果提供了作者ID，检查权限
    if (commentAuthorId) {
      const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
      const isOwner = userInfo.id === commentAuthorId
      const isAdmin = roleUtils.isAdmin()
      if (!isOwner && !isAdmin) {
        return Promise.reject(new Error('权限不足：只能编辑自己的评论或管理员权限'))
      }
    }
    // 改为 POST /comment/update，JSON 传 id + data
    return api.post('/comment/update', { id: Number(id), ...data })
  },
  // 删除评论（仅作者或管理员）
  deleteComment(id, commentAuthorId = null) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    // 如果提供了作者ID，检查权限
    if (commentAuthorId) {
      const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
      const isOwner = userInfo.id === commentAuthorId
      const isAdmin = roleUtils.isAdmin()
      if (!isOwner && !isAdmin) {
        return Promise.reject(new Error('权限不足：只能删除自己的评论或管理员权限'))
      }
    }
    // 改为 POST /comment/delete，JSON 传 id
    return api.post('/comment/delete', { id: Number(id) })
  },
  // 点赞评论
  likeComment(id) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    // 改为 POST /comment/like，JSON 传 id
    return api.post('/comment/like', { id: Number(id) })
  },
  // 取消点赞评论
  unlikeComment(id) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    // 改为 POST /comment/unlike，JSON 传 id
    return api.post('/comment/unlike', { id: Number(id) })
  },
  // 检查评论操作权限
  canEditComment(commentAuthorId) {
    const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
    const isOwner = userInfo.id === commentAuthorId
    const isAdmin = roleUtils.isAdmin()
    return isOwner || isAdmin
  },
  canDeleteComment(commentAuthorId) {
    const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
    const isOwner = userInfo.id === commentAuthorId
    const isAdmin = roleUtils.isAdmin()
    return isOwner || isAdmin
  }
}

// 搜索相关API
export const searchAPI = {
  // 搜索（旧）
  search(params) {
    // 兼容老调用，默认按笔记搜索
    return api.post('/search/notes', params || {})
  },
  // 简单笔记搜索（关键词 LIKE）
  searchNotes(params) {
    return api.post('/search/notes', params || {})
  }
}

// 首页相关API
export const homeAPI = {
  // 获取热门笔记
  getPopularNotes() {
    return api.post('/home/popular-notes')
  },
  // 获取最新课程
  getLatestCourses() {
    return api.post('/home/latest-courses')
  },
  // 获取公开统计数据
  getPublicStats() {
    return api.post('/home/stats')
  }
}

// AI相关API
export const aiAPI = {
  // 获取笔记AI元数据（摘要/关键词）
  getNoteMeta(id) {
    // 改为 POST /ai/notes/meta，JSON 传 id
    return api.post('/ai/notes/meta', { id: Number(id) })
  },
  // 生成摘要（需要登录；可传note_id或content）
  summarize(payload = {}) {
    return api.post('/ai/summarize', payload)
  }
}

export default api