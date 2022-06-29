import styles from './Filter.module.scss';
import { Input } from 'antd';
import { observer } from 'mobx-react';
import clientListStore from '@/store/modules/clientList';
import filFind from '@/assets/images/home/filFind.png';
const { Search } = Input;
const subTitle =
  'A discovery platform for SPs and Clients in Filecoin Ecosystem';

const onSearch = (value: string) => {
  clientListStore.SET_CLIENTSPARAM({
    ...clientListStore.params,
    page: 1,
    search: value,
  });
  clientListStore._getClients(clientListStore.params);
};
const Filter = () => {
  return (
    <div className={styles.filter}>
      <div className={styles.filterContent}>
        <div className={styles.filFind}>
          <img src={filFind} alt="" />
        </div>
        <div className={styles.subTitle}>{subTitle}</div>
        <div className={styles.searchRow}>
          <Search
            placeholder="Client ID/Name/Location"
            enterButton="Search"
            onSearch={onSearch}
          />
        </div>
      </div>
    </div>
  );
};

export default observer(Filter);
