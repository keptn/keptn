import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceListComponent } from './ktb-sequence-list.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbSequenceListModule } from './ktb-sequence-list.module';
import { RouterTestingModule } from '@angular/router/testing';

describe('KtbSequenceListComponent', () => {
  let component: KtbSequenceListComponent;
  let fixture: ComponentFixture<KtbSequenceListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbSequenceListModule, HttpClientTestingModule, RouterTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
