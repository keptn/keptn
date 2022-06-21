import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RETRY_ON_HTTP_ERROR } from '../../_utils/app.utils';
import { KtbProjectListComponent } from './ktb-project-list.component';
import { KtbProjectListModule } from './ktb-project-list.module';

describe('KtbProjectListComponent', () => {
  let component: KtbProjectListComponent;
  let fixture: ComponentFixture<KtbProjectListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbProjectListModule, HttpClientTestingModule],
      providers: [{ provide: RETRY_ON_HTTP_ERROR, useValue: false }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbProjectListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
