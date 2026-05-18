import type { Config } from 'tailwindcss'

export default {
  content: [
    './app.vue',
    './error.vue',
    './pages/**/*.{vue,ts}',
    './layouts/**/*.{vue,ts}',
    './components/**/*.{vue,ts}',
    './composables/**/*.{ts}',
    './plugins/**/*.{ts}',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        parchment:         'rgb(var(--ww-parchment) / <alpha-value>)',
        'parchment-deep':  'rgb(var(--ww-parchment-deep) / <alpha-value>)',
        'parchment-stain': 'rgb(var(--ww-parchment-stain) / <alpha-value>)',
        ink:               'rgb(var(--ww-ink) / <alpha-value>)',
        'ink-shade':       'rgb(var(--ww-ink-shade) / <alpha-value>)',
        gold:              'rgb(var(--ww-gold) / <alpha-value>)',
        'gold-bright':     'rgb(var(--ww-gold-bright) / <alpha-value>)',
        'gold-deep':       'rgb(var(--ww-gold-deep) / <alpha-value>)',
        vermilion:         'rgb(var(--ww-vermilion) / <alpha-value>)',
        'vermilion-deep':  'rgb(var(--ww-vermilion-deep) / <alpha-value>)',
      },
      fontFamily: {
        display: ['Fraunces', 'EB Garamond', 'Georgia', 'serif'],
        body:    ['EB Garamond', 'Georgia', 'serif'],
        caps:    ['Cormorant SC', 'EB Garamond', 'Georgia', 'serif'],
      },
      letterSpacing: {
        codex: '0.32em',
      },
      transitionTimingFunction: {
        'codex-out': 'cubic-bezier(.22, 1, .36, 1)',
        'codex-soft': 'cubic-bezier(.65, 0, .35, 1)',
      },
      keyframes: {
        rise: {
          '0%': { opacity: '0', transform: 'translateY(20px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        twinkle: {
          '0%, 100%': { opacity: '.4', transform: 'scale(.9)' },
          '50%': { opacity: '1', transform: 'scale(1.1)' },
        },
        slowSpin: {
          '0%': { transform: 'rotate(0deg)' },
          '100%': { transform: 'rotate(360deg)' },
        },
        dawn: {
          '0%': { opacity: '0', transform: 'scale(.94)' },
          '100%': { opacity: '1', transform: 'scale(1)' },
        },
      },
      animation: {
        rise:     'rise 1.1s cubic-bezier(.22,1,.36,1) forwards',
        twinkle:  'twinkle 4s ease-in-out infinite',
        'slow-spin': 'slowSpin 90s linear infinite',
        dawn:     'dawn 1.6s cubic-bezier(.22,1,.36,1) forwards',
      },
    },
  },
  plugins: [],
} satisfies Config
