import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 5173,
    proxy: {
		// Toutes les requêtes qui commencent par /api seront redirigées vers le backend Go
		'/api': {
        target: 'http://backend:8081',
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, '') // Décommente si ton backend n'a pas le préfixe /api dans ses routes
      }
    }
  }
})

// export default defineConfig({
//   plugins: [react(), tailwindcss()],
//   server: {
//     port: 5173,
//     https: {
//       key: '/etc/nginx/ssl/key.pem',
//       cert: '/etc/nginx/ssl/cert.pem',
//     },
//     proxy: {
//       '/api': {
//         target: 'https://backend:8081',
//         changeOrigin: true,
//         secure: false,
//         // rewrite: (path) => path.replace(/^\/api/, '') // Décommente si ton backend n'a pas le préfixe /api dans ses routes
//       }
//     }
//   }
// })