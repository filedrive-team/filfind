import React from 'react';
import PropTypes from 'prop-types';
import style from './style.module.css';
import classNames from 'classnames';
import { ChatContact, ChatMessage } from '@/socket/ws.types';

type MessageType = 'text' | 'image';

export type TPureMsg = {
  type: MessageType;
  content: string;
};

export type TMessage = {
  _id: string;
  date: number;
  user: ChatContact;
  message: TPureMsg;
};

interface IProps {
  data: ChatMessage;
  isMe: boolean;
}

const MsgBubble = ({ data, isMe }: IProps) => {
  const renderContent = (message: ChatMessage) => {
    switch (message.type) {
      case 0:
        return message.content;
      case 1:
        return (
          <img
            className={classNames(style.img_content)}
            src={message.content}
          />
        );
      default:
        break;
    }
  };

  return (
    <div
      className={classNames(
        style.text_content,
        isMe ? style.arrow_right : style.arrow_left,
        style.arrow,
      )}
    >
      {renderContent(data)}
    </div>
  );
};

export default MsgBubble;
