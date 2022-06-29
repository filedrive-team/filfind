import {
  action,
  runInAction,
  makeObservable,
  observable,
  computed,
} from 'mobx';
import { getProviders } from '@/api/modules';
import Decimal from 'decimal.js';
interface ParamsRequest {
  page: number;
  page_size: number;
  sort_by: string;
  order: string;
  sps_status: string;
  region: string;
  raw_power_range: Range;
  storage_success_rate_range: Range;
  reputation_score_range: Range;
  review_score_range: Range;
  search: string;
}
interface Range {
  min: string;
  max: string;
}
class HomeStore {
  @observable loading: boolean = true;
  @observable params: ParamsRequest = {
    page: 1,
    page_size: 10,
    sort_by: 'reputation_score',
    order: 'desc',
    sps_status: 'all',
    region: 'all',
    raw_power_range: { min: '', max: '' },
    storage_success_rate_range: { min: '', max: '' },
    reputation_score_range: { min: '', max: '' },
    review_score_range: { min: '', max: '' },
    search: '',
  };
  @observable list: [] = [];
  @observable total: number = 0;

  constructor() {
    makeObservable(this);
  }
  @computed get raw_power_range_min() {
    const { min } = this.params.raw_power_range;
    return min === '' ? '-1' : formatValue(min);
  }
  @computed get raw_power_range_max() {
    const { max } = this.params.raw_power_range;
    return max === '' ? '-1' : formatValue(max);
  }

  @computed get storage_success_rate_range_min() {
    const { min } = this.params.storage_success_rate_range;
    return min === '' ? '-1' : min;
  }
  @computed get storage_success_rate_range_max() {
    const { max } = this.params.storage_success_rate_range;
    return max === '' ? '-1' : max;
  }

  @computed get reputation_score_range_min() {
    const { min } = this.params.reputation_score_range;
    return min === '' ? '-1' : min;
  }
  @computed get reputation_score_range_max() {
    const { max } = this.params.reputation_score_range;
    return max === '' ? '-1' : max;
  }
  @computed get review_score_range_min() {
    const { min } = this.params.review_score_range;
    return min === '' ? '-1' : min;
  }
  @computed get review_score_range_max() {
    const { max } = this.params.review_score_range;
    return max === '' ? '-1' : max;
  }

  @action
  SET_LOADING = (loading: boolean) => {
    runInAction(() => {
      this.loading = loading;
    });
  };
  @action
  SET_PROVIDERSPARAM = (params: ParamsRequest) => {
    runInAction(() => {
      this.params = params;
    });
  };

  @action.bound
  _getProviders(data: ParamsRequest) {
    const formatParams = {
      ...data,
      raw_power_range: JSON.stringify({
        min: this.raw_power_range_min,
        max: this.raw_power_range_max,
      }),
      storage_success_rate_range: JSON.stringify({
        min: this.storage_success_rate_range_min,
        max: this.storage_success_rate_range_max,
      }),
      reputation_score_range: JSON.stringify({
        min: this.reputation_score_range_min,
        max: this.reputation_score_range_max,
      }),
      review_score_range: JSON.stringify({
        min: this.review_score_range_min,
        max: this.review_score_range_max,
      }),
    };

    return new Promise((resolve) => {
      getProviders(formatParams).then((res) => {
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
const formatValue = (value) => {
  if (!value) return -1;
  const unit = Decimal.pow(1024, 4).toNumber();
  const _value = new Decimal(value);
  const res = _value.times(unit).toNumber();
  return res;
};
const homeStore = new HomeStore();
export default homeStore;
