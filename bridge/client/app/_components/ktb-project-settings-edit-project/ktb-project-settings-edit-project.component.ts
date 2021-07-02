import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'ktb-project-settings-edit-project',
  templateUrl: './ktb-project-settings-edit-project.component.html',
  styleUrls: ['./ktb-project-settings-edit-project.component.scss']
})
export class KtbProjectSettingsEditProjectComponent implements OnInit {

  constructor() { }

  ngOnInit(): void {
  }

  public updateFile(files: any) {
    console.log(files);
  }
}
