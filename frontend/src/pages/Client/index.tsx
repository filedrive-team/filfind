import CVHeader from '@/components/CVHeader';
import Footer from '@/components/Footer';
import styles from './style.module.scss';
import Info from '@/components/Client/Info';
import List from '@/components/Client/List';
import Detail from '@/components/Client/Detail';
import ClientReviews from '@/components/StorageProvide/ClientReviews';
import EditProfile from '@/components/Client/EditProfile';
import EditInfo from '@/components/Client/EditInfo';
import { useEffect } from 'react';
import { observer } from 'mobx-react';
import clientStore from '@/store/modules/client';
import headerStore from '@/store/modules/header';
import { urlInterface } from '@/utils/interface';
import { useParams } from 'react-router';
import { Cookies } from '@/data/cookies';
const Client = (props: any) => {
  const _url: urlInterface = useParams();
  const userInfo = Cookies.getUserInfo();
  const _addressId = _url.id;
  if (userInfo) {
    const { address_id } = userInfo;
    if (_addressId === address_id) {
      headerStore.SET_ACTIVEINDEX(2);
    } else {
      headerStore.SET_ACTIVEINDEX(-1);
    }
  }
  useEffect(() => {
    const address_id = _url.id;
    clientStore._getClientDetail({ address_id: address_id });
    clientStore.SET_HISTORYDEALSTATSPARAM({
      page: 1,
      page_size: 10,
      address_id: address_id,
    });
    clientStore._getClientHistoryDealStats(clientStore.historyDealStatsParam);
    clientStore._getClientProfile({ address_id: address_id });
    clientStore.SET_REVIEWSPARAM({
      page: 1,
      page_size: 10,
      address_id: address_id,
    });
    clientStore._getClientReviews(clientStore.reviewsParam);
  }, [_url.id]);

  const clientReviewsNode = () => {
    if (clientStore.reviewsList.length) {
      return (
        <div className={styles.clientReviewsWrap}>
          <ClientReviews
            title={'Client Reviews'}
            commentList={clientStore.reviewsList}
            onChange={(page, pageSize) => {
              clientStore.SET_REVIEWSPARAM({
                ...clientStore.reviewsParam,
                page: 1,
                page_size: 10,
              });
              clientStore._getClientReviews(clientStore.reviewsParam);
            }}
            total={clientStore.reviewsTotal}
          ></ClientReviews>
        </div>
      );
    } else {
      return <></>;
    }
  };

  return (
    <div className={styles.client}>
      <div className={styles.top}>
        <CVHeader></CVHeader>
      </div>
      <div className={styles.infoWrap}>
        <Info></Info>
      </div>
      <div className={styles.DetailWrap}>
        <Detail></Detail>
      </div>
      <div className={styles.listWrap}>
        <List></List>
      </div>
      {clientReviewsNode()}
      <EditProfile></EditProfile>
      <EditInfo></EditInfo>
      <Footer></Footer>
    </div>
  );
};
export default observer(Client);
