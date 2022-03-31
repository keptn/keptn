import { SequencesPage } from '../support/pageobjects/SequencesPage';

describe('Sequences', () => {
  const sequencePage = new SequencesPage();

  beforeEach(() => {
    sequencePage.interceptRemediationSequences();
  });

  it('should show remediation in regular state while running', () => {
    sequencePage
      .visit('sockshop')
      .selectSequence('cfaadbb1-3c47-46e5-a230-2e312cf1828a')
      .assertTaskFailed('remediation', false)
      .assertTaskSuccessful('remediation', false);
  });

  it('should show remediation green when successful', () => {
    sequencePage
      .visit('sockshop')
      .selectSequence('cfaadbb1-3c47-46e5-a230-d0f055f4f518')
      .assertTaskFailed('remediation', false)
      .assertTaskSuccessful('remediation', true);
  });

  it('should show remediation red when failed', () => {
    sequencePage
      .visit('sockshop')
      .selectSequence('29355a07-7b65-47fa-896e-06f656283c5d')
      .assertTaskFailed('remediation', true)
      .assertTaskSuccessful('remediation', false);
  });
});
