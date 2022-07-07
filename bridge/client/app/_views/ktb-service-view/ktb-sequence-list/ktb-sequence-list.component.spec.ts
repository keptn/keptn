import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceListComponent } from './ktb-sequence-list.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { KtbServiceViewModule } from '../ktb-service-view.module';

describe('KtbSequenceListComponent', () => {
  let component: KtbSequenceListComponent;
  let fixture: ComponentFixture<KtbSequenceListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbServiceViewModule, HttpClientTestingModule, RouterTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
