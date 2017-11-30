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

import { Injectable, SecurityContext } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

import { XDSAgentService } from '../services/xdsagent.service';

import 'rxjs/add/observable/throw';

export interface ISdk {
    id: string;
    profile: string;
    version: string;
    arch: number;
    path: string;
}

@Injectable()
export class SdkService {
    public Sdks$: Observable<ISdk[]>;

    private _sdksList = [];
    private current: ISdk;
    private sdksSubject = <BehaviorSubject<ISdk[]>>new BehaviorSubject(this._sdksList);

    constructor(private xdsSvr: XDSAgentService) {
        this.current = null;
        this.Sdks$ = this.sdksSubject.asObservable();

        this.xdsSvr.XdsConfig$.subscribe(cfg => {
            if (!cfg || cfg.servers.length < 1) {
                return;
            }
            // FIXME support multiple server
            // cfg.servers.forEach(svr => {
            this.xdsSvr.getSdks(cfg.servers[0].id).subscribe((s) => {
                this._sdksList = s;
                this.sdksSubject.next(s);
            });
        });
    }

    public setCurrent(s: ISdk) {
        this.current = s;
    }

    public getCurrent(): ISdk {
        return this.current;
    }

    public getCurrentId(): string {
        if (this.current && this.current.id) {
            return this.current.id;
        }
        return '';
    }

    public add(sdk: ISdk): Observable<ISdk> {
      // TODO SEB
      return Observable.throw('Not implement yet');
    }

    public delete(sdk: ISdk): Observable<ISdk> {
      // TODO SEB
      return Observable.throw('Not implement yet');
    }
}
