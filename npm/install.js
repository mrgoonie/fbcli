const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const https = require("https");

const REPO = "mrgoonie/fbcli";
const BIN_DIR = path.join(__dirname, "bin");
const VERSION = require("./package.json").version;

const PLATFORM_MAP = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const ARCH_MAP = {
  x64: "amd64",
  arm64: "arm64",
};

function getBinaryName() {
  const platform = PLATFORM_MAP[process.platform];
  const arch = ARCH_MAP[process.arch];

  if (!platform || !arch) {
    console.error(
      `Unsupported platform: ${process.platform} ${process.arch}`
    );
    process.exit(1);
  }

  const ext = platform === "windows" ? ".zip" : ".tar.gz";
  return {
    archive: `fbcli_${platform}_${arch}${ext}`,
    binary: platform === "windows" ? "fbcli.exe" : "fbcli",
    platform,
  };
}

function downloadFile(url) {
  return new Promise((resolve, reject) => {
    const follow = (url) => {
      https
        .get(url, (res) => {
          if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
            follow(res.headers.location);
            return;
          }
          if (res.statusCode !== 200) {
            reject(new Error(`Download failed: HTTP ${res.statusCode}`));
            return;
          }
          const chunks = [];
          res.on("data", (chunk) => chunks.push(chunk));
          res.on("end", () => resolve(Buffer.concat(chunks)));
          res.on("error", reject);
        })
        .on("error", reject);
    };
    follow(url);
  });
}

async function install() {
  const { archive, binary, platform } = getBinaryName();
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${archive}`;

  console.log(`Downloading fbcli v${VERSION} for ${process.platform}/${process.arch}...`);

  try {
    const data = await downloadFile(url);
    const tmpFile = path.join(__dirname, archive);
    fs.writeFileSync(tmpFile, data);

    fs.mkdirSync(BIN_DIR, { recursive: true });

    if (platform === "windows") {
      // Use PowerShell to extract zip on Windows
      execSync(
        `powershell -command "Expand-Archive -Path '${tmpFile}' -DestinationPath '${BIN_DIR}' -Force"`,
        { stdio: "pipe" }
      );
    } else {
      execSync(`tar -xzf "${tmpFile}" -C "${BIN_DIR}"`, { stdio: "pipe" });
    }

    const binPath = path.join(BIN_DIR, binary);
    if (platform !== "windows") {
      fs.chmodSync(binPath, 0o755);
    }

    fs.unlinkSync(tmpFile);
    console.log(`fbcli v${VERSION} installed successfully.`);
  } catch (err) {
    console.error(`Failed to install fbcli: ${err.message}`);
    console.error(`You can download manually from: https://github.com/${REPO}/releases`);
    process.exit(1);
  }
}

install();
