import { Link, useHistory } from 'react-router-dom';
import styles from './style.module.scss';
import { Button, Form, Input } from 'antd';
import { LoginParam } from '@/api/modules/interface';
import { postUserLogin } from '@/api/modules';
import loginPng from '@/assets/images/login/login.png';
import { Cookies } from '@/data/cookies';
import { useState } from 'react';
import { ValidateStatus } from '@/utils/enum';
import logo from '@/assets/images/header/logo.png';
import { RouterPath } from '@/router/RouterConfig';
import classnames from 'classnames';
export enum signUpUser {
  Client = 'Client',
  SP = 'SP',
}

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

const Auth = () => {
  const history = useHistory();
  const [vilad, setVilad] = useState<ValidateStatus>(ValidateStatus.none);
  const [message, setMessage] = useState<string>('');
  const onSignIn = (values: LoginParam) => {
    postUserLogin(values)
      .then((res) => {
        if (res.data) {
          Cookies.setUserInfo(res.data);
          window.location.href = `/`;
        } else {
          setVilad(ValidateStatus.error);
          setMessage(res.msg);
        }
      })
      .catch((error) => {});
  };

  const goHome = () => {
    history.push(RouterPath.home);
  };

  const skipTo = (item) => {
    if (item.url) {
      history.push(item.url);
    }
  };

  return (
    <div className={styles['sign-in']}>
      <div className={styles.headerContent}>
        <div className={styles.left}>
          <div className={styles.logo} onClick={goHome}>
            <img src={logo} alt="" />
          </div>
          <div className={styles.menuList}>
            {menu.map((item, index) => {
              return (
                <div
                  className={classnames(styles.menuItem)}
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
      </div>
      <div className={styles['sign-in_content']}>
        <div className={styles.bgWrap}>
          <img src={loginPng} alt="" />
        </div>
        <div className={styles['sign-in__form']}>
          <h1 className={styles.title}>Sign In</h1>
          <h2 className={styles.description}>
            Register your account?
            <Link
              to={{
                pathname: '/signUp',
                state: { user: signUpUser.Client },
              }}
            >
              I'm a Client!
            </Link>
            <Link
              to={{
                pathname: '/signUp',
                state: { user: signUpUser.SP },
              }}
            >
              I'm a Storage Provider!
            </Link>
          </h2>
          <Form
            name={styles.address}
            layout="vertical"
            onFinish={(values: LoginParam) => {
              onSignIn(values);
            }}
          >
            <Form.Item
              name={'email'}
              label={'E-mail'}
              rules={[
                { type: 'email', message: 'Please enter the correct email.' },
              ]}
            >
              <Input placeholder="Enter your e-mail" />
            </Form.Item>
            <Form.Item
              validateStatus={vilad}
              help={message}
              name={'password'}
              label={'Password'}
            >
              <Input.Password placeholder="At least 8 characters with a mix of letters, numbers" />
            </Form.Item>
            <Form.Item
              wrapperCol={{ span: 24 }}
              shouldUpdate={(prevValues, curValues) => prevValues !== curValues}
            >
              {({ getFieldsValue }) => {
                const { email, password } = getFieldsValue(true);
                const isDisabled = !(email && password);

                return (
                  <Button
                    type="primary"
                    htmlType="submit"
                    disabled={isDisabled}
                  >
                    Sign In
                  </Button>
                );
              }}
            </Form.Item>
            <Form.Item wrapperCol={{ span: 24 }} noStyle>
              <Link className={styles.link} to="/changePassword" replace>
                Forget password？
              </Link>
            </Form.Item>
          </Form>
          <div className={styles.agreement}>
            By creating an account to use FileDrive Datasets, you
            unconditionally agree to our
            <a href="/"> Terms of Service</a>.
          </div>
        </div>
      </div>
    </div>
  );
};

export default Auth;
