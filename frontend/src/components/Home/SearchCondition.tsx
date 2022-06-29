import styles from './SearchCondition.module.scss';
import classnames from 'classnames';
import { RegionType, SPSStatusType } from '@/utils/enum';
import homeStore from '@/store/modules/home';
import { observer } from 'mobx-react';
const statusList = [
  {
    key: SPSStatusType.All,
    text: 'All',
  },
  {
    key: SPSStatusType.Registered,
    text: 'Registered',
  },
  {
    key: SPSStatusType.Autofilled,
    text: 'Autofilled',
  },
];
const regionList = [
  {
    key: RegionType.All,
    text: 'All',
  },
  {
    key: RegionType.Africa,
    text: 'Africa',
  },
  {
    key: RegionType.Asia,
    text: 'Asia',
  },
  {
    key: RegionType.Europe,
    text: 'Europe',
  },
  {
    key: RegionType.NorthAmerica,
    text: 'North America',
  },
  {
    key: RegionType.Oceania,
    text: 'Oceania',
  },
  {
    key: RegionType.SouthAmerica,
    text: 'South America',
  },
];
const SearchCondition = () => {
  const selectStatus = (status: SPSStatusType) => {
    homeStore.SET_PROVIDERSPARAM({
      ...homeStore.params,
      page: 1,
      sps_status: status,
    });
    homeStore._getProviders(homeStore.params);
  };

  const selectRegion = (region: RegionType) => {
    homeStore.SET_PROVIDERSPARAM({
      ...homeStore.params,
      page: 1,
      region: region,
    });
    homeStore._getProviders(homeStore.params);
  };

  return (
    <div className={styles.searchCondition}>
      <div className={styles.statusRow}>
        <div className={styles.label}>SPS Status：</div>
        <div className={styles.list}>
          {statusList.map((item, index) => {
            return (
              <div
                className={classnames(
                  styles.item,
                  homeStore.params.sps_status === item.key ? styles.active : '',
                )}
                key={index}
                onClick={() => {
                  selectStatus(item.key);
                }}
              >
                {item.text}
              </div>
            );
          })}
        </div>
      </div>
      <div className={styles.regionRow}>
        <div className={styles.label}>Region：</div>
        <div className={styles.list}>
          {regionList.map((item, index) => {
            return (
              <div
                className={classnames(
                  styles.item,
                  homeStore.params.region === item.key ? styles.active : '',
                )}
                onClick={() => {
                  selectRegion(item.key);
                }}
                key={index}
              >
                {item.text}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};

export default observer(SearchCondition);
