import { AppPage } from './app.po';
import { browser, by, element, logging, protractor } from 'protractor';

import { takeScreenshot } from './utils';

// run with ng e2e --dev-server-target= --base-url=http://localhost:3000/
describe('Bridge Navigation', () => {
  let page: AppPage;

  beforeEach(() => {
    page = new AppPage();
  });

  it('should click on the project dropdown', async () => {
    await page.navigateTo();

    await element(
      by.xpath('//*[@id="projectSelect"]//div')
    ).click();
    takeScreenshot('project-menu-open.png');

    // find first option
    await element(
      by.id('dt-option-0')
    ).click();

    takeScreenshot('project-menu-clicked.png');
  });

  afterEach(async () => {
    // Assert that there are no errors emitted from the browser
    const logs = await browser.manage().logs().get(logging.Type.BROWSER);
    expect(logs).not.toContain(jasmine.objectContaining({
      level: logging.Level.SEVERE,
    } as logging.Entry));
  });
});
