import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { WebsocketService } from './websocket.service';
import { HttpClient } from "@angular/common/http";

export interface MissedPostsNumber{
	following_id: number
	number: number
}

export interface MissedMessagesNumber{
	chat_id: number
	number: number
}

@Injectable({
  providedIn: 'root'
})
export class SharedService {
	webSocketService: WebsocketService

	badgeValueSource = new BehaviorSubject<number>(0);
	badgeValue$ = this.badgeValueSource.asObservable();
	
	msgBadgeValueSource = new BehaviorSubject<number>(0);
	msgBadgeValue = this.msgBadgeValueSource.asObservable();
	
	missed = new Map<number, number>();
	missedMsg = new Map<number, number>();
	
	roomId = -1;
	login = '';
	bol = true;
	
	updateBadgeValue(newValue: number) {
		let value = this.badgeValueSource.getValue();
		this.badgeValueSource.next(value - newValue);
	}
	
	addBadgeValue(newValue: number){
		let value = this.badgeValueSource.getValue();
		this.badgeValueSource.next(value + newValue);
	}
	
	updateMsgBadgeValue(newValue: number){
		let value = this.msgBadgeValueSource.getValue();
		this.msgBadgeValueSource.next(value - newValue);
	}
	
	addMsgBadgeValue(newValue: number){
		let value = this.msgBadgeValueSource.getValue();
		this.msgBadgeValueSource.next(value + newValue);
		console.log('addMsgBadgeValue', this.msgBadgeValueSource.getValue())
	}

	receivedChatMessages(chatID: number){
		const missedChatMsg = this.missedMsg.get(chatID);
		this.updateMsgBadgeValue(missedChatMsg!!);
		this.missedMsg.delete(chatID);
	}
	

	constructor(webSocketService: WebsocketService, private http: HttpClient) {
		this.webSocketService = webSocketService;
	}

	start(){
		this.http.get<MissedPostsNumber[]>(`http://localhost:8082/posts/missed`, {withCredentials: true}).subscribe({
			next: result => {
				console.log("missed posts result: ", result);
				let sum = 0;
				for (let i=0; i<result.length; i++) {
					var missedPostNumber = result[i]
					this.missed.set(missedPostNumber.following_id, missedPostNumber.number);
					sum += missedPostNumber.number;
				}
				this.badgeValueSource.next(sum);
			}
		});
		this.http.get<MissedMessagesNumber[]>(`http://localhost:8082/chats/messages/missed`, {withCredentials: true}).subscribe({
			next: result => {
				console.log("missed messages: ", result);
				let sum = 0;
				for (let i=0; i<result.length; i++) {
						var missedMessagesNumber = result[i]
						this.missedMsg.set(missedMessagesNumber.chat_id, missedMessagesNumber.number);
						sum += missedMessagesNumber.number;
				}
				console.log(sum);
				this.msgBadgeValueSource.next(sum);
			},
			error: err => {
				console.log(err);
			}
		});
		this.webSocketService.postReceived.subscribe(
	      (message: any) => {
			console.log('from shared service');
			console.log(message);
			if (message.login !== this.login){
				let value = this.missed.get(message.login);
				if (value !== undefined){
					this.missed.set(message.login, value+1);
				} else{
					this.missed.set(message.login, 1);
				}
				this.addBadgeValue(1);
			}
	      }
	    );
		this.webSocketService.messageReceived.subscribe(
	      (message: any) => {
			console.log('from shared service');
			console.log(message);
			if (message.roomId !== this.roomId){
				let value = this.missedMsg.get(message.roomId);
				if (value !== undefined){
					this.missedMsg.set(message.roomId, value+1);
				} else{
					this.missedMsg.set(message.roomId, 1);
				}
				this.addMsgBadgeValue(1);
			}
	      }
	    );
	}
}
