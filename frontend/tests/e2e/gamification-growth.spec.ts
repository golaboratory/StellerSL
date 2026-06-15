import { test, expect } from '@playwright/test';

test.describe('Gamification and Growth System', () => {
  let email: string;
  const password = 'password123';

  test.beforeEach(async ({ page }) => {
    email = `growth-${Date.now()}-${Math.random().toString(36).slice(2, 8)}@example.com`;
    // 1. Sign up a new user
    await page.goto('/signup');
    await page.fill('#name', 'Growth Tester');
    await page.fill('#email', email);
    await page.fill('#password input', password);
    await page.click('button:has-text("Sign Up")');
    // Wait for the signup redirect; navigating away earlier aborts the register request
    await expect(page).toHaveURL(/\/login/);

    // 2. Login
    await page.fill('#email', email);
    await page.fill('#password input', password);
    await page.click('button:has-text("Login")');
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should increase EXP when a task is completed', async ({ page }) => {
    const title = `Reward Task ${Date.now()}`;
    // 1. Check initial state in Dashboard
    await expect(page.locator('text=Level 1')).toBeVisible();
    await expect(page.locator('text=EXP: 0')).toBeVisible();

    // 2. Create a task
    await page.goto('/tasks');
    await page.click('button:has-text("Bulk Add")');
    await page.fill('textarea[placeholder*="Task 1"]', title);
    await page.click('button:has-text("Create Tasks")');

    // 3. Mark task as DONE
    const taskCard = page.locator('div.p-card', { hasText: title });
    // First click: todo -> doing
    await taskCard.locator('button.p-button-success').click(); 
    await expect(taskCard.locator('.p-tag:has-text("DOING")')).toBeVisible();
    // Second click: doing -> done
    await taskCard.locator('button.p-button-success').click();
    await expect(taskCard.locator('.p-tag:has-text("DONE")')).toBeVisible();

    // 4. Return to Dashboard and check EXP
    await page.goto('/dashboard');
    // Each task completion adds 10 EXP
    await expect(page.locator('text=EXP: 10')).toBeVisible();
  });

  test('should show achievement toast when first task is completed', async ({ page }) => {
    const title = `Badge Task ${Date.now()}`;
    // 1. Go to tasks and create one
    await page.goto('/tasks');
    await page.click('button:has-text("Bulk Add")');
    await page.fill('textarea[placeholder*="Task 1"]', title);
    await page.click('button:has-text("Create Tasks")');

    // 2. Complete the task, waiting for each transition to commit
    const taskCard = page.locator('div.p-card', { hasText: title });
    await taskCard.locator('button.p-button-success').click(); // to DOING
    await expect(taskCard.locator('.p-tag:has-text("DOING")')).toBeVisible();
    await taskCard.locator('button.p-button-success').click(); // to DONE
    await expect(taskCard.locator('.p-tag:has-text("DONE")')).toBeVisible();
    
    // 3. Go to Dashboard to trigger achievement check
    await page.goto('/dashboard');
    
    // 4. Verify Achievement Toast (PrimeVue Toast)
    await expect(page.locator('.p-toast-summary:has-text("Achievement Unlocked!")')).toBeVisible();
    await expect(page.locator('.p-toast-detail')).toContainText('earned');
    
    // 5. Verify badge appears in Achievements section
    await expect(page.locator('text=Achievements')).toBeVisible();
    // Badge icons are SVG assets (img.badge-icon) with a PrimeIcons fallback
    await expect(page.locator('div.p-card', { hasText: 'Achievements' }).locator('img.badge-icon, i.pi').first()).toBeVisible();
  });
});
