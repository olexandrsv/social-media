import { GetProjectResponse } from "../components/projects/projects.component";

export class Project{
    id: number = 0;
    login: string = "";
    name: string = "";
    description: string = "";
    stack: string = "";

    constructor(){

    }

    parseHTTPResponse(r: GetProjectResponse){
        this.id = r.id;
        this.login = r.login;
        this.name = r.name;
        this.description = r.description;
        this.stack = r.stack;
    }

    generateCreateForm(): FormData{
        const form = new FormData();
        form.append("name", this.name);
        form.append("description", this.description);
        form.append("stack", this.stack);
        return form;
    }

    generateUpdateForm(): FormData{
        const form = new FormData();
        form.append("name", this.name);
        form.append("description", this.description);
        form.append("stack", this.stack);
        return form;
    }
}