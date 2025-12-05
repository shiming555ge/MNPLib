<script setup>
import { ref, watch, computed, onUnmounted, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import MoleculeViewer from './MoleculeViewer.vue'
import { useAuth } from '../composables/useAuth'

const { t } = useI18n()
const { isAuthenticated, getAuthHeader, isAdmin } = useAuth()

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
const protectedData = ref(null)

// 复制文本到剪贴板
const copyToClipboard = async (text) => {
  if (!text || text === 'N/A') return
  
  try {
    await navigator.clipboard.writeText(text)
    // 可以添加一个简单的提示，这里使用Bootstrap的toast或alert
    showToast(t("superadmin.copied"))
  } catch (err) {
    console.error('复制失败:', err)
    showToast(t("superadmin.copy_failed"), 'error')
  }
}

const clickToDownload = async (type, filename = null) => {  
  try {
    // 准备请求头，如果需要认证的话
    const headers = {
      'Content-Type': 'application/json',
      ...getAuthHeader()
    };
    
    axios.get(`/api/data/${props.compound.id}/${type}`, { headers }).then(res => {
      const text = res.data; // 后端返回的纯文本
      
      const blob = new Blob([text], {
        type: 'text/plain;charset=utf-8'
      });

      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename || `${props.compound.id}.${type === 'structure' ? 'mol' : 'log'}`;   // 保存的文件名
      link.click();

      window.URL.revokeObjectURL(url);
    });

    showToast(t("details.downloaded"))
  } catch (err) {
    console.error('下载失败:', err)
    showToast(t('details.download_failed'), 'error')
  }
}

// 显示提示消息
const showToast = (message, type = 'success') => {
  // 创建一个简单的提示元素
  const toastEl = document.createElement('div')
  toastEl.className = `alert alert-${type === 'error' ? 'danger' : 'success'} alert-dismissible fade show position-fixed`
  toastEl.style.cssText = 'top: 20px; right: 20px; z-index: 9999; min-width: 200px;'
  toastEl.innerHTML = `
    ${message}
    <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
  `
  
  document.body.appendChild(toastEl)
  
  // 3秒后自动移除
  setTimeout(() => {
    if (toastEl.parentNode) {
      toastEl.remove()
    }
  }, 3000)
}

// 截断文本并添加省略号
const truncateText = (text, maxLength = 30) => {
  if (!text || text === 'N/A') return text
  
  // 确保text是字符串
  const textStr = String(text)
  if (textStr.length <= maxLength) return textStr
  return textStr.substring(0, maxLength) + '...'
}

// 监听compound变化，获取详细数据
watch(() => props.compound, async (newCompound) => {
  if (newCompound && newCompound.id) {
    await fetchDetailedData(newCompound.id)
    await fetchSpecialData(newCompound.id)
  } else {
    detailedData.value = null
    protectedData.value = null
  }
}, { immediate: true })

// 监听认证状态变化，重新获取受保护的数据
watch(isAuthenticated, (newVal) => {
  if (newVal && props.compound && props.compound.id) {
    // 用户登录后，重新获取受保护的数据
    fetchSpecialData(props.compound.id)
  }
})

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
      detailedData.value = result.data || null
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
  if (!isAuthenticated) return
  
  // 获取受保护的数据
  try {
    const headers = {
      'Content-Type': 'application/json',
      ...getAuthHeader()
    }
    
    const response = await fetch(`/api/data/${id}/protected`, {
      headers
    })
    
    if (response.ok) {
      const result = await response.json()
      protectedData.value = result.data || null
    } else {
      protectedData.value = null
      console.warn('Failed to fetch protected data:', response.status)
    }
  } catch (error) {
    console.error('Error fetching protected data:', error)
    protectedData.value = null
  }
}

// 手动关闭方法
const handleClose = () => {
  emit('update:show', false)
}

// 初始化Bootstrap tooltip
const initTooltips = () => {
  if (typeof bootstrap !== 'undefined' && bootstrap.Tooltip) {
    // 初始化所有带有data-bs-toggle="tooltip"的元素
    const tooltipTriggerList = document.querySelectorAll('[data-bs-toggle="tooltip"]')
    tooltipTriggerList.forEach(tooltipTriggerEl => {
      new bootstrap.Tooltip(tooltipTriggerEl)
    })
  }
}

// 监听show变化，当offcanvas显示时初始化tooltip
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
      
      // 初始化tooltip
      nextTick(() => {
        initTooltips()
      })
    }
  }
})

// 清理函数
onUnmounted(() => {
  if (offcanvasElement.value) {
    offcanvasElement.value.removeEventListener('hidden.bs.offcanvas', () => {
      emit('update:show', false)
    })
  }
  
  // 清理tooltip
  if (typeof bootstrap !== 'undefined' && bootstrap.Tooltip) {
    const tooltipTriggerList = document.querySelectorAll('[data-bs-toggle="tooltip"]')
    tooltipTriggerList.forEach(tooltipTriggerEl => {
      const tooltip = bootstrap.Tooltip.getInstance(tooltipTriggerEl)
      if (tooltip) {
        tooltip.dispose()
      }
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
        {{ t('browse.compound_details') }}
      </h5>
      <button type="button" class="btn-close" @click="handleClose" aria-label="Close"></button>
    </div>
    <div class="offcanvas-body">
      <div v-if="detailedData" class="compound-detail">
        <!-- 加载状态 -->
        <div v-if="loadingDetail" class="text-center py-3">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">{{ t("browse.loading") }}</span>
          </div>
          <p class="text-muted mt-2">{{ t("browse.loading") }}</p>
        </div>
        
        <!-- 绘制分子结构 -->
        <div class="card mb-3">
          <div class="card-header bg-primary text-white d-flex justify-content-between align-items-center">
            <h6 class="card-title mb-0"><i class="bi bi-info-circle"></i> {{ t('details.structure') }}</h6>
            <div>
              <i class="bi bi-download me-2" style="cursor: pointer;" 
                 @click="clickToDownload('structure', `${detailedData.id}.mol`)"
                 :title="t('details.click2download')"></i>
            </div>
          </div>
          <div class="card-body text-center">
            <MoleculeViewer 
              :id="detailedData.id" 
              :smiles="detailedData.smiles" 
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
                <p>
                  <strong>{{ t('details.id') }}:</strong> 
                  <span 
                    class="copyable-text" 
                    :title="detailedData.id || 'N/A'"
                    @click="copyToClipboard(detailedData.id)"
                    data-bs-toggle="tooltip" 
                    data-bs-placement="top"
                  >
                    {{ truncateText(detailedData.id, 20) || 'N/A' }}
                  </span>
                  <i class="bi bi-clipboard ms-1 text-muted small" style="cursor: pointer;" 
                     @click="copyToClipboard(detailedData.id)"
                     :title="t('details.click2copy')"></i>
                </p>
                <p>
                  <strong>{{ t('details.compound_name') }}:</strong> 
                  <span 
                    :title="detailedData.item_name || 'N/A'"
                    data-bs-toggle="tooltip" 
                    data-bs-placement="top"
                  >
                    {{ truncateText(detailedData.item_name, 25) || 'N/A' }}
                  </span>
                </p>
                <p>
                  <strong>{{ t('details.type') }}:</strong> 
                  {{ detailedData.item_type || 'N/A' }}
                </p>
                <p>
                  <strong>{{ t('details.weight') }}:</strong> 
                  <span 
                    class="copyable-text" 
                    :title="detailedData.weight || 'N/A'"
                    @click="copyToClipboard(detailedData.weight)"
                    data-bs-toggle="tooltip" 
                    data-bs-placement="top"
                  >
                    {{ detailedData.weight || 'N/A' }}
                  </span>
                  <i class="bi bi-clipboard ms-1 text-muted small" style="cursor: pointer;" 
                     @click="copyToClipboard(detailedData.weight)"
                     :title="t('details.click2copy')"></i>
                </p>
              </div>
              <div class="col-md-12">
                <p>
                  <strong>{{ t('details.description') }}:</strong> 
                  <span 
                    :title="detailedData.description || 'N/A'"
                    data-bs-toggle="tooltip" 
                    data-bs-placement="top"
                  >
                    {{ truncateText(detailedData.description, 40) || 'N/A' }}
                  </span>
                </p>
                <p>
                  <strong>{{ t('details.cas_number') }}:</strong> 
                  <span 
                    class="copyable-text" 
                    :title="detailedData.cas_number || 'N/A'"
                    @click="copyToClipboard(detailedData.cas_number)"
                    data-bs-toggle="tooltip" 
                    data-bs-placement="top"
                  >
                    {{ truncateText(detailedData.cas_number, 15) || 'N/A' }}
                  </span>
                  <i class="bi bi-clipboard ms-1 text-muted small" style="cursor: pointer;" 
                     @click="copyToClipboard(detailedData.cas_number)"
                     :title="t('details.click2copy')"></i>
                </p>
                <p>
                  <strong>{{ t('details.formula') }}:</strong> 
                  <span 
                    class="copyable-text" 
                    :title="detailedData.formula || 'N/A'"
                    @click="copyToClipboard(detailedData.formula)"
                    data-bs-toggle="tooltip" 
                    data-bs-placement="top"
                  >
                    {{ truncateText(detailedData.formula, 20) || 'N/A' }}
                  </span>
                  <i class="bi bi-clipboard ms-1 text-muted small" style="cursor: pointer;" 
                     @click="copyToClipboard(detailedData.formula)"
                     :title="t('details.click2copy')"></i>
                </p>
                <p>
                  <strong>{{ t('details.smiles') }}:</strong> 
                  <span 
                    class="copyable-text" 
                    :title="detailedData.smiles || 'N/A'"
                    @click="copyToClipboard(detailedData.smiles)"
                    data-bs-toggle="tooltip" 
                    data-bs-placement="top"
                  >
                    {{ truncateText(detailedData.smiles, 30) || 'N/A' }}
                  </span>
                  <i class="bi bi-clipboard ms-1 text-muted small" style="cursor: pointer;" 
                     @click="copyToClipboard(detailedData.smiles)"
                     :title="t('details.click2copy')"></i>
                </p>
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
            <!-- MS1-H 和 MS1-Na 同行显示 -->
            <div class="row mb-3">
              <div class="col-md-6">
                <div class="d-flex align-items-center">
                  <strong class="me-2">MS1-H:</strong>
                  <span 
                    class="copyable-text flex-grow-1"
                    :title="detailedData.ms1_h || 'N/A'"
                    @click="copyToClipboard(detailedData.ms1_h)"
                    data-bs-toggle="tooltip" 
                    data-bs-placement="top"
                  >
                    {{ detailedData.ms1_h || 'N/A' }}
                  </span>
                  <i class="bi bi-clipboard ms-1 text-muted small" style="cursor: pointer;" 
                     @click="copyToClipboard(detailedData.ms1_h)"
                     :title="t('details.click2copy')"></i>
                </div>
              </div>
              <div class="col-md-6">
                <div class="d-flex align-items-center">
                  <strong class="me-2">MS1-Na:</strong>
                  <span 
                    class="copyable-text flex-grow-1"
                    :title="detailedData.ms1_na || 'N/A'"
                    @click="copyToClipboard(detailedData.ms1_na)"
                    data-bs-toggle="tooltip" 
                    data-bs-placement="top"
                  >
                    {{ detailedData.ms1_na || 'N/A' }}
                  </span>
                  <i class="bi bi-clipboard ms-1 text-muted small" style="cursor: pointer;" 
                     @click="copyToClipboard(detailedData.ms1_na)"
                     :title="t('details.click2copy')"></i>
                </div>
              </div>
            </div>
            
            <!-- Bioactivity 完整显示 -->
            <div class="row mt-3">
              <div class="col-12">
                <div class="bioactivity-section p-3 bg-light rounded">
                  <div class="d-flex justify-content-between align-items-center mb-2">
                    <strong>{{ t('details.bioactivity') }}:</strong>
                    <div>
                      <i class="bi bi-clipboard text-muted small" style="cursor: pointer;" 
                         @click="copyToClipboard(protectedData ? protectedData.bioactivity : '')"
                         :title="t('details.click2copy')"></i>
                    </div>
                  </div>
                  <div v-if="!isAuthenticated" class="text-muted">{{ t("details.not_auth") }}</div>
                  <div v-else-if="protectedData && protectedData.bioactivity" 
                       class="bioactivity-content p-2 bg-white rounded border">
                    {{ truncateText(protectedData.bioactivity, 100) }}
                  </div>
                  <div v-else class="text-muted">N/A</div>
                </div>
              </div>
            </div>
            
            <!-- C13-NMR 完整显示 -->
            <div class="row mt-3">
              <div class="col-12">
                <div class="c13nmr-section p-3 bg-light rounded">
                  <div class="d-flex justify-content-between align-items-center mb-2">
                    <strong>C13-NMR:</strong>
                    <div>
                      <i class="bi bi-clipboard text-muted small" style="cursor: pointer;" 
                         @click="copyToClipboard(protectedData ? protectedData.nmr_13c_data : '')"
                         :title="t('details.click2copy')"></i>
                    </div>
                  </div>
                  <div v-if="!isAuthenticated" class="text-muted">{{ t("details.not_auth") }}</div>
                  <div v-else-if="protectedData && protectedData.nmr_13c_data" 
                       class="c13nmr-content p-2 bg-white rounded border">
                    {{ protectedData.nmr_13c_data }}
                  </div>
                  <div v-else class="text-muted">N/A</div>
                </div>
              </div>
            </div>
            
            <!-- MS2 完整显示 -->
            <div class="row mt-3">
              <div class="col-12">
                <div class="ms2-section p-3 bg-light rounded">
                  <div class="d-flex justify-content-between align-items-center mb-2">
                    <strong>MS2:</strong>
                    <div>
                      <i class="bi bi-download me-2 text-primary" style="cursor: pointer;" 
                         @click="clickToDownload('ms2-full')"
                         :title="t('details.click2download')"></i>
                      <i class="bi bi-clipboard text-muted small" style="cursor: pointer;" 
                         @click="copyToClipboard(protectedData ? protectedData.ms2 : '')"
                         :title="t('details.click2copy')"></i>
                    </div>
                  </div>
                  <div v-if="!isAuthenticated" class="text-muted">{{ t("details.not_auth") }}</div>
                  <div v-else-if="protectedData && protectedData.ms2" 
                       class="ms2-content p-2 bg-white rounded border">
                    {{ protectedData.ms2 }}
                  </div>
                  <div v-else class="text-muted">N/A</div>
                </div>
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
            <div class="col-6 p-1">
              <strong>{{ t('details.tag') }}:</strong> 
              <span 
                :title="detailedData.item_tag || 'N/A'"
                data-bs-toggle="tooltip" 
                data-bs-placement="top"
              >
                {{ truncateText(detailedData.item_tag, 15) || 'N/A' }}
              </span>
            </div>
            <div class="col-6 p-1">
              <strong>{{ t('details.source') }}:</strong> 
              <span 
                :title="detailedData.source || 'N/A'"
                data-bs-toggle="tooltip" 
                data-bs-placement="top"
              >
                {{ truncateText(detailedData.source, 15) || 'N/A' }}
              </span>
            </div>
            <div class="col-6 p-1">
              <strong>FP:</strong> 
              <span 
                :title="detailedData.fp || 'N/A'"
                data-bs-toggle="tooltip" 
                data-bs-placement="top"
              >
                {{ truncateText(detailedData.fp, 12) || 'N/A' }}
              </span>
            </div>
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

/* 可复制文本的样式 */
.copyable-text {
  cursor: pointer;
  transition: all 0.2s ease;
  border-bottom: 1px dotted #ccc;
  padding: 0 2px;
}

.copyable-text:hover {
  background-color: #f0f8ff;
  border-bottom-color: #007bff;
}

/* 复制图标样式 */
.bi-clipboard {
  opacity: 0.6;
  transition: opacity 0.2s ease;
}

.bi-clipboard:hover {
  opacity: 1;
  color: #007bff;
}
</style>
