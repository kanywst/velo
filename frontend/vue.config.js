let cssConfig = {};

if (process.env.NODE_ENV == "production") {
  cssConfig = {
    extract: {
      filename: "app.css",
      chunkFilename: "app.css"
    }
  };
}

module.exports = {
  chainWebpack: config => {
    // Vue CLI 5 uses Webpack 5 which uses Asset Modules instead of url-loader
    // We want to inline everything to keep the binary self-contained
    const limit = 10 * 1024 * 1024 * 1024; // 10GB, effectively unlimited

    config.module
      .rule('images')
      .set('parser', {
        dataUrlCondition: {
          maxSize: limit
        }
      });

    config.module
      .rule('fonts')
      .set('parser', {
        dataUrlCondition: {
          maxSize: limit
        }
      });
      
    config.module
      .rule('svg')
      .set('parser', {
        dataUrlCondition: {
          maxSize: limit
        }
      });
  },
  css: cssConfig,
  configureWebpack: {
    output: {
      filename: "app.js"
    },
    optimization: {
      splitChunks: false
    }
  },
  devServer: {
    allowedHosts: 'auto'
  }
};