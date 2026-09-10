import { Component } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ProjectComponent } from '../project/project.component';
import { HttpClient } from '@angular/common/http';
import { Project } from 'src/app/objects/project';
import { CreateProject } from 'src/app/handlers/create_project';
import { UpdateProject } from 'src/app/handlers/update_project';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { RouterModule } from '@angular/router';


export interface GetProjectResponse{
  id: number,
  login: string,
  name: string,
  description: string,
  stack: string
}

@Component({
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	  MatInputModule, RouterModule],
  selector: 'app-projects',
  templateUrl: './projects.component.html',
  styleUrls: ['./projects.component.css']
})
export class ProjectsComponent {
  projects: Project[] = [];

  constructor(
    public dialog: MatDialog,
    private http: HttpClient
  ){
    this.http.get<GetProjectResponse[]>(`http://localhost:8086/projects`, {withCredentials: true}).subscribe({
      next: response => {
        for (let i=0; i<response.length; i++){
          let project = new Project();
          project.parseHTTPResponse(response[i]);
          this.projects.push(project);
        }
      },
      error: err => {
        console.log(err);
      }
    })
  }

  addProject(){
    new CreateProject(this.dialog, this.http, (project: Project) => {
      this.projects.push(project);
    });
      // const dialogRef = this.dialog.open(ProjectComponent, {
      //   height: '600px',
      //   width: '500px',
      // });
  }

  update(index: number){
    new UpdateProject(this.dialog, this.http, this.projects[index], (project: Project) => {
      this.projects[index] = project;
    });
  }

  delete(index: number){
    let project = this.projects[index];
    this.http.delete(`http://localhost:8086/projects/${project.id}`, {withCredentials: true}).subscribe({
      next: response => {
        this.projects = this.projects.filter((el, i) => {
          return i !== index;
        });
      },
      error: err => {
        
      }
    });
  }
}
