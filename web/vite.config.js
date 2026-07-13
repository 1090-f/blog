import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => ({
  plugins: [vue()],
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
