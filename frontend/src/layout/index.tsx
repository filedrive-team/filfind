import { Redirect, Route } from 'react-router-dom';
import styles from './style.module.scss';
import { RouterPath } from '@/router/RouterConfig';
import classNames from 'classnames';
import { ReactComponent as ChatButtonIcon } from '@/assets/images/chat/chat_bottom.svg';
import Chat from '@/components/Chat';
import chatStore from '@/store/modules/chat';
import { observer } from 'mobx-react';
import userStore from '@/store/modules/user';
import { useHistory } from 'react-router';
import { useEffect } from 'react';
//TODO: replace with project variable ⭐️⭐️⭐️
const IS_LOGGED = false;

/**
 *
 * @param {component} props
 * @returns
 */
const Layout = (props: any) => {
  const { component: Com, auth, ...rest } = props;
  const history = useHistory();
  const onChatTap = () => {
    if (userStore.isLogin) {
      chatStore.SET_SHOW(!chatStore.show);
    } else {
      history.push(RouterPath.auth);
    }
  };

  useEffect(() => {
    if (userStore.isLogin === false) {
      chatStore.SET_SHOW(false);
    }
  }, [userStore.isLogin]);

  return (
    <Route
      {...rest}
      render={(props: any) => (
        <div className={classNames(styles.bodyWrap)}>
          {!auth || (IS_LOGGED && auth) ? (
            <div className={styles.body}>
              <Com {...props} />
            </div>
          ) : (
            <Redirect to={RouterPath.auth} />
          )}
          <Chat />
          {userStore.isLogin ? (
            <div
              className={classNames(styles.chatButtonWrap)}
              onClick={onChatTap}
            >
              <ChatButtonIcon />
            </div>
          ) : (
            <></>
          )}
        </div>
      )}
    />
  );
};

export default observer(Layout);
