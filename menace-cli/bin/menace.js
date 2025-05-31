#!/usr/bin/env node

const { spawn } = require("child_process");
const path = require("path");
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

const flaskBinary = path.join(__dirname, '..', 'dist', 'reposerver'); // or 'reposerver.exe' for Windows

// Spin up Flask server
console.log("Starting up... (this may take a few seconds)")
const flaskServer = spawn(flaskBinary, {
  stdio: ['inherit', 'ignore', 'ignore'],
});

// Handle any errors
child.on("error", (err) => {
    console.error(`Failed to start process: ${err.message}`);
    process.exit(1);
});

// Handle when the child process exits
child.on('close', (code) => {
    console.log(`Child process exited with code ${code}`);
    flaskServer.kill();
    process.exit(1);
});

// Handle Flask server errors - uncomment below for debugging and set 
// flaskServer.stderr.on('data', (data) => {
//     const msg = data.toString();
//     process.stderr.write(msg);
//     console.log(msg);  // This will show us the actual error
// });


flaskServer.on('error', (err) => {
  process.env.FLASK_READY = "false";
  console.error(`Failed to ping flask server: ${err.message}`);
});

flaskServer.on('close', (code) => {
  console.log(`Flask server exited with code ${code}`);
  child.kill();
  process.exit(1);
});