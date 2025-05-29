#!/usr/bin/env node

const { spawn } = require("child_process");
const path = require("path");
let flaskReady = false;

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

function startChildProcess() {
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
    });
}

// Grab the bin and spin it up!
const flaskBinary = path.join(__dirname, '..', 'dist', 'reposerver'); // or 'reposerver.exe' for Windows

console.log("Starting up... (this may take a few seconds)")
const flaskServer = spawn(flaskBinary, {
  stdio: ['inherit', 'pipe', 'pipe'],
});

flaskServer.stdout.on('data', (data) => {
    const msg = data.toString();
    process.stdout.write(msg);
    if (msg.includes("FLASK SERVER READY")) {
      flaskReady = true;
      startChildProcess();
      // close this listener after the first message
      flaskServer.stdout.removeListener('data', onData);
    }
});

flaskServer.on('close', (code) => {
  console.log(`Flask server exited with code ${code}`);
});