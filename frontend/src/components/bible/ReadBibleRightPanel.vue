<template>
    <aside class="right-panel">

        <div class="right-header">
            <h3 class="right-title">{{ ui.rightPanelContext.props.title ?? '' }}</h3>
            <button class="close-btn" @click="ui.closeRightPanel">✕</button>
        </div>

        <div class="right-content">
            <div class="rp-section">
                <div class="rp-section-label">Verse</div>
                <div class="cross-ref">
                    <!-- <span class="cross-ref-tag">1</span> -->
                    <span class="verse-text">{{selectedVerse.text}}</span>
                </div>
            </div>
            <div class="rp-section">
                <div class="rp-section-label" @click="printy">Cross-References</div>
                <div class="cross-ref" v-for="refer in crossReferences" :key="refer.id">
                    <span class="cross-ref-tag">{{ refer.to_book }} {{ refer.to_chapter }}:{{ refer.to_verse_start }}-{{ refer.to_verse_end }}</span>
                    <span class="cross-ref-text">{{ refer.to_verse_text }}</span>
                </div>
                <!-- <div class="cross-ref" v-for="ref in crossRefs" :key="ref.ref">
                    <span class="cross-ref-tag">{{ ref.ref }}</span>
                    <span class="cross-ref-text">{{ ref.snippet }}</span>
                </div> -->
            </div>

            <div class="rp-section">
                <div class="rp-section-label">Commentary</div>
                <p class="commentary-text">
                    This passage is foundational to understanding the covenant narrative that runs
                    through all of Scripture. Scholars note the parallel structure between the
                    creation account and the tabernacle instructions.
                </p>
            </div>
            <div class="rp-section">
                <div class="rp-section-label">Original Language</div>
                <div class="orig-word" v-for="w in origWords" :key="w.word">
                    <span class="script">{{ w.word }}</span>
                    <span class="trans">{{ w.meaning }}</span>
                </div>
            </div>
        </div>
    </aside>
</template>

<script setup>
import { useUiStore } from '@/stores/ui';

const ui = useUiStore()

const props = defineProps({
    selectedVerse: {
        type: Object,
        required: true,
        default: () => ({ num: 0, text: 'No verse selected' })
    },
    crossReferences: {
        type: Array,
        required: true,
        default: () => []
    }
})

const printy = () => {
    console.log("Cross References: ", props.crossReferences)
}

// const crossRefs = [
//     { ref: 'John 1:1', snippet: '"In the beginning was the Word…"' },
//     { ref: 'Heb 11:3', snippet: 'By faith we understand the universe was created…' },
//     { ref: 'Col 1:16', snippet: 'For in him all things were created…' },
// ]
const origWords = [
    { word: 'בְּרֵאשִׁית', meaning: "b'reshit — In the beginning" },
    { word: 'בָּרָא', meaning: 'bara — created (divine act)' },
    { word: 'אֱלֹהִים', meaning: 'Elohim — God (plural majesty)' },
]

</script>

<style lang="css" scoped>
.right-panel {
  grid-column: 3;
  background: var(--white);
  border-left: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  width: var(--right-panel-w);
  min-height: 0;
  height: 100%;
  overflow: hidden;
}

.right-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: var(--header-h);
  border-bottom: 1px solid var(--border-light);
  flex-shrink: 0;
}

.right-title {
  font-size: 1rem;
  color: var(--ink);
  letter-spacing: 0.04em;
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--ink-light);
  font-size: 0.85rem;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background var(--transition);
}

.close-btn:hover {
  background: var(--parchment-mid);
  color: var(--ink);
}

.right-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}






/* ── Tablet: fixed overlay ──────────────────── */
@media (max-width: 1024px) {
  .right-panel {
    position: fixed;
    right: 0;
    top: 0;
    bottom: 0;
    z-index: 100;
    box-shadow: -4px 0 20px rgba(44, 31, 14, 0.15);
  }
}

/* ── Mobile: full-width overlay ─────────────── */
@media (max-width: 640px) {
  .right-panel {
    width: 100% !important;
    max-width: 100%;
  }
}
/* ── Section ────────────────────────────────── */
.rp-section {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.rp-section-label {
    font-size: 0.68rem;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: var(--ink-light);
    font-family: sans-serif;
    font-weight: 700;
    border-bottom: 1px solid var(--border-light);
    padding-bottom: 6px;
}

/* ── Cross-refs ─────────────────────────────── */
.cross-ref {
    display: flex;
    flex-direction: column;
    gap: 3px;
    padding: 8px 0;
    border-bottom: 1px solid var(--border-light);
}

.cross-ref:last-child {
    border-bottom: none;
}

.cross-ref-tag {
    font-size: 0.75rem;
    font-weight: 700;
    color: var(--gold);
    font-family: sans-serif;
}

.cross-ref-text {
    font-size: 0.82rem;
    color: var(--ink-light);
    font-style: italic;
}

/* ── Commentary ─────────────────────────────── */
.verse-text {
    font-size: 0.88rem;
    line-height: 1.7;
    color: var(--ink-mid);
}
/* ── Commentary ─────────────────────────────── */
.commentary-text {
    font-size: 0.88rem;
    line-height: 1.7;
    color: var(--ink-mid);
}

/* ── Original language ──────────────────────── */
.orig-word {
    display: flex;
    gap: 14px;
    align-items: baseline;
    padding: 6px 0;
}

.script {
    font-size: 1rem;
    color: var(--ink);
}

.trans {
    font-size: 0.8rem;
    color: var(--ink-light);
    font-family: sans-serif;
    font-style: italic;
}
</style>