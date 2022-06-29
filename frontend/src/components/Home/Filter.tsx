import styles from './Filter.module.scss';
import { Input, Popover } from 'antd';
import { DownOutlined } from '@ant-design/icons';
import { useState } from 'react';
import homeStore from '@/store/modules/home';
import { observer } from 'mobx-react';
import filFind from '@/assets/images/home/filFind.png';
const { Search } = Input;
const subTitle =
  'A discovery platform for SPs and Clients in Filecoin Ecosystem';

const onSearch = (value: string) => {
  homeStore.SET_PROVIDERSPARAM({
    ...homeStore.params,
    page: 1,
    search: value,
  });
  homeStore._getProviders(homeStore.params);
};
const Filter = () => {
  const [conditionVisible, setConditionVisible] = useState(false);

  const openCondition = () => {
    setConditionVisible(true);
  };
  const confirm = () => {
    setConditionVisible(false);
    homeStore.SET_LOADING(true);
    homeStore._getProviders(homeStore.params);
  };
  const clear = () => {
    homeStore.SET_PROVIDERSPARAM({
      ...homeStore.params,
      page: 1,
      raw_power_range: { min: '', max: '' },
      storage_success_rate_range: { min: '', max: '' },
      reputation_score_range: { min: '', max: '' },
      review_score_range: { min: '', max: '' },
    });
    setConditionVisible(false);
    homeStore._getProviders(homeStore.params);
  };
  return (
    <div className={styles.filter}>
      <div className={styles.filFind}>
        <img src={filFind} alt="" />
      </div>
      <div className={styles.subTitle}>{subTitle}</div>
      <div className={styles.searchRow}>
        <Search
          placeholder="Miner ID/Name/Location"
          enterButton="Search"
          onSearch={onSearch}
        />
        <div className={styles.condition}>
          <Popover
            trigger="click"
            placement="bottom"
            visible={conditionVisible}
            content={
              <div className="filterCondition">
                <div className="item">
                  <div className="from">Raw Power: from</div>
                  <Input
                    placeholder=""
                    value={homeStore.params.raw_power_range.min}
                    onChange={(event) => {
                      const _value = event.target.value;
                      homeStore.SET_PROVIDERSPARAM({
                        ...homeStore.params,
                        page: 1,
                        raw_power_range: {
                          min: _value,
                          max: homeStore.params.raw_power_range.max,
                        },
                      });
                    }}
                  />
                  <div className="to">to</div>
                  <Input
                    placeholder=""
                    value={homeStore.params.raw_power_range.max}
                    onChange={(event) => {
                      const _value = event.target.value;
                      homeStore.SET_PROVIDERSPARAM({
                        ...homeStore.params,
                        page: 1,
                        raw_power_range: {
                          min: homeStore.params.raw_power_range.min,
                          max: _value,
                        },
                      });
                    }}
                  />
                  <div className="unit">(Min: 0 TiB)</div>
                </div>
                <div className="item">
                  <div className="from">Deal Success Rate: from</div>
                  <Input
                    placeholder=""
                    value={homeStore.params.storage_success_rate_range.min}
                    onChange={(event) => {
                      const _value = event.target.value;
                      homeStore.SET_PROVIDERSPARAM({
                        ...homeStore.params,
                        page: 1,
                        storage_success_rate_range: {
                          min: _value,
                          max: homeStore.params.storage_success_rate_range.max,
                        },
                      });
                    }}
                  />
                  <div className="to">to</div>
                  <Input
                    placeholder=""
                    value={homeStore.params.storage_success_rate_range.max}
                    onChange={(event) => {
                      const _value = event.target.value;
                      homeStore.SET_PROVIDERSPARAM({
                        ...homeStore.params,
                        page: 1,
                        storage_success_rate_range: {
                          min: homeStore.params.storage_success_rate_range.min,
                          max: _value,
                        },
                      });
                    }}
                  />
                  <div className="unit">(Min: 0, Max: 100)</div>
                </div>
                <div className="item">
                  <div className="from">Reputation Score: from</div>
                  <Input
                    placeholder=""
                    value={homeStore.params.reputation_score_range.min}
                    onChange={(event) => {
                      const _value = event.target.value;
                      homeStore.SET_PROVIDERSPARAM({
                        ...homeStore.params,
                        page: 1,
                        reputation_score_range: {
                          min: _value,
                          max: homeStore.params.reputation_score_range.max,
                        },
                      });
                    }}
                  />
                  <div className="to">to</div>
                  <Input
                    placeholder=""
                    value={homeStore.params.reputation_score_range.max}
                    onChange={(event) => {
                      const _value = event.target.value;
                      homeStore.SET_PROVIDERSPARAM({
                        ...homeStore.params,
                        page: 1,
                        reputation_score_range: {
                          min: homeStore.params.reputation_score_range.min,
                          max: _value,
                        },
                      });
                    }}
                  />
                  <div className="unit">(Min: 0, Max: 10.0)</div>
                </div>
                <div className="item">
                  <div className="from">Client Review Score: from</div>
                  <Input
                    placeholder=""
                    value={homeStore.params.review_score_range.min}
                    onChange={(event) => {
                      const _value = event.target.value;
                      homeStore.SET_PROVIDERSPARAM({
                        ...homeStore.params,
                        page: 1,
                        review_score_range: {
                          min: _value,
                          max: homeStore.params.review_score_range.max,
                        },
                      });
                    }}
                  />
                  <div className="to">to</div>
                  <Input
                    placeholder=""
                    value={homeStore.params.review_score_range.max}
                    onChange={(event) => {
                      const _value = event.target.value;
                      homeStore.SET_PROVIDERSPARAM({
                        ...homeStore.params,
                        page: 1,
                        review_score_range: {
                          min: homeStore.params.review_score_range.min,
                          max: _value,
                        },
                      });
                    }}
                  />
                  <div className="unit">(Min: 0, Max: 5.0)</div>
                </div>
                <div className="action">
                  <div className="clear" onClick={clear}>
                    Clear
                  </div>
                  <div className="confirm" onClick={confirm}>
                    Confirm
                  </div>
                </div>
              </div>
            }
          >
            <div className="button" onClick={openCondition}>
              <div className="text">Filter</div>
              <div className="sda">
                <DownOutlined />
              </div>
            </div>
          </Popover>
        </div>
      </div>
    </div>
  );
};

export default observer(Filter);
