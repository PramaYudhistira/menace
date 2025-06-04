const { execSync } = require('child_process');
const path = require('path');

function getKitPath() {
    const venvPath = process.env.MENACE_VENV_PATH;
    if (!venvPath) {
        throw new Error('MENACE_VENV_PATH not set');
    }
    return path.join(venvPath, 'bin', 'kit');
}

function runKitCommand(command, args = []) {
    const kitPath = getKitPath();
    const fullCommand = `"${kitPath}" ${command} ${args.join(' ')}`;
    try {
        return execSync(fullCommand, { 
            stdio: ['inherit', 'pipe', 'pipe'],
            encoding: 'utf-8'
        });
    } catch (error) {
        console.error(`Kit command failed: ${error.message}`);
        throw error;
    }
}

// Utility functions for common kit commands
const kit = {
    fileTree: (path) => runKitCommand('file-tree', [path]),
    fileContent: (repoPath, filePath) => runKitCommand('file-content', [repoPath, filePath]),
    findSymbols: (repoPath, symbol) => runKitCommand('find-symbols', [repoPath, symbol]),
    search: (query) => runKitCommand('search', [query])
};

module.exports = kit; 