const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");

const buildTargets = [
    { os: "darwin", arch: "amd64", output: "menace-go-darwin" },
    { os: "darwin", arch: "arm64", output: "menace-go-darwin-arm64" }
];

const goModDir = path.resolve(__dirname, "../src");
const binDir = path.resolve(__dirname, "../bin");

for (const { os: GOOS, arch: GOARCH, output } of buildTargets) {
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
        console.log(`✅ ${output} build complete.`);
    } catch (err) {
        console.error(`❌ ${output} build failed`);
        process.exit(1);
    }
} 