import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import {KtbTaskItemComponent} from './ktb-task-item.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbEventItemComponent', () => {
  let component: KtbTaskItemComponent;
  let fixture: ComponentFixture<KtbTaskItemComponent>;

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
    fixture = TestBed.createComponent(KtbTaskItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
