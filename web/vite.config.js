import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => ({
  plugins: [vue()],
  cacheDir: mode === 'admin' ? 'node_modules/.vite-admin' : 'node_modules/.vite',
  server: {
    port: mode === 'admin' ? 3001 : 3000,
    strictPort: true,
    proxy: {
      '/api': 'http://localhost:8080',
      '/admin-api': {
        target: 'http://localhost:8081',
        rewrite: path => path.replace(/^\/admin-api/, '/api')
      },
      '/uploads': 'http://localhost:8080'
    }
  }
}))
