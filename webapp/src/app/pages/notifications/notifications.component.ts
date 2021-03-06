/**
* @license
* Copyright (C) 2017 "IoT.bzh"
* Author Sebastien Douheret <sebastien@iot.bzh>
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*   http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

import { Component } from '@angular/core';
import { ToasterService, ToasterConfig, Toast, BodyOutputType } from 'angular2-toaster';
import { Observable } from 'rxjs/Observable';
import { AlertService, IAlert } from '../../@core-xds/services/alert.service';

import 'style-loader!angular2-toaster/toaster.css';

@Component({
  selector: 'ngx-notifications',
  styleUrls: ['./notifications.component.scss'],
  template: '<toaster-container [toasterconfig]="config"></toaster-container>',
})
export class NotificationsComponent {

  config: ToasterConfig;

  private position = 'toast-top-full-width';
  private animationType = 'slideDown';
  private toastsLimit = 10;
  private toasterService: ToasterService;
  private alerts$: Observable<IAlert[]>;

  constructor(
    toasterService: ToasterService,
    private alertSvr: AlertService,
  ) {
    this.toasterService = toasterService;

    this.alertSvr.alerts.subscribe(alerts => {
      if (alerts.length === 0) {
        this.clearToasts();
      } else {
        alerts.forEach(al => {
          const title = al.type.toUpperCase();
          this.showToast(al.type, title, al.msg, al.dismissTimeout);
        });
      }
    });
  }

  private showToast(type: string, title: string, body: string, tmo: number) {
    this.config = new ToasterConfig({
      positionClass: this.position,
      timeout: tmo,
      newestOnTop: true,
      tapToDismiss: true, // is Hide OnClick
      preventDuplicates: false,
      animation: this.animationType,
      limit: this.toastsLimit,
    });
    const toast: Toast = {
      type: type,
      title: title,
      body: body,
      timeout: tmo,
      showCloseButton: true,
      bodyOutputType: BodyOutputType.TrustedHtml,
    };
    this.toasterService.popAsync(toast);
  }

  clearToasts() {
    this.toasterService.clear();
  }
}
