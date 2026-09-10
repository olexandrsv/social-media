import { Component } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Router, RouterModule } from '@angular/router';
import { CookieService } from 'ngx-cookie-service';
import { WebsocketService } from '../../services/websocket.service';
import { SharedService } from '../../services/shared.service'
import { MatDialog } from '@angular/material/dialog';
import { ErrorComponent, ErrorMessage } from '../error/error.component';
import { ErrorWindow } from 'src/app/handlers/error';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule, MatLabel } from '@angular/material/form-field';
import { MatButtonModule } from '@angular/material/button';
import {
  FormsModule,
  ReactiveFormsModule,
} from '@angular/forms';
import { MatInputModule } from '@angular/material/input';

export interface Response {
	id: number
    token: string
}

@Component({
  standalone: true,
  selector: 'app-login',
  imports: [MatIconModule, MatFormFieldModule, FormsModule, MatButtonModule, 
	MatInputModule, RouterModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {
	hide: boolean = true;
	login: string = "";
	password: string = "";
	
	constructor(
		private http: HttpClient, 
		private router: Router, 
		private cookieService: CookieService, 
		private websocket: WebsocketService,
		private sharedService: SharedService,
		public dialog: MatDialog,
	){}
	
	onSubmit(){
		const postData = new FormData();
		const login = this.login;
		postData.append('login',  login);
		postData.append('password',  this.password);
		
		this.http.post<Response>('http://localhost:8081/users/login', postData).subscribe({
			next: response => {
				this.saveData(response.token, login, response.id);

				this.websocket.start();
				this.sharedService.start();
				this.sharedService.bol = true;
				this.router.navigate(['/']);
			},
			error: err => {
				new ErrorWindow(err, this.dialog);
				console.log(err);
			}
		});
		this.clearData();
	}

	saveData(token: string, login: string, id: number){
		const expiryDate = new Date();
		expiryDate.setHours(expiryDate.getHours() + 1);
		this.cookieService.set('token', token, expiryDate, '/', 'localhost', true, 'None');

		localStorage.setItem("login", login);
		localStorage.setItem("id", id.toString());
	}

	clearData(){
		this.login = "";
		this.password = "";
	}
}
