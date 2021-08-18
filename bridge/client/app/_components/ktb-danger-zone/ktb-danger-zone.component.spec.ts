import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbDangerZoneComponent } from './ktb-danger-zone.component';
import { AppModule } from '../../app.module';
import { DeleteType } from '../../_interfaces/delete';
import { KtbDeletionDialogComponent } from '../_dialogs/ktb-deletion-dialog/ktb-deletion-dialog.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbDeletionComponent', () => {
  let component: KtbDangerZoneComponent;
  let fixture: ComponentFixture<KtbDangerZoneComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbDangerZoneComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should open the deletion dialog', () => {
    // given
    component.data = {name: 'sockshop', type: DeleteType.PROJECT};
    fixture.detectChanges();
    const button = fixture.nativeElement.querySelector('.danger-button');
    const spy = jest.spyOn(component.dialog, 'open');

    // when
    button.click();
    fixture.detectChanges();

    // then
    expect(spy).toHaveBeenCalled();
    expect(spy).toHaveBeenCalledWith(KtbDeletionDialogComponent, {data: {name: 'sockshop', type: 'project'}});
  });
});
