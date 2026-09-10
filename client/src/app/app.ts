import { Component, signal } from '@angular/core';
import { Router, RouterModule, RouterOutlet } from '@angular/router';
import { delay, Observable, of } from 'rxjs';
import { Missed, PostsStoreService } from './services/posts-store.service';
import { WebsocketService } from './services/websocket.service';
import { CookieService } from 'ngx-cookie-service';
import { SharedService } from './services/shared.service';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatBadgeModule } from '@angular/material/badge';
import { MatSelectModule } from '@angular/material/select';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	  MatInputModule, MatBadgeModule, RouterModule],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected readonly title = signal('client');

  count = 0;
	badgeValue$: Observable<number> = of(this.count).pipe(delay(1000));
	missedPosts: Missed<number>;
	missedMessages: Missed<number>
	
	constructor(
	  private router: Router, 
	  webSocketService: WebsocketService, 
	  cookieService: CookieService,
	  public postsStore: PostsStoreService,
	  public sharedService: SharedService) {
		this.missedPosts = postsStore.postsStore.missed;
		this.missedMessages = postsStore.messagesStore.missed
		if (cookieService.check('token')) {
			webSocketService.start();
			//sharedService.start();
		} else {
			this.sharedService.bol = false;
			router.navigate(['/login']);
		}
		console.log('from main app');
	}
	
	navigateTo(path: string){
		this.router.navigate([path]);
	}
	
}
