import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AppHeaderComponent } from './app-header.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../app.module';
import { RETRY_ON_HTTP_ERROR } from '../_utils/app.utils';

describe('AppHeaderComponent', () => {
  let component: AppHeaderComponent;
  let fixture: ComponentFixture<AppHeaderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [{ provide: RETRY_ON_HTTP_ERROR, useValue: false }],
    }).compileComponents();

    fixture = TestBed.createComponent(AppHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
