import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbLogoutViewComponent } from './ktb-logout-view.component';
import { AppModule } from '../../app.module';

describe('KtbLogoutViewComponent', () => {
  let component: KtbLogoutViewComponent;
  let fixture: ComponentFixture<KtbLogoutViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbLogoutViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
