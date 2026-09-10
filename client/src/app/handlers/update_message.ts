import { MatDialog, MatDialogRef } from "@angular/material/dialog";
import { MessageComponent, MessageData } from "../components/message/message.component";
import { HttpClient } from "@angular/common/http";
import { Message } from "../objects/message";
import { GetMessageReposnse } from "../components/messages/messages.component";
import { Chat } from "../objects/chat";

export class UpdateMessage{
    messageData: MessageData;
    dialogRef: MatDialogRef<MessageComponent, any>;
    
    constructor(public dialog: MatDialog, http: HttpClient, chat: Chat, m: Message, updateMessage: (m: Message)=>void){
        this.messageData = {
            message: m,
            title: "Update message",
            buttonTitle: "Update",
            btnOnClick: this.update(http, chat, updateMessage),
        }
        this.dialogRef = this.dialog.open(MessageComponent, {
            data: this.messageData,
            height: '400px',
            width: '600px',
        })
    }

    update(http: HttpClient, chat: Chat, updateMessage: (m: Message)=>void): (m:Message, dialogRef: MatDialogRef<MessageComponent>) => void {
        return (m:Message, dialogRef: MatDialogRef<MessageComponent>) => {
            m.chatID = chat.id
            const form = m.generateUpdateForm();
            console.log(form)
            console.log("message id: "+m.id);

            http.put<GetMessageReposnse>(`http://localhost:8082/chats/messages/`+m.id, form, {withCredentials: true}).subscribe({
                next: result => {
                    console.log(result);
                    m.parseHttpResponse(result);
                    updateMessage(m);
                    dialogRef.close();
                },
                error: err => {
                    console.log(err);
                }
            })
        }
    }
}