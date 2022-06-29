import { useEffect, useRef, useState } from 'react';
import { useHistory, useLocation } from 'react-router';
import { Link } from 'react-router-dom';
import { Button, Form, Input, notification, Steps } from 'antd';
import classNames from 'classnames';
import logo from '@/assets/images/header/logo.png';
import styles from './style.module.scss';
import loginPng from '@/assets/images/login/login.png';
import { signUpUser } from '@/pages/Auth';
import { SignUpParam } from '@/api/modules/interface';
import { postUserSignUp } from '@/api/modules';
import { strToHex } from '@/utils/transcode';
import { copyText } from '@/utils/clipboard';
import { RouterPath } from '@/router/RouterConfig';
import classnames from 'classnames';

const { Step } = Steps;

interface DataType extends SignUpParam {
  confirm?: string;
}

interface LocationParams {
  user: keyof typeof signUpUser;
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

const SignUp = () => {
  const history = useHistory();
  const {
    state: { user },
  } = useLocation<LocationParams>();

  const [baseInfoForm] = Form.useForm();
  const [addressForm] = Form.useForm();
  const [step, setStep] = useState<number>(0);
  const [message, setMessage] = useState<string>('');
  const [code, setCode] = useState<string>('');

  const addressToCode = (address: string | undefined): string => {
    const str =
      'Signature for filfind\r\n' + address + '\r\n' + new Date().toISOString();
    return address ? strToHex(str) : '';
  };

  const onSignUp = () => {
    const _address = addressForm.getFieldsValue();
    const _base = baseInfoForm.getFieldsValue();

    const param = {
      type: user === signUpUser.SP ? 'sp_owner' : 'data_client',
      ..._address,
      ..._base,
    };
    console.log(param, 'param_form');

    postUserSignUp(param)
      .then((res) => {
        if (res.code === 200) {
          notification.success({
            message: 'Register successfully',
          });
          history.push(RouterPath.auth);
        } else {
          notification.error({
            message: res.msg,
          });
        }
      })
      .catch((error) => console.log(error));
  };

  const layout = {
    labelCol: { span: 8 },
    wrapperCol: { span: 24 },
  };

  const baseInfoStep = (
    <Form
      style={{ display: step === 1 ? 'block' : 'none' }}
      form={baseInfoForm}
      {...layout}
      layout="vertical"
      name={styles.baseInfo}
      onFinish={() => {
        onSignUp();
      }}
    >
      <Form.Item name={'name'} label="Name" rules={[{ max: 128 }]}>
        <Input placeholder="Enter your name" />
      </Form.Item>
      <Form.Item name={'email'} label="E-mail" rules={[{ type: 'email' }]}>
        <Input placeholder="Enter your e-mail" />
      </Form.Item>
      <Form.Item name={'password'} label="Password" rules={[{ min: 8 }]}>
        <Input.Password placeholder="At least 8 characters with a mix of letters, numbers" />
      </Form.Item>
      <Form.Item
        name="confirm"
        label="Password Confirm"
        dependencies={['password']}
        rules={[
          ({ getFieldValue }) => ({
            validator(_, value) {
              if (!value || getFieldValue('password') === value) {
                return Promise.resolve();
              }
              return Promise.reject(
                new Error('The two passwords that you entered do not match!'),
              );
            },
          }),
        ]}
      >
        <Input.Password placeholder="At least 8 characters with a mix of letters, numbers & symbols" />
      </Form.Item>
      <span
        className={classNames(styles.link, 'pointer')}
        onClick={() => setStep(0)}
      >
        Previous step
      </span>
      <Form.Item
        shouldUpdate={(prevValues, curValues) => prevValues !== curValues}
        wrapperCol={{ ...layout.wrapperCol }}
      >
        {({ getFieldsValue }) => {
          const { name, email, password, confirm } = getFieldsValue(true);
          const isDisabled = !(name && email && password && confirm);

          return (
            <Button
              className={'mb-17 mt-17'}
              type="primary"
              htmlType="submit"
              disabled={isDisabled}
            >
              Sign Up
            </Button>
          );
        }}
      </Form.Item>
      <Link className={classNames(styles.link, 'mb-19')} to="/auth" replace>
        Log in now!
      </Link>
      <span className={styles.agreement}>
        By creating an account to use FileDrive Datasets, you unconditionally
        agree to our
        <a href="/"> Terms of Service</a>.
      </span>
    </Form>
  );

  const addressStep = (
    <div style={{ display: step === 0 ? 'block' : 'none' }}>
      <div className={classNames(styles.description, 'mb-25', 'pt-10')}>
        <span className={styles['description__title']}>
          Address Verification：
        </span>
        <br />
        To verify your address, please sign the following message and submit the
        signature.
        <br />
        E.g. sign {'\u007B'}address{'\u007d'} {'\u007B'}message{'\u007d'} on a
        lotus node
      </div>
      <Form
        {...layout}
        form={addressForm}
        name={'address'}
        layout="vertical"
        onFinish={(values) => {
          setStep(1);
        }}
      >
        <Form.Item
          name={'address'}
          label={`${user === 'SP' ? 'Owner' : 'Client'} Address`}
        >
          <Input
            onChange={(e) => {
              const address = e.target.value;
              const message = addressToCode(address);
              const code = 'lotus wallet sign ' + address + ' ' + message;
              setMessage(message);
              setCode(code);
              addressForm.setFieldsValue({ message: message });
            }}
            placeholder={`Enter your ${
              user === 'SP' ? 'SP Owner' : 'Client'
            } Node Address`}
          />
        </Form.Item>
        <Form.Item
          label="Code"
          name="message"
          // shouldUpdate={(prevValues, curValues) =>
          //   prevValues.address !== curValues.address
          // }
        >
          <div className={styles.code}>
            <Input.TextArea value={code} disabled />
            <span
              className="pointer"
              onClick={() => {
                copyText(code)
                  .then(() => notification.success({ message: 'Copied' }))
                  .catch(() => notification.error({ message: 'Copy failure' }));
              }}
            >
              Copy
            </span>
          </div>
        </Form.Item>
        <Form.Item label="Signature">
          <Form.Item name={'signature'}>
            <Input placeholder="Enter your signature" />
          </Form.Item>
          <div className={'mb-42'} />
        </Form.Item>
        <Form.Item
          shouldUpdate={(prevValues, curValues) => prevValues !== curValues}
          wrapperCol={{ ...layout.wrapperCol }}
        >
          {({ getFieldsValue }) => {
            const { address, signature } = getFieldsValue(true);
            const isDisabled = !(address && signature);

            return (
              <Button type="primary" htmlType="submit" disabled={isDisabled}>
                Next Step
              </Button>
            );
          }}
        </Form.Item>
      </Form>
    </div>
  );

  const goHome = () => {
    history.push(RouterPath.home);
  };

  const skipTo = (item) => {
    if (item.url) {
      history.push(item.url);
    }
  };

  const signIn = () => {
    history.push(RouterPath.auth);
  };

  const signUp = () => {
    if (user === signUpUser.Client) {
      history.push({
        pathname: RouterPath.signUp,
        state: { user: signUpUser.SP },
      });
    } else {
      history.push({
        pathname: RouterPath.signUp,
        state: { user: signUpUser.Client },
      });
    }
  };

  return (
    <div className={styles['sign-up']}>
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
        <div className={styles.right}>
          <div className={styles.signIn} onClick={signIn}>
            Sign in
          </div>
          <div className={styles.signUp} onClick={signUp}>
            {user === signUpUser.SP
              ? 'Client Sign Up'
              : 'Storage Provider Sign Up'}
          </div>
        </div>
      </div>
      <div className={styles['sign-up_content']}>
        <div className={styles.bgWrap}>
          <img src={loginPng} alt="" />
        </div>
        <div className={styles.form}>
          <div className={styles.title}>
            {`Sign Up as A ${user === 'SP' ? 'Storage Provider' : 'Client'} `}
          </div>
          <Steps className={styles.step} progressDot current={step}>
            <Step title="Step 01" />
            <Step title="Step 02" />
          </Steps>
          {addressStep} {baseInfoStep}
        </div>
      </div>
    </div>
  );
};

export default SignUp;
