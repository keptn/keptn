import { AppPage } from './app.po';
import { browser, by, element, logging, protractor } from 'protractor';
import { takeScreenshot } from './utils';

// run with ng e2e --dev-server-target= --base-url=http://localhost:3000/
describe('Loading Bridge Example', () => {
  let page: AppPage;

  beforeEach(() => {
    page = new AppPage();
  });

  it('should match the title', () => {
    page.navigateTo();

    expect(browser.getTitle()).toEqual('keptn');

    takeScreenshot('entry-page.png');
  });

  afterEach(async () => {
    // Assert that there are no errors emitted from the browser
    const logs = await browser.manage().logs().get(logging.Type.BROWSER);
    expect(logs).not.toContain(jasmine.objectContaining({
      level: logging.Level.SEVERE,
    } as logging.Entry));
  });
});
