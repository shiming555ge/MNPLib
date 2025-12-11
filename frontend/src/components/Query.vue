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
const defaultNMR = '173.7, 170.6, 170.3, 169.5, 169.5, 157.8, 150.5, 134.5, 129.2, 128.5, 128.5, 126.7, 123.2, 121.0, 119.9, 113.6, 113.6, 57.5, 56.4, 55.8, 54.9, 42.9, 40.0, 38.0, 34.8, 28.4, 27.4, 22.5, 22.5, 19.1, 16.8, 14.1'

// MS2搜索相关
const ms2Text = ref('')
const ms2File = ref(null)
const prefilterThreshold = ref(0.25) // 默认值为0.5的一半
const defaultMS2 = 'energy0\n263.08150 3.50\n265.09715 100.00\nenergy1\n169.07602 14.84\n247.08659 11.96\n265.09715 100.00'

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
        nmrText.value = defaultNMR
      }

      nmrText.value = nmrText.value.replace("，", ",")
      
      console.log(`/api/rdkit/nmr-search?query_nmr=${encodeURIComponent(nmrText.value)}&threshold=${threshold.value}&tolerance=${tolerance.value}`)
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
      // MS2相似度搜索 - 按能量级别分别搜索
      if (!ms2Text.value.trim()) {
        ms2Text.value = defaultMS2
      }
      
      try {
        // 动态导入MS2处理器（避免在不需要时加载）
        const ms2Processor = await import('../utils/ms2Processor.js')
        
        // 解析MS2数据，获取能量级别信息
        const ms2Data = ms2Processor.parseMS2Data(ms2Text.value, 5000, true)
        console.log('解析的MS2数据能量级别:', Object.keys(ms2Data))
        
        // 确定要搜索的能量级别
        let energyLevelsToSearch = []
        const hasEnergyMarkers = Object.keys(ms2Data).some(key => key.startsWith('energy'))
        
        if (hasEnergyMarkers) {
          // 有energy标识的块，使用所有检测到的能量级别
          energyLevelsToSearch = Object.keys(ms2Data).filter(key => key.startsWith('energy'))
          console.log('检测到能量级别块:', energyLevelsToSearch)
        } else {
          // 没有energy标识的块，从energy0到energy2各查找一次，再加上总体搜索
          energyLevelsToSearch = ['energy0', 'energy1', 'energy2', 'all']
          console.log('未检测到能量级别，使用默认能量级别和总体搜索:', energyLevelsToSearch)
        }
        
        if (energyLevelsToSearch.length === 0) {
          throw new Error('无法确定能量级别：请检查MS2数据格式')
        }
        
        // 为每个能量级别准备搜索请求
        const searchPromises = energyLevelsToSearch.map(async (energyLevel) => {
          console.log(`发送 ${energyLevel} 能量级别搜索请求...`)
          
          if (energyLevel === 'all') {
            // 总体搜索 - 使用原始的MS2搜索API
            const response = await fetch('/api/rdkit/ms2-search', {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({
                query_ms2: ms2Text.value,
                threshold: threshold.value,
                tolerance: tolerance.value,
                prefilter_threshold: prefilterThreshold.value
              })
            })
            
            if (!response.ok) {
              throw new Error(`总体API请求失败: ${response.status}`)
            }
            
            const result = await response.json()
            console.log('总体搜索API响应结果:', result)
            
            // 处理API响应格式
            if (result.code === 200200 && result.data) {
              try {
                const parsedData = JSON.parse(result.data)
                console.log('总体搜索解析出的MS2搜索结果:', parsedData)
                
                // 提取化合物ID和相似度分数
                let compoundResults = []
                if (Array.isArray(parsedData) && parsedData.length > 0) {
                  compoundResults = parsedData.map(item => ({
                    id: item[0],
                    similarity: item[1]
                  }))
                }
                
                return {
                  energyLevel: 'all',
                  success: true,
                  compoundResults
                }
              } catch (parseError) {
                console.error('解析总体搜索数据失败:', parseError)
                return {
                  energyLevel: 'all',
                  success: false,
                  compoundResults: []
                }
              }
            } else {
              return {
                energyLevel: 'all',
                success: false,
                compoundResults: []
              }
            }
          } else {
            // 能量级别搜索
            const response = await fetch('/api/rdkit/ms2-search-by-energy-level', {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({
                query_ms2: ms2Text.value,
                threshold: threshold.value,
                tolerance: tolerance.value,
                prefilter_threshold: prefilterThreshold.value,
                energy_level: energyLevel
              })
            })
            
            if (!response.ok) {
              throw new Error(`${energyLevel} API请求失败: ${response.status}`)
            }
            
            const result = await response.json()
            console.log(`${energyLevel} 搜索API响应结果:`, result)
            
            // 处理API响应格式
            if (result.code === 200200 && result.data) {
              try {
                const parsedData = JSON.parse(result.data)
                console.log(`${energyLevel} 解析出的MS2搜索结果:`, parsedData)
                
                // 提取化合物ID和相似度分数
                let compoundResults = []
                if (Array.isArray(parsedData) && parsedData.length > 0) {
                  compoundResults = parsedData.map(item => ({
                    id: item[0],
                    similarity: item[1]
                  }))
                }
                
                return {
                  energyLevel,
                  success: true,
                  compoundResults
                }
              } catch (parseError) {
                console.error(`解析 ${energyLevel} 数据失败:`, parseError)
                return {
                  energyLevel,
                  success: false,
                  compoundResults: []
                }
              }
            } else {
              return {
                energyLevel,
                success: false,
                compoundResults: []
              }
            }
          }
        })
        
        // 等待所有能量级别的搜索结果
        const allResults = await Promise.allSettled(searchPromises)
        console.log('所有能量级别搜索结果:', allResults)
        
        // 收集所有成功的搜索结果
        const successfulResults = []
        for (const result of allResults) {
          if (result.status === 'fulfilled' && result.value.success) {
            successfulResults.push(result.value)
          }
        }
        
        if (successfulResults.length === 0) {
          throw new Error('所有能量级别搜索都失败了')
        }
        
        // 根据是否有能量标记决定取交集还是并集
        let finalCompoundIds = []
        
        if (hasEnergyMarkers) {
          // 有能量标记：多个能量级别结果取交集
          if (successfulResults.length === 1) {
            // 只有一个能量级别，直接使用其结果
            finalCompoundIds = successfulResults[0].compoundResults.map(item => item.id)
          } else {
            // 多个能量级别，取交集
            // 首先将每个能量级别的结果转换为ID集合
            const idSets = successfulResults.map(result => 
              new Set(result.compoundResults.map(item => item.id))
            )
            
            // 取所有集合的交集
            if (idSets.length > 0) {
              const firstSet = idSets[0]
              const intersection = Array.from(firstSet).filter(id => 
                idSets.every(set => set.has(id))
              )
              finalCompoundIds = intersection
            }
          }
          console.log('有能量标记，最终化合物ID（交集）:', finalCompoundIds)
        } else {
          // 没有能量标记：所有搜索结果取并集
          const allIds = new Set()
          for (const result of successfulResults) {
            for (const compound of result.compoundResults) {
              allIds.add(compound.id)
            }
          }
          finalCompoundIds = Array.from(allIds)
          console.log('无能量标记，最终化合物ID（并集）:', finalCompoundIds)
        }
        
        // 根据ID获取完整的化合物数据
        if (finalCompoundIds.length > 0) {
          const compoundPromises = finalCompoundIds.map(id => fetchCompoundById(id))
          const compounds = await Promise.all(compoundPromises)
          searchResults.value = compounds
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
                :placeholder="t('query.enter_ms2_data_or_upload')+defaultNMR"
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
                :placeholder="t('query.enter_ms2_data_or_upload')+defaultMS2"
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
