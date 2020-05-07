import { Injectable } from '@angular/core';
import {
    HttpInterceptor,
    HttpRequest,
    HttpHandler,
    HttpEvent
} from '@angular/common/http';

import { finalize } from 'rxjs/operators';
import { Observable } from 'rxjs';

import { Router } from '@angular/router';
import { LoadingService } from 'src/services/loafing.service';


@Injectable()
export class Watcher implements HttpInterceptor {
    constructor(private router: Router, private loadingService: LoadingService) { }
    intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
        let request = req;
        this.loadingService.addRequestCount();
        return next.handle(request).pipe(finalize(() => this.loadingService.subtractRequestCount()));
    }
}
