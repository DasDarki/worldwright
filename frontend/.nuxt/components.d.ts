
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


export const AppHeader: typeof import("../components/AppHeader.vue")['default']
export const AppSidebar: typeof import("../components/AppSidebar.vue")['default']
export const AssetUpload: typeof import("../components/AssetUpload.vue")['default']
export const BodyNode: typeof import("../components/BodyNode.vue")['default']
export const BodyView: typeof import("../components/BodyView.vue")['default']
export const DateInput: typeof import("../components/DateInput.vue")['default']
export const EntityForm: typeof import("../components/EntityForm.vue")['default']
export const EntityPicker: typeof import("../components/EntityPicker.vue")['default']
export const EntityTreeNode: typeof import("../components/EntityTreeNode.vue")['default']
export const EventForm: typeof import("../components/EventForm.vue")['default']
export const EventsForEntityPanel: typeof import("../components/EventsForEntityPanel.vue")['default']
export const LangToggle: typeof import("../components/LangToggle.vue")['default']
export const MapCanvas: typeof import("../components/MapCanvas.vue")['default']
export const MobileDrawer: typeof import("../components/MobileDrawer.vue")['default']
export const MoonCycleTimeline: typeof import("../components/MoonCycleTimeline.vue")['default']
export const MoonDisplay: typeof import("../components/MoonDisplay.vue")['default']
export const MoonGlyph: typeof import("../components/MoonGlyph.vue")['default']
export const Ornament: typeof import("../components/Ornament.vue")['default']
export const PWAStatus: typeof import("../components/PWAStatus.client.vue")['default']
export const PaperOverlays: typeof import("../components/PaperOverlays.vue")['default']
export const PedigreeChart: typeof import("../components/PedigreeChart.vue")['default']
export const RelationshipEditor: typeof import("../components/RelationshipEditor.vue")['default']
export const RelationshipPanel: typeof import("../components/RelationshipPanel.vue")['default']
export const SearchBox: typeof import("../components/SearchBox.vue")['default']
export const ShareButton: typeof import("../components/ShareButton.vue")['default']
export const ShortcutHelp: typeof import("../components/ShortcutHelp.vue")['default']
export const SpoilerBlock: typeof import("../components/SpoilerBlock.vue")['default']
export const StarSparkle: typeof import("../components/StarSparkle.vue")['default']
export const TagInput: typeof import("../components/TagInput.vue")['default']
export const ThemeToggle: typeof import("../components/ThemeToggle.vue")['default']
export const TimelineCurator: typeof import("../components/TimelineCurator.vue")['default']
export const TimelineForm: typeof import("../components/TimelineForm.vue")['default']
export const ToastContainer: typeof import("../components/ToastContainer.vue")['default']
export const WikilinkPreview: typeof import("../components/WikilinkPreview.vue")['default']
export const Editor: typeof import("../components/editor/Editor.client.vue")['default']
export const EditorToolbar: typeof import("../components/editor/EditorToolbar.vue")['default']
export const ImageExtension: typeof import("../components/editor/ImageExtension")['default']
export const ImagePicker: typeof import("../components/editor/ImagePicker.client.vue")['default']
export const SpoilerExtension: typeof import("../components/editor/SpoilerExtension")['default']
export const WikilinkExtension: typeof import("../components/editor/WikilinkExtension")['default']
export const WikilinkPicker: typeof import("../components/editor/WikilinkPicker.client.vue")['default']
export const NuxtWelcome: typeof import("../node_modules/nuxt/dist/app/components/welcome.vue")['default']
export const NuxtLayout: typeof import("../node_modules/nuxt/dist/app/components/nuxt-layout")['default']
export const NuxtErrorBoundary: typeof import("../node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']
export const ClientOnly: typeof import("../node_modules/nuxt/dist/app/components/client-only")['default']
export const DevOnly: typeof import("../node_modules/nuxt/dist/app/components/dev-only")['default']
export const ServerPlaceholder: typeof import("../node_modules/nuxt/dist/app/components/server-placeholder")['default']
export const NuxtLink: typeof import("../node_modules/nuxt/dist/app/components/nuxt-link")['default']
export const NuxtLoadingIndicator: typeof import("../node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']
export const NuxtTime: typeof import("../node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']
export const NuxtRouteAnnouncer: typeof import("../node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']
export const NuxtImg: typeof import("../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtImg']
export const NuxtPicture: typeof import("../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtPicture']
export const NuxtLinkLocale: typeof import("../node_modules/@nuxtjs/i18n/dist/runtime/components/NuxtLinkLocale")['default']
export const SwitchLocalePathLink: typeof import("../node_modules/@nuxtjs/i18n/dist/runtime/components/SwitchLocalePathLink")['default']
export const ColorScheme: typeof import("../node_modules/@nuxtjs/color-mode/dist/runtime/component.vue3.vue")['default']
export const VitePwaManifest: typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/VitePwaManifest")['default']
export const NuxtPwaManifest: typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/VitePwaManifest")['default']
export const NuxtPwaAssets: typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/NuxtPwaAssets")['default']
export const PwaAppleImage: typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaAppleImage.vue")['default']
export const PwaAppleSplashScreenImage: typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaAppleSplashScreenImage.vue")['default']
export const PwaFaviconImage: typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaFaviconImage.vue")['default']
export const PwaMaskableImage: typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaMaskableImage.vue")['default']
export const PwaTransparentImage: typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaTransparentImage.vue")['default']
export const NuxtPage: typeof import("../node_modules/nuxt/dist/pages/runtime/page")['default']
export const NoScript: typeof import("../node_modules/nuxt/dist/head/runtime/components")['NoScript']
export const Link: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Link']
export const Base: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Base']
export const Title: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Title']
export const Meta: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Meta']
export const Style: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Style']
export const Head: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Head']
export const Html: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Html']
export const Body: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Body']
export const NuxtIsland: typeof import("../node_modules/nuxt/dist/app/components/nuxt-island")['default']
export const LazyAppHeader: LazyComponent<typeof import("../components/AppHeader.vue")['default']>
export const LazyAppSidebar: LazyComponent<typeof import("../components/AppSidebar.vue")['default']>
export const LazyAssetUpload: LazyComponent<typeof import("../components/AssetUpload.vue")['default']>
export const LazyBodyNode: LazyComponent<typeof import("../components/BodyNode.vue")['default']>
export const LazyBodyView: LazyComponent<typeof import("../components/BodyView.vue")['default']>
export const LazyDateInput: LazyComponent<typeof import("../components/DateInput.vue")['default']>
export const LazyEntityForm: LazyComponent<typeof import("../components/EntityForm.vue")['default']>
export const LazyEntityPicker: LazyComponent<typeof import("../components/EntityPicker.vue")['default']>
export const LazyEntityTreeNode: LazyComponent<typeof import("../components/EntityTreeNode.vue")['default']>
export const LazyEventForm: LazyComponent<typeof import("../components/EventForm.vue")['default']>
export const LazyEventsForEntityPanel: LazyComponent<typeof import("../components/EventsForEntityPanel.vue")['default']>
export const LazyLangToggle: LazyComponent<typeof import("../components/LangToggle.vue")['default']>
export const LazyMapCanvas: LazyComponent<typeof import("../components/MapCanvas.vue")['default']>
export const LazyMobileDrawer: LazyComponent<typeof import("../components/MobileDrawer.vue")['default']>
export const LazyMoonCycleTimeline: LazyComponent<typeof import("../components/MoonCycleTimeline.vue")['default']>
export const LazyMoonDisplay: LazyComponent<typeof import("../components/MoonDisplay.vue")['default']>
export const LazyMoonGlyph: LazyComponent<typeof import("../components/MoonGlyph.vue")['default']>
export const LazyOrnament: LazyComponent<typeof import("../components/Ornament.vue")['default']>
export const LazyPWAStatus: LazyComponent<typeof import("../components/PWAStatus.client.vue")['default']>
export const LazyPaperOverlays: LazyComponent<typeof import("../components/PaperOverlays.vue")['default']>
export const LazyPedigreeChart: LazyComponent<typeof import("../components/PedigreeChart.vue")['default']>
export const LazyRelationshipEditor: LazyComponent<typeof import("../components/RelationshipEditor.vue")['default']>
export const LazyRelationshipPanel: LazyComponent<typeof import("../components/RelationshipPanel.vue")['default']>
export const LazySearchBox: LazyComponent<typeof import("../components/SearchBox.vue")['default']>
export const LazyShareButton: LazyComponent<typeof import("../components/ShareButton.vue")['default']>
export const LazyShortcutHelp: LazyComponent<typeof import("../components/ShortcutHelp.vue")['default']>
export const LazySpoilerBlock: LazyComponent<typeof import("../components/SpoilerBlock.vue")['default']>
export const LazyStarSparkle: LazyComponent<typeof import("../components/StarSparkle.vue")['default']>
export const LazyTagInput: LazyComponent<typeof import("../components/TagInput.vue")['default']>
export const LazyThemeToggle: LazyComponent<typeof import("../components/ThemeToggle.vue")['default']>
export const LazyTimelineCurator: LazyComponent<typeof import("../components/TimelineCurator.vue")['default']>
export const LazyTimelineForm: LazyComponent<typeof import("../components/TimelineForm.vue")['default']>
export const LazyToastContainer: LazyComponent<typeof import("../components/ToastContainer.vue")['default']>
export const LazyWikilinkPreview: LazyComponent<typeof import("../components/WikilinkPreview.vue")['default']>
export const LazyEditor: LazyComponent<typeof import("../components/editor/Editor.client.vue")['default']>
export const LazyEditorToolbar: LazyComponent<typeof import("../components/editor/EditorToolbar.vue")['default']>
export const LazyImageExtension: LazyComponent<typeof import("../components/editor/ImageExtension")['default']>
export const LazyImagePicker: LazyComponent<typeof import("../components/editor/ImagePicker.client.vue")['default']>
export const LazySpoilerExtension: LazyComponent<typeof import("../components/editor/SpoilerExtension")['default']>
export const LazyWikilinkExtension: LazyComponent<typeof import("../components/editor/WikilinkExtension")['default']>
export const LazyWikilinkPicker: LazyComponent<typeof import("../components/editor/WikilinkPicker.client.vue")['default']>
export const LazyNuxtWelcome: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/welcome.vue")['default']>
export const LazyNuxtLayout: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-layout")['default']>
export const LazyNuxtErrorBoundary: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']>
export const LazyClientOnly: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/client-only")['default']>
export const LazyDevOnly: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/dev-only")['default']>
export const LazyServerPlaceholder: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/server-placeholder")['default']>
export const LazyNuxtLink: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-link")['default']>
export const LazyNuxtLoadingIndicator: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']>
export const LazyNuxtTime: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']>
export const LazyNuxtRouteAnnouncer: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']>
export const LazyNuxtImg: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtImg']>
export const LazyNuxtPicture: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtPicture']>
export const LazyNuxtLinkLocale: LazyComponent<typeof import("../node_modules/@nuxtjs/i18n/dist/runtime/components/NuxtLinkLocale")['default']>
export const LazySwitchLocalePathLink: LazyComponent<typeof import("../node_modules/@nuxtjs/i18n/dist/runtime/components/SwitchLocalePathLink")['default']>
export const LazyColorScheme: LazyComponent<typeof import("../node_modules/@nuxtjs/color-mode/dist/runtime/component.vue3.vue")['default']>
export const LazyVitePwaManifest: LazyComponent<typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/VitePwaManifest")['default']>
export const LazyNuxtPwaManifest: LazyComponent<typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/VitePwaManifest")['default']>
export const LazyNuxtPwaAssets: LazyComponent<typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/NuxtPwaAssets")['default']>
export const LazyPwaAppleImage: LazyComponent<typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaAppleImage.vue")['default']>
export const LazyPwaAppleSplashScreenImage: LazyComponent<typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaAppleSplashScreenImage.vue")['default']>
export const LazyPwaFaviconImage: LazyComponent<typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaFaviconImage.vue")['default']>
export const LazyPwaMaskableImage: LazyComponent<typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaMaskableImage.vue")['default']>
export const LazyPwaTransparentImage: LazyComponent<typeof import("../node_modules/@vite-pwa/nuxt/dist/runtime/components/PwaTransparentImage.vue")['default']>
export const LazyNuxtPage: LazyComponent<typeof import("../node_modules/nuxt/dist/pages/runtime/page")['default']>
export const LazyNoScript: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['NoScript']>
export const LazyLink: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Link']>
export const LazyBase: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Base']>
export const LazyTitle: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Title']>
export const LazyMeta: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Meta']>
export const LazyStyle: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Style']>
export const LazyHead: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Head']>
export const LazyHtml: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Html']>
export const LazyBody: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Body']>
export const LazyNuxtIsland: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-island")['default']>

export const componentNames: string[]
