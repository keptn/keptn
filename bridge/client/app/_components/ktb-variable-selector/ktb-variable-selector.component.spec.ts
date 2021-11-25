import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { FormControl } from '@angular/forms';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { KtbVariableSelectorComponent } from './ktb-variable-selector.component';

describe('KtbSecretSelectorComponent', () => {
  const variablePath = '.secret.SecretA.key1';
  let component: KtbVariableSelectorComponent;
  let fixture: ComponentFixture<KtbVariableSelectorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
      providers: [{ provide: DataService, useClass: DataServiceMock }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbVariableSelectorComponent);
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
    component.setVariable(variablePath);

    // then
    expect(component.control.value).toEqual(`{{${variablePath}}}`);
  });

  it('should insert the processed string as value to the control at the given position', () => {
    // given
    component.control = new FormControl('');
    component.control.setValue('https://example.com?somestringtoinsert');
    component.selectionStart = 30;

    // when
    component.setVariable(variablePath);

    // then
    expect(component.control.value).toEqual(`https://example.com?somestring{{${variablePath}}}toinsert`);
  });
});
