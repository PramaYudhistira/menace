const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");

const buildTarget = {
    os: "windows",
    arch: "amd64",
    output: "menace-go-win.exe"
};

const goModDir = path.resolve(__dirname, "../src");
const binDir = path.resolve(__dirname, "../bin");

// build windows target
const { os: GOOS, arch: GOARCH, output } = buildTarget;
const outPath = path.join(binDir, output);
console.log(`Building for ${GOOS}/${GOARCH} → ${output}`);

try {
    execSync(
        `go build -o "${outPath}" "${goModDir}" .`,
        {
            stdio: "inherit",
            cwd: goModDir,
            env: { ...process.env, GOOS, GOARCH }
        }
    );
    console.log("✅ Windows build complete.");
} catch (err) {
    console.error(`❌ Windows build failed`);
    process.exit(1);
}
