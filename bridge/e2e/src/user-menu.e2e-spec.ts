import { AppPage } from './app.po';
import { browser, by, element, logging, protractor } from 'protractor';

import { takeScreenshot } from './utils';

// run with: ng e2e --dev-server-target= --base-url=http://localhost:3000/
describe('User Menu', () => {
  let page: AppPage;

  beforeEach(() => {
    page = new AppPage();
  });

  it('should open the user menu validate api token', async () => {
    await page.navigateTo();

    await element(
      by.xpath('//*[@uitestid="keptn-nav-userMenu"]')
    ).click();

    await element(
      by.xpath('//*[@uitestid="keptn-nav-copyKeptnApiToken"]/ktb-copy-to-clipboard/div/div[3]/button')
    ).click();

    takeScreenshot('user-menu-open-with-api-token-revealed.png');
  });

  afterEach(async () => {
    // Assert that there are no errors emitted from the browser
    const logs = await browser.manage().logs().get(logging.Type.BROWSER);
    expect(logs).not.toContain(jasmine.objectContaining({
      level: logging.Level.SEVERE,
    } as logging.Entry));
  });
});
