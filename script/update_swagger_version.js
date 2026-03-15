const fs = require("fs");

function readFileIfExists(filePath) {
  if (!fs.existsSync(filePath)) {
    return "";
  }

  return fs.readFileSync(filePath, "utf8");
}

function parseCodeVersion(source) {
  const match = source.match(/\/\/\s*@version\s+(\d+)(?:\.(\d+))?/);
  if (!match) {
    throw new Error("@version コメントが見つかりません");
  }

  return {
    major: Number.parseInt(match[1], 10),
    minor: Number.parseInt(match[2] ?? "0", 10),
  };
}

function parseSwaggerVersion(source) {
  if (!source) {
    return null;
  }

  const json = JSON.parse(source);
  const version = String(json?.info?.version ?? "");
  const match = version.match(/^(\d+)\.(\d+)$/);
  if (!match) {
    return null;
  }

  return {
    major: Number.parseInt(match[1], 10),
    minor: Number.parseInt(match[2], 10),
  };
}

function main() {
  const [, , mainGoPath, swaggerPath, previousSwaggerPath] = process.argv;
  if (!mainGoPath || !swaggerPath || !previousSwaggerPath) {
    throw new Error("usage: node update_swagger_version.js <main.go> <swagger.json> <previous_swagger.json>");
  }

  const codeVersion = parseCodeVersion(fs.readFileSync(mainGoPath, "utf8"));
  const previousVersion = parseSwaggerVersion(readFileIfExists(previousSwaggerPath));
  const swagger = JSON.parse(fs.readFileSync(swaggerPath, "utf8"));

  const nextMinor =
    previousVersion && previousVersion.major === codeVersion.major
      ? previousVersion.minor + 1
      : codeVersion.minor;

  swagger.info.version = `${codeVersion.major}.${nextMinor}`;
  fs.writeFileSync(swaggerPath, `${JSON.stringify(swagger, null, 4)}\n`, "utf8");
  console.log(`swagger version updated: ${swagger.info.version}`);
}

main();
