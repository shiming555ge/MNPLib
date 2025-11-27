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
    <div class="col-8 col-sm-6 col-xl-3">
        <div class="card h-100 p-3">
            <!-- 分子结构显示区域 -->
            <div class="structure-container text-center">
                <MoleculeCanvas :smiles="compound.smiles"></MoleculeCanvas>
            </div>
            <div class="card-body d-flex flex-column">
                <h5 class="card-title">{{ compound.item_name ? compound.item_name.replace(/"/g, "") : '未命名化合物' }}</h5>
                <p class="card-subtitle text-muted no-warp text-truncate">{{ compound.cas_number || '无CAS号' }}</p>
                <a class="card-link mt-auto" data-bs-toggle="offcanvas" href="#detailOffcanvas" role="button" aria-controls="detailOffcanvas" @click="showDetail">
                    {{ t("compound_card.details") }}
                </a>
            </div>
        </div>
    </div>
</template>

<style scoped>


.structure-canvas {
    max-width: 100%;
    max-height: 100%;
    border: 1px solid #e9ecef;
    border-radius: 4px;
}

</style>
