<script setup>
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import compoundCard from './compoundCard.vue';
import CompoundDetail from './CompoundDetail.vue';
const { t } = useI18n();

// 响应式数据
const compounds = ref([]);
const loading = ref(false);
const currentPage = ref(1);
const itemsPerPage = ref(12);
const totalItems = ref(0);
const selectedCompound = ref(null);
const showDetail = ref(false);

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
  showDetail.value = true;
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
            <!-- 化合物详情组件 -->
            <CompoundDetail 
                :compound="selectedCompound"
                :show="showDetail"
                @update:show="showDetail = $event"
            />
        </div>
    </div>
</template>
