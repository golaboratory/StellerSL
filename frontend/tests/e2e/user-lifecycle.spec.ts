import { test, expect } from '@playwright/test';

test.describe('User Lifecycle', () => {
  const email = `test-${Date.now()}@example.com`;
  const password = 'password123';

  test('should register a new user and login successfully', async ({ page }) => {
    // 1. Go to signup page
    await page.goto('/signup');
    await expect(page).toHaveTitle(/StellerSL/);

    // 2. Fill registration form
    await page.fill('#name', 'Test User');
    await page.fill('#email', email);
    await page.fill('#password input', password);
    await page.click('button:has-text("Sign Up")');

    // 3. Should be redirected to login
    await expect(page).toHaveURL(/\/login/);

    // 4. Fill login form
    await page.fill('#email', email);
    await page.fill('#password input', password);
    await page.click('button:has-text("Login")');

    // 5. Should be redirected to dashboard
    await expect(page).toHaveURL(/\/dashboard/);
    await expect(page.locator('text=Welcome')).toBeVisible();
  });
});
