<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue';
import * as THREE from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js';

const props = defineProps({
  smiles: {
    type: String,
    default: ''
  },
  width: {
    type: [Number, String],
    default: 400
  },
  height: {
    type: [Number, String],
    default: 300
  },
  backgroundColor: {
    type: String,
    default: '#ffffff'
  }
});

const containerRef = ref(null);
const hasError = ref(false);
const isLoading = ref(false);

// Three.js 相关变量
let scene = null;
let camera = null;
let renderer = null;
let controls = null;
let moleculeGroup = null;

// 原子颜色映射
const atomColors = {
  'C': 0x909090, // 碳 - 灰色
  'H': 0xffffff, // 氢 - 白色
  'O': 0xff0000, // 氧 - 红色
  'N': 0x0000ff, // 氮 - 蓝色
  'S': 0xffff00, // 硫 - 黄色
  'P': 0xff8000, // 磷 - 橙色
  'F': 0x00ff00, // 氟 - 绿色
  'Cl': 0x00ff00, // 氯 - 绿色
  'Br': 0x800000, // 溴 - 深红色
  'I': 0x800080  // 碘 - 紫色
};

// 原子半径映射 (单位: Å)
const atomRadii = {
  'C': 0.7,
  'H': 0.3,
  'O': 0.66,
  'N': 0.65,
  'S': 1.04,
  'P': 1.0,
  'F': 0.57,
  'Cl': 1.0,
  'Br': 1.14,
  'I': 1.33
};

// 键半径
const bondRadius = 0.1;

/** 初始化 Three.js 场景 */
const initScene = () => {
  if (!containerRef.value) return;

  // 清理现有场景
  if (renderer) {
    containerRef.value.removeChild(renderer.domElement);
  }

  // 创建场景
  scene = new THREE.Scene();
  scene.background = new THREE.Color(props.backgroundColor);

  // 创建相机
  const aspect = parseFloat(props.width) / parseFloat(props.height);
  camera = new THREE.PerspectiveCamera(45, aspect, 0.1, 1000);
  camera.position.z = 10;

  // 创建渲染器
  renderer = new THREE.WebGLRenderer({ antialias: true });
  renderer.setSize(parseFloat(props.width), parseFloat(props.height));
  renderer.setPixelRatio(window.devicePixelRatio);
  containerRef.value.appendChild(renderer.domElement);

  // 添加轨道控制器
  controls = new OrbitControls(camera, renderer.domElement);
  controls.enableDamping = true;
  controls.dampingFactor = 0.05;

  // 添加环境光
  const ambientLight = new THREE.AmbientLight(0x404040, 0.6);
  scene.add(ambientLight);

  // 添加方向光
  const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8);
  directionalLight.position.set(1, 1, 1);
  scene.add(directionalLight);

  // 创建分子组
  moleculeGroup = new THREE.Group();
  scene.add(moleculeGroup);
};

/** 创建原子球体 */
const createAtom = (element, position) => {
  const radius = atomRadii[element] || 0.7;
  const color = atomColors[element] || 0x808080;
  
  const geometry = new THREE.SphereGeometry(radius, 32, 32);
  const material = new THREE.MeshPhongMaterial({ 
    color: color,
    shininess: 30
  });
  
  const atom = new THREE.Mesh(geometry, material);
  atom.position.set(position.x, position.y, position.z);
  
  return atom;
};

/** 创建键圆柱体 */
const createBond = (start, end) => {
  const direction = new THREE.Vector3().subVectors(end, start);
  const length = direction.length();
  
  const geometry = new THREE.CylinderGeometry(bondRadius, bondRadius, length, 16);
  const material = new THREE.MeshPhongMaterial({ 
    color: 0x666666,
    shininess: 30
  });
  
  const bond = new THREE.Mesh(geometry, material);
  
  // 设置位置和方向
  bond.position.copy(start).add(end).multiplyScalar(0.5);
  bond.lookAt(end);
  bond.rotateX(Math.PI / 2);
  
  return bond;
};

/** 解析PDB文件内容 */
const parsePDB = (pdbContent) => {
  const atoms = [];
  const bonds = [];
  const connections = [];
  
  const lines = pdbContent.split('\n');
  
  // 解析原子坐标
  lines.forEach(line => {
    if (line.startsWith('ATOM') || line.startsWith('HETATM')) {
      // 从HETATM记录中提取原子信息
      // 格式: HETATM    1  C1  UNL     1      -1.140  -0.804   0.101  1.00  0.00           C
      const atomId = parseInt(line.substring(6, 11).trim());
      const element = line.substring(76, 78).trim();
      const x = parseFloat(line.substring(30, 38));
      const y = parseFloat(line.substring(38, 46));
      const z = parseFloat(line.substring(46, 54));
      
      atoms.push({
        id: atomId,
        element: element,
        position: new THREE.Vector3(x, y, z)
      });
    } else if (line.startsWith('CONECT')) {
      // 解析连接记录
      // 格式: CONECT    1    2    2    6
      const parts = line.trim().split(/\s+/);
      const centralAtom = parseInt(parts[1]);
      const connectedAtoms = parts.slice(2).map(Number);
      
      connections.push({
        central: centralAtom,
        connected: connectedAtoms
      });
    }
  });
  
  // 使用CONECT记录创建键
  connections.forEach(conn => {
    const centralAtom = atoms.find(atom => atom.id === conn.central);
    if (centralAtom) {
      conn.connected.forEach(connectedId => {
        const connectedAtom = atoms.find(atom => atom.id === connectedId);
        if (connectedAtom) {
          // 避免重复创建键
          const bondExists = bonds.some(bond => 
            (bond.startAtomId === centralAtom.id && bond.endAtomId === connectedAtom.id) ||
            (bond.startAtomId === connectedAtom.id && bond.endAtomId === centralAtom.id)
          );
          
          if (!bondExists) {
            bonds.push({
              start: centralAtom.position,
              end: connectedAtom.position,
              startAtomId: centralAtom.id,
              endAtomId: connectedAtom.id
            });
          }
        }
      });
    }
  });
  
  // 如果没有CONECT记录，使用距离阈值作为备选方案
  if (bonds.length === 0) {
    const bondThreshold = 1.8; // Å
    
    for (let i = 0; i < atoms.length; i++) {
      for (let j = i + 1; j < atoms.length; j++) {
        const distance = atoms[i].position.distanceTo(atoms[j].position);
        if (distance < bondThreshold) {
          bonds.push({
            start: atoms[i].position,
            end: atoms[j].position,
            startAtomId: atoms[i].id,
            endAtomId: atoms[j].id
          });
        }
      }
    }
  }
  
  return { atoms, bonds };
};

/** 从API获取PDB数据 */
const fetchPDBFromAPI = async (smiles) => {
  try {
    const response = await fetch(`/api/rdkit/smiles-to-pdb?smiles=${encodeURIComponent(smiles)}`);
    
    if (!response.ok) {
      throw new Error(`API请求失败: ${response.status}`);
    }
    
    const result = await response.json()
    console.log('API响应结果:', result)
    
    // 处理API响应格式
    if (result.code === 200200 && result.data) {
      return result.data
    }
  } catch (error) {
    console.error('获取PDB数据失败:', error);
    throw error;
  }
};

/** 生成3D分子结构 */
const generate3DStructure = async (smiles) => {
  // 从API获取真实的3D结构
  const pdbData = await fetchPDBFromAPI(smiles);
  console.log('获取到的PDB数据:', pdbData);
  
  // 解析PDB数据
  const structure = parsePDB(pdbData);
  console.log('解析后的结构:', structure);
  
  return structure;
};

/** 渲染3D分子结构 */
const render3DStructure = async () => {
  hasError.value = false;
  isLoading.value = true;

  if (!props.smiles || typeof props.smiles !== 'string' || props.smiles.trim() === '') {
    clearScene();
    isLoading.value = false;
    return;
  }

  try {
    // 初始化场景
    initScene();
    
    // 生成3D结构
    const structure = await generate3DStructure(props.smiles);
    
    // 清除现有分子
    moleculeGroup.clear();
    
    // 创建原子
    structure.atoms.forEach(atomData => {
      const atom = createAtom(atomData.element, atomData.position);
      moleculeGroup.add(atom);
    });
    
    // 创建键
    structure.bonds.forEach(bondData => {
      const bond = createBond(bondData.start, bondData.end);
      moleculeGroup.add(bond);
    });
    
    // 开始动画循环
    animate();
    
  } catch (error) {
    console.error('3D Render Error:', error);
    hasError.value = true;
  } finally {
    isLoading.value = false;
  }
};

/** 清除场景 */
const clearScene = () => {
  if (moleculeGroup) {
    moleculeGroup.clear();
  }
};

/** 动画循环 */
const animate = () => {
  requestAnimationFrame(animate);
  
  if (controls) {
    controls.update();
  }
  
  if (renderer && scene && camera) {
    renderer.render(scene, camera);
  }
};

/** 处理窗口大小变化 */
const handleResize = () => {
  if (camera && renderer) {
    const width = parseFloat(props.width);
    const height = parseFloat(props.height);
    
    camera.aspect = width / height;
    camera.updateProjectionMatrix();
    renderer.setSize(width, height);
  }
};

// 监听属性变化
watch(() => props.smiles, () => {
  render3DStructure();
}, { immediate: true });

watch([() => props.width, () => props.height, () => props.backgroundColor], () => {
  handleResize();
  if (scene) {
    scene.background = new THREE.Color(props.backgroundColor);
  }
});

onMounted(() => {
  render3DStructure();
});

onUnmounted(() => {
  // 清理资源
  if (renderer) {
    renderer.dispose();
  }
});
</script>

<template>
  <div class="molecule-viewer-wrapper">
    <div 
      ref="containerRef" 
      class="viewer-container"
      :style="{
        width: typeof width === 'number' ? width + 'px' : width,
        height: typeof height === 'number' ? height + 'px' : height
      }"
    ></div>
    
    <div v-if="isLoading" class="loading-overlay">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">加载中...</span>
      </div>
      <p class="mt-2">正在生成3D结构...</p>
    </div>
    
    <div v-else-if="hasError" class="error-overlay">
      <i class="bi bi-exclamation-triangle text-danger"></i>
      <p class="mt-2 text-danger">3D结构生成失败</p>
    </div>
    
    <div v-else-if="!smiles" class="no-data-overlay">
      <i class="bi bi-molecule text-muted"></i>
      <p class="mt-2 text-muted">无分子数据</p>
    </div>
  </div>
</template>

<style scoped>
.molecule-viewer-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  overflow: hidden;
  background-color: v-bind('backgroundColor');
}

.viewer-container {
  width: 100%;
  height: 100%;
  min-height: 200px;
}

.loading-overlay,
.error-overlay,
.no-data-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: rgba(255, 255, 255, 0.9);
  z-index: 10;
}

.loading-overlay i,
.error-overlay i,
.no-data-overlay i {
  font-size: 2rem;
}

.error-overlay i {
  color: #dc3545;
}

.no-data-overlay i {
  opacity: 0.5;
}
</style>
