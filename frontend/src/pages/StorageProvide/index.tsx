import CVHeader from '@/components/CVHeader';
import Footer from '@/components/Footer';
import styles from './style.module.scss';
import Info from '@/components/StorageProvide/Info';
import List from '@/components/StorageProvide/List';
import Rank from '@/components/StorageProvide/Rank';
import ClientReviews from '@/components/StorageProvide/ClientReviews';
import AddNewReveiews from '@/components/StorageProvide/AddNewReveiews';
import storageProvideStore from '@/store/modules/storageProvide';
import { useEffect } from 'react';
import { observer } from 'mobx-react';
import headerStore from '@/store/modules/header';
import EditProfile from '@/components/StorageProvide/EditProfile';
import EditInfo from '@/components/StorageProvide/EditInfo';
import { UserType } from '@/utils/enum';
import { Cookies } from '@/data/cookies';
import { useParams } from 'react-router';
import { urlInterface } from '@/utils/interface';

const StorageProvide = (props: any) => {
  const _params: urlInterface = useParams();
  const userInfo = Cookies.getUserInfo();
  const _addressId = _params.id;
  if (userInfo) {
    const { address_id } = userInfo;
    if (_addressId === address_id) {
      headerStore.SET_ACTIVEINDEX(2);
    } else {
      headerStore.SET_ACTIVEINDEX(-1);
    }
  }
  useEffect(() => {
    const address_id = _params.id;
    storageProvideStore._getSpOwnerProfile({ address_id: address_id });
    storageProvideStore.SET_SPOWNERREVIEWSPARAM({
      page: 1,
      page_size: 10,
      address_id: address_id,
    });
    storageProvideStore._getSpOwnerReviews(storageProvideStore.reviewsParam);
    storageProvideStore._getSpServiceDetail({ address_id: address_id });
  }, [_params.id]);

  const addNewReveiewsNode = () => {
    const _userInfo = Cookies.getUserInfo();
    if (_userInfo?.type === UserType.client) {
      return (
        <div className={styles.addNewReveiewsWrap}>
          <AddNewReveiews></AddNewReveiews>
        </div>
      );
    } else {
      return <></>;
    }
  };

  const clientReviewsNode = () => {
    if (storageProvideStore.reviewsList.length) {
      return (
        <div className={styles.clientReviewsWrap}>
          <ClientReviews
            title={'Client Reviews'}
            commentList={storageProvideStore.reviewsList}
            onChange={(page, pageSize) => {
              storageProvideStore.SET_SPOWNERREVIEWSPARAM({
                ...storageProvideStore.reviewsParam,
                page: 1,
                page_size: 10,
              });
              storageProvideStore._getSpOwnerReviews(
                storageProvideStore.reviewsParam,
              );
            }}
            total={storageProvideStore.reviewsTotal}
          ></ClientReviews>
        </div>
      );
    } else {
      return <></>;
    }
  };

  return (
    <div className={styles.storageProvide}>
      <div className={styles.top}>
        <CVHeader></CVHeader>
      </div>
      <div className={styles.infoWrap}>
        <Info></Info>
      </div>
      <div className={styles.listWrap}>
        <List></List>
      </div>
      <div className={styles.rankWrap}>
        <Rank></Rank>
      </div>
      {clientReviewsNode()}
      {addNewReveiewsNode()}
      <EditProfile></EditProfile>
      <EditInfo></EditInfo>
      <Footer></Footer>
    </div>
  );
};

export default observer(StorageProvide);
