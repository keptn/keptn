import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbHttpLoadingBarComponent } from './ktb-http-loading-bar.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';

describe('HttpLoadingBarComponent', () => {
  let component: KtbHttpLoadingBarComponent;
  let fixture: ComponentFixture<KtbHttpLoadingBarComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbHttpLoadingBarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
