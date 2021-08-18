import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';
import { KtbEventsListComponent } from './ktb-events-list.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RETRY_ON_HTTP_ERROR } from '../../_utils/app.utils';

describe('KtbEventsListComponent', () => {
  let component: KtbEventsListComponent;
  let fixture: ComponentFixture<KtbEventsListComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
      providers: [
        {provide: RETRY_ON_HTTP_ERROR, useValue: false}
      ]
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbEventsListComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
