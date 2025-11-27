<template>
    <nav ref="navbarRef" class="navbar navbar-expend-lg text-white navbar-bg-primary">
      <div class="container-fluid justify-content-between navbar-bg-primary">
        <a class="navbar-brand text-light p-3 fs-1" href="#">{{ t("navbar.brand_for_short") }}<span class="fs-6"> {{ t("navbar.brand_for_full") }}</span></a>
        <div class="d-flex align-items-center">
          <!-- 语言切换下拉菜单 -->
          <div class="dropdown me-3">
            <button class="btn btn-outline-light dropdown-toggle" type="button" id="languageDropdown" data-bs-toggle="dropdown" aria-expanded="false">
              <i class="bi bi-translate me-1"></i>
              {{ currentLanguage }}
            </button>
            <ul class="dropdown-menu z-3" aria-labelledby="languageDropdown">
              <li><a class="dropdown-item" href="#" @click.prevent="switchLanguage('zh_cn')">中文</a></li>
              <li><a class="dropdown-item" href="#" @click.prevent="switchLanguage('en_us')">English</a></li>
            </ul>
          </div>
          <img class='navbar-brand' src="/logo.png" alt="ZJUT logo" />
        </div>
      </div>
    </nav>
    <div class="container-fluid sticky-top z-2" :class="{'navbar-bg-primary': isScrolled, 'navbar-bg-secondary': !isScrolled}">
        <nav class="nav nav-fill p-2 fs-6 text-light">
            <RouterLink class="nav-link" aria-current="page" to="/" active-class="active">{{ t("navbar.home") }}</RouterLink>
            <RouterLink class="nav-link" to="/browse" active-class="active">{{ t("navbar.browse") }}</RouterLink>
            <RouterLink class="nav-link" to="/query" active-class="active">{{ t("navbar.query") }}</RouterLink>
            <RouterLink class="nav-link" to="/help" active-class="active">{{ t("navbar.help") }}</RouterLink>
        </nav>
    </div>
</template>
<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useI18n } from 'vue-i18n';
import { RouterLink } from 'vue-router';
const { t, locale } = useI18n();
const isScrolled = ref(false);
const navbarRef = ref(null);

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
  
  window.addEventListener('scroll', handleScroll)
})

const handleScroll = () => {
    isScrolled.value = (window.scrollY > navbarRef.value.offsetHeight);
}

onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll)
})
</script>
