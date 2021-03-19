import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbKeptnServicesListComponent } from './ktb-keptn-services-list.component';

describe('KtbKeptnServicesListComponent', () => {
  let component: KtbKeptnServicesListComponent;
  let fixture: ComponentFixture<KtbKeptnServicesListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbKeptnServicesListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbKeptnServicesListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
