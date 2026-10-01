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

  // Apple tab
  btnScanApple: document.getElementById('btn-scan-apple'),
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
  monitor: { title: 'Donanım & Kaynak Monitörü', sub: 'İşlemci, Bellek (RAM) ve sistem çalışma sürelerinin canlı görünümü' },
  startup: { title: 'Başlangıç Öğeleri & Arka Plan Hizmetleri', sub: 'Oturum açma uygulamaları, LaunchAgent ve arka plan servislerini yönetin' },
};

let monitorInterval = null;

// Initialize Application
document.addEventListener('DOMContentLoaded', () => {
  checkAuthStatus();
  setupNavigation();
  setupEventHandlers();
  fetchSystemStats();
  startGlobalScan();
  startLeftoversScan(true); // background initial scan for badge
  loadSmartCare(); // Initial smart care health scan
  loadStartupItems(true); // background initial scan for startup badge
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

  const eventSource = new EventSource('/api/scan?stream=true');

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
  if (state.smartCareData && state.smartCareData.totalCleanable > 0) {
    executeSmartCareClean();
    return;
  }
  if (!state.scanData) return;
  const recommendedIds = state.scanData.targets
    .filter(t => t.exists && (t.risk === 'safe' || t.risk === 'recommended') && t.size > 0)
    .map(t => t.id);

  if (recommendedIds.length === 0) {
    showToast('Temizlenebilecek önerilen dosya bulunamadı!', 'info');
    return;
  }

  confirmAndClean(recommendedIds, [], "Önerilen Güvenli Dosyalar", elements.quickCleanSize.textContent);
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

  const eventSource = new EventSource(url);
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
  elements.appleSnapshotsList.innerHTML = '<div class="loading-state">APFS anlık görüntüleri listeleniyor...</div>';
  elements.appleBackupsList.innerHTML = '<div class="loading-state">iOS yedekleri taranıyor...</div>';
  elements.appleSimulatorsList.innerHTML = '<div class="loading-state">Simülatörler taranıyor...</div>';

  try {
    const res = await fetch('/api/apple');
    if (!res.ok) throw new Error('Apple sistem verileri alınamadı');
    const data = await res.json();
    state.appleData = data;

    renderAppleSystem(data);
    showToast('Apple ve Sistem verileri tarandı.', 'success');
  } catch (err) {
    console.error(err);
    elements.appleSnapshotsList.innerHTML = `<div class="loading-state" style="color:var(--accent-rose);">${escapeHtml(err.message)}</div>`;
    showToast('Apple tarama hatası: ' + err.message, 'error');
  } finally {
    elements.btnScanApple.disabled = false;
  }
}

function renderAppleSystem(data) {
  // 1. APFS Snapshots
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

  // 2. iOS Backups
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

  // 3. Simulators
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

    if (elements.quickCleanSize && data.totalCleanableStr) {
      elements.quickCleanSize.textContent = data.totalCleanableStr;
    }

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
    if (!res.ok) return;
    const data = await res.json();
    state.authStatus = data;

    if (data.enabled) {
      if (!data.authenticated) {
        if (elements.authModal && !elements.authModal.open) {
          elements.authModal.showModal();
        }
        if (elements.authStatusBar) elements.authStatusBar.style.display = 'none';
      } else {
        if (elements.authModal && elements.authModal.open) {
          elements.authModal.close();
        }
        if (elements.authStatusBar) elements.authStatusBar.style.display = 'flex';
      }
    } else {
      if (elements.authModal && elements.authModal.open) {
        elements.authModal.close();
      }
      if (elements.authStatusBar) elements.authStatusBar.style.display = 'none';
    }
  } catch (err) {
    console.error('Auth check error:', err);
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
      elements.authModal.close();
      elements.authPasswordInput.value = '';
      if (elements.authStatusBar) elements.authStatusBar.style.display = 'flex';
      showToast('Giriş başarılı! Oturum açıldı.', 'success');
      // Refresh system storage and stats
      fetchSystemStats();
      startGlobalScan();
      startLeftoversScan(true);
      loadSmartCare();
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

    renderStartupItems();
  } catch (err) {
    console.error('Startup items error:', err);
    if (!isSilent && elements.startupItemsList) {
      elements.startupItemsList.innerHTML = `<div class="empty-state">Hata: ${escapeHtml(err.message)}</div>`;
    }
  }
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

// Setup Event Handlers
function setupEventHandlers() {
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

  // Apple tab handlers
  if (elements.btnScanApple) {
    elements.btnScanApple.addEventListener('click', () => startAppleScan());
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
