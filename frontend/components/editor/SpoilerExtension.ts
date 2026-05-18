import { Node, mergeAttributes, wrappingInputRule } from '@tiptap/core'

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    spoiler: {
      toggleSpoiler: () => ReturnType
    }
  }
}

export const Spoiler = Node.create({
  name: 'spoiler',
  group: 'block',
  content: 'block+',
  defining: true,

  parseHTML() {
    return [{ tag: 'div[data-spoiler]' }]
  },

  renderHTML({ HTMLAttributes }) {
    return ['div', mergeAttributes(HTMLAttributes, {
      'data-spoiler': '',
      class: 'spoiler-block',
    }), 0]
  },

  addCommands() {
    return {
      toggleSpoiler:
        () =>
        ({ commands }) => commands.toggleWrap(this.name),
    }
  },

  addInputRules() {
    return [
      wrappingInputRule({
        find: /^!!!\s$/,
        type: this.type,
      }),
    ]
  },
})
