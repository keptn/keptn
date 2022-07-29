import { interceptHeatmapComponent } from '../intercept';

export class LineChartComponentPage {
  public intercept(): this {
    interceptHeatmapComponent();
    return this;
  }

  public visitPageWithHeatmapComponent(): this {
    cy.visit('/project/sockshop/service/carts/context/da740469-9920-4e0c-b304-0fd4b18d17c2/stage/staging').wait(
      '@heatmapEvaluations'
    );
    return this;
  }

  public selectLineChart(): this {
    cy.byTestId('keptn-evaluation-details-contextButtons-chart').click();
    return this;
  }
}
