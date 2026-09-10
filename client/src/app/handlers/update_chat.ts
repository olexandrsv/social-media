import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { ChatComponent, ChatData } from "../components/chat/chat.component";
import { HttpClient } from "@angular/common/http";
import { Chat } from "../objects/chat";
import { User } from "../objects/user";
import { GetChatResponse } from "./create_chat";
import { ErrorWindow } from "./error";

export class UpdateChat{
    chatData: ChatData

    constructor(
        public dialog: MatDialog, 
        http: HttpClient, 
        chat: Chat,
        updateChat: (c: Chat) => void)
    {
        this.chatData = {
            title: "Update chat",
            buttonTitle: "Update",
            btnOnClick: this.create(http, updateChat),
            chat: chat,
        }
        const dialogRef = this.dialog.open(ChatComponent, {
            data: this.chatData,
            height: '500px',
            width: '400px',
        })
    }

    create(http: HttpClient, updateChat: (c: Chat) => void): 
        (c: Chat, dialogRef: MatDialogRef<ChatComponent>)=> void 
    {
        return (c: Chat, dialogRef: MatDialogRef<ChatComponent>)=>{
            const form = c.generateCreateForm();
            
            http.put(
                `http://localhost:8083/chats/${c.id}`, 
                form, 
                {withCredentials: true}
            ).subscribe({
                next: result => {                                        
                    updateChat(c);
                    dialogRef.close();
                },
                error: err => {
                    new ErrorWindow(err, this.dialog);
                }
            })
        }
    }    
}