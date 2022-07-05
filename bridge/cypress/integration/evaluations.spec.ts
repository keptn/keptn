import { HeatmapComponentPage } from '../support/pageobjects/HeatmapComponentPage';
import { EvaluationBoardPage } from '../support/pageobjects/EvaluationBoardPage';
import { ResultTypes } from '../../shared/models/result-types';

const heatmap = new HeatmapComponentPage();
const evaluationBoard = new EvaluationBoardPage();

describe('evaluations', () => {
  beforeEach(() => {
    heatmap.intercept();
  });

  it('should load the heatmap with sli breakdown in service screen', () => {
    heatmap.visitPageWithHeatmapComponent();
    heatmap.assertComponentExists();
  });

  it('should truncate score to 2 decimals', () => {
    heatmap.visitPageWithHeatmapComponent();
    evaluationBoard.assertScoreInfo(33.99, '<', 75);
  });
});

describe('evaluations with key sli', () => {
  beforeEach(() => {
    heatmap.interceptWithKeySli();
  });

  it('should show key sli info', () => {
    heatmap.visitPageWithHeatmapComponent();

    evaluationBoard.assertScoreInfo(50, '<', 75);
    evaluationBoard.assertResultInfo(ResultTypes.FAILED);
    evaluationBoard.assertKeySliInfo('passed');

    heatmap.clickScore('52b4b2c7-fa49-41f3-9b5c-b9aea2370bb4');
    evaluationBoard.assertScoreInfo(75, '>=', 75);
    evaluationBoard.assertResultInfo(ResultTypes.FAILED);
    evaluationBoard.assertKeySliInfo('failed');

    heatmap.clickScore('182d10b8-b68d-49d4-86cd-5521352d7a42');
    evaluationBoard.assertScoreInfo(100, '>=', 90);
    evaluationBoard.assertResultInfo(ResultTypes.PASSED);
    evaluationBoard.assertKeySliInfo('passed');
  });
});
