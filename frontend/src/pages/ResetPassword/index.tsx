import { Link, useHistory } from 'react-router-dom';
import styles from './style.module.scss';
import { Button, Col, Form, Input, Row, notification } from 'antd';
import { ResetPwdParam } from '@/api/modules/interface';
import { postUserResetPwd, postVcodeByEmailToResetPwd } from '@/api/modules';
import loginPng from '@/assets/images/login/login.png';
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

const ResetPassword = () => {
  let timer;

  const [form] = Form.useForm();
  const history = useHistory();
  const [isSend, setIsSend] = useState<boolean>(false);
  const [second, setSecond] = useState<number>(60);
  const onSignIn = (values: ResetPwdParam) => {
    postUserResetPwd(values)
      .then((res) => {
        if (res.data) {
          window.location.href = `/auth`;
        } else {
          notification.error({
            key: 'pictureError',
            message: res.msg || 'Request error',
          });
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

  const sendEmail = () => {
    const _email = form.getFieldValue('email');
    if (_email) {
      postVcodeByEmailToResetPwd({ email: _email })
        .then((data) => {
          if (data.code === 200) {
            setIsSend(true);
            timer = setInterval(clock, 1000);
          } else {
            notification.error({
              key: 'pictureError',
              message: data.msg || 'Request error',
            });
          }
        })
        .catch((error) => {
          notification.error({
            key: 'pictureError',
            message: error,
          });
        });
    } else {
      notification.error({
        key: 'pictureError',
        message: 'Please enter E-mail.',
      });
    }
  };

  const clock = () => {
    setSecond((m) => {
      if (m > 0) {
        return m - 1;
      } else {
        clearInterval(timer);
        setIsSend(false);
        return 60;
      }
    });
  };

  const codeNode = () => {
    if (isSend) {
      return <div className="second">{second}</div>;
    } else {
      return (
        <div className="send" onClick={sendEmail}>
          Send
        </div>
      );
    }
  };

  return (
    <div className={styles['resetPassword']}>
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
      <div className={styles['resetPassword_content']}>
        <div className={styles.bgWrap}>
          <img src={loginPng} alt="" />
        </div>
        <div className={styles['resetPassword__form']}>
          <h1 className={styles.title}>Reset Password</h1>
          <h2 className={styles.description}>
            To reset your password, please use the same e-mail you register your
            account.
          </h2>
          <Form
            form={form}
            name={styles.address}
            layout="vertical"
            onFinish={(values: ResetPwdParam) => {
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

            <Form.Item label="Verification Code">
              <Row gutter={8}>
                <Col span={18}>
                  <Form.Item
                    name="vcode"
                    rules={[
                      {
                        required: true,
                        message: 'Please input verification code!',
                      },
                    ]}
                  >
                    <Input placeholder="Enter verification code sent to your e-mail" />
                  </Form.Item>
                </Col>
                <Col span={6} className="codeWrap">
                  {codeNode()}
                </Col>
              </Row>
            </Form.Item>

            <Form.Item name={'new_password'} label={'New Password'}>
              <Input.Password placeholder="At least 8 characters with a mix of letters, numbers" />
            </Form.Item>
            <Form.Item
              wrapperCol={{ span: 24 }}
              shouldUpdate={(prevValues, curValues) => prevValues !== curValues}
            >
              {({ getFieldsValue }) => {
                const { email, new_password, vcode } = getFieldsValue(true);
                const isDisabled = !(email && new_password && vcode);

                return (
                  <Button
                    type="primary"
                    htmlType="submit"
                    disabled={isDisabled}
                  >
                    Confirm
                  </Button>
                );
              }}
            </Form.Item>
            <Form.Item wrapperCol={{ span: 24 }} noStyle>
              <Link className={styles.link} to="/auth" replace>
                Log in now!
              </Link>
            </Form.Item>
          </Form>
          <div className={styles.agreement}>
            By creating an account to use FilFind Discovery Platform, you
            unconditionally agree to our
            <a href="/"> Terms of Service</a>.
          </div>
        </div>
      </div>
    </div>
  );
};

export default ResetPassword;
