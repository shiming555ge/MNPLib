import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

// 创建router
const router = useRouter()

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

// 验证token是否有效
const verifyToken = async () => {
  try {
    // 如果没有token，直接返回false
    if (!authToken.value) {
      return false;
    }
    
    const headers = {
      'Content-Type': 'application/json',
      ...getAuthHeader()
    };
    
    const response = await fetch('/api/auth/verify', {
      headers
    });
    
    // 如果token无效，自动清理
    if (!response.ok && response.status === 401) {
      console.warn('Token已过期，自动登出');
      logout();
    }
    
    return response.ok;
  } catch (err) {
    console.error('验证token失败:', err);
    return false;
  }
}

// 初始化时验证token（如果存在）
const initAuthValidation = async () => {
  if (authToken.value) {
    console.log('初始化验证token...');
    await verifyToken();
  }
}

// 立即执行初始化验证
initAuthValidation();

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
    verifyPasskeyModifiable,
    verifyToken
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
  verifyPasskeyModifiable,
  verifyToken
}
