import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { ChatComponent, ChatData } from "../components/chat/chat.component";
import { HttpClient } from "@angular/common/http";
import { Chat, GetBasicChatReposnse } from "../objects/chat";
import { User } from "../objects/user";
import { MessageComponent, MessageData } from "../components/message/message.component";
import { Message } from "../objects/message";
import { GetChatResponse } from "./create_chat";
import { GetMessageReposnse } from "../components/messages/messages.component";


export class CreateMessage{
    messageData: MessageData

    constructor(
        public dialog: MatDialog, 
        http: HttpClient,
        chat: Chat,
        addMessage: (msg: Message) => void)
    {
        this.messageData = {
            title: "Create message",
            buttonTitle: "Create",
            btnOnClick: this.create(http, chat, addMessage),
        }
        const dialogRef = this.dialog.open(MessageComponent, {
            data: this.messageData,
            height: '400px',
            width: '600px',
        })
    }

    create(http: HttpClient, chat: Chat, addMessage: (msg: Message) => void): 
        (msg: Message, dialogRef: MatDialogRef<MessageComponent>)=> void 
    {
        return (msg: Message, dialogRef: MatDialogRef<MessageComponent>)=>{
            msg.chatID = chat.id
            const form = msg.generateCreateForm();
            
            http.post<GetMessageReposnse>(
                `http://localhost:8082/chats/${chat.id}/messages`, 
                form, 
                {withCredentials: true}
            ).subscribe({
                next: result => {
                    let message = new Message();
                    message.parseHttpResponse(result);
                    addMessage(message);
                    dialogRef.close();
                },
                error: err => {
                    console.log(err);
                }
            })
        }
    }    
}