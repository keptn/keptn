import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute} from '@angular/router';
import {Observable} from 'rxjs';
import {UniformRegistration} from "../../_models/uniform-registration";

@Component({
  selector: 'ktb-uniform-view',
  templateUrl: './ktb-uniform-view.component.html',
  styleUrls: ['./ktb-uniform-view.component.scss']
})
export class KtbUniformViewComponent implements OnInit {
  public selectedService: UniformRegistration;
  public uniformRegistrations$: Observable<UniformRegistration[]>;

  constructor(private dataService: DataService, private route: ActivatedRoute, private _changeDetectorRef: ChangeDetectorRef) { }

  ngOnInit(): void {
    this.uniformRegistrations$ = this.dataService.getUniformRegistrations();
  }

  selectService(service: UniformRegistration) {
    this.selectedService = service;
  }
}
