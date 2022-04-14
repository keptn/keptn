import { ErrorPage } from './ErrorPage';

export class LogoutPage extends ErrorPage {
  public visit(): this {
    cy.visit('/logoutsession');
    return this;
  }
}
