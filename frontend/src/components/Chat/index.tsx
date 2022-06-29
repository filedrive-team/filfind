import styles from './style.module.scss';
import classNames from 'classnames';
import chatStore from '@/store/modules/chat';
import { observer } from 'mobx-react';
import { useEffect, useState } from 'react';

import ChatLeft from '@/components/Chat/ChatLeft';
import ChatHead from '@/components/Chat/ChatHead';
import ChatFooter from '@/components/Chat/ChatFooter';
import ChatRecordList from '@/components/Chat/ChatRecordList';
import ScrollWrapper from '@/components/Chat/ScrollWrapper';
import ws from '@/socket/ws';
import { ChannelTypeEnum, ChatContact, ChatMessage } from '@/socket/ws.types';
import { Cookies } from '@/data/cookies';

const WrappedChatRecordList = ScrollWrapper(ChatRecordList);

const Chat = () => {
  const userInfo = Cookies.getUserInfo();

  let chatState = useState(false);
  const [data, setData] = useState<ChatMessage[]>([]);
  const [me, setMe] = useState<ChatContact>();
  useEffect(() => {
    if (userInfo) {
      setMe({
        ...userInfo,
        address_id: userInfo.address,
        contact_email: userInfo.email,
      });
    }
  }, []);

  useEffect(() => {
    if (chatStore.partnersList.length !== 0) {
      ws.send({
        channel: ChannelTypeEnum.ChatHistory,
        subscribe: true,
        body: {
          partner: chatStore.currentContact?.uid,
          limit: 10,
        },
      });

      ws.send({
        channel: ChannelTypeEnum.ChatReceive,
        subscribe: true,
        body: {},
      });

      ws.send({
        channel: ChannelTypeEnum.ChatUncheckedList,
        subscribe: true,
        body: {},
      });
    }
  }, [chatStore.partnersList, chatStore.currentContact, chatStore.isChatAuth]);

  return chatStore.show === false ? (
    <div></div>
  ) : (
    <div
      className={classNames(
        styles.chatWrap,
        !chatStore.show ? styles.hidden : chatState ? styles.show : '',
      )}
    >
      <div className={classNames(styles.chatContent)}>
        <div className={styles.chatSide}>
          <ChatLeft />
        </div>
        <div className={styles.chatRight}>
          <ChatHead />
          <div className={styles.dialogContent}>
            <WrappedChatRecordList
              scrollToBottom
              onIsAutoToScroll={(value) => {
                chatStore.SET_AUTO_SCROLL(value);
              }}
              isAutoToScroll={chatStore.isAutoScroll}
              data={chatStore.chatMessages}
              me={me}
              other={chatStore.currentContact}
              onEarlier={() => {
                chatStore.getChatHistory().then();
              }}
            ></WrappedChatRecordList>
          </div>
          <ChatFooter />
        </div>
      </div>
    </div>
  );
};

export default observer(Chat);
