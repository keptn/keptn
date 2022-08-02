/// <reference types="cypress" />

import { interceptSecrets } from '../intercept';

class SecretsPage {
  private readonly REMOVE_SECRET_KEY_VALUE_PAIR_ID = 'keptn-secret-remove-pair-button';
  private readonly SECRET_NAME_ID = 'keptn-secret-name-input';
  private readonly SECRET_SCOPE_ID = 'keptn-secret-scope-input';
  private readonly CREATE_SECRET_ID = 'keptn-secret-create-button';
  private readonly ADD_SECRET_ID = 'keptn-add-secret';
  private readonly SECRET_KEY_ID = 'keptn-secret-key-input';
  private readonly SECRET_VALUE_ID = 'keptn-secret-value-input';
  private readonly ADD_KEY_VALUE_PAIR_ID = 'keptn-secret-add-pair-button';

  public intercept(): this {
    interceptSecrets();
    return this;
  }

  public visit(projectName: string): this {
    cy.visit(`/project/${projectName}/settings/uniform/secrets`);
    return this;
  }

  public visitCreate(projectName: string): this {
    cy.visit(`/project/${projectName}/settings/uniform/secrets/add`);
    return this;
  }

  public assertKeyValuePairLength(length: number): this {
    cy.byTestId(this.REMOVE_SECRET_KEY_VALUE_PAIR_ID).should('have.length', length);
    return this;
  }

  public assertKeyValuePairEnabled(index: number, status: boolean): this {
    cy.byTestId(this.REMOVE_SECRET_KEY_VALUE_PAIR_ID)
      .eq(index)
      .should(status ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertSecretInList(index: number, name: string, scope: string, key: string): this {
    return this.assertSecretName(index, name).assertScope(index, scope).assertSecretKey(index, key);
  }

  public assertSecretName(index: number, secretName: string): this {
    cy.get('dt-row').eq(index).find('dt-cell').eq(0).find('p').should('have.text', secretName);
    return this;
  }

  public assertScopesEnabled(status: boolean): this {
    cy.byTestId(this.SECRET_SCOPE_ID).should(status ? 'not.have.class' : 'have.class', 'dt-select-disabled');
    return this;
  }

  public assertCreateButtonEnabled(status: boolean): this {
    cy.byTestId(this.CREATE_SECRET_ID).should(status ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertScope(index: number, scope: string): this {
    cy.get('dt-row').eq(index).find('dt-cell').eq(1).find('p').should('have.text', scope);
    return this;
  }

  public assertSecretKey(index: number, key: string): this {
    cy.get('dt-row').eq(index).find('dt-cell').eq(2).find('p').should('contain.text', key);
    return this;
  }

  public secretExistsInList(secretName: string, index: number): this {
    cy.get('dt-row').eq(index).find('dt-cell').eq(0).find('p').should('not.have.text', secretName);
    return this;
  }

  public appendSecretName(name: string): this {
    cy.byTestId(this.SECRET_NAME_ID).type(name);
    return this;
  }

  public setSecretScope(scope: string): this {
    cy.byTestId(this.SECRET_SCOPE_ID).type(scope).type('{enter}');
    return this;
  }

  public appendSecretKey(index: number, key: string): this {
    cy.byTestId(this.SECRET_KEY_ID).eq(index).type(key);
    return this;
  }

  public appendSecretValue(index: number, value: string): this {
    cy.byTestId(this.SECRET_VALUE_ID).eq(index).type(value);
    return this;
  }

  public createSecret(): this {
    cy.byTestId(this.CREATE_SECRET_ID).click();
    return this;
  }

  public clickAddSecret(): this {
    cy.byTestId(this.ADD_SECRET_ID).click();
    return this;
  }

  public setSecret(name: string, scope: string, key: string, value: string): this {
    return this.appendSecretName(name).setSecretScope(scope).appendSecretKey(0, key).appendSecretValue(0, value);
  }

  public deleteSecret(SECRET_NAME: string | number | RegExp): this {
    cy.get('dt-row.dt-row > dt-cell > p')
      .contains(SECRET_NAME)
      .parent()
      .nextAll('.dt-table-column-action')
      .children('button')
      .click();
    cy.get('span.dt-button-label').contains('Delete').click();
    return this;
  }

  public addKeyValuePair(): this {
    cy.byTestId(this.ADD_KEY_VALUE_PAIR_ID).click();
    return this;
  }
}

export default SecretsPage;
