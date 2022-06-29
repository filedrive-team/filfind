import { action, computed, makeObservable, observable } from 'mobx';

class TestStore {
  @observable count: number = 0;

  constructor() {
    makeObservable(this);
  }

  @computed get doubleCount() {
    return this.count * 2;
  }

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

const testStore = new TestStore();

export default testStore;
