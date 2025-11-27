<script setup>
import { ref, watch } from 'vue'
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

// 手动关闭方法
const handleClose = () => {
  emit('update:show', false)
}
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
        {{ compound ? compound.item_name?.replace(/"/g, "") || compound.ItemName?.replace(/"/g, "") || t('browse.compound_details') : t('browse.compound_details') }}
      </h5>
      <button type="button" class="btn-close" @click="handleClose" aria-label="Close"></button>
    </div>
    <div class="offcanvas-body">
      <div v-if="compound" class="compound-detail">
        <!-- 绘制分子结构 -->
        <div class="card mb-3">
          <div class="card-header bg-primary text-white">
            <h6 class="card-title mb-0"><i class="bi bi-info-circle"></i> {{ t('details.structure') }}</h6>
          </div>
          <div class="card-body text-center">
            <MoleculeViewer 
              :smiles="compound.smiles" 
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
                <p><strong>{{ t('details.id') }}:</strong> {{ compound.id || compound.ID || 'N/A' }}</p>
                <p><strong>{{ t('details.source') }}:</strong> {{ compound.source || compound.Source || 'N/A' }}</p>
                <p><strong>{{ t('details.compound_name') }}:</strong> {{ compound.item_name || compound.ItemName || 'N/A' }}</p>
                <p><strong>{{ t('details.type') }}:</strong> {{ compound.item_type || compound.ItemType || 'N/A' }}</p>
                <p><strong>{{ t('details.iupac_name') }}:</strong> {{ compound.iupac_name || compound.IUPACName || 'N/A' }}</p>
              </div>
              <div class="col-md-12">
                <p><strong>{{ t('details.description') }}:</strong> {{ compound.description || compound.Description || 'N/A' }}</p>
                <p><strong>{{ t('details.cas_number') }}:</strong> {{ compound.cas_number || compound.CASNumber || 'N/A' }}</p>
                <p><strong>{{ t('details.formula') }}:</strong> {{ compound.formula || compound.Formula || 'N/A' }}</p>
                <p><strong>{{ t('details.smiles') }}:</strong> {{ compound.smiles || compound.Smiles || 'N/A' }}</p>
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
                <p><strong>MS1:</strong> {{ compound.ms1 || compound.MS1 || 'N/A' }}</p>
                <p><strong>MS2:</strong> {{ compound.ms2 || compound.MS2 || 'N/A' }}</p>
              </div>
              <div class="col-md-6">
                <p><strong>{{ t('details.bioactivity') }}:</strong> {{ compound.bioactivity || compound.Bioactivity || 'N/A' }}</p>
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
            <div class="col-6 p-1"><strong>{{ t('details.tag') }}:</strong> {{ compound.item_tag || compound.ItemTag || 'N/A' }}</div>
            <div class="col-6 p-1"><strong>{{ t('details.structure') }}:</strong> {{ compound.structure || compound.Structure || 'N/A' }}</div>
            <div class="col-6 p-1"><strong>{{ t('details.created_at') }}:</strong> {{ compound.created_at || compound.CreatedAt || 'N/A' }}</div>
            <div class="col-6 p-1"><strong>{{ t('details.updated_at') }}:</strong> {{ compound.updated_at || compound.UpdatedAt || 'N/A' }}</div>
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
