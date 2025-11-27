<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import MoleculeCanvas from './MoleculeCanvas.vue'
import CompoundDetail from './CompoundDetail.vue'

const { t } = useI18n()

// 搜索模式状态
const searchMode = ref('structure') // structure, substructure, similarity
const ketcherRef = ref(null)
const searchResults = ref([])
const loading = ref(false)
const errorMessage = ref('')
const currentSmiles = ref('')
const selectedCompound = ref(null)
const showDetail = ref(false)

// 分页相关
const currentPage = ref(1)
const itemsPerPage = ref(20)
const totalItems = ref(0)

// 处理搜索模式切换
const setSearchMode = (mode) => {
  searchMode.value = mode
}

// 获取Ketcher实例
const getKetcher = () => {
  return ketcherRef.value?.contentWindow?.ketcher
}

// 从Ketcher获取SMILES
const getSmilesFromKetcher = async () => {
  try {
    const ketcher = getKetcher()
    if (ketcher) {
      const smiles = await ketcher.getSmiles()
      currentSmiles.value = smiles
      return smiles
    }
  } catch (error) {
    console.error('获取SMILES失败:', error)
    errorMessage.value = t('query.get_structure_failed')
  }
  return null
}

// 根据ID获取完整化合物数据
const fetchCompoundById = async (id) => {
  try {
    const response = await fetch(`/api/data/${id}`)
    if (!response.ok) {
      throw new Error(`获取化合物数据失败: ${response.status}`)
    }
    const result = await response.json()
    // API返回格式：{"code": 200200, "data": {...}, "msg": "success"}
    if (result.code === 200200 && result.data) {
      return result.data
    } else {
      throw new Error('API返回数据格式错误')
    }
  } catch (error) {
    console.error(`获取化合物 ${id} 数据失败:`, error)
    // 返回基本数据作为fallback
    return {
      id: id,
      item_name: `化合物 ${id}`,
      smiles: 'C1=CC=CC=C1'
    }
  }
}

// 处理搜索操作
const handleSearch = async () => {
  loading.value = true
  errorMessage.value = ''
  searchResults.value = []
  currentPage.value = 1

  try {
    const smiles = await getSmilesFromKetcher()
    if (!smiles) {
      errorMessage.value = t('query.draw_structure_first')
      return
    }

    let response
    switch (searchMode.value) {
      case 'structure':
        // 精确结构匹配
        response = await fetch(`/api/rdkit/exact-match?smiles=${encodeURIComponent(smiles)}`)
        break
      case 'substructure':
        // 子结构搜索
        response = await fetch(`/api/rdkit/substructure-search?smarts_pattern=${encodeURIComponent(smiles)}`)
        break
      case 'similarity':
        // 相似度搜索
        // 先获取指纹，然后进行相似度搜索
        const fpResponse = await fetch(`/api/rdkit/smiles-to-fingerprint?smiles=${encodeURIComponent(smiles)}`)
        const fpData = await fpResponse.json()
        response = await fetch(`/api/rdkit/similarity?qfp=${encodeURIComponent(fpData.fingerprint)}&threshold=0.5`)
        break
    }

    if (!response.ok) {
      throw new Error(`API请求失败: ${response.status}`)
    }

    const result = await response.json()
    console.log('API响应结果:', result)
    
    // 处理API响应格式
    if (result.code === 200200 && result.data) {
      try {
        const parsedData = JSON.parse(result.data)
        console.log('解析出的数据:', parsedData)
        
        let compoundIds = []
        
        if (searchMode.value === 'similarity') {
          // 相似度搜索：数据格式为 [["CMP0005", 1.0], ["CMP0054", 1.0]]
          if (Array.isArray(parsedData) && parsedData.length > 0) {
            compoundIds = parsedData.map(item => item[0]) // 提取ID
          }
        } else {
          // 精确匹配和子结构搜索：数据格式为 ["CMP0002", "CMP0003"]
          if (Array.isArray(parsedData) && parsedData.length > 0) {
            compoundIds = parsedData
          }
        }
        
        console.log('提取的化合物ID:', compoundIds)
        
        // 根据ID获取完整的化合物数据
        if (compoundIds.length > 0) {
          const compoundPromises = compoundIds.map(id => fetchCompoundById(id))
          const compounds = await Promise.all(compoundPromises)
          searchResults.value = compounds
        } else {
          searchResults.value = []
        }
      } catch (parseError) {
        console.error('解析数据失败:', parseError)
        searchResults.value = []
      }
    } else {
      searchResults.value = []
    }
    
    totalItems.value = searchResults.value.length
    
    // 如果没有结果，显示提示信息
    if (searchResults.value.length === 0) {
      errorMessage.value = t('query.no_results')
    }
    
  } catch (error) {
    console.error('搜索失败:', error)
    errorMessage.value = error.message || t('query.search_failed')
  } finally {
    loading.value = false
  }
}

// 下载结构
const handleDownloadStructure = async () => {
  const ketcher = getKetcher()
  if (ketcher) {
    try {
      const molfile = await ketcher.getMolfile()
      const blob = new Blob([molfile], { type: 'chemical/x-mdl-molfile' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = 'structure.mol'
      a.click()
      URL.revokeObjectURL(url)
    } catch (error) {
      console.error('下载结构失败:', error)
      errorMessage.value = t('query.download_structure_failed')
    }
  } else {
    errorMessage.value = t('query.draw_structure_first')
  }
}

// 分页相关方法
const totalPages = () => {
  return Math.ceil(totalItems.value / itemsPerPage.value)
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages()) {
    currentPage.value = page
  }
}

// 获取当前页的数据
const getCurrentPageResults = () => {
  const startIndex = (currentPage.value - 1) * itemsPerPage.value
  const endIndex = startIndex + itemsPerPage.value
  return searchResults.value.slice(startIndex, endIndex)
}

// 获取要显示的页码数组（智能分页）
const getDisplayPages = () => {
  const total = totalPages()
  const current = currentPage.value
  const pages = []
  
  // 如果总页数小于等于7，显示所有页码
  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
    return pages
  }
  
  // 显示当前页附近的页码
  if (current <= 4) {
    // 前几页
    for (let i = 1; i <= 5; i++) {
      pages.push(i)
    }
    pages.push('...')
    pages.push(total)
  } else if (current >= total - 3) {
    // 后几页
    pages.push(1)
    pages.push('...')
    for (let i = total - 4; i <= total; i++) {
      pages.push(i)
    }
  } else {
    // 中间页
    pages.push(1)
    pages.push('...')
    for (let i = current - 1; i <= current + 1; i++) {
      pages.push(i)
    }
    pages.push('...')
    pages.push(total)
  }
  
  return pages
}

// 查看化合物详情
const handleShowCompoundDetail = (compound) => {
  selectedCompound.value = compound
  showDetail.value = true
}

// 获取搜索模式文本
const getSearchModeText = () => {
  const modes = {
    structure: t('query.exact_structure'),
    substructure: t('query.substructure'),
    similarity: t('query.similarity')
  }
  return modes[searchMode.value]
}

// 组件挂载后设置Ketcher引用
onMounted(() => {
  // Ketcher iframe加载完成后可以执行初始化操作
})
</script>

<template>
  <div class="container-fluid py-4">
    <!-- 页面标题 -->
    <div class="row mb-4">
      <div class="col-12">
        <h1 class="display-6 text-primary fw-bold">{{ t('query.title') }}</h1>
        <p class="text-muted">{{ t('query.description') }}</p>
      </div>
    </div>

    <div class="row">
      <!-- 左侧：控制面板 -->
      <div class="col-lg-3 col-md-4 mb-4">
        <div class="card shadow-sm border-0">
          <div class="card-header bg-primary text-white">
            <h5 class="card-title mb-0">
              <i class="bi bi-sliders"></i> {{ t('query.search_control') }}
            </h5>
          </div>
          <div class="card-body">
            <!-- 搜索模式按钮组 -->
            <div class="mb-4">
              <label class="form-label fw-semibold">{{ t('query.search_mode') }}</label>
              <div class="btn-group-vertical w-100" role="group">
                <button
                  type="button"
                  class="btn btn-outline-primary text-start"
                  :class="{ 'active': searchMode === 'structure' }"
                  @click="setSearchMode('structure')"
                >
                  <i class="bi bi-square"></i> {{ t('query.exact_structure') }}
                </button>
                <button
                  type="button"
                  class="btn btn-outline-primary text-start"
                  :class="{ 'active': searchMode === 'substructure' }"
                  @click="setSearchMode('substructure')"
                >
                  <i class="bi bi-square-half"></i> {{ t('query.substructure') }}
                </button>
                <button
                  type="button"
                  class="btn btn-outline-primary text-start"
                  :class="{ 'active': searchMode === 'similarity' }"
                  @click="setSearchMode('similarity')"
                >
                  <i class="bi bi-arrow-left-right"></i> {{ t('query.similarity') }}
                </button>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="d-grid gap-2">
              <button
                type="button"
                class="btn btn-success btn-lg"
                @click="handleSearch"
              >
                <i class="bi bi-search"></i> {{ t('query.execute_search') }}
              </button>
              <button
                type="button"
                class="btn btn-info btn-lg"
                @click="handleDownloadStructure"
              >
                <i class="bi bi-download"></i> {{ t('query.save_structure') }}
              </button>
            </div>

            <!-- 当前模式显示 -->
            <div class="mt-4 p-3 bg-light rounded">
              <small class="text-muted">{{ t('query.current_mode') }}:</small>
              <div class="fw-bold text-primary">{{ getSearchModeText() }}{{ t('query.search') }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：Ketcher 编辑器和搜索结果 -->
      <div class="col-lg-9 col-md-8">
        <!-- Ketcher 编辑器 -->
        <div class="card shadow-sm border-0 mb-4">
          <div class="card-header bg-primary text-white d-flex justify-content-between align-items-center">
            <h5 class="card-title mb-0">
              <i class="bi bi-pencil-square"></i> {{ t('query.chemical_editor') }}
            </h5>
            <span class="badge bg-light text-primary">{{ t('query.ketcher') }}</span>
          </div>
          <div class="card-body p-0">
            <iframe 
              ref="ketcherRef"
              src="/Ketcher/index.html" 
              class="w-100 border-0" 
              style="height: 400px; min-height: 300px;"
              :title="t('query.chemical_editor')"
            ></iframe>
          </div>
        </div>

        <!-- 错误信息显示 -->
        <div v-if="errorMessage" class="alert alert-danger alert-dismissible fade show" role="alert">
          <i class="bi bi-exclamation-triangle"></i> {{ errorMessage }}
          <button type="button" class="btn-close" @click="errorMessage = ''"></button>
        </div>

        <!-- 搜索结果区域 -->
        <div v-if="searchResults.length > 0" class="card shadow-sm border-0">
          <div class="card-header bg-success text-white d-flex justify-content-between align-items-center">
            <h5 class="card-title mb-0">
              <i class="bi bi-search"></i> {{ t('query.search') }} ({{ searchResults.length }} {{ t('query.results_unit') }})
            </h5>
            <span class="badge bg-light text-success">{{ t('query.page_info', { current: currentPage, total: totalPages() }) }}</span>
          </div>
          <div class="card-body">
            <div v-if="loading" class="text-center py-3">
              <div class="spinner-border text-success" role="status">
                <span class="visually-hidden">{{ t('query.searching') }}</span>
              </div>
              <p class="mt-2 text-muted">{{ t('query.searching') }}</p>
            </div>
            
            <div v-else class="row g-3">
              <div v-for="result in getCurrentPageResults()" :key="result.id || result.ID" class="col-12">
                <div class="card border">
                  <div class="card-body">
                    <div class="row align-items-center">
                      <div class="col-md-3">
                        <div class="text-center">
                          <div class="molecule-preview bg-light rounded p-2">
                            <MoleculeCanvas 
                              :smiles="result.smiles" 
                              :width="150" 
                              :height="100"
                            />
                          </div>
                        </div>
                      </div>
                      <div class="col-md-9">
                        <h6 class="card-title">{{ result.item_name || result.ItemName || t('query.unnamed_compound') }}</h6>
                        <div class="row small text-muted">
                          <div class="col-6">
                            <strong>{{ t('details.id') }}:</strong> {{ result.id || result.ID || 'N/A' }}
                          </div>
                          <div class="col-6">
                            <strong>{{ t('details.cas_number') }}:</strong> {{ result.cas_number || result.CASNumber || 'N/A' }}
                          </div>
                          <div class="col-6">
                            <strong>{{ t('details.formula') }}:</strong> {{ result.formula || result.Formula || 'N/A' }}
                          </div>
                          <div class="col-6">
                            <strong>{{ t('details.source') }}:</strong> {{ result.source || result.Source || 'N/A' }}
                          </div>
                        </div>
                        <div class="mt-2">
                          <button 
                            class="btn btn-outline-primary btn-sm"
                            @click="handleShowCompoundDetail(result)"
                          >
                            <i class="bi bi-eye"></i> {{ t('compound_card.details') }}
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 分页导航 -->
            <nav v-if="totalPages() > 1" aria-label="搜索结果分页" class="mt-4">
              <ul class="pagination justify-content-center">
                <li class="page-item" :class="{ 'disabled': currentPage === 1 }">
                  <a class="page-link" href="#" @click.prevent="goToPage(currentPage - 1)">
                    <i class="bi bi-chevron-left"></i>
                  </a>
                </li>
                <li 
                  v-for="page in getDisplayPages()" 
                  :key="page"
                  class="page-item" 
                  :class="{ 
                    'active': page === currentPage,
                    'disabled': page === '...'
                  }"
                >
                  <a 
                    class="page-link" 
                    href="#" 
                    @click.prevent="page !== '...' && goToPage(page)"
                  >
                    {{ page }}
                  </a>
                </li>
                <li class="page-item" :class="{ 'disabled': currentPage === totalPages() }">
                  <a class="page-link" href="#" @click.prevent="goToPage(currentPage + 1)">
                    <i class="bi bi-chevron-right"></i>
                  </a>
                </li>
              </ul>
            </nav>
          </div>
        </div>

        <!-- 无结果提示 -->
        <div v-else-if="!loading && searchResults.length === 0 && currentSmiles" class="card shadow-sm border-0">
          <div class="card-body text-center py-5">
            <i class="bi bi-search text-muted" style="font-size: 3rem;"></i>
            <p class="text-muted mt-2">{{ t('query.no_results') }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 化合物详情组件 -->
    <CompoundDetail 
      :compound="selectedCompound"
      :show="showDetail"
      @update:show="showDetail = $event"
    />
  </div>
</template>

<style scoped>
.card {
  border-radius: 12px;
}

.btn-group-vertical .btn {
  border-radius: 6px;
  margin-bottom: 4px;
}

.btn-group-vertical .btn:last-child {
  margin-bottom: 0;
}

.btn.active {
  background-color: var(--bs-primary);
  border-color: var(--bs-primary);
  color: white;
}

.iframe-container {
  border-radius: 8px;
  overflow: hidden;
}

.badge {
  font-size: 0.75rem;
  padding: 0.375rem 0.75rem;
}
</style>
