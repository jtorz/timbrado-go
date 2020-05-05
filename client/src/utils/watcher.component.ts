import { Injectable } from '@angular/core';
import {
    HttpInterceptor,
    HttpRequest,
    HttpErrorResponse,
    HttpHandler,
    HttpEvent
} from '@angular/common/http';

import { catchError, finalize } from 'rxjs/operators';
import { Observable, throwError } from 'rxjs';

import { Router } from '@angular/router';
import { LoadingService } from 'src/services/loafing.service';


@Injectable()
export class Watcher implements HttpInterceptor {
    constructor(private router: Router, private loadingService: LoadingService) { }
    /**
     * intercepta las peticiones para verificar que el usuario tenga una sesion activa
     * agrega el CSRFKey a todas las peticiones
     */
    intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
        let request = req;
        this.loadingService.addRequestCount();
        console.log('intercepted')
        return next.handle(request).pipe(finalize(() => this.loadingService.subtractRequestCount()));
    }
}
