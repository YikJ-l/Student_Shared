import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Profile from '../views/Profile.vue'
import Courses from '../views/Courses.vue'
import CourseDetail from '../views/CourseDetail.vue'
import Notes from '../views/Notes.vue'
import NoteDetail from '../views/NoteDetail.vue'
import CreateNote from '../views/CreateNote.vue'
import EditNote from '../views/EditNote.vue'
import MyNotes from '../views/MyNotes.vue'
import Search from '../views/Search.vue'
import Admin from '../views/Admin.vue'
import Teacher from '../views/Teacher.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { hideSidebar: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { hideSidebar: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: { requiresAuth: true }
  },
  {
    path: '/courses',
    name: 'Courses',
    component: Courses
  },
  {
    path: '/courses/:id',
    name: 'CourseDetail',
    component: CourseDetail
  },
  {
    path: '/notes',
    name: 'Notes',
    component: Notes
  },
  {
    path: '/notes/:id',
    name: 'NoteDetail',
    component: NoteDetail
  },
  {
    path: '/create-note',
    name: 'CreateNote',
    component: CreateNote,
    meta: { requiresAuth: true }
  },
  {
    path: '/notes/:id/edit',
    name: 'EditNote',
    component: EditNote,
    meta: { requiresAuth: true }
  },
  {
    path: '/my-notes',
    name: 'MyNotes',
    component: MyNotes,
    meta: { requiresAuth: true }
  },
  {
    path: '/search',
    name: 'Search',
    component: Search
  },
  {
    path: '/admin',
    name: 'Admin',
    component: Admin,
    meta: { 
      requiresAuth: true,
      requiresRole: 'admin'
    }
  },
  {
    path: '/teacher',
    name: 'Teacher',
    component: Teacher,
    meta: { 
      requiresAuth: true,
      requiresRole: ['teacher', 'admin']
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
  const userRole = userInfo.role || 'student'
  
  // 检查是否需要登录
  if (to.meta.requiresAuth && !token) {
    next('/login')
    return
  }
  
  // 检查角色权限
  if (to.meta.requiresRole) {
    if (!token) {
      next('/login')
      return
    }
    
    const requiredRoles = Array.isArray(to.meta.requiresRole) ? to.meta.requiresRole : [to.meta.requiresRole]
    if (!requiredRoles.includes(userRole)) {
      // 权限不足，重定向到首页并显示错误信息
      next('/')
      return
    }
  }
  
  next()
})

export default router