import { KtbUserComponent } from './ktb-user.component';
import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('ktbUserComponentTest', () => {

  let component: KtbUserComponent;
  let fixture: ComponentFixture<KtbUserComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbUserComponent);
        component = fixture.componentInstance;
      });
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeDefined();
  });

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
