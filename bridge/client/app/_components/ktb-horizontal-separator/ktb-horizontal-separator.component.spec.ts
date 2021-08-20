import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbHorizontalSeparatorComponent } from './ktb-horizontal-separator.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';

describe('KtbHorizontalSeparatorComponent', () => {
  let component: KtbHorizontalSeparatorComponent;
  let fixture: ComponentFixture<KtbHorizontalSeparatorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbHorizontalSeparatorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
