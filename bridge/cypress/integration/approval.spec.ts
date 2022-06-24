import { SequencesPage } from '../support/pageobjects/SequencesPage';

describe('Test sequence screen approval', () => {
  const sequencesPage = new SequencesPage();
  const approvalId = '97b4a9a3-610a-4697-b37b-cec522a6b42f';

  beforeEach(() => {
    sequencesPage.intercept();
  });

  it('should show an loading indicator while the deployed image is fetched', () => {
    cy.intercept('/api/controlPlane/v1/project/sockshop/stage/production/service/carts', {
      body: {
        deployedImage: 'myImage:0.0.1',
      },
      delay: 10_000,
    }).as('approvalImage');
    sequencesPage
      .visit('sockshop')
      .selectSequence('99a20ef4-d822-4185-bbee-0d7a364c213b')
      .clickEvent(approvalId)
      .assertIsApprovalLoading(true);
  });

  it('should show an approval with the latest deployed image', () => {
    sequencesPage
      .visit('sockshop')
      .selectSequence('99a20ef4-d822-4185-bbee-0d7a364c213b')
      .clickEvent(approvalId)
      .assertIsApprovalLoading(false)
      .assertApprovalDeployedImage('myImage:0.0.1');
  });

  it('should not trigger an API call for the latest deployed image if it is not expanded', () => {
    sequencesPage.visit('sockshop').selectSequence('99a20ef4-d822-4185-bbee-0d7a364c213b');
    cy.get('@approvalImage.all').should('have.length', 0);
  });

  it('should show an approval without the latest deployed image', () => {
    cy.intercept('/api/controlPlane/v1/project/sockshop/stage/production/service/carts', {}).as('approvalImage');
    sequencesPage
      .visit('sockshop')
      .selectSequence('99a20ef4-d822-4185-bbee-0d7a364c213b')
      .clickEvent(approvalId)
      .assertIsApprovalLoading(false)
      .assertApprovalDeployedImage('N/A');
  });

  it('should not trigger latest deployed image call twice', () => {
    sequencesPage
      .visit('sockshop')
      .selectSequence('99a20ef4-d822-4185-bbee-0d7a364c213b')
      .clickEvent(approvalId) // open
      .clickEvent(approvalId) // close
      .clickEvent(approvalId); // open
    cy.get('@approvalImage.all').should('have.length', 1);
  });

  it('should show loading indicator while evaluation is loading', () => {
    sequencesPage
      .interceptEvaluationOfApproval(false, 10_000)
      .visit('sockshop')
      .selectSequence('99a20ef4-d822-4185-bbee-0d7a364c213b')
      .clickEvent(approvalId)
      .assertIsApprovalEvaluationLoading(true)
      .assertApprovalEvaluationBubbleExists(false);
  });

  it('should not show an evaluation for an approval', () => {
    sequencesPage
      .interceptEvaluationOfApproval(false)
      .visit('sockshop')
      .selectSequence('99a20ef4-d822-4185-bbee-0d7a364c213b')
      .clickEvent(approvalId)
      .assertIsApprovalEvaluationLoading(false)
      .assertApprovalEvaluationBubbleExists(false);
  });

  it('should show an evaluation for an approval', () => {
    sequencesPage
      .interceptEvaluationOfApproval(true)
      .visit('sockshop')
      .selectSequence('99a20ef4-d822-4185-bbee-0d7a364c213b')
      .clickEvent(approvalId)
      .assertIsApprovalEvaluationLoading(false)
      .assertApprovalEvaluationBubbleExists(true)
      .assertApprovalEvaluationBubble(50, 'error');
  });
});
