<template>
  <div class="h-screen w-full flex items-center justify-center p-2 md:p-4">
    <AnimatedBorderCard
      :borderColor="'#f0ada0'"
      :borderWidth="3"
      class="w-full max-w-lg px-8 py-10 flex flex-col overflow-y-auto"
    >
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-2xl md:text-3xl font-bold mb-4 font-[ZCOOL_KuaiLe]">
          欢迎使用情侣纪念站点
        </h1>
        <p class="text-md font-[Ma_Shan_Zheng]">初始化你们的专属纪念空间</p>
      </div>

      <!-- Step Indicator -->
      <div class="flex items-center justify-between mb-6 text-sm text-on-surface">
        <span
          v-for="i in totalSteps"
          :key="i"
          :class="[
            'flex-1 text-center py-1',
            i === currentStep ? 'text-primary font-bold bg-primary-light/30 rounded-lg' : '',
          ]"
        >
          {{ i }}/{{ totalSteps }}
        </span>
      </div>
      <!-- Fixed height form container -->
      <div class="flex-grow h-100 pt-4 pb-4 flex flex-col">
        <!-- Step 1: Site -->
        <div v-if="currentStep === 1" class="flex flex-col gap-4 flex-grow">
          <div>
            <label class="win11-label">站点名称</label>
            <input
              v-model="form.siteName"
              placeholder="例如：鹿与星的纪念站"
              class="win11-input w-full"
            />
          </div>

          <div>
            <label class="win11-label">站点描述（可选）</label>
            <input
              v-model="form.siteDescription"
              placeholder="例如：记录我们的美好时光"
              class="win11-input w-full"
            />
          </div>

          <div>
            <label class="win11-label">故事开始的日期</label>
            <input v-model="form.startDate" type="date" class="win11-input w-full" />
          </div>

          <!-- Error Message Display -->
          <div class="min-h-[20px] w-full">
            <div v-if="getCurrentError" class="win11-error text-left">
              {{ getCurrentError }}
            </div>
          </div>
        </div>

        <!-- Step 2: User A -->
        <div v-if="currentStep === 2" class="flex flex-col gap-4 flex-grow">
          <div class="flex justify-center">
            <div
              class="w-16 h-16 rounded-full bg-surface flex items-center justify-center text-lg text-on-surface-variant cursor-pointer overflow-hidden ring-1 ring-border-color hover:ring-primary transition"
            >
              <img
                v-if="avatarAPreview"
                :src="avatarAPreview"
                class="w-full h-full object-cover"
                draggable="false"
              />
              <span v-else>{{ form.userAName?.[0] || 'A' }}</span>
            </div>
          </div>

          <div>
            <label class="win11-label">昵称</label>
            <input v-model="form.userAName" placeholder="昵称" class="win11-input w-full" />
          </div>

          <div>
            <label class="win11-label">角色</label>
            <select v-model="form.userARole" class="win11-input w-full">
              <option value="" disabled>请选择角色</option>
              <option value="boy">男生</option>
              <option value="girl">女生</option>
            </select>
          </div>

          <div>
            <label class="win11-label">邮箱（可选）</label>
            <input
              v-model="form.userAEmail"
              type="email"
              placeholder="邮箱（可选）"
              class="win11-input w-full"
            />
          </div>

          <div>
            <label class="win11-label">手机号（可选）</label>
            <input
              v-model="form.userAPhone"
              type="tel"
              placeholder="手机号（可选）"
              class="win11-input w-full"
            />
          </div>

          <!-- Error Message Display -->
          <div class="min-h-[20px] w-full">
            <div v-if="getCurrentError" class="win11-error text-left">
              {{ getCurrentError }}
            </div>
          </div>
        </div>

        <!-- Step 3: User B -->
        <div v-if="currentStep === 3" class="flex flex-col gap-4 flex-grow">
          <div class="flex justify-center">
            <div
              class="w-16 h-16 rounded-full bg-surface flex items-center justify-center text-lg text-on-surface-variant cursor-pointer overflow-hidden ring-1 ring-border-color hover:ring-primary transition"
            >
              <img
                v-if="avatarBPreview"
                :src="avatarBPreview"
                class="w-full h-full object-cover"
                draggable="false"
              />
              <span v-else>{{ form.userBName?.[0] || 'B' }}</span>
            </div>
          </div>

          <div>
            <label class="win11-label">昵称</label>
            <input v-model="form.userBName" placeholder="昵称" class="win11-input w-full" />
          </div>

          <div>
            <label class="win11-label">角色</label>
            <select v-model="form.userBRole" class="win11-input w-full">
              <option value="" disabled>请选择角色</option>
              <option value="boy">男生</option>
              <option value="girl">女生</option>
            </select>
          </div>

          <div>
            <label class="win11-label">邮箱（可选）</label>
            <input
              v-model="form.userBEmail"
              type="email"
              placeholder="邮箱（可选）"
              class="win11-input w-full"
            />
          </div>

          <div>
            <label class="win11-label">手机号（可选）</label>
            <input
              v-model="form.userBPhone"
              type="tel"
              placeholder="手机号（可选）"
              class="win11-input w-full"
            />
          </div>

          <!-- Error Message Display -->
          <div class="min-h-[20px] w-full">
            <div v-if="getCurrentError" class="win11-error text-left">
              {{ getCurrentError }}
            </div>
          </div>
        </div>

        <!-- Step 4: Password -->
        <div v-if="currentStep === 4" class="flex flex-col gap-4 flex-grow">
          <div>
            <label class="win11-label">站点访问密码</label>
            <input
              v-model="form.sitePassword"
              type="password"
              placeholder="站点访问密码"
              class="win11-input w-full"
            />
          </div>

          <div>
            <label class="win11-label">确认密码</label>
            <input
              v-model="form.sitePasswordConfirm"
              type="password"
              placeholder="确认密码"
              class="win11-input w-full"
            />
          </div>

          <!-- Error Message Display -->
          <div class="min-h-[20px] w-full">
            <div v-if="getCurrentError" class="win11-error text-left">
              {{ getCurrentError }}
            </div>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex justify-between mt-8">
        <button
          @click="prevStep"
          :class="[
            'px-4 py-2 rounded-full text-sm transition',
            currentStep === 1 ? 'invisible' : 'win11-button outline',
          ]"
        >
          上一步
        </button>

        <button type="button" @click.prevent="nextStep" class="win11-button">
          {{ currentStep === totalSteps ? '完成' : '下一步' }}
        </button>
      </div>
    </AnimatedBorderCard>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'

import AnimatedBorderCard from '@/components/ui/AnimatedBorderCard.vue'
import router from '@/router'
import { systemApi } from '@/services/system'
import { useSystemStore } from '@/stores/system'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/utils/toastUtils'

const currentStep = ref(1)
const totalSteps = 4
const uiStore = useUIStore()

const avatarAPreview = ref<string | null>(null)
const avatarBPreview = ref<string | null>(null)

const form = reactive({
  siteName: '',
  siteDescription: '',
  startDate: '',
  userAName: '',
  userARole: 'boy' as 'boy' | 'girl',
  userAEmail: '',
  userAPhone: '',
  userBName: '',
  userBRole: 'girl' as 'boy' | 'girl',
  userBEmail: '',
  userBPhone: '',
  sitePassword: '',
  sitePasswordConfirm: '',
})

const showToast = useToast()

const errors = reactive({
  siteName: '',
  startDate: '',
  userAName: '',
  userARole: '',
  userBName: '',
  userBRole: '',
  sitePassword: '',
  sitePasswordConfirm: '',
  passwordMismatch: '',
})

// 监听用户A角色变化，自动切换用户B角色
watch(
  () => form.userARole,
  (newRole: 'boy' | 'girl') => {
    form.userBRole = newRole === 'boy' ? 'girl' : 'boy'
  }
)

// 计算当前步骤的错误信息
const getCurrentError = computed(() => {
  switch (currentStep.value) {
    case 1:
      return errors.siteName || errors.startDate
    case 2:
      return errors.userAName || errors.userARole
    case 3:
      return errors.userBName || errors.userBRole
    case 4:
      return errors.sitePassword || errors.sitePasswordConfirm || errors.passwordMismatch
    default:
      return ''
  }
})

function clearErrors() {
  const errorKeys = Object.keys(errors) as Array<keyof typeof errors>
  errorKeys.forEach(key => {
    errors[key] = ''
  })
}

function validateStep(): boolean {
  clearErrors()

  if (currentStep.value === 1) {
    if (!form.siteName) {
      errors.siteName = '请输入站点名称'
      return false
    }
    if (!form.startDate) {
      errors.startDate = '请选择在一起日期'
      return false
    }
  }

  if (currentStep.value === 2) {
    if (!form.userAName) {
      errors.userAName = '请输入用户 A 昵称'
      return false
    }
    if (!form.userARole) {
      errors.userARole = '请选择用户 A 角色'
      return false
    }
  }

  if (currentStep.value === 3) {
    if (!form.userBName) {
      errors.userBName = '请输入用户 B 昵称'
      return false
    }
    if (!form.userBRole) {
      errors.userBRole = '请选择用户 B 角色'
      return false
    }
  }

  if (currentStep.value === 4) {
    if (!form.sitePassword) {
      errors.sitePassword = '请输入站点访问密码'
      return false
    }
    if (!form.sitePasswordConfirm) {
      errors.sitePasswordConfirm = '请确认站点访问密码'
      return false
    }
    if (form.sitePassword !== form.sitePasswordConfirm) {
      errors.passwordMismatch = '两次输入的密码不一致'
      return false
    }
  }

  return true
}

function nextStep() {
  if (!validateStep()) return
  if (currentStep.value < totalSteps) currentStep.value++
  else submit()
}

function prevStep() {
  if (currentStep.value > 1) currentStep.value--
}

function submit() {
  uiStore.setLoading(true)

  systemApi
    .initSystem(form)
    .then(res => {
      if (res.data.code === 0) {
        // 初始化成功后，立即更新系统状态
        const systemStore = useSystemStore()
        systemStore.setInitialized(true)

        showToast('系统初始化成功！', 'success')
        setTimeout(() => router.push('/'), 300)
      } else {
        showToast(res.data.msg || '初始化失败', 'error')
      }
    })
    .catch(err => {
      console.error(err)
      showToast('网络错误，请稍后重试', 'error')
    })
    .finally(() => {
      uiStore.setLoading(false)
    })
}
</script>

<style scoped>
/* Win11风格标签 */
.win11-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--on-surface-variant);
  margin-bottom: 4px;
  display: block;
}

/* Win11风格错误信息 */
.win11-error {
  color: #ff4757;
  font-size: 12px;
  margin-top: 4px;
  min-height: 18px;
}
</style>
