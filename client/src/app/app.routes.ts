import { Routes } from '@angular/router';
import { RegisterComponent } from 'src/app/components/register/register.component';
import { LoginComponent } from 'src/app/components/login/login.component';
import { ProfileComponent } from 'src/app/components/profile/profile.component';
import { FindComponent } from 'src/app/components/find/find.component';
import { PostsComponent } from 'src/app/components/posts/posts.component';
import { FollowingComponent } from 'src/app/components/following/following.component';
import { MessagesComponent } from 'src/app/components/messages/messages.component';
import { ProjectsComponent } from 'src/app/components/projects/projects.component';
import { AuthGuard } from 'src/auth.guard';

export const routes: Routes = [
	{path:"register", component: RegisterComponent},
	{path:"login", component: LoginComponent},
	{path:"profile", component: ProfileComponent, canActivate: [AuthGuard]},
	{path:"find", component: FindComponent, canActivate: [AuthGuard]},
	{path:"posts", component: PostsComponent, canActivate: [AuthGuard]},
	{path:"following", component: FollowingComponent, canActivate: [AuthGuard]},
	{path:"messages", component: MessagesComponent, canActivate: [AuthGuard]},
	{path:"projects", component: ProjectsComponent, canActivate: [AuthGuard]}
];