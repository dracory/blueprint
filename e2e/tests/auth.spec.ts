import { test, expect } from '@playwright/test';

// Credentials of the user seeded by e2e/browser (seed.go).
const EMAIL = 'e2e-browser@blueprint.local';
const PASSWORD = 'password123';

test.describe('authentication', () => {
  test('login page loads', async ({ page }) => {
    await page.goto('/auth/login');
    await expect(page.locator('input[name="email"]')).toBeVisible();
    await expect(page.locator('input[name="password"]')).toBeVisible();
  });

  test('seeded admin can log in and reach the dashboard', async ({ page }) => {
    await page.goto('/auth/login');
    await page.fill('input[name="email"]', EMAIL);
    await page.fill('input[name="password"]', PASSWORD);
    await page.locator('button.ButtonLogin, button[type="submit"]').first().click();

    // The app redirects authenticated admins to the admin area.
    await expect(page).toHaveURL(/\/admin/, { timeout: 10_000 });
  });
});
