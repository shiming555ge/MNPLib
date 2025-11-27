<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <div class="container-fluid py-5">
    <!-- 英雄区域 -->
    <div class="row align-items-center min-vh-75">
      <div class="col-lg-6">
        <div class="hero-content">
          <h1 class="display-4 fw-bold text-primary mb-4">
            {{ t('home.title') }}
          </h1>
          <p class="lead text-muted mb-4">
            {{ t('home.description') }}
          </p>
          <div class="d-flex gap-3 flex-wrap">
            <router-link to="/browse" class="btn btn-primary btn-lg">
              <i class="bi bi-collection"></i> {{ t('home.browse_compounds') }}
            </router-link>
            <router-link to="/query" class="btn btn-outline-primary btn-lg">
              <i class="bi bi-search"></i> {{ t('home.chemical_query') }}
            </router-link>
          </div>
        </div>
      </div>
      <div class="col-lg-6">
        <div class="hero-image text-center">
          <div class="molecule-visualization bg-light rounded p-4">
            <div class="text-center">
              <i class="bi bi-hexagon text-primary" style="font-size: 8rem;"></i>
              <p class="text-muted mt-3">{{ t('home.structure_visualization') }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 功能特性 -->
    <div class="row py-5 mt-5">
      <div class="col-12 text-center mb-5">
        <h2 class="fw-bold text-primary">{{ t('home.features_title') }}</h2>
        <p class="text-muted">{{ t('home.features_subtitle') }}</p>
      </div>
      
      <div class="col-md-4 mb-4">
        <div class="card border-0 shadow-sm h-100">
          <div class="card-body text-center p-4">
            <div class="feature-icon bg-primary text-white rounded-circle mx-auto mb-3 d-flex align-items-center justify-content-center" style="width: 80px; height: 80px;">
              <i class="bi bi-collection" style="font-size: 2rem;"></i>
            </div>
            <h5 class="card-title fw-bold">{{ t('home.feature1_title') }}</h5>
            <p class="card-text text-muted">
              {{ t('home.feature1_description') }}
            </p>
          </div>
        </div>
      </div>

      <div class="col-md-4 mb-4">
        <div class="card border-0 shadow-sm h-100">
          <div class="card-body text-center p-4">
            <div class="feature-icon bg-success text-white rounded-circle mx-auto mb-3 d-flex align-items-center justify-content-center" style="width: 80px; height: 80px;">
              <i class="bi bi-search" style="font-size: 2rem;"></i>
            </div>
            <h5 class="card-title fw-bold">{{ t('home.feature2_title') }}</h5>
            <p class="card-text text-muted">
              {{ t('home.feature2_description') }}
            </p>
          </div>
        </div>
      </div>

      <div class="col-md-4 mb-4">
        <div class="card border-0 shadow-sm h-100">
          <div class="card-body text-center p-4">
            <div class="feature-icon bg-info text-white rounded-circle mx-auto mb-3 d-flex align-items-center justify-content-center" style="width: 80px; height: 80px;">
              <i class="bi bi-graph-up" style="font-size: 2rem;"></i>
            </div>
            <h5 class="card-title fw-bold">{{ t('home.feature3_title') }}</h5>
            <p class="card-text text-muted">
              {{ t('home.feature3_description') }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- 统计数据 -->
    <div class="row py-5 bg-light rounded-3">
      <div class="col-12 text-center mb-4">
        <h3 class="fw-bold text-primary">{{ t('home.stats_title') }}</h3>
      </div>
      <div class="col-md-3 text-center">
        <div class="stat-number text-primary fw-bold display-6">1,000+</div>
        <div class="stat-label text-muted">{{ t('home.stats1_label') }}</div>
      </div>
      <div class="col-md-3 text-center">
        <div class="stat-number text-success fw-bold display-6">500+</div>
        <div class="stat-label text-muted">{{ t('home.stats2_label') }}</div>
      </div>
      <div class="col-md-3 text-center">
        <div class="stat-number text-info fw-bold display-6">200+</div>
        <div class="stat-label text-muted">{{ t('home.stats3_label') }}</div>
      </div>
      <div class="col-md-3 text-center">
        <div class="stat-number text-warning fw-bold display-6">50+</div>
        <div class="stat-label text-muted">{{ t('home.stats4_label') }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.min-vh-75 {
  min-height: 75vh;
}

.hero-content {
  max-width: 500px;
}

.feature-icon {
  transition: transform 0.3s ease;
}

.card:hover .feature-icon {
  transform: scale(1.1);
}

.stat-number {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}
</style>
