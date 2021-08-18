import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing';

import { KtbSequenceListComponent } from './ktb-sequence-list.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbSequenceListComponent', () => {
  let component: KtbSequenceListComponent;
  let fixture: ComponentFixture<KtbSequenceListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbSequenceListComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
