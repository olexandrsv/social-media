import { HttpClient } from '@angular/common/http';
import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Chat } from 'src/app/objects/chat';
import { User } from 'src/app/objects/user';
import { ErrorComponent } from '../error/error.component';
import { ErrorWindow } from 'src/app/handlers/error';
import { ClientService } from 'src/app/services/client.service';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatBadgeModule } from '@angular/material/badge';
import { MatSelectModule } from '@angular/material/select';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';

export interface ChatData{
  chat?: Chat
  title: string
  buttonTitle: string
  btnOnClick: (c: Chat, dialogRef: MatDialogRef<ChatComponent>) => void;
}

@Component({
  standalone: true,
  selector: 'app-chat',
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	  MatInputModule, MatSelectModule, MatAutocompleteModule, MatBadgeModule,
	  ReactiveFormsModule],
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.css']
})
export class ChatComponent {
  chatData: ChatData;
  text: string = "";
  chat: Chat = new Chat();
  options: User[] = [];

  constructor(
    public dialogRef: MatDialogRef<ChatComponent>,
    private http: HttpClient,  
    public client: ClientService,
    @Inject(MAT_DIALOG_DATA) public data: ChatData,
    public dialog: MatDialog
  ){
    this.chatData = data;
    if (data.chat) {
      this.chat = data.chat;
    } else {
      const user = new User();
      const id = localStorage.getItem("id");
      user.id = parseInt(id!!, 10);
      this.chat.owner = user;
    }
    console.log(this.chat);
  }

  addUser(){
    const user = this.options.find(v => v.login === this.text)
    if (user){
      this.chat.addUser(user);
    }
    this.text = "";
    console.log(this.chat.users);
  }

  deleteUser(i: number){
    this.chat.deleteUser(i)
  }

  onClick(){
    this.chatData.btnOnClick(this.chat, this.dialogRef);
  }

  async onButtonPressed(event: KeyboardEvent){
    if (event.key === "Enter"){
      this.onClick()
      return
    }
    const [response, err] = await this.client.getUsersByInfo(this.text)
    if (err){
      new ErrorWindow(err, this.dialog)
      return
    }
    const id = localStorage.getItem("id")
    this.options = response.filter(u => u.id.toString() !== id);
  }

  onOptionClicked(login: string){
		this.text = login;
		this.onOptionChoosed();
	}

  onOptionChoosed(){
    let id = -1;
		for (let i=0; i<this.options.length; i++){
			if (this.options[i].login == this.text){
				id = this.options[i].id;
			}
		}
		if (id === -1){
			return
		}

    this.addUser();
  }
}
