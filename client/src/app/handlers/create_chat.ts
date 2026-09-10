import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { ChatComponent, ChatData } from "../components/chat/chat.component";
import { HttpClient } from "@angular/common/http";
import { Chat, GetBasicChatReposnse } from "../objects/chat";
import { User } from "../objects/user";

export interface GetChatResponse{
    id: number
    name: string
    owner: User
    users: User[]
}

export class CreateChat{
    chatData: ChatData

    constructor(
        public dialog: MatDialog, 
        http: HttpClient, 
        parentID: string,
        addChat: (c: Chat) => void)
    {
        this.chatData = {
            title: "Create chat",
            buttonTitle: "Create",
            btnOnClick: this.create(http, parentID, addChat),
        }
        const dialogRef = this.dialog.open(ChatComponent, {
            data: this.chatData,
            height: '500px',
            width: '400px',
        })
    }

    create(http: HttpClient, parentID: string, addChat: (c: Chat) => void): 
        (c: Chat, dialogRef: MatDialogRef<ChatComponent>)=> void 
    {
        return (c: Chat, dialogRef: MatDialogRef<ChatComponent>)=>{
            const form = c.generateCreateForm();
            
            http.post<GetChatResponse>(
                `http://localhost:8083/chats`, 
                form, 
                {withCredentials: true}
            ).subscribe({
                next: result => {
                    const chat = new Chat();
                    chat.parseHttpResponse(result)
                    
                    addChat(chat);
                    dialogRef.close();
                },
                error: err => {
    
                }
            })
        }
    }    
}