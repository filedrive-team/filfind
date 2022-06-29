import {
  action,
  computed,
  makeObservable,
  observable,
  runInAction,
} from 'mobx';
import { Cookies } from '@/data/cookies';
import { UserInfo } from '@/models/userInfo';

class UserStore {
  @observable count: number = 0;
  @observable userInfo: UserInfo | undefined;

  constructor() {
    makeObservable(this);
  }

  @computed get isLogin(): boolean {
    return this.userInfo !== undefined;
  }

  @action
  SET_USERINFO = (data: UserInfo | undefined) => {
    runInAction(() => {
      this.userInfo = data;
    });
  };

  @action.bound
  async add() {
    this.count++;
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve('');
      }, 5000);
    });
  }
}

const userStore = new UserStore();

export default userStore;
