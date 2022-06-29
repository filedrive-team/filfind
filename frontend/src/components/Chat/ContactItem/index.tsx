import React, { Component, CSSProperties } from 'react';
import PropTypes from 'prop-types';
import style from './style.module.css';
import classNames from 'classnames';

interface IProps {
  styles?: CSSProperties;
  selected: boolean;
  border: boolean;
  contact: any;
  onClick: Function;
}

const ContactItem = ({
  styles,
  selected,
  border,
  contact,
  onClick,
}: IProps) => {
  return (
    <div
      style={styles}
      className={classNames([
        style.content,
        border && style.bottom_border,
        selected && style.selected,
      ])}
      onClick={onClick.bind(null, contact)}
    >
      <img className={style.icon} src={contact.avatar} />
      <div className={style.info_area}>
        <span className={`${style.nickname} ${style.ellipsis}`}>
          {contact.nickname}
        </span>
        <span className={`${style.desc} ${style.ellipsis}`}>
          {contact.message}
        </span>
      </div>
      <span className={style.date_area}>{contact.date}</span>
    </div>
  );
};

export default ContactItem;
