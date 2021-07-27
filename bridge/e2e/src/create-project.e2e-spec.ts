import {AppPage} from './app.po';
import {browser, logging} from 'protractor';
import {takeScreenshot} from './utils';

// run with: ng e2e --dev-server-target= --base-url=http://localhost:3000/
describe('Create Project', () => {
  let page: AppPage;

  beforeEach(() => {
    page = new AppPage();
    page.navigateToPath('create/project');
  });

  it('should create a new project', async () => {
    takeScreenshot('create-project-page.png');
  });

  afterEach(async () => {
    // Assert that there are no errors emitted from the browser
    const logs = await browser.manage().logs().get(logging.Type.BROWSER);
    expect(logs).not.toContain(jasmine.objectContaining({
      level: logging.Level.SEVERE,
    } as logging.Entry));
  });
});
