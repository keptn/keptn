import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { FormControl } from '@angular/forms';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { Secret } from '../../_models/secret';
import { SecretScope } from '../../../../shared/interfaces/secret-scope';
import { KtbSecretSelectorComponent } from './ktb-secret-selector.component';

describe('KtbSecretSelectorComponent', () => {
  const secretPath = 'SecretA.key1';
  let component: KtbSecretSelectorComponent;
  let fixture: ComponentFixture<KtbSecretSelectorComponent>;

  const secretDataSource = [
    {
      name: 'SecretA',
      keys: [
        {
          name: 'key1',
          path: 'SecretA.key1',
        },
        {
          name: 'key2',
          path: 'SecretA.key2',
        },
        {
          name: 'key3',
          path: 'SecretA.key3',
        },
      ],
    },
    {
      name: 'SecretB',
      keys: [
        {
          name: 'key1',
          path: 'SecretB.key1',
        },
        {
          name: 'key2',
          path: 'SecretB.key2',
        },
        {
          name: 'key3',
          path: 'SecretB.key3',
        },
      ],
    },
  ];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
      providers: [{ provide: DataService, useClass: DataServiceMock }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSecretSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should insert the processed string as value to the control', () => {
    // given
    component.control = new FormControl('');
    component.control.setValue('');
    component.selectionStart = 0;

    // when
    component.setSecret(secretPath);

    // then
    expect(component.control.value).toEqual(`{{.secret.${secretPath}}}`);
  });

  it('should insert the processed string as value to the control at the given position', () => {
    // given
    component.control = new FormControl('');
    component.control.setValue('https://example.com?somestringtoinsert');
    component.selectionStart = 30;

    // when
    component.setSecret(secretPath);

    // then
    expect(component.control.value).toEqual(`https://example.com?somestring{{.secret.${secretPath}}}toinsert`);
  });

  it('should map secrets to a tree when set', () => {
    // given, when
    const secrets = [new Secret(), new Secret()];
    secrets[0].name = 'SecretA';
    secrets[0].scope = SecretScope.WEBHOOK;
    secrets[0].keys = ['key1', 'key2', 'key3'];
    secrets[1].name = 'SecretB';
    secrets[1].scope = SecretScope.WEBHOOK;
    secrets[1].keys = ['key1', 'key2', 'key3'];

    // when
    component.secrets = secrets;

    // then
    expect(component.secretDataSource).toEqual(secretDataSource);
  });
});
