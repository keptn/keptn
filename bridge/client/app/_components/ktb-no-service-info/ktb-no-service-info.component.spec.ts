import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbNoServiceInfoComponent } from './ktb-no-service-info.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';

describe('KtbNoServiceInfoComponent', () => {
  let component: KtbNoServiceInfoComponent;
  let fixture: ComponentFixture<KtbNoServiceInfoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
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
    expect(component.routerLink).toBe('/project/sockshop/settings/services/create');
  });
});
