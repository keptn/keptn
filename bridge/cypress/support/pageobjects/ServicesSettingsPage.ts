import { interceptServiceSettings } from '../intercept';

export class ServicesSettingsPage {
  public intercept(): this {
    interceptServiceSettings();
    return this;
  }

  public visitService(project: string, service: string): this {
    cy.visit(`/project/${project}/settings/services/edit/${service}`).wait('@metadata').wait('@projectPlain');
    return this;
  }

  public inputService(serviceName: string): this {
    cy.get('input[formcontrolname="serviceName"]').type(serviceName);
    return this;
  }

  public createService(serviceName?: string): this {
    if (serviceName) {
      this.inputService(serviceName);
    }
    cy.byTestId('createServiceButton').click();
    return this;
  }

  public assertNoFilesMessageExists(status: boolean): this {
    cy.byTestId('ktb-no-files-for-file-tree').should(status ? 'exist' : 'not.exist');
    return this;
  }
}
