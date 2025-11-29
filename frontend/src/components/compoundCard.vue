<script setup>
import { useI18n } from 'vue-i18n';
import MoleculeCanvas from './MoleculeCanvas.vue';

const { t } = useI18n();

const props = defineProps({
  compound: {
    type: Object,
    required: true,
    default: () => ({
      ID: 'CMP000',
      ItemName: 'Compound Name',
      CASNumber: 'N/A',
      smiles: null
    })
  }
});

const emit = defineEmits(['show-detail']);
const showDetail = () => emit('show-detail');
</script>

<template>
    <div class="col-12 col-sm-6 col-lg-4 col-xl-3">
        <div class="card h-100 shadow-sm border-0">
            <!-- 分子结构显示区域 -->
            <div class="structure-container text-center bg-light rounded-top">
                <MoleculeCanvas :smiles="compound.smiles"></MoleculeCanvas>
            </div>
            <div class="card-body d-flex flex-column p-3">
                <h6 class="card-title fw-bold text-truncate mb-1" :title="compound.item_name || '未命名化合物'">
                    {{ compound.item_name ? compound.item_name.replace(/"/g, "") : '未命名化合物' }}
                </h6>
                <p class="card-subtitle text-muted small mb-2 text-truncate">
                    <i class="bi bi-tag me-1"></i>{{ compound.cas_number || '无CAS号' }}
                </p>
                <div class="mt-auto">
                    <button class="btn btn-outline-primary btn-sm w-100" @click="showDetail">
                        <i class="bi bi-eye me-1"></i>{{ t("compound_card.details") }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.card {
    transition: all 0.3s ease;
    border-radius: 12px;
}

.card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

.structure-container {
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f8f9fa;
    overflow: hidden;
    padding: 0.5rem !important;
}

.card-title {
    font-size: 0.95rem;
    line-height: 1.3;
}

.card-subtitle {
    font-size: 0.8rem;
}

.btn {
    border-radius: 8px;
    font-size: 0.85rem;
    padding: 0.4rem 0.75rem;
}

/* 小屏设备优化 */
@media (max-width: 576px) {
    .col-12 {
        margin-bottom: 1rem;
    }
    
    .card {
        border-radius: 10px;
    }
    
    .structure-container {
        min-height: 100px;
        padding: 1rem !important;
    }
    
    .card-body {
        padding: 1rem !important;
    }
    
    .card-title {
        font-size: 0.9rem;
    }
    
    .btn {
        font-size: 0.8rem;
        padding: 0.35rem 0.5rem;
    }
}

/* 中等屏幕优化 */
@media (min-width: 576px) and (max-width: 768px) {
    .structure-container {
        min-height: 110px;
    }
}

.structure-canvas {
    max-width: 100%;
    max-height: 100%;
    border-radius: 6px;
}
</style>
