import { KtbUserComponent } from './ktb-user.component';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { EndSessionData } from '../../../../shared/interfaces/end-session-data';
import { Location } from '@angular/common';

describe('ktbUserComponentTest', () => {
  let component: KtbUserComponent;
  let fixture: ComponentFixture<KtbUserComponent>;
  let httpMock: HttpTestingController;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbUserComponent);
    component = fixture.componentInstance;
    httpMock = TestBed.inject(HttpTestingController);
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeDefined();
  });

  it('should send POST to logout', () => {
    const submitForm = { target: { submit: (): void => {} } };
    const submitSpy = jest.spyOn(submitForm.target, 'submit');
    component.logout(submitForm);
    httpMock.expectOne('./logout').flush({
      id_token_hint: '',
      end_session_endpoint: '',
      post_logout_redirect_uri: '',
      state: '',
    } as EndSessionData);
    expect(submitSpy).toHaveBeenCalled();
  });

  it('should redirect to root if no data for logout is returned', () => {
    const submitForm = { target: { submit: (): void => {} } };
    const location = TestBed.inject(Location);
    const locationSpy = jest.spyOn(location, 'prepareExternalUrl');
    component.logout(submitForm);
    httpMock.expectOne('./logout').flush(null);
    expect(locationSpy).toBeCalledWith('/logout?status=true');
  });
});
