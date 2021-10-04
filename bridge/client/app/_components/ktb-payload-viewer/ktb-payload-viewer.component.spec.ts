import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbPayloadViewerComponent } from './ktb-payload-viewer.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataServiceMock } from '../../_services/data.service.mock';

describe('KtbPayloadViewerComponent', () => {
  let component: KtbPayloadViewerComponent;
  let fixture: ComponentFixture<KtbPayloadViewerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [DataServiceMock],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbPayloadViewerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
