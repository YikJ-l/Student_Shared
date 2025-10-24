import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 10000
})

api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      // 后端鉴权从请求头读取 "token"，不是 "Authorization"
      config.headers['Token'] = token
      // 同时移除可能遗留的 Authorization 以避免误判
      delete config.headers['Authorization']
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    if (error.response) {
      const { status, data } = error.response
      if (status === 401) {
        ElMessage.error('登录已过期或未授权，请重新登录')
      } else if (data && data.error) {
        ElMessage.error(data.error)
      } else {
        ElMessage.error('请求失败，请稍后重试')
      }
    } else {
      ElMessage.error('网络错误，请检查您的连接')
    }
    return Promise.reject(error)
  }
)

export const roleUtils = {
  isLoggedIn() {
    const token = localStorage.getItem('token')
    return !!token
  },
  getCurrentUserRole() {
    const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
    return userInfo.role || 'guest'
  },
  isAdmin() {
    return this.getCurrentUserRole() === 'admin'
  },
  isTeacher() {
    return this.getCurrentUserRole() === 'teacher'
  },
  isStudent() {
    return this.getCurrentUserRole() === 'student'
  },
  hasRole(roles) {
    const role = this.getCurrentUserRole()
    return Array.isArray(roles) ? roles.includes(role) : role === roles
  },
  hasManagePermission() {
    return this.isAdmin() || this.isTeacher()
  }
}

export const userAPI = {
  register(userData) {
    return api.post('/users/register', userData)
  },
  login(credentials) {
    return api.post('/users/login', credentials).then(res => {
      if (res && res.token) {
        localStorage.setItem('token', res.token)
        localStorage.setItem('userInfo', JSON.stringify(res.user || {}))
      }
      return res
    })
  },
  getProfile() {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('未登录'))
    }
    return api.post('/users/profile')
  },
  updateProfile(userIdOrProfileData, profileData) {
    let payload = {}
    if (typeof userIdOrProfileData === 'object') {
      payload = userIdOrProfileData
    } else {
      payload = { id: Number(userIdOrProfileData), ...(profileData || {}) }
    }
    return api.post('/users/profile/update', payload)
  },
  getAllUsers(params = {}) {
    return api.post('/admin/users/list', params)
  },
  updateUserRole(userId, roleData) {
    return api.post('/admin/users/update-role', { id: Number(userId), ...(roleData || {}) })
  },
  deleteUser(userId) {
    return api.post('/admin/users/delete', { id: Number(userId) })
  },
  getAdminStats() {
    return api.post('/admin/stats')
  },
  uploadAvatar(file) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/upload/avatar', formData)
  },
  deleteAvatar(path) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/upload/avatar/delete', { path })
  },
  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    return Promise.resolve()
  }
}

export const courseAPI = {
  getCourses(params) {
    return api.post('/courses', params || {})
  },
  getCourseDetail(id) {
    return api.post('/courses/detail', { id: Number(id) })
  },
  createCourse(data) {
    return api.post('/courses/create', data)
  },
  updateCourse(id, data) {
    return api.post('/courses/update', { id: Number(id), ...(data || {}) })
  },
  deleteCourse(id) {
    return api.post('/courses/delete', { id: Number(id) })
  },
  joinCourse(id) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/courses/join', { id: Number(id) })
  },
  getMyCourses(params) {
    return api.post('/courses/my', params || {})
  },
  canCreateCourse() {
    return roleUtils.isTeacher() || roleUtils.isAdmin()
  },
  canUpdateCourse() {
    return roleUtils.isTeacher() || roleUtils.isAdmin()
  },
  canDeleteCourse() {
    return roleUtils.isAdmin()
  },
  getLatestCourses() {
    return api.post('/home/latest-courses')
  }
}

export const noteAPI = {
  getNotes(params) {
    return api.post('/notes', params || {})
  },
  getNoteDetail(id, params = {}) {
    return api.post('/notes/detail', { id: Number(id), ...(params || {}) })
  },
  createNote(data) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/notes/create', data)
  },
  updateNote(id, data, noteAuthorId = null) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    if (noteAuthorId != null) {
      const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
      const isOwner = userInfo.id === noteAuthorId
      const isAdmin = roleUtils.isAdmin()
      if (!isOwner && !isAdmin) {
        return Promise.reject(new Error('您没有权限更新此笔记'))
      }
    }
    return api.post('/notes/update', { id: Number(id), ...(data || {}) })
  },
  deleteNote(id, noteAuthorId = null) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    if (noteAuthorId != null) {
      const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
      const isOwner = userInfo.id === noteAuthorId
      const isAdmin = roleUtils.isAdmin()
      if (!isOwner && !isAdmin) {
        return Promise.reject(new Error('您没有权限删除此笔记'))
      }
    }
    return api.post('/notes/delete', { id: Number(id) })
  },
  likeNote(id) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/notes/like', { id: Number(id) })
  },
  unlikeNote(id) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/notes/unlike', { id: Number(id) })
  },
  favoriteNote(id) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/notes/favorite', { id: Number(id) })
  },
  unfavoriteNote(id) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/notes/unfavorite', { id: Number(id) })
  },
  getMyFavorites(params = {}) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/notes/favorites', params)
  },
  getMyLikes(params = {}) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/notes/likes', params)
  },
  getNotesByAuthor(authorId) {
    return api.post('/notes', { user_id: Number(authorId) })
  },
  getRelatedNotes(id) {
    return api.post('/notes', { page_size: 5 })
  },
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
  getPopularNotes() {
    return api.post('/home/popular-notes')
  }
}

export const commentAPI = {
  getCommentsByNote(noteId, params) {
    return api.post('/comment/list', { note_id: Number(noteId), ...(params || {}) })
  },
  createComment(data) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    return api.post('/comment/create', data)
  },
  updateComment(id, data, commentAuthorId = null) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    if (commentAuthorId != null) {
      const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
      const isOwner = userInfo.id === commentAuthorId
      const isAdmin = roleUtils.isAdmin()
      if (!isOwner && !isAdmin) {
        return Promise.reject(new Error('您没有权限更新此评论'))
      }
    }
    return api.post('/comment/update', { id: Number(id), ...(data || {}) })
  },
  deleteComment(id, commentAuthorId = null) {
    const token = localStorage.getItem('token')
    if (!token) {
      return Promise.reject(new Error('请先登录'))
    }
    if (commentAuthorId != null) {
      const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
      const isOwner = userInfo.id === commentAuthorId
      const isAdmin = roleUtils.isAdmin()
      if (!isOwner && !isAdmin) {
        return Promise.reject(new Error('您没有权限删除此评论'))
      }
    }
    return api.post('/comment/delete', { id: Number(id) })
  },
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
  },
  // 语义笔记搜索（向量检索 + 排序）
  searchNotesSemantic(params) {
    return api.post('/search/notes/semantic', params || {})
  },
  // 语义课程搜索（向量检索 + 排序）
  searchCoursesSemantic(params) {
    return api.post('/search/courses/semantic', params || {})
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