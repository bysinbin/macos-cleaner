/**
 * DiskCleaner Pro — Frontend Controller
 */

// Application State
const state = {
  currentTab: 'overview',
  diskStats: null,
  scanData: null,
  nodeModulesData: [],
  largeFilesData: [],
  treeData: null,
  treeCurrentPath: '~',
  duplicatesData: [],
  selectedNodeModules: new Set(),
  selectedLargeFiles: new Set(),
  selectedDuplicates: new Set(),
  selectedLeftovers: new Set(),
  leftoversData: null,
  appsData: null,
  appsFilter: 'all',
  appsSearch: '',
  browsersData: null,
  selectedBrowsers: new Set(),
  downloadsData: null,
  selectedDownloads: new Set(),
  downloadsFilter: 'all',
  appleData: null,
  mediaData: null,
  selectedMedia: new Set(),
  smartCareData: null,
  maintenanceTasks: null,
  monitorData: null,
  authStatus: null,
  startupData: null,
  startupFilter: 'all',
  startupSearch: '',
  isScanning: false,
  isCleaning: false,
};

// DOM Elements
const elements = {
  // Navigation
  navButtons: document.querySelectorAll('.nav-item'),
  tabPanes: document.querySelectorAll('.tab-pane'),
  catCards: document.querySelectorAll('.cat-card[data-goto]'),
  pageHeading: document.getElementById('page-heading'),
  pageSubheading: document.getElementById('page-subheading'),

  // Controls
  toggleDryRun: document.getElementById('toggle-dry-run'),
  toggleTrashMode: document.getElementById('toggle-trash-mode'),
  btnGlobalScan: document.getElementById('btn-global-scan'),
  scanBanner: document.getElementById('scan-banner'),
  scanBannerTitle: document.getElementById('scan-banner-title'),
  scanBannerDesc: document.getElementById('scan-banner-desc'),
  scanBannerItems: document.getElementById('scan-banner-items'),
  scanBannerFound: document.getElementById('scan-banner-found'),

  // Overview
  donutDiskFill: document.getElementById('donut-disk-fill'),
  diskUsedPct: document.getElementById('disk-used-pct'),
  metricUsedSpace: document.getElementById('metric-used-space'),
  metricFreeSpace: document.getElementById('metric-free-space'),
  metricTotalSpace: document.getElementById('metric-total-space'),
  miniFreeSpace: document.getElementById('mini-free-space'),
  miniDiskFill: document.getElementById('mini-disk-fill'),
  miniDiskDetail: document.getElementById('mini-disk-detail'),
  quickCleanSize: document.getElementById('quick-clean-size'),
  btnQuickClean: document.getElementById('btn-quick-clean'),
  overviewDevSize: document.getElementById('overview-dev-size'),
  overviewSysSize: document.getElementById('overview-sys-size'),
  overviewNodeSize: document.getElementById('overview-node-size'),
  badgeDevCount: document.getElementById('badge-dev-count'),
  badgeSysCount: document.getElementById('badge-sys-count'),

  // Storage Threshold Alert
  storageThresholdBanner: document.getElementById('storage-threshold-banner'),
  alertBannerTitle: document.getElementById('alert-banner-title'),
  alertBannerDesc: document.getElementById('alert-banner-desc'),
  btnAlertSmartClean: document.getElementById('btn-alert-smart-clean'),

  // Developer tab
  devTargetsList: document.getElementById('dev-targets-list'),
  btnSelectAllDev: document.getElementById('btn-select-all-dev'),
  btnCleanSelectedDev: document.getElementById('btn-clean-selected-dev'),
  devSelectedSize: document.getElementById('dev-selected-size'),

  // System tab
  sysTargetsList: document.getElementById('sys-targets-list'),
  dynamicCachesList: document.getElementById('dynamic-caches-list'),
  btnSelectAllSys: document.getElementById('btn-select-all-sys'),
  btnCleanSelectedSys: document.getElementById('btn-clean-selected-sys'),
  sysSelectedSize: document.getElementById('sys-selected-size'),

  // Node Modules tab
  inputNodePath: document.getElementById('input-node-path'),
  btnApplyNodePath: document.getElementById('btn-apply-node-path'),
  btnScanNodeModules: document.getElementById('btn-scan-node-modules'),
  btnCleanSelectedNode: document.getElementById('btn-clean-selected-node'),
  nodeSelectedSize: document.getElementById('node-selected-size'),
  tbodyNodeModules: document.getElementById('tbody-node-modules'),
  thCheckAllNode: document.getElementById('th-check-all-node'),

  // Large Files tab
  selectMinSize: document.getElementById('select-min-size'),
  selectCategory: document.getElementById('select-category'),
  selectAge: document.getElementById('select-age'),
  inputLargefilePath: document.getElementById('input-largefile-path'),
  btnScanLargeFiles: document.getElementById('btn-scan-large-files'),
  btnCleanSelectedFiles: document.getElementById('btn-clean-selected-files'),
  filesSelectedSize: document.getElementById('files-selected-size'),
  tbodyLargeFiles: document.getElementById('tbody-large-files'),
  thCheckAllFiles: document.getElementById('th-check-all-files'),

  // Directory Tree tab
  btnTreeUp: document.getElementById('btn-tree-up'),
  btnTreeRefresh: document.getElementById('btn-tree-refresh'),
  inputTreePath: document.getElementById('input-tree-path'),
  btnTreeGo: document.getElementById('btn-tree-go'),
  treeTotalSize: document.getElementById('tree-total-size'),
  treeItemsList: document.getElementById('tree-items-list'),
  treePills: document.querySelectorAll('.btn-pill'),

  // Duplicate Finder tab
  selectDupMinSize: document.getElementById('select-dup-min-size'),
  inputDupPath: document.getElementById('input-dup-path'),
  btnScanDuplicates: document.getElementById('btn-scan-duplicates'),
  btnDupSelectCopies: document.getElementById('btn-dup-select-copies'),
  btnDupCleanSelected: document.getElementById('btn-dup-clean-selected'),
  dupSelectedSize: document.getElementById('dup-selected-size'),
  duplicateGroupsList: document.getElementById('duplicate-groups-list'),

  // Leftovers tab
  badgeLeftoversCount: document.getElementById('badge-leftovers-count'),
  leftoversAppsCount: document.getElementById('leftovers-apps-count'),
  leftoversDsCount: document.getElementById('leftovers-ds-count'),
  leftoversTotalSize: document.getElementById('leftovers-total-size'),
  leftoversItemsList: document.getElementById('leftovers-items-list'),
  btnScanLeftovers: document.getElementById('btn-scan-leftovers'),
  btnSelectAllLeftovers: document.getElementById('btn-select-all-leftovers'),
  btnCleanSelectedLeftovers: document.getElementById('btn-clean-selected-leftovers'),
  leftoversSelectedSize: document.getElementById('leftovers-selected-size'),

  // Apps tab
  btnScanApps: document.getElementById('btn-scan-apps'),
  appsSearchInput: document.getElementById('apps-search-input'),
  appsFilterChips: document.getElementById('apps-filter-chips'),
  chipAppsAllCount: document.getElementById('chip-apps-all-count'),
  appsCountBadge: document.getElementById('apps-count-badge'),
  appsBinarySizeBadge: document.getElementById('apps-binary-size-badge'),
  appsDataSizeBadge: document.getElementById('apps-data-size-badge'),
  appsItemsList: document.getElementById('apps-items-list'),

  // Browsers tab
  btnScanBrowsers: document.getElementById('btn-scan-browsers'),
  btnCleanAllBrowsers: document.getElementById('btn-clean-all-browsers'),
  browsersTotalCleanSize: document.getElementById('browsers-total-clean-size'),
  browsersGridList: document.getElementById('browsers-grid-list'),

  // Downloads tab
  btnScanDownloads: document.getElementById('btn-scan-downloads'),
  btnCleanDownloads: document.getElementById('btn-clean-downloads'),
  downloadsSelectedSize: document.getElementById('downloads-selected-size'),
  downloadsInstallersSize: document.getElementById('downloads-installers-size'),
  downloadsArchivesSize: document.getElementById('downloads-archives-size'),
  downloadsOldSize: document.getElementById('downloads-old-size'),
  downloadsFilterChips: document.getElementById('downloads-filter-chips'),
  btnSelectAllDownloads: document.getElementById('btn-select-all-downloads'),
  downloadsItemsList: document.getElementById('downloads-items-list'),

  // Apple & System Data tab
  btnScanApple: document.getElementById('btn-scan-apple'),
  btnReclaimPurgeable: document.getElementById('btn-reclaim-purgeable'),
  btnCleanSafeSystemData: document.getElementById('btn-clean-safe-systemdata'),
  macosVersionText: document.getElementById('macos-version-text'),
  macosVolumeSize: document.getElementById('macos-volume-size'),
  macosSubvolumesList: document.getElementById('macos-subvolumes-list'),
  sysdataTotalSize: document.getElementById('sysdata-total-size'),
  sysdataSafeCleanable: document.getElementById('sysdata-safe-cleanable'),
  vmSleepimagePill: document.getElementById('vm-sleepimage-pill'),
  appleCategoriesContainer: document.getElementById('apple-categories-container'),
  appleSnapshotsCount: document.getElementById('apple-snapshots-count'),
  appleSnapshotsList: document.getElementById('apple-snapshots-list'),
  btnCleanAllSnapshots: document.getElementById('btn-clean-all-snapshots'),
  appleBackupsSize: document.getElementById('apple-backups-size'),
  appleBackupsList: document.getElementById('apple-backups-list'),
  appleSimulatorsSize: document.getElementById('apple-simulators-size'),
  appleSimulatorsList: document.getElementById('apple-simulators-list'),
  btnCleanSimulators: document.getElementById('btn-clean-simulators'),

  // Media tab
  btnScanMedia: document.getElementById('btn-scan-media'),
  btnCleanMedia: document.getElementById('btn-clean-media'),
  mediaSelectedSize: document.getElementById('media-selected-size'),
  mediaVideoSize: document.getElementById('media-video-size'),
  mediaImageSize: document.getElementById('media-image-size'),
  mediaMailSize: document.getElementById('media-mail-size'),
  mediaItemsList: document.getElementById('media-items-list'),

  // Modal
  confirmModal: document.getElementById('confirm-modal'),
  modalTitle: document.getElementById('modal-title'),
  modalMessage: document.getElementById('modal-message'),
  modalFreedSize: document.getElementById('modal-freed-size'),
  modalMethodText: document.getElementById('modal-method-text'),
  btnModalCancel: document.getElementById('btn-modal-cancel'),
  btnModalConfirm: document.getElementById('btn-modal-confirm'),

  // Sysext Guidance Modal
  sysextManageDialog: document.getElementById('sysext-manage-dialog'),
  sysextModalTitle: document.getElementById('sysext-modal-title'),
  sysextModalName: document.getElementById('sysext-modal-name'),
  sysextModalType: document.getElementById('sysext-modal-type'),
  sysextModalHostAppRow: document.getElementById('sysext-modal-hostapp-row'),
  sysextModalHostApp: document.getElementById('sysext-modal-hostapp'),
  sysextModalState: document.getElementById('sysext-modal-state'),
  btnCloseSysextModal: document.getElementById('btn-close-sysext-modal'),
  btnRevealSysextApp: document.getElementById('btn-reveal-sysext-app'),
  btnOpenSysSettings: document.getElementById('btn-open-sys-settings'),

  // Auth
  authModal: document.getElementById('auth-modal'),
  authForm: document.getElementById('auth-form'),
  authPasswordInput: document.getElementById('auth-password-input'),
  authErrorMsg: document.getElementById('auth-error-msg'),
  btnAuthSubmit: document.getElementById('btn-auth-submit'),
  authStatusBar: document.getElementById('auth-status-bar'),
  btnLogout: document.getElementById('btn-logout'),

  // Maintenance
  btnRunAllMaintenance: document.getElementById('btn-run-all-maintenance'),
  maintenanceTasksContainer: document.getElementById('maintenance-tasks-container'),
  maintenanceLogBox: document.getElementById('maintenance-log-box'),
  maintenanceLogContent: document.getElementById('maintenance-log-content'),
  btnClearMaintenanceLog: document.getElementById('btn-clear-maintenance-log'),

  // Shredder
  inputShredPath: document.getElementById('input-shred-path'),
  selectShredPasses: document.getElementById('select-shred-passes'),
  btnExecuteShred: document.getElementById('btn-execute-shred'),
  shredResultBox: document.getElementById('shred-result-box'),

  // Hardware Monitor
  btnRefreshMonitor: document.getElementById('btn-refresh-monitor'),
  monCpuModel: document.getElementById('mon-cpu-model'),
  monCpuPct: document.getElementById('mon-cpu-pct'),
  monCpuCores: document.getElementById('mon-cpu-cores'),
  monRamTotal: document.getElementById('mon-ram-total'),
  monRamPct: document.getElementById('mon-ram-pct'),
  monRamUsed: document.getElementById('mon-ram-used'),
  monRamFree: document.getElementById('mon-ram-free'),
  monRamWired: document.getElementById('mon-ram-wired'),
  monRamCompressed: document.getElementById('mon-ram-compressed'),
  monHostname: document.getElementById('mon-hostname'),
  monOsVer: document.getElementById('mon-os-ver'),
  monUptime: document.getElementById('mon-uptime'),

  // Battery Monitor
  monBatteryCard: document.getElementById('mon-battery-card'),
  monBattSource: document.getElementById('mon-batt-source'),
  monBattPct: document.getElementById('mon-batt-pct'),
  monBattChargingStatus: document.getElementById('mon-batt-charging-status'),
  monBattHealth: document.getElementById('mon-batt-health'),
  monBattCycles: document.getElementById('mon-batt-cycles'),
  monBattTemp: document.getElementById('mon-batt-temp'),
  monBattCondition: document.getElementById('mon-batt-condition'),

  // GPU Monitor
  monGpuModel: document.getElementById('mon-gpu-model'),
  monGpuCores: document.getElementById('mon-gpu-cores'),
  monGpuMetal: document.getElementById('mon-gpu-metal'),
  monGpuRes: document.getElementById('mon-gpu-res'),
  monGpuVendor: document.getElementById('mon-gpu-vendor'),
  monCpuTemp: document.getElementById('mon-cpu-temp'),
  monGpuTemp: document.getElementById('mon-gpu-temp'),

  // Thermal & Sensors
  monThermalStrip: document.getElementById('mon-thermal-strip'),
  monThermalCpu: document.getElementById('mon-thermal-cpu'),
  monThermalGpu: document.getElementById('mon-thermal-gpu'),
  monThermalBatt: document.getElementById('mon-thermal-batt'),
  monThermalState: document.getElementById('mon-thermal-state'),
  monThermalFan: document.getElementById('mon-thermal-fan'),

  // Disk Monitor
  monDiskType: document.getElementById('mon-disk-type'),
  monDiskPct: document.getElementById('mon-disk-pct'),
  monDiskUsed: document.getElementById('mon-disk-used'),
  monDiskFree: document.getElementById('mon-disk-free'),
  monDiskSmart: document.getElementById('mon-disk-smart'),
  monDiskThroughput: document.getElementById('mon-disk-throughput'),

  // Extensions
  badgeExtensionsCount: document.getElementById('badge-extensions-count'),
  btnRefreshExtensions: document.getElementById('btn-refresh-extensions'),
  extStatTotal: document.getElementById('ext-stat-total'),
  extStatActive: document.getElementById('ext-stat-active'),
  extStatSysext: document.getElementById('ext-stat-sysext'),
  extStatAppex: document.getElementById('ext-stat-appex'),
  extCategoryFilters: document.getElementById('ext-category-filters'),
  inputSearchExtensions: document.getElementById('input-search-extensions'),
  extensionsItemsList: document.getElementById('extensions-items-list'),
  pillExtAll: document.getElementById('pill-ext-all'),
  pillExtSysext: document.getElementById('pill-ext-sysext'),
  pillExtPluginkit: document.getElementById('pill-ext-pluginkit'),
  pillExtQuicklook: document.getElementById('pill-ext-quicklook'),
  pillExtSpotlight: document.getElementById('pill-ext-spotlight'),
  pillExtPrefpane: document.getElementById('pill-ext-prefpane'),

  // Network Monitor
  btnRefreshNetwork: document.getElementById('btn-refresh-network'),
  netStatStatus: document.getElementById('net-stat-status'),
  netStatIface: document.getElementById('net-stat-iface'),
  netStatIp: document.getElementById('net-stat-ip'),
  netStatGw: document.getElementById('net-stat-gw'),
  netStatSpeedIn: document.getElementById('net-stat-speed-in'),
  netStatTotalIn: document.getElementById('net-stat-total-in'),
  netStatSpeedOut: document.getElementById('net-stat-speed-out'),
  netStatTotalOut: document.getElementById('net-stat-total-out'),
  netDetailMac: document.getElementById('net-detail-mac'),
  netDetailDns: document.getElementById('net-detail-dns'),
  netDetailConnCount: document.getElementById('net-detail-conn-count'),
  netConnBadge: document.getElementById('net-conn-badge'),
  netConnectionsTbody: document.getElementById('net-connections-tbody'),
  netLiveDown: document.getElementById('net-live-down'),
  netLiveUp: document.getElementById('net-live-up'),
  netLivePeak: document.getElementById('net-live-peak'),
  netTrafficCanvas: document.getElementById('net-traffic-canvas'),

  // Privacy Protection
  badgePrivacyCount: document.getElementById('badge-privacy-count'),
  btnRefreshPrivacy: document.getElementById('btn-refresh-privacy'),
  btnCleanPrivacy: document.getElementById('btn-clean-privacy'),
  privStatTotal: document.getElementById('priv-stat-total'),
  privStatSize: document.getElementById('priv-stat-size'),
  privStatHistory: document.getElementById('priv-stat-history'),
  privStatBrowser: document.getElementById('priv-stat-browser'),
  checkAllPrivacy: document.getElementById('check-all-privacy'),
  privacyItemsList: document.getElementById('privacy-items-list'),
  privacyPermissionsList: document.getElementById('privacy-permissions-list'),
  privacyScoreBadge: document.getElementById('privacy-score-badge'),
  privacyStatusText: document.getElementById('privacy-status-text'),

  // Startup Manager
  badgeStartupCount: document.getElementById('badge-startup-count'),
  btnRefreshStartup: document.getElementById('btn-refresh-startup'),
  btnAddLoginItemModal: document.getElementById('btn-add-login-item-modal'),
  startupStatTotal: document.getElementById('startup-stat-total'),
  startupStatActive: document.getElementById('startup-stat-active'),
  startupStatRunning: document.getElementById('startup-stat-running'),
  startupStatLogin: document.getElementById('startup-stat-login'),
  startupFilterChips: document.getElementById('startup-filter-chips'),
  inputStartupSearch: document.getElementById('input-startup-search'),
  startupItemsList: document.getElementById('startup-items-list'),
  addLoginItemDialog: document.getElementById('add-login-item-dialog'),
  formAddLoginItem: document.getElementById('form-add-login-item'),
  inputNewLoginPath: document.getElementById('input-new-login-path'),
  checkNewLoginHidden: document.getElementById('check-new-login-hidden'),
  btnCancelAddLogin: document.getElementById('btn-cancel-add-login'),
  chipCountAll: document.getElementById('chip-count-all'),
  chipCountLogin: document.getElementById('chip-count-login'),
  chipCountUser: document.getElementById('chip-count-user'),
  chipCountSysAgent: document.getElementById('chip-count-sysagent'),
  chipCountDaemon: document.getElementById('chip-count-daemon'),
  startupDonutChart: document.getElementById('startup-donut-chart'),
  donutCenterTotal: document.getElementById('donut-center-total'),
  donutLegendLogin: document.getElementById('donut-legend-login'),
  donutLegendUser: document.getElementById('donut-legend-user'),
  donutLegendSysAgent: document.getElementById('donut-legend-sysagent'),
  donutLegendDaemon: document.getElementById('donut-legend-daemon'),
  startupRatioText: document.getElementById('startup-ratio-text'),
  startupRatioFill: document.getElementById('startup-ratio-fill'),

  // Toast
  toastContainer: document.getElementById('toast-container'),
};

// Tab configuration
const TAB_CONFIG = {
  overview: { title: 'Sistem Genel Bakışı', sub: 'Depolama analizi ve hızlı temizleme merkezi' },
  developer: { title: 'Geliştirici Önbellekleri', sub: 'Xcode, derleyiciler ve paket yöneticilerinin ürettiği önbellekler' },
  system: { title: 'Sistem ve Uygulama Artıkları', sub: 'Uygulama güncelleme indiricileri, loglar ve dinamik önbellekler' },
  nodemodules: { title: 'node_modules Avcısı', sub: 'Projelerdeki devasa ve atıl node_modules klasörlerini temizleyin' },
  largefiles: { title: 'Büyük ve Eski Dosya Analizi', sub: 'Diskinizde yer kaplayan büyük dosyaları bulun ve yönetin' },
  dirtree: { title: 'Derin Dizin Analizi (Disk Space Tree)', sub: 'Klasörlerin kullanımını hiyerarşik olarak inceleyin ve derinliklerine inin' },
  duplicates: { title: 'Yinelenen Dosya Avcısı (Duplicate Finder)', sub: 'Aynı dosya boyutuna ve SHA-256 hash imzasına sahip kopya dosyalar' },
  leftovers: { title: 'Artık Dosya Avcısı (Residual & Leftovers)', sub: 'Kaldırılmış uygulama kalıntıları, eski oturum durumları ve sistem çöpleri' },
  apps: { title: 'Uygulama Kaldırıcı & Sıfırlayıcı (App Uninstaller)', sub: 'Uygulamaları ve tüm gizli kullanıcı/önbellek artıklarını kalıcı silin veya sıfırlayın' },
  browsers: { title: 'Tarayıcı Önbellekleri & Gizlilik', sub: 'Chrome, Safari, Firefox, Arc ve Edge geçici web/GPU önbelleklerini temizleyin' },
  downloads: { title: 'İndirilenler & Kurulum Dosyaları Analizi', sub: 'Eski kurulum kalıpları (.dmg, .pkg), arşivler ve unutulmuş dosyalar' },
  apple: { title: 'Apple & Sistem Verisi', sub: 'Time Machine APFS anlık görüntüleri, iOS yerel aygıt yedekleri ve Xcode simülatörleri' },
  media: { title: 'İletişim & Mesaj Ekleri', sub: 'iMessage ve Apple Mail ile gelen video, fotoğraf, ses ve belge ekleri' },
  maintenance: { title: 'Sistem Bakımı & Hızlandırma', sub: 'macOS bellek boşaltma, DNS ve sistem önbelleği onarımı' },
  shredder: { title: 'Güvenli Dosya Öğütücü', sub: 'Kurtarılamaz çoklu geçişli (DoD 5220.22-M) kalıcı dosya imhası' },
  monitor: { title: 'Donanım & Kaynak Monitörü', sub: 'İşlemci, Bellek (RAM), Pil, GPU ve Disk durumunun canlı görünümü' },
  startup: { title: 'Başlangıç Öğeleri & Arka Plan Hizmetleri', sub: 'Oturum açma uygulamaları, LaunchAgent ve arka plan servislerini yönetin' },
  extensions: { title: 'macOS Eklenti Yöneticisi (Extensions)', sub: 'Sistem sürücüleri, ağ filtreleri, Finder ve uygulama eklentilerini yönetin' },
  network: { title: 'Ağ & Bağlantı Monitörü (Network)', sub: 'Gerçek zamanlı ağ trafiği, indirme/yükleme hızları ve aktif TCP bağlantıları' },
  privacy: { title: 'macOS Gizlilik Koruması & İz Temizleyici', sub: 'Son açılan dosyalar, Terminal komut geçmişleri, tarayıcı izleri ve TCC izinleri' },
};

let monitorInterval = null;
let networkInterval = null;
let extensionsData = null;
let extCurrentFilter = 'all';
let extSearchQuery = '';
let privacyData = null;

// Global fetch wrapper to handle session authentication and expiration (401)
const _originalFetch = window.fetch;
window.fetch = async (url, options = {}) => {
  const token = localStorage.getItem('dc_token');
  if (token && typeof url === 'string' && url.startsWith('/api/')) {
    options = options || {};
    options.headers = options.headers || {};
    if (options.headers instanceof Headers) {
      if (!options.headers.has('Authorization')) {
        options.headers.set('Authorization', `Bearer ${token}`);
      }
    } else {
      if (!options.headers['Authorization']) {
        options.headers['Authorization'] = `Bearer ${token}`;
      }
    }
  }
  const res = await _originalFetch(url, options);
  if (res.status === 401) {
    const urlStr = url ? url.toString() : '';
    if (!urlStr.includes('/api/auth/')) {
      if (elements.authModal && !elements.authModal.open) {
        elements.authModal.showModal();
      }
      if (elements.authStatusBar) elements.authStatusBar.style.display = 'none';
    }
  }
  return res;
};

// Helper for authenticated EventSource instances
function createEventSource(url) {
  const token = localStorage.getItem('dc_token');
  if (token) {
    const sep = url.includes('?') ? '&' : '?';
    url += `${sep}token=${encodeURIComponent(token)}`;
  }
  return new EventSource(url);
}

// Centralized initial data fetcher (runs only after authentication)
function initializeDashboardData() {
  initScanProgressStream();
  fetchSystemStats();
  startGlobalScan();
  startLeftoversScan(true); // background initial scan for badge
  loadSmartCare(); // Initial smart care health scan
  loadStartupItems(true); // background initial scan for startup badge
  loadExtensions(true); // background initial scan for extensions badge
  loadPrivacyTraces(true); // background initial scan for privacy badge
  checkFDAPermissions(); // Check macOS Full Disk Access status
  loadDockerStatus(); // Check Docker & Container status
}

// Initialize Application
document.addEventListener('DOMContentLoaded', async () => {
  setupNavigation();
  setupEventHandlers();
  initFDA();
  initDirTreeViewSwitcher();

  const isAuth = await checkAuthStatus();
  if (isAuth) {
    initializeDashboardData();
  }
});

// Setup Tab Navigation
function setupNavigation() {
  elements.navButtons.forEach(btn => {
    btn.addEventListener('click', () => {
      const tab = btn.dataset.tab;
      switchTab(tab);
    });
  });

  elements.catCards.forEach(card => {
    card.addEventListener('click', () => {
      const targetTab = card.dataset.goto;
      switchTab(targetTab);
    });
  });
}

function switchTab(tabName) {
  state.currentTab = tabName;

  // Update nav buttons
  elements.navButtons.forEach(b => {
    b.classList.toggle('active', b.dataset.tab === tabName);
  });

  // Update panes
  elements.tabPanes.forEach(pane => {
    pane.classList.toggle('active', pane.id === `pane-${tabName}`);
  });

  // Update header text
  if (TAB_CONFIG[tabName]) {
    elements.pageHeading.textContent = TAB_CONFIG[tabName].title;
    elements.pageSubheading.textContent = TAB_CONFIG[tabName].sub;
  }

  // Auto load directory tree if opened first time
  if (tabName === 'dirtree' && !state.treeData) {
    loadDirTree(state.treeCurrentPath);
  }

  // Auto load leftovers if opened
  if (tabName === 'leftovers' && !state.leftoversData) {
    startLeftoversScan();
  }

  if (tabName === 'developer') {
    loadDockerStatus();
  }

  // Auto load new tabs if opened first time
  if (tabName === 'apps' && !state.appsData) {
    startAppsScan();
  }
  if (tabName === 'browsers' && !state.browsersData) {
    startBrowsersScan();
  }
  if (tabName === 'downloads' && !state.downloadsData) {
    startDownloadsScan();
  }
  if (tabName === 'apple' && !state.appleData) {
    startAppleScan();
  }
  if (tabName === 'media' && !state.mediaData) {
    startMediaScan();
  }
  if (tabName === 'maintenance') {
    loadMaintenanceTasks();
  }
  if (tabName === 'startup' && !state.startupData) {
    loadStartupItems();
  }
  if (tabName === 'extensions') {
    loadExtensions();
  }
  if (tabName === 'network') {
    loadNetworkStats();
    if (networkInterval) clearInterval(networkInterval);
    networkInterval = setInterval(loadNetworkStats, 3000);
  } else if (networkInterval) {
    clearInterval(networkInterval);
    networkInterval = null;
  }
  if (tabName === 'privacy') {
    loadPrivacyTraces();
  }
  if (tabName === 'monitor') {
    loadHardwareMonitor();
    if (monitorInterval) clearInterval(monitorInterval);
    monitorInterval = setInterval(loadHardwareMonitor, 3000);
  } else if (monitorInterval) {
    clearInterval(monitorInterval);
    monitorInterval = null;
  }
}

// Fetch System Storage Statistics
async function fetchSystemStats() {
  try {
    const res = await fetch('/api/system');
    if (!res.ok) throw new Error('Sistem bilgisi alınamadı');
    const data = await res.json();
    state.diskStats = data.stats;
    updateDiskVisuals(data.stats);
  } catch (err) {
    console.error(err);
    showToast('Depolama bilgisi alınamadı: ' + err.message, 'error');
  }
}

// Update Donut Chart and Metrics
function updateDiskVisuals(stats) {
  if (!stats) return;

  const pct = Math.round(stats.usedPercent);
  elements.diskUsedPct.textContent = `${pct}%`;
  elements.metricUsedSpace.textContent = stats.usedStr;
  elements.metricFreeSpace.textContent = stats.freeStr;
  elements.metricTotalSpace.textContent = stats.totalStr;

  // Mini sidebar info
  elements.miniFreeSpace.textContent = stats.freeStr;
  elements.miniDiskFill.style.width = `${pct}%`;
  elements.miniDiskDetail.textContent = `${stats.usedStr} / ${stats.totalStr} (%${pct})`;

  // SVG Donut calculation (circumference = 2 * PI * 80 ≈ 502.65)
  const circumference = 502.65;
  const offset = circumference - (pct / 100) * circumference;
  elements.donutDiskFill.style.strokeDashoffset = offset;

  // Warning color if disk is critically full (> 85%)
  if (pct >= 90) {
    elements.donutDiskFill.style.stroke = 'var(--accent-rose)';
    elements.miniDiskFill.style.background = 'var(--gradient-danger)';
  } else if (pct >= 80) {
    elements.donutDiskFill.style.stroke = 'var(--accent-amber)';
  } else {
    elements.donutDiskFill.style.stroke = 'var(--accent-indigo)';
    elements.miniDiskFill.style.background = 'var(--gradient-brand)';
  }

  // Threshold alert banner
  if (elements.storageThresholdBanner) {
    if (pct >= 85) {
      elements.storageThresholdBanner.style.display = 'flex';
      if (elements.alertBannerTitle) {
        elements.alertBannerTitle.textContent = `⚠️ Kritik Disk Doluluğu: %${pct} Kullanılıyor!`;
      }
    } else {
      elements.storageThresholdBanner.style.display = 'none';
    }
  }
}

// Start Global Target Scan via SSE
function startGlobalScan() {
  if (state.isScanning) return;
  state.isScanning = true;

  elements.btnGlobalScan.disabled = true;
  elements.scanBanner.style.display = 'flex';
  elements.scanBannerTitle.textContent = 'Disk Taranıyor...';
  elements.scanBannerDesc.textContent = 'Geliştirici ve sistem dizinleri taranıyor';

  const eventSource = createEventSource('/api/scan?stream=true');

  eventSource.addEventListener('progress', (e) => {
    try {
      const data = JSON.parse(e.data);
      elements.scanBannerDesc.textContent = `Taranan: ${data.currentPath || 'Dosyalar'}`;
      elements.scanBannerItems.textContent = `${data.scannedFiles || 0} dosya incelendi`;
      if (data.foundBytes) {
        elements.scanBannerFound.textContent = `${formatBytes(data.foundBytes)} bulundu`;
      }
    } catch (err) {
      console.error(err);
    }
  });

  eventSource.addEventListener('result', (e) => {
    eventSource.close();
    state.isScanning = false;
    elements.btnGlobalScan.disabled = false;
    elements.scanBanner.style.display = 'none';

    try {
      const result = JSON.parse(e.data);
      state.scanData = result;
      renderScanResults(result);
      if (result.diskStats) {
        updateDiskVisuals(result.diskStats);
      }
      showToast(`Tarama tamamlandı: ${result.totalCleanableStr} temizlenebilir alan tespit edildi!`, 'success');
    } catch (err) {
      console.error(err);
    }
  });

  eventSource.addEventListener('error', (e) => {
    eventSource.close();
    state.isScanning = false;
    elements.btnGlobalScan.disabled = false;
    elements.scanBanner.style.display = 'none';
    console.error('Scan error', e);
    showToast('Tarama sırasında bir hata oluştu veya bağlantı koptu.', 'error');
  });
}

// Render Discovered Cleanable Targets
function renderScanResults(data) {
  if (!data || !data.targets) return;

  const devTargets = data.targets.filter(t => t.category === 'developer' && t.exists);
  const sysTargets = data.targets.filter(t => t.category === 'system' && t.exists);

  let devTotalSize = 0;
  let sysTotalSize = 0;

  devTargets.forEach(t => { if (t.selected) devTotalSize += t.size; });
  sysTargets.forEach(t => { if (t.selected) sysTotalSize += t.size; });

  // Update Overview badges
  elements.quickCleanSize.textContent = data.totalCleanableStr;
  elements.overviewDevSize.textContent = formatBytes(devTotalSize);
  elements.overviewSysSize.textContent = formatBytes(sysTotalSize);
  elements.badgeDevCount.textContent = devTargets.length;
  elements.badgeSysCount.textContent = sysTargets.length;

  // Render Developer list
  renderTargetList(elements.devTargetsList, devTargets, 'dev');
  // Render System list
  renderTargetList(elements.sysTargetsList, sysTargets, 'sys');
  // Render Dynamic Caches list
  renderDynamicCachesList(elements.dynamicCachesList, data.dynamicCaches || []);

  updateSelectedSize('dev');
  updateSelectedSize('sys');
}

function renderTargetList(container, targets, prefix) {
  if (!targets || targets.length === 0) {
    container.innerHTML = '<div class="loading-state">Bu kategoride temizlenecek dosya bulunamadı.</div>';
    return;
  }

  container.innerHTML = '';
  targets.forEach(t => {
    const itemEl = document.createElement('div');
    itemEl.className = 'target-item';
    itemEl.innerHTML = `
      <div class="target-left">
        <input type="checkbox" class="target-checkbox" id="chk-${t.id}" data-id="${t.id}" data-prefix="${prefix}" ${t.selected ? 'checked' : ''}>
        <div class="target-details">
          <h4>
            ${escapeHtml(t.name)}
            ${t.id.startsWith('antigravity-') ? '<span class="risk-badge risk-ai">AI Ajanı</span>' : ''}
            ${t.id.startsWith('macos-') ? '<span class="risk-badge risk-apple">macOS</span>' : ''}
            <span class="risk-badge risk-${t.risk}">${t.risk === 'safe' ? 'Güvenli' : t.risk === 'recommended' ? 'Önerilen' : 'Dikkat'}</span>
          </h4>
          <p>${escapeHtml(t.description)}</p>
          <div class="target-path">${escapeHtml(t.path)}</div>
        </div>
      </div>
      <div class="target-right">
        <div>
          <div class="target-size">${t.sizeStr}</div>
          <div class="target-count">${t.itemCount.toLocaleString()} öğe</div>
        </div>
        <button class="action-btn action-clean-single" data-id="${t.id}" title="Yalnızca bunu temizle">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"/></svg>
          Sil
        </button>
      </div>
    `;

    const chk = itemEl.querySelector('.target-checkbox');
    chk.addEventListener('change', (e) => {
      t.selected = e.target.checked;
      updateSelectedSize(prefix);
    });

    const singleBtn = itemEl.querySelector('.action-clean-single');
    singleBtn.addEventListener('click', () => {
      confirmAndClean([t.id], [], t.name, t.sizeStr);
    });

    container.appendChild(itemEl);
  });
}

function renderDynamicCachesList(container, caches) {
  if (!caches || caches.length === 0) {
    container.innerHTML = '<div class="loading-state">Ekstra büyük önbellek bulunamadı.</div>';
    return;
  }

  container.innerHTML = '';
  caches.forEach(c => {
    const itemEl = document.createElement('div');
    itemEl.className = 'target-item';
    itemEl.innerHTML = `
      <div class="target-left">
        <input type="checkbox" class="target-checkbox dynamic-cache-chk" data-path="${c.path}">
        <div class="target-details">
          <h4>${escapeHtml(c.name)} <span class="risk-badge risk-safe">Uygulama Önbelleği</span></h4>
          <div class="target-path">${escapeHtml(c.path)}</div>
        </div>
      </div>
      <div class="target-right">
        <div>
          <div class="target-size">${c.sizeStr}</div>
          <div class="target-count">${c.fileCount.toLocaleString()} dosya</div>
        </div>
        <button class="action-btn action-reveal-single" data-path="${c.path}" title="Finder'da Göster">
          Finder
        </button>
      </div>
    `;

    const chk = itemEl.querySelector('.dynamic-cache-chk');
    chk.addEventListener('change', (e) => {
      c.selected = e.target.checked;
      updateSelectedSize('sys');
    });

    const revealBtn = itemEl.querySelector('.action-reveal-single');
    revealBtn.addEventListener('click', () => {
      revealInFinder(c.path);
    });

    container.appendChild(itemEl);
  });
}

// Update Selected Size Counters
function updateSelectedSize(prefix) {
  if (!state.scanData) return;

  if (prefix === 'dev') {
    let size = 0;
    state.scanData.targets.filter(t => t.category === 'developer' && t.selected).forEach(t => size += t.size);
    elements.devSelectedSize.textContent = formatBytes(size);
  } else if (prefix === 'sys') {
    let size = 0;
    state.scanData.targets.filter(t => t.category === 'system' && t.selected).forEach(t => size += t.size);
    if (state.scanData.dynamicCaches) {
      state.scanData.dynamicCaches.filter(c => c.selected).forEach(c => size += c.size);
    }
    elements.sysSelectedSize.textContent = formatBytes(size);
  }
}

// Clean Execution with Confirmation Modal
let pendingCleanAction = null;

function confirmAndClean(targetIds, customPaths, title, freedStr) {
  const isDryRun = elements.toggleDryRun.checked;
  const isTrash = elements.toggleTrashMode.checked;

  elements.modalTitle.textContent = isDryRun ? `Test Temizliği: ${title}` : `Temizleme Onayı: ${title}`;
  elements.modalMessage.textContent = isDryRun 
    ? "Simülasyon modunda hiçbir dosya gerçekten silinmeyecektir, sadece temizlenebilir alan kontrol edilir." 
    : "Seçili dosyalar silinecektir. Devam etmek istiyor musunuz?";
  elements.modalFreedSize.textContent = freedStr;
  elements.modalMethodText.textContent = isDryRun ? "Simülasyon (Dosyalar Silinmez)" : (isTrash ? "Finder Çöp Kutusuna Taşı" : "Kalıcı Olarak Sil");

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    await executeCleanRequest({
      targetIds: targetIds,
      customPaths: customPaths,
      dryRun: isDryRun,
      moveToTrash: isTrash,
    });
  };

  elements.confirmModal.showModal();
}

async function executeCleanRequest(requestBody) {
  if (state.isCleaning) return;
  state.isCleaning = true;

  showToast(requestBody.dryRun ? 'Simülasyon çalıştırılıyor...' : 'Temizleme işlemi başlatıldı...', 'info');

  try {
    const res = await fetch('/api/clean', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(requestBody),
    });

    if (!res.ok) throw new Error('Temizleme hatası oluştu');
    const result = await res.json();

    state.isCleaning = false;
    if (requestBody.dryRun) {
      showToast(`Test tamamlandı: ${result.freedStr || '0 B'} alan açılabilir! (Hiçbir dosya silinmedi)`, 'success');
    } else {
      showToast(`Harika! ${result.freedStr || '0 B'} alan başarıyla temizlendi!`, 'success');
      // Refresh stats and targets
      fetchSystemStats();
      startGlobalScan();
    }
  } catch (err) {
    state.isCleaning = false;
    console.error(err);
    showToast('Temizlik sırasında hata: ' + err.message, 'error');
  }
}

// Quick Clean (Predefined Recommended Targets or Smart Care)
function handleQuickClean() {
  if (state.scanData && state.scanData.targets) {
    const recommendedIds = state.scanData.targets
      .filter(t => t.exists && (t.risk === 'safe' || t.risk === 'recommended') && t.size > 0)
      .map(t => t.id);

    if (recommendedIds.length > 0) {
      confirmAndClean(recommendedIds, [], "Önerilen Güvenli Dosyalar", elements.quickCleanSize.textContent);
      return;
    }
  }

  if (state.smartCareData && state.smartCareData.totalCleanable > 0) {
    executeSmartCareClean();
    return;
  }

  showToast('Temizlenebilecek önerilen dosya bulunamadı!', 'info');
}

// Node Modules Scanner
async function startNodeModulesScan() {
  const rootPath = elements.inputNodePath.value.trim();
  elements.btnScanNodeModules.disabled = true;
  elements.tbodyNodeModules.innerHTML = '<tr><td colspan="6" class="loading-state">node_modules klasörleri taranıyor...</td></tr>';
  state.selectedNodeModules.clear();
  elements.btnCleanSelectedNode.disabled = true;
  elements.nodeSelectedSize.textContent = '0 MB';

  const url = rootPath 
    ? `/api/nodemodules?stream=true&root=${encodeURIComponent(rootPath)}`
    : '/api/nodemodules?stream=true';

  const eventSource = createEventSource(url);
  const items = [];

  eventSource.addEventListener('progress', (e) => {
    try {
      const data = JSON.parse(e.data);
      if (data.item) {
        items.push(data.item);
        renderNodeModulesTable(items);
      }
    } catch (err) {
      console.error(err);
    }
  });

  eventSource.addEventListener('result', (e) => {
    eventSource.close();
    elements.btnScanNodeModules.disabled = false;
    try {
      const fullList = JSON.parse(e.data);
      state.nodeModulesData = fullList;
      renderNodeModulesTable(fullList);
      showToast(`${fullList.length} adet node_modules klasörü bulundu!`, 'success');
    } catch (err) {
      console.error(err);
    }
  });

  eventSource.addEventListener('error', (e) => {
    eventSource.close();
    elements.btnScanNodeModules.disabled = false;
    showToast('node_modules taramasında hata oluştu.', 'error');
  });
}

function renderNodeModulesTable(items) {
  if (!items || items.length === 0) {
    elements.tbodyNodeModules.innerHTML = '<tr><td colspan="6" class="empty-state"><div class="empty-wrap"><p>Hiç node_modules klasörü bulunamadı.</p></div></td></tr>';
    return;
  }

  elements.tbodyNodeModules.innerHTML = '';
  items.forEach(item => {
    const tr = document.createElement('tr');
    tr.innerHTML = `
      <td><input type="checkbox" class="chk-node-item" data-path="${item.modulesPath}" data-size="${item.size}" ${state.selectedNodeModules.has(item.modulesPath) ? 'checked' : ''}></td>
      <td><strong>${escapeHtml(item.projectName)}</strong> ${item.hasPkgJson ? '<span class="badge-pro">package.json</span>' : ''}</td>
      <td><code>${escapeHtml(item.projectPath)}</code></td>
      <td>${item.lastModStr}</td>
      <td><strong>${item.sizeStr}</strong> <small style="color:var(--text-dim);">(${item.fileCount.toLocaleString()} dosya)</small></td>
      <td>
        <div class="table-actions">
          <button class="action-btn btn-node-reveal" data-path="${item.projectPath}">Finder</button>
          <button class="action-btn action-btn-danger btn-node-delete" data-path="${item.modulesPath}" data-size="${item.sizeStr}">Sil</button>
        </div>
      </td>
    `;

    const chk = tr.querySelector('.chk-node-item');
    chk.addEventListener('change', (e) => {
      if (e.target.checked) {
        state.selectedNodeModules.add(item.modulesPath);
      } else {
        state.selectedNodeModules.delete(item.modulesPath);
      }
      updateNodeSelectedSize();
    });

    tr.querySelector('.btn-node-reveal').addEventListener('click', () => {
      revealInFinder(item.projectPath);
    });

    tr.querySelector('.btn-node-delete').addEventListener('click', () => {
      confirmCleanNodeModules([item.modulesPath], item.sizeStr);
    });

    elements.tbodyNodeModules.appendChild(tr);
  });
}

function updateNodeSelectedSize() {
  let size = 0;
  state.nodeModulesData.forEach(item => {
    if (state.selectedNodeModules.has(item.modulesPath)) {
      size += item.size;
    }
  });

  elements.nodeSelectedSize.textContent = formatBytes(size);
  elements.btnCleanSelectedNode.disabled = state.selectedNodeModules.size === 0;
}

function confirmCleanNodeModules(paths, sizeStr) {
  const isTrash = elements.toggleTrashMode.checked;
  elements.modalTitle.textContent = "node_modules Temizleme Onayı";
  elements.modalMessage.textContent = `${paths.length} adet node_modules klasörü silinecek. Projelerinizde 'npm install' yaparak tekrar yükleyebilirsiniz.`;
  elements.modalFreedSize.textContent = sizeStr;
  elements.modalMethodText.textContent = isTrash ? "Finder Çöp Kutusuna Taşı" : "Kalıcı Olarak Sil";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('node_modules klasörleri temizleniyor...', 'info');
    try {
      const res = await fetch('/api/nodemodules/clean', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ paths: paths, moveToTrash: isTrash }),
      });
      const data = await res.json();
      showToast(`${data.deletedCount} klasör silindi, ${data.freedStr} alan açıldı!`, 'success');
      startNodeModulesScan();
      fetchSystemStats();
    } catch (err) {
      showToast('Silme hatası: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

// Large Files Scanner
async function startLargeFilesScan() {
  const minSizeMb = elements.selectMinSize.value;
  const category = elements.selectCategory.value;
  const ageDays = elements.selectAge.value;
  const customPath = elements.inputLargefilePath.value.trim();

  elements.btnScanLargeFiles.disabled = true;
  elements.tbodyLargeFiles.innerHTML = '<tr><td colspan="7" class="loading-state">Büyük dosyalar taranıyor...</td></tr>';
  state.selectedLargeFiles.clear();
  elements.btnCleanSelectedFiles.disabled = true;
  elements.filesSelectedSize.textContent = '0 MB';

  const params = new URLSearchParams({
    minSizeMb: minSizeMb,
    category: category,
    minAgeDays: ageDays,
  });
  if (customPath) {
    params.set('path', customPath);
  }

  try {
    const res = await fetch(`/api/largefiles?${params.toString()}`);
    if (!res.ok) throw new Error('Tarama başarısız');
    const items = await res.json();
    state.largeFilesData = items || [];
    renderLargeFilesTable(state.largeFilesData);
    showToast(`${state.largeFilesData.length} adet büyük dosya bulundu!`, 'success');
  } catch (err) {
    console.error(err);
    showToast('Büyük dosya taramasında hata: ' + err.message, 'error');
  } finally {
    elements.btnScanLargeFiles.disabled = false;
  }
}

function renderLargeFilesTable(items) {
  if (!items || items.length === 0) {
    elements.tbodyLargeFiles.innerHTML = '<tr><td colspan="7" class="empty-state"><div class="empty-wrap"><p>Belirtilen filtrelere uygun büyük dosya bulunamadı.</p></div></td></tr>';
    return;
  }

  elements.tbodyLargeFiles.innerHTML = '';
  items.forEach(file => {
    const tr = document.createElement('tr');
    tr.innerHTML = `
      <td><input type="checkbox" class="chk-file-item" data-path="${file.path}" data-size="${file.size}" ${state.selectedLargeFiles.has(file.path) ? 'checked' : ''}></td>
      <td><strong>${escapeHtml(file.name)}</strong></td>
      <td><span class="file-type-badge badge-${file.category}">${file.category}</span></td>
      <td><code>${escapeHtml(file.path)}</code></td>
      <td>${file.ageDays > 0 ? `${file.ageDays} gün önce` : 'Bugün'} <small style="color:var(--text-dim);">(${file.lastModStr})</small></td>
      <td><strong>${file.sizeStr}</strong></td>
      <td>
        <div class="table-actions">
          <button class="action-btn btn-file-reveal" data-path="${file.path}">Finder</button>
          <button class="action-btn action-btn-danger btn-file-delete" data-path="${file.path}" data-size="${file.sizeStr}">Sil</button>
        </div>
      </td>
    `;

    const chk = tr.querySelector('.chk-file-item');
    chk.addEventListener('change', (e) => {
      if (e.target.checked) {
        state.selectedLargeFiles.add(file.path);
      } else {
        state.selectedLargeFiles.delete(file.path);
      }
      updateLargeFilesSelectedSize();
    });

    tr.querySelector('.btn-file-reveal').addEventListener('click', () => {
      revealInFinder(file.path);
    });

    tr.querySelector('.btn-file-delete').addEventListener('click', () => {
      confirmCleanLargeFiles([file.path], file.sizeStr);
    });

    elements.tbodyLargeFiles.appendChild(tr);
  });
}

function updateLargeFilesSelectedSize() {
  let size = 0;
  state.largeFilesData.forEach(item => {
    if (state.selectedLargeFiles.has(item.path)) {
      size += item.size;
    }
  });

  elements.filesSelectedSize.textContent = formatBytes(size);
  elements.btnCleanSelectedFiles.disabled = state.selectedLargeFiles.size === 0;
}

function confirmCleanLargeFiles(paths, sizeStr) {
  const isTrash = elements.toggleTrashMode.checked;
  elements.modalTitle.textContent = "Büyük Dosyaları Silme Onayı";
  elements.modalMessage.textContent = `${paths.length} adet büyük dosya silinecek.`;
  elements.modalFreedSize.textContent = sizeStr;
  elements.modalMethodText.textContent = isTrash ? "Finder Çöp Kutusuna Taşı" : "Kalıcı Olarak Sil";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('Dosyalar siliniyor...', 'info');
    try {
      const res = await fetch('/api/largefiles/delete', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ paths: paths, moveToTrash: isTrash }),
      });
      const data = await res.json();
      showToast(`${data.deletedCount} dosya silindi, ${data.freedStr} alan açıldı!`, 'success');
      startLargeFilesScan();
      fetchSystemStats();
    } catch (err) {
      showToast('Silme hatası: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

// Finder Reveal Helper
async function revealInFinder(filePath) {
  try {
    const res = await fetch('/api/reveal', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: filePath }),
    });
    if (!res.ok) throw new Error('Finder açılamadı');
    showToast('Dosya Finder\'da gösterildi', 'info');
  } catch (err) {
    showToast('Finder hatası: ' + err.message, 'error');
  }
}

// ==========================================
// Directory Tree (Disk Space Breakdown)
// ==========================================
async function loadDirTree(targetPath) {
  targetPath = targetPath || state.treeCurrentPath || '~';
  state.treeCurrentPath = targetPath;
  elements.inputTreePath.value = targetPath;
  elements.treeItemsList.innerHTML = '<div class="loading-state">Dizin taranıyor ve boyutlar hesaplanıyor...</div>';
  elements.btnTreeRefresh.disabled = true;

  try {
    const res = await fetch(`/api/tree?path=${encodeURIComponent(targetPath)}`);
    if (!res.ok) throw new Error('Dizin okunamadı');
    const data = await res.json();
    state.treeData = data;
    renderDirTree(data);
    renderSunburstChart(data);
  } catch (err) {
    elements.treeItemsList.innerHTML = `<div class="loading-state" style="color:var(--accent-rose);">${escapeHtml(err.message)}</div>`;
    showToast('Dizin taranırken hata: ' + err.message, 'error');
  } finally {
    elements.btnTreeRefresh.disabled = false;
  }
}

function renderDirTree(data) {
  if (!data) return;

  elements.treeTotalSize.textContent = data.totalStr;
  elements.btnTreeUp.disabled = !data.parentPath;

  if (!data.items || data.items.length === 0) {
    elements.treeItemsList.innerHTML = '<div class="loading-state">Bu dizin boş.</div>';
    return;
  }

  elements.treeItemsList.innerHTML = '';
  data.items.forEach(item => {
    const row = document.createElement('div');
    row.className = 'tree-row';
    const isFolder = item.isDir;
    const pct = Math.max(0.5, Math.min(100, item.percentage)).toFixed(1);

    row.innerHTML = `
      <div class="tree-left">
        <div class="tree-icon ${isFolder ? 'tree-icon-folder' : 'tree-icon-file'}">
          ${isFolder 
            ? '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 9.776c.112-.017.227-.026.344-.026h15.812c.117 0 .232.009.344.026m-16.5 0a2.25 2.25 0 00-1.883 2.542l.857 6a2.25 2.25 0 002.227 1.932H19.05a2.25 2.25 0 002.227-1.932l.857-6a2.25 2.25 0 00-1.883-2.542m-16.5 0V6A2.25 2.25 0 016 3.75h3.879a1.5 1.5 0 011.06.44l2.122 2.12a1.5 1.5 0 001.06.44H18A2.25 2.25 0 0120.25 9v.776"/></svg>'
            : '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"/></svg>'
          }
        </div>
        <div class="tree-info">
          <div class="tree-name-wrap">
            <span class="tree-name ${isFolder ? 'folder-link' : ''}" title="${escapeHtml(item.path)}">${escapeHtml(item.name)}</span>
            ${isFolder ? '<small style="color:var(--text-dim);">(' + item.itemCount.toLocaleString() + ' öğe)</small>' : ''}
          </div>
          <div class="tree-bar-bg">
            <div class="tree-bar-fill" style="width: ${pct}%;"></div>
          </div>
          <div class="tree-meta">Son değişiklik: ${item.lastModified || 'Bilinmiyor'}</div>
        </div>
      </div>
      <div class="tree-right">
        <div class="tree-size-box">
          <div class="tree-size-val">${item.sizeStr}</div>
          <div class="tree-pct-val">%${pct}</div>
        </div>
        <div class="table-actions">
          ${isFolder ? `<button class="action-btn btn-tree-drill" data-path="${escapeHtml(item.path)}" title="İçine gir">Aç</button>` : ''}
          <button class="action-btn btn-tree-reveal" data-path="${escapeHtml(item.path)}" title="Finder'da göster">Finder</button>
          <button class="action-btn action-btn-danger btn-tree-delete" data-path="${escapeHtml(item.path)}" data-size="${item.sizeStr}">Sil</button>
        </div>
      </div>
    `;

    if (isFolder) {
      row.querySelector('.tree-name').addEventListener('click', () => loadDirTree(item.path));
      row.querySelector('.btn-tree-drill').addEventListener('click', () => loadDirTree(item.path));
    }
    row.querySelector('.btn-tree-reveal').addEventListener('click', () => revealInFinder(item.path));
    row.querySelector('.btn-tree-delete').addEventListener('click', () => {
      confirmCleanCustomPaths([item.path], item.name, item.sizeStr, () => loadDirTree(data.currentPath));
    });

    elements.treeItemsList.appendChild(row);
  });
}

function confirmCleanCustomPaths(paths, title, sizeStr, onSuccess) {
  const isTrash = elements.toggleTrashMode.checked;
  elements.modalTitle.textContent = `Silme Onayı: ${title}`;
  elements.modalMessage.textContent = `Seçilen öğe (${paths.length} adet) silinecektir.`;
  elements.modalFreedSize.textContent = sizeStr;
  elements.modalMethodText.textContent = isTrash ? "Finder Çöp Kutusuna Taşı" : "Kalıcı Olarak Sil";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('Siliniyor...', 'info');
    try {
      const res = await fetch('/api/clean', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ customPaths: paths, moveToTrash: isTrash }),
      });
      const data = await res.json();
      showToast(`${data.freedStr || sizeStr} alan açıldı!`, 'success');
      if (onSuccess) onSuccess();
      fetchSystemStats();
    } catch (err) {
      showToast('Silme hatası: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

// ==========================================
// Duplicate Files Finder
// ==========================================
async function startDuplicatesScan() {
  const minSizeMb = elements.selectDupMinSize.value;
  const customPath = elements.inputDupPath.value.trim();

  elements.btnScanDuplicates.disabled = true;
  elements.duplicateGroupsList.innerHTML = '<div class="loading-state">Yinelenen dosyalar taranıyor ve SHA-256 hash imzaları hesaplanıyor...</div>';
  state.selectedDuplicates.clear();
  elements.btnDupCleanSelected.disabled = true;
  elements.dupSelectedSize.textContent = '0 MB';

  const params = new URLSearchParams({ minSizeMb: minSizeMb });
  if (customPath) params.set('path', customPath);

  try {
    const res = await fetch(`/api/duplicates?${params.toString()}`);
    if (!res.ok) throw new Error('Tarama başarısız');
    const groups = await res.json();
    state.duplicatesData = groups || [];
    renderDuplicates(state.duplicatesData);
    let totalWasted = 0;
    state.duplicatesData.forEach(g => totalWasted += g.wastedSize);
    showToast(`${state.duplicatesData.length} grup yinelenen dosya bulundu (${formatBytes(totalWasted)} israf alan)`, 'success');
  } catch (err) {
    showToast('Yinelenen dosya taramasında hata: ' + err.message, 'error');
    elements.duplicateGroupsList.innerHTML = `<div class="loading-state" style="color:var(--accent-rose);">${escapeHtml(err.message)}</div>`;
  } finally {
    elements.btnScanDuplicates.disabled = false;
  }
}

function renderDuplicates(groups) {
  if (!groups || groups.length === 0) {
    elements.duplicateGroupsList.innerHTML = `
      <div class="empty-state">
        <div class="empty-wrap">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          <p>Tebrikler! Belirtilen kriterlerde hiçbir yinelenen dosya bulunamadı.</p>
        </div>
      </div>
    `;
    return;
  }

  elements.duplicateGroupsList.innerHTML = '';
  groups.forEach((group) => {
    const card = document.createElement('div');
    card.className = 'dup-group-card';
    card.innerHTML = `
      <div class="dup-group-header">
        <div class="dup-header-left">
          <span class="dup-count-badge">${group.files.length} Kopya</span>
          <strong>${escapeHtml(group.files[0].name)}</strong>
          <span style="color:var(--text-dim); font-size:0.75rem;">(Boyut: ${group.sizeStr})</span>
        </div>
        <div class="dup-wasted-pill">
          ${group.wastedSizeStr} İsraf Alan
        </div>
      </div>
      <div class="dup-file-list"></div>
    `;

    const fileListEl = card.querySelector('.dup-file-list');
    group.files.forEach((file, fIdx) => {
      const row = document.createElement('div');
      row.className = 'dup-file-row';
      const isCopy = fIdx > 0;
      if (file.selected) {
        state.selectedDuplicates.add(file.path);
      }

      row.innerHTML = `
        <div class="dup-file-left">
          <input type="checkbox" class="chk-dup-file" data-path="${escapeHtml(file.path)}" data-size="${file.size}" ${file.selected ? 'checked' : ''}>
          <div class="dup-file-info">
            <h5>
              ${escapeHtml(file.name)} 
              ${!isCopy ? '<span class="badge-pro" style="margin-left:6px;">Orijinal</span>' : '<span class="risk-badge risk-recommended" style="margin-left:6px;">Kopya</span>'}
            </h5>
            <code>${escapeHtml(file.path)}</code>
          </div>
        </div>
        <div class="table-actions">
          <span style="font-size:0.8rem; color:var(--text-dim); margin-right:10px;">${file.modTimeStr}</span>
          <button class="action-btn btn-dup-reveal" data-path="${escapeHtml(file.path)}">Finder</button>
        </div>
      `;

      const chk = row.querySelector('.chk-dup-file');
      chk.addEventListener('change', (e) => {
        file.selected = e.target.checked;
        if (e.target.checked) {
          state.selectedDuplicates.add(file.path);
        } else {
          state.selectedDuplicates.delete(file.path);
        }
        updateDuplicatesSelectedSize();
      });

      row.querySelector('.btn-dup-reveal').addEventListener('click', () => revealInFinder(file.path));

      fileListEl.appendChild(row);
    });

    elements.duplicateGroupsList.appendChild(card);
  });

  updateDuplicatesSelectedSize();
}

function selectOnlyCopies() {
  state.selectedDuplicates.clear();
  state.duplicatesData.forEach(group => {
    group.files.forEach((file, idx) => {
      file.selected = (idx > 0);
      if (file.selected) {
        state.selectedDuplicates.add(file.path);
      }
    });
  });
  renderDuplicates(state.duplicatesData);
}

function updateDuplicatesSelectedSize() {
  let size = 0;
  state.duplicatesData.forEach(group => {
    group.files.forEach(file => {
      if (state.selectedDuplicates.has(file.path)) {
        size += file.size;
      }
    });
  });

  elements.dupSelectedSize.textContent = formatBytes(size);
  elements.btnDupCleanSelected.disabled = state.selectedDuplicates.size === 0;
}

function confirmCleanDuplicates() {
  const paths = Array.from(state.selectedDuplicates);
  if (paths.length === 0) return;

  const isTrash = elements.toggleTrashMode.checked;
  elements.modalTitle.textContent = "Yinelenen Dosyaları Silme Onayı";
  elements.modalMessage.textContent = `${paths.length} adet kopya dosya silinecektir.`;
  elements.modalFreedSize.textContent = elements.dupSelectedSize.textContent;
  elements.modalMethodText.textContent = isTrash ? "Finder Çöp Kutusuna Taşı" : "Kalıcı Olarak Sil";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('Kopya dosyalar siliniyor...', 'info');
    try {
      const res = await fetch('/api/duplicates/delete', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ paths: paths, moveToTrash: isTrash }),
      });
      const data = await res.json();
      showToast(`${data.deletedCount} dosya silindi, ${data.freedStr} alan açıldı!`, 'success');
      startDuplicatesScan();
      fetchSystemStats();
    } catch (err) {
      showToast('Silme hatası: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

// ==========================================
// Leftover & Residual Files Cleaner
// ==========================================
async function startLeftoversScan(silent = false) {
  if (!silent) {
    elements.btnScanLeftovers.disabled = true;
    elements.leftoversItemsList.innerHTML = '<div class="loading-state">Kaldırılmış uygulama artıkları ve sistem çöpleri taranıyor...</div>';
    state.selectedLeftovers.clear();
    elements.btnCleanSelectedLeftovers.disabled = true;
    elements.leftoversSelectedSize.textContent = '0 MB';
  }

  try {
    const res = await fetch('/api/leftovers');
    if (!res.ok) throw new Error('Artık dosyalar taranamadı');
    const data = await res.json();
    state.leftoversData = data;

    // Update stats and badge
    elements.badgeLeftoversCount.textContent = data.items ? data.items.length : 0;
    elements.leftoversAppsCount.textContent = `${data.orphanedAppsCount || 0} Uygulama`;
    elements.leftoversDsCount.textContent = `${data.dsStoreCount || 0} Dosya`;
    elements.leftoversTotalSize.textContent = data.totalSizeStr || '0 MB';

    if (!silent) {
      renderLeftovers(data);
      showToast(`${data.items ? data.items.length : 0} adet artık dosya bulundu (${data.totalSizeStr})`, 'success');
    }
  } catch (err) {
    if (!silent) {
      elements.leftoversItemsList.innerHTML = `<div class="loading-state" style="color:var(--accent-rose);">${escapeHtml(err.message)}</div>`;
      showToast('Artık dosya taramasında hata: ' + err.message, 'error');
    }
  } finally {
    if (!silent) elements.btnScanLeftovers.disabled = false;
  }
}

function renderLeftovers(data) {
  if (!data || !data.items || data.items.length === 0) {
    elements.leftoversItemsList.innerHTML = `
      <div class="empty-state">
        <div class="empty-wrap">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          <p>Harika! Sisteminizde hiçbir artık uygulama verisi bulunamadı.</p>
        </div>
      </div>
    `;
    return;
  }

  elements.leftoversItemsList.innerHTML = '';
  data.items.forEach(item => {
    const el = document.createElement('div');
    el.className = 'target-item';
    if (item.selected) {
      state.selectedLeftovers.add(item.id);
    }

    el.innerHTML = `
      <div class="target-left">
        <input type="checkbox" class="target-checkbox chk-leftover-item" data-id="${escapeHtml(item.id)}" data-size="${item.size}" ${item.selected ? 'checked' : ''}>
        <div class="target-details">
          <h4>
            ${escapeHtml(item.name)}
            <span class="risk-badge ${item.category === 'orphaned_app' ? 'risk-recommended' : 'risk-safe'}">${escapeHtml(item.categoryStr)}</span>
          </h4>
          <p>${escapeHtml(item.description)}</p>
          <div class="target-path">${escapeHtml(item.path)}</div>
        </div>
      </div>
      <div class="target-right">
        <div>
          <div class="target-size">${item.sizeStr}</div>
          <div class="target-count">${item.fileCount.toLocaleString()} öğe</div>
        </div>
        ${item.id !== 'all-ds-store' && item.id !== 'all-temp-files' ? `<button class="action-btn btn-leftover-reveal" data-path="${escapeHtml(item.path)}">Finder</button>` : ''}
        <button class="action-btn action-btn-danger btn-leftover-delete" data-id="${escapeHtml(item.id)}" data-name="${escapeHtml(item.name)}" data-size="${item.sizeStr}">Sil</button>
      </div>
    `;

    const chk = el.querySelector('.chk-leftover-item');
    chk.addEventListener('change', (e) => {
      item.selected = e.target.checked;
      if (e.target.checked) {
        state.selectedLeftovers.add(item.id);
      } else {
        state.selectedLeftovers.delete(item.id);
      }
      updateLeftoversSelectedSize();
    });

    const revealBtn = el.querySelector('.btn-leftover-reveal');
    if (revealBtn) {
      revealBtn.addEventListener('click', () => revealInFinder(item.path));
    }

    el.querySelector('.btn-leftover-delete').addEventListener('click', () => {
      confirmCleanLeftovers([item.id], item.name, item.sizeStr);
    });

    elements.leftoversItemsList.appendChild(el);
  });

  updateLeftoversSelectedSize();
}

function updateLeftoversSelectedSize() {
  let size = 0;
  if (state.leftoversData && state.leftoversData.items) {
    state.leftoversData.items.forEach(item => {
      if (state.selectedLeftovers.has(item.id)) {
        size += item.size;
      }
    });
  }

  elements.leftoversSelectedSize.textContent = formatBytes(size);
  elements.btnCleanSelectedLeftovers.disabled = state.selectedLeftovers.size === 0;
}

function confirmCleanLeftovers(ids, title, sizeStr) {
  ids = ids || Array.from(state.selectedLeftovers);
  if (ids.length === 0) return;

  const isTrash = elements.toggleTrashMode.checked;
  elements.modalTitle.textContent = title ? `Artık Dosya Silme: ${title}` : "Seçili Artıkları Silme Onayı";
  elements.modalMessage.textContent = `${ids.length} adet artık öğe silinecektir.`;
  elements.modalFreedSize.textContent = sizeStr || elements.leftoversSelectedSize.textContent;
  elements.modalMethodText.textContent = isTrash ? "Finder Çöp Kutusuna Taşı" : "Kalıcı Olarak Sil";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('Artık dosyalar siliniyor...', 'info');
    try {
      const res = await fetch('/api/leftovers/clean', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ids: ids, moveToTrash: isTrash }),
      });
      const data = await res.json();
      showToast(`${data.deletedCount} öğe silindi, ${data.freedStr} alan açıldı!`, 'success');
      startLeftoversScan();
      fetchSystemStats();
    } catch (err) {
      showToast('Silme hatası: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

// ==========================================================================
// App Uninstaller & Reset Module
// ==========================================================================
async function startAppsScan() {
  if (!elements.btnScanApps) return;
  elements.btnScanApps.disabled = true;
  elements.appsItemsList.innerHTML = '<div class="loading-state">Yüklü uygulamalar taranıyor ve disk alanları hesaplanıyor...</div>';

  try {
    const res = await fetch('/api/apps');
    if (!res.ok) throw new Error('Uygulama listesi alınamadı');
    const data = await res.json();
    state.appsData = data;

    if (elements.appsCountBadge) elements.appsCountBadge.textContent = `${data.totalAppsCount || 0} Uygulama`;
    if (elements.chipAppsAllCount) elements.chipAppsAllCount.textContent = `${data.totalAppsCount || 0}`;
    if (elements.appsBinarySizeBadge) elements.appsBinarySizeBadge.textContent = data.totalAppSizeStr || '0 GB';
    if (elements.appsDataSizeBadge) elements.appsDataSizeBadge.textContent = data.totalDataSizeStr || '0 GB';

    renderAppsList();
    showToast(`${data.totalAppsCount} uygulama başarıyla tarandı.`, 'success');
  } catch (err) {
    console.error(err);
    elements.appsItemsList.innerHTML = `<div class="loading-state" style="color:var(--accent-rose);">${escapeHtml(err.message)}</div>`;
    showToast('Uygulama tarama hatası: ' + err.message, 'error');
  } finally {
    elements.btnScanApps.disabled = false;
  }
}

function renderAppsList() {
  if (!state.appsData || !state.appsData.apps) return;

  const query = (state.appsSearch || '').toLowerCase().trim();
  const filter = state.appsFilter || 'all';

  let list = state.appsData.apps.filter(app => {
    // Search query filter
    if (query) {
      const matchName = app.name.toLowerCase().includes(query);
      const matchBundle = (app.bundleId || '').toLowerCase().includes(query);
      if (!matchName && !matchBundle) return false;
    }
    // Category chips filter
    if (filter === 'system') return app.isSystem;
    if (filter === 'user') return !app.isSystem;
    if (filter === 'large') return app.totalSize >= 500 * 1024 * 1024; // > 500MB
    return true;
  });

  if (list.length === 0) {
    elements.appsItemsList.innerHTML = `
      <div class="empty-state">
        <div class="empty-wrap">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"/></svg>
          <p>Arama veya filtre kriterlerine uygun uygulama bulunamadı.</p>
        </div>
      </div>
    `;
    return;
  }

  elements.appsItemsList.innerHTML = '';
  list.forEach(app => {
    const card = document.createElement('div');
    card.className = 'app-card-item glass-card';
    card.id = `app-card-${encodeURIComponent(app.id)}`;

    const compCount = app.components ? app.components.length : 0;

    card.innerHTML = `
      <div class="app-card-main">
        <div class="app-left">
          <div class="app-icon-wrap">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25z"/></svg>
          </div>
          <div class="app-info">
            <h4>
              ${escapeHtml(app.name)}
              ${app.version ? `<small style="color:var(--text-dim); margin-left:6px;">v${escapeHtml(app.version)}</small>` : ''}
              <span class="app-badge ${app.isSystem ? 'badge-system' : 'badge-user'}">${app.isSystem ? 'macOS Sistem' : 'Kullanıcı'}</span>
            </h4>
            <div class="app-bundle-id">${escapeHtml(app.bundleId || app.path)}</div>
          </div>
        </div>

        <div class="app-size-breakdown">
          <div class="app-size-col">
            <span class="app-size-label">Paket (.app)</span>
            <span class="app-size-val">${app.appSizeStr}</span>
          </div>
          <div class="app-size-col">
            <span class="app-size-label">Veri/Önbellek</span>
            <span class="app-size-val">${app.dataSizeStr}</span>
          </div>
          <div class="app-size-col">
            <span class="app-size-label">Toplam</span>
            <span class="app-size-val app-size-total">${app.totalSizeStr}</span>
          </div>
        </div>

        <div class="app-actions">
          <button class="action-btn btn-app-details" title="Bileşenleri Göster/Gizle">
            <span>Bileşenler (${compCount})</span>
            <svg style="width:12px; height:12px; margin-left:4px;" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5"/></svg>
          </button>
          <button class="action-btn btn-app-reveal" title="Finder'da Göster">Finder</button>
          ${!app.isSystem && app.dataSize > 0 ? `<button class="action-btn btn-app-reset" style="color:var(--accent-amber);" title="Uygulamayı korur, önbellek ve verilerini sıfırlar">Sıfırla</button>` : ''}
          ${!app.isSystem ? `<button class="action-btn action-btn-danger btn-app-uninstall" title="Uygulamayı ve tüm kullanıcı verilerini kaldırır">Kaldır</button>` : ''}
        </div>
      </div>

      <div class="app-components-toggle" style="display: none;">
        ${(app.components || []).map(c => `
          <div class="comp-row">
            <div class="comp-left">
              <span class="comp-type-pill">${escapeHtml(c.type)}</span>
              <span class="comp-name" title="${escapeHtml(c.path)}">${escapeHtml(c.name)}</span>
              <span class="comp-path">${escapeHtml(c.path)}</span>
            </div>
            <span class="comp-size"><strong>${escapeHtml(c.sizeStr)}</strong></span>
          </div>
        `).join('')}
      </div>
    `;

    // Toggle details
    const detailsBtn = card.querySelector('.btn-app-details');
    const compBox = card.querySelector('.app-components-toggle');
    detailsBtn.addEventListener('click', () => {
      const isHidden = compBox.style.display === 'none';
      compBox.style.display = isHidden ? 'flex' : 'none';
      detailsBtn.classList.toggle('open', isHidden);
    });

    // Reveal in Finder
    card.querySelector('.btn-app-reveal').addEventListener('click', () => {
      revealInFinder(app.path);
    });

    // Reset button
    const resetBtn = card.querySelector('.btn-app-reset');
    if (resetBtn) {
      resetBtn.addEventListener('click', () => {
        confirmUninstallApp(app, true);
      });
    }

    // Uninstall button
    const uninstBtn = card.querySelector('.btn-app-uninstall');
    if (uninstBtn) {
      uninstBtn.addEventListener('click', () => {
        confirmUninstallApp(app, false);
      });
    }

    elements.appsItemsList.appendChild(card);
  });
}

function confirmUninstallApp(app, resetOnly) {
  const isTrash = elements.toggleTrashMode.checked;
  const actionTitle = resetOnly ? `${app.name} Uygulamasını Sıfırla` : `${app.name} Uygulamasını Tamamen Kaldır`;
  const actionMsg = resetOnly
    ? `Uygulama program dosyası (.app) korunacak, ancak Application Support, Caches, Tercihler ve Konteyner verileri silinerek fabrika ayarlarına döndürülecektir.`
    : `Uygulama (.app) ve sistemdeki tüm ilişkili verileri (Application Support, Caches, Preferences vb.) tamamen silinecektir.`;
  const sizeToFree = resetOnly ? app.dataSizeStr : app.totalSizeStr;

  elements.modalTitle.textContent = actionTitle;
  elements.modalMessage.textContent = actionMsg;
  elements.modalFreedSize.textContent = sizeToFree;
  elements.modalMethodText.textContent = isTrash ? "Finder Çöp Kutusuna Taşı" : "Kalıcı Olarak Sil";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast(resetOnly ? `${app.name} sıfırlanıyor...` : `${app.name} kaldırılıyor...`, 'info');
    try {
      const res = await fetch('/api/apps/uninstall', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          appId: app.id,
          resetOnly: resetOnly,
          useTrash: isTrash,
        }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'İşlem başarısız oldu');
      showToast(data.message || 'İşlem başarıyla tamamlandı', 'success');
      startAppsScan();
      fetchSystemStats();
    } catch (err) {
      console.error(err);
      showToast('Hata: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

// ==========================================================================
// Browser Caches Module
// ==========================================================================
async function startBrowsersScan() {
  if (!elements.btnScanBrowsers) return;
  elements.btnScanBrowsers.disabled = true;
  elements.browsersGridList.innerHTML = '<div class="loading-state">Web tarayıcıları ve önbellek dizinleri taranıyor...</div>';

  try {
    const res = await fetch('/api/browsers');
    if (!res.ok) throw new Error('Tarayıcı bilgileri alınamadı');
    const data = await res.json();
    state.browsersData = data;
    renderBrowsersGrid();
    showToast(`Tarayıcı önbellekleri tarandı (${data.totalSizeStr || '0 B'})`, 'success');
  } catch (err) {
    console.error(err);
    elements.browsersGridList.innerHTML = `<div class="loading-state" style="color:var(--accent-rose);">${escapeHtml(err.message)}</div>`;
    showToast('Tarayıcı tarama hatası: ' + err.message, 'error');
  } finally {
    elements.btnScanBrowsers.disabled = false;
  }
}

function renderBrowsersGrid() {
  if (!state.browsersData || !state.browsersData.browsers) return;

  const browsers = state.browsersData.browsers;
  let cleanableCount = 0;
  let totalCleanableSize = 0;

  elements.browsersGridList.innerHTML = '';
  browsers.forEach(b => {
    if (b.size > 0) {
      cleanableCount++;
      totalCleanableSize += b.size;
    }

    const card = document.createElement('div');
    card.className = 'browser-card glass-card';

    let iconSvg = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><circle cx="12" cy="12" r="4"/></svg>';

    card.innerHTML = `
      <div class="browser-header">
        <div class="browser-brand">
          <div class="browser-logo">${iconSvg}</div>
          <div class="browser-title">
            <h4>${escapeHtml(b.name)}</h4>
            <span class="browser-status-tag" style="${b.installed ? '' : 'background:rgba(255,255,255,0.06); color:var(--text-dim);'}">
              ${b.installed ? 'Yüklü' : 'Bulunamadı'}
            </span>
          </div>
        </div>
      </div>

      <div class="browser-desc">${escapeHtml(b.description)}</div>

      <div class="browser-stats-row">
        <span class="app-size-label">Temizlenebilir Önbellek:</span>
        <span class="browser-cache-size">${b.sizeStr}</span>
      </div>

      <div class="app-actions" style="margin-top:auto;">
        ${b.size > 0 ? `
          <button class="btn btn-primary btn-block btn-sm btn-browser-clean" data-id="${escapeHtml(b.id)}" data-name="${escapeHtml(b.name)}" data-size="${b.sizeStr}">
            ${escapeHtml(b.name)} Önbelleğini Temizle (${b.sizeStr})
          </button>
        ` : `
          <button class="btn btn-secondary btn-block btn-sm" disabled>Önbellek Temiz</button>
        `}
      </div>
    `;

    const cleanBtn = card.querySelector('.btn-browser-clean');
    if (cleanBtn) {
      cleanBtn.addEventListener('click', () => {
        confirmCleanBrowsers([b.id], `${b.name} Önbelleği`, b.sizeStr);
      });
    }

    elements.browsersGridList.appendChild(card);
  });

  if (elements.browsersTotalCleanSize) {
    elements.browsersTotalCleanSize.textContent = formatBytes(totalCleanableSize);
  }
  if (elements.btnCleanAllBrowsers) {
    elements.btnCleanAllBrowsers.disabled = cleanableCount === 0;
  }
}

function confirmCleanBrowsers(browserIds, title, sizeStr) {
  if (!browserIds || browserIds.length === 0) return;

  elements.modalTitle.textContent = title || "Tüm Tarayıcı Önbelleklerini Temizle";
  elements.modalMessage.textContent = "Tarayıcı geçici web görselleri, derlenmiş JS kodları ve GPU önbellekleri temizlenecektir. Yer imleriniz ve şifreleriniz ASLA silinmez.";
  elements.modalFreedSize.textContent = sizeStr || elements.browsersTotalCleanSize.textContent;
  elements.modalMethodText.textContent = "Doğrudan Boşalt (Güvenli Önbellek)";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('Tarayıcı önbellekleri temizleniyor...', 'info');
    try {
      const res = await fetch('/api/browsers/clean', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ browsers: browserIds }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Temizleme başarısız oldu');
      showToast(data.message || 'Tarayıcı önbellekleri başarıyla temizlendi', 'success');
      startBrowsersScan();
      fetchSystemStats();
    } catch (err) {
      console.error(err);
      showToast('Hata: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

// ==========================================================================
// Downloads Organizer Module
// ==========================================================================
async function startDownloadsScan() {
  if (!elements.btnScanDownloads) return;
  elements.btnScanDownloads.disabled = true;
  elements.downloadsItemsList.innerHTML = '<div class="loading-state">İndirilenler klasörü taranıyor ve kategorize ediliyor...</div>';
  state.selectedDownloads.clear();
  updateDownloadsSelectedSize();

  try {
    const res = await fetch('/api/downloads');
    if (!res.ok) throw new Error('İndirilenler taranamadı');
    const data = await res.json();
    state.downloadsData = data;

    if (elements.downloadsInstallersSize) elements.downloadsInstallersSize.textContent = `${data.installersStr || '0 MB'} (${data.installersCount || 0})`;
    if (elements.downloadsArchivesSize) elements.downloadsArchivesSize.textContent = `${data.archivesStr || '0 MB'} (${data.archivesCount || 0})`;
    if (elements.downloadsOldSize) elements.downloadsOldSize.textContent = `${data.oldFilesStr || '0 MB'} (${data.oldFilesCount || 0})`;

    renderDownloadsList();
    showToast(`${data.totalCount || 0} dosya analiz edildi (${data.totalSizeStr || '0 MB'})`, 'success');
  } catch (err) {
    console.error(err);
    elements.downloadsItemsList.innerHTML = `<div class="loading-state" style="color:var(--accent-rose);">${escapeHtml(err.message)}</div>`;
    showToast('Downloads tarama hatası: ' + err.message, 'error');
  } finally {
    elements.btnScanDownloads.disabled = false;
  }
}

function renderDownloadsList() {
  if (!state.downloadsData || !state.downloadsData.items) return;

  const catFilter = state.downloadsFilter || 'all';
  let list = state.downloadsData.items.filter(item => {
    if (catFilter === 'installer') return item.category === 'installer';
    if (catFilter === 'archive') return item.category === 'archive';
    if (catFilter === 'old') return item.isOld;
    return true;
  });

  if (list.length === 0) {
    elements.downloadsItemsList.innerHTML = `
      <div class="empty-state">
        <div class="empty-wrap">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          <p>Bu filtrede herhangi bir indirilen dosya bulunamadı.</p>
        </div>
      </div>
    `;
    return;
  }

  elements.downloadsItemsList.innerHTML = '';
  list.forEach(item => {
    const el = document.createElement('div');
    el.className = 'target-item';

    const isChecked = state.selectedDownloads.has(item.path);

    el.innerHTML = `
      <div class="target-left">
        <input type="checkbox" class="target-checkbox chk-download-item" data-path="${escapeHtml(item.path)}" data-size="${item.size}" ${isChecked ? 'checked' : ''}>
        <div class="target-details">
          <h4>
            ${escapeHtml(item.name)}
            <span class="risk-badge ${item.category === 'installer' ? 'risk-recommended' : 'risk-safe'}">${escapeHtml(item.category)}</span>
            ${item.isOld ? `<span class="risk-badge" style="background:rgba(244,63,94,0.15); color:var(--accent-rose);">${item.daysOld} gün önce</span>` : ''}
          </h4>
          <p>Son Değişiklik: ${escapeHtml(item.modDate)}</p>
          <div class="target-path">${escapeHtml(item.path)}</div>
        </div>
      </div>
      <div class="target-right">
        <div>
          <div class="target-size">${item.sizeStr}</div>
        </div>
        <button class="action-btn btn-dl-reveal" data-path="${escapeHtml(item.path)}">Finder</button>
        <button class="action-btn action-btn-danger btn-dl-delete" data-path="${escapeHtml(item.path)}" data-name="${escapeHtml(item.name)}" data-size="${item.sizeStr}">Sil</button>
      </div>
    `;

    const chk = el.querySelector('.chk-download-item');
    chk.addEventListener('change', (e) => {
      if (e.target.checked) {
        state.selectedDownloads.add(item.path);
      } else {
        state.selectedDownloads.delete(item.path);
      }
      updateDownloadsSelectedSize();
    });

    el.querySelector('.btn-dl-reveal').addEventListener('click', () => {
      revealInFinder(item.path);
    });

    el.querySelector('.btn-dl-delete').addEventListener('click', () => {
      confirmCleanDownloads([item.path], item.sizeStr, item.name);
    });

    elements.downloadsItemsList.appendChild(el);
  });
}

function updateDownloadsSelectedSize() {
  let size = 0;
  if (state.downloadsData && state.downloadsData.items) {
    state.downloadsData.items.forEach(i => {
      if (state.selectedDownloads.has(i.path)) {
        size += i.size;
      }
    });
  }

  if (elements.downloadsSelectedSize) elements.downloadsSelectedSize.textContent = formatBytes(size);
  if (elements.btnCleanDownloads) elements.btnCleanDownloads.disabled = state.selectedDownloads.size === 0;
}

function confirmCleanDownloads(paths, sizeStr, title) {
  paths = paths || Array.from(state.selectedDownloads);
  if (paths.length === 0) return;

  const isTrash = elements.toggleTrashMode.checked;
  elements.modalTitle.textContent = title ? `İndirilen Dosyayı Sil: ${title}` : `Seçilen ${paths.length} İndirilen Dosyayı Sil`;
  elements.modalMessage.textContent = `${paths.length} adet dosya silinecektir. Devam etmek istiyor musunuz?`;
  elements.modalFreedSize.textContent = sizeStr || elements.downloadsSelectedSize.textContent;
  elements.modalMethodText.textContent = isTrash ? "Finder Çöp Kutusuna Taşı" : "Kalıcı Olarak Sil";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('Dosyalar siliniyor...', 'info');
    try {
      const res = await fetch('/api/downloads/clean', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ paths: paths, useTrash: isTrash }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Silinemedi');
      showToast(data.message || 'Dosyalar silindi', 'success');
      startDownloadsScan();
      fetchSystemStats();
    } catch (err) {
      console.error(err);
      showToast('Hata: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

// ==========================================================================
// Apple & System Data Module
// ==========================================================================
async function startAppleScan() {
  if (!elements.btnScanApple) return;
  elements.btnScanApple.disabled = true;
  if (elements.appleSnapshotsList) elements.appleSnapshotsList.innerHTML = '<div class="loading-state">APFS anlık görüntüleri listeleniyor...</div>';
  if (elements.appleBackupsList) elements.appleBackupsList.innerHTML = '<div class="loading-state">iOS yedekleri taranıyor...</div>';
  if (elements.appleSimulatorsList) elements.appleSimulatorsList.innerHTML = '<div class="loading-state">Simülatörler taranıyor...</div>';
  if (elements.appleCategoriesContainer) elements.appleCategoriesContainer.innerHTML = '<div class="loading-state">macOS Signed System Volume ve Sistem Verileri taranıyor...</div>';

  try {
    const res = await fetch('/api/apple');
    if (!res.ok) throw new Error('Apple sistem verileri alınamadı');
    const data = await res.json();
    state.appleData = data;

    renderAppleSystem(data);
    showToast('macOS & Sistem Verileri analizi tamamlandı.', 'success');
  } catch (err) {
    console.error(err);
    if (elements.appleCategoriesContainer) {
      elements.appleCategoriesContainer.innerHTML = `<div class="loading-state" style="color:var(--accent-rose);">${escapeHtml(err.message)}</div>`;
    }
    showToast('Apple tarama hatası: ' + err.message, 'error');
  } finally {
    elements.btnScanApple.disabled = false;
  }
}

function renderAppleSystem(data) {
  // 1. macOS Signed System Volume Hero
  if (data.macOSInfo) {
    if (elements.macosVolumeSize) elements.macosVolumeSize.textContent = data.macOSInfo.totalSizeStr || '22.9 GB';
    if (elements.macosVersionText) elements.macosVersionText.textContent = data.macOSInfo.version || 'macOS';

    if (elements.macosSubvolumesList) {
      elements.macosSubvolumesList.innerHTML = '';
      const vols = data.macOSInfo.volumes || [];
      if (vols.length > 0) {
        vols.forEach(v => {
          const pill = document.createElement('span');
          pill.className = 'subvol-pill';
          pill.textContent = `${v.name} (${v.role}): ${v.sizeStr}`;
          elements.macosSubvolumesList.appendChild(pill);
        });
      } else {
        elements.macosSubvolumesList.innerHTML = `
          <span class="subvol-pill">System: ~12.6 GB</span>
          <span class="subvol-pill">Preboot: ~7.8 GB</span>
          <span class="subvol-pill">Recovery: ~1.4 GB</span>
          <span class="subvol-pill">VM: ~1.1 GB</span>
        `;
      }
    }
  }

  // 2. Sistem Verileri Hero
  if (elements.sysdataTotalSize) {
    elements.sysdataTotalSize.textContent = data.totalSystemDataStr || '0 B';
  }
  if (elements.sysdataSafeCleanable) {
    elements.sysdataSafeCleanable.textContent = data.safeCleanableStr || '0 B';
  }
  if (elements.vmSleepimagePill) {
    elements.vmSleepimagePill.textContent = data.vmSleepimageStr ? `Sleepimage: ${data.vmSleepimageStr}` : 'Sleepimage: ~2.15 GB';
  }
  if (elements.btnCleanSafeSystemData) {
    elements.btnCleanSafeSystemData.disabled = !(data.safeCleanableSize && data.safeCleanableSize > 0);
  }

  // 3. Render Categorized System Data
  renderSystemDataCategories(data.systemDataCategories || []);

  // 4. APFS Snapshots
  const snaps = data.snapshots || [];
  if (elements.appleSnapshotsCount) elements.appleSnapshotsCount.textContent = `${snaps.length} Görüntü`;
  if (elements.btnCleanAllSnapshots) elements.btnCleanAllSnapshots.disabled = snaps.length === 0;

  if (snaps.length === 0) {
    elements.appleSnapshotsList.innerHTML = '<div class="apple-item-row" style="color:var(--accent-emerald);">Yerel APFS anlık görüntüsü bulunmuyor (Temiz).</div>';
  } else {
    elements.appleSnapshotsList.innerHTML = '';
    snaps.forEach(snap => {
      const row = document.createElement('div');
      row.className = 'apple-item-row';
      row.innerHTML = `
        <div style="flex:1;">
          <div style="font-weight:600; font-size:0.85rem;">${escapeHtml(snap.snapshotId)}</div>
          <div style="font-size:0.75rem; color:var(--text-dim);">${escapeHtml(snap.dateStr)}</div>
        </div>
        <button class="action-btn action-btn-danger btn-snap-del" data-id="${escapeHtml(snap.snapshotId)}">Sil</button>
      `;
      row.querySelector('.btn-snap-del').addEventListener('click', () => {
        deleteSingleSnapshot(snap.snapshotId);
      });
      elements.appleSnapshotsList.appendChild(row);
    });
  }

  // 5. iOS Backups
  const backups = data.backups || [];
  if (elements.appleBackupsSize) elements.appleBackupsSize.textContent = data.totalBackupStr || '0 GB';

  if (backups.length === 0) {
    elements.appleBackupsList.innerHTML = '<div class="apple-item-row" style="color:var(--text-dim);">Hiçbir yerel iPhone/iPad aygıt yedeklemesi bulunmuyor.</div>';
  } else {
    elements.appleBackupsList.innerHTML = '';
    backups.forEach(bk => {
      const row = document.createElement('div');
      row.className = 'apple-item-row';
      row.innerHTML = `
        <div style="flex:1;">
          <div style="font-weight:700; font-size:0.88rem; color:#fff;">${escapeHtml(bk.deviceName || 'iOS Cihazı')} <small style="color:var(--text-dim);">(${escapeHtml(bk.deviceModel)})</small></div>
          <div style="font-size:0.75rem; color:var(--text-dim); margin-top:2px;">Son Yedek: ${escapeHtml(bk.lastDate)} • ${bk.fileCount.toLocaleString()} dosya</div>
        </div>
        <div style="text-align:right; margin-right:12px;">
          <strong style="color:var(--accent-amber);">${escapeHtml(bk.sizeStr)}</strong>
        </div>
        <button class="action-btn action-btn-danger btn-backup-del">Sil</button>
      `;
      row.querySelector('.btn-backup-del').addEventListener('click', () => {
        deleteBackupItem(bk.path, bk.sizeStr, bk.deviceName);
      });
      elements.appleBackupsList.appendChild(row);
    });
  }

  // 6. Simulators
  const sims = data.simulators || [];
  if (elements.appleSimulatorsSize) elements.appleSimulatorsSize.textContent = data.totalSimStr || '0 MB';

  if (sims.length === 0) {
    elements.appleSimulatorsList.innerHTML = '<div class="apple-item-row" style="color:var(--accent-emerald);">Kullanılmayan simülatör veya runtime önbelleği yok.</div>';
  } else {
    elements.appleSimulatorsList.innerHTML = '';
    sims.forEach(sim => {
      const row = document.createElement('div');
      row.className = 'apple-item-row';
      row.innerHTML = `
        <div style="flex:1;">
          <div style="font-weight:600; font-size:0.85rem;">${escapeHtml(sim.name)}</div>
          <div style="font-size:0.75rem; color:var(--text-dim);">${escapeHtml(sim.runtime || sim.path)}</div>
        </div>
        <div style="text-align:right; margin-right:8px;">
          <strong>${escapeHtml(sim.sizeStr)}</strong>
        </div>
      `;
      elements.appleSimulatorsList.appendChild(row);
    });
  }
}

function renderSystemDataCategories(categories) {
  if (!elements.appleCategoriesContainer) return;
  if (!categories || categories.length === 0) {
    elements.appleCategoriesContainer.innerHTML = '<div class="apple-item-row" style="color:var(--accent-emerald); padding: 16px;">Sistem verilerinde temizlenecek önbellek veya ek bileşen tespit edilmedi.</div>';
    return;
  }

  elements.appleCategoriesContainer.innerHTML = '';

  categories.forEach(cat => {
    const card = document.createElement('div');
    card.className = 'sysdata-cat-card';

    let catIcon = '📁';
    if (cat.id === 'package_managers') catIcon = '🍺';
    else if (cat.id === 'app_support') catIcon = '🎮';
    else if (cat.id === 'developer') catIcon = '🛠️';
    else if (cat.id === 'containers') catIcon = '📦';
    else if (cat.id === 'system_frameworks') catIcon = '⚙️';
    else if (cat.id === 'logs') catIcon = '📋';

    const header = document.createElement('div');
    header.className = 'sysdata-cat-header';
    header.innerHTML = `
      <div class="sysdata-cat-title-wrap">
        <span style="font-size:1.35rem; line-height:1; display:flex; align-items:center;">${catIcon}</span>
        <div>
          <h4>${escapeHtml(cat.title)}</h4>
          <div class="sysdata-cat-desc">${escapeHtml(cat.description)}</div>
        </div>
      </div>
      <div style="display:flex; align-items:center; gap:8px;">
        ${cat.safeSize > 0 ? `<span class="stat-pill" style="background:rgba(16,185,129,0.15); color:#34d399; font-size:0.75rem;">Güvenli: ${cat.safeSizeStr}</span>` : ''}
        <span class="sysdata-cat-size-badge">${escapeHtml(cat.totalSizeStr)}</span>
      </div>
    `;
    card.appendChild(header);

    const list = document.createElement('div');
    list.className = 'sysdata-items-list';

    (cat.items || []).forEach(item => {
      const row = document.createElement('div');
      row.className = 'sysdata-item-row';

      let riskBadge = '';
      if (item.risk === 'safe') {
        riskBadge = '<span class="risk-badge risk-safe">Güvenli</span>';
      } else if (item.risk === 'recommended') {
        riskBadge = '<span class="risk-badge risk-recommended">Önerilen</span>';
      } else {
        riskBadge = '<span class="risk-badge risk-caution">Dikkat</span>';
      }

      const isBrewCleanup = item.id === 'sysdata-homebrew-cleanup';
      const cleanBtnLabel = isBrewCleanup ? 'Brew Temizle' : 'Temizle';
      const cleanBtnClass = isBrewCleanup ? 'action-btn action-btn-primary btn-sysdata-clean' : 'action-btn action-btn-danger btn-sysdata-clean';

      row.innerHTML = `
        <div class="sysdata-item-meta">
          <div class="sysdata-item-name">
            <span>${escapeHtml(item.name)}</span>
            ${riskBadge}
          </div>
          <div class="sysdata-item-desc">${escapeHtml(item.description)}</div>
          <div class="sysdata-item-path" title="${escapeHtml(item.path)}">${escapeHtml(item.path)}</div>
        </div>
        <div class="sysdata-item-size">${escapeHtml(item.sizeStr)}</div>
        <div class="sysdata-item-actions">
          <button class="action-btn btn-sysdata-reveal" title="Finder'da Göster">Finder</button>
          ${item.cleanable ? `<button class="${cleanBtnClass}">${cleanBtnLabel}</button>` : ''}
        </div>
      `;

      row.querySelector('.btn-sysdata-reveal').addEventListener('click', () => {
        revealInFinder(item.path);
      });

      const cleanBtn = row.querySelector('.btn-sysdata-clean');
      if (cleanBtn) {
        cleanBtn.addEventListener('click', () => {
          cleanSingleSystemDataItem(item);
        });
      }

      list.appendChild(row);
    });

    card.appendChild(list);
    elements.appleCategoriesContainer.appendChild(card);
  });
}

function cleanSingleSystemDataItem(item) {
  const isTrash = elements.toggleTrashMode ? elements.toggleTrashMode.checked : false;
  const isBrewCleanup = item.id === 'sysdata-homebrew-cleanup';

  elements.modalTitle.textContent = isBrewCleanup
    ? "Homebrew Cleanup Çalıştırılsın mı?"
    : `${item.name} Temizlensin mi?`;

  elements.modalMessage.textContent = isBrewCleanup
    ? "Homebrew tarafından artık gereksiz görülen eski formül sürümleri, bağımlılıklar ve indirme önbellekleri temizlenecektir.\n\nÇalıştırılacak komut: brew cleanup -s"
    : `${item.description}\n\nKonum: ${item.path}`;

  elements.modalFreedSize.textContent = item.sizeStr;
  elements.modalMethodText.textContent = isBrewCleanup
    ? "Homebrew Bakım & Temizlik Rutini"
    : (isTrash ? "Finder Çöp Kutusuna Taşı" : "Kalıcı Olarak Temizle");

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast(`${item.name} temizleniyor...`, 'info');
    try {
      const endpoint = isBrewCleanup ? '/api/apple/brew/cleanup' : '/api/apple/systemdata/clean';
      const body = isBrewCleanup ? {} : { id: item.id, path: item.path, useTrash: isTrash };
      const res = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Temizlenemedi');
      showToast(data.message || 'Öğe temizlendi', 'success');
      startAppleScan();
      fetchSystemStats();
    } catch (err) {
      console.error(err);
      showToast('Hata: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

function cleanAllSafeSystemData() {
  if (!state.appleData || !state.appleData.safeCleanableSize) return;

  elements.modalTitle.textContent = "Tüm Güvenli Sistem Verileri Temizlensin mi?";
  elements.modalMessage.textContent = "Tespit edilen tüm güvenli önbellekler, indirme staging dosyaları, container önbellekleri ve günlükler temizlenecektir. Kişisel dosyalarınız veya uygulama kayıtlarınız etkilenmez.";
  elements.modalFreedSize.textContent = state.appleData.safeCleanableStr || "0 B";
  elements.modalMethodText.textContent = "Güvenli Sistem Verisi Temizliği";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('Tüm güvenli sistem verileri temizleniyor...', 'info');
    try {
      const res = await fetch('/api/apple/systemdata/clean-safe', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Temizlenemedi');
      showToast(data.message || 'Güvenli sistem verileri temizlendi', 'success');
      startAppleScan();
      fetchSystemStats();
    } catch (err) {
      console.error(err);
      showToast('Hata: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

async function reclaimPurgeableSpace() {
  if (elements.btnReclaimPurgeable) elements.btnReclaimPurgeable.disabled = true;
  showToast('Boşaltılabilir alan ve sistem önbellekleri serbest bırakılıyor...', 'info');

  try {
    const res = await fetch('/api/apple/purgeable/reclaim', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'İşlem başarısız');
    showToast(data.message || 'Boşaltılabilir alan serbest bırakıldı.', 'success');
    fetchSystemStats();
  } catch (err) {
    console.error(err);
    showToast('Boşaltılabilir alan hatası: ' + err.message, 'error');
  } finally {
    if (elements.btnReclaimPurgeable) elements.btnReclaimPurgeable.disabled = false;
  }
}

function deleteSingleSnapshot(snapId) {
  elements.modalTitle.textContent = "APFS Anlık Görüntüsü Silinsin mi?";
  elements.modalMessage.textContent = `${snapId} anlık görüntüsü silinecektir.`;
  elements.modalFreedSize.textContent = "Disk Alanı";
  elements.modalMethodText.textContent = "tmutil deletelocalsnapshots";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('APFS anlık görüntüsü siliniyor...', 'info');
    try {
      const res = await fetch('/api/apple/snapshots/delete', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ snapshotId: snapId }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Silinemedi');
      showToast(data.message || 'Anlık görüntü silindi', 'success');
      startAppleScan();
      fetchSystemStats();
    } catch (err) {
      console.error(err);
      showToast('Hata: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

function deleteAllSnapshots() {
  elements.modalTitle.textContent = "Tüm APFS Anlık Görüntüleri Silinsin mi?";
  elements.modalMessage.textContent = "Time Machine'in diskinizde oluşturduğu tüm yerel anlık görüntüler silinecek ve kullanılan disk alanı anında serbest kalacaktır.";
  elements.modalFreedSize.textContent = "Tüm Anlık Görüntüler";
  elements.modalMethodText.textContent = "APFS Snapshot Temizliği";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('Tüm anlık görüntüler siliniyor...', 'info');
    try {
      const res = await fetch('/api/apple/snapshots/delete', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ snapshotId: 'all' }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Silinemedi');
      showToast(data.message || 'Tüm anlık görüntüler silindi', 'success');
      startAppleScan();
      fetchSystemStats();
    } catch (err) {
      console.error(err);
      showToast('Hata: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

function deleteBackupItem(path, sizeStr, name) {
  const isTrash = elements.toggleTrashMode.checked;
  elements.modalTitle.textContent = `Aygıt Yedeğini Sil: ${name || 'iPhone/iPad'}`;
  elements.modalMessage.textContent = "Bu cihaza ait yerel yedekleme silinecektir. Gerekirse iCloud yedeklemeniz olduğunu doğrulayın.";
  elements.modalFreedSize.textContent = sizeStr;
  elements.modalMethodText.textContent = isTrash ? "Finder Çöp Kutusuna Taşı" : "Kalıcı Olarak Sil";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('Aygıt yedeği siliniyor...', 'info');
    try {
      const res = await fetch('/api/apple/backups/delete', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ path: path, useTrash: isTrash }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Silinemedi');
      showToast('Yedek başarıyla silindi', 'success');
      startAppleScan();
      fetchSystemStats();
    } catch (err) {
      console.error(err);
      showToast('Hata: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

function cleanSimulators() {
  elements.modalTitle.textContent = "Kullanılmayan Simülatörleri Temizle";
  elements.modalMessage.textContent = "Kullanılmayan iOS simülatör aygıtları, eski runtime önbellekleri ve CoreSimulator artıkları xcrun simctl ile temizlenecektir.";
  elements.modalFreedSize.textContent = elements.appleSimulatorsSize ? elements.appleSimulatorsSize.textContent : "Simülatör Alanı";
  elements.modalMethodText.textContent = "simctl delete unavailable";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('Simülatörler temizleniyor...', 'info');
    try {
      const res = await fetch('/api/apple/simulators/clean', {
        method: 'POST',
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Temizlenemedi');
      showToast(data.message || 'Simülatörler temizlendi', 'success');
      startAppleScan();
      fetchSystemStats();
    } catch (err) {
      console.error(err);
      showToast('Hata: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

// ==========================================================================
// Media & Messages Attachments Module
// ==========================================================================
async function startMediaScan() {
  if (!elements.btnScanMedia) return;
  elements.btnScanMedia.disabled = true;
  elements.mediaItemsList.innerHTML = '<div class="loading-state">iMessage ve Apple Mail ekleri taranıyor...</div>';
  state.selectedMedia.clear();
  updateMediaSelectedSize();

  try {
    const res = await fetch('/api/media/attachments');
    if (!res.ok) throw new Error('Mesaj ekleri taranamadı');
    const data = await res.json();
    state.mediaData = data;

    if (elements.mediaVideoSize) elements.mediaVideoSize.textContent = data.videoSizeStr || '0 MB';
    if (elements.mediaImageSize) elements.mediaImageSize.textContent = data.imageSizeStr || '0 MB';
    if (elements.mediaMailSize) elements.mediaMailSize.textContent = data.mailSizeStr || '0 MB';

    renderMediaList();
    showToast(`${data.totalCount || 0} ek bulundu (${data.totalSizeStr || '0 MB'})`, 'success');
  } catch (err) {
    console.error(err);
    elements.mediaItemsList.innerHTML = `<div class="loading-state" style="color:var(--accent-rose);">${escapeHtml(err.message)}</div>`;
    showToast('Ekleri tarama hatası: ' + err.message, 'error');
  } finally {
    elements.btnScanMedia.disabled = false;
  }
}

function renderMediaList() {
  if (!state.mediaData || !state.mediaData.items) return;

  const items = state.mediaData.items;
  if (items.length === 0) {
    elements.mediaItemsList.innerHTML = `
      <div class="empty-state">
        <div class="empty-wrap">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
          <p>Harika! Mesajlaşma veya posta kutunuzda yer kaplayan büyük ek bulunamadı.</p>
        </div>
      </div>
    `;
    return;
  }

  elements.mediaItemsList.innerHTML = '';
  items.forEach(item => {
    const el = document.createElement('div');
    el.className = 'target-item';

    const isChecked = state.selectedMedia.has(item.id);

    el.innerHTML = `
      <div class="target-left">
        <input type="checkbox" class="target-checkbox chk-media-item" data-id="${escapeHtml(item.id)}" data-size="${item.size}" ${isChecked ? 'checked' : ''}>
        <div class="target-details">
          <h4>
            ${escapeHtml(item.name)}
            <span class="risk-badge risk-safe">${escapeHtml(item.category)}</span>
            <span class="app-badge badge-user">${escapeHtml(item.source)}</span>
          </h4>
          <p>Tarih: ${escapeHtml(item.modDate)} (${item.daysOld} gün önce)</p>
          <div class="target-path">${escapeHtml(item.path)}</div>
        </div>
      </div>
      <div class="target-right">
        <div>
          <div class="target-size">${item.sizeStr}</div>
        </div>
        <button class="action-btn btn-media-reveal" data-path="${escapeHtml(item.path)}">Finder</button>
        <button class="action-btn action-btn-danger btn-media-delete">Sil</button>
      </div>
    `;

    const chk = el.querySelector('.chk-media-item');
    chk.addEventListener('change', (e) => {
      if (e.target.checked) {
        state.selectedMedia.add(item.id);
      } else {
        state.selectedMedia.delete(item.id);
      }
      updateMediaSelectedSize();
    });

    el.querySelector('.btn-media-reveal').addEventListener('click', () => {
      revealInFinder(item.path);
    });

    el.querySelector('.btn-media-delete').addEventListener('click', () => {
      confirmCleanMedia([item.id], item.sizeStr, item.name);
    });

    elements.mediaItemsList.appendChild(el);
  });
}

function updateMediaSelectedSize() {
  let size = 0;
  if (state.mediaData && state.mediaData.items) {
    state.mediaData.items.forEach(i => {
      if (state.selectedMedia.has(i.id)) {
        size += i.size;
      }
    });
  }

  if (elements.mediaSelectedSize) elements.mediaSelectedSize.textContent = formatBytes(size);
  if (elements.btnCleanMedia) elements.btnCleanMedia.disabled = state.selectedMedia.size === 0;
}

function confirmCleanMedia(ids, sizeStr, title) {
  ids = ids || Array.from(state.selectedMedia);
  if (ids.length === 0) return;

  const isTrash = elements.toggleTrashMode.checked;
  elements.modalTitle.textContent = title ? `Eki Sil: ${title}` : `${ids.length} Adet Eki Sil`;
  elements.modalMessage.textContent = "Seçilen mesaj ve posta ekleri silinecektir. Devam etmek istiyor musunuz?";
  elements.modalFreedSize.textContent = sizeStr || elements.mediaSelectedSize.textContent;
  elements.modalMethodText.textContent = isTrash ? "Finder Çöp Kutusuna Taşı" : "Kalıcı Olarak Sil";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('Ekler siliniyor...', 'info');
    try {
      const res = await fetch('/api/media/attachments/clean', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ids: ids, useTrash: isTrash }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Silinemedi');
      showToast(data.message || 'Ekler temizlendi', 'success');
      startMediaScan();
      fetchSystemStats();
    } catch (err) {
      console.error(err);
      showToast('Hata: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

// ==========================================================================
// Smart Care & Disk Health Module
// ==========================================================================
async function loadSmartCare() {
  try {
    const res = await fetch('/api/smartcare');
    if (!res.ok) return;
    const data = await res.json();
    state.smartCareData = data;

    // elements.quickCleanSize is managed by the comprehensive global scan (renderScanResults)
    // to include full developer/system targets (e.g. /cores, VS Code duplicates, Xcode, etc.)

    if (elements.storageThresholdBanner) {
      if (data.isCritical) {
        elements.storageThresholdBanner.style.display = 'flex';
        if (elements.alertBannerTitle) {
          elements.alertBannerTitle.textContent = `⚠️ Kritik Disk Doluluğu: %${Math.round(data.diskUsedPct)} Kullanılıyor!`;
        }
        if (elements.alertBannerDesc) {
          elements.alertBannerDesc.textContent = `${data.totalCleanableStr} boyutunda gereksiz dosya tek tıkla temizlenmeye hazır.`;
        }
      } else {
        elements.storageThresholdBanner.style.display = 'none';
      }
    }
  } catch (err) {
    console.error('Smart care check error:', err);
  }
}

function executeSmartCareClean() {
  const cleanableStr = (state.smartCareData && state.smartCareData.totalCleanableStr) || elements.quickCleanSize.textContent;
  elements.modalTitle.textContent = "Tek Tıkla Akıllı Sistem Bakımı";
  elements.modalMessage.textContent = "Tüm güvenli Xcode önbellekleri, derleyici/paket yöneticisi artıkları, sistem logları, geçici çöpler ve tarayıcı önbellekleri tek seferde güvenle temizlenecektir.";
  elements.modalFreedSize.textContent = cleanableStr;
  elements.modalMethodText.textContent = "Güvenli Otomatik Bakım";

  pendingCleanAction = async () => {
    elements.confirmModal.close();
    showToast('Akıllı bakım temizliği çalıştırılıyor...', 'info');
    try {
      const res = await fetch('/api/smartcare/clean', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({}),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Temizlik başarısız oldu');
      showToast(data.message || 'Akıllı bakım başarıyla tamamlandı!', 'success');
      loadSmartCare();
      fetchSystemStats();
      startGlobalScan();
    } catch (err) {
      console.error(err);
      showToast('Hata: ' + err.message, 'error');
    }
  };

  elements.confirmModal.showModal();
}

// ==========================================================================
// Authentication Module (config.json Password Protection)
// ==========================================================================
async function checkAuthStatus() {
  try {
    const res = await fetch('/api/auth/status');
    if (!res.ok) return false;
    const data = await res.json();
    state.authStatus = data;

    const authEnabled = (data.authEnabled !== undefined) ? data.authEnabled : data.enabled;

    if (authEnabled) {
      if (!data.authenticated) {
        if (elements.authModal && !elements.authModal.open) {
          elements.authModal.showModal();
        }
        if (elements.authStatusBar) elements.authStatusBar.style.display = 'none';
        return false;
      } else {
        if (data.token) {
          localStorage.setItem('dc_token', data.token);
        }
        if (elements.authModal && elements.authModal.open) {
          elements.authModal.close();
        }
        if (elements.authStatusBar) elements.authStatusBar.style.display = 'flex';
        return true;
      }
    } else {
      if (elements.authModal && elements.authModal.open) {
        elements.authModal.close();
      }
      if (elements.authStatusBar) elements.authStatusBar.style.display = 'none';
      return true;
    }
  } catch (err) {
    console.error('Auth check error:', err);
    return false;
  }
}

async function handleLoginSubmit() {
  const pwd = elements.authPasswordInput.value;
  if (!pwd) {
    elements.authErrorMsg.textContent = 'Lütfen şifrenizi girin.';
    elements.authErrorMsg.style.display = 'block';
    return;
  }
  elements.btnAuthSubmit.disabled = true;
  elements.authErrorMsg.style.display = 'none';

  try {
    const res = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password: pwd })
    });
    const data = await res.json();
    if (res.ok && data.success) {
      if (data.token) {
        localStorage.setItem('dc_token', data.token);
      }
      elements.authModal.close();
      elements.authPasswordInput.value = '';
      if (elements.authStatusBar) elements.authStatusBar.style.display = 'flex';
      showToast('Giriş başarılı! Oturum açıldı.', 'success');
      initializeDashboardData();
    } else {
      elements.authErrorMsg.textContent = data.error || 'Geçersiz şifre.';
      elements.authErrorMsg.style.display = 'block';
      elements.authPasswordInput.select();
    }
  } catch (err) {
    elements.authErrorMsg.textContent = 'Bağlantı hatası: ' + err.message;
    elements.authErrorMsg.style.display = 'block';
  } finally {
    elements.btnAuthSubmit.disabled = false;
  }
}

async function handleLogout() {
  try {
    localStorage.removeItem('dc_token');
    if (sseSource) {
      try { sseSource.close(); } catch (_) {}
      sseSource = null;
    }
    await fetch('/api/auth/logout', { method: 'POST' });
    showToast('Oturum kapatıldı.', 'info');
    checkAuthStatus();
  } catch (err) {
    console.error('Logout error:', err);
  }
}

// ==========================================================================
// Maintenance & Speedup Toolkit (CleanMyMac Inspired)
// ==========================================================================
async function loadMaintenanceTasks() {
  if (!elements.maintenanceTasksContainer) return;
  try {
    elements.maintenanceTasksContainer.innerHTML = '<div class="loading-state">Bakım görevleri taranıyor...</div>';
    const res = await fetch('/api/maintenance');
    if (!res.ok) throw new Error('Bakım görevleri listelenemedi');
    const tasks = await res.json();
    state.maintenanceTasks = tasks;
    renderMaintenanceTasks(tasks);
  } catch (err) {
    elements.maintenanceTasksContainer.innerHTML = `<div class="empty-state">Hata: ${escapeHtml(err.message)}</div>`;
  }
}

function renderMaintenanceTasks(tasks) {
  if (!elements.maintenanceTasksContainer) return;
  if (!tasks || tasks.length === 0) {
    elements.maintenanceTasksContainer.innerHTML = '<div class="empty-state">Bakım görevi bulunamadı.</div>';
    return;
  }

  const iconInfo = {
    'free-ram': { cls: 'icon-ram', text: 'RAM' },
    'flush-dns': { cls: 'icon-dns', text: 'DNS' },
    'speed-mail': { cls: 'icon-mail', text: 'MAIL' },
    'rebuild-launchservices': { cls: 'icon-launchservices', text: 'LS' },
    'reindex-spotlight': { cls: 'icon-spotlight', text: 'SPOT' },
    'run-periodic': { cls: 'icon-periodic', text: 'SCRP' },
    'verify-disk': { cls: 'icon-verify', text: 'DISK' }
  };

  elements.maintenanceTasksContainer.innerHTML = tasks.map(task => {
    const info = iconInfo[task.id] || { cls: 'icon-' + (task.icon || 'verify'), text: (task.icon || 'SYS').toUpperCase() };

    let statusClass = '';
    let statusText = 'Hazır';
    if (task.status === 'running') {
      statusClass = 'status-running';
      statusText = 'Çalışıyor...';
    } else if (task.status === 'done') {
      statusClass = 'status-done';
      statusText = '✓ Tamamlandı';
    }

    return `
      <div class="glass-card maintenance-task-card" id="mcard-${escapeHtml(task.id)}">
        <div class="m-task-left">
          <div class="m-task-icon ${info.cls}">
            ${info.text}
          </div>
          <div class="m-task-details">
            <h4>${escapeHtml(task.title || task.name)}</h4>
            <p>${escapeHtml(task.description)}</p>
          </div>
        </div>
        <div class="m-task-right">
          <div class="mcard-status-slot">
            <span class="task-status-pill ${statusClass}">${statusText}</span>
          </div>
          <button class="btn btn-secondary btn-sm btn-run-single-task" data-id="${escapeHtml(task.id)}">
            <span>Çalıştır</span>
          </button>
        </div>
      </div>
    `;
  }).join('');

  elements.maintenanceTasksContainer.querySelectorAll('.btn-run-single-task').forEach(btn => {
    btn.addEventListener('click', () => {
      runSingleMaintenanceTask(btn.dataset.id);
    });
  });
}

async function runSingleMaintenanceTask(taskId) {
  const card = document.getElementById(`mcard-${taskId}`);
  const btn = card ? card.querySelector('.btn-run-single-task') : null;
  const statusSlot = card ? card.querySelector('.mcard-status-slot') : null;

  if (btn) btn.disabled = true;
  if (statusSlot) statusSlot.innerHTML = '<span class="task-status-pill status-running">Çalışıyor...</span>';

  try {
    const res = await fetch('/api/maintenance/run', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ taskId: taskId, id: taskId })
    });
    const result = await res.json();
    const taskTitle = result.title || taskId;

    if (res.ok && (result.success || result.status === 'done')) {
      if (statusSlot) statusSlot.innerHTML = '<span class="task-status-pill status-done">✓ Tamamlandı</span>';
      showToast(`${taskTitle} başarıyla tamamlandı`, 'success');
    } else {
      if (statusSlot) statusSlot.innerHTML = '<span class="task-status-pill" style="color:var(--accent-rose);background:rgba(244,63,94,0.15)">✕ Hata</span>';
      showToast(`${taskTitle} sırasında sorun: ${result.error || 'Başarısız'}`, 'error');
    }

    if (elements.maintenanceLogContent && elements.maintenanceLogBox) {
      elements.maintenanceLogBox.style.display = 'block';
      const timestamp = new Date().toLocaleTimeString();
      elements.maintenanceLogContent.textContent += `[${timestamp}] ${taskTitle} (${result.duration || '0s'}):\n${result.output || '(Tamamlandı, ek çıktı yok)'}\n\n`;
      elements.maintenanceLogContent.scrollTop = elements.maintenanceLogContent.scrollHeight;
    }
  } catch (err) {
    if (statusSlot) statusSlot.innerHTML = '<span class="task-status-pill" style="color:var(--accent-rose);background:rgba(244,63,94,0.15)">✕ Hata</span>';
    showToast(`Hata: ${err.message}`, 'error');
  } finally {
    if (btn) btn.disabled = false;
  }
}

async function runAllMaintenanceTasks() {
  if (!state.maintenanceTasks || state.maintenanceTasks.length === 0) {
    await loadMaintenanceTasks();
  }
  if (!state.maintenanceTasks || state.maintenanceTasks.length === 0) return;

  elements.btnRunAllMaintenance.disabled = true;
  showToast('Tüm bakım görevleri sırayla çalıştırılıyor...', 'info');

  for (const task of state.maintenanceTasks) {
    await runSingleMaintenanceTask(task.id);
  }

  elements.btnRunAllMaintenance.disabled = false;
  showToast('Tüm bakım rutinleri tamamlandı!', 'success');
}

// ==========================================================================
// Secure File Shredder (DoD 5220.22-M Multi-Pass)
// ==========================================================================
async function executeShred() {
  const targetPath = elements.inputShredPath.value.trim();
  const passes = parseInt(elements.selectShredPasses.value, 10) || 3;

  if (!targetPath) {
    showToast('Lütfen kalıcı imha edilecek dosya veya klasör yolunu girin.', 'error');
    return;
  }

  const confirmMsg = `DİKKAT: Kalıcı ve Kurtarılamaz İmha!\n\n` +
    `"${targetPath}" yolu ve içeriğindeki tüm dosyalar ${passes} geçişli güvenli öğütücü (DoD standardı) ile ezilerek yok edilecektir.\n\n` +
    `Bu işlem GERİ ALINAMAZ. Onaylıyor musunuz?`;

  if (!confirm(confirmMsg)) return;

  elements.btnExecuteShred.disabled = true;
  elements.btnExecuteShred.innerHTML = '<span>Öğütülüyor, Lütfen Bekleyin...</span>';
  if (elements.shredResultBox) elements.shredResultBox.style.display = 'none';

  try {
    const res = await fetch('/api/shred', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: targetPath, passes: passes })
    });
    const result = await res.json();

    if (res.ok && !result.error) {
      const fileCount = result.totalFiles ?? result.shreddedFiles ?? 1;
      const sizeStr = result.totalStr ?? result.freedStr ?? '';
      showToast(`Güvenli imha tamamlandı! ${fileCount} dosya ezildi.`, 'success');
      elements.inputShredPath.value = '';
      if (elements.shredResultBox) {
        elements.shredResultBox.style.display = 'block';
        elements.shredResultBox.innerHTML = `
          <div class="shred-success-banner">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            <div>
              <h4>Güvenli İmha Başarıyla Tamamlandı</h4>
              <p><strong>${fileCount} adet dosya</strong> (${sizeStr}) DoD standardında ${passes} geçiş yapılarak rastgele bitlerle ezildi, sıfırlandı ve diskten kalıcı olarak silindi (${result.duration || '0s'}).</p>
            </div>
          </div>
        `;
      }
      fetchSystemStats();
    } else {
      showToast(`İmha hatası: ${result.error || 'İşlem başarısız'}`, 'error');
      if (elements.shredResultBox) {
        elements.shredResultBox.style.display = 'block';
        elements.shredResultBox.innerHTML = `
          <div class="shred-error-banner">
            <strong>Hata:</strong> ${escapeHtml(result.error || 'Bilinmeyen hata')}
          </div>
        `;
      }
    }
  } catch (err) {
    showToast(`İmha bağlantı hatası: ${err.message}`, 'error');
  } finally {
    elements.btnExecuteShred.disabled = false;
    elements.btnExecuteShred.innerHTML = `
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"/></svg>
      <span>Dosyayı Güvenle Öğüt ve Yok Et</span>
    `;
  }
}

// ==========================================================================
// Hardware & Resource Monitor (Sensei Inspired)
// ==========================================================================
async function loadHardwareMonitor() {
  if (!elements.monCpuPct) return;
  try {
    const res = await fetch('/api/monitor');
    if (!res.ok) return;
    const data = await res.json();
    state.monitorData = data;

    // CPU Info
    const cpuModel = data.cpuModel || (data.cpu && data.cpu.model) || 'Apple Silicon';
    const cpuPct = data.cpuUsagePercent ?? (data.cpu && data.cpu.usagePercent) ?? 0;
    const cpuCores = data.cpuCores || (data.cpu && data.cpu.cores) || 8;

    if (elements.monCpuModel) elements.monCpuModel.textContent = cpuModel;
    if (elements.monCpuPct) elements.monCpuPct.textContent = `${Math.round(cpuPct)}%`;
    if (elements.monCpuCores) elements.monCpuCores.textContent = `${cpuCores} Çekirdek`;

    // Memory Info
    const mem = data.memory || data.ram;
    if (mem) {
      if (elements.monRamTotal) elements.monRamTotal.textContent = `${mem.totalStr} Toplam`;
      if (elements.monRamPct) elements.monRamPct.textContent = `${Math.round(mem.usedPercent)}%`;
      if (elements.monRamUsed) elements.monRamUsed.textContent = mem.usedStr;
      if (elements.monRamFree) elements.monRamFree.textContent = mem.freeStr;
      if (elements.monRamWired) elements.monRamWired.textContent = mem.wiredStr || '--';
      if (elements.monRamCompressed) elements.monRamCompressed.textContent = mem.compressedStr || '--';
    }

    // OS & Uptime
    const hostname = data.hostname || (data.os && data.os.hostname) || 'Mac';
    const osVer = data.osVersion || (data.os && data.os.os) || 'macOS';
    const uptime = data.uptimeStr || (data.os && data.os.uptime) || '--';

    if (elements.monHostname) elements.monHostname.textContent = hostname;
    if (elements.monOsVer) elements.monOsVer.textContent = osVer;
    if (elements.monUptime) elements.monUptime.textContent = uptime;

    // Battery Info
    const batt = data.battery;
    if (batt && batt.hasBattery) {
      if (elements.monBatteryCard) elements.monBatteryCard.style.display = 'block';
      if (elements.monBattPct) elements.monBattPct.textContent = `${batt.percentage}%`;
      if (elements.monBattChargingStatus) {
        elements.monBattChargingStatus.textContent = batt.isCharging ? (batt.fullyCharged ? 'Tam Dolu' : 'Şarj Ediliyor') : (batt.remainingTime || 'Pilden Çalışıyor');
      }
      if (elements.monBattSource) elements.monBattSource.textContent = batt.powerSource || 'Şebeke Gücü (AC)';
      if (elements.monBattHealth) elements.monBattHealth.textContent = `${batt.healthPercent}%`;
      if (elements.monBattCycles) elements.monBattCycles.textContent = `${batt.cycleCount}`;
      if (elements.monBattTemp) elements.monBattTemp.textContent = `${batt.temperatureCelsius ? batt.temperatureCelsius.toFixed(1) : '--'} °C`;
      if (elements.monBattCondition) {
        elements.monBattCondition.textContent = batt.condition || 'Normal';
        elements.monBattCondition.style.color = batt.condition && batt.condition.includes('Servis') ? 'var(--accent-rose)' : 'var(--accent-emerald)';
      }
    } else if (elements.monBatteryCard) {
      if (elements.monBattPct) elements.monBattPct.textContent = 'Masaüstü';
      if (elements.monBattChargingStatus) elements.monBattChargingStatus.textContent = 'Harici Güç';
      if (elements.monBattSource) elements.monBattSource.textContent = 'Masaüstü Mac Gücü';
      if (elements.monBattHealth) elements.monBattHealth.textContent = 'N/A';
      if (elements.monBattCycles) elements.monBattCycles.textContent = 'N/A';
      if (elements.monBattTemp) elements.monBattTemp.textContent = '--';
      if (elements.monBattCondition) elements.monBattCondition.textContent = 'Normal';
    }

    // GPU Info
    const gpu = data.gpu;
    if (gpu) {
      if (elements.monGpuModel) elements.monGpuModel.textContent = gpu.model || 'Apple Silicon GPU';
      if (elements.monGpuCores) elements.monGpuCores.textContent = `${gpu.cores || 8}`;
      if (elements.monGpuMetal) elements.monGpuMetal.textContent = gpu.metalSupport || 'Metal Destekli';
      if (elements.monGpuRes) elements.monGpuRes.textContent = gpu.displayResolution || 'Retina Ekran';
      if (elements.monGpuVendor) elements.monGpuVendor.textContent = gpu.vendor || 'Apple';
    }

    // Disk Detail Info
    const disk = data.diskDetail || data.disk;
    if (disk) {
      if (elements.monDiskPct) elements.monDiskPct.textContent = `${Math.round(disk.usedPercent || 0)}%`;
      if (elements.monDiskUsed) elements.monDiskUsed.textContent = disk.usedStr || '--';
      if (elements.monDiskFree) elements.monDiskFree.textContent = disk.freeStr || '--';
      if (elements.monDiskSmart) elements.monDiskSmart.textContent = disk.smartStatus || 'Doğrulandı';
      if (elements.monDiskThroughput) elements.monDiskThroughput.textContent = disk.throughput || '0 MB/s';
      if (elements.monDiskType) elements.monDiskType.textContent = `${disk.solidState ? 'NVMe SSD' : 'Depolama'} (${disk.fileSystem || 'APFS'})`;
    }

    // Thermal & Temperature Metrics
    const thermal = data.thermal;
    if (thermal) {
      const cpuT = thermal.cpuTempCelsius || 34.8;
      const gpuT = thermal.gpuTempCelsius || Math.max(30.0, cpuT - 1.4);
      const battT = thermal.batteryTempCelsius || (batt && batt.temperatureCelsius) || 30.8;

      const formatT = (t) => `${t.toFixed(1)} °C`;
      const applyTempClass = (el, t) => {
        if (!el) return;
        el.classList.remove('temp-warm', 'temp-hot');
        if (t >= 80) el.classList.add('temp-hot');
        else if (t >= 65) el.classList.add('temp-warm');
      };

      if (elements.monCpuTemp) {
        elements.monCpuTemp.textContent = formatT(cpuT);
        applyTempClass(elements.monCpuTemp, cpuT);
      }
      if (elements.monGpuTemp) {
        elements.monGpuTemp.textContent = formatT(gpuT);
        applyTempClass(elements.monGpuTemp, gpuT);
      }
      if (elements.monThermalCpu) {
        elements.monThermalCpu.textContent = formatT(cpuT);
        applyTempClass(elements.monThermalCpu, cpuT);
      }
      if (elements.monThermalGpu) {
        elements.monThermalGpu.textContent = formatT(gpuT);
        applyTempClass(elements.monThermalGpu, gpuT);
      }
      if (elements.monThermalBatt) {
        elements.monThermalBatt.textContent = formatT(battT);
        applyTempClass(elements.monThermalBatt, battT);
      }
      if (elements.monThermalState) {
        elements.monThermalState.textContent = thermal.thermalState || 'Nominal (Serin & Kararlı)';
      }
      if (elements.monThermalFan) {
        elements.monThermalFan.textContent = thermal.fanSpeedRPM > 0 ? `${thermal.fanSpeedRPM} RPM` : '0 RPM (Sessiz)';
      }
    }
  } catch (err) {
    console.error('Hardware monitor fetch error:', err);
  }
}

// ==========================================================================
// Startup & Background Items Manager
// ==========================================================================
async function loadStartupItems(isSilent = false) {
  if (!isSilent && elements.startupItemsList) {
    elements.startupItemsList.innerHTML = '<div class="loading-state">Başlangıç öğeleri ve arka plan servisleri taranıyor...</div>';
  }

  try {
    const res = await fetch('/api/startup');
    if (!res.ok) throw new Error('Başlangıç öğeleri listelenemedi');
    const data = await res.json();
    state.startupData = data;

    // Update badges and counters
    if (elements.badgeStartupCount) {
      elements.badgeStartupCount.textContent = data.activeCount || data.totalCount || 0;
      elements.badgeStartupCount.style.display = data.totalCount > 0 ? 'inline-block' : 'none';
    }

    if (elements.startupStatTotal) elements.startupStatTotal.textContent = data.totalCount || 0;
    if (elements.startupStatActive) elements.startupStatActive.textContent = data.activeCount || 0;
    if (elements.startupStatRunning) elements.startupStatRunning.textContent = data.runningCount || 0;
    if (elements.startupStatLogin) elements.startupStatLogin.textContent = data.loginItemsCount || 0;

    if (elements.chipCountAll) elements.chipCountAll.textContent = data.totalCount || 0;
    if (elements.chipCountLogin) elements.chipCountLogin.textContent = data.loginItemsCount || 0;
    if (elements.chipCountUser) elements.chipCountUser.textContent = data.userAgentsCount || 0;
    if (elements.chipCountSysAgent) elements.chipCountSysAgent.textContent = data.systemAgentsCount || 0;
    if (elements.chipCountDaemon) elements.chipCountDaemon.textContent = data.daemonsCount || 0;

    renderStartupChart(data);
    renderStartupItems();
  } catch (err) {
    console.error('Startup items error:', err);
    if (!isSilent && elements.startupItemsList) {
      elements.startupItemsList.innerHTML = `<div class="empty-state">Hata: ${escapeHtml(err.message)}</div>`;
    }
  }
}

function renderStartupChart(data) {
  if (!elements.startupDonutChart || !data) return;
  const canvas = elements.startupDonutChart;
  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  const login = data.loginItemsCount || 0;
  const user = data.userAgentsCount || 0;
  const sys = data.systemAgentsCount || 0;
  const daemon = data.daemonsCount || 0;
  const total = data.totalCount || (login + user + sys + daemon);
  const active = data.activeCount || 0;

  if (elements.donutCenterTotal) elements.donutCenterTotal.textContent = total;
  if (elements.donutLegendLogin) elements.donutLegendLogin.textContent = login;
  if (elements.donutLegendUser) elements.donutLegendUser.textContent = user;
  if (elements.donutLegendSysAgent) elements.donutLegendSysAgent.textContent = sys;
  if (elements.donutLegendDaemon) elements.donutLegendDaemon.textContent = daemon;

  if (elements.startupRatioText && elements.startupRatioFill) {
    const ratio = total > 0 ? Math.round((active / total) * 100) : 0;
    elements.startupRatioText.textContent = `${ratio}% Aktif (${active} / ${total})`;
    elements.startupRatioFill.style.width = `${ratio}%`;
  }

  // High-DPI canvas scaling
  const dpr = window.devicePixelRatio || 1;
  const size = 180;
  if (canvas.width !== size * dpr) {
    canvas.width = size * dpr;
    canvas.height = size * dpr;
  }
  ctx.save();
  ctx.scale(dpr, dpr);
  ctx.clearRect(0, 0, size, size);

  const cx = size / 2;
  const cy = size / 2;
  const radius = 68;
  const thickness = 14;

  // Background ring
  ctx.beginPath();
  ctx.arc(cx, cy, radius, 0, Math.PI * 2);
  ctx.strokeStyle = 'rgba(255, 255, 255, 0.06)';
  ctx.lineWidth = thickness;
  ctx.stroke();

  if (total === 0) {
    ctx.restore();
    return;
  }

  const segments = [
    { count: login, color: '#06b6d4' },
    { count: user, color: '#6366f1' },
    { count: sys, color: '#f59e0b' },
    { count: daemon, color: '#10b981' }
  ].filter(s => s.count > 0);

  let startAngle = -Math.PI / 2;
  const gap = segments.length > 1 ? 0.05 : 0;

  segments.forEach(seg => {
    const sliceAngle = (seg.count / total) * (Math.PI * 2);
    const endAngle = startAngle + sliceAngle - gap;

    if (endAngle > startAngle) {
      ctx.beginPath();
      ctx.arc(cx, cy, radius, startAngle, endAngle);
      ctx.strokeStyle = seg.color;
      ctx.lineWidth = thickness;
      ctx.lineCap = 'round';
      ctx.stroke();
    }
    startAngle += sliceAngle;
  });

  ctx.restore();
}

function renderStartupItems() {
  if (!elements.startupItemsList || !state.startupData) return;

  let items = state.startupData.items || [];

  // Filter by category
  if (state.startupFilter !== 'all') {
    items = items.filter(item => item.type === state.startupFilter);
  }

  // Filter by search query
  if (state.startupSearch) {
    const q = state.startupSearch.toLowerCase();
    items = items.filter(item =>
      (item.name && item.name.toLowerCase().includes(q)) ||
      (item.label && item.label.toLowerCase().includes(q)) ||
      (item.vendor && item.vendor.toLowerCase().includes(q)) ||
      (item.program && item.program.toLowerCase().includes(q)) ||
      (item.path && item.path.toLowerCase().includes(q))
    );
  }

  if (items.length === 0) {
    elements.startupItemsList.innerHTML = `
      <div class="empty-state">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/></svg>
        <p>Seçili filtre veya aramaya uygun başlangıç öğesi bulunamadı.</p>
      </div>
    `;
    return;
  }

  const typeConfig = {
    login_item: { iconCls: 'type-login', tagCls: 's-tag-login', typeName: 'Oturum Açma' },
    user_agent: { iconCls: 'type-user', tagCls: 's-tag-user', typeName: 'Kullanıcı Ajanı' },
    system_agent: { iconCls: 'type-sysagent', tagCls: 's-tag-sysagent', typeName: 'Sistem Ajanı' },
    system_daemon: { iconCls: 'type-daemon', tagCls: 's-tag-daemon', typeName: 'Arka Plan Servisi' },
  };

  elements.startupItemsList.innerHTML = items.map(item => {
    const cfg = typeConfig[item.type] || { iconCls: 'type-user', tagCls: 's-tag-user', typeName: item.typeName };
    const disabledClass = !item.enabled ? 'is-disabled' : '';

    let runningBadge = '';
    if (item.running) {
      runningBadge = `<span class="s-tag s-tag-running">Çalışıyor ${item.pid > 0 ? `(PID: ${item.pid})` : ''}</span>`;
    }

    let statusTag = '';
    if (!item.enabled) {
      statusTag = `<span class="s-tag" style="background:rgba(244,63,94,0.15);color:var(--accent-rose);">Devre Dışı</span>`;
    }

    const pathToShow = item.path || item.program || '';

    return `
      <div class="glass-card startup-item-card ${disabledClass}" id="scard-${escapeHtml(item.id)}">
        <div class="s-item-left">
          <div class="s-item-icon ${cfg.iconCls}">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M15.59 14.37a6 6 0 01-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 006.16-12.12A14.98 14.98 0 009.631 8.41m5.96 5.96a14.926 14.926 0 01-5.841 2.58m-.119-8.54a6 6 0 00-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 00-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 01-2.448-2.448 14.9 14.9 0 01.06-.312m-2.24 2.39a4.493 4.493 0 00-1.757 4.306 4.493 4.493 0 004.306-1.758M16.5 9a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0z"/></svg>
          </div>
          <div class="s-item-details">
            <div class="s-item-title-row">
              <h4 title="${escapeHtml(item.name)}">${escapeHtml(item.name)}</h4>
              <span class="s-tag ${cfg.tagCls}">${cfg.typeName}</span>
              ${item.vendor ? `<span class="s-tag s-tag-vendor">${escapeHtml(item.vendor)}</span>` : ''}
              ${runningBadge}
              ${statusTag}
            </div>
            <div class="s-item-meta">
              <span class="s-item-path" title="${escapeHtml(pathToShow)}">${escapeHtml(pathToShow)}</span>
            </div>
          </div>
        </div>

        <div class="s-item-right">
          <!-- Toggle Switch -->
          <label class="toggle-switch" title="${item.enabled ? 'Devre Dışı Bırak' : 'Etkinleştir'}">
            <input type="checkbox" class="toggle-startup-switch" data-id="${escapeHtml(item.id)}" ${item.enabled ? 'checked' : ''} />
            <span class="toggle-slider"></span>
          </label>

          <!-- Finder Reveal -->
          ${item.path ? `
            <button class="btn btn-secondary btn-icon-only btn-reveal-startup" data-path="${escapeHtml(item.path)}" title="Finder'da Göster">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 003 8.25v10.5A2.25 2.25 0 005.25 21h10.5A2.25 2.25 0 0018 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25"/></svg>
            </button>
          ` : ''}

          <!-- Delete / Remove Button -->
          <button class="btn btn-secondary btn-icon-only btn-delete-startup" data-id="${escapeHtml(item.id)}" data-name="${escapeHtml(item.name)}" title="Başlangıçtan Kaldır">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"/></svg>
          </button>
        </div>
      </div>
    `;
  }).join('');

  // Attach event handlers
  elements.startupItemsList.querySelectorAll('.toggle-startup-switch').forEach(sw => {
    sw.addEventListener('change', (e) => {
      const id = sw.dataset.id;
      const enable = sw.checked;
      toggleStartupItem(id, enable);
    });
  });

  elements.startupItemsList.querySelectorAll('.btn-reveal-startup').forEach(btn => {
    btn.addEventListener('click', () => {
      revealInFinder(btn.dataset.path);
    });
  });

  elements.startupItemsList.querySelectorAll('.btn-delete-startup').forEach(btn => {
    btn.addEventListener('click', () => {
      deleteStartupItem(btn.dataset.id, btn.dataset.name);
    });
  });
}

async function toggleStartupItem(id, enabled) {
  try {
    const res = await fetch('/api/startup/toggle', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: id, enabled: enabled })
    });
    const result = await res.json();
    if (res.ok && result.success) {
      showToast(result.message || 'Öğe durumu güncellendi.', 'success');
      loadStartupItems(true);
    } else {
      showToast('İşlem başarısız: ' + (result.error || 'Bilinmeyen hata'), 'error');
      renderStartupItems(); // revert toggle visually
    }
  } catch (err) {
    showToast('Bağlantı hatası: ' + err.message, 'error');
    renderStartupItems();
  }
}

async function deleteStartupItem(id, name) {
  if (!confirm(`"${name}" başlangıç öğesi listeden kaldırılacaktır. Onaylıyor musunuz?`)) return;

  try {
    const res = await fetch('/api/startup/delete', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: id, useTrash: true })
    });
    const result = await res.json();
    if (res.ok && result.success) {
      showToast(`"${name}" başarıyla kaldırıldı.`, 'success');
      loadStartupItems(true);
    } else {
      showToast('Kaldırma hatası: ' + (result.error || 'Başarısız'), 'error');
    }
  } catch (err) {
    showToast('Hata: ' + err.message, 'error');
  }
}

async function handleAddLoginItemSubmit() {
  const path = elements.inputNewLoginPath.value.trim();
  const hidden = elements.checkNewLoginHidden.checked;

  if (!path) {
    showToast('Lütfen uygulamanın tam yolunu girin.', 'error');
    return;
  }

  try {
    const res = await fetch('/api/startup/add', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: path, hidden: hidden })
    });
    const result = await res.json();
    if (res.ok && result.success) {
      showToast(result.message || 'Yeni başlangıç öğesi eklendi!', 'success');
      elements.addLoginItemDialog.close();
      elements.inputNewLoginPath.value = '';
      elements.checkNewLoginHidden.checked = false;
      loadStartupItems();
    } else {
      showToast('Ekleme hatası: ' + (result.error || 'Başarısız'), 'error');
    }
  } catch (err) {
    showToast('Hata: ' + err.message, 'error');
  }
}

// ==========================================================================
// Extensions Manager
// ==========================================================================
async function loadExtensions(isSilent = false) {
  if (!isSilent && elements.extensionsItemsList) {
    elements.extensionsItemsList.innerHTML = '<div class="loading-state">Eklentiler taranıyor...</div>';
  }

  try {
    const res = await fetch('/api/extensions');
    if (!res.ok) throw new Error('Eklentiler alınamadı');
    const data = await res.json();
    extensionsData = data;

    // Update badge & stats
    if (elements.badgeExtensionsCount) elements.badgeExtensionsCount.textContent = data.totalCount;
    if (elements.extStatTotal) elements.extStatTotal.textContent = data.totalCount;
    if (elements.extStatActive) elements.extStatActive.textContent = data.activeCount;
    if (elements.extStatSysext) elements.extStatSysext.textContent = data.systemExtCount;
    if (elements.extStatAppex) elements.extStatAppex.textContent = data.pluginKitCount;

    // Filter pill counts
    if (elements.pillExtAll) elements.pillExtAll.textContent = data.totalCount;
    if (elements.pillExtSysext) elements.pillExtSysext.textContent = data.systemExtCount;
    if (elements.pillExtPluginkit) elements.pillExtPluginkit.textContent = data.pluginKitCount;
    if (elements.pillExtQuicklook) elements.pillExtQuicklook.textContent = data.quickLookCount;
    if (elements.pillExtSpotlight) elements.pillExtSpotlight.textContent = data.spotlightCount;
    if (elements.pillExtPrefpane) elements.pillExtPrefpane.textContent = data.prefPaneCount;

    renderExtensionsList();
  } catch (err) {
    console.error('Extensions scan error:', err);
    if (elements.extensionsItemsList) {
      elements.extensionsItemsList.innerHTML = `<div class="error-state">Hata: ${escapeHtml(err.message)}</div>`;
    }
  }
}

function renderExtensionsList() {
  if (!elements.extensionsItemsList || !extensionsData) return;

  let items = extensionsData.items || [];

  // Filter by category
  if (extCurrentFilter !== 'all') {
    switch (extCurrentFilter) {
      case 'sysext':
        items = items.filter(it => it.type === 'system_ext');
        break;
      case 'pluginkit':
        items = items.filter(it => it.type === 'pluginkit');
        break;
      case 'quicklook':
        items = items.filter(it => it.type === 'quicklook');
        break;
      case 'spotlight':
        items = items.filter(it => it.type === 'spotlight');
        break;
      case 'prefpane':
        items = items.filter(it => it.type === 'prefpane');
        break;
    }
  }

  // Filter by search query
  if (extSearchQuery) {
    const q = extSearchQuery.toLowerCase();
    items = items.filter(it =>
      (it.name && it.name.toLowerCase().includes(q)) ||
      (it.bundleId && it.bundleId.toLowerCase().includes(q)) ||
      (it.vendor && it.vendor.toLowerCase().includes(q)) ||
      (it.description && it.description.toLowerCase().includes(q))
    );
  }

  if (items.length === 0) {
    elements.extensionsItemsList.innerHTML = `
      <div class="empty-state" style="padding: 40px; text-align: center; color: var(--text-muted);">
        <p>Arama kriterine uygun eklenti bulunamadı.</p>
      </div>
    `;
    return;
  }

  let html = '';
  items.forEach(it => {
    let tagClass = 's-tag-pluginkit';
    let typeName = it.typeName || 'Eklenti';
    if (it.type === 'system_ext') tagClass = 's-tag-sysext';
    else if (it.type === 'quicklook') tagClass = 's-tag-quicklook';
    else if (it.type === 'spotlight') tagClass = 's-tag-spotlight';
    else if (it.type === 'prefpane') tagClass = 's-tag-prefpane';

    const isSystemExt = it.type === 'system_ext';

    html += `
      <div class="startup-item-card glass-card">
        <div class="s-item-left">
          <div class="s-item-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M14.25 6.087c0-.355.186-.676.401-.959A3.75 3.75 0 0010.5 2.25c-1.28 0-2.417.643-3.109 1.628-.215.283-.401.604-.401.959v2.163H4.5A2.25 2.25 0 002.25 9.25v2.25h2.163c.355 0 .676.186.959.401.985.692 1.628 1.829 1.628 3.109 0 1.28-.643 2.417-1.628 3.109-.283.215-.604.401-.959.401H2.25v2.25A2.25 2.25 0 004.5 22.5h2.25v-2.163c0-.355.186-.676.401-.959a3.75 3.75 0 016.218 0c.215.283.401.604.401.959V22.5h2.25a2.25 2.25 0 002.25-2.25v-2.25h-2.163c-.355 0-.676-.186-.959-.401a3.75 3.75 0 010-6.218c.283-.215.604-.401.959-.401H21.75V9.25A2.25 2.25 0 0019.5 7h-2.25V4.837z"/>
            </svg>
          </div>
          <div class="s-item-details">
            <div class="s-item-title-row">
              <h4>${escapeHtml(it.name || it.bundleId)}</h4>
              <span class="s-tag ${tagClass}">${escapeHtml(typeName)}</span>
              ${it.vendor ? `<span class="s-tag s-tag-vendor">${escapeHtml(it.vendor)}</span>` : ''}
              ${it.version ? `<span class="s-tag s-tag-vendor">v${escapeHtml(it.version)}</span>` : ''}
              ${it.active ? `<span class="s-tag s-tag-running">Kullanımda</span>` : ''}
              ${it.isOrphaned ? `<span class="s-tag" style="background: rgba(239, 68, 68, 0.15); color: #f87171; border: 1px solid rgba(239, 68, 68, 0.3);">Artık / Kaldırılabilir</span>` : ''}
            </div>
            <div class="s-item-meta">
              <span class="s-item-path" title="${escapeHtml(it.path || it.bundleId)}">${escapeHtml(it.path || it.bundleId)}</span>
              <span style="font-size:0.75rem; color:var(--text-muted);">${escapeHtml(it.description || '')}</span>
              ${it.hostAppPath ? `<div style="font-size:0.72rem; color:var(--accent-indigo); margin-top:2px;">Ana Uygulama: ${escapeHtml(it.hostAppPath)}</div>` : ''}
            </div>
          </div>
        </div>
        <div class="s-item-right">
          ${it.path ? `
            <button class="btn btn-secondary btn-xs btn-reveal-ext" data-path="${escapeHtml(it.path)}" title="Finder'da Göster">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:12px; height:12px;"><path d="M15 3h6v6M10 14L21 3M18 13v6a2 2 0 01-2 2H5a2 2 0 01-2-2V8a2 2 0 012-2h6"/></svg>
              <span>Bul</span>
            </button>
          ` : ''}
          ${isSystemExt ? `
            <button class="btn btn-secondary btn-xs btn-manage-sysext" 
                    data-id="${escapeHtml(it.id)}" 
                    data-name="${escapeHtml(it.name || it.bundleId)}" 
                    data-type="${escapeHtml(typeName)}" 
                    data-hostapp="${escapeHtml(it.hostAppPath || '')}" 
                    data-desc="${escapeHtml(it.description || '')}"
                    title="Sistem Ayarları'nda Yönet / Devre Dışı Bırak">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:12px; height:12px;"><path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
              <span>Yönet</span>
            </button>
            <label class="toggle-switch" title="${it.enabled ? 'Devre Dışı Bırak' : 'Etkinleştir'}">
              <input type="checkbox" class="check-toggle-ext" data-id="${escapeHtml(it.id)}" ${it.enabled ? 'checked' : ''} />
              <span class="toggle-slider"></span>
            </label>
          ` : `
            <label class="toggle-switch" title="${it.enabled ? 'Devre Dışı Bırak' : 'Etkinleştir'}">
              <input type="checkbox" class="check-toggle-ext" data-id="${escapeHtml(it.id)}" ${it.enabled ? 'checked' : ''} />
              <span class="toggle-slider"></span>
            </label>
          `}
        </div>
      </div>
    `;
  });

  elements.extensionsItemsList.innerHTML = html;

  // Wire toggles
  elements.extensionsItemsList.querySelectorAll('.check-toggle-ext').forEach(chk => {
    chk.addEventListener('change', (e) => {
      const id = e.target.dataset.id;
      const enabled = e.target.checked;
      toggleExtension(id, enabled);
    });
  });

  // Wire manage sysext
  elements.extensionsItemsList.querySelectorAll('.btn-manage-sysext').forEach(btn => {
    btn.addEventListener('click', () => {
      const id = btn.dataset.id;
      const name = btn.dataset.name;
      const type = btn.dataset.type;
      const hostApp = btn.dataset.hostapp;
      const desc = btn.dataset.desc;
      openSysextModal({ id, name, type, hostApp, desc });
    });
  });

  // Wire reveal
  elements.extensionsItemsList.querySelectorAll('.btn-reveal-ext').forEach(btn => {
    btn.addEventListener('click', () => {
      const p = btn.dataset.path;
      if (p) revealInFinder(p);
    });
  });
}

function openSysextModal(data) {
  if (!elements.sysextManageDialog) return;
  elements.sysextModalName.textContent = data.name || '--';
  elements.sysextModalType.textContent = data.type || '--';
  elements.sysextModalState.textContent = data.desc || '--';

  if (data.hostApp) {
    elements.sysextModalHostAppRow.style.display = 'flex';
    elements.sysextModalHostApp.textContent = data.hostApp;
    elements.btnRevealSysextApp.style.display = 'inline-block';
    elements.btnRevealSysextApp.onclick = () => revealInFinder(data.hostApp);
  } else {
    elements.sysextModalHostAppRow.style.display = 'none';
    elements.btnRevealSysextApp.style.display = 'none';
  }

  elements.btnOpenSysSettings.onclick = async () => {
    try {
      await fetch('/api/system/open-settings');
      showToast('macOS Sistem Ayarları açıldı.', 'success');
    } catch (e) {
      console.error(e);
    }
  };

  elements.btnCloseSysextModal.onclick = () => {
    elements.sysextManageDialog.close();
  };

  // Open settings right away
  fetch('/api/system/open-settings').catch(() => {});
  elements.sysextManageDialog.showModal();
}

async function toggleExtension(id, enabled) {
  try {
    const res = await fetch('/api/extensions/toggle', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: id, enabled: enabled })
    });
    const result = await res.json();
    if (res.ok && result.success) {
      showToast(result.message || 'Eklenti durumu güncellendi.', 'success');
      loadExtensions(true);

      if (id.startsWith('sysext:')) {
        const item = state.extensionsData && state.extensionsData.items 
          ? state.extensionsData.items.find(x => x.id === id) 
          : null;
        if (item) {
          openSysextModal({
            id: item.id,
            name: item.name || item.bundleId,
            type: item.typeName,
            hostApp: item.hostAppPath,
            desc: item.description
          });
        }
      }
    } else {
      showToast('İşlem başarısız: ' + (result.error || 'Bilinmeyen hata'), 'error');
      renderExtensionsList();
    }
  } catch (err) {
    showToast('Bağlantı hatası: ' + err.message, 'error');
    renderExtensionsList();
  }
}

// ==========================================================================
// Network Monitor
// ==========================================================================
async function loadNetworkStats(isSilent = false) {
  try {
    const res = await fetch('/api/network');
    if (!res.ok) return;
    const stats = await res.json();
    state.networkStats = stats;
    renderNetworkStats(stats);
  } catch (err) {
    console.error('Network stats error:', err);
  }
}

function renderNetworkStats(stats) {
  if (!stats) return;

  // Always update live traffic chart first
  renderNetworkChart(stats);

  if (elements.netStatStatus) elements.netStatStatus.textContent = stats.status || 'Bağlı';
  if (elements.netStatIface) elements.netStatIface.textContent = stats.interfaceName || 'en0';
  if (elements.netStatIp) elements.netStatIp.textContent = stats.ipv4Address || 'Bilinmiyor';
  if (elements.netStatGw) elements.netStatGw.textContent = `Ağ Geçidi: ${stats.gateway || '--'}`;
  if (elements.netStatSpeedIn) elements.netStatSpeedIn.textContent = stats.downloadSpeedStr || '0 KB/s';
  if (elements.netStatTotalIn) elements.netStatTotalIn.textContent = `Toplam Gelen: ${stats.bytesInFormatted || '0 B'}`;
  if (elements.netStatSpeedOut) elements.netStatSpeedOut.textContent = stats.uploadSpeedStr || '0 KB/s';
  if (elements.netStatTotalOut) elements.netStatTotalOut.textContent = `Toplam Giden: ${stats.bytesOutFormatted || '0 B'}`;

  if (elements.netDetailMac) elements.netDetailMac.textContent = stats.macAddress || '--';
  if (elements.netDetailDns) elements.netDetailDns.textContent = (stats.dns && stats.dns.length > 0) ? stats.dns.join(', ') : '--';
  if (elements.netDetailConnCount) elements.netDetailConnCount.textContent = stats.connectionCount || 0;
  if (elements.netConnBadge) elements.netConnBadge.textContent = `${stats.connectionCount || 0} Aktif Soket`;

  if (!elements.netConnectionsTbody) return;

  const conns = stats.activeConnections || [];
  if (conns.length === 0) {
    elements.netConnectionsTbody.innerHTML = `<tr><td colspan="7" class="empty-state">Aktif TCP bağlantısı bulunamadı.</td></tr>`;
    return;
  }

  let html = '';
  conns.forEach(c => {
    html += `
      <tr>
        <td><strong>${escapeHtml(c.command)}</strong></td>
        <td><code>${c.pid}</code></td>
        <td>${escapeHtml(c.user)}</td>
        <td>${escapeHtml(c.protocol)}</td>
        <td><code style="font-size:0.75rem;">${escapeHtml(c.localAddr)}</code></td>
        <td><code style="font-size:0.75rem; color:var(--accent-cyan);">${escapeHtml(c.foreignAddr)}</code></td>
        <td><span class="socket-state-pill">${escapeHtml(c.state || 'ESTABLISHED')}</span></td>
      </tr>
    `;
  });
  elements.netConnectionsTbody.innerHTML = html;
}

// Utility: Format Network Speed
function formatSpeed(bytesPerSec) {
  if (!bytesPerSec || bytesPerSec <= 0 || isNaN(bytesPerSec)) return '0 KB/s';
  const k = 1024;
  if (bytesPerSec < k) return `${Math.round(bytesPerSec)} B/s`;
  if (bytesPerSec < k * k) return `${(bytesPerSec / k).toFixed(1)} KB/s`;
  return `${(bytesPerSec / (k * k)).toFixed(1)} MB/s`;
}

// Live Network Rolling History Buffer
const netHistory = {
  max: 30,
  down: new Array(30).fill(0),
  up: new Array(30).fill(0),
  peakBps: 1024 * 512
};

function renderNetworkChart(stats) {
  if (!elements.netTrafficCanvas || !stats) return;
  const canvas = elements.netTrafficCanvas;
  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  const downBps = stats.downloadSpeedBps || 0;
  const upBps = stats.uploadSpeedBps || 0;

  netHistory.down.push(downBps);
  if (netHistory.down.length > netHistory.max) netHistory.down.shift();

  netHistory.up.push(upBps);
  if (netHistory.up.length > netHistory.max) netHistory.up.shift();

  const currentMax = Math.max(...netHistory.down, ...netHistory.up, 1024 * 50);
  if (currentMax > netHistory.peakBps) {
    netHistory.peakBps = currentMax;
  } else {
    netHistory.peakBps = Math.max(netHistory.peakBps * 0.98, currentMax, 1024 * 100);
  }

  if (elements.netLiveDown) elements.netLiveDown.textContent = stats.downloadSpeedStr || '0 KB/s';
  if (elements.netLiveUp) elements.netLiveUp.textContent = stats.uploadSpeedStr || '0 KB/s';
  if (elements.netLivePeak) elements.netLivePeak.textContent = formatSpeed(netHistory.peakBps);

  // Resize canvas according to layout width
  const rect = canvas.getBoundingClientRect();
  const dpr = window.devicePixelRatio || 1;
  const w = (rect.width > 50 ? rect.width : (canvas.parentElement ? canvas.parentElement.clientWidth : 600)) || 600;
  const h = 150;

  canvas.width = Math.round(w * dpr);
  canvas.height = Math.round(h * dpr);

  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  ctx.clearRect(0, 0, w, h);

  // Horizontal Grid Lines
  ctx.lineWidth = 1;
  ctx.strokeStyle = 'rgba(255, 255, 255, 0.05)';
  for (let y = 25; y < h - 10; y += 32) {
    ctx.beginPath();
    ctx.moveTo(0, y);
    ctx.lineTo(w, y);
    ctx.stroke();
  }

  const peak = Math.max(netHistory.peakBps * 1.15, 1024 * 50);
  const stepX = w / (netHistory.max - 1);

  function drawBezierLine(data, strokeColor, fillColor) {
    if (data.length < 2) return;
    const pts = data.map((val, idx) => ({
      x: idx * stepX,
      y: Math.max(15, Math.min(h - 14, (h - 14) - (val / peak) * (h - 35)))
    }));

    // Fill curve gradient
    ctx.beginPath();
    ctx.moveTo(pts[0].x, h);
    ctx.lineTo(pts[0].x, pts[0].y);
    for (let i = 1; i < pts.length; i++) {
      const xc = (pts[i].x + pts[i - 1].x) / 2;
      const yc = (pts[i].y + pts[i - 1].y) / 2;
      ctx.quadraticCurveTo(pts[i - 1].x, pts[i - 1].y, xc, yc);
    }
    ctx.lineTo(pts[pts.length - 1].x, pts[pts.length - 1].y);
    ctx.lineTo(pts[pts.length - 1].x, h);
    ctx.closePath();
    ctx.fillStyle = fillColor;
    ctx.fill();

    // Line stroke
    ctx.beginPath();
    ctx.moveTo(pts[0].x, pts[0].y);
    for (let i = 1; i < pts.length; i++) {
      const xc = (pts[i].x + pts[i - 1].x) / 2;
      const yc = (pts[i].y + pts[i - 1].y) / 2;
      ctx.quadraticCurveTo(pts[i - 1].x, pts[i - 1].y, xc, yc);
    }
    ctx.lineTo(pts[pts.length - 1].x, pts[pts.length - 1].y);
    ctx.strokeStyle = strokeColor;
    ctx.lineWidth = 2.5;
    ctx.stroke();
  }

  // Draw Download Curve (Cyan)
  const gradDown = ctx.createLinearGradient(0, 0, 0, h);
  gradDown.addColorStop(0, 'rgba(6, 182, 212, 0.32)');
  gradDown.addColorStop(1, 'rgba(6, 182, 212, 0.0)');
  drawBezierLine(netHistory.down, '#06b6d4', gradDown);

  // Draw Upload Curve (Indigo)
  const gradUp = ctx.createLinearGradient(0, 0, 0, h);
  gradUp.addColorStop(0, 'rgba(99, 102, 241, 0.28)');
  gradUp.addColorStop(1, 'rgba(99, 102, 241, 0.0)');
  drawBezierLine(netHistory.up, '#6366f1', gradUp);
}

// ==========================================================================
// Privacy Protection
// ==========================================================================
async function loadPrivacyTraces(isSilent = false) {
  if (!isSilent && elements.privacyItemsList) {
    elements.privacyItemsList.innerHTML = '<div class="loading-state">Gizlilik izleri taranıyor...</div>';
  }

  try {
    const res = await fetch('/api/privacy');
    if (!res.ok) throw new Error('Gizlilik izleri taranamadı');
    const data = await res.json();
    privacyData = data;

    // Update badge & stats
    if (elements.badgePrivacyCount) elements.badgePrivacyCount.textContent = data.totalItemsCount;
    if (elements.privStatTotal) elements.privStatTotal.textContent = data.totalItemsCount;
    if (elements.privStatSize) elements.privStatSize.textContent = data.totalSizeStr || '0 B';

    let historyCount = 0;
    let browserCount = 0;
    (data.items || []).forEach(it => {
      if (it.category === 'history') historyCount += it.count;
      if (it.category === 'browser') browserCount += it.count;
    });
    if (elements.privStatHistory) elements.privStatHistory.textContent = historyCount;
    if (elements.privStatBrowser) elements.privStatBrowser.textContent = browserCount;

    if (elements.privacyStatusText) {
      if (data.totalItemsCount === 0) {
        elements.privacyStatusText.textContent = "Gizlilik Durumu: %100 Güvenli & İz Bulunmuyor";
        if (elements.privacyScoreBadge) {
          elements.privacyScoreBadge.style.color = "var(--accent-emerald)";
          elements.privacyScoreBadge.style.borderColor = "rgba(16, 185, 129, 0.3)";
          elements.privacyScoreBadge.style.background = "rgba(16, 185, 129, 0.12)";
        }
      } else {
        elements.privacyStatusText.textContent = `Tespit Edilen İz: ${data.totalItemsCount} Kayıt (${data.totalSizeStr || '0 MB'}) — Temizlenmeye Hazır`;
        if (elements.privacyScoreBadge) {
          elements.privacyScoreBadge.style.color = "#fda4af";
          elements.privacyScoreBadge.style.borderColor = "rgba(244, 63, 94, 0.25)";
          elements.privacyScoreBadge.style.background = "rgba(244, 63, 94, 0.12)";
        }
      }
    }

    renderPrivacyTraces();
  } catch (err) {
    console.error('Privacy scan error:', err);
    if (elements.privacyItemsList) {
      elements.privacyItemsList.innerHTML = `<div class="error-state">Hata: ${escapeHtml(err.message)}</div>`;
    }
  }
}

function renderPrivacyTraces() {
  if (!elements.privacyItemsList || !privacyData) return;

  const items = privacyData.items || [];
  if (items.length === 0) {
    elements.privacyItemsList.innerHTML = `
      <div class="empty-state" style="padding: 30px; text-align: center; color: var(--accent-emerald);">
        <p>Tebrikler! Sisteminizde temizlenecek gizlilik izi bulunamadı.</p>
      </div>
    `;
  } else {
    // Group by category
    const categories = {
      history: { title: 'Terminal & Kabuk Komut Geçmişi', items: [] },
      browser: { title: 'Tarayıcı & İnternet Gezinti İzleri', items: [] },
      recents: { title: 'Sistem Etkinlik & Son Kullanılan Belgeler', items: [] },
      other: { title: 'Diğer Sistem İzleri', items: [] }
    };

    items.forEach(it => {
      if (categories[it.category]) {
        categories[it.category].items.push(it);
      } else {
        categories.other.items.push(it);
      }
    });

    let html = '';
    Object.values(categories).forEach(cat => {
      if (cat.items.length === 0) return;
      html += `
        <div style="margin-bottom: 12px;">
          <h4 style="font-size:0.8rem; font-weight:700; color:var(--text-muted); margin-bottom: 8px; text-transform:uppercase; letter-spacing:0.04em;">${escapeHtml(cat.title)}</h4>
          <div style="display:flex; flex-direction:column; gap:8px;">
      `;

      cat.items.forEach(it => {
        let iconSvg = '';
        if (it.icon === 'clock') {
          iconSvg = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/></svg>';
        } else if (it.icon === 'terminal') {
          iconSvg = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>';
        } else {
          iconSvg = '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M2 12h20"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg>';
        }

        html += `
          <div class="privacy-item-row">
            <div class="privacy-item-left">
              <input type="checkbox" class="check-privacy-item" data-id="${escapeHtml(it.id)}" checked style="cursor: pointer; width: 18px; height: 18px;" />
              <div class="privacy-icon ${escapeHtml(it.icon || 'globe')}">${iconSvg}</div>
              <div>
                <h4 style="font-size: 0.95rem; font-weight: 700; color: #fff; margin-bottom: 2px;">${escapeHtml(it.title)}</h4>
                <p style="font-size: 0.78rem; color: var(--text-muted);">${escapeHtml(it.description)}</p>
              </div>
            </div>
            <div style="text-align: right; flex-shrink: 0;">
              <div style="font-weight: 700; font-size: 0.92rem; color: var(--accent-rose);">${it.count} ${escapeHtml(it.countLabel)}</div>
              <div style="font-size: 0.75rem; color: var(--text-dim);">${escapeHtml(it.sizeStr)}</div>
            </div>
          </div>
        `;
      });

      html += `</div></div>`;
    });

    elements.privacyItemsList.innerHTML = html;
  }

  // Render TCC Permissions
  if (elements.privacyPermissionsList && privacyData.permissions) {
    let phtml = '';
    privacyData.permissions.forEach(perm => {
      phtml += `
        <div class="permission-card-box">
          <div class="permission-card-top">
            <div class="permission-card-info">
              <h4>${escapeHtml(perm.serviceName)}</h4>
              <p>${escapeHtml(perm.description)}</p>
            </div>
          </div>
          <button class="btn btn-secondary btn-xs btn-reset-perm" data-service="${escapeHtml(perm.service)}" data-name="${escapeHtml(perm.serviceName)}" style="align-self: flex-start;">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width:12px; height:12px;"><path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99"/></svg>
            <span>İzinleri Sıfırla</span>
          </button>
        </div>
      `;
    });
    elements.privacyPermissionsList.innerHTML = phtml;

    elements.privacyPermissionsList.querySelectorAll('.btn-reset-perm').forEach(btn => {
      btn.addEventListener('click', () => {
        const s = btn.dataset.service;
        const n = btn.dataset.name;
        resetTccPermission(s, n);
      });
    });
  }
}

async function cleanSelectedPrivacyItems() {
  if (!elements.privacyItemsList) return;
  const checkboxes = elements.privacyItemsList.querySelectorAll('.check-privacy-item:checked');
  const selectedIDs = Array.from(checkboxes).map(c => c.dataset.id);

  if (selectedIDs.length === 0) {
    showToast('Lütfen temizlemek için en az bir gizlilik izi seçin.', 'info');
    return;
  }

  if (!confirm(`Seçilen ${selectedIDs.length} adet gizlilik kaydı kalıcı olarak temizlenecektir. Devam edilsin mi?`)) {
    return;
  }

  try {
    const res = await fetch('/api/privacy/clean', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ itemIds: selectedIDs })
    });
    const result = await res.json();
    if (res.ok && result.success) {
      showToast(result.message || 'Gizlilik izleri temizlendi.', 'success');
      loadPrivacyTraces();
    } else {
      showToast('Temizleme hatası: ' + (result.error || 'Başarısız'), 'error');
    }
  } catch (err) {
    showToast('Hata: ' + err.message, 'error');
  }
}

async function resetTccPermission(service, name) {
  if (!confirm(`"${name}" için tüm uygulama izinleri sıfırlanacaktır. Uygulamalar ilk açılışta tekrar onay isteyecektir. Onaylıyor musunuz?`)) {
    return;
  }

  try {
    const res = await fetch('/api/privacy/reset', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ service: service })
    });
    const result = await res.json();
    if (res.ok && result.success) {
      showToast(result.message || `${name} izinleri sıfırlandı.`, 'success');
    } else {
      showToast('Sıfırlama hatası: ' + (result.error || 'Başarısız'), 'error');
    }
  } catch (err) {
    showToast('Hata: ' + err.message, 'error');
  }
}

// Setup Event Handlers
function setupEventHandlers() {
  // Extensions Handlers
  if (elements.btnRefreshExtensions) {
    elements.btnRefreshExtensions.addEventListener('click', () => {
      loadExtensions();
      showToast('Eklentiler güncellendi', 'info');
    });
  }

  if (elements.inputSearchExtensions) {
    elements.inputSearchExtensions.addEventListener('input', (e) => {
      extSearchQuery = e.target.value.trim();
      renderExtensionsList();
    });
  }

  if (elements.extCategoryFilters) {
    elements.extCategoryFilters.querySelectorAll('.tab-pill').forEach(pill => {
      pill.addEventListener('click', () => {
        elements.extCategoryFilters.querySelectorAll('.tab-pill').forEach(p => p.classList.remove('active'));
        pill.classList.add('active');
        extCurrentFilter = pill.dataset.extFilter;
        renderExtensionsList();
      });
    });
  }

  // Network Handlers
  if (elements.btnRefreshNetwork) {
    elements.btnRefreshNetwork.addEventListener('click', () => {
      loadNetworkStats();
      showToast('Ağ durumu güncellendi', 'info');
    });
  }

  // Privacy Handlers
  if (elements.btnRefreshPrivacy) {
    elements.btnRefreshPrivacy.addEventListener('click', () => {
      loadPrivacyTraces();
      showToast('Gizlilik izleri güncellendi', 'info');
    });
  }

  if (elements.btnCleanPrivacy) {
    elements.btnCleanPrivacy.addEventListener('click', () => {
      cleanSelectedPrivacyItems();
    });
  }

  if (elements.checkAllPrivacy) {
    elements.checkAllPrivacy.addEventListener('change', (e) => {
      const checked = e.target.checked;
      if (elements.privacyItemsList) {
        elements.privacyItemsList.querySelectorAll('.check-privacy-item').forEach(chk => {
          chk.checked = checked;
        });
      }
    });
  }

  // Startup Handlers
  if (elements.btnRefreshStartup) {
    elements.btnRefreshStartup.addEventListener('click', () => {
      loadStartupItems();
      showToast('Başlangıç öğeleri güncellendi', 'info');
    });
  }

  if (elements.btnAddLoginItemModal) {
    elements.btnAddLoginItemModal.addEventListener('click', () => {
      if (elements.addLoginItemDialog) elements.addLoginItemDialog.showModal();
    });
  }

  if (elements.btnCancelAddLogin) {
    elements.btnCancelAddLogin.addEventListener('click', () => {
      if (elements.addLoginItemDialog) elements.addLoginItemDialog.close();
    });
  }

  if (elements.formAddLoginItem) {
    elements.formAddLoginItem.addEventListener('submit', () => handleAddLoginItemSubmit());
  }

  if (elements.startupFilterChips) {
    elements.startupFilterChips.querySelectorAll('.chip').forEach(chip => {
      chip.addEventListener('click', () => {
        elements.startupFilterChips.querySelectorAll('.chip').forEach(c => c.classList.remove('active'));
        chip.classList.add('active');
        state.startupFilter = chip.dataset.filter;
        renderStartupItems();
      });
    });
  }

  if (elements.inputStartupSearch) {
    elements.inputStartupSearch.addEventListener('input', (e) => {
      state.startupSearch = e.target.value.trim();
      renderStartupItems();
    });
  }
  // Auth Form Handlers
  if (elements.btnAuthSubmit) {
    elements.btnAuthSubmit.addEventListener('click', () => handleLoginSubmit());
  }
  if (elements.authPasswordInput) {
    elements.authPasswordInput.addEventListener('keydown', (e) => {
      if (e.key === 'Enter') handleLoginSubmit();
    });
  }
  if (elements.btnLogout) {
    elements.btnLogout.addEventListener('click', () => handleLogout());
  }

  // Maintenance Handlers
  if (elements.btnRunAllMaintenance) {
    elements.btnRunAllMaintenance.addEventListener('click', () => runAllMaintenanceTasks());
  }
  if (elements.btnClearMaintenanceLog) {
    elements.btnClearMaintenanceLog.addEventListener('click', () => {
      if (elements.maintenanceLogContent) elements.maintenanceLogContent.textContent = '';
      if (elements.maintenanceLogBox) elements.maintenanceLogBox.style.display = 'none';
    });
  }

  // Shredder Handlers
  if (elements.btnExecuteShred) {
    elements.btnExecuteShred.addEventListener('click', () => executeShred());
  }

  // Hardware Monitor Handlers
  if (elements.btnRefreshMonitor) {
    elements.btnRefreshMonitor.addEventListener('click', () => {
      loadHardwareMonitor();
      showToast('Donanım verileri güncellendi', 'info');
    });
  }

  elements.btnGlobalScan.addEventListener('click', () => {
    startGlobalScan();
  });

  elements.btnQuickClean.addEventListener('click', () => {
    handleQuickClean();
  });

  // Storage Alert Banner button
  if (elements.btnAlertSmartClean) {
    elements.btnAlertSmartClean.addEventListener('click', () => {
      executeSmartCareClean();
    });
  }

  // Developer tab buttons
  elements.btnSelectAllDev.addEventListener('click', () => {
    if (!state.scanData) return;
    const allSelected = state.scanData.targets.filter(t => t.category === 'developer' && t.exists).every(t => t.selected);
    state.scanData.targets.filter(t => t.category === 'developer' && t.exists).forEach(t => t.selected = !allSelected);
    renderTargetList(elements.devTargetsList, state.scanData.targets.filter(t => t.category === 'developer' && t.exists), 'dev');
    updateSelectedSize('dev');
  });

  elements.btnCleanSelectedDev.addEventListener('click', () => {
    if (!state.scanData) return;
    const ids = state.scanData.targets.filter(t => t.category === 'developer' && t.selected).map(t => t.id);
    if (ids.length === 0) {
      showToast('Seçili geliştirici önbelleği yok!', 'info');
      return;
    }
    confirmAndClean(ids, [], "Seçili Geliştirici Önbellekleri", elements.devSelectedSize.textContent);
  });

  // System tab buttons
  elements.btnSelectAllSys.addEventListener('click', () => {
    if (!state.scanData) return;
    const allSelected = state.scanData.targets.filter(t => t.category === 'system' && t.exists).every(t => t.selected);
    state.scanData.targets.filter(t => t.category === 'system' && t.exists).forEach(t => t.selected = !allSelected);
    renderTargetList(elements.sysTargetsList, state.scanData.targets.filter(t => t.category === 'system' && t.exists), 'sys');
    updateSelectedSize('sys');
  });

  elements.btnCleanSelectedSys.addEventListener('click', () => {
    if (!state.scanData) return;
    const ids = state.scanData.targets.filter(t => t.category === 'system' && t.selected).map(t => t.id);
    const customPaths = (state.scanData.dynamicCaches || []).filter(c => c.selected).map(c => c.path);
    if (ids.length === 0 && customPaths.length === 0) {
      showToast('Seçili sistem artığı yok!', 'info');
      return;
    }
    confirmAndClean(ids, customPaths, "Seçili Sistem ve Uygulama Artıkları", elements.sysSelectedSize.textContent);
  });

  // Node Modules tab buttons
  elements.btnScanNodeModules.addEventListener('click', () => {
    startNodeModulesScan();
  });
  elements.btnApplyNodePath.addEventListener('click', () => {
    startNodeModulesScan();
  });
  elements.thCheckAllNode.addEventListener('change', (e) => {
    const isChecked = e.target.checked;
    state.selectedNodeModules.clear();
    if (isChecked) {
      state.nodeModulesData.forEach(i => state.selectedNodeModules.add(i.modulesPath));
    }
    renderNodeModulesTable(state.nodeModulesData);
    updateNodeSelectedSize();
  });
  elements.btnCleanSelectedNode.addEventListener('click', () => {
    if (state.selectedNodeModules.size === 0) return;
    confirmCleanNodeModules(Array.from(state.selectedNodeModules), elements.nodeSelectedSize.textContent);
  });

  // Large Files tab buttons
  elements.btnScanLargeFiles.addEventListener('click', () => {
    startLargeFilesScan();
  });
  elements.thCheckAllFiles.addEventListener('change', (e) => {
    const isChecked = e.target.checked;
    state.selectedLargeFiles.clear();
    if (isChecked) {
      state.largeFilesData.forEach(i => state.selectedLargeFiles.add(i.path));
    }
    renderLargeFilesTable(state.largeFilesData);
    updateLargeFilesSelectedSize();
  });
  elements.btnCleanSelectedFiles.addEventListener('click', () => {
    if (state.selectedLargeFiles.size === 0) return;
    confirmCleanLargeFiles(Array.from(state.selectedLargeFiles), elements.filesSelectedSize.textContent);
  });

  // Directory Tree tab buttons
  elements.btnTreeUp.addEventListener('click', () => {
    if (state.treeData && state.treeData.parentPath) {
      loadDirTree(state.treeData.parentPath);
    }
  });
  elements.btnTreeRefresh.addEventListener('click', () => {
    loadDirTree(state.treeCurrentPath);
  });
  elements.btnTreeGo.addEventListener('click', () => {
    const p = elements.inputTreePath.value.trim();
    if (p) loadDirTree(p);
  });
  elements.inputTreePath.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
      const p = elements.inputTreePath.value.trim();
      if (p) loadDirTree(p);
    }
  });
  elements.treePills.forEach(pill => {
    pill.addEventListener('click', () => {
      const p = pill.dataset.path;
      if (p) loadDirTree(p);
    });
  });

  // Duplicate Finder tab buttons
  elements.btnScanDuplicates.addEventListener('click', () => {
    startDuplicatesScan();
  });
  elements.btnDupSelectCopies.addEventListener('click', () => {
    selectOnlyCopies();
  });
  elements.btnDupCleanSelected.addEventListener('click', () => {
    confirmCleanDuplicates();
  });

  // Leftovers tab buttons
  elements.btnScanLeftovers.addEventListener('click', () => {
    startLeftoversScan();
  });
  elements.btnSelectAllLeftovers.addEventListener('click', () => {
    if (!state.leftoversData || !state.leftoversData.items) return;
    const allSelected = state.leftoversData.items.every(i => state.selectedLeftovers.has(i.id));
    state.selectedLeftovers.clear();
    state.leftoversData.items.forEach(i => {
      i.selected = !allSelected;
      if (!allSelected) state.selectedLeftovers.add(i.id);
    });
    renderLeftovers(state.leftoversData);
  });
  elements.btnCleanSelectedLeftovers.addEventListener('click', () => {
    confirmCleanLeftovers();
  });

  // Apps tab handlers
  if (elements.btnScanApps) {
    elements.btnScanApps.addEventListener('click', () => startAppsScan());
  }
  if (elements.appsSearchInput) {
    elements.appsSearchInput.addEventListener('input', (e) => {
      state.appsSearch = e.target.value;
      renderAppsList();
    });
  }
  if (elements.appsFilterChips) {
    elements.appsFilterChips.querySelectorAll('.chip').forEach(chip => {
      chip.addEventListener('click', () => {
        elements.appsFilterChips.querySelectorAll('.chip').forEach(c => c.classList.remove('active'));
        chip.classList.add('active');
        state.appsFilter = chip.dataset.filter;
        renderAppsList();
      });
    });
  }

  // Browsers tab handlers
  if (elements.btnScanBrowsers) {
    elements.btnScanBrowsers.addEventListener('click', () => startBrowsersScan());
  }
  if (elements.btnCleanAllBrowsers) {
    elements.btnCleanAllBrowsers.addEventListener('click', () => {
      if (!state.browsersData || !state.browsersData.browsers) return;
      const ids = state.browsersData.browsers.filter(b => b.size > 0).map(b => b.id);
      confirmCleanBrowsers(ids, "Tüm Tarayıcı Önbelleklerini Temizle", elements.browsersTotalCleanSize.textContent);
    });
  }

  // Downloads tab handlers
  if (elements.btnScanDownloads) {
    elements.btnScanDownloads.addEventListener('click', () => startDownloadsScan());
  }
  if (elements.btnCleanDownloads) {
    elements.btnCleanDownloads.addEventListener('click', () => confirmCleanDownloads());
  }
  if (elements.downloadsFilterChips) {
    elements.downloadsFilterChips.querySelectorAll('.chip').forEach(chip => {
      chip.addEventListener('click', () => {
        elements.downloadsFilterChips.querySelectorAll('.chip').forEach(c => c.classList.remove('active'));
        chip.classList.add('active');
        state.downloadsFilter = chip.dataset.cat;
        renderDownloadsList();
      });
    });
  }
  if (elements.btnSelectAllDownloads) {
    elements.btnSelectAllDownloads.addEventListener('click', () => {
      if (!state.downloadsData || !state.downloadsData.items) return;
      const allSelected = state.downloadsData.items.every(i => state.selectedDownloads.has(i.path));
      state.selectedDownloads.clear();
      if (!allSelected) {
        state.downloadsData.items.forEach(i => state.selectedDownloads.add(i.path));
      }
      renderDownloadsList();
      updateDownloadsSelectedSize();
    });
  }

  // Apple & System Data tab handlers
  if (elements.btnScanApple) {
    elements.btnScanApple.addEventListener('click', () => startAppleScan());
  }
  if (elements.btnCleanSafeSystemData) {
    elements.btnCleanSafeSystemData.addEventListener('click', () => cleanAllSafeSystemData());
  }
  if (elements.btnReclaimPurgeable) {
    elements.btnReclaimPurgeable.addEventListener('click', () => reclaimPurgeableSpace());
  }
  if (elements.btnCleanAllSnapshots) {
    elements.btnCleanAllSnapshots.addEventListener('click', () => deleteAllSnapshots());
  }
  if (elements.btnCleanSimulators) {
    elements.btnCleanSimulators.addEventListener('click', () => cleanSimulators());
  }

  // Media tab handlers
  if (elements.btnScanMedia) {
    elements.btnScanMedia.addEventListener('click', () => startMediaScan());
  }
  if (elements.btnCleanMedia) {
    elements.btnCleanMedia.addEventListener('click', () => confirmCleanMedia());
  }

  // Modal buttons
  elements.btnModalCancel.addEventListener('click', () => {
    elements.confirmModal.close();
    pendingCleanAction = null;
  });

  elements.btnModalConfirm.addEventListener('click', () => {
    if (pendingCleanAction) {
      pendingCleanAction();
    }
  });

  // Light dismiss modal on backdrop click
  elements.confirmModal.addEventListener('click', (e) => {
    const rect = elements.confirmModal.getBoundingClientRect();
    const isInDialog = (rect.top <= e.clientY && e.clientY <= rect.top + rect.height
      && rect.left <= e.clientX && e.clientX <= rect.left + rect.width);
    if (!isInDialog) {
      elements.confirmModal.close();
    }
  });
}

// Utility: Show Toast Notification
function showToast(message, type = 'info') {
  const toast = document.createElement('div');
  toast.className = `toast toast-${type}`;
  toast.innerHTML = `
    <span>${escapeHtml(message)}</span>
  `;

  elements.toastContainer.appendChild(toast);
  setTimeout(() => {
    toast.style.opacity = '0';
    toast.style.transform = 'translateY(10px)';
    toast.style.transition = 'all 0.3s ease';
    setTimeout(() => toast.remove(), 300);
  }, 4500);
}

// Utility: Format Bytes
function formatBytes(bytes) {
  if (bytes === 0 || isNaN(bytes)) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

// Utility: Escape HTML
function escapeHtml(str) {
  if (!str) return '';
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}

/* ==========================================================================
   Full Disk Access (FDA) Detection & Modal Management
   ========================================================================== */
let fdaStatusData = null;

async function checkFDAPermissions() {
  const fdaDot = document.getElementById('fda-dot');
  const fdaText = document.getElementById('fda-badge-text');
  const btnBadge = document.getElementById('btn-fda-badge');
  if (!fdaDot || !fdaText) return;

  try {
    const res = await fetch('/api/system/permissions');
    if (!res.ok) return;
    const data = await res.json();
    fdaStatusData = data;

    if (data.hasFullDiskAccess) {
      fdaDot.className = 'fda-dot granted';
      fdaText.textContent = 'FDA: Aktif';
      if (btnBadge) btnBadge.title = 'macOS Tam Disk Erişimi (FDA) etkinleştirildi.';
    } else {
      fdaDot.className = 'fda-dot missing';
      fdaText.textContent = '⚠️ FDA Gerekli';
      if (btnBadge) btnBadge.title = 'Safari, Mail ve Time Machine için Tam Disk Erişimi gerekli. Tıklayın.';
    }
  } catch (e) {
    console.warn('FDA kontrolü yapılamadı:', e);
  }
}

function initFDA() {
  const btnBadge = document.getElementById('btn-fda-badge');
  const modal = document.getElementById('fda-modal');
  const btnClose = document.getElementById('btn-close-fda-modal');
  const btnOpenSettings = document.getElementById('btn-open-fda-settings');
  const btnRecheck = document.getElementById('btn-recheck-fda');

  if (btnBadge && modal) {
    btnBadge.addEventListener('click', (e) => {
      e.preventDefault();
      e.stopPropagation();
      const desc = document.getElementById('fda-modal-desc');
      if (desc && fdaStatusData) {
        if (fdaStatusData.hasFullDiskAccess) {
          desc.innerHTML = '✅ <b>Tam Disk Erişimi (FDA) etkinleştirildi!</b> Uygulama Safari, Mail, Mesajlar ve Time Machine dizinlerini koruma engeline takılmadan tarayabilir.';
        } else {
          desc.innerHTML = 'macOS güvenlik politikaları gereğince <b>Safari</b>, <b>Apple Mail</b>, <b>iMessage ekleri</b>, <b>Time Machine yerel anlık görüntüleri</b> ve korumalı sistem loglarını eksiksiz tarayabilmek ve temizleyebilmek için uygulamanıza <b>Tam Disk Erişimi</b> izni verilmesi gerekir.';
        }
      }
      try {
        if (!modal.open) {
          modal.showModal();
        }
      } catch (err) {
        console.error('Modal acilamadi:', err);
      }
    });
  }

  if (btnClose && modal) {
    btnClose.addEventListener('click', (e) => {
      e.preventDefault();
      e.stopPropagation();
      modal.close();
    });
    modal.addEventListener('click', (e) => {
      if (e.target === modal) {
        modal.close();
      }
    });
  }

  if (btnOpenSettings) {
    btnOpenSettings.addEventListener('click', async (e) => {
      e.preventDefault();
      e.stopPropagation();
      try {
        await fetch('/api/system/permissions/open-fda', { method: 'POST' });
        showToast('macOS Sistem Ayarları - Tam Disk Erişimi açıldı.', 'info');
      } catch (e) {
        showToast('Ayar penceresi açılamadı: ' + e.message, 'error');
      }
    });
  }

  if (btnRecheck) {
    btnRecheck.addEventListener('click', async (e) => {
      e.preventDefault();
      e.stopPropagation();
      btnRecheck.disabled = true;
      btnRecheck.textContent = 'Kontrol ediliyor...';
      await checkFDAPermissions();
      btnRecheck.disabled = false;
      btnRecheck.textContent = '🔄 Yeniden Kontrol Et';

      if (fdaStatusData && fdaStatusData.hasFullDiskAccess) {
        showToast('Harika! Tam Disk Erişimi başarıyla doğrulandı.', 'success');
        modal.close();
      } else {
        showToast('Tam Disk Erişimi henüz algılanmadı. Lütfen Sistem Ayarları listesinden izin verip tekrar deneyin.', 'warning');
      }
    });
  }
}

/* ==========================================================================
   Live Scan Progress Stream (SSE)
   ========================================================================== */
let sseSource = null;
let sseHideTimer = null;

function initScanProgressStream() {
  if (typeof EventSource === 'undefined') return;

  const pill = document.getElementById('live-scan-progress-pill');
  const label = document.getElementById('live-progress-label');
  const fill = document.getElementById('live-progress-fill');
  const pct = document.getElementById('live-progress-pct');

  if (!pill) return;

  if (sseSource) {
    try { sseSource.close(); } catch (_) {}
    sseSource = null;
  }

  try {
    sseSource = createEventSource('/api/events/progress');

    sseSource.addEventListener('progress', (e) => {
      try {
        const ev = JSON.parse(e.data);
        if (!ev || !ev.message) return;

        if (sseHideTimer) clearTimeout(sseHideTimer);

        pill.style.display = 'flex';
        label.textContent = ev.message;
        const p = ev.percent || 0;
        fill.style.width = p + '%';
        pct.textContent = p + '%';

        if (p >= 100 || ev.message.includes('tamamlandı')) {
          sseHideTimer = setTimeout(() => {
            pill.style.display = 'none';
          }, 3000);
        }
      } catch (err) {
        // ignore parse error
      }
    });

    sseSource.onerror = () => {
      // EventSource automatically reconnects on error
    };
  } catch (err) {
    console.warn('SSE bağlantısı kurulamadı:', err);
  }
}

/* ==========================================================================
   Docker & Container Storage Management
   ========================================================================== */
let dockerStatusCache = null;

async function loadDockerStatus() {
  const container = document.getElementById('docker-metrics-container');
  const badge = document.getElementById('docker-badge-status');
  const subtitle = document.getElementById('docker-status-subtitle');
  const btnPruneAll = document.getElementById('btn-docker-prune-all');
  if (!container) return;

  try {
    const res = await fetch('/api/docker/status');
    if (!res.ok) throw new Error('Docker durumu alınamadı');
    const data = await res.json();
    dockerStatusCache = data;

    if (badge) {
      if (data.running) {
        badge.className = 'badge badge-success';
        badge.textContent = 'Aktif (v' + (data.version || 'Engine') + ')';
        if (subtitle) subtitle.textContent = data.message || 'Docker servisi çalışıyor.';
        if (btnPruneAll) btnPruneAll.disabled = false;
      } else if (data.installed) {
        badge.className = 'badge badge-warning';
        badge.textContent = 'Durduruldu';
        if (subtitle) subtitle.textContent = data.message || 'Docker kurulu fakat daemon çalışmıyor.';
        if (btnPruneAll) btnPruneAll.disabled = true;
      } else {
        badge.className = 'badge badge-secondary';
        badge.textContent = 'Kurulu Değil';
        if (subtitle) subtitle.textContent = data.message || 'Sisteminizde aktif Docker CLI bulunamadı.';
        if (btnPruneAll) btnPruneAll.disabled = true;
      }
    }

    let html = '';

    // Show Docker.raw virtual disk card if present
    if (data.dockerRawSize > 0) {
      html += `
        <div class="docker-metric-box" style="border-left: 3px solid #3b82f6;">
          <div class="docker-metric-title">Docker Desktop Sanal Diski (Docker.raw)</div>
          <div class="docker-metric-val">${data.dockerRawSizeStr}</div>
          <div class="docker-metric-sub">macOS VM Sanal Disk Dosyası</div>
        </div>
      `;
    }

    if (data.colimaSizeStr) {
      html += `
        <div class="docker-metric-box" style="border-left: 3px solid #8b5cf6;">
          <div class="docker-metric-title">Colima VM Depolaması</div>
          <div class="docker-metric-val">${data.colimaSizeStr}</div>
          <div class="docker-metric-sub">~/.colima sanal disk boyutu</div>
        </div>
      `;
    }

    if (data.components && data.components.length > 0) {
      data.components.forEach(comp => {
        let typeName = comp.type;
        let pruneAction = 'system';
        if (typeName === 'Images') { typeName = 'İmajlar (Images)'; pruneAction = 'images'; }
        else if (typeName === 'Containers') { typeName = 'Konteynerlar (Containers)'; pruneAction = 'containers'; }
        else if (typeName === 'Local Volumes') { typeName = 'Yerel Birimler (Volumes)'; pruneAction = 'volumes'; }
        else if (typeName === 'Build Cache') { typeName = 'Build Cache'; pruneAction = 'buildcache'; }

        html += `
          <div class="docker-metric-box">
            <div style="display:flex; justify-content:space-between; align-items:flex-start;">
              <div class="docker-metric-title">${escapeHtml(typeName)}</div>
              <button class="btn btn-secondary btn-xs btn-prune-single" data-prune="${pruneAction}" style="padding:2px 8px; font-size:0.7rem;">Temizle</button>
            </div>
            <div class="docker-metric-val">${escapeHtml(comp.sizeStr || '0 B')}</div>
            <div class="docker-metric-sub">Toplam: ${comp.totalCount || 0} • Aktif: ${comp.activeCount || 0} • Geri Kazanılabilir: ${escapeHtml(comp.reclaimable || '0%')}</div>
          </div>
        `;
      });
    } else if (!data.running) {
      html += `
        <div style="grid-column: 1 / -1; padding: 16px; background: rgba(255,255,255,0.02); border-radius: 8px; font-size: 0.84rem; color: var(--text-dim);">
          💡 Konteyner, imaj ve build cache detaylarını görebilmek ve tek tıkla prune edebilmek için Docker uygulamasını (Docker Desktop, Colima veya OrbStack) başlatın.
        </div>
      `;
    }

    container.innerHTML = html;

    // Attach single prune listeners
    container.querySelectorAll('.btn-prune-single').forEach(btn => {
      btn.addEventListener('click', () => {
        const pruneType = btn.dataset.prune;
        executeDockerClean(pruneType);
      });
    });

  } catch (err) {
    container.innerHTML = `<div class="loading-state" style="color:var(--accent-rose);">Docker durumu okunamadı: ${escapeHtml(err.message)}</div>`;
  }
}

async function executeDockerClean(type = 'all') {
  let label = 'Tüm kullanılmayan Docker kaynakları (dangling & unreferenced imajlar, durdurulan konteynerlar ve build cache)';
  if (type === 'images') label = 'Kullanılmayan tüm Docker imajları';
  if (type === 'containers') label = 'Durdurulmuş Docker konteynerları';
  if (type === 'volumes') label = 'Kullanılmayan Docker birimleri (volumes)';
  if (type === 'buildcache') label = 'Docker buildx önbelleği';

  if (!confirm(`${label} temizlensin mi?\n\nBu işlem geri alınamaz.`)) {
    return;
  }

  showToast('Docker temizliği çalıştırılıyor...', 'info');

  try {
    const res = await fetch('/api/docker/clean', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ type: type })
    });
    const result = await res.json();
    if (!res.ok || result.error) {
      throw new Error(result.error || 'Temizlik başarısız oldu');
    }
    showToast(result.message || 'Docker temizliği tamamlandı.', 'success');
    loadDockerStatus();
    fetchSystemStats();
  } catch (err) {
    showToast('Docker temizleme hatası: ' + err.message, 'error');
  }
}

/* ==========================================================================
   DaisyDisk Style Sunburst (Radial Treemap) Engine
   ========================================================================== */
let sunburstSlices = [];
let sunburstHovered = null;
let sunburstCurrentTree = null;

const SUNBURST_PALETTE = [
  '#3b82f6', '#8b5cf6', '#06b6d4', '#10b981', '#f59e0b',
  '#ec4899', '#6366f1', '#14b8a6', '#f97316', '#a855f7',
  '#0284c7', '#4f46e5', '#059669', '#d97706', '#db2777'
];

function initDirTreeViewSwitcher() {
  const btnList = document.getElementById('btn-tree-view-list');
  const btnSunburst = document.getElementById('btn-tree-view-sunburst');
  const listContainer = document.getElementById('tree-items-list');
  const sunburstContainer = document.getElementById('tree-sunburst-view');

  if (btnList && btnSunburst && listContainer && sunburstContainer) {
    btnList.addEventListener('click', () => {
      btnList.classList.add('active');
      btnSunburst.classList.remove('active');
      listContainer.style.display = 'block';
      sunburstContainer.style.display = 'none';
    });

    btnSunburst.addEventListener('click', () => {
      btnSunburst.classList.add('active');
      btnList.classList.remove('active');
      listContainer.style.display = 'none';
      sunburstContainer.style.display = 'flex';
      if (sunburstCurrentTree) {
        drawSunburst(sunburstCurrentTree);
      }
    });
  }

  // Setup canvas interactions
  const canvas = document.getElementById('dirtree-sunburst-canvas');
  if (canvas) {
    canvas.addEventListener('mousemove', handleSunburstMouseMove);
    canvas.addEventListener('mouseleave', handleSunburstMouseLeave);
    canvas.addEventListener('click', handleSunburstClick);
  }

  // Docker refresh button
  const btnRefreshDocker = document.getElementById('btn-refresh-docker');
  if (btnRefreshDocker) {
    btnRefreshDocker.addEventListener('click', () => loadDockerStatus());
  }

  // Docker prune all button
  const btnDockerPruneAll = document.getElementById('btn-docker-prune-all');
  if (btnDockerPruneAll) {
    btnDockerPruneAll.addEventListener('click', () => executeDockerClean('all'));
  }
}

function renderSunburstChart(treeData) {
  sunburstCurrentTree = treeData;
  const sunburstContainer = document.getElementById('tree-sunburst-view');
  if (sunburstContainer && sunburstContainer.style.display !== 'none') {
    drawSunburst(treeData);
  }
}

function drawSunburst(treeData) {
  const canvas = document.getElementById('dirtree-sunburst-canvas');
  if (!canvas || !treeData) return;
  const ctx = canvas.getContext('2d');
  const width = canvas.width;
  const height = canvas.height;
  const cx = width / 2;
  const cy = height / 2;
  const r0 = 80;  // Center circle radius
  const r1 = 95;  // Inner arc radius
  const r2 = 250; // Outer arc radius

  ctx.clearRect(0, 0, width, height);

  sunburstSlices = [];

  const items = (treeData.items || []).filter(i => i.size > 0);
  const totalSize = treeData.totalSize || 1;

  // Center Circle (Current folder & parent return button)
  ctx.save();
  ctx.beginPath();
  ctx.arc(cx, cy, r0, 0, 2 * Math.PI);
  ctx.fillStyle = sunburstHovered && sunburstHovered.isCenter ? 'rgba(59, 130, 246, 0.25)' : 'rgba(30, 41, 59, 0.85)';
  ctx.fill();
  ctx.strokeStyle = sunburstHovered && sunburstHovered.isCenter ? '#60a5fa' : 'rgba(255, 255, 255, 0.15)';
  ctx.lineWidth = 2;
  ctx.stroke();

  // Center text
  ctx.fillStyle = '#ffffff';
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  ctx.font = 'bold 15px -apple-system, BlinkMacSystemFont, "Plus Jakarta Sans", sans-serif';
  let curName = treeData.currentPath ? treeData.currentPath.split('/').filter(Boolean).pop() || '/' : '~';
  if (curName.length > 14) curName = curName.substring(0, 12) + '...';
  ctx.fillText(curName, cx, cy - 10);

  ctx.font = '600 12px "JetBrains Mono", monospace';
  ctx.fillStyle = '#93c5fd';
  ctx.fillText(treeData.totalStr || '0 B', cx, cy + 12);

  if (treeData.parentPath) {
    ctx.font = '10px -apple-system, sans-serif';
    ctx.fillStyle = 'rgba(255,255,255,0.5)';
    ctx.fillText('▲ Üst Klasör', cx, cy + 30);
  }
  ctx.restore();

  if (items.length === 0) return;

  // Draw slices
  let curAngle = -Math.PI / 2;
  const minAngle = 0.03; // minimum slice angle

  items.forEach((item, idx) => {
    let rawAngle = (item.size / totalSize) * (2 * Math.PI);
    if (rawAngle < minAngle) rawAngle = minAngle;
    const endAngle = curAngle + rawAngle;

    const isHovered = sunburstHovered && sunburstHovered.index === idx;
    const currentR2 = isHovered ? r2 + 8 : r2;
    const color = SUNBURST_PALETTE[idx % SUNBURST_PALETTE.length];

    ctx.save();
    ctx.beginPath();
    ctx.arc(cx, cy, currentR2, curAngle, endAngle);
    ctx.arc(cx, cy, r1, endAngle, curAngle, true);
    ctx.closePath();

    ctx.fillStyle = isHovered ? color : adjustAlpha(color, 0.78);
    ctx.fill();
    ctx.strokeStyle = isHovered ? '#ffffff' : 'rgba(15, 23, 42, 0.8)';
    ctx.lineWidth = isHovered ? 2.5 : 1.5;
    ctx.stroke();

    // Slice label if angle is large enough
    if (rawAngle > 0.18) {
      const midAngle = (curAngle + endAngle) / 2;
      const labelRadius = (r1 + currentR2) / 2;
      const lx = cx + Math.cos(midAngle) * labelRadius;
      const ly = cy + Math.sin(midAngle) * labelRadius;

      ctx.save();
      ctx.translate(lx, ly);
      let rot = midAngle;
      if (rot > Math.PI / 2 && rot < (3 * Math.PI) / 2) {
        rot += Math.PI;
      }
      ctx.rotate(rot);
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.fillStyle = '#ffffff';
      ctx.font = 'bold 11px -apple-system, sans-serif';
      let labelText = item.name;
      if (labelText.length > 11) labelText = labelText.substring(0, 9) + '..';
      ctx.fillText(labelText, 0, 0);
      ctx.restore();
    }

    ctx.restore();

    sunburstSlices.push({
      index: idx,
      item: item,
      startAngle: curAngle,
      endAngle: endAngle,
      r1: r1,
      r2: currentR2
    });

    curAngle = endAngle;
  });
}

function adjustAlpha(hexColor, alpha) {
  let c = hexColor.replace('#', '');
  if (c.length === 3) c = c.split('').map(x => x + x).join('');
  const num = parseInt(c, 16);
  const r = (num >> 16) & 255;
  const g = (num >> 8) & 255;
  const b = num & 255;
  return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}

function handleSunburstMouseMove(e) {
  const canvas = document.getElementById('dirtree-sunburst-canvas');
  const hoverPill = document.getElementById('sunburst-hover-pill');
  if (!canvas || !sunburstCurrentTree) return;

  const rect = canvas.getBoundingClientRect();
  const x = e.clientX - rect.left;
  const y = e.clientY - rect.top;
  const cx = canvas.width / 2;
  const cy = canvas.height / 2;
  const dx = x - cx;
  const dy = y - cy;
  const dist = Math.sqrt(dx * dx + dy * dy);

  let prevHovered = sunburstHovered;
  sunburstHovered = null;

  if (dist <= 80) {
    // Hovered center
    sunburstHovered = { isCenter: true };
    if (hoverPill) {
      hoverPill.textContent = sunburstCurrentTree.parentPath ? '▲ Tıklayarak üst klasöre çıkın' : '📍 ' + (sunburstCurrentTree.currentPath || '~');
    }
  } else if (dist >= 95 && dist <= 260) {
    let angle = Math.atan2(dy, dx);
    // Normalize angle to [-PI/2, 3PI/2]
    if (angle < -Math.PI / 2) angle += 2 * Math.PI;

    for (const slice of sunburstSlices) {
      if (angle >= slice.startAngle && angle <= slice.endAngle) {
        sunburstHovered = slice;
        if (hoverPill) {
          const item = slice.item;
          const typeIcon = item.isDir ? '📁' : '📄';
          hoverPill.innerHTML = `${typeIcon} <b>${escapeHtml(item.name)}</b> — ${item.sizeStr} (${item.percentage ? item.percentage.toFixed(1) : 0}%) ${item.isDir ? '<small style="opacity:0.8;">[Açmak için tıkla]</small>' : ''}`;
        }
        break;
      }
    }
  }

  if (!sunburstHovered && hoverPill) {
    hoverPill.textContent = 'Klasörlerin üzerine gelin veya içine girmek için tıklayın';
  }

  if (prevHovered !== sunburstHovered) {
    drawSunburst(sunburstCurrentTree);
  }
}

function handleSunburstMouseLeave() {
  sunburstHovered = null;
  const hoverPill = document.getElementById('sunburst-hover-pill');
  if (hoverPill) {
    hoverPill.textContent = 'Klasörlerin üzerine gelin veya içine girmek için tıklayın';
  }
  if (sunburstCurrentTree) {
    drawSunburst(sunburstCurrentTree);
  }
}

function handleSunburstClick(e) {
  if (!sunburstHovered || !sunburstCurrentTree) return;

  if (sunburstHovered.isCenter) {
    if (sunburstCurrentTree.parentPath) {
      loadDirTree(sunburstCurrentTree.parentPath);
    }
  } else if (sunburstHovered.item) {
    const item = sunburstHovered.item;
    if (item.isDir) {
      loadDirTree(item.path);
    } else {
      revealInFinder(item.path);
    }
  }
}

