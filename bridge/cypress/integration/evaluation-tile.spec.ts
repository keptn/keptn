import { HeatmapComponentPage } from '../support/pageobjects/HeatmapComponentPage';
import { EvaluationTileComponentPage } from '../support/pageobjects/EvaluationTileComponentPage';

describe('SLO file content', () => {
  const heatmap = new HeatmapComponentPage();
  const evaluationTileComponentPage = new EvaluationTileComponentPage();

  it('should show enabled "Show SLO" button', () => {
    heatmap
      .interceptWithSLO(
        'LS0tDQpzcGVjX3ZlcnNpb246ICcwLjEuMCcNCmNvbXBhcmlzb246DQogIGNvbXBhcmVfd2l0aDogInNpbmdsZV9yZXN1bHQiDQogIGluY2x1ZGVfcmVzdWx0X3dpdGhfc2NvcmU6ICJwYXNzIg0KICBhZ2dyZWdhdGVfZnVuY3Rpb246IGF2Zw0Kb2JqZWN0aXZlczoNCiAgLSBzbGk6IGh0dHBfcmVzcG9uc2VfdGltZV9zZWNvbmRzX21haW5fcGFnZV9zdW0NCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PTEiDQogICAgd2FybmluZzoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0wLjUiDQogIC0gc2xpOiByZXF1ZXN0X3Rocm91Z2hwdXQNCiAgICBwYXNzOg0KICAgICAgLSBjcml0ZXJpYToNCiAgICAgICAgICAtICI8PSsxMDAlIg0KICAgICAgICAgIC0gIj49LTgwJSINCiAgLSBzbGk6IGdvX3JvdXRpbmVzDQogICAgcGFzczoNCiAgICAgIC0gY3JpdGVyaWE6DQogICAgICAgICAgLSAiPD0xMDAiDQp0b3RhbF9zY29yZToNCiAgcGFzczogIjkwJSINCiAgd2FybmluZzogIjc1JSI='
      )
      .visitPageWithHeatmapComponent();
    evaluationTileComponentPage.assertShowSLOButtonEnabled(true).assertSLOButtonOverlayExists(false);
  });

  it('should not show "Show SLO" button', () => {
    heatmap.interceptWithSLO().visitPageWithHeatmapComponent();
    evaluationTileComponentPage.assertShowSLOButtonExists(false);
  });

  it('should have overlay and disabled "Show SLO" button', () => {
    heatmap.interceptWithSLO('_invalid_').visitPageWithHeatmapComponent();
    evaluationTileComponentPage.assertShowSLOButtonEnabled(false).assertSLOButtonOverlayExists(true);
  });
});
