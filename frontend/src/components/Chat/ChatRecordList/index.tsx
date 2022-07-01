import React, { CSSProperties, MouseEventHandler } from 'react';
import style from './style.module.css';
import classNames from 'classnames';
import MsgItem from '@/components/Chat/MsgItem';
import { ChatContact } from '@/socket/ws.types';

interface IProps {
  onEarlier?: MouseEventHandler;
  data: any[];
  me?: ChatContact;
  other?: ChatContact;
  style?: CSSProperties;
}

const ChatRecordList = (props: IProps) => {
  return (
    <div className={classNames([style.list_area])}>
      <div>
        <button className={style.load_more} onClick={props.onEarlier}>
          Load more···
        </button>
      </div>
      {props.data.map((bubble, index) => (
        <MsgItem {...props} data={bubble} key={index} />
      ))}
    </div>
  );
};

export default ChatRecordList;
