import React from 'react';
import style from './style.module.css';
import MsgBubble, { TMessage } from '@/components/Chat/MsgBubble';
import dayjs from 'dayjs';
import classNames from 'classnames';
import { ChatContact, ChatMessage } from '@/socket/ws.types';
import HashImage from '@/components/HashImage';
import avatar from '@/assets/images/default_avatar.png';

interface IProps {
  data: ChatMessage;
  me?: ChatContact;
  other?: ChatContact;
}

export default function MsgItem({ data, me, other }: IProps) {
  const isMe = data.sender === me?.uid;

  return (
    <div
      className={classNames(style.content, 'flex')}
      style={{ flexDirection: isMe ? 'row-reverse' : 'row' }}
    >
      <div className={style.avatar}>
        <HashImage
          placeholder={avatar}
          hash={isMe ? me?.avatar : other?.avatar}
        />
      </div>
      <div
        className={style.text_area}
        style={{ alignItems: isMe ? 'flex-end' : 'flex-start' }}
      >
        <MsgBubble isMe={isMe} data={data} />
        <div className={style.comment_area}>
          <span className={style.date_text}>
            {dayjs.unix(data.timestamp).format('MM-DD HH:mm')}
          </span>
        </div>
      </div>
    </div>
  );
}
