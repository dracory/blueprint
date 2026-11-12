import { defineConfig } from '@playwright/test';

const BASE_URL = 'http://127.0.0.1:34757';

export default defineConfig({
  testDir: './tests',
  timeout: 30_000,
  retries: process.env.CI ? 1 : 0,
  use: {
    baseURL: BASE_URL,
    channel: 'chromium', // bundled headless Chromium
    screenshot: 'only-on-failure',
    trace: 'retain-on-failure',
  },
  webServer: {
    // `go run ./e2e/browser` run from the repository root.
    command: 'go run ./e2e/browser',
    cwd: '..',
    url: BASE_URL,
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
  },
});
