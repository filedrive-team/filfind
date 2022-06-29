import styles from './AddNewReveiews.module.scss';
import { Input, Button, Select, Rate } from 'antd';
import { observer } from 'mobx-react';
import storageProvideStore from '@/store/modules/storageProvide';
const { TextArea } = Input;
const { Option } = Select;
const AddNewReveiews = () => {
  const clear = () => {
    storageProvideStore.SET_ADDREVIEWS({
      title: '',
      content: '',
      provider: '',
      score: 0,
    });
  };
  const submit = () => {
    storageProvideStore
      ._postClientReview(storageProvideStore.addReveiews)
      .then((res) => {
        clear();
        window.location.reload();
      });
  };
  const rateChange = (value) => {
    storageProvideStore.SET_ADDREVIEWS({
      ...storageProvideStore.addReveiews,
      score: value,
    });
  };

  return (
    <div className={styles.addNewReveiews}>
      <div className={styles.title}>Add New Reviews</div>
      <div className={styles.content}>
        <div className={styles.info}>
          <Select
            style={{ width: 120 }}
            onChange={(_value) => {
              storageProvideStore.SET_ADDREVIEWS({
                ...storageProvideStore.addReveiews,
                provider: _value,
              });
            }}
          >
            {storageProvideStore.MinerIDList.map((item, index) => {
              return (
                <Option value={item} key={index}>
                  {item}
                </Option>
              );
            })}
          </Select>
          <div className={styles.sorce}>
            <div className={styles.rateLabel}>Rating Score：</div>
            <Rate onChange={rateChange} />
          </div>
        </div>
        <div className={styles.label}>Title：</div>
        <div className={styles.inputRow}>
          <Input
            placeholder="Enter review title"
            value={storageProvideStore.addReveiews.title}
            maxLength={128}
            onChange={(event) => {
              const _value = event.target.value;
              storageProvideStore.SET_ADDREVIEWS({
                ...storageProvideStore.addReveiews,
                title: _value,
              });
            }}
          />
          <div className={styles.tips}>
            ( Attention：you can only submit reviews for storage providers who
            have stored yours deals.)
          </div>
        </div>
        <TextArea
          rows={3}
          maxLength={1024}
          value={storageProvideStore.addReveiews.content}
          onChange={(event) => {
            const _value = event.target.value;
            storageProvideStore.SET_ADDREVIEWS({
              ...storageProvideStore.addReveiews,
              content: _value,
            });
          }}
        />
        <div className={styles.action}>
          <Button onClick={clear}>CLEAR</Button>
          <Button
            type="primary"
            disabled={storageProvideStore.valid}
            onClick={submit}
          >
            SUBMIT
          </Button>
        </div>
      </div>
    </div>
  );
};

export default observer(AddNewReveiews);
