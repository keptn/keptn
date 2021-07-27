import {browser} from 'protractor';

export class AppPage {
  navigateTo() {
    return browser.get(browser.baseUrl) as Promise<any>;
  }

  navigateToPath(path: string) {
    return browser.get(`${browser.baseUrl}${path}`) as Promise<any>;
  }
}
