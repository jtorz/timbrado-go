import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { TimbradoComponent } from './timbrado/timbrado.component';


const routes: Routes = [
  {path:"",component:TimbradoComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
