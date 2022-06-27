import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KtbUniformRegistrationLogsComponent } from './ktb-uniform-registration-logs.component';
import { RouterModule } from '@angular/router';

@NgModule({
  declarations: [KtbUniformRegistrationLogsComponent],
  imports: [CommonModule, RouterModule],
  exports: [KtbUniformRegistrationLogsComponent],
})
export class KtbUniformRegistrationLogsModule {}
