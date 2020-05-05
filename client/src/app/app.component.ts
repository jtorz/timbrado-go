import { Component } from '@angular/core';
import { Subscription } from 'rxjs';
import { LoadingService } from 'src/services/loafing.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  loading: boolean = false;
  loadingSubscription: Subscription;

  constructor(private loadingService: LoadingService) { }

  ngOnInit(): void {
    this.loadingSubscription = this.loadingService.getLoadingStatus().subscribe((value) => {
      this.loading = value;
    });

  }
  ngOnDestroy() {
    this.loadingSubscription.unsubscribe();
  }

}
