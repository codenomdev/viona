const path = require('path');
const fs = require('fs');
const yaml = require('js-yaml');

const pluginPath = path.join(__dirname, '../src/plugins');
const pluginFolders = fs.readdirSync(pluginPath);

function pascalize(str) {
  return str.split(/[_-]/).map((part) => part.charAt(0).toUpperCase() + part.slice(1)).join('');
}

function resetIndexTs() {
  const indexTsPath = path.join(pluginPath, 'index.ts');
  fs.writeFileSync(indexTsPath, '');
}

function addPluginToIndexTs(packageName, pluginFolder) {
  const indexTsPath = path.join(pluginPath, 'index.ts');
  const indexTsContent = fs.readFileSync(indexTsPath, 'utf-8');
  const lines = indexTsContent.split('\n');
  const ComponentName = pascalize(packageName);

  const importLine = `const load${ComponentName} = () => import('${packageName}').then(module => module.default);`;
  const info = yaml.load(fs.readFileSync(path.join(pluginFolder, 'info.yaml'), 'utf8'));
  const exportLine = `export const ${info.slug_name} = load${ComponentName}`;

  if (!lines.includes(exportLine)) {
    lines.push(importLine);
    lines.push(exportLine);
  }

  fs.writeFileSync(indexTsPath, lines.join('\n'));
}

const pluginLength = pluginFolders.filter((folder) => {
  const pluginFolder = path.join(pluginPath, folder);
  const stat = fs.statSync(pluginFolder);
  return stat.isDirectory() && folder !== 'builtin';
}).length;

if (pluginLength > 0) {
  resetIndexTs();
}

pluginFolders.forEach((folder) => {
  const pluginFolder = path.join(pluginPath, folder);
  const stat = fs.statSync(pluginFolder);

  if (stat.isDirectory() && folder !== 'builtin') {
    if (!fs.existsSync(path.join(pluginFolder, 'index.ts'))) {
      return;
    }
    const packageJson = require(path.join(pluginFolder, 'package.json'));
    const packageName = packageJson.name;

    addPluginToIndexTs(packageName, pluginFolder);
  }
});