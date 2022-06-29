import { action, runInAction, makeObservable, observable } from 'mobx';
import { postUserModifyPassword } from '@/api/modules';
import { UserModifyPasswordParam } from '@/api/modules/interface';
interface FormItem {
  value: string;
  valid: boolean;
  message: string;
}
class ChangePasswordStore {
  @observable recentPassword: FormItem = {
    value: '',
    valid: false,
    message: '',
  };

  @observable newPassword: FormItem = {
    value: '',
    valid: false,
    message: '',
  };

  @observable confirmPassword: FormItem = {
    value: '',
    valid: false,
    message: '',
  };

  constructor() {
    makeObservable(this);
  }
  @action
  SET_RECENTPASSWORD = (data: FormItem) => {
    runInAction(() => {
      this.recentPassword = data;
    });
  };
  @action
  SET_NEWPASSWORD = (data: FormItem) => {
    runInAction(() => {
      this.newPassword = data;
    });
  };

  @action
  SET_CONFIRMPASSWORD = (data: FormItem) => {
    runInAction(() => {
      this.confirmPassword = data;
    });
  };

  @action.bound
  _postUserModifyPassword(data: UserModifyPasswordParam) {
    return new Promise((resolve, reject) => {
      postUserModifyPassword(data).then((res) => {
        runInAction(() => {
          if (res.data) {
            resolve(res.data);
          } else {
            reject(res.data);
          }
        });
      });
    });
  }
}

const changePasswordStore = new ChangePasswordStore();
export default changePasswordStore;
