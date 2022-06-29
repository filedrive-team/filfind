import styles from './Info.module.scss';
import { observer } from 'mobx-react';
import { useEffect, useState } from 'react';
import email from '@/assets/images/storageProvide/email.png';
import location from '@/assets/images/storageProvide/location.png';
import Icon_1 from '@/assets/images/storageProvide/icon_1.png';
import Icon_2 from '@/assets/images/storageProvide/icon_2.png';
import Icon_3 from '@/assets/images/storageProvide/icon_3.png';
import Icon_4 from '@/assets/images/storageProvide/icon_4.png';
import Icon_5 from '@/assets/images/storageProvide/icon_5.png';
import Icon_6 from '@/assets/images/storageProvide/icon_6.png';
import edit from '@/assets/images/client/edit.png';
import addChat from '@/assets/images/add_chat.png';
import HashImage from '@/components/HashImage';
import default_avatar from '@/assets/images/default_avatar.png';
import default_logo from '@/assets/images/default_logo.png';
import { Cookies } from '@/data/cookies';
import storageProvideStore from '@/store/modules/storageProvide';
import { UserType } from '@/utils/enum';
import { urlInterface } from '@/utils/interface';
import { useParams } from 'react-router';
import { Storages } from '@/data/storages';
import { ChatContact } from '@/socket/ws.types';
import chatStore from '@/store/modules/chat';
import userStore from '@/store/modules/user';

const arr1 = [
  {
    icon: Icon_1,
    key: 'reputation_score',
  },
  {
    icon: Icon_2,
    key: 'review_score',
  },
  {
    icon: Icon_3,
    key: 'reviews',
  },
];
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
  const [editable, setEditable] = useState<boolean>(false);
  const [owner, setOwner] = useState<string>('');
  const _url: urlInterface = useParams();
  useEffect(() => {
    const userInfo = Cookies.getUserInfo();
    const _addressId = _url.id;
    setOwner(_addressId);
    if (userInfo) {
      const { address_id, type } = userInfo;
      if (_addressId === address_id && type === UserType.sp) {
        setEditable(true);
      }
    }
  }, [_url.id]);

  const editProfile = () => {
    storageProvideStore.SET_PROFILEVISIBLE(true);
  };
  const editNode = () => {
    if (editable) {
      return (
        <div className={styles.editWrap} onClick={() => editProfile()}>
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
      return value.uid === storageProvideStore.spOwnerProfile.uid;
    });
    const contact = {
      ...storageProvideStore.spOwnerProfile,
    } as ChatContact;
    if (filterData.length === 0) {
      data.unshift(contact);
    }
    chatStore.SET_PARTNER_DATA(data);
    chatStore.SET_CURRENT_CONTACT({
      ...storageProvideStore.spOwnerProfile,
    } as ChatContact);
    chatStore.SET_SHOW(true);
  };

  return (
    <div className={styles.info}>
      <div className={styles.top}>
        <div className={styles.logoWrap}>
          <HashImage
            placeholder={default_avatar}
            hash={storageProvideStore.spOwnerProfile.avatar}
          />
        </div>
        <div className={styles.editRow}>
          <div className={styles.avatarWrap}>
            <HashImage
              placeholder={default_logo}
              hash={storageProvideStore.spOwnerProfile.logo}
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
                {storageProvideStore.spOwnerProfile.name === ''
                  ? 'Storage Provider'
                  : storageProvideStore.spOwnerProfile.name}
              </div>
              <div className={styles.subTitle}>Owner:{' ' + owner}</div>
            </div>
            <div
              className={styles.addChat}
              onClick={onChatButton}
              style={{
                display:
                  storageProvideStore.spOwnerProfile.uid === '' ||
                  storageProvideStore.spOwnerProfile.uid ===
                    userStore.userInfo?.uid
                    ? 'none'
                    : undefined,
              }}
            >
              <img src={addChat} alt="" />
            </div>
          </div>

          <div className={styles.describe}>
            {storageProvideStore.spOwnerProfile.description === ''
              ? '-'
              : storageProvideStore.spOwnerProfile.description}
          </div>
          <div className={styles.iconTextlist}>
            {arr3.map((item, index) => {
              return (
                <div className={styles.item} key={'arr3' + index}>
                  <div className={styles.iconWrap}>
                    <img src={item.icon} alt="" />
                  </div>
                  <div className={styles.text}>
                    {storageProvideStore.spOwnerProfile[item.key] === ''
                      ? '-'
                      : storageProvideStore.spOwnerProfile[item.key]}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
        <div className={styles.right}>
          <div className={styles.list}>
            {arr1.map((item, index) => {
              return (
                <div className={styles.item} key={'arr1' + index}>
                  <div className={styles.iconWrap}>
                    <img src={item.icon} alt="" />
                  </div>
                  <div className={styles.count}>
                    {storageProvideStore.spOwnerProfile[item.key]}
                  </div>
                </div>
              );
            })}
          </div>

          <div className={styles.iconTextlist}>
            {arr4.map((item, index) => {
              return (
                <div className={styles.item} key={'arr4' + index}>
                  <div className={styles.iconWrap}>
                    <img src={item.icon} alt="" />
                  </div>
                  <div className={styles.count}>
                    {storageProvideStore.spOwnerProfile[item.key] === ''
                      ? '-'
                      : storageProvideStore.spOwnerProfile[item.key]}
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
