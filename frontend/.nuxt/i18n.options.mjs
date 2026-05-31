
// @ts-nocheck
import locale_en_46json_d20ffe24 from "#nuxt-i18n/d20ffe24";
import locale_de_46json_3a94d79e from "#nuxt-i18n/3a94d79e";

export const localeCodes =  [
  "en",
  "de"
]

export const localeLoaders = {
  en: [
    {
      key: "locale_en_46json_d20ffe24",
      load: () => Promise.resolve(locale_en_46json_d20ffe24),
      cache: true
    }
  ],
  de: [
    {
      key: "locale_de_46json_3a94d79e",
      load: () => Promise.resolve(locale_de_46json_3a94d79e),
      cache: true
    }
  ]
}

export const vueI18nConfigs = []

export const nuxtI18nOptions = {
  restructureDir: "i18n",
  experimental: {
    localeDetector: "",
    switchLocalePathLinkSSR: false,
    autoImportTranslationFunctions: false,
    typedPages: true,
    typedOptionsAndMessages: false,
    generatedLocaleFilePathFormat: "absolute",
    alternateLinkCanonicalQueries: false,
    hmr: true
  },
  bundle: {
    compositionOnly: true,
    runtimeOnly: false,
    fullInstall: true,
    dropMessageCompiler: false,
    optimizeTranslationDirective: true
  },
  compilation: {
    strictMessage: true,
    escapeHtml: false
  },
  customBlocks: {
    defaultSFCLang: "json",
    globalSFCScope: false
  },
  locales: [
    {
      code: "en",
      name: "English",
      files: [
        {
          path: "/home/dasdarki/Development/PnP/worldwright/frontend/i18n/locales/en.json",
          cache: undefined
        }
      ]
    },
    {
      code: "de",
      name: "Deutsch",
      files: [
        {
          path: "/home/dasdarki/Development/PnP/worldwright/frontend/i18n/locales/de.json",
          cache: undefined
        }
      ]
    }
  ],
  defaultLocale: "en",
  defaultDirection: "ltr",
  routesNameSeparator: "___",
  trailingSlash: false,
  defaultLocaleRouteNameSuffix: "default",
  strategy: "no_prefix",
  lazy: false,
  langDir: "locales",
  rootRedirect: undefined,
  detectBrowserLanguage: {
    alwaysRedirect: false,
    cookieCrossOrigin: false,
    cookieDomain: null,
    cookieKey: "worldwright_locale",
    cookieSecure: false,
    fallbackLocale: "",
    redirectOn: "no prefix",
    useCookie: true
  },
  differentDomains: false,
  baseUrl: "",
  customRoutes: "page",
  pages: {},
  skipSettingLocaleOnNavigate: false,
  types: "composition",
  debug: false,
  parallelPlugin: false,
  multiDomainLocales: false,
  i18nModules: []
}

export const normalizedLocales = [
  {
    code: "en",
    name: "English",
    files: [
      {
        path: "/home/dasdarki/Development/PnP/worldwright/frontend/i18n/locales/en.json",
        cache: undefined
      }
    ]
  },
  {
    code: "de",
    name: "Deutsch",
    files: [
      {
        path: "/home/dasdarki/Development/PnP/worldwright/frontend/i18n/locales/de.json",
        cache: undefined
      }
    ]
  }
]

export const NUXT_I18N_MODULE_ID = "@nuxtjs/i18n"
export const parallelPlugin = false
export const isSSG = false
export const hasPages = true

export const DEFAULT_COOKIE_KEY = "i18n_redirected"
export const DEFAULT_DYNAMIC_PARAMS_KEY = "nuxtI18nInternal"
export const SWITCH_LOCALE_PATH_LINK_IDENTIFIER = "nuxt-i18n-slp"
/** client **/

/** client-end **/