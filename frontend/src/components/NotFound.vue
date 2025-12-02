<template>
  <div class="container-fluid p-0">
    <div style="min-height: 100vh;">
      <!-- Top Section: 404 Error with Gallery Background -->
      <div class="d-flex align-items-center justify-content-center text-center text-white p-4">
        <div class="position-relative z-1">
          <h1 class="display-1 fw-bold mb-4">404</h1>
          <h2 class="display-5 fw-bold mb-4">{{ t('not_found.title') }}</h2>
          <p class="fs-4 mb-5" style="color: #000000AA;">{{ t('not_found.description') }}</p>
          <router-link to="/" class="btn btn-primary btn-lg px-5 py-3">
            <i class="bi bi-house-door me-2"></i>{{ t('not_found.back_home') }}
          </router-link>
        </div>
      </div>

      <!-- Bottom Section: Helpful Links -->
      <div class="container py-5">
        <div class="row justify-content-center">
          <div class="col-12 text-center mb-4">
            <h3 class="fw-bold text-primary">{{ t('not_found.helpful_links') }}</h3>
            <p class="text-muted">{{ t('not_found.helpful_description') }}</p>
          </div>
          
          <!-- Helpful Links Cards -->
          <div class="col-lg-3 col-md-4 col-sm-6 mb-4" v-for="link in helpfulLinks" :key="link.name">
            <div class="card border-0 shadow-sm h-100 helpful-card">
              <div class="card-body text-center p-4">
                <div class="mb-3">
                  <i :class="link.icon" class="fs-1 text-primary"></i>
                </div>
                <h4 class="fw-bold text-dark mb-3">{{ link.title }}</h4>
                <p class="text-muted mb-4">{{ link.description }}</p>
                <router-link :to="link.path" class="btn btn-outline-primary">
                  {{ link.buttonText }}
                </router-link>
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
import { useRouter } from 'vue-router';

const { t } = useI18n();

// Helpful links for navigation
const helpfulLinks = computed(() => [
  {
    name: 'home',
    title: t('not_found.links.home.title'),
    description: t('not_found.links.home.description'),
    icon: 'bi bi-house-door',
    path: '/',
    buttonText: t('not_found.links.home.button')
  },
  {
    name: 'browse',
    title: t('not_found.links.browse.title'),
    description: t('not_found.links.browse.description'),
    icon: 'bi bi-search',
    path: '/browse',
    buttonText: t('not_found.links.browse.button')
  },
  {
    name: 'query',
    title: t('not_found.links.query.title'),
    description: t('not_found.links.query.description'),
    icon: 'bi bi-question-circle',
    path: '/query',
    buttonText: t('not_found.links.query.button')
  }
]);
</script>

<style scoped>
h1.display-1 {
  text-shadow: 0 4px 16px rgba(0,0,0,0.75);
  font-weight: 700;
  letter-spacing: 2px;
  color: #fff;
}

h2.display-5 {
  text-shadow: 0 2px 8px rgba(0,0,0,0.75);
  color: #fff;
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

/* Helpful cards hover effects */
.helpful-card {
  transition: all 0.3s ease;
  border-radius: 12px;
}

.helpful-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 25px rgba(0,0,0,0.15) !important;
}

.btn-primary {
  background: linear-gradient(135deg, #0d6efd 0%, #0a58ca 100%);
  border: none;
  transition: all 0.3s ease;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(13, 110, 253, 0.3);
}

.btn-outline-primary:hover {
  background: linear-gradient(135deg, #0d6efd 0%, #0a58ca 100%);
  color: white;
  border: none;
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
  
  .helpful-card:hover {
    transform: translateY(-3px);
  }
  
  h1.display-1 {
    font-size: 4rem;
  }
  
  h2.display-5 {
    font-size: 1.8rem;
  }
}
</style>
