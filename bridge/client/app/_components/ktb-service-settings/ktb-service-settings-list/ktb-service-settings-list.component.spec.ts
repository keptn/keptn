import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbServiceSettingsListComponent } from './ktb-service-settings-list.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbServiceSettingsModule } from '../ktb-service-settings.module';

describe('KtbServiceSettingsListComponent', () => {
  let component: KtbServiceSettingsListComponent;
  let fixture: ComponentFixture<KtbServiceSettingsListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbServiceSettingsModule, HttpClientTestingModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbServiceSettingsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should return empty service names on default', () => {
    expect(component.serviceNames).toStrictEqual([]);
  });

  it('should contain services', () => {
    // given
    component.serviceNames = ['carts-db', 'carts'];

    // when, then
    expect(component.dataSource.data.sort()).toEqual(['carts', 'carts-db']);
  });
});
