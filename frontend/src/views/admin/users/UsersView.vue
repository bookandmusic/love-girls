<template>
  <div class="admin-users">
    <h2 class="text-2xl font-bold text-gray-800 mb-6">情侣信息</h2>

    <!-- PC端表格视图 -->
    <div class="hidden md:block">
      <UserTable :users="users" @edit="openEditDialog" />
    </div>

    <!-- 移动端列表视图 -->
    <div class="md:hidden">
      <UserMobileList :users="users" @edit="openEditDialog" />
    </div>

    <!-- 编辑用户信息对话框 -->
    <UserEditDialog
      :open="showEditDialog"
      :user="editingUser"
      :loading="loading"
      @confirm="saveUser"
      @cancel="closeEditDialog"
      @uploadAvatar="handleAvatarUpload"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { uploadApi } from '@/services/upload'
import { type User, userApi, type UserFormData } from '@/services/userApi'
import { useUIStore } from '@/stores/ui'
import { calculateFileHash } from '@/utils/fileUtils'
import { useToast } from '@/utils/toastUtils'

import UserEditDialog from './components/UserEditDialog.vue'
import UserMobileList from './components/UserMobileList.vue'
import UserTable from './components/UserTable.vue'

// 情侣用户数据
const users = ref<User[]>([])
const uiStore = useUIStore()
// Toast 通知状态
const showToast = useToast()
const loading = ref(false)
// 获取用户列表
const fetchUsers = async () => {
  loading.value = true
  uiStore.setLoading(true)
  try {
    const response = await userApi.getUsers()
    if (response.code === 0) {
      users.value = response.data
    } else {
      showToast('获取用户列表失败', 'error')
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
    showToast('获取用户列表失败', 'error')
  } finally {
    uiStore.setLoading(false)
    loading.value = false
  }
}

// 编辑对话框相关状态
const showEditDialog = ref(false)
const editingUser = ref<UserFormData>({
  id: 0,
  name: '',
  email: '',
  role: '',
  joinDate: '',
  url: '',
  thumbnailUrl: '',
  avatarId: 0,
})

// 打开编辑对话框
const openEditDialog = (user: User) => {
  editingUser.value = {
    id: user.id,
    name: user.name,
    email: user.email,
    role: user.role,
    joinDate: user.joinDate,
    url: '',
    thumbnailUrl: '',
    avatarId: 0,
  }
  showEditDialog.value = true
}

// 关闭编辑对话框
const closeEditDialog = () => {
  showEditDialog.value = false
  editingUser.value = {
    id: 0,
    name: '',
    email: '',
    role: '',
    joinDate: '',
    url: '',
    thumbnailUrl: '',
    avatarId: 0,
  }
}

// 保存用户信息
const saveUser = async (user: UserFormData) => {
  if (!editingUser.value) {
    showToast('编辑用户信息为空', 'error')
    return
  }

  try {
    uiStore.setLoading(true)
    loading.value = true

    // 更新请求，不传递role字段，使用后端默认值
    const response = await userApi.updateUser(user.id, {
      name: user.name,
      email: user.email,
      avatarId: user.avatarId,
      newPassword: user.newPassword,
      role: user.role,
    })

    if (response.data.code === 0) {
      const userIndex = users.value.findIndex(u => u.id === editingUser!.value!.id)
      if (userIndex !== -1) {
        users.value[userIndex] = { ...user }
      }

      showToast('用户信息已保存', 'success')
    } else {
      showToast(response.data.message || '保存失败', 'error')
    }
  } catch (error) {
    console.error('保存用户信息失败:', error)
    showToast('保存失败，请重试', 'error')
  } finally {
    uiStore.setLoading(false)
    loading.value = false
    closeEditDialog()
  }
}

// 头像上传函数
const handleAvatarUpload = async (event: Event) => {
  uiStore.setLoading(true)
  loading.value = true
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    if (!file) {
      showToast('请选择文件', 'error')
      return
    }

    try {
      // 计算文件哈希值
      const hash = await calculateFileHash(file)

      // 创建 FormData 对象，用于上传文件
      const formData = new FormData()
      formData.append('file', file)
      formData.append('hash', hash)
      // 添加path字段，设置为avatar/{userId}格式
      const path = `avatar/${editingUser.value.id}`
      formData.append('path', path)

      // 调用上传接口，上传头像文件
      const { data } = await uploadApi.uploadImage(formData)

      // 处理响应
      if (data.code === 0) {
        const { file } = data.data

        // 更新头像的 URL 和 ID
        if (editingUser.value) {
          editingUser.value.url = file.url
          editingUser.value.thumbnailUrl = file.thumbnail
          editingUser.value.avatarId = file.id
        }
        showToast('头像上传成功', 'success')
      } else {
        // 处理上传失败的情况
        showToast('头像上传失败', 'error')
      }
    } catch {
      showToast('头像上传失败，请稍后重试', 'error')
    } finally {
      uiStore.setLoading(false)
      loading.value = false
    }
  } else {
    uiStore.setLoading(false)
    loading.value = false
  }
}

// 组件挂载时获取用户列表
onMounted(() => {
  fetchUsers()
})
</script>
