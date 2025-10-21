<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-form-wrapper">
        <div class="login-header">
          <el-icon class="login-icon"><User /></el-icon>
          <h2>用户登录</h2>
          <p>欢迎回到高校课程互助平台</p>
        </div>
        
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          class="login-form"
          @submit.prevent="handleLogin"
        >
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名"
              size="large"
              prefix-icon="User"
            />
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              prefix-icon="Lock"
              show-password
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          
          <el-form-item>
            <div class="login-options">
              <el-checkbox v-model="rememberMe">记住我</el-checkbox>
              <el-link type="primary" underline="never">忘记密码？</el-link>
            </div>
          </el-form-item>
          
          <el-form-item>
            <el-button
              type="primary"
              size="large"
              class="login-button"
              :loading="loading"
              @click="handleLogin"
            >
              登录
            </el-button>
          </el-form-item>
          
          <el-form-item>
            <div class="register-link">
              还没有账号？
              <el-link type="primary" @click="$router.push('/register')">立即注册</el-link>
            </div>
          </el-form-item>
        </el-form>
      </div>
      
      <div class="login-banner">
        <div class="banner-content">
          <h3>加入我们的学习社区</h3>
          <p>与同学分享知识，共同进步</p>
          <ul class="feature-list">
            <li>
              <el-icon><Check /></el-icon>
              丰富的课程资源
            </li>
            <li>
              <el-icon><Check /></el-icon>
              优质的学习笔记
            </li>
            <li>
              <el-icon><Check /></el-icon>
              活跃的学习社区
            </li>
            <li>
              <el-icon><Check /></el-icon>
              便捷的知识搜索
            </li>
          </ul>
        </div>
        <div class="banner-image">
          <el-icon class="banner-icon"><Reading /></el-icon>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { userAPI } from '../api'
import { User, Lock, Check, Reading } from '@element-plus/icons-vue'

export default {
  name: 'Login',
  components: {
    User,
    Lock,
    Check,
    Reading
  },
  setup() {
    const router = useRouter()
    const loginFormRef = ref()
    const loading = ref(false)
    const rememberMe = ref(false)
    
    const loginForm = reactive({
      username: '',
      password: ''
    })
    
    const loginRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
      ]
    }
    
    const handleLogin = async () => {
      if (!loginFormRef.value) return
      
      try {
        await loginFormRef.value.validate()
        loading.value = true
        
        const response = await userAPI.login(loginForm)
        
        if (response.token) {
          // 保存登录信息
          localStorage.setItem('token', response.token)
          localStorage.setItem('userInfo', JSON.stringify(response.user))
          
          ElMessage.success('登录成功')
          
          // 跳转到首页或之前的页面
          const redirect = router.currentRoute.value.query.redirect || '/'
          router.push(redirect)
          
          // 重新加载页面以更新用户状态和权限
          setTimeout(() => {
            window.location.reload()
          }, 100)
        } else {
          ElMessage.error(response.message || '登录失败')
        }
      } catch (error) {
        if (error.response && error.response.status === 401) {
          ElMessage.error('用户名或密码错误')
        } else {
          ElMessage.error('登录失败，请稍后重试')
        }
      } finally {
        loading.value = false
      }
    }
    
    return {
      loginFormRef,
      loginForm,
      loginRules,
      loading,
      rememberMe,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, var(--bg-primary), var(--bg-secondary));
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-container {
  display: flex;
  max-width: 1000px;
  width: 100%;
  background: white;
  border-radius: 12px;
  box-shadow: var(--shadow-heavy);
  overflow: hidden;
}

.login-form-wrapper {
  flex: 1;
  padding: 60px 50px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-icon {
  font-size: 4rem;
  color: var(--primary-color);
  margin-bottom: 20px;
}

.login-header h2 {
  font-size: 2rem;
  color: var(--text-primary);
  margin-bottom: 10px;
}

.login-header p {
  color: var(--text-secondary);
  font-size: 1rem;
}

.login-form {
  width: 100%;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.login-button {
  width: 100%;
  height: 50px;
  font-size: 1.1rem;
}

.register-link {
  text-align: center;
  width: 100%;
  color: var(--text-secondary);
}

.login-banner {
  flex: 1;
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: white;
  padding: 60px 50px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: relative;
}

.banner-content h3 {
  font-size: 2rem;
  margin-bottom: 15px;
}

.banner-content p {
  font-size: 1.1rem;
  margin-bottom: 30px;
  opacity: 0.9;
}

.feature-list {
  list-style: none;
  padding: 0;
}

.feature-list li {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  font-size: 1rem;
}

.feature-list li .el-icon {
  margin-right: 10px;
  font-size: 1.2rem;
}

.banner-image {
  position: absolute;
  bottom: 30px;
  right: 30px;
  opacity: 0.2;
}

.banner-icon {
  font-size: 8rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-container {
    flex-direction: column;
    max-width: 400px;
  }
  
  .login-form-wrapper {
    padding: 40px 30px;
  }
  
  .login-banner {
    padding: 40px 30px;
    text-align: center;
  }
  
  .banner-image {
    position: static;
    margin-top: 20px;
    opacity: 0.3;
  }
  
  .banner-icon {
    font-size: 4rem;
  }
  
  .login-header h2 {
    font-size: 1.5rem;
  }
  
  .banner-content h3 {
    font-size: 1.5rem;
  }
}

@media (max-width: 480px) {
  .login-page {
    padding: 10px;
  }
  
  .login-form-wrapper {
    padding: 30px 20px;
  }
  
  .login-banner {
    padding: 30px 20px;
  }
}
</style>