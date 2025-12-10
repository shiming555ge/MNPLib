<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import MoleculeCanvas from './MoleculeCanvas.vue'
import CompoundDetail from './CompoundDetail.vue'

const { t } = useI18n()

// 搜索模式状态
const searchMode = ref('structure') // structure, substructure, similarity, c-nmr, ms2
const ketcherRef = ref(null)
const searchResults = ref([])
const loading = ref(false)
const errorMessage = ref('')
const currentSmiles = ref('')
const selectedCompound = ref(null)
const showDetail = ref(false)

// 核磁谱搜索相关
const nmrText = ref('')
const nmrFile = ref(null)
const threshold = ref(0.5)
const tolerance = ref(0.5)

// MS2搜索相关
const ms2Text = ref('')
const ms2File = ref(null)
const prefilterThreshold = ref(0.25) // 默认值为0.5的一半

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
    if (!ketcher) {
      console.error('Ketcher实例未找到')
      errorMessage.value = t('query.ketcher_not_ready')
      return null
    }
    
    const smiles = await ketcher.getSmiles()
    console.log('从Ketcher获取的SMILES:', smiles)
    
    // 验证SMILES是否有效
    if (!smiles || smiles.trim() === '') {
      console.error('获取的SMILES为空')
      errorMessage.value = t('query.empty_structure')
      return null
    }
    
    // 检查是否为有效的SMILES格式（简单检查）
    if (typeof smiles !== 'string' || smiles.length < 2) {
      console.error('SMILES格式无效:', smiles)
      errorMessage.value = t('query.invalid_structure')
      return null
    }
    
    currentSmiles.value = smiles
    return smiles
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

// 读取文本文件
const readTextFile = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => resolve(e.target.result)
    reader.onerror = (e) => reject(e)
    reader.readAsText(file)
  })
}

// 处理文件上传
const handleFileUpload = async (event, type = 'nmr') => {
  const file = event.target.files[0]
  if (!file) return
  
  try {
    const content = await readTextFile(file)
    if (type === 'nmr') {
      nmrText.value = content
    } else if (type === 'ms2') {
      ms2Text.value = content
    }
  } catch (error) {
    console.error('读取文件失败:', error)
    errorMessage.value = t('query.file_read_failed', { error: error.message })
  }
}

// 处理搜索操作
const handleSearch = async () => {
  loading.value = true
  errorMessage.value = ''
  searchResults.value = []
  currentPage.value = 1

  try {
    // 验证MS2搜索参数
    if (searchMode.value === 'ms2') {
      // 验证Prefilter Threshold和Similarity threshold的关系
      if (prefilterThreshold.value > threshold.value) {
        errorMessage.value = t('query.prefilter_threshold_error', { 
          prefilter: prefilterThreshold.value.toFixed(2), 
          similarity: threshold.value.toFixed(2) 
        })
        loading.value = false
        return
      }
      
      // 验证阈值范围
      if (threshold.value < 0 || threshold.value > 1) {
        errorMessage.value = t('query.threshold_range_error')
        loading.value = false
        return
      }
      
      if (prefilterThreshold.value < 0 || prefilterThreshold.value > 1) {
        errorMessage.value = t('query.prefilter_range_error')
        loading.value = false
        return
      }
    }

    if (searchMode.value === 'c-nmr') {
      // 核磁谱搜索
      if (!nmrText.value.trim()) {
        errorMessage.value = t('query.enter_nmr_data_or_upload')
        loading.value = false
        return
      }
      
      const response = await fetch(`/api/rdkit/nmr-search?query_nmr=${encodeURIComponent(nmrText.value)}&threshold=${threshold.value}&tolerance=${tolerance.value}`)
      
      if (!response.ok) {
        throw new Error(`API请求失败: ${response.status}`)
      }

      const result = await response.json()
      console.log('核磁谱搜索API响应结果:', result)
      
      // 处理API响应格式
      if (result.code === 200200 && result.data) {
        try {
          const parsedData = JSON.parse(result.data)
          console.log('解析出的核磁谱搜索结果:', parsedData)
          
          let compoundIds = []
          
          // 核磁谱搜索：数据格式为 [["CMP0005", 0.8], ["CMP0054", 0.7]]
          if (Array.isArray(parsedData) && parsedData.length > 0) {
            compoundIds = parsedData.map(item => item[0]) // 提取ID
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
          console.error('解析核磁谱数据失败:', parseError)
          searchResults.value = []
        }
      } else {
        searchResults.value = []
      }
    } else if (searchMode.value === 'ms2') {
      // MS2相似度搜索 - 使用客户端指纹计算
      if (!ms2Text.value.trim()) {
        ms2Text.value = "energy0\n271.06010 15.81 45 (2.9501)\n273.07575 50.41 49 53 64 62 61 60 55 58 57 56 92 (7.0903 1.1007 0.33349 0.31693 0.24094 0.11569 0.11206 0.052208 0.03833 0.0080716 1.1611e-05)\n275.09140 35.02 65 99 71 14 96 73 94 88 93 77 82 91 83 76 86 (4.645 1.3746 0.20302 0.070415 0.068653 0.036737 0.030711 0.027029 0.025537 0.013149 0.011772 0.01096 0.0083248 0.005589 0.0044003)\n277.10705 19.46 101 (3.6316)\n279.12270 65.26 103 118 108 117 115 114 116 109 112 110 111 166 (11.901 0.12046 0.059233 0.041516 0.025318 0.01447 0.0041627 0.0037119 0.0032001 0.0031289 0.0028009 3.376e-06)\n281.13835 60.04 119 135 132 134 124 125 128 130 131 126 127 (10.786 0.36525 0.011662 0.010759 0.0082765 0.0074821 0.0051479 0.0040924 0.0031439 0.0019514 0.0013896)\n299.09140 29.56 140 174 (5.5149 0.0025352)\n317.10196 30.63 173 175 183 182 177 179 178 (4.6922 0.78196 0.14486 0.044353 0.027232 0.01716 0.0086157)\n355.15400 12.07 136 143 192 16 257 18 141 (1.5294 0.41971 0.12532 0.087646 0.050999 0.034418 0.0060273)\n373.16456 100.00 0 193 287 275 198 224 256 259 254 216 249 252 213 226 (16.058 1.2157 0.81227 0.31953 0.064396 0.046603 0.03934 0.031753 0.025731 0.024811 0.0087004 0.0075604 0.0066244 0.0016569)\nenergy1\n55.05423 5.42 32 147 (0.86092 0.42807)\n65.03858 6.17 44 (1.4683)\n69.06988 7.44 31 (1.7702)\n81.06988 14.95 37 (3.557)\n83.08553 9.13 30 (2.1728)\n85.10118 14.19 22 (3.3757)\n93.03349 4.88 43 (1.1605)\n135.11683 7.03 265 (1.6733)\n219.02880 5.13 95 269 (1.2159 0.0040504)\n221.04445 4.40 100 (1.0461)\n243.06519 19.56 6 87 (4.6498 0.0050997)\n247.06010 14.86 63 273 (3.4605 0.075913)\n247.09649 5.57 66 84 75 (1.019 0.21666 0.090447)\n257.08084 6.25 70 17 67 (0.64099 0.50921 0.33669)\n259.06010 13.10 2 68 69 (1.9458 0.85692 0.31498)\n273.07575 23.04 49 53 64 56 60 62 57 61 58 55 92 (4.2745 0.48848 0.19544 0.16148 0.10389 0.097638 0.058492 0.055929 0.025423 0.018002 0.0040935)\n275.09140 100.00 65 99 96 71 76 94 93 73 77 83 91 88 82 86 14 (17.92 2.9315 0.78767 0.45717 0.35829 0.32573 0.23679 0.21329 0.13089 0.099135 0.08903 0.067998 0.066015 0.065826 0.046886)\n277.10705 4.41 101 (1.0501)\n279.12270 12.03 103 118 115 117 110 116 108 112 114 109 111 166 (2.6587 0.060794 0.042565 0.029404 0.028902 0.012736 0.011731 0.0071067 0.0058836 0.0031198 0.0013054 0.0009384)\n281.13835 9.32 119 135 134 125 132 124 128 130 131 126 127 (2.0616 0.055308 0.046042 0.015413 0.01341 0.011671 0.0045704 0.004069 0.0031131 0.0026323 0.00063113)\n299.09140 6.39 140 174 (1.5087 0.011822)\n305.10196 6.96 151 157 155 154 153 152 (1.0824 0.2701 0.15763 0.11724 0.015745 0.012406)\n311.09140 3.81 186 (0.9058)\n315.12270 4.54 190 221 (1.0759 0.0048568)\n317.10196 5.40 173 175 183 182 177 179 178 (0.77711 0.21349 0.17155 0.07804 0.025173 0.011967 0.0075643)\n345.13326 5.02 279 (1.1937)\n373.16456 13.88 0 193 275 259 287 224 213 252 254 198 216 249 226 256 (1.484 0.5231 0.43611 0.24937 0.16677 0.14815 0.098375 0.061532 0.056654 0.033009 0.020285 0.011845 0.0069774 0.006148)\nenergy2\n39.02293 22.07 171 (2.6539)\n41.03858 100.00 33 (12.027)\n43.05423 17.20 27 (2.0682)\n53.03858 23.73 148 (2.8545)\n55.05423 54.00 147 32 (4.6013 1.8934)\n57.03349 7.65 81 34 (0.66921 0.25071)\n57.06988 14.51 146 23 (1.58 0.16543)\n69.06988 10.90 31 (1.3115)\n81.06988 51.94 37 (6.2476)\n83.08553 34.18 30 (4.1111)\n125.05971 10.95 85 (1.3165)\n243.06519 8.55 6 87 (1.0079 0.020658)\n245.08084 14.53 59 11 (1.46 0.28734)\n247.06010 18.42 63 273 (2.2148 0.0011023)\n253.04954 7.18 47 (0.86386)\n255.06519 8.65 51 50 (0.60718 0.43283)\n257.08084 12.16 17 70 67 (0.58136 0.44624 0.43513)\n259.06010 8.47 2 69 68 (0.54946 0.32542 0.14384)\n273.07575 32.48 64 49 57 60 61 56 55 62 58 92 53 (0.79898 0.72379 0.53677 0.42095 0.38328 0.31433 0.27206 0.23707 0.097845 0.061922 0.060014)\n275.09140 16.03 65 14 99 77 96 82 88 93 73 83 86 71 91 76 94 (0.71118 0.34821 0.2212 0.21317 0.11493 0.072851 0.03795 0.036783 0.034656 0.034093 0.033468 0.027057 0.015765 0.014873 0.011478)\n277.07066 9.18 156 (1.1042)\n279.12270 10.19 166 115 111 109 114 108 116 118 110 117 103 112 (0.25777 0.1503 0.12477 0.12188 0.10482 0.080141 0.078731 0.077139 0.073108 0.073018 0.066898 0.016458)\n289.10705 11.57 138 159 181 (1.1828 0.17822 0.030163)\n289.14344 9.66 196 205 206 208 199 (0.60481 0.21704 0.16288 0.13588 0.040776)\n291.08631 8.52 3 160 184 (0.79832 0.2053 0.021057)\n305.10196 13.88 152 157 153 154 155 151 (0.53806 0.35267 0.33738 0.26837 0.10512 0.067448)\n307.11761 10.69 164 162 169 163 161 167 168 158 165 (0.30888 0.27666 0.18987 0.1687 0.14948 0.13251 0.031785 0.02005 0.0079015)\n317.13835 14.31 195 (1.7208)\n325.14344 8.40 238 (1.0107)\n329.13835 9.89 144 258 250 228 (0.49375 0.35819 0.22699 0.11026)\n357.13326 7.81 191 19 20 1 (0.66319 0.16606 0.10971 0.00092106)\n"
      }
      
      try {
        // 动态导入MS2处理器（避免在不需要时加载）
        const ms2Processor = await import('../utils/ms2Processor.js')
        
        // 在客户端计算指纹，限制最大峰数量为5000，保留强度最高的峰，避免处理大文件时性能问题
        const fingerprintJson = ms2Processor.calculateMS2FingerprintJson(ms2Text.value, 5000, true)
        
        if (!fingerprintJson) {
          throw new Error('无法计算MS2指纹：数据格式可能不正确')
        }
        
        console.log('客户端计算的MS2指纹完成，发送到后端进行搜索...')
        
        // 使用新的指纹搜索API
        const response = await fetch('/api/rdkit/ms2-search-by-fingerprint', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            fingerprint_json: fingerprintJson,
            threshold: threshold.value,
            tolerance: tolerance.value,
            prefilter_threshold: prefilterThreshold.value
          })
        })
        
        if (!response.ok) {
          throw new Error(`API请求失败: ${response.status}`)
        }

        const result = await response.json()
        console.log('MS2指纹搜索API响应结果:', result)
        
        // 处理API响应格式
        if (result.code === 200200 && result.data) {
          try {
            const parsedData = JSON.parse(result.data)
            console.log('解析出的MS2搜索结果:', parsedData)
            
            let compoundIds = []
            
            // MS2搜索：数据格式为 [["CMP0005", 0.8], ["CMP0054", 0.7]]
            if (Array.isArray(parsedData) && parsedData.length > 0) {
              compoundIds = parsedData.map(item => item[0]) // 提取ID
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
            console.error('解析MS2数据失败:', parseError)
            searchResults.value = []
          }
        } else {
          searchResults.value = []
        }
      } catch (error) {
        console.error('MS2搜索失败:', error)
        errorMessage.value = error.message || t('query.search_failed')
      }
    } else {
      // 原有的结构搜索
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
          
          if (!fpResponse.ok) {
            throw new Error(`指纹生成失败: ${fpResponse.status}`)
          }
          
          const fpData = await fpResponse.json()
          
          // 检查API响应格式
          if (!fpData) {
            console.error('指纹API响应数据:', fpData)
            throw new Error('指纹生成失败: 响应中未找到data字段')
          }
          
          // 确保fingerprint不是undefined或null
          if (!fpData) {
            throw new Error('指纹生成失败: 指纹为空')
          }
          
          response = await fetch(`/api/rdkit/similarity?qfp=${encodeURIComponent(fpData)}&threshold=${threshold.value}`)
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
    similarity: t('query.similarity'),
    'c-nmr': t('query.c-nmr'),
    'ms2': t('query.ms2')
  }
  return modes[searchMode.value] || searchMode.value
}

// 调整Prefilter threshold的最大值
const adjustPrefilterMax = () => {
  // 如果prefilterThreshold大于threshold，则将其调整为threshold
  if (prefilterThreshold.value > threshold.value) {
    prefilterThreshold.value = threshold.value
  }
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
                <button
                  type="button"
                  class="btn btn-outline-primary text-start"
                  :class="{ 'active': searchMode === 'c-nmr' }"
                  @click="setSearchMode('c-nmr')"
                >
                  <i class="bi bi-magnet"></i> {{ t('query.c-nmr') }}
                </button>
                <button
                  type="button"
                  class="btn btn-outline-primary text-start"
                  :class="{ 'active': searchMode === 'ms2' }"
                  @click="setSearchMode('ms2')"
                >
                  <i class="bi bi-graph-up"></i> {{ t('query.ms2') }}
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
                v-if="searchMode != 'c-nmr'"
              >
                <i class="bi bi-download"></i> {{ t('query.save_structure') }}
              </button>
            </div>

            <!-- 当前模式显示 -->
            <div class="mt-4 p-3 bg-light rounded">
              <small class="text-muted">{{ t('query.current_mode') }}:</small>
              <div class="fw-bold text-primary">{{ getSearchModeText() }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：Ketcher 编辑器和搜索结果 -->
      <div class="col-lg-9 col-md-8">
        <!-- Ketcher 编辑器 -->
        <div class="card shadow-sm border-0 mb-4" v-if="searchMode != 'c-nmr' && searchMode != 'ms2'">
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
        
        <!-- 相似度搜索参数 -->
        <div class="card shadow-sm border-0 mb-4" v-if="searchMode === 'similarity'">
          <div class="card-header bg-info text-white">
            <h5 class="card-title mb-0">
              <i class="bi bi-sliders"></i> {{ t('query.similarity_settings') }}
            </h5>
          </div>
          <div class="card-body">
            <div class="row">
              <div class="col-md-12 mb-3">
                <label for="similarityThreshold" class="form-label fw-semibold">{{ t('query.similarity_threshold') }}</label>
                <input 
                  type="range" 
                  id="similarityThreshold" 
                  class="form-range" 
                  v-model.number="threshold" 
                  min="0" 
                  max="1" 
                  step="0.05"
                >
                <div class="d-flex justify-content-between">
                  <small>0</small>
                  <small class="fw-bold">{{ Number(threshold).toFixed(2) }}</small>
                  <small>1</small>
                </div>
                <div class="form-text">{{ t('query.similarity_threshold_help') }}</div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="card shadow-sm border-0 mb-4" v-else-if="searchMode === 'c-nmr'">
          <div class="card-header bg-primary text-white d-flex justify-content-between align-items-center">
            <h5 class="card-title mb-0">
              <i class="bi bi-pencil-square"></i> {{ t('query.cnmr_editor') }}
            </h5>
          </div>
              <div class="card-body">
            <!-- 核磁谱文本输入 -->
            <div class="mb-3">
              <label for="nmrText" class="form-label fw-semibold">{{ t('query.nmr_data') }}</label>
              <textarea 
                id="nmrText"
                class="form-control" 
                v-model="nmrText" 
                rows="8" 
                :placeholder="t('query.enter_nmr_data_or_upload')"
              ></textarea>
              <div class="form-text">{{ t('query.supported_formats') }}</div>
            </div>
            
            <!-- 文件上传 -->
            <div class="mb-3">
              <label for="nmrFile" class="form-label fw-semibold">{{ t('query.upload_from_file') }}</label>
              <input 
                type="file" 
                id="nmrFile" 

                class="form-control" 
                @change="(event) => handleFileUpload(event, 'nmr')"
                accept=".txt,.csv,.json,.log"
              >
              <div class="form-text">{{ t('query.supported_file_formats') }}</div>
            </div>
            
            <!-- 搜索参数 -->
            <div class="row">
              <div class="col-md-6 mb-3">
                <label for="threshold" class="form-label fw-semibold">{{ t('query.similarity_threshold') }}</label>
                <input 
                  type="range" 
                  id="threshold" 
                  class="form-range" 
                  v-model.number="threshold" 
                  min="0" 
                  max="1" 
                  step="0.05"
                >
                <div class="d-flex justify-content-between">
                  <small>0</small>
                  <small class="fw-bold">{{ Number(threshold).toFixed(2) }}</small>
                  <small>1</small>
                </div>
              </div>
              <div class="col-md-6 mb-3">
                <label for="tolerance" class="form-label fw-semibold">{{ t('query.tolerance_ppm') }}</label>
                <input 
                  type="range" 
                  id="tolerance" 
                  class="form-range" 
                  v-model.number="tolerance" 
                  min="0.1" 
                  max="2" 
                  step="0.1"
                >
                <div class="d-flex justify-content-between">
                  <small>0.1</small>
                  <small class="fw-bold">{{ Number(tolerance).toFixed(1) }}</small>
                  <small>2.0</small>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- MS2搜索界面 -->
        <div class="card shadow-sm border-0 mb-4" v-else-if="searchMode === 'ms2'">
          <div class="card-header bg-primary text-white d-flex justify-content-between align-items-center">
            <h5 class="card-title mb-0">
              <i class="bi bi-pencil-square"></i> {{ t('query.ms2_editor') }}
            </h5>
          </div>
              <div class="card-body">
            <!-- MS2文本输入 -->
            <div class="mb-3">
              <label for="ms2Text" class="form-label fw-semibold">{{ t('query.ms2_data') }}</label>
              <textarea 
                id="ms2Text"
                class="form-control" 
                v-model="ms2Text" 
                rows="8" 
                :placeholder="t('query.enter_ms2_data_or_upload')"
              ></textarea>
              <div class="form-text">{{ t('query.ms2_supported_formats') }}</div>
            </div>
            
            <!-- 文件上传 -->
            <div class="mb-3">
              <label for="ms2File" class="form-label fw-semibold">{{ t('query.upload_from_file') }}</label>
              <input 
                type="file" 
                id="ms2File" 

                class="form-control" 
                @change="(event) => handleFileUpload(event, 'ms2')"
                accept=".txt,.csv,.json,.log"
              >
              <div class="form-text">{{ t('query.supported_file_formats') }}</div>
            </div>
            
            <!-- 搜索参数 -->
            <div class="row">
              <div class="col-md-4 mb-3">
                <label for="threshold" class="form-label fw-semibold">{{ t('query.similarity_threshold') }}</label>
                <input 
                  type="range" 
                  id="threshold" 
                  class="form-range" 
                  v-model.number="threshold" 
                  min="0" 
                  max="1" 
                  step="0.05"
                  @input="adjustPrefilterMax"
                >
                <div class="d-flex justify-content-between">
                  <small>0</small>
                  <small class="fw-bold">{{ Number(threshold).toFixed(2) }}</small>
                  <small>1</small>
                </div>
                <div class="form-text">{{ t('query.similarity_threshold_help_ms2') }}</div>
              </div>
              <div class="col-md-4 mb-3">
                <label for="tolerance" class="form-label fw-semibold">{{ t('query.tolerance_da') }}</label>
                <input 
                  type="range" 
                  id="tolerance" 
                  class="form-range" 
                  v-model.number="tolerance" 
                  min="0.1" 
                  max="2" 
                  step="0.1"
                >
                <div class="d-flex justify-content-between">
                  <small>0.1</small>
                  <small class="fw-bold">{{ Number(tolerance).toFixed(1) }}</small>
                  <small>2.0</small>
                </div>
                <div class="form-text">{{ t('query.tolerance_da_help') }}</div>
              </div>
              <div class="col-md-4 mb-3">
                <label for="prefilterThreshold" class="form-label fw-semibold">{{ t('query.prefilter_threshold') }}</label>
                <input 
                  type="range" 
                  id="prefilterThreshold" 
                  class="form-range" 
                  v-model.number="prefilterThreshold" 
                  :min="0" 
                  :max="threshold" 
                  step="0.05"
                >
                <div class="d-flex justify-content-between">
                  <small>0</small>
                  <small class="fw-bold">{{ Number(prefilterThreshold).toFixed(2) }}</small>
                  <small>{{ Number(threshold).toFixed(2) }}</small>
                </div>
                <div class="form-text">{{ t('query.prefilter_threshold_help', { max: threshold.toFixed(2) }) }}</div>
              </div>
            </div>
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
