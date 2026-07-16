import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => ({
  plugins: [
    vue(),
    {
      name: 'runtime-config',
      configureServer(server) {
        const appMode = mode === 'admin' ? 'admin' : 'public'
        server.middlewares.use('/runtime-config.js', (req, res) => {
          res.setHeader('Content-Type', 'application/javascript; charset=utf-8')
          res.end(`window.__BLOG_APP_MODE__ = '${appMode}';`)
        })
      }
    }
  ],
  cacheDir: mode === 'admin' ? 'node_modules/.vite-admin' : 'node_modules/.vite',
  server: {
    port: mode === 'admin' ? 3001 : 3000,
    strictPort: true,
    proxy: {
      '/api': mode === 'admin' ? 'http://localhost:8081' : 'http://localhost:8080',
      '/uploads': 'http://localhost:8080'
    }
  }
}))
