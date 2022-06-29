import styles from './style.module.scss';
import logo from '@/assets/images/header/logo.png';
import { Menu, Dropdown } from 'antd';
import classnames from 'classnames';
import { useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import { RouterPath } from '@/router/RouterConfig';
import { UserType } from '@/utils/enum';
import { signUpUser } from '@/pages/Auth';
import home_white from '@/assets/images/header/icon-home-white.png';
import cv_white from '@/assets/images/header/icon-cv-white.png';
import message_white from '@/assets/images/header/icon-message-white.png';
import user_white from '@/assets/images/header/icon-user-white.png';
import home_blue from '@/assets/images/header/icon-home-blue.png';
import cv_blue from '@/assets/images/header/icon-cv-blue.png';
import message_blue from '@/assets/images/header/icon-message-blue.png';
import user_blue from '@/assets/images/header/icon-user-blue.png';
import home_black from '@/assets/images/header/icon-home-black.png';
import cv_black from '@/assets/images/header/icon-cv-black.png';
import message_black from '@/assets/images/header/icon-message-black.png';
import user_black from '@/assets/images/header/icon-user-black.png';
import default_avatar from '@/assets/images/default_avatar.png';
import client_dashboard_white from '@/assets/images/header/icon-clientDashboard-white.png';
import client_dashboard_blue from '@/assets/images/header/icon-clientDashboard-blue.png';
import client_dashboard_black from '@/assets/images/header/icon-clientDashboard-black.png';

import headerStore from '@/store/modules/header';
import { Cookies } from '@/data/cookies';
import HashImage from '@/components/HashImage';
import chatStore from '@/store/modules/chat';
import userStore from '@/store/modules/user';
import { DownOutlined } from '@ant-design/icons';

enum menuPage {
  messages = 'Messages',
  CV = 'CV',
}
const loginMenu = [
  {
    text: 'SPs Dashboard',
    url: RouterPath.home,
    white: home_white,
    blue: home_blue,
    black: home_black,
  },
  {
    text: 'Clients Dashboard',
    url: RouterPath.clientList,
    white: client_dashboard_white,
    blue: client_dashboard_blue,
    black: client_dashboard_black,
  },
  {
    text: 'CV',
    url: 'CV',
    white: cv_white,
    blue: cv_blue,
    black: cv_black,
  },
  {
    text: 'Messages',
    url: 'Messages',
    white: message_white,
    blue: message_blue,
    black: message_black,
  },
  {
    text: 'Account',
    url: RouterPath.changePassword,
    white: user_white,
    blue: user_blue,
    black: user_black,
  },
];

const unLoginMenu = [
  {
    text: 'SPs Dashboard',
    url: RouterPath.home,
    white: home_white,
    blue: home_blue,
    black: home_black,
  },
  {
    text: 'Clients Dashboard',
    url: RouterPath.clientList,
    white: client_dashboard_white,
    blue: client_dashboard_blue,
    black: client_dashboard_black,
  },
];

const CVHeader = (props: any) => {
  const access_token = Cookies.getAuthorization();
  const avatar = Cookies.getUserInfo()?.avatar ?? '';
  const history = useHistory();
  const [isFixed, setIsFixed] = useState<Boolean>(false);
  useEffect(() => {
    window.addEventListener('scroll', _onScroll);
    return () => {
      window.removeEventListener('scroll', _onScroll);
    };
  }, []);

  const _onScroll = () => {
    const _top = document.documentElement.scrollTop;
    if (_top > 80) {
      setIsFixed(true);
    } else {
      setIsFixed(false);
    }
  };

  const menuClick = (obj) => {
    const { key } = obj;
    history.push({
      pathname: RouterPath.signUp,
      state: { user: key },
    });
  };
  const _dropDown = (
    <Menu
      onClick={menuClick}
      items={[
        {
          key: signUpUser.Client,
          label: 'Client',
        },
        {
          key: signUpUser.SP,
          label: 'Storage Provider',
        },
      ]}
    />
  );

  const skipToSignIn = () => {
    history.push(RouterPath.auth);
  };
  const goHome = () => {
    history.push(RouterPath.home);
  };
  const logOff = () => {
    Cookies.removeUserInfo();
    history.push(RouterPath.home);
  };

  const skipTo = (item) => {
    console.log(item, 'iiiii');

    if (item.url === menuPage.messages) {
      if (userStore.isLogin) {
        chatStore.SET_SHOW(!chatStore.show);
      } else {
        history.push(RouterPath.auth);
      }
    } else if (item.url === menuPage.CV) {
      const userInfo = Cookies.getUserInfo();
      if (userInfo?.type === UserType.client) {
        window.location.href = `/client/${userInfo?.address_id}`;
      } else {
        window.location.href = `/storageProvide/${userInfo?.address_id}`;
      }
    } else {
      history.push(item.url);
    }
  };

  const getIcon = (item, isActive) => {
    if (isFixed) {
      if (isActive) {
        return <img src={item.blue} alt="" />;
      } else {
        return <img src={item.black} alt="" />;
      }
    } else {
      return <img src={item.white} alt="" />;
    }
  };

  const goCv = () => {
    const userInfo = Cookies.getUserInfo();
    if (userInfo?.type === UserType.client) {
      window.location.href = `/client/${userInfo?.address_id}`;
    } else {
      window.location.href = `/storageProvide/${userInfo?.address_id}`;
    }
  };

  const loginNode = () => {
    if (access_token) {
      return (
        <div className={styles.login}>
          <div className={styles.logOff} onClick={logOff}>
            Log Out
          </div>
          <div className={styles.userLogo} onClick={goCv}>
            <HashImage placeholder={default_avatar} hash={avatar} />
          </div>
        </div>
      );
    } else {
      return (
        <div className={styles.unlogin}>
          <div className={styles.signIn} onClick={skipToSignIn}>
            Sign In
          </div>
          <Dropdown overlay={_dropDown} trigger={['click']}>
            <div className="menuRow">
              <div className="menuText">Sign Up</div>
              <DownOutlined />
            </div>
          </Dropdown>
        </div>
      );
    }
  };

  return (
    <div className={classnames(styles.header, isFixed ? styles.fixed : '')}>
      <div className={styles.headerContent}>
        <div className={styles.left}>
          <div className={styles.logo} onClick={goHome}>
            <img src={logo} alt="" />
          </div>
          <div className={styles.menuList}>
            {userStore.isLogin
              ? loginMenu.map((item, index) => {
                  return (
                    <div
                      className={classnames(
                        styles.menuItem,
                        headerStore.activeIndex === index ? styles.active : '',
                      )}
                      key={'CVHeader' + index}
                      onClick={() => {
                        skipTo(item);
                      }}
                    >
                      <div className={styles.iconWrap}>
                        {getIcon(item, headerStore.activeIndex === index)}
                      </div>
                      <div className={styles.text}>{item.text}</div>
                    </div>
                  );
                })
              : unLoginMenu.map((item, index) => {
                  return (
                    <div
                      className={classnames(
                        styles.menuItem,
                        headerStore.activeIndex === index ? styles.active : '',
                      )}
                      key={'CVHeader' + index}
                      onClick={() => {
                        skipTo(item);
                      }}
                    >
                      <div className={styles.iconWrap}>
                        {getIcon(item, headerStore.activeIndex === index)}
                      </div>
                      <div className={styles.text}>{item.text}</div>
                    </div>
                  );
                })}
          </div>
        </div>
        <div className={styles.right}>
          {/* <div className={styles.searchRow}>
            <Search
              placeholder="Miner ID/Name/Location"
              allowClear
              enterButton="Search"
              onSearch={onSearch}
            />
          </div> */}
          {loginNode()}
        </div>
      </div>
    </div>
  );
};

export default CVHeader;
