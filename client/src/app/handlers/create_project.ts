import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { ProjectComponent, ProjectData } from "../components/project/project.component";
import { HttpClient } from "@angular/common/http";
import { Project } from "../objects/project";

export interface CreateProjectResponse {
    id: number;
}

export class CreateProject {
    projectData: ProjectData;

    constructor(public dialog: MatDialog, http: HttpClient, add: (project: Project) => void){
        this.projectData = {
            title: "Create project",
            buttonTitle: "Create",
            btnOnClick: this.create(http, add),
        }
        const dialogRef = this.dialog.open(ProjectComponent, {
            data: this.projectData,
            height: '600px',
            width: '600px',
        })
    }

    create(http: HttpClient, add: (project: Project) => void): (project: Project, dialogRef: MatDialogRef<ProjectComponent>)=> void{
        return (project: Project, dialogRef: MatDialogRef<ProjectComponent>) => {
            const form = project.generateCreateForm();
    
            http.post<CreateProjectResponse>(`http://localhost:8086/projects`, form, {withCredentials: true}).subscribe({
                next: result => {
                    console.log(result);
                    project.id = result.id;
                    add(project);
                    dialogRef.close();
                },
                error: err => {
                    
                }
                })
            }
        }
}