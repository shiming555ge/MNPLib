<script setup>
import { ref, watch, onMounted, nextTick } from 'vue';
// 引入 OpenChemLib 完整版 (包含坐标生成算法)
import OCL from 'openchemlib';

const props = defineProps({
  smiles: {
    type: String,
    default: ''
  }
});

const containerRef = ref(null);
const svgContent = ref('');
const hasError = ref(false);

/** 绘制逻辑 */
const renderStructure = () => {
  hasError.value = false;
  svgContent.value = '';

  if (!props.smiles || typeof props.smiles !== 'string' || props.smiles.trim() === '') {
    return;
  }

  try {
    // 获取容器实际尺寸
    const container = containerRef.value;
    if (!container) return;
    
    const width = container.clientWidth || 200;
    const height = container.clientHeight || 120;
    
    // 1. 解析 SMILES
    const mol = OCL.Molecule.fromSmiles(props.smiles);
    
    // 移除未知手性文字和不显示？原子
    const svg = mol.toSVG(
      width, 
      height, 
      null, 
      {
        strokeWidth: 0.8,
        fontWeight: 'bold',
        suppressChiralText: true,  // 移除unknown chirality文字
        suppressCIPParity: true    // 移除手性标记
      }
    );

    svgContent.value = svg;

  } catch (error) {
    console.error('OCL Render Error:', error);
    hasError.value = true;
  }
};

// 监听 SMILES 变化
watch(() => props.smiles, () => {
  nextTick(() => {
    renderStructure();
  });
}, { immediate: true });

onMounted(() => {
  nextTick(() => {
    renderStructure();
  });
});
</script>

<template>
  <div class="molecule-wrapper" ref="containerRef">
    <div 
      v-if="svgContent && !hasError" 
      class="svg-container" 
      v-html="svgContent"
    ></div>

    <div v-else class="no-structure-message">
      <i class="bi bi-question-circle" :class="{ 'text-danger': hasError }"></i>
      <p class="small text-muted mt-1">
        {{ hasError ? '结构解析失败' : '无结构数据' }}
      </p>
    </div>
  </div>
</template>

<style scoped>
.molecule-wrapper {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background-color: #fff;
  border-radius: 4px;
  overflow: hidden;
  aspect-ratio: 7/3; /* 保持宽高比 7:3 */
}

.svg-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 穿透控制 SVG 内部样式 (可选) */
:deep(svg) {
  width: 100%;
  height: 100%;
  /* 保持比例 */
  object-fit: contain; 
}

.no-structure-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #6c757d;
  width: 100%;
  height: 100%;
}

.no-structure-message i {
  font-size: 2rem;
  opacity: 0.5;
}

.text-danger {
  color: #dc3545;
}
</style>
