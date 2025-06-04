#!/usr/bin/env node

const { spawn } = require("child_process");
const path = require("path");
const { execSync } = require("child_process");
process.env.FLASK_READY = "false";
process.env.PORT = 5974;

// check if theres an global env var named "MENACE_VENV_PATH", if not then create a venv in bink
if (!process.env.MENACE_VENV_PATH) {
    const venvPath = path.join(__dirname, ".venv");
    execSync(`python -m venv "${venvPath}"`, { stdio: ["pipe", "inherit", "inherit"] });
    process.env.MENACE_VENV_PATH = venvPath;
    
    // Install deps using in venv pip
    try {
        const pipPath = path.join(process.env.MENACE_VENV_PATH, 
            process.platform === "win32" ? "Scripts\\pip" : "bin/pip");
        const requirementsPath = path.join(__dirname, "..", "requirements.txt");
        
        // First upgrade pip
        execSync(`"${pipPath}" install --upgrade pip`, 
            { stdio: ["inherit", "inherit", "inherit"] });
            
        // Install build tools first
        execSync(`"${pipPath}" install --upgrade setuptools wheel`, 
            { stdio: ["inherit", "inherit", "inherit"] });

        // Install requirements from requirements.txt
        execSync(`"${pipPath}" install -r "${requirementsPath}"`, 
            { stdio: ["inherit", "inherit", "inherit"] });

        // install completion using the full path to kit
        const kitPath = path.join(process.env.MENACE_VENV_PATH, 'bin', 'kit');
        execSync(`"${kitPath}" --install-completion`, 
            { stdio: ["inherit", "inherit", "inherit"] });
            
        // Add kit to PATH
        process.env.PATH = `${path.join(process.env.MENACE_VENV_PATH, 'bin')}:${process.env.PATH}`;
    } catch (e) {
        console.error("Failed to install Python packages:", e);
        process.exit(1);
    }
    
    // Add to shell config
    const shellConfig = process.env.SHELL?.includes('zsh') ? '~/.zshrc' : '~/.bashrc';
    const exportCmd = `echo 'export MENACE_VENV_PATH="${venvPath}"' >> ${shellConfig}`;
    execSync(exportCmd);
}

// Instead of using source (which doesn't work in Node), add the venv bin to PATH
const venvBinPath = path.join(process.env.MENACE_VENV_PATH, 'bin');
process.env.PATH = `${venvBinPath}:${process.env.PATH}`;

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
const child = spawn(binPath, [], {
    stdio: 'inherit',
    detached: false,
    shell: true
});

// Handle any errors
child.on("error", (err) => {
    console.error(`Failed to start process: ${err.message}`);
    process.exit(1);
});

// Handle when the child process exits
child.on('close', (code) => {
    console.log(`Child process exited with code ${code}`);
    process.exit(code);
});