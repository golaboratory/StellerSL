import { test, expect } from '@playwright/test';

test.describe('Gamification and Growth System', () => {
  const email = `growth-${Date.now()}@example.com`;
  const password = 'password123';

  test.beforeEach(async ({ page }) => {
    // 1. Sign up a new user
    await page.goto('/signup');
    await page.fill('#name', 'Growth Tester');
    await page.fill('#email', email);
    await page.fill('#password input', password);
    await page.click('button:has-text("Sign Up")');
    
    // 2. Login
    await page.goto('/login');
    await page.fill('#email', email);
    await page.fill('#password input', password);
    await page.click('button:has-text("Login")');
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should increase EXP when a task is completed', async ({ page }) => {
    // 1. Check initial state in Dashboard
    await expect(page.locator('text=Level 1')).toBeVisible();
    await expect(page.locator('text=EXP: 0')).toBeVisible();

    // 2. Create a task
    await page.goto('/tasks');
    await page.click('button:has-text("Bulk Add")');
    await page.fill('textarea[placeholder*="Task 1"]', 'Reward Task');
    await page.click('button:has-text("Create Tasks")');

    // 3. Mark task as DONE
    const taskCard = page.locator('div.p-card', { hasText: 'Reward Task' });
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
    // 1. Go to tasks and create one
    await page.goto('/tasks');
    await page.click('button:has-text("Bulk Add")');
    await page.fill('textarea[placeholder*="Task 1"]', 'Badge Task');
    await page.click('button:has-text("Create Tasks")');

    // 2. Complete the task
    const taskCard = page.locator('div.p-card', { hasText: 'Badge Task' });
    await taskCard.locator('button.p-button-success').click(); // to DOING
    await taskCard.locator('button.p-button-success').click(); // to DONE
    
    // 3. Go to Dashboard to trigger achievement check
    await page.goto('/dashboard');
    
    // 4. Verify Achievement Toast (PrimeVue Toast)
    await expect(page.locator('.p-toast-summary:has-text("Achievement Unlocked!")')).toBeVisible();
    await expect(page.locator('.p-toast-detail')).toContainText('earned');
    
    // 5. Verify badge appears in Achievements section
    await expect(page.locator('text=Achievements')).toBeVisible();
    // The first badge might be 'Task Starter' or similar
    await expect(page.locator('div.p-card', { hasText: 'Achievements' }).locator('i.pi')).toBeVisible();
  });
});
