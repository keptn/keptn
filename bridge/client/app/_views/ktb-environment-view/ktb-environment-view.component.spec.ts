import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbEnvironmentViewComponent } from './ktb-environment-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbEnvironmentViewModule } from './ktb-environment-view.module';
import { RouterTestingModule } from '@angular/router/testing';

describe('KtbEnvironmentViewComponent', () => {
  let component: KtbEnvironmentViewComponent;
  let fixture: ComponentFixture<KtbEnvironmentViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbEnvironmentViewModule, RouterTestingModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbEnvironmentViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
