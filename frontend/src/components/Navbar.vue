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
              <RouterLink class="nav-link text-light py-2 border-bottom" to="/help" active-class="active">{{ t("navbar.help") }}</RouterLink>
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
              <RouterLink class="nav-link text-light d-flex align-items-center" to="/help" active-class="active">
                <i class="bi bi-info-circle me-2"></i>
                {{ t("navbar.help") }}
              </RouterLink>
            </nav>
            
            <!-- 语言切换下拉菜单 -->
            <div class="dropdown">
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
</template>
<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useI18n } from 'vue-i18n';
import { RouterLink } from 'vue-router';
const { t, locale } = useI18n();

// 计算当前显示的语言名称
const currentLanguage = computed(() => {
  return locale.value === 'zh_cn' ? '中文' : 'English';
});

// 切换语言函数
const switchLanguage = (lang) => {
  locale.value = lang;
  // 可选：将语言设置保存到localStorage
  localStorage.setItem('user-language', lang);
}

// 从localStorage加载用户的语言偏好
onMounted(() => {
  const savedLanguage = localStorage.getItem('user-language');
  if (savedLanguage) {
    locale.value = savedLanguage;
  }
})

</script>
