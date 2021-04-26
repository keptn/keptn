import {KtbUserComponent} from './ktb-user.component';
import {ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

describe('ktbUserComponentTest', () => {

  let component: KtbUserComponent;
  let fixture: ComponentFixture<KtbUserComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [KtbUserComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbUserComponent);
    component = fixture.componentInstance;
  });

  it('should be created successfully.', () => {
    fixture.detectChanges();
    expect(component).toBeDefined();
  });

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));

});
