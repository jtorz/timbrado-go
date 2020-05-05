import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { MatTabsModule } from '@angular/material/tabs';
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from '@angular/material/input';
import { MatCardModule } from "@angular/material/card";
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatSelectModule } from '@angular/material/select';
import { MatProgressBarModule } from '@angular/material/progress-bar';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { CertComponent } from './cert/cert.component';
import { WsComponent } from './ws/ws.component';
import { CfdiComponent } from './cfdi/cfdi.component';
import { TimbradoComponent } from './timbrado/timbrado.component';

import { XmlPipe } from '../utils/XmlPipe'
import { Watcher } from 'src/utils/watcher.component';

@NgModule({
  declarations: [
    AppComponent,
    CertComponent,
    WsComponent,
    CfdiComponent,
    TimbradoComponent,
    XmlPipe
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    MatTabsModule,
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatCardModule,
    MatTooltipModule,
    MatSelectModule,
    MatProgressBarModule,
  ],
  providers: [
    { provide: HTTP_INTERCEPTORS, useClass: Watcher, multi: true },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
