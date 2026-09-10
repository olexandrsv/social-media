import { HttpClient } from "@angular/common/http";
import { GetPostResponse } from "../components/posts/posts.component";
import { Message } from "./message";
import { Post } from "./post";

export class MessagesState {
    missedPosts = new Map<number, number> ()
    userID: number = -1;
    posts: Post[] = [];

    missedMessages = new Map<number, number>()
    chatID: number = -1;
    messages: Message[] = [];

    constructor(
    ){
        
    }

    postCreated(data: any){
        var post = new Post();
        post.parseHttpResponse(data as GetPostResponse);
        post.parseText();
        if (this.userID == post.userID){
            this.posts.push(post);
            return
        }
    }

    changeUser(userID: number, postsResponse: GetPostResponse[]){
        if (userID == this.userID) {
            return
        }
        var posts: Post[] = [];
        for (let i = 0; i<postsResponse.length; i++){
            var post = new Post();
            post.parseHttpResponse(postsResponse[i]);
            posts.push(post);
        }
        this.userID = userID;
        this.posts = posts;
    }
}