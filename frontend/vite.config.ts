import { sveltekit } from "@sveltejs/kit/vite";
import type { UserConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";

const config: UserConfig = {
	plugins: [
		tailwindcss(),
		sveltekit(),
	],
	server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  },
};

export default config;