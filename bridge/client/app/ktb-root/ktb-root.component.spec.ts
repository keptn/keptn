import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbRootComponent } from './ktb-root.component';
import { AppModule } from '../app.module';

describe('KtbRootComponent', () => {
  let component: KtbRootComponent;
  let fixture: ComponentFixture<KtbRootComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbRootComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
