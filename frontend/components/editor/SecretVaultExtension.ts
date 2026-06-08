import { Node, mergeAttributes } from '@tiptap/core'

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    secretVault: {
      insertSecretVault: (attrs: { vault_id: string; author_id: number }) => ReturnType
      toggleSecretVault: (attrs: { vault_id: string; author_id: number }) => ReturnType
    }
  }
}

// Block-level container with paragraph-or-block content. Each vault has a
// stable vault_id (uuid) so the backend can match it across edits and an
// author_id that pins read access to the user who created it. The
// "redacted" flag is set by the server when serving the body to anyone who
// is not the author; the editor uses it to render a sealed placeholder and
// strip the flag on save (the backend re-asserts the redaction state).
export const SecretVault = Node.create({
  name: 'secretVault',
  group: 'block',
  content: 'block+',
  defining: true,
  isolating: true,

  addAttributes() {
    return {
      vault_id: {
        default: '',
        parseHTML: (el) => el.getAttribute('data-vault-id') || '',
        renderHTML: (attrs) => ({ 'data-vault-id': attrs.vault_id }),
      },
      author_id: {
        default: 0,
        parseHTML: (el) => Number(el.getAttribute('data-author-id')) || 0,
        renderHTML: (attrs) => ({ 'data-author-id': String(attrs.author_id || 0) }),
      },
      redacted: {
        default: false,
        parseHTML: () => false, // never trust HTML; runtime flag only
        renderHTML: () => ({}),  // never serialize
      },
    }
  },

  parseHTML() {
    return [{ tag: 'aside[data-secret-vault]' }]
  },

  renderHTML({ HTMLAttributes }) {
    return [
      'aside',
      mergeAttributes(HTMLAttributes, {
        'data-secret-vault': '',
        class: 'ww-vault',
      }),
      0,
    ]
  },

  addCommands() {
    return {
      insertSecretVault:
        (attrs) =>
        ({ commands }) =>
          commands.insertContent({
            type: this.name,
            attrs,
            content: [{ type: 'paragraph' }],
          }),
      toggleSecretVault:
        (attrs) =>
        ({ commands }) =>
          commands.toggleWrap(this.name, attrs),
    }
  },
})
