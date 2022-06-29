import {
  action,
  runInAction,
  makeObservable,
  observable,
  computed,
} from 'mobx';
import {
  getSpOwnerProfile,
  getSpOwnerReviews,
  getSpServiceDetail,
  postClientReview,
  postProviderDetail,
  postUserProfile,
} from '@/api/modules';
import {
  SpOwnerProfileParam,
  SpOwnerReviewsParam,
  SpServiceDetailParam,
  PostClientReviewParam,
  ProviderDetailParam,
} from '@/api/modules/interface';
import { notification } from 'antd';
import { ProfileResponse } from '@/models/types';

interface ListParams {
  page: number;
  page_size: number;
  address_id: string;
}

class StorageProvideStore {
  @observable editProfileVisible: boolean = false;
  @observable editInfoVisible: boolean = false;
  @observable reviewsList: [] = [];
  @observable reviewsTotal: number = 0;
  @observable reviewsParam: ListParams = {
    page: 1,
    page_size: 10,
    address_id: '',
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

  @observable spOwnerProfile: ProfileResponse = {
    address_id: '',
    address_robust: '',
    avatar: '',
    contact_email: '',
    description: '',
    github: '',
    location: '',
    logo: '',
    name: '',
    reputation_score: 0,
    review_score: 0,
    reviews: 0,
    slack: '',
    twitter: '',
    type: '',
  };

  @observable spServiceDetailList: [] = [];

  @observable providerDetail: ProviderDetailParam = {
    address: '',
    available_deals: '',
    bandwidth: '',
    certification: '',
    experience: '',
    is_member: '',
    parallel_deals: '',
    renewable_energy: '',
    sealing_speed: '',
  };

  @observable addReveiews: PostClientReviewParam = {
    title: '',
    content: '',
    provider: '',
    score: 0,
  };

  @computed get MinerIDList() {
    const MinerIDList = this.spServiceDetailList.map((n) => {
      return n['address'];
    });
    return MinerIDList;
  }

  @computed get valid() {
    const _bool =
      this.addReveiews.content === '' ||
      this.addReveiews.provider === '' ||
      this.addReveiews.title === '' ||
      this.addReveiews.score === 0;
    return _bool;
  }

  constructor() {
    makeObservable(this);
  }

  @action
  SET_PROVIDERDETAIL = (data: ProviderDetailParam) => {
    runInAction(() => {
      this.providerDetail = data;
    });
  };

  @action
  SET_SPOWNERREVIEWSPARAM = (data: ListParams) => {
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
  SET_USERPROFILE = (data: ProfileResponse) => {
    runInAction(() => {
      this.userProfile = data;
    });
  };

  SET_ADDREVIEWS = (data: PostClientReviewParam) => {
    runInAction(() => {
      this.addReveiews = data;
    });
  };

  @action.bound
  _getSpOwnerProfile(data: SpOwnerProfileParam) {
    return new Promise((resolve) => {
      getSpOwnerProfile(data).then((res) => {
        runInAction(() => {
          if (res.data) {
            const {
              avatar,
              logo,
              contact_email,
              description,
              github,
              location,
              name,
              slack,
              twitter,
            } = res.data;
            this.userProfile = {
              avatar,
              logo,
              contact_email,
              description,
              github,
              location,
              name,
              slack,
              twitter,
            };
            this.spOwnerProfile = res.data;
          }
        });
      });
    });
  }

  @action.bound
  _getSpOwnerReviews(data: SpOwnerReviewsParam) {
    return new Promise((resolve) => {
      getSpOwnerReviews(data).then((res) => {
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
  _getSpServiceDetail(data: SpServiceDetailParam) {
    return new Promise((resolve) => {
      getSpServiceDetail(data).then((res) => {
        runInAction(() => {
          if (res.data) {
            this.spServiceDetailList = res.data;
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

  @action.bound
  _postClientReview(data: PostClientReviewParam) {
    return new Promise((resolve, reject) => {
      postClientReview(data).then((res) => {
        runInAction(() => {
          if (res.data) {
            resolve(res.data);
          } else {
            const _message = res.msg;
            notification.error({
              key: 'pictureError',
              message: _message,
            });
          }
        });
      });
    });
  }

  @action.bound
  _postProviderDetail(data: ProviderDetailParam) {
    return new Promise((resolve, reject) => {
      postProviderDetail(data).then((res) => {
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

const storageProvideStore = new StorageProvideStore();
export default storageProvideStore;
