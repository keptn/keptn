import { SequencesPage } from '../support/pageobjects/SequencesPage';

describe('Test sequence screen approval', () => {
  const sequencesPage = new SequencesPage();

  it('should show an loading indicator while the deployed image is fetched', () => {
    // sequencesPage.intercept().visit('sockshop');
  });

  it('should show an approval with the latest deployed image', () => {});

  it('should show an approval without the latest deployed image', () => {});
});
