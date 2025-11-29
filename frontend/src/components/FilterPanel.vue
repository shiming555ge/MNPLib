<script setup>
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const props = defineProps({
  filters: { type: Object, required: true },
  itemTypes: { type: Array, default: () => [] },
  descriptions: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false }
});

const emit = defineEmits(['apply', 'reset']);

// --- 折叠/展开 逻辑 ---
const ITEM_LIMIT = 10; // 默认显示的个数
const showAllTypes = ref(false);
const showAllDescs = ref(false);

// 计算显示的 ItemTypes
const displayedItemTypes = computed(() => {
  if (showAllTypes.value) return props.itemTypes;
  return props.itemTypes.slice(0, ITEM_LIMIT);
});

// 计算显示的 Descriptions
const displayedDescriptions = computed(() => {
  if (showAllDescs.value) return props.descriptions;
  return props.descriptions.slice(0, ITEM_LIMIT);
});

// 切换选择
const toggleFilterItem = (array, item) => {
  const index = array.indexOf(item);
  if (index > -1) array.splice(index, 1);
  else array.push(item);
};

const hasActiveFilters = computed(() => {
  return props.filters.item_type.length > 0 || 
         props.filters.description.length > 0 || 
         props.filters.min_weight || 
         props.filters.max_weight;
});
</script>

<template>
  <div class="filter-panel">
    <div class="mb-4">
      <div class="d-flex justify-content-between align-items-center mb-2">
        <label class="form-label fw-bold small text-uppercase text-secondary mb-0">
          <i class="bi bi-diagram-3 me-1"></i> {{ t('browse.compound_type') }}
        </label>
        <button 
          v-if="itemTypes.length > ITEM_LIMIT"
          class="btn btn-link btn-sm p-0 text-decoration-none"
          style="font-size: 0.8rem;"
          @click="showAllTypes = !showAllTypes"
        >
          {{ showAllTypes ? t('browse.collapse') || '收起' : t('browse.expand') || '展开' }}
          <i class="bi" :class="showAllTypes ? 'bi-chevron-up' : 'bi-chevron-down'"></i>
        </button>
      </div>
      
      <div class="d-flex flex-wrap gap-2">
        <div 
          v-for="type in displayedItemTypes" 
          :key="type"
          class="filter-chip"
          :class="{ 'active bg-primary text-white border-primary': filters.item_type.includes(type) }"
          @click="toggleFilterItem(filters.item_type, type)"
        >
          {{ type }}
          <i v-if="filters.item_type.includes(type)" class="bi bi-check-lg ms-1"></i>
        </div>
        <span v-if="!showAllTypes && itemTypes.length > ITEM_LIMIT" class="badge text-secondary bg-light border align-self-center">
            +{{ itemTypes.length - ITEM_LIMIT }}
        </span>
      </div>
    </div>

    <hr class="border-secondary opacity-10 my-4">

    <div class="mb-4">
      <label class="form-label fw-bold small text-uppercase text-secondary mb-2">
        <i class="bi bi-speedometer2 me-1"></i> {{ t('browse.molecular_weight_range') }}
      </label>
      <div class="input-group input-group-sm">
        <input type="number" class="form-control" :placeholder="t('browse.min_weight')" v-model="filters.min_weight">
        <span class="input-group-text bg-light text-muted border-start-0 border-end-0">-</span>
        <input type="number" class="form-control" :placeholder="t('browse.max_weight')" v-model="filters.max_weight">
      </div>
    </div>

    <hr class="border-secondary opacity-10 my-4">

    <div class="mb-4">
      <div class="d-flex justify-content-between align-items-center mb-2">
        <label class="form-label fw-bold small text-uppercase text-secondary mb-0">
          <i class="bi bi-tags me-1"></i> {{ t('browse.description_category') }}
        </label>
        <button 
          v-if="descriptions.length > ITEM_LIMIT"
          class="btn btn-link btn-sm p-0 text-decoration-none"
          style="font-size: 0.8rem;"
          @click="showAllDescs = !showAllDescs"
        >
          {{ showAllDescs ? t('browse.collapse') || '收起' : t('browse.expand') || '展开' }}
          <i class="bi" :class="showAllDescs ? 'bi-chevron-up' : 'bi-chevron-down'"></i>
        </button>
      </div>

      <div class="d-flex flex-wrap gap-2">
        <div 
          v-for="desc in displayedDescriptions" 
          :key="desc"
          class="filter-chip desc-chip"
          :class="{ 'active bg-success text-white border-success': filters.description.includes(desc) }"
          @click="toggleFilterItem(filters.description, desc)"
        >
          {{ desc }}
          <i v-if="filters.description.includes(desc)" class="bi bi-check-lg ms-1"></i>
        </div>
         <span v-if="!showAllDescs && descriptions.length > ITEM_LIMIT" class="badge text-secondary bg-light border align-self-center">
            +{{ descriptions.length - ITEM_LIMIT }}
        </span>
      </div>
    </div>

    <div class="d-grid gap-2 mt-5 pb-2">
      <button class="btn btn-primary shadow-sm fw-medium" @click="$emit('apply')" :disabled="loading" data-bs-dismiss="offcanvas">
        <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
        <i v-else class="bi bi-search me-1"></i> {{ t('browse.apply_filters') }}
      </button>
      <button class="btn btn-outline-secondary border-0" @click="$emit('reset')" :disabled="loading || !hasActiveFilters">
        {{ t('browse.reset_filters') }}
      </button>
    </div>
  </div>
</template>

<style scoped>
/* 保持之前的 chip 样式 */
.filter-chip {
  cursor: pointer;
  padding: 0.35rem 0.75rem;
  border-radius: 50rem;
  border: 1px solid #dee2e6;
  background-color: #fff;
  color: #495057;
  font-size: 0.85rem;
  transition: all 0.2s ease;
  user-select: none;
}
.filter-chip:hover {
  background-color: #f8f9fa;
  border-color: #c5cdd4;
  transform: translateY(-1px);
}
.filter-chip.active:hover {
  opacity: 0.9;
  transform: translateY(0);
}
</style>
