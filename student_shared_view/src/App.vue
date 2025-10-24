<template>
  <div id="app">
    <el-container class="layout-root">
      <!-- 侧边导航栏 -->
      <el-aside v-if="!$route.meta.hideSidebar" :width="isCollapsed ? '64px' : '220px'" class="app-aside" :class="{ 'is-collapsed': isCollapsed, 'is-animating': isAnimating }">
        <div class="aside-logo"></div>
        <el-menu
          mode="vertical"
          :default-active="$route.path"
          class="side-nav-menu"
          router
          unique-opened
          :collapse="isCollapsed"
          :collapse-transition="false"
        >
          <el-menu-item index="/">
            <el-icon><House /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="/courses">
            <el-icon><Reading /></el-icon>
            <span>课程</span>
          </el-menu-item>
          <el-menu-item index="/notes">
            <el-icon><Document /></el-icon>
            <span>笔记</span>
          </el-menu-item>
          <el-menu-item index="/search">
            <el-icon><SearchIcon /></el-icon>
            <span>搜索</span>
          </el-menu-item>
          <el-menu-item v-if="isLoggedIn" index="/profile">
            <el-icon><User /></el-icon>
            <span>个人中心</span>
          </el-menu-item>
          <el-menu-item v-if="isAdmin" index="/admin">
            <el-icon><Setting /></el-icon>
            <span>管理员控制台</span>
          </el-menu-item>
          <el-menu-item v-if="isTeacher" index="/teacher">
            <el-icon><School /></el-icon>
            <span>教师工作台</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <!-- 右侧主内容区域 -->
      <el-container>
        <!-- 顶部仅保留用户操作，不再显示菜单 -->
        <el-header class="app-header" style="height: 64px;">
          <div class="header-content">
            <div class="left-tools">
        <el-button text class="collapse-toggle" @click="toggleAside">
          <el-icon v-if="isCollapsed"><Expand /></el-icon>
          <el-icon v-else><Fold /></el-icon>
        </el-button>
      </div>
      <div class="logo" @click="$router.push('/')">
        <el-icon class="logo-icon"><School /></el-icon>
        <span class="logo-text">高校课程互助与笔记分享平台</span>
      </div>
      <div class="spacer"></div>
      <div class="user-actions">
              <template v-if="isLoggedIn">
                <RoleNavigation class="role-nav" />
                
                <el-dropdown @command="handleCommand">
                  <span class="user-info">
                    <el-icon><User /></el-icon>
                    {{ userInfo.nickname || userInfo.username || '用户' }}
                    <el-icon class="el-icon--right"><arrow-down /></el-icon>
                  </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                      <!-- <el-dropdown-item command="create-note">发布笔记</el-dropdown-item> -->
                      <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
                
                <el-button 
                  v-if="isAdmin"
                  type="primary" 
                  @click="goToAdmin"
                  :class="{ 'admin-active': $route.path === '/admin' }"
                  class="admin-console-btn"
                >
                  <el-icon><Setting /></el-icon>
                  管理员控制台
                </el-button>
                <el-button 
                  v-if="isTeacher"
                  type="warning" 
                  @click="goToTeacher"
                  :class="{ 'teacher-active': $route.path === '/teacher' }"
                  class="teacher-console-btn"
                >
                  <el-icon><School /></el-icon>
                  教师工作台
                </el-button>
              </template>
              <template v-else>
                <el-button link @click="$router.push('/login')">登录</el-button>
                <el-button type="primary" @click="$router.push('/register')">注册</el-button>
              </template>
            </div>
          </div>
        </el-header>
        
        <el-main class="app-main">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </el-main>
      
      </el-container>
    </el-container>
  </div>
</template>

<script>
// 添加侧栏图标到组件
import { House, Reading, Document, Search as SearchIcon, Fold, Expand } from '@element-plus/icons-vue'
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { School, User, ArrowDown, Setting } from '@element-plus/icons-vue'
import RoleNavigation from './components/RoleNavigation.vue'
import { roleUtils } from '@/api/index.js'

export default {
  name: 'App',
  components: {
    School,
    User,
    ArrowDown,
    Setting,
    RoleNavigation,
    House,
    Reading,
    Document,
    SearchIcon,
    Fold,
    Expand
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const isLoggedIn = ref(false)
    const userInfo = ref({})
    const isCollapsed = ref(false)
    const isAnimating = ref(false)
    const toggleAside = () => {
      isAnimating.value = true
      isCollapsed.value = !isCollapsed.value
      setTimeout(() => { isAnimating.value = false }, 420)
    }
    
    // 管理员权限检查
    const isAdmin = computed(() => roleUtils.isAdmin())
    const isTeacher = computed(() => roleUtils.isTeacher())
    
    // 检查登录状态
    const checkLoginStatus = () => {
      const token = localStorage.getItem('token')
      const user = localStorage.getItem('userInfo')
      
      if (token && user) {
        isLoggedIn.value = true
        userInfo.value = JSON.parse(user)
      } else {
        isLoggedIn.value = false
        userInfo.value = {}
      }
    }
    
    // 处理用户菜单命令
    const handleCommand = (command) => {
      switch (command) {
        case 'profile':
          router.push('/profile')
          break
        case 'create-note':
          router.push('/notes/create')
          break
        case 'logout':
          logout()
          break
      }
    }
    
    // 退出登录
    const logout = () => {
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      isLoggedIn.value = false
      userInfo.value = {}
      router.push('/')
      ElMessage.success('退出登录成功')
    }
    
    // 管理员控制台导航
    const goToAdmin = () => {
      if (isAdmin.value) {
        router.push('/admin')
      } else {
        ElMessage.error('权限不足')
      }
    }
    const goToTeacher = () => {
      if (isTeacher.value) {
        router.push('/teacher')
      } else {
        ElMessage.error('仅教师可访问')
      }
    }
    
    // 监听路由变化，每次路由变化时检查登录状态
    watch(() => route.path, () => {
      checkLoginStatus()
    })
    
    // 监听localStorage变化（用于跨标签页同步）
    const handleStorageChange = (e) => {
      if (e.key === 'token' || e.key === 'userInfo') {
        checkLoginStatus()
      }
    }
    
    onMounted(() => {
      checkLoginStatus()
      // 添加storage事件监听器，用于跨标签页同步登录状态
      window.addEventListener('storage', handleStorageChange)
    })
    
    // 组件卸载时移除事件监听器
    onUnmounted(() => {
      window.removeEventListener('storage', handleStorageChange)
    })
    
    return {
      isLoggedIn,
      userInfo,
      isAdmin,
      isTeacher,
      handleCommand,
      goToAdmin,
      goToTeacher,
      isCollapsed,
      isAnimating,
      toggleAside
    }
  }
}
</script>

<style>
/* 布局 */
.layout-root { height: 100vh; }
.app-aside {
  background-image: linear-gradient(180deg,
    color-mix(in srgb, var(--primary-color) 10%, white),
    color-mix(in srgb, var(--secondary-color) 10%, white)
  );
  border-right: 1px solid color-mix(in srgb, var(--primary-color) 15%, white);
  backdrop-filter: blur(4px);
  position: sticky;
  top: 64px;
  height: calc(100vh - 64px);
  padding: 8px 6px;
  overflow: auto;
  transition: width 0.36s cubic-bezier(0.22, 1, 0.36, 1),
              padding 0.36s cubic-bezier(0.22, 1, 0.36, 1),
              background 0.36s ease,
              border-color 0.36s ease;
  will-change: width, padding;
}
.app-aside.is-animating { backdrop-filter: none; }
.aside-logo {
  padding: 6px 8px;
  margin-bottom: 4px;
  border-bottom: none;
}
/* 移除侧栏中的系统名称，仅保留菜单 */
.side-nav-menu { border-right: none; padding: 8px 0; transition: padding 0.36s cubic-bezier(0.22, 1, 0.36, 1), width 0.36s cubic-bezier(0.22, 1, 0.36, 1); }
.app-aside.is-collapsed .side-nav-menu { padding: 4px 0; }
.side-nav-menu .el-menu-item { 
  height: 44px; 
  line-height: 44px; 
  border-radius: 8px; 
  margin: 4px 8px; 
  transition: all .2s ease;
}
.side-nav-menu .el-menu-item .el-icon { margin-right: 8px; }
.side-nav-menu .el-menu-item.is-active { 
  background: color-mix(in srgb, var(--primary-color) 14%, white);
  color: var(--primary-color);
  box-shadow: inset 2px 0 0 var(--primary-color);
}
.side-nav-menu .el-menu-item:hover { background: color-mix(in srgb, var(--secondary-color) 12%, white); }
.spacer { flex: 1; }

/* 保留原有样式 */
.app-header {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  padding: 0;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0 32px;
  height: 100%;
  min-height: 64px;
  gap: 20px;
}

.logo {
  display: flex;
  align-items: center;
  color: white;
  font-size: 1.6rem;
  font-weight: 700;
  text-decoration: none;
  transition: all 0.3s ease;
  padding: 8px 0;
  width: 280px;
  flex-shrink: 0;
}

.logo:hover {
  transform: scale(1.05);
  text-shadow: 0 2px 8px rgba(255, 255, 255, 0.3);
}

.logo-icon {
  font-size: 2.2rem;
  margin-right: 12px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
}

.logo-text {
  font-size: 1.0rem;
  font-weight: 600;
}

.nav-menu {
  width: 400px;
  flex-shrink: 0;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: none;
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  transition: all 0.3s ease;
  padding: 0;
  display: flex;
}

.el-menu--horizontal.el-menu {
  border-bottom: none !important;
}

.nav-menu .el-menu-item {
  border-bottom: 2px solid transparent;
  margin: 0;
  padding: 0;
  border-radius: 0;
  background-color: rebeccapurple;
  font-weight: 500;
  font-size: 0.95rem;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  background: transparent;
  backdrop-filter: none;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 44px;
}



.nav-menu .el-menu-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0.12));
  opacity: 0;
  transition: opacity 0.3s ease;
}

.nav-menu .el-menu-item:hover::before {
  opacity: 1;
}

.nav-menu .el-menu-item:hover {
  /* background: rgba(243, 231, 231, 0.08); */
  background-color: rgba(142, 10, 10, 0.626);
  border-bottom-color: rgba(255, 255, 255, 0.6);
  color: white;
  transform: none;
}

.nav-menu .el-menu-item:first-child {
}

.nav-menu .el-menu-item:last-child {
}

.nav-menu .el-menu-item.is-active {
  background: rgba(255, 255, 255, 0.1);
  border-bottom-color: rgba(255, 255, 255, 0.8);
  transform: none;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  width: 280px;
  justify-content: flex-end;
  flex-shrink: 0;
}

.role-nav {
  margin-right: 4px;
}

/* 管理员控制台按钮 */
.admin-console-btn {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: white;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.admin-console-btn:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-1px);
}

.admin-console-btn:active {
  transform: translateY(0);
}

/* 教师工作台按钮，仅教师显示 */
.teacher-console-btn {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: white;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}
.teacher-console-btn:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-1px);
}
.teacher-console-btn:active {
  transform: translateY(0);
}
.role-nav .el-button {
  color: white;
  border-color: rgba(255, 255, 255, 0.4);
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 8px 20px;
  font-weight: 500;
  transition: all 0.3s ease;
  backdrop-filter: blur(5px);
}

.role-nav .el-button:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.6);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.role-nav .el-tag {
  margin-left: 10px;
}

.user-info {
  display: flex;
  align-items: center;
  color: white;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 6px;
  transition: background-color 0.3s;
  margin-right: 16px;
  font-weight: 500;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  gap: 6px;
  white-space: nowrap;
}

.user-info:hover {
  background: rgba(255, 255, 255, 0.1);
}

.user-avatar {
  margin-right: 10px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.user-avatar:hover {
  border-color: rgba(255, 255, 255, 0.6);
  transform: scale(1.05);
}

.user-dropdown .el-dropdown-link {
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  padding: 8px 16px;
  border-radius: 20px;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(5px);
}

.user-dropdown .el-dropdown-link:hover {
  color: white;
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.auth-buttons {
  display: flex;
  gap: 12px;
}

.auth-buttons .el-button {
  color: white;
  border-color: rgba(255, 255, 255, 0.4);
  border-radius: 20px;
  padding: 8px 20px;
  font-weight: 500;
  transition: all 0.3s ease;
  backdrop-filter: blur(5px);
}

.auth-buttons .el-button:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.6);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.auth-buttons .el-button--primary {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.5);
}

.auth-buttons .el-button--primary:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.7);
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.user-actions .el-button--text {
  color: white;
}

.user-actions .el-button--text:hover {
  background: rgba(255, 255, 255, 0.1);
}

.app-main {
  flex: 1;
  padding: 64px 0 0 0;
  background-color: var(--bg-primary);
}

.app-footer {
  background-color: var(--text-primary);
  color: white;
  text-align: center;
  padding: 20px 0;
}

.footer-content p {
  margin: 5px 0;
}

.footer-links a {
  color: white;
  text-decoration: none;
  margin: 0 5px;
}

.footer-links a:hover {
  color: var(--primary-light);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-main {
    padding-top: 120px;
  }

  .header-content {
    padding: 0 20px;
    min-height: 56px;
  }

  .nav-menu {
    display: none;
  }

  .logo {
    font-size: 1.3rem;
    padding: 6px 0;
  }

  .logo-icon {
    font-size: 1.8rem;
    margin-right: 8px;
  }

  .user-actions {
    gap: 12px;
  }

  .role-nav .el-button {
    padding: 6px 16px;
    font-size: 0.9rem;
    border-radius: 16px;
  }

  .auth-buttons .el-button {
    padding: 6px 16px;
    font-size: 0.9rem;
    border-radius: 16px;
  }

  .user-dropdown .el-dropdown-link {
    padding: 6px 12px;
    border-radius: 16px;
    font-size: 0.9rem;
  }

  .user-info {
    margin-right: 8px;
    font-size: 0.9rem;
  }
}
.collapse-toggle { color: #fff; }
.collapse-toggle .el-icon { font-size: 18px; }
.left-tools { display: flex; align-items: center; }
/* 菜单文字平滑出现/隐藏 */
.side-nav-menu .el-menu-item span {
  display: inline-block !important;
  overflow: hidden;
  white-space: nowrap;
  opacity: 1;
  max-width: 160px;
  transition: opacity 0.36s cubic-bezier(0.22, 1, 0.36, 1),
              max-width 0.36s cubic-bezier(0.22, 1, 0.36, 1),
              transform 0.36s cubic-bezier(0.22, 1, 0.36, 1);
}
.app-aside.is-collapsed .side-nav-menu .el-menu-item span {
  opacity: 0;
  max-width: 0;
  transform: translateX(6px);
}
.app-aside.is-collapsed .side-nav-menu .el-menu-item .el-icon { margin-right: 0; }
.app-aside.is-collapsed .side-nav-menu .el-menu-item { justify-content: center; }
</style>
