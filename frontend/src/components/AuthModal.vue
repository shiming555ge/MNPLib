<template>
  <div class="modal fade" id="authModal" tabindex="-1" aria-labelledby="authModalLabel" aria-hidden="true" ref="modal">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="authModalLabel">{{ t('auth.auth_required') }}</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <p class="text-muted mb-4">{{ t('auth.auth_description') }}</p>
          
          <form @submit.prevent="handleSubmit">
            <div class="mb-3">
              <label for="passkeyInput" class="form-label">{{ t('auth.passkey') }}</label>
              <input 
                type="password" 
                class="form-control" 
                id="passkeyInput" 
                v-model="passkey"
                :placeholder="t('auth.enter_passkey')"
                required
                :disabled="loading"
              />
              <div class="form-text text-danger" v-if="error">
                {{ error }}
              </div>
            </div>
            
            <div class="d-flex justify-content-end gap-2">
              <button type="button" class="btn btn-outline-secondary" data-bs-dismiss="modal" :disabled="loading">
                {{ t('auth.cancel') }}
              </button>
              <button type="submit" class="btn btn-primary" :disabled="loading">
                <span v-if="loading" class="spinner-border spinner-border-sm me-1" role="status"></span>
                {{ loading ? t('auth.logging_in') : t('auth.submit') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { Modal } from 'bootstrap';

const { t } = useI18n();

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:show', 'login-success']);

const modal = ref(null);
const modalInstance = ref(null);
const passkey = ref('');
const loading = ref(false);
const error = ref('');

const handleSubmit = async () => {
  if (!passkey.value.trim()) {
    error.value = t('auth.invalid_passkey');
    return;
  }

  loading.value = true;
  error.value = '';

  try {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ passkey: passkey.value })
    });

    if (response.ok) {
      const data = await response.json();
      console.log('Login successful:', data);
      
      // 根据返回体结构解析数据
      // 返回体: { operator: "string", token: "string", description: "string" }
      const info = data.data || null;
      
      if (info) {
        // 存储data到localStorage（Navbar组件中的login函数会处理这个）
        localStorage.setItem('userInfo', info);
        
        // 隐藏模态框
        if (modalInstance.value) {
          modalInstance.value.hide();
        }
        
        // 清空表单
        passkey.value = '';
        
        // 触发成功事件，传递完整的响应数据
        emit('login-success', info);
        
        // 显示成功消息
        alert(t('auth.login_success'));        
      } else {
        // 如果响应中没有token，显示错误
        error.value = t('auth.login_failed') + ': No token received';
      }
    } else {
      const errorData = await response.json().catch(() => ({}));
      error.value = errorData.message || t('auth.login_failed');
    }
  } catch (err) {
    console.error('Login error:', err);
    error.value = t('auth.login_failed');
  } finally {
    loading.value = false;
  }
};

const showModal = () => {
  if (modalInstance.value) {
    modalInstance.value.show();
  }
};

const hideModal = () => {
  if (modalInstance.value) {
    modalInstance.value.hide();
  }
};

// 监听show prop变化
onMounted(() => {
  if (modal.value) {
    modalInstance.value = new Modal(modal.value);
    
    // 监听模态框隐藏事件
    modal.value.addEventListener('hidden.bs.modal', () => {
      emit('update:show', false);
      passkey.value = '';
      error.value = '';
      loading.value = false;
    });
    
    // 监听模态框显示事件
    modal.value.addEventListener('shown.bs.modal', () => {
      emit('update:show', true);
    });
  }
});

// 当show prop变化时控制模态框显示/隐藏
watch(() => props.show, (newVal) => {
  if (modalInstance.value) {
    if (newVal) {
      modalInstance.value.show();
    } else {
      modalInstance.value.hide();
    }
  }
});

onUnmounted(() => {
  if (modalInstance.value) {
    modalInstance.value.dispose();
  }
});

// 暴露方法给父组件
defineExpose({
  showModal,
  hideModal
});
</script>

<style scoped>
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

.modal-title {
  color: #2c3e50;
  font-weight: 600;
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
