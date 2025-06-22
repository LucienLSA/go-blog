module.exports = {
    assetsDir: "static",
    devServer: {
      
      // 更换ip
      host: '127.0.0.1',
  
      // 更换端口号
      port: 8001,
        proxy: {
            '/api/v1': {
              target: 'http://127.0.0.1:8081',
              changeOrigin: true,
            },
            '/api/v2': {
              target: 'http://localhost:8081',
              changeOrigin: true,
              pathRewrite: { '^/api/v2': '/api/v2' }
            }
        }
    }
  }
