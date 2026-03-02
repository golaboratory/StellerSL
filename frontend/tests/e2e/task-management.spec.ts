import { test, expect } from '@playwright/test';

test.describe('Task Management (GTD)', () => {
  const email = `test-task-${Date.now()}@example.com`;
  const password = 'password123';

  test.beforeEach(async ({ page }) => {
    // Register and Login before each test
    await page.goto('/signup');
    await page.fill('#name', 'Task Tester');
    await page.fill('#email', email);
    await page.fill('#password input', password);
    await page.click('button:has-text("Sign Up")');
    
    await page.goto('/login');
    await page.fill('#email', email);
    await page.fill('#password input', password);
    await page.click('button:has-text("Login")');
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should perform bulk operations and GTD status updates', async ({ page }) => {
    // 1. Navigate to Tasks
    await page.goto('/tasks');
    await expect(page.locator('text=Tasks (GTD)')).toBeVisible();

    // 2. Bulk Add Tasks
    await page.click('button:has-text("Bulk Add")');
    await page.fill('textarea[placeholder*="Task 1"]', `Task A
Task B
Task C`);
    await page.click('button:has-text("Create Tasks")');

    // 3. Verify tasks are created
    await expect(page.locator('text=Task A')).toBeVisible();
    await expect(page.locator('text=Task B')).toBeVisible();
    await expect(page.locator('text=Task C')).toBeVisible();

    // 4. Update individual status (Task A: todo -> doing)
    const taskACard = page.locator('div.p-card', { hasText: 'Task A' });
    await taskACard.locator('button.p-button-success').click(); // First check marks as 'doing'
    await expect(taskACard.locator('.p-tag:has-text("DOING")')).toBeVisible();

    // 5. Bulk update status (Task B & C: todo -> done)
    // Find checkboxes for Task B and C
    await page.locator('div.p-card', { hasText: 'Task B' }).locator('.p-checkbox-box').click();
    await page.locator('div.p-card', { hasText: 'Task C' }).locator('.p-checkbox-box').click();
    
    await page.click('button:has-text("Set Done")');
    
    // Verify Task B and C are DONE
    await expect(page.locator('div.p-card', { hasText: 'Task B' }).locator('.p-tag:has-text("DONE")')).toBeVisible();
    await expect(page.locator('div.p-card', { hasText: 'Task C' }).locator('.p-tag:has-text("DONE")')).toBeVisible();

    // 6. Bulk Delete (Task B & C)
    // Handle the window.confirm dialog
    page.on('dialog', dialog => dialog.accept());
    
    await page.click('button.p-button-danger'); // Bulk delete button
    
    // Verify Task B and C are gone
    await expect(page.locator('text=Task B')).not.toBeVisible();
    await expect(page.locator('text=Task C')).not.toBeVisible();
    await expect(page.locator('text=Task A')).toBeVisible();
  });
});
