import { Node, mergeAttributes } from '@tiptap/core'

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    callout: {
      setCallout: (attrs?: { variant?: 'info' | 'warn' | 'note' | 'lore' }) => ReturnType
      toggleCallout: (attrs?: { variant?: 'info' | 'warn' | 'note' | 'lore' }) => ReturnType
      unsetCallout: () => ReturnType
      cycleCalloutVariant: () => ReturnType
    }
  }
}

export type CalloutVariant = 'info' | 'warn' | 'note' | 'lore'

const VARIANTS: CalloutVariant[] = ['info', 'warn', 'note', 'lore']

export const Callout = Node.create({
  name: 'callout',
  group: 'block',
  content: 'block+',
  defining: true,

  addAttributes() {
    return {
      variant: {
        default: 'info' as CalloutVariant,
        parseHTML: (el) => (el.getAttribute('data-variant') as CalloutVariant) || 'info',
        renderHTML: (attrs) => ({ 'data-variant': attrs.variant }),
      },
    }
  },

  parseHTML() {
    return [{ tag: 'aside[data-callout]' }]
  },

  renderHTML({ HTMLAttributes }) {
    return ['aside', mergeAttributes(HTMLAttributes, { 'data-callout': '', class: 'ww-callout' }), 0]
  },

  addCommands() {
    return {
      setCallout:
        (attrs) =>
        ({ commands }) =>
          commands.wrapIn(this.name, attrs || { variant: 'info' }),
      toggleCallout:
        (attrs) =>
        ({ commands }) =>
          commands.toggleWrap(this.name, attrs || { variant: 'info' }),
      unsetCallout:
        () =>
        ({ commands }) =>
          commands.lift(this.name),
      cycleCalloutVariant:
        () =>
        ({ editor, commands }) => {
          const cur = (editor.getAttributes(this.name).variant as CalloutVariant) || 'info'
          const next = VARIANTS[(VARIANTS.indexOf(cur) + 1) % VARIANTS.length]
          return commands.updateAttributes(this.name, { variant: next })
        },
    }
  },

  addKeyboardShortcuts() {
    return {
      'Mod-Alt-c': () => this.editor.commands.toggleCallout({ variant: 'info' }),
    }
  },
})
