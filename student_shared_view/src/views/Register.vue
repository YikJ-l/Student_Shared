<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-form-wrapper">
        <div class="register-header">
          <el-icon class="register-icon"><UserFilled /></el-icon>
          <h2>用户注册</h2>
          <p>加入高校课程互助平台</p>
        </div>
        
        <el-form
          ref="registerFormRef"
          :model="registerForm"
          :rules="registerRules"
          class="register-form"
          @submit.prevent="handleRegister"
        >
          <el-form-item prop="username">
            <el-input
              v-model="registerForm.username"
              placeholder="请输入用户名"
              size="large"
              prefix-icon="User"
            />
          </el-form-item>
          
          <el-form-item prop="email">
            <el-input
              v-model="registerForm.email"
              placeholder="请输入邮箱"
              size="large"
              prefix-icon="Message"
            />
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="registerForm.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          
          <el-form-item prop="confirmPassword">
            <el-input
              v-model="registerForm.confirmPassword"
              type="password"
              placeholder="请确认密码"
              size="large"
              prefix-icon="Lock"
              show-password
              @keyup.enter="handleRegister"
            />
          </el-form-item>
          
          <el-form-item prop="agreement">
            <el-checkbox v-model="registerForm.agreement">
              我已阅读并同意
              <el-link type="primary" underline="never">《用户协议》</el-link>
              和
              <el-link type="primary" underline="never">《隐私政策》</el-link>
            </el-checkbox>
          </el-form-item>
          
          <el-form-item>
            <el-button
              type="primary"
              size="large"
              class="register-button"
              :loading="loading"
              @click="handleRegister"
            >
              注册
            </el-button>
          </el-form-item>
          
          <el-form-item>
            <div class="login-link">
              已有账号？
              <el-link type="primary" @click="$router.push('/login')">立即登录</el-link>
            </div>
          </el-form-item>
        </el-form>
      </div>
      
      <div class="register-banner">
        <div class="banner-content">
          <h3>开启你的学习之旅</h3>
          <p>与千万同学一起分享知识，共同成长</p>
          
          <div class="benefits">
            <div class="benefit-item">
              <el-icon><Reading /></el-icon>
              <div>
                <h4>海量课程资源</h4>
                <p>涵盖各个专业领域的优质课程</p>
              </div>
            </div>
            
            <div class="benefit-item">
              <el-icon><EditPen /></el-icon>
              <div>
                <h4>分享学习笔记</h4>
                <p>记录学习心得，帮助他人成长</p>
              </div>
            </div>
            
            <div class="benefit-item">
              <el-icon><ChatDotRound /></el-icon>
              <div>
                <h4>互动学习社区</h4>
                <p>与同学交流讨论，解决学习难题</p>
              </div>
            </div>
            
            <div class="benefit-item">
              <el-icon><TrophyBase /></el-icon>
              <div>
                <h4>提升学习效率</h4>
                <p>智能推荐，个性化学习体验</p>
              </div>
            </div>
          </div>
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
import { UserFilled, User, Message, Lock, Reading, EditPen, ChatDotRound, TrophyBase } from '@element-plus/icons-vue'

export default {
  name: 'Register',
  components: {
    UserFilled,
    User,
    Message,
    Lock,
    Reading,
    EditPen,
    ChatDotRound,
    TrophyBase
  },
  setup() {
    const router = useRouter()
    const registerFormRef = ref()
    const loading = ref(false)
    
    const registerForm = reactive({
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      agreement: false
    })
    
    // 自定义验证规则
    const validateConfirmPassword = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== registerForm.password) {
        callback(new Error('两次输入密码不一致'))
      } else {
        callback()
      }
    }
    
    const validateAgreement = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请阅读并同意用户协议'))
      } else {
        callback()
      }
    }
    
    const registerRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' },
        { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
        { pattern: /^[a-zA-Z0-9_\u4e00-\u9fa5]+$/, message: '用户名只能包含字母、数字、下划线和中文', trigger: 'blur' }
      ],
      email: [
        { required: true, message: '请输入邮箱地址', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' },
        { pattern: /^(?=.*[a-zA-Z])(?=.*\d)/, message: '密码必须包含字母和数字', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, validator: validateConfirmPassword, trigger: 'blur' }
      ],
      agreement: [
        { validator: validateAgreement, trigger: 'change' }
      ]
    }
    
    const handleRegister = async () => {
      if (!registerFormRef.value) return
      
      try {
        await registerFormRef.value.validate()
        loading.value = true
        
        const { confirmPassword, agreement, ...registerData } = registerForm
        const response = await userAPI.register(registerData)
        
        if (response.token) {
          ElMessage.success('注册成功，请登录')
          router.push('/login')
        } else {
          ElMessage.error(response.message || '注册失败')
        }
      } catch (error) {
        if (error.response && error.response.data) {
          ElMessage.error(error.response.data.message || '注册失败')
        } else {
          ElMessage.error('注册失败，请稍后重试')
        }
      } finally {
        loading.value = false
      }
    }
    
    return {
      registerFormRef,
      registerForm,
      registerRules,
      loading,
      handleRegister
    }
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  background: linear-gradient(135deg, var(--bg-primary), var(--bg-secondary));
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.register-container {
  display: flex;
  max-width: 1100px;
  width: 100%;
  background: white;
  border-radius: 12px;
  box-shadow: var(--shadow-heavy);
  overflow: hidden;
}

.register-form-wrapper {
  flex: 1;
  padding: 50px 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-icon {
  font-size: 4rem;
  color: var(--primary-color);
  margin-bottom: 20px;
}

.register-header h2 {
  font-size: 2rem;
  color: var(--text-primary);
  margin-bottom: 10px;
}

.register-header p {
  color: var(--text-secondary);
  font-size: 1rem;
}

.register-form {
  width: 100%;
}

.register-button {
  width: 100%;
  height: 50px;
  font-size: 1.1rem;
}

.login-link {
  text-align: center;
  width: 100%;
  color: var(--text-secondary);
}

.register-banner {
  flex: 1.2;
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: white;
  padding: 50px 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.banner-content h3 {
  font-size: 2.2rem;
  margin-bottom: 15px;
  font-weight: 600;
}

.banner-content > p {
  font-size: 1.2rem;
  margin-bottom: 40px;
  opacity: 0.9;
}

.benefits {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.benefit-item {
  display: flex;
  align-items: flex-start;
  gap: 15px;
}

.benefit-item .el-icon {
  font-size: 2rem;
  margin-top: 5px;
  opacity: 0.9;
}

.benefit-item h4 {
  font-size: 1.1rem;
  margin-bottom: 5px;
  font-weight: 600;
}

.benefit-item p {
  font-size: 0.95rem;
  opacity: 0.8;
  line-height: 1.4;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .register-container {
    flex-direction: column;
    max-width: 450px;
  }
  
  .register-form-wrapper {
    padding: 40px 30px;
  }
  
  .register-banner {
    padding: 40px 30px;
  }
  
  .register-header h2 {
    font-size: 1.5rem;
  }
  
  .banner-content h3 {
    font-size: 1.8rem;
  }
  
  .benefits {
    gap: 20px;
  }
  
  .benefit-item {
    gap: 12px;
  }
  
  .benefit-item .el-icon {
    font-size: 1.5rem;
  }
}

@media (max-width: 480px) {
  .register-page {
    padding: 10px;
  }
  
  .register-form-wrapper {
    padding: 30px 20px;
  }
  
  .register-banner {
    padding: 30px 20px;
  }
  
  .register-header {
    margin-bottom: 25px;
  }
  
  .banner-content h3 {
    font-size: 1.5rem;
  }
  
  .banner-content > p {
    font-size: 1rem;
    margin-bottom: 30px;
  }
}
</style>