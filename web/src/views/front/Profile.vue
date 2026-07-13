<template>
  <div class="profile-page">
    <div class="profile-hero">
      <div class="profile-avatar-wrap">
        <img v-if="form.avatar" :src="form.avatar" alt="avatar preview" class="profile-avatar" />
        <div v-else class="profile-avatar profile-avatar-fallback">
          {{ (form.nickname || userStore.user?.nickname || 'U').slice(0, 1).toUpperCase() }}
        </div>
      </div>
      <div class="profile-hero-copy">
        <h1>个人中心</h1>
        <p>修改昵称和头像后，导航栏会立即同步。</p>
      </div>
    </div>

    <div class="profile-card card">
      <div class="form-group">
        <label class="form-label">用户名</label>
        <input :value="userStore.user?.username || ''" class="input" type="text" disabled />
      </div>

      <div class="form-group">
        <label class="form-label">昵称</label>
        <input v-model="form.nickname" class="input" type="text" maxlength="50" placeholder="请输入昵称" />
      </div>

      <div class="form-group">
        <label class="form-label">头像</label>
        <div class="avatar-row">
          <input
            v-model="form.avatar"
            class="input"
            type="text"
            maxlength="255"
            placeholder="粘贴图片地址，或点击右侧上传"
          />
          <label class="upload-btn">
            上传头像
            <input type="file" accept="image/*" hidden @change="handleAvatarUpload" />
          </label>
        </div>
      </div>

      <div class="profile-actions">
        <button class="btn btn-primary" :disabled="saving" @click="handleSave">
          {{ saving ? '保存中...' : '保存资料' }}
        </button>
        <button class="btn btn-outline" type="button" @click="resetForm">重置</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { uploadFile } from '../../api/article'
import { updateProfile } from '../../api/user'
import { useUserStore } from '../../stores/user'
import { message } from '../../utils/message'

const userStore = useUserStore()
const saving = ref(false)
const form = ref(buildForm())

function buildForm() {
  return {
    nickname: userStore.user?.nickname || '',
    avatar: userStore.user?.avatar || ''
  }
}

function resetForm() {
  form.value = buildForm()
}

async function handleAvatarUpload(event) {
  const file = event.target.files?.[0]
  if (!file) return

  try {
    const res = await uploadFile(file)
    form.value.avatar = res.data?.url || res.data || ''
    message.success('头像上传成功')
  } catch (error) {
    message.error('头像上传失败')
  } finally {
    event.target.value = ''
  }
}

async function handleSave() {
  if (!form.value.nickname.trim()) {
    message.warning('请输入昵称')
    return
  }

  saving.value = true
  try {
    const res = await updateProfile({
      nickname: form.value.nickname,
      avatar: form.value.avatar
    })
    userStore.updateUserProfile(res.data)
    form.value = buildForm()
    message.success('个人资料已更新')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.profile-page {
  max-width: 760px;
  margin: 0 auto;
}

.profile-hero {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 24px 0 28px;
}

.profile-avatar {
  width: 88px;
  height: 88px;
  border-radius: 24px;
  object-fit: cover;
  border: 3px solid rgba(56, 189, 248, 0.2);
  box-shadow: 0 20px 45px rgba(15, 23, 42, 0.22);
}

.profile-avatar-fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #22c55e, #06b6d4);
  color: #04131d;
  font-size: 30px;
  font-weight: 700;
}

.profile-hero-copy h1 {
  margin-bottom: 8px;
  font-size: 32px;
}

.profile-hero-copy p {
  color: var(--text-muted);
  font-size: 15px;
}

.profile-card {
  padding: 32px;
}

.form-group {
  margin-bottom: 24px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
}

.avatar-row {
  display: flex;
  gap: 10px;
}

.avatar-row .input {
  flex: 1;
}

.upload-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0 18px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--bg-secondary);
  color: var(--text-secondary);
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s ease;
}

.upload-btn:hover {
  color: var(--accent);
  border-color: var(--accent);
}

.profile-actions {
  display: flex;
  gap: 12px;
  padding-top: 8px;
}

@media (max-width: 768px) {
  .profile-hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .avatar-row,
  .profile-actions {
    flex-direction: column;
  }

  .profile-card {
    padding: 24px;
  }
}
</style>
