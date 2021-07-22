import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbDangerZoneComponent } from './ktb-danger-zone.component';
import { AppModule } from '../../app.module';
import { DeleteType } from '../../_interfaces/delete';

describe('KtbDeletionComponent', () => {
  let component: KtbDangerZoneComponent;
  let fixture: ComponentFixture<KtbDangerZoneComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KtbDangerZoneComponent],
      imports: [AppModule],
    })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbDangerZoneComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should open the deletion dialog', () => {
    // given
    const data = {name: 'sockshop', type: DeleteType.PROJECT};
    component.data = data;
    fixture.detectChanges();
    const button = fixture.nativeElement.querySelector('.danger-button');
    const spy = spyOn(component.dialog, 'open');

    // when
    button.click();
    fixture.detectChanges();

    // then
    expect(spy).toHaveBeenCalled();
    const spyData = spy.calls.mostRecent().args[1]?.data;
    expect(spyData).toEqual(data);
  });
});
