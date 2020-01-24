import {Injectable} from '@angular/core';
import {BehaviorSubject, from, Observable, timer} from "rxjs";
import {debounce, map, mergeMap, toArray} from "rxjs/operators";

import {Root} from "../_models/root";
import {Trace} from "../_models/trace";
import {Stage} from "../_models/stage";
import {Project} from "../_models/project";
import {Service} from "../_models/service";

import {ApiService} from "./api.service";

@Injectable({
  providedIn: 'root'
})
export class DataService {

  private _projects = new BehaviorSubject<Project[]>([]);
  private _projectsLastUpdated: Date;

  constructor(private apiService: ApiService) {
    this.loadProjects();
  }

  get projects(): Observable<Project[]> {
    return this._projects.asObservable();
  }

  public loadProjects() {
    this.apiService.getProjects()
      .pipe(
        debounce(() => timer(10000)),
        map((projects) => projects.filter(project => project.projectName != 'lost+found')), // TODO: API_FIX: don't provide lost+found project in result, see https://github.com/keptn/keptn/issues/1275
        mergeMap((projects) =>
          from(projects).pipe(
            mergeMap(
              project => this.apiService.getStages(project.projectName)
                .pipe(
                  mergeMap((stages) =>
                    from(stages).pipe(
                      mergeMap(
                        stage => this.apiService.getServices(project.projectName, stage.stageName)
                          .pipe(
                            map(services => services.map(service => Service.fromJSON(service))),
                            map(services => ({ ...stage, services}))
                          )
                      ),
                      toArray()
                    )
                  ),
                  map(stages => stages.map(stage => Stage.fromJSON(stage))),
                  map(stages => ({ ...project, stages}))
                )
            ),
            toArray()
          )
        ),
        map(projects => projects.map(project => Project.fromJSON(project)))
      ).subscribe((projects: Project[]) => {
        this._projects.next([...this._projects.getValue(), ...projects]);
      }, (err) => {
        this._projects.error(err);
      });
  }

  public loadRoots(project: Project, service: Service) {
    this.apiService.getRoots(project.projectName, service.serviceName)
      .pipe(
        debounce(() => timer(10000)),
        mergeMap((roots) =>
          from(roots).pipe(
            mergeMap(
              root => this.apiService.getTraces(root.shkeptncontext)
                .pipe(
                  map(traces => traces.map(trace => Trace.fromJSON(trace))),
                  map(traces => ({ ...root, traces}))
                )
            ),
            toArray()
          )
        ),
        map(roots => roots.map(root => Root.fromJSON(root)))
      )
      .subscribe((roots: Root[]) => {
        // TODO: investigate why is the sorting changed?
        service.roots = roots.sort((a, b) => new Date(b.time).getTime() - new Date(a.time).getTime());;
        // TODO: return Subject with proper value handling
        // this._projects.next([...this._projects.getValue(), ...projects]);
      }, (err) => {
        // TODO: return Subject with proper error handling
        // this._projects.error(err);
      });
  }

  public loadTraces(root: Root) {
    this.apiService.getTraces(root.shkeptncontext)
      .pipe(map(traces => traces.map(trace => Trace.fromJSON(trace))))
      .subscribe((traces: Trace[]) => {
        root.traces = traces;
      });
  }
}
