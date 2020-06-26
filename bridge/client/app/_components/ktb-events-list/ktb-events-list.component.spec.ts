import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEventsListComponent } from './ktb-events-list.component';
import { AppModule } from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbEventsListComponent', () => {
  let component: KtbEventsListComponent;
  let fixture: ComponentFixture<KtbEventsListComponent>;

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
    fixture = TestBed.createComponent(KtbEventsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create an instance', () => {
    expect(component).toBeTruthy();
  });
});
