<template>
  <div class="container-fluid p-0">
    <div style="min-height: 100vh;">
      <!-- Top Section: Title and Description with Gallery Background -->
      <div class="gallery-wrapper d-flex flex-column align-items-center justify-content-center text-center text-white">
        <!-- left / center / right images are DOM elements so we can bind style dynamically -->
        <div class="gallery-img z-0 left" :style="leftStyle"></div>
        <div class="gallery-img z-1 center" :style="centerStyle"></div>
        <div class="gallery-img z-0 right" :style="rightStyle"></div>

        <div class="position-relative z-1 content-wrapper">
          <h1 class="display-4 fw-bold mb-4">{{ t('home.welcome_title') }}</h1>
          <p class="fs-5 mb-0 description-text">{{ t('home.welcome_description') }}</p>
        </div>
      </div>

      <!-- Data Statistics Section -->
      <div class="container py-4">
        <div class="row justify-content-center">
          <div class="col-12 text-center mb-4">
            <h2 class="fw-bold text-primary">{{ t('home.data_statistics') }}</h2>
            <p class="text-muted">{{ t('home.data_info') }}</p>
          </div>
          
          <!-- Statistics Cards -->
          <div class="col-lg-2 col-md-4 col-sm-6 col-6 mb-4" v-for="stat in statistics" :key="stat.key">
            <div class="card border-0 shadow-sm h-100 statistics-card">
              <div class="card-body text-center p-4">
                <div class="mb-3">
                  <i :class="getStatIcon(stat.key)" class="fs-1 text-primary"></i>
                </div>
                <h3 class="fw-bold text-dark mb-2">{{ stat.value }}</h3>
                <p class="text-muted mb-0">{{ stat.label }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

// Data statistics
const statisticsData = ref({
  total_compounds: 0,
  total_species: 0,
  bioactivity_data: 0,
  ms_data: 0,
  nmr_data: 0
});

const loading = ref(false);

const images = [
  '/home_pics/pic1.jpg',
  '/home_pics/pic2.jpg',
  '/home_pics/pic3.jpg',
  '/home_pics/pic4.jpg',
  '/home_pics/pic5.jpg',
  '/home_pics/pic6.jpg',
  '/home_pics/pic7.jpg',
  '/home_pics/pic8.jpg',
  '/home_pics/pic9.jpg'
];

const currentIndex = ref(1);
let timer = null;

// 跟踪已加载的图片和对应的Image对象
const loadedImages = ref(new Set());
const imageCache = ref(new Map()); // 存储Image对象

// Fetch statistics data from API
const fetchStatistics = async () => {
  loading.value = true;
  try {
    const response = await fetch('/api/data/statistics');
    if (response.ok) {
      const data = await response.json();
      statisticsData.value = data.data;
      statisticsData.value.total_species = "500+"
    } else {
      console.error('Failed to fetch statistics data');
      // Set default values if API fails
      statisticsData.value = {
        total_compounds: 1000,
        total_species: 500,
        bioactivity_data: 300,
        ms_data: 200,
        nmr_data: 150
      };
    }
  } catch (error) {
    console.error('Error fetching statistics:', error);
    // Set default values on error
    statisticsData.value = {
      total_compounds: "3000+",
      total_species: "500+",
      bioactivity_data: "700+",
      ms_data: "800+",
      nmr_data: "2100+"
    };
  } finally {
    loading.value = false;
  }
};

// Computed property for statistics display
const statistics = computed(() => [
  {
    key: 'compounds',
    label: t('home.compounds'),
    value: statisticsData.value.total_compounds.toLocaleString()
  },
  {
    key: 'organisms',
    label: t('home.organisms'),
    value: statisticsData.value.total_species.toLocaleString()
  },
  {
    key: 'bioactivities',
    label: t('home.bioactivities'),
    value: statisticsData.value.bioactivity_data.toLocaleString()
  },
  {
    key: 'ms',
    label: t('home.ms'),
    value: statisticsData.value.ms_data.toLocaleString()
  },
  {
    key: 'nmr',
    label: t('home.nmr'),
    value: statisticsData.value.nmr_data.toLocaleString()
  }
]);

// 改进的图片预加载函数，返回Promise并跟踪加载状态，使用缓存
const preloadImages = (imageUrls) => {
  const promises = imageUrls.map(url => {
    // 如果已经在缓存中，直接返回
    if (imageCache.value.has(url)) {
      return Promise.resolve(url);
    }
    
    return new Promise((resolve, reject) => {
      const img = new Image();
      img.onload = () => {
        loadedImages.value.add(url);
        imageCache.value.set(url, img); // 缓存Image对象
        resolve(url);
      };
      img.onerror = () => {
        console.warn(`Failed to load image: ${url}`);
        // 即使加载失败，也标记为已尝试加载
        loadedImages.value.add(url);
        imageCache.value.set(url, img); // 仍然缓存，即使加载失败
        resolve(url); // 不reject，继续执行
      };
      img.src = url;
    });
  });
  return Promise.all(promises);
};

// 图片加载状态
const imagesLoaded = ref(false);

// 检查图片是否已加载
const isImageLoaded = (url) => {
  return loadedImages.value.has(url);
};

onMounted(async () => {
  try {
    // 预加载所有图片，等待加载完成
    await preloadImages(images);
    imagesLoaded.value = true;
    console.log('All images preloaded successfully');
  } catch (error) {
    console.error('Error preloading images:', error);
    // 即使预加载失败，也继续显示，但标记为已加载
    imagesLoaded.value = true;
  }
  
  // 延迟启动图片轮播，确保图片完全加载
  setTimeout(() => {
    // Start image rotation timer - 增加间隔时间，让过渡更完整
    timer = setInterval(() => {
      currentIndex.value = (currentIndex.value + 1) % images.length;
    }, 4000); // 增加到4秒，让过渡更完整
  }, 500); // 延迟500ms，确保过渡平滑
  
  // Fetch statistics data
  fetchStatistics();
});

onBeforeUnmount(() => {
  if (timer) clearInterval(timer);
});

// 使用更高效的图片切换逻辑 - 避免重复下载
const leftStyle = computed(() => {
  const index = (currentIndex.value - 1 + images.length) % images.length;
  const url = images[index];
  const loaded = isImageLoaded(url);
  // 只有图片已加载时才设置backgroundImage
  return loaded ? {
    backgroundImage: `url(${url})`,
    opacity: 0.9,
    backgroundColor: 'transparent'
  } : {
    backgroundImage: 'none',
    opacity: 0,
    backgroundColor: '#f8f9fa'
  };
});

const centerStyle = computed(() => {
  const url = images[currentIndex.value];
  const loaded = isImageLoaded(url);
  return loaded ? {
    backgroundImage: `url(${url})`,
    opacity: 1,
    backgroundColor: 'transparent'
  } : {
    backgroundImage: 'none',
    opacity: 0,
    backgroundColor: '#f8f9fa'
  };
});

const rightStyle = computed(() => {
  const index = (currentIndex.value + 1) % images.length;
  const url = images[index];
  const loaded = isImageLoaded(url);
  return loaded ? {
    backgroundImage: `url(${url})`,
    opacity: 0.9,
    backgroundColor: 'transparent'
  } : {
    backgroundImage: 'none',
    opacity: 0,
    backgroundColor: '#f8f9fa'
  };
});

// Get appropriate icon for each statistic
const getStatIcon = (key) => {
  const icons = {
    compounds: 'bi bi-capsule',
    organisms: 'bi bi-tree',
    bioactivities: 'bi bi-heart-pulse',
    ms: 'bi bi-graph-up',
    nmr: 'bi bi-magnet'
  };
  return icons[key] || 'bi bi-graph-up';
};
</script>

<style scoped>
.gallery-wrapper {
  margin-top: 50px; /* prevent overlap with navbar */
  position: relative;
  height: 60vh; /* 减小高度，使页面更紧凑 */
  overflow: hidden;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
}

.content-wrapper {
  max-width: 40vw; /* 比中央图片稍窄，确保不超出 */
  padding: 0 20px;
  width: 100%;
}

.description-text {
  text-shadow: 0 2px 8px rgba(0,0,0,0.75);
  line-height: 1.5; /* 减少行高 */
  background-color: rgba(0, 0, 0, 0.4);
  padding: 15px 20px; /* 减少上下padding */
  border-radius: 12px;
  backdrop-filter: blur(5px);
  max-width: 100%;
  box-sizing: border-box;
  width: fit-content;
  margin: 0 auto;
  word-wrap: break-word;
  overflow-wrap: break-word;
  max-height: 80%; /* 限制最大高度 */
  overflow: hidden; /* 防止内容溢出 */
}

.gallery-img {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  border-radius: 16px;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  transition: all 1.2s cubic-bezier(0.4, 0, 0.2, 1), background-size 0.5s ease, opacity 0.8s ease;
  box-shadow: 0 10px 40px rgba(0,0,0,0.45);
  opacity: 0.9;
  background-color: #f8f9fa; /* 添加默认背景色，避免空白 */
  will-change: transform, opacity, filter, background-size; /* 提示浏览器优化 */
}

.gallery-img.left {
  left: 5%;
  width: 25%; /* 减小宽度 */
  height: 55%; /* 减小高度 */
  filter: blur(4px) brightness(0.55);
  transform: translateY(-50%) scale(0.85);
  transition-delay: 0.1s; /* 左侧图片延迟切换，创造层次感 */
}

.gallery-img.right {
  right: 5%;
  width: 25%; /* 减小宽度 */
  height: 55%; /* 减小高度 */
  filter: blur(4px) brightness(0.55);
  transform: translateY(-50%) scale(0.85);
  transition-delay: 0.2s; /* 右侧图片更晚切换 */
}

.gallery-img.center {
  left: 50%;
  width: 45%; /* 减小宽度 */
  height: 75%; /* 减小高度 */
  transform: translate(-50%, -50%) scale(1.08);
  opacity: 1;
  border-radius: 20px;
  box-shadow: 0 15px 50px rgba(0,0,0,0.6);
  transition-duration: 1.5s; /* 中央图片过渡时间更长 */
  transition-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1); /* 更弹性的缓动函数 */
}

h1.display-4 {
  text-shadow: 0 4px 16px rgba(0,0,0,0.75);
  font-weight: 700;
  letter-spacing: 2px;
  margin-bottom: 1.5rem !important;
}

/* 移除冲突的fadeIn动画，使用transition代替 */
/* .gallery-img {
  animation: fadeIn 1.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-50%) scale(0.95); }
  to { opacity: 1; transform: translateY(-50%) scale(1); }
} */

/* Statistics cards hover effects */
.statistics-card {
  transition: all 0.3s ease;
  border-radius: 12px;
}

.statistics-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 25px rgba(0,0,0,0.15) !important;
}

/* 调整统计信息部分的间距 */
.container.py-4 {
  padding-top: 2rem !important;
  padding-bottom: 2rem !important;
}

/* responsive */
@media (max-width: 768px) {
  .gallery-wrapper {
    height: 50vh; /* 移动端进一步减小高度 */
  }
  
  .gallery-img.left,
  .gallery-img.right {
    display: none;
  }
  .gallery-img.center {
    width: 85%;
    height: 50%;
  }
  
  .content-wrapper {
    max-width: 80vw; /* 在移动端与图片宽度匹配 */
  }
  
  .description-text {
    padding: 12px 15px; /* 减少上下padding */
    font-size: 1rem !important;
    line-height: 1.4; /* 进一步减少行高 */
    max-height: 70%; /* 调整最大高度 */
  }
  
  h1.display-4 {
    font-size: 2.5rem;
  }
  
  .statistics-card:hover {
    transform: translateY(-3px);
  }
}

@media (max-width: 576px) {
  .gallery-wrapper {
    height: 45vh;
  }
  
  h1.display-4 {
    font-size: 2rem;
  }
  
  .description-text {
    padding: 10px 12px; /* 进一步减少padding */
    font-size: 0.9rem !important;
    line-height: 1.3; /* 进一步减少行高 */
  }
}
</style>
