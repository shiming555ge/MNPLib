<template>
  <div class="container-fluid py-5 bg-light min-vh-100">
    <div class="row justify-content-center">
      <div class="col-12 col-lg-10 col-xl-8">
        <!-- 页面标题 -->
        <div class="text-center mb-5">
          <h1 class="display-5 text-primary fw-bold mb-3">{{ t('about.title') }}</h1>
          <p class="lead text-muted">{{ t('about.subtitle') }}</p>
        </div>

        <!-- 关于内容 -->
        <div class="card shadow-sm border-0 mb-5">
          <div class="card-body p-4 p-lg-5">
            <div class="row align-items-center">
              <div class="col-12 col-lg-6 mb-4 mb-lg-0">
                <img 
                  src="/logo.png" 
                  alt="MNPLib Logo" 
                  class="img-fluid rounded shadow"
                  style="max-height: 300px; object-fit: contain;"
                >
              </div>
              <div class="col-12 col-lg-6">
                <h3 class="h4 text-primary mb-3">{{ t('about.what_is_mnplib') }}</h3>
                <p class="text-muted mb-4">{{ t('about.description') }}</p>
                
                <div class="d-flex align-items-center mb-3">
                  <i class="bi bi-check-circle-fill text-success me-3 fs-5"></i>
                  <span class="text-muted">{{ t('about.feature1') }}</span>
                </div>
                <div class="d-flex align-items-center mb-3">
                  <i class="bi bi-check-circle-fill text-success me-3 fs-5"></i>
                  <span class="text-muted">{{ t('about.feature2') }}</span>
                </div>
                <div class="d-flex align-items-center mb-3">
                  <i class="bi bi-check-circle-fill text-success me-3 fs-5"></i>
                  <span class="text-muted">{{ t('about.feature3') }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 使命与愿景 -->
        <div class="row g-4 mb-5">
          <div class="col-12 col-md-6">
            <div class="card h-100 border-0 shadow-sm">
              <div class="card-body p-4">
                <div class="d-flex align-items-center mb-3">
                  <div class="bg-primary bg-opacity-10 p-3 rounded-circle me-3">
                    <i class="bi bi-bullseye text-primary fs-4"></i>
                  </div>
                  <h3 class="h5 mb-0">{{ t('about.mission_title') }}</h3>
                </div>
                <p class="text-muted">{{ t('about.mission_description') }}</p>
              </div>
            </div>
          </div>
          
          <div class="col-12 col-md-6">
            <div class="card h-100 border-0 shadow-sm">
              <div class="card-body p-4">
                <div class="d-flex align-items-center mb-3">
                  <div class="bg-info bg-opacity-10 p-3 rounded-circle me-3">
                    <i class="bi bi-eye text-info fs-4"></i>
                  </div>
                  <h3 class="h5 mb-0">{{ t('about.vision_title') }}</h3>
                </div>
                <p class="text-muted">{{ t('about.vision_description') }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 技术栈 -->
        <div class="card border-0 shadow-sm mb-5">
          <div class="card-header bg-white border-0 py-3">
            <h3 class="h5 mb-0 text-primary">
              <i class="bi bi-gear me-2"></i>{{ t('about.tech_stack_title') }}
            </h3>
          </div>
          <div class="card-body p-4">
            <div class="row g-3">
              <div class="col-6 col-md-4 col-lg-3" v-for="tech in techStack" :key="tech.name">
                <div class="d-flex flex-column align-items-center p-3 border rounded bg-white">
                  <div class="mb-2" :style="{ color: tech.color, fontSize: '2rem' }">
                    <i :class="tech.icon"></i>
                  </div>
                  <span class="fw-medium text-center">{{ tech.name }}</span>
                  <small class="text-muted text-center">{{ tech.description }}</small>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 团队信息 -->
        <div class="card border-0 shadow-sm">
          <div class="card-header bg-white border-0 py-3">
            <h3 class="h5 mb-0 text-primary">
              <i class="bi bi-people me-2"></i>{{ t('about.team_title') }}
            </h3>
          </div>
          <div class="card-body p-4">
            <p class="text-muted mb-4">{{ t('about.team_description') }}</p>
            
            <div class="row g-4">
              <div class="col-12 col-md-6 col-lg-4" v-for="member in teamMembers" :key="member.name">
                <div class="card border h-100">
                  <div class="card-body text-center p-4">
                    <div class="mb-3">
                      <div class="bg-primary bg-opacity-10 rounded-circle d-inline-flex align-items-center justify-content-center" style="width: 80px; height: 80px;">
                        <i class="bi bi-person fs-3 text-primary"></i>
                      </div>
                    </div>
                    <h5 class="h6 fw-bold mb-1">{{ member.name }}</h5>
                    <p class="text-muted small mb-2">{{ member.role }}</p>
                    <p class="text-muted small">{{ member.description }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 联系信息 -->
        <div class="text-center mt-5 pt-4 border-top">
          <h4 class="h5 text-primary mb-3">{{ t('about.contact_title') }}</h4>
          <p class="text-muted mb-4">{{ t('about.contact_description') }}</p>
          <div class="d-flex justify-content-center gap-3">
            <a href="mailto:contact@mnplib.org" class="btn btn-outline-primary">
              <i class="bi bi-envelope me-2"></i>{{ t('about.contact_email') }}
            </a>
            <a href="https://github.com/mnplib" target="_blank" class="btn btn-outline-dark">
              <i class="bi bi-github me-2"></i>{{ t('about.contact_github') }}
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

// 技术栈数据
const techStack = ref([
  { name: 'Vue.js', icon: 'bi bi-vue', color: '#42b883', description: '前端框架' },
  { name: 'Bootstrap', icon: 'bi bi-bootstrap', color: '#7952b3', description: 'CSS框架' },
  { name: 'Vite', icon: 'bi bi-lightning', color: '#646cff', description: '构建工具' },
  { name: 'RDKit', icon: 'bi bi-flask', color: '#007bff', description: '化学信息学' },
  { name: 'Ketcher', icon: 'bi bi-pencil', color: '#28a745', description: '化学编辑器' },
  { name: 'Node.js', icon: 'bi bi-node-plus', color: '#339933', description: '后端运行环境' }
]);

// 团队成员数据
const teamMembers = ref([
  { 
    name: '张三', 
    role: '项目负责人', 
    description: '负责项目整体架构和化学信息学算法开发' 
  },
  { 
    name: '李四', 
    role: '前端开发', 
    description: '负责用户界面设计和前端功能实现' 
  },
  { 
    name: '王五', 
    role: '后端开发', 
    description: '负责数据库设计和API接口开发' 
  },
  { 
    name: '赵六', 
    role: '数据科学家', 
    description: '负责数据分析和机器学习模型开发' 
  },
  { 
    name: '钱七', 
    role: '化学专家', 
    description: '负责化学数据验证和质量控制' 
  },
  { 
    name: '孙八', 
    role: 'UI/UX设计师', 
    description: '负责用户体验和界面设计' 
  }
]);
</script>

<style scoped>
.card {
  transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
}

.card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1) !important;
}

.bg-primary.bg-opacity-10 {
  background-color: rgba(13, 110, 253, 0.1) !important;
}

.bg-info.bg-opacity-10 {
  background-color: rgba(13, 202, 240, 0.1) !important;
}
</style>
