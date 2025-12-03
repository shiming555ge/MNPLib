<template>
  <div class="super-admin-container">
    <div class="container-fluid py-4">
      <!-- 页面标题 -->
      <div class="row mb-4">
        <div class="col-12">
          <h1 class="display-6 fw-bold text-primary">
            <i class="bi bi-shield-lock me-2"></i>
            {{ t('superadmin.title') }}
          </h1>
          <p class="text-muted">{{ t('superadmin.description') }}</p>
        </div>
      </div>

      <!-- 操作按钮区域 -->
      <div class="row mb-4">
        <div class="col-12 d-flex justify-content-between align-items-center">
          <div>
            <button class="btn btn-primary" @click="showCreateModal = true">
              <i class="bi bi-plus-circle me-1"></i>
              {{ t('superadmin.create_passkey') }}
            </button>
          </div>
          <div>
            <button class="btn btn-outline-secondary" @click="refreshPasskeys">
              <i class="bi bi-arrow-clockwise me-1"></i>
              {{ t('superadmin.refresh') }}
            </button>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="row">
        <div class="col-12 text-center py-5">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">加载中...</span>
          </div>
          <p class="text-muted mt-2">{{ t('superadmin.loading') }}</p>
        </div>
      </div>

      <!-- 错误提示 -->
      <div v-else-if="error" class="row">
        <div class="col-12">
          <div class="alert alert-danger" role="alert">
            <i class="bi bi-exclamation-triangle me-2"></i>
            {{ error }}
            <button class="btn btn-sm btn-outline-danger ms-2" @click="refreshPasskeys">
              {{ t('superadmin.retry') }}
            </button>
          </div>
        </div>
      </div>

      <!-- 数据表格 -->
      <div v-else class="row">
        <div class="col-12">
          <div class="card shadow-sm">
            <div class="card-header bg-light">
              <h5 class="card-title mb-0">
                <i class="bi bi-key me-2"></i>
                {{ t('superadmin.passkey_list') }}
                <span class="badge bg-primary ms-2">{{ passkeys.length }}</span>
              </h5>
            </div>
            <div class="card-body p-0">
              <div class="table-responsive">
                <table class="table table-hover mb-0">
                  <thead class="table-light">
                    <tr>
                      <th>{{ t('superadmin.passkey') }}</th>
                      <th>{{ t('superadmin.pdescription') }}</th>
                      <th>{{ t('superadmin.operator') }}</th>
                      <th>{{ t('superadmin.created_at') }}</th>
                      <th>{{ t('superadmin.status') }}</th>
                      <th>{{ t('superadmin.creator') }}</th>
                      <th>{{ t('superadmin.actions') }}</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="passkey in passkeys" :key="passkey.passkey">
                      <td>
                        <code class="text-truncate d-inline-block" style="max-width: 150px;" 
                              :title="passkey.passkey">
                          {{ passkey.passkey }}
                        </code>
                        <button class="btn btn-sm btn-outline-secondary ms-1" 
                                @click="copyToClipboard(passkey.passkey)"
                                :title="t('superadmin.copy')">
                          <i class="bi bi-clipboard"></i>
                        </button>
                      </td>
                      <td>{{ passkey.description || '-' }}</td>
                      <td>{{ passkey.operator || '-' }}</td>
                      <td>{{ formatDate(passkey.created_at) }}</td>
                      <td>
                        <span :class="['badge', passkey.is_active ? 'bg-success' : 'bg-danger']">
                          {{ passkey.is_active ? t('superadmin.active') : t('superadmin.inactive') }}
                        </span>
                      </td>
                      <td>{{ passkey.extends || '-' }}</td>
                      <td>
                        <div class="btn-group btn-group-sm" role="group">
                          <button class="btn btn-outline-primary" 
                                  @click="editPasskey(passkey)"
                                  :title="t('superadmin.edit')">
                            <i class="bi bi-pencil"></i>
                          </button>
                          <button class="btn btn-outline-warning" 
                                  @click="togglePasskeyStatus(passkey)"
                                  :title="passkey.is_active ? t('superadmin.deactivate') : t('superadmin.activate')">
                            <i class="bi" :class="passkey.is_active ? 'bi-toggle-off' : 'bi-toggle-on'"></i>
                          </button>
                          <button class="btn btn-outline-danger" 
                                  @click="deletePasskey(passkey)"
                                  :title="t('superadmin.delete')">
                            <i class="bi bi-trash"></i>
                          </button>
                        </div>
                      </td>
                    </tr>
                    <tr v-if="passkeys.length === 0">
                      <td colspan="7" class="text-center py-4 text-muted">
                        <i class="bi bi-inbox display-6 d-block mb-2"></i>
                        {{ t('superadmin.no_passkeys') }}
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑模态框 -->
    <div class="modal fade" id="passkeyModal" tabindex="-1" aria-labelledby="passkeyModalLabel" aria-hidden="true" ref="modal">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="passkeyModalLabel">
              {{ isEditing ? t('superadmin.edit_passkey') : t('superadmin.create_passkey') }}
            </h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="handleSubmit">
              <div class="mb-3">
                <label for="operator" class="form-label">{{ t('superadmin.operator') }} *</label>
                <input type="text" class="form-control" id="operator" v-model="form.operator" required
                       :placeholder="t('superadmin.operator_placeholder')">
              </div>
              <div class="mb-3">
                <label for="description" class="form-label">{{ t('superadmin.description') }}</label>
                <textarea class="form-control" id="description" v-model="form.description" rows="3"
                          :placeholder="t('superadmin.description_placeholder')"></textarea>
              </div>
              <div class="mb-3">
                <div class="form-check">
                  <input class="form-check-input" type="checkbox" id="is_active" v-model="form.is_active">
                  <label class="form-check-label" for="is_active">
                    {{ t('superadmin.active') }}
                  </label>
                </div>
              </div>
              <div class="d-flex justify-content-end gap-2">
                <button type="button" class="btn btn-outline-secondary" data-bs-dismiss="modal">
                  {{ t('superadmin.cancel') }}
                </button>
                <button type="submit" class="btn btn-primary" :disabled="submitting">
                  <span v-if="submitting" class="spinner-border spinner-border-sm me-1" role="status"></span>
                  {{ isEditing ? t('superadmin.update') : t('superadmin.create') }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { Modal } from 'bootstrap';
import { useAuth } from '../composables/useAuth';
import { useRouter } from 'vue-router'

const router = useRouter()

const { t } = useI18n();
const { getAuthHeader, verifyPasskeyModifiable } = useAuth();

// 状态管理
const passkeys = ref([]);
const loading = ref(false);
const error = ref('');
const showCreateModal = ref(false);
const isEditing = ref(false);
const submitting = ref(false);
const modal = ref(null);
const modalInstance = ref(null);

// 表单数据
const form = ref({
  passkey: '',
  operator: '',
  description: '',
  is_active: true
});

// 当前编辑的passkey
const currentPasskey = ref(null);

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-';
  try {
    const date = new Date(dateString);
    return date.toLocaleString();
  } catch (e) {
    return dateString;
  }
};

// 复制到剪贴板
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text);
    alert(t('superadmin.copied'));
  } catch (err) {
    console.error('复制失败:', err);
    alert(t('superadmin.copy_failed'));
  }
};

// 获取所有passkey
const fetchPasskeys = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    const headers = {
      'Content-Type': 'application/json',
      ...getAuthHeader()
    };
    
    const response = await fetch('/api/passkeys', {
      headers
    });
    
    if (response.ok) {
      const data = await response.json();
      passkeys.value = data.data || [];
    } else {
      const errorData = await response.json().catch(() => ({}));
      error.value = errorData.message || t('superadmin.fetch_failed');
    }
  } catch (err) {
    console.error('获取passkey失败:', err);
    error.value = t('superadmin.network_error');
    
  } finally {
    loading.value = false;
  }
};

// 刷新数据
const refreshPasskeys = () => {
  fetchPasskeys();
};

// 创建passkey
const createPasskey = async () => {
  submitting.value = true;
  
  try {
    const headers = {
      'Content-Type': 'application/json',
      ...getAuthHeader()
    };
    
    const response = await fetch('/api/passkeys', {
      method: 'POST',
      headers,
      body: JSON.stringify({
        description: form.value.description,
        operator: form.value.operator,
        is_active: form.value.is_active
      })
    });

    if (response.ok) {
      const data = await response.json();
      passkeys.value.unshift(data);
      hideModal();
      resetForm();
      alert(t('superadmin.create_success'));
      // 刷新数据以确保获取最新状态
      fetchPasskeys();
    } else {
      const errorData = await response.json().catch(() => ({}));
      alert(errorData.message || t('superadmin.create_failed'));
    }
  } catch (err) {
    console.error('创建passkey失败:', err);
    alert(t('superadmin.network_error'));
  } finally {
    submitting.value = false;
  }
};

// 更新passkey
const updatePasskey = async () => {
  submitting.value = true;
  
  try {
    const headers = {
      'Content-Type': 'application/json',
      ...getAuthHeader()
    };
    
    const response = await fetch(`/api/passkeys/${currentPasskey.value.passkey}`, {
      method: 'PUT',
      headers,
      body: JSON.stringify({
        description: form.value.description,
        operator: form.value.operator,
        is_active: form.value.is_active
      })
    });
    
    if (response.ok) {
      const data = await response.json();
      const index = passkeys.value.findIndex(p => p.passkey === data.passkey);
      if (index !== -1) {
        passkeys.value[index] = data;
      }
      hideModal();
      resetForm();
      alert(t('superadmin.update_success'));
      // 刷新数据以确保获取最新状态
      fetchPasskeys();
    } else {
      const errorData = await response.json().catch(() => ({}));
      alert(errorData.message || t('superadmin.update_failed'));
    }
  } catch (err) {
    console.error('更新passkey失败:', err);
    alert(t('superadmin.network_error'));
  } finally {
    submitting.value = false;
  }
};

// 删除passkey
const deletePasskey = async (passkey) => {
  if (!confirm(t('superadmin.delete_confirm'))) {
    return;
  }
  
  try {
    const headers = {
      'Content-Type': 'application/json',
      ...getAuthHeader()
    };
    
    const response = await fetch(`/api/passkeys/${passkey.passkey}`, {
      method: 'DELETE',
      headers
    });
    
    if (response.ok) {
      passkeys.value = passkeys.value.filter(p => p.passkey !== passkey.passkey);
      alert(t('superadmin.delete_success'));
      // 刷新数据以确保获取最新状态
      fetchPasskeys();
    } else {
      const errorData = await response.json().catch(() => ({}));
      alert(errorData.message || t('superadmin.delete_failed'));
    }
  } catch (err) {
    console.error('删除passkey失败:', err);
    alert(t('superadmin.network_error'));
  }
};

// 切换passkey状态
const togglePasskeyStatus = async (passkey) => {
  try {
    const headers = {
      'Content-Type': 'application/json',
      ...getAuthHeader()
    };
    
    const response = await fetch(`/api/passkeys/${passkey.passkey}/toggle`, {
      method: 'POST',
      headers
    });
    
    if (response.ok) {
      const data = await response.json();
      const index = passkeys.value.findIndex(p => p.passkey === data.passkey);
      if (index !== -1) {
        passkeys.value[index] = data;
      }
      alert(t('superadmin.toggle_success'));
      // 刷新数据以确保获取最新状态
      fetchPasskeys();
    } else {
      const errorData = await response.json().catch(() => ({}));
      alert(errorData.message || t('superadmin.toggle_failed'));
    }
  } catch (err) {
    console.error('切换状态失败:', err);
    alert(t('superadmin.network_error'));
  }
};

// 编辑passkey
const editPasskey = (passkey) => {
  isEditing.value = true;
  currentPasskey.value = passkey;
  form.value = {
    passkey: passkey.passkey,
    operator: passkey.operator,
    description: passkey.description,
    is_active: passkey.is_active
  };
  showModal();
};

// 显示模态框
const showModal = () => {
  if (modalInstance.value) {
    modalInstance.value.show();
  }
};

// 隐藏模态框
const hideModal = () => {
  if (modalInstance.value) {
    modalInstance.value.hide();
  }
};

// 重置表单
const resetForm = () => {
  form.value = {
    passkey: '',
    operator: '',
    description: '',
    is_active: true
  };
  isEditing.value = false;
  currentPasskey.value = null;
};

// 表单提交
const handleSubmit = () => {
  if (isEditing.value) {
    updatePasskey();
  } else {
    createPasskey();
  }
};

// 监听showCreateModal变化
watch(() => showCreateModal.value, (newVal) => {
  if (newVal) {
    resetForm();
    showModal();
  }
});

// 组件挂载
onMounted(async () => {
  if(!(await verifyPasskeyModifiable())){
    router.push({name:"home"})
  }
  // 初始化模态框
  if (modal.value) {
    modalInstance.value = new Modal(modal.value);
    
    // 监听模态框隐藏事件
    modal.value.addEventListener('hidden.bs.modal', () => {
      showCreateModal.value = false;
      resetForm();
    });
  }
  
  // 获取passkey列表
  fetchPasskeys();
});
</script>

<style scoped>
.super-admin-container {
  min-height: calc(100vh - 200px);
}

.table th {
  font-weight: 600;
  color: #495057;
  background-color: #f8f9fa;
}

.table td {
  vertical-align: middle;
}

.btn-group-sm > .btn {
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
}

.badge {
  font-size: 0.75em;
  font-weight: 500;
}

.modal-content {
  border-radius: 12px;
  border: none;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
}

.modal-header {
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  border-radius: 12px 12px 0 0;
}

.form-control:focus {
  border-color: #007bff;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

.btn-primary {
  background-color: #007bff;
  border-color: #007bff;
  font-weight: 500;
}

.btn-primary:hover {
  background-color: #0056b3;
  border-color: #0056b3;
}

.btn-primary:disabled {
  background-color: #6c757d;
  border-color: #6c757d;
}
</style>
