//cross platform build script

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");

const buildTargets = [
  { os: "linux", arch: "amd64", output: "menace-go-linux" },
  { os: "darwin", arch: "amd64", output: "menace-go-darwin" },
  { os: "darwin", arch: "arm64", output: "menace-go-darwin-arm64" },
  { os: "windows", arch: "amd64", output: "menace-go-win.exe" },
];

const goSrc     = path.resolve(__dirname, "../src/main.go");
const goModDir  = path.resolve(__dirname, "../src");      // NEW
const binDir    = path.resolve(__dirname, "../bin");

// build all targets
for (const {os: GOOS, arch: GOARCH, output} of buildTargets) {
    const outPath = path.join(binDir, output);
    console.log(`Building for ${GOOS}/${GOARCH} → ${output}`);

    try {
        execSync(
          // note: still pointing at the same src/main.go
          `go build -o "${outPath}" "${goModDir}" .`,
          {
            stdio:  "inherit",
            cwd:    goModDir,                      // NEW: run in src/
            env:    { ...process.env, GOOS, GOARCH }
          }
        );
    } catch (err) {
        console.error(`❌ build failed for ${GOOS}/${GOARCH}`);
        process.exit(1);
    }
}

console.log("✅ All builds complete.");
