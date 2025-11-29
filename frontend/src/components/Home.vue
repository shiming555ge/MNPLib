<template>
  <div class="container-fluid p-0">
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
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

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

onMounted(() => {
  timer = setInterval(() => {
    currentIndex.value = (currentIndex.value + 1) % images.length;
  }, 3000);
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
</script>

<style scoped>
.gallery-wrapper {
  margin-top: 70px; /* prevent overlap with navbar */
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
}
</style>
