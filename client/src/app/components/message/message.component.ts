import { HttpClient } from '@angular/common/http';
import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Message } from 'src/app/objects/message';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FormsModule } from '@angular/forms';
import { MatBadgeModule } from '@angular/material/badge';
import { MatButtonModule } from '@angular/material/button';
import { RouterModule } from '@angular/router';
import { MatInputModule } from '@angular/material/input';


export interface MessageData{
  message?: Message
  title: string
  buttonTitle: string
  btnOnClick: (message: Message, dialogRef: MatDialogRef<MessageComponent>) => void;
}

@Component({
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	  MatInputModule, RouterModule, MatBadgeModule],
  selector: 'app-message',
  templateUrl: './message.component.html',
  styleUrls: ['./message.component.css']
})
export class MessageComponent {
  message: Message = new Message();
  messageData: MessageData;

  constructor(
    public dialogRef: MatDialogRef<MessageComponent>,
    private http: HttpClient,  
    @Inject(MAT_DIALOG_DATA) public data: MessageData
  ){
    this.messageData = data;
    if (data.message) {
      this.message = data.message;
    }
  }

  onImgSelected(event: any, input: HTMLInputElement){
    this.message.onImgSelected(event);
    input.value = '';
  }

  onFileSelected(event: any, input: HTMLInputElement){
    this.message.onFileSelected(event);
    input.value = '';
  }

  removeImage(i: number){
    this.message.onImgRemoved(i);
  }

  removeFile(i: number){
    this.message.onFileRemoved(i);
  }

  onClick(){
    this.messageData.btnOnClick(this.message, this.dialogRef);
  }
}
