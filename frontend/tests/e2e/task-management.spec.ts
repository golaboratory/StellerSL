import { test, expect } from '@playwright/test';

test.describe('Task Management (GTD)', () => {
  let email: string;
  const password = 'password123';

  test.beforeEach(async ({ page }) => {
    email = `test-task-${Date.now()}-${Math.random().toString(36).slice(2, 8)}@example.com`;
    // Register and Login before each test
    await page.goto('/signup');
    await page.fill('#name', 'Task Tester');
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

  test('should perform bulk operations and GTD status updates', async ({ page }) => {
    // 1. Navigate to Tasks
    await page.goto('/tasks');
    await expect(page.locator('text=Tasks (GTD)')).toBeVisible();

    // Unique titles: tasks are tenant-wide, so fixed names collide with
    // seed data ("Task CRUD") and with parallel/previous runs
    const stamp = Date.now();
    const taskA = `Task A ${stamp}`;
    const taskB = `Task B ${stamp}`;
    const taskC = `Task C ${stamp}`;

    // 2. Bulk Add Tasks
    await page.click('button:has-text("Bulk Add")');
    await page.fill('textarea[placeholder*="Task 1"]', `${taskA}
${taskB}
${taskC}`);
    await page.click('button:has-text("Create Tasks")');

    // 3. Verify tasks are created
    await expect(page.locator(`text=${taskA}`)).toBeVisible();
    await expect(page.locator(`text=${taskB}`)).toBeVisible();
    await expect(page.locator(`text=${taskC}`)).toBeVisible();

    // 4. Update individual status (Task A: todo -> doing)
    const taskACard = page.locator('div.p-card', { hasText: taskA });
    await taskACard.locator('button.p-button-success').click(); // First check marks as 'doing'
    await expect(taskACard.locator('.p-tag:has-text("DOING")')).toBeVisible();

    // 5. Bulk update status (Task B & C: todo -> done)
    // Find checkboxes for Task B and C
    await page.locator('div.p-card', { hasText: taskB }).locator('input.p-checkbox-input').click();
    await page.locator('div.p-card', { hasText: taskC }).locator('input.p-checkbox-input').click();
    
    await page.click('button:has-text("Set Done")');
    
    // Verify Task B and C are DONE
    await expect(page.locator('div.p-card', { hasText: taskB }).locator('.p-tag:has-text("DONE")')).toBeVisible();
    await expect(page.locator('div.p-card', { hasText: taskC }).locator('.p-tag:has-text("DONE")')).toBeVisible();

    // 6. Bulk Delete (Task B & C)
    // The selection resets after the list refresh, so select the tasks again
    await page.locator('div.p-card', { hasText: taskB }).locator('input.p-checkbox-input').click();
    await page.locator('div.p-card', { hasText: taskC }).locator('input.p-checkbox-input').click();
    // Handle the window.confirm dialog
    page.on('dialog', dialog => dialog.accept());
    
    await page.click('button.p-button-danger'); // Bulk delete button
    
    // Verify Task B and C are gone
    await expect(page.locator(`text=${taskB}`)).not.toBeVisible();
    await expect(page.locator(`text=${taskC}`)).not.toBeVisible();
    await expect(page.locator(`text=${taskA}`)).toBeVisible();
  });
});
