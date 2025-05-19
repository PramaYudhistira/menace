#!/usr/bin/env node

const { spawn } = require("child_process");
const path = require("path");

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

// Spawn the executable
const child = spawn(binPath, [], { stdio: "inherit" });

child.on("error", (err) => {
    console.error(`Failed to start process: ${err.message}`);
    process.exit(1);
});
