<script setup>
import { ref, watch, computed, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import MoleculeViewer from './MoleculeViewer.vue'

const { t } = useI18n()

const props = defineProps({
  compound: {
    type: Object,
    default: null
  },
  show: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:show'])

const offcanvasElement = ref(null)
const detailedData = ref(null)
const loadingDetail = ref(false)
const ms2Data = ref(null)
const cnmrData = ref(null)
const bioactivityData = ref(null)
const loadingSpecialData = ref({
  ms2: false,
  cnmr: false,
  bioactivity: false
})

// 监听compound变化，获取详细数据
watch(() => props.compound, async (newCompound) => {
  if (newCompound && newCompound.id) {
    await fetchDetailedData(newCompound.id)
    await fetchSpecialData(newCompound.id)
  } else {
    detailedData.value = null
    ms2Data.value = null
    cnmrData.value = null
    bioactivityData.value = null
  }
}, { immediate: true })

// 监听show变化，控制offcanvas显示
watch(() => props.show, (newVal) => {
  if (newVal && offcanvasElement.value) {
    // 使用Bootstrap的JavaScript API来显示offcanvas
    if (typeof bootstrap !== 'undefined' && bootstrap.Offcanvas) {
      const offcanvas = new bootstrap.Offcanvas(offcanvasElement.value)
      offcanvas.show()
      
      // 监听隐藏事件
      offcanvasElement.value.addEventListener('hidden.bs.offcanvas', () => {
        emit('update:show', false)
      })
    }
  }
})

// 获取详细数据
const fetchDetailedData = async (id) => {
  if (!id) return
  
  loadingDetail.value = true
  try {
    const response = await fetch(`/api/data/${id}`)
    if (response.ok) {
      const result = await response.json()
      detailedData.value = result.data || result
    } else {
      console.error('Failed to fetch detailed data:', response.status)
      detailedData.value = null
    }
  } catch (error) {
    console.error('Error fetching detailed data:', error)
    detailedData.value = null
  } finally {
    loadingDetail.value = false
  }
}

// 获取特殊数据（MS2, C-NMR, Bioactivity）
const fetchSpecialData = async (id) => {
  if (!id) return
  
  // 获取MS2数据
  loadingSpecialData.value.ms2 = true
  try {
    const response = await fetch(`/api/data/${id}/ms2`)
    if (response.ok) {
      const result = await response.json()
      // 检查是否返回200400（表示无数据）
      if (result.code !== 200400) {
        ms2Data.value = result.data || result
      } else {
        ms2Data.value = null
      }
    } else {
      ms2Data.value = null
    }
  } catch (error) {
    console.error('Error fetching MS2 data:', error)
    ms2Data.value = null
  } finally {
    loadingSpecialData.value.ms2 = false
  }
  
  // 获取C-NMR数据
  loadingSpecialData.value.cnmr = true
  try {
    const response = await fetch(`/api/data/${id}/cnmr`)
    if (response.ok) {
      const result = await response.json()
      if (result.code !== 200400) {
        cnmrData.value = result.data || result
      } else {
        cnmrData.value = null
      }
    } else {
      cnmrData.value = null
    }
  } catch (error) {
    console.error('Error fetching C-NMR data:', error)
    cnmrData.value = null
  } finally {
    loadingSpecialData.value.cnmr = false
  }
  
  // 获取Bioactivity数据
  loadingSpecialData.value.bioactivity = true
  try {
    const response = await fetch(`/api/data/${id}/bioactivity`)
    if (response.ok) {
      const result = await response.json()
      if (result.code !== 200400) {
        bioactivityData.value = result.data || result
      } else {
        bioactivityData.value = null
      }
    } else {
      bioactivityData.value = null
    }
  } catch (error) {
    console.error('Error fetching Bioactivity data:', error)
    bioactivityData.value = null
  } finally {
    loadingSpecialData.value.bioactivity = false
  }
}

// 合并数据：优先使用详细数据，如果没有则使用基本数据
const mergedData = computed(() => {
  if (detailedData.value) {
    return { ...props.compound, ...detailedData.value }
  }
  return props.compound
})

// 手动关闭方法
const handleClose = () => {
  emit('update:show', false)
}

// 清理函数
onUnmounted(() => {
  if (offcanvasElement.value) {
    offcanvasElement.value.removeEventListener('hidden.bs.offcanvas', () => {
      emit('update:show', false)
    })
  }
})
</script>

<template>
  <div 
    v-if="show"
    class="offcanvas offcanvas-end show" 
    tabindex="-1" 
    ref="offcanvasElement"
    id="compoundDetailOffcanvas" 
    aria-labelledby="compoundDetailOffcanvasLabel"
    data-bs-backdrop="true"
    style="visibility: visible;"
  >
    <div class="offcanvas-header">
      <h5 class="offcanvas-title" id="compoundDetailOffcanvasLabel">
        {{ mergedData ? mergedData.item_name?.replace(/"/g, "") || mergedData.ItemName?.replace(/"/g, "") || t('browse.compound_details') : t('browse.compound_details') }}
      </h5>
      <button type="button" class="btn-close" @click="handleClose" aria-label="Close"></button>
    </div>
    <div class="offcanvas-body">
      <div v-if="mergedData" class="compound-detail">
        <!-- 加载状态 -->
        <div v-if="loadingDetail" class="text-center py-3">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">加载中...</span>
          </div>
          <p class="text-muted mt-2">正在加载详细数据...</p>
        </div>
        
        <!-- 绘制分子结构 -->
        <div class="card mb-3">
          <div class="card-header bg-primary text-white">
            <h6 class="card-title mb-0"><i class="bi bi-info-circle"></i> {{ t('details.structure') }}</h6>
          </div>
          <div class="card-body text-center">
            <MoleculeViewer 
              :smiles="mergedData.smiles" 
              :width="300" 
              :height="200"
            />
          </div>
        </div>
        
        <!-- 基本信息卡片 -->
        <div class="card mb-3">
          <div class="card-header bg-primary text-white">
            <h6 class="card-title mb-0"><i class="bi bi-info-circle"></i> {{ t('details.basic_info') }}</h6>
          </div>
          <div class="card-body">
            <div class="row">
              <div class="col-md-12">
                <p><strong>{{ t('details.id') }}:</strong> {{ mergedData.id || mergedData.ID || 'N/A' }}</p>
                <p><strong>{{ t('details.source') }}:</strong> {{ mergedData.source || mergedData.Source || 'N/A' }}</p>
                <p><strong>{{ t('details.compound_name') }}:</strong> {{ mergedData.item_name || mergedData.ItemName || 'N/A' }}</p>
                <p><strong>{{ t('details.type') }}:</strong> {{ mergedData.item_type || mergedData.ItemType || 'N/A' }}</p>
                <p><strong>{{ t('details.iupac_name') }}:</strong> {{ mergedData.iupac_name || mergedData.IUPACName || 'N/A' }}</p>
              </div>
              <div class="col-md-12">
                <p><strong>{{ t('details.description') }}:</strong> {{ mergedData.description || mergedData.Description || 'N/A' }}</p>
                <p><strong>{{ t('details.cas_number') }}:</strong> {{ mergedData.cas_number || mergedData.CASNumber || 'N/A' }}</p>
                <p><strong>{{ t('details.formula') }}:</strong> {{ mergedData.formula || mergedData.Formula || 'N/A' }}</p>
                <p><strong>{{ t('details.smiles') }}:</strong> {{ mergedData.smiles || mergedData.Smiles || 'N/A' }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 分析数据卡片 -->
        <div class="card mb-3">
          <div class="card-header bg-info text-white">
            <h6 class="card-title mb-0"><i class="bi bi-graph-up"></i> {{ t('details.analysis_data') }}</h6>
          </div>
          <div class="card-body">
            <div class="row">
              <div class="col-md-6">
                <p><strong>MS1:</strong> {{ mergedData.ms1 || mergedData.MS1 || 'N/A' }}</p>
                <p><strong>MS2:</strong> 
                  <span v-if="loadingSpecialData.ms2" class="text-muted">加载中...</span>
                  <span v-else-if="ms2Data">{{ ms2Data }}</span>
                  <span v-else>N/A</span>
                </p>
              </div>
              <div class="col-md-6">
                <p><strong>C-NMR:</strong> 
                  <span v-if="loadingSpecialData.cnmr" class="text-muted">加载中...</span>
                  <span v-else-if="cnmrData">{{ cnmrData }}</span>
                  <span v-else>N/A</span>
                </p>
                <p><strong>{{ t('details.bioactivity') }}:</strong> 
                  <span v-if="loadingSpecialData.bioactivity" class="text-muted">加载中...</span>
                  <span v-else-if="bioactivityData">{{ bioactivityData }}</span>
                  <span v-else>{{ mergedData.bioactivity || mergedData.Bioactivity || 'N/A' }}</span>
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- 其他信息卡片 -->
        <div class="card mb-3">
          <div class="card-header bg-warning text-dark">
            <h6 class="card-title mb-0"><i class="bi bi-card-text"></i> {{ t('details.other_info') }}</h6>
          </div>
          <div class="card-body row">
            <div class="col-6 p-1"><strong>{{ t('details.tag') }}:</strong> {{ mergedData.item_tag || mergedData.ItemTag || 'N/A' }}</div>
            <div class="col-6 p-1"><strong>{{ t('details.structure') }}:</strong> {{ mergedData.structure || mergedData.Structure || 'N/A' }}</div>
            <div class="col-6 p-1"><strong>{{ t('details.created_at') }}:</strong> {{ mergedData.created_at || mergedData.CreatedAt || 'N/A' }}</div>
            <div class="col-6 p-1"><strong>{{ t('details.updated_at') }}:</strong> {{ mergedData.updated_at || mergedData.UpdatedAt || 'N/A' }}</div>
          </div>
        </div>
      </div>
      <div v-else class="text-center py-5">
        <p class="text-muted">{{ t('browse.select_compound') }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.compound-detail {
  max-height: calc(100vh - 120px);
  overflow-y: auto;
  width: 100%;
}

/* 确保卡片内容不会超出容器 */
.compound-detail .card {
  width: 100%;
  max-width: 100%;
}

.compound-detail .card-body {
  width: 100%;
  overflow-wrap: break-word;
  word-wrap: break-word;
}

/* 处理长文本内容 */
.compound-detail p,
.compound-detail .col-6 {
  word-break: break-word;
  overflow-wrap: break-word;
  max-width: 100%;
}

/* 确保分子查看器不会超出容器 */
.compound-detail .text-center {
  width: 100%;
}

/* 移除水平滚动条 */
.compound-detail::-webkit-scrollbar {
  width: 8px;
}

.compound-detail::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.compound-detail::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.compound-detail::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
