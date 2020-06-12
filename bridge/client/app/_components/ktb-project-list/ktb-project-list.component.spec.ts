import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectListComponent } from './ktb-project-list.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";

describe('KtbProjectListComponent', () => {
  let component: KtbProjectListComponent;
  let fixture: ComponentFixture<KtbProjectListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
      ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbProjectListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create an instance', () => {
    expect(component).toBeTruthy();
  });
});
