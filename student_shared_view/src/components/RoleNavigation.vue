<template>
  <div class="role-navigation">
    <!-- 角色标识 -->
    <div class="role-badge">
      <el-tag 
        :type="getRoleTagType()"
        size="small"
        effect="dark"
      >
        {{ getRoleDisplayName() }}
      </el-tag>
    </div>
  </div>
</template>

<script setup>
import { roleUtils } from '@/api/index.js'

// 获取角色显示名称
const getRoleDisplayName = () => {
  const role = roleUtils.getCurrentUserRole()
  const roleNames = {
    admin: '管理员',
    teacher: '教师',
    student: '学生'
  }
  return roleNames[role] || '未知角色'
}

// 获取角色标签类型
const getRoleTagType = () => {
  const role = roleUtils.getCurrentUserRole()
  const tagTypes = {
    admin: 'danger',
    teacher: 'warning',
    student: 'success'
  }
  return tagTypes[role] || 'info'
}
</script>

<style scoped>
.role-navigation {
  display: flex;
  align-items: center;
}

.role-badge {
  display: flex;
  align-items: center;
}
.role-badge .el-tag {
  white-space: nowrap;
}
</style>