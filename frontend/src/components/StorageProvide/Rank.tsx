import storageProvideStore from '@/store/modules/storageProvide';
import styles from './Rank.module.scss';
import { observer } from 'mobx-react';
import edit from '@/assets/images/client/edit.png';
import { spServiceDetailResponse } from '@/api/types';
import { useEffect, useState } from 'react';
import { Cookies } from '@/data/cookies';
import { UserType } from '@/utils/enum';
import { urlInterface } from '@/utils/interface';
import { useParams } from 'react-router';
const labels = [
  {
    label: '·Available FIL+ deals',
    key: 'available_deals',
  },
  {
    label: '·Data transfer bandwidth',
    key: 'bandwidth',
  },
  {
    label: '·Sealing speed',
    key: 'sealing_speed',
  },
  {
    label: '·Renewable energy',
    key: 'renewable_energy',
  },
  {
    label: '·Number of parallel deals',
    key: 'parallel_deals',
  },
  {
    label: '·Certification',
    key: 'certification',
  },
  {
    label: '·Member of Enterprise SP Group',
    key: 'is_member',
  },
  {
    label: '·Onboarding Products Experience',
    key: 'experience',
  },
];
const Rank = () => {
  const _url: urlInterface = useParams();
  const [editable, setEditable] = useState<Boolean>(false);
  useEffect(() => {
    const userInfo = Cookies.getUserInfo();
    const _addressId = _url.id;
    if (userInfo) {
      const { address_id, type } = userInfo;
      if (_addressId === address_id && type === UserType.sp) {
        setEditable(true);
      }
    }
  }, [_url.id]);
  const editNode = (item) => {
    if (editable) {
      return (
        <div className={styles.editWrap} onClick={() => editInfo(item)}>
          <img src={edit} alt="" />
        </div>
      );
    } else {
      return <></>;
    }
  };

  const editInfo = (item) => {
    const {
      address,
      available_deals,
      bandwidth,
      certification,
      experience,
      is_member,
      parallel_deals,
      renewable_energy,
      sealing_speed,
    } = item;

    storageProvideStore.SET_PROVIDERDETAIL({
      address,
      available_deals,
      bandwidth,
      certification,
      experience,
      is_member,
      parallel_deals,
      renewable_energy,
      sealing_speed,
    });
    storageProvideStore.SET_INFOVISIBLE(true);
  };

  if (storageProvideStore.spServiceDetailList.length) {
    return (
      <div className={styles.rank}>
        {storageProvideStore.spServiceDetailList.map(
          (item: spServiceDetailResponse, index) => {
            return (
              <div className={styles.rankItem} key={'rank' + index}>
                <div className={styles.rankTitleRow}>
                  <div className={styles.rankAddress}>
                    <div className={styles.no}>Node{index + 1}:</div>
                    <div className={styles.address}>{item.address}</div>
                  </div>
                  {editNode(item)}
                </div>
                <div className={styles.list}>
                  {labels.map((subItem, subIndex) => {
                    return (
                      <div className={styles.item} key={'sub' + subIndex}>
                        <div className={styles.itemLable}>
                          {subItem.label}：
                        </div>
                        <div className={styles.itemValue}>
                          {item[subItem.key] === '' ? '-' : item[subItem.key]}
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            );
          },
        )}
      </div>
    );
  } else {
    return <></>;
  }
};

export default observer(Rank);
