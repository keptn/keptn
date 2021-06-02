import { browser } from 'protractor';


// takes a screenshot and saves it into e2e-screenshots
export function takeScreenshot(filename) {
  browser.takeScreenshot().then(
    (image) => {
      require('fs').writeFile('e2e/screenshots/' + filename, image, 'base64', function(err) {
        if (err) { console.log(err); }
      });
    }
  );
}
