import { ref, computed } from 'vue'

// 创建全局认证状态
const authToken = ref(localStorage.getItem('authToken') || '')
const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || 'null'))
const isAuthenticated = computed(() => !!authToken.value)

// 登录函数
const login = (token, userData = null) => {
  authToken.value = token
  userInfo.value = userData
  localStorage.setItem('authToken', token)
  if (userData) {
    localStorage.setItem('userInfo', JSON.stringify(userData))
  }
  // 触发自定义事件，通知其他组件
  window.dispatchEvent(new CustomEvent('auth-changed', { detail: { isAuthenticated: true, userInfo: userData } }))
}

// 登出函数
const logout = () => {
  authToken.value = ''
  userInfo.value = null
  localStorage.removeItem('authToken')
  localStorage.removeItem('userInfo')
  // 触发自定义事件，通知其他组件
  window.dispatchEvent(new CustomEvent('auth-changed', { detail: { isAuthenticated: false, userInfo: null } }))
}

// 获取认证头
const getAuthHeader = () => {
  if (!authToken.value) return {}
  return {
    'Authorization': `Bearer ${authToken.value}`
  }
}

// 获取用户信息
const getUserInfo = () => userInfo.value

// 检查是否是管理员（Extends字段为空）
const isAdmin = computed(() => {
  const info = getUserInfo();
  return info && info.extends === null;
})

// 验证是否有权限修改passkey
const verifyPasskeyModifiable = async () => {
  try {
    const headers = {
      'Content-Type': 'application/json',
      ...getAuthHeader()
    };
    
    const response = await fetch('/api/auth/verify-passkey-modifiable', {
      headers
    });
    
    return response.ok;
  } catch (err) {
    console.error('验证权限失败:', err);
    return false;
  }
}

// 可组合函数
export function useAuth() {
  return {
    authToken,
    userInfo,
    isAuthenticated,
    isAdmin,
    login,
    logout,
    getAuthHeader,
    getUserInfo,
    verifyPasskeyModifiable
  }
}

// 全局实例
export const auth = {
  authToken,
  userInfo,
  isAuthenticated,
  isAdmin,
  login,
  logout,
  getAuthHeader,
  getUserInfo,
  verifyPasskeyModifiable
}
