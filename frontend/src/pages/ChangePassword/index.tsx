import styles from './style.module.scss';
import { Input, Button, message } from 'antd';
import CVHeader from '@/components/CVHeader';
import Footer from '@/components/Footer';
import changePasswordStore from '@/store/modules/changePassword';
import { observer } from 'mobx-react';
import headerStore from '@/store/modules/header';
const ChangePassword = () => {
  headerStore.SET_ACTIVEINDEX(4);
  const validValue = (value) => {
    const reg = /^(?=.*\d)(?=.*?[a-z])(?=.*?[A-Z])[\da-zA-Z]{8,20}$/;
    const valid = reg.test(value);
    return valid;
  };
  const confirm = () => {
    changePasswordStore
      ._postUserModifyPassword({
        new_password: changePasswordStore.newPassword.value,
        password: changePasswordStore.recentPassword.value,
      })
      .then((res) => {
        message.success('修改密码成功！');
        setTimeout(() => {
          window.location.href = `/auth`;
        }, 3000);
      });
  };

  return (
    <div className={styles.changePassword}>
      <div>
        <div className={styles.top}>
          <CVHeader></CVHeader>
        </div>

        <div className={styles.content}>
          <div className={styles.title}>Change Password</div>
          <div className={styles.form}>
            <div className={styles.inputRow}>
              <div className={styles.inputItem}>
                <div className={styles.label}>Recent Password</div>
                <Input.Password
                  value={changePasswordStore.recentPassword.value}
                  type="password"
                  placeholder="Enter your password"
                  onChange={(event) => {
                    const _value = event.target.value;
                    if (_value) {
                      const valid = validValue(_value);
                      changePasswordStore.SET_RECENTPASSWORD({
                        valid: valid,
                        value: _value,
                        message: valid
                          ? ''
                          : 'At least 8 characters with a mix of letters, numbers & symbols',
                      });
                    } else {
                      changePasswordStore.SET_RECENTPASSWORD({
                        valid: false,
                        value: _value,
                        message: '',
                      });
                    }
                  }}
                ></Input.Password>
                <div className={styles.message}>
                  {changePasswordStore.recentPassword.message}
                </div>
              </div>
              <div className={styles.inputItem}>
                <div className={styles.label}>New Password</div>
                <Input.Password
                  value={changePasswordStore.newPassword.value}
                  type="password"
                  placeholder="At least 8 characters with a mix of letters, numbers & symbols"
                  onChange={(event) => {
                    const _value = event.target.value;

                    if (_value) {
                      const valid = validValue(_value);
                      changePasswordStore.SET_NEWPASSWORD({
                        valid: valid,
                        value: _value,
                        message: valid
                          ? ''
                          : 'At least 8 characters with a mix of letters, numbers & symbols',
                      });
                    } else {
                      changePasswordStore.SET_NEWPASSWORD({
                        valid: false,
                        value: _value,
                        message: '',
                      });
                    }
                  }}
                ></Input.Password>
                <div className={styles.message}>
                  {changePasswordStore.newPassword.message}
                </div>
              </div>
            </div>
            <div className={styles.inputRow}>
              <div className={styles.inputItem}>
                <div className={styles.label}>Confirm New Password</div>
                <Input.Password
                  value={changePasswordStore.confirmPassword.value}
                  type="password"
                  placeholder="At least 8 characters with a mix of letters, numbers & symbols"
                  onChange={(event) => {
                    const _value = event.target.value;
                    if (_value) {
                      const valid =
                        _value === changePasswordStore.newPassword.value;
                      changePasswordStore.SET_CONFIRMPASSWORD({
                        valid: valid,
                        value: _value,
                        message: valid
                          ? ''
                          : 'The two passwords that you entered do not match!',
                      });
                    } else {
                      changePasswordStore.SET_CONFIRMPASSWORD({
                        valid: false,
                        value: _value,
                        message: '',
                      });
                    }
                  }}
                ></Input.Password>
                <div className={styles.message}>
                  {changePasswordStore.confirmPassword.message}
                </div>
              </div>
            </div>
          </div>

          <div className={styles.action}>
            <Button
              type="primary"
              disabled={
                !changePasswordStore.recentPassword.valid ||
                !changePasswordStore.newPassword.valid ||
                !changePasswordStore.confirmPassword.valid
              }
              onClick={confirm}
            >
              Confirm
            </Button>
          </div>
        </div>
      </div>
      <Footer></Footer>
    </div>
  );
};

export default observer(ChangePassword);
