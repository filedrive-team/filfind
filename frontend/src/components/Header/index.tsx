import styles from './style.module.scss';
import logo from '@/assets/images/header/logo.png';
import { Dropdown, Menu } from 'antd';
import classnames from 'classnames';
import { useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import { RouterPath } from '@/router/RouterConfig';
import { UserType } from '@/utils/enum';
import { signUpUser } from '@/pages/Auth';
import { Cookies } from '@/data/cookies';
import HashImage from '@/components/HashImage';
import headerStore from '@/store/modules/header';
import default_avatar from '@/assets/images/default_avatar.png';
import { DownOutlined } from '@ant-design/icons';
const menu = [
  {
    text: 'SPs Dashboard',
    url: RouterPath.home,
  },
  {
    text: 'Clients Dashboard',
    url: RouterPath.clientList,
  },
];

const Header = (props: any) => {
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
          label: 'Client Sign Up',
        },
        {
          key: signUpUser.SP,
          label: 'Storage Provider Sign Up',
        },
      ]}
    />
  );

  const _onScroll = () => {
    const _top = document.documentElement.scrollTop;
    if (_top > 80) {
      setIsFixed(true);
    } else {
      setIsFixed(false);
    }
  };

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
    if (item.url) {
      history.push(item.url);
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
            {menu.map((item, index) => {
              return (
                <div
                  className={classnames(
                    styles.menuItem,
                    headerStore.activeIndex === index ? styles.active : '',
                  )}
                  key={index}
                  onClick={() => {
                    skipTo(item);
                  }}
                >
                  {item.text}
                </div>
              );
            })}
          </div>
        </div>
        <div className={styles.right}>{loginNode()}</div>
      </div>
    </div>
  );
};

export default Header;
