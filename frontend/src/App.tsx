import React, { useEffect } from 'react';
import './App.css';
import MRouter from '@/router';
import { Provider } from 'mobx-react';
import * as stores from './store';
import ws from '@/socket/ws';
import { ChannelTypeEnum } from '@/socket/ws.types';
import { Cookies } from '@/data/cookies';
import chatStore from '@/store/modules/chat';
import { observer } from 'mobx-react';
import userStore from '@/store/modules/user';
import { StoragesKey } from '@/data/storages';
import md5 from 'md5';
import { envType } from './utils/enum';

function App() {
  useEffect(() => {
    if (process.env.NODE_ENV === envType.production) {
      window.addEventListener('load', disabledDevTools);
    }
    window.onstorage = (value) => {
      if (value.key === md5(StoragesKey.LoginStatus)) {
        window.location.reload();
      }
    };
    if (Cookies.getUserInfo()) {
      userStore.SET_USERINFO(Cookies.getUserInfo());
    }
  });

  useEffect(() => {
    if (!chatStore.isLogin) {
      ws.send({
        channel: ChannelTypeEnum.Login,
        subscribe: true,
        body: {
          token: Cookies.getUserInfo()?.access_token,
        },
      });
    }
  }, [chatStore.isLogin]);

  useEffect(() => {
    if (chatStore.chatAuth) {
      ws.send({
        channel: ChannelTypeEnum.ChatPartners,
        subscribe: true,
      });
    }
  }, [chatStore.isChatAuth]);

  const disabledDevTools = () => {
    const noop = () => undefined;
    const DEV_TOOLS = window['__REACT_DEVTOOLS_GLOBAL_HOOK__'];

    if (typeof DEV_TOOLS === 'object') {
      for (const [key, value] of Object.entries(DEV_TOOLS)) {
        DEV_TOOLS[key] = typeof value === 'function' ? noop : null;
      }
    }
  };

  return <Provider {...stores}>{<MRouter />}</Provider>;
}

export default observer(App);
