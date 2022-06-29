import styles from './Info.module.scss';
import { observer } from 'mobx-react';
import email from '@/assets/images/storageProvide/email.png';
import location from '@/assets/images/storageProvide/location.png';
import Icon_4 from '@/assets/images/storageProvide/icon_4.png';
import Icon_5 from '@/assets/images/storageProvide/icon_5.png';
import Icon_6 from '@/assets/images/storageProvide/icon_6.png';
import edit from '@/assets/images/client/edit.png';
import addChat from '@/assets/images/add_chat.png';
import clientStore from '@/store/modules/client';
import { Cookies } from '@/data/cookies';
import { useEffect, useState } from 'react';
import HashImage from '@/components/HashImage';
import default_avatar from '@/assets/images/default_avatar.png';
import default_logo from '@/assets/images/default_logo.png';
import { UserType } from '@/utils/enum';
import { useParams } from 'react-router';
import { urlInterface } from '@/utils/interface';
import userStore from '@/store/modules/user';
import { Storages } from '@/data/storages';
import { ChatContact } from '@/socket/ws.types';
import chatStore from '@/store/modules/chat';

const arr3 = [
  {
    icon: email,
    key: 'contact_email',
  },
  {
    icon: location,
    key: 'location',
  },
];

const arr4 = [
  {
    icon: Icon_4,
    key: 'slack',
  },
  {
    icon: Icon_5,
    key: 'twitter',
  },
  {
    icon: Icon_6,
    key: 'github',
  },
];
const Info = (props: any) => {
  const _url: urlInterface = useParams();
  const [editable, setEditable] = useState<Boolean>(false);
  useEffect(() => {
    const userInfo = Cookies.getUserInfo();
    const _addressId = _url.id;
    if (userInfo) {
      const { address_id, type } = userInfo;
      if (_addressId === address_id && type === UserType.client) {
        setEditable(true);
      }
    }
  }, [_url.id]);

  const editInfo = () => {
    clientStore.SET_PROFILEVISIBLE(true);
  };
  const editNode = () => {
    if (editable) {
      return (
        <div className={styles.editWrap} onClick={() => editInfo()}>
          <img src={edit} alt="" />
        </div>
      );
    } else {
      return <></>;
    }
  };

  const onChatButton = () => {
    let data = Storages.getPartners();
    const filterData = data.filter((value) => {
      return value.uid === clientStore.clientProfile.uid;
    });
    const contact = {
      ...clientStore.clientProfile,
    } as ChatContact;
    if (filterData.length === 0) {
      data.unshift(contact);
    }
    chatStore.SET_PARTNER_DATA(data);
    chatStore.SET_CURRENT_CONTACT({
      ...clientStore.clientProfile,
    } as ChatContact);
    chatStore.SET_SHOW(true);
  };

  return (
    <div className={styles.info}>
      <div className={styles.top}>
        <div className={styles.imageWrap} onClick={editInfo}>
          <HashImage
            placeholder={default_avatar}
            hash={clientStore.clientProfile.avatar}
          />
        </div>
        <div className={styles.editRow}>
          <div className={styles.avatarWrap}>
            <HashImage
              placeholder={default_logo}
              hash={clientStore.clientProfile.logo}
            />
          </div>
          {editNode()}
        </div>
      </div>
      <div className={styles.content}>
        <div className={styles.left}>
          <div className={styles.titleChat}>
            <div className={styles.titleWrap}>
              <div className={styles.title}>
                {clientStore.clientProfile.name === ''
                  ? 'Client'
                  : clientStore.clientProfile.name}
              </div>
              <div className={styles.subTitle}>
                {clientStore.clientProfile.address_id}
              </div>
            </div>
            <div
              className={styles.addChat}
              onClick={onChatButton}
              style={{
                display:
                  clientStore.clientProfile.uid === userStore.userInfo?.uid
                    ? 'none'
                    : undefined,
              }}
            >
              <img src={addChat} alt="" />
            </div>
          </div>
          <div className={styles.describe}>
            {clientStore.clientProfile.description}
          </div>
          <div className={styles.iconTextlist}>
            {arr3.map((item, index) => {
              return (
                <div className={styles.item} key={'arr2' + index}>
                  <div className={styles.iconWrap}>
                    <img src={item.icon} alt="" />
                  </div>
                  <div className={styles.text}>
                    {clientStore.clientProfile[item.key] === ''
                      ? '-'
                      : clientStore.clientProfile[item.key]}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
        <div className={styles.right}>
          <div className={styles.iconTextlist}>
            {arr4.map((item, index) => {
              return (
                <div className={styles.item} key={'arr5' + index}>
                  <div className={styles.iconWrap}>
                    <img src={item.icon} alt="" />
                  </div>
                  <div className={styles.count}>
                    {clientStore.clientProfile[item.key] === ''
                      ? '-'
                      : clientStore.clientProfile[item.key]}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
};

export default observer(Info);
