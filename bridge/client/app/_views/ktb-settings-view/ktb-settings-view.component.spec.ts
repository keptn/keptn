import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbSettingsViewComponent } from './ktb-settings-view.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {DataService} from "../../_services/data.service";
import {DataServiceMock} from "../../_services/data.service.mock";

describe('KtbSettingsViewComponent', () => {
  let component: KtbSettingsViewComponent;
  let fixture: ComponentFixture<KtbSettingsViewComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbSettingsViewComponent ],
      imports: [ AppModule, HttpClientTestingModule ],
      providers: [
        {provide: DataService, useClass: DataServiceMock}
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSettingsViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
