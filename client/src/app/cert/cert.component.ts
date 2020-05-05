import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { SatCertificate } from 'src/model/models';
import { AppService } from 'src/services/app.service';

@Component({
  selector: 'app-cert',
  templateUrl: './cert.component.html',
  styleUrls: ['./cert.component.scss']
})
export class CertComponent implements OnInit {

  cert: SatCertificate ;
  constructor(private appservice: AppService) { }

  @ViewChild("certFile") certFile: ElementRef;
  @ViewChild("keyFile") keyFile: ElementRef;
  @ViewChild("keypass") keyPass: ElementRef;
  ngOnInit(): void {
    this.appservice.getCert().subscribe((c)=> {
      this.cert = c
    },
    (error) =>{
      console.log(error);
    })
  }
  loadCer(event) {
    let thiz = this;
    if(this.certFile.nativeElement.files[0] == null){
      alert('Seleccione el archivo .cer del certificado.');
      return;
    }
    this.animate(event.target);
    this.appservice.uploadCert(this.certFile.nativeElement.files[0]).subscribe({
      next(c) { thiz.cert = c },
      error(err) { alert('Error: ' + JSON.stringify(err)); },
      complete() { thiz.deanimate(event.target); }
    })
  }

  loadKey(event) {
    let thiz = this;
    if(this.keyFile.nativeElement.files[0] == null){
      alert('Seleccione el archivo .key del certificado.');
      return;
    }
    if(this.keyPass.nativeElement.value == ''){
      alert('Ingrese la contraseña de la llave certificado.');
      return;
    }
    this.animate(event.target);
    this.appservice.uploadKey(this.keyFile.nativeElement.files[0], this.keyPass.nativeElement.value).subscribe({
      next(response) { console.log(response); },
      error(err) { alert('Error: ' + JSON.stringify(err)); },
      complete() { thiz.deanimate(event.target); }
    })
  }

  private animate(target) {
    target.classList.remove('spin-animation');
    target.classList.add('spin-animation');
  }
  private deanimate(target) {
    setTimeout(() => {
      target.classList.remove('spin-animation');
    }, 500);
  }
}
