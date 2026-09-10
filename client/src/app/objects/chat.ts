import { GetChatResponse } from "../handlers/create_chat";
import { User } from "./user";

export interface GetBasicChatReposnse{
	id: number
	name: string
}

export class Chat{
    id: number = 0;
    name: string = "";
    owner: User = new User();
    users: User[] = [];
    
    constructor(){

    }

    deleteUser(i: number){
        this.users.splice(i, 1)
    }

    addUser(user: User){
        this.users.push(user);
    }

    parseBasicChatResponse(r: GetBasicChatReposnse){
        this.id = r.id;
        this.name = r.name;
    }

    parseHttpResponse(r: GetChatResponse){
        const rawID = localStorage.getItem("id")
        if (!rawID){
            return
        }
        const userID = parseInt(rawID, 10)
        this.id = r.id;
        this.name = r.name;
        this.owner = r.owner;
        this.users = r.users.filter(v => v.id !== userID);
    }

    generateCreateForm(): FormData {
        const form = new FormData();
        form.append("name", this.name);
        for (let i=0; i<this.users.length; i++){
            form.append("users_ids[]", this.users[i].id.toString());
        }
        return form
    }

    generateUpdateForm(): FormData {
        return this.generateUpdateForm();
    }

    parseBasicResponse(r: GetBasicChatReposnse){
        this.id = r.id;
        this.name = r.name;
    }
}