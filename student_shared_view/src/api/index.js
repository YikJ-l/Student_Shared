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
    return api.get('/users/profile').then(response => {
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
      return api.put(`/admin/users/${userIdOrProfileData}`, profileData).then(response => {
        return response
      })
    }
    // 否则是用户更新自己的资料，第一个参数就是profileData
    const actualProfileData = arguments.length === 1 ? userIdOrProfileData : profileData
    return api.put('/users/profile', actualProfileData).then(response => {
      // 更新本地用户信息
      if (response) {
        localStorage.setItem('userInfo', JSON.stringify(response))
      }
      return response
    })
  },
  // 管理员获取所有用户
  getAllUsers(params = {}) {
    return api.get('/admin/users', { params })
  },
  // 管理员更新用户角色
  updateUserRole(userId, roleData) {
    return api.put(`/admin/users/${userId}/role`, roleData)
  },
  // 管理员删除用户
  deleteUser(userId) {
    return api.delete(`/admin/users/${userId}`)
  },
  // 获取管理员统计数据
  getAdminStats() {
    return api.get('/admin/stats')
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
    return api.delete('/upload/avatar', { data: { path } })
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
    return api.get('/courses', { params })
  },
  // 获取课程详情
  getCourseDetail(id) {
    return api.get(`/courses/${id}`)
  },
  // 创建课程（仅教师和管理员）
  createCourse(data) {
    if (!roleUtils.hasManagePermission()) {
      return Promise.reject(new Error('权限不足：只有教师和管理员可以创建课程'))
    }
    return api.post('/courses', data)
  },
  // 更新课程（仅教师和管理员）
  updateCourse(id, data) {
    if (!roleUtils.hasManagePermission()) {
      return Promise.reject(new Error('权限不足：只有教师和管理员可以更新课程'))
    }
    return api.put(`/courses/${id}`, data)
  },
  // 删除课程（仅管理员）
  deleteCourse(id) {
    if (!roleUtils.isAdmin()) {
      return Promise.reject(new Error('权限不足：只有管理员可以删除课程'))
    }
    return api.delete(`/courses/${id}`)
  },
  // 加入课程
  joinCourse(id) {
    return api.post(`/courses/${id}/join`)
  },

  getMyCourses(params) {
    return api.get('/courses/my', { params })
  },
  // 检查课程操作权限
  canCreateCourse() {
    return roleUtils.hasManagePermission()
  },
  canUpdateCourse() {
    return roleUtils.hasManagePermission()
  },
  canDeleteCourse() {
    return roleUtils.isAdmin()
  },
  // 获取最新课程（首页用）
  getLatestCourses() {
    return api.get('/home/latest-courses')
  }
}

// 笔记相关API
export const noteAPI = {
  // 获取笔记列表
  getNotes(params) {
    return api.get('/notes', { params })
  },
  // 获取笔记详情
  getNoteDetail(id, params = {}) {
    // 添加强制缓存控制参数
    const cacheControlParams = {
      ...params,
      t: Date.now(), // 时间戳
      _: Math.random().toString(36).substr(2, 9), // 随机字符串
      nocache: 'true' // 明确的无缓存标识
    }
    
    // 添加强制缓存控制头部
    const config = {
      params: cacheControlParams,
      headers: {
        'Cache-Control': 'no-cache, no-store, must-revalidate',
        'Pragma': 'no-cache',
        'Expires': '0',
        'If-Modified-Since': 'Mon, 26 Jul 1997 05:00:00 GMT'
      }
    }
    
    return api.get(`/notes/${id}`, config)
  },
  // 创建笔记（需要登录）
  createNote(data) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/notes', data)
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
    return api.put(`/notes/${id}`, data)
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
    return api.delete(`/notes/${id}`)
  },
  // 点赞笔记
  likeNote(id) {
    return api.post(`/notes/${id}/like`)
  },
  // 取消点赞
  unlikeNote(id) {
    return api.delete(`/notes/${id}/like`)
  },
  // 收藏笔记
  favoriteNote(id) {
    return api.post(`/notes/${id}/favorite`)
  },
  // 取消收藏
  unfavoriteNote(id) {
    return api.delete(`/notes/${id}/favorite`)
  },
  // 获取我的收藏笔记
  getMyFavorites() {
    return api.get('/notes/favorites')
  },
  // 获取我的点赞笔记
  getMyLikes() {
    return api.get('/notes/likes')
  },
  // 获取作者其他笔记
  getNotesByAuthor(authorId) {
    return api.get(`/notes?author_id=${authorId}`)
  },
  // 获取相关笔记
  getRelatedNotes(id) {
    return api.get(`/notes/${id}/related`)
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
    return api.get('/home/popular-notes')
  }
}

// 评论相关API
export const commentAPI = {
  // 获取评论列表
  getCommentsByNote(noteId, params) {
    return api.get(`/comment/note/${noteId}`, { params })
  },
  // 创建评论（需要登录）
  createComment(data) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/comment/', data)
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
    return api.put(`/comment/${id}`, data)
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
    return api.delete(`/comment/${id}`)
  },
  // 点赞评论
  likeComment(id) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post(`/comment/${id}/like`)
  },
  // 取消点赞评论
  unlikeComment(id) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.delete(`/comment/${id}/like`)
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
  // 搜索
  search(params) {
    return api.get('/search', { params })
  }
}

// 首页相关API
export const homeAPI = {
  // 获取热门笔记
  getPopularNotes() {
    return api.get('/home/popular-notes')
  },
  // 获取最新课程
  getLatestCourses() {
    return api.get('/home/latest-courses')
  },
  // 获取公开统计数据
  getPublicStats() {
    return api.get('/home/stats')
  }
}

export default api