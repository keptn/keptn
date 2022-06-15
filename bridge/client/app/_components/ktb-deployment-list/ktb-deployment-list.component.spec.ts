import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbDeploymentListComponent } from './ktb-deployment-list.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbDeploymentListModule } from './ktb-deployment-list.module';
import { RouterTestingModule } from '@angular/router/testing';

describe('KtbDeploymentListComponent', () => {
  let component: KtbDeploymentListComponent;
  let fixture: ComponentFixture<KtbDeploymentListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbDeploymentListModule, HttpClientTestingModule, RouterTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbDeploymentListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
