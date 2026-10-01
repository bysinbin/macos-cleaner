# 🍏 DiskCleaner Pro (macOS)

Go ile geliştirilmiş, macOS için modern, yüksek performanslı ve görsel disk temizleme ve geliştirici önbellek analiz aracı.

---

## ✨ Özellikler

1. **Macintosh HD Depolama Analizi & Genel Bakış**:
   - Diskinizin toplam, kullanılan ve boş alanını canlı gösteren dairesel (donut) depolama grafiği.
   - Tek tıkla "Önerilen Güvenli Temizlik" (Quick Clean).
2. **🗑️ Artık Dosya Avcısı (Residual & Leftover Cleaner - AppCleaner Tarzı)**:
   - **Kaldırılmış Uygulama Artıkları**: Çöp kutusuna atılıp silinmiş fakat `~/Library/Application Support` ve `Saved Application State` dizinlerinde unutulmuş tüm dosya ve klasörleri tespit eder.
   - **.DS_Store Temizleyici**: macOS Finder'ın klasörlerde bıraktığı gereksiz `.DS_Store` dosyalarını toplu olarak temizler.
   - **Kırık Sembolik Linkler**: Hedefi bulunmayan geçersiz linkleri yakalar.
   - **Geçici & Yedek Dosyalar**: Editörlerin geride bıraktığı `*.tmp`, `*.swp`, `*~` artıkları.
3. **macOS & Geliştirici Önbellekleri (Developer Caches)**:
   - **Xcode SwiftUI Previews**: Diskinizde gigabaytlarca yer kaplayan önizleme önbellekleri.
   - **Xcode DerivedData** & iOS Simülatör geçici dosyaları.
   - **Arduino IDE İndirme ve Paket Deposu** (`~/Library/Arduino15/staging` — ~2.9 GB).
   - **Android Emülatör Sistem İmajları** (`~/Library/Android/sdk/system-images` — ~4.1 GB).
   - **VS Code Eklenti Kurulum Paketleri** (`CachedExtensionVSIXs` — ~660 MB).
   - **Gradle Önbellekleri** (`~/.gradle/caches`).
   - **Homebrew İndirme Arşivleri** (`~/Library/Caches/Homebrew`).
   - **Go Derleme Önbelleği** (`go-build` & modül cache).
   - **npm, Yarn, pnpm, CocoaPods, Rust Cargo, Python pip** önbellekleri.
4. **Sistem ve Uygulama Artıkları**:
   - **VS Code ShipIt** eski güncelleme paketleri.
   - **Google Chrome** Service Worker ve GPUCache.
   - **Discord, Steam, Postman** medya ve geçici önbellekleri.
   - **Kullanıcı ve Uygulama Logları** (`~/Library/Logs`).
   - **macOS Çöp Kutusu** (`~/.Trash`).
   - `~/Library/Caches` altındaki diğer tüm 10 MB+ uygulamaların dinamik tespiti.
5. **📂 Derin Dizin Analizi (Disk Space Tree)**:
   - DaisyDisk / ncdu benzeri hiyerarşik klasör boyut analizi.
   - Oransal renkli ilerleme çubukları ile hangi klasörün ne kadar yer kapladığını anında görme.
   - Tıklayarak içine girme (drill-down), üst klasöre çıkma, Finder'da açma ve doğrudan silme.
6. **👯 Yinelenen Dosya Avcısı (Duplicate Finder)**:
   - Seçilen klasörlerdeki (Masaüstü, İndirilenler, Belgeler) aynı boyuta ve SHA-256 hash imzasına sahip kopya dosyaları bulur.
   - "Orijinali Koru, Kopyaları Seç" akıllı seçimi ile gereksiz kopyaları tek tıkla temizleme.
7. **node_modules Avcısı**:
   - Masaüstü, Belgeler veya özel proje klasörünü tarar.
   - Atıl kalmış devasa `node_modules` klasörlerini listeler, Finder'da açabilir veya silebilirsiniz.
8. **Büyük & Eski Dosya Gezgini**:
   - Boyut (>50MB, >100MB, >500MB, >1GB), dosya türü (Videolar, Kalıplar/DMG, Arşivler, Belgeler) ve yaş filtreleri.
   - Dosyayı Finder'da doğrudan gösterme (`open -R`).
   - Güvenli şekilde Çöp Kutusuna taşıma veya kalıcı silme.
9. **🚀 Uygulama Kaldırıcı & Sıfırlayıcı (App Uninstaller & App Reset)**:
   - macOS üzerindeki tüm yüklü uygulamaları (.app ikili dosyası) ve ilişkili kullanıcı verilerini (Application Support, Caches, Preferences, Saved State, Containers, WebKit vb.) ayrı ayrı hesaplar.
   - **Sıfırla (Reset) Modu**: Uygulamayı korurken önbellek ve kullanıcı ayarlarını sıfırlayarak fabrika ayarlarına döndürür.
   - **Kaldır (Uninstall) Modu**: Uygulama paketini ve arkasındaki tüm gizli artıkları tek tıkla Finder Çöp Kutusuna taşır veya kalıcı olarak kaldırır.
   - Sistem uygulamaları otomatik olarak koruma altındadır.
10. **🌐 Tarayıcı Önbellekleri & Gizlilik (Browser Caches)**:
    - Chrome, Safari, Firefox, Edge, Brave ve Arc tarayıcılarının web görseli, derlenmiş JavaScript kodları ve GPU önbelleklerini temizler.
    - Şifreler ve yer imleri asla silinmez.
11. **📥 İndirilenler & Kurulum Kalıpları Analizi (Downloads Organizer)**:
    - `~/Downloads` dizinindeki disk kalıplarını (`.dmg`, `.pkg`, `.iso`), büyük arşivleri (`.zip`, `.rar`, `.tar.gz`) ve 30 günden eski dosyaları otomatik kategorize eder.
12. **🍏 Apple & Sistem Verisi (APFS Snapshots, iOS Yedekleri & Simülatörler)**:
    - **Time Machine APFS Anlık Görüntüleri**: Diski dolduran yerel APFS snapshot'larını listeleme ve `tmutil` ile temizleme.
    - **iPhone & iPad Yerel Yedekleri**: `MobileSync/Backup` altında unutulan devasa aygıt yedeklerini cihaz adı ve tarihiyle tespit etme.
    - **Xcode iOS Simülatörleri**: Kullanılmayan simülatör aygıtlarını ve runtime önbelleklerini `xcrun simctl` ile temizleme.
13. **💬 İletişim & Mesaj Ekleri (Messages & Mail Attachments)**:
    - iMessage ve Apple Mail ile diskte biriken videolar, yüksek çözünürlüklü fotoğraflar, sesler ve indirilen belgeleri listeler ve temizler.
15. **⚡ Akıllı Bakım (Smart Care) & Disk Eşik Uyarısı (Threshold Alert)**:
    - Diskin doluluk oranı %85'i aştığında Genel Bakış panelinde nabız atan uyarı barı (`⚠️ Kritik Disk Doluluğu`) gösterilir.
    - Tek tıkla güvenli tüm önbellek, çöp, log ve artıkları temizleyen Akıllı Bakım motoru.
16. **🛠️ Sistem Bakımı & Hızlandırma (CleanMyMac Maintenance Toolkit)**:
    - **Inaktif RAM Belleği Boşaltma**: `purge` ile bellekte tutulan atıl önbellekleri anında serbest bırakır.
    - **DNS Önbelleğini Sıfırlama**: `dscacheutil` ve `mDNSResponder` servislerini yenileyerek internet takılmalarını çözer.
    - **Apple Mail Veritabanı Optimizasyonu**: Mail `Envelope Index` SQLite veritabanlarını `VACUUM` ile sıkıştırıp hızlandırır.
    - **Launch Services & 'Birlikte Aç' Menüsü Onarımı**: Çift veya bozuk kayıtları sıfırlayarak onarır.
    - **Spotlight Arama İndeksini Yeniden Oluşturma**: `mdutil -E /` ile arama dizinini baştan oluşturur.
    - **Periyodik macOS Sistem Betikleri**: Günlük, haftalık ve aylık bakım betiklerini tetikler.
    - **Başlangıç Diskini Doğrulama**: `diskutil` ile APFS dosya sistemi sağlığını kontrol eder.
17. **🔥 Güvenli Dosya Öğütücü (File Shredder - DoD 5220.22-M)**:
    - Disk Drill, PhotoRec gibi kurtarma yazılımlarıyla dahi asla geri getirilemeyecek şekilde çoklu geçişle (DoD 3-pass & 1-pass) rastgele bayt yazma, sıfırlama ve dosya tablosu parçalama.
18. **📊 Donanım & Kaynak Monitörü (Sensei Tarzı Hardware Monitor)**:
    - Canlı İşlemci (CPU) modeli, çekirdek sayısı ve kullanım yüzdesi.
    - Canlı Bellek (RAM) toplam, kullanılan, boş, kablolu (wired) ve sıkıştırılmış (compressed) detayları.
    - macOS sürümü, çalışma süresi (uptime) ve sistem kimliği.
19. **🔒 Güvenlik & Yalnızca Yerel Erişim Garantisi**:
    - **Tamamen Yerel (Localhost Isolation)**: Sunucu yalnızca `127.0.0.1` adresine bağlanır. Güvenlik katmanı harici ağlardan (`192.168.x.x`, `10.x.x.x` vb.) gelen tüm istekleri `403 Forbidden` ile reddeder.
    - **Şifre Korumalı Giriş (`config.json`)**: Uygulamaya girişte parola koruması sağlanır; güvenli oturum tokenı (`dc_token`, HttpOnly) ile korunur.
    - **Test (Dry Run) & Çöp Kutusu Modları**: Dosyaları kalıcı silmek yerine Finder Çöp Kutusuna taşıma veya sadece test çalıştırma opsiyonu.
    - **Kritik Sistem Dizin Koruması**: Sistem kök dizinleri (`/System`, `/usr`, `/Library`, `/Applications`, `~/.ssh`) otomatik koruma altındadır.
20. **🚀 Başlangıç Öğeleri & Arka Plan Hizmetleri Yöneticisi (CleanMyMac Optimization Tarzı)**:
    - **Kullanıcı Oturum Açma Öğeleri (Login Items)**: Mac açıldığında otomatik başlayan uygulamaları listeleme, gizli başlatma, listeden kaldırma, kalıcı pasife alma ve yeni `.app` ekleme. (TouchSwitcher ve benzeri masaüstü uygulamaları için kalıcı durum yönetimi).
    - **Kullanıcı & Sistem Başlatma Ajanları (LaunchAgents)**: Oturumla başlayan ajanları modern Aç/Kapa anahtarlarıyla anında etkinleştirme veya devre dışı bırakma (AdobeCreativeCloud vb. izin hataları olmadan doğrudan `launchctl` desteği ile).
    - **Sistem Arka Plan Hizmetleri (LaunchDaemons)**: Sistem düzeyinde çalışan servisleri (CleanMyMac, PostgreSQL, Warp vb.) listeleme ve yönetme.
21. **🧩 Sistem ve Uygulama Eklenti Yöneticisi (Extensions Manager)**:
    - **Sistem Sürücü & Ağ Eklentileri (DriverKit & NetworkExtension)**: `systemextensionsctl` entegrasyonu ile VPN filtreleri, USB/Seri sürücüleri ve durumlarını görüntüleme.
    - **Finder & Paylaşım Eklentileri (PluginKit Appex)**: Finder senkronizasyon (Creative Cloud, UnzipOne vb.), paylaşım menüsü ve widget eklentilerini tek tıkla Aç/Kapat yapabilme.
    - **Hızlı Bakış (QuickLook) & Spotlight**: Sistem ve kullanıcı dizinlerindeki `.qlgenerator` ve `.mdimporter` eklentilerini yönetme.
22. **🔋 Genişletilmiş Donanım Monitörü (Pil, GPU, Disk I/O)**:
    - **Pil & Güç Durumu**: Pil şarj yüzdesi, şarj durumu, pil sağlığı (Health %), devir sayısı (Cycle Count), pil sıcaklığı (°C), güç kaynağı (AC/Pil) ve kalan süre.
    - **Grafik İşlemci (GPU)**: Apple Silicon GPU modeli, çekirdek sayısı, Metal sürümü ve ekran çözünürlüğü.
    - **Depolama & Disk I/O**: NVMe SSD dosya sistemi (APFS), S.M.A.R.T. sağlık durumu (Verified), anlık disk okuma/yazma I/O hızı ve TPS.
23. **🌐 Ağ Monitörü (Network Monitor & Active Connections)**:
    - **Ağ Adaptörü & IP Yapılandırması**: Aktif arayüz (en0 / Wi-Fi), Yerel IPv4 adresi, Donanım (MAC) adresi, Ağ Geçidi (Gateway) ve DNS sunucuları.
    - **Anlık Trafik & Hız**: Gerçek zamanlı İndirme (Download) ve Yükleme (Upload) hız göstergesi ve toplam aktarılan veri sayaçları.
    - **Aktif TCP Bağlantıları Tablosu**: Sistemdeki uygulamaların (AnyDesk, Chrome, Node vb.) aktif TCP soketleri, PID, port ve uzak sunucu bağlantı listesi.
24. **🛡️ Gizlilik Koruması & İz Temizleyici (Privacy Protection)**:
    - **Son Açılan Belgeler & Sunucular**: `sfltool` entegrasyonuyla Finder ve sistem son kullanılanlar listesini temizleme.
    - **Terminal Komut Geçmişi**: `.zsh_history`, `.bash_history`, Python ve Node REPL komut kayıtlarını görüntüleme ve tek tıkla güvenle temizleme.
    - **Tarayıcı Gizlilik Verileri**: Safari, Chrome, Arc, Brave ve Firefox yerel gezinme veritabanları ve çerez izlerini kaldırma.
    - **macOS Uygulama İzin Sıfırlayıcı (TCC)**: Kamera, Mikrofon, Ekran Kaydı, Erişilebilirlik ve Tam Disk yetkilerini tek tıkla sıfırlama (`tccutil reset`).

---

## ⚙️ Yapılandırma (`config.json`)

Uygulama kök dizinindeki `config.json` ile port ve erişim şifresini belirleyebilirsiniz:

```json
{
  "port": 8089,
  "bindAddress": "127.0.0.1",
  "auth": {
    "enabled": true,
    "password": "admin"
  }
}
```

* `enabled`: Şifre korumasını etkinleştirir veya devre dışı bırakır (`true` / `false`).
* `password`: Panele erişim parolası (varsayılan: `admin`).
* `bindAddress`: Güvenlik için daima `127.0.0.1` olmalıdır.

---

## 🚀 Çalıştırma

Terminalde proje dizinindeyken:

```bash
# Doğrudan çalıştırma (Otomatik olarak tarayıcınızda açılır):
./disk-cleaner

# Veya Go ile çalıştırma:
go run .

# Farklı port belirleme:
./disk-cleaner -port 9090
```

Çalıştırıldığında tarayıcınızda otomatik olarak **`http://127.0.0.1:8089`** adresinde açılacaktır.

---

## 🛠️ Yeniden Derleme

Tüm web arayüzü (HTML, CSS, JS) Go'nun `embed.FS` özelliği ile tek bir çalıştırılabilir ikili dosyaya (binary) gömülmüştür:

```bash
go build -o disk-cleaner .
```

