import clientStore from '@/store/modules/client';
import styles from './Detail.module.scss';
import { observer } from 'mobx-react';
import { Cookies } from '@/data/cookies';
import edit from '@/assets/images/client/edit.png';
import { useEffect, useState } from 'react';
import { bytesToSize } from '@/utils';
import { UserType } from '@/utils/enum';
import { useParams } from 'react-router';
import { urlInterface } from '@/utils/interface';

const Detail = () => {
  const _params: urlInterface = useParams();
  const [editable, setEditable] = useState<boolean>(false);
  useEffect(() => {
    const userInfo = Cookies.getUserInfo();
    const _addressId = _params.id;
    if (userInfo) {
      const { address_id, type } = userInfo;
      if (_addressId === address_id && type === UserType.client) {
        setEditable(true);
      }
    }
  }, [_params.id]);

  const editInfo = () => {
    clientStore.SET_INFOVISIBLE(true);
  };

  const editNode = () => {
    if (editable) {
      return (
        <div className={styles.editWrap} onClick={() => editInfo()}>
          <img src={edit} alt="" />
        </div>
      );
    }
  };

  const requirementNode = () => {
    if (clientStore.clientDetail.service_requirement) {
      return (
        <>
          <div className={styles.title}>Storage Service Requirement</div>
          <div className={styles.requirement}>
            {clientStore.clientDetail.service_requirement}
          </div>
        </>
      );
    } else {
      return <></>;
    }
  };
  return (
    <div className={styles.detail}>
      <div className={styles.title}>
        <div className={styles.text}>Client Detailed Info</div>
        {editNode()}
      </div>
      <div className={styles.list}>
        <div className={styles.item}>
          ·Monthly storage capacity：
          {clientStore.clientDetail.monthly_storage === ''
            ? '-'
            : clientStore.clientDetail.monthly_storage}
        </div>
        <div className={styles.item}>
          ·Data transfer bandwidth：
          {clientStore.clientDetail.bandwidth === ''
            ? '-'
            : clientStore.clientDetail.bandwidth}
        </div>
        <div className={styles.item}>
          ·Available DataCap：
          {clientStore.clientDetail.data_cap === ''
            ? '-'
            : bytesToSize(clientStore.clientDetail.data_cap)}
        </div>
        <div className={styles.item}>
          ·Data use case：
          {clientStore.clientDetail.use_case === ''
            ? '-'
            : clientStore.clientDetail.use_case}
        </div>
      </div>

      {requirementNode()}
    </div>
  );
};

export default observer(Detail);
