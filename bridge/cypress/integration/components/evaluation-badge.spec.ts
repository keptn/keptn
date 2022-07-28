import EnvironmentPage from '../../support/pageobjects/EnvironmentPage';
import { EvaluationBadgeVariant } from '../../../client/app/_components/ktb-evaluation-badge/ktb-evaluation-badge.utils';
import ServicesPage from '../../support/pageobjects/ServicesPage';
import { SequencesPage } from '../../support/pageobjects/SequencesPage';

describe('Test evaluation badge border and fill', () => {
  const project = 'sockshop';

  describe('Test evaluation badge in environment screen', () => {
    const environmentPage = new EnvironmentPage();

    beforeEach(() => {
      environmentPage.intercept();
    });
    it('should have filled evaluation in stage overview', () => {
      environmentPage
        .visit(project)
        .assertEvaluationInOverview('dev', 'carts', 100, 'success', EvaluationBadgeVariant.FILL)
        .assertEvaluationInOverview('staging', 'carts', 50, 'error', EvaluationBadgeVariant.FILL)
        .assertEvaluationInOverview('production', 'carts', 50, 'warning', EvaluationBadgeVariant.FILL)
        .assertEvaluationInOverview('dev', 'carts-db', '-', undefined, EvaluationBadgeVariant.NONE);
    });

    it('should have bordered evaluation history in stage details', () => {
      environmentPage
        .interceptEvaluationHistory('sockshop', 'dev', 'carts', 6)
        .visit(project)
        .selectStage('dev')
        .waitForEvaluationHistory('carts', 'dev', 6)
        .waitForEvaluationHistory('carts-db', 'dev', 5)
        .assertEvaluationHistory('carts', [
          { score: 50, status: 'warning', variant: EvaluationBadgeVariant.BORDER },
          { score: 0, status: 'error', variant: EvaluationBadgeVariant.BORDER },
          { score: 0, status: 'error', variant: EvaluationBadgeVariant.BORDER },
          { score: 0, status: 'error', variant: EvaluationBadgeVariant.BORDER },
          { score: 100, status: 'success', variant: EvaluationBadgeVariant.BORDER },
          { score: 100, status: 'success', variant: EvaluationBadgeVariant.FILL },
        ]);
    });

    it('should have success evaluation badge in the stage details', () => {
      environmentPage
        .visit(project)
        .selectStage('dev')
        .waitForEvaluationHistory('carts', 'dev', 6)
        .waitForEvaluationHistory('carts-db', 'dev', 5)
        .assertEvaluationInDetails('carts', 100, 'success', EvaluationBadgeVariant.FILL)
        .assertEvaluationInDetails('carts-db', '-', undefined, EvaluationBadgeVariant.NONE);
    });

    it('should have error evaluation badge in stage details', () => {
      environmentPage
        .visit(project)
        .selectStage('staging')
        .waitForEvaluationHistory('carts', 'staging', 6)
        .waitForEvaluationHistory('carts-db', 'staging', 5)
        .assertEvaluationInDetails('carts', 50, 'error', EvaluationBadgeVariant.FILL)
        .assertEvaluationInDetails('carts-db', '-', undefined, EvaluationBadgeVariant.NONE);
    });

    it('should have warning evaluation badge in stage details', () => {
      environmentPage
        .visit(project)
        .selectStage('production')
        .waitForEvaluationHistory('carts', 'production', 6)
        .waitForEvaluationHistory('carts-db', 'production', 5)
        .assertEvaluationInDetails('carts', 50, 'warning', EvaluationBadgeVariant.FILL)
        .assertEvaluationInDetails('carts-db', '-', undefined, EvaluationBadgeVariant.NONE);
    });
  });

  describe('Test evaluation badge in service screen timeline', () => {
    const servicesPage = new ServicesPage();

    beforeEach(() => {
      servicesPage
        .interceptAll()
        .interceptForEvaluationBadge()
        .visitServicePage(project)
        .selectService('carts', 'v0.1.2');
    });

    describe('Bordered evaluation badge', () => {
      it('should have successful evaluation badge', () => {
        servicesPage
          .assertStageEvaluationBadge('evaluation-success', 'success', 33, EvaluationBadgeVariant.BORDER)
          .assertStageEvaluationBadge('evaluation-none', 'success', '-', EvaluationBadgeVariant.BORDER);
      });

      it('should have error evaluation badge', () => {
        servicesPage.assertStageEvaluationBadge('evaluation-error', 'error', 33, EvaluationBadgeVariant.BORDER);
      });

      it('should have warning evaluation badge', () => {
        servicesPage.assertStageEvaluationBadge('evaluation-warning', 'warning', 33, EvaluationBadgeVariant.BORDER);
      });
    });

    describe('Filled evaluation badge', () => {
      it('should have successful evaluation badge', () => {
        servicesPage
          .assertStageEvaluationBadge('evaluation-success-deployed', 'success', 33, EvaluationBadgeVariant.FILL)
          .assertStageEvaluationBadge('evaluation-none-deployed', 'success', '-', EvaluationBadgeVariant.FILL);
      });

      it('should have error evaluation badge', () => {
        servicesPage.assertStageEvaluationBadge('evaluation-error-deployed', 'error', 33, EvaluationBadgeVariant.FILL);
      });

      it('should have warning evaluation badge', () => {
        servicesPage.assertStageEvaluationBadge(
          'evaluation-warning-deployed',
          'warning',
          33,
          EvaluationBadgeVariant.FILL
        );
      });
    });
  });

  describe('Test evaluation badge in sequence screen', () => {
    const sequencePage = new SequencesPage();

    beforeEach(() => {
      sequencePage.intercept().visit(project);
    });

    it('should have error evaluation badge', () => {
      sequencePage.assertSequenceEvaluationBadge(
        'ce6f3686-90a2-497e-ac6a-e8ed95a845c5',
        'staging',
        'error',
        0,
        EvaluationBadgeVariant.FILL
      );
    });

    it('should have warning evaluation badge', () => {
      sequencePage.assertSequenceEvaluationBadge(
        '99a20ef4-d822-4185-bbee-0d7a364c213a',
        'dev',
        'warning',
        50,
        EvaluationBadgeVariant.FILL
      );
    });

    it('should have successful evaluation badge', () => {
      sequencePage.assertSequenceEvaluationBadge(
        '99a20ef4-d822-4185-bbee-0d7a364c213a',
        'staging',
        'success',
        0,
        EvaluationBadgeVariant.FILL
      );
    });

    it('should have empty evaluation badge', () => {
      sequencePage.assertSequenceEvaluationBadge(
        '99a20ef4-d822-4185-bbee-0d7a364c213a',
        'production',
        undefined,
        '-',
        EvaluationBadgeVariant.NONE
      );
    });
  });
});
