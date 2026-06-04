<template>
  <div class="app-shell" :class="{
    'sidebar-collapsed': ui.sidebarCollapsed,
    'panel-closed': !ui.rightPanelOpen
  }">

    <LeftSidebar />

    <!-- CENTER — routed views live here -->
    <main class="center-panel">
      <div class="center-header">
        <div class="header-left">
          <h1 class="page-title">{{ routeMeta.title }}</h1>
          <p class="page-subtitle">{{ routeMeta.subtitle }}</p>
        </div>
        <div class="header-actions">
          <button
            v-if="!ui.disableRightPanel"
            class="action-btn"
            @click="ui.toggleRightPanel"
            :title="ui.rightPanelOpen ? 'Hide panel' : 'Show panel'">
            <span>{{ ui.rightPanelOpen ? '▶' : '◀' }}</span>
          </button>
        </div>
      </div>

      <div class="center-content">
        <RouterView v-slot="{ Component }">
          <!-- <Transition name="view-fade" mode="out-in"> -->
          <Transition v-if="ui.rightPanelOpen && !ui.disableRightPanel">
            <component :is="Component" :key="$route.path" />
          </Transition>
          <Transition name="view-fade" mode="out-in" v-else>
            <component :is="Component" :key="$route.path" />
          </Transition>
        </RouterView>
      </div>
    </main>

    <!-- RIGHT PANEL — context-aware, collapses via store -->
    <Transition name="right-panel-slide">
      <RightPanel v-if="ui.rightPanelOpen && !ui.disableRightPanel" />
    </Transition>

  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUiStore } from '@/stores/ui'
import RightPanel from '@/components/ui/RightPanel.vue'
import LeftSidebar from '@/components/ui/LeftSidebar.vue'
// import RightPanel from './components/ui/RightPanel.vue'
// import LeftSidebar from './components/ui/LeftSidebar.vue'
// import { useUiStore } from '@/stores/ui'

// import LeftSidebar from '@/components/LeftSidebar.vue'
// import RightPanel  from '@/components/RightPanel.vue'

const route = useRoute()
const ui = useUiStore()

const routeMeta = computed(() => ({
  title: route.meta?.title ?? '',
  subtitle: route.meta?.subtitle ?? '',
}))
</script>

<style scoped>
/* ── App shell grid ─────────────────────────── */
.app-shell {
  display: grid;
  grid-template-columns: var(--sidebar-w) 1fr var(--right-panel-w);
  grid-template-rows: 100vh;
  height: 100vh;
  overflow: hidden;
  background: var(--parchment);
  transition: grid-template-columns var(--transition);
}

.app-shell.sidebar-collapsed {
  grid-template-columns: var(--sidebar-w-min) 1fr var(--right-panel-w);
}

.app-shell.panel-closed {
  grid-template-columns: var(--sidebar-w) 1fr 0;
}

.app-shell.sidebar-collapsed.panel-closed {
  grid-template-columns: var(--sidebar-w-min) 1fr 0;
}

/* ── Center panel ───────────────────────────── */
.center-panel {
  grid-column: 2;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--parchment);
  min-width: 0;
}

.center-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 28px;
  height: var(--header-h);
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.page-title {
  font-size: 1.3rem;
  font-weight: 700;
  color: var(--ink);
  letter-spacing: 0.01em;
}

.page-subtitle {
  font-size: 0.8rem;
  color: var(--ink-light);
  margin-top: 1px;
  font-style: italic;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  background: none;
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 6px 10px;
  cursor: pointer;
  color: var(--ink-mid);
  font-size: 0.8rem;
  transition: background var(--transition), color var(--transition);
}

.action-btn:hover {
  background: var(--parchment-mid);
  color: var(--ink);
}

.center-content {
  flex: 1;
  overflow-y: auto;
  padding: 28px;
}

/* ── View transition ────────────────────────── */
.view-fade-enter-active,
.view-fade-leave-active {
  transition: opacity 0.18s ease, transform 0.18s ease;
}

.view-fade-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.view-fade-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

/* ── Right panel slide ──────────────────────── */
.right-panel-slide-enter-active,
.right-panel-slide-leave-active {
  transition: width var(--transition), opacity var(--transition);
  overflow: hidden;
}

.right-panel-slide-enter-from,
.right-panel-slide-leave-to {
  width: 0;
  opacity: 0;
}

.right-panel-slide-enter-to,
.right-panel-slide-leave-from {
  width: var(--right-panel-w);
  opacity: 1;
}

/* ── Tablet ─────────────────────────────────── */
@media (max-width: 1024px) {

  .app-shell,
  .app-shell.sidebar-collapsed,
  .app-shell.panel-closed,
  .app-shell.sidebar-collapsed.panel-closed {
    grid-template-columns: var(--sidebar-w) 1fr;
  }

  .app-shell.sidebar-collapsed,
  .app-shell.sidebar-collapsed.panel-closed {
    grid-template-columns: var(--sidebar-w-min) 1fr;
  }
}

/* ── Mobile ─────────────────────────────────── */
@media (max-width: 640px) {

  .app-shell,
  .app-shell.sidebar-collapsed {
    grid-template-columns: 1fr;
    grid-template-rows: 1fr 60px;
    height: 100dvh;
  }

  .center-panel {
    grid-column: 1;
    grid-row: 1;
  }

  .center-header {
    padding: 0 16px;
  }

  .center-content {
    padding: 16px;
  }
}
</style>