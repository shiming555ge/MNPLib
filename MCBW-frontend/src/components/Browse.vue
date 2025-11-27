<script setup>
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import compoundCard from './compoundCard.vue';
const { t } = useI18n();

// 响应式数据
const compounds = ref([]);
const loading = ref(false);
const currentPage = ref(1);
const itemsPerPage = ref(12);
const totalItems = ref(0);
const selectedCompound = ref(null);

// 从后端API获取数据
const fetchCompounds = async (page = 1) => {
  loading.value = true;
  try {
    const limit = itemsPerPage.value;
    const offset = (page - 1) * itemsPerPage.value;
    
    const response = await fetch(`/api/data?limit=${limit}&offset=${offset}`);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const result = await response.json();
    console.log(result)
    compounds.value = result.data.data || [];
    totalItems.value = result.data.total || 0;
    currentPage.value = page;
  } catch (error) {
    console.error('Failed to fetch compounds:', error);
    // 如果API调用失败，使用模拟数据作为fallback
    // compounds.value = generateMockData();
    // totalItems.value = compounds.value.length;
  } finally {
    loading.value = false;
  }
};

// // 生成模拟数据作为fallback
// const generateMockData = () => {
//   const mockData = [];
//   for (let i = 1; i <= itemsPerPage.value; i++) {
//     mockData.push({
//       ID: `CMP${i}`,
//       Source: 'Mock Source',
//       ItemName: `Compound ${i}`,
//       ItemType: 'Natural Product',
//       IUPACName: `IUPAC Name ${i}`,
//       Description: `This is a description for compound ${i}`,
//       CASNumber: `123-45-${i}`,
//       ItemTag: 'Tag',
//       Formula: 'C10H12O2',
//       Structure: 'Structure data',
//       MS1: 'MS1',
//       MS2: 'MS2',
//       Bioactivity: 'Active',
//       Smiles: 'C1=CC=CC=C1',
//       CreatedAt: new Date(),
//       UpdatedAt: new Date()
//     });
//   }
//   return mockData;
// };

// 分页相关方法
const totalPages = () => {
  return Math.ceil(totalItems.value / itemsPerPage.value);
};

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages()) {
    fetchCompounds(page);
  }
};

// 获取要显示的页码数组（智能分页）
const getDisplayPages = () => {
  const total = totalPages();
  const current = currentPage.value;
  const pages = [];
  
  // 如果总页数小于等于7，显示所有页码
  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i);
    }
    return pages;
  }
  
  // 显示当前页附近的页码
  if (current <= 4) {
    // 前几页
    for (let i = 1; i <= 5; i++) {
      pages.push(i);
    }
    pages.push('...');
    pages.push(total);
  } else if (current >= total - 3) {
    // 后几页
    pages.push(1);
    pages.push('...');
    for (let i = total - 4; i <= total; i++) {
      pages.push(i);
    }
  } else {
    // 中间页
    pages.push(1);
    pages.push('...');
    for (let i = current - 1; i <= current + 1; i++) {
      pages.push(i);
    }
    pages.push('...');
    pages.push(total);
  }
  
  return pages;
};

const setSearchMode = (mode) => {
  // 这里可以添加搜索模式切换的逻辑
  console.log('Search mode set to:', mode);
};

// 显示化合物详情
const showCompoundDetail = (compound) => {
  selectedCompound.value = compound;
  // 使用Bootstrap的JavaScript API来显示offcanvas
  const offcanvasElement = document.getElementById('detailOffcanvas');
  if (offcanvasElement) {
    // 尝试使用Bootstrap 5的API
    if (typeof bootstrap !== 'undefined' && bootstrap.Offcanvas) {
      const offcanvas = new bootstrap.Offcanvas(offcanvasElement);
      offcanvas.show();
    } else {
      // 备用方案：直接设置属性
      offcanvasElement.classList.add('show');
      offcanvasElement.style.visibility = 'visible';
      document.body.classList.add('offcanvas-open');
    }
  }
};

// 组件挂载时获取数据
onMounted(() => {
  fetchCompounds(1);
});
</script>

<template>
    <div class="container-fluid py-4">
        <!-- 页面标题和筛选按钮 -->
        <div class="rowmb-4  align-items-center ">
            <!-- 页面标题 -->
            <div class="col-12">
                <h1 class="display-6 text-primary fw-bold">{{ t('browse.title') }}</h1>
                <p class="text-muted">{{ t('browse.description') }}</p>
            </div>
            <!-- 标题行 -->    
            </div>
            <div class="fixed-top d-md-none start-0" style="top: 30%;">
                <button class="btn btn-outline-primary" type="button" data-bs-toggle="offcanvas" data-bs-target="#filterOffcanvas" aria-controls="filterOffcanvas">
                    <i class="bi bi-filter-left"></i>
                </button>
        </div>

        <!-- 移动端筛选侧边栏 -->
        <div class="offcanvas offcanvas-start d-md-none" tabindex="-1" id="filterOffcanvas" aria-labelledby="filterOffcanvasLabel">
            <div class="offcanvas-header">
                <h5 class="offcanvas-title" id="filterOffcanvasLabel">{{ t("browse.filter") }}</h5>
                <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close"></button>
            </div>
            <div class="offcanvas-body">
                <div class="filter-section">
                    <h6 class="fw-bold">{{ t("browse.categories") }}</h6>
                    <!-- 这里可以添加筛选选项 -->
                    <div class="form-check">
                        <input class="form-check-input" type="checkbox" id="category1">
                        <label class="form-check-label" for="category1">
                            Category 1
                        </label>
                    </div>
                    <div class="form-check">
                        <input class="form-check-input" type="checkbox" id="category2">
                        <label class="form-check-label" for="category2">
                            Category 2
                        </label>
                    </div>
                </div>
            </div>
        </div>

        <!-- 主要内容区域 -->
        <div class="row">
            <!-- 桌面端筛选列 -->
            <div class="card shadow-sm border-0 mt-3 col-md-3 d-none d-md-block">
                <div class="card-header bg-primary text-white">
                    <h6 class="card-title mb-0">
                    <i class="bi bi-filter-left"></i> {{ t('browse.filter') }}
                    </h6>
                </div>
                <!-- 筛选组 -->
                <div class="card-body"> 
                    <!-- 筛选组类型    -->
                    <label class="form-label fw-semibold">{{ t('browse.compound_type') }}</label>
                    <!-- 筛选组按钮组 -->
                    <div class="btn-group-vertical w-100" role="group">
                        <button
                            type="button"
                            class="btn btn-outline-primary text-start"
                            :class="{ 'active': searchMode === 'structure' }"
                            @click="setSearchMode('structure')"
                        >
                            Piptide
                        </button>
                        <button
                            type="button"
                            class="btn btn-outline-primary text-start"
                            :class="{ 'active': searchMode === 'substructure' }"
                            @click="setSearchMode('substructure')"
                        >
                            Alkloid
                        </button>
                        <button
                            type="button"
                            class="btn btn-outline-primary text-start"
                            :class="{ 'active': searchMode === 'similarity' }"
                            @click="setSearchMode('similarity')"
                        >
                            Piperazine
                        </button>
                    </div>
                
                </div>
            </div>

            <!-- 卡片展示区域 -->
            <div class="col-12 col-md-9">
                <!-- 加载状态 -->
                <div v-if="loading" class="text-center py-5">
                    <div class="spinner-border text-primary" role="status">
                        <span class="visually-hidden">{{ t('browse.loading') }}</span>
                    </div>
                    <p class="mt-2 text-muted">{{ t('browse.loading') }}</p>
                </div>
                
                <!-- 调试信息 -->
                <div v-if="!loading && totalItems === 0" class="text-center py-5">
                    <p class="text-muted">{{ t('browse.not_found') }}</p>
                </div>
                
                <!-- 调试信息：显示实际数据 -->
                <div v-else-if="!loading && totalItems > 0" class="mb-3">
                    <p class="text-muted small">{{ t('browse.found').replace("%d", compounds.length) }}</p>
                </div>
                
                <!-- 化合物卡片 -->
                <div v-if="!loading && totalItems > 0" class="row g-3">
                    <compoundCard 
                        v-for="compound in compounds" 
                        :key="compound.ID" 
                        :compound="compound"
                        @show-detail="showCompoundDetail(compound)"
                    />
                </div>
                
                <!-- 页尾导航 -->
                <nav v-if="!loading && totalPages() > 1" aria-label="Page navigation justify-self-center">
                    <ul class="pagination justify-content-center p-3">
                        <li class="page-item" :class="{ 'disabled': currentPage === 1 }">
                            <a class="page-link" href="#" @click.prevent="goToPage(currentPage - 1)">
                              <i class="bi bi-chevron-double-left"></i> 
                            </a>
                        </li>
                        <li 
                            v-for="page in getDisplayPages()" 
                            :key="page"
                            class="page-item" 
                            :class="{ 'active': page === currentPage }"
                        >
                            <a class="page-link" href="#" @click.prevent="goToPage(page)">{{ page }}</a>
                        </li>
                        <li class="page-item" :class="{ 'disabled': currentPage === totalPages() }">
                            <a class="page-link" href="#" @click.prevent="goToPage(currentPage + 1)">
                              <i class="bi bi-chevron-double-right"></i>
                            </a>
                        </li>
                    </ul>
                </nav>
            </div>
            <!-- 卡片侧拉数据区域 -->
            <div class="offcanvas offcanvas-end" tabindex="-1" id="detailOffcanvas" aria-labelledby="detailOffcanvasLabel">
                <div class="offcanvas-header">
                    <h5 class="offcanvas-title" id="detailOffcanvasLabel">
                        {{ selectedCompound ? selectedCompound.ItemName?.replace(/"/g, "") || t('browse.compound_details') : t('browse.compound_details') }}
                    </h5>
                    <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close"></button>
                </div>
                <div class="offcanvas-body">
                    <div v-if="selectedCompound" class="compound-detail">
                        <!-- 绘制分子结构 -->
                        <div class="card mb-3">
                            <div class="card-header bg-primary text-white">
                                <h6 class="card-title mb-0"><i class="bi bi-info-circle"></i> {{ t('browse.details.structure') }}</h6>
                            </div>
                            <canvas class="card-body" id="structure-canvas"></canvas>
                        </div>
                        <!-- 基本信息卡片 -->
                        <div class="card mb-3">
                            <div class="card-header bg-primary text-white">
                                <h6 class="card-title mb-0"><i class="bi bi-info-circle"></i> {{ t('browse.details.basic_info') }}</h6>
                            </div>
                            <div class="card-body">
                                <div class="row">
                                    <div class="col-md-12">
                                        <p><strong>{{ t('browse.details.id') }}:</strong> {{ selectedCompound.id || 'N/A' }}</p>
                                        <p><strong>{{ t('browse.details.source') }}:</strong> {{ selectedCompound.source || 'N/A' }}</p>
                                        <p><strong>{{ t('browse.details.compound_name') }}:</strong> {{ selectedCompound.item_name || 'N/A' }}</p>
                                        <p><strong>{{ t('browse.details.type') }}:</strong> {{ selectedCompound.item_type || 'N/A' }}</p>
                                        <p><strong>{{ t('browse.details.iupac_name') }}:</strong> {{ selectedCompound.iupac_name || 'N/A' }}</p>
                                    </div>
                                    <div class="col-md-12">
                                        <p><strong>{{ t('browse.details.description') }}:</strong> {{ selectedCompound.description || 'N/A' }}</p>
                                        <p><strong>{{ t('browse.details.cas_number') }}:</strong> {{ selectedCompound.cas_number || 'N/A' }}</p>
                                        <p><strong>{{ t('browse.details.formula') }}:</strong> {{ selectedCompound.formula || 'N/A' }}</p>
                                        <p><strong>{{ t('browse.details.smiles') }}:</strong> {{ selectedCompound.smiles || 'N/A' }}</p>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- 分析数据卡片 -->
                        <div class="card mb-3">
                            <div class="card-header bg-info text-white">
                                <h6 class="card-title mb-0"><i class="bi bi-graph-up"></i> {{ t('browse.details.analysis_data') }}</h6>
                            </div>
                            <div class="card-body">
                                <div class="row">
                                    <div class="col-md-6">
                                        <p><strong>MS1:</strong> {{ selectedCompound.ms1 || 'N/A' }}</p>
                                        <p><strong>MS2:</strong> {{ selectedCompound.ms2 || 'N/A' }}</p>
                                    </div>
                                    <div class="col-md-6">
                                        <p><strong>{{ t('browse.details.bioactivity') }}:</strong> {{ selectedCompound.Bioactivity || 'N/A' }}</p>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- 其他信息卡片 -->
                        <div class="card mb-3">
                            <div class="card-header bg-warning text-dark">
                                <h6 class="card-title mb-0"><i class="bi bi-card-text"></i> {{ t('browse.details.other_info') }}</h6>
                            </div>
                            <div class="card-body row">
                                <div class="col-6 p-1"><strong>{{ t('browse.details.tag') }}:</strong> {{ selectedCompound.ItemTag || 'N/A' }}</div>
                                <div class="col-6 p-1"><strong>{{ t('browse.details.structure') }}:</strong> {{ selectedCompound.Structure || 'N/A' }}</div>
                                <div class="col-6 p-1"><strong>{{ t('browse.details.created_at') }}:</strong> {{ selectedCompound.CreatedAt || 'N/A' }}</div>
                                <div class="col-6 p-1"><strong>{{ t('browse.details.updated_at') }}:</strong> {{ selectedCompound.UpdatedAt || 'N/A' }}</div>
                            </div>
                        </div>
                    </div>
                    <div v-else class="text-center py-5">
                        <p class="text-muted">{{ t('browse.select_compound') }}</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
