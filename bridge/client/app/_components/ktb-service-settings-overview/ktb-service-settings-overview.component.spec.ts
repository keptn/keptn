import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbServiceSettingsOverviewComponent } from './ktb-service-settings-overview.component';
import { AppModule } from '../../app.module';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';
import { By } from '@angular/platform-browser';

describe('KtbServiceSettingsOverviewComponent', () => {
  let component: KtbServiceSettingsOverviewComponent;
  let fixture: ComponentFixture<KtbServiceSettingsOverviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
      ],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(convertToParamMap({
              projectName: 'sockshop',
            })),
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

  it('should show createButton', () => {
    // given
    const createButton = fixture.debugElement.query(By.css('button'));
    expect(createButton).toBeTruthy();
    expect(createButton.nativeElement.textContent).toEqual('Create service');
  });
});
