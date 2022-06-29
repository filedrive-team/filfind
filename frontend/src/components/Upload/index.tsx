import styles from './style.module.scss';
import _upload from '@/assets/images/client/upload.png';
import { Web3Storage } from 'web3.storage/dist/bundle.esm.min.js';
import { observer } from 'mobx-react';
import { web3StorageToken } from '@/config';
import { useEffect, useState } from 'react';
import { notification, Spin } from 'antd';
import { getWeb3Image } from '@/utils';
const Upload = (props: any) => {
  const [url, setUrl] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const client = new Web3Storage({ token: web3StorageToken });
  useEffect(() => {
    previewUrl();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [props.hash]);
  const handleChange = (event) => {
    const _files = event.target.files;
    if (_files.length) {
      const size = _files[0].size;
      const _name = _files[0].name;
      const _max_size = 1024 * 1024 * 1;

      if (size > _max_size) {
        notification.error({
          key: 'pictureError',
          message: `The uploaded picture must not be greater than 1M`,
        });
      } else {
        setLoading(true);
        client
          .put(_files)
          .then((res) => {
            setLoading(false);
            props.valueChange(res + '/' + _name);
          })
          .catch((error) => {
            setLoading(false);
          });
      }
    }
  };
  const previewNode = () => {
    if (props.hash) {
      return (
        <div className={styles.previewWrap}>
          <img src={url} alt="" />
        </div>
      );
    } else {
      return (
        <div className={styles.imageWrap}>
          <img className={styles.photo} src={_upload} alt="" />
        </div>
      );
    }
  };

  const previewUrl = async () => {
    if (props.hash) {
      const url = getWeb3Image(props.hash);
      setUrl(url);
    }
  };

  return (
    <div className={styles.uploadWrap}>
      <input className={styles.file} type="file" onChange={handleChange} />
      <Spin spinning={loading}>{previewNode()}</Spin>
      <div className={styles.text}>{props.label}</div>
      <div className={styles.tips}>Size: max 1M</div>
    </div>
  );
};

export default observer(Upload);
