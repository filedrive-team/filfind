import styles from './ClientReviews.module.scss';
import { Rate, Pagination } from 'antd';
import { observer } from 'mobx-react';
import Icon_7 from '@/assets/images/storageProvide/icon_7.png';
import Icon_8 from '@/assets/images/storageProvide/icon_8.png';
import dayjs from 'dayjs';

const ClientReviews = (props) => {
  const _listNode = () => {
    const _length = props.commentList.length;
    if (_length) {
      return props.commentList.map((item, index) => {
        return (
          <div className={styles.commentItem} key={index}>
            <div className={styles.name}>{item.title}</div>
            <div className={styles.info}>
              <div className={styles.cards}>
                <div className={styles.cardsItem}>{item.provider}</div>
                <div className={styles.cardsItem}>FIL+Deals</div>
              </div>
              <Rate allowHalf defaultValue={item.score} disabled />
            </div>
            <div className={styles.comment}>{item.content}</div>
            <div className={styles.commentInfo}>
              <div className={styles.item}>
                <div className={styles.iconWrap}>
                  <img src={Icon_7} alt="" />
                </div>
                <div className={styles.text}>{item.client}</div>
              </div>
              <div className={styles.item}>
                <div className={styles.iconWrap}>
                  <img src={Icon_8} alt="" />
                </div>
                <div className={styles.text}>
                  {dayjs.unix(item.created_at).format('YYYY-MM-DD')}
                </div>
              </div>
            </div>
          </div>
        );
      });
    } else {
      return <></>;
    }
  };
  return (
    <div className={styles.list}>
      <div className={styles.title}>{props.title}</div>
      <div className={styles.content}>
        <div className={styles.commentList}>{_listNode()}</div>
        <Pagination
          showSizeChanger
          onChange={props.onChange}
          total={props.total}
        />
      </div>
    </div>
  );
};

export default observer(ClientReviews);
