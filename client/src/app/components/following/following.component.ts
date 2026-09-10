import { Component, ViewChild, EventEmitter, Output, OnDestroy } from '@angular/core';
import { HttpClient } from "@angular/common/http";
import {FormControl, FormsModule, ReactiveFormsModule} from '@angular/forms';
import {Observable} from 'rxjs';
import {map, startWith} from 'rxjs/operators';
import {MatDialog, MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';
import {CommentsWindow, GetPostResponse} from '../posts/posts.component';
import {WebsocketService} from '../../services/websocket.service';
import { SharedService } from '../../services/shared.service';
import { User } from '../find/find.component';
import { Post } from 'src/app/objects/post';
import { ErrorWindow } from 'src/app/handlers/error';
import { PostsStoreService } from 'src/app/services/posts-store.service';
import { ClientService } from 'src/app/services/client.service';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatBadgeModule } from '@angular/material/badge';
import { MatSelectModule } from '@angular/material/select';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';


@Component({
  standalone: true,
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	MatInputModule, MatSelectModule, MatAutocompleteModule, MatBadgeModule,
	ReactiveFormsModule],
  selector: 'app-following',
  templateUrl: './following.component.html',
  styleUrls: ['./following.component.css']
})

export class FollowingComponent implements OnDestroy {
	postList: Post[] = [];
	text: string = '';
	following: User[] = [];
  	filteredFollowing: User[] = [];
	input = new FormControl('');
	missed: Map<number, number>

	choosedUser: User | null = null;
	client: ClientService;
	dialog: MatDialog;
	
	constructor(
	  private http: HttpClient, 
	  dialog: MatDialog, 
	  public webSocketService: WebsocketService, 
	  public postsStore: PostsStoreService,
	  client: ClientService,
	  private sharedService: SharedService){
		this.client = client
		this.dialog = dialog
		 this.missed = this.postsStore.missedPosts.userMissedItemsNumber
		this.init()
	}

	async init(){
		const [users, err] = await this.postsStore.getFollowing()
		if (err){
			new ErrorWindow(err, this.dialog)
			return
		}
		console.log("filteredFollowing: ", users)

		this.following = users;
		this.filteredFollowing = users;

		this.input.valueChanges.subscribe(newValue => {
			this.filteredFollowing = this.filter(newValue || '')
		})
	}
	
	ngOnDestroy(){
		//this.login = '';
		//this.sharedService.login = '';
	}
	
	private filter(value: string): User[] {
    	const filterValue = value.toLowerCase();
    	return this.following.filter(option => option.login.toLowerCase().includes(filterValue));
    }
	
	onUserChoosed(user: User){
		console.log("onUserChoosed")
		this.postsStore.changeUser(user.id);
	}
	
	onButtonPressed(event: KeyboardEvent){
    	if(event.key == "Enter"){
	    	const users = this.filter(this.text)
			if (users.length == 0){
				return
			}
			this.onUserChoosed(users[0]);
		}
    }
    
    openDialog(post: Post){
		console.log(post.id)
	    const dialogRef = this.dialog.open(CommentsWindow, {
	        data: {postId: post.id},
			height: '600px',
  			width: '700px',
	    });
	
	    dialogRef.afterClosed().subscribe(result => {
	      console.log('result');
	    });
    }

}
