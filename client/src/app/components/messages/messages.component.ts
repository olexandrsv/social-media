import { Component, OnDestroy, Inject } from '@angular/core';
import { HttpClient } from "@angular/common/http";
import {WebsocketService} from '../../services/websocket.service';
import { SharedService } from '../../services/shared.service';
import {MatDialog, MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import { Chat, GetBasicChatReposnse } from 'src/app/objects/chat';
import { CreateChat, GetChatResponse } from 'src/app/handlers/create_chat';
import { UpdateChat } from 'src/app/handlers/update_chat';
import { ErrorWindow } from 'src/app/handlers/error';
import { CreateMessage } from 'src/app/handlers/create_message';
import { Message } from 'src/app/objects/message';
import { UpdateMessage } from 'src/app/handlers/update_message';
import { CookieService } from 'ngx-cookie-service';
import { ClientService } from 'src/app/services/client.service';
import { Missed, PostsStoreService } from 'src/app/services/posts-store.service';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FormsModule } from '@angular/forms';
import { MatBadgeModule } from '@angular/material/badge';
import { RouterModule } from '@angular/router';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';


export interface CreatePostRequest{
	post: PostModel
	message: MessageModel
}

export interface PostModel {
	id: number
	userID: number
	text: string
}

export interface MessageModel{
	text: string
}

export interface GetMessageReposnse{
	id: string
	chat_id: number
	user_id: number
	user_name: string
	user_surname: string
	text: string
	files: string[]
	images: string[]
}

@Component({
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	MatInputModule, RouterModule, MatBadgeModule],
  selector: 'app-messages',
  templateUrl: './messages.component.html',
  styleUrls: ['./messages.component.css']
})
export class MessagesComponent implements OnDestroy {
	cookieService: CookieService
	text: string = '';
	selectedImgs: File[] = [];
	selectedFiles: File[] = [];
	posts = [];
	messageList: any[] = [];
	rooms: any = [];
	roomId = -1;
	chats: Chat[] = [];
	choosedChat: Chat | undefined = undefined;
	message: Message = new Message();
	messages: Message[] = [];
	missed: Map<number, number>

	constructor(
		private http: HttpClient,
		private webSocketService: WebsocketService,
		public sharedService: SharedService,
		public postsStore: PostsStoreService,
		cookieService: CookieService,
		public client: ClientService,
		public dialog: MatDialog
	){
		this.cookieService = cookieService;
		this.missed = postsStore.messagesStore.missed.userMissedItemsNumber
		this.messageList = [];
		this.webSocketService.messageReceived.subscribe(
	      (message: any) => {
			console.log('new message received !!!');
	        this.newMessage(message);
	      }
	    );
		this.init()	
	}

	async init(){
		const [chats, err] = await this.client.getChatsBasicInfo()
		if (err){
			new ErrorWindow(err, this.dialog)
			return
		}
		this.chats = chats
	}

	updateMessage(message: Message){
		if (!this.choosedChat){
			return
		}
		new UpdateMessage(this.dialog, this.http, this.choosedChat, message, (m: Message)=>{
			m.parseText();
			this.postsStore.messagesStore.itemUpdated(message)
		})
	}

	deleteMessage(message: Message){
		this.http.delete(`http://localhost:8082/chats/messages/${message.id}`, {withCredentials: true}).subscribe({
			next: response => {
				
			},
			error: err => {
				new ErrorWindow(err, this.dialog);
			}
		})
	}

	createMessage(){
		if (!this.choosedChat){
			return
		}
		new CreateMessage(this.dialog, this.http, this.choosedChat, (m: Message)=>{
			m.parseText();
			this.messages.push(m);
		});
	}
	
	ngOnDestroy(){
		this.roomId = -1;
		this.sharedService.roomId = -1;
	}
	
	onFileSelected(event: any){
		this.selectedFiles = event.target.files;
	}
	
	onImgSelected(event: any){
		this.selectedImgs = event.target.files;
	}

	getChatMessages(chat: Chat){
		this.postsStore.getMessages(chat.id)
	}

	markReadMessages(){
		if (this.messages.length == 0){
			return
		}
		var form = new FormData();
		form.append("last_read_message", this.messages[this.messages.length-1].id)
		this.http.put<any>(`http://localhost:8081/users/chats/`+this.choosedChat?.id+`/read`, form, {withCredentials: true}).subscribe({
			error: error => {
				new ErrorWindow(error, this.dialog);
			}
		})
	}

	updateReadMessages(){
		const form = new FormData();
		const id = this.messages[this.messages.length-1].id
		console.log("message_id: "+id);
		form.append("message_id", id);

		this.http.put<any>(`http://localhost:8081/users/chats/${this.choosedChat!!.id}/read`, form, {withCredentials: true}).subscribe({
			next: response => {
				console.log("Here")
				console.log(response);
			},
			error: err => {
				console.log(err);
			}
		})
	}
	
	changeChat(chat: Chat){
		this.choosedChat = chat;
		this.postsStore.changeChat(chat.id)
		this.http.get<GetChatResponse>(`http://localhost:8083/chats/${chat.id}`, {withCredentials: true}).subscribe({
			next: response => {
				chat.parseHttpResponse(response);
				this.getChatMessages(chat);
			},
			error: err => {
				new ErrorWindow(err, this.dialog);
			}
		})

		// this.roomId = id;
		// this.sharedService.roomId = this.roomId;
		// this.http.get(`http://localhost:8080/msg/${id}`).subscribe((response: any) => {
		// 	console.log('change room', response);
		// 	this.messageList = response;
		// 	let value = this.sharedService.missedMsg.get(id);
		// 	if (value !== undefined){
		// 		this.sharedService.updateMsgBadgeValue(value)
		// 	}
		// 	this.missedMsg.delete(id);
		// });
	}
	
	newMessage(message: any){
		this.addMessage(message);
	}
	
	addMessage(msg: any){
		console.log('from add msg', msg);
		if(msg.roomId === this.roomId){
			this.messageList.push(msg);
		}    	
    }
    
    openDialog(){
		new CreateChat(this.dialog, this.http, "", (c: Chat)=>{
			this.chats.push(c)
		})
    }

	updateChat(){
		if (!this.choosedChat){
			return
		}
		let idx = -1;
		for (let i=0; i<this.chats.length; i++){
			if (this.chats[i].id == this.choosedChat.id){
				idx = i;
			}
		}
		new UpdateChat(this.dialog, this.http, this.choosedChat, (c: Chat)=>{
			this.chats[idx] = c;
		})
	}

	deleteChat(){
		if (!this.choosedChat){
			return
		}
		this.http.delete(`http://localhost:8083/chats/${this.choosedChat.id}`, {withCredentials: true}).subscribe({
			next: response => {
				this.chats = this.chats.filter(v => v.id !== this.choosedChat!!.id)
				this.choosedChat = undefined;
			},
			error: err => {
				new ErrorWindow(err, this.dialog);
			}
		})
	}
}

@Component({
	imports: [MatIconModule, MatFormFieldModule, FormsModule],
	selector: 'room',
	templateUrl: './room.html',
	styleUrls: ['./room.css']
})

export class NewRoomDialog {
	name: string = '';
	users: string = '';
	
	constructor(
	  private http: HttpClient,
	  public dialogRef: MatDialogRef<NewRoomDialog>,
	  @Inject(MAT_DIALOG_DATA) public data: any,
	) {
		
	}
	
	onNoClick(): void {
	  this.dialogRef.close();
	}
	
	createRoom(){
		const postData = new FormData();
		postData.append('name',  this.name);
		postData.append('users', this.users);
		
		this.http.post('http://localhost:8080/room', postData).subscribe((response: any) => {
			this.dialogRef.close(response);
		})
	}

}
