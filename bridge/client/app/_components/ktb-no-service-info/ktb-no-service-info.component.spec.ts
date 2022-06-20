import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbNoServiceInfoComponent } from './ktb-no-service-info.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';
import { KtbNoServiceInfoModule } from './ktb-no-service-info.module';

describe('KtbNoServiceInfoComponent', () => {
  let component: KtbNoServiceInfoComponent;
  let fixture: ComponentFixture<KtbNoServiceInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbNoServiceInfoModule, HttpClientTestingModule],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({}),
            params: of({}),
            paramMap: of(
              convertToParamMap({
                projectName: 'sockshop',
              })
            ),
            queryParams: of({}),
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbNoServiceInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should have the right routerLink', () => {
    expect(component.createServiceLink).toBe('/project/sockshop/settings/services/create');
  });
});
