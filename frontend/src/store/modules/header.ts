import { action, runInAction, makeObservable, observable } from 'mobx';

class HeaderStore {
  @observable activeIndex: number = 0;

  constructor() {
    makeObservable(this);
  }

  @action
  SET_ACTIVEINDEX = (data: number) => {
    runInAction(() => {
      this.activeIndex = data;
    });
  };
}

const headerStore = new HeaderStore();
export default headerStore;
