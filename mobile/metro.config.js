const { getDefaultConfig } = require('expo/metro-config');
const path = require('path');

const config = getDefaultConfig(__dirname);

// Add shared folder to watchFolders
config.watchFolders = [
  path.resolve(__dirname, '../shared'),
];

// Add shared to node_modules resolution
config.resolver.nodeModulesPaths = [
  path.resolve(__dirname, 'node_modules'),
  path.resolve(__dirname, '../shared'),
];

module.exports = config;
