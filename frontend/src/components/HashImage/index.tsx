import styles from './style.module.scss';
import { useEffect, useState } from 'react';
import React from 'react';
import { getWeb3Image } from '@/utils';
interface PropsHash {
  hash: string | undefined;
  placeholder?: string | undefined;
}

const HashImage = React.memo(
  (props: PropsHash) => {
    const [url, setUrl] = useState<string>('');
    const [done, setDone] = useState<boolean>(false);
    useEffect(() => {
      previewUrl().then();
      // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [props.hash]);

    const previewUrl = async () => {
      if (props.hash) {
        const url = getWeb3Image(props.hash);
        console.log(url, 'urlurlurl123');

        setUrl(url);
      }
    };

    useEffect(() => {
      const img = new Image();
      img.src = url;
      img.onload = () => {
        setDone(true);
      };
    }, [url]);

    return done ? (
      <img className={styles.hashImage} src={url} alt="" />
    ) : (
      <div>
        <img className={styles.hashImage} src={props.placeholder} alt={''} />
      </div>
    );
  },
  // (prevProps, nextProps) => {
  //   return prevProps.hash !== nextProps.hash;
  // },
);

export default HashImage;
