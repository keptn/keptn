import { browser } from 'protractor';
import ErrnoException = NodeJS.ErrnoException;


// takes a screenshot and saves it into e2e-screenshots
export function takeScreenshot(filename: string) {
  browser.takeScreenshot().then(
    (image) => {
      require('fs').writeFile('e2e/screenshots/' + filename, image, 'base64', (err: ErrnoException | null) => {
        if (err) {
          console.log(err);
        }
      });
    }
  );
}
