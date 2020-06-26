import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import {KtbEventItemComponent} from './ktb-event-item.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbEventItemComponent', () => {
  let component: KtbEventItemComponent;
  let fixture: ComponentFixture<KtbEventItemComponent>;

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
    fixture = TestBed.createComponent(KtbEventItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create an instance', () => {
    expect(component).toBeTruthy();
  });
});
