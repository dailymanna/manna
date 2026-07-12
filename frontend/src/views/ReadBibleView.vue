<template>
  <div class="view-read">

    <div class="bible-nav">
      <select class="bible-select" @change="updateChaptersAndVerses" v-model="selectedBook">
        <option v-for="b in books" :key="b" :value="b">{{ b }}</option>
      </select>
      <select class="bible-select" @change="updateChaptersAndVerses" v-model="selectedChapter">
        <option v-for="c in chapters" :key="c" :value="c">Chapter {{ c }}</option>
      </select>
      <select class="bible-select" @change="updateChaptersAndVerses" v-model="selectedTranslation">
        <option v-for="(translation, tindex) in listOfTranslationsAvailable" :key="tindex">{{ translation }}</option>

        <!-- <option>KJV</option>
        <option>NIV</option>
        <option>ESV</option>
        <option>NLT</option> -->
      </select>
    </div>

    <div class="scripture-body">
      <h2 class="scripture-heading">{{ selectedBook }} · Chapter {{ selectedChapter }}</h2>
      <div class="verses">
        <BibleVersesList :verses="versesInReadView" @verseNumEvent="onVerseClick"></BibleVersesList>
      </div>
    </div>

  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useUiStore } from '@/stores/ui'
import { GetCrossReferences, GetBooksOfTheBible, GetCountOfChaptersInTheBook, GetListOfTranslationsAvailable, GetVerses, ShowError, ShowWarning } from '@@/bindings/github.com/dailymanna/manna/internal/services/bible/bibleservice'
import { Application } from '@wailsio/runtime'
import BibleVersesList from '@/components/bible/BibleVersesList.vue'
import ReadBibleRightPanel from '@/components/bible/ReadBibleRightPanel.vue'

const ui = useUiStore()

let selectedBook = ref('Genesis')
let selectedChapter = ref(1)
let selectedVerse = ref({ num: 0, text: 'No verse selected' })
let selectedTranslation = ref('NKJV')
let defaultTranslation = ref('NKJV')
let listOfTranslationsAvailable = ref([])

// const books    = ['Genesis','Exodus','Psalms','Proverbs','Matthew','John','Romans','Revelation']
let books = ref([])
// const chapters = Array.from({ length: 50 }, (_, i) => i + 1)
let chapters = ref([])

let numOfChapters = ref(0)

let versesInReadView = ref([])

let crossReferences = ref([])



const updateChaptersAndVerses = async () => {
  await GetCountOfChaptersInTheBook(selectedBook.value).then(async (data) => {
    numOfChapters.value = data
    chapters.value = Array.from({ length: data }, (_, i) => i + 1)
    await GetVerses(selectedTranslation.value, selectedBook.value, selectedChapter.value).then(async (verses) => {
      versesInReadView.value = [...verses] // Doing this ensures that the ref variable is updated
      await onVerseClick(selectedVerse.value.num)
    }).catch(async (err) => {
      console.log("Error in getting verses:", err)
      await ShowWarning("Bible version not found", `The bible version ${selectedTranslation.value} is not available. Switching to default(${defaultTranslation.value}) version.`)
      selectedTranslation.value = 'KJV'
    })
  }).catch(async (err) => {
    await ShowError("Error occured", err)
  })
}

const fetchReadViewData = async () => {

  // Fetch all the available translations
  await GetListOfTranslationsAvailable().then((translations) => {
    listOfTranslationsAvailable.value = translations
  }).catch(async (err) => {
    await ShowError("Error occured", err)
  })

  // Fetch all the books of the bible
  await GetBooksOfTheBible().then(async (data) => {
    books.value = data
    selectedBook.value = data[0]
    await updateChaptersAndVerses()
  }).catch(async (err) => {
    await ShowError("Error occured", err)
  })
}

const onVerseClick = async (num) => {
  const verse = versesInReadView.value.find((v) => v.num === num)
  if (!verse) {
    console.warn(`Verse not found: ${num}`)
    return
  }
  selectedVerse.value = verse
  const crInpt = {
    book: selectedBook.value,
    chapter: selectedChapter.value,
    verse_number: Number(selectedVerse.value.num)
  }
  await GetCrossReferences(crInpt).then((result) => {
    crossReferences.value = result.cross_references
    console.log("References:", crossReferences.value)
  }).catch(async (err) => {
    console.log("No cross references found")
  })
  await ui.setRightPanel(ReadBibleRightPanel, {
    title: 'Study Tools',
    book: selectedBook.value,
    chapter: selectedChapter.value,
    selectedVerse: selectedVerse.value,   // live data from center → right
    crossReferences: crossReferences.value,
  })
}

onMounted(async () => {
  await fetchReadViewData()

  await ui.setRightPanel(ReadBibleRightPanel, {
    title: 'Study Tools',
    book: selectedBook.value,
    chapter: selectedChapter.value,
    selectedVerse: selectedVerse.value,   // live data from center → right
    crossReferences: crossReferences.value,
  })
})
</script>

<style scoped>
.view-read {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ── Nav bar ────────────────────────────────── */
.bible-nav {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.bible-select {
  appearance: none;
  background: var(--white);
  border: 1px solid var(--border);
  border-radius: 7px;
  padding: 8px 14px;
  font-family: inherit;
  font-size: 0.88rem;
  color: var(--ink);
  cursor: pointer;
}

.bible-select:focus {
  outline: 2px solid var(--gold);
  outline-offset: 1px;
}

/* ── Scripture ──────────────────────────────── */
.scripture-heading {
  font-size: 1.4rem;
  color: var(--ink);
  margin-bottom: 20px;
  padding-bottom: 14px;
  border-bottom: 1px solid var(--border-light);
  letter-spacing: 0.02em;
}
</style>