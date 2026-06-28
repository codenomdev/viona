const path = require('path');
const fs = require('fs');

const pluginPath = path.join(__dirname, '../src/plugins');
const pluginFolders = fs.readdirSync(pluginPath);

function resetPackageJson() {
    const packageJsonPath = path.join(__dirname, '..', 'package.json');
    const packageJsonContent = require(packageJsonPath);
    const dependencies = packageJsonContent.dependencies;
    for (const key in dependencies) {
        if (dependencies[key].startsWith('workspace')) {
            delete dependencies[key];
        }
    }
    fs.writeFileSync(
        packageJsonPath,
        JSON.stringify(packageJsonContent, null, 2),
    );
}

function addPluginToPackageJson(packageName) {
    const packageJsonPath = path.join(__dirname, '..', 'package.json');
    const packageJsonContent = require(packageJsonPath);
    packageJsonContent.dependencies[packageName] = 'workspace:*';

    fs.writeFileSync(
        packageJsonPath,
        JSON.stringify(packageJsonContent, null, 2),
    );
}


resetPackageJson();

pluginFolders.forEach((folder) => {
    const pluginFolder = path.join(pluginPath, folder);
    const stat = fs.statSync(pluginFolder);

    if (stat.isDirectory() && folder !== 'builtin') {
        if (!fs.existsSync(path.join(pluginFolder, 'index.ts'))) {
            return;
        }
        const packageJson = require(path.join(pluginFolder, 'package.json'));
        const packageName = packageJson.name;

        addPluginToPackageJson(packageName);
    }
});