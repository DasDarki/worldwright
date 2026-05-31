
import type { DefineComponent, SlotsType } from 'vue'
type IslandComponent<T> = DefineComponent<{}, {refresh: () => Promise<void>}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, SlotsType<{ fallback: { error: unknown } }>> & T

type HydrationStrategies = {
  hydrateOnVisible?: IntersectionObserverInit | true
  hydrateOnIdle?: number | true
  hydrateOnInteraction?: keyof HTMLElementEventMap | Array<keyof HTMLElementEventMap> | true
  hydrateOnMediaQuery?: string
  hydrateAfter?: number
  hydrateWhen?: boolean
  hydrateNever?: true
}
type LazyComponent<T> = DefineComponent<HydrationStrategies, {}, {}, {}, {}, {}, {}, { hydrated: () => void }> & T

interface _GlobalComponents {
  AppHeader: typeof import("../../components/AppHeader.vue")['default']
  AppSidebar: typeof import("../../components/AppSidebar.vue")['default']
  AssetUpload: typeof import("../../components/AssetUpload.vue")['default']
  BodyNode: typeof import("../../components/BodyNode.vue")['default']
  BodyView: typeof import("../../components/BodyView.vue")['default']
  DateInput: typeof import("../../components/DateInput.vue")['default']
  EntityForm: typeof import("../../components/EntityForm.vue")['default']
  EntityPicker: typeof import("../../components/EntityPicker.vue")['default']
  EntityTreeNode: typeof import("../../components/EntityTreeNode.vue")['default']
  EventForm: typeof import("../../components/EventForm.vue")['default']
  EventsForEntityPanel: typeof import("../../components/EventsForEntityPanel.vue")['default']
  LangToggle: typeof import("../../components/LangToggle.vue")['default']
  MapCanvas: typeof import("../../components/MapCanvas.vue")['default']
  MobileDrawer: typeof import("../../components/MobileDrawer.vue")['default']
  MoonCycleTimeline: typeof import("../../components/MoonCycleTimeline.vue")['default']
  MoonDisplay: typeof import("../../components/MoonDisplay.vue")['default']
  MoonGlyph: typeof import("../../components/MoonGlyph.vue")['default']
  Ornament: typeof import("../../components/Ornament.vue")['default']
  PWAStatus: typeof import("../../components/PWAStatus.client.vue")['default']
  PaperOverlays: typeof import("../../components/PaperOverlays.vue")['default']
  PedigreeChart: typeof import("../../components/PedigreeChart.vue")['default']
  RelationshipEditor: typeof import("../../components/RelationshipEditor.vue")['default']
  RelationshipPanel: typeof import("../../components/RelationshipPanel.vue")['default']
  SearchBox: typeof import("../../components/SearchBox.vue")['default']
  ShareButton: typeof import("../../components/ShareButton.vue")['default']
  ShortcutHelp: typeof import("../../components/ShortcutHelp.vue")['default']
  SpoilerBlock: typeof import("../../components/SpoilerBlock.vue")['default']
  StarSparkle: typeof import("../../components/StarSparkle.vue")['default']
  TagInput: typeof import("../../components/TagInput.vue")['default']
  ThemeToggle: typeof import("../../components/ThemeToggle.vue")['default']
  TimelineCurator: typeof import("../../components/TimelineCurator.vue")['default']
  TimelineForm: typeof import("../../components/TimelineForm.vue")['default']
  ToastContainer: typeof import("../../components/ToastContainer.vue")['default']
  WikilinkPreview: typeof import("../../components/WikilinkPreview.vue")['default']
  Editor: typeof import("../../components/editor/Editor.client.vue")['default']
  EditorToolbar: typeof import("../../components/editor/EditorToolbar.vue")['default']
  ImageExtension: typeof import("../../components/editor/ImageExtension")['default']
  ImagePicker: typeof import("../../components/editor/ImagePicker.client.vue")['default']
  SpoilerExtension: typeof import("../../components/editor/SpoilerExtension")['default']
  WikilinkExtension: typeof import("../../components/editor/WikilinkExtension")['default']
  WikilinkPicker: typeof import("../../components/editor/WikilinkPicker.client.vue")['default']
  NuxtWelcome: typeof import("../../node_modules/nuxt/dist/app/components/welcome.vue")['default']
  NuxtLayout: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-layout")['default']
  NuxtErrorBoundary: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']
  ClientOnly: typeof import("../../node_modules/nuxt/dist/app/components/client-only")['default']
  DevOnly: typeof import("../../node_modules/nuxt/dist/app/components/dev-only")['default']
  ServerPlaceholder: typeof import("../../node_modules/nuxt/dist/app/components/server-placeholder")['default']
  NuxtLink: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-link")['default']
  NuxtLoadingIndicator: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']
  NuxtTime: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']
  NuxtRouteAnnouncer: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']
  NuxtImg: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtImg']
  NuxtPicture: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtPicture']
  NuxtLinkLocale: typeof import("../../node_modules/@nuxtjs/i18n/dist/runtime/components/NuxtLinkLocale")['default']
  SwitchLocalePathLink: typeof import("../../node_modules/@nuxtjs/i18n/dist/runtime/components/SwitchLocalePathLink")['default']
  ColorScheme: typeof import("../../node_modules/@nuxtjs/color-mode/dist/runtime/component.vue3.vue")['default']
  VitePwaManifest: typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/VitePwaManifest")['default']
  NuxtPwaManifest: typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/VitePwaManifest")['default']
  NuxtPwaAssets: typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/NuxtPwaAssets")['default']
  PwaAppleImage: typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaAppleImage.vue")['default']
  PwaAppleSplashScreenImage: typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaAppleSplashScreenImage.vue")['default']
  PwaFaviconImage: typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaFaviconImage.vue")['default']
  PwaMaskableImage: typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaMaskableImage.vue")['default']
  PwaTransparentImage: typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaTransparentImage.vue")['default']
  NuxtPage: typeof import("../../node_modules/nuxt/dist/pages/runtime/page")['default']
  NoScript: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['NoScript']
  Link: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Link']
  Base: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Base']
  Title: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Title']
  Meta: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Meta']
  Style: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Style']
  Head: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Head']
  Html: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Html']
  Body: typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Body']
  NuxtIsland: typeof import("../../node_modules/nuxt/dist/app/components/nuxt-island")['default']
  LazyAppHeader: LazyComponent<typeof import("../../components/AppHeader.vue")['default']>
  LazyAppSidebar: LazyComponent<typeof import("../../components/AppSidebar.vue")['default']>
  LazyAssetUpload: LazyComponent<typeof import("../../components/AssetUpload.vue")['default']>
  LazyBodyNode: LazyComponent<typeof import("../../components/BodyNode.vue")['default']>
  LazyBodyView: LazyComponent<typeof import("../../components/BodyView.vue")['default']>
  LazyDateInput: LazyComponent<typeof import("../../components/DateInput.vue")['default']>
  LazyEntityForm: LazyComponent<typeof import("../../components/EntityForm.vue")['default']>
  LazyEntityPicker: LazyComponent<typeof import("../../components/EntityPicker.vue")['default']>
  LazyEntityTreeNode: LazyComponent<typeof import("../../components/EntityTreeNode.vue")['default']>
  LazyEventForm: LazyComponent<typeof import("../../components/EventForm.vue")['default']>
  LazyEventsForEntityPanel: LazyComponent<typeof import("../../components/EventsForEntityPanel.vue")['default']>
  LazyLangToggle: LazyComponent<typeof import("../../components/LangToggle.vue")['default']>
  LazyMapCanvas: LazyComponent<typeof import("../../components/MapCanvas.vue")['default']>
  LazyMobileDrawer: LazyComponent<typeof import("../../components/MobileDrawer.vue")['default']>
  LazyMoonCycleTimeline: LazyComponent<typeof import("../../components/MoonCycleTimeline.vue")['default']>
  LazyMoonDisplay: LazyComponent<typeof import("../../components/MoonDisplay.vue")['default']>
  LazyMoonGlyph: LazyComponent<typeof import("../../components/MoonGlyph.vue")['default']>
  LazyOrnament: LazyComponent<typeof import("../../components/Ornament.vue")['default']>
  LazyPWAStatus: LazyComponent<typeof import("../../components/PWAStatus.client.vue")['default']>
  LazyPaperOverlays: LazyComponent<typeof import("../../components/PaperOverlays.vue")['default']>
  LazyPedigreeChart: LazyComponent<typeof import("../../components/PedigreeChart.vue")['default']>
  LazyRelationshipEditor: LazyComponent<typeof import("../../components/RelationshipEditor.vue")['default']>
  LazyRelationshipPanel: LazyComponent<typeof import("../../components/RelationshipPanel.vue")['default']>
  LazySearchBox: LazyComponent<typeof import("../../components/SearchBox.vue")['default']>
  LazyShareButton: LazyComponent<typeof import("../../components/ShareButton.vue")['default']>
  LazyShortcutHelp: LazyComponent<typeof import("../../components/ShortcutHelp.vue")['default']>
  LazySpoilerBlock: LazyComponent<typeof import("../../components/SpoilerBlock.vue")['default']>
  LazyStarSparkle: LazyComponent<typeof import("../../components/StarSparkle.vue")['default']>
  LazyTagInput: LazyComponent<typeof import("../../components/TagInput.vue")['default']>
  LazyThemeToggle: LazyComponent<typeof import("../../components/ThemeToggle.vue")['default']>
  LazyTimelineCurator: LazyComponent<typeof import("../../components/TimelineCurator.vue")['default']>
  LazyTimelineForm: LazyComponent<typeof import("../../components/TimelineForm.vue")['default']>
  LazyToastContainer: LazyComponent<typeof import("../../components/ToastContainer.vue")['default']>
  LazyWikilinkPreview: LazyComponent<typeof import("../../components/WikilinkPreview.vue")['default']>
  LazyEditor: LazyComponent<typeof import("../../components/editor/Editor.client.vue")['default']>
  LazyEditorToolbar: LazyComponent<typeof import("../../components/editor/EditorToolbar.vue")['default']>
  LazyImageExtension: LazyComponent<typeof import("../../components/editor/ImageExtension")['default']>
  LazyImagePicker: LazyComponent<typeof import("../../components/editor/ImagePicker.client.vue")['default']>
  LazySpoilerExtension: LazyComponent<typeof import("../../components/editor/SpoilerExtension")['default']>
  LazyWikilinkExtension: LazyComponent<typeof import("../../components/editor/WikilinkExtension")['default']>
  LazyWikilinkPicker: LazyComponent<typeof import("../../components/editor/WikilinkPicker.client.vue")['default']>
  LazyNuxtWelcome: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/welcome.vue")['default']>
  LazyNuxtLayout: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-layout")['default']>
  LazyNuxtErrorBoundary: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']>
  LazyClientOnly: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/client-only")['default']>
  LazyDevOnly: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/dev-only")['default']>
  LazyServerPlaceholder: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/server-placeholder")['default']>
  LazyNuxtLink: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-link")['default']>
  LazyNuxtLoadingIndicator: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']>
  LazyNuxtTime: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']>
  LazyNuxtRouteAnnouncer: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']>
  LazyNuxtImg: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtImg']>
  LazyNuxtPicture: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtPicture']>
  LazyNuxtLinkLocale: LazyComponent<typeof import("../../node_modules/@nuxtjs/i18n/dist/runtime/components/NuxtLinkLocale")['default']>
  LazySwitchLocalePathLink: LazyComponent<typeof import("../../node_modules/@nuxtjs/i18n/dist/runtime/components/SwitchLocalePathLink")['default']>
  LazyColorScheme: LazyComponent<typeof import("../../node_modules/@nuxtjs/color-mode/dist/runtime/component.vue3.vue")['default']>
  LazyVitePwaManifest: LazyComponent<typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/VitePwaManifest")['default']>
  LazyNuxtPwaManifest: LazyComponent<typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/VitePwaManifest")['default']>
  LazyNuxtPwaAssets: LazyComponent<typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/NuxtPwaAssets")['default']>
  LazyPwaAppleImage: LazyComponent<typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaAppleImage.vue")['default']>
  LazyPwaAppleSplashScreenImage: LazyComponent<typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaAppleSplashScreenImage.vue")['default']>
  LazyPwaFaviconImage: LazyComponent<typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaFaviconImage.vue")['default']>
  LazyPwaMaskableImage: LazyComponent<typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaMaskableImage.vue")['default']>
  LazyPwaTransparentImage: LazyComponent<typeof import("../../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaTransparentImage.vue")['default']>
  LazyNuxtPage: LazyComponent<typeof import("../../node_modules/nuxt/dist/pages/runtime/page")['default']>
  LazyNoScript: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['NoScript']>
  LazyLink: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Link']>
  LazyBase: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Base']>
  LazyTitle: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Title']>
  LazyMeta: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Meta']>
  LazyStyle: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Style']>
  LazyHead: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Head']>
  LazyHtml: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Html']>
  LazyBody: LazyComponent<typeof import("../../node_modules/nuxt/dist/head/runtime/components")['Body']>
  LazyNuxtIsland: LazyComponent<typeof import("../../node_modules/nuxt/dist/app/components/nuxt-island")['default']>
}

declare module 'vue' {
  export interface GlobalComponents extends _GlobalComponents { }
}

export {}
