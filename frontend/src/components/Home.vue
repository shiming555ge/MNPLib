<template>
  <div class="container-fluid p-0">
    <div style="min-height: 100vh;">
      <!-- Top Section: Title with Gallery Background -->
      <div class="gallery-wrapper d-flex align-items-center justify-content-center text-center text-white">
        <!-- left / center / right images are DOM elements so we can bind style dynamically -->
        <div class="gallery-img z-0 left" :style="leftStyle"></div>
        <div class="gallery-img z-1 center" :style="centerStyle"></div>
        <div class="gallery-img z-0 right" :style="rightStyle"></div>

        <h1 class="display-4 fw-bold position-relative z-3">{{ t('home.welcome_title') }}</h1>
      </div>

      <!-- Bottom Section: Description -->
      <div class="py-5 text-center">
        <p class="fs-4">{{ t('home.welcome_description') }}</p>
      </div>
    </div>

    <!-- Data Statistics Section -->
    <div class="container py-5">
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

const currentIndex = ref(0);
let timer = null;

// Fetch statistics data from API
const fetchStatistics = async () => {
  loading.value = true;
  try {
    const response = await fetch('/api/data/statistics');
    if (response.ok) {
      const data = await response.json();
      statisticsData.value = data;
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
      total_compounds: 1000,
      total_species: 500,
      bioactivity_data: 300,
      ms_data: 200,
      nmr_data: 150
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

onMounted(() => {
  // Start image rotation timer
  timer = setInterval(() => {
    currentIndex.value = (currentIndex.value + 1) % images.length;
  }, 3000);
  
  // Fetch statistics data
  fetchStatistics();
});

onBeforeUnmount(() => {
  if (timer) clearInterval(timer);
});

const prevIndex = computed(() => (currentIndex.value - 1 + images.length) % images.length);
const nextIndex = computed(() => (currentIndex.value + 1) % images.length);

const leftStyle = computed(() => ({
  backgroundImage: `url(${images[prevIndex.value]})`
}));
const centerStyle = computed(() => ({
  backgroundImage: `url(${images[currentIndex.value]})`
}));
const rightStyle = computed(() => ({
  backgroundImage: `url(${images[nextIndex.value]})`
}));

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
  height: 70vh;
  overflow: hidden;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
}

.gallery-img {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  border-radius: 16px;
  background-size: cover;
  background-position: center;
  transition: all 1s ease;
  box-shadow: 0 10px 40px rgba(0,0,0,0.45);
  opacity: 0.9;
}

.gallery-img.left {
  left: 5%;
  width: 28%;
  height: 65%;
  filter: blur(4px) brightness(0.55);
  transform: translateY(-50%) scale(0.85);
}

.gallery-img.right {
  right: 5%;
  width: 28%;
  height: 65%;
  filter: blur(4px) brightness(0.55);
  transform: translateY(-50%) scale(0.85);
}

.gallery-img.center {
  left: 50%;
  width: 50%;
  height: 85%;
  transform: translate(-50%, -50%) scale(1.08);
  opacity: 1;
  border-radius: 20px;
  box-shadow: 0 15px 50px rgba(0,0,0,0.6);
}

h1.z-3 {
  z-index: 3;
  text-shadow: 0 4px 16px rgba(0,0,0,0.75);
  font-weight: 700;
  letter-spacing: 2px;
}

/* bottom section */
.py-5 p {
  color: #333;
  max-width: 800px;
  margin: auto;
  line-height: 1.7;
}

/* subtle fade-in animation */
.gallery-img {
  animation: fadeIn 1.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-50%) scale(0.95); }
  to { opacity: 1; transform: translateY(-50%) scale(1); }
}

/* Statistics cards hover effects */
.statistics-card {
  transition: all 0.3s ease;
  border-radius: 12px;
}

.statistics-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 25px rgba(0,0,0,0.15) !important;
}

/* responsive */
@media (max-width: 768px) {
  .gallery-img.left,
  .gallery-img.right {
    display: none;
  }
  .gallery-img.center {
    width: 90%;
    height: 60%;
  }
  
  .statistics-card:hover {
    transform: translateY(-3px);
  }
}
</style>
