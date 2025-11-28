<script setup>
import { ref, watch, onMounted, computed } from 'vue';
// 引入 OpenChemLib 完整版 (包含坐标生成算法)
import OCL from 'openchemlib';

const props = defineProps({
  smiles: {
    type: String,
    default: ''
  },
  // 可选：允许外部控制尺寸，但默认自动适应容器
  width: {
    type: [Number, String],
    default: 300
  },
  height: {
    type: [Number, String],
    default: 200
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
    // 1. 解析 SMILES
    const mol = OCL.Molecule.fromSmiles(props.smiles);
    
    const svg = mol.toSVG(
      parseFloat(props.width), 
      parseFloat(props.height), 
      null, 
      {
        strokeWidth: 0.8,
        fontWeight: 'bold'
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
  renderStructure();
}, { immediate: true });

// 监听尺寸变化 (可选)
watch([() => props.width, () => props.height], () => {
  renderStructure();
});

onMounted(() => {
  renderStructure();
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
  min-height: 150px; /* 保持最小高度 */
  background-color: #fff;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  overflow: hidden;
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
}

.no-structure-message i {
  font-size: 2rem;
  opacity: 0.5;
}

.text-danger {
  color: #dc3545;
}
</style>