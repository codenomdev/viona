const fs = require('fs');
const path = require('path');
const yaml = require('js-yaml');

const configFilePath = path.resolve(__dirname, '../../config/config.yaml');
const envFilePath = path.resolve(__dirname, '../.env.production');

// Read config.yaml file
const config = yaml.load(fs.readFileSync(configFilePath, 'utf8'));

// Generate .env file content
let envContent = 'TSC_COMPILE_ON_ERROR=true\nESLINT_NO_DEV_ERRORS=true\n';
for (const key in config.UI) {
  const value = config.UI[key];
  envContent += `${key !== 'PUBLIC_URL' ? 'REACT_APP_' : ''}${key.toUpperCase()}=${value}\n`;
}

// Write .env file
fs.writeFileSync(envFilePath, envContent);
