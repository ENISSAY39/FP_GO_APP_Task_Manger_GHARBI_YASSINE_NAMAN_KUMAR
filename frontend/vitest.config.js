import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

// Test config is kept out of vite.config.js: the app build is a multi-page
// setup with its own inputs and a dev proxy, none of which the runner needs.
// Only the React plugin is shared, so components compile the same way here as
// they do in the app.
export default defineConfig({
  plugins: [react()],

  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.js'],
    include: ['src/**/*.test.{js,jsx}'],
  },
})
