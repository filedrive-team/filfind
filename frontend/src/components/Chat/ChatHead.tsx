import styles from './ChatHead.module.scss';
import classNames from 'classnames';
import { ReactComponent as CloseButtonIcon } from '@/assets/images/chat/close_black.svg';
import chatStore from '@/store/modules/chat';

const ChatHead = () => {
  const onChatTap = () => {
    chatStore.SET_SHOW(!chatStore.show);
  };
  return (
    <div className={styles.chatHeader}>
      <div className={classNames(styles.name)}>
        {chatStore.currentContact?.name}
      </div>
      <div onClick={onChatTap}>
        <CloseButtonIcon />
      </div>
    </div>
  );
};

export default ChatHead;
