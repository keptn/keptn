import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbConfirmationDialogComponent } from './ktb-confirmation-dialog.component';
import { AppModule } from '../../../app.module';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { DeleteType } from '../../../_interfaces/delete';


describe('KtbDeletionDialogComponent', () => {
  let component: KtbConfirmationDialogComponent;
  let fixture: ComponentFixture<KtbConfirmationDialogComponent>;
  const dialogData = {
    sequence: {
      shkeptncontext: 'f6a38eb6-e99d-4d14-ab4c-3e94ed288b45',
      name: 'delivery'
    },
    confirmCallback: (params: any) => { }
  };
  const dialogRefMock = {
    close: () => { }
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KtbConfirmationDialogComponent],
      imports: [AppModule, MatDialogModule],
      providers: [
        {provide: MAT_DIALOG_DATA, useValue: dialogData},
        {provide: MatDialogRef, useValue: dialogRefMock}
      ]
    })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbConfirmationDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
