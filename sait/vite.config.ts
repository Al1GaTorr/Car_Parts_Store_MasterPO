import path from 'path';
import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, '.', '');
  const devPort = Number(env.VITE_DEV_PORT || 5174);
  const serverPort = Number.isFinite(devPort) ? devPort : 5174;
  const serverHost = env.VITE_DEV_HOST || '0.0.0.0';
  const apiProxyTarget = env.VITE_API_PROXY_TARGET || 'http://localhost:8090';

  return {
    server: {
      port: serverPort,
      host: serverHost,
      proxy: {
        '/api': {
          target: apiProxyTarget,
          changeOrigin: true,
        }
      }
    },
    plugins: [react()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, '.'),
      }
    }
  };
});
