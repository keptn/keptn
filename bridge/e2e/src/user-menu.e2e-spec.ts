import { AppPage } from './app.po';
import { browser, by, element, logging, protractor } from 'protractor';

import { takeScreenshot } from './utils';

// run with: ng e2e --dev-server-target= --base-url=http://localhost:3000/
describe('User Menu', () => {
  let page: AppPage;

  beforeEach(() => {
    page = new AppPage();
  });

  it('should open the user menu and copy api token', async () => {
    await page.navigateTo();

    await element(
      by.xpath('//*[@uitestid="keptn-nav-userMenu"]')
    ).click();
    takeScreenshot('user-menu-open.png');

    await element(
      by.xpath('//*[@uitestid="keptn-nav-copyKeptnApiToken"]/ktb-copy-to-clipboard/div/div[2]/dt-copy-to-clipboard/button')
    ).click();

    await browser.executeScript(() => {
      let el = document.createElement('input');
      el.setAttribute('id', 'apiTokenTempInput');
      document.getElementsByTagName('body')[0].appendChild(el);
    });

    const apiTokenTempInput = element(by.id('apiTokenTempInput'));
    await apiTokenTempInput.sendKeys(protractor.Key.chord(protractor.Key.CONTROL, 'v'));

    const val = await apiTokenTempInput.getAttribute('value');

    const apiToken = process.env.KEPTN_API_TOKEN;

    await expect(val).toEqual(apiToken);
  });

  afterEach(async () => {
    // Assert that there are no errors emitted from the browser
    const logs = await browser.manage().logs().get(logging.Type.BROWSER);
    expect(logs).not.toContain(jasmine.objectContaining({
      level: logging.Level.SEVERE,
    } as logging.Entry));
  });
});
