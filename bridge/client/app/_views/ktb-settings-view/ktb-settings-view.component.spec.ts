import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSettingsViewComponent } from './ktb-settings-view.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbSettingsViewComponent', () => {
  let component: KtbSettingsViewComponent;
  let fixture: ComponentFixture<KtbSettingsViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbSettingsViewComponent ],
      imports: [ AppModule, HttpClientTestingModule ],
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
