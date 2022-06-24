export class EvaluationTileComponentPage {
  public assertShowSLOButtonEnabled(enabled: boolean): this {
    cy.byTestId('ktb-show-slo-button').should(enabled ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertShowSLOButtonExists(exists: boolean): this {
    cy.byTestId('ktb-show-slo-button').should(exists ? 'exist' : 'not.exist');
    return this;
  }

  public assertSLOButtonOverlayExists(exists: boolean): this {
    cy.byTestId('ktb-show-slo-button').parent().trigger('mouseenter');
    cy.byTestId('ktb-invalid-slo-overlay').should(exists ? 'exist' : 'not.exist');
    cy.byTestId('ktb-show-slo-button').parent().trigger('mouseleave');
    return this;
  }
}
