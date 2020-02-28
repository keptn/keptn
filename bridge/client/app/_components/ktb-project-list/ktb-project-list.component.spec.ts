import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectListComponent } from './ktb-project-list.component';

describe('KtbProjectListComponent', () => {
  let component: KtbProjectListComponent;
  let fixture: ComponentFixture<KtbProjectListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbProjectListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbProjectListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
