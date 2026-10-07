#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-v2.0.0}"
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%d")

echo "=================================================="
echo " 🍏 Building DiskCleaner Pro Universal ($VERSION)"
echo "=================================================="

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="$ROOT_DIR/build"
APP_DIR="$BUILD_DIR/DiskCleaner.app"

rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR" "$APP_DIR/Contents/MacOS" "$APP_DIR/Contents/Resources"

LDFLAGS="-s -w -X disk-cleaner/internal/version.Version=$VERSION -X disk-cleaner/internal/version.GitCommit=$COMMIT -X disk-cleaner/internal/version.BuildDate=$BUILD_DATE"

echo "-> Compiling for darwin/arm64 (Apple Silicon)..."
CGO_ENABLED=1 CC="clang -target arm64-apple-macos11.0" GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/disk-cleaner-arm64" "$ROOT_DIR"

echo "-> Compiling for darwin/amd64 (Intel Mac)..."
CGO_ENABLED=1 CC="clang -target x86_64-apple-macos11.0" GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BUILD_DIR/disk-cleaner-amd64" "$ROOT_DIR"

echo "-> Creating universal macOS binary with lipo..."
lipo -create -output "$APP_DIR/Contents/MacOS/DiskCleaner" "$BUILD_DIR/disk-cleaner-arm64" "$BUILD_DIR/disk-cleaner-amd64"
chmod +x "$APP_DIR/Contents/MacOS/DiskCleaner"

# Also copy universal binary to build directory root
cp "$APP_DIR/Contents/MacOS/DiskCleaner" "$BUILD_DIR/disk-cleaner"

if [ -f "$ROOT_DIR/assets/AppIcon.icns" ]; then
    echo "-> Installing macOS AppIcon.icns..."
    cp "$ROOT_DIR/assets/AppIcon.icns" "$APP_DIR/Contents/Resources/AppIcon.icns"
fi

echo "-> Generating Info.plist..."
cat << 'EOF' > "$APP_DIR/Contents/Info.plist"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>DiskCleaner</string>
    <key>CFBundleIdentifier</key>
    <string>com.bysinbin.diskcleaner</string>
    <key>CFBundleName</key>
    <string>DiskCleaner Pro</string>
    <key>CFBundleDisplayName</key>
    <string>DiskCleaner Pro</string>
    <key>CFBundleIconFile</key>
    <string>AppIcon</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>2.0.0</string>
    <key>CFBundleVersion</key>
    <string>2.0.0</string>
    <key>LSMinimumSystemVersion</key>
    <string>11.0</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSRequiresAquaSystemAppearance</key>
    <false/>
    <key>NSSystemExtensionUsageDescription</key>
    <string>DiskCleaner Pro requires system access to inspect caches and clean residual files.</string>
</dict>
</plist>
EOF

# Update version in Info.plist
sed -i '' "s/2.0.0/${VERSION#v}/g" "$APP_DIR/Contents/Info.plist"

echo "-> Packaging distribution archives..."
cd "$BUILD_DIR"
zip -r -q "DiskCleaner-${VERSION}-macos-app.zip" "DiskCleaner.app"
tar -czf "disk-cleaner-${VERSION}-darwin-universal.tar.gz" "disk-cleaner"
tar -czf "disk-cleaner-${VERSION}-darwin-arm64.tar.gz" -C "$BUILD_DIR" "disk-cleaner-arm64"
tar -czf "disk-cleaner-${VERSION}-darwin-amd64.tar.gz" -C "$BUILD_DIR" "disk-cleaner-amd64"

# Generate SHA256 checksums
shasum -a 256 *.zip *.tar.gz > "checksums.txt"

echo "=================================================="
echo " ✅ Release build completed successfully!"
echo " Output files in: $BUILD_DIR"
ls -lh "$BUILD_DIR"
echo "=================================================="
