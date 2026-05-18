<script setup lang="ts">
const props = defineProps<{ modelValue: string[] }>()
const emit = defineEmits<{ 'update:modelValue': [v: string[]] }>()

const draft = ref('')
const inputRef = ref<HTMLInputElement | null>(null)

function addTag() {
  const t = draft.value.trim().toLowerCase().replace(/\s+/g, '-')
  if (!t) return
  if (props.modelValue.includes(t)) {
    draft.value = ''
    return
  }
  emit('update:modelValue', [...props.modelValue, t])
  draft.value = ''
}

function removeTag(t: string) {
  emit('update:modelValue', props.modelValue.filter((x) => x !== t))
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' || e.key === ',') {
    e.preventDefault()
    addTag()
  } else if (e.key === 'Backspace' && !draft.value && props.modelValue.length) {
    e.preventDefault()
    emit('update:modelValue', props.modelValue.slice(0, -1))
  }
}
</script>

<template>
  <div class="tag-input" @click="inputRef?.focus()">
    <TransitionGroup name="chip" tag="span" class="chips">
      <span v-for="t in modelValue" :key="t" class="chip">
        {{ t }}
        <button type="button" class="x" :aria-label="`Remove ${t}`" @click.stop="removeTag(t)">×</button>
      </span>
    </TransitionGroup>
    <input
      ref="inputRef"
      v-model="draft"
      class="entry"
      placeholder="add tag…"
      @keydown="onKeydown"
      @blur="addTag"
    />
  </div>
</template>

<style scoped lang="scss">
.tag-input {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  padding: 10px 4px;
  border-bottom: 1px solid var(--ww-ink-hairline);
  cursor: text;
  transition: border-color .3s ease;
}
.tag-input:focus-within { border-color: rgb(var(--ww-gold)); }
.chips {
  display: contents;
}
.chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  font-family: 'Cormorant SC', serif;
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: .22em;
  color: var(--ww-ink-faint);
  border: 1px solid var(--ww-ink-hairline);
  background: rgb(var(--ww-gold) / .04);
  transition: background-color .25s ease, transform .25s ease, color .25s ease;
}
.chip:hover { background: rgb(var(--ww-gold) / .15); color: rgb(var(--ww-gold-deep)); }
.x {
  background: none;
  border: 0;
  font: inherit;
  font-size: 14px;
  color: inherit;
  line-height: 1;
  padding: 0 2px;
  cursor: pointer;
  opacity: .6;
  transition: opacity .25s ease;
}
.x:hover { opacity: 1; color: rgb(var(--ww-vermilion)); }
.entry {
  flex: 1;
  min-width: 120px;
  border: 0;
  background: transparent;
  outline: none;
  font-family: 'EB Garamond', serif;
  font-size: 14px;
  color: rgb(var(--ww-ink));
  padding: 4px 0;
}

.chip-enter-active, .chip-leave-active { transition: opacity .3s ease, transform .35s cubic-bezier(.22,1,.36,1); }
.chip-enter-from { opacity: 0; transform: scale(.85) translateY(4px); }
.chip-leave-to { opacity: 0; transform: scale(.85); }
</style>
