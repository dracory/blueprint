import { test, expect } from '@playwright/test';

test.describe('home page', () => {
  test('loads and shows the welcome heading', async ({ page }) => {
    await page.goto('/');
    // The home controller renders "Welcome to <appName>"; sandbox sets
    // APP_NAME=Blueprint.
    await expect(page.locator('h1')).toContainText('Blueprint');
  });
});
