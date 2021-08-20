import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbUniformSubscriptionsComponent } from './ktb-uniform-subscriptions.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';

describe('KtbUniformSubscriptionsComponent', () => {
  let component: KtbUniformSubscriptionsComponent;
  let fixture: ComponentFixture<KtbUniformSubscriptionsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
      providers: [
        {provide: DataService, useClass: DataServiceMock},
        {
          provide: ActivatedRoute,
          useValue: {
            data: of({}),
            params: of({}),
            paramMap: of(convertToParamMap({
              projectName: 'sockshop'
            })),
            queryParams: of({})
          }
        }
      ]
    }).compileComponents();
    fixture = TestBed.createComponent(KtbUniformSubscriptionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
