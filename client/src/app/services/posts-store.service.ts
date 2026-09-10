import { Injectable } from '@angular/core';
import { Post } from '../objects/post';
import { Message } from '../objects/message';
import { GetPostResponse } from '../components/posts/posts.component';
import { NotFoundError } from '../objects/errors';
import { Router } from '@angular/router';
import { ClientService, Err, Result } from './client.service';
import { User } from '../objects/user';
import { GetMessageReposnse } from '../components/messages/messages.component';

@Injectable({
  providedIn: 'root'
})
export class PostsStoreService {
    router: Router;
    client: ClientService;
    totalMissedMessagesNumber: number = 0;
  
    addedPosts = new Map<number, Post[]>();

    currentUserID: number = -1;
    missedPosts: Missed<number>;
    isPostsLoaded: boolean = false
    posts: Post[] = [];

    postsStore: Store<number, Post, string>
    messagesStore: Store<number, Message, string>

    following: User[] = [];


    constructor(
      router: Router,
      client: ClientService
    ) {
      this.router = router;
      this.client = client;
      this.missedPosts = new Missed()

      this.messagesStore = new Store(router, "/messages",
        (userID: number, postID: string) => {return this.client.markMessagesAsRead(userID, postID) },
        () => { return this.getMissedMessagesNumber()},
        (id: number) => {return client.getMessages(id)}, 
        (t: Message) => { return t.chatID }
      )
      
      this.postsStore = new Store(router, "/following", 
        (userID: number, postID: string) => {return this.client.markPostsAsRead(userID, postID) },
        () => { return this.getMissedPostsNumber()},
        (id: number) => {return client.getPosts(id)}, 
        (t: Post)=> { return t.userID }
      )

      this.addedPosts = this.postsStore.itemsStore
      this.missedPosts = this.postsStore.missed
      this.posts = this.postsStore.currentItems      
    }

    async getMissedPostsNumber(): Promise<Result<MissedNumber<number>[]>>{
      const [response, err] = await this.client.getMissedPostNumber()
      if (err){
        return [null, err]
      }
      const missedNumber = response.map(item => ({
        id: item.following_id,
        number: item.number
      }))
      return [missedNumber, null]
    }

    async getMissedMessagesNumber(): Promise<Result<MissedNumber<number>[]>>{
      const [response, err] = await this.client.getMissedMessagesNumber()
      if (err){
        return [null, err]
      }
      const missedNumber = response.map(item => ({
        id: item.chat_id,
        number: item.number
      }))
      return [missedNumber, null]
    } 

    async subscribe(following: User): Promise<Err>{
      let err = await this.client.subscribe(following.id)
      if (err){
        return err
      }
      this.following.push(following)
      return null
    }
    
    async getFollowing(): Promise<Result<User[]>>{
      if (this.following.length != 0){
        return [this.following, null]
      }
      const [users, err] = await this.client.getFollowing()
      if (err) {
        return [null, err]
      }
      this.following = users;
      return  [users, null]
    }

    async getPosts(userID: number): Promise<Err>{
      return this.postsStore.getItems(userID)
    }

    async changeUser(userID: number): Promise<Err>{
      return this.postsStore.changeKey(userID)
    }

    postCreated(data: GetPostResponse){
      var post = new Post();
      post.parseHttpResponse(data);
      return this.postsStore.itemCreated(post)
    }

    postUpdated(data: GetPostResponse){
      var post = new Post()
      post.parseHttpResponse(data)
      this.postsStore.itemUpdated(post)
    }

    postDeleted(data: GetPostResponse){
      var post = new Post()
      post.parseHttpResponse(data)
      this.postsStore.itemDeleted(post)
    }

    messageCreated(data: GetMessageReposnse){
      var message = new Message()
      message.parseHttpResponse(data)

      this.messagesStore.itemCreated(message)
      console.log("messages: ", this.messagesStore.currentItems)
    }

    async getMessages(chatID: number): Promise<Err>{
      return this.messagesStore.getItems(chatID)
    }

    async changeChat(chatID: number): Promise<Err>{
      return this.messagesStore.changeKey(chatID)
    }

    updateMessage(data: GetMessageReposnse){
      var message = new Message()
      message.parseHttpResponse(data)
      this.messagesStore.itemUpdated(message)
    }

    async deleteMessage(data: GetMessageReposnse){
      var message = new Message()
      message.parseHttpResponse(data)
      this.messagesStore.itemDeleted(message)
    }

    getLastPost() {
      const l = this.posts.length;
      if (l === 0){
        throw new NotFoundError();
      }
      return this.posts[l-1]
    }

    getPostID(i: number): string {
      return this.posts[i].id
    }
}

export interface Storable{
  id: any
  userID: any
}

export interface MissedNumber<K>{
  id: K
  number: number
}

export class Store<KeyType, ItemType extends Storable, ItemIdType>{
  itemsStore = new Map<KeyType, ItemType[]>()
  loadedMark = new Map<KeyType, boolean>()
  markAsRead: (key: KeyType, id: ItemIdType) => Promise<Err>
  getItemsMissedNumber: () => Promise<Result<MissedNumber<KeyType>[]>>
  getKey: (t: ItemType) => KeyType
  loadItems: (key: KeyType) => Promise<Result<ItemType[]>>
  currentKey: KeyType | null = null
  currentItems: ItemType[] = []
  missed: Missed<KeyType>
  router: Router
  url: string

  constructor(
    router: Router, 
    url: string,
    markAsRead: (key: KeyType, id: ItemIdType) => Promise<Err>,
    getItemsMissedNumber: () => Promise<Result<MissedNumber<KeyType>[]>>,
    getItems: (key: KeyType) => Promise<Result<ItemType[]>>, 
    getKey: (t: ItemType) => KeyType
  ){
    this.router = router
    this.url = url
    this.missed = new Missed<KeyType>()
    this.markAsRead = markAsRead
    this.getKey = getKey
    this.loadItems = getItems
    this.getItemsMissedNumber = getItemsMissedNumber
    getItemsMissedNumber()

    this.init()
  }

  async init(){
    const [response, err] = await this.getItemsMissedNumber()
    if (err){
      console.log(err)
      return
    }

    for (const item of response){
      this.itemsMissed(item.id, item.number)
    }
  }

  async getItems(key: KeyType): Promise<Err>{
    if (this.loadedMark.get(key)){
      const items = this.itemsStore.get(key)!
      if (items === undefined){
        return new Error("posts are empty")
      }
      this.currentItems = items
      return null
    }

    const [items, err] = await this.loadItems(key);
    if (err){
      return err
    }

    this.itemsStore.set(key, items)
    this.currentItems = items
    this.loadedMark.set(key, true)
    return null
  }

  async changeKey(key: KeyType): Promise<Err>{
    let err = await this.getItems(key)
    if (err){
      return err
    }
    this.missed.itemsRead(key);
    this.currentKey = key

    const len = this.currentItems.length
    if (len == 0){
      return null
    }
    err = await this.markAsRead(key, this.currentItems[len-1].id)
    if (err){
      return err
    }

    return null
  }

  itemsMissed(key: KeyType, count: number){
    for (let i = 0; i < count; i++){
      this.missed.itemCreated(key)
    }
  }

  itemCreated(item: ItemType){    
    const key = this.getKey(item)
    const saved = this.itemsStore.get(key)
    if (saved){
      saved.push(item)
    } else {
      const items = [item];
      this.itemsStore.set(key, items);
    }
    if (localStorage.getItem("id") === item.userID.toString()){
      return
    }
    if (this.postsVisible(key)){
      this.markAsRead(key, item.id)
      return
    }

    this.missed.itemCreated(key);
  }

  itemUpdated(item: ItemType){
    const key = this.getKey(item)
    const saved = this.itemsStore.get(key)
    if (!saved) {
      this.itemsStore.set(key, [item])
      return
    }
    const idx = saved.findIndex(t => t.id === item.id)
    saved[idx] = item
  }

  itemDeleted(item: ItemType){
    const key = this.getKey(item)
    const saved = this.itemsStore.get(key)
    if (!saved){
      return
    }
    const idx = saved.findIndex(t => t.id === item.id)
    saved.splice(idx, 1)
  }

  postsVisible(key: KeyType){
    console.log("postsVisible()", this.router.url)
    return key == this.currentKey && this.router.url.endsWith(this.url)
  }
}

export class Missed<K>{
  totalMissedItems: number = 0;

  userMissedItemsNumber = new Map<K, number>();
  key: K | null = null;

  constructor(){
  }

  changeKey(key: K){
    this.key = key;
  }

  itemCreated(key: K){
      this.totalMissedItems++
      const missedNumber = this.userMissedItemsNumber.get(key)
      if (missedNumber === undefined){
        this.userMissedItemsNumber.set(key, 1)
      } else {
        this.userMissedItemsNumber.set(key, missedNumber+1)
      }
  }

  itemsRead(key: K){
    this.changeKey(key);
    const readPostsNumber = this.userMissedItemsNumber.get(key)

    console.log("readPostsNumber", readPostsNumber)
    console.log("totalMissedPosts", this.totalMissedItems)

    if (readPostsNumber === undefined){
      return
    }
    this.userMissedItemsNumber.set(key, 0);
    this.totalMissedItems -= readPostsNumber;
  }
}
