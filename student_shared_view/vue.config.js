const { defineConfig } = require('@vue/cli-service')
const webpack = require('webpack')

module.exports = defineConfig({
  transpileDependencies: true,
  // PWA配置
  pwa: {
    // 在开发环境中禁用Service Worker
    workboxPluginMode: process.env.NODE_ENV === 'production' ? 'GenerateSW' : 'InjectManifest',
    skipWaiting: true,
    clientsClaim: true,
    // 开发环境配置
    ...(process.env.NODE_ENV === 'development' && {
      workboxOptions: {
        skipWaiting: false,
        clientsClaim: false
      }
    })
  },
  configureWebpack: {
    plugins: [
      new webpack.DefinePlugin({
        __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: JSON.stringify(false)
      })
    ]
  }
})
