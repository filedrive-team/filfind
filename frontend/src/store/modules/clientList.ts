import { action, runInAction, makeObservable, observable } from 'mobx';
import { getClients } from '@/api/modules';
import { clientsParam } from '@/api/modules/interface';

class ClientListStore {
  @observable loading: boolean = true;
  @observable params: clientsParam = {
    page: 1,
    page_size: 10,
    sort_by: 'data_cap',
    order: 'desc',
    search: '',
  };
  @observable list: [] = [];
  @observable total: number = 0;

  constructor() {
    makeObservable(this);
  }

  @action
  SET_LOADING = (loading: boolean) => {
    runInAction(() => {
      this.loading = loading;
    });
  };
  @action
  SET_CLIENTSPARAM = (params: clientsParam) => {
    runInAction(() => {
      this.params = params;
    });
  };

  @action.bound
  _getClients(data: clientsParam) {
    return new Promise((resolve) => {
      getClients(data).then((res) => {
        runInAction(() => {
          if (res.data) {
            const { total = 0, list = [] } = res.data;
            this.list = list;
            this.total = total;
          }
        });
      });
    });
  }
}

const clientListStore = new ClientListStore();
export default clientListStore;
