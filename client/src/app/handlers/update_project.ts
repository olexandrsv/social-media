import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { ProjectComponent, ProjectData } from "../components/project/project.component";
import { HttpClient } from "@angular/common/http";
import { Project } from "../objects/project";
import { GetPostResponse } from "../components/posts/posts.component";

export class UpdateProject {
    projectData: ProjectData;

    constructor(public dialog: MatDialog, http: HttpClient, project: Project, updateProject: (p: Project) => void) {
        this.projectData = {
            project: project,
            title: "Update project",
            buttonTitle: "Update",
            btnOnClick: this.update(http, updateProject),
        }
        const dialogRef = this.dialog.open(ProjectComponent, {
            data: this.projectData,
            height: '600px',
            width: '600px',
        })
    }

    update(http: HttpClient, updateProject: (p: Project)=> void): (project: Project, dialogRef: MatDialogRef<ProjectComponent>)=> void{
        return (project: Project, dialogRef: MatDialogRef<ProjectComponent>) => {
            const form = project.generateUpdateForm();
    
            http.put(`http://localhost:8086/projects/`+project.id, form, {withCredentials: true}).subscribe({
                next: result => {
                    console.log(result);
                    dialogRef.close();
                    updateProject(project);
                },
                error: err => {
    
                }
                })
            }
        }
}