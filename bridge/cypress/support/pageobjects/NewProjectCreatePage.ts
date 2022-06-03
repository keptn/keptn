/// <reference types="cypress" />

import {
  interceptCreateProject,
  interceptMain,
  interceptMainResourceApEnabled,
  interceptMainResourceEnabled,
  interceptProjectBoard,
  interceptProjectSettings,
} from '../intercept';

class NewProjectCreatePage {
  private validCertificateInput = '-----BEGIN CERTIFICATE-----\nmyCertificate\n-----END CERTIFICATE-----';
  private validPrivateKeyInput = '-----BEGIN OPENSSH PRIVATE KEY-----\nmyPrivateKey\n-----END OPENSSH PRIVATE KEY-----';

  public intercept(resourceServiceEnabled = false, automaticProvisioningEnabled = false): this {
    if (resourceServiceEnabled) {
      if (!automaticProvisioningEnabled) {
        interceptMainResourceEnabled();
      } else {
        interceptMainResourceApEnabled();
      }
    } else {
      interceptMain();
    }
    interceptCreateProject();
    return this;
  }

  public interceptSettings(resourceServiceEnabled = false, automaticProvisioningEnabled = false): this {
    interceptProjectBoard();
    if (resourceServiceEnabled) {
      if (!automaticProvisioningEnabled) {
        interceptMainResourceEnabled();
      } else {
        interceptMainResourceApEnabled();
      }
    } else {
      interceptMain();
    }
    interceptProjectSettings();
    return this;
  }

  public visit(): this {
    cy.visit('/create/project').wait('@metadata');
    return this;
  }

  public visitSettings(project: string): this {
    cy.visit(`/project/${project}/settings/project`).wait('@metadata');
    return this;
  }

  public setShipyardFile(): this {
    cy.byTestId('ktb-shipyard-file-input').attachFile('shipyard.yaml');
    return this;
  }

  public typeProjectName(projectName: string): this {
    cy.byTestId('ktb-project-name-input').type(projectName);
    return this;
  }

  public typeGitUrl(url: string): this {
    cy.byTestId('ktb-git-url-input').type(url);
    return this;
  }

  public assertGitUrl(url: string): this {
    cy.byTestId('ktb-git-url-input').should('have.value', url);
    return this;
  }

  public typeGitUrlSsh(url: string): this {
    cy.byTestId('ktb-ssh-git-url-input').type(url);
    return this;
  }

  public assertGitUrlSsh(url: string): this {
    cy.byTestId('ktb-ssh-git-url-input').should('have.value', url);
    return this;
  }

  public typeGitUsername(username: string): this {
    cy.byTestId('ktb-git-username-input').type(username);
    return this;
  }

  public typeGitUsernameSsh(username: string): this {
    cy.byTestId('ktb-ssh-git-username-input').type(username);
    return this;
  }

  public assertGitUsernameSsh(username: string): this {
    cy.byTestId('ktb-ssh-git-username-input').should('have.value', username);
    return this;
  }

  public clearGitUsernameSsh(): this {
    cy.byTestId('ktb-ssh-git-username-input').clear();
    return this;
  }

  public typeGitToken(token: string): this {
    cy.byTestId('ktb-git-token-input').type(token);
    return this;
  }

  public assertGitToken(token: string): this {
    cy.byTestId('ktb-git-token-input').should('have.value', token);
    return this;
  }

  public clearGitToken(): this {
    cy.byTestId('ktb-git-token-input').clear();
    return this;
  }

  public typeCertificate(certificate: string): this {
    cy.byTestId('ktb-certificate-input').type(certificate);
    return this;
  }

  public assertCertificate(certificate: string): this {
    cy.byTestId('ktb-certificate-input').should('have.value', certificate);
    return this;
  }

  public typeValidCertificate(): this {
    return this.typeCertificate(this.validCertificateInput);
  }

  public typeInvalidCertificate(): this {
    return this.typeCertificate('myInvalidCertificate-----END CERTIFICATE-----');
  }

  public typeProxyPort(port: string | number): this {
    cy.byTestId('ktb-proxy-port-input').type(port.toString());
    return this;
  }

  public clearProxyPort(): this {
    cy.byTestId('ktb-proxy-port-input').clear();
    return this;
  }

  public assertProxyPort(port: string | number): this {
    cy.byTestId('ktb-proxy-port-input').should('have.value', port.toString());
    return this;
  }

  public assertProxyUrl(url: string): this {
    cy.byTestId('ktb-proxy-url-input').should('have.value', url);
    return this;
  }

  public typeProxyUrl(url: string): this {
    cy.byTestId('ktb-proxy-url-input').type(url);
    return this;
  }

  public clearProxyUrl(): this {
    cy.byTestId('ktb-proxy-url-input').clear();
    return this;
  }

  public typeProxyUsername(username: string): this {
    cy.byTestId('ktb-proxy-username-input').type(username);
    return this;
  }

  public clearProxyUsername(): this {
    cy.byTestId('ktb-proxy-username-input').clear();
    return this;
  }

  public assertProxyUsername(username: string): this {
    cy.byTestId('ktb-proxy-username-input').should('have.value', username);
    return this;
  }

  public typeProxyPassword(password: string): this {
    cy.byTestId('ktb-proxy-password-input').type(password);
    return this;
  }

  public clearProxyPassword(): this {
    cy.byTestId('ktb-proxy-password-input').clear();
    return this;
  }

  public assertProxyPassword(password: string): this {
    cy.byTestId('ktb-proxy-password-input').should('have.value', password);
    return this;
  }

  public typeSshPrivateKey(key: string): this {
    cy.byTestId('ktb-ssh-private-key-input').type(key);
    return this;
  }

  public assertSshPrivateKey(key: string): this {
    cy.byTestId('ktb-ssh-private-key-input').should('have.value', key);
    return this;
  }

  public clearSshPrivateKey(): this {
    cy.byTestId('ktb-ssh-private-key-input').clear();
    return this;
  }

  public typeSshPrivateKeyPassphrase(passphrase: string): this {
    cy.byTestId('ktb-ssh-private-key-passphrase-input').type(passphrase);
    return this;
  }

  public clearSshPrivateKeyPassphrase(): this {
    cy.byTestId('ktb-ssh-private-key-passphrase-input').clear();
    return this;
  }

  public assertSshPrivateKeyPassphrase(passphrase: string): this {
    cy.byTestId('ktb-ssh-private-key-passphrase-input').should('have.value', passphrase);
    return this;
  }

  public setProxyInsecure(status: boolean): this {
    cy.byTestId('ktb-proxy-insecure').dtCheck(status);
    return this;
  }

  public assertProxyInsecure(status: boolean): this {
    cy.byTestId('ktb-proxy-insecure')
      .find('input')
      .should(status ? 'be.checked' : 'not.be.checked');
    return this;
  }

  public selectProxyScheme(scheme: 'HTTP' | 'HTTPS'): this {
    cy.byTestId('ktb-proxy-scheme').dtSelect(scheme);
    return this;
  }

  public assertProxyScheme(scheme: 'HTTPS' | 'HTTP'): this {
    cy.byTestId('ktb-proxy-scheme').should('have.text', scheme);
    return this;
  }

  public typeValidSshPrivateKey(): this {
    return this.typeSshPrivateKey(this.validPrivateKeyInput);
  }

  public typeInvalidSshPrivateKey(): this {
    return this.typeSshPrivateKey('my-invalid-input');
  }

  public setValidSshPrivateKeyFile(): this {
    cy.byTestId('ktb-ssh-private-key-file-input').attachFile('files/ssh-private-key.pem');
    return this;
  }

  public setInvalidSshPrivateKeyFile(): this {
    cy.byTestId('ktb-ssh-private-key-file-input').attachFile('files/ssh-private-key-invalid.pem');
    return this;
  }

  public setValidCertificateFile(): this {
    cy.byTestId('ktb-certificate-file-input').attachFile('files/certificate.pem');
    return this;
  }

  public setInvalidCertificateFile(): this {
    cy.byTestId('ktb-certificate-file-input').attachFile('files/certificate-invalid.pem');
    return this;
  }

  public setEnableProxy(status: boolean): this {
    cy.byTestId('ktb-enable-git-proxy').dtCheck(status);
    return this;
  }

  public assertProxyEnabled(status: boolean): this {
    cy.byTestId('ktb-enable-git-proxy')
      .find('input')
      .should(status ? 'be.checked' : 'not.be.checked');
    return this;
  }

  public enterBasicValidProjectWithoutGitUpstream(): this {
    return this.typeProjectName('my-project').setShipyardFile();
  }

  public enterBasicSsh(): this {
    return this.typeGitUrlSsh('ssh://example.com').typeValidSshPrivateKey();
  }

  public enterBasicValidProjectSsh(fillPrivateKey = true): this {
    this.typeProjectName('my-project').setShipyardFile().typeGitUrlSsh('ssh://example.com');
    if (fillPrivateKey) {
      return this.typeValidSshPrivateKey();
    } else {
      return this;
    }
  }

  public enterFullValidProjectSsh(): this {
    return this.enterBasicValidProjectSsh()
      .typeGitUsernameSsh('myUsername')
      .typeSshPrivateKeyPassphrase('myPassphrase');
  }

  public validateFullValidProjectSsh(): this {
    return this.assertGitUrlSsh('ssh://example.com')
      .assertGitUsernameSsh('myUsername')
      .assertSshPrivateKey(this.validPrivateKeyInput)
      .assertSshPrivateKeyPassphrase('myPassphrase');
  }

  public enterBasicHttps(): this {
    return this.typeGitUrl('https://example.com').typeGitToken('myToken');
  }

  public enterBasicValidProjectHttps(): this {
    return this.typeProjectName('my-project').setShipyardFile().enterBasicHttps();
  }

  public enterFullValidProjectHttps(): this {
    return this.enterBasicValidProjectHttps()
      .typeGitUsername('myUsername')
      .typeValidCertificate()
      .setEnableProxy(true)
      .typeProxyUrl('0.0.0.0')
      .typeProxyPort(5000)
      .typeProxyUsername('myProxyUsername')
      .typeProxyPassword('myPassword')
      .setProxyInsecure(true)
      .selectProxyScheme('HTTP');
  }

  public validateFullValidProjectHttps(): this {
    return this.assertGitUrl('https://example.com')
      .assertGitToken('myToken')
      .assertGitUsername('myUsername')
      .assertCertificate(this.validCertificateInput)
      .assertProxyEnabled(true)
      .assertProxyUrl('0.0.0.0')
      .assertProxyPort(5000)
      .assertProxyUsername('myProxyUsername')
      .assertProxyPassword('myPassword')
      .assertProxyInsecure(true)
      .assertProxyScheme('HTTP');
  }

  public assertGitUsername(username: string): this {
    cy.byTestId('ktb-git-username-input').should('have.value', username);
    return this;
  }

  public assertCreateButtonEnabled(status: boolean): this {
    cy.byTestId('ktb-create-project').should(status ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertGitUpstreamHeadlineExistsOnce(): this {
    cy.get('ktb-project-settings').find('h2').filter(':contains("Git upstream repository")').should('have.length', 1);
    return this;
  }

  public assertSshFormExists(status: boolean): this {
    cy.get('ktb-project-settings-git-ssh').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertHttpsFormExists(status: boolean): this {
    cy.get('ktb-project-settings-git-https').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public assertSshFormVisible(status: boolean): this {
    cy.get('ktb-project-settings-git-ssh').should(status ? 'be.visible' : 'not.be.visible');
    return this;
  }

  public assertHttpsFormVisible(status: boolean): this {
    cy.get('ktb-project-settings-git-https').should(status ? 'be.visible' : 'not.be.visible');
    return this;
  }

  public assertProxyFormVisible(status: boolean): this {
    cy.get('ktb-proxy-input').should(status ? 'be.visible' : 'not.be.visible');
    return this;
  }

  public assertNoUpstreamSelected(status: boolean): this {
    cy.byTestId('ktb-no-upstream-form-button').should(status ? 'have.class' : 'not.have.class', 'dt-radio-checked');
    return this;
  }

  public assertNoUpstreamEnabled(status: boolean): this {
    cy.byTestId('ktb-no-upstream-form-button')
      .get('input')
      .should(status ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public selectHttpsForm(): this {
    cy.byTestId('ktb-https-form-button').click();
    return this;
  }

  public selectSshForm(): this {
    cy.byTestId('ktb-ssh-form-button').click();
    return this;
  }

  public selectNoUpstreamForm(): this {
    cy.byTestId('ktb-no-upstream-form-button').click();
    return this;
  }

  public clickCreateProject(): this {
    cy.byTestId('ktb-create-project').click();
    return this;
  }

  public assertUpdateButtonExists(status: boolean): this {
    cy.byTestId('ktb-project-update-button').should(status ? 'exist' : 'not.exist');
    return this;
  }

  public updateProject(): this {
    cy.byTestId('ktb-project-update-button').click();
    return this;
  }

  public assertUpdateButtonEnabled(status: boolean): this {
    cy.byTestId('ktb-project-update-button').should(status ? 'be.enabled' : 'be.disabled');
    return this;
  }

  public assertGitUpstreamMessageContains(message: string): this {
    cy.byTestId('ktb-settings-git-upstream-message').should('contain', message);
    return this;
  }

  public assertErrorVisible(status: boolean): this {
    cy.get('ktb-error-view').should(status ? 'be.visible' : 'not.be.visible');
    return this;
  }
}

export default NewProjectCreatePage;
