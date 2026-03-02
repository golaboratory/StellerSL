import { chromium } from '@playwright/test';

(async () => {
  const browser = await chromium.launch();
  // No custom context with Host header here, config handles X-Tenant-Host
  const page = await browser.newPage();
  
  page.on('console', msg => console.log('BROWSER CONSOLE:', msg.text()));
  
  page.on('response', async response => {
    if (response.url().includes('/api/auth/register')) {
      console.log('API RESPONSE status:', response.status());
      try {
        const text = await response.text();
        console.log('API RESPONSE body:', text);
      } catch (e) {
        console.log('API RESPONSE body: (could not read)');
      }
    }
  });

  try {
    console.log('Navigating to http://localhost:5173/signup...');
    await page.goto('http://localhost:5173/signup', { waitUntil: 'networkidle' });
    
    await page.fill('#name', 'Debug User');
    await page.fill('#email', `debug-${Date.now()}@example.com`);
    await page.fill('#password input', 'password123');
    
    console.log('Clicking Sign Up...');
    await page.click('button:has-text("Sign Up")');
    await page.waitForTimeout(5000);
    
  } catch (err) {
    console.error('Error during debug:', err);
  } finally {
    await browser.close();
  }
})();
