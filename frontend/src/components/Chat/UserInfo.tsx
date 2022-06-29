import classNames from 'classnames';
import styles from './UserInfo.module.scss';
import { ReactComponent as EmailIcon } from '@/assets/images/chat/email.svg';
import { ReactComponent as LocationIcon } from '@/assets/images/chat/location.svg';
import { Input, Tooltip } from 'antd';
import { ReactComponent as SearchIcon } from '@/assets/images/chat/search.svg';
import avatar from '@/assets/images/default_avatar.png';
import { Cookies } from '@/data/cookies';
import { desensitization } from '@/utils';
import HashImage from '@/components/HashImage';
import chatStore from '@/store/modules/chat';

const UserInfo = () => {
  const userInfo = Cookies.getUserInfo();

  const type = () => {
    return userInfo?.type === 'data_client'
      ? 'Filecoin Client'
      : 'Filecoin Storage Provider';
  };

  return (
    <div className={classNames(styles.userInfoWrap)}>
      <div className={styles.userInfoLogo}>
        <div className={styles.avatar}>
          <HashImage placeholder={avatar} hash={'userInfo?.avatar'} />
        </div>
        <div className={styles.icon}>
          <Tooltip title={userInfo?.email}>
            <EmailIcon />
          </Tooltip>
          <Tooltip title={userInfo?.location}>
            <LocationIcon />
          </Tooltip>
        </div>
      </div>
      <div className={styles.userInfoText}>
        {type() + '\nOwner: ' + desensitization(userInfo?.address || '', 6, 6)}
      </div>
      <div className={'mt-24'} />
      <Input
        onChange={(value) => {
          chatStore.SET_CONTACT_SEARCH(value.target.value);
        }}
        prefix={<SearchIcon />}
      />
    </div>
  );
};

export default UserInfo;
