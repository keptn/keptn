module.exports = {
  devServer: {
    proxy: 'http://localhost:3000',
  },
  configureWebpack: {
    resolve: {
      alias: {
        '@': `${__dirname}/client`,
      },
    },
    entry: {
      app: './client/main.js',
    },
  },
};
