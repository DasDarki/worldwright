import { Node, mergeAttributes } from '@tiptap/core'

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    wwimage: {
      insertWWImage: (attrs: { src: string; alt?: string; assetId?: number | null }) => ReturnType
    }
  }
}

export const WWImage = Node.create({
  name: 'wwimage',
  group: 'block',
  atom: true,
  selectable: true,
  draggable: true,

  addAttributes() {
    return {
      src:     { default: '' },
      alt:     { default: '' },
      assetId: { default: null, parseHTML: (el) => Number(el.getAttribute('data-asset-id')) || null, renderHTML: (a) => a.assetId ? { 'data-asset-id': a.assetId } : {} },
    }
  },

  parseHTML() {
    return [{ tag: 'img[src]' }]
  },

  renderHTML({ HTMLAttributes }) {
    return ['img', mergeAttributes(HTMLAttributes, { class: 'ww-img-embed' })]
  },

  addCommands() {
    return {
      insertWWImage:
        (attrs) =>
        ({ chain }) => chain().focus().insertContent({ type: this.name, attrs }).run(),
    }
  },
})
