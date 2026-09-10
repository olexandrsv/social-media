import { HttpClient } from '@angular/common/http';
import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Project } from 'src/app/objects/project';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { RouterModule } from '@angular/router';
import { MatInputModule } from '@angular/material/input';


export interface ProjectData{
  project?: Project;
  title: string;
  buttonTitle: string;
  btnOnClick: (project: Project, dialogRef: MatDialogRef<ProjectComponent>) => void;
}

@Component({
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	  MatInputModule, RouterModule],
  selector: 'app-project',
  templateUrl: './project.component.html',
  styleUrls: ['./project.component.css']
})
export class ProjectComponent {
  project: Project = new Project();
  projectData: ProjectData;

  constructor(
    public dialogRef: MatDialogRef<ProjectComponent>,
    private http: HttpClient,
    @Inject(MAT_DIALOG_DATA) public data: ProjectData
  ){
    this.projectData = data;
    if (data.project) {
      this.project = data.project;
    } else {
      this.project = new Project();
    }
  }

  onClick(){
    this.projectData.btnOnClick(this.project, this.dialogRef);
  }
}
