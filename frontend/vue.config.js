const { defineConfig } = require('@vue/cli-service')
const fs = require('fs')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    host: 'dev.local',
    port: 8080,
    // server: {
    //   type: 'https',
    //   options: {
    //     key: fs.readFileSync('./certs/dev.local+4-key.pem'),
    //     cert: fs.readFileSync('./certs/dev.local+4.pem'),
    //   },
    // },
  }
})