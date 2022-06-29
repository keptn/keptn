import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { StoreModule } from '@ngrx/store';
import { rootStoreReducer } from './root.store.reducer';
import { EffectsModule } from '@ngrx/effects';
import { RootStoreEffects } from './root.store.effects';
import { Features } from '../store';

@NgModule({
  imports: [
    CommonModule,
    StoreModule.forRoot({ [Features.ROOT]: rootStoreReducer }),
    EffectsModule.forRoot([RootStoreEffects]),
  ],
})
export class RootStoreModule {}
