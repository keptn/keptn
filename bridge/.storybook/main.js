module.exports = {
  stories: ['../stories/**/*.stories.mdx', '../stories/**/*.stories.@(js|jsx|ts|tsx)'],
  addons: ['@storybook/addon-links', '@storybook/addon-essentials', '@storybook/addon-interactions'],
  framework: '@storybook/angular',
  core: {
    builder: '@storybook/builder-webpack5',
  },
  staticDirs: [
    { from: '../node_modules/@dynatrace/barista-icons', to: '/assets/icons' },
    { from: '../node_modules/@dynatrace/barista-fonts/fonts', to: '/%5E./assets/fonts' },
    { from: '../client/assets', to: '/assets' },
  ],
};
