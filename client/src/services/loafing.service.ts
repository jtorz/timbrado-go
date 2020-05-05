import { Injectable } from '@angular/core';
import { Subject, Observable } from "rxjs";

@Injectable({
    providedIn: 'root'
})
export class LoadingService {

    constructor() { }


    private _loading: boolean = false;
    private _requestNumber: number = 0;

    loadingStatus: Subject<boolean> = new Subject();

    get loading(): boolean {
        return this._loading;
    }

    set loading(value) {
        this._loading = value;
        this.loadingStatus.next(value);
    }

    getLoadingStatus(): Observable<any> {
        return this.loadingStatus.asObservable();
    }

    startLoading() {
        this.loading = true;
    }

    stopLoading() {
        this.loading = false;
    }

    addRequestCount() {
        this._requestNumber++;
        if (this._requestNumber == 1)
            this.startLoading();
    }

    subtractRequestCount() {
        this._requestNumber--;
        if (this._requestNumber < 1)
            this.stopLoading();
    }
}
