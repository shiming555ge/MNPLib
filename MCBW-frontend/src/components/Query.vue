<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// 搜索模式状态
const searchMode = ref('structure') // structure, substructure, similarity

// 处理搜索模式切换
const setSearchMode = (mode) => {
  searchMode.value = mode
}

// 处理搜索操作
const handleSearch = () => {
  // 这里可以添加与 Ketcher iframe 交互的逻辑
  console.log(`执行${getSearchModeText()}搜索`)
}

// 处理获取结构操作
const handleGetStructure = () => {
  // 这里可以添加从 Ketcher iframe 获取结构的逻辑
  console.log('获取当前结构')
}

// 获取搜索模式文本
const getSearchModeText = () => {
  const modes = {
    structure: t('query.exact_structure'),
    substructure: t('query.substructure'),
    similarity: t('query.similarity')
  }
  return modes[searchMode.value]
}
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
                @click="handleGetStructure"
              >
                <i class="bi bi-download"></i> {{ t('query.get_structure') }}
              </button>
            </div>

            <!-- 当前模式显示 -->
            <div class="mt-4 p-3 bg-light rounded">
              <small class="text-muted">{{ t('query.current_mode') }}:</small>
              <div class="fw-bold text-primary">{{ getSearchModeText() }}{{ t('query.search') }}</div>
            </div>
          </div>
        </div>

        <!-- 快捷操作卡片 -->
        <div class="card shadow-sm border-0 mt-3">
          <div class="card-header bg-secondary text-white">
            <h6 class="card-title mb-0">
              <i class="bi bi-lightning"></i> {{ t('query.quick_actions') }}
            </h6>
          </div>
          <div class="card-body">
            <div class="d-grid gap-2">
              <button class="btn btn-outline-secondary btn-sm">
                <i class="bi bi-arrow-clockwise"></i> {{ t('query.clear_canvas') }}
              </button>
              <button class="btn btn-outline-secondary btn-sm">
                <i class="bi bi-floppy"></i> {{ t('query.save_structure') }}
              </button>
              <button class="btn btn-outline-secondary btn-sm">
                <i class="bi bi-upload"></i> {{ t('query.import_structure') }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：Ketcher 编辑器 -->
      <div class="col-lg-9 col-md-8">
        <div class="card shadow-sm border-0 h-100">
          <div class="card-header bg-primary text-white d-flex justify-content-between align-items-center">
            <h5 class="card-title mb-0">
              <i class="bi bi-pencil-square"></i> {{ t('query.chemical_editor') }}
            </h5>
            <span class="badge bg-light text-primary">{{ t('query.ketcher') }}</span>
          </div>
          <div class="card-body p-0">
            <iframe 
              src="/Ketcher/index.html" 
              class="w-100 border-0" 
              style="height: 600px; min-height: 500px;"
              :title="t('query.chemical_editor')"
            ></iframe>
          </div>
        </div>
      </div>
    </div>
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
