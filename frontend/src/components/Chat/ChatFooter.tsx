import styles from './ChatFooter.module.scss';
import { Input } from 'antd';
import { ReactComponent as SendIcon } from '@/assets/images/chat/send.svg';
import { ReactComponent as DisableSendIcon } from '@/assets/images/chat/disable_send.svg';

import {
  ChangeEvent,
  ChangeEventHandler,
  KeyboardEventHandler,
  useState,
} from 'react';
import dayjs from 'dayjs';
import classNames from 'classnames';
import chatStore from '@/store/modules/chat';
import { Cookies } from '@/data/cookies';
import { ChannelTypeEnum, ChatMessage } from '@/socket/ws.types';
import ws from '@/socket/ws';
import { observer } from 'mobx-react';

const { TextArea } = Input;

const ChatFooter = () => {
  const [text, setText] = useState('');
  const [isAllowSend, setIsAllowSend] = useState(false);
  const [isShift, setIsShift] = useState(false);
  const userInfo = Cookies.getUserInfo();

  const textChangeHandle: ChangeEventHandler = (
    e: ChangeEvent<HTMLInputElement>,
  ) => {
    const isAllowSend = !!e.target.value.trim();
    const text = e.target.value;
    setText(text);
    setIsAllowSend(isAllowSend);
  };

  const keyDownHandle: KeyboardEventHandler = (e) => {
    if (e.keyCode === 16) {
      setIsShift(true);
    }

    if (e.keyCode === 13 && !isShift) {
      e.preventDefault();
      sendHandle();
    }
  };

  const keyUpHandle: KeyboardEventHandler = (e) => {
    if (e.keyCode === 16) {
      setIsShift(false);
    }
  };

  const sendHandle = () => {
    if (!isAllowSend) {
      return;
    }

    chatStore.SET_AUTO_SCROLL(true);

    const timestamp = dayjs().unix();
    const message: ChatMessage = {
      recipient: chatStore.currentContact?.uid,
      sender: userInfo?.uid,
      type: 0,
      content: text,
      checked: true,
      timestamp,
    };
    ws.send({
      channel: ChannelTypeEnum.ChatSend,
      subscribe: true,
      body: message,
    });
    resetText();
  };

  const resetText = () => {
    setText('');
    setIsAllowSend(false);
  };

  return (
    <div className={styles.chatFooter}>
      <TextArea
        placeholder={'Input content…'}
        bordered={false}
        autoSize={{ minRows: 1, maxRows: 3 }}
        onChange={textChangeHandle}
        onKeyUp={keyUpHandle}
        onKeyDown={keyDownHandle}
        value={text}
      />
      <div
        className={classNames(
          styles.sendButton,
          isAllowSend ? styles.allow : '',
        )}
        onClick={sendHandle}
      >
        {isAllowSend ? <SendIcon /> : <DisableSendIcon />}
      </div>
    </div>
  );
};

export default observer(ChatFooter);
