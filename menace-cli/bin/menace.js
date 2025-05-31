#!/usr/bin/env node

const { spawn } = require("child_process");
const path = require("path");
const fs = require("fs");
const { execSync } = require("child_process");
process.env.FLASK_READY = "false";
process.env.PORT = 5974;

//get binary name
let binName;
switch (process.platform) {
    case "win32":
        binName = "menace-go-win.exe";
        break;
    case "darwin":
        binName = process.arch === "arm64" ? "menace-go-darwin-arm64" : "menace-go-darwin";
        break;
    case "linux":
        binName = "menace-go-linux";
        break;
    default:
        console.error(`Unsupported platform: ${process.platform}`);
        process.exit(1);
}

//get path to binary
const binPath = path.join(__dirname, binName);

// Spawn the executable for the main terminal agent
const child = spawn(binPath, {
    stdio: 'inherit',
});

// Handle any errors
child.on("error", (err) => {
    console.error(`Failed to start process: ${err.message}`);
    process.exit(1);
});

// Handle when the child process exits
child.on('close', (code) => {
    console.log(`Child process exited with code ${code}`);
    process.exit(1);
});


// check if a venv is present
const venvPath = path.join(__dirname, '..', 'venv');
if (!fs.existsSync(venvPath)) {
    console.log("Venv not found, creating...")
    execSync('python -m venv venv');
}

// 2. Install deps using venv pip
const pipPath = path.join(venvPath, process.platform === "win32" ? "Scripts" : "bin", "pip"); // Windows: use "Scripts" instead of "bin" for pip
try {
  execSync(`${pipPath} install -r requirements.txt`, { stdio: ["inherit", "ignore", "ignore"] });
} catch (e) {
  console.error("Failed to install Python packages:", e);
  process.exit(1);
}

// test if the venv is working
// const kitPath = path.join(venvPath, "bin", "kit");
// try {
//   execSync(`${kitPath} file-tree ${__dirname}`, { stdio: "inherit" });
// } catch (e) {
//   console.error("Failed to run kit file-tree:", e);
// }