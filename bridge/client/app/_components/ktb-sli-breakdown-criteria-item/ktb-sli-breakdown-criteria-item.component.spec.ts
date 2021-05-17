import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbSliBreakdownCriteriaItemComponent } from './ktb-sli-breakdown-criteria-item.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {KtbServicesListComponent} from "../ktb-services-list/ktb-services-list.component";

describe('KtbSliBreakdownCriteriaItemComponent', () => {
  let component: KtbSliBreakdownCriteriaItemComponent;
  let fixture: ComponentFixture<KtbSliBreakdownCriteriaItemComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbSliBreakdownCriteriaItemComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
