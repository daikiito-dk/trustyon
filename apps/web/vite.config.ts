import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  base: process.env.GITHUB_ACTIONS ? '/trustyon/' : '/',
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
