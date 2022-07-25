import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSecretsListComponent } from './ktb-secrets-list.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbSecretsListModule } from './ktb-secrets-list.module';
import { RouterTestingModule } from '@angular/router/testing';
import { ApiService } from '../../../_services/api.service';
import { ApiServiceMock } from '../../../_services/api.service.mock';
import { finalize, take } from 'rxjs/operators';
import { SecretsResponseMock } from '../../../_services/_mockData/api-responses/secrets-response.mock';
import { firstValueFrom } from 'rxjs';

describe('KtbUniformViewComponent', () => {
  let component: KtbSecretsListComponent;
  let fixture: ComponentFixture<KtbSecretsListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbSecretsListModule, HttpClientTestingModule, RouterTestingModule],
      providers: [{ provide: ApiService, useClass: ApiServiceMock }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSecretsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load secrets', (done) => {
    // given
    const secretsLength = 3;

    // when
    component.ngOnInit();

    // then
    component.secrets$.pipe(take(1), finalize(done)).subscribe((value) => {
      expect(secretsLength).toBe(value.length);
    });
  });

  it('should confirm a secret deletion', async () => {
    // given
    const secretToDelete = SecretsResponseMock.Secrets[0];

    // when
    component.triggerDeleteSecret(secretToDelete);

    // then
    expect(component.deleteState).toBe('confirm');
    expect(component.currentSecret).toBe(secretToDelete);
  });

  it('should delete a secret', async () => {
    // given
    const secretToDelete = SecretsResponseMock.Secrets[0];
    await firstValueFrom(component.secrets$.pipe(take(1)));
    component.ngOnInit();

    // when
    component.deleteSecret(secretToDelete);

    // then
    const secrets = await firstValueFrom(component.secrets$.pipe(take(1)));
    expect(component.deleteState).toBe('success');
    expect(secrets.length).toBe(2);
    expect(secrets.some((s) => s.name === secretToDelete.name)).toBe(false);
  });
});
