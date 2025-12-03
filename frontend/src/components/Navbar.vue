<template>
    <nav class="navbar navbar-expand-lg text-white navbar-bg-primary sticky-top z-2">
      <div class="container-fluid">
        <a class="navbar-brand text-light p-3" href="#">
            <span class="fs-1 fw-bold">{{ t("navbar.brand_for_short") }}</span>
            <span class="fs-6 p-3">{{ t("navbar.brand_for_full") }}</span>
        </a>
        
        <!-- 移动端汉堡菜单按钮 -->
        <button class="navbar-toggler border-light" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        
        <!-- 右侧功能区域 -->
        <div class="collapse navbar-collapse" id="navbarNav">
          <!-- 小屏设备：垂直排列的功能按钮 -->
          <div class="d-lg-none">
            <div class="d-flex flex-column py-2">
              <RouterLink class="nav-link text-light py-2 border-bottom" aria-current="page" to="/" active-class="active">{{ t("navbar.home") }}</RouterLink>
              <RouterLink class="nav-link text-light py-2 border-bottom" to="/browse" active-class="active">{{ t("navbar.browse") }}</RouterLink>
              <RouterLink class="nav-link text-light py-2 border-bottom" to="/query" active-class="active">{{ t("navbar.query") }}</RouterLink>
              <RouterLink class="nav-link text-light py-2 border-bottom" to="/about" active-class="active">{{ t("navbar.about") }}</RouterLink>
            </div>
            
            <!-- 移动端认证按钮 -->
            <div class="mt-2" v-if="!isAuthenticated">
              <button class="btn btn-outline-light w-100" type="button" @click="openAuthModal">
                <i class="bi bi-key me-1"></i>
                {{ t('auth.login') }}
              </button>
            </div>
            
            <!-- 移动端用户菜单（已登录状态） -->
            <div class="mt-2" v-else>
              <button class="btn btn-outline-light w-100" type="button" @click="handleLogout">
                <i class="bi bi-person-circle me-1"></i>
                {{ userDisplayName }}
              </button>
            </div>
            
            <!-- 语言切换下拉菜单 -->
            <div class="dropdown mt-2">
              <button class="btn btn-outline-light w-100 dropdown-toggle" type="button" id="languageDropdown" data-bs-toggle="dropdown" aria-expanded="false">
                <i class="bi bi-translate me-1"></i>
                {{ currentLanguage }}
              </button>
              <ul class="dropdown-menu w-100" aria-labelledby="languageDropdown">
                <li><a class="dropdown-item" href="#" @click.prevent="switchLanguage('zh_cn')">中文</a></li>
                <li><a class="dropdown-item" href="#" @click.prevent="switchLanguage('en_us')">English</a></li>
              </ul>
            </div>
          </div>
          
          <!-- 大屏设备：水平排列的功能按钮 -->
          <div class="d-none d-lg-flex align-items-center ms-auto">
            <!-- 功能导航链接 -->
            <nav class="nav me-3">
              <RouterLink class="nav-link text-light d-flex align-items-center" aria-current="page" to="/" active-class="active">
                <i class="bi bi-house me-2"></i>
                {{ t("navbar.home") }}
              </RouterLink>
              <RouterLink class="nav-link text-light d-flex align-items-center" to="/browse" active-class="active">
                <i class="bi bi-search me-2"></i>
                {{ t("navbar.browse") }}
              </RouterLink>
              <RouterLink class="nav-link text-light d-flex align-items-center" to="/query" active-class="active">
                <i class="bi bi-question-circle me-2"></i>
                {{ t("navbar.query") }}
              </RouterLink>
              <RouterLink class="nav-link text-light d-flex align-items-center" to="/about" active-class="active">
                <i class="bi bi-info-circle me-2"></i>
                {{ t("navbar.about") }}
              </RouterLink>
            </nav>
            
            <!-- 认证按钮 -->
            <div class="dropdown ms-3" v-if="!isAuthenticated">
              <button class="btn btn-outline-light" type="button" @click="openAuthModal">
                <i class="bi bi-key me-1"></i>
                {{ t('auth.login') }}
              </button>
            </div>
            
            <!-- 用户菜单（已登录状态） -->
            <div class="dropdown ms-3" v-else>
              <button class="btn btn-outline-light dropdown-toggle" type="button" id="userDropdown" data-bs-toggle="dropdown" aria-expanded="false">
                <i class="bi bi-person-circle me-1"></i>
                {{ userDisplayName }}
              </button>
              <ul class="dropdown-menu z-3" aria-labelledby="userDropdown">
                <li v-if="hasAdminPermission">
                  <RouterLink class="dropdown-item" to="/superadmin">
                    <i class="bi bi-shield-lock me-2"></i>
                    {{ t('navbar.superadmin') }}
                  </RouterLink>
                </li>
                <li><hr class="dropdown-divider" v-if="hasAdminPermission"></li>
                <li><a class="dropdown-item" href="#" @click.prevent="handleLogout">
                  <i class="bi bi-box-arrow-right me-2"></i>
                  {{ t('auth.logout') }}
                </a></li>
              </ul>
            </div>
            
            <!-- 语言切换下拉菜单 -->
            <div class="dropdown ms-3">
              <button class="btn btn-outline-light dropdown-toggle" type="button" id="languageDropdown" data-bs-toggle="dropdown" aria-expanded="false">
                <i class="bi bi-translate me-1"></i>
                {{ currentLanguage }}
              </button>
              <ul class="dropdown-menu z-3" aria-labelledby="languageDropdown">
                <li><a class="dropdown-item" href="#" @click.prevent="switchLanguage('zh_cn')">中文</a></li>
                <li><a class="dropdown-item" href="#" @click.prevent="switchLanguage('en_us')">English</a></li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </nav>
    
    <!-- 认证模态框 -->
    <AuthModal 
      :show="showAuthModal" 
      @update:show="showAuthModal = $event"
      @login-success="handleLoginSuccess"
    />
</template>
<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n';
import { RouterLink } from 'vue-router';
import AuthModal from './AuthModal.vue';
import { useAuth } from '../composables/useAuth';
const { t, locale } = useI18n();

// 使用全局认证状态
const { isAuthenticated, userInfo, isAdmin, login, logout, getUserInfo, verifyPasskeyModifiable } = useAuth();
const showAuthModal = ref(false);
const hasAdminPermission = ref(false);

// 计算当前显示的语言名称
const currentLanguage = computed(() => {
  return locale.value === 'zh_cn' ? '中文' : 'English';
});

// 获取用户显示名称
const userDisplayName = computed(() => {
  const info = getUserInfo();
  if (!info) return t('auth.logout');
  
  // 显示operator属性
  return info.operator || t('auth.logout');
});

// 切换语言函数
const switchLanguage = (lang) => {
  locale.value = lang;
  // 可选：将语言设置保存到localStorage
  localStorage.setItem('user-language', lang);
}

// 打开认证模态框
const openAuthModal = () => {
  showAuthModal.value = true;
}

// 处理登录成功
const handleLoginSuccess = (data) => {
  // 根据返回体结构解析数据
  // 返回体: { operator: "string", token: "string", description: "string" }
  if (data && data.token) {
    const token = data.token;
    const userData = {
      operator: data.operator,
      description: data.description
    };
    login(token, userData);
  }
  showAuthModal.value = false;
  console.log('Login successful, user data:', data);
}

// 退出登录
const handleLogout = () => {
  logout();
  // 可以添加退出成功的提示
  alert(t('auth.logout'));
}

// 检查管理员权限
const checkAdminPermission = async () => {
  if (isAuthenticated.value) {
    hasAdminPermission.value = await verifyPasskeyModifiable();
  } else {
    hasAdminPermission.value = false;
  }
};

// 监听认证状态变化
watch(isAuthenticated, (newVal) => {
  if (newVal) {
    checkAdminPermission();
  } else {
    hasAdminPermission.value = false;
  }
});

// 从localStorage加载用户的语言偏好
onMounted(() => {
  const savedLanguage = localStorage.getItem('user-language');
  if (savedLanguage) {
    locale.value = savedLanguage;
  }
  
  // 初始检查管理员权限
  if (isAuthenticated.value) {
    checkAdminPermission();
  }
})

</script>
