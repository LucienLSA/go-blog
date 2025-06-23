module.exports = {
    assetsDir: "static",
    devServer: {
      
      // 更换ip
      host: '0.0.0.0',
  
      // 更换端口号
      port: 8001,
        proxy: {
            '/api/v2': {
              target: 'http://127.0.0.1:8081',
              changeOrigin: true,
            }
        }
    }
  }
