import styles from './ChatLeft.module.scss';
import classNames from 'classnames';
import UserInfo from '@/components/Chat/UserInfo';
import { Menu, MenuProps } from 'antd';
import React from 'react';
import chatStore from '@/store/modules/chat';
import { observer } from 'mobx-react';
import { Cookies } from '@/data/cookies';
import { Storages } from '@/data/storages';
import HashImage from '@/components/HashImage';
import avatar from '@/assets/images/default_avatar.png';
import ws from '@/socket/ws';
import { ChannelTypeEnum } from '@/socket/ws.types';

const ChatLeft = () => {
  const uncheckCount = (uid: string) => {
    const value = chatStore.chatUncheckedData.filter(
      (value) => value.partner === uid,
    );
    return value.length === 0 ? 0 : value[0].number;
  };

  const lastMessage = (uid: string): string => {
    const value = Storages.getChatMessage(uid);
    // @ts-ignore
    return value.length === 0 ? '-' : value.reverse()[0].content;
  };

  const items: MenuProps['items'] = chatStore.partnersList.map(
    (value, index) => {
      return {
        label: (
          <>
            <div className={'friend-avatar'}>
              <HashImage placeholder={avatar} hash={value.avatar} />
            </div>
            <div className={'content'}>
              <div className={'name'}>{value.name}</div>
              <div className={'speak'}>{lastMessage(value.uid)}</div>
            </div>
            <div
              className={'count'}
              style={{ display: uncheckCount(value.uid) === 0 ? 'none' : '' }}
            >
              {uncheckCount(value.uid)}
            </div>
          </>
        ),
        key: value.uid,
      };
    },
  );

  const onItemClick: MenuProps['onClick'] = (e) => {
    const userInfo = Cookies.getUserInfo();
    if (!userInfo) return;
    const currentContact = chatStore.partnersList.filter((value, index) => {
      return value.uid === e.key;
    });
    if (currentContact.length !== 0) {
      ws.send({
        channel: ChannelTypeEnum.ChatCheck,
        subscribe: true,
        body: {
          partner: currentContact[0].uid,
        },
      });
      chatStore.SET_CURRENT_CONTACT(currentContact[0]);
    }
  };

  return (
    <div className={classNames(styles.chatLeftWrap)}>
      <UserInfo />
      <div className={'mt-26'} />
      {chatStore.currentContact ? (
        <div>
          <div className={styles.listWrap}>
            <Menu
              mode="inline"
              items={items}
              defaultSelectedKeys={[chatStore.currentContact.uid]}
              onClick={onItemClick}
            ></Menu>
          </div>
        </div>
      ) : (
        <></>
      )}
    </div>
  );
};

export default observer(ChatLeft);
