import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import { DashboardComponent } from './dashboard.component';
import { AppModule } from '../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppHeaderComponent} from "../app-header/app-header.component";

describe('DashboardComponent', () => {
  let component: DashboardComponent;
  let fixture: ComponentFixture<DashboardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(DashboardComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
