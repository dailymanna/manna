<template>
  <div class="section">
    <S type="group" label="Profile">
      <div class="avatar-row">
        <div class="avatar">{{ initials }}</div>
        <div class="avatar-info">
          <span class="avatar-name">{{ displayName }}</span>
          <span class="avatar-email">{{ email }}</span>
        </div>
        <button class="btn-edit" @click="editing = !editing">
          {{ editing ? 'Cancel' : 'Edit profile' }}
        </button>
      </div>
      <template v-if="editing">
        <S type="row" label="Display name">
          <input class="s-input" style="width:200px" v-model="displayName" />
        </S>
        <S type="row" label="Bio">
          <input class="s-input" style="width:200px" v-model="bio" placeholder="A short bio…" />
        </S>
        <S type="action" label="" button-label="Save changes" @click="saveProfile" />
      </template>
    </S>

    <S type="group" label="Account">
      <S type="row" label="Email address">
        <span class="value-text">{{ email }}</span>
      </S>
      <S type="action" label="Change password"
         hint="You will receive a reset link at your email address."
         button-label="Change password"
         @click="changePassword" />
      <S type="action" label="Sign out of Manna Hub"
         hint="Your local data will not be affected."
         button-label="Sign out"
         @click="signOut" />
    </S>

    <S type="group" label="Upgrade">
      <S type="action" label="Upgrade Manna"
         :hint="upgradeStatus"
         button-label="Check for updates"
         @click="upgradeManna" />
    </S>
    <S type="group" label="Danger Zone">
      <S type="action" label="Delete account"
         hint="Permanently deletes your Hub account, posts and shared notes. Local data is kept."
         button-label="Delete account"
         :danger="true"
         @click="deleteAccount" />
    </S>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import S from '@/components/settings/SettingPrimitives.vue'
import { UpgradeMannaApp } from '@@/bindings/github.com/dailymanna/manna/internal/services/settings/settingsservice'
import { Events } from '@wailsio/runtime'
const displayName = ref('James Reuben')
const email       = ref('james@example.com')
const bio         = ref('')
const editing     = ref(false)
const upgradeStatus = ref('Check for a new version of Manna.')
const isUpgrading = ref(false)
const initials    = computed(() =>
  displayName.value.split(' ').map(w => w[0]).join('').toUpperCase().slice(0,2)
)
function saveProfile()    { editing.value = false }
function changePassword() { console.log('Change password') }
function signOut()        { console.log('Sign out') }
function deleteAccount()  { console.log('Delete account') }
const upgradeManna = async () => {

  console.log('Checking for upgrades')
  await UpgradeMannaApp()
    .then((resp) => {})
    .catch((err) => {
      console.log("Found error when upgrading")
    })
}

onMounted(() => {
  const upgradeStatusOfManna = Events.On('manna-upgrade-event', (event) => {
    upgradeStatus.value = '${event.data[status]}'
  })
})
</script>
<style scoped>
.section { display: flex; flex-direction: column; }
.avatar-row {
  display: flex; align-items: center; gap: 16px; padding: 16px 18px;
}
.avatar {
  width: 52px; height: 52px; border-radius: 50%;
  background: var(--gold); color: var(--sidebar-bg);
  display: flex; align-items: center; justify-content: center;
  font-size: 1.1rem; font-weight: 700; flex-shrink: 0;
}
.avatar-info { flex: 1; display: flex; flex-direction: column; gap: 3px; }
.avatar-name  { font-size: 0.95rem; color: var(--ink); font-weight: 600; }
.avatar-email { font-size: 0.78rem; color: var(--ink-light); font-family: sans-serif; }
.btn-edit {
  padding: 6px 14px; border-radius: 7px; border: 1px solid var(--border);
  background: var(--parchment-mid); font-family: inherit; font-size: 0.82rem;
  cursor: pointer; color: var(--ink);
}
.value-text { font-size: 0.85rem; color: var(--ink-light); font-family: sans-serif; }
</style>