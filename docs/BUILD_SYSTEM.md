# Menace CLI Build System Documentation

## Overview

The Menace CLI build system is designed for cross-platform distribution, combining Go's excellent cross-compilation capabilities with Node.js packaging for easy installation and distribution.

## Architecture

### Hybrid Build Approach

The project uses a two-tier architecture:

1. **Go Binary Layer**: Core application compiled for multiple platforms
2. **Node.js Wrapper Layer**: Platform detection and process management

```
User runs `menace`
        ↓
Node.js wrapper (menace.js)
        ↓ 
Platform detection
        ↓
Spawn appropriate Go binary
        ↓
menace-go-{platform}
```

### Benefits of This Approach

- **Simple Distribution**: Single npm package works on all platforms
- **No Runtime Dependencies**: Go binaries are statically linked
- **Automatic Platform Detection**: Users don't need to choose the right binary
- **Familiar Installation**: Uses npm, familiar to developers
- **Small Package Size**: Only includes necessary binaries per platform

## Build Scripts

### Primary Build Scripts

Located in `scripts/` directory:

#### 1. `build-all.js` - Complete Multi-Platform Build

```javascript
const buildTargets = [
  { os: "linux", arch: "amd64", output: "menace-go-linux" },
  { os: "darwin", arch: "amd64", output: "menace-go-darwin" }, 
  { os: "darwin", arch: "arm64", output: "menace-go-darwin-arm64" },
  { os: "windows", arch: "amd64", output: "menace-go-win.exe" },
];
```

**Usage:**
```bash
npm run build
```

**Process:**
1. Iterates through all build targets
2. Sets `GOOS` and `GOARCH` environment variables
3. Executes `go build` for each platform
4. Outputs binaries to `bin/` directory
5. Exits with error code 1 if any build fails

#### 2. `build-mac.js` - macOS-Specific Build

Builds both Intel and Apple Silicon binaries:
- `menace-go-darwin` (Intel)
- `menace-go-darwin-arm64` (Apple Silicon)

**Usage:**
```bash
npm run build-mac
```

#### 3. `build-windows.js` - Windows-Specific Build

Builds Windows executable with `.exe` extension:
- `menace-go-win.exe`

**Usage:**
```bash
npm run build-windows
```

### Build Configuration

#### Go Module Configuration

**File:** `src/go.mod`

```go
module menace-go

go 1.24.1

require (
    github.com/charmbracelet/bubbletea v1.3.5
    github.com/charmbracelet/lipgloss v1.1.0
    // ... other dependencies
)
```

#### Build Environment Variables

| Variable | Purpose | Example |
|----------|---------|---------|
| `GOOS` | Target operating system | `linux`, `darwin`, `windows` |
| `GOARCH` | Target architecture | `amd64`, `arm64` |
| `CGO_ENABLED` | CGO support (default: 0 for static linking) | `0`, `1` |

### Platform-Specific Considerations

#### Linux Build
- **Target**: `linux/amd64`
- **Output**: `menace-go-linux` 
- **Compatible with**: Most Linux distributions (Ubuntu, CentOS, Debian, etc.)
- **Static linking**: No external dependencies required

#### macOS Build
- **Targets**: 
  - `darwin/amd64` → `menace-go-darwin`
  - `darwin/arm64` → `menace-go-darwin-arm64`
- **Compatibility**: 
  - Intel: macOS 10.15+
  - Apple Silicon: macOS 11.0+
- **Code signing**: Not currently implemented (consider for future)

#### Windows Build
- **Target**: `windows/amd64`
- **Output**: `menace-go-win.exe`
- **Compatibility**: Windows 10+, Windows Server 2019+
- **Console support**: Windows Terminal recommended

## Node.js Wrapper

### Wrapper Script (`bin/menace.js`)

The Node.js wrapper handles:

1. **Platform Detection**:
   ```javascript
   switch (process.platform) {
       case "win32": binName = "menace-go-win.exe"; break;
       case "darwin": binName = process.arch === "arm64" ? 
           "menace-go-darwin-arm64" : "menace-go-darwin"; break;
       case "linux": binName = "menace-go-linux"; break;
   }
   ```

2. **Process Spawning**:
   ```javascript
   const child = spawn(binPath, [], { stdio: "inherit" });
   ```

3. **Error Handling**:
   ```javascript
   child.on("error", (err) => {
       console.error(`Failed to start process: ${err.message}`);
       process.exit(1);
   });
   ```

### Package.json Configuration

**Binary Declaration:**
```json
{
  "bin": {
    "menace": "./bin/menace.js"
  },
  "scripts": {
    "build": "node scripts/build-all.js",
    "build-windows": "node scripts/build-windows.js", 
    "build-mac": "node scripts/build-mac.js"
  }
}
```

## Build Process Deep Dive

### Local Development Build

```bash
# Option 1: Run from source (fastest for development)
cd src
go run .

# Option 2: Build for current platform only
cd src
go build -o ../bin/menace-dev .

# Option 3: Full multi-platform build
npm run build
```

### Production Build Process

1. **Clean Previous Builds**:
   ```bash
   rm -rf bin/menace-go-*
   ```

2. **Build All Platforms**:
   ```bash
   npm run build
   ```

3. **Verify Binaries**:
   ```bash
   ls -la bin/
   # Should show:
   # menace-go-darwin
   # menace-go-darwin-arm64
   # menace-go-linux
   # menace-go-win.exe
   ```

4. **Test Wrapper**:
   ```bash
   node bin/menace.js
   ```

### Build Optimization

#### Binary Size Optimization

Current approach:
```bash
go build -o output .
```

**Optimized approach** (consider implementing):
```bash
go build -ldflags="-s -w" -o output .
```

Flags explained:
- `-s`: Strip symbol table and debug info
- `-w`: Strip DWARF debug info
- Results in ~30% smaller binaries

#### Build Speed Optimization

**Parallel builds** (potential improvement):
```javascript
// Instead of sequential builds
for (const target of buildTargets) {
    buildTarget(target);
}

// Parallel builds
const buildPromises = buildTargets.map(target => 
    buildTargetAsync(target)
);
await Promise.all(buildPromises);
```

## CI/CD Integration

### GitHub Actions Configuration

```yaml
name: Build and Test
on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.24.1]
        node-version: [18.x, 20.x]
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Set up Node.js
      uses: actions/setup-node@v3
      with:
        node-version: ${{ matrix.node-version }}
    
    - name: Install dependencies
      run: |
        cd src
        go mod download
    
    - name: Build all platforms
      run: npm run build
    
    - name: Upload binaries
      uses: actions/upload-artifact@v3
      with:
        name: binaries
        path: bin/
```

### Release Process

```yaml
name: Release
on:
  release:
    types: [created]

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Build
      run: npm run build
    
    - name: Publish to npm
      run: |
        echo "//registry.npmjs.org/:_authToken=${{ secrets.NPM_TOKEN }}" > ~/.npmrc
        npm publish
    
    - name: Upload release assets
      uses: actions/upload-release-asset@v1
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      with:
        upload_url: ${{ github.event.release.upload_url }}
        asset_path: ./bin/
        asset_name: menace-cli-binaries.zip
        asset_content_type: application/zip
```

## Troubleshooting

### Common Build Issues

#### 1. "go: cannot find main module"
**Cause**: Running build from wrong directory
**Solution**: Ensure you're in the project root or `src/` directory

#### 2. "GOOS/GOARCH not supported"
**Cause**: Invalid platform combination
**Solution**: Check supported combinations:
```bash
go tool dist list
```

#### 3. "permission denied" (Linux/macOS)
**Cause**: Binary not executable
**Solution**: 
```bash
chmod +x bin/menace-go-*
```

#### 4. Build fails with "module not found"
**Cause**: Dependencies not downloaded
**Solution**:
```bash
cd src
go mod download
go mod tidy
```

### Platform-Specific Issues

#### macOS Code Signing
**Issue**: "cannot be opened because the developer cannot be verified"
**Temporary Solution**: 
```bash
sudo spctl --add bin/menace-go-darwin
```
**Long-term Solution**: Implement code signing in build process

#### Windows Antivirus
**Issue**: Antivirus flags the executable
**Solution**: 
- Build with official Go compiler
- Consider code signing certificate
- Add documentation about false positives

#### Linux Distribution Compatibility
**Issue**: GLIBC version conflicts
**Current Approach**: Static linking avoids this issue
**Verification**: Test on different distributions

## Future Improvements

### Build System Enhancements

1. **Automated Testing**: Run tests on all platforms before build
2. **Code Signing**: Implement for macOS and Windows
3. **Compression**: Use UPX or similar for smaller binaries
4. **Incremental Builds**: Only rebuild changed platforms
5. **Build Caching**: Cache dependencies and intermediate artifacts

### Distribution Improvements

1. **Multiple Package Managers**: Support Homebrew, Chocolatey, etc.
2. **Docker Images**: Containerized distribution
3. **Snap Packages**: Linux universal packages
4. **AppImage**: Linux portable applications

### Monitoring and Analytics

1. **Build Metrics**: Track build times and success rates
2. **Binary Size Tracking**: Monitor binary size growth
3. **Platform Usage**: Analytics on which platforms are most used
4. **Performance Metrics**: Startup time and memory usage per platform

## Build Script Reference

### Environment Setup

```bash
# Set up development environment
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Install build dependencies
go install github.com/goreleaser/goreleaser@latest
npm install -g npm@latest
```

### Manual Build Commands

```bash
# Build specific platform
GOOS=linux GOARCH=amd64 go build -o bin/menace-go-linux src/

# Build with optimization
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/menace-go-linux src/

# Build with debug info
GOOS=linux GOARCH=amd64 go build -gcflags="all=-N -l" -o bin/menace-go-linux-debug src/

# Cross-compile all platforms manually
for GOOS in darwin linux windows; do
    for GOARCH in amd64 arm64; do
        if [[ $GOOS == "windows" && $GOARCH == "arm64" ]]; then continue; fi
        EXT=""
        if [[ $GOOS == "windows" ]]; then EXT=".exe"; fi
        GOOS=$GOOS GOARCH=$GOARCH go build -o bin/menace-go-${GOOS}-${GOARCH}${EXT} src/
    done
done
```

This build system provides a robust foundation for cross-platform distribution while maintaining simplicity and reliability.