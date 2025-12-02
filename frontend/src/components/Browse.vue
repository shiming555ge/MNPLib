<script setup>
import { ref, onMounted, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import compoundCard from './compoundCard.vue';
import CompoundDetail from './CompoundDetail.vue';
// 引入新组件
import FilterPanel from './FilterPanel.vue';

const { t } = useI18n();

// 响应式数据
const compounds = ref([]);
const loading = ref(false);
const currentPage = ref(1);
const itemsPerPage = ref(12);
const totalItems = ref(0);
const selectedCompound = ref(null);
const showDetail = ref(false);

// 筛选条件 (严格对应 API 参数)
const filters = ref({
  item_type: [],     // 数组
  description: [],   // 数组
  source: [],        // 数组，新增source筛选
  min_weight: '',
  max_weight: ''
});

// 分类数据
const itemTypes = ref([]);
const descriptions = ref([]);
const sources = ref([]); // 新增sources数据

// 1. 获取分类选项（获取descriptions和sources，item_type使用固定值）
const fetchOptions = async () => {
  try {
    const [descRes, sourcesRes] = await Promise.all([
      fetch('/api/data/descriptions'),
      fetch('/api/data/sources')
    ]);
    
    if (descRes.ok) {
      const result = await descRes.json();
      descriptions.value = result.data || result;
    }
    
    if (sourcesRes.ok) {
      const result = await sourcesRes.json();
      sources.value = result.data || result;
    }
  } catch (error) {
    console.error('Fetch options error:', error);
  }
};

// 2. 核心：获取数据 (解决数组传参问题)
const fetchCompounds = async (page = 1) => {
  loading.value = true;
  try {
    const params = new URLSearchParams();
    
    // 基本分页参数
    params.append('limit', itemsPerPage.value);
    params.append('offset', (page - 1) * itemsPerPage.value);
    
    // --- 关键：数组参数处理 ---
    // 遍历数组，生成 ?item_type=A&item_type=B 格式
    if (filters.value.item_type && filters.value.item_type.length) {
      filters.value.item_type.forEach(val => params.append('item_type', val));
    }
    
    if (filters.value.description && filters.value.description.length) {
      filters.value.description.forEach(val => params.append('description', val));
    }
    
    // 新增source筛选参数
    if (filters.value.source && filters.value.source.length) {
      filters.value.source.forEach(val => params.append('source', val));
    }
    
    // 普通参数
    if (filters.value.min_weight) params.append('min_weight', filters.value.min_weight);
    if (filters.value.max_weight) params.append('max_weight', filters.value.max_weight);
    
    // 调用 API
    const response = await fetch(`/api/data/filter?${params.toString()}`);
    
    if (!response.ok) throw new Error(`Status: ${response.status}`);
    
    const result = await response.json();
    
    // 新的API响应结构：直接返回数组或包含data字段的对象
    // 预期返回格式: { data: [{id, item_name, smiles, cas_number}, ...], total: number }
    const dataNode = result.data || result;
    compounds.value = Array.isArray(dataNode) ? dataNode : (dataNode.data || []);
    totalItems.value = result.total || dataNode.total || compounds.value.length;
    currentPage.value = page;

  } catch (error) {
    console.error('Fetch error:', error);
    compounds.value = [];
    totalItems.value = 0;
  } finally {
    loading.value = false;
  }
};

// 分页与交互逻辑
const totalPages = () => Math.ceil(totalItems.value / itemsPerPage.value);
const goToPage = (page) => {
  if (page >= 1 && page <= totalPages()) fetchCompounds(page);
};

const applyFilters = () => {
  // 移动端点击应用后，触发筛选并自动关闭 Offcanvas
  currentPage.value = 1;
  fetchCompounds(1);
};

const resetFilters = () => {
  filters.value = {
    item_type: [],
    description: [],
    source: [],
    min_weight: '',
    max_weight: ''
  };
  currentPage.value = 1;
  fetchCompounds(1);
};

// 统计当前筛选数量 (用于移动端按钮上的红点)
const activeFilterCount = computed(() => {
  let count = filters.value.item_type.length + filters.value.description.length + (filters.value.source?.length || 0);
  if (filters.value.min_weight || filters.value.max_weight) count++;
  return count;
});

const getDisplayPages = () => { /* ...保持原有的分页逻辑... */ 
    const total = totalPages();
    const current = currentPage.value;
    const pages = [];
    if (total <= 7) {
        for (let i = 1; i <= total; i++) pages.push(i);
    } else {
        if (current <= 4) {
            for (let i = 1; i <= 5; i++) pages.push(i);
            pages.push('...', total);
        } else if (current >= total - 3) {
            pages.push(1, '...');
            for (let i = total - 4; i <= total; i++) pages.push(i);
        } else {
            pages.push(1, '...', current - 1, current, current + 1, '...', total);
        }
    }
    return pages;
};

onMounted(() => {
  fetchOptions();
  fetchCompounds(1);
});
</script>

<template>
  <div class="container-fluid py-4 bg-light min-vh-100">
    
    <div class="row align-items-end mb-4">
      <div class="col">
        <h1 class="display-6 text-primary fw-bold mb-1">{{ t('browse.title') }}</h1>
        <p class="text-muted mb-0">{{ t('browse.description') }}</p>
      </div>
      
      <div class="col-auto d-md-none">
        <button 
          class="btn btn-white border shadow-sm position-relative rounded-pill px-3" 
          type="button" 
          data-bs-toggle="offcanvas" 
          data-bs-target="#filterOffcanvas"
        >
          <i class="bi bi-funnel"></i> {{ t('browse.filter') }}
          <span 
            v-if="activeFilterCount > 0" 
            class="position-absolute top-0 start-100 translate-middle badge rounded-pill bg-danger"
          >
            {{ activeFilterCount }}
          </span>
        </button>
      </div>
    </div>

    <div class="row g-4" style="position: relative;">
      
      <div class="col-md-3 d-none d-md-block">
        <div class="card shadow-sm border-0 sticky-top custom-sticky-sidebar" style="top: 120px; z-index: 900;">
          <div class="card-header bg-white pt-3 pb-2 border-bottom-0">
             <h6 class="fw-bold mb-0"><i class="bi bi-sliders me-2"></i>{{ t('browse.filter') }}</h6>
          </div>
          <div class="card-body pt-0">
            <FilterPanel 
              :filters="filters"
              :descriptions="descriptions"
              :sources="sources"
              :loading="loading"
              @apply="applyFilters"
              @reset="resetFilters"
            />
          </div>
        </div>
      </div>

      <div class="offcanvas offcanvas-start d-md-none" tabindex="-1" id="filterOffcanvas" data-bs-backdrop="static">
        <div class="offcanvas-header bg-light border-bottom">
          <h5 class="offcanvas-title fw-bold">
            <i class="bi bi-funnel-fill text-primary me-2"></i>{{ t("browse.filter") }}
          </h5>
          <button type="button" class="btn-close" data-bs-dismiss="offcanvas"></button>
        </div>
        <div class="offcanvas-body">
          <FilterPanel 
              :filters="filters"
              :descriptions="descriptions"
              :sources="sources"
              :loading="loading"
              @apply="applyFilters"
              @reset="resetFilters"
            />
        </div>
      </div>

      <div class="col-12 col-md-9">
        
        <div class="d-flex justify-content-between align-items-center mb-3">
          <span class="text-muted small" v-if="!loading">
            {{ t('browse.found').replace("%d", totalItems) }}
          </span>
          <button 
             v-if="activeFilterCount > 0" 
             class="btn btn-link btn-sm text-decoration-none text-muted d-none d-md-block"
             @click="resetFilters"
          >
             <i class="bi bi-x-circle"></i> {{ t('browse.reset_filters') }}
          </button>
        </div>

        <div v-if="loading" class="text-center py-5">
           <div class="spinner-border text-primary" role="status"></div>
        </div>

        <div v-else-if="totalItems > 0" class="row g-3">
          <compoundCard 
            v-for="compound in compounds" 
            :key="compound.ID" 
            :compound="compound"
            @show-detail="selectedCompound = compound; showDetail = true"
          />
        </div>

        <div v-else class="text-center py-5 bg-white rounded shadow-sm">
           <p class="text-muted mb-0">{{ t('browse.not_found') }}</p>
        </div>

        <nav v-if="!loading && totalPages() > 1" class="mt-5">
            <ul class="pagination justify-content-center">
                <li class="page-item" :class="{ disabled: currentPage === 1 }">
                    <a class="page-link rounded-circle mx-1" href="#" @click.prevent="goToPage(currentPage - 1)"><i class="bi bi-chevron-left"></i></a>
                </li>
                <li v-for="p in getDisplayPages()" :key="p" class="page-item" :class="{ active: p === currentPage }">
                    <a class="page-link rounded-circle mx-1" href="#" @click.prevent="goToPage(p)">{{ p }}</a>
                </li>
                <li class="page-item" :class="{ disabled: currentPage === totalPages() }">
                    <a class="page-link rounded-circle mx-1" href="#" @click.prevent="goToPage(currentPage + 1)"><i class="bi bi-chevron-right"></i></a>
                </li>
            </ul>
        </nav>

      </div>
    </div>
    
    <CompoundDetail 
      :compound="selectedCompound" 
      :show="showDetail" 
      @update:show="showDetail = $event" 
    />
  </div>
</template>

<style scoped>
.custom-sticky-sidebar {
  /* 限制最大高度，防止展开时覆盖navbar */
  max-height: calc(100vh - 140px);
  
  /* 开启垂直滚动 */
  overflow-y: auto;
  
  /* 美化滚动条 (Webkit browsers) */
  scrollbar-width: thin;
  scrollbar-color: #dee2e6 transparent;
}

/* 针对 Chrome/Safari/Edge 的滚动条样式优化 */
.custom-sticky-sidebar::-webkit-scrollbar {
  width: 6px;
}
.custom-sticky-sidebar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-sticky-sidebar::-webkit-scrollbar-thumb {
  background-color: #dee2e6;
  border-radius: 20px;
}
</style>
