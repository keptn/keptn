import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbDeploymentListComponent } from './ktb-deployment-list.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';

describe('KtbArtifactListComponent', () => {
  let component: KtbDeploymentListComponent;
  let fixture: ComponentFixture<KtbDeploymentListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbDeploymentListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
