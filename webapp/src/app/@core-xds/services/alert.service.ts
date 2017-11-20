import { Injectable, SecurityContext } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';


export type AlertType = 'error' | 'warning' | 'info' | 'success';

export interface IAlert {
  type: AlertType;
  msg: string;
  show?: boolean;
  dismissible?: boolean;
  dismissTimeout?: number;     // close alert after this time (in seconds)
  id?: number;
}

@Injectable()
export class AlertService {
  public alerts: Observable<IAlert[]>;

  private _alerts: IAlert[];
  private alertsSubject = <Subject<IAlert[]>>new Subject();
  private uid = 0;
  private defaultDismissTmo = 5; // in seconds

  constructor() {
    this.alerts = this.alertsSubject.asObservable();
    this._alerts = [];
    this.uid = 0;
  }

  public error(msg: string, dismissTime?: number) {
    this.add({
      type: 'error', msg: msg, dismissible: true, dismissTimeout: dismissTime
    });
  }

  public warning(msg: string, dismissible?: boolean) {
    this.add({ type: 'warning', msg: msg, dismissible: true, dismissTimeout: (dismissible ? this.defaultDismissTmo : 0) });
  }

  public info(msg: string) {
    this.add({ type: 'info', msg: msg, dismissible: true, dismissTimeout: this.defaultDismissTmo });
  }

  public add(al: IAlert) {
    const msg = String(al.msg).replace('\n', '<br>');
    // this._alerts.push({
    this._alerts = [{
      show: true,
      type: al.type,
      msg: msg,
      dismissible: al.dismissible || true,
      dismissTimeout: (al.dismissTimeout * 1000) || 0,
      id: this.uid,
    }];
    this.uid += 1;
    this.alertsSubject.next(this._alerts);

  }

  public del(al: IAlert) {
    /*
    const idx = this._alerts.findIndex((a) => a.id === al.id);
    if (idx > -1) {
      this._alerts.splice(idx, 1);
      this.alertsSubject.next(this._alerts);
    }
    */
    this._alerts = [];
    this.alertsSubject.next(this._alerts);
  }
}
