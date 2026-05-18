import { Node, mergeAttributes } from '@tiptap/core'

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    wikilink: {
      insertWikilink: (attrs: { slug: string; label?: string }) => ReturnType
    }
  }
}

export const Wikilink = Node.create({
  name: 'wikilink',
  group: 'inline',
  inline: true,
  atom: true,
  selectable: true,

  addAttributes() {
    return {
      slug:  { default: '' },
      label: { default: '' },
    }
  },

  parseHTML() {
    return [{ tag: 'a[data-wikilink]' }]
  },

  renderHTML({ HTMLAttributes }) {
    const label = HTMLAttributes.label || HTMLAttributes.slug || ''
    return [
      'a',
      mergeAttributes(HTMLAttributes, {
        'data-wikilink': '',
        'data-slug': HTMLAttributes.slug,
        class: 'wikilink-chip',
        contenteditable: 'false',
      }),
      label,
    ]
  },

  addCommands() {
    return {
      insertWikilink:
        (attrs) =>
        ({ chain }) => chain()
          .insertContent({ type: this.name, attrs })
          .insertContent(' ')
          .run(),
    }
  },
})
