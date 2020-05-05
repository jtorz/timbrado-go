import { Component, OnInit } from '@angular/core';
import { FormGroup, Validators, FormBuilder } from '@angular/forms';
import { AppService } from 'src/services/app.service';
import { WS } from 'src/model/models';

@Component({
  selector: 'app-ws',
  templateUrl: './ws.component.html',
  styleUrls: ['./ws.component.scss']
})
export class WsComponent implements OnInit {
  formGroup: FormGroup;
  catalogoWS: WS[];
  constructor(private formBuilder: FormBuilder, private appservice: AppService) { }

  ngOnInit(): void {
    this.createForm();
  }
  createForm() {
    this.formGroup = this.formBuilder.group({
      'usuario': [null, /* [Validators.required] */],
      'password': [null, /* [Validators.required] */],
      'ws': [null, [Validators.required]],
    });
    this.appservice.catalogoWS().subscribe(
      (c) => {
        this.catalogoWS = c;
      },
      (err) => {
        alert('Error: ' + JSON.stringify(err));
        console.log(err);
      }
    )
    this.appservice.getInfoWS().subscribe(
      (c) => {
        this.formGroup.setValue({
          'usuario': c.usuario,
          'password': c.password,
          'ws': c.ws.id,
        });
      },
      (err) => {
        alert('Error: ' + JSON.stringify(err));
        console.log(err);
      }
    )
  }
  onSubmit() {
    if (!this.formGroup.valid) return;
    let v = this.formGroup.value;
    this.appservice.setInfoWS( v['ws'], v['usuario'], v['password']).subscribe(
      ()=>{},
      (err) => {
        alert('Error: ' + JSON.stringify(err));
        console.log(err);
      }
    )
  }

  getErrorEmail() {
    return this.formGroup.get('usuario').hasError('required') ? 'Requerido' : '';
  }
  getErrorPassword() {
    return this.formGroup.get('password').hasError('required') ? 'Requerido' : '';
  }
  getErrorWS() {
    return this.formGroup.get('ws').hasError('required') ? 'Requerido' : '';
  }
}
