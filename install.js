#!/usr/bin/env node

/**
 * XyPriss CLI Installer
 *
 * This script automatically downloads the appropriate XyPCLI binary
 * for the current platform from GitHub releases.
 */

const https = require('https');
const http = require('http');
const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

// GitHub repository information
const GITHUB_REPO = 'Nehonix-Team/XyPCLI';
const GITHUB_API_URL = `https://api.github.com/repos/${GITHUB_REPO}/releases/latest`;

// Colors for output
const colors = {
  reset: '\x1b[0m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m'
};

function log(message, color = 'reset') {
  console.log(`${colors[color]}${message}${colors.reset}`);
}

function error(message) {
  log(`❌ Error: ${message}`, 'red');
}

function success(message) {
  log(`✅ ${message}`, 'green');
}

function info(message) {
  log(`ℹ️  ${message}`, 'blue');
}

// Detect platform and architecture
function detectPlatform() {
  const platform = os.platform();
  const arch = os.arch();

  let osName;
  let binaryName;

  switch (platform) {
    case 'darwin':
      osName = 'darwin';
      binaryName = arch === 'arm64' ? 'xypcli-darwin-arm64' : 'xypcli-darwin-amd64';
      break;
    case 'linux':
      osName = 'linux';
      binaryName = arch === 'arm64' ? 'xypcli-linux-arm64' : 'xypcli-linux-amd64';
      break;
    case 'win32':
      osName = 'windows';
      binaryName = arch === 'arm' ? 'xypcli-windows-arm.exe' : 'xypcli-windows-amd64.exe';
      break;
    default:
      throw new Error(`Unsupported platform: ${platform}`);
  }

  return { os: osName, arch, binaryName };
}

// Download file from URL
function downloadFile(url, destPath) {
  return new Promise((resolve, reject) => {
    const protocol = url.startsWith('https:') ? https : http;

    const request = protocol.get(url, (response) => {
      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download: HTTP ${response.statusCode}`));
        return;
      }

      const file = fs.createWriteStream(destPath);
      response.pipe(file);

      file.on('finish', () => {
        file.close();
        resolve();
      });

      file.on('error', (err) => {
        fs.unlink(destPath, () => {}); // Delete the file on error
        reject(err);
      });
    });

    request.on('error', (err) => {
      fs.unlink(destPath, () => {}); // Delete the file on error
      reject(err);
    });

    // Set a timeout
    request.setTimeout(30000, () => {
      request.destroy();
      reject(new Error('Download timeout'));
    });
  });
}

// Make binary executable (Unix-like systems)
function makeExecutable(filePath) {
  if (os.platform() !== 'win32') {
    try {
      fs.chmodSync(filePath, '755');
    } catch (err) {
      // Ignore chmod errors on some systems
    }
  }
}

// Get the binary installation path
function getBinaryPath() {
  // Use the same directory as this script
  return path.join(__dirname, 'xypcli' + (os.platform() === 'win32' ? '.exe' : ''));
}

// Check if binary already exists and is working
function isBinaryInstalled() {
  const binaryPath = getBinaryPath();
  if (!fs.existsSync(binaryPath)) {
    return false;
  }

  try {
    // Try to execute the binary with --version
    execSync(`"${binaryPath}" --version`, { timeout: 5000 });
    return true;
  } catch (err) {
    return false;
  }
}

// Main installation function
async function install() {
  try {
    info('XyPriss CLI Installer');
    info('Detecting platform...');

    const { os: osName, arch, binaryName } = detectPlatform();
    info(`Platform detected: ${osName}/${arch}`);
    info(`Binary to download: ${binaryName}`);

    // Check if binary is already installed
    if (isBinaryInstalled()) {
      success('XyPCLI is already installed and working!');
      return;
    }

    const binaryPath = getBinaryPath();

    // Get latest release information from GitHub
    info('Fetching latest release information...');

    const releaseInfo = await new Promise((resolve, reject) => {
      https.get(GITHUB_API_URL, {
        headers: {
          'User-Agent': 'XyPCLI-Installer'
        }
      }, (response) => {
        if (response.statusCode !== 200) {
          reject(new Error(`Failed to fetch release info: HTTP ${response.statusCode}`));
          return;
        }

        let data = '';
        response.on('data', (chunk) => data += chunk);
        response.on('end', () => {
          try {
            resolve(JSON.parse(data));
          } catch (err) {
            reject(new Error('Failed to parse release info'));
          }
        });
      }).on('error', reject);
    });

    // Find the correct asset
    const asset = releaseInfo.assets.find(a => a.name === binaryName);
    if (!asset) {
      error(`Binary ${binaryName} not found in latest release`);
      error('Available assets:');
      releaseInfo.assets.forEach(a => console.log(`  - ${a.name}`));
      process.exit(1);
    }

    const downloadUrl = asset.browser_download_url;
    info(`Downloading from: ${downloadUrl}`);

    // Download the binary
    await downloadFile(downloadUrl, binaryPath);

    // Make executable on Unix-like systems
    makeExecutable(binaryPath);

    success(`XyPCLI installed successfully!`);
    info(`Binary location: ${binaryPath}`);

    // Test the installation
    try {
      const version = execSync(`"${binaryPath}" --version`, { encoding: 'utf8' }).trim();
      success(`Installation verified: ${version}`);
    } catch (err) {
      error('Installation verification failed');
      throw err;
    }

  } catch (err) {
    error(err.message);
    process.exit(1);
  }
}

// If this script is run directly (not as a module)
if (require.main === module) {
  // Check if we should run the CLI or install
  const args = process.argv.slice(2);

  if (args.length === 0) {
    // No arguments, run installation
    install();
  } else {
    // Arguments provided, try to run the CLI
    const binaryPath = getBinaryPath();

    if (!isBinaryInstalled()) {
      log('XyPCLI not found. Installing...');
      install().then(() => {
        // After installation, execute the CLI with the provided arguments
        const { spawn } = require('child_process');
        const child = spawn(binaryPath, args, {
          stdio: 'inherit',
          shell: true
        });

        child.on('exit', (code) => {
          process.exit(code);
        });
      });
    } else {
      // CLI is already installed, execute it directly
      const { spawn } = require('child_process');
      const child = spawn(binaryPath, args, {
        stdio: 'inherit',
        shell: true
      });

      child.on('exit', (code) => {
        process.exit(code);
      });
    }
  }
}

module.exports = { install, detectPlatform };