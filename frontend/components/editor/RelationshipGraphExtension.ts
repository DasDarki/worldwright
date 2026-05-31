import { Node, mergeAttributes } from '@tiptap/core'

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    relationshipGraph: {
      insertRelationshipGraph: (attrs: { entityIds: number[] }) => ReturnType
    }
  }
}

// Block-level atom that stores a list of entity IDs. Rendering of the actual
// graph is handled by the RelationshipGraph Vue component on the read view;
// inside the editor we render a static placeholder card with the count and a
// hint to re-pick the contents.
export const RelationshipGraphNode = Node.create({
  name: 'relationshipGraph',
  group: 'block',
  atom: true,
  selectable: true,
  draggable: true,

  addAttributes() {
    return {
      entityIds: {
        default: [] as number[],
        parseHTML: (el) => {
          const raw = el.getAttribute('data-entity-ids') || ''
          if (!raw) return []
          try { return (JSON.parse(raw) as unknown[]).map(Number).filter((n) => Number.isFinite(n)) }
          catch { return [] }
        },
        renderHTML: (attrs) => ({
          'data-entity-ids': JSON.stringify((attrs.entityIds as number[]) || []),
        }),
      },
    }
  },

  parseHTML() {
    return [{ tag: 'div[data-relationship-graph]' }]
  },

  renderHTML({ HTMLAttributes }) {
    const ids: number[] = JSON.parse((HTMLAttributes['data-entity-ids'] as string) || '[]')
    return [
      'div',
      mergeAttributes(HTMLAttributes, { 'data-relationship-graph': '', class: 'ww-graph-placeholder' }),
      ['span', { class: 'ww-graph-placeholder__head' }, 'Relationship graph'],
      ['span', { class: 'ww-graph-placeholder__count' }, `${ids.length} entries`],
    ]
  },

  addCommands() {
    return {
      insertRelationshipGraph:
        (attrs) =>
        ({ chain }) => chain().focus().insertContent({ type: this.name, attrs }).run(),
    }
  },
})
