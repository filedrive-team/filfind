import _Cookies from 'js-cookie';

import { UserInfo } from '@/models/userInfo';
import { jwt_decode } from '@/utils';
import userStore from '@/store/modules/user';
import { Storages } from '@/data/storages';

enum CookieKeys {
  UserInfo = 'access_token',
  Avatar = 'avatar',
  Account = 'account',
}

export class Cookies {
  static setUserInfo(value: UserInfo): boolean {
    const { access_token } = value;
    const _token = jwt_decode(access_token);
    const { exp } = _token;
    const _expires = exp * 1000;
    userStore.SET_USERINFO(value);
    Storages.setLoginStatus(JSON.stringify(value));
    return _Cookies.set(CookieKeys.UserInfo, JSON.stringify(value), {
      expires: new Date(_expires),
    });
  }

  static getUserInfo(): UserInfo | undefined {
    const _userInfo = _Cookies.get(CookieKeys.UserInfo);
    if (!_userInfo) return undefined;
    return JSON.parse(_userInfo);
  }

  static removeUserInfo(): boolean {
    userStore.SET_USERINFO(undefined);
    Storages.setLoginStatus(undefined);
    return _Cookies.remove(CookieKeys.UserInfo);
  }

  static getAuthorization(): any {
    return this.getUserInfo()
      ? 'Bearer ' + this.getUserInfo()?.access_token
      : undefined;
  }
}
