import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbServiceSettingsOverviewComponent } from './ktb-service-settings-overview.component';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';
import { By } from '@angular/platform-browser';
import { KtbServiceSettingsModule } from '../ktb-service-settings.module';
import { RouterTestingModule } from '@angular/router/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbServiceSettingsOverviewComponent', () => {
  let component: KtbServiceSettingsOverviewComponent;
  let fixture: ComponentFixture<KtbServiceSettingsOverviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbServiceSettingsModule, RouterTestingModule, HttpClientTestingModule],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(
              convertToParamMap({
                projectName: 'sockshop',
              })
            ),
          },
        },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbServiceSettingsOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show create button', () => {
    // given
    const createButton = fixture.debugElement.query(By.css('button'));
    expect(createButton).toBeTruthy();
    expect(createButton.nativeElement.textContent.trim()).toEqual('Create service');
  });
});
