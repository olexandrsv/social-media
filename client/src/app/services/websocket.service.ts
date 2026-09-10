import { Injectable, EventEmitter  } from '@angular/core';
import { webSocket } from 'rxjs/webSocket';
import {map} from 'rxjs/operators';
import { GetPostResponse } from '../components/posts/posts.component';
import { FollowingComponent } from '../components/following/following.component';
import { Post } from '../objects/post';
import { SharedService } from './shared.service';
import { MessagesState } from '../objects/global';
import { PostsStoreService } from './posts-store.service';

export interface LastReadMessage {
	type: string
	chat_id: number
	last_read_message_id: string
}

export interface WebsocketMessage {
	type: string
	data: any
}

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
	postReceived = new EventEmitter<any>();
	messageReceived = new EventEmitter<any>();
	socket: WebSocket | null = null;
	postsStore: PostsStoreService

    constructor(postsStore: PostsStoreService) {
		this.postsStore = postsStore;
	}

	start(){
		this.socket = new WebSocket('http://localhost:8081/ws');
		this.socket.onopen = () => {
			//this.send("Hello, World");
			this.sendLastReadMessageID("123", 1);
		}
		this.socket.onmessage = (ev: MessageEvent) => {
			
			console.log("ev data: ", ev.data);
			var data = JSON.parse(ev.data) as WebsocketMessage
			console.log(data.type, data.data);
			this.callFunction(data.type, data.data);
		}
		this.socket.onerror = (ev: Event) => {
			console.error("WebSocket error observed:", ev);
		}
		
	
		console.log('websocket service started');
	}

	callFunction(name: string, data: any){
		switch (name) {
			case "postCreated":
				console.log("postCreated")
				this.postsStore.postCreated(data)
				break
			case "postUpdated":
				console.log("postUpdated")
				this.postsStore.postUpdated(data)
				break
			case "postDeleted":
				console.log("postDeleted")
				this.postsStore.postDeleted(data)
				break
			case "messageCreated":
				console.log("message created")
				this.postsStore.messageCreated(data)
				break
			case "messageUpdated":
				console.log("message updated")
				this.postsStore.updateMessage(data)
				break
			case "messageDeleted":
				console.log("message deleted")
				this.postsStore.deleteMessage(data)
				break
		}
	}

	

	messageCreated(data: any){

	}

	sendLastReadMessageID(id: string, chatID: number){
		const data: LastReadMessage = {
			type: "message",
			chat_id: chatID,
			last_read_message_id: id,
		}
		this.sendJSON(data)
	}

	send(data: string){
		if (!this.socket){
			return
		}
		this.socket.send(data);
	}

	sendJSON(data: any){
		this.send(JSON.stringify(data));
	}
}
