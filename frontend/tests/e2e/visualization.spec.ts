import { test, expect } from '@playwright/test';

test.describe('Project Visualization', () => {
  let email: string;
  const password = 'password123';

  test.beforeEach(async ({ page }) => {
    email = `viz-${Date.now()}-${Math.random().toString(36).slice(2, 8)}@example.com`;
    // 1. Signup and Login
    await page.goto('/signup');
    await page.fill('#name', 'Visualization Tester');
    await page.fill('#email', email);
    await page.fill('#password input', password);
    await page.click('button:has-text("Sign Up")');
    // Wait for the signup redirect; navigating away earlier aborts the register request
    await expect(page).toHaveURL(/\/login/);
    await page.fill('#email', email);
    await page.fill('#password input', password);
    await page.click('button:has-text("Login")');
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should display tasks in Gantt and Calendar views', async ({ page }) => {
    const title = `Visualization Task ${Date.now()}`;
    // 1. Create a task in Tasks page
    await page.goto('/tasks');
    await page.click('button:has-text("Bulk Add")');
    await page.fill('textarea[placeholder*="Task 1"]', title);
    await page.click('button:has-text("Create Tasks")');
    // Wait until the task is persisted and listed; navigating away earlier aborts the request
    await expect(page.locator(`text=${title}`)).toBeVisible();

    // 2. Open Gantt view
    await page.goto('/gantt');
    await expect(page.locator('text=Gantt Chart')).toBeVisible();
    
    // 3. Verify task appears in Gantt (Frappe Gantt uses SVG elements)
    await expect(page.locator('#gantt-target svg').first()).toBeVisible();
    await expect(page.locator(`#gantt-target text:has-text("${title}")`).first()).toBeVisible();

    // 4. Open Calendar view
    await page.goto('/calendar');
    await expect(page.locator('text=Calendar')).toBeVisible();
    
    // 5. Verify task appears in Calendar (v-calendar components)
    // Note: Since the mock start/end is Today, it should appear on current day
    const today = new Date().getDate().toString();
    // V-calendar day cell with today's date
    const todayCell = page.locator(`.vc-day.is-today`);
    await expect(todayCell.locator('.vc-dot').first()).toBeVisible();
    
    // 6. Verify task list below calendar
    await expect(page.locator(`text=${title}`).first()).toBeVisible();
  });
});
