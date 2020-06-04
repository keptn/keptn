import { Project } from './project';
import {async} from "@angular/core/testing";
import {Trace} from "./trace";

describe('Project', () => {
  it('should create an instance', () => {
    expect(new Project()).toBeTruthy();
  });

  it('should create instances from json', async(() => {
    let projects: Project[] = [{"projectName":"sockshop","stages":[{"services":[{"serviceName":"carts"},{"serviceName":"carts-db"}],"stageName":"dev"},{"services":[{"serviceName":"carts"},{"serviceName":"carts-db"}],"stageName":"staging"},{"services":[{"serviceName":"carts"},{"serviceName":"carts-db"}],"stageName":"production"}]}].map(project => Project.fromJSON(project));

    expect(projects[0] instanceof Project).toBe(true, 'instance of Project');

    expect(projects[0].projectName).toBe('sockshop');
    expect(projects[0].getServices().length).toBe(2, 'Project "sockshop" should have 2 services');
    expect(projects[0].stages.length).toBe(3, 'Project "sockshop" should have 3 stages');
  }));
});
