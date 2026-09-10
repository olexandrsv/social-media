import { Injectable } from '@angular/core';
import { empty, firstValueFrom, Observable } from 'rxjs';
import { GetCommentsReq, GetPostResponse } from '../components/posts/posts.component';
import { HttpClient } from '@angular/common/http';
import { Post } from '../objects/post';
import { Comment } from '../objects/comment';
import { newUser, User } from '../objects/user';
import { Tone } from '../objects/tone';
import { Chat, GetBasicChatReposnse } from '../objects/chat';
import { Message } from '../objects/message';
import { GetMessageReposnse } from '../components/messages/messages.component';
import { MissedMessagesNumber, MissedPostsNumber } from './shared.service';

export type Result<T> = [data: T, error: null] | [data: null, error: Error];
export type BiResult<T, K> = [t: T, k: K, error: null] | [t: null, k: null, error: Error];
export type Err = Error | null;

@Injectable({
  providedIn: 'root'
})
export class ClientService {
  http: HttpClient

  constructor(
    http: HttpClient
  ) {
    this.http = http;
    }

  async catchError<T>(o: Observable<T>): Promise<Result<T>> {
    try {
      const response = await firstValueFrom(o);
      return [response, null];
    } catch (error) {
      return [null, error as Error];
    }
  }

  async getFollowing(): Promise<Result<User[]>>{
    const [response, err] = await this.catchError(
      this.http.get<User[]>(`http://localhost:8081/users/followed`, {withCredentials: true})
    );
    if (err){
      return [null, err]
    }

    const users = response.map(model => {
      const user = newUser(model.id, model.login);
      return user
    })

    return [users, null]
  }

  async getPosts(userID: number): Promise<Result<Post[]>> {
    const [response, err] = await this.catchError(
      this.http.get<GetPostResponse[]>(`http://localhost:8082/users/${userID}/posts`, { withCredentials: true })
    );
    if (err) {
      return [null, err];
    }

    const posts = response.map(model => {
      const post = new Post();
      post.parseHttpResponse(model);
      return post;
    });

    return [posts, null];
  }

  async subscribe(followingID: number): Promise<Err> {
    const postData = new FormData();
    const [response, err] = await this.catchError(
      this.http.post(`http://localhost:8081/users/${followingID}/followers`, postData, {withCredentials: true})
    )
    return err
  }

  async getPostComments(postID: string): Promise<BiResult<Comment[], Tone>>{
    return await this.getComments(`http://localhost:8082/posts/${postID}/comments`)
  }

  async getCommentComments(commentID: string): Promise<BiResult<Comment[], Tone>>{
    return await this.getComments(`http://localhost:8082/posts/comments/${commentID}/comments`)
  }

  private async getComments(url: string): Promise<BiResult<Comment[], Tone>>{
    const [response, err] = await this.catchError(
      this.http.get<GetCommentsReq>(url, {withCredentials: true})
    )
    if (err){
      return [null, null, err]
    }

    let comments: Comment[] = [];
    const tone = new Tone(response.tone.positive_percentage*100, response.tone.negative_percentage*100);
    for (let commentModel of response.comments){
      const comment = new Comment();
      comment.parseHttpResponse(commentModel);
      comment.parseText();
      comments.push(comment);
    }

    return [comments, tone, null]              
  }

  async deletePostComment(postID: string, commentID: string): Promise<Err>{
    const [response, err] = await this.catchError(
        this.http.delete(`http://localhost:8082/users/posts/${postID}/comments/${commentID}`, {withCredentials:true})
    )
    return err
  }

  async deleteCommentComment(parentCommentID: string, commentID: string): Promise<Err>{
    const [response, err] = await this.catchError(
        this.http.delete(
          `http://localhost:8082/users/posts/comments/${parentCommentID}/comments/${commentID}`, {withCredentials:true}
        )
    )
    return err
  }

  async getUsersByInfo(info: string): Promise<Result<User[]>>{
    return this.catchError(
      this.http.get<User[]>(`http://localhost:8081/users?info=${info}`, {withCredentials: true})
    )
  }

  async getChatsBasicInfo(): Promise<Result<Chat[]>>{
    const [response, err] = await this.catchError(
      this.http.get<GetBasicChatReposnse[]>('http://localhost:8083/chats', {withCredentials: true})
    )
    if (err){
      return [null, err]
    }
    const chats: Chat[] = [];
    for (let i=0; i < response.length; i++){
      const chat = new Chat();
      chat.parseBasicChatResponse(response[i]);
      chats.push(chat);
    }
    return [chats, null]
  }

  async getMessages(chatID: number): Promise<Result<Message[]>>{
    const [response, err] = await this.catchError(
      this.http.get<GetMessageReposnse[]>(`http://localhost:8082/chats/${chatID}/messages`, {withCredentials: true})
    )
    if (err){
      return [null, err]
    }
    
    const messages = response.map(messageModel => {
      const message = new Message();
      message.parseHttpResponse(messageModel);
      return message
    })
    return [messages, null]
  }

  async getMissedPostNumber(): Promise<Result<MissedPostsNumber[]>>{
    return await this.catchError(
      this.http.get<MissedPostsNumber[]>(`http://localhost:8082/posts/missed`, {withCredentials: true})
    )
  }

  async getMissedMessagesNumber(): Promise<Result<MissedMessagesNumber[]>>{
    return await this.catchError(
      this.http.get<MissedMessagesNumber[]>(`http://localhost:8082/chats/messages/missed`, {withCredentials: true})
    )
  }

  async markPostsAsRead(postOwnerID: number, lastReadPostID: string): Promise<Err>{
    const form = new FormData()
    form.append("last_read_post", lastReadPostID)
    const [response, err] = await this.catchError(
      this.http.put<any>(`http://localhost:8081/users/posts/${postOwnerID}/read`, form, {withCredentials: true})
    )
    return err
  }

  async markMessagesAsRead(chatID: number, lastReadMessageID: string): Promise<Err>{
    const form = new FormData()
    form.append("last_read_message", lastReadMessageID)
    const [response, err] = await this.catchError(
      this.http.put<any>(`http://localhost:8083/users/chats/${chatID}/read`, form, {withCredentials: true})
    )
    return err
  }
}
