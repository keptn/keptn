import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbServicesListComponent } from './ktb-services-list.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';

describe('KtbServicesListComponent', () => {
  let component: KtbServicesListComponent;
  let fixture: ComponentFixture<KtbServicesListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbServicesListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
