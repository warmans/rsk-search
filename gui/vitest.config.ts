import { defineConfig } from 'vitest/config';
import angular from '@analogjs/vite-plugin-angular';

export default defineConfig(({ mode }) => {
  return {
    resolve: {
      tsconfigPaths: true,
    },
    plugins: [
      angular({
        tsconfig: 'tsconfig.spec.json',
      }),
    ],
    test: {
      globals: true,
      setupFiles: ['src/test-setup.ts'],
      environment: 'jsdom',
      reporters: ['default'],
      passWithNoTests: true,
    },
    define: {
      'import.meta.vitest': mode !== 'production',
    },
  };
});
