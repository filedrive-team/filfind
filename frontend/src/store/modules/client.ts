import { action, runInAction, makeObservable, observable } from 'mobx';
import {
  getClientDetail,
  getClientHistoryDealStats,
  getClientProfile,
  getClientReviews,
  postClientDetail,
  postUserProfile,
} from '@/api/modules';
import {
  ClientDetailParam,
  ClientHistoryDealStatsParam,
  ClientProfileParam,
  ClientReviewsParam,
  PostClientDetailParam,
} from '@/api/modules/interface';
import { ProfileResponse } from '@/models/types';
interface ClientDetailResponse {
  address_id: string;
  bandwidth: string;
  data_cap: string;
  monthly_storage: string;
  service_requirement: string;
  use_case: string;
}

interface ListParams {
  page: number;
  page_size: number;
  address_id: string;
}

class ClientStore {
  @observable editProfileVisible: boolean = false;
  @observable editInfoVisible: boolean = false;
  @observable historyDealStatsList: [] = [];
  @observable historyDealStatsTotal: number = 0;
  @observable historyDealStatsParam: ListParams = {
    page: 1,
    page_size: 10,
    address_id: '',
  };

  @observable reviewsList: [] = [];
  @observable reviewsTotal: number = 0;
  @observable reviewsParam: ListParams = {
    page: 1,
    page_size: 10,
    address_id: '',
  };

  @observable clientDetail: ClientDetailResponse = {
    address_id: '',
    bandwidth: '',
    data_cap: '',
    monthly_storage: '',
    service_requirement: '',
    use_case: '',
  };

  @observable clientProfile: ProfileResponse = {
    address_id: '',
    address_robust: '',
    avatar: '',
    contact_email: '',
    description: '',
    github: '',
    location: '',
    logo: '',
    name: '',
    slack: '',
    twitter: '',
    type: '',
  };

  @observable postClientDetail: PostClientDetailParam = {
    bandwidth: '',
    monthly_storage: '',
    service_requirement: '',
    use_case: '',
  };

  @observable userProfile: ProfileResponse = {
    avatar: '',
    logo: '',
    contact_email: '',
    description: '',
    github: '',
    location: '',
    name: '',
    slack: '',
    twitter: '',
  };

  constructor() {
    makeObservable(this);
  }

  @action
  SET_HISTORYDEALSTATSPARAM = (data: ListParams) => {
    runInAction(() => {
      this.historyDealStatsParam = data;
    });
  };

  @action
  SET_REVIEWSPARAM = (data: ListParams) => {
    runInAction(() => {
      this.reviewsParam = data;
    });
  };

  @action
  SET_PROFILEVISIBLE = (editProfileVisible: boolean) => {
    runInAction(() => {
      this.editProfileVisible = editProfileVisible;
    });
  };

  @action
  SET_INFOVISIBLE = (editInfoVisible: boolean) => {
    runInAction(() => {
      this.editInfoVisible = editInfoVisible;
    });
  };

  @action
  SET_MODIFYCLIENTDETAIL = (data: PostClientDetailParam) => {
    runInAction(() => {
      this.postClientDetail = data;
    });
  };

  @action
  SET_CLIENTPROFILE = (data: ProfileResponse) => {
    runInAction(() => {
      this.clientProfile = data;
    });
  };

  @action.bound
  _getClientDetail(data: ClientDetailParam) {
    return new Promise((resolve) => {
      getClientDetail(data).then((res) => {
        runInAction(() => {
          if (res.data) {
            this.clientDetail = res.data;
          }
        });
      });
    });
  }

  @action.bound
  _getClientHistoryDealStats(data: ClientHistoryDealStatsParam) {
    return new Promise((resolve) => {
      getClientHistoryDealStats(data).then((res) => {
        runInAction(() => {
          if (res.data) {
            const { total = 0, list = [] } = res.data;
            this.historyDealStatsList = list;
            this.historyDealStatsTotal = total;
          }
        });
      });
    });
  }

  @action.bound
  _getClientProfile(data: ClientProfileParam) {
    return new Promise((resolve) => {
      getClientProfile(data).then((res) => {
        runInAction(() => {
          if (res.data) {
            this.clientProfile = res.data;
          }
        });
      });
    });
  }

  @action.bound
  _getClientReviews(data: ClientReviewsParam) {
    return new Promise((resolve) => {
      getClientReviews(data).then((res) => {
        runInAction(() => {
          const { list, total } = res.data;
          if (res.data && list) {
            this.reviewsList = list;
            this.reviewsTotal = total;
          }
        });
      });
    });
  }

  @action.bound
  _postClientDetail(data: PostClientDetailParam) {
    return new Promise((resolve, reject) => {
      postClientDetail(data).then((res) => {
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

  @action.bound
  _postUserProfile(data: ProfileResponse) {
    return new Promise((resolve) => {
      postUserProfile(data).then((res) => {
        runInAction(() => {
          resolve(res);
        });
      });
    });
  }
}

const clientStore = new ClientStore();
export default clientStore;
