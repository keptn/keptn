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

  it('should show key sli info result failed 2 key slis', () => {
    heatmap.visitPageWithHeatmapComponent();
    evaluationBoard.assertScoreInfo(25, '<', 75).assertResultInfo(ResultTypes.FAILED).assertKeySliInfo(2);
  });

  it('should show key sli info result failed 0 key slis', () => {
    heatmap.visitPageWithHeatmapComponent().clickScore('1500b971-bfc3-4e20-8dc1-a624e0faf961');
    evaluationBoard.assertScoreInfo(50, '<', 75).assertResultInfo(ResultTypes.FAILED).assertKeySliInfo(0);
  });

  it('should show key sli info result failed 1 key sli', () => {
    heatmap.visitPageWithHeatmapComponent().clickScore('52b4b2c7-fa49-41f3-9b5c-b9aea2370bb4');
    evaluationBoard.assertScoreInfo(75, '>=', 75).assertResultInfo(ResultTypes.FAILED).assertKeySliInfo(1);
  });

  it('should show key sli info result passed 0 key sli', () => {
    heatmap.visitPageWithHeatmapComponent().clickScore('182d10b8-b68d-49d4-86cd-5521352d7a42');
    evaluationBoard.assertScoreInfo(100, '>=', 90).assertResultInfo(ResultTypes.PASSED).assertKeySliInfo(0);
  });
});
